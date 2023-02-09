package gttools

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/roemer/gotaskr/execr"
	"github.com/roemer/gotaskr/goext"
)

// DockerImageTool provides access to the helper methods for Docker Images.
type DockerImageTool struct {
}

type DockerBuildSettings struct {
	ToolSettingsBase
	Dockerfile  string
	ContextPath string
	Tags        []string
	Labels      []string
	BuildArgs   []string
}

func (settings *DockerBuildSettings) AddTags(tags ...string) *DockerBuildSettings {
	for _, entry := range tags {
		settings.Tags = goext.AppendIfMissing(settings.Tags, entry)
	}
	return settings
}

func (settings *DockerBuildSettings) AddLabels(labels ...string) *DockerBuildSettings {
	for _, entry := range labels {
		settings.Labels = goext.AppendIfMissing(settings.Labels, entry)
	}
	return settings
}

func (settings *DockerBuildSettings) AddBuildArgs(buildArgs ...string) *DockerBuildSettings {
	for _, entry := range buildArgs {
		settings.BuildArgs = goext.AppendIfMissing(settings.BuildArgs, entry)
	}
	return settings
}

func (tool *DockerImageTool) Build(settings *DockerBuildSettings) error {
	args := []string{
		"build",
	}
	args = goext.AddIf(args, settings.Dockerfile != "", "--file", settings.Dockerfile)
	for _, entry := range settings.Tags {
		args = append(args, "--tag", entry)
	}
	for _, entry := range settings.Labels {
		args = append(args, "--label", entry)
	}
	for _, entry := range settings.BuildArgs {
		args = append(args, "--build-arg", entry)
	}
	args = append(args, settings.CustomArguments...)
	args = append(args, goext.Ternary(settings.ContextPath == "", ".", settings.ContextPath))

	cmd := exec.Command("docker", goext.RemoveEmpty(args)...)
	cmd.Dir = settings.WorkingDirectory
	cmd.Stdin = os.Stdin
	return execr.RunCommand(settings.OutputToConsole, cmd)
}

type DockerSaveSettings struct {
	ToolSettingsBase
	OutputFile     string
	ImageReference string
}

func (tool *DockerImageTool) Save(settings *DockerSaveSettings) error {
	args := []string{
		"save",
	}
	args = goext.AddIf(args, settings.OutputFile != "", "--output", settings.OutputFile)
	args = append(args, settings.CustomArguments...)
	args = append(args, settings.ImageReference)

	// Make sure the directory exists
	if settings.OutputFile != "" {
		if err := os.MkdirAll(filepath.Dir(settings.OutputFile), os.ModePerm); err != nil {
			return err
		}
	}

	cmd := exec.Command("docker", goext.RemoveEmpty(args)...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(settings.OutputToConsole, cmd)
}

type DockerLoadSettings struct {
	ToolSettingsBase
	InputFile string
}

func (tool *DockerImageTool) Load(settings *DockerLoadSettings) ([]string, error) {
	args := []string{
		"load",
	}
	args = goext.AddIf(args, settings.InputFile != "", "--input", settings.InputFile)
	args = append(args, settings.CustomArguments...)

	cmd := exec.Command("docker", goext.RemoveEmpty(args)...)
	cmd.Dir = settings.WorkingDirectory
	stdout, _, err := execr.RunCommandGetOutput(settings.OutputToConsole, cmd)

	// Parse out all loaded images
	loadedImages := []string{}
	re := regexp.MustCompile(`Loaded image: (.*)`)
	matched := re.FindAllStringSubmatch(stdout, -1)
	for _, match := range matched {
		loadedImages = append(loadedImages, match[1])
	}

	// Return the loaded images
	return loadedImages, err
}

type DockerPushSettings struct {
	ToolSettingsBase
	AllTags        bool
	ImageReference string
}

func (tool *DockerImageTool) Push(settings *DockerPushSettings) error {
	args := []string{
		"push",
	}
	args = goext.AddIf(args, settings.AllTags, "--all-tags")
	args = append(args, settings.CustomArguments...)
	args = append(args, settings.ImageReference)

	cmd := exec.Command("docker", goext.RemoveEmpty(args)...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(settings.OutputToConsole, cmd)
}
