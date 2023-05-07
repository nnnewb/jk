package slices

// Filter returns a new slice containing only the elements of the original slice for which the predicate function returns true
// The predicate function takes an element of type T as input and returns a boolean value
func (s Slice[T]) Filter(predicate func(T) bool) Slice[T] {
	if len(s) == 0 {
		return make(Slice[T], 0)
	}

	var filteredSlice Slice[T]
	for _, elem := range s {
		if predicate(elem) {
			filteredSlice = append(filteredSlice, elem)
		}
	}

	return filteredSlice
}
