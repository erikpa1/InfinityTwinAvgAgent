package tools

import "math/rand"

func RandomFloatInRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func RandomInt64InRange(min, max int64) int64 {

	if min > max {
		min, max = max, min
	}

	delta := max - min + 1

	var randomNum int64
	if delta <= 0 {
		randomNum = min + rand.Int63()%(max-min+1)
	} else {
		randomNum = min + rand.Int63n(delta)
	}

	return randomNum
}
