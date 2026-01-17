// Package auth provides authentication and authorization utilities for USC platform services.
package auth

import (
	"regexp"
	"strings"
	"unicode"
)

// AuthValidationError represents an authentication validation error
type AuthValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// Error implements the error interface
func (e *AuthValidationError) Error() string {
	return e.Message
}

// ValidationResult represents the result of validation
type ValidationResult struct {
	Valid    bool                  `json:"valid"`
	Errors   []AuthValidationError `json:"errors"`
	Warnings []AuthValidationError `json:"warnings"`
}

// AuthValidator provides authentication validation utilities
type AuthValidator struct {
	emailRegex    *regexp.Regexp
	usernameRegex *regexp.Regexp
	passwordRegex *regexp.Regexp
}

// NewAuthValidator creates a new authentication validator
func NewAuthValidator() *AuthValidator {
	return &AuthValidator{
		emailRegex:    regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
		usernameRegex: regexp.MustCompile(`^[a-zA-Z0-9_-]{3,20}$`),
		passwordRegex: regexp.MustCompile(`^.{8,}$`),
	}
}

// ValidateEmail validates an email address
func (av *AuthValidator) ValidateEmail(email string) *ValidationResult {
	result := &ValidationResult{Valid: true, Errors: []AuthValidationError{}, Warnings: []AuthValidationError{}}

	if email == "" {
		result.Valid = false
		result.Errors = append(result.Errors, AuthValidationError{
			Field:   "email",
			Message: "Email is required",
			Code:    "EMAIL_REQUIRED",
		})
		return result
	}

	if len(email) > 254 {
		result.Valid = false
		result.Errors = append(result.Errors, AuthValidationError{
			Field:   "email",
			Message: "Email is too long",
			Code:    "EMAIL_TOO_LONG",
		})
		return result
	}

	if !av.emailRegex.MatchString(email) {
		result.Valid = false
		result.Errors = append(result.Errors, AuthValidationError{
			Field:   "email",
			Message: "Invalid email format",
			Code:    "EMAIL_INVALID_FORMAT",
		})
		return result
	}

	// Check for common email issues
	if strings.Contains(email, "..") {
		result.Warnings = append(result.Warnings, AuthValidationError{
			Field:   "email",
			Message: "Email contains consecutive dots",
			Code:    "EMAIL_CONSECUTIVE_DOTS",
		})
	}

	return result
}

// ValidateUsername validates a username
func (av *AuthValidator) ValidateUsername(username string) *ValidationResult {
	result := &ValidationResult{Valid: true, Errors: []AuthValidationError{}, Warnings: []AuthValidationError{}}

	if username == "" {
		result.Valid = false
		result.Errors = append(result.Errors, AuthValidationError{
			Field:   "username",
			Message: "Username is required",
			Code:    "USERNAME_REQUIRED",
		})
		return result
	}

	if len(username) < 3 {
		result.Valid = false
		result.Errors = append(result.Errors, AuthValidationError{
			Field:   "username",
			Message: "Username must be at least 3 characters long",
			Code:    "USERNAME_TOO_SHORT",
		})
	}

	if len(username) > 20 {
		result.Valid = false
		result.Errors = append(result.Errors, AuthValidationError{
			Field:   "username",
			Message: "Username must be at most 20 characters long",
			Code:    "USERNAME_TOO_LONG",
		})
	}

	if !av.usernameRegex.MatchString(username) {
		result.Valid = false
		result.Errors = append(result.Errors, AuthValidationError{
			Field:   "username",
			Message: "Username can only contain letters, numbers, underscores, and hyphens",
			Code:    "USERNAME_INVALID_CHARACTERS",
		})
	}

	// Check for reserved usernames
	reservedUsernames := []string{"admin", "root", "user", "guest", "test", "api", "www", "mail", "support"}
	for _, reserved := range reservedUsernames {
		if strings.EqualFold(username, reserved) {
			result.Warnings = append(result.Warnings, AuthValidationError{
				Field:   "username",
				Message: "Username is reserved",
				Code:    "USERNAME_RESERVED",
			})
		}
	}

	return result
}

// ValidatePassword validates a password
func (av *AuthValidator) ValidatePassword(password string) *ValidationResult {
	result := &ValidationResult{Valid: true, Errors: []AuthValidationError{}, Warnings: []AuthValidationError{}}

	if password == "" {
		result.Valid = false
		result.Errors = append(result.Errors, AuthValidationError{
			Field:   "password",
			Message: "Password is required",
			Code:    "PASSWORD_REQUIRED",
		})
		return result
	}

	if len(password) < 8 {
		result.Valid = false
		result.Errors = append(result.Errors, AuthValidationError{
			Field:   "password",
			Message: "Password must be at least 8 characters long",
			Code:    "PASSWORD_TOO_SHORT",
		})
	}

	if len(password) > 128 {
		result.Valid = false
		result.Errors = append(result.Errors, AuthValidationError{
			Field:   "password",
			Message: "Password must be at most 128 characters long",
			Code:    "PASSWORD_TOO_LONG",
		})
	}

	// Check password strength
	strength := av.calculatePasswordStrength(password)
	if strength < 3 {
		result.Warnings = append(result.Warnings, AuthValidationError{
			Field:   "password",
			Message: "Password is weak",
			Code:    "PASSWORD_WEAK",
		})
	}

	// Check for common passwords
	if av.isCommonPassword(password) {
		result.Warnings = append(result.Warnings, AuthValidationError{
			Field:   "password",
			Message: "Password is commonly used",
			Code:    "PASSWORD_COMMON",
		})
	}

	return result
}

// ValidateRole validates a role
func (av *AuthValidator) ValidateRole(role string) *ValidationResult {
	result := &ValidationResult{Valid: true, Errors: []AuthValidationError{}, Warnings: []AuthValidationError{}}

	if role == "" {
		result.Valid = false
		result.Errors = append(result.Errors, AuthValidationError{
			Field:   "role",
			Message: "Role is required",
			Code:    "ROLE_REQUIRED",
		})
		return result
	}

	validRoles := []string{"admin", "moderator", "user", "guest"}
	valid := false
	for _, validRole := range validRoles {
		if role == validRole {
			valid = true
			break
		}
	}

	if !valid {
		result.Valid = false
		result.Errors = append(result.Errors, AuthValidationError{
			Field:   "role",
			Message: "Invalid role",
			Code:    "ROLE_INVALID",
		})
	}

	return result
}

// ValidatePermission validates a permission
func (av *AuthValidator) ValidatePermission(permission string) *ValidationResult {
	result := &ValidationResult{Valid: true, Errors: []AuthValidationError{}, Warnings: []AuthValidationError{}}

	if permission == "" {
		result.Valid = false
		result.Errors = append(result.Errors, AuthValidationError{
			Field:   "permission",
			Message: "Permission is required",
			Code:    "PERMISSION_REQUIRED",
		})
		return result
	}

	// Check permission format (resource:action)
	parts := strings.Split(permission, ":")
	if len(parts) != 2 {
		result.Valid = false
		result.Errors = append(result.Errors, AuthValidationError{
			Field:   "permission",
			Message: "Permission must be in format 'resource:action'",
			Code:    "PERMISSION_INVALID_FORMAT",
		})
		return result
	}

	resource := parts[0]
	action := parts[1]

	if resource == "" || action == "" {
		result.Valid = false
		result.Errors = append(result.Errors, AuthValidationError{
			Field:   "permission",
			Message: "Permission resource and action cannot be empty",
			Code:    "PERMISSION_EMPTY_PARTS",
		})
	}

	return result
}

// ValidateToken validates a JWT token format
func (av *AuthValidator) ValidateToken(token string) *ValidationResult {
	result := &ValidationResult{Valid: true, Errors: []AuthValidationError{}, Warnings: []AuthValidationError{}}

	if token == "" {
		result.Valid = false
		result.Errors = append(result.Errors, AuthValidationError{
			Field:   "token",
			Message: "Token is required",
			Code:    "TOKEN_REQUIRED",
		})
		return result
	}

	// Check if token has 3 parts (header.payload.signature)
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		result.Valid = false
		result.Errors = append(result.Errors, AuthValidationError{
			Field:   "token",
			Message: "Invalid token format",
			Code:    "TOKEN_INVALID_FORMAT",
		})
		return result
	}

	// Check if each part is not empty
	for i, part := range parts {
		if part == "" {
			result.Valid = false
			result.Errors = append(result.Errors, AuthValidationError{
				Field:   "token",
				Message: "Token part is empty",
				Code:    "TOKEN_EMPTY_PART",
			})
			return result
		}

		// Check if part is valid base64
		if !av.isValidBase64(part) {
			result.Valid = false
			result.Errors = append(result.Errors, AuthValidationError{
				Field:   "token",
				Message: "Invalid token encoding",
				Code:    "TOKEN_INVALID_ENCODING",
			})
			return result
		}

		// Add part index to error for debugging
		if i == 0 {
			result.Errors[len(result.Errors)-1].Field = "token.header"
		} else if i == 1 {
			result.Errors[len(result.Errors)-1].Field = "token.payload"
		} else {
			result.Errors[len(result.Errors)-1].Field = "token.signature"
		}
	}

	return result
}

// calculatePasswordStrength calculates password strength (0-5)
func (av *AuthValidator) calculatePasswordStrength(password string) int {
	strength := 0

	// Length check
	if len(password) >= 8 {
		strength++
	}
	if len(password) >= 12 {
		strength++
	}

	// Character type checks
	hasLower := false
	hasUpper := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if hasLower {
		strength++
	}
	if hasUpper {
		strength++
	}
	if hasDigit {
		strength++
	}
	if hasSpecial {
		strength++
	}

	return strength
}

// isCommonPassword checks if password is commonly used
func (av *AuthValidator) isCommonPassword(password string) bool {
	commonPasswords := []string{
		"password", "123456", "123456789", "qwerty", "abc123",
		"password123", "admin", "letmein", "welcome", "monkey",
		"1234567890", "password1", "qwerty123", "dragon", "master",
	}

	passwordLower := strings.ToLower(password)
	for _, common := range commonPasswords {
		if passwordLower == common {
			return true
		}
	}

	return false
}

// isValidBase64 checks if string is valid base64
func (av *AuthValidator) isValidBase64(s string) bool {
	// Basic base64 validation
	base64Regex := regexp.MustCompile(`^[A-Za-z0-9+/]*={0,2}$`)
	return base64Regex.MatchString(s)
}

// ValidateUserInput validates user input for authentication
func (av *AuthValidator) ValidateUserInput(email, username, password, role string) *ValidationResult {
	result := &ValidationResult{Valid: true, Errors: []AuthValidationError{}, Warnings: []AuthValidationError{}}

	// Validate email
	emailResult := av.ValidateEmail(email)
	if !emailResult.Valid {
		result.Valid = false
		result.Errors = append(result.Errors, emailResult.Errors...)
	}
	result.Warnings = append(result.Warnings, emailResult.Warnings...)

	// Validate username
	usernameResult := av.ValidateUsername(username)
	if !usernameResult.Valid {
		result.Valid = false
		result.Errors = append(result.Errors, usernameResult.Errors...)
	}
	result.Warnings = append(result.Warnings, usernameResult.Warnings...)

	// Validate password
	passwordResult := av.ValidatePassword(password)
	if !passwordResult.Valid {
		result.Valid = false
		result.Errors = append(result.Errors, passwordResult.Errors...)
	}
	result.Warnings = append(result.Warnings, passwordResult.Warnings...)

	// Validate role
	roleResult := av.ValidateRole(role)
	if !roleResult.Valid {
		result.Valid = false
		result.Errors = append(result.Errors, roleResult.Errors...)
	}
	result.Warnings = append(result.Warnings, roleResult.Warnings...)

	return result
}

// ValidateLoginInput validates login input
func (av *AuthValidator) ValidateLoginInput(identifier, password string) *ValidationResult {
	result := &ValidationResult{Valid: true, Errors: []AuthValidationError{}, Warnings: []AuthValidationError{}}

	if identifier == "" {
		result.Valid = false
		result.Errors = append(result.Errors, AuthValidationError{
			Field:   "identifier",
			Message: "Email or username is required",
			Code:    "IDENTIFIER_REQUIRED",
		})
	}

	if password == "" {
		result.Valid = false
		result.Errors = append(result.Errors, AuthValidationError{
			Field:   "password",
			Message: "Password is required",
			Code:    "PASSWORD_REQUIRED",
		})
	}

	return result
}

// ValidateTokenInput validates token input
func (av *AuthValidator) ValidateTokenInput(token string) *ValidationResult {
	return av.ValidateToken(token)
}
