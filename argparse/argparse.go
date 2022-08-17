package argparse

import (
	"os"
	"strings"
)

func ParseArgs() map[string]string {
	return ParseArgString(os.Args[1:])
}

func ParseArgString(args []string) map[string]string {
	var argsMap map[string]string = make(map[string]string)

	var lastKey = ""
	var nextIsValue = false
	var nextCanBeValue = false
	for _, arg := range args {
		// The current value must be a value to a previous key
		if nextIsValue {
			nextIsValue = false
			argsMap[lastKey] = arg
			continue
		}

		// Long Options
		if strings.HasPrefix(arg, "--") {
			if strings.Contains(arg, "=") {
				parts := strings.SplitN(arg, "=", 2)
				keyPart := parts[0]
				valuePart := parts[1]
				key := keyPart[2:]
				argsMap[key] = valuePart
			} else {
				lastKey = arg[2:]
				nextIsValue = true
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
			} else {
				// Only one flag, so the next might be a value
				lastKey = optionChars
				nextCanBeValue = true
			}
			continue
		}

		// The current value is no key so it seems to be a value to a previous key
		if nextCanBeValue {
			nextCanBeValue = false
			argsMap[lastKey] = arg
			continue
		}
	}

	return argsMap
}
