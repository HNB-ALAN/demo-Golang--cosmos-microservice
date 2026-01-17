// Package validation provides validation utilities for USC platform services.
package validation

import (
	"fmt"
	"reflect"
	"strings"
)

// StructValidator provides struct validation functionality
type StructValidator struct {
	validator *Validator
}

// NewStructValidator creates a new struct validator
func NewStructValidator() *StructValidator {
	return &StructValidator{
		validator: NewValidator(),
	}
}

// ValidateStruct validates a struct
func (sv *StructValidator) ValidateStruct(s interface{}) ValidationErrors {
	return sv.validator.Validate(s)
}

// ValidateField validates a specific field of a struct
func (sv *StructValidator) ValidateField(s interface{}, fieldName string) ValidationErrors {
	var errors ValidationErrors

	// Get field value
	fieldValue, err := sv.getFieldValue(s, fieldName)
	if err != nil {
		errors = append(errors, ValidationError{
			Field:   fieldName,
			Tag:     "field_not_found",
			Value:   nil,
			Message: "Field not found",
		})
		return errors
	}

	// Get field tags
	fieldTags, err := sv.getFieldTags(s, fieldName)
	if err != nil {
		errors = append(errors, ValidationError{
			Field:   fieldName,
			Tag:     "tag_not_found",
			Value:   fieldValue,
			Message: "Field tags not found",
		})
		return errors
	}

	// Validate field against tags
	for _, tag := range fieldTags {
		// Simple validation - in production, implement proper field validation
		if fieldValue == nil || fieldValue == "" {
			errors = append(errors, ValidationError{
				Field:   fieldName,
				Tag:     tag,
				Value:   fieldValue,
				Message: fieldName + " is required",
			})
		}
	}

	return errors
}

// ValidateFields validates specific fields of a struct
func (sv *StructValidator) ValidateFields(s interface{}, fieldNames []string) ValidationErrors {
	var errors ValidationErrors

	for _, fieldName := range fieldNames {
		fieldErrors := sv.ValidateField(s, fieldName)
		errors = append(errors, fieldErrors...)
	}

	return errors
}

// ValidateExcept validates all fields except specified ones
func (sv *StructValidator) ValidateExcept(s interface{}, excludeFields []string) ValidationErrors {
	var errors ValidationErrors

	// Get all field names
	fieldNames := sv.getFieldNames(s)

	// Filter out excluded fields
	var fieldsToValidate []string
	for _, fieldName := range fieldNames {
		excluded := false
		for _, excludeField := range excludeFields {
			if fieldName == excludeField {
				excluded = true
				break
			}
		}
		if !excluded {
			fieldsToValidate = append(fieldsToValidate, fieldName)
		}
	}

	// Validate remaining fields
	for _, fieldName := range fieldsToValidate {
		fieldErrors := sv.ValidateField(s, fieldName)
		errors = append(errors, fieldErrors...)
	}

	return errors
}

// getFieldValue gets the value of a field
func (sv *StructValidator) getFieldValue(s interface{}, fieldName string) (interface{}, error) {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("not a struct")
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if field.Name == fieldName {
			return v.Field(i).Interface(), nil
		}
	}

	return nil, fmt.Errorf("field not found")
}

// getFieldTags gets the validation tags of a field
func (sv *StructValidator) getFieldTags(s interface{}, fieldName string) ([]string, error) {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("not a struct")
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if field.Name == fieldName {
			tag := field.Tag.Get("validate")
			if tag == "" {
				return []string{}, nil
			}
			return strings.Split(tag, ","), nil
		}
	}

	return nil, fmt.Errorf("field not found")
}

// getFieldNames gets all field names of a struct
func (sv *StructValidator) getFieldNames(s interface{}) []string {
	var fieldNames []string

	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fieldNames
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldNames = append(fieldNames, field.Name)
	}

	return fieldNames
}

// StructValidationRule represents a struct validation rule
type StructValidationRule struct {
	Field    string   `json:"field"`
	Tags     []string `json:"tags"`
	Required bool     `json:"required"`
	Message  string   `json:"message,omitempty"`
}

// StructValidationSchema represents a struct validation schema
type StructValidationSchema struct {
	Rules []StructValidationRule `json:"rules"`
}

// ValidateWithSchema validates a struct against a schema
func (sv *StructValidator) ValidateWithSchema(s interface{}, schema StructValidationSchema) ValidationErrors {
	var errors ValidationErrors

	// Convert struct to map for easier access
	dataMap := make(map[string]interface{})
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.Struct {
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			value := v.Field(i)
			dataMap[field.Name] = value.Interface()
		}
	}

	// Validate each rule
	for _, rule := range schema.Rules {
		value, exists := dataMap[rule.Field]
		if !exists && rule.Required {
			message := rule.Message
			if message == "" {
				message = fmt.Sprintf("%s is required", rule.Field)
			}
			errors = append(errors, ValidationError{
				Field:   rule.Field,
				Tag:     "required",
				Value:   nil,
				Message: message,
			})
			continue
		}

		if exists {
			for _, tag := range rule.Tags {
				// Simple validation - in production, implement proper field validation
				if value == nil || value == "" {
					message := rule.Message
					if message == "" {
						message = rule.Field + " is required"
					}

					errors = append(errors, ValidationError{
						Field:   rule.Field,
						Tag:     tag,
						Value:   value,
						Message: message,
					})
				}
			}
		}
	}

	return errors
}

// StructValidationBuilder provides a builder for struct validation rules
type StructValidationBuilder struct {
	rules []StructValidationRule
}

// NewStructValidationBuilder creates a new struct validation builder
func NewStructValidationBuilder() *StructValidationBuilder {
	return &StructValidationBuilder{
		rules: make([]StructValidationRule, 0),
	}
}

// AddRule adds a validation rule
func (svb *StructValidationBuilder) AddRule(rule StructValidationRule) *StructValidationBuilder {
	svb.rules = append(svb.rules, rule)
	return svb
}

// Required adds a required rule
func (svb *StructValidationBuilder) Required(field string, message ...string) *StructValidationBuilder {
	rule := StructValidationRule{
		Field:    field,
		Tags:     []string{"required"},
		Required: true,
	}
	if len(message) > 0 {
		rule.Message = message[0]
	}
	svb.rules = append(svb.rules, rule)
	return svb
}

// Min adds a minimum value rule
func (svb *StructValidationBuilder) Min(field string, min string, message ...string) *StructValidationBuilder {
	rule := StructValidationRule{
		Field: field,
		Tags:  []string{"min=" + min},
	}
	if len(message) > 0 {
		rule.Message = message[0]
	}
	svb.rules = append(svb.rules, rule)
	return svb
}

// Max adds a maximum value rule
func (svb *StructValidationBuilder) Max(field string, max string, message ...string) *StructValidationBuilder {
	rule := StructValidationRule{
		Field: field,
		Tags:  []string{"max=" + max},
	}
	if len(message) > 0 {
		rule.Message = message[0]
	}
	svb.rules = append(svb.rules, rule)
	return svb
}

// Email adds an email validation rule
func (svb *StructValidationBuilder) Email(field string, message ...string) *StructValidationBuilder {
	rule := StructValidationRule{
		Field: field,
		Tags:  []string{"email"},
	}
	if len(message) > 0 {
		rule.Message = message[0]
	}
	svb.rules = append(svb.rules, rule)
	return svb
}

// URL adds a URL validation rule
func (svb *StructValidationBuilder) URL(field string, message ...string) *StructValidationBuilder {
	rule := StructValidationRule{
		Field: field,
		Tags:  []string{"url"},
	}
	if len(message) > 0 {
		rule.Message = message[0]
	}
	svb.rules = append(svb.rules, rule)
	return svb
}

// UUID adds a UUID validation rule
func (svb *StructValidationBuilder) UUID(field string, message ...string) *StructValidationBuilder {
	rule := StructValidationRule{
		Field: field,
		Tags:  []string{"uuid"},
	}
	if len(message) > 0 {
		rule.Message = message[0]
	}
	svb.rules = append(svb.rules, rule)
	return svb
}

// Alphanumeric adds an alphanumeric validation rule
func (svb *StructValidationBuilder) Alphanumeric(field string, message ...string) *StructValidationBuilder {
	rule := StructValidationRule{
		Field: field,
		Tags:  []string{"alphanumeric"},
	}
	if len(message) > 0 {
		rule.Message = message[0]
	}
	svb.rules = append(svb.rules, rule)
	return svb
}

// StrongPassword adds a strong password validation rule
func (svb *StructValidationBuilder) StrongPassword(field string, message ...string) *StructValidationBuilder {
	rule := StructValidationRule{
		Field: field,
		Tags:  []string{"strongpassword"},
	}
	if len(message) > 0 {
		rule.Message = message[0]
	}
	svb.rules = append(svb.rules, rule)
	return svb
}

// Build builds the struct validation schema
func (svb *StructValidationBuilder) Build() StructValidationSchema {
	return StructValidationSchema{
		Rules: svb.rules,
	}
}
