package maps

// EqualFunc is like Equal, but compares values using eq. Keys are still compared with ==.
func EqualFunc[M1 ~map[K]V1, M2 ~map[K]V2, K comparable, V1, V2 any](m1 M1, m2 M2, eq func(V1, V2) bool) bool {
	for k1, v1 := range m1 {
		if v2, ok := m2[k1]; ok {
			if !eq(v1, v2) {
				return false
			}
		} else {
			return false
		}
	}
	return true
}
