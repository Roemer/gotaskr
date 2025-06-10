package execr

// This type holds the options for running an executable or command.
type RunOptions struct {
	arguments             []string
	outputToConsole       bool
	workingDirectory      string
	skipPostProcessOutput bool
}

// Appends the given arguments.
func WithArguments(arguments ...string) func(*RunOptions) {
	return func(options *RunOptions) {
		options.arguments = append(options.arguments, arguments...)
	}
}

// Splits the given argument string into separate arguments and appends them.
func WithArgumentsSplitted(arguments string) func(*RunOptions) {
	return func(options *RunOptions) {
		options.arguments = append(options.arguments, SplitArguments(arguments)...)
	}
}

// Sets the working directory for the command.
func WithWorkingDirectory(workingDirectory string) func(*RunOptions) {
	return func(options *RunOptions) {
		options.workingDirectory = workingDirectory
	}
}

// Allows enabling or disabling console output.
func WithConsoleOutput(outputToConsole bool) func(*RunOptions) {
	return func(options *RunOptions) {
		options.outputToConsole = outputToConsole
	}
}

// Allows skipping the post-processing of output.
func WithSkipPostProcessOutput(skipPostProcessOutput bool) func(*RunOptions) {
	return func(options *RunOptions) {
		options.skipPostProcessOutput = skipPostProcessOutput
	}
}
