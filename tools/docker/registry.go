package docker

import (
	"os/exec"

	"github.com/roemer/gotaskr/execr"
	"github.com/roemer/gotaskr/goext"
)

type LoginSettings struct {
	Registry string
	Username string
	Password string
}

func Login(settings *LoginSettings, outputToConsole bool) error {
	args := []string{
		"login",
	}
	args = append(args, "--username", settings.Username)
	args = append(args, "--password", settings.Password)
	args = goext.AddIf(args, settings.Registry != "", settings.Registry)

	cmd := exec.Command("docker", goext.RemoveEmpty(args)...)
	return execr.RunCommand(cmd, outputToConsole)
}

type LogoutSettings struct {
	Registry string
}

func Logout(settings *LogoutSettings, outputToConsole bool) error {
	args := []string{
		"logout",
	}
	args = goext.AddIf(args, settings.Registry != "", settings.Registry)

	cmd := exec.Command("docker", goext.RemoveEmpty(args)...)
	return execr.RunCommand(cmd, outputToConsole)
}
