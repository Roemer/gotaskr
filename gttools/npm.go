package gttools

import (
	"os/exec"

	"github.com/roemer/gotaskr/execr"
	"github.com/roemer/gotaskr/goext"
)

// NpmTool provides access to the helper methods for npm.
type NpmTool struct {
}

func CreateNpmTool() *NpmTool {
	return &NpmTool{}
}

// Init
type NpmInitSettings struct {
	WorkingDirectory string
}

func (tool *NpmTool) Init(outputToConsole bool, settings *NpmInitSettings) error {
	if settings == nil {
		settings = &NpmInitSettings{}
	}
	args := []string{
		"init",
		"-y",
	}
	cmd := exec.Command("npm", args...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(outputToConsole, cmd)
}

// Run
type NpmRunSettings struct {
	WorkingDirectory string
	Script           string
}

func (tool *NpmTool) Run(outputToConsole bool, settings *NpmRunSettings) error {
	if settings == nil {
		settings = &NpmRunSettings{}
	}
	args := []string{
		"run",
		settings.Script,
	}
	cmd := exec.Command("npm", args...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(outputToConsole, cmd)
}

func (tool *NpmTool) RunScript(outputToConsole bool, script string) error {
	return tool.Run(outputToConsole, &NpmRunSettings{Script: script})
}

// CleanInstall
type NpmCleanInstallSettings struct {
	WorkingDirectory string
	CacheDir         string
	NoAudit          bool
	PreferOffline    bool
}

func (tool *NpmTool) CleanInstall(outputToConsole bool, settings *NpmCleanInstallSettings) error {
	if settings == nil {
		settings = &NpmCleanInstallSettings{}
	}
	args := []string{
		"ci",
	}
	args = goext.AddIf(args, settings.NoAudit, "--no-audit")
	args = goext.AddIf(args, settings.PreferOffline, "--prefer-offline")
	args = tool.addCache(args, settings.CacheDir)
	cmd := exec.Command("npm", goext.RemoveEmpty(args)...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(outputToConsole, cmd)
}

// Install
type NpmInstallSettings struct {
	WorkingDirectory string
	CacheDir         string
	PackageSpec      string
	NoAudit          bool
	PreferOffline    bool
	SaveProd         bool
	SaveDev          bool
	SaveOptional     bool
	SaveExact        bool
}

func (tool *NpmTool) Install(outputToConsole bool, settings *NpmInstallSettings) error {
	if settings == nil {
		settings = &NpmInstallSettings{}
	}
	args := []string{
		"install",
		settings.PackageSpec,
	}
	args = goext.AddIf(args, settings.SaveProd, "--save-prod")
	args = goext.AddIf(args, settings.SaveDev, "--save-dev")
	args = goext.AddIf(args, settings.SaveOptional, "--save-optional")
	args = goext.AddIf(args, settings.SaveExact, "--save-exact")
	args = goext.AddIf(args, settings.NoAudit, "--no-audit")
	args = goext.AddIf(args, settings.PreferOffline, "--prefer-offline")
	args = tool.addCache(args, settings.CacheDir)
	cmd := exec.Command("npm", goext.RemoveEmpty(args)...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(outputToConsole, cmd)
}

// Bin
type NpmBinSettings struct {
	WorkingDirectory string
	Global           bool
}

func (tool *NpmTool) Bin(settings *NpmBinSettings) (string, error) {
	if settings == nil {
		settings = &NpmBinSettings{}
	}
	args := []string{
		"bin",
	}
	args = goext.AddIf(args, settings.Global, "--global")
	cmd := exec.Command("npm", goext.RemoveEmpty(args)...)
	cmd.Dir = settings.WorkingDirectory

	stdout, _, err := execr.RunCommandGetOutput(false, cmd)
	return stdout, err
}

// Internal Methods
func (tool *NpmTool) addCache(args []string, cacheDir string) []string {
	return goext.AddIf(args, cacheDir != "", "--cache", cacheDir)
}
