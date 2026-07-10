package tools

import (
	"fmt"
	"reflect"
	"strings"
	"turtle/interfaces"
)

func GetUidOrIdBson(entity interface{}) string {

	_, ok := entity.(interfaces.UidMap)

	if ok {
		return "uid"
	}

	t := reflect.TypeOf(entity)

	if t.Kind() == reflect.Ptr {
		t = t.Elem() // Dereference if pointer
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		bsonTag := field.Tag.Get("bson")

		// Split bson tag on ',' to handle options like "_id,omitempty"
		bsonParts := strings.Split(bsonTag, ",")
		for _, part := range bsonParts {
			if strings.TrimSpace(part) == "_id" {
				return "_id" // Found _id, return immediately
			}
		}
	}

	return "uid" // Default to "uid" if _id is not found
}

func GetStringField(obj interface{}, fieldName string) (string, bool) {
	val, ok := GetField[string](obj, fieldName)
	return val, ok
}

func GetField[T any](obj interface{}, fieldName string) (T, bool) {
	var zero T

	// Get the value of the object
	val := reflect.ValueOf(obj)

	// If it's a pointer, get the value it points to
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Check if the object is a struct
	if val.Kind() != reflect.Struct {
		return zero, false
	}

	// Try to get the field
	field := val.FieldByName(fieldName)

	// Check if the field exists and is accessible
	if !field.IsValid() {
		return zero, false
	}

	// Get the field value as interface{}
	fieldValue := field.Interface()

	// Try to cast the field value to type T
	if castedValue, ok := fieldValue.(T); ok {
		return castedValue, true
	}

	return zero, false
}

func GetMapKeys(mapValue interface{}) []string {
	// Use reflection to handle any map type
	v := reflect.ValueOf(mapValue)

	// Check if the input is a map
	if v.Kind() != reflect.Map {
		return nil // Return nil if not a map
	}

	// Get all keys
	keys := v.MapKeys()
	result := make([]string, len(keys))

	// Convert each key to string
	for i, key := range keys {
		// If key is already a string, use it directly
		if key.Kind() == reflect.String {
			result[i] = key.String()
		} else {
			// Otherwise convert to string using fmt.Sprint
			result[i] = fmt.Sprint(key.Interface())
		}
	}

	return result
}
