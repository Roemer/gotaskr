package gttools

import (
	"os/exec"

	"github.com/roemer/gotaskr/execr"
	"github.com/roemer/gotaskr/goext"
)

// DockerRegistryTool provides access to the helper methods for Docker Registries.
type DockerRegistryTool struct {
}

type DockerLoginSettings struct {
	ToolSettingsBase
	Registry string
	Username string
	Password string
}

func (tool *DockerRegistryTool) Login(settings *DockerLoginSettings) error {
	args := []string{
		"login",
	}
	args = append(args, "--username", settings.Username)
	args = append(args, "--password", settings.Password)
	args = append(args, settings.CustomArguments...)
	args = goext.AddIf(args, settings.Registry != "", settings.Registry)

	cmd := exec.Command("docker", goext.RemoveEmpty(args)...)
	return execr.RunCommand(settings.OutputToConsole, cmd)
}

type DockerLogoutSettings struct {
	ToolSettingsBase
	Registry string
}

func (tool *DockerRegistryTool) Logout(settings *DockerLogoutSettings) error {
	args := []string{
		"logout",
	}
	args = append(args, settings.CustomArguments...)
	args = goext.AddIf(args, settings.Registry != "", settings.Registry)

	cmd := exec.Command("docker", goext.RemoveEmpty(args)...)
	return execr.RunCommand(settings.OutputToConsole, cmd)
}
