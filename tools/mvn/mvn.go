package mvn

import (
	"os/exec"
	"strings"

	"github.com/roemer/gotaskr/execr"
	"github.com/roemer/gotaskr/goext"
)

type MvnRunSettings struct {
	WorkingDirectory   string
	Phases             []string
	Projects           []string
	AlsoMake           bool
	BatchMode          bool
	Debug              bool
	NoTransferProgress bool
	Offline            bool
	Quiet              bool
	Version            bool
	ShowVersion        bool
}

func Run(outputToConsole bool, settings MvnRunSettings) error {
	args := []string{}
	args = append(args, settings.Phases...)
	args = append(args, "--projects", strings.Join(settings.Projects, ","))
	args = goext.AddIf(args, settings.AlsoMake, "--also-make")
	args = goext.AddIf(args, settings.BatchMode, "--batch-mode")
	args = goext.AddIf(args, settings.Debug, "--debug")
	args = goext.AddIf(args, settings.NoTransferProgress, "--no-transfer-progress")
	args = goext.AddIf(args, settings.Offline, "--offline")
	args = goext.AddIf(args, settings.Quiet, "--quiet")
	args = goext.AddIf(args, settings.Version, "--version")
	args = goext.AddIf(args, settings.ShowVersion, "--show-version")

	cmd := exec.Command("mvn", args...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(outputToConsole, cmd)
}
