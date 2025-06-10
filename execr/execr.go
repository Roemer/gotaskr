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

var argumentsRegex = regexp.MustCompile(`(?m)[\w]+|"[\w\\"\s]*"`)

func NewCmd(executable string, arguments ...string) *exec.Cmd {
	return exec.Command(executable, arguments...)
}

func NewCmdSplitted(executable string, arguments string) *exec.Cmd {
	return exec.Command(executable, SplitArguments(arguments)...)
}

func Run(executable string, arguments []string, options ...func(*RunOptions)) error {
	cmd := NewCmd(executable, arguments...)
	return RunCommand(cmd, options...)
}

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

func RunGetOutput(executable string, arguments []string, options ...func(*RunOptions)) (string, string, error) {
	cmd := NewCmd(executable, arguments...)
	return RunCommandGetOutput(cmd, options...)
}

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

func RunGetCombinedOutput(executable string, arguments []string, options ...func(*RunOptions)) (string, error) {
	cmd := NewCmd(executable, arguments...)
	return RunCommandGetCombinedOutput(cmd, options...)
}

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

func SplitArguments(arguments string) []string {
	return argumentsRegex.FindAllString(arguments, -1)
}

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
