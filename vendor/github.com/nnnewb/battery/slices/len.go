package slices

// Len returns the length of the slice.
func (s Slice[T]) Len() int {
	return len(s)
}
