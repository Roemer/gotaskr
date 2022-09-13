package gotaskr

import (
	"fmt"
	"os/exec"
	"runtime"
	"testing"

	"github.com/roemer/gotaskr/execr"
	"github.com/roemer/gotaskr/goext"
	"github.com/stretchr/testify/assert"
)

func TestNoErrorTask(t *testing.T) {
	assert := assert.New(t)

	// Prepare
	clear()
	task := Task("NoErrorTask", goext.Noop)
	argumentsMap = map[string]string{"target": task.name}

	// Execute
	exitCode := Execute()

	// Validate
	assert.Equal(0, exitCode)
	assert.Equal(1, len(taskRun))
	assert.Nil(taskRun[0].err)
}

func TestErrorTask(t *testing.T) {
	assert := assert.New(t)

	// Prepare
	clear()
	desiredExitCode := 10
	task := Task("ErrorTask", func() error { return getExitError(desiredExitCode) })
	argumentsMap = map[string]string{"target": task.name}

	// Execute
	exitCode := Execute()

	// Validate
	assert.Equal(1, len(taskRun))
	assert.Equal(desiredExitCode, exitCode)
	assertExitError(assert, task.err, desiredExitCode)
}

func TestDependencyErrorWithDefer(t *testing.T) {
	assert := assert.New(t)

	// Prepare
	clear()
	desiredExitCode := 10
	task1 := Task("Test1", func() error { return nil })
	task2 := Task("Test2", func() error { return getExitError(desiredExitCode) })
	task3 := Task("Test3", func() error { return nil })
	taskAll := Task("All", func() error { return nil }).DependsOn(task1.name).DependsOn(task2.name).DependsOn(task3.name).DeferOnError()
	argumentsMap = map[string]string{"target": taskAll.name}

	// Execute
	exitCode := Execute()

	// Validate
	assert.Equal(desiredExitCode, exitCode)
	assert.Equal(4, len(taskRun))
	assertExitError(assert, task2.err, desiredExitCode)
	assertExitError(assert, taskAll.deferedErr, desiredExitCode)
}

func TestDependencyErrorWithoutDefer(t *testing.T) {
	assert := assert.New(t)

	// Prepare
	clear()
	desiredExitCode := 10
	task1 := Task("Test1", func() error { return nil })
	task2 := Task("Test2", func() error { return getExitError(desiredExitCode) })
	task3 := Task("Test3", func() error { return nil })
	taskAll := Task("All", func() error { return nil }).DependsOn(task1.name).DependsOn(task2.name).DependsOn(task3.name)
	argumentsMap = map[string]string{"target": taskAll.name}

	// Execute
	exitCode := Execute()

	// Validate
	assert.Equal(desiredExitCode, exitCode)
	assert.Equal(2, len(taskRun))
	assertExitError(assert, task2.err, desiredExitCode)
}

func TestDependencyErrorWithDeferOnTask(t *testing.T) {
	assert := assert.New(t)

	// Prepare
	clear()
	desiredExitCode := 10
	task1 := Task("Test1", func() error { return nil })
	task2 := Task("Test2", func() error { return getExitError(desiredExitCode) }).DeferOnError()
	task3 := Task("Test3", func() error { return nil })
	taskAll := Task("All", func() error { return nil }).DependsOn(task1.name).DependsOn(task2.name).DependsOn(task3.name)
	argumentsMap = map[string]string{"target": taskAll.name}

	// Execute
	exitCode := Execute()

	// Validate
	assert.Equal(desiredExitCode, exitCode)
	assert.Equal(2, len(taskRun))
	assertExitError(assert, task2.deferedErr, desiredExitCode)
}

////////////////////
// Helpers
////////////////////

func getExitError(errorCode int) error {
	return execr.Run(false, getBashPath(), "-c", fmt.Sprintf("exit %d", errorCode))
}

func getBashPath() string {
	if runtime.GOOS == "windows" {
		return `C:\Program Files\Git\bin\bash.exe`
	}
	return "/bin/bash"
}

func assertExitError(assert *assert.Assertions, err error, exitCode int) {
	assert.NotNil(err)
	if err != nil {
		ierr, isError := err.(*exec.ExitError)
		assert.True(isError)
		assert.Equal(exitCode, ierr.ExitCode())
	}
}
