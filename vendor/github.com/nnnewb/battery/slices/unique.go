package slices

// Unique returns a new slice containing only the unique elements of the original slice, in the order they first appear.
func (s Slice[T]) Unique(equal func(x, y T) bool) Slice[T] {
	uniqueSlice := make(Slice[T], 0, len(s))
	for _, element := range s {
		isUnique := true
		for _, uniqueElement := range uniqueSlice {
			if equal(element, uniqueElement) {
				isUnique = false
				break
			}
		}
		if isUnique {
			uniqueSlice = append(uniqueSlice, element)
		}
	}
	return uniqueSlice
}

// Unique returns a new slice containing only the unique elements of the original slice, in the order they first appear.
// It uses a map to improve performance compared to the original Unique function.
func Unique[T comparable](s Slice[T]) Slice[T] {
	uniqueSlice := make(Slice[T], 0, len(s))
	seen := make(map[T]bool)
	for _, elem := range s {
		if !seen[elem] {
			uniqueSlice = append(uniqueSlice, elem)
			seen[elem] = true
		}
	}
	return uniqueSlice
}
