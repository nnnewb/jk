package maps

// Equal reports whether two maps contain the same key/value pairs. Values are compared using ==.
func Equal[M1, M2 ~map[K]V, K, V comparable](m1 M1, m2 M2) bool {
	for k, v1 := range m1 {
		if v2, ok := m2[k]; !ok || v1 != v2 {
			return false
		}
	}
	for k := range m2 {
		if _, ok := m1[k]; !ok {
			return false
		}
	}
	return true
}
