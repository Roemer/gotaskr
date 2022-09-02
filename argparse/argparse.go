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

// ParseArgString parses the arguments from a given array.
func ParseArgString(args []string) map[string]string {
	var argsMap map[string]string = make(map[string]string)

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
				argsMap[key] = valuePart
				lastKey = ""
			} else {
				// Key only
				key := arg[2:]
				lastKey = key
				nextCanBeValue = true
				argsMap[key] = ""
			}
			continue
		}

		// Short options
		if strings.HasPrefix(arg, "-") {
			optionChars := arg[1:]
			if len(optionChars) > 1 {
				// Multiple flags, so they do not have a value
				for _, c := range optionChars {
					argsMap[string(c)] = ""
				}
				lastKey = ""
			} else {
				// Only one flag, so the next might be a value
				key := optionChars
				lastKey = key
				nextCanBeValue = true
				argsMap[key] = ""
			}
			continue
		}

		// The current value is no key so it seems to be a value to a previous key
		if nextCanBeValue && len(lastKey) > 0 {
			nextCanBeValue = false
			argsMap[lastKey] = arg
			continue
		}
	}

	// Make sure to add the last key when no value was provided
	if nextCanBeValue && len(lastKey) > 0 {
		argsMap[lastKey] = ""
	}

	return argsMap
}
