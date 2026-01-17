// Package utils provides utility functions for USC platform services.
package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
)

// JSONUtils provides JSON utility functions
type JSONUtils struct{}

// NewJSONUtils creates a new JSON utils instance
func NewJSONUtils() *JSONUtils {
	return &JSONUtils{}
}

// Marshal marshals a value to JSON
func (ju *JSONUtils) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// MarshalIndent marshals a value to JSON with indentation
func (ju *JSONUtils) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

// Unmarshal unmarshals JSON data to a value
func (ju *JSONUtils) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// UnmarshalFromReader unmarshals JSON data from a reader
func (ju *JSONUtils) UnmarshalFromReader(reader io.Reader, v interface{}) error {
	return json.NewDecoder(reader).Decode(v)
}

// MarshalToWriter marshals a value to JSON and writes to a writer
func (ju *JSONUtils) MarshalToWriter(v interface{}, writer io.Writer) error {
	return json.NewEncoder(writer).Encode(v)
}

// PrettyPrint pretty prints a value as JSON
func (ju *JSONUtils) PrettyPrint(v interface{}) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// IsValidJSON checks if a string is valid JSON
func (ju *JSONUtils) IsValidJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

// MinifyJSON minifies a JSON string
func (ju *JSONUtils) MinifyJSON(jsonStr string) (string, error) {
	var v interface{}
	if err := json.Unmarshal([]byte(jsonStr), &v); err != nil {
		return "", err
	}

	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// MergeJSON merges two JSON objects
func (ju *JSONUtils) MergeJSON(json1, json2 string) (string, error) {
	var obj1, obj2 map[string]interface{}

	if err := json.Unmarshal([]byte(json1), &obj1); err != nil {
		return "", err
	}

	if err := json.Unmarshal([]byte(json2), &obj2); err != nil {
		return "", err
	}

	// Merge obj2 into obj1
	for key, value := range obj2 {
		obj1[key] = value
	}

	result, err := json.Marshal(obj1)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// GetJSONValue gets a value from JSON using a path
func (ju *JSONUtils) GetJSONValue(jsonStr, path string) (interface{}, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return nil, err
	}

	return ju.getValueByPath(data, path)
}

// SetJSONValue sets a value in JSON using a path
func (ju *JSONUtils) SetJSONValue(jsonStr, path string, value interface{}) (string, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return "", err
	}

	if err := ju.setValueByPath(data, path, value); err != nil {
		return "", err
	}

	result, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// RemoveJSONValue removes a value from JSON using a path
func (ju *JSONUtils) RemoveJSONValue(jsonStr, path string) (string, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return "", err
	}

	if err := ju.removeValueByPath(data, path); err != nil {
		return "", err
	}

	result, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// getValueByPath gets a value from a nested structure using a path
func (ju *JSONUtils) getValueByPath(data interface{}, path string) (interface{}, error) {
	keys := strings.Split(path, ".")
	current := data

	for _, key := range keys {
		switch v := current.(type) {
		case map[string]interface{}:
			if val, exists := v[key]; exists {
				current = val
			} else {
				return nil, fmt.Errorf("key %s not found", key)
			}
		case []interface{}:
			// Handle array index
			if key == "" {
				return current, nil
			}
			// For now, we'll assume it's a map key
			return nil, fmt.Errorf("cannot access key %s on array", key)
		default:
			return nil, fmt.Errorf("cannot access key %s on type %T", key, current)
		}
	}

	return current, nil
}

// setValueByPath sets a value in a nested structure using a path
func (ju *JSONUtils) setValueByPath(data interface{}, path string, value interface{}) error {
	keys := strings.Split(path, ".")
	current := data

	for i, key := range keys {
		if i == len(keys)-1 {
			// Last key, set the value
			switch v := current.(type) {
			case map[string]interface{}:
				v[key] = value
			default:
				return fmt.Errorf("cannot set key %s on type %T", key, current)
			}
		} else {
			// Navigate deeper
			switch v := current.(type) {
			case map[string]interface{}:
				if val, exists := v[key]; exists {
					current = val
				} else {
					// Create new map
					newMap := make(map[string]interface{})
					v[key] = newMap
					current = newMap
				}
			default:
				return fmt.Errorf("cannot access key %s on type %T", key, current)
			}
		}
	}

	return nil
}

// removeValueByPath removes a value from a nested structure using a path
func (ju *JSONUtils) removeValueByPath(data interface{}, path string) error {
	keys := strings.Split(path, ".")
	current := data

	for i, key := range keys {
		if i == len(keys)-1 {
			// Last key, remove the value
			switch v := current.(type) {
			case map[string]interface{}:
				delete(v, key)
			default:
				return fmt.Errorf("cannot remove key %s from type %T", key, current)
			}
		} else {
			// Navigate deeper
			switch v := current.(type) {
			case map[string]interface{}:
				if val, exists := v[key]; exists {
					current = val
				} else {
					return fmt.Errorf("key %s not found", key)
				}
			default:
				return fmt.Errorf("cannot access key %s on type %T", key, current)
			}
		}
	}

	return nil
}

// JSONSchema represents a JSON schema
type JSONSchema struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Required   []string               `json:"required,omitempty"`
	Items      *JSONSchema            `json:"items,omitempty"`
	MinItems   int                    `json:"minItems,omitempty"`
	MaxItems   int                    `json:"maxItems,omitempty"`
	MinLength  int                    `json:"minLength,omitempty"`
	MaxLength  int                    `json:"maxLength,omitempty"`
	Minimum    float64                `json:"minimum,omitempty"`
	Maximum    float64                `json:"maximum,omitempty"`
	Pattern    string                 `json:"pattern,omitempty"`
	Format     string                 `json:"format,omitempty"`
	Enum       []interface{}          `json:"enum,omitempty"`
	Default    interface{}            `json:"default,omitempty"`
}

// ValidateJSON validates JSON against a schema
func (ju *JSONUtils) ValidateJSON(jsonStr string, schema *JSONSchema) error {
	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return err
	}

	return ju.validateValue(data, schema)
}

// validateValue validates a value against a schema
func (ju *JSONUtils) validateValue(value interface{}, schema *JSONSchema) error {
	if schema == nil {
		return nil
	}

	// Check type
	if err := ju.validateType(value, schema.Type); err != nil {
		return err
	}

	// Validate based on type
	switch schema.Type {
	case "object":
		return ju.validateObject(value, schema)
	case "array":
		return ju.validateArray(value, schema)
	case "string":
		return ju.validateString(value, schema)
	case "number":
		return ju.validateNumber(value, schema)
	case "integer":
		return ju.validateInteger(value, schema)
	case "boolean":
		return ju.validateBoolean(value, schema)
	}

	return nil
}

// validateType validates the type of a value
func (ju *JSONUtils) validateType(value interface{}, expectedType string) error {
	actualType := reflect.TypeOf(value).Kind().String()

	switch expectedType {
	case "string":
		if actualType != "string" {
			return fmt.Errorf("expected string, got %s", actualType)
		}
	case "number":
		if actualType != "float64" {
			return fmt.Errorf("expected number, got %s", actualType)
		}
	case "integer":
		if actualType != "float64" {
			return fmt.Errorf("expected integer, got %s", actualType)
		}
	case "boolean":
		if actualType != "bool" {
			return fmt.Errorf("expected boolean, got %s", actualType)
		}
	case "array":
		if actualType != "slice" {
			return fmt.Errorf("expected array, got %s", actualType)
		}
	case "object":
		if actualType != "map" {
			return fmt.Errorf("expected object, got %s", actualType)
		}
	}

	return nil
}

// validateObject validates an object
func (ju *JSONUtils) validateObject(value interface{}, schema *JSONSchema) error {
	obj, ok := value.(map[string]interface{})
	if !ok {
		return fmt.Errorf("expected object")
	}

	// Check required fields
	for _, required := range schema.Required {
		if _, exists := obj[required]; !exists {
			return fmt.Errorf("required field %s is missing", required)
		}
	}

	// Validate properties
	for key, val := range obj {
		if propSchema, exists := schema.Properties[key]; exists {
			if propSchemaMap, ok := propSchema.(map[string]interface{}); ok {
				// Convert to JSONSchema
				propSchemaBytes, err := json.Marshal(propSchemaMap)
				if err != nil {
					return err
				}
				var propSchema JSONSchema
				if err := json.Unmarshal(propSchemaBytes, &propSchema); err != nil {
					return err
				}

				if err := ju.validateValue(val, &propSchema); err != nil {
					return fmt.Errorf("validation failed for field %s: %v", key, err)
				}
			}
		}
	}

	return nil
}

// validateArray validates an array
func (ju *JSONUtils) validateArray(value interface{}, schema *JSONSchema) error {
	arr, ok := value.([]interface{})
	if !ok {
		return fmt.Errorf("expected array")
	}

	// Check length constraints
	if schema.MinItems > 0 && len(arr) < schema.MinItems {
		return fmt.Errorf("array must have at least %d items", schema.MinItems)
	}
	if schema.MaxItems > 0 && len(arr) > schema.MaxItems {
		return fmt.Errorf("array must have at most %d items", schema.MaxItems)
	}

	// Validate items
	if schema.Items != nil {
		for i, item := range arr {
			if err := ju.validateValue(item, schema.Items); err != nil {
				return fmt.Errorf("validation failed for array item %d: %v", i, err)
			}
		}
	}

	return nil
}

// validateString validates a string
func (ju *JSONUtils) validateString(value interface{}, schema *JSONSchema) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string")
	}

	// Check length constraints
	if schema.MinLength > 0 && len(str) < schema.MinLength {
		return fmt.Errorf("string must be at least %d characters long", schema.MinLength)
	}
	if schema.MaxLength > 0 && len(str) > schema.MaxLength {
		return fmt.Errorf("string must be at most %d characters long", schema.MaxLength)
	}

	// Check pattern
	if schema.Pattern != "" {
		// Simple pattern matching - in a real implementation, you would use regex
		if !strings.Contains(str, schema.Pattern) {
			return fmt.Errorf("string does not match pattern %s", schema.Pattern)
		}
	}

	// Check enum
	if len(schema.Enum) > 0 {
		found := false
		for _, enumVal := range schema.Enum {
			if str == enumVal {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("string must be one of %v", schema.Enum)
		}
	}

	return nil
}

// validateNumber validates a number
func (ju *JSONUtils) validateNumber(value interface{}, schema *JSONSchema) error {
	num, ok := value.(float64)
	if !ok {
		return fmt.Errorf("expected number")
	}

	// Check range constraints
	if schema.Minimum != 0 && num < schema.Minimum {
		return fmt.Errorf("number must be at least %f", schema.Minimum)
	}
	if schema.Maximum != 0 && num > schema.Maximum {
		return fmt.Errorf("number must be at most %f", schema.Maximum)
	}

	return nil
}

// validateInteger validates an integer
func (ju *JSONUtils) validateInteger(value interface{}, schema *JSONSchema) error {
	num, ok := value.(float64)
	if !ok {
		return fmt.Errorf("expected integer")
	}

	// Check if it's actually an integer
	if num != float64(int64(num)) {
		return fmt.Errorf("expected integer, got float")
	}

	// Check range constraints
	if schema.Minimum != 0 && num < schema.Minimum {
		return fmt.Errorf("integer must be at least %f", schema.Minimum)
	}
	if schema.Maximum != 0 && num > schema.Maximum {
		return fmt.Errorf("integer must be at most %f", schema.Maximum)
	}

	return nil
}

// validateBoolean validates a boolean
func (ju *JSONUtils) validateBoolean(value interface{}, schema *JSONSchema) error {
	_, ok := value.(bool)
	if !ok {
		return fmt.Errorf("expected boolean")
	}

	return nil
}
