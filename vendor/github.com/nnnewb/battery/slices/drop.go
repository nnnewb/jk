package slices

// Drop returns a new slice with the first n elements removed. If n is greater than or equal to the length of the slice,
// an empty slice of the same type is returned. The original slice is not modified.
func Drop[T any](s []T, n int) []T {
	// Check if n is greater than or equal to the length of the slice
	if n >= len(s) {
		// If n is greater than or equal to the length of the slice, return an empty slice of the same type
		return make([]T, 0)
	}

	// Return a new slice starting from index n to the end of the original slice
	return s[n:]
}
