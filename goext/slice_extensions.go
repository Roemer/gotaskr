package goext

// AddIf adds the given values to a slice if the condition is fulfilled.
func AddIf[T any](slice []T, cond bool, values ...T) []T {
	if cond {
		return append(slice, values...)
	}
	return slice
}

// AppendIfMissing appends the given element if it is missing in the slice.
func AppendIfMissing[T comparable](slice []T, value T) []T {
	for _, ele := range slice {
		if ele == value {
			return slice
		}
	}
	return append(slice, value)
}

// RemoveEmpty removes all empty strings from a slice.
func RemoveEmpty(slice []string) []string {
	var r []string
	for _, str := range slice {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
