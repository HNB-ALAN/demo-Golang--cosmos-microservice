package errors

import (
	"fmt"
	"time"

	"github.com/usc-platform/shared/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// USC Blockchain Core Service specific error codes
const (
	// Service-specific error codes (1000-1999)
	ErrServiceUnavailable = 1000 + iota
	ErrInvalidRequest
	ErrResourceNotFound
	ErrResourceAlreadyExists
	ErrPermissionDenied
	ErrRateLimitExceeded
	ErrValidationFailed
	ErrInternalError
)

// Error messages for USC Blockchain Core Service
var errorMessages = map[int]string{
	ErrServiceUnavailable:    "USC Blockchain Core Service is currently unavailable",
	ErrInvalidRequest:        "Invalid request format or parameters",
	ErrResourceNotFound:      "Requested resource not found",
	ErrResourceAlreadyExists: "Resource already exists",
	ErrPermissionDenied:      "Permission denied for this operation",
	ErrRateLimitExceeded:     "Rate limit exceeded, please try again later",
	ErrValidationFailed:      "Request validation failed",
	ErrInternalError:         "Internal server error occurred",
}

// ErrorManager manages errors for USC Blockchain Core Service
type ErrorManager struct {
	serviceName string
}

// NewErrorManager creates a new error manager
func NewErrorManager() *ErrorManager {
	return &ErrorManager{
		serviceName: "{{SERVICE_NAME}}",
	}
}

// NewServiceError creates a new service-specific error
func (em *ErrorManager) NewServiceError(code int, details ...interface{}) error {
	message, exists := errorMessages[code]
	if !exists {
		message = "Unknown error"
	}

	if len(details) > 0 {
		message = fmt.Sprintf("%s: %v", message, details[0])
	}

	return &errors.DomainError{
		Code:      errors.ErrorCode(fmt.Sprintf("SERVICE_%d", code)),
		Message:   message,
		Details:   fmt.Sprintf("%v", details),
		Severity:  errors.SeverityMedium,
		Category:  errors.CategorySystem,
		Timestamp: time.Now(),
		Context:   map[string]interface{}{"service": em.serviceName},
		Retryable: false,
	}
}

// NewValidationError creates a validation error
func (em *ErrorManager) NewValidationError(field string, reason string) error {
	return em.NewServiceError(ErrValidationFailed, fmt.Sprintf("field '%s': %s", field, reason))
}

// NewNotFoundError creates a not found error
func (em *ErrorManager) NewNotFoundError(resource string) error {
	return em.NewServiceError(ErrResourceNotFound, fmt.Sprintf("resource '%s' not found", resource))
}

// NewAlreadyExistsError creates an already exists error
func (em *ErrorManager) NewAlreadyExistsError(resource string) error {
	return em.NewServiceError(ErrResourceAlreadyExists, fmt.Sprintf("resource '%s' already exists", resource))
}

// NewPermissionDeniedError creates a permission denied error
func (em *ErrorManager) NewPermissionDeniedError(operation string) error {
	return em.NewServiceError(ErrPermissionDenied, fmt.Sprintf("permission denied for operation '%s'", operation))
}

// NewRateLimitError creates a rate limit error
func (em *ErrorManager) NewRateLimitError(limit int, window string) error {
	return em.NewServiceError(ErrRateLimitExceeded, fmt.Sprintf("rate limit of %d requests per %s exceeded", limit, window))
}

// ToGRPCError converts service error to gRPC error
func (em *ErrorManager) ToGRPCError(err error) error {
	if err == nil {
		return nil
	}

	// Check if it's already a gRPC error
	if _, ok := status.FromError(err); ok {
		return err
	}

	// Check if it's a domain error
	if serviceErr, ok := err.(*errors.DomainError); ok {
		// Convert domain error to gRPC error based on error code
		switch string(serviceErr.Code) {
		case "SERVICE_1001", "SERVICE_1006": // ErrInvalidRequest, ErrValidationFailed
			return status.Error(codes.InvalidArgument, serviceErr.Message)
		case "SERVICE_1002": // ErrResourceNotFound
			return status.Error(codes.NotFound, serviceErr.Message)
		case "SERVICE_1003": // ErrResourceAlreadyExists
			return status.Error(codes.AlreadyExists, serviceErr.Message)
		case "SERVICE_1004": // ErrPermissionDenied
			return status.Error(codes.PermissionDenied, serviceErr.Message)
		case "SERVICE_1005": // ErrRateLimitExceeded
			return status.Error(codes.ResourceExhausted, serviceErr.Message)
		case "SERVICE_1000": // ErrServiceUnavailable
			return status.Error(codes.Unavailable, serviceErr.Message)
		default:
			return status.Error(codes.Internal, serviceErr.Message)
		}
	}

	// Default to internal error
	return status.Error(codes.Internal, "Internal server error")
}

// FromGRPCError converts gRPC error to service error
func (em *ErrorManager) FromGRPCError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return em.NewServiceError(ErrInternalError, err.Error())
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return em.NewServiceError(ErrInvalidRequest, st.Message())
	case codes.NotFound:
		return em.NewServiceError(ErrResourceNotFound, st.Message())
	case codes.AlreadyExists:
		return em.NewServiceError(ErrResourceAlreadyExists, st.Message())
	case codes.PermissionDenied:
		return em.NewServiceError(ErrPermissionDenied, st.Message())
	case codes.ResourceExhausted:
		return em.NewServiceError(ErrRateLimitExceeded, st.Message())
	case codes.Unavailable:
		return em.NewServiceError(ErrServiceUnavailable, st.Message())
	default:
		return em.NewServiceError(ErrInternalError, st.Message())
	}
}

// LogError logs an error with service context
func (em *ErrorManager) LogError(err error, context map[string]interface{}) {
	if serviceErr, ok := err.(*errors.DomainError); ok {
		// Log structured error with service context
		errorContext := map[string]interface{}{
			"service":   em.serviceName,
			"code":      string(serviceErr.Code),
			"message":   serviceErr.Message,
			"details":   serviceErr.Details,
			"severity":  string(serviceErr.Severity),
			"category":  string(serviceErr.Category),
			"timestamp": serviceErr.Timestamp,
			"retryable": serviceErr.Retryable,
		}

		// Merge additional context if provided
		for k, v := range context {
			errorContext[k] = v
		}

		// Log the error (this would integrate with actual logging system)
		_ = errorContext
	} else {
		// Log generic error
		errorContext := map[string]interface{}{
			"service": em.serviceName,
			"error":   err.Error(),
		}

		// Merge additional context if provided
		for k, v := range context {
			errorContext[k] = v
		}

		// Log the error (this would integrate with actual logging system)
		_ = errorContext
	}
}
