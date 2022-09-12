package docker

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/roemer/gotaskr/execr"
	"github.com/roemer/gotaskr/goext"
)

type BuildSettings struct {
	WorkingDirectory string
	Dockerfile       string
	ContextPath      string
	Tags             []string
	Labels           []string
	BuildArgs        []string
}

func (settings *BuildSettings) AddTags(tags ...string) *BuildSettings {
	for _, entry := range tags {
		settings.Tags = goext.AppendIfMissing(settings.Tags, entry)
	}
	return settings
}

func (settings *BuildSettings) AddLabels(labels ...string) *BuildSettings {
	for _, entry := range labels {
		settings.Labels = goext.AppendIfMissing(settings.Labels, entry)
	}
	return settings
}

func (settings *BuildSettings) AddBuildArgs(buildArgs ...string) *BuildSettings {
	for _, entry := range buildArgs {
		settings.BuildArgs = goext.AppendIfMissing(settings.BuildArgs, entry)
	}
	return settings
}

func Build(settings *BuildSettings, outputToConsole bool) error {
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
	args = append(args, goext.Ternary(settings.ContextPath == "", ".", settings.ContextPath))

	cmd := exec.Command("docker", goext.RemoveEmpty(args)...)
	cmd.Dir = settings.WorkingDirectory
	cmd.Stdin = os.Stdin
	return execr.RunCommand(cmd, outputToConsole)
}

type SaveSettings struct {
	WorkingDirectory string
	OutputFile       string
	ImageReference   string
}

func Save(settings *SaveSettings, outputToConsole bool) error {
	args := []string{
		"save",
	}
	args = goext.AddIf(args, settings.OutputFile != "", "--output", settings.OutputFile)
	args = append(args, settings.ImageReference)

	// Make sure the directory exists
	if settings.OutputFile != "" {
		os.MkdirAll(filepath.Dir(settings.OutputFile), os.ModePerm)
	}

	cmd := exec.Command("docker", goext.RemoveEmpty(args)...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(cmd, outputToConsole)
}

type LoadSettings struct {
	WorkingDirectory string
	InputFile        string
}

func Load(settings *LoadSettings, outputToConsole bool) ([]string, error) {
	args := []string{
		"load",
	}
	args = goext.AddIf(args, settings.InputFile != "", "--input", settings.InputFile)

	cmd := exec.Command("docker", goext.RemoveEmpty(args)...)
	cmd.Dir = settings.WorkingDirectory
	stdout, _, err := execr.RunCommandGetOutput(cmd, outputToConsole)

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

type PushSettings struct {
	WorkingDirectory string
	ImageReference   string
}

func Push(settings *PushSettings, outputToConsole bool) error {
	args := []string{
		"push",
	}
	args = append(args, settings.ImageReference)

	cmd := exec.Command("docker", goext.RemoveEmpty(args)...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(cmd, outputToConsole)
}
