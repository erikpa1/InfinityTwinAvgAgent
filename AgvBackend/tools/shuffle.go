package tools

import (
	"math/rand"
)

// shuffle shuffles an array of float64 in place.
func ShuffleArray[T any](arr_ptr *[]T) {
	arr := *arr_ptr

	for i := range arr {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i] // Swap elements
	}
}
