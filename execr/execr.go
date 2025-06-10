// Package execr is a wrapper to run exec commands.
package execr

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

// Runs an executable with the given arguments.
func Run(outputToConsole bool, executable string, arguments ...string) error {
	cmd := exec.Command(executable, arguments...)
	return RunCommand(outputToConsole, cmd)
}

// Runs an executable with the given arguments and returns the output.
func RunGetOutput(outputToConsole bool, executable string, arguments ...string) (string, string, error) {
	cmd := exec.Command(executable, arguments...)
	return RunCommandGetOutput(outputToConsole, cmd)
}

// Runs an executable with the given arguments and returns the output.
func RunGetCombinedOutput(outputToConsole bool, executable string, arguments ...string) (string, error) {
	cmd := exec.Command(executable, arguments...)
	return RunCommandGetCombinedOutput(outputToConsole, cmd)
}

// Runs a command and writes the stdout and stderr into the console in realtime.
func RunCommand(outputToConsole bool, cmd *exec.Cmd) error {
	logArguments(cmd)

	if outputToConsole {
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

func RunCommandGetOutput(outputToConsole bool, cmd *exec.Cmd) (string, string, error) {
	logArguments(cmd)

	var stdoutBuf, stderrBuf bytes.Buffer
	if outputToConsole {
		cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	} else {
		cmd.Stdout = &stdoutBuf
		cmd.Stderr = &stderrBuf
	}
	err := cmd.Run()
	outStr, errStr := processOutputString(stdoutBuf.String()), processOutputString(stderrBuf.String())
	return outStr, errStr, err
}

func RunCommandGetCombinedOutput(outputToConsole bool, cmd *exec.Cmd) (string, error) {
	logArguments(cmd)

	var outBuf bytes.Buffer
	if outputToConsole {
		cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)
		cmd.Stderr = io.MultiWriter(os.Stderr, &outBuf)
	} else {
		cmd.Stdout = &outBuf
		cmd.Stderr = &outBuf
	}
	err := cmd.Run()
	outStr := processOutputString(outBuf.String())
	return outStr, err
}
