package slices

func (s Slice[T]) Concat(other ...Slice[T]) Slice[T] {
	resultCap := len(s)
	for i := 0; i < len(other); i++ {
		resultCap += len(other[i])
	}
	result := make(Slice[T], 0, resultCap)
	result = append(result, s...)
	for i := 0; i < len(other); i++ {
		result = append(result, other[i]...)
	}
	return result
}
