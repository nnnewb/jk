package slices

// TakeWhile returns a new slice containing the longest prefix of the original slice for which the given function returns true.
func (s Slice[T]) TakeWhile(predicate func(T) bool) Slice[T] {
	// Initialize an empty slice to hold the prefix elements
	prefix := make([]T, 0)

	// Iterate over the original slice
	for _, elem := range s {

		// If the predicate returns false, break out of the loop
		if !predicate(elem) {
			break
		}

		// Append the current element to the prefix slice
		prefix = append(prefix, elem)
	}

	// Return the prefix slice as a new Slice[T]
	return prefix
}
