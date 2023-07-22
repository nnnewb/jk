package slices

// Chunk returns a new slice containing slices of size n from the original slice.
func Chunk[T any](s []T, n int) [][]T {
	// Initialize an empty slice of slices
	var result [][]T

	// Iterate over the original slice in steps of size n
	for i := 0; i < len(s); i += n {
		// Calculate the end index of the current chunk
		end := i + n
		// If the end index is greater than the length of the slice, set it to the length of the slice
		if end > len(s) {
			end = len(s)
		}
		// Append the current chunk to the result slice
		result = append(result, s[i:end])
	}

	// Return the result slice
	return result
}
