package slices

import (
	"sort"
)

// SortLessFunc sorts the slice s using the provided less function and returns the sorted slice.
func (s Slice[T]) SortLessFunc(less func(i, j int) bool) Slice[T] {
	sort.Slice(s, less)
	return s
}
