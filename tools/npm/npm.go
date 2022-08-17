package npm

import (
	"os/exec"

	"github.com/roemer/gotaskr/execr"
	"github.com/roemer/gotaskr/goext"
)

type InitSettings struct {
	WorkingDirectory string
}

func Init(settings InitSettings) error {
	args := []string{
		"init",
		"-y",
	}
	cmd := exec.Command("npm", args...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(cmd)
}

type RunSettings struct {
	WorkingDirectory string
	Script           string
}

func Run(settings RunSettings) error {
	args := []string{
		"run",
		settings.Script,
	}
	cmd := exec.Command("npm", args...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(cmd)
}

func RunScript(script string) error {
	return Run(RunSettings{Script: script})
}

type CleanInstallSettings struct {
	WorkingDirectory string
	CacheDir         string
	NoAudit          bool
	PreferOffline    bool
}

func CleanInstall(settings CleanInstallSettings) error {
	args := []string{
		"ci",
		goext.Ternary(settings.NoAudit, "--no-audit", ""),
		goext.Ternary(settings.PreferOffline, "--prefer-offline", ""),
	}
	args = addCache(args, settings.CacheDir)
	cmd := exec.Command("npm", goext.RemoveEmpty(args)...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(cmd)
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

func Install(settings InstallSettings) error {
	args := []string{
		"install",
		settings.PackageSpec,
		goext.Ternary(settings.SaveProd, "--save-prod", ""),
		goext.Ternary(settings.SaveDev, "--save-dev", ""),
		goext.Ternary(settings.SaveOptional, "--save-optional", ""),
		goext.Ternary(settings.SaveExact, "--save-exact", ""),
		goext.Ternary(settings.NoAudit, "--no-audit", ""),
		goext.Ternary(settings.PreferOffline, "--prefer-offline", ""),
	}
	args = addCache(args, settings.CacheDir)
	cmd := exec.Command("npm", goext.RemoveEmpty(args)...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(cmd)
}

func addCache(args []string, cacheDir string) []string {
	if cacheDir != "" {
		args = append(args, "--cache", cacheDir)
	}
	return args
}

func addLogLevel(args []string, logLevel string) []string {
	if logLevel != "" {
		args = append(args, "--loglevel", logLevel)
	}
	return args
}
