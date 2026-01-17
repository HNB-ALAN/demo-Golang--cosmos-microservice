// Package validation provides validation utilities for USC platform services.
package validation

import (
	"fmt"
	"strings"
)

// ValidationErrorType represents the type of validation error
type ValidationErrorType string

const (
	// ErrorTypeRequired represents a required field error
	ErrorTypeRequired ValidationErrorType = "required"
	// ErrorTypeInvalid represents an invalid value error
	ErrorTypeInvalid ValidationErrorType = "invalid"
	// ErrorTypeTooShort represents a too short error
	ErrorTypeTooShort ValidationErrorType = "too_short"
	// ErrorTypeTooLong represents a too long error
	ErrorTypeTooLong ValidationErrorType = "too_long"
	// ErrorTypeTooSmall represents a too small error
	ErrorTypeTooSmall ValidationErrorType = "too_small"
	// ErrorTypeTooLarge represents a too large error
	ErrorTypeTooLarge ValidationErrorType = "too_large"
	// ErrorTypeFormat represents a format error
	ErrorTypeFormat ValidationErrorType = "format"
	// ErrorTypeBusinessRule represents a business rule error
	ErrorTypeBusinessRule ValidationErrorType = "business_rule"
	// ErrorTypeCustom represents a custom error
	ErrorTypeCustom ValidationErrorType = "custom"
)

// ValidationErrorSeverity represents the severity of a validation error
type ValidationErrorSeverity string

const (
	// SeverityError represents an error severity
	SeverityError ValidationErrorSeverity = "error"
	// SeverityWarning represents a warning severity
	SeverityWarning ValidationErrorSeverity = "warning"
	// SeverityInfo represents an info severity
	SeverityInfo ValidationErrorSeverity = "info"
)

// ValidationError represents a validation error with additional metadata
type ValidationError struct {
	Field    string                  `json:"field"`
	Tag      string                  `json:"tag"`
	Value    interface{}             `json:"value"`
	Message  string                  `json:"message"`
	Type     ValidationErrorType     `json:"type"`
	Severity ValidationErrorSeverity `json:"severity"`
	Code     string                  `json:"code,omitempty"`
	Context  map[string]interface{}  `json:"context,omitempty"`
}

// Error implements the error interface
func (e ValidationError) Error() string {
	return e.Message
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

// Error implements the error interface
func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return ""
	}

	var messages []string
	for _, err := range e {
		messages = append(messages, err.Message)
	}

	return strings.Join(messages, "; ")
}

// HasErrors returns true if there are validation errors
func (e ValidationErrors) HasErrors() bool {
	return len(e) > 0
}

// HasWarnings returns true if there are validation warnings
func (e ValidationErrors) HasWarnings() bool {
	for _, err := range e {
		if err.Severity == SeverityWarning {
			return true
		}
	}
	return false
}

// GetErrors returns only error severity validation errors
func (e ValidationErrors) GetErrors() ValidationErrors {
	var errors ValidationErrors
	for _, err := range e {
		if err.Severity == SeverityError {
			errors = append(errors, err)
		}
	}
	return errors
}

// GetWarnings returns only warning severity validation errors
func (e ValidationErrors) GetWarnings() ValidationErrors {
	var warnings ValidationErrors
	for _, err := range e {
		if err.Severity == SeverityWarning {
			warnings = append(warnings, err)
		}
	}
	return warnings
}

// GetInfos returns only info severity validation errors
func (e ValidationErrors) GetInfos() ValidationErrors {
	var infos ValidationErrors
	for _, err := range e {
		if err.Severity == SeverityInfo {
			infos = append(infos, err)
		}
	}
	return infos
}

// GetFieldErrors returns errors for a specific field
func (e ValidationErrors) GetFieldErrors(field string) ValidationErrors {
	var errors ValidationErrors
	for _, err := range e {
		if err.Field == field {
			errors = append(errors, err)
		}
	}
	return errors
}

// GetErrorsByType returns errors of a specific type
func (e ValidationErrors) GetErrorsByType(errorType ValidationErrorType) ValidationErrors {
	var errors ValidationErrors
	for _, err := range e {
		if err.Type == errorType {
			errors = append(errors, err)
		}
	}
	return errors
}

// GetFirstError returns the first validation error
func (e ValidationErrors) GetFirstError() *ValidationError {
	if len(e) > 0 {
		return &e[0]
	}
	return nil
}

// GetFirstErrorByField returns the first error for a specific field
func (e ValidationErrors) GetFirstErrorByField(field string) *ValidationError {
	for _, err := range e {
		if err.Field == field {
			return &err
		}
	}
	return nil
}

// GetFirstErrorByType returns the first error of a specific type
func (e ValidationErrors) GetFirstErrorByType(errorType ValidationErrorType) *ValidationError {
	for _, err := range e {
		if err.Type == errorType {
			return &err
		}
	}
	return nil
}

// Add adds a validation error
func (e *ValidationErrors) Add(err ValidationError) {
	*e = append(*e, err)
}

// AddError adds an error severity validation error
func (e *ValidationErrors) AddError(field, tag, message string, value interface{}) {
	e.Add(ValidationError{
		Field:    field,
		Tag:      tag,
		Value:    value,
		Message:  message,
		Type:     ErrorTypeInvalid,
		Severity: SeverityError,
	})
}

// AddWarning adds a warning severity validation error
func (e *ValidationErrors) AddWarning(field, tag, message string, value interface{}) {
	e.Add(ValidationError{
		Field:    field,
		Tag:      tag,
		Value:    value,
		Message:  message,
		Type:     ErrorTypeInvalid,
		Severity: SeverityWarning,
	})
}

// AddInfo adds an info severity validation error
func (e *ValidationErrors) AddInfo(field, tag, message string, value interface{}) {
	e.Add(ValidationError{
		Field:    field,
		Tag:      tag,
		Value:    value,
		Message:  message,
		Type:     ErrorTypeInvalid,
		Severity: SeverityInfo,
	})
}

// Clear clears all validation errors
func (e *ValidationErrors) Clear() {
	*e = ValidationErrors{}
}

// Count returns the number of validation errors
func (e ValidationErrors) Count() int {
	return len(e)
}

// CountByField returns the number of errors for a specific field
func (e ValidationErrors) CountByField(field string) int {
	count := 0
	for _, err := range e {
		if err.Field == field {
			count++
		}
	}
	return count
}

// CountByType returns the number of errors of a specific type
func (e ValidationErrors) CountByType(errorType ValidationErrorType) int {
	count := 0
	for _, err := range e {
		if err.Type == errorType {
			count++
		}
	}
	return count
}

// CountBySeverity returns the number of errors of a specific severity
func (e ValidationErrors) CountBySeverity(severity ValidationErrorSeverity) int {
	count := 0
	for _, err := range e {
		if err.Severity == severity {
			count++
		}
	}
	return count
}

// ToMap converts validation errors to a map
func (e ValidationErrors) ToMap() map[string]interface{} {
	result := make(map[string]interface{})

	// Group errors by field
	fieldErrors := make(map[string][]ValidationError)
	for _, err := range e {
		fieldErrors[err.Field] = append(fieldErrors[err.Field], err)
	}

	result["errors"] = fieldErrors
	result["count"] = len(e)
	result["has_errors"] = e.HasErrors()
	result["has_warnings"] = e.HasWarnings()

	return result
}

// ToJSON converts validation errors to JSON string
func (e ValidationErrors) ToJSON() (string, error) {
	// In a real implementation, you would use json.Marshal
	// For now, we'll return a simple string representation
	return fmt.Sprintf("ValidationErrors(count=%d)", len(e)), nil
}

// ValidationErrorBuilder provides a builder for validation errors
type ValidationErrorBuilder struct {
	error ValidationError
}

// NewValidationErrorBuilder creates a new validation error builder
func NewValidationErrorBuilder() *ValidationErrorBuilder {
	return &ValidationErrorBuilder{
		error: ValidationError{
			Type:     ErrorTypeInvalid,
			Severity: SeverityError,
			Context:  make(map[string]interface{}),
		},
	}
}

// Field sets the field name
func (veb *ValidationErrorBuilder) Field(field string) *ValidationErrorBuilder {
	veb.error.Field = field
	return veb
}

// Tag sets the validation tag
func (veb *ValidationErrorBuilder) Tag(tag string) *ValidationErrorBuilder {
	veb.error.Tag = tag
	return veb
}

// Value sets the field value
func (veb *ValidationErrorBuilder) Value(value interface{}) *ValidationErrorBuilder {
	veb.error.Value = value
	return veb
}

// Message sets the error message
func (veb *ValidationErrorBuilder) Message(message string) *ValidationErrorBuilder {
	veb.error.Message = message
	return veb
}

// Type sets the error type
func (veb *ValidationErrorBuilder) Type(errorType ValidationErrorType) *ValidationErrorBuilder {
	veb.error.Type = errorType
	return veb
}

// Severity sets the error severity
func (veb *ValidationErrorBuilder) Severity(severity ValidationErrorSeverity) *ValidationErrorBuilder {
	veb.error.Severity = severity
	return veb
}

// Code sets the error code
func (veb *ValidationErrorBuilder) Code(code string) *ValidationErrorBuilder {
	veb.error.Code = code
	return veb
}

// Context sets the error context
func (veb *ValidationErrorBuilder) Context(context map[string]interface{}) *ValidationErrorBuilder {
	veb.error.Context = context
	return veb
}

// AddContext adds a context value
func (veb *ValidationErrorBuilder) AddContext(key string, value interface{}) *ValidationErrorBuilder {
	if veb.error.Context == nil {
		veb.error.Context = make(map[string]interface{})
	}
	veb.error.Context[key] = value
	return veb
}

// Build builds the validation error
func (veb *ValidationErrorBuilder) Build() ValidationError {
	return veb.error
}

// ValidationErrorFormatter provides formatting for validation errors
type ValidationErrorFormatter struct {
	format string
}

// NewValidationErrorFormatter creates a new validation error formatter
func NewValidationErrorFormatter(format string) *ValidationErrorFormatter {
	return &ValidationErrorFormatter{
		format: format,
	}
}

// Format formats a validation error
func (vef *ValidationErrorFormatter) Format(err ValidationError) string {
	switch vef.format {
	case "json":
		return fmt.Sprintf(`{"field":"%s","tag":"%s","message":"%s"}`, err.Field, err.Tag, err.Message)
	case "xml":
		return fmt.Sprintf(`<error field="%s" tag="%s">%s</error>`, err.Field, err.Tag, err.Message)
	case "yaml":
		return fmt.Sprintf("field: %s\ntag: %s\nmessage: %s", err.Field, err.Tag, err.Message)
	default:
		return fmt.Sprintf("%s: %s", err.Field, err.Message)
	}
}

// FormatErrors formats multiple validation errors
func (vef *ValidationErrorFormatter) FormatErrors(errors ValidationErrors) string {
	var formatted []string
	for _, err := range errors {
		formatted = append(formatted, vef.Format(err))
	}
	return strings.Join(formatted, "\n")
}

// ValidationErrorTranslator provides translation for validation errors
type ValidationErrorTranslator struct {
	translations map[string]map[string]string
}

// NewValidationErrorTranslator creates a new validation error translator
func NewValidationErrorTranslator() *ValidationErrorTranslator {
	return &ValidationErrorTranslator{
		translations: make(map[string]map[string]string),
	}
}

// AddTranslation adds a translation for a field and tag
func (vet *ValidationErrorTranslator) AddTranslation(field, tag, translation string) {
	if vet.translations[field] == nil {
		vet.translations[field] = make(map[string]string)
	}
	vet.translations[field][tag] = translation
}

// Translate translates a validation error
func (vet *ValidationErrorTranslator) Translate(err ValidationError) string {
	if fieldTranslations, exists := vet.translations[err.Field]; exists {
		if translation, exists := fieldTranslations[err.Tag]; exists {
			return translation
		}
	}
	return err.Message
}

// TranslateErrors translates multiple validation errors
func (vet *ValidationErrorTranslator) TranslateErrors(errors ValidationErrors) ValidationErrors {
	var translated ValidationErrors
	for _, err := range errors {
		translatedErr := err
		translatedErr.Message = vet.Translate(err)
		translated = append(translated, translatedErr)
	}
	return translated
}
