package tools

import "fmt"

func CloneMap[K comparable, V any](original map[K]V) map[K]V {
	// Create a new map
	cloned := make(map[K]V, len(original))

	// Copy each key-value pair from the original map to the new map
	for key, value := range original {
		cloned[key] = value
	}

	return cloned
}

// PopFirst removes and returns the first element of a slice
func PopFirst[T any](s []T) ([]T, T, error) {
	if len(s) == 0 {
		var zero T
		return s, zero, fmt.Errorf("cannot pop from an empty slice")
	}
	return s[1:], s[0], nil
}

// PopLast removes and returns the last element of a slice
func PopLast[T any](s []T) ([]T, T, error) {
	if len(s) == 0 {
		var zero T
		return s, zero, fmt.Errorf("cannot pop from an empty slice")
	}
	return s[:len(s)-1], s[len(s)-1], nil
}

func FetchLast[T any](s []T) (T, error) {
	if len(s) == 0 {
		var zero T
		return zero, fmt.Errorf("cannot fetch from an empty slice")
	}
	return s[len(s)-1], nil
}

func MergeVec[T any](vectors ...[]T) []T {

	size := 0

	for _, vec := range vectors {
		size += len(vec)
	}

	result := make([]T, size)

	index := 0
	for _, vec := range vectors {
		for _, element := range vec {
			result[index] = element

			index += 1
		}
	}

	return result
}

func ReverseVector[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func ReverseVectorCopy[T any](s []T) []T {
	reversed := make([]T, len(s))

	for i, v := range s {
		reversed[len(s)-1-i] = v
	}

	return reversed
}

func MapValuesToVec[K comparable, V any](m map[K]V) []V {
	// Create a slice to store the values
	values := make([]V, 0, len(m))

	// Iterate through the map and append values to the slice
	for _, v := range m {
		values = append(values, v)
	}

	return values
}
