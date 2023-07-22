package maps

// Clone returns a copy of m. This is a shallow clone: the new keys and values are set using ordinary assignment.
func Clone[M ~map[K]V, K comparable, V any](m M) M {
	// create a new map of the same type as the input map
	newMap := make(M, len(m))
	// iterate over the input map and copy each key-value pair to the new map using ordinary assignment
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}
