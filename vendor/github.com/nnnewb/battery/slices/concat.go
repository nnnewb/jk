package slices

func Concat[T any](s []T, other ...[]T) []T {
	resultCap := len(s)
	for i := 0; i < len(other); i++ {
		resultCap += len(other[i])
	}
	result := make([]T, 0, resultCap)
	result = append(result, s...)
	for i := 0; i < len(other); i++ {
		result = append(result, other[i]...)
	}
	return result
}
