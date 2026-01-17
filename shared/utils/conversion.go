// Package utils provides utility functions for USC platform services.
package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// ConversionUtils provides conversion utility functions
type ConversionUtils struct{}

// NewConversionUtils creates a new conversion utils instance
func NewConversionUtils() *ConversionUtils {
	return &ConversionUtils{}
}

// ToString converts a value to string
func (cu *ConversionUtils) ToString(value interface{}) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case time.Time:
		return v.Format(time.RFC3339)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// ToInt converts a value to int
func (cu *ConversionUtils) ToInt(value interface{}) (int, error) {
	if value == nil {
		return 0, fmt.Errorf("cannot convert nil to int")
	}

	switch v := value.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to int", value)
	}
}

// ToInt64 converts a value to int64
func (cu *ConversionUtils) ToInt64(value interface{}) (int64, error) {
	if value == nil {
		return 0, fmt.Errorf("cannot convert nil to int64")
	}

	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to int64", value)
	}
}

// ToFloat64 converts a value to float64
func (cu *ConversionUtils) ToFloat64(value interface{}) (float64, error) {
	if value == nil {
		return 0, fmt.Errorf("cannot convert nil to float64")
	}

	switch v := value.(type) {
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float64:
		return v, nil
	case string:
		return strconv.ParseFloat(v, 64)
	case bool:
		if v {
			return 1.0, nil
		}
		return 0.0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", value)
	}
}

// ToBool converts a value to bool
func (cu *ConversionUtils) ToBool(value interface{}) (bool, error) {
	if value == nil {
		return false, fmt.Errorf("cannot convert nil to bool")
	}

	switch v := value.(type) {
	case bool:
		return v, nil
	case int:
		return v != 0, nil
	case int64:
		return v != 0, nil
	case float64:
		return v != 0, nil
	case string:
		return strconv.ParseBool(v)
	default:
		return false, fmt.Errorf("cannot convert %T to bool", value)
	}
}

// ToTime converts a value to time.Time
func (cu *ConversionUtils) ToTime(value interface{}) (time.Time, error) {
	if value == nil {
		return time.Time{}, fmt.Errorf("cannot convert nil to time")
	}

	switch v := value.(type) {
	case time.Time:
		return v, nil
	case string:
		layouts := []string{
			time.RFC3339,
			time.RFC3339Nano,
			time.RFC1123,
			time.RFC1123Z,
			time.RFC822,
			time.RFC822Z,
			time.RFC850,
			time.ANSIC,
			time.UnixDate,
			time.RubyDate,
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05",
			"2006-01-02",
			"15:04:05",
		}

		for _, layout := range layouts {
			if t, err := time.Parse(layout, v); err == nil {
				return t, nil
			}
		}

		return time.Time{}, fmt.Errorf("cannot parse time string: %s", v)
	case int64:
		return time.Unix(v, 0), nil
	case int:
		return time.Unix(int64(v), 0), nil
	default:
		return time.Time{}, fmt.Errorf("cannot convert %T to time", value)
	}
}

// ToSlice converts a value to a slice
func (cu *ConversionUtils) ToSlice(value interface{}) ([]interface{}, error) {
	if value == nil {
		return nil, fmt.Errorf("cannot convert nil to slice")
	}

	switch v := value.(type) {
	case []interface{}:
		return v, nil
	case []string:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = item
		}
		return result, nil
	case []int:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = item
		}
		return result, nil
	case []int64:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = item
		}
		return result, nil
	case []float64:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = item
		}
		return result, nil
	case []bool:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = item
		}
		return result, nil
	default:
		// Try to convert using reflection
		rv := reflect.ValueOf(value)
		if rv.Kind() == reflect.Slice {
			result := make([]interface{}, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				result[i] = rv.Index(i).Interface()
			}
			return result, nil
		}
		return nil, fmt.Errorf("cannot convert %T to slice", value)
	}
}

// ToMap converts a value to a map
func (cu *ConversionUtils) ToMap(value interface{}) (map[string]interface{}, error) {
	if value == nil {
		return nil, fmt.Errorf("cannot convert nil to map")
	}

	switch v := value.(type) {
	case map[string]interface{}:
		return v, nil
	case map[string]string:
		result := make(map[string]interface{})
		for key, val := range v {
			result[key] = val
		}
		return result, nil
	case map[string]int:
		result := make(map[string]interface{})
		for key, val := range v {
			result[key] = val
		}
		return result, nil
	case map[string]int64:
		result := make(map[string]interface{})
		for key, val := range v {
			result[key] = val
		}
		return result, nil
	case map[string]float64:
		result := make(map[string]interface{})
		for key, val := range v {
			result[key] = val
		}
		return result, nil
	case map[string]bool:
		result := make(map[string]interface{})
		for key, val := range v {
			result[key] = val
		}
		return result, nil
	default:
		// Try to convert using reflection
		rv := reflect.ValueOf(value)
		if rv.Kind() == reflect.Map {
			result := make(map[string]interface{})
			for _, key := range rv.MapKeys() {
				result[cu.ToString(key.Interface())] = rv.MapIndex(key).Interface()
			}
			return result, nil
		}
		return nil, fmt.Errorf("cannot convert %T to map", value)
	}
}

// ConvertTo converts a value to a specific type
func (cu *ConversionUtils) ConvertTo(value interface{}, targetType reflect.Type) (interface{}, error) {
	if value == nil {
		return nil, fmt.Errorf("cannot convert nil to %s", targetType)
	}

	// Check if already the correct type
	if reflect.TypeOf(value) == targetType {
		return value, nil
	}

	// Convert based on target type
	switch targetType.Kind() {
	case reflect.String:
		return cu.ToString(value), nil
	case reflect.Int:
		return cu.ToInt(value)
	case reflect.Int64:
		return cu.ToInt64(value)
	case reflect.Float64:
		return cu.ToFloat64(value)
	case reflect.Bool:
		return cu.ToBool(value)
	case reflect.Slice:
		return cu.ToSlice(value)
	case reflect.Map:
		return cu.ToMap(value)
	default:
		return nil, fmt.Errorf("unsupported target type: %s", targetType)
	}
}

// ConvertSlice converts a slice to a different type
func (cu *ConversionUtils) ConvertSlice(slice interface{}, targetType reflect.Type) (interface{}, error) {
	if slice == nil {
		return nil, fmt.Errorf("cannot convert nil slice")
	}

	// Get the slice as []interface{}
	interfaceSlice, err := cu.ToSlice(slice)
	if err != nil {
		return nil, err
	}

	// Create a new slice of the target type
	result := reflect.MakeSlice(reflect.SliceOf(targetType), len(interfaceSlice), len(interfaceSlice))

	// Convert each element
	for i, item := range interfaceSlice {
		converted, err := cu.ConvertTo(item, targetType)
		if err != nil {
			return nil, fmt.Errorf("failed to convert element %d: %v", i, err)
		}
		result.Index(i).Set(reflect.ValueOf(converted))
	}

	return result.Interface(), nil
}

// ConvertMap converts a map to a different type
func (cu *ConversionUtils) ConvertMap(m map[string]interface{}, targetType reflect.Type) (interface{}, error) {
	if m == nil {
		return nil, fmt.Errorf("cannot convert nil map")
	}

	// Create a new map of the target type
	result := reflect.MakeMap(reflect.MapOf(reflect.TypeOf(""), targetType))

	// Convert each value
	for key, value := range m {
		converted, err := cu.ConvertTo(value, targetType)
		if err != nil {
			return nil, fmt.Errorf("failed to convert value for key %s: %v", key, err)
		}
		result.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(converted))
	}

	return result.Interface(), nil
}

// ConvertStruct converts a struct to a map
func (cu *ConversionUtils) ConvertStruct(s interface{}) (map[string]interface{}, error) {
	if s == nil {
		return nil, fmt.Errorf("cannot convert nil struct")
	}

	rv := reflect.ValueOf(s)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %s", rv.Kind())
	}

	result := make(map[string]interface{})
	rt := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		field := rt.Field(i)
		value := rv.Field(i)

		// Skip unexported fields
		if !value.CanInterface() {
			continue
		}

		// Get field name (use json tag if available)
		fieldName := field.Name
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "" {
				fieldName = parts[0]
			}
		}

		result[fieldName] = value.Interface()
	}

	return result, nil
}

// ConvertMapToStruct converts a map to a struct
func (cu *ConversionUtils) ConvertMapToStruct(m map[string]interface{}, targetStruct interface{}) error {
	if m == nil {
		return fmt.Errorf("cannot convert nil map")
	}

	rv := reflect.ValueOf(targetStruct)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer to struct")
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to struct")
	}

	rt := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		field := rt.Field(i)
		value := rv.Field(i)

		// Skip unexported fields
		if !value.CanSet() {
			continue
		}

		// Get field name (use json tag if available)
		fieldName := field.Name
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "" {
				fieldName = parts[0]
			}
		}

		// Get value from map
		mapValue, exists := m[fieldName]
		if !exists {
			continue
		}

		// Convert and set the value
		converted, err := cu.ConvertTo(mapValue, field.Type)
		if err != nil {
			return fmt.Errorf("failed to convert field %s: %v", fieldName, err)
		}

		value.Set(reflect.ValueOf(converted))
	}

	return nil
}

// ConvertStringToType converts a string to a specific type
func (cu *ConversionUtils) ConvertStringToType(str string, targetType reflect.Type) (interface{}, error) {
	if str == "" {
		return nil, fmt.Errorf("cannot convert empty string")
	}

	switch targetType.Kind() {
	case reflect.String:
		return str, nil
	case reflect.Int:
		return strconv.Atoi(str)
	case reflect.Int64:
		return strconv.ParseInt(str, 10, 64)
	case reflect.Float64:
		return strconv.ParseFloat(str, 64)
	case reflect.Bool:
		return strconv.ParseBool(str)
	default:
		return nil, fmt.Errorf("unsupported target type: %s", targetType)
	}
}
