// Package errors provides error handling utilities for USC platform services.
package errors

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// HTTPError represents an HTTP-specific error
type HTTPError struct {
	Code      int                    `json:"code"`
	Message   string                 `json:"message"`
	Details   string                 `json:"details,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Path      string                 `json:"path,omitempty"`
	Method    string                 `json:"method,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
	UserID    string                 `json:"user_id,omitempty"`
	Service   string                 `json:"service,omitempty"`
	Domain    *DomainError           `json:"domain,omitempty"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

// Error implements the error interface
func (e *HTTPError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("HTTP %d: %s - %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("HTTP %d: %s", e.Code, e.Message)
}

// NewHTTPError creates a new HTTP error
func NewHTTPError(code int, message string, details ...string) *HTTPError {
	err := &HTTPError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
	}

	if len(details) > 0 {
		err.Details = details[0]
	}

	return err
}

// NewHTTPErrorFromDomain creates an HTTP error from a domain error
func NewHTTPErrorFromDomain(domainErr *DomainError) *HTTPError {
	code := mapDomainErrorToHTTPCode(domainErr.Code)

	return &HTTPError{
		Code:      code,
		Message:   domainErr.Message,
		Details:   domainErr.Details,
		Timestamp: time.Now(),
		Domain:    domainErr,
		Context:   domainErr.Context,
		UserID:    domainErr.UserID,
		RequestID: domainErr.RequestID,
		Service:   domainErr.Service,
	}
}

// WriteResponse writes the error as an HTTP response
func (e *HTTPError) WriteResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Code)

	json.NewEncoder(w).Encode(e)
}

// mapDomainErrorToHTTPCode maps domain error codes to HTTP status codes
func mapDomainErrorToHTTPCode(errorCode ErrorCode) int {
	switch errorCode {
	case ErrCodeInternal:
		return http.StatusInternalServerError
	case ErrCodeInvalidInput:
		return http.StatusBadRequest
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeConflict:
		return http.StatusConflict
	case ErrCodeTimeout:
		return http.StatusRequestTimeout
	case ErrCodeRateLimited:
		return http.StatusTooManyRequests
	case ErrCodeServiceUnavailable:
		return http.StatusServiceUnavailable
	case ErrCodeDatabaseConnection:
		return http.StatusServiceUnavailable
	case ErrCodeDatabaseQuery:
		return http.StatusInternalServerError
	case ErrCodeDatabaseTransaction:
		return http.StatusInternalServerError
	case ErrCodeDatabaseConstraint:
		return http.StatusBadRequest
	case ErrCodeInvalidCredentials:
		return http.StatusUnauthorized
	case ErrCodeTokenExpired:
		return http.StatusUnauthorized
	case ErrCodeTokenInvalid:
		return http.StatusUnauthorized
	case ErrCodeAccountLocked:
		return http.StatusForbidden
	case ErrCodeAccountDisabled:
		return http.StatusForbidden
	case ErrCodeValidationFailed:
		return http.StatusBadRequest
	case ErrCodeRequiredField:
		return http.StatusBadRequest
	case ErrCodeInvalidFormat:
		return http.StatusBadRequest
	case ErrCodeOutOfRange:
		return http.StatusBadRequest
	case ErrCodeDuplicateValue:
		return http.StatusConflict
	case ErrCodeInsufficientFunds:
		return http.StatusBadRequest
	case ErrCodeResourceExhausted:
		return http.StatusTooManyRequests
	case ErrCodeOperationNotAllowed:
		return http.StatusForbidden
	case ErrCodeQuotaExceeded:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}

// mapHTTPCodeToDomainError maps HTTP status codes to domain error codes
func mapHTTPCodeToDomainError(httpCode int) ErrorCode {
	switch httpCode {
	case http.StatusBadRequest:
		return ErrCodeInvalidInput
	case http.StatusUnauthorized:
		return ErrCodeUnauthorized
	case http.StatusForbidden:
		return ErrCodeForbidden
	case http.StatusNotFound:
		return ErrCodeNotFound
	case http.StatusConflict:
		return ErrCodeConflict
	case http.StatusRequestTimeout:
		return ErrCodeTimeout
	case http.StatusTooManyRequests:
		return ErrCodeRateLimited
	case http.StatusInternalServerError:
		return ErrCodeInternal
	case http.StatusServiceUnavailable:
		return ErrCodeServiceUnavailable
	default:
		return ErrCodeInternal
	}
}

// ConvertToHTTPError converts any error to an HTTP error
func ConvertToHTTPError(err error) *HTTPError {
	if err == nil {
		return nil
	}

	// If it's already an HTTP error, return it
	if httpErr, ok := err.(*HTTPError); ok {
		return httpErr
	}

	// If it's a domain error, convert it
	if domainErr, ok := err.(*DomainError); ok {
		return NewHTTPErrorFromDomain(domainErr)
	}

	// Default to internal server error
	return NewHTTPError(http.StatusInternalServerError, err.Error())
}

// ConvertFromHTTPError converts an HTTP error to a domain error
func ConvertFromHTTPError(httpErr *HTTPError) *DomainError {
	if httpErr == nil {
		return nil
	}

	errorCode := mapHTTPCodeToDomainError(httpErr.Code)

	domainErr := NewDomainError(errorCode, httpErr.Message)

	if httpErr.Details != "" {
		domainErr.Details = httpErr.Details
	}

	if httpErr.Context != nil {
		domainErr.Context = httpErr.Context
	}

	domainErr.UserID = httpErr.UserID
	domainErr.RequestID = httpErr.RequestID
	domainErr.Service = httpErr.Service

	return domainErr
}

// HTTPErrorHandler handles HTTP errors
type HTTPErrorHandler struct {
	service string
}

// NewHTTPErrorHandler creates a new HTTP error handler
func NewHTTPErrorHandler(service string) *HTTPErrorHandler {
	return &HTTPErrorHandler{
		service: service,
	}
}

// HandleError handles an error and returns an HTTP error
func (h *HTTPErrorHandler) HandleError(ctx context.Context, err error, r *http.Request) *HTTPError {
	if err == nil {
		return nil
	}

	httpErr := ConvertToHTTPError(err)

	// Add request context
	if r != nil {
		httpErr.Path = r.URL.Path
		httpErr.Method = r.Method
	}

	// Add service context
	httpErr.Service = h.service

	return httpErr
}

// HandleDomainError handles a domain error and returns an HTTP error
func (h *HTTPErrorHandler) HandleDomainError(ctx context.Context, domainErr *DomainError, r *http.Request) *HTTPError {
	if domainErr == nil {
		return nil
	}

	// Add service context
	domainErr.Service = h.service

	httpErr := NewHTTPErrorFromDomain(domainErr)

	// Add request context
	if r != nil {
		httpErr.Path = r.URL.Path
		httpErr.Method = r.Method
	}

	return httpErr
}

// HandlePanic handles a panic and returns an HTTP error
func (h *HTTPErrorHandler) HandlePanic(ctx context.Context, panic interface{}, r *http.Request) *HTTPError {
	message := "Internal server error"
	if panic != nil {
		message = fmt.Sprintf("Panic: %v", panic)
	}

	domainErr := NewInternalError(message, WithService(h.service), WithStackTrace())
	httpErr := NewHTTPErrorFromDomain(domainErr)

	// Add request context
	if r != nil {
		httpErr.Path = r.URL.Path
		httpErr.Method = r.Method
	}

	return httpErr
}

// HTTPErrorMiddleware provides HTTP error middleware
type HTTPErrorMiddleware struct {
	handler *HTTPErrorHandler
}

// NewHTTPErrorMiddleware creates a new HTTP error middleware
func NewHTTPErrorMiddleware(service string) *HTTPErrorMiddleware {
	return &HTTPErrorMiddleware{
		handler: NewHTTPErrorHandler(service),
	}
}

// ErrorHandler returns an HTTP error handler middleware
func (em *HTTPErrorMiddleware) ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if panic := recover(); panic != nil {
				// Handle panic
				httpErr := em.handler.HandlePanic(r.Context(), panic, r)
				httpErr.WriteResponse(w)
			}
		}()

		// Create a custom response writer to capture errors
		errorWriter := &ErrorResponseWriter{
			ResponseWriter: w,
			request:        r,
		}

		next.ServeHTTP(errorWriter, r)

		// Check if an error occurred
		if errorWriter.error != nil {
			httpErr := em.handler.HandleError(r.Context(), errorWriter.error, r)
			httpErr.WriteResponse(w)
		}
	})
}

// ErrorResponseWriter captures errors during request processing
type ErrorResponseWriter struct {
	http.ResponseWriter
	request *http.Request
	error   error
}

// WriteHeader captures the status code
func (erw *ErrorResponseWriter) WriteHeader(code int) {
	// Only write header if it's not an error status
	if code < 400 {
		erw.ResponseWriter.WriteHeader(code)
	}
}

// Write captures the response body
func (erw *ErrorResponseWriter) Write(data []byte) (int, error) {
	// Only write if no error occurred
	if erw.error == nil {
		return erw.ResponseWriter.Write(data)
	}
	return 0, nil
}

// SetError sets an error to be handled
func (erw *ErrorResponseWriter) SetError(err error) {
	erw.error = err
}

// HTTPErrorLogger logs HTTP errors
type HTTPErrorLogger struct {
	service string
}

// NewHTTPErrorLogger creates a new HTTP error logger
func NewHTTPErrorLogger(service string) *HTTPErrorLogger {
	return &HTTPErrorLogger{
		service: service,
	}
}

// LogError logs an HTTP error
func (el *HTTPErrorLogger) LogError(ctx context.Context, err error, r *http.Request) {
	if err == nil {
		return
	}

	httpErr := ConvertToHTTPError(err)

	// Add request context
	if r != nil {
		httpErr.Path = r.URL.Path
		httpErr.Method = r.Method
	}

	// Log based on status code
	switch {
	case httpErr.Code >= 500:
		// Log server errors
		el.logServerError(ctx, httpErr)
	case httpErr.Code >= 400:
		// Log client errors
		el.logClientError(ctx, httpErr)
	default:
		// Log other errors
		el.logOtherError(ctx, httpErr)
	}
}

// logServerError logs server errors
func (el *HTTPErrorLogger) logServerError(ctx context.Context, err *HTTPError) {
	// In a real implementation, you would use a proper logger
	fmt.Printf("SERVER ERROR [%s] %s %s: %s\n", el.service, err.Method, err.Path, err.Error())
}

// logClientError logs client errors
func (el *HTTPErrorLogger) logClientError(ctx context.Context, err *HTTPError) {
	// In a real implementation, you would use a proper logger
	fmt.Printf("CLIENT ERROR [%s] %s %s: %s\n", el.service, err.Method, err.Path, err.Error())
}

// logOtherError logs other errors
func (el *HTTPErrorLogger) logOtherError(ctx context.Context, err *HTTPError) {
	// In a real implementation, you would use a proper logger
	fmt.Printf("ERROR [%s] %s %s: %s\n", el.service, err.Method, err.Path, err.Error())
}

// HTTPErrorMetrics collects HTTP error metrics
type HTTPErrorMetrics struct {
	service string
}

// NewHTTPErrorMetrics creates a new HTTP error metrics collector
func NewHTTPErrorMetrics(service string) *HTTPErrorMetrics {
	return &HTTPErrorMetrics{
		service: service,
	}
}

// RecordError records an error metric
func (em *HTTPErrorMetrics) RecordError(ctx context.Context, err error, r *http.Request) {
	if err == nil {
		return
	}

	httpErr := ConvertToHTTPError(err)

	// Record metrics based on status code
	// In a real implementation, you would use a metrics library like Prometheus
	fmt.Printf("METRIC: HTTP error recorded [%s] %s %s: %d\n", em.service, httpErr.Method, httpErr.Path, httpErr.Code)
}

// HTTPErrorResponse represents a standardized HTTP error response
type HTTPErrorResponse struct {
	Error HTTPError `json:"error"`
}

// WriteErrorResponse writes a standardized error response
func WriteErrorResponse(w http.ResponseWriter, err error, r *http.Request) {
	httpErr := ConvertToHTTPError(err)

	// Add request context
	if r != nil {
		httpErr.Path = r.URL.Path
		httpErr.Method = r.Method
	}

	response := HTTPErrorResponse{
		Error: *httpErr,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpErr.Code)
	json.NewEncoder(w).Encode(response)
}

// WriteDomainErrorResponse writes a domain error as an HTTP response
func WriteDomainErrorResponse(w http.ResponseWriter, domainErr *DomainError, r *http.Request) {
	httpErr := NewHTTPErrorFromDomain(domainErr)

	// Add request context
	if r != nil {
		httpErr.Path = r.URL.Path
		httpErr.Method = r.Method
	}

	response := HTTPErrorResponse{
		Error: *httpErr,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpErr.Code)
	json.NewEncoder(w).Encode(response)
}
