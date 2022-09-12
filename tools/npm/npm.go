package npm

import (
	"os/exec"

	"github.com/roemer/gotaskr/execr"
	"github.com/roemer/gotaskr/goext"
)

type InitSettings struct {
	WorkingDirectory string
}

func Init(settings *InitSettings, outputToConsole bool) error {
	if settings == nil {
		settings = &InitSettings{}
	}
	args := []string{
		"init",
		"-y",
	}
	cmd := exec.Command("npm", args...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(cmd, outputToConsole)
}

type RunSettings struct {
	WorkingDirectory string
	Script           string
}

func Run(settings *RunSettings, outputToConsole bool) error {
	if settings == nil {
		settings = &RunSettings{}
	}
	args := []string{
		"run",
		settings.Script,
	}
	cmd := exec.Command("npm", args...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(cmd, outputToConsole)
}

func RunScript(script string, outputToConsole bool) error {
	return Run(&RunSettings{Script: script}, outputToConsole)
}

type CleanInstallSettings struct {
	WorkingDirectory string
	CacheDir         string
	NoAudit          bool
	PreferOffline    bool
}

func CleanInstall(settings *CleanInstallSettings, outputToConsole bool) error {
	if settings == nil {
		settings = &CleanInstallSettings{}
	}
	args := []string{
		"ci",
	}
	args = goext.AddIf(args, settings.NoAudit, "--no-audit")
	args = goext.AddIf(args, settings.PreferOffline, "--prefer-offline")
	args = addCache(args, settings.CacheDir)
	cmd := exec.Command("npm", goext.RemoveEmpty(args)...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(cmd, outputToConsole)
}

type InstallSettings struct {
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

func Install(settings *InstallSettings, outputToConsole bool) error {
	if settings == nil {
		settings = &InstallSettings{}
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
	args = addCache(args, settings.CacheDir)
	cmd := exec.Command("npm", goext.RemoveEmpty(args)...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(cmd, outputToConsole)
}

func addCache(args []string, cacheDir string) []string {
	return goext.AddIf(args, cacheDir != "", "--cache", cacheDir)
}
