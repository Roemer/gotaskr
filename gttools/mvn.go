package gttools

import (
	"os/exec"
	"strings"

	"github.com/roemer/gotaskr/execr"
	"github.com/roemer/gotaskr/goext"
)

// MvnTool provides access to the helper methods for Maven.
type MvnTool struct {
}

func CreateMvnTool() *MvnTool {
	return &MvnTool{}
}

// MvnRunSettings is the settings object to tell mvn what to do.
type MvnRunSettings struct {
	ToolSettingsBase
	Phases             []string
	Projects           []string
	ActivateProfiles   []string
	AlsoMake           bool
	BatchMode          bool
	Debug              bool
	Help               bool
	NoTransferProgress bool
	Offline            bool
	Quiet              bool
	Settings           string
	Version            bool
	ShowVersion        bool
}

// Run runs mvn according to the settings.
func (tool *MvnTool) Run(settings MvnRunSettings) error {
	args := []string{}
	args = append(args, settings.Phases...)
	args = append(args, "--projects", strings.Join(settings.Projects, ","))
	args = append(args, "--activate-profiles", strings.Join(settings.ActivateProfiles, ","))
	args = goext.AddIf(args, settings.AlsoMake, "--also-make")
	args = goext.AddIf(args, settings.BatchMode, "--batch-mode")
	args = goext.AddIf(args, settings.Debug, "--debug")
	args = goext.AddIf(args, settings.Help, "--help")
	args = goext.AddIf(args, settings.NoTransferProgress, "--no-transfer-progress")
	args = goext.AddIf(args, settings.Offline, "--offline")
	args = goext.AddIf(args, settings.Quiet, "--quiet")
	args = goext.AddIf(args, settings.Settings != "", "--settings", settings.Settings)
	args = goext.AddIf(args, settings.Version, "--version")
	args = goext.AddIf(args, settings.ShowVersion, "--show-version")
	args = append(args, settings.CustomArguments...)

	cmd := exec.Command("mvn", args...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(settings.OutputToConsole, cmd)
}
