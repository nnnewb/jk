package slices

// Contains returns true if the given value v is present in the slice s, using the provided equal function to compare values.
// If the value is not found, it returns false.
func (s Slice[T]) Contains(v T, equal func(T, T) bool) bool {
	for _, val := range s {
		if equal(val, v) {
			return true
		}
	}
	return false
}
