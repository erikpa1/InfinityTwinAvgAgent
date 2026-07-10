package tools

import (
	"fmt"
	"math"
)

// PercentageFormatting struct (not strictly necessary for just the function, but mimicking the class concept)
type PercentageFormatting struct{}

// ToString method to format a float as a percentage
func (pf PercentageFormatting) ToString(numba float64) string {
	// Convert number to percentage and round to two decimal places
	percentage := numba * 100
	roundedPercentage := math.Round(percentage*100) / 100

	// Format the rounded percentage with a "%" symbol
	formattedPercentage := fmt.Sprintf("%.2f%%", roundedPercentage)

	return formattedPercentage
}
