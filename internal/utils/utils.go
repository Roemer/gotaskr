package utils

import (
	"fmt"
	"sort"
	"strings"
)

// PanicOnError panics if the given error is set.
func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// Noop is a no operation function that can be used for tasks.
// For example if they are just used for multiple dependencies.
func Noop() error {
	return nil
}

// Pass is a no-op function that can be used to set variables to used.
// Useful during development but must be removed afterwards!
func Pass(i ...any) {
	// No-Op
}

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

// ProcessMapSorted calls the given function on all entries of the map, sorted by their key value (by string).
func ProcessMapSorted[TK comparable, TV any](m map[TK]TV, pf func(TK, TV) error) error {
	// Get all keys
	keys := GetMapKeys(m)

	// Call the method for each function
	for _, k := range keys {
		if err := pf(k, m[k]); err != nil {
			return err
		}
	}
	return nil
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
