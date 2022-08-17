package execr

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/roemer/gotaskr/goext"
)

type CmdError struct {
	msg      string
	ExitCode int
}

func (e *CmdError) Error() string { return e.msg }

func Run(executable string, arguments ...string) error {
	cmd := exec.Command(executable, arguments...)
	return RunCommand(cmd)
}

// RunCommand runs a command and writes the stdout and stderr into the console in realtime.
func RunCommand(cmd *exec.Cmd) error {
	//stdReader, _ := cmd.StdoutPipe()
	//cmd.Stderr = cmd.Stdout

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	goext.Printfln("Executing %s with arguments: %s", cmd.Path, cmd.Args[1:])

	err := cmd.Start()
	if err != nil {
		return err
	}

	/*scanner := bufio.NewScanner(stdReader)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}*/
	if err := cmd.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode := exitError.ExitCode()
			return &CmdError{msg: fmt.Sprintf("Cmd failed with exit code %d", exitCode), ExitCode: exitCode}
		}
	}
	return nil
}
