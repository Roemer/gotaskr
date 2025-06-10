package gttools

import (
	"bufio"
	"fmt"
	"os/exec"
	"slices"
	"strings"

	"github.com/roemer/gotaskr/execr"
	"github.com/roemer/gotaskr/goext"
)

// NxTool provides access to the helper methods for Nx.
type NxTool struct {
}

func CreateNxTool() *NxTool {
	return &NxTool{}
}

// NxRunType is used to define how Nx should be invoked.
type NxRunType int

const (
	NxRunTypeDirect = iota // Directly invokes Nx, for example when it is globally installed.
	NxRunTypeNpx           // Runs Nx via npx.
	NxRunTypeYarn          // Runs Nx via yarn.
	NxRunTypePnpx          // Runs Nx via pnpx.
)

// NxOutputStyle is used to define how Nx emits outputs tasks logs.
type NxOutputStyle string

func (c NxOutputStyle) String() string {
	return string(c)
}

const (
	NxOutputStyleDynamic               = "dynamic"
	NxOutputStyleStatic                = "static"
	NxOutputStyleStream                = "stream"
	NxOutputStyleStreamWithoutPrefixes = "stream-without-prefixes"
)

// Settings for affected: https://nx.dev/nx-api/nx/documents/affected
type NxAffectedSettings struct {
	ToolSettingsBase
	Base           string        // Base of the current branch (usually main).
	Batch          bool          // Run task(s) in batches for executors which support batches.
	Configuration  string        // This is the configuration to use when performing tasks on projects.
	Exclude        []string      // Exclude certain projects from being processed.
	Files          []string      // Change the way Nx is calculating the affected command by providing directly changed files.
	Head           string        // Latest commit of the current branch (usually HEAD).
	NxBail         bool          // Stop command execution after the first failed task.
	NxIgnoreCycles bool          // Ignore cycles in the task graph.
	OutputStyle    NxOutputStyle // Defines how Nx emits outputs tasks logs.
	Parallel       string        // Max number of parallel processes. Set to "false" to disable.
	Runner         string        // This is the name of the tasks runner configured in nx.json.
	SkipNxCache    bool          // Rerun the tasks even when the results are available in the cache.
	Targets        []string      // Tasks to run for affected projects.
	Uncommitted    bool          // Use uncommitted changes.
	Untracked      bool          // Use untracked changes.
	Verbose        bool          // Prints additional information about the commands (e.g., stack traces).
}

// Affected allows running target(s) for affected projects.
func (tool *NxTool) Affected(runType NxRunType, settings NxAffectedSettings) error {
	args := []string{}
	args = goext.AppendIf(args, len(settings.Base) > 0, "--base="+settings.Base)
	args = goext.AppendIf(args, settings.Batch, "--batch")
	args = goext.AppendIf(args, len(settings.Configuration) > 0, "--configuration="+settings.Configuration)
	args = goext.AppendIf(args, len(settings.Exclude) > 0, "--exclude="+strings.Join(settings.Exclude, ","))
	args = goext.AppendIf(args, len(settings.Files) > 0, "--files="+strings.Join(settings.Files, ","))
	args = goext.AppendIf(args, len(settings.Head) > 0, "--head="+settings.Head)
	args = goext.AppendIf(args, settings.NxBail, "--nxBail")
	args = goext.AppendIf(args, settings.NxIgnoreCycles, "--nxIgnoreCycles")
	args = goext.AppendIf(args, len(settings.OutputStyle) > 0, "--output-style="+settings.OutputStyle.String())
	args = goext.AppendIf(args, len(settings.Parallel) > 0, "--parallel="+settings.Parallel)
	args = goext.AppendIf(args, len(settings.Runner) > 0, "--runner="+settings.Runner)
	args = goext.AppendIf(args, settings.SkipNxCache, "--skipNxCache")
	args = goext.AppendIf(args, len(settings.Targets) > 0, "--targets="+strings.Join(settings.Targets, ","))
	args = goext.AppendIf(args, settings.Uncommitted, "--uncommitted")
	args = goext.AppendIf(args, settings.Untracked, "--untracked")
	args = goext.AppendIf(args, settings.Verbose, "--verbose")
	return tool.RunCommand(runType, &settings.ToolSettingsBase, "affected", args...)
}

// Settings for run: https://nx.dev/nx-api/nx/documents/run
type NxRunSettings struct {
	ToolSettingsBase
	Project       string
	Target        string
	Configuration string
}

// Runs a target defined for your project.
func (tool *NxTool) Run(runType NxRunType, settings NxRunSettings) error {
	args := []string{}
	if len(settings.Project) > 0 && len(settings.Target) > 0 {
		args = append(args, fmt.Sprintf("%s:%s", settings.Project, settings.Target))
	} else if len(settings.Project) > 0 {
		args = append(args, settings.Project)
	} else if len(settings.Target) > 0 {
		args = append(args, settings.Target)
	}
	args = goext.AppendIf(args, len(settings.Configuration) > 0, "--configuration="+settings.Configuration)
	return tool.RunCommand(runType, &settings.ToolSettingsBase, "run", args...)
}

// Settings for run-many: https://nx.dev/nx-api/nx/documents/run-many
type NxRunManySettings struct {
	ToolSettingsBase
	Batch          bool          // Run task(s) in batches for executors which support batches.
	Configuration  string        // This is the configuration to use when performing tasks on projects.
	Exclude        []string      // Exclude certain projects from being processed.
	NxBail         bool          // Stop command execution after the first failed task.
	NxIgnoreCycles bool          // Ignore cycles in the task graph.
	OutputStyle    NxOutputStyle // Defines how Nx emits outputs tasks logs.
	Parallel       string        // Max number of parallel processes. Set to "false" to disable.
	Projects       []string      // Projects to run.
	Runner         string        // This is the name of the tasks runner configured in nx.json.
	SkipNxCache    bool          // Rerun the tasks even when the results are available in the cache.
	Targets        []string      // Tasks to run for affected projects.
	Verbose        bool          // Prints additional information about the commands (e.g., stack traces).
}

// RunMany runs target(s) for multiple listed projects.
func (tool *NxTool) RunMany(runType NxRunType, settings NxRunManySettings) error {
	args := []string{}
	args = goext.AppendIf(args, settings.Batch, "--batch")
	args = goext.AppendIf(args, len(settings.Configuration) > 0, "--configuration="+settings.Configuration)
	args = goext.AppendIf(args, len(settings.Exclude) > 0, "--exclude="+strings.Join(settings.Exclude, ","))
	args = goext.AppendIf(args, settings.NxBail, "--nxBail")
	args = goext.AppendIf(args, settings.NxIgnoreCycles, "--nxIgnoreCycles")
	args = goext.AppendIf(args, len(settings.OutputStyle) > 0, "--output-style="+settings.OutputStyle.String())
	args = goext.AppendIf(args, len(settings.Parallel) > 0, "--parallel="+settings.Parallel)
	args = goext.AppendIf(args, len(settings.Projects) > 0, "--projects="+strings.Join(settings.Projects, ","))
	args = goext.AppendIf(args, len(settings.Runner) > 0, "--runner="+settings.Runner)
	args = goext.AppendIf(args, settings.SkipNxCache, "--skipNxCache")
	args = goext.AppendIf(args, len(settings.Targets) > 0, "--targets="+strings.Join(settings.Targets, ","))
	args = goext.AppendIf(args, settings.Verbose, "--verbose")
	return tool.RunCommand(runType, &settings.ToolSettingsBase, "run-many", args...)
}

// Settings for show project: https://nx.dev/nx-api/nx/documents/show#project
type NxShowProjectSettings struct {
	ToolSettingsBase
	ProjectName string
}

// ShowProject shows the projects configuration and returns it as a string.
func (tool *NxTool) ShowProject(runType NxRunType, settings NxShowProjectSettings) (string, error) {
	args := []string{"project", settings.ProjectName}
	stdout, stderr, err := tool.RunCommandGetOutput(runType, &settings.ToolSettingsBase, "show", args...)
	if err != nil {
		return stdout, fmt.Errorf("nx show project failed: %s", stdout+" "+stderr)
	}
	return stdout, nil
}

// Settings for show projects: https://nx.dev/nx-api/nx/documents/show#projects
type NxShowProjectsSettings struct {
	ToolSettingsBase
	Affected    bool     // Show only affected projects.
	Base        string   // Base of the current branch (usually main).
	Exclude     []string // Exclude certain projects from being processed.
	Files       []string // Change the way Nx is calculating the affected command by providing directly changed files.
	Head        string   // Latest commit of the current branch (usually HEAD).
	Projects    []string // Show only projects that match a given pattern.
	Type        string   // Select only projects of the given type. Choices: [app, lib, e2e]
	Uncommitted bool     // Use uncommitted changes.
	Untracked   bool     // Use untracked changes.
	WithTarget  string   // Show only projects that have a specific target.
}

// ShowProjects returns projects according to the given criteria.
func (tool *NxTool) ShowProjects(runType NxRunType, settings NxShowProjectsSettings) ([]string, error) {
	args := []string{"projects"}
	args = goext.AppendIf(args, settings.Affected, "--affected")
	args = goext.AppendIf(args, len(settings.Base) > 0, "--base="+settings.Base)
	args = goext.AppendIf(args, len(settings.Exclude) > 0, "--exclude="+strings.Join(settings.Exclude, ","))
	args = goext.AppendIf(args, len(settings.Files) > 0, "--files="+strings.Join(settings.Files, ","))
	args = goext.AppendIf(args, len(settings.Head) > 0, "--head="+settings.Head)
	args = goext.AppendIf(args, len(settings.Projects) > 0, "--projects="+strings.Join(settings.Projects, ","))
	args = goext.AppendIf(args, len(settings.Type) > 0, "--type="+settings.Type)
	args = goext.AppendIf(args, settings.Uncommitted, "--uncommitted")
	args = goext.AppendIf(args, settings.Untracked, "--untracked")
	args = goext.AppendIf(args, len(settings.WithTarget) > 0, "--withTarget="+settings.WithTarget)

	stdout, stderr, err := tool.RunCommandGetOutput(runType, &settings.ToolSettingsBase, "show", args...)
	if err != nil {
		return nil, fmt.Errorf("nx show projects failed: %s", stdout+" "+stderr)
	}

	projects := []string{}
	sc := bufio.NewScanner(strings.NewReader(stdout))
	for sc.Scan() {
		projects = append(projects, strings.TrimSpace(sc.Text()))
	}
	slices.Sort(projects)

	return projects, nil
}

// RunCommand is a generic runner for any Nx command.
func (tool *NxTool) RunCommand(runType NxRunType, settings *ToolSettingsBase, command string, args ...string) error {
	cmd := tool.prepareCommand(runType, settings, command, args...)
	return execr.RunCommand(cmd, execr.WithConsoleOutput(settings.OutputToConsole))
}

// RunCommandGetOutput is a generic runner for any Nx command which also returns the stdout and stderr.
func (tool *NxTool) RunCommandGetOutput(runType NxRunType, settings *ToolSettingsBase, command string, args ...string) (string, string, error) {
	cmd := tool.prepareCommand(runType, settings, command, args...)
	return execr.RunCommandGetOutput(cmd, execr.WithConsoleOutput(settings.OutputToConsole))
}

func (tool *NxTool) prepareCommand(runType NxRunType, settings *ToolSettingsBase, command string, args ...string) *exec.Cmd {
	if settings == nil {
		settings = &ToolSettingsBase{}
	}

	// Add the command
	args = goext.Prepend(args, command)

	// Choose the runner type
	var bin string
	switch runType {
	default:
	case NxRunTypeDirect:
		bin = "nx"
	case NxRunTypeNpx:
		bin = "npx"
		args = goext.Prepend(args, "nx")
	case NxRunTypeYarn:
		bin = "yarn"
		args = goext.Prepend(args, "nx")
	case NxRunTypePnpx:
		bin = "pnpx"
		args = goext.Prepend(args, "nx")
	}

	args = append(args, settings.CustomArguments...)

	// Add the arguments and run the command
	cmd := exec.Command(bin, args...)
	cmd.Dir = settings.WorkingDirectory
	return cmd
}
