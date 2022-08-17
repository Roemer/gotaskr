package mvn

import (
	"os/exec"

	"github.com/roemer/gotaskr/execr"
)

type MvnSettings struct {
	WorkingDirectory string
}

type CleanInstallSettings struct {
	MvnSettings
	Project string
}

func CleanInstall(settings CleanInstallSettings) error {
	args := []string{
		"clean",
		"install",
		"--projects",
		settings.Project,
		"--also-make",
		"--no-transfer-progress",
		"--batch-mode",
	}
	cmd := exec.Command("mvn", args...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(cmd)
}
