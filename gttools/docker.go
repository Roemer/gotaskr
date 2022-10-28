package gttools

// DockerTool provides access to the helper methods for tools around Docker.
type DockerTool struct {
	Image    *DockerImageTool
	Registry *DockerRegistryTool
}

func CreateDockerTool() *DockerTool {
	return &DockerTool{
		Image:    &DockerImageTool{},
		Registry: &DockerRegistryTool{},
	}
}
