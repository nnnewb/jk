package slices

// Last returns the last element in the slice that satisfies the given predicate function.
// If no element satisfies the predicate, it returns the zero value of the slice's element type and false.
func (s Slice[T]) Last(predicate func(T) bool) (T, bool) {
	var zero T
	for i := len(s) - 1; i >= 0; i-- {
		if predicate(s[i]) {
			return s[i], true
		}
	}
	return zero, false
}
