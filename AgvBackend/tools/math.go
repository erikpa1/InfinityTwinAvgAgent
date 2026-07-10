package tools

import (
	"math"
)

type F64Vec3 [3]float64

func CalculateNextPosition(currentX, currentZ, destX, destZ, speed float64) ([3]float64, float64) {
	// Calculate distance between current position and destination
	distanceX := destX - currentX
	distanceZ := destZ - currentZ

	// Calculate the total distance to the destination
	totalDistance := math.Sqrt(distanceX*distanceX + distanceZ*distanceZ)

	// If speed exceeds the distance, cap it to reach the destination
	if speed >= totalDistance {
		overflowSpeed := speed - totalDistance
		return [3]float64{destX, 0, destZ}, overflowSpeed // Speed for next point is zero, as we are at destination
	}

	// Normalize direction vector
	unitX := distanceX / totalDistance
	unitZ := distanceZ / totalDistance

	// Calculate next position
	nextX := currentX + unitX*speed
	nextZ := currentZ + unitZ*speed

	return [3]float64{nextX, 0, nextZ}, 0
}

func (self *F64Vec3) IsZero() bool {
	return self[0] == 0 && self[1] == 0 && self[2] == 0
}

type F64Vec3Struct struct {
	X float64 `json:"x" bson:"x" mapstructure:"x"`
	Y float64 `json:"y" bson:"y" mapstructure:"y"`
	Z float64 `json:"z" bson:"z" mapstructure:"z"`
}

func StaticF64Vec3Struct() F64Vec3Struct {
	return F64Vec3Struct{
		X: 1,
		Y: 1,
		Z: 1,
	}
}
