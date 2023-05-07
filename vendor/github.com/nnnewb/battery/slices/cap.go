package slices

// Cap returns the capacity of the slice s.
func (s Slice[T]) Cap() int {
	return cap(s)
}
