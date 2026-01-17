// Package errors provides error handling utilities for USC platform services.
package errors

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// ErrorCode represents a standardized error code
type ErrorCode string

const (
	// General errors
	ErrCodeInternal           ErrorCode = "INTERNAL_ERROR"
	ErrCodeInvalidInput       ErrorCode = "INVALID_INPUT"
	ErrCodeNotFound           ErrorCode = "NOT_FOUND"
	ErrCodeUnauthorized       ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden          ErrorCode = "FORBIDDEN"
	ErrCodeConflict           ErrorCode = "CONFLICT"
	ErrCodeTimeout            ErrorCode = "TIMEOUT"
	ErrCodeRateLimited        ErrorCode = "RATE_LIMITED"
	ErrCodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"

	// Database errors
	ErrCodeDatabaseConnection  ErrorCode = "DATABASE_CONNECTION_ERROR"
	ErrCodeDatabaseQuery       ErrorCode = "DATABASE_QUERY_ERROR"
	ErrCodeDatabaseTransaction ErrorCode = "DATABASE_TRANSACTION_ERROR"
	ErrCodeDatabaseConstraint  ErrorCode = "DATABASE_CONSTRAINT_ERROR"

	// Authentication errors
	ErrCodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
	ErrCodeTokenExpired       ErrorCode = "TOKEN_EXPIRED"
	ErrCodeTokenInvalid       ErrorCode = "TOKEN_INVALID"
	ErrCodeAccountLocked      ErrorCode = "ACCOUNT_LOCKED"
	ErrCodeAccountDisabled    ErrorCode = "ACCOUNT_DISABLED"

	// Validation errors
	ErrCodeValidationFailed ErrorCode = "VALIDATION_FAILED"
	ErrCodeRequiredField    ErrorCode = "REQUIRED_FIELD"
	ErrCodeInvalidFormat    ErrorCode = "INVALID_FORMAT"
	ErrCodeOutOfRange       ErrorCode = "OUT_OF_RANGE"
	ErrCodeDuplicateValue   ErrorCode = "DUPLICATE_VALUE"

	// Business logic errors
	ErrCodeInsufficientFunds   ErrorCode = "INSUFFICIENT_FUNDS"
	ErrCodeResourceExhausted   ErrorCode = "RESOURCE_EXHAUSTED"
	ErrCodeOperationNotAllowed ErrorCode = "OPERATION_NOT_ALLOWED"
	ErrCodeQuotaExceeded       ErrorCode = "QUOTA_EXCEEDED"

	// Blockchain errors
	ErrCodeBlockchainConnection ErrorCode = "BLOCKCHAIN_CONNECTION_ERROR"
	ErrCodeTransactionFailed    ErrorCode = "TRANSACTION_FAILED"
	ErrCodeInvalidAddress       ErrorCode = "INVALID_ADDRESS"
	ErrCodeContractError        ErrorCode = "CONTRACT_ERROR"

	// File errors
	ErrCodeFileNotFound       ErrorCode = "FILE_NOT_FOUND"
	ErrCodeFileUploadFailed   ErrorCode = "FILE_UPLOAD_FAILED"
	ErrCodeFileDownloadFailed ErrorCode = "FILE_DOWNLOAD_FAILED"
	ErrCodeFileSizeExceeded   ErrorCode = "FILE_SIZE_EXCEEDED"
	ErrCodeInvalidFileType    ErrorCode = "INVALID_FILE_TYPE"

	// Network errors
	ErrCodeNetworkError      ErrorCode = "NETWORK_ERROR"
	ErrCodeConnectionTimeout ErrorCode = "CONNECTION_TIMEOUT"
	ErrCodeDNSResolution     ErrorCode = "DNS_RESOLUTION_ERROR"
	ErrCodeSSLHandshake      ErrorCode = "SSL_HANDSHAKE_ERROR"

	// Additional database errors
	ErrCodeDatabaseMigration ErrorCode = "DATABASE_MIGRATION_ERROR"

	// Additional authentication errors
	ErrCodeInvalidToken   ErrorCode = "INVALID_TOKEN"
	ErrCodeInvalidRequest ErrorCode = "INVALID_REQUEST"
)

// ErrorSeverity represents the severity level of an error
type ErrorSeverity string

const (
	SeverityLow      ErrorSeverity = "LOW"
	SeverityMedium   ErrorSeverity = "MEDIUM"
	SeverityHigh     ErrorSeverity = "HIGH"
	SeverityCritical ErrorSeverity = "CRITICAL"
)

// ErrorCategory represents the category of an error
type ErrorCategory string

const (
	CategorySystem     ErrorCategory = "SYSTEM"
	CategoryDatabase   ErrorCategory = "DATABASE"
	CategoryNetwork    ErrorCategory = "NETWORK"
	CategoryAuth       ErrorCategory = "AUTHENTICATION"
	CategoryValidation ErrorCategory = "VALIDATION"
	CategoryBusiness   ErrorCategory = "BUSINESS"
	CategoryExternal   ErrorCategory = "EXTERNAL"

	// USC-specific categories
	CategoryUSCBlockchain ErrorCategory = "usc_blockchain"
	CategoryUSCWallet     ErrorCategory = "usc_wallet"
	CategoryUSCNFT        ErrorCategory = "usc_nft"
	CategoryUSCStaking    ErrorCategory = "usc_staking"
	CategoryUSCRewards    ErrorCategory = "usc_rewards"

	// Service-specific categories
	CategorySocial      ErrorCategory = "social"
	CategoryVideo       ErrorCategory = "video"
	CategoryCommerce    ErrorCategory = "commerce"
	CategoryAI          ErrorCategory = "ai"
	CategorySearch      ErrorCategory = "search"
	CategoryAnalytics   ErrorCategory = "analytics"
	CategoryModeration  ErrorCategory = "moderation"
	CategoryAdvertising ErrorCategory = "advertising"
	CategoryMessaging   ErrorCategory = "messaging"
	CategoryGateway     ErrorCategory = "gateway"
	CategoryAdmin       ErrorCategory = "admin"
)

// DomainError represents a domain-specific error
type DomainError struct {
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	Details    string                 `json:"details,omitempty"`
	Severity   ErrorSeverity          `json:"severity"`
	Category   ErrorCategory          `json:"category"`
	Timestamp  time.Time              `json:"timestamp"`
	StackTrace string                 `json:"stack_trace,omitempty"`
	Context    map[string]interface{} `json:"context,omitempty"`
	Retryable  bool                   `json:"retryable"`
	UserID     string                 `json:"user_id,omitempty"`
	RequestID  string                 `json:"request_id,omitempty"`
	Service    string                 `json:"service,omitempty"`
}

// Error implements the error interface
func (e *DomainError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s - %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewDomainError creates a new domain error
func NewDomainError(code ErrorCode, message string, options ...ErrorOption) *DomainError {
	err := &DomainError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		Severity:  SeverityMedium,
		Category:  CategorySystem,
		Retryable: false,
	}

	// Apply options
	for _, option := range options {
		option(err)
	}

	// Capture stack trace for high severity errors
	if err.Severity == SeverityHigh || err.Severity == SeverityCritical {
		err.StackTrace = captureStackTrace()
	}

	return err
}

// ErrorOption represents an option for configuring domain errors
type ErrorOption func(*DomainError)

// WithDetails adds details to the error
func WithDetails(details string) ErrorOption {
	return func(e *DomainError) {
		e.Details = details
	}
}

// WithSeverity sets the error severity
func WithSeverity(severity ErrorSeverity) ErrorOption {
	return func(e *DomainError) {
		e.Severity = severity
	}
}

// WithCategory sets the error category
func WithCategory(category ErrorCategory) ErrorOption {
	return func(e *DomainError) {
		e.Category = category
	}
}

// WithContext adds context to the error
func WithContext(context map[string]interface{}) ErrorOption {
	return func(e *DomainError) {
		e.Context = context
	}
}

// WithRetryable sets whether the error is retryable
func WithRetryable(retryable bool) ErrorOption {
	return func(e *DomainError) {
		e.Retryable = retryable
	}
}

// WithUserID sets the user ID associated with the error
func WithUserID(userID string) ErrorOption {
	return func(e *DomainError) {
		e.UserID = userID
	}
}

// WithRequestID sets the request ID associated with the error
func WithRequestID(requestID string) ErrorOption {
	return func(e *DomainError) {
		e.RequestID = requestID
	}
}

// WithService sets the service name
func WithService(service string) ErrorOption {
	return func(e *DomainError) {
		e.Service = service
	}
}

// WithStackTrace adds a stack trace to the error
func WithStackTrace() ErrorOption {
	return func(e *DomainError) {
		e.StackTrace = captureStackTrace()
	}
}

// captureStackTrace captures the current stack trace
func captureStackTrace() string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

// IsRetryable checks if the error is retryable
func (e *DomainError) IsRetryable() bool {
	return e.Retryable
}

// IsCritical checks if the error is critical
func (e *DomainError) IsCritical() bool {
	return e.Severity == SeverityCritical
}

// IsHighSeverity checks if the error is high severity
func (e *DomainError) IsHighSeverity() bool {
	return e.Severity == SeverityHigh || e.Severity == SeverityCritical
}

// GetContextValue retrieves a value from the error context
func (e *DomainError) GetContextValue(key string) (interface{}, bool) {
	if e.Context == nil {
		return nil, false
	}
	value, exists := e.Context[key]
	return value, exists
}

// SetContextValue sets a value in the error context
func (e *DomainError) SetContextValue(key string, value interface{}) {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
}

// Wrap wraps an existing error with additional context
func (e *DomainError) Wrap(err error) *DomainError {
	if err == nil {
		return e
	}

	// If the wrapped error is already a DomainError, merge contexts
	if domainErr, ok := err.(*DomainError); ok {
		// Merge contexts
		if e.Context == nil {
			e.Context = make(map[string]interface{})
		}
		for key, value := range domainErr.Context {
			e.Context[key] = value
		}

		// Use the wrapped error's details if not set
		if e.Details == "" {
			e.Details = domainErr.Details
		}
	} else {
		// Wrap a standard error
		if e.Details == "" {
			e.Details = err.Error()
		} else {
			e.Details = fmt.Sprintf("%s: %s", e.Details, err.Error())
		}
	}

	return e
}

// ToMap converts the error to a map representation
func (e *DomainError) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"code":      e.Code,
		"message":   e.Message,
		"severity":  e.Severity,
		"category":  e.Category,
		"timestamp": e.Timestamp,
		"retryable": e.Retryable,
	}

	if e.Details != "" {
		result["details"] = e.Details
	}

	if e.Context != nil {
		result["context"] = e.Context
	}

	if e.UserID != "" {
		result["user_id"] = e.UserID
	}

	if e.RequestID != "" {
		result["request_id"] = e.RequestID
	}

	if e.Service != "" {
		result["service"] = e.Service
	}

	if e.StackTrace != "" {
		result["stack_trace"] = e.StackTrace
	}

	return result
}

// Predefined error constructors

// NewInternalError creates an internal server error
func NewInternalError(message string, options ...ErrorOption) *DomainError {
	return NewDomainError(ErrCodeInternal, message, append(options, WithSeverity(SeverityHigh), WithCategory(CategorySystem))...)
}

// NewInvalidInputError creates an invalid input error
func NewInvalidInputError(message string, options ...ErrorOption) *DomainError {
	return NewDomainError(ErrCodeInvalidInput, message, append(options, WithSeverity(SeverityMedium), WithCategory(CategoryValidation))...)
}

// NewNotFoundError creates a not found error
func NewNotFoundError(message string, options ...ErrorOption) *DomainError {
	return NewDomainError(ErrCodeNotFound, message, append(options, WithSeverity(SeverityMedium), WithCategory(CategoryBusiness))...)
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(message string, options ...ErrorOption) *DomainError {
	return NewDomainError(ErrCodeUnauthorized, message, append(options, WithSeverity(SeverityMedium), WithCategory(CategoryAuth))...)
}

// NewForbiddenError creates a forbidden error
func NewForbiddenError(message string, options ...ErrorOption) *DomainError {
	return NewDomainError(ErrCodeForbidden, message, append(options, WithSeverity(SeverityMedium), WithCategory(CategoryAuth))...)
}

// NewConflictError creates a conflict error
func NewConflictError(message string, options ...ErrorOption) *DomainError {
	return NewDomainError(ErrCodeConflict, message, append(options, WithSeverity(SeverityMedium), WithCategory(CategoryBusiness))...)
}

// NewTimeoutError creates a timeout error
func NewTimeoutError(message string, options ...ErrorOption) *DomainError {
	return NewDomainError(ErrCodeTimeout, message, append(options, WithSeverity(SeverityMedium), WithCategory(CategorySystem), WithRetryable(true))...)
}

// NewDatabaseError creates a database error
func NewDatabaseError(message string, options ...ErrorOption) *DomainError {
	return NewDomainError(ErrCodeDatabaseQuery, message, append(options, WithSeverity(SeverityHigh), WithCategory(CategoryDatabase), WithRetryable(true))...)
}

// NewValidationError creates a validation error
func NewValidationError(message string, options ...ErrorOption) *DomainError {
	return NewDomainError(ErrCodeValidationFailed, message, append(options, WithSeverity(SeverityLow), WithCategory(CategoryValidation))...)
}

// NewBusinessError creates a business logic error
func NewBusinessError(message string, options ...ErrorOption) *DomainError {
	return NewDomainError(ErrCodeOperationNotAllowed, message, append(options, WithSeverity(SeverityMedium), WithCategory(CategoryBusiness))...)
}

// Error utilities

// IsDomainError checks if an error is a domain error
func IsDomainError(err error) bool {
	_, ok := err.(*DomainError)
	return ok
}

// GetErrorCode extracts the error code from an error
func GetErrorCode(err error) ErrorCode {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.Code
	}
	return ErrCodeInternal
}

// GetErrorSeverity extracts the error severity from an error
func GetErrorSeverity(err error) ErrorSeverity {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.Severity
	}
	return SeverityMedium
}

// GetErrorCategory extracts the error category from an error
func GetErrorCategory(err error) ErrorCategory {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.Category
	}
	return CategorySystem
}

// IsRetryableError checks if an error is retryable
func IsRetryableError(err error) bool {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.IsRetryable()
	}
	return false
}

// IsCriticalError checks if an error is critical
func IsCriticalError(err error) bool {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.IsCritical()
	}
	return false
}

// ErrorAggregator aggregates multiple errors
type ErrorAggregator struct {
	errors []*DomainError
}

// NewErrorAggregator creates a new error aggregator
func NewErrorAggregator() *ErrorAggregator {
	return &ErrorAggregator{
		errors: make([]*DomainError, 0),
	}
}

// Add adds an error to the aggregator
func (ea *ErrorAggregator) Add(err error) {
	if err == nil {
		return
	}

	if domainErr, ok := err.(*DomainError); ok {
		ea.errors = append(ea.errors, domainErr)
	} else {
		ea.errors = append(ea.errors, NewInternalError(err.Error()))
	}
}

// HasErrors checks if there are any errors
func (ea *ErrorAggregator) HasErrors() bool {
	return len(ea.errors) > 0
}

// GetErrors returns all errors
func (ea *ErrorAggregator) GetErrors() []*DomainError {
	return ea.errors
}

// GetFirstError returns the first error
func (ea *ErrorAggregator) GetFirstError() *DomainError {
	if len(ea.errors) > 0 {
		return ea.errors[0]
	}
	return nil
}

// GetLastError returns the last error
func (ea *ErrorAggregator) GetLastError() *DomainError {
	if len(ea.errors) > 0 {
		return ea.errors[len(ea.errors)-1]
	}
	return nil
}

// GetCriticalErrors returns all critical errors
func (ea *ErrorAggregator) GetCriticalErrors() []*DomainError {
	critical := make([]*DomainError, 0)
	for _, err := range ea.errors {
		if err.IsCritical() {
			critical = append(critical, err)
		}
	}
	return critical
}

// GetRetryableErrors returns all retryable errors
func (ea *ErrorAggregator) GetRetryableErrors() []*DomainError {
	retryable := make([]*DomainError, 0)
	for _, err := range ea.errors {
		if err.IsRetryable() {
			retryable = append(retryable, err)
		}
	}
	return retryable
}

// ToError converts the aggregator to a single error
func (ea *ErrorAggregator) ToError() error {
	if !ea.HasErrors() {
		return nil
	}

	if len(ea.errors) == 1 {
		return ea.errors[0]
	}

	// Create a summary error
	messages := make([]string, len(ea.errors))
	for i, err := range ea.errors {
		messages[i] = err.Error()
	}

	return NewInternalError(fmt.Sprintf("Multiple errors occurred: %s", strings.Join(messages, "; ")))
}
