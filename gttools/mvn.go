package gttools

import (
	"strings"

	"github.com/roemer/goext"
)

// MvnTool provides access to the helper methods for Maven.
type MvnTool struct {
	ToolBase
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
	// Skips executing the tests
	SkipTests bool
	// Skips compiling the tests
	MavenTestSkip bool
}

// Run runs mvn according to the settings.
func (tool *MvnTool) Run(settings *MvnRunSettings) error {
	args := []string{}
	args = append(args, settings.Phases...)
	args = goext.SliceAppendIf(args, len(settings.Projects) > 0, "--projects", strings.Join(settings.Projects, ","))
	args = goext.SliceAppendIf(args, len(settings.ActivateProfiles) > 0, "--activate-profiles", strings.Join(settings.ActivateProfiles, ","))
	args = goext.SliceAppendIf(args, settings.AlsoMake, "--also-make")
	args = goext.SliceAppendIf(args, settings.BatchMode, "--batch-mode")
	args = goext.SliceAppendIf(args, settings.Debug, "--debug")
	args = goext.SliceAppendIf(args, settings.Help, "--help")
	args = goext.SliceAppendIf(args, settings.NoTransferProgress, "--no-transfer-progress")
	args = goext.SliceAppendIf(args, settings.Offline, "--offline")
	args = goext.SliceAppendIf(args, settings.Quiet, "--quiet")
	args = goext.SliceAppendIf(args, settings.Settings != "", "--settings", settings.Settings)
	args = goext.SliceAppendIf(args, settings.Version, "--version")
	args = goext.SliceAppendIf(args, settings.ShowVersion, "--show-version")
	args = goext.SliceAppendIf(args, settings.SkipTests, "-DskipTests")
	args = goext.SliceAppendIf(args, settings.MavenTestSkip, "-Dmaven.test.skip")
	args = append(args, settings.CustomArguments...)

	return tool.run("mvn", args, settings.ToolSettingsBase)
}
