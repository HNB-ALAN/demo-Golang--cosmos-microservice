package validation

import (
	"testing"
)

func TestNewValidator(t *testing.T) {
	validator := NewValidator()

	if validator == nil {
		t.Fatal("Expected validator to be created")
	}
}

func TestValidator_Validate(t *testing.T) {
	validator := NewValidator()

	// Test valid struct
	validStruct := struct {
		Name  string `validate:"required"`
		Email string `validate:"required,email"`
		Age   int    `validate:"min=18"`
	}{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   25,
	}

	errors := validator.Validate(validStruct)
	if len(errors) > 0 {
		t.Errorf("Expected no errors for valid struct, got %v", errors)
	}
}

func TestValidator_ValidateWithInvalidStruct(t *testing.T) {
	validator := NewValidator()

	// Test invalid struct
	invalidStruct := struct {
		Name  string `validate:"required"`
		Email string `validate:"required,email"`
		Age   int    `validate:"min=18"`
	}{
		Name:  "",              // Invalid: required field is empty
		Email: "invalid-email", // Invalid: not a valid email
		Age:   15,              // Invalid: below minimum age
	}

	errors := validator.Validate(invalidStruct)
	if len(errors) == 0 {
		t.Error("Expected errors for invalid struct, got none")
	}
}

func TestValidator_ValidateWithComplexStruct(t *testing.T) {
	validator := NewValidator()

	// Test with complex struct
	type Address struct {
		Street string `validate:"required"`
		City   string `validate:"required"`
		Zip    string `validate:"required,numeric"`
	}

	type User struct {
		Name    string  `validate:"required,min=2,max=50"`
		Email   string  `validate:"required,email"`
		Age     int     `validate:"min=18,max=120"`
		Address Address `validate:"required"`
	}

	validUser := User{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   25,
		Address: Address{
			Street: "123 Main St",
			City:   "New York",
			Zip:    "10001",
		},
	}

	errors := validator.Validate(validUser)
	if len(errors) > 0 {
		t.Errorf("Expected no errors for valid complex struct, got %v", errors)
	}

	// Test with invalid complex struct
	invalidUser := User{
		Name:  "",              // Invalid: required field is empty
		Email: "invalid-email", // Invalid: not a valid email
		Age:   15,              // Invalid: below minimum age
		Address: Address{
			Street: "",    // Invalid: required field is empty
			City:   "",    // Invalid: required field is empty
			Zip:    "abc", // Invalid: not numeric
		},
	}

	errors = validator.Validate(invalidUser)
	if len(errors) == 0 {
		t.Error("Expected errors for invalid complex struct, got none")
	}
}

func TestValidator_ValidateWithPointer(t *testing.T) {
	validator := NewValidator()

	// Test with pointer to valid struct
	validStruct := &struct {
		Name  string `validate:"required"`
		Email string `validate:"required,email"`
	}{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	errors := validator.Validate(validStruct)
	if len(errors) > 0 {
		t.Errorf("Expected no errors for valid pointer struct, got %v", errors)
	}

	// Test with nil pointer
	var nilStruct *struct {
		Name string `validate:"required"`
	}

	errors = validator.Validate(nilStruct)
	if len(errors) > 0 {
		t.Errorf("Expected no errors for nil pointer, got %v", errors)
	}
}

func TestValidator_ValidateWithNonStruct(t *testing.T) {
	validator := NewValidator()

	// Test with non-struct value
	errors := validator.Validate("not a struct")
	if len(errors) > 0 {
		t.Errorf("Expected no errors for non-struct value, got %v", errors)
	}
}

func TestValidator_ValidateWithEmptyStruct(t *testing.T) {
	validator := NewValidator()

	// Test with empty struct
	emptyStruct := struct{}{}
	errors := validator.Validate(emptyStruct)
	if len(errors) > 0 {
		t.Errorf("Expected no errors for empty struct, got %v", errors)
	}
}

func TestValidator_ValidateWithUnexportedFields(t *testing.T) {
	validator := NewValidator()

	// Test with struct containing unexported fields
	unexportedStruct := struct {
		Name  string `validate:"required"`
		email string `validate:"required,email"` // unexported field
		Age   int    `validate:"min=18"`
	}{
		Name:  "John Doe",
		email: "john@example.com",
		Age:   25,
	}

	errors := validator.Validate(unexportedStruct)
	if len(errors) > 0 {
		t.Errorf("Expected no errors for struct with unexported fields, got %v", errors)
	}
}

func TestValidator_ValidateWithCustomTag(t *testing.T) {
	validator := NewValidatorWithTag("custom")

	// Test with custom tag
	customStruct := struct {
		Name string `custom:"required"`
	}{
		Name: "John Doe",
	}

	errors := validator.Validate(customStruct)
	if len(errors) > 0 {
		t.Errorf("Expected no errors for custom tag, got %v", errors)
	}

	// Test with invalid custom tag
	invalidCustomStruct := struct {
		Name string `custom:"required"`
	}{
		Name: "", // Invalid: required field is empty
	}

	errors = validator.Validate(invalidCustomStruct)
	if len(errors) == 0 {
		t.Error("Expected errors for invalid custom tag, got none")
	}
}

// Benchmark tests
func BenchmarkValidator_Validate(b *testing.B) {
	validator := NewValidator()
	validStruct := struct {
		Name  string `validate:"required"`
		Email string `validate:"required,email"`
		Age   int    `validate:"min=18"`
	}{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   25,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.Validate(validStruct)
	}
}

func BenchmarkValidator_ValidateComplex(b *testing.B) {
	validator := NewValidator()

	type Address struct {
		Street string `validate:"required"`
		City   string `validate:"required"`
		Zip    string `validate:"required,numeric"`
	}

	type User struct {
		Name    string  `validate:"required,min=2,max=50"`
		Email   string  `validate:"required,email"`
		Age     int     `validate:"min=18,max=120"`
		Address Address `validate:"required"`
	}

	validUser := User{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   25,
		Address: Address{
			Street: "123 Main St",
			City:   "New York",
			Zip:    "10001",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.Validate(validUser)
	}
}
