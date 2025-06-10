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
	assertExitError(assert, taskAll.deferredErr, desiredExitCode)
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
	assertExitError(assert, task2.deferredErr, desiredExitCode)
}

func TestLifeTimeNoError(t *testing.T) {
	assert := assert.New(t)

	// Prepare
	clear()
	setupCalled := false
	teardownCalled := false
	taskSetupCalled := false
	taskTeardownCalled := false
	taskCalled := false
	Setup(func() error { setupCalled = true; return nil })
	Teardown(func() error { teardownCalled = true; return nil })
	TaskSetup(func() error { taskSetupCalled = true; return nil })
	TaskTeardown(func() error { taskTeardownCalled = true; return nil })
	dummyTask := Task("Test1", func() error { taskCalled = true; return nil })
	argumentsMap = map[string]string{"target": dummyTask.name}

	// Execute
	exitCode := Execute()

	// Validate
	assert.Equal(0, exitCode)
	assert.True(setupCalled)
	assert.True(teardownCalled)
	assert.True(taskSetupCalled)
	assert.True(taskTeardownCalled)
	assert.True(taskCalled)
}

func TestLifeTimeSetupError(t *testing.T) {
	assert := assert.New(t)

	// Prepare
	clear()
	setupCalled := false
	teardownCalled := false
	taskSetupCalled := false
	taskTeardownCalled := false
	taskCalled := false
	Setup(func() error { setupCalled = true; return getExitError(10) })
	Teardown(func() error { teardownCalled = true; return nil })
	TaskSetup(func() error { taskSetupCalled = true; return nil })
	TaskTeardown(func() error { taskTeardownCalled = true; return nil })
	dummyTask := Task("Test1", func() error { taskCalled = true; return nil })
	argumentsMap = map[string]string{"target": dummyTask.name}

	// Execute
	exitCode := Execute()

	// Validate
	assert.Equal(10, exitCode)
	assert.True(setupCalled)
	assert.True(teardownCalled)
	assert.False(taskSetupCalled)
	assert.False(taskTeardownCalled)
	assert.False(taskCalled)
}

func TestLifeTimeTeardownError(t *testing.T) {
	assert := assert.New(t)

	// Prepare
	clear()
	setupCalled := false
	teardownCalled := false
	taskSetupCalled := false
	taskTeardownCalled := false
	taskCalled := false
	Setup(func() error { setupCalled = true; return nil })
	Teardown(func() error { teardownCalled = true; return getExitError(20) })
	TaskSetup(func() error { taskSetupCalled = true; return nil })
	TaskTeardown(func() error { taskTeardownCalled = true; return nil })
	dummyTask := Task("Test1", func() error { taskCalled = true; return nil })
	argumentsMap = map[string]string{"target": dummyTask.name}

	// Execute
	exitCode := Execute()

	// Validate
	assert.Equal(20, exitCode)
	assert.True(setupCalled)
	assert.True(teardownCalled)
	assert.True(taskSetupCalled)
	assert.True(taskTeardownCalled)
	assert.True(taskCalled)
}

func TestLifeTimeTeardownErrorWithFailedTask(t *testing.T) {
	assert := assert.New(t)

	// Prepare
	clear()
	setupCalled := false
	teardownCalled := false
	taskSetupCalled := false
	taskTeardownCalled := false
	taskCalled := false
	Setup(func() error { setupCalled = true; return nil })
	Teardown(func() error { teardownCalled = true; return getExitError(20) })
	TaskSetup(func() error { taskSetupCalled = true; return nil })
	TaskTeardown(func() error { taskTeardownCalled = true; return nil })
	dummyTask := Task("Test1", func() error { taskCalled = true; return getExitError(30) })
	argumentsMap = map[string]string{"target": dummyTask.name}

	// Execute
	exitCode := Execute()

	// Validate
	assert.Equal(30, exitCode)
	assert.True(setupCalled)
	assert.True(teardownCalled)
	assert.True(taskSetupCalled)
	assert.True(taskTeardownCalled)
	assert.True(taskCalled)
}

func TestLifeTimeSetupAndTeardownError(t *testing.T) {
	assert := assert.New(t)

	// Prepare
	clear()
	setupCalled := false
	teardownCalled := false
	taskSetupCalled := false
	taskTeardownCalled := false
	taskCalled := false
	Setup(func() error { setupCalled = true; return getExitError(10) })
	Teardown(func() error { teardownCalled = true; return getExitError(20) })
	TaskSetup(func() error { taskSetupCalled = true; return nil })
	TaskTeardown(func() error { taskTeardownCalled = true; return nil })
	dummyTask := Task("Test1", func() error { taskCalled = true; return nil })
	argumentsMap = map[string]string{"target": dummyTask.name}

	// Execute
	exitCode := Execute()

	// Validate
	assert.Equal(10, exitCode)
	assert.True(setupCalled)
	assert.True(teardownCalled)
	assert.False(taskSetupCalled)
	assert.False(taskTeardownCalled)
	assert.False(taskCalled)
}

func TestLifeTimeTaskSetupError(t *testing.T) {
	assert := assert.New(t)

	// Prepare
	clear()
	setupCalled := false
	teardownCalled := false
	taskSetupCalled := false
	taskTeardownCalled := false
	taskCalled := false
	Setup(func() error { setupCalled = true; return nil })
	Teardown(func() error { teardownCalled = true; return nil })
	TaskSetup(func() error { taskSetupCalled = true; return getExitError(100) })
	TaskTeardown(func() error { taskTeardownCalled = true; return nil })
	dummyTask := Task("Test1", func() error { taskCalled = true; return nil })
	argumentsMap = map[string]string{"target": dummyTask.name}

	// Execute
	exitCode := Execute()

	// Validate
	assert.Equal(100, exitCode)
	assert.True(setupCalled)
	assert.True(teardownCalled)
	assert.True(taskSetupCalled)
	assert.True(taskTeardownCalled)
	assert.False(taskCalled)
}

func TestLifeTimeTaskTeardownError(t *testing.T) {
	assert := assert.New(t)

	// Prepare
	clear()
	setupCalled := false
	teardownCalled := false
	taskSetupCalled := false
	taskTeardownCalled := false
	taskCalled := false
	Setup(func() error { setupCalled = true; return nil })
	Teardown(func() error { teardownCalled = true; return nil })
	TaskSetup(func() error { taskSetupCalled = true; return nil })
	TaskTeardown(func() error { taskTeardownCalled = true; return getExitError(200) })
	dummyTask := Task("Test1", func() error { taskCalled = true; return nil })
	argumentsMap = map[string]string{"target": dummyTask.name}

	// Execute
	exitCode := Execute()

	// Validate
	assert.Equal(200, exitCode)
	assert.True(setupCalled)
	assert.True(teardownCalled)
	assert.True(taskSetupCalled)
	assert.True(taskTeardownCalled)
	assert.True(taskCalled)
}

func TestLifeTimeTaskTeardownErrorWithFailedTask(t *testing.T) {
	assert := assert.New(t)

	// Prepare
	clear()
	setupCalled := false
	teardownCalled := false
	taskSetupCalled := false
	taskTeardownCalled := false
	taskCalled := false
	Setup(func() error { setupCalled = true; return nil })
	Teardown(func() error { teardownCalled = true; return nil })
	TaskSetup(func() error { taskSetupCalled = true; return nil })
	TaskTeardown(func() error { taskTeardownCalled = true; return getExitError(200) })
	dummyTask := Task("Test1", func() error { taskCalled = true; return getExitError(30) })
	argumentsMap = map[string]string{"target": dummyTask.name}

	// Execute
	exitCode := Execute()

	// Validate
	assert.Equal(30, exitCode)
	assert.True(setupCalled)
	assert.True(teardownCalled)
	assert.True(taskSetupCalled)
	assert.True(taskTeardownCalled)
	assert.True(taskCalled)
}

func TestLifeTimeTaskSetupAndTaskTeardownError(t *testing.T) {
	assert := assert.New(t)

	// Prepare
	clear()
	setupCalled := false
	teardownCalled := false
	taskSetupCalled := false
	taskTeardownCalled := false
	taskCalled := false
	Setup(func() error { setupCalled = true; return nil })
	Teardown(func() error { teardownCalled = true; return nil })
	TaskSetup(func() error { taskSetupCalled = true; return getExitError(100) })
	TaskTeardown(func() error { taskTeardownCalled = true; return getExitError(200) })
	dummyTask := Task("Test1", func() error { taskCalled = true; return nil })
	argumentsMap = map[string]string{"target": dummyTask.name}

	// Execute
	exitCode := Execute()

	// Validate
	assert.Equal(100, exitCode)
	assert.True(setupCalled)
	assert.True(teardownCalled)
	assert.True(taskSetupCalled)
	assert.True(taskTeardownCalled)
	assert.False(taskCalled)
}

////////////////////
// Helpers
////////////////////

func getExitError(errorCode int) error {
	return execr.Run(getBashPath(), []string{"-c", fmt.Sprintf("exit %d", errorCode)})
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
