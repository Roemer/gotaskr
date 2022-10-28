package gttools

// JFrogTool provides access to the helper methods for tools from JFrog.
type JFrogTool struct {
	Artifactory *ArtifactoryTool
}

func CreateJFrogTool() *JFrogTool {
	return &JFrogTool{
		Artifactory: &ArtifactoryTool{},
	}
}
