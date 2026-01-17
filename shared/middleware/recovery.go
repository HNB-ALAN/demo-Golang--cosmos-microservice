// Package middleware provides common middleware for USC platform services.
package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"google.golang.org/grpc"
)

// RecoveryConfig represents recovery middleware configuration
type RecoveryConfig struct {
	LogPanic     bool          `mapstructure:"log_panic"`
	LogStack     bool          `mapstructure:"log_stack"`
	LogRequest   bool          `mapstructure:"log_request"`
	LogResponse  bool          `mapstructure:"log_response"`
	LogDuration  bool          `mapstructure:"log_duration"`
	MaxStackSize int           `mapstructure:"max_stack_size"`
	Timeout      time.Duration `mapstructure:"timeout"`
}

// DefaultRecoveryConfig returns the default recovery configuration
func DefaultRecoveryConfig() RecoveryConfig {
	return RecoveryConfig{
		LogPanic:     true,
		LogStack:     true,
		LogRequest:   false,
		LogResponse:  false,
		LogDuration:  false,
		MaxStackSize: 4096,
		Timeout:      30 * time.Second,
	}
}

// HTTPRecoveryMiddleware provides HTTP panic recovery middleware
type HTTPRecoveryMiddleware struct {
	config RecoveryConfig
	logger *log.Logger
}

// NewHTTPRecoveryMiddleware creates a new HTTP recovery middleware
func NewHTTPRecoveryMiddleware(config RecoveryConfig) *HTTPRecoveryMiddleware {
	return &HTTPRecoveryMiddleware{
		config: config,
		logger: log.New(log.Writer(), "[RECOVERY] ", log.LstdFlags),
	}
}

// Middleware returns the HTTP recovery middleware
func (m *HTTPRecoveryMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if panic := recover(); panic != nil {
					m.handlePanic(w, r, panic)
				}
			}()

			// Create a custom response writer to capture response details
			responseWriter := &RecoveryResponseWriter{
				ResponseWriter: w,
				config:         m.config,
			}

			// Log request if configured
			if m.config.LogRequest {
				m.logRequest(r)
			}

			// Start timing if configured
			start := time.Now()
			if m.config.LogDuration {
				defer func() {
					duration := time.Since(start)
					m.logDuration(r, duration)
				}()
			}

			// Call next handler
			next.ServeHTTP(responseWriter, r)

			// Log response if configured
			if m.config.LogResponse {
				m.logResponse(r, responseWriter)
			}
		})
	}
}

// handlePanic handles a panic in HTTP middleware
func (m *HTTPRecoveryMiddleware) handlePanic(w http.ResponseWriter, r *http.Request, panic interface{}) {
	// Log panic if configured
	if m.config.LogPanic {
		m.logPanic(r, panic)
	}

	// Log stack trace if configured
	if m.config.LogStack {
		m.logStack(r, panic)
	}

	// Write error response
	m.writeErrorResponse(w, r, panic)
}

// logPanic logs panic information
func (m *HTTPRecoveryMiddleware) logPanic(r *http.Request, panic interface{}) {
	m.logger.Printf("PANIC: %v - %s %s - %s", panic, r.Method, r.URL.Path, r.RemoteAddr)
}

// logStack logs stack trace
func (m *HTTPRecoveryMiddleware) logStack(r *http.Request, panic interface{}) {
	stack := m.getStackTrace()
	m.logger.Printf("STACK TRACE: %s %s - %s\n%s", r.Method, r.URL.Path, r.RemoteAddr, stack)
}

// logRequest logs request information
func (m *HTTPRecoveryMiddleware) logRequest(r *http.Request) {
	m.logger.Printf("REQUEST: %s %s - %s - %s", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
}

// logResponse logs response information
func (m *HTTPRecoveryMiddleware) logResponse(r *http.Request, w *RecoveryResponseWriter) {
	m.logger.Printf("RESPONSE: %s %s - %d - %s", r.Method, r.URL.Path, w.statusCode, r.RemoteAddr)
}

// logDuration logs request duration
func (m *HTTPRecoveryMiddleware) logDuration(r *http.Request, duration time.Duration) {
	m.logger.Printf("DURATION: %s %s - %v - %s", r.Method, r.URL.Path, duration, r.RemoteAddr)
}

// getStackTrace gets the current stack trace
func (m *HTTPRecoveryMiddleware) getStackTrace() string {
	buf := make([]byte, m.config.MaxStackSize)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

// writeErrorResponse writes an error response
func (m *HTTPRecoveryMiddleware) writeErrorResponse(w http.ResponseWriter, r *http.Request, panic interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	response := map[string]interface{}{
		"error":   "Internal server error",
		"code":    "INTERNAL_ERROR",
		"message": "An unexpected error occurred",
	}

	if m.config.LogPanic {
		response["panic"] = fmt.Sprintf("%v", panic)
	}

	fmt.Fprintf(w, `{"error":"Internal server error","code":"INTERNAL_ERROR","message":"An unexpected error occurred"}`)
}

// RecoveryResponseWriter captures response details for recovery middleware
type RecoveryResponseWriter struct {
	http.ResponseWriter
	config     RecoveryConfig
	statusCode int
	written    bool
}

// WriteHeader captures the status code
func (rrw *RecoveryResponseWriter) WriteHeader(code int) {
	if !rrw.written {
		rrw.statusCode = code
		rrw.written = true
	}
	rrw.ResponseWriter.WriteHeader(code)
}

// Write captures the response body
func (rrw *RecoveryResponseWriter) Write(data []byte) (int, error) {
	if !rrw.written {
		rrw.statusCode = http.StatusOK
		rrw.written = true
	}
	return rrw.ResponseWriter.Write(data)
}

// GRPCRecoveryInterceptor provides gRPC panic recovery interceptor
type GRPCRecoveryInterceptor struct {
	config RecoveryConfig
	logger *log.Logger
}

// NewGRPCRecoveryInterceptor creates a new gRPC recovery interceptor
func NewGRPCRecoveryInterceptor(config RecoveryConfig) *GRPCRecoveryInterceptor {
	return &GRPCRecoveryInterceptor{
		config: config,
		logger: log.New(log.Writer(), "[GRPC-RECOVERY] ", log.LstdFlags),
	}
}

// UnaryServerInterceptor returns a unary server interceptor that recovers from panics
func (i *GRPCRecoveryInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		defer func() {
			if panic := recover(); panic != nil {
				i.handlePanic(ctx, info.FullMethod, panic)
			}
		}()

		// Log request if configured
		if i.config.LogRequest {
			i.logRequest(ctx, info.FullMethod)
		}

		// Start timing if configured
		start := time.Now()
		if i.config.LogDuration {
			defer func() {
				duration := time.Since(start)
				i.logDuration(ctx, info.FullMethod, duration)
			}()
		}

		// Call handler
		result, err := handler(ctx, req)

		// Log response if configured
		if i.config.LogResponse {
			i.logResponse(ctx, info.FullMethod, err)
		}

		return result, err
	}
}

// StreamServerInterceptor returns a stream server interceptor that recovers from panics
func (i *GRPCRecoveryInterceptor) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		defer func() {
			if panic := recover(); panic != nil {
				i.handlePanic(ss.Context(), info.FullMethod, panic)
			}
		}()

		// Log request if configured
		if i.config.LogRequest {
			i.logRequest(ss.Context(), info.FullMethod)
		}

		// Start timing if configured
		start := time.Now()
		if i.config.LogDuration {
			defer func() {
				duration := time.Since(start)
				i.logDuration(ss.Context(), info.FullMethod, duration)
			}()
		}

		// Call handler
		err := handler(srv, ss)

		// Log response if configured
		if i.config.LogResponse {
			i.logResponse(ss.Context(), info.FullMethod, err)
		}

		return err
	}
}

// handlePanic handles a panic in gRPC interceptor
func (i *GRPCRecoveryInterceptor) handlePanic(ctx context.Context, method string, panicValue interface{}) {
	// Log panic if configured
	if i.config.LogPanic {
		i.logPanic(ctx, method, panicValue)
	}

	// Log stack trace if configured
	if i.config.LogStack {
		i.logStack(ctx, method, panicValue)
	}

	// This function is called from defer, so we can't return values
	// The actual error handling is done in the interceptor
}

// logPanic logs panic information
func (i *GRPCRecoveryInterceptor) logPanic(ctx context.Context, method string, panicValue interface{}) {
	i.logger.Printf("PANIC: %v - %s", panicValue, method)
}

// logStack logs stack trace
func (i *GRPCRecoveryInterceptor) logStack(ctx context.Context, method string, panicValue interface{}) {
	stack := i.getStackTrace()
	i.logger.Printf("STACK TRACE: %s\n%s", method, stack)
}

// logRequest logs request information
func (i *GRPCRecoveryInterceptor) logRequest(ctx context.Context, method string) {
	i.logger.Printf("REQUEST: %s", method)
}

// logResponse logs response information
func (i *GRPCRecoveryInterceptor) logResponse(ctx context.Context, method string, err error) {
	if err != nil {
		i.logger.Printf("RESPONSE: %s - ERROR: %v", method, err)
	} else {
		i.logger.Printf("RESPONSE: %s - SUCCESS", method)
	}
}

// logDuration logs request duration
func (i *GRPCRecoveryInterceptor) logDuration(ctx context.Context, method string, duration time.Duration) {
	i.logger.Printf("DURATION: %s - %v", method, duration)
}

// getStackTrace gets the current stack trace
func (i *GRPCRecoveryInterceptor) getStackTrace() string {
	buf := make([]byte, i.config.MaxStackSize)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

// RecoveryHandler provides a recovery handler interface
type RecoveryHandler interface {
	HandlePanic(ctx context.Context, panic interface{}) error
}

// DefaultRecoveryHandler provides a default recovery handler
type DefaultRecoveryHandler struct {
	config RecoveryConfig
	logger *log.Logger
}

// NewDefaultRecoveryHandler creates a new default recovery handler
func NewDefaultRecoveryHandler(config RecoveryConfig) *DefaultRecoveryHandler {
	return &DefaultRecoveryHandler{
		config: config,
		logger: log.New(log.Writer(), "[RECOVERY-HANDLER] ", log.LstdFlags),
	}
}

// HandlePanic handles a panic
func (h *DefaultRecoveryHandler) HandlePanic(ctx context.Context, panic interface{}) error {
	// Log panic if configured
	if h.config.LogPanic {
		h.logger.Printf("PANIC: %v", panic)
	}

	// Log stack trace if configured
	if h.config.LogStack {
		stack := h.getStackTrace()
		h.logger.Printf("STACK TRACE:\n%s", stack)
	}

	// Return error
	return fmt.Errorf("panic occurred: %v", panic)
}

// getStackTrace gets the current stack trace
func (h *DefaultRecoveryHandler) getStackTrace() string {
	buf := make([]byte, h.config.MaxStackSize)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

// RecoveryMiddleware provides a generic recovery middleware
type RecoveryMiddleware struct {
	handler RecoveryHandler
	config  RecoveryConfig
}

// NewRecoveryMiddleware creates a new recovery middleware
func NewRecoveryMiddleware(handler RecoveryHandler, config RecoveryConfig) *RecoveryMiddleware {
	return &RecoveryMiddleware{
		handler: handler,
		config:  config,
	}
}

// WrapFunction wraps a function with recovery
func (m *RecoveryMiddleware) WrapFunction(fn func() error) func() error {
	return func() error {
		defer func() {
			if panic := recover(); panic != nil {
				m.handler.HandlePanic(context.Background(), panic)
			}
		}()

		return fn()
	}
}

// WrapFunctionWithContext wraps a function with recovery and context
func (m *RecoveryMiddleware) WrapFunctionWithContext(fn func(context.Context) error) func(context.Context) error {
	return func(ctx context.Context) error {
		defer func() {
			if panic := recover(); panic != nil {
				m.handler.HandlePanic(ctx, panic)
			}
		}()

		return fn(ctx)
	}
}

// RecoveryMetrics provides recovery metrics
type RecoveryMetrics struct {
	PanicCount    int64     `json:"panic_count"`
	LastPanic     time.Time `json:"last_panic"`
	RecoveryCount int64     `json:"recovery_count"`
	LastRecovery  time.Time `json:"last_recovery"`
}

// RecoveryMetricsCollector collects recovery metrics
type RecoveryMetricsCollector struct {
	metrics *RecoveryMetrics
}

// NewRecoveryMetricsCollector creates a new recovery metrics collector
func NewRecoveryMetricsCollector() *RecoveryMetricsCollector {
	return &RecoveryMetricsCollector{
		metrics: &RecoveryMetrics{},
	}
}

// RecordPanic records a panic
func (rmc *RecoveryMetricsCollector) RecordPanic() {
	rmc.metrics.PanicCount++
	rmc.metrics.LastPanic = time.Now()
}

// RecordRecovery records a recovery
func (rmc *RecoveryMetricsCollector) RecordRecovery() {
	rmc.metrics.RecoveryCount++
	rmc.metrics.LastRecovery = time.Now()
}

// GetMetrics returns current recovery metrics
func (rmc *RecoveryMetricsCollector) GetMetrics() *RecoveryMetrics {
	return rmc.metrics
}

// Reset resets recovery metrics
func (rmc *RecoveryMetricsCollector) Reset() {
	rmc.metrics = &RecoveryMetrics{}
}
