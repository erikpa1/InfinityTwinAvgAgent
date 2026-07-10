package timeexpr

import (
	"fmt"
	"strconv"
	"strings"
)

// MakeFromSecondsFull converts seconds to a formatted time string with full details (days, hours, minutes, seconds).
func MakeFromSecondsFull(seconds float64) string {
	return MakeFromMillisFull(seconds * 1000)
}

// MakeFromSeconds converts seconds to a formatted time string (hours, minutes, seconds).
func MakeFromSeconds(seconds float64) string {
	return MakeFromMillis(seconds * 1000)
}

// MakeFromMillis converts milliseconds to a formatted time string (hours, minutes, seconds).
func MakeFromMillis(millis float64) string {
	if millis == 0 {
		return "00:00"
	}

	seconds := millis / 1000
	days, rem := divmod(seconds, 86400)
	hours, rem := divmod(rem, 3600)
	minutes, seconds := divmod(rem, 60)

	var daysStr, hoursStr, minutesStr, secondsStr string

	if days > 0 {
		daysStr = fmt.Sprintf("%d:", int(days))
		hoursStr = fmt.Sprintf("%02d:", int(hours))
		minutesStr = fmt.Sprintf("%02d:", int(minutes))
		secondsStr = fmt.Sprintf("%02d", int(seconds))
	} else if hours > 0 {
		hoursStr = fmt.Sprintf("%d:", int(hours))
		minutesStr = fmt.Sprintf("%02d:", int(minutes))
		secondsStr = fmt.Sprintf("%02d", int(seconds))
	} else if minutes > 0 {
		minutesStr = fmt.Sprintf("%d:", int(minutes))
		secondsStr = fmt.Sprintf("%02d", int(seconds))
	} else {
		secondsStr = fmt.Sprintf("%02d", int(seconds))
	}

	if days > 0 {
		return fmt.Sprintf("%s%s%s%s", daysStr, hoursStr, minutesStr, secondsStr)
	} else if hours > 0 {
		return fmt.Sprintf("%s%s%s", hoursStr, minutesStr, secondsStr)
	} else if minutes > 0 {
		return fmt.Sprintf("%s%s", minutesStr, secondsStr)
	}
	return fmt.Sprintf("00:%s", secondsStr)
}

func MakeFromMillisPretty(millis float64) string {
	if millis == 0 {
		return "0s"
	}

	seconds := millis / 1000
	days, rem := divmod(seconds, 86400)
	hours, rem := divmod(rem, 3600)
	minutes, seconds := divmod(rem, 60)

	var parts []string

	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", int(days)))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", int(hours)))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", int(minutes)))
	}
	// Always include seconds
	parts = append(parts, fmt.Sprintf("%ds", int(seconds)))

	return fmt.Sprintf("%s", join(parts, " "))
}

func join(parts []string, sep string) string {
	if len(parts) == 0 {
		return ""
	}
	result := parts[0]
	for _, part := range parts[1:] {
		result += sep + part
	}
	return result
}

// MakeFromMillisFull converts milliseconds to a full formatted time string with days, hours, minutes, and seconds.
func MakeFromMillisFull(millis float64) string {
	if millis == 0 {
		return "00:00"
	}

	seconds := millis / 1000
	days, rem := divmod(seconds, 86400)
	hours, rem := divmod(rem, 3600)
	minutes, seconds := divmod(rem, 60)

	daysStr := fmt.Sprintf("%02d:", int(days))
	hoursStr := fmt.Sprintf("%02d:", int(hours))
	minutesStr := fmt.Sprintf("%02d:", int(minutes))
	secondsStr := fmt.Sprintf("%02d", int(seconds))

	return fmt.Sprintf("%s%s%s%s", daysStr, hoursStr, minutesStr, secondsStr)
}

// SecondsFromTimeString converts a time string (e.g. "1:02:30") to seconds.
func SecondsFromTimeString(timeString string) int64 {
	return MillisFromTimeString(timeString) / 1000
}

// MillisFromTimeString converts a time string (e.g. "1:02:30") to milliseconds.
func MillisFromTimeString(timeString string) int64 {
	components := strings.Split(timeString, ":")
	// Initialize variables for days, hours, minutes, and seconds
	var days, hours, minutes, seconds int

	// Start parsing from the last component (seconds)
	if len(components) > 0 {
		seconds, _ = strconv.Atoi(components[len(components)-1])
		components = components[:len(components)-1]
	}
	if len(components) > 0 {
		minutes, _ = strconv.Atoi(components[len(components)-1])
		components = components[:len(components)-1]
	}
	if len(components) > 0 {
		hours, _ = strconv.Atoi(components[len(components)-1])
		components = components[:len(components)-1]
	}
	if len(components) > 0 {
		days, _ = strconv.Atoi(components[len(components)-1])
	}

	// Calculate total milliseconds
	totalMillis := (days*24*60*60 + hours*60*60 + minutes*60 + seconds) * 1000

	return int64(totalMillis)
}

// divmod function for Go (returns quotient and remainder)
func divmod(a, b float64) (quotient, remainder float64) {
	quotient = float64(int(a) / int(b))
	remainder = a - quotient*b
	return quotient, remainder
}
