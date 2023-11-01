package goext

// AppendIf adds the given values to a slice if the condition is fulfilled.
func AppendIf[T any](slice []T, cond bool, values ...T) []T {
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

// ChunkBy creates chunks of the given size from the given array and returns them.
func ChunkBy[T any](items []T, chunkSize int) (chunks [][]T) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}

// Prepend prepends the given elements to the given array
func Prepend[T any](slice []T, elems ...T) []T {
	return append(elems, slice...)
}
