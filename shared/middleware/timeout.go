// Package middleware provides common middleware for USC platform services.
package middleware

import (
	"context"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TimeoutConfig represents timeout configuration
type TimeoutConfig struct {
	DefaultTimeout time.Duration `mapstructure:"default_timeout"`
	MaxTimeout     time.Duration `mapstructure:"max_timeout"`
	MinTimeout     time.Duration `mapstructure:"min_timeout"`
	TimeoutHeader  string        `mapstructure:"timeout_header"`
	TimeoutQuery   string        `mapstructure:"timeout_query"`
}

// DefaultTimeoutConfig returns the default timeout configuration
func DefaultTimeoutConfig() TimeoutConfig {
	return TimeoutConfig{
		DefaultTimeout: 30 * time.Second,
		MaxTimeout:     5 * time.Minute,
		MinTimeout:     1 * time.Second,
		TimeoutHeader:  "X-Timeout",
		TimeoutQuery:   "timeout",
	}
}

// HTTPTimeoutMiddleware provides HTTP timeout middleware
type HTTPTimeoutMiddleware struct {
	config TimeoutConfig
}

// NewHTTPTimeoutMiddleware creates a new HTTP timeout middleware
func NewHTTPTimeoutMiddleware(config TimeoutConfig) *HTTPTimeoutMiddleware {
	return &HTTPTimeoutMiddleware{
		config: config,
	}
}

// Middleware returns the HTTP timeout middleware
func (m *HTTPTimeoutMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get timeout from request
			timeout := m.getTimeout(r)

			// Create context with timeout
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			// Create request with timeout context
			req := r.WithContext(ctx)

			// Create response writer that handles timeout
			responseWriter := &TimeoutResponseWriter{
				ResponseWriter: w,
				timeout:        timeout,
			}

			// Call next handler
			next.ServeHTTP(responseWriter, req)

			// Check if context was cancelled due to timeout
			if ctx.Err() == context.DeadlineExceeded {
				m.handleTimeout(w, r, timeout)
			}
		})
	}
}

// getTimeout gets the timeout from the request
func (m *HTTPTimeoutMiddleware) getTimeout(r *http.Request) time.Duration {
	// Try to get timeout from header
	if m.config.TimeoutHeader != "" {
		if timeoutStr := r.Header.Get(m.config.TimeoutHeader); timeoutStr != "" {
			if timeout, err := time.ParseDuration(timeoutStr); err == nil {
				return m.validateTimeout(timeout)
			}
		}
	}

	// Try to get timeout from query parameter
	if m.config.TimeoutQuery != "" {
		if timeoutStr := r.URL.Query().Get(m.config.TimeoutQuery); timeoutStr != "" {
			if timeout, err := time.ParseDuration(timeoutStr); err == nil {
				return m.validateTimeout(timeout)
			}
		}
	}

	// Use default timeout
	return m.config.DefaultTimeout
}

// validateTimeout validates and adjusts timeout
func (m *HTTPTimeoutMiddleware) validateTimeout(timeout time.Duration) time.Duration {
	// Check minimum timeout
	if timeout < m.config.MinTimeout {
		timeout = m.config.MinTimeout
	}

	// Check maximum timeout
	if timeout > m.config.MaxTimeout {
		timeout = m.config.MaxTimeout
	}

	return timeout
}

// handleTimeout handles timeout errors
func (m *HTTPTimeoutMiddleware) handleTimeout(w http.ResponseWriter, r *http.Request, timeout time.Duration) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusRequestTimeout)

	_ = map[string]interface{}{
		"error":   "Request timeout",
		"code":    "REQUEST_TIMEOUT",
		"message": "The request timed out",
		"timeout": timeout.String(),
	}

	// In a real implementation, you would use json.Marshal
	w.Write([]byte(`{"error":"Request timeout","code":"REQUEST_TIMEOUT","message":"The request timed out"}`))
}

// TimeoutResponseWriter captures response details for timeout middleware
type TimeoutResponseWriter struct {
	http.ResponseWriter
	timeout time.Duration
	written bool
}

// WriteHeader captures the status code
func (trw *TimeoutResponseWriter) WriteHeader(code int) {
	if !trw.written {
		trw.written = true
	}
	trw.ResponseWriter.WriteHeader(code)
}

// Write captures the response body
func (trw *TimeoutResponseWriter) Write(data []byte) (int, error) {
	if !trw.written {
		trw.written = true
	}
	return trw.ResponseWriter.Write(data)
}

// GRPCTimeoutInterceptor provides gRPC timeout interceptor
type GRPCTimeoutInterceptor struct {
	config TimeoutConfig
}

// NewGRPCTimeoutInterceptor creates a new gRPC timeout interceptor
func NewGRPCTimeoutInterceptor(config TimeoutConfig) *GRPCTimeoutInterceptor {
	return &GRPCTimeoutInterceptor{
		config: config,
	}
}

// UnaryServerInterceptor returns a unary server interceptor that applies timeout
func (i *GRPCTimeoutInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Get timeout from context
		timeout := i.getTimeout(ctx)

		// Create context with timeout
		timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		// Call handler with timeout context
		result, err := handler(timeoutCtx, req)

		// Check if context was cancelled due to timeout
		if timeoutCtx.Err() == context.DeadlineExceeded {
			return nil, status.Error(codes.DeadlineExceeded, "Request timeout")
		}

		return result, err
	}
}

// StreamServerInterceptor returns a stream server interceptor that applies timeout
func (i *GRPCTimeoutInterceptor) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Get timeout from context
		timeout := i.getTimeout(ss.Context())

		// Create context with timeout
		timeoutCtx, cancel := context.WithTimeout(ss.Context(), timeout)
		defer cancel()

		// Create stream with timeout context
		timeoutStream := &TimeoutServerStream{
			ServerStream: ss,
			ctx:          timeoutCtx,
		}

		// Call handler with timeout context
		err := handler(srv, timeoutStream)

		// Check if context was cancelled due to timeout
		if timeoutCtx.Err() == context.DeadlineExceeded {
			return status.Error(codes.DeadlineExceeded, "Request timeout")
		}

		return err
	}
}

// getTimeout gets the timeout from the context
func (i *GRPCTimeoutInterceptor) getTimeout(ctx context.Context) time.Duration {
	// Try to get timeout from context
	if timeout, ok := ctx.Value("timeout").(time.Duration); ok {
		return i.validateTimeout(timeout)
	}

	// Use default timeout
	return i.config.DefaultTimeout
}

// validateTimeout validates and adjusts timeout
func (i *GRPCTimeoutInterceptor) validateTimeout(timeout time.Duration) time.Duration {
	// Check minimum timeout
	if timeout < i.config.MinTimeout {
		timeout = i.config.MinTimeout
	}

	// Check maximum timeout
	if timeout > i.config.MaxTimeout {
		timeout = i.config.MaxTimeout
	}

	return timeout
}

// TimeoutServerStream wraps a gRPC ServerStream with timeout context
type TimeoutServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context returns the timeout context
func (tss *TimeoutServerStream) Context() context.Context {
	return tss.ctx
}

// TimeoutHandler provides a timeout handler
type TimeoutHandler struct {
	config TimeoutConfig
}

// NewTimeoutHandler creates a new timeout handler
func NewTimeoutHandler(config TimeoutConfig) *TimeoutHandler {
	return &TimeoutHandler{
		config: config,
	}
}

// HandleTimeout handles timeout for a request
func (h *TimeoutHandler) HandleTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	// Validate timeout
	timeout = h.validateTimeout(timeout)

	// Create context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)

	return timeoutCtx, cancel
}

// validateTimeout validates and adjusts timeout
func (h *TimeoutHandler) validateTimeout(timeout time.Duration) time.Duration {
	// Check minimum timeout
	if timeout < h.config.MinTimeout {
		timeout = h.config.MinTimeout
	}

	// Check maximum timeout
	if timeout > h.config.MaxTimeout {
		timeout = h.config.MaxTimeout
	}

	return timeout
}

// TimeoutValidator provides timeout validation
type TimeoutValidator struct {
	config TimeoutConfig
}

// NewTimeoutValidator creates a new timeout validator
func NewTimeoutValidator(config TimeoutConfig) *TimeoutValidator {
	return &TimeoutValidator{
		config: config,
	}
}

// ValidateTimeout validates a timeout duration
func (v *TimeoutValidator) ValidateTimeout(timeout time.Duration) bool {
	return timeout >= v.config.MinTimeout && timeout <= v.config.MaxTimeout
}

// AdjustTimeout adjusts a timeout to be within valid range
func (v *TimeoutValidator) AdjustTimeout(timeout time.Duration) time.Duration {
	// Check minimum timeout
	if timeout < v.config.MinTimeout {
		timeout = v.config.MinTimeout
	}

	// Check maximum timeout
	if timeout > v.config.MaxTimeout {
		timeout = v.config.MaxTimeout
	}

	return timeout
}

// TimeoutConfigBuilder provides a builder for timeout configuration
type TimeoutConfigBuilder struct {
	config TimeoutConfig
}

// NewTimeoutConfigBuilder creates a new timeout config builder
func NewTimeoutConfigBuilder() *TimeoutConfigBuilder {
	return &TimeoutConfigBuilder{
		config: DefaultTimeoutConfig(),
	}
}

// WithDefaultTimeout sets the default timeout
func (b *TimeoutConfigBuilder) WithDefaultTimeout(timeout time.Duration) *TimeoutConfigBuilder {
	b.config.DefaultTimeout = timeout
	return b
}

// WithMaxTimeout sets the maximum timeout
func (b *TimeoutConfigBuilder) WithMaxTimeout(timeout time.Duration) *TimeoutConfigBuilder {
	b.config.MaxTimeout = timeout
	return b
}

// WithMinTimeout sets the minimum timeout
func (b *TimeoutConfigBuilder) WithMinTimeout(timeout time.Duration) *TimeoutConfigBuilder {
	b.config.MinTimeout = timeout
	return b
}

// WithTimeoutHeader sets the timeout header name
func (b *TimeoutConfigBuilder) WithTimeoutHeader(header string) *TimeoutConfigBuilder {
	b.config.TimeoutHeader = header
	return b
}

// WithTimeoutQuery sets the timeout query parameter name
func (b *TimeoutConfigBuilder) WithTimeoutQuery(query string) *TimeoutConfigBuilder {
	b.config.TimeoutQuery = query
	return b
}

// Build builds the timeout configuration
func (b *TimeoutConfigBuilder) Build() TimeoutConfig {
	return b.config
}

// TimeoutMetrics provides timeout metrics
type TimeoutMetrics struct {
	TotalRequests   int64         `json:"total_requests"`
	TimeoutRequests int64         `json:"timeout_requests"`
	AverageTimeout  time.Duration `json:"average_timeout"`
	MaxTimeout      time.Duration `json:"max_timeout"`
	MinTimeout      time.Duration `json:"min_timeout"`
	LastTimeout     time.Time     `json:"last_timeout"`
}

// TimeoutMetricsCollector collects timeout metrics
type TimeoutMetricsCollector struct {
	metrics *TimeoutMetrics
}

// NewTimeoutMetricsCollector creates a new timeout metrics collector
func NewTimeoutMetricsCollector() *TimeoutMetricsCollector {
	return &TimeoutMetricsCollector{
		metrics: &TimeoutMetrics{},
	}
}

// RecordRequest records a request
func (tmc *TimeoutMetricsCollector) RecordRequest(timeout time.Duration) {
	tmc.metrics.TotalRequests++

	// Update timeout statistics
	if tmc.metrics.AverageTimeout == 0 {
		tmc.metrics.AverageTimeout = timeout
		tmc.metrics.MaxTimeout = timeout
		tmc.metrics.MinTimeout = timeout
	} else {
		// Update average
		total := tmc.metrics.AverageTimeout * time.Duration(tmc.metrics.TotalRequests-1)
		tmc.metrics.AverageTimeout = (total + timeout) / time.Duration(tmc.metrics.TotalRequests)

		// Update max
		if timeout > tmc.metrics.MaxTimeout {
			tmc.metrics.MaxTimeout = timeout
		}

		// Update min
		if timeout < tmc.metrics.MinTimeout {
			tmc.metrics.MinTimeout = timeout
		}
	}
}

// RecordTimeout records a timeout
func (tmc *TimeoutMetricsCollector) RecordTimeout() {
	tmc.metrics.TimeoutRequests++
	tmc.metrics.LastTimeout = time.Now()
}

// GetMetrics returns current timeout metrics
func (tmc *TimeoutMetricsCollector) GetMetrics() *TimeoutMetrics {
	return tmc.metrics
}

// Reset resets timeout metrics
func (tmc *TimeoutMetricsCollector) Reset() {
	tmc.metrics = &TimeoutMetrics{}
}
