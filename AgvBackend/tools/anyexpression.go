package tools

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"turtle/tools/timeexpr"
)

func AnyExpr_CompileMilis(expr string, defaultValue Milliseconds) Milliseconds {

	if expr == "" {
		return defaultValue
	} else {
		tmp := MathExpr_Execute(expr)

		if tmp != nil {
			return Milliseconds(*tmp * 1000)
		} else {
			return Milliseconds(timeexpr.SecondsFromTimeString(expr) * 1000)
		}
	}

}

func AnyExpr_CompileSeconds(expr string, defaultValue float64) float64 {

	if expr == "" {
		return defaultValue
	} else {
		tmp := MathExpr_Execute(expr)

		if tmp != nil {
			return *tmp
		} else {
			return float64(timeexpr.SecondsFromTimeString(expr))
		}
	}

}

func MathExpr_Execute(expr string) *float64 {
	// Check if expr is a number
	if num, err := strconv.ParseFloat(expr, 64); err == nil {
		return &num
	}

	uniformRegex := regexp.MustCompile(`uniform\(([^,]+),\s*([^)]+)\)`)
	standardRangeRegex := regexp.MustCompile(`standard\(([^,]+),\s*([^)]+)\)`)
	standardMaxRegex := regexp.MustCompile(`standard\(([^,]+)\)`)

	if match := uniformRegex.FindStringSubmatch(expr); match != nil {

		timeexpr.SecondsFromTimeString(match[1])

		minVal := float64(timeexpr.SecondsFromTimeString(match[1]))
		maxVal := float64(timeexpr.SecondsFromTimeString(match[2]))
		result := minVal + rand.Float64()*(maxVal-minVal)
		return &result
	}

	if match := standardRangeRegex.FindStringSubmatch(expr); match != nil {
		minVal := float64(timeexpr.SecondsFromTimeString(match[1]))
		maxVal := float64(timeexpr.SecondsFromTimeString(match[2]))
		value := rand.NormFloat64()
		result := value*(maxVal-minVal) + minVal
		return &result
	}

	if match := standardMaxRegex.FindStringSubmatch(expr); match != nil {
		maxVal := float64(timeexpr.SecondsFromTimeString(match[1]))
		value := rand.NormFloat64()
		result := value * maxVal
		return &result
	}

	// Attempt to evaluate the mathematical expression (using a simple parser or eval library)
	// For simplicity, this part is skipped. You can use third-party libraries for expression evaluation.
	// E.g., https://github.com/Knetic/govaluate
	return nil
}

func AnyNumberConvert[T any](value interface{}) (T, error) {

	tmp, ok := value.(*T)

	if ok {
		return *tmp, nil
	} else {
		var tmp T
		return tmp, fmt.Errorf("Failed to convert")
	}

}
