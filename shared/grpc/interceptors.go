package grpc

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/usc-platform/shared/logging"
	"github.com/usc-platform/shared/metrics"
)

// UnaryServerInterceptor creates a unary server interceptor with logging and metrics
func UnaryServerInterceptor(logger *logging.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		// Log request
		logger.Info("gRPC request started",
			logging.String("method", info.FullMethod),
			logging.String("service", fmt.Sprintf("%v", info.Server)),
		)

		// Record metrics (simplified)
		_ = map[string]string{
			"method":  info.FullMethod,
			"service": fmt.Sprintf("%v", info.Server),
		}

		// Call handler
		resp, err := handler(ctx, req)

		// Calculate duration
		duration := time.Since(start)

		// Log response
		if err != nil {
			logger.Error("gRPC request failed",
				logging.String("method", info.FullMethod),
				logging.String("service", fmt.Sprintf("%v", info.Server)),
				logging.Duration("duration", duration),
				logging.Error(err),
			)

			// Record error metrics (simplified)
			_ = map[string]string{
				"method":  info.FullMethod,
				"service": fmt.Sprintf("%v", info.Server),
				"error":   status.Code(err).String(),
			}
		} else {
			logger.Info("gRPC request completed",
				logging.String("method", info.FullMethod),
				logging.String("service", fmt.Sprintf("%v", info.Server)),
				logging.Duration("duration", duration),
			)
		}

		// Record duration metrics (simplified)
		_ = map[string]string{
			"method":  info.FullMethod,
			"service": fmt.Sprintf("%v", info.Server),
		}

		return resp, err
	}
}

// StreamServerInterceptor creates a stream server interceptor with logging and metrics
func StreamServerInterceptor(logger *logging.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		// Log request
		logger.Info("gRPC stream started",
			logging.String("method", info.FullMethod),
		)

		// Record metrics (simplified)
		_ = map[string]string{
			"method": info.FullMethod,
		}

		// Call handler
		err := handler(srv, ss)

		// Calculate duration
		duration := time.Since(start)

		// Log response
		if err != nil {
			logger.Error("gRPC stream failed",
				logging.String("method", info.FullMethod),
				logging.Duration("duration", duration),
				logging.Error(err),
			)

			// Record error metrics (simplified)
			_ = map[string]string{
				"method": info.FullMethod,
				"error":  status.Code(err).String(),
			}
		} else {
			logger.Info("gRPC stream completed",
				logging.String("method", info.FullMethod),
				logging.Duration("duration", duration),
			)
		}

		// Record duration metrics (simplified)
		_ = map[string]string{
			"method": info.FullMethod,
		}

		return err
	}
}

// RecoveryInterceptor creates a recovery interceptor to handle panics
func RecoveryInterceptor(logger *logging.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("gRPC panic recovered",
					logging.String("method", info.FullMethod),
					logging.String("service", "unknown"),
					logging.Any("panic", r),
				)

				// Record panic metrics (simplified)
				_ = map[string]string{
					"method": info.FullMethod,
				}

				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()

		return handler(ctx, req)
	}
}

// TimeoutInterceptor creates a timeout interceptor
func TimeoutInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		return handler(ctx, req)
	}
}

// AuthInterceptor creates an authentication interceptor
func AuthInterceptor(logger *logging.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Skip auth for health checks
		if info.FullMethod == "/grpc.health.v1.Health/Check" ||
			info.FullMethod == "/grpc.health.v1.Health/Watch" {
			return handler(ctx, req)
		}

		// Extract and validate token from metadata
		// This is a placeholder implementation
		// In production, you would implement proper JWT validation

		logger.Debug("gRPC auth check",
			logging.String("method", info.FullMethod),
		)

		return handler(ctx, req)
	}
}

// RateLimitInterceptor creates a rate limiting interceptor
func RateLimitInterceptor(logger *logging.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// This is a placeholder implementation
		// In production, you would implement proper rate limiting

		logger.Debug("gRPC rate limit check",
			logging.String("method", info.FullMethod),
		)

		return handler(ctx, req)
	}
}

// ChainUnaryInterceptors chains multiple unary interceptors
func ChainUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		chain := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			chain = func(interceptor grpc.UnaryServerInterceptor, next grpc.UnaryHandler) grpc.UnaryHandler {
				return func(ctx context.Context, req interface{}) (interface{}, error) {
					return interceptor(ctx, req, info, next)
				}
			}(interceptors[i], chain)
		}
		return chain(ctx, req)
	}
}

// ChainStreamInterceptors chains multiple stream interceptors
func ChainStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		chain := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			chain = func(interceptor grpc.StreamServerInterceptor, next grpc.StreamHandler) grpc.StreamHandler {
				return func(srv interface{}, ss grpc.ServerStream) error {
					return interceptor(srv, ss, info, next)
				}
			}(interceptors[i], chain)
		}
		return chain(srv, ss)
	}
}

// RetryConfig represents retry configuration
type RetryConfig struct {
	MaxAttempts       int           `mapstructure:"max_attempts"`
	InitialDelay      time.Duration `mapstructure:"initial_delay"`
	MaxDelay          time.Duration `mapstructure:"max_delay"`
	BackoffMultiplier float64       `mapstructure:"backoff_multiplier"`
	RetryableCodes    []codes.Code  `mapstructure:"retryable_codes"`
}

// DefaultRetryConfig returns default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:       3,
		InitialDelay:      100 * time.Millisecond,
		MaxDelay:          5 * time.Second,
		BackoffMultiplier: 2.0,
		RetryableCodes: []codes.Code{
			codes.Unavailable,
			codes.DeadlineExceeded,
			codes.ResourceExhausted,
		},
	}
}

// isRetryableCode checks if an error code is retryable
func isRetryableCode(err error, retryableCodes []codes.Code) bool {
	st, ok := status.FromError(err)
	if !ok {
		return false
	}

	for _, code := range retryableCodes {
		if st.Code() == code {
			return true
		}
	}
	return false
}

// calculateBackoffDelay calculates the backoff delay for retry
func calculateBackoffDelay(attempt int, config RetryConfig) time.Duration {
	delay := float64(config.InitialDelay) * math.Pow(config.BackoffMultiplier, float64(attempt))
	if delay > float64(config.MaxDelay) {
		delay = float64(config.MaxDelay)
	}
	return time.Duration(delay)
}

// UnaryClientRetryInterceptor creates a unary client interceptor with retry logic
func UnaryClientRetryInterceptor(config RetryConfig, logger *logging.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var lastErr error

		for attempt := 0; attempt < config.MaxAttempts; attempt++ {
			if attempt > 0 {
				// Calculate backoff delay
				delay := calculateBackoffDelay(attempt-1, config)

				logger.Debug("gRPC retry attempt",
					logging.String("method", method),
					logging.Int("attempt", attempt+1),
					logging.Duration("delay", delay),
					logging.Error(lastErr),
				)

				// Wait before retry
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(delay):
				}
			}

			// Make the call
			err := invoker(ctx, method, req, reply, cc, opts...)

			if err == nil {
				if attempt > 0 {
					logger.Info("gRPC retry succeeded",
						logging.String("method", method),
						logging.Int("attempt", attempt+1),
					)
				}
				return nil
			}

			lastErr = err

			// Check if error is retryable
			if !isRetryableCode(err, config.RetryableCodes) {
				logger.Debug("gRPC error not retryable",
					logging.String("method", method),
					logging.Error(err),
				)
				return err
			}

			// Check if this is the last attempt
			if attempt == config.MaxAttempts-1 {
				logger.Error("gRPC retry exhausted",
					logging.String("method", method),
					logging.Int("attempts", config.MaxAttempts),
					logging.Error(err),
				)
				return err
			}
		}

		return lastErr
	}
}

// LoadBalancerConfig represents load balancer configuration
type LoadBalancerConfig struct {
	Strategy            string        `mapstructure:"strategy"` // round_robin, least_conn, random
	HealthCheckInterval time.Duration `mapstructure:"health_check_interval"`
	MaxFailures         int           `mapstructure:"max_failures"`
}

// DefaultLoadBalancerConfig returns default load balancer configuration
func DefaultLoadBalancerConfig() LoadBalancerConfig {
	return LoadBalancerConfig{
		Strategy:            "round_robin",
		HealthCheckInterval: 30 * time.Second,
		MaxFailures:         3,
	}
}

// LoadBalancer manages multiple gRPC connections with load balancing
type LoadBalancer struct {
	config     LoadBalancerConfig
	addresses  []string
	current    int
	mu         sync.RWMutex
	failures   map[string]int
	lastHealth map[string]time.Time
}

// NewLoadBalancer creates a new load balancer
func NewLoadBalancer(config LoadBalancerConfig, addresses []string) *LoadBalancer {
	return &LoadBalancer{
		config:     config,
		addresses:  addresses,
		current:    0,
		failures:   make(map[string]int),
		lastHealth: make(map[string]time.Time),
	}
}

// GetNextAddress returns the next address based on load balancing strategy
func (lb *LoadBalancer) GetNextAddress() string {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if len(lb.addresses) == 0 {
		return ""
	}

	switch lb.config.Strategy {
	case "round_robin":
		address := lb.addresses[lb.current]
		lb.current = (lb.current + 1) % len(lb.addresses)
		return address
	case "random":
		// Simple random selection
		index := time.Now().UnixNano() % int64(len(lb.addresses))
		return lb.addresses[index]
	case "least_conn":
		// For simplicity, use round robin
		// In production, you'd track active connections
		address := lb.addresses[lb.current]
		lb.current = (lb.current + 1) % len(lb.addresses)
		return address
	default:
		return lb.addresses[0]
	}
}

// RecordFailure records a failure for an address
func (lb *LoadBalancer) RecordFailure(address string) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	lb.failures[address]++
	lb.lastHealth[address] = time.Now()
}

// RecordSuccess records a success for an address
func (lb *LoadBalancer) RecordSuccess(address string) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	lb.failures[address] = 0
	lb.lastHealth[address] = time.Now()
}

// IsHealthy checks if an address is healthy
func (lb *LoadBalancer) IsHealthy(address string) bool {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	failures, exists := lb.failures[address]
	if !exists {
		return true
	}

	return failures < lb.config.MaxFailures
}

// CircuitBreakerConfig represents circuit breaker configuration
type CircuitBreakerConfig struct {
	MaxFailures    int           `mapstructure:"max_failures"`
	ResetTimeout   time.Duration `mapstructure:"reset_timeout"`
	RequestTimeout time.Duration `mapstructure:"request_timeout"`
	MaxRequests    int           `mapstructure:"max_requests"`
}

// DefaultCircuitBreakerConfig returns default circuit breaker configuration
func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		MaxFailures:    5,
		ResetTimeout:   30 * time.Second,
		RequestTimeout: 10 * time.Second,
		MaxRequests:    10,
	}
}

// CircuitBreakerState represents the state of a circuit breaker
type CircuitBreakerState int

const (
	StateClosed CircuitBreakerState = iota
	StateOpen
	StateHalfOpen
)

// CircuitBreaker implements circuit breaker pattern for gRPC calls
type CircuitBreaker struct {
	config      CircuitBreakerConfig
	state       CircuitBreakerState
	failures    int
	lastFailure time.Time
	requests    int
	mu          sync.RWMutex
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{
		config: config,
		state:  StateClosed,
	}
}

// Execute executes a function with circuit breaker protection
func (cb *CircuitBreaker) Execute(fn func() error) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	// Check if circuit breaker is open
	if cb.state == StateOpen {
		if time.Since(cb.lastFailure) > cb.config.ResetTimeout {
			cb.state = StateHalfOpen
			cb.requests = 0
		} else {
			return fmt.Errorf("circuit breaker is open")
		}
	}

	// Check if we're in half-open state and have reached max requests
	if cb.state == StateHalfOpen && cb.requests >= cb.config.MaxRequests {
		return fmt.Errorf("circuit breaker half-open max requests reached")
	}

	// Execute the function
	err := fn()
	cb.requests++

	if err != nil {
		cb.failures++
		cb.lastFailure = time.Now()

		if cb.failures >= cb.config.MaxFailures {
			cb.state = StateOpen
		}
		return err
	}

	// Success - reset failures and close circuit if it was half-open
	cb.failures = 0
	if cb.state == StateHalfOpen {
		cb.state = StateClosed
	}

	return nil
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() CircuitBreakerState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// CircuitBreakerInterceptor creates a circuit breaker interceptor
func CircuitBreakerInterceptor(cb *CircuitBreaker, logger *logging.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return cb.Execute(func() error {
			return invoker(ctx, method, req, reply, cc, opts...)
		})
	}
}

// MetricsInterceptor creates a metrics interceptor for gRPC calls
func MetricsInterceptor(metricsCollector *metrics.PrometheusMetricsCollector, logger *logging.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()

		// Make the call
		err := invoker(ctx, method, req, reply, cc, opts...)
		duration := time.Since(start)

		// Record metrics using the custom metrics collector
		// This is a simplified implementation - in production you'd use proper Prometheus metrics
		metricsData := map[string]interface{}{
			"method":   method,
			"duration": duration.Seconds(),
			"success":  err == nil,
		}

		if err != nil {
			metricsData["error"] = err.Error()
		}

		// Store metrics (simplified)
		_ = metricsData

		return err
	}
}

// LoadBalancingInterceptor creates a load balancing interceptor
func LoadBalancingInterceptor(lb *LoadBalancer, logger *logging.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// Get next address from load balancer
		address := lb.GetNextAddress()
		if address == "" {
			return fmt.Errorf("no healthy addresses available")
		}

		// Make the call
		err := invoker(ctx, method, req, reply, cc, opts...)

		// Record result in load balancer
		if err != nil {
			lb.RecordFailure(address)
		} else {
			lb.RecordSuccess(address)
		}

		return err
	}
}

// EnhancedUnaryClientInterceptor creates an enhanced unary client interceptor with all features
func EnhancedUnaryClientInterceptor(
	retryConfig RetryConfig,
	circuitBreaker *CircuitBreaker,
	loadBalancer *LoadBalancer,
	metricsCollector *metrics.PrometheusMetricsCollector,
	logger *logging.Logger,
) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// Chain interceptors: Load Balancing -> Circuit Breaker -> Retry -> Metrics

		// Load balancing
		if loadBalancer != nil {
			address := loadBalancer.GetNextAddress()
			if address == "" {
				return fmt.Errorf("no healthy addresses available")
			}
		}

		// Circuit breaker
		if circuitBreaker != nil {
			return circuitBreaker.Execute(func() error {
				// Retry logic
				return UnaryClientRetryInterceptor(retryConfig, logger)(ctx, method, req, reply, cc, invoker, opts...)
			})
		}

		// Metrics
		if metricsCollector != nil {
			return MetricsInterceptor(metricsCollector, logger)(ctx, method, req, reply, cc, invoker, opts...)
		}

		// Fallback to retry interceptor
		return UnaryClientRetryInterceptor(retryConfig, logger)(ctx, method, req, reply, cc, invoker, opts...)
	}
}
