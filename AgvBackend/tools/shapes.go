package tools

import (
	"math"
	"turtle/core/lgr"
)

type Rectangle struct {
	posX, posZ     float64
	scaleX, scaleZ float64
	x1, z1, x2, z2 float64
}

// Constructor
func NewRectangle() *Rectangle {
	return &Rectangle{
		posX:   0,
		posZ:   0,
		scaleX: 0,
		scaleZ: 0,
		x1:     0,
		z1:     0,
		x2:     0,
		z2:     0,
	}
}

// SetFromMiddle method
func (r *Rectangle) SetFromMiddle(midX, midZ, scaleX, scaleZ float64) {
	r.posX = midX
	r.posZ = midZ
	r.scaleX = scaleX
	r.scaleZ = scaleZ

	r.RecalculateX1Y1X2Z2()
}

// SetFrom2PointsRectangle method
func (r *Rectangle) SetFrom2PointsRectangle(x1, z1, x2, z2 float64) {
	r.scaleX = math.Abs(x1 - x2)
	r.scaleZ = math.Abs(z1 - z2)

	r.posX = x1 + (r.scaleX * 0.5)
	r.posZ = z1 - (r.scaleZ * 0.5)

	r.RecalculateXZAndScale()
}

// RecalculateX1Y1X2Z2 method (placeholder)
func (r *Rectangle) RecalculateX1Y1X2Z2() {
	lgr.Info("TODOOO")
}

// RecalculateXZAndScale method (placeholder)
func (r *Rectangle) RecalculateXZAndScale() {
	lgr.Info("TODOOO")
}
