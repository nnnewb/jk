package slices

// Map applies the function f to each element of the slice s and returns a new slice with the results in the same order.
// The type of the elements in the resulting slice is determined by the return type of the function f.
func Map[T, R any](s []T, f func(T) R) Slice[R] {
	result := make([]R, len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}
