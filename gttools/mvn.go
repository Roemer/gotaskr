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
	SkipTests          bool
}

// Run runs mvn according to the settings.
func (tool *MvnTool) Run(settings MvnRunSettings) error {
	args := []string{}
	args = append(args, settings.Phases...)
	args = goext.AppendIf(args, len(settings.Projects) > 0, "--projects", strings.Join(settings.Projects, ","))
	args = goext.AppendIf(args, len(settings.ActivateProfiles) > 0, "--activate-profiles", strings.Join(settings.ActivateProfiles, ","))
	args = goext.AppendIf(args, settings.AlsoMake, "--also-make")
	args = goext.AppendIf(args, settings.BatchMode, "--batch-mode")
	args = goext.AppendIf(args, settings.Debug, "--debug")
	args = goext.AppendIf(args, settings.Help, "--help")
	args = goext.AppendIf(args, settings.NoTransferProgress, "--no-transfer-progress")
	args = goext.AppendIf(args, settings.Offline, "--offline")
	args = goext.AppendIf(args, settings.Quiet, "--quiet")
	args = goext.AppendIf(args, settings.Settings != "", "--settings", settings.Settings)
	args = goext.AppendIf(args, settings.Version, "--version")
	args = goext.AppendIf(args, settings.ShowVersion, "--show-version")
	args = goext.AppendIf(args, settings.SkipTests, "-DskipTests")
	args = append(args, settings.CustomArguments...)

	cmd := exec.Command("mvn", args...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(settings.OutputToConsole, cmd)
}
