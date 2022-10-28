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
	Registry string
	Username string
	Password string
}

func (tool *DockerRegistryTool) Login(outputToConsole bool, settings *DockerLoginSettings) error {
	args := []string{
		"login",
	}
	args = append(args, "--username", settings.Username)
	args = append(args, "--password", settings.Password)
	args = goext.AddIf(args, settings.Registry != "", settings.Registry)

	cmd := exec.Command("docker", goext.RemoveEmpty(args)...)
	return execr.RunCommand(outputToConsole, cmd)
}

type DockerLogoutSettings struct {
	Registry string
}

func (tool *DockerRegistryTool) Logout(outputToConsole bool, settings *DockerLogoutSettings) error {
	args := []string{
		"logout",
	}
	args = goext.AddIf(args, settings.Registry != "", settings.Registry)

	cmd := exec.Command("docker", goext.RemoveEmpty(args)...)
	return execr.RunCommand(outputToConsole, cmd)
}
