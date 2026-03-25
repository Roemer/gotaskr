package utils

import (
	"fmt"
	"sort"
	"strings"
)

// Joins any slice with the given separator.
func StringsJoinAny[T any](elems []T, sep string) string {
	values := []string{}
	for _, element := range elems {
		values = append(values, fmt.Sprintf("%v", element))
	}
	return strings.Join(values, sep)
}

// GetMapKeys gets all the keys of the given map in alpahnummeric sorted order.
func GetMapKeys[TK comparable, TV any](m map[TK]TV) []TK {
	// Get all keys
	keys := make([]TK, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	// Sort the keys alphanumeric
	sort.Slice(keys, func(i, j int) bool {
		return fmt.Sprint(keys[i]) < fmt.Sprint(keys[j])
	})

	return keys
}

// ConvertMapToSingleString converts the given map into a joined string with the given separators
func ConvertMapToSingleString(m map[string]string, kvpSeparator string, entrySeparator string) string {
	var sb strings.Builder
	mapKeys := GetMapKeys(m)
	isFirst := true
	for _, key := range mapKeys {
		if !isFirst {
			sb.WriteString(entrySeparator)
		}
		sb.WriteString(key)
		sb.WriteString(kvpSeparator)
		sb.WriteString(m[key])
		isFirst = false
	}
	return sb.String()
}
