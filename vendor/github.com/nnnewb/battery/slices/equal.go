package slices

// Equal compares two slices for equality. Pass in a comparison function eq, which returns true if the elements are the same.
// The slices are considered equal if they have the same number of elements and the same order. Note that a nil slice and an empty slice are considered equal.
func Equal[T comparable](s, other []T) bool {
	if len(s) != len(other) {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] != other[i] {
			return false
		}
	}
	return true
}

// EqualFunc compares two slices for equality. Pass in a comparison function eq, which returns true if the elements are the same.
// The slices are considered equal if they have the same number of elements and the same order. Note that a nil slice and an empty slice are considered equal.
func EqualFunc[T1, T2 any](s []T1, other []T2, eq func(T1, T2) bool) bool {
	if len(s) != len(other) {
		return false
	}
	for i, v1 := range s {
		v2 := other[i]
		if !eq(v1, v2) {
			return false
		}
	}
	return true
}
