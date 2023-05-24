// Package gttools provides helper methods for various tools.
package gttools

// ToolsClient provides typed access to the different tools.
type ToolsClient struct {
	Cypress *CypressTool
	Docker  *DockerTool
	DotNet  *DotNetTool
	EsLint  *EsLintTool
	GitLab  *GitLabTool
	JFrog   *JFrogTool
	Mvn     *MvnTool
	Npm     *NpmTool
}

// CreateToolsClient creates a new client to access the different tools.
func CreateToolsClient() *ToolsClient {
	return &ToolsClient{
		Cypress: CreateCypressTool(),
		Docker:  CreateDockerTool(),
		DotNet:  CreateDotNetTool(),
		EsLint:  CreateEsLintTool(),
		GitLab:  CreateGitLabTool(),
		JFrog:   CreateJFrogTool(),
		Mvn:     CreateMvnTool(),
		Npm:     CreateNpmTool(),
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
