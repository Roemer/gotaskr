package gttools

import (
	"github.com/roemer/goext"
)

// DockerRegistryTool provides access to the helper methods for Docker Registries.
type DockerRegistryTool struct {
	ToolBase
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
	args = goext.SliceAppendIf(args, settings.Registry != "", settings.Registry)

	return tool.run("docker", args, settings.ToolSettingsBase)
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
	args = goext.SliceAppendIf(args, settings.Registry != "", settings.Registry)

	return tool.run("docker", args, settings.ToolSettingsBase)
}
