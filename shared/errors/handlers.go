// Package errors provides error handling utilities for USC platform services.
package errors

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

// ErrorHandler provides comprehensive error handling capabilities
type ErrorHandler struct {
	service   string
	registry  *ErrorCodeRegistry
	logger    ErrorLogger
	metrics   ErrorMetrics
	reporter  ErrorReporter
	recoverer ErrorRecoverer
}

// ErrorLogger interface for logging errors
type ErrorLogger interface {
	LogError(ctx context.Context, err *DomainError)
	LogPanic(ctx context.Context, panic interface{}, stackTrace string)
}

// ErrorMetrics interface for collecting error metrics
type ErrorMetrics interface {
	RecordError(ctx context.Context, err *DomainError)
	RecordPanic(ctx context.Context, panic interface{})
}

// ErrorReporter interface for reporting errors to external services
type ErrorReporter interface {
	ReportError(ctx context.Context, err *DomainError)
	ReportPanic(ctx context.Context, panic interface{}, stackTrace string)
}

// ErrorRecoverer interface for error recovery strategies
type ErrorRecoverer interface {
	CanRecover(err *DomainError) bool
	Recover(ctx context.Context, err *DomainError) error
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(service string, options ...ErrorHandlerOption) *ErrorHandler {
	handler := &ErrorHandler{
		service:  service,
		registry: GetGlobalErrorCodeRegistry(),
	}

	// Apply options
	for _, option := range options {
		option(handler)
	}

	return handler
}

// ErrorHandlerOption represents an option for configuring error handlers
type ErrorHandlerOption func(*ErrorHandler)

// WithErrorLogger sets the error logger
func WithErrorLogger(logger ErrorLogger) ErrorHandlerOption {
	return func(h *ErrorHandler) {
		h.logger = logger
	}
}

// WithErrorMetrics sets the error metrics collector
func WithErrorMetrics(metrics ErrorMetrics) ErrorHandlerOption {
	return func(h *ErrorHandler) {
		h.metrics = metrics
	}
}

// WithErrorReporter sets the error reporter
func WithErrorReporter(reporter ErrorReporter) ErrorHandlerOption {
	return func(h *ErrorHandler) {
		h.reporter = reporter
	}
}

// WithErrorRecoverer sets the error recoverer
func WithErrorRecoverer(recoverer ErrorRecoverer) ErrorHandlerOption {
	return func(h *ErrorHandler) {
		h.recoverer = recoverer
	}
}

// HandleError handles an error with comprehensive processing
func (h *ErrorHandler) HandleError(ctx context.Context, err error) *DomainError {
	if err == nil {
		return nil
	}

	// Convert to domain error if needed
	var domainErr *DomainError
	if de, ok := err.(*DomainError); ok {
		domainErr = de
	} else {
		domainErr = NewInternalError(err.Error(), WithService(h.service))
	}

	// Add service context
	domainErr.Service = h.service

	// Validate error code (simplified - just check if code exists)
	if domainErr.Code == "" {
		// Log warning about empty error code
		if h.logger != nil {
			h.logger.LogError(ctx, NewInternalError("Empty error code", WithService(h.service)))
		}
	}

	// Log error
	if h.logger != nil {
		h.logger.LogError(ctx, domainErr)
	}

	// Record metrics
	if h.metrics != nil {
		h.metrics.RecordError(ctx, domainErr)
	}

	// Report error if critical
	if domainErr.IsCritical() && h.reporter != nil {
		h.reporter.ReportError(ctx, domainErr)
	}

	// Attempt recovery if possible
	if h.recoverer != nil && h.recoverer.CanRecover(domainErr) {
		if recoveryErr := h.recoverer.Recover(ctx, domainErr); recoveryErr != nil {
			// Log recovery failure
			if h.logger != nil {
				h.logger.LogError(ctx, NewInternalError("Recovery failed: "+recoveryErr.Error(), WithService(h.service)))
			}
		}
	}

	return domainErr
}

// HandlePanic handles a panic with comprehensive processing
func (h *ErrorHandler) HandlePanic(ctx context.Context, panic interface{}) *DomainError {
	if panic == nil {
		return nil
	}

	// Capture stack trace
	stackTrace := captureStackTraceHandler()

	// Create domain error for panic
	domainErr := NewInternalError(
		fmt.Sprintf("Panic: %v", panic),
		WithService(h.service),
		WithSeverity(SeverityCritical),
		WithStackTrace(),
	)

	// Log panic
	if h.logger != nil {
		h.logger.LogPanic(ctx, panic, stackTrace)
	}

	// Record metrics
	if h.metrics != nil {
		h.metrics.RecordPanic(ctx, panic)
	}

	// Report panic
	if h.reporter != nil {
		h.reporter.ReportPanic(ctx, panic, stackTrace)
	}

	return domainErr
}

// HandleWithRecovery handles an error with recovery attempt
func (h *ErrorHandler) HandleWithRecovery(ctx context.Context, err error, recoveryFunc func() error) *DomainError {
	domainErr := h.HandleError(ctx, err)
	if domainErr == nil {
		return nil
	}

	// Attempt recovery
	if recoveryFunc != nil {
		if recoveryErr := recoveryFunc(); recoveryErr != nil {
			// Log recovery failure
			if h.logger != nil {
				h.logger.LogError(ctx, NewInternalError("Recovery failed: "+recoveryErr.Error(), WithService(h.service)))
			}
		}
	}

	return domainErr
}

// HandleWithRetry handles an error with retry logic
func (h *ErrorHandler) HandleWithRetry(ctx context.Context, err error, maxRetries int, retryFunc func() error) *DomainError {
	domainErr := h.HandleError(ctx, err)
	if domainErr == nil {
		return nil
	}

	// Check if error is retryable
	if !domainErr.IsRetryable() {
		return domainErr
	}

	// Attempt retry
	for i := 0; i < maxRetries; i++ {
		if retryErr := retryFunc(); retryErr == nil {
			// Success, return nil
			return nil
		}

		// Wait before next retry with context cancellation support
		select {
		case <-time.After(time.Duration(i+1) * time.Second):
			// Continue with retry
		case <-ctx.Done():
			return NewInternalError("context cancelled", WithService(h.service))
		}
	}

	// All retries failed
	return NewInternalError("All retry attempts failed", WithService(h.service))
}

// DefaultErrorLogger provides a default error logger implementation
type DefaultErrorLogger struct {
	service string
}

// NewDefaultErrorLogger creates a new default error logger
func NewDefaultErrorLogger(service string) *DefaultErrorLogger {
	return &DefaultErrorLogger{
		service: service,
	}
}

// LogError logs an error
func (l *DefaultErrorLogger) LogError(ctx context.Context, err *DomainError) {
	// In a real implementation, you would use a proper logger
	fmt.Printf("[%s] ERROR: %s\n", l.service, err.Error())
}

// LogPanic logs a panic
func (l *DefaultErrorLogger) LogPanic(ctx context.Context, panic interface{}, stackTrace string) {
	// In a real implementation, you would use a proper logger
	fmt.Printf("[%s] PANIC: %v\n%s\n", l.service, panic, stackTrace)
}

// DefaultErrorMetrics provides a default error metrics implementation
type DefaultErrorMetrics struct {
	service string
}

// NewDefaultErrorMetrics creates a new default error metrics collector
func NewDefaultErrorMetrics(service string) *DefaultErrorMetrics {
	return &DefaultErrorMetrics{
		service: service,
	}
}

// RecordError records an error metric
func (m *DefaultErrorMetrics) RecordError(ctx context.Context, err *DomainError) {
	// In a real implementation, you would use a metrics library like Prometheus
	fmt.Printf("[%s] METRIC: Error recorded - %s\n", m.service, err.Code)
}

// RecordPanic records a panic metric
func (m *DefaultErrorMetrics) RecordPanic(ctx context.Context, panic interface{}) {
	// In a real implementation, you would use a metrics library like Prometheus
	fmt.Printf("[%s] METRIC: Panic recorded - %v\n", m.service, panic)
}

// DefaultErrorReporter provides a default error reporter implementation
type DefaultErrorReporter struct {
	service string
}

// NewDefaultErrorReporter creates a new default error reporter
func NewDefaultErrorReporter(service string) *DefaultErrorReporter {
	return &DefaultErrorReporter{
		service: service,
	}
}

// ReportError reports an error to external services
func (r *DefaultErrorReporter) ReportError(ctx context.Context, err *DomainError) {
	// In a real implementation, you would send to external services like Sentry
	fmt.Printf("[%s] REPORT: Error reported - %s\n", r.service, err.Error())
}

// ReportPanic reports a panic to external services
func (r *DefaultErrorReporter) ReportPanic(ctx context.Context, panic interface{}, stackTrace string) {
	// In a real implementation, you would send to external services like Sentry
	fmt.Printf("[%s] REPORT: Panic reported - %v\n", r.service, panic)
}

// DefaultErrorRecoverer provides a default error recoverer implementation
type DefaultErrorRecoverer struct {
	service string
}

// NewDefaultErrorRecoverer creates a new default error recoverer
func NewDefaultErrorRecoverer(service string) *DefaultErrorRecoverer {
	return &DefaultErrorRecoverer{
		service: service,
	}
}

// CanRecover checks if an error can be recovered
func (r *DefaultErrorRecoverer) CanRecover(err *DomainError) bool {
	// Only recover from retryable errors
	return err.IsRetryable()
}

// Recover attempts to recover from an error
func (r *DefaultErrorRecoverer) Recover(ctx context.Context, err *DomainError) error {
	// In a real implementation, you would implement specific recovery strategies
	fmt.Printf("[%s] RECOVER: Attempting recovery from %s\n", r.service, err.Code)
	return nil
}

// ErrorHandlerMiddleware provides middleware for error handling
type ErrorHandlerMiddleware struct {
	handler *ErrorHandler
}

// NewErrorHandlerMiddleware creates a new error handler middleware
func NewErrorHandlerMiddleware(handler *ErrorHandler) *ErrorHandlerMiddleware {
	return &ErrorHandlerMiddleware{
		handler: handler,
	}
}

// WrapFunction wraps a function with error handling
func (m *ErrorHandlerMiddleware) WrapFunction(fn func() error) func() error {
	return func() error {
		defer func() {
			if r := recover(); r != nil {
				m.handler.HandlePanic(context.Background(), r)
			}
		}()

		if err := fn(); err != nil {
			return m.handler.HandleError(context.Background(), err)
		}

		return nil
	}
}

// WrapFunctionWithContext wraps a function with error handling and context
func (m *ErrorHandlerMiddleware) WrapFunctionWithContext(fn func(context.Context) error) func(context.Context) error {
	return func(ctx context.Context) error {
		defer func() {
			if r := recover(); r != nil {
				m.handler.HandlePanic(ctx, r)
			}
		}()

		if err := fn(ctx); err != nil {
			return m.handler.HandleError(ctx, err)
		}

		return nil
	}
}

// ErrorHandlerChain provides a chain of error handlers
type ErrorHandlerChain struct {
	handlers []*ErrorHandler
}

// NewErrorHandlerChain creates a new error handler chain
func NewErrorHandlerChain(handlers ...*ErrorHandler) *ErrorHandlerChain {
	return &ErrorHandlerChain{
		handlers: handlers,
	}
}

// AddHandler adds a handler to the chain
func (c *ErrorHandlerChain) AddHandler(handler *ErrorHandler) {
	c.handlers = append(c.handlers, handler)
}

// HandleError handles an error through the chain
func (c *ErrorHandlerChain) HandleError(ctx context.Context, err error) *DomainError {
	var lastErr *DomainError

	for _, handler := range c.handlers {
		lastErr = handler.HandleError(ctx, err)
	}

	return lastErr
}

// HandlePanic handles a panic through the chain
func (c *ErrorHandlerChain) HandlePanic(ctx context.Context, panic interface{}) *DomainError {
	var lastErr *DomainError

	for _, handler := range c.handlers {
		lastErr = handler.HandlePanic(ctx, panic)
	}

	return lastErr
}

// captureStackTrace captures the current stack trace
func captureStackTraceHandler() string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}
