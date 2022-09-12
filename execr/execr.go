// Package execr is a wapper to run exec commands.
package execr

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/roemer/gotaskr/log"
)

// Run runs an executable with the given arguments.
func Run(outputToConsole bool, executable string, arguments ...string) error {
	cmd := exec.Command(executable, arguments...)
	return RunCommand(outputToConsole, cmd)
}

// Run runs an executable with the given arguments and returns the output.
func RunGetOutput(outputToConsole bool, executable string, arguments ...string) (string, string, error) {
	cmd := exec.Command(executable, arguments...)
	return RunCommandGetOutput(outputToConsole, cmd)
}

// Run runs an executable with the given arguments and returns the output.
func RunGetCombinedOutput(outputToConsole bool, executable string, arguments ...string) (string, error) {
	cmd := exec.Command(executable, arguments...)
	return RunCommandGetCombinedOutput(outputToConsole, cmd)
}

// RunCommand runs a command and writes the stdout and stderr into the console in realtime.
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
	outStr, errStr := stdoutBuf.String(), stderrBuf.String()
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
	outStr := outBuf.String()
	return outStr, err
}

// SplitArgumentString splits the given string by spaces (preserving quotes).
func SplitArgumentString(s string) []string {
	quoted := false
	return strings.FieldsFunc(s, func(r rune) bool {
		if r == '"' {
			quoted = !quoted
		}
		return !quoted && r == ' '
	})
}

// SplitByNewLine splits the given value by newlines.
func SplitByNewLine(value string) []string {
	return strings.Split(strings.ReplaceAll(value, "\r\n", "\n"), "\n")
}

func logArguments(cmd *exec.Cmd) {
	log.Debugf("Executing '%s' with arguments: %s", cmd.Path, cmd.Args[1:])
}
