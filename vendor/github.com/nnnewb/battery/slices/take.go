package slices

// Take returns the first n elements of the slice s.
// If n is greater than the length of s, Take returns the entire slice.
func (s Slice[T]) Take(n int) Slice[T] {
	if n > len(s) {
		n = len(s)
	}
	return s[:n]
}
