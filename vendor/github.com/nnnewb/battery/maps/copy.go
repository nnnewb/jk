package maps

// Copy copies all key/value pairs in src adding them to dst.
// When a key in src is already present in dst, the value in
// dst will be overwritten by the value associated with the key in src.
func Copy[M1 ~map[K]V, M2 ~map[K]V, K comparable, V any](dst M1, src M2) {
	if src == nil || dst == nil {
		return
	}

	for k, v := range src {
		dst[k] = v
	}
}
