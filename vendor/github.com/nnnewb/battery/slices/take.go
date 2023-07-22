package slices

// Take returns the first n elements of the slice s.
// If n is greater than the length of s, Take returns the entire slice.
func Take[T any](s []T, n int) []T {
	if n > len(s) {
		n = len(s)
	}
	return s[:n]
}
