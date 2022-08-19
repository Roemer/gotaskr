// Package gotaskr provides the basic methods to register and run tasks.
// It also provides the main entrypoint for gotaskr.
package gotaskr

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/roemer/gotaskr/argparse"
	"github.com/roemer/gotaskr/execr"
	"github.com/roemer/gotaskr/log"
)

// Generate a map that holds all passed arguments from the cli
var argumentsMap = argparse.ParseArgs()

// Prepare a map for all the task objects
var taskMap map[string]*TaskObject = make(map[string]*TaskObject)

// Prepare an array for the tasks that were run (in run order)
var taskRun []*TaskObject

// Execute is the entry point of gotaskr.
func Execute() int {
	log.Initialize(HasArgument("verbose") || HasArgument("v"))

	target, hasTarget := GetArgument("target")
	if !hasTarget {
		printTasks()
		return 0
	}

	log.Information("Running gotaskr")
	printArguments()
	err := RunTarget(target)
	exitCode := 0
	if err != nil {
		if err, ok := err.(*execr.CmdError); ok {
			// Custom exit code form CmdErrors
			exitCode = err.ExitCode
		} else {
			exitCode = 1
		}
	}
	printTaskRuns()
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

// HasArgument returns true if an arument was set and false otherwise.
func HasArgument(argName string) bool {
	_, exist := GetArgument(argName)
	return exist
}

// RunTarget runs the given task and all the needed dependencies.
func RunTarget(target string) error {
	var currentTask = taskMap[target]
	// Early exit if the target does not exist
	if currentTask == nil {
		return fmt.Errorf("Target does not exist: %s", target)
	}
	// Early exit if the task did already run
	if currentTask.didRun {
		return currentTask.err
	}
	// Run dependencies
	if len(currentTask.dependencies) > 0 {
		for _, dependency := range currentTask.dependencies {
			err := RunTarget(dependency)
			if err != nil {
				return err
			}
		}
	}
	printTaskHeader(target)
	start := time.Now()
	err := runTaskFunc(currentTask)
	elapsed := time.Since(start)
	currentTask.didRun = true
	currentTask.duration = elapsed
	currentTask.err = err
	taskRun = append(taskRun, currentTask)
	if err != nil {
		color.Red("Failed with error: %v", err)
	}
	log.Information() // Make sure to add a newline to separate from the next task
	return err
}

func runTaskFunc(currentTask *TaskObject) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Task panicked: %v", r)
			log.Information(err)
		}
	}()
	err = currentTask.taskFunc()
	return err
}

// Task registers the given function with the name so it can be executed.
func Task(name string, taskFunc func() error) *TaskObject {
	task := TaskObject{}
	task.name = name
	task.taskFunc = taskFunc
	taskMap[name] = &task
	return &task
}

// TaskObject represents a registered task.
type TaskObject struct {
	name         string        // The name of the task.
	description  string        // The description of the task.
	taskFunc     func() error  // The function of the task.
	dependencies []string      // A list of dependecy tasks.
	didRun       bool          // A flag to indicate if the task did already run.
	duration     time.Duration // A runtime duration of the task if it ran already.
	err          error         // The error (if any) of the task when it ran.
}

// DependsOn adds dependencies in the given order. Duplicate dependencies are removed.
func (taskObject *TaskObject) DependsOn(taskName ...string) *TaskObject {
	keys := make(map[string]bool)
	for _, entry := range taskName {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			taskObject.dependencies = append(taskObject.dependencies, entry)
		}
	}
	return taskObject
}

// Description sets the description of a task. Will be shown when the help is displayed.
func (taskObject *TaskObject) Description(description string) *TaskObject {
	taskObject.description = description
	return taskObject
}

func printTasks() {
	log.Information("Please specify one of the following targets:")
	var sb strings.Builder
	for _, task := range taskMap {
		fmt.Fprintf(&sb, " - %s", task.name)
		if task.description != "" {
			fmt.Fprintf(&sb, ": %s", task.description)
		}
		sb.WriteString(log.Newline)
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
	log.Information(strings.Repeat("=", 40))
	log.Informationf("%-50s", taskName)
	log.Information(strings.Repeat("=", 40))
}

func printTaskRuns() {
	color.Set(color.FgGreen)
	defer color.Unset()
	log.Informationf("%-30s%-20s", "Task", "Duration")
	log.Information(strings.Repeat("-", 50))
	totalDuration := time.Duration(0)
	for _, run := range taskRun {
		text := fmt.Sprintf("%-30s%-20s", run.name, formatDuration(run.duration))
		if run.err != nil {
			color.Red(text)
			color.Set(color.FgGreen)
		} else {
			log.Information(text)
		}
		totalDuration += run.duration
	}
	log.Information(strings.Repeat("-", 50))
	log.Informationf("%-30s%-20s", "Total", formatDuration(totalDuration))
}

func formatDuration(duration time.Duration) string {
	hour := int(duration.Seconds() / 3600)
	minute := int(duration.Seconds()/60) % 60
	second := int(duration.Seconds()) % 60
	micro := duration.Microseconds() - (int64(duration.Seconds()) * 1000000)
	return fmt.Sprintf("%02d:%02d:%02d.%06d", hour, minute, second, micro)
}
