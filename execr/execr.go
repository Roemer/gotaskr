// Package execr is a wapper to run exec commands.
package execr

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/roemer/gotaskr/log"
)

// CmdError is an error which is returned when a command failed to execute.
type CmdError struct {
	msg      string // Contains the error message
	ExitCode int    // Contains the exit code of the command
}

func (e *CmdError) Error() string { return e.msg }

// Run runs an executable with the given arguments.
func Run(executable string, arguments ...string) error {
	cmd := exec.Command(executable, arguments...)
	return RunCommand(cmd)
}

// RunCommand runs a command and writes the stdout and stderr into the console in realtime.
func RunCommand(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Debugf("Executing %s with arguments: %s", cmd.Path, cmd.Args[1:])

	err := cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode := exitError.ExitCode()
			return &CmdError{msg: fmt.Sprintf("Cmd failed with exit code %d", exitCode), ExitCode: exitCode}
		}
	}
	return err
}
