// Package gttools provides helper methods for various tools.
package gttools

// ToolsClient provides typed access to the different tools.
type ToolsClient struct {
	Cypress *CypressTool
	Docker  *DockerTool
	DotNet  *DotNetTool
	EsLint  *EsLintTool
	Flyway  *FlywayTool
	GitLab  *GitLabTool
	JFrog   *JFrogTool
	Mvn     *MvnTool
	Npm     *NpmTool
	Nx      *NxTool
}

// CreateToolsClient creates a new client to access the different tools.
func CreateToolsClient() *ToolsClient {
	return &ToolsClient{
		Cypress: CreateCypressTool(),
		Docker:  CreateDockerTool(),
		DotNet:  CreateDotNetTool(),
		EsLint:  CreateEsLintTool(),
		Flyway:  CreateFlywayTool(),
		GitLab:  CreateGitLabTool(),
		JFrog:   CreateJFrogTool(),
		Mvn:     CreateMvnTool(),
		Npm:     CreateNpmTool(),
		Nx:      CreateNxTool(),
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
