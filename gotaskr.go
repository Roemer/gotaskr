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
)

// Generate a map that holds all passed arguments from the cli
var argumentsMap = argparse.ParseArgs()

// Prepare a map for all the task objects
var taskMap map[string]*TaskObject = make(map[string]*TaskObject)

// Prepare an array for the tasks that were run (in run order)
var taskRun []*TaskObject

// Execute is the entry point of gotaskr.
func Execute() int {
	if !HasArgument("target") {
		printTasks()
		return 0
	}

	fmt.Print("Running gotaskr")
	printArguments()
	fmt.Println()
	target := GetArgument("target", "Usage")
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
// or the defaultValue otherwise.
func GetArgument(argName string, defaultValue string) string {
	v, exist := argumentsMap[argName]
	if exist {
		return v
	}
	return defaultValue
}

// HasArgument returns true if an arument was set and false otherwise.
func HasArgument(argName string) bool {
	_, exist := argumentsMap[argName]
	return exist
}

// RunTarget runs the given task and all the needed dependencies.
func RunTarget(target string) error {
	var currentTask = taskMap[target]
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
	fmt.Println()
	return err
}

func runTaskFunc(currentTask *TaskObject) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Task panicked: %v", r)
			fmt.Println(err)
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
	name         string
	description  string
	taskFunc     func() error
	dependencies []string
	didRun       bool
	duration     time.Duration
	err          error
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
	fmt.Println("Please specify one of the following targets:")
	for _, task := range taskMap {
		fmt.Printf(" - %s", task.name)
		if task.description != "" {
			fmt.Printf(": %s", task.description)
		}
		fmt.Println()
	}
}

func printArguments() {
	if len(argumentsMap) > 0 {
		fmt.Println(" with arguments:")
		isFirst := true
		for key, val := range argumentsMap {
			if !isFirst {
				fmt.Print(", ")
			}
			if isFirst {
				isFirst = false
			}
			s := fmt.Sprintf("%s=\"%s\"", key, val)
			fmt.Print(s)
		}
	}
	fmt.Println()
}
func printTaskHeader(taskName string) {
	fmt.Println(strings.Repeat("=", 40))
	fmt.Printf("%-50s", taskName)
	fmt.Println()
	fmt.Println(strings.Repeat("=", 40))
}

func printTaskRuns() {
	color.Set(color.FgGreen)
	defer color.Unset()
	fmt.Printf("%-30s%-20s", "Task", "Duration")
	fmt.Println()
	fmt.Println(strings.Repeat("-", 50))
	totalDuration := time.Duration(0)
	for _, run := range taskRun {
		text := fmt.Sprintf("%-30s%-20s", run.name, formatDuration(run.duration))
		if run.err != nil {
			color.Red(text)
			color.Set(color.FgGreen)
		} else {
			fmt.Println(text)
		}
		totalDuration += run.duration
	}
	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("%-30s%-20s", "Total", formatDuration(totalDuration))
	fmt.Println()
}

func formatDuration(duration time.Duration) string {
	hour := int(duration.Seconds() / 3600)
	minute := int(duration.Seconds()/60) % 60
	second := int(duration.Seconds()) % 60
	micro := duration.Microseconds() - (int64(duration.Seconds()) * 1000000)
	return fmt.Sprintf("%02d:%02d:%02d.%06d", hour, minute, second, micro)
}
