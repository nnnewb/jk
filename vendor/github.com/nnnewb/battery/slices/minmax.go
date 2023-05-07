package slices

// Min returns the minimum element in the slice based on the provided less function.
// If the slice is empty, it returns false as the second value.
func (s Slice[T]) Min(less func(i, j T) bool) (T, bool) {
	if len(s) == 0 {
		var zero T
		return zero, false
	}
	min := s[0]
	for _, t := range s {
		if less(t, min) {
			min = t
		}
	}
	return min, true
}

// Max returns the maximum element in the slice based on the provided less function.
// If the slice is empty, it returns false as the second value.
func (s Slice[T]) Max(less func(i, j T) bool) (T, bool) {
	if len(s) == 0 {
		var zero T
		return zero, false
	}
	max := s[0]
	for _, t := range s {
		if less(max, t) {
			max = t
		}
	}
	return max, true
}
