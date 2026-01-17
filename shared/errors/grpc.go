// Package errors provides error handling utilities for USC platform services.
package errors

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCError represents a gRPC-specific error
type GRPCError struct {
	Code    codes.Code   `json:"code"`
	Message string       `json:"message"`
	Details string       `json:"details,omitempty"`
	Domain  *DomainError `json:"domain,omitempty"`
}

// Error implements the error interface
func (e *GRPCError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("gRPC %s: %s - %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("gRPC %s: %s", e.Code, e.Message)
}

// NewGRPCError creates a new gRPC error
func NewGRPCError(code codes.Code, message string, details ...string) *GRPCError {
	err := &GRPCError{
		Code:    code,
		Message: message,
	}

	if len(details) > 0 {
		err.Details = details[0]
	}

	return err
}

// NewGRPCErrorFromDomain creates a gRPC error from a domain error
func NewGRPCErrorFromDomain(domainErr *DomainError) *GRPCError {
	code := mapDomainErrorToGRPCCode(domainErr.Code)

	return &GRPCError{
		Code:    code,
		Message: domainErr.Message,
		Details: domainErr.Details,
		Domain:  domainErr,
	}
}

// ToStatus converts the gRPC error to a gRPC status
func (e *GRPCError) ToStatus() *status.Status {
	return status.New(e.Code, e.Message)
}

// mapDomainErrorToGRPCCode maps domain error codes to gRPC codes
func mapDomainErrorToGRPCCode(errorCode ErrorCode) codes.Code {
	switch errorCode {
	case ErrCodeInternal:
		return codes.Internal
	case ErrCodeInvalidInput:
		return codes.InvalidArgument
	case ErrCodeNotFound:
		return codes.NotFound
	case ErrCodeUnauthorized:
		return codes.Unauthenticated
	case ErrCodeForbidden:
		return codes.PermissionDenied
	case ErrCodeConflict:
		return codes.AlreadyExists
	case ErrCodeTimeout:
		return codes.DeadlineExceeded
	case ErrCodeRateLimited:
		return codes.ResourceExhausted
	case ErrCodeServiceUnavailable:
		return codes.Unavailable
	case ErrCodeDatabaseConnection:
		return codes.Unavailable
	case ErrCodeDatabaseQuery:
		return codes.Internal
	case ErrCodeDatabaseTransaction:
		return codes.Internal
	case ErrCodeDatabaseConstraint:
		return codes.InvalidArgument
	case ErrCodeInvalidCredentials:
		return codes.Unauthenticated
	case ErrCodeTokenExpired:
		return codes.Unauthenticated
	case ErrCodeTokenInvalid:
		return codes.Unauthenticated
	case ErrCodeAccountLocked:
		return codes.PermissionDenied
	case ErrCodeAccountDisabled:
		return codes.PermissionDenied
	case ErrCodeValidationFailed:
		return codes.InvalidArgument
	case ErrCodeRequiredField:
		return codes.InvalidArgument
	case ErrCodeInvalidFormat:
		return codes.InvalidArgument
	case ErrCodeOutOfRange:
		return codes.OutOfRange
	case ErrCodeDuplicateValue:
		return codes.AlreadyExists
	case ErrCodeInsufficientFunds:
		return codes.FailedPrecondition
	case ErrCodeResourceExhausted:
		return codes.ResourceExhausted
	case ErrCodeOperationNotAllowed:
		return codes.PermissionDenied
	case ErrCodeQuotaExceeded:
		return codes.ResourceExhausted
	case ErrCodeBlockchainConnection:
		return codes.Unavailable
	case ErrCodeTransactionFailed:
		return codes.Internal
	case ErrCodeInvalidAddress:
		return codes.InvalidArgument
	case ErrCodeContractError:
		return codes.Internal
	case ErrCodeFileNotFound:
		return codes.NotFound
	case ErrCodeFileUploadFailed:
		return codes.Internal
	case ErrCodeFileDownloadFailed:
		return codes.Internal
	case ErrCodeFileSizeExceeded:
		return codes.InvalidArgument
	case ErrCodeInvalidFileType:
		return codes.InvalidArgument
	case ErrCodeNetworkError:
		return codes.Unavailable
	case ErrCodeConnectionTimeout:
		return codes.DeadlineExceeded
	case ErrCodeDNSResolution:
		return codes.Unavailable
	case ErrCodeSSLHandshake:
		return codes.Unavailable
	case ErrCodeDatabaseMigration:
		return codes.Internal
	case ErrCodeInvalidToken:
		return codes.Unauthenticated
	case ErrCodeInvalidRequest:
		return codes.InvalidArgument
	default:
		return codes.Internal
	}
}

// mapGRPCCodeToDomainError maps gRPC codes to domain error codes
func mapGRPCCodeToDomainError(grpcCode codes.Code) ErrorCode {
	switch grpcCode {
	case codes.Internal:
		return ErrCodeInternal
	case codes.InvalidArgument:
		return ErrCodeInvalidInput
	case codes.NotFound:
		return ErrCodeNotFound
	case codes.Unauthenticated:
		return ErrCodeUnauthorized
	case codes.PermissionDenied:
		return ErrCodeForbidden
	case codes.AlreadyExists:
		return ErrCodeConflict
	case codes.DeadlineExceeded:
		return ErrCodeTimeout
	case codes.ResourceExhausted:
		return ErrCodeResourceExhausted
	case codes.Unavailable:
		return ErrCodeServiceUnavailable
	case codes.OutOfRange:
		return ErrCodeOutOfRange
	case codes.FailedPrecondition:
		return ErrCodeOperationNotAllowed
	default:
		return ErrCodeInternal
	}
}

// ConvertToGRPCError converts any error to a gRPC error
func ConvertToGRPCError(err error) *GRPCError {
	if err == nil {
		return nil
	}

	// If it's already a gRPC error, return it
	if grpcErr, ok := err.(*GRPCError); ok {
		return grpcErr
	}

	// If it's a domain error, convert it
	if domainErr, ok := err.(*DomainError); ok {
		return NewGRPCErrorFromDomain(domainErr)
	}

	// If it's a gRPC status error, extract information
	if st, ok := status.FromError(err); ok {
		return &GRPCError{
			Code:    st.Code(),
			Message: st.Message(),
			Details: st.Details()[0].(string),
		}
	}

	// Default to internal error
	return NewGRPCError(codes.Internal, err.Error())
}

// ConvertFromGRPCError converts a gRPC error to a domain error
func ConvertFromGRPCError(grpcErr *GRPCError) *DomainError {
	if grpcErr == nil {
		return nil
	}

	errorCode := mapGRPCCodeToDomainError(grpcErr.Code)

	domainErr := NewDomainError(errorCode, grpcErr.Message)

	if grpcErr.Details != "" {
		domainErr.Details = grpcErr.Details
	}

	if grpcErr.Domain != nil {
		domainErr.Context = grpcErr.Domain.Context
		domainErr.Severity = grpcErr.Domain.Severity
		domainErr.Category = grpcErr.Domain.Category
		domainErr.Retryable = grpcErr.Domain.Retryable
		domainErr.UserID = grpcErr.Domain.UserID
		domainErr.RequestID = grpcErr.Domain.RequestID
		domainErr.Service = grpcErr.Domain.Service
	}

	return domainErr
}

// GRPCErrorHandler handles gRPC errors
type GRPCErrorHandler struct {
	service string
}

// NewGRPCErrorHandler creates a new gRPC error handler
func NewGRPCErrorHandler(service string) *GRPCErrorHandler {
	return &GRPCErrorHandler{
		service: service,
	}
}

// HandleError handles an error and returns a gRPC status
func (h *GRPCErrorHandler) HandleError(ctx context.Context, err error) *status.Status {
	if err == nil {
		return status.New(codes.OK, "success")
	}

	grpcErr := ConvertToGRPCError(err)

	// Add service context
	if grpcErr.Domain != nil {
		grpcErr.Domain.Service = h.service
	}

	return grpcErr.ToStatus()
}

// HandleDomainError handles a domain error and returns a gRPC status
func (h *GRPCErrorHandler) HandleDomainError(ctx context.Context, domainErr *DomainError) *status.Status {
	if domainErr == nil {
		return status.New(codes.OK, "success")
	}

	// Add service context
	domainErr.Service = h.service

	grpcErr := NewGRPCErrorFromDomain(domainErr)
	return grpcErr.ToStatus()
}

// HandlePanic handles a panic and returns a gRPC status
func (h *GRPCErrorHandler) HandlePanic(ctx context.Context, panic interface{}) *status.Status {
	message := "Internal server error"
	if panic != nil {
		message = fmt.Sprintf("Panic: %v", panic)
	}

	domainErr := NewInternalError(message, WithService(h.service), WithStackTrace())
	grpcErr := NewGRPCErrorFromDomain(domainErr)

	return grpcErr.ToStatus()
}

// GRPCErrorInterceptor provides gRPC error interceptors
type GRPCErrorInterceptor struct {
	handler *GRPCErrorHandler
}

// NewGRPCErrorInterceptor creates a new gRPC error interceptor
func NewGRPCErrorInterceptor(service string) *GRPCErrorInterceptor {
	return &GRPCErrorInterceptor{
		handler: NewGRPCErrorHandler(service),
	}
}

// UnaryServerInterceptor returns a unary server interceptor that handles errors
func (ei *GRPCErrorInterceptor) UnaryServerInterceptor() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		defer func() {
			if r := recover(); r != nil {
				// Handle panic
				st := ei.handler.HandlePanic(ctx, r)
				panic(st.Err())
			}
		}()

		resp, err := handler(ctx, req)
		if err != nil {
			st := ei.handler.HandleError(ctx, err)
			return nil, st.Err()
		}

		return resp, nil
	}
}

// StreamServerInterceptor returns a stream server interceptor that handles errors
func (ei *GRPCErrorInterceptor) StreamServerInterceptor() func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		defer func() {
			if r := recover(); r != nil {
				// Handle panic
				st := ei.handler.HandlePanic(ss.Context(), r)
				panic(st.Err())
			}
		}()

		err := handler(srv, ss)
		if err != nil {
			st := ei.handler.HandleError(ss.Context(), err)
			return st.Err()
		}

		return nil
	}
}

// GRPCErrorLogger logs gRPC errors
type GRPCErrorLogger struct {
	service string
}

// NewGRPCErrorLogger creates a new gRPC error logger
func NewGRPCErrorLogger(service string) *GRPCErrorLogger {
	return &GRPCErrorLogger{
		service: service,
	}
}

// LogError logs a gRPC error
func (el *GRPCErrorLogger) LogError(ctx context.Context, err error, method string) {
	if err == nil {
		return
	}

	grpcErr := ConvertToGRPCError(err)

	// Log based on severity
	if grpcErr.Domain != nil {
		switch grpcErr.Domain.Severity {
		case SeverityCritical:
			// Log critical errors
			el.logCritical(ctx, grpcErr, method)
		case SeverityHigh:
			// Log high severity errors
			el.logHigh(ctx, grpcErr, method)
		case SeverityMedium:
			// Log medium severity errors
			el.logMedium(ctx, grpcErr, method)
		case SeverityLow:
			// Log low severity errors
			el.logLow(ctx, grpcErr, method)
		}
	} else {
		// Log unknown errors
		el.logUnknown(ctx, grpcErr, method)
	}
}

// logCritical logs critical errors
func (el *GRPCErrorLogger) logCritical(ctx context.Context, err *GRPCError, method string) {
	// In a real implementation, you would use a proper logger
	fmt.Printf("CRITICAL ERROR [%s] %s: %s\n", el.service, method, err.Error())
}

// logHigh logs high severity errors
func (el *GRPCErrorLogger) logHigh(ctx context.Context, err *GRPCError, method string) {
	// In a real implementation, you would use a proper logger
	fmt.Printf("HIGH ERROR [%s] %s: %s\n", el.service, method, err.Error())
}

// logMedium logs medium severity errors
func (el *GRPCErrorLogger) logMedium(ctx context.Context, err *GRPCError, method string) {
	// In a real implementation, you would use a proper logger
	fmt.Printf("MEDIUM ERROR [%s] %s: %s\n", el.service, method, err.Error())
}

// logLow logs low severity errors
func (el *GRPCErrorLogger) logLow(ctx context.Context, err *GRPCError, method string) {
	// In a real implementation, you would use a proper logger
	fmt.Printf("LOW ERROR [%s] %s: %s\n", el.service, method, err.Error())
}

// logUnknown logs unknown errors
func (el *GRPCErrorLogger) logUnknown(ctx context.Context, err *GRPCError, method string) {
	// In a real implementation, you would use a proper logger
	fmt.Printf("UNKNOWN ERROR [%s] %s: %s\n", el.service, method, err.Error())
}

// GRPCErrorMetrics collects gRPC error metrics
type GRPCErrorMetrics struct {
	service string
}

// NewGRPCErrorMetrics creates a new gRPC error metrics collector
func NewGRPCErrorMetrics(service string) *GRPCErrorMetrics {
	return &GRPCErrorMetrics{
		service: service,
	}
}

// RecordError records an error metric
func (em *GRPCErrorMetrics) RecordError(ctx context.Context, err error, method string) {
	if err == nil {
		return
	}

	grpcErr := ConvertToGRPCError(err)

	// Record metrics based on error code and severity
	// In a real implementation, you would use a metrics library like Prometheus
	fmt.Printf("METRIC: gRPC error recorded [%s] %s: %s\n", em.service, method, grpcErr.Code)
}
