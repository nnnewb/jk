package slices

// GroupBy groups the elements of a slice by a key returned by the keyFunc function.
// It returns a map where the keys are the unique values returned by keyFunc and the values are slices of elements that share the same key.
func GroupBy[T any, K comparable](s []T, keyFunc func(T) K) map[K][]T {
	groups := make(map[K][]T)
	for _, elem := range s {
		key := keyFunc(elem)
		groups[key] = append(groups[key], elem)
	}
	return groups
}
