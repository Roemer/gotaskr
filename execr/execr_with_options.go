package execr

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

// Runs an executable with the given options.
func RunO(executable string, options ...func(*RunOptions)) error {
	cmd := NewCmd(executable)
	return RunCommandO(cmd, options...)
}

// Runs a command with the given options.
func RunCommandO(cmd *exec.Cmd, options ...func(*RunOptions)) error {
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

// Runs an executable with the given options and returns the separate output from stdout and stderr.
func RunOGetOutput(executable string, options ...func(*RunOptions)) (string, string, error) {
	cmd := NewCmd(executable)
	return RunCommandOGetOutput(cmd, options...)
}

// Runs a command with the given options and returns the separate output from stdout and stderr.
func RunCommandOGetOutput(cmd *exec.Cmd, options ...func(*RunOptions)) (string, string, error) {
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

// Runs an executable with the given options and returns the output from stdout and stderr combined.
func RunOGetCombinedOutput(executable string, options ...func(*RunOptions)) (string, error) {
	cmd := NewCmd(executable)
	return RunCommandOGetCombinedOutput(cmd, options...)
}

// Runs a command with the given options and returns the output from stdout and stderr combined.
func RunCommandOGetCombinedOutput(cmd *exec.Cmd, options ...func(*RunOptions)) (string, error) {
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
