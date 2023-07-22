package slices

import (
	"sort"

	"github.com/nnnewb/battery/internal/constraints"
)

// Sort sorts the slice s and returns the sorted slice.
func Sort[T constraints.Ordered](s []T) []T {
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})
	return s
}

// SortLessFunc sorts the slice s using the provided less function and returns the sorted slice.
func SortLessFunc[T any](s []T, less func(i, j int) bool) []T {
	sort.Slice(s, less)
	return s
}
