package slices

// Contains takes in a slice of type T and a value of type T as
// arguments.
//
// It then iterates over the slice and checks if the value is equal to any of
// the elements in the slice.
//
// If it finds a match, it returns true, otherwise it returns false.
//
// This function can be used to check if a specific value exists in a slice.
func Contains[T comparable](s []T, v T) bool {
	for _, val := range s {
		if val == v {
			return true
		}
	}
	return false
}

// ContainsFunc takes in a slice of any type T, a value of type T, and a
// function that compares two values of type T and returns a boolean.
//
// It iterates through the slice and checks if any of the values in the
// slice are equal to the given value using the provided comparison function.
//
// If a match is found, it returns true. If no match is found, it returns false.
//
// This function can be used to check if a value exists in a slice using a
// custom comparison function.
func ContainsFunc[T any](s []T, v T, equal func(T, T) bool) bool {
	for _, val := range s {
		if equal(val, v) {
			return true
		}
	}
	return false
}
