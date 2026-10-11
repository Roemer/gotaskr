// Package argparse is a simple parser for command line arguments.
package argparse

import (
	"os"
	"strings"
)

// ParseArgs parses the arguments from the os.Args (ignoring the first one).
func ParseArgs() map[string]string {
	return ParseArgString(os.Args[1:])
}

// ParseArgsMulti parses the arguments from the os.Args (ignoring the first one)
// and returns all values for each key, preserving repeated flags.
func ParseArgsMulti() map[string][]string {
	return ParseArgStringMulti(os.Args[1:])
}

// ParseArgString parses the arguments from a given array.
func ParseArgString(args []string) map[string]string {
	multi := ParseArgStringMulti(args)
	result := make(map[string]string, len(multi))
	for key, values := range multi {
		result[key] = values[len(values)-1]
	}
	return result
}

// ParseArgStringMulti parses the arguments from a given array and returns all
// values for each key, preserving repeated flags (e.g. --target A --target B).
func ParseArgStringMulti(args []string) map[string][]string {
	argsMap := make(map[string][]string)

	var lastKey = ""
	var nextCanBeValue = false
	for _, arg := range args {
		// Long Options
		if strings.HasPrefix(arg, "--") {
			if strings.Contains(arg, "=") {
				// Key with value
				parts := strings.SplitN(arg, "=", 2)
				keyPart := parts[0]
				valuePart := parts[1]
				key := keyPart[2:]
				argsMap[key] = append(argsMap[key], valuePart)
				lastKey = ""
				nextCanBeValue = false
			} else {
				// Key only
				key := arg[2:]
				lastKey = key
				nextCanBeValue = true
				argsMap[key] = append(argsMap[key], "")
			}
			continue
		}

		// Short options
		if strings.HasPrefix(arg, "-") {
			optionChars := arg[1:]
			if len(optionChars) > 1 {
				// Multiple flags, so they do not have a value
				for _, c := range optionChars {
					argsMap[string(c)] = append(argsMap[string(c)], "")
				}
				lastKey = ""
				nextCanBeValue = false
			} else {
				// Only one flag, so the next might be a value
				key := optionChars
				lastKey = key
				nextCanBeValue = true
				argsMap[key] = append(argsMap[key], "")
			}
			continue
		}

		// The current value is no key so it seems to be a value to a previous key
		if nextCanBeValue && len(lastKey) > 0 {
			nextCanBeValue = false
			values := argsMap[lastKey]
			values[len(values)-1] = arg
			continue
		}
	}

	return argsMap
}
