package gttools

import (
	"github.com/roemer/goext"
)

// DotNetVerbosity can be used to set the verbosity of the log output.
type DotNetVerbosity int

const (
	DotNetVerbosityUndefined DotNetVerbosity = iota
	DotNetVerbosityQuiet
	DotNetVerbosityMinimal
	DotNetVerbosityNormal
	DotNetVerbosityDetailed
	DotNetVerbosityDiagnostic
)

func (s DotNetVerbosity) String() string {
	switch s {
	case DotNetVerbosityQuiet:
		return "quiet"
	case DotNetVerbosityMinimal:
		return "minimal"
	case DotNetVerbosityNormal:
		return "normal"
	case DotNetVerbosityDetailed:
		return "detailed"
	case DotNetVerbosityDiagnostic:
		return "diagnostic"
	}
	return "unknown"
}

type DotNetTool struct {
	ToolBase
}

func CreateDotNetTool() *DotNetTool {
	return &DotNetTool{}
}

type DotNetSettings struct {
	Verbosity        DotNetVerbosity
	DiagnosticOutput bool
}

////////////////////////////////////////////////////////////
// Build
////////////////////////////////////////////////////////////

type DotNetBuildSettings struct {
	ToolSettingsBase
	DotNetSettings
	Configuration   string
	OutputDirectory string
	Runtime         string
	Framework       string
	NoIncremental   bool
	NoDependencies  bool
	NoRestore       bool
	NoLogo          bool
}

func (tool *DotNetTool) DotNetBuild(project string, settings *DotNetBuildSettings) error {
	if settings == nil {
		settings = &DotNetBuildSettings{}
	}
	args := []string{
		"build",
	}
	args = goext.SliceAppendIf(args, project != "", project)
	args = goext.SliceAppendIf(args, settings.Configuration != "", "--configuration", settings.Configuration)
	args = goext.SliceAppendIf(args, settings.OutputDirectory != "", "--output", settings.OutputDirectory)
	args = goext.SliceAppendIf(args, settings.Runtime != "", "--runtime", settings.Runtime)
	args = goext.SliceAppendIf(args, settings.Framework != "", "--framework", settings.Framework)
	args = goext.SliceAppendIf(args, settings.NoIncremental, "--no-incremental")
	args = goext.SliceAppendIf(args, settings.NoDependencies, "--no-dependencies")
	args = goext.SliceAppendIf(args, settings.NoRestore, "--no-restore")
	args = goext.SliceAppendIf(args, settings.NoLogo, "--nologo")
	args = append(args, settings.CustomArguments...)
	args = tool.addDotNetSettings(args, settings.DotNetSettings)

	return tool.run("dotnet", args, settings.ToolSettingsBase)
}

////////////////////////////////////////////////////////////
// Clean
////////////////////////////////////////////////////////////

type DotNetCleanSettings struct {
	ToolSettingsBase
	DotNetSettings
	Configuration   string
	OutputDirectory string
	Runtime         string
	Framework       string
	NoLogo          bool
}

func (tool *DotNetTool) DotNetClean(project string, settings *DotNetCleanSettings) error {
	if settings == nil {
		settings = &DotNetCleanSettings{}
	}
	args := []string{
		"clean",
	}
	args = goext.SliceAppendIf(args, project != "", project)
	args = goext.SliceAppendIf(args, settings.Configuration != "", "--configuration", settings.Configuration)
	args = goext.SliceAppendIf(args, settings.OutputDirectory != "", "--output", settings.OutputDirectory)
	args = goext.SliceAppendIf(args, settings.Runtime != "", "--runtime", settings.Runtime)
	args = goext.SliceAppendIf(args, settings.Framework != "", "--framework", settings.Framework)
	args = goext.SliceAppendIf(args, settings.NoLogo, "--nologo")
	args = append(args, settings.CustomArguments...)
	args = tool.addDotNetSettings(args, settings.DotNetSettings)

	return tool.run("dotnet", args, settings.ToolSettingsBase)
}

////////////////////////////////////////////////////////////
// Pack
////////////////////////////////////////////////////////////

type DotNetPackSettings struct {
	ToolSettingsBase
	DotNetSettings
	Configuration   string
	OutputDirectory string
	Runtime         string
	NoBuild         bool
	NoDependencies  bool
	NoRestore       bool
	NoLogo          bool
	IncludeSymbols  bool
	IncludeSource   bool
}

func (tool *DotNetTool) DotNetPack(project string, settings *DotNetPackSettings) error {
	if settings == nil {
		settings = &DotNetPackSettings{}
	}
	args := []string{
		"pack",
	}
	args = goext.SliceAppendIf(args, project != "", project)
	args = goext.SliceAppendIf(args, settings.Configuration != "", "--configuration", settings.Configuration)
	args = goext.SliceAppendIf(args, settings.OutputDirectory != "", "--output", settings.OutputDirectory)
	args = goext.SliceAppendIf(args, settings.Runtime != "", "--runtime", settings.Runtime)
	args = goext.SliceAppendIf(args, settings.NoBuild, "--no-build")
	args = goext.SliceAppendIf(args, settings.NoDependencies, "--no-dependencies")
	args = goext.SliceAppendIf(args, settings.NoRestore, "--no-restore")
	args = goext.SliceAppendIf(args, settings.NoLogo, "--nologo")
	args = goext.SliceAppendIf(args, settings.IncludeSymbols, "--include-symbols")
	args = goext.SliceAppendIf(args, settings.IncludeSource, "--include-source")
	args = append(args, settings.CustomArguments...)
	args = tool.addDotNetSettings(args, settings.DotNetSettings)

	return tool.run("dotnet", args, settings.ToolSettingsBase)
}

////////////////////////////////////////////////////////////
// Publish
////////////////////////////////////////////////////////////

type DotNetPublishSettings struct {
	ToolSettingsBase
	DotNetSettings
	Configuration   string
	OutputDirectory string
	Runtime         string
	Framework       string
	NoBuild         bool
	NoDependencies  bool
	NoRestore       bool
	NoLogo          bool
	Force           bool
}

func (tool *DotNetTool) DotNetPublish(path string, settings *DotNetPublishSettings) error {
	if settings == nil {
		settings = &DotNetPublishSettings{}
	}
	args := []string{
		"publish",
	}
	args = goext.SliceAppendIf(args, path != "", path)
	args = goext.SliceAppendIf(args, settings.Configuration != "", "--configuration", settings.Configuration)
	args = goext.SliceAppendIf(args, settings.OutputDirectory != "", "--output", settings.OutputDirectory)
	args = goext.SliceAppendIf(args, settings.Runtime != "", "--runtime", settings.Runtime)
	args = goext.SliceAppendIf(args, settings.Framework != "", "--framework", settings.Framework)
	args = goext.SliceAppendIf(args, settings.NoBuild, "--no-build")
	args = goext.SliceAppendIf(args, settings.NoDependencies, "--no-dependencies")
	args = goext.SliceAppendIf(args, settings.NoRestore, "--no-restore")
	args = goext.SliceAppendIf(args, settings.NoLogo, "--nologo")
	args = goext.SliceAppendIf(args, settings.Force, "--force")
	args = append(args, settings.CustomArguments...)
	args = tool.addDotNetSettings(args, settings.DotNetSettings)

	return tool.run("dotnet", args, settings.ToolSettingsBase)
}

////////////////////////////////////////////////////////////
// Restore
////////////////////////////////////////////////////////////

type DotNetRestoreSettings struct {
	ToolSettingsBase
	DotNetSettings
	PackagesDirectory   string
	Runtime             string
	NoCache             bool
	DisableParallel     bool
	IgnoreFailedSources bool
	NoDependencies      bool
	Force               bool
	Interactive         bool
}

func (tool *DotNetTool) DotNetRestore(root string, settings *DotNetRestoreSettings) error {
	if settings == nil {
		settings = &DotNetRestoreSettings{}
	}
	args := []string{
		"restore",
	}
	args = goext.SliceAppendIf(args, root != "", root)
	args = goext.SliceAppendIf(args, settings.PackagesDirectory != "", "--packages", settings.PackagesDirectory)
	args = goext.SliceAppendIf(args, settings.Runtime != "", "--runtime", settings.Runtime)
	args = goext.SliceAppendIf(args, settings.NoCache, "--no-cache")
	args = goext.SliceAppendIf(args, settings.DisableParallel, "--disable-parallel")
	args = goext.SliceAppendIf(args, settings.IgnoreFailedSources, "--ignore-failed-sources")
	args = goext.SliceAppendIf(args, settings.NoDependencies, "--no-dependencies")
	args = goext.SliceAppendIf(args, settings.Force, "--force")
	args = goext.SliceAppendIf(args, settings.Interactive, "--interactive")
	args = append(args, settings.CustomArguments...)
	args = tool.addDotNetSettings(args, settings.DotNetSettings)

	return tool.run("dotnet", args, settings.ToolSettingsBase)
}

////////////////////////////////////////////////////////////
// Run
////////////////////////////////////////////////////////////

type DotNetRunSettings struct {
	ToolSettingsBase
	DotNetSettings
	Configuration string
	Runtime       string
	Framework     string
	NoBuild       bool
	NoRestore     bool
}

func (tool *DotNetTool) DotNetRun(project string, settings *DotNetRunSettings) error {
	if settings == nil {
		settings = &DotNetRunSettings{}
	}
	args := []string{
		"run",
	}
	args = goext.SliceAppendIf(args, project != "", project)
	args = goext.SliceAppendIf(args, settings.Configuration != "", "--configuration", settings.Configuration)
	args = goext.SliceAppendIf(args, settings.Runtime != "", "--runtime", settings.Runtime)
	args = goext.SliceAppendIf(args, settings.Framework != "", "--framework", settings.Framework)
	args = goext.SliceAppendIf(args, settings.NoBuild, "--no-build")
	args = goext.SliceAppendIf(args, settings.NoRestore, "--no-restore")
	args = append(args, settings.CustomArguments...)
	args = tool.addDotNetSettings(args, settings.DotNetSettings)

	return tool.run("dotnet", args, settings.ToolSettingsBase)
}

////////////////////////////////////////////////////////////
// Test
////////////////////////////////////////////////////////////

type DotNetTestSettings struct {
	ToolSettingsBase
	DotNetSettings
	Configuration   string
	OutputDirectory string
	Runtime         string
	Framework       string
	Filter          string
	NoBuild         bool
	NoRestore       bool
	NoLogo          bool
	Blame           bool
}

func (tool *DotNetTool) DotNeTest(project string, settings *DotNetTestSettings) error {
	if settings == nil {
		settings = &DotNetTestSettings{}
	}
	args := []string{
		"test",
	}
	args = goext.SliceAppendIf(args, project != "", project)
	args = goext.SliceAppendIf(args, settings.Configuration != "", "--configuration", settings.Configuration)
	args = goext.SliceAppendIf(args, settings.OutputDirectory != "", "--output", settings.OutputDirectory)
	args = goext.SliceAppendIf(args, settings.Runtime != "", "--runtime", settings.Runtime)
	args = goext.SliceAppendIf(args, settings.Framework != "", "--framework", settings.Framework)
	args = goext.SliceAppendIf(args, settings.Filter != "", "--filter", settings.Filter)
	args = goext.SliceAppendIf(args, settings.NoBuild, "--no-build")
	args = goext.SliceAppendIf(args, settings.NoRestore, "--no-restore")
	args = goext.SliceAppendIf(args, settings.NoLogo, "--nologo")
	args = goext.SliceAppendIf(args, settings.Blame, "--blame")
	args = append(args, settings.CustomArguments...)
	args = tool.addDotNetSettings(args, settings.DotNetSettings)

	return tool.run("dotnet", args, settings.ToolSettingsBase)
}

////////////////////////////////////////////////////////////
// Internal Methods
////////////////////////////////////////////////////////////

func (tool *DotNetTool) addDotNetSettings(args []string, settings DotNetSettings) []string {
	args = goext.SliceAppendIf(args, settings.DiagnosticOutput, "--diagnostics")
	args = goext.SliceAppendIf(args, settings.Verbosity != DotNetVerbosityUndefined, "--verbosity", settings.Verbosity.String())
	return args
}
