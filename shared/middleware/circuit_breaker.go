// Package middleware provides common middleware for USC platform services.
package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"google.golang.org/grpc"
)

// CircuitState represents the state of a circuit breaker
type CircuitState int

const (
	// CircuitClosed represents a closed circuit (normal operation)
	CircuitClosed CircuitState = iota
	// CircuitOpen represents an open circuit (failing fast)
	CircuitOpen
	// CircuitHalfOpen represents a half-open circuit (testing)
	CircuitHalfOpen
)

// String returns the string representation of the circuit state
func (cs CircuitState) String() string {
	switch cs {
	case CircuitClosed:
		return "CLOSED"
	case CircuitOpen:
		return "OPEN"
	case CircuitHalfOpen:
		return "HALF_OPEN"
	default:
		return "UNKNOWN"
	}
}

// CircuitBreakerConfig represents circuit breaker configuration
type CircuitBreakerConfig struct {
	FailureThreshold int                                                   `mapstructure:"failure_threshold"`
	SuccessThreshold int                                                   `mapstructure:"success_threshold"`
	Timeout          time.Duration                                         `mapstructure:"timeout"`
	MaxRequests      int                                                   `mapstructure:"max_requests"`
	Interval         time.Duration                                         `mapstructure:"interval"`
	ReadyToTrip      func(counts Counts) bool                              `mapstructure:"-"`
	OnStateChange    func(name string, from CircuitState, to CircuitState) `mapstructure:"-"`
}

// Counts represents the counts of requests and failures
type Counts struct {
	Requests             uint32 `json:"requests"`
	TotalSuccesses       uint32 `json:"total_successes"`
	TotalFailures        uint32 `json:"total_failures"`
	ConsecutiveSuccesses uint32 `json:"consecutive_successes"`
	ConsecutiveFailures  uint32 `json:"consecutive_failures"`
}

// CircuitBreaker represents a circuit breaker
type CircuitBreaker struct {
	name          string
	maxRequests   uint32
	interval      time.Duration
	timeout       time.Duration
	readyToTrip   func(counts Counts) bool
	onStateChange func(name string, from CircuitState, to CircuitState)

	mutex      sync.Mutex
	state      CircuitState
	generation uint64
	counts     Counts
	expiry     time.Time
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(name string, config CircuitBreakerConfig) *CircuitBreaker {
	cb := &CircuitBreaker{
		name:          name,
		maxRequests:   uint32(config.MaxRequests),
		interval:      config.Interval,
		timeout:       config.Timeout,
		readyToTrip:   config.ReadyToTrip,
		onStateChange: config.OnStateChange,
	}

	if cb.readyToTrip == nil {
		cb.readyToTrip = func(counts Counts) bool {
			return counts.ConsecutiveFailures >= uint32(config.FailureThreshold)
		}
	}

	cb.toNewGeneration(time.Now())

	return cb
}

// Execute executes a function with circuit breaker protection
func (cb *CircuitBreaker) Execute(req func() (interface{}, error)) (interface{}, error) {
	generation, err := cb.beforeRequest()
	if err != nil {
		return nil, err
	}

	defer func() {
		e := recover()
		if e != nil {
			cb.afterRequest(generation, false)
			// Convert panic to error instead of re-panicking
			err = fmt.Errorf("circuit breaker panic: %v", e)
		}
	}()

	result, err := req()
	cb.afterRequest(generation, err == nil)
	return result, err
}

// ExecuteWithContext executes a function with context and circuit breaker protection
func (cb *CircuitBreaker) ExecuteWithContext(ctx context.Context, req func(context.Context) (interface{}, error)) (interface{}, error) {
	generation, err := cb.beforeRequest()
	if err != nil {
		return nil, err
	}

	defer func() {
		e := recover()
		if e != nil {
			cb.afterRequest(generation, false)
			// Convert panic to error instead of re-panicking
			err = fmt.Errorf("circuit breaker panic: %v", e)
		}
	}()

	result, err := req(ctx)
	cb.afterRequest(generation, err == nil)
	return result, err
}

// beforeRequest checks if a request should be allowed
func (cb *CircuitBreaker) beforeRequest() (uint64, error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	now := time.Now()
	state, generation := cb.currentState(now)

	if state == CircuitOpen {
		return generation, errors.New("circuit breaker is open")
	} else if state == CircuitHalfOpen && cb.counts.Requests >= cb.maxRequests {
		return generation, errors.New("circuit breaker is half-open and max requests reached")
	}

	cb.counts.onRequest()
	return generation, nil
}

// afterRequest records the result of a request
func (cb *CircuitBreaker) afterRequest(before uint64, success bool) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	now := time.Now()
	state, generation := cb.currentState(now)
	if generation != before {
		return
	}

	if success {
		cb.onSuccess(state, now)
	} else {
		cb.onFailure(state, now)
	}
}

// currentState returns the current state and generation
func (cb *CircuitBreaker) currentState(now time.Time) (CircuitState, uint64) {
	switch cb.state {
	case CircuitClosed:
		if !cb.expiry.IsZero() && cb.expiry.Before(now) {
			cb.toNewGeneration(now)
		}
	case CircuitOpen:
		if cb.expiry.Before(now) {
			cb.setState(CircuitHalfOpen, now)
		}
	}
	return cb.state, cb.generation
}

// onSuccess handles a successful request
func (cb *CircuitBreaker) onSuccess(state CircuitState, now time.Time) {
	switch state {
	case CircuitClosed:
		cb.counts.onSuccess()
	case CircuitHalfOpen:
		cb.counts.onSuccess()
		if cb.counts.ConsecutiveSuccesses >= uint32(cb.maxRequests) {
			cb.setState(CircuitClosed, now)
		}
	}
}

// onFailure handles a failed request
func (cb *CircuitBreaker) onFailure(state CircuitState, now time.Time) {
	switch state {
	case CircuitClosed:
		cb.counts.onFailure()
		if cb.readyToTrip(cb.counts) {
			cb.setState(CircuitOpen, now)
		}
	case CircuitHalfOpen:
		cb.setState(CircuitOpen, now)
	}
}

// setState sets the circuit breaker state
func (cb *CircuitBreaker) setState(state CircuitState, now time.Time) {
	if cb.state == state {
		return
	}

	prev := cb.state
	cb.state = state

	cb.toNewGeneration(now)

	if cb.onStateChange != nil {
		cb.onStateChange(cb.name, prev, state)
	}
}

// toNewGeneration creates a new generation
func (cb *CircuitBreaker) toNewGeneration(now time.Time) {
	cb.generation++
	cb.counts.clear()

	var zero time.Time
	switch cb.state {
	case CircuitClosed:
		if cb.interval == 0 {
			cb.expiry = zero
		} else {
			cb.expiry = now.Add(cb.interval)
		}
	case CircuitOpen:
		cb.expiry = now.Add(cb.timeout)
	default: // CircuitHalfOpen
		cb.expiry = zero
	}
}

// State returns the current state of the circuit breaker
func (cb *CircuitBreaker) State() CircuitState {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	now := time.Now()
	state, _ := cb.currentState(now)
	return state
}

// Counts returns the current counts
func (cb *CircuitBreaker) Counts() Counts {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	return cb.counts
}

// Name returns the name of the circuit breaker
func (cb *CircuitBreaker) Name() string {
	return cb.name
}

// onRequest increments the request count
func (c *Counts) onRequest() {
	c.Requests++
}

// onSuccess increments the success count
func (c *Counts) onSuccess() {
	c.TotalSuccesses++
	c.ConsecutiveSuccesses++
	c.ConsecutiveFailures = 0
}

// onFailure increments the failure count
func (c *Counts) onFailure() {
	c.TotalFailures++
	c.ConsecutiveFailures++
	c.ConsecutiveSuccesses = 0
}

// clear resets all counts
func (c *Counts) clear() {
	c.Requests = 0
	c.TotalSuccesses = 0
	c.TotalFailures = 0
	c.ConsecutiveSuccesses = 0
	c.ConsecutiveFailures = 0
}

// CircuitBreakerManager manages multiple circuit breakers
type CircuitBreakerManager struct {
	breakers map[string]*CircuitBreaker
	mu       sync.RWMutex
}

// NewCircuitBreakerManager creates a new circuit breaker manager
func NewCircuitBreakerManager() *CircuitBreakerManager {
	return &CircuitBreakerManager{
		breakers: make(map[string]*CircuitBreaker),
	}
}

// GetBreaker gets or creates a circuit breaker
func (cbm *CircuitBreakerManager) GetBreaker(name string, config CircuitBreakerConfig) *CircuitBreaker {
	cbm.mu.Lock()
	defer cbm.mu.Unlock()

	breaker, exists := cbm.breakers[name]
	if !exists {
		breaker = NewCircuitBreaker(name, config)
		cbm.breakers[name] = breaker
	}

	return breaker
}

// GetBreakerState returns the state of a circuit breaker
func (cbm *CircuitBreakerManager) GetBreakerState(name string) (CircuitState, bool) {
	cbm.mu.RLock()
	defer cbm.mu.RUnlock()

	breaker, exists := cbm.breakers[name]
	if !exists {
		return CircuitClosed, false
	}

	return breaker.State(), true
}

// GetBreakerCounts returns the counts of a circuit breaker
func (cbm *CircuitBreakerManager) GetBreakerCounts(name string) (Counts, bool) {
	cbm.mu.RLock()
	defer cbm.mu.RUnlock()

	breaker, exists := cbm.breakers[name]
	if !exists {
		return Counts{}, false
	}

	return breaker.Counts(), true
}

// ListBreakers returns all circuit breaker names
func (cbm *CircuitBreakerManager) ListBreakers() []string {
	cbm.mu.RLock()
	defer cbm.mu.RUnlock()

	names := make([]string, 0, len(cbm.breakers))
	for name := range cbm.breakers {
		names = append(names, name)
	}

	return names
}

// HTTPCircuitBreakerMiddleware provides HTTP circuit breaker middleware
type HTTPCircuitBreakerMiddleware struct {
	manager *CircuitBreakerManager
	config  CircuitBreakerConfig
}

// NewHTTPCircuitBreakerMiddleware creates a new HTTP circuit breaker middleware
func NewHTTPCircuitBreakerMiddleware(config CircuitBreakerConfig) *HTTPCircuitBreakerMiddleware {
	return &HTTPCircuitBreakerMiddleware{
		manager: NewCircuitBreakerManager(),
		config:  config,
	}
}

// Middleware returns the HTTP circuit breaker middleware
func (m *HTTPCircuitBreakerMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get circuit breaker for this endpoint
			breaker := m.manager.GetBreaker(r.URL.Path, m.config)

			// Execute request with circuit breaker protection
			_, err := breaker.Execute(func() (interface{}, error) {
				// Create a custom response writer to capture status
				responseWriter := &CircuitBreakerResponseWriter{
					ResponseWriter: w,
				}

				next.ServeHTTP(responseWriter, r)

				// Check if the response indicates an error
				if responseWriter.statusCode >= 500 {
					return nil, fmt.Errorf("server error: %d", responseWriter.statusCode)
				}

				return nil, nil
			})

			// If circuit breaker is open, return error response
			if err != nil {
				m.writeCircuitBreakerResponse(w, r, err)
			}
		})
	}
}

// writeCircuitBreakerResponse writes a circuit breaker error response
func (m *HTTPCircuitBreakerMiddleware) writeCircuitBreakerResponse(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusServiceUnavailable)

	_ = map[string]interface{}{
		"error":   "Service temporarily unavailable",
		"code":    "CIRCUIT_BREAKER_OPEN",
		"message": err.Error(),
	}

	fmt.Fprintf(w, `{"error":"Service temporarily unavailable","code":"CIRCUIT_BREAKER_OPEN","message":"%s"}`, err.Error())
}

// CircuitBreakerResponseWriter captures response status for circuit breaker
type CircuitBreakerResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (cbrw *CircuitBreakerResponseWriter) WriteHeader(code int) {
	cbrw.statusCode = code
	cbrw.ResponseWriter.WriteHeader(code)
}

// GRPCCircuitBreakerInterceptor provides gRPC circuit breaker interceptor
type GRPCCircuitBreakerInterceptor struct {
	manager *CircuitBreakerManager
	config  CircuitBreakerConfig
}

// NewGRPCCircuitBreakerInterceptor creates a new gRPC circuit breaker interceptor
func NewGRPCCircuitBreakerInterceptor(config CircuitBreakerConfig) *GRPCCircuitBreakerInterceptor {
	return &GRPCCircuitBreakerInterceptor{
		manager: NewCircuitBreakerManager(),
		config:  config,
	}
}

// UnaryServerInterceptor returns a unary server interceptor that applies circuit breaker
func (i *GRPCCircuitBreakerInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Get circuit breaker for this method
		breaker := i.manager.GetBreaker(info.FullMethod, i.config)

		// Execute request with circuit breaker protection
		result, err := breaker.ExecuteWithContext(ctx, func(ctx context.Context) (interface{}, error) {
			return handler(ctx, req)
		})

		return result, err
	}
}

// StreamServerInterceptor returns a stream server interceptor that applies circuit breaker
func (i *GRPCCircuitBreakerInterceptor) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Get circuit breaker for this method
		breaker := i.manager.GetBreaker(info.FullMethod, i.config)

		// Execute request with circuit breaker protection
		_, err := breaker.ExecuteWithContext(ss.Context(), func(ctx context.Context) (interface{}, error) {
			return nil, handler(srv, ss)
		})

		return err
	}
}
