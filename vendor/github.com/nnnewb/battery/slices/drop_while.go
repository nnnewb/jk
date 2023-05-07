package slices

// DropWhile returns a new slice containing all but the longest prefix of the original slice for which the given function returns true.
func (s Slice[T]) DropWhile(predicate func(T) bool) Slice[T] {
	// Initialize a variable i to 0
	i := 0

	// Iterate over the slice s using a for loop
	for i < len(s) && predicate(s[i]) {
		i++
	}

	// Return a new slice containing all but the longest prefix of the original slice for which the given function returns true
	return s[i:]
}
