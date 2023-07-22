package slices

// Last returns the index of the element's last appearance in the slice.
// If no element satisfies the predicate, the second return value is false.
func Last[T comparable](s []T, v T) (int, bool) {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == v {
			return i, true
		}
	}
	return 0, false
}

// LastFunc returns the last element in the slice that satisfies the given predicate function.
// If no element satisfies the predicate, it returns the zero value of the slice's element type and false.
func LastFunc[T any](s []T, predicate func(T) bool) (int, bool) {
	for i := len(s) - 1; i >= 0; i-- {
		if predicate(s[i]) {
			return i, true
		}
	}
	return 0, false
}
