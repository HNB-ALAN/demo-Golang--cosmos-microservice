// Package validation provides validation utilities for USC platform services.
package validation

import (
	"reflect"
	"strings"
)

// Validator provides validation functionality
type Validator struct {
	tagName string
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{
		tagName: "validate",
	}
}

// NewValidatorWithTag creates a new validator with custom tag name
func NewValidatorWithTag(tagName string) *Validator {
	return &Validator{
		tagName: tagName,
	}
}

// Validate validates a struct using struct tags
func (v *Validator) Validate(obj interface{}) ValidationErrors {
	var validationErrors ValidationErrors

	val := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)

	// Handle pointers
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return validationErrors
		}
		val = val.Elem()
		typ = typ.Elem()
	}

	// Only validate structs
	if val.Kind() != reflect.Struct {
		return validationErrors
	}

	// Iterate through struct fields
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}

		// Get validation tag
		tag := fieldType.Tag.Get(v.tagName)
		if tag == "" {
			continue
		}

		// Validate field
		fieldErrors := v.validateField(field, fieldType, tag)
		validationErrors = append(validationErrors, fieldErrors...)
	}

	return validationErrors
}

// validateField validates a single field
func (v *Validator) validateField(field reflect.Value, fieldType reflect.StructField, tag string) ValidationErrors {
	var validationErrors ValidationErrors

	// Parse validation rules
	rules := strings.Split(tag, ",")

	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}

		// Parse rule (e.g., "required", "min=5", "max=10")
		parts := strings.Split(rule, "=")
		ruleName := parts[0]
		ruleValue := ""
		if len(parts) > 1 {
			ruleValue = parts[1]
		}

		// Apply validation rule
		if err := v.applyRule(field, fieldType, ruleName, ruleValue); err != nil {
			validationErrors = append(validationErrors, *err)
		}
	}

	return validationErrors
}

// applyRule applies a specific validation rule
func (v *Validator) applyRule(field reflect.Value, fieldType reflect.StructField, ruleName, ruleValue string) *ValidationError {
	switch ruleName {
	case "required":
		return v.validateRequired(field, fieldType)
	case "min":
		return v.validateMin(field, fieldType, ruleValue)
	case "max":
		return v.validateMax(field, fieldType, ruleValue)
	case "email":
		return v.validateEmail(field, fieldType)
	case "url":
		return v.validateURL(field, fieldType)
	case "numeric":
		return v.validateNumeric(field, fieldType)
	case "alpha":
		return v.validateAlpha(field, fieldType)
	case "alphanumeric":
		return v.validateAlphanumeric(field, fieldType)
	default:
		return nil
	}
}

// validateRequired validates that a field is not empty
func (v *Validator) validateRequired(field reflect.Value, fieldType reflect.StructField) *ValidationError {
	if v.isEmpty(field) {
		return &ValidationError{
			Field:   fieldType.Name,
			Tag:     "required",
			Value:   field.Interface(),
			Message: fieldType.Name + " is required",
		}
	}
	return nil
}

// validateMin validates minimum value/length
func (v *Validator) validateMin(field reflect.Value, fieldType reflect.StructField, min string) *ValidationError {
	if v.isEmpty(field) {
		return nil // Skip empty fields
	}

	switch field.Kind() {
	case reflect.String:
		if len(field.String()) < v.parseInt(min) {
			return &ValidationError{
				Field:   fieldType.Name,
				Tag:     "min",
				Value:   field.Interface(),
				Message: fieldType.Name + " must be at least " + min + " characters",
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Int() < int64(v.parseInt(min)) {
			return &ValidationError{
				Field:   fieldType.Name,
				Tag:     "min",
				Value:   field.Interface(),
				Message: fieldType.Name + " must be at least " + min,
			}
		}
	case reflect.Float32, reflect.Float64:
		if field.Float() < float64(v.parseFloat(min)) {
			return &ValidationError{
				Field:   fieldType.Name,
				Tag:     "min",
				Value:   field.Interface(),
				Message: fieldType.Name + " must be at least " + min,
			}
		}
	}
	return nil
}

// validateMax validates maximum value/length
func (v *Validator) validateMax(field reflect.Value, fieldType reflect.StructField, max string) *ValidationError {
	if v.isEmpty(field) {
		return nil // Skip empty fields
	}

	switch field.Kind() {
	case reflect.String:
		if len(field.String()) > v.parseInt(max) {
			return &ValidationError{
				Field:   fieldType.Name,
				Tag:     "max",
				Value:   field.Interface(),
				Message: fieldType.Name + " must be at most " + max + " characters",
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Int() > int64(v.parseInt(max)) {
			return &ValidationError{
				Field:   fieldType.Name,
				Tag:     "max",
				Value:   field.Interface(),
				Message: fieldType.Name + " must be at most " + max,
			}
		}
	case reflect.Float32, reflect.Float64:
		if field.Float() > float64(v.parseFloat(max)) {
			return &ValidationError{
				Field:   fieldType.Name,
				Tag:     "max",
				Value:   field.Interface(),
				Message: fieldType.Name + " must be at most " + max,
			}
		}
	}
	return nil
}

// validateEmail validates email format
func (v *Validator) validateEmail(field reflect.Value, fieldType reflect.StructField) *ValidationError {
	if v.isEmpty(field) {
		return nil // Skip empty fields
	}

	if field.Kind() != reflect.String {
		return nil
	}

	email := field.String()
	if !v.isValidEmail(email) {
		return &ValidationError{
			Field:   fieldType.Name,
			Tag:     "email",
			Value:   field.Interface(),
			Message: fieldType.Name + " must be a valid email address",
		}
	}
	return nil
}

// validateURL validates URL format
func (v *Validator) validateURL(field reflect.Value, fieldType reflect.StructField) *ValidationError {
	if v.isEmpty(field) {
		return nil // Skip empty fields
	}

	if field.Kind() != reflect.String {
		return nil
	}

	url := field.String()
	if !v.isValidURL(url) {
		return &ValidationError{
			Field:   fieldType.Name,
			Tag:     "url",
			Value:   field.Interface(),
			Message: fieldType.Name + " must be a valid URL",
		}
	}
	return nil
}

// validateNumeric validates that a field contains only numeric characters
func (v *Validator) validateNumeric(field reflect.Value, fieldType reflect.StructField) *ValidationError {
	if v.isEmpty(field) {
		return nil // Skip empty fields
	}

	if field.Kind() != reflect.String {
		return nil
	}

	str := field.String()
	if !v.isNumeric(str) {
		return &ValidationError{
			Field:   fieldType.Name,
			Tag:     "numeric",
			Value:   field.Interface(),
			Message: fieldType.Name + " must contain only numeric characters",
		}
	}
	return nil
}

// validateAlpha validates that a field contains only alphabetic characters
func (v *Validator) validateAlpha(field reflect.Value, fieldType reflect.StructField) *ValidationError {
	if v.isEmpty(field) {
		return nil // Skip empty fields
	}

	if field.Kind() != reflect.String {
		return nil
	}

	str := field.String()
	if !v.isAlpha(str) {
		return &ValidationError{
			Field:   fieldType.Name,
			Tag:     "alpha",
			Value:   field.Interface(),
			Message: fieldType.Name + " must contain only alphabetic characters",
		}
	}
	return nil
}

// validateAlphanumeric validates that a field contains only alphanumeric characters
func (v *Validator) validateAlphanumeric(field reflect.Value, fieldType reflect.StructField) *ValidationError {
	if v.isEmpty(field) {
		return nil // Skip empty fields
	}

	if field.Kind() != reflect.String {
		return nil
	}

	str := field.String()
	if !v.isAlphanumeric(str) {
		return &ValidationError{
			Field:   fieldType.Name,
			Tag:     "alphanumeric",
			Value:   field.Interface(),
			Message: fieldType.Name + " must contain only alphanumeric characters",
		}
	}
	return nil
}

// Helper methods

// isEmpty checks if a field is empty
func (v *Validator) isEmpty(field reflect.Value) bool {
	switch field.Kind() {
	case reflect.String:
		return field.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() == 0
	case reflect.Float32, reflect.Float64:
		return field.Float() == 0
	case reflect.Bool:
		return !field.Bool()
	case reflect.Slice, reflect.Array, reflect.Map:
		return field.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return field.IsNil()
	default:
		return false
	}
}

// parseInt parses a string to int
func (v *Validator) parseInt(s string) int {
	// Simple implementation - in production, use strconv.Atoi
	if s == "" {
		return 0
	}
	result := 0
	for _, char := range s {
		if char >= '0' && char <= '9' {
			result = result*10 + int(char-'0')
		} else {
			break
		}
	}
	return result
}

// parseFloat parses a string to float64
func (v *Validator) parseFloat(s string) float64 {
	// Simple implementation - in production, use strconv.ParseFloat
	if s == "" {
		return 0
	}
	result := 0.0
	divisor := 1.0
	decimal := false
	for _, char := range s {
		if char >= '0' && char <= '9' {
			if decimal {
				divisor *= 10
				result += float64(char-'0') / divisor
			} else {
				result = result*10 + float64(char-'0')
			}
		} else if char == '.' && !decimal {
			decimal = true
		} else {
			break
		}
	}
	return result
}

// isValidEmail checks if a string is a valid email
func (v *Validator) isValidEmail(email string) bool {
	// Simple email validation - in production, use a proper regex
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// isValidURL checks if a string is a valid URL
func (v *Validator) isValidURL(url string) bool {
	// Simple URL validation - in production, use a proper regex
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

// isNumeric checks if a string contains only numeric characters
func (v *Validator) isNumeric(s string) bool {
	for _, char := range s {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

// isAlpha checks if a string contains only alphabetic characters
func (v *Validator) isAlpha(s string) bool {
	for _, char := range s {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')) {
			return false
		}
	}
	return true
}

// isAlphanumeric checks if a string contains only alphanumeric characters
func (v *Validator) isAlphanumeric(s string) bool {
	for _, char := range s {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
			return false
		}
	}
	return true
}
