package tools

// Copied on https://turtleSim/blob/main/TurtleSimCpp/src/serialization/safejson.cpp

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"turtle/core/lgr"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SafeJson is just a map. Methods are attached to a named type so we can keep
// the convenient `sj.GetString(...)` syntax, but you can still treat it as a
// plain map: `sj["foo"] = 42`, range over it, pass it to json.Marshal, etc.
type SafeJson map[string]interface{}

// OptimizeSpace controls whether Write* helpers skip zero/empty values.
// It used to be a per-instance flag; now it's package-wide since SafeJson
// is just a map and has no instance state.
var OptimizeSpace = true

// NewSafeJson initializes an empty SafeJson.
func NewSafeJson() SafeJson {
	return make(SafeJson)
}

// UnmarshalToSafeJson parses a JSON string into a new SafeJson.
func UnmarshalToSafeJson(str string) SafeJson {
	data := make(SafeJson)
	if err := json.Unmarshal([]byte(str), &data); err != nil {
		lgr.Error("UnmarshalToSafeJson: %v", err)
	}
	return data
}

// ParseString parses a JSON string into the SafeJson, replacing its contents.
func (s SafeJson) ParseString(str string) error {
	// Clear existing entries so the receiver reflects only the new data.
	for k := range s {
		delete(s, k)
	}
	tmp := make(map[string]interface{})
	if err := json.Unmarshal([]byte(str), &tmp); err != nil {
		lgr.Error("ParseString: %v", err)
		return err
	}
	for k, v := range tmp {
		s[k] = v
	}
	return nil
}

// ---------------------------------------------------------------------------
// Number casting helper
// ---------------------------------------------------------------------------

// castToFloat64 tries very hard to turn anything numeric-ish into a float64.
// It accepts every Go numeric kind, json.Number, and numeric strings.
// Returns (value, true) on success, (0, false) on failure.
func castToFloat64(v interface{}) (float64, bool) {
	if v == nil {
		return 0, false
	}

	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int8:
		return float64(n), true
	case int16:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case uint:
		return float64(n), true
	case uint8:
		return float64(n), true
	case uint16:
		return float64(n), true
	case uint32:
		return float64(n), true
	case uint64:
		return float64(n), true
	case json.Number:
		if f, err := n.Float64(); err == nil {
			return f, true
		}
	case string:
		if f, err := strconv.ParseFloat(n, 64); err == nil {
			return f, true
		}
	case bool:
		if n {
			return 1, true
		}
		return 0, true
	}

	// Last-resort reflection path for named numeric types
	// (e.g. `type Seconds int64`, custom enums backed by ints, etc.).
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Float32, reflect.Float64:
		return rv.Float(), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(rv.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(rv.Uint()), true
	}

	return 0, false
}

// ---------------------------------------------------------------------------
// String / object id / interface
// ---------------------------------------------------------------------------

// GetString retrieves a string value or a default if the key is not found.
func (s SafeJson) GetString(key, notFound string) string {
	val, hasVal := s[key]
	if !hasVal {
		return notFound
	}
	if str, ok := val.(string); ok {
		return str
	}
	lgr.Error("GetString: value at key %q is not a string (got %T)", key, val)
	return notFound
}

// GetPrimitiveObjectId retrieves a Mongo ObjectID stored as a hex string.
func (s SafeJson) GetPrimitiveObjectId(key string) primitive.ObjectID {
	val, hasVal := s[key]
	if !hasVal {
		return primitive.ObjectID{}
	}
	str, ok := val.(string)
	if !ok {
		lgr.Error("GetPrimitiveObjectId: value at key %q is not a string (got %T)", key, val)
		return primitive.ObjectID{}
	}
	oid, err := primitive.ObjectIDFromHex(str)
	if err != nil {
		lgr.Error("GetPrimitiveObjectId: invalid hex at key %q: %v", key, err)
		return primitive.ObjectID{}
	}
	return oid
}

// GetInterface returns the raw value at key or notFound if absent.
func (s SafeJson) GetInterface(key string, notFound interface{}) interface{} {
	if val, ok := s[key]; ok {
		return val
	}
	return notFound
}

// GetStringMatrix returns a [][]string at the given key, or an empty slice.
func (s SafeJson) GetStringMatrix(key string) [][]string {
	data, ok := s[key]
	if !ok {
		return [][]string{}
	}
	marshaled, err := json.Marshal(data)
	if err != nil {
		lgr.Error("GetStringMatrix: marshal at key %q failed: %v", key, err)
		return [][]string{}
	}
	tmp := make([][]string, 0)
	if err := json.Unmarshal(marshaled, &tmp); err != nil {
		lgr.Error("GetStringMatrix: unmarshal at key %q failed: %v", key, err)
		return [][]string{}
	}
	return tmp
}

// ---------------------------------------------------------------------------
// Numeric getters — all use castToFloat64 so any numeric format works
// ---------------------------------------------------------------------------

// GetInt retrieves an int value, casting from any numeric type.
func (s SafeJson) GetInt(key string, notFound int) int {
	val, hasVal := s[key]
	if !hasVal {
		return notFound
	}
	if f, ok := castToFloat64(val); ok {
		return int(f)
	}
	lgr.Error("GetInt: value at key %q is not numeric (got %T)", key, val)
	return notFound
}

// GetInt64 retrieves an int64 value, casting from any numeric type.
func (s SafeJson) GetInt64(key string, notFound int64) int64 {
	val, hasVal := s[key]
	if !hasVal {
		return notFound
	}
	if f, ok := castToFloat64(val); ok {
		return int64(f)
	}
	lgr.Error("GetInt64: value at key %q is not numeric (got %T)", key, val)
	return notFound
}

// GetInt8 retrieves an int8 value, casting from any numeric type.
func (s SafeJson) GetInt8(key string, notFound int8) int8 {
	val, hasVal := s[key]
	if !hasVal {
		return notFound
	}
	if f, ok := castToFloat64(val); ok {
		return int8(f)
	}
	lgr.Error("GetInt8: value at key %q is not numeric (got %T)", key, val)
	return notFound
}

// GetSeconds retrieves a Seconds value, casting from any numeric type.
func (s SafeJson) GetSeconds(key string, notFound Seconds) Seconds {
	val, hasVal := s[key]
	if !hasVal {
		return notFound
	}
	if f, ok := castToFloat64(val); ok {
		return Seconds(f)
	}
	lgr.Error("GetSeconds: value at key %q is not numeric (got %T)", key, val)
	return notFound
}

// GetFloat64 retrieves a float64 value, casting from any numeric type.
func (s SafeJson) GetFloat64(key string, notFound float64) float64 {
	val, hasVal := s[key]
	if !hasVal {
		return notFound
	}
	if f, ok := castToFloat64(val); ok {
		return f
	}
	lgr.Error("GetFloat64: value at key %q is not numeric (got %T)", key, val)
	return notFound
}

// GetLong retrieves an int64 value, casting from any numeric type.
func (s SafeJson) GetLong(key string, notFound int64) int64 {
	val, hasVal := s[key]
	if !hasVal {
		return notFound
	}
	if f, ok := castToFloat64(val); ok {
		return int64(f)
	}
	lgr.Error("GetLong: value at key %q is not numeric (got %T)", key, val)
	return notFound
}

// GetDouble retrieves a float64 value, casting from any numeric type.
func (s SafeJson) GetDouble(key string, notFound float64) float64 {
	val, hasVal := s[key]
	if !hasVal {
		return notFound
	}
	if f, ok := castToFloat64(val); ok {
		return f
	}
	lgr.Error("GetDouble: value at key %q is not numeric (got %T)", key, val)
	return notFound
}

// GetBool retrieves a bool value or a default if the key is not found.
func (s SafeJson) GetBool(key string, notFound bool) bool {
	val, hasVal := s[key]
	if !hasVal {
		return notFound
	}
	if b, ok := val.(bool); ok {
		return b
	}
	// Accept numeric truthiness too (0 -> false, anything else -> true)
	// and the strings "true"/"false"/"1"/"0".
	if f, ok := castToFloat64(val); ok {
		return f != 0
	}
	if str, ok := val.(string); ok {
		if b, err := strconv.ParseBool(str); err == nil {
			return b
		}
	}
	lgr.Error("GetBool: value at key %q is not boolean-ish (got %T)", key, val)
	return notFound
}

// ---------------------------------------------------------------------------
// Nested objects & arrays
// ---------------------------------------------------------------------------

// GetSafeJson retrieves a nested object as a SafeJson. Returns an empty
// SafeJson (not nil) if the key is missing or wrong type, so callers can
// chain without nil checks.
func (s SafeJson) GetSafeJson(key string) SafeJson {
	val, hasVal := s[key]
	if !hasVal {
		return NewSafeJson()
	}
	switch m := val.(type) {
	case map[string]interface{}:
		return SafeJson(m)
	case SafeJson:
		return m
	}
	lgr.Error("GetSafeJson: value at key %q is not an object (got %T)", key, val)
	return NewSafeJson()
}

// GetObjectArray retrieves an array of objects as []SafeJson.
func (s SafeJson) GetObjectArray(key string) []SafeJson {
	raw, hasVal := s[key]
	if !hasVal {
		return nil
	}

	// Tento kus kodu tu musi byt kvoli subtypovaniu golangovych struktur
	value := reflect.ValueOf(raw)
	if value.Kind() != reflect.Slice {
		lgr.Error("GetObjectArray: value at key %q is not a slice (got %T)", key, raw)
		return nil
	}

	result := make([]SafeJson, 0, value.Len())
	mapType := reflect.TypeOf(map[string]interface{}{})

	for i := 0; i < value.Len(); i++ {
		item := value.Index(i)

		// Unwrap interface{} if necessary.
		if item.Kind() == reflect.Interface {
			item = item.Elem()
		}

		switch item.Kind() {
		case reflect.Map:
			if item.Type().ConvertibleTo(mapType) {
				converted := item.Convert(mapType).Interface().(map[string]interface{})
				result = append(result, SafeJson(converted))
			} else {
				lgr.Error("GetObjectArray: map at index %d in key %q not convertible to map[string]interface{}", i, key)
			}
		case reflect.Ptr:
			if !item.IsNil() && item.Elem().Kind() == reflect.Map && item.Elem().Type().ConvertibleTo(mapType) {
				converted := item.Elem().Convert(mapType).Interface().(map[string]interface{})
				result = append(result, SafeJson(converted))
			} else {
				lgr.Error("GetObjectArray: pointer at index %d in key %q not a *map[string]interface{}", i, key)
			}
		default:
			lgr.Error("GetObjectArray: item at index %d in key %q is not an object (got %s)", i, key, item.Kind())
		}
	}
	return result
}

// GetFloats64Array retrieves a []float64, casting each element from any numeric type.
func (s SafeJson) GetFloats64Array(key string, defaultValue []float64) []float64 {
	raw, hasVal := s[key]
	if !hasVal {
		return defaultValue
	}

	value := reflect.ValueOf(raw)
	if value.Kind() != reflect.Slice {
		lgr.Error("GetFloats64Array: value at key %q is not a slice (got %T)", key, raw)
		return defaultValue
	}

	result := make([]float64, value.Len())
	for i := 0; i < value.Len(); i++ {
		item := value.Index(i).Interface()
		if f, ok := castToFloat64(item); ok {
			result[i] = f
		} else {
			lgr.Error("GetFloats64Array: element %d at key %q is not numeric (got %T)", i, key, item)
		}
	}
	return result
}

// GetStringsArray retrieves a []string from the JSON.
func (s SafeJson) GetStringsArray(key string) []string {
	raw, hasVal := s[key]
	if !hasVal {
		return []string{}
	}

	value := reflect.ValueOf(raw)
	if value.Kind() != reflect.Slice {
		lgr.Error("GetStringsArray: value at key %q is not a slice (got %T)", key, raw)
		return []string{}
	}

	result := make([]string, value.Len())
	for i := 0; i < value.Len(); i++ {
		item := value.Index(i).Interface()
		if str, ok := item.(string); ok {
			result[i] = str
		} else {
			lgr.Error("GetStringsArray: element %d at key %q is not a string (got %T)", i, key, item)
			result[i] = fmt.Sprintf("%v", item) // best-effort
		}
	}
	return result
}

// GetFloatsPositionArray retrieves a fixed-size [3]float64, casting each element.
func (s SafeJson) GetFloatsPositionArray(key string, defaultValue [3]float64) [3]float64 {
	result := defaultValue

	raw, hasVal := s[key]
	if !hasVal {
		return result
	}

	value := reflect.ValueOf(raw)
	if value.Kind() != reflect.Slice {
		lgr.Error("GetFloatsPositionArray: value at key %q is not a slice (got %T)", key, raw)
		return result
	}

	n := value.Len()
	if n > 3 {
		n = 3
	}
	for i := 0; i < n; i++ {
		item := value.Index(i).Interface()
		if f, ok := castToFloat64(item); ok {
			result[i] = f
		} else {
			lgr.Error("GetFloatsPositionArray: element %d at key %q is not numeric (got %T)", i, key, item)
			result[i] = 0
		}
	}
	return result
}

// ---------------------------------------------------------------------------
// Writers
// ---------------------------------------------------------------------------

// WriteString writes a string value to the JSON.
func (s SafeJson) WriteString(key, value string) {
	if !OptimizeSpace || value != "" {
		s[key] = value
	}
}

// WriteFloat writes a float32 value to the JSON.
func (s SafeJson) WriteFloat(key string, value float32) {
	if !OptimizeSpace || value != 0.0 {
		s[key] = value
	}
}

// WriteFloat3 writes an array of three float32 values to the JSON.
func (s SafeJson) WriteFloat3(key string, x, y, z float32) {
	s[key] = []float32{x, y, z}
}

// WriteDouble writes a float64 value to the JSON.
func (s SafeJson) WriteDouble(key string, value float64) {
	if !OptimizeSpace || value != 0 {
		s[key] = value
	}
}

// WriteDouble3 writes an array of three float64 values to the JSON.
func (s SafeJson) WriteDouble3(key string, x, y, z float64) {
	s[key] = []float64{x, y, z}
}

// WriteInt writes an int value to the JSON.
func (s SafeJson) WriteInt(key string, value int) {
	if !OptimizeSpace || value != 0 {
		s[key] = value
	}
}

// WriteLong writes an int64 value to the JSON.
func (s SafeJson) WriteLong(key string, value int64) {
	if !OptimizeSpace || value != 0 {
		s[key] = value
	}
}

// WriteBoolean writes a bool value to the JSON.
func (s SafeJson) WriteBoolean(key string, value bool) {
	if !OptimizeSpace || value {
		s[key] = value
	}
}

// WriteJson writes another SafeJson as a nested object.
func (s SafeJson) WriteJson(key string, another SafeJson) {
	if another != nil {
		s[key] = map[string]interface{}(another)
	}
}

// WriteJsonArray writes an array of SafeJson objects.
func (s SafeJson) WriteJsonArray(key string, array []SafeJson) {
	if OptimizeSpace && len(array) == 0 {
		return
	}
	jsonArray := make([]map[string]interface{}, 0, len(array))
	for _, item := range array {
		jsonArray = append(jsonArray, map[string]interface{}(item))
	}
	s[key] = jsonArray
}

// GetKeysLength returns the number of top-level keys.
func (s SafeJson) GetKeysLength() int {
	return len(s)
}

// Dump returns the JSON as a pretty-printed string.
func (s SafeJson) Dump() string {
	bytes, err := json.MarshalIndent(map[string]interface{}(s), "", "  ")
	if err != nil {
		lgr.Error("Dump: %v", err)
		return "{}"
	}
	return string(bytes)
}

// ---------------------------------------------------------------------------
// Free functions (unchanged in behavior, just kept here)
// ---------------------------------------------------------------------------

func MarshalWithoutDefaults(v interface{}) ([]byte, error) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a struct or pointer to a struct")
	}

	nonZeroFields := make(map[string]interface{})
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		if !field.CanInterface() {
			continue
		}
		zeroValue := reflect.Zero(field.Type()).Interface()
		if !reflect.DeepEqual(field.Interface(), zeroValue) {
			nonZeroFields[fieldType.Name] = field.Interface()
		}
	}
	return json.Marshal(nonZeroFields)
}

func FindFirstJsonString(s string) (string, int, error) {
	start := -1
	for i, char := range s {
		if char == '{' || char == '[' {
			start = i
			break
		}
	}
	if start == -1 {
		return "", -1, fmt.Errorf("no JSON start found")
	}

	openBraces := 0
	openBrackets := 0
	inString := false
	escaped := false

	for i := start; i < len(s); i++ {
		char := rune(s[i])

		if escaped {
			escaped = false
			continue
		}
		if char == '\\' && inString {
			escaped = true
			continue
		}
		if char == '"' {
			inString = !inString
			continue
		}
		if !inString {
			switch char {
			case '{':
				openBraces++
			case '}':
				openBraces--
			case '[':
				openBrackets++
			case ']':
				openBrackets--
			}
			if openBraces == 0 && openBrackets == 0 {
				candidate := s[start : i+1]
				var temp interface{}
				if json.Unmarshal([]byte(candidate), &temp) == nil {
					return candidate, start, nil
				}
			}
		}
	}

	return "", -1, fmt.Errorf("no valid JSON found")
}
