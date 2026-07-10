package tools

import "fmt"

func GetDistanceStringFromMeters(meters float64) string {
	return FormatKm(meters / 1000)
}

// FormatKmh formats a speed in km/h to two decimal places with (km/h) suffix.
func FormatKmh(kmh float64) string {
	return fmt.Sprintf("%.2f (km/h)", kmh)
}

// FormatKm formats a distance in kilometers to two decimal places with (km) suffix.
func FormatKm(km float64) string {
	return fmt.Sprintf("%.2f (km)", km)
}

func FormatPercentage(number float64) string {
	percentage := number * 100
	formattedPercentage := fmt.Sprintf("%.2f%%", percentage)
	return formattedPercentage
}

func To4NumberFormatter(numba int) string {
	return fmt.Sprintf("%0*d", 4, numba)
}
