package execr

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"regexp"

	"github.com/roemer/gotaskr/goext"
	"github.com/roemer/gotaskr/log"
)

var argumentsRegex = regexp.MustCompile(`[^\s"]+|"((\\"|[^"])*)"`)

// Creates a new exec.Cmd with the given executable and arguments.
func NewCmd(executable string, arguments ...string) *exec.Cmd {
	return exec.Command(executable, arguments...)
}

// Creates a new exec.Cmd with the given executable and the argument string splitted into separate arguments.
func NewCmdSplitted(executable string, arguments string) *exec.Cmd {
	return exec.Command(executable, SplitArguments(arguments)...)
}

// Runs an executable with the given arguments.
func Run(executable string, arguments []string, options ...func(*RunOptions)) error {
	cmd := NewCmd(executable, arguments...)
	return RunCommand(cmd, options...)
}

// Runs a command with the given arguments.
func RunCommand(cmd *exec.Cmd, options ...func(*RunOptions)) error {
	runOptions := prepare(cmd, options...)
	logArguments(cmd)
	if runOptions.outputToConsole {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	return err
}

// Runs an executable with the given arguments and returns the separate output from stdout and stderr.
func RunGetOutput(executable string, arguments []string, options ...func(*RunOptions)) (string, string, error) {
	cmd := NewCmd(executable, arguments...)
	return RunCommandGetOutput(cmd, options...)
}

// Runs a command with the given arguments and returns the separate output from stdout and stderr.
func RunCommandGetOutput(cmd *exec.Cmd, options ...func(*RunOptions)) (string, string, error) {
	runOptions := prepare(cmd, options...)
	logArguments(cmd)
	var stdoutBuf, stderrBuf bytes.Buffer
	if runOptions.outputToConsole {
		cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	} else {
		cmd.Stdout = &stdoutBuf
		cmd.Stderr = &stderrBuf
	}
	err := cmd.Run()
	if runOptions.skipPostProcessOutput {
		return stdoutBuf.String(), stderrBuf.String(), err
	}
	return processOutputString(stdoutBuf.String()), processOutputString(stderrBuf.String()), err
}

// Runs an executable with the given arguments and returns the output from stdout and stderr combined.
func RunGetCombinedOutput(executable string, arguments []string, options ...func(*RunOptions)) (string, error) {
	cmd := NewCmd(executable, arguments...)
	return RunCommandGetCombinedOutput(cmd, options...)
}

// Runs a command with the given arguments and returns the output from stdout and stderr combined.
func RunCommandGetCombinedOutput(cmd *exec.Cmd, options ...func(*RunOptions)) (string, error) {
	runOptions := prepare(cmd, options...)
	logArguments(cmd)
	var outBuf bytes.Buffer
	if runOptions.outputToConsole {
		cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)
		cmd.Stderr = io.MultiWriter(os.Stderr, &outBuf)
	} else {
		cmd.Stdout = &outBuf
		cmd.Stderr = &outBuf
	}
	err := cmd.Run()
	if runOptions.skipPostProcessOutput {
		return outBuf.String(), err
	}
	return processOutputString(outBuf.String()), err
}

// Splits a string of arguments into a slice of strings, handling quoted arguments correctly.
func SplitArguments(arguments string) []string {
	return argumentsRegex.FindAllString(arguments, -1)
}

// Returns a slice of strings containing the given arguments.
func Arguments(arguments ...string) []string {
	return arguments
}

////////////////////////////////////////////////////////////
// Internal
////////////////////////////////////////////////////////////

func prepare(cmd *exec.Cmd, options ...func(*RunOptions)) *RunOptions {
	// Build the options
	runOptions := &RunOptions{}
	for _, option := range options {
		option(runOptions)
	}
	// Arguments
	if len(runOptions.arguments) > 0 {
		cmd.Args = append(cmd.Args, runOptions.arguments...)
	}
	// Working directory
	if runOptions.workingDirectory != "" {
		cmd.Dir = runOptions.workingDirectory
	}
	return runOptions
}

func logArguments(cmd *exec.Cmd) {
	log.Debugf("Executing '%s' with arguments: %s", cmd.Path, cmd.Args[1:])
}

func processOutputString(value string) string {
	return goext.TrimNewlineSuffix(value)
}
