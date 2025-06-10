package execr

import (
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
	return exec.Command(executable, SplitArgumentString(arguments)...)
}

// Splits a string of arguments into a slice of strings, handling quoted arguments correctly.
func SplitArgumentString(arguments string) []string {
	return argumentsRegex.FindAllString(arguments, -1)
}

////////////////////////////////////////////////////////////
// Internal
////////////////////////////////////////////////////////////

func logArguments(cmd *exec.Cmd) {
	log.Debugf("Executing '%s' with arguments: %s", cmd.Path, cmd.Args[1:])
}

func processOutputString(value string) string {
	return goext.TrimNewlineSuffix(value)
}
