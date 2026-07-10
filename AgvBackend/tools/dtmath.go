package tools

import "math"

// CalculateDistance calculates the distance between two points (x1, y1) and (x2, y2).
func CalculateDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}

// IsInZone checks if a point (posX, posZ) is inside a rectangular zone centered at (centerX, centerZ) with dimensions scaleX, scaleY, and scaleZ.
func IsInZone(posX, posZ, centerX, centerZ, scaleX, scaleZ float64) bool {
	halfWidth := scaleX / 2
	halfLength := scaleZ / 2

	minX := centerX - halfWidth
	maxX := centerX + halfWidth
	minZ := centerZ - halfLength
	maxZ := centerZ + halfLength

	return posX >= minX && posX <= maxX && posZ >= minZ && posZ <= maxZ
}

// IsInPolygonZone checks if a point (posX, posZ) is inside a polygon defined by a list of vertices (polygonPoints).
func IsInPolygonZone(posX, posZ float64, polygonPoints [][2]float64) bool {
	numVertices := len(polygonPoints)
	inside := false

	for i := 0; i < numVertices; i++ {
		p1 := polygonPoints[i]
		p2 := polygonPoints[(i+1)%numVertices]

		if posZ > math.Min(p1[1], p2[1]) {
			if posZ <= math.Max(p1[1], p2[1]) {
				if posX <= math.Max(p1[0], p2[0]) {
					xIntersection := (posZ-p1[1])*(p2[0]-p1[0])/(p2[1]-p1[1]) + p1[0]
					if p1[0] == p2[0] || posX <= xIntersection {
						inside = !inside
					}
				}
			}
		}
	}

	return inside
}

// LineIntersectsZone checks if a line segment between (x1, z1) and (x2, z2) intersects with a rectangular zone.
func LineIntersectsZone(x1, z1, x2, z2, centerX, centerZ, scaleX, scaleZ float64) bool {
	halfWidth := scaleX / 2
	halfLength := scaleZ / 2

	minX := centerX - halfWidth
	maxX := centerX + halfWidth
	minZ := centerZ - halfLength
	maxZ := centerZ + halfLength

	if IsInZone(x1, z1, centerX, centerZ, scaleX, scaleZ) || IsInZone(x2, z2, centerX, centerZ, scaleX, scaleZ) {
		return true
	}

	slope := (z2 - z1) / (x2 - x1)
	intercept := z1 - slope*x1

	if (minX <= x1 && x1 <= maxX && minZ <= z1 && z1 <= maxZ) || (minX <= x2 && x2 <= maxX && minZ <= z2 && z2 <= maxZ) {
		return true
	}

	if minX <= x1 && x1 <= maxX {
		zIntersection := slope*x1 + intercept
		if minZ <= zIntersection && zIntersection <= maxZ {
			return true
		}
	}

	if minZ <= z1 && z1 <= maxZ {
		xIntersection := (z1 - intercept) / slope
		if minX <= xIntersection && xIntersection <= maxX {
			return true
		}
	}

	return false
}

func IsInRectangle(posX, posZ float64, rectangle Rectangle) bool {
	return IsInZone(posX, posZ, rectangle.posX, rectangle.posZ, rectangle.scaleX, rectangle.scaleZ)
}

func MaxInt64() int64 {
	return (int64(^uint64(0) >> 1))
}

func MaxSeconds() Seconds {
	return Seconds(MaxInt64())
}

func MaxInt32() int32 {
	return int32(^uint32(0) >> 1)
}

func MaxFloat64() float64 {
	return math.MaxFloat64
}

func MaxFloat32() float32 {
	return math.MaxFloat32
}
