package slices

// Reverse create a new slice with the same length as the original slice
// Traverse the original slice from the end to the beginning and copy each element to the new slice
// Return the new slice
func Reverse[T any](s []T) []T {
	reversed := make([]T, len(s))
	for i, j := len(s)-1, 0; i >= 0; i, j = i-1, j+1 {
		reversed[j] = s[i]
	}
	return reversed
}
