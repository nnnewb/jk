package slices

// CountBy returns a map of keys to the number of elements in the slice that map to that key.
func CountBy[T any, K comparable](s []T, keyFunc func(T) K) map[K]int {
	counts := make(map[K]int)
	for _, elem := range s {
		key := keyFunc(elem)
		counts[key]++
	}
	return counts
}
