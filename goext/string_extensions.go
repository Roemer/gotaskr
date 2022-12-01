package goext

import "strings"

func TrimAllPrefix(s, prefix string) string {
	for strings.HasPrefix(s, prefix) {
		s = s[len(prefix):]
	}
	return s
}

func TrimAllSuffix(s, suffix string) string {
	for strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func TrimNewlineSuffix(v string) string {
	return TrimAllSuffix(TrimAllSuffix(v, "\r\n"), "\n")
}

// SplitByNewLine splits the given value by newlines and returns a slice.
func SplitByNewLine(value string) []string {
	return strings.Split(strings.ReplaceAll(value, "\r\n", "\n"), "\n")
}
