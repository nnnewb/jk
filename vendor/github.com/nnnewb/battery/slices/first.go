package slices

// First returns the index of first element in the slice that satisfies the
// given predicate function.
//
// If no element satisfies the predicate, the second return value is false.
func First[T any](s []T, predicate func(T) bool) (int, bool) {
	for idx, elem := range s {
		if predicate(elem) {
			return idx, true
		}
	}
	return 0, false
}
