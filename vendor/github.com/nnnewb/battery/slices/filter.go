package slices

// Filter returns a new slice containing only the elements of the original slice for which the predicate function returns true
// The predicate function takes an element of type T as input and returns a boolean value
func Filter[T any](s []T, predicate func(T) bool) []T {
	if len(s) == 0 {
		return make([]T, 0)
	}

	var filteredSlice = make([]T, 0)
	for _, elem := range s {
		if predicate(elem) {
			filteredSlice = append(filteredSlice, elem)
		}
	}

	return filteredSlice
}
