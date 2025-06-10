package execr

type RunOptions struct {
	arguments             []string
	outputToConsole       bool
	workingDirectory      string
	skipPostProcessOutput bool
}

func WithArguments(arguments ...string) func(*RunOptions) {
	return func(options *RunOptions) {
		options.arguments = append(options.arguments, arguments...)
	}
}

func WithArgumentsSplitted(arguments string) func(*RunOptions) {
	return func(options *RunOptions) {
		options.arguments = append(options.arguments, SplitArguments(arguments)...)
	}
}

func WithWorkingDirectory(workingDirectory string) func(*RunOptions) {
	return func(options *RunOptions) {
		options.workingDirectory = workingDirectory
	}
}

func WithConsoleOutput(outputToConsole bool) func(*RunOptions) {
	return func(options *RunOptions) {
		options.outputToConsole = outputToConsole
	}
}

func WithSkipPostProcessOutput(skipPostProcessOutput bool) func(*RunOptions) {
	return func(options *RunOptions) {
		options.skipPostProcessOutput = skipPostProcessOutput
	}
}
