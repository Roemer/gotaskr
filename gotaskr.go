// Package gotaskr provides the basic methods to register and run tasks.
// It also provides the main entrypoint for gotaskr.
package gotaskr

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/roemer/gotaskr/argparse"
	"github.com/roemer/gotaskr/goext"
	"github.com/roemer/gotaskr/gttools"
	"github.com/roemer/gotaskr/log"
)

// Generate a map that holds all passed arguments from the cli
var argumentsMap = argparse.ParseArgs()

// Prepare a map for all the task objects
var taskMap map[string]*TaskObject = make(map[string]*TaskObject)

// Prepare a list of the task names. Used to print the tasks in order
var taskList []string

// Prepare an array for the tasks that were run (in run order)
var taskRun []*TaskObject

// The task object of the currently running task
var currentRunningTask *TaskObject

// A context for the current gotaskr run
var context gotaskrContext = gotaskrContext{}

// Tools provides typed access to the various tools supported.
var Tools *gttools.ToolsClient = gttools.CreateToolsClient()

// Execute is the entry point of gotaskr.
func Execute() int {
	log.Initialize(HasArgument("verbose") || HasArgument("v"))

	target, hasTarget := GetArgument("target")
	if !hasTarget {
		printTasks()
		return 0
	}

	// Log start
	log.Information(strings.Repeat("-", 60))
	log.Information("Running gotaskr")
	log.Information(strings.Repeat("-", 60))
	printArguments()
	log.Information()
	// Validate dependencies and convert dependees to dependencies
	for _, task := range taskMap {
		for _, followup := range task.followups {
			followupTask := taskMap[followup]
			if followupTask == nil {
				color.Red("Followup '%s' for '%s' does not exist.", followup, task.name)
				return 1
			}
		}
		for _, dependency := range task.dependencies {
			dependencyTask := taskMap[dependency]
			if dependencyTask == nil {
				color.Red("Dependency '%s' for '%s' does not exist.", dependency, task.name)
				return 1
			}
		}
		for _, dependee := range task.dependees {
			dependeeTask := taskMap[dependee]
			if dependeeTask == nil {
				color.Red("Dependee '%s' for '%s' does not exist.", dependee, task.name)
				return 1
			}
			dependeeTask.DependsOn(task.name)
		}
	}

	// Run the setup method
	setupErr := runLifetimeFunc("Setup", context.SetupFunc)

	// In case of a setup error, run the teardown and exit
	if setupErr != nil {
		// We can ignore a possible teardown error
		_ = runLifetimeFunc("Teardown", context.TeardownFunc)
		return getExitCodeFromError(setupErr)
	}

	// Run the main target only if the setup succeeded
	taskErr := RunTarget(target)

	// Run the teardown method
	teardownErr := runLifetimeFunc("Teardown", context.TeardownFunc)

	// Run finished
	log.Information()
	log.Information(strings.Repeat("-", 60))
	log.Informationf("Finished gotaskr at %s", time.Now().Format("2006-01-02 15:04:05.000"))
	exitCode := getExitCodeFromError(taskErr)

	// Print errors and check the deferred errors
	for _, run := range taskRun {
		printTaskError(run, true)
		if exitCode == 0 {
			exitCode = getExitCodeFromTaskRun(run)
		}
	}
	log.Information(strings.Repeat("-", 60))
	log.Information()
	printTaskRuns()

	// If the teardown failed but nothing else, still fail with the teardowns error
	if teardownErr != nil && exitCode == 0 {
		exitCode = getExitCodeFromError(teardownErr)
	}

	return exitCode
}

// GetArgument returns the value of the argument with the given name
// and also a flag, if the argument was present or not.
func GetArgument(argName string) (string, bool) {
	return GetArgumentOrDefault(argName, "")
}

// GetArgumentOrDefault returns the value of the argument with the given name
// or the given default value if the value was not present
// and also a flag, if the argument was present or not.
func GetArgumentOrDefault(argName string, defaultValue string) (string, bool) {
	value, exists := argumentsMap[argName]
	if exists {
		return value, true
	}
	return defaultValue, false
}

// HasArgument returns true if an arument was set and false otherwise, regardless of the value.
func HasArgument(argName string) bool {
	_, exist := GetArgument(argName)
	return exist
}

// RunTarget runs the given task and all the needed dependencies.
func RunTarget(target string) error {
	var currentTask = taskMap[target]
	currentRunningTask = currentTask
	// Early exit if the target does not exist
	if currentTask == nil {
		err := fmt.Errorf("target does not exist: %s", target)
		color.Red("%v", err)
		return err
	}
	// Early exit if the task did already run
	if currentTask.didRun {
		return currentTask.err
	}
	// Get the flag for exclusive runs
	exclusive := HasArgument("exclusive") || HasArgument("e")
	// Run dependencies
	if !exclusive && len(currentTask.dependencies) > 0 {
		for _, dependency := range currentTask.dependencies {
			dependencyErr := RunTarget(dependency)
			if dependencyErr != nil {
				if currentTask.deferOnError {
					// Handle deferred errors
					currentTask.deferredErr = dependencyErr
				} else {
					return dependencyErr
				}
			}
		}
	}

	// Run the task setup method
	setupErr := runLifetimeFunc("TaskSetup", context.TaskSetupFunc)

	// In case of a setup error, run the teardown and exit
	if setupErr != nil {
		// We can ignore a possible teardown error
		_ = runLifetimeFunc("TaskTeardown", context.TaskTeardownFunc)
		return setupErr
	}

	// Run the task itself
	currentRunningTask = currentTask
	printTaskHeader(target)
	start := time.Now()
	taskErr := runTaskFunc(currentTask)
	elapsed := time.Since(start)
	// Handle error deferring
	if taskErr != nil && currentTask.deferOnError {
		currentTask.deferredErr = taskErr
		taskErr = nil
	}
	// Handle error skipping
	if taskErr != nil && currentTask.continueOnError {
		currentTask.ignoredErr = taskErr
		taskErr = nil
	}
	currentTask.didRun = true
	currentTask.duration = elapsed
	currentTask.err = taskErr
	taskRun = append(taskRun, currentTask)
	printTaskFooter(currentTask)

	// Run the task teardown method
	teardownErr := runLifetimeFunc("TaskTeardown", context.TaskTeardownFunc)

	// If the teardown failed but nothing else, still fail with the teardowns error
	if teardownErr != nil && taskErr == nil {
		return teardownErr
	}

	// Abort execution if the task failed
	if taskErr != nil {
		return taskErr
	}
	// Run followup tasks
	if !exclusive && len(currentTask.followups) > 0 {
		for _, followup := range currentTask.followups {
			followupErr := RunTarget(followup)
			if followupErr != nil {
				if currentTask.deferOnError {
					// Handle deferred errors
					currentTask.deferredErr = followupErr
				} else {
					return followupErr
				}
			}
		}
	}
	if currentTask.deferredErr != nil {
		return currentTask.deferredErr
	}
	return nil
}

func runLifetimeFunc(lifetimeStage string, function func() error) error {
	if function == nil {
		return nil
	}
	log.Informationf("--- %s %s", lifetimeStage, strings.Repeat("-", 60-5-len(lifetimeStage)))
	err := runFuncRecover(function)
	if err != nil {
		log.Informationf("Error occured: %v", err)
		return err
	}
	return nil
}

func runTaskFunc(currentTask *TaskObject) (err error) {
	return runFuncRecover(currentTask.taskFunc)
}

func runFuncRecover(function func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("task panicked: %v", r)
		}
	}()
	err = function()
	return err
}

// Task registers the given function with the name so it can be executed.
func Task(name string, taskFunc func() error) *TaskObject {
	task := TaskObject{}
	task.name = name
	task.taskFunc = taskFunc
	taskMap[name] = &task
	taskList = append(taskList, name)
	return &task
}

func Setup(setupFunc func() error) {
	context.SetupFunc = setupFunc
}

func Teardown(taskFunc func() error) {
	context.TeardownFunc = taskFunc
}

func TaskSetup(taskFunc func() error) {
	context.TaskSetupFunc = taskFunc
}

func TaskTeardown(taskFunc func() error) {
	context.TaskTeardownFunc = taskFunc
}

type argument struct {
	name        string
	description string
	optional    bool
}

type gotaskrContext struct {
	SetupFunc        func() error
	TeardownFunc     func() error
	TaskSetupFunc    func() error
	TaskTeardownFunc func() error
}

type timeMeasurement struct {
	name      string
	startTime time.Time
	duration  time.Duration
}

// TaskObject represents a registered task.
type TaskObject struct {
	name             string        // The name of the task.
	description      string        // The description of the task.
	arguments        []argument    // The arguments of the task.
	taskFunc         func() error  // The function of the task.
	dependencies     []string      // A list of dependecy tasks.
	dependees        []string      // A list of dependee tasks.
	followups        []string      // A list of followup tasks.
	continueOnError  bool          // A flag to incdicate if the run should continue when an error occured.
	deferOnError     bool          // A flag to indicate if the error should be deferred until the end.
	didRun           bool          // A flag to indicate if the task did already run.
	duration         time.Duration // A runtime duration of the task if it ran already.
	err              error         // The error (if any) of the task when it ran.
	ignoredErr       error         // The error (if any) which is ignored.
	deferredErr      error         // The deferred error (if any) of the task when it ran.
	timeMeasurements []*timeMeasurement
}

// GetName gets the name of the task.
func (taskObject *TaskObject) GetName() string {
	return taskObject.name
}

// DependsOn adds dependencies in the given order. Duplicate dependencies are removed.
func (taskObject *TaskObject) DependsOn(taskName ...string) *TaskObject {
	for _, entry := range taskName {
		if entry == taskObject.name {
			// Skip itself
			continue
		}
		taskObject.dependencies = goext.AppendIfMissing(taskObject.dependencies, entry)
	}
	return taskObject
}

// DependsOnTask adds dependencies in the given order. Duplicate dependencies are removed.
func (taskObject *TaskObject) DependsOnTask(task ...*TaskObject) *TaskObject {
	for _, entry := range task {
		if entry.name == taskObject.name {
			// Skip itself
			continue
		}
		taskObject.dependencies = goext.AppendIfMissing(taskObject.dependencies, entry.name)
	}
	return taskObject
}

// DependeeOf adds dependees in the given order. Duplicate dependees are removed.
func (taskObject *TaskObject) DependeeOf(taskName ...string) *TaskObject {
	for _, entry := range taskName {
		if entry == taskObject.name {
			// Skip itself
			continue
		}
		taskObject.dependees = goext.AppendIfMissing(taskObject.dependees, entry)
	}
	return taskObject
}

// DependeeOfTask adds dependees in the given order. Duplicate dependees are removed.
func (taskObject *TaskObject) DependeeOfTask(task ...*TaskObject) *TaskObject {
	for _, entry := range task {
		if entry.name == taskObject.name {
			// Skip itself
			continue
		}
		taskObject.dependees = goext.AppendIfMissing(taskObject.dependees, entry.name)
	}
	return taskObject
}

// Then adds followup tasks in the given order.
func (taskObject *TaskObject) Then(taskName ...string) *TaskObject {
	for _, entry := range taskName {
		if entry == taskObject.name {
			// Skip itself
			continue
		}
		taskObject.followups = goext.AppendIfMissing(taskObject.followups, entry)
	}
	return taskObject
}

// Then adds followup tasks in the given order.
func (taskObject *TaskObject) ThenTask(task ...*TaskObject) *TaskObject {
	for _, entry := range task {
		if entry.name == taskObject.name {
			// Skip itself
			continue
		}
		taskObject.followups = goext.AppendIfMissing(taskObject.followups, entry.name)
	}
	return taskObject
}

// ContinueOnError will continue with dependencies or dependees even when the task returned an error.
func (taskObject *TaskObject) ContinueOnError() *TaskObject {
	taskObject.continueOnError = true
	return taskObject
}

// DeferOnError will continue with dependencies or dependees even when the task returned an error.
func (taskObject *TaskObject) DeferOnError() *TaskObject {
	taskObject.deferOnError = true
	return taskObject
}

// Description sets the description of a task. Will be shown when the help is displayed.
func (taskObject *TaskObject) Description(description string) *TaskObject {
	taskObject.description = description
	return taskObject
}

// Argument adds a description for an argument. Will be shown when the help is displayed.
func (taskObject *TaskObject) Argument(argumentName string, argumentDescription string, optional bool) *TaskObject {
	newArgument := argument{
		name:        argumentName,
		description: argumentDescription,
		optional:    optional,
	}
	taskObject.arguments = append(taskObject.arguments, newArgument)
	return taskObject
}

// AddFollowupTask allows to add one or more tasks that should run after the current finished.
func AddFollowupTask(taskName ...string) {
	currentRunningTask.Then(taskName...)
}

func MeasureTime(measurementName string, f func() error) error {
	// Execute the function
	start := time.Now()
	err := f()
	elapsed := time.Since(start)

	// Add the time measurement
	currentRunningTask.timeMeasurements = append(currentRunningTask.timeMeasurements, &timeMeasurement{
		name:      measurementName,
		startTime: start,
		duration:  elapsed,
	})

	return err
}

func printTasks() {
	log.Information("Please specify one of the following targets:")
	var sb strings.Builder
	for _, taskName := range taskList {
		task := taskMap[taskName]
		fmt.Fprintf(&sb, "- %s", task.name)
		sb.WriteString(log.Newline)
		if task.description != "" {
			lines := goext.SplitByNewLine(task.description)
			for _, line := range lines {
				fmt.Fprintf(&sb, "  %s", line)
				sb.WriteString(log.Newline)
			}
		}
		if len(task.arguments) > 0 {
			fmt.Fprintln(&sb, "  Arguments:")
			for _, arg := range task.arguments {
				fmt.Fprintf(&sb, "    %s: %s%s", arg.name, arg.description, goext.Ternary(arg.optional, " (optional)", ""))
				sb.WriteString(log.Newline)
			}
		}
	}
	log.Information(sb.String())
}

func printArguments() {
	if len(argumentsMap) > 0 {
		log.Debug("Arguments:")
		var sb strings.Builder
		isFirst := true
		for key, val := range argumentsMap {
			if !isFirst {
				sb.WriteString(", ")
			}
			if isFirst {
				isFirst = false
			}
			fmt.Fprintf(&sb, "%s=\"%s\"", key, val)
		}
		sb.WriteString(log.Newline)
		log.Debug(sb.String())
	}
}

func printTaskHeader(taskName string) {
	log.Informationf("=== %s %s", taskName, strings.Repeat("=", 60-5-len(taskName)))
}

func printTaskFooter(task *TaskObject) {
	log.Informationf("=== /%s %s", task.name, strings.Repeat("=", 60-5-1-len(task.name)))
	log.Informationf("Duration: %s", formatDuration(task.duration))
	printTaskError(task, false)
}

func printTaskError(task *TaskObject, withTaskName bool) {
	taskString := goext.Ternary(withTaskName, fmt.Sprintf(" in '%s'", task.name), "")
	if task.err != nil {
		color.Red("Task error%s: %v", taskString, task.err)
	}
	if task.ignoredErr != nil {
		color.Red("Ignored error%s: %v", taskString, task.ignoredErr)
	}
	if task.deferredErr != nil {
		color.Red("Deferred error%s: %v", taskString, task.deferredErr)
	}
}

func printTaskRuns() {
	if len(taskRun) == 0 {
		return
	}
	color.Set(color.FgGreen)
	defer color.Unset()
	log.Informationf("%-50s%-13s%-17s", "Task", "Exit Code", "Duration")
	log.Information(strings.Repeat("-", 80))
	totalDuration := time.Duration(0)
	for _, run := range taskRun {
		text := fmt.Sprintf("%-50s%-13d%-17s", run.name, getExitCodeFromTaskRun(run), formatDuration(run.duration))
		if run.err != nil || run.deferredErr != nil {
			color.Red(text)

		} else {
			log.Information(text)
		}
		color.Set(color.FgWhite)
		for i, measurement := range run.timeMeasurements {
			prefix := goext.Ternary(i == len(run.timeMeasurements)-1, "└─", "├─")
			measurementText := fmt.Sprintf("%s %-60s%-17s", prefix, measurement.name, formatDuration(measurement.duration))
			log.Information(measurementText)
		}
		color.Set(color.FgGreen)
		totalDuration += run.duration
	}
	log.Information(strings.Repeat("-", 80))
	log.Informationf("%-63s%-18s", "Total", formatDuration(totalDuration))
}

func formatDuration(duration time.Duration) string {
	hour := int(duration.Seconds() / 3600)
	minute := int(duration.Seconds()/60) % 60
	second := int(duration.Seconds()) % 60
	micro := duration.Microseconds() - (int64(duration.Seconds()) * 1000000)
	return fmt.Sprintf("%02d:%02d:%02d.%06d", hour, minute, second, micro)
}

func getExitCodeFromTaskRun(run *TaskObject) int {
	if ec := getExitCodeFromError(run.err); ec != 0 {
		return ec
	}
	if ec := getExitCodeFromError(run.deferredErr); ec != 0 {
		return ec
	}
	// No Error
	return 0
}

func getExitCodeFromError(err error) int {
	if err != nil {
		if ierr, ok := err.(*exec.ExitError); ok {
			// Exit code from exec
			return ierr.ExitCode()
		} else {
			// Any other error
			return 1
		}
	}
	// No Error
	return 0
}

func clear() {
	taskMap = map[string]*TaskObject{}
	taskList = []string{}
	taskRun = []*TaskObject{}
	context.SetupFunc = nil
	context.TeardownFunc = nil
	context.TaskSetupFunc = nil
	context.TaskTeardownFunc = nil
}
