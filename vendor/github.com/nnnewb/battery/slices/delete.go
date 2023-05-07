package slices

// Delete removes the element at the specified index from the slice and returns the resulting slice.
func (s Slice[T]) Delete(idx int) Slice[T] {
	if idx < 0 || idx >= len(s) {
		// index out of bounds
		return s
	}
	return append(s[:idx], s[idx+1:]...)
}
