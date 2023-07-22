//go:build go1.18

package slices

// All tests whether all elements in the slice satisfy the predicate function.
// If the length of s is 0, returns false.
func All[T any](s []T, predicate func(T) bool) bool {
	if len(s) == 0 {
		return false
	}
	for _, elem := range s {
		if !predicate(elem) {
			return false
		}
	}
	return true
}
