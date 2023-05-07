package slices

// First returns the first element in the slice that satisfies the given predicate function.
// If no element satisfies the predicate, the second return value is false.
func (s Slice[T]) First(predicate func(T) bool) (T, bool) {
	var zeroValT T
	for _, elem := range s {
		if predicate(elem) {
			return elem, true
		}
	}
	return zeroValT, false
}
