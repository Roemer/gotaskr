package gttools

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/roemer/gotaskr/execr"
	"github.com/roemer/gotaskr/goext"
)

type CypressTool struct {
}

func CreateCypressTool() *CypressTool {
	return &CypressTool{}
}

const (
	CypressBrowserChrome   = "chrome"
	CypressBrowserChromium = "chromium"
	CypressBrowserEdge     = "edge"
	CypressBrowserElectron = "electron"
	CypressBrowserFirefox  = "firefox"
)

// CypressRunSettings defines the settings to use for running Cypress.
// Also see https://docs.cypress.io/guides/guides/command-line#cypress-run
type CypressRunSettings struct {
	ToolSettingsBase
	Browser         string            // defines the browser to launch like chrome, chromium, edge, electron, firefox. Alternatively a path to an executable.
	CiBuildId       string            // the unique id to group tests together.
	Component       bool              // flag to define if component tests should run.
	Config          string            // specify the config to use. Defined as key value pairs, comma separated. Can alsop be a stringified json object.
	ConfigFile      string            // the path to a config file to use.
	E2e             bool              // flag to define if end to end tests should run (default).
	Env             map[string]string // environment variables to use.
	Group           string            // name of the group to group tests together.
	Headed          bool              // flag to indicate if the browser should be displayed.
	Headless        bool              // flag to indicate if the browser should be hidden (default).
	RecordKey       string            // the record key used to record the results to the Cypress Cloud.
	NoExit          bool              // flag to indicate if Cypress should stay open after the run.
	Parallel        bool              // flag to indicate if the tests should run in parallel across multiple machines.
	Port            int               // override the default port.
	Project         string            // the path to a specific project to run.
	Quiet           bool              // flag to indicate the quite mode where no output is passed to stdout.
	Record          bool              // flag to indicate if the tests shouldbe recorded or not.
	Reporter        string            // define the reporter to use. Can be any of the mocha, cypress or a custom reporter.
	ReporterOptions string            // specify the reporter options to use as key value pairs, comma separated. Can also be a stringified json object.
	Specs           []string          // define the spec file(s) to run.
	Tags            []string          // add tags to identify a run.
}

// AddEnv adds an environment value to the Cypress settings.
func (settings *CypressRunSettings) AddEnv(key string, value string) *CypressRunSettings {
	settings.Env[key] = value
	return settings
}

// AddSpecs adds one or more specs to run to the Cypress settings.
func (settings *CypressRunSettings) AddSpecs(specs ...string) *CypressRunSettings {
	for _, entry := range specs {
		settings.Specs = goext.AppendIfMissing(settings.Specs, entry)
	}
	return settings
}

// AddTags adds one or more tagst to the Cypress settings.
func (settings *CypressRunSettings) AddTags(tags ...string) *CypressRunSettings {
	for _, entry := range tags {
		settings.Tags = goext.AppendIfMissing(settings.Tags, entry)
	}
	return settings
}

// CypressRun allows you to run Cypress with a defined binary. Usefull for example when using cy2.
func (tool *CypressTool) CypressRun(cypressBinPath string, settings *CypressRunSettings) error {
	args := settings.buildCliArguments()
	cmd := exec.Command(cypressBinPath, goext.RemoveEmpty(args)...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(settings.OutputToConsole, cmd)
}

// CypressRunWithNpx runs Cypress with npx.
func (tool *CypressTool) CypressRunWithNpx(settings *CypressRunSettings) error {
	args := []string{
		"cypress",
	}
	args = append(args, settings.buildCliArguments()...)
	cmd := exec.Command("npx", goext.RemoveEmpty(args)...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(settings.OutputToConsole, cmd)
}

// CypressRunWitYarn runs Cypress with Yarn.
func (tool *CypressTool) CypressRunWitYarn(settings *CypressRunSettings) error {
	args := []string{
		"cypress",
	}
	args = append(args, settings.buildCliArguments()...)
	cmd := exec.Command("yarn", goext.RemoveEmpty(args)...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(settings.OutputToConsole, cmd)
}

////////////////////////////////////////////////////////////
// Internal Methods
////////////////////////////////////////////////////////////

func (settings *CypressRunSettings) buildCliArguments() []string {
	args := []string{
		"run",
	}
	args = goext.AddIf(args, settings.Browser != "", "--browser", settings.Browser)
	args = goext.AddIf(args, settings.CiBuildId != "", "--ci-build-id", settings.CiBuildId)
	args = goext.AddIf(args, settings.Component, "--component")
	args = goext.AddIf(args, settings.Config != "", "--config", settings.Config)
	args = goext.AddIf(args, settings.ConfigFile != "", "--config-file", settings.ConfigFile)
	args = goext.AddIf(args, settings.E2e, "--e2e")
	args = goext.AddIf(args, len(settings.Env) > 0, "--env", goext.ConvertMapToSingleString(settings.Env, "=", ","))
	args = goext.AddIf(args, settings.Group != "", "--group", settings.Group)
	args = goext.AddIf(args, settings.Headed, "--headed")
	args = goext.AddIf(args, settings.Headless, "--headless")
	args = goext.AddIf(args, settings.RecordKey != "", "--key", settings.RecordKey)
	args = goext.AddIf(args, settings.NoExit, "--no-exit")
	args = goext.AddIf(args, settings.Parallel, "--parallel")
	args = goext.AddIf(args, settings.Port > 0, "--port", strconv.Itoa(settings.Port))
	args = goext.AddIf(args, settings.Project != "", "--project", settings.Project)
	args = goext.AddIf(args, settings.Quiet, "--quiet")
	args = goext.AddIf(args, settings.Record, "--record")
	args = goext.AddIf(args, settings.Reporter != "", "--reporter", settings.Reporter)
	args = goext.AddIf(args, settings.ReporterOptions != "", "--reporter-options", settings.ReporterOptions)
	args = goext.AddIf(args, len(settings.Specs) > 0, "--spec", strings.Join(settings.Specs, ","))
	args = goext.AddIf(args, len(settings.Tags) > 0, "--tag", strings.Join(settings.Tags, ","))
	args = append(args, settings.CustomArguments...)
	return args
}
