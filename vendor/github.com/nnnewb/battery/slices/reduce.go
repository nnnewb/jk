package slices

// Reduce applies a function to each element of a slice and returns a single value.
// The function takes two arguments: an accumulator and the current element.
// The accumulator is initialized to the initial value.
// The function is applied to the accumulator and the first element of the slice,
// then to the result and the second element, and so on.
// The final result is returned.
func Reduce[T, R any](s []T, initial R, f func(R, T) R) R {
	result := initial
	for _, v := range s {
		result = f(result, v)
	}
	return result
}
