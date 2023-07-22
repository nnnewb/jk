package slices

import (
	"math/rand"
	"time"
)

// Choices returns a new Slice[T] containing n randomly chosen elements from the original Slice[T].
// Returned slices may contain duplicated element.
// The original Slice[T] is not modified.
func Choices[T any](s []T, n int) []T {
	// First, we need to check if n is greater than the length of the slice.
	// If it is, we return a copy of the original slice.
	if n >= len(s) {
		return Copy(s)
	}

	// We create a new slice with capacity n.
	// This is more efficient than appending to a slice with a length of 0.
	result := make([]T, 0, n)

	// We create a new rand.Source using the current time as the seed.
	// This is not cryptographically secure, but it's good enough for our purposes.
	source := rand.NewSource(time.Now().UnixNano())

	// We create a new rand.Rand using the source we just created.
	random := rand.New(source)

	// We loop n times.
	for i := 0; i < n; i++ {
		// We generate a random index between 0 and the length of the slice.
		index := random.Intn(len(s))

		// We append the element at the random index to the result slice.
		result = append(result, s[index])
	}

	// We return the result slice.
	return result
}
