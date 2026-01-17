// Package validation provides validation utilities for USC platform services.
package validation

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// BusinessValidator provides business logic validation
type BusinessValidator struct {
	validator *Validator
}

// NewBusinessValidator creates a new business validator
func NewBusinessValidator() *BusinessValidator {
	return &BusinessValidator{
		validator: NewValidator(),
	}
}

// ValidateUser validates user business rules
func (bv *BusinessValidator) ValidateUser(user interface{}) ValidationErrors {
	var errors ValidationErrors

	// Validate basic user fields
	basicErrors := bv.validator.Validate(user)
	errors = append(errors, basicErrors...)

	// Add business-specific validations
	// These would be custom business rules specific to your application

	return errors
}

// ValidateOrder validates order business rules
func (bv *BusinessValidator) ValidateOrder(order interface{}) ValidationErrors {
	var errors ValidationErrors

	// Validate basic order fields
	basicErrors := bv.validator.Validate(order)
	errors = append(errors, basicErrors...)

	// Add business-specific validations
	// These would be custom business rules specific to your application

	return errors
}

// ValidateProduct validates product business rules
func (bv *BusinessValidator) ValidateProduct(product interface{}) ValidationErrors {
	var errors ValidationErrors

	// Validate basic product fields
	basicErrors := bv.validator.Validate(product)
	errors = append(errors, basicErrors...)

	// Add business-specific validations
	// These would be custom business rules specific to your application

	return errors
}

// ValidateContent validates content business rules
func (bv *BusinessValidator) ValidateContent(content interface{}) ValidationErrors {
	var errors ValidationErrors

	// Validate basic content fields
	basicErrors := bv.validator.Validate(content)
	errors = append(errors, basicErrors...)

	// Add business-specific validations
	// These would be custom business rules specific to your application

	return errors
}

// BusinessValidationRule represents a business validation rule
type BusinessValidationRule struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Condition   func(interface{}) bool `json:"-"`
	Message     string                 `json:"message"`
	Severity    string                 `json:"severity"`
}

// BusinessValidationSchema represents a business validation schema
type BusinessValidationSchema struct {
	Rules []BusinessValidationRule `json:"rules"`
}

// ValidateWithBusinessRules validates data against business rules
func (bv *BusinessValidator) ValidateWithBusinessRules(data interface{}, schema BusinessValidationSchema) ValidationErrors {
	var errors ValidationErrors

	for _, rule := range schema.Rules {
		if !rule.Condition(data) {
			errors = append(errors, ValidationError{
				Field:   rule.Name,
				Tag:     "business_rule",
				Value:   data,
				Message: rule.Message,
			})
		}
	}

	return errors
}

// Common business validation rules

// ValidateEmailDomain validates email domain
func (bv *BusinessValidator) ValidateEmailDomain(email string, allowedDomains []string) ValidationErrors {
	var errors ValidationErrors

	if email == "" {
		errors = append(errors, ValidationError{
			Field:   "email",
			Tag:     "required",
			Value:   email,
			Message: "Email is required",
		})
		return errors
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		errors = append(errors, ValidationError{
			Field:   "email",
			Tag:     "invalid_format",
			Value:   email,
			Message: "Invalid email format",
		})
		return errors
	}

	domain := parts[1]
	allowed := false
	for _, allowedDomain := range allowedDomains {
		if domain == allowedDomain {
			allowed = true
			break
		}
	}

	if !allowed {
		errors = append(errors, ValidationError{
			Field:   "email",
			Tag:     "domain_not_allowed",
			Value:   email,
			Message: fmt.Sprintf("Email domain %s is not allowed", domain),
		})
	}

	return errors
}

// ValidatePasswordStrength validates password strength
func (bv *BusinessValidator) ValidatePasswordStrength(password string, minLength int) ValidationErrors {
	var errors ValidationErrors

	if password == "" {
		errors = append(errors, ValidationError{
			Field:   "password",
			Tag:     "required",
			Value:   password,
			Message: "Password is required",
		})
		return errors
	}

	if len(password) < minLength {
		errors = append(errors, ValidationError{
			Field:   "password",
			Tag:     "min_length",
			Value:   password,
			Message: fmt.Sprintf("Password must be at least %d characters long", minLength),
		})
	}

	// Check for uppercase letter
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		errors = append(errors, ValidationError{
			Field:   "password",
			Tag:     "uppercase_required",
			Value:   password,
			Message: "Password must contain at least one uppercase letter",
		})
	}

	// Check for lowercase letter
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		errors = append(errors, ValidationError{
			Field:   "password",
			Tag:     "lowercase_required",
			Value:   password,
			Message: "Password must contain at least one lowercase letter",
		})
	}

	// Check for digit
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		errors = append(errors, ValidationError{
			Field:   "password",
			Tag:     "digit_required",
			Value:   password,
			Message: "Password must contain at least one digit",
		})
	}

	// Check for special character
	if !regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password) {
		errors = append(errors, ValidationError{
			Field:   "password",
			Tag:     "special_char_required",
			Value:   password,
			Message: "Password must contain at least one special character",
		})
	}

	return errors
}

// ValidateAge validates age
func (bv *BusinessValidator) ValidateAge(birthDate time.Time, minAge, maxAge int) ValidationErrors {
	var errors ValidationErrors

	if birthDate.IsZero() {
		errors = append(errors, ValidationError{
			Field:   "birth_date",
			Tag:     "required",
			Value:   birthDate,
			Message: "Birth date is required",
		})
		return errors
	}

	age := time.Now().Year() - birthDate.Year()
	if time.Now().YearDay() < birthDate.YearDay() {
		age--
	}

	if age < minAge {
		errors = append(errors, ValidationError{
			Field:   "age",
			Tag:     "min_age",
			Value:   age,
			Message: fmt.Sprintf("Age must be at least %d years", minAge),
		})
	}

	if age > maxAge {
		errors = append(errors, ValidationError{
			Field:   "age",
			Tag:     "max_age",
			Value:   age,
			Message: fmt.Sprintf("Age must be at most %d years", maxAge),
		})
	}

	return errors
}

// ValidatePhoneNumber validates phone number
func (bv *BusinessValidator) ValidatePhoneNumber(phone string, countryCode string) ValidationErrors {
	var errors ValidationErrors

	if phone == "" {
		errors = append(errors, ValidationError{
			Field:   "phone",
			Tag:     "required",
			Value:   phone,
			Message: "Phone number is required",
		})
		return errors
	}

	// Remove all non-digit characters
	digits := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")

	// Validate based on country code
	switch countryCode {
	case "US":
		if len(digits) != 10 {
			errors = append(errors, ValidationError{
				Field:   "phone",
				Tag:     "invalid_length",
				Value:   phone,
				Message: "US phone number must be 10 digits",
			})
		}
	case "UK":
		if len(digits) != 11 {
			errors = append(errors, ValidationError{
				Field:   "phone",
				Tag:     "invalid_length",
				Value:   phone,
				Message: "UK phone number must be 11 digits",
			})
		}
	default:
		if len(digits) < 7 || len(digits) > 15 {
			errors = append(errors, ValidationError{
				Field:   "phone",
				Tag:     "invalid_length",
				Value:   phone,
				Message: "Phone number must be between 7 and 15 digits",
			})
		}
	}

	return errors
}

// ValidateCreditCard validates credit card number
func (bv *BusinessValidator) ValidateCreditCard(cardNumber string) ValidationErrors {
	var errors ValidationErrors

	if cardNumber == "" {
		errors = append(errors, ValidationError{
			Field:   "card_number",
			Tag:     "required",
			Value:   cardNumber,
			Message: "Credit card number is required",
		})
		return errors
	}

	// Remove all non-digit characters
	digits := regexp.MustCompile(`\D`).ReplaceAllString(cardNumber, "")

	// Check length
	if len(digits) < 13 || len(digits) > 19 {
		errors = append(errors, ValidationError{
			Field:   "card_number",
			Tag:     "invalid_length",
			Value:   cardNumber,
			Message: "Credit card number must be between 13 and 19 digits",
		})
		return errors
	}

	// Luhn algorithm validation
	if !bv.validateLuhn(digits) {
		errors = append(errors, ValidationError{
			Field:   "card_number",
			Tag:     "invalid_checksum",
			Value:   cardNumber,
			Message: "Invalid credit card number",
		})
	}

	return errors
}

// validateLuhn validates credit card number using Luhn algorithm
func (bv *BusinessValidator) validateLuhn(digits string) bool {
	sum := 0
	alternate := false

	// Process digits from right to left
	for i := len(digits) - 1; i >= 0; i-- {
		digit := int(digits[i] - '0')

		if alternate {
			digit *= 2
			if digit > 9 {
				digit = (digit % 10) + 1
			}
		}

		sum += digit
		alternate = !alternate
	}

	return sum%10 == 0
}

// ValidateDateRange validates date range
func (bv *BusinessValidator) ValidateDateRange(startDate, endDate time.Time) ValidationErrors {
	var errors ValidationErrors

	if startDate.IsZero() {
		errors = append(errors, ValidationError{
			Field:   "start_date",
			Tag:     "required",
			Value:   startDate,
			Message: "Start date is required",
		})
	}

	if endDate.IsZero() {
		errors = append(errors, ValidationError{
			Field:   "end_date",
			Tag:     "required",
			Value:   endDate,
			Message: "End date is required",
		})
	}

	if !startDate.IsZero() && !endDate.IsZero() {
		if startDate.After(endDate) {
			errors = append(errors, ValidationError{
				Field:   "date_range",
				Tag:     "invalid_range",
				Value:   fmt.Sprintf("%v to %v", startDate, endDate),
				Message: "Start date must be before end date",
			})
		}
	}

	return errors
}

// ValidateFileSize validates file size
func (bv *BusinessValidator) ValidateFileSize(fileSize int64, maxSize int64) ValidationErrors {
	var errors ValidationErrors

	if fileSize <= 0 {
		errors = append(errors, ValidationError{
			Field:   "file_size",
			Tag:     "invalid_size",
			Value:   fileSize,
			Message: "File size must be greater than 0",
		})
		return errors
	}

	if fileSize > maxSize {
		errors = append(errors, ValidationError{
			Field:   "file_size",
			Tag:     "max_size_exceeded",
			Value:   fileSize,
			Message: fmt.Sprintf("File size must be at most %d bytes", maxSize),
		})
	}

	return errors
}

// ValidateFileType validates file type
func (bv *BusinessValidator) ValidateFileType(fileName string, allowedTypes []string) ValidationErrors {
	var errors ValidationErrors

	if fileName == "" {
		errors = append(errors, ValidationError{
			Field:   "file_name",
			Tag:     "required",
			Value:   fileName,
			Message: "File name is required",
		})
		return errors
	}

	// Get file extension
	parts := strings.Split(fileName, ".")
	if len(parts) < 2 {
		errors = append(errors, ValidationError{
			Field:   "file_type",
			Tag:     "invalid_extension",
			Value:   fileName,
			Message: "File must have an extension",
		})
		return errors
	}

	extension := strings.ToLower(parts[len(parts)-1])
	allowed := false
	for _, allowedType := range allowedTypes {
		if extension == strings.ToLower(allowedType) {
			allowed = true
			break
		}
	}

	if !allowed {
		errors = append(errors, ValidationError{
			Field:   "file_type",
			Tag:     "type_not_allowed",
			Value:   fileName,
			Message: fmt.Sprintf("File type .%s is not allowed", extension),
		})
	}

	return errors
}
