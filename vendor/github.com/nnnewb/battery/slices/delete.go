package slices

// Delete removes the element at the specified index from the slice and returns the resulting slice.
func Delete[T any](s []T, idx int) []T {
	if idx < 0 || idx >= len(s) {
		// index out of bounds
		return s
	}
	return append(s[:idx], s[idx+1:]...)
}
