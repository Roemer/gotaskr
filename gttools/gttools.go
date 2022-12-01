// Package gttools provides helper methods for various tools.
package gttools

// ToolsClient provides typed access to the different tools.
type ToolsClient struct {
	Cypress *CypressTool
	Docker  *DockerTool
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
		EsLint:  CreateEsLintTool(),
		GitLab:  CreateGitLabTool(),
		JFrog:   CreateJFrogTool(),
		Mvn:     CreateMvnTool(),
		Npm:     CreateNpmTool(),
	}
}
