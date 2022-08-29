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

func logArguments(cmd *exec.Cmd) {
	log.Debugf("Executing '%s' with arguments: %s", cmd.Path, cmd.Args[1:])
}
