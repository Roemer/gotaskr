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
func Run(executable string, arguments ...string) error {
	cmd := exec.Command(executable, arguments...)
	return RunCommand(cmd)
}

// Run runs an executable with the given arguments and returns the output.
func RunGetOutput(executable string, arguments ...string) (string, string, error) {
	cmd := exec.Command(executable, arguments...)
	return RunCommandGetOutput(cmd, false)
}

// Run runs an executable with the given arguments and returns the output.
func RunGetCombinedOutput(executable string, arguments ...string) (string, error) {
	cmd := exec.Command(executable, arguments...)
	return RunCommandGetCombinedOutput(cmd, false)
}

// RunCommand runs a command and writes the stdout and stderr into the console in realtime.
func RunCommand(cmd *exec.Cmd) error {
	logArguments(cmd)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	return err
}

func RunCommandGetOutput(cmd *exec.Cmd, alsoOutputToOs bool) (string, string, error) {
	logArguments(cmd)

	var stdoutBuf, stderrBuf bytes.Buffer
	if alsoOutputToOs {
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

func RunCommandGetCombinedOutput(cmd *exec.Cmd, alsoOutputToOs bool) (string, error) {
	logArguments(cmd)

	var outBuf bytes.Buffer
	if alsoOutputToOs {
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

// SplitArgumentString splits the given string by spaces (preserving quotes)
func SplitArgumentString(s string) []string {
	quoted := false
	return strings.FieldsFunc(s, func(r rune) bool {
		if r == '"' {
			quoted = !quoted
		}
		return !quoted && r == ' '
	})
}

func SplitByNewLine(value string) []string {
	return strings.Split(strings.ReplaceAll(value, "\r\n", "\n"), "\n")
}

func logArguments(cmd *exec.Cmd) {
	log.Debugf("Executing '%s' with arguments: %s", cmd.Path, cmd.Args[1:])
}
