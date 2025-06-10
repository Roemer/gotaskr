package gttools

import (
	"os/exec"

	"github.com/roemer/gotaskr/execr"
	"github.com/roemer/gotaskr/goext"
)

// This tool allows interacting with the Dev Container CLI.
type DevContainerCliTool struct {
}

func CreateDevContainerCliTool() *DevContainerCliTool {
	return &DevContainerCliTool{}
}

////////// Build

type DevContainerCliBuildSettings struct {
	ToolSettingsBase
	// Host path to a directory that is intended to be persisted and share state between sessions.
	UserDataFolder string
	// Docker CLI path.
	DockerPath string
	// Docker Compose CLI path.
	DockerComposePath string
	// Workspace folder path. The devcontainer.json will be looked up relative to this path.
	WorkspaceFolder string
	// devcontainer.json path.
	Config string
	// Log level.
	LogLevel DevContainerCliLogLevel
	// Log format.
	LogFormat DevContainerCliLogFormat
	// Builds the image with `--no-cache`.
	NoCache *bool
	// Image name.
	ImageNames []string
	// Additional image to use as potential layer cache.
	CachesFrom []string
	// A destination of buildx cache.
	CacheTo string
	// Control whether BuildKit should be used.
	Buildkit DevContainerCliBuildkit
	// Set target platforms.
	Platform string
	// Push to a container registry.
	Push *bool
	// Provide key and value configuration that adds metadata to an image
	Labels []string
	//  Overrides the default behavior to load built images into the local docker registry. Valid options are the same ones provided to the --output option of docker buildx build.
	Output string
	// Additional features to apply to the dev container (JSON as per "features" section in devcontainer.json).
	AdditionalFeatures string
}

func (tool *DevContainerCliTool) Build(settings *DevContainerCliBuildSettings) error {
	args := []string{
		"build",
	}

	args = addString(args, settings.UserDataFolder, addSettings{prependElements: []string{"--user-data-folder"}})
	args = addString(args, settings.DockerPath, addSettings{prependElements: []string{"--docker-path"}})
	args = addString(args, settings.DockerComposePath, addSettings{prependElements: []string{"--docker-compose-path"}})
	args = addString(args, settings.WorkspaceFolder, addSettings{prependElements: []string{"--workspace-folder"}})
	args = addString(args, settings.Config, addSettings{prependElements: []string{"--config"}})
	args = addString(args, string(settings.LogLevel), addSettings{prependElements: []string{"--log-level"}})
	args = addString(args, string(settings.LogFormat), addSettings{prependElements: []string{"--log-format"}})
	args = addBoolean(args, settings.NoCache, addSettings{prependElements: []string{"--no-cache"}})
	args = addStringList(args, settings.ImageNames, addSettings{prependElements: []string{"--image-name"}, handleEachListItemSeparately: true})
	args = addStringList(args, settings.CachesFrom, addSettings{prependElements: []string{"--cache-from"}, handleEachListItemSeparately: true})
	args = addString(args, settings.CacheTo, addSettings{prependElements: []string{"--cache-to"}})
	args = addString(args, string(settings.Buildkit), addSettings{prependElements: []string{"--buildkit"}})
	args = addString(args, settings.Platform, addSettings{prependElements: []string{"--platform"}})
	args = addBoolean(args, settings.Push, addSettings{prependElements: []string{"--push"}})
	args = addStringList(args, settings.Labels, addSettings{prependElements: []string{"--label"}, handleEachListItemSeparately: true})
	args = addString(args, settings.Output, addSettings{prependElements: []string{"--output"}})
	args = addString(args, settings.AdditionalFeatures, addSettings{prependElements: []string{"--additional-features"}})

	args = append(args, settings.CustomArguments...)
	cmd := exec.Command("devcontainer", goext.RemoveEmpty(args)...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommandO(cmd, execr.WithConsoleOutput(settings.OutputToConsole))
}

////////// Up

type DevContainerCliUpSettings struct {
	ToolSettingsBase
	// Docker CLI path.
	DockerPath string
	// Docker Compose CLI path.
	DockerComposePath string
	// Container data folder where user data inside the container will be stored.
	ContainerDataFolder string
	// Container system data folder where system data inside the container will be stored.
	ContainerSystemDataFolder string
	// Workspace folder path. The devcontainer.json will be looked up relative to this path.
	WorkspaceFolder string
	// Workspace mount consistency.
	WorkspaceMountConsistency DevContainerCliWorkspaceMountConsistency
	// Availability of GPUs in case the dev container requires any. `all` expects a GPU to be available.
	GpuAvailability DevContainerCliGpuAvailability
	// Mount the workspace using its Git root.
	MountWorkspaceGitRoot *bool
	// Id label(s) of the format name=value. These will be set on the container and used to query for an existing container. If no --id-label is given, one will be inferred from the --workspace-folder path.
	IdLabels []string
	// devcontainer.json path. The default is to use .devcontainer/devcontainer.json or, if that does not exist, .devcontainer.json in the workspace folder.
	Config string
	// devcontainer.json path to override any devcontainer.json in the workspace folder (or built-in configuration). This is required when there is no devcontainer.json otherwise.
	OverrideConfig string
	// Log level for the --terminal-log-file. When set to trace, the log level for --log-file will also be set to trace.
	LogLevel DevContainerCliLogLevel
	// Log format.
	LogFormat DevContainerCliLogFormat
	// Number of columns to render the output for. This is required for some of the subprocesses to correctly render their output.
	TerminalColumns *int
	// Number of rows to render the output for. This is required for some of the subprocesses to correctly render their output.
	TerminalRows *int
	// Default value for the devcontainer.json's "userEnvProbe".
	DefaultUserEnvProbe DevContainerCliDefaultUserEnvProbe
	// Default for updating the remote user's UID and GID to the local user's one.
	UpdateRemoteUserUidDefault DevContainerCliUpdateRemoteUserUidDefault
	// Removes the dev container if it already exists.
	RemoveExistingContainer *bool
	// Builds the image with `--no-cache` if the container does not exist.
	BuildNoCache *bool
	// Fail if the container does not exist.
	ExpectExistingContainer *bool
	// Do not run onCreateCommand, updateContentCommand, postCreateCommand, postStartCommand or postAttachCommand and do not install dotfiles.
	SkipPostCreate *bool
	// Stop running user commands after running the command configured with waitFor or the updateContentCommand by default.
	SkipNonBlockingCommands *bool
	// Stop after onCreateCommand and updateContentCommand, rerunning updateContentCommand if it has run before.
	Prebuild *bool
	// Host path to a directory that is intended to be persisted and share state between sessions.
	UserDataFolder string
	// Additional mount point(s). Format: type=<bind|volume>,source=<source>,target=<target>[,external=<true|false>]
	Mounts []string
	// Remote environment variables of the format name=value. These will be added when executing the user commands.
	RemoteEnvs []string
	// Additional image to use as potential layer cache during image building
	CachesFrom []string
	// Additional image to use as potential layer cache during image building
	CacheTo string
	// Control whether BuildKit should be used
	Buildkit DevContainerCliBuildkit
	// Additional features to apply to the dev container (JSON as per "features" section in devcontainer.json)
	AdditionalFeatures string
	// Do not run postAttachCommand.
	SkipPostAttach *bool
	// URL of a dotfiles Git repository (e.g., https://github.com/owner/repository.git)
	DotfilesRepository string
	// The command to run after cloning the dotfiles repository. Defaults to run the first file of `install.sh`, `install`, `bootstrap.sh`, `bootstrap`, `setup.sh` and `setup` found in the dotfiles repository`s root folder.
	DotfilesInstallCommand string
	// The path to clone the dotfiles repository to. Defaults to `~/dotfiles`.
	DotfilesTargetPath string
	// Folder to cache CLI data, for example userEnvProbe results
	ContainerSessionDataFolder string
	// Omit remoteEnv from devcontainer.json for container metadata label
	OmitConfigRemoteEnvFromMetadata *bool
	// Path to a json file containing secret environment variables as key-value pairs.
	SecretsFile string
	// Include configuration in result.
	IncludeConfiguration *bool
	// Include merged configuration in result.
	IncludeMergedConfiguration *bool
}

func (tool *DevContainerCliTool) Run(settings *DevContainerCliUpSettings) error {
	args := []string{
		"up",
	}

	args = addString(args, settings.DockerPath, addSettings{prependElements: []string{"--docker-path"}})
	args = addString(args, settings.DockerComposePath, addSettings{prependElements: []string{"--docker-compose-path"}})
	args = addString(args, settings.ContainerDataFolder, addSettings{prependElements: []string{"--container-data-folder"}})
	args = addString(args, settings.ContainerSystemDataFolder, addSettings{prependElements: []string{"--container-system-data-folder"}})
	args = addString(args, settings.WorkspaceFolder, addSettings{prependElements: []string{"--workspace-folder"}})
	args = addString(args, string(settings.WorkspaceMountConsistency), addSettings{prependElements: []string{"--workspace-mount-consistency"}})
	args = addString(args, string(settings.GpuAvailability), addSettings{prependElements: []string{"--gpu-availability"}})
	args = addBoolean(args, settings.MountWorkspaceGitRoot, addSettings{prependElements: []string{"--mount-workspace-git-root"}})
	args = addStringList(args, settings.IdLabels, addSettings{prependElements: []string{"--id-label"}, handleEachListItemSeparately: true})
	args = addString(args, settings.Config, addSettings{prependElements: []string{"--config"}})
	args = addString(args, settings.OverrideConfig, addSettings{prependElements: []string{"--override-config"}})
	args = addString(args, string(settings.LogLevel), addSettings{prependElements: []string{"--log-level"}})
	args = addString(args, string(settings.LogFormat), addSettings{prependElements: []string{"--log-format"}})
	args = addInt(args, settings.TerminalColumns, addSettings{prependElements: []string{"--terminal-columns"}})
	args = addInt(args, settings.TerminalRows, addSettings{prependElements: []string{"--terminal-rows"}})
	args = addString(args, string(settings.DefaultUserEnvProbe), addSettings{prependElements: []string{"--default-user-env-probe"}})
	args = addString(args, string(settings.UpdateRemoteUserUidDefault), addSettings{prependElements: []string{"--update-remote-user-uid-default"}})
	args = addBoolean(args, settings.RemoveExistingContainer, addSettings{prependElements: []string{"--remove-existing-container"}})
	args = addBoolean(args, settings.BuildNoCache, addSettings{prependElements: []string{"--build-no-cache"}})
	args = addBoolean(args, settings.ExpectExistingContainer, addSettings{prependElements: []string{"--expect-existing-container"}})
	args = addBoolean(args, settings.SkipPostCreate, addSettings{prependElements: []string{"--skip-post-create"}})
	args = addBoolean(args, settings.SkipNonBlockingCommands, addSettings{prependElements: []string{"--skip-non-blocking-commands"}})
	args = addBoolean(args, settings.Prebuild, addSettings{prependElements: []string{"--prebuild"}})
	args = addString(args, settings.UserDataFolder, addSettings{prependElements: []string{"--user-data-folder"}})
	args = addStringList(args, settings.Mounts, addSettings{prependElements: []string{"--mount"}, handleEachListItemSeparately: true})
	args = addStringList(args, settings.RemoteEnvs, addSettings{prependElements: []string{"--remote-env"}, handleEachListItemSeparately: true})
	args = addStringList(args, settings.CachesFrom, addSettings{prependElements: []string{"--cache-from"}, handleEachListItemSeparately: true})
	args = addString(args, settings.CacheTo, addSettings{prependElements: []string{"--cache-to"}})
	args = addString(args, string(settings.Buildkit), addSettings{prependElements: []string{"--buildkit"}})
	args = addString(args, settings.AdditionalFeatures, addSettings{prependElements: []string{"--additional-features"}})
	args = addBoolean(args, settings.SkipPostAttach, addSettings{prependElements: []string{"--skip-post-attach"}})
	args = addString(args, settings.DotfilesRepository, addSettings{prependElements: []string{"--dotfiles-repository"}})
	args = addString(args, settings.DotfilesInstallCommand, addSettings{prependElements: []string{"--dotfiles-install-command"}})
	args = addString(args, settings.DotfilesTargetPath, addSettings{prependElements: []string{"--dotfiles-target-path"}})
	args = addString(args, settings.ContainerSessionDataFolder, addSettings{prependElements: []string{"--container-session-data-folder"}})
	args = addBoolean(args, settings.OmitConfigRemoteEnvFromMetadata, addSettings{prependElements: []string{"--omit-config-remote-env-from-metadata"}})
	args = addString(args, settings.SecretsFile, addSettings{prependElements: []string{"--secrets-file"}})
	args = addBoolean(args, settings.IncludeConfiguration, addSettings{prependElements: []string{"--include-configuration"}})
	args = addBoolean(args, settings.IncludeMergedConfiguration, addSettings{prependElements: []string{"--include-merged-configuration"}})

	args = append(args, settings.CustomArguments...)
	cmd := exec.Command("devcontainer", goext.RemoveEmpty(args)...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommandO(cmd, execr.WithConsoleOutput(settings.OutputToConsole))
}

////////// Exec

type DevContainerCliExecSettings struct {
	ToolSettingsBase
	// Host path to a directory that is intended to be persisted and share state between sessions.
	UserDataFolder string
	// Docker CLI path.
	DockerPath string
	// Docker Compose CLI path.
	DockerComposePath string
	// Container data folder where user data inside the container will be stored.
	ContainerDataFolder string
	// Container system data folder where system data inside the container will be stored.
	ContainerSystemDataFolder string
	// Workspace folder path. The devcontainer.json will be looked up relative to this path.
	WorkspaceFolder string
	// Mount the workspace using its Git root.
	MountWorkspaceGitRoot *bool
	// Id of the container to run the user commands for.
	ContainerId string
	// Id label(s) of the format name=value. If no --container-id is given the id labels will be used to look up the container. If no --id-label is given, one will be inferred from the --workspace-folder path.
	IdLabels []string
	// devcontainer.json path. The default is to use .devcontainer/devcontainer.json or, if that does not exist, .devcontainer.json in the workspace folder.
	Config string
	// devcontainer.json path to override any devcontainer.json in the workspace folder (or built-in configuration). This is required when there is no devcontainer.json otherwise.
	OverrideConfig string
	// Log level for the --terminal-log-file. When set to trace, the log level for --log-file will also be set to trace.
	LogLevel DevContainerCliLogLevel
	// Log format.
	LogFormat DevContainerCliLogFormat
	// Number of columns to render the output for. This is required for some of the subprocesses to correctly render their output.
	TerminalColumns *int
	// Number of rows to render the output for. This is required for some of the subprocesses to correctly render their output.
	TerminalRows *int
	// Default value for the devcontainer.json's "userEnvProbe".
	DefaultUserEnvProbe DevContainerCliDefaultUserEnvProbe
	// Remote environment variables of the format name=value. These will be added when executing the user commands.
	RemoteEnvs []string
}

func (tool *DevContainerCliTool) Exec(command string, settings *DevContainerCliExecSettings) error {
	args := []string{
		"exec",
		command,
	}

	args = addString(args, settings.UserDataFolder, addSettings{prependElements: []string{"user-data-folder"}})
	args = addString(args, settings.DockerPath, addSettings{prependElements: []string{"docker-path"}})
	args = addString(args, settings.DockerComposePath, addSettings{prependElements: []string{"docker-compose-path"}})
	args = addString(args, settings.ContainerDataFolder, addSettings{prependElements: []string{"container-data-folder"}})
	args = addString(args, settings.ContainerSystemDataFolder, addSettings{prependElements: []string{"container-system-data-folder"}})
	args = addString(args, settings.WorkspaceFolder, addSettings{prependElements: []string{"workspace-folder"}})
	args = addBoolean(args, settings.MountWorkspaceGitRoot, addSettings{prependElements: []string{"mount-workspace-git-root"}})
	args = addString(args, settings.ContainerId, addSettings{prependElements: []string{"container-id"}})
	args = addStringList(args, settings.IdLabels, addSettings{prependElements: []string{"id-label"}, handleEachListItemSeparately: true})
	args = addString(args, settings.Config, addSettings{prependElements: []string{"config"}})
	args = addString(args, settings.OverrideConfig, addSettings{prependElements: []string{"override-config"}})
	args = addString(args, string(settings.LogLevel), addSettings{prependElements: []string{"log-level"}})
	args = addString(args, string(settings.LogFormat), addSettings{prependElements: []string{"log-format"}})
	args = addInt(args, settings.TerminalColumns, addSettings{prependElements: []string{"terminal-columns"}})
	args = addInt(args, settings.TerminalRows, addSettings{prependElements: []string{"terminal-rows"}})
	args = addString(args, string(settings.DefaultUserEnvProbe), addSettings{prependElements: []string{"default-user-env-probe"}})
	args = addStringList(args, settings.RemoteEnvs, addSettings{prependElements: []string{"remote-env"}, handleEachListItemSeparately: true})

	args = append(args, settings.CustomArguments...)
	cmd := exec.Command("devcontainer", goext.RemoveEmpty(args)...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommandO(cmd, execr.WithConsoleOutput(settings.OutputToConsole))
}

////////// Features Test

type DevContainerCliFeaturesTestSettings struct {
	ToolSettingsBase
	// Path to folder containing 'src' and 'test' sub-folders.
	ProjectFolder string
	// Feature(s) to test .
	Features []string
	// Filter current tests to only run scenarios containing this string.
	Filter string
	//  Run only scenario tests under 'tests/_global'
	GlobalScenariosOnly *bool
	// Skip all 'scenario' style tests.
	SkipScenarios *bool
	// Skip all 'autogenerated' style tests (test.sh)
	SkipAutogenerated *bool
	// Skip all 'duplicate' style tests (duplicate.sh).
	SkipDuplicated *bool
	// Allow an element of randomness in test cases.
	PermitRandomization *bool
	// Base Image. Not used for scenarios.
	BaseImage string
	// Remote user.  Not used for scenarios.
	RemoteUser string
	// Log level.
	LogLevel DevContainerCliLogLevel
	// Do not remove test containers after running tests.
	PreserveTestContainers *bool
	// Quiets output.
	Quiet *bool
}

func (tool *DevContainerCliTool) FeaturesTest(settings *DevContainerCliFeaturesTestSettings) error {
	args := []string{
		"features",
		"test",
	}

	args = addString(args, settings.ProjectFolder, addSettings{prependElements: []string{"--project-folder"}})
	args = addStringList(args, settings.Features, addSettings{prependElements: []string{"--features"}, listSeparator: " "})
	args = addString(args, settings.Filter, addSettings{prependElements: []string{"--filter"}})
	args = addBoolean(args, settings.GlobalScenariosOnly, addSettings{prependElements: []string{"--global-scenarios-only"}})
	args = addBoolean(args, settings.SkipScenarios, addSettings{prependElements: []string{"--skip-scenarios"}})
	args = addBoolean(args, settings.SkipAutogenerated, addSettings{prependElements: []string{"--skip-autogenerated"}})
	args = addBoolean(args, settings.SkipDuplicated, addSettings{prependElements: []string{"--skip-duplicated"}})
	args = addBoolean(args, settings.PermitRandomization, addSettings{prependElements: []string{"--permit-randomization"}})
	args = addString(args, settings.BaseImage, addSettings{prependElements: []string{"--base-image"}})
	args = addString(args, settings.RemoteUser, addSettings{prependElements: []string{"--remote-user"}})
	args = addString(args, string(settings.LogLevel), addSettings{prependElements: []string{"--log-level"}})
	args = addBoolean(args, settings.PreserveTestContainers, addSettings{prependElements: []string{"--preserve-test-containers"}})
	args = addBoolean(args, settings.Quiet, addSettings{prependElements: []string{"--quiet"}})

	args = append(args, settings.CustomArguments...)
	cmd := exec.Command("devcontainer", goext.RemoveEmpty(args)...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommandO(cmd, execr.WithConsoleOutput(settings.OutputToConsole))
}

////////// Features Package

type DevContainerCliFeaturesPackageSettings struct {
	ToolSettingsBase
	// Path to a feature or folder of features to package.
	Target string
	// Path to output directory. Will create directories as needed.
	OutputFolder string
	// Automatically delete previous output directory before packaging.
	ForceCleanOutputFolder *bool
	// Log level.
	LogLevel DevContainerCliLogLevel
}

func (tool *DevContainerCliTool) FeaturesPackage(settings *DevContainerCliFeaturesPackageSettings) error {
	args := []string{
		"features",
		"package",
		settings.Target,
	}

	args = addString(args, settings.OutputFolder, addSettings{prependElements: []string{"--output-folder"}})
	args = addBoolean(args, settings.ForceCleanOutputFolder, addSettings{prependElements: []string{"--force-clean-output-folder"}})
	args = addString(args, string(settings.LogLevel), addSettings{prependElements: []string{"--log-level"}})

	args = append(args, settings.CustomArguments...)
	cmd := exec.Command("devcontainer", goext.RemoveEmpty(args)...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommandO(cmd, execr.WithConsoleOutput(settings.OutputToConsole))
}

////////// Features Publish

type DevContainerCliFeaturesPublishSettings struct {
	ToolSettingsBase
	// Path to a feature or folder of features to publish.
	Target string
	// Name of the OCI registry.
	Registry string
	// Unique identifier for the collection of features.
	Namespace string
	// Log level.
	LogLevel DevContainerCliLogLevel
}

func (tool *DevContainerCliTool) FeaturesPublish(settings *DevContainerCliFeaturesPublishSettings) error {
	args := []string{
		"features",
		"publish",
		settings.Target,
	}

	args = addString(args, settings.Registry, addSettings{prependElements: []string{"--registry"}})
	args = addString(args, settings.Namespace, addSettings{prependElements: []string{"--namespace"}})
	args = addString(args, string(settings.LogLevel), addSettings{prependElements: []string{"--log-level"}})

	args = append(args, settings.CustomArguments...)
	cmd := exec.Command("devcontainer", goext.RemoveEmpty(args)...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommandO(cmd, execr.WithConsoleOutput(settings.OutputToConsole))
}

////////// Types

type DevContainerCliLogLevel string

const (
	DEV_CONTAINER_CLI_LOG_LEVEL_DEFAULT DevContainerCliLogLevel = ""
	DEV_CONTAINER_CLI_LOG_LEVEL_INFO    DevContainerCliLogLevel = "info"
	DEV_CONTAINER_CLI_LOG_LEVEL_DEBUG   DevContainerCliLogLevel = "debug"
	DEV_CONTAINER_CLI_LOG_LEVEL_TRACE   DevContainerCliLogLevel = "trace"
)

type DevContainerCliLogFormat string

const (
	DEV_CONTAINER_CLI_LOG_FORMAT_DEFAULT DevContainerCliLogFormat = ""
	DEV_CONTAINER_CLI_LOG_FORMAT_TEXT    DevContainerCliLogFormat = "text"
	DEV_CONTAINER_CLI_LOG_FORMAT_JSON    DevContainerCliLogFormat = "json"
)

type DevContainerCliBuildkit string

const (
	DEV_CONTAINER_CLI_BUILDKIT_DEFAULT DevContainerCliBuildkit = ""
	DEV_CONTAINER_CLI_BUILDKIT_AUTO    DevContainerCliBuildkit = "auto"
	DEV_CONTAINER_CLI_BUILDKIT_NEVER   DevContainerCliBuildkit = "never"
)

type DevContainerCliWorkspaceMountConsistency string

const (
	DEV_CONTAINER_CLI_WORKSPACE_MOUNT_CONSISTENCY_DEFAULT    DevContainerCliWorkspaceMountConsistency = ""
	DEV_CONTAINER_CLI_WORKSPACE_MOUNT_CONSISTENCY_CONSISTENT DevContainerCliWorkspaceMountConsistency = "consistent"
	DEV_CONTAINER_CLI_WORKSPACE_MOUNT_CONSISTENCY_CACHED     DevContainerCliWorkspaceMountConsistency = "cached"
	DEV_CONTAINER_CLI_WORKSPACE_MOUNT_CONSISTENCY_DELEGATED  DevContainerCliWorkspaceMountConsistency = "delegated"
)

type DevContainerCliGpuAvailability string

const (
	DEV_CONTAINER_CLI_GPU_AVAILABILITY_DEFAULT DevContainerCliGpuAvailability = ""
	DEV_CONTAINER_CLI_GPU_AVAILABILITY_ALL     DevContainerCliGpuAvailability = "all"
	DEV_CONTAINER_CLI_GPU_AVAILABILITY_DETECT  DevContainerCliGpuAvailability = "detect"
	DEV_CONTAINER_CLI_GPU_AVAILABILITY_NONE    DevContainerCliGpuAvailability = "none"
)

type DevContainerCliDefaultUserEnvProbe string

const (
	DEV_CONTAINER_CLI_DEFAULT_USER_ENV_PROBE_DEFAULT               DevContainerCliDefaultUserEnvProbe = ""
	DEV_CONTAINER_CLI_DEFAULT_USER_ENV_PROBE_NONE                  DevContainerCliDefaultUserEnvProbe = "none"
	DEV_CONTAINER_CLI_DEFAULT_USER_ENV_PROBE_LOGININTERACTIVESHELL DevContainerCliDefaultUserEnvProbe = "loginInteractiveShell"
	DEV_CONTAINER_CLI_DEFAULT_USER_ENV_PROBE_INTERACTIVESHELL      DevContainerCliDefaultUserEnvProbe = "interactiveShell"
	DEV_CONTAINER_CLI_DEFAULT_USER_ENV_PROBE_LOGINSHELL            DevContainerCliDefaultUserEnvProbe = "loginShell"
)

type DevContainerCliUpdateRemoteUserUidDefault string

const (
	DEV_CONTAINER_CLI_UPDATE_REMOTE_USER_UID_DEFAULT_DEFAULT DevContainerCliUpdateRemoteUserUidDefault = ""
	DEV_CONTAINER_CLI_UPDATE_REMOTE_USER_UID_DEFAULT_NEVER   DevContainerCliUpdateRemoteUserUidDefault = "never"
	DEV_CONTAINER_CLI_UPDATE_REMOTE_USER_UID_DEFAULT_ON      DevContainerCliUpdateRemoteUserUidDefault = "on"
	DEV_CONTAINER_CLI_UPDATE_REMOTE_USER_UID_DEFAULT_OFF     DevContainerCliUpdateRemoteUserUidDefault = "off"
)
