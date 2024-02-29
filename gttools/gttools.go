// Package gttools provides helper methods for various tools.
package gttools

import (
	"fmt"
	"strings"
)

// ToolsClient provides typed access to the different tools.
type ToolsClient struct {
	Cypress         *CypressTool
	DevContainerCli *DevContainerCliTool
	Docker          *DockerTool
	DotNet          *DotNetTool
	EsLint          *EsLintTool
	Flyway          *FlywayTool
	GitLab          *GitLabTool
	JFrog           *JFrogTool
	Mvn             *MvnTool
	Npm             *NpmTool
	Nx              *NxTool
}

// CreateToolsClient creates a new client to access the different tools.
func CreateToolsClient() *ToolsClient {
	return &ToolsClient{
		Cypress:         CreateCypressTool(),
		DevContainerCli: CreateDevContainerCliTool(),
		Docker:          CreateDockerTool(),
		DotNet:          CreateDotNetTool(),
		EsLint:          CreateEsLintTool(),
		Flyway:          CreateFlywayTool(),
		GitLab:          CreateGitLabTool(),
		JFrog:           CreateJFrogTool(),
		Mvn:             CreateMvnTool(),
		Npm:             CreateNpmTool(),
		Nx:              CreateNxTool(),
	}
}

// ToolSettingsBase are common settings usefull for all tools that run executables.
type ToolSettingsBase struct {
	WorkingDirectory string   // the path to use as working directory when running the tool
	OutputToConsole  bool     // flag to define if the output of the tool should be written into the console or not.
	CustomArguments  []string // list with custom arguments passed to the tool
}

// Customize adds a custom argument to the settings object.
func (s *ToolSettingsBase) Customize(setting ...string) *ToolSettingsBase {
	s.CustomArguments = append(s.CustomArguments, setting...)
	return s
}

// Ptr is a helper returns a pointer to v.
func Ptr[T any](v T) *T {
	return &v
}

// True value to use when a nullable bool is needed.
var True *bool = Ptr(true)

// False value to use when a nullable bool is needed.
var False *bool = Ptr(false)

//////////
// Internal helper methods
//////////

type addSettings struct {
	// Prefix before the effective value.
	prefix string
	// Suffix after the effective value.
	suffix string
	// Elements to add after the effective value.
	appendElements []string
	// Elements to add before the effective value.
	prependElements []string
}

// Adds a nullable boolean to the list if it is not nil
func addBoolean(slice []string, value *bool, settings addSettings) []string {
	if value != nil {
		if len(settings.prependElements) > 0 {
			slice = append(slice, settings.prependElements...)
		}
		slice = append(slice, fmt.Sprintf("%s%t%s", settings.prefix, *value, settings.suffix))
		if len(settings.appendElements) > 0 {
			slice = append(slice, settings.appendElements...)
		}
	}
	return slice
}

// Adds a nullable int to the list if it is not nil
func addInt(slice []string, value *int, settings addSettings) []string {
	if value != nil {
		if len(settings.prependElements) > 0 {
			slice = append(slice, settings.prependElements...)
		}
		slice = append(slice, fmt.Sprintf("%s%d%s", settings.prefix, *value, settings.suffix))
		if len(settings.appendElements) > 0 {
			slice = append(slice, settings.appendElements...)
		}
	}
	return slice
}

// Adds a string to the list if it has a length > 0
func addString(slice []string, value string, settings addSettings) []string {
	if len(value) > 0 {
		if len(settings.prependElements) > 0 {
			slice = append(slice, settings.prependElements...)
		}
		slice = append(slice, fmt.Sprintf("%s%s%s", settings.prefix, value, settings.suffix))
		if len(settings.appendElements) > 0 {
			slice = append(slice, settings.appendElements...)
		}
	}
	return slice
}

// Adds a string list to the list, separated by the separator
func addStringList(slice []string, values []string, separator string, settings addSettings) []string {
	if len(values) > 0 {
		if len(settings.prependElements) > 0 {
			slice = append(slice, settings.prependElements...)
		}
		slice = append(slice, fmt.Sprintf("%s%s%s", settings.prefix, strings.Join(values, separator), settings.suffix))
		if len(settings.appendElements) > 0 {
			slice = append(slice, settings.appendElements...)
		}
	}
	return slice
}
