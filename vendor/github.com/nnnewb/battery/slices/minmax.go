package slices

import "github.com/nnnewb/battery/internal/constraints"

// Min returns the minimum element in the slice.
// If the slice is empty, it returns false as the second value.
func Min[T constraints.Ordered](s []T) (int, T, bool) {
	if len(s) == 0 {
		var zero T
		return 0, zero, false
	}

	minIdx := 0
	min := s[0]
	for idx, t := range s {
		if t < min {
			minIdx = idx
			min = t
		}
	}
	return minIdx, min, true
}

// Max returns the maximum element in the slice based on the provided less function.
// If the slice is empty, it returns false as the second value.
func Max[T constraints.Ordered](s []T) (int, T, bool) {
	if len(s) == 0 {
		var zero T
		return 0, zero, false
	}

	maxIdx := 0
	max := s[0]
	for i, t := range s {
		if max < t {
			maxIdx = i
			max = t
		}
	}
	return maxIdx, max, true
}

// MinFunc returns the minimum element in the slice based on the provided less function.
// If the slice is empty, it returns false as the second value.
func MinFunc[T any](s []T, less func(i, j T) bool) (T, bool) {
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

// MaxFunc returns the maximum element in the slice based on the provided less function.
// If the slice is empty, it returns false as the second value.
func MaxFunc[T any](s []T, less func(i, j T) bool) (T, bool) {
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
