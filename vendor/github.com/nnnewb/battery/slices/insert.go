package slices

// Insert inserts an element at the specified index in the slice and returns the modified slice.
// If the index is out of range, it appends the element to the end of the slice.
func Insert[T any](s []T, idx int, elem T) []T {
	// First, we check if the index is out of range or not
	if idx > len(s) {
		// If it is out of range, we append the element to the end of the slice
		return append(s, elem)
	} else {
		// If it is within range, we insert the element at the specified index
		// We create a new slice with the same capacity as the original slice
		newSlice := make([]T, len(s)+1)
		// We copy the elements before the specified index to the new slice
		copy(newSlice, s[:idx])
		// We insert the element at the specified index
		newSlice[idx] = elem
		// We copy the elements after the specified index to the new slice
		copy(newSlice[idx+1:], s[idx:])
		return newSlice
	}
}
