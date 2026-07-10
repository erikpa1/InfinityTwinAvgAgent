package tools

import (
	"fmt"
	"strconv"
	"strings"
	"turtle/core/lgr"
)

func ToIArray[T any](data []T) []interface{} {
	interfaceArray := make([]interface{}, len(data))

	for i, element := range data {
		interfaceArray[i] = element
	}
	return interfaceArray
}

func ToAnotherArray[SRC any, DEST any](data []SRC) ([]DEST, error) {
	interfaceArray := make([]DEST, len(data))

	for i, element := range data {
		destElement, ok := any(element).(DEST)
		if !ok {
			return nil, fmt.Errorf("element at index %d cannot be converted to type", i)
		}
		interfaceArray[i] = destElement
	}
	return interfaceArray, nil
}

func MineValueFromDictList[MAP_T comparable, DATA_T any](data []map[MAP_T]interface{}, key MAP_T) []DATA_T {
	result := make([]DATA_T, len(data))

	for i, element := range data {
		result[i] = element[key].(DATA_T)
	}

	return result

}

func F64fromString(value string) float64 {
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		lgr.Error("%d", f)
		return 0
	}
	return f
}

func F32fromString(value string) float32 {
	f, err := strconv.ParseFloat(value, 32)
	if err != nil {
		lgr.Error("%d", f)
		return 0
	}
	return float32(f)
}

func Int32fromString(value string) int32 {
	if value == "" {
		return 0
	}
	f, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		lgr.Error(value, err)
	}
	return int32(f)
}

func Int64fromString(value string) int64 {
	if value == "" {
		return 0
	}
	f, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		lgr.Error(value, err)
		return 0
	}
	return f
}

func Int8FromString(value string) int8 {
	if value == "" {
		return 0
	}
	f, err := strconv.ParseInt(value, 10, 8)
	if err != nil {
		lgr.Error(value, err)
		return 0
	}
	return int8(f)
}

func Float64FromString(value string) float64 {
	if value == "" {
		return 0
	}
	f, err := strconv.ParseFloat(strings.Replace(value, ",", ".", -1), 64)
	if err != nil {
		lgr.Error(value, err)
		return 0
	}
	return f
}

func InterfaceToInt64(value interface{}) int64 {
	switch v := value.(type) {

	case int:
		return int64(v)
	case int64:
		return v
	case float64:
		return int64(v)
	case string:
		return Int64fromString(v)
	default:
		lgr.Error("Unsupported type:", v)
		return 0
	}
}
