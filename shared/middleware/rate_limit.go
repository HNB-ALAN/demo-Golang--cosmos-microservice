// Package middleware provides common middleware for USC platform services.
package middleware

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RateLimiter represents a rate limiter
type RateLimiter struct {
	limiter  *rate.Limiter
	burst    int
	interval time.Duration
}

// RateLimitConfig represents rate limiting configuration
type RateLimitConfig struct {
	RequestsPerSecond float64                    `mapstructure:"requests_per_second"`
	Burst             int                        `mapstructure:"burst"`
	Interval          time.Duration              `mapstructure:"interval"`
	KeyFunc           func(*http.Request) string `mapstructure:"-"`
}

// RateLimitManager manages multiple rate limiters
type RateLimitManager struct {
	limiters map[string]*RateLimiter
	config   RateLimitConfig
	mu       sync.RWMutex
}

// NewRateLimitManager creates a new rate limit manager
func NewRateLimitManager(config RateLimitConfig) *RateLimitManager {
	return &RateLimitManager{
		limiters: make(map[string]*RateLimiter),
		config:   config,
	}
}

// GetLimiter gets or creates a rate limiter for a key
func (rlm *RateLimitManager) GetLimiter(key string) *RateLimiter {
	rlm.mu.Lock()
	defer rlm.mu.Unlock()

	limiter, exists := rlm.limiters[key]
	if !exists {
		limiter = &RateLimiter{
			limiter:  rate.NewLimiter(rate.Limit(rlm.config.RequestsPerSecond), rlm.config.Burst),
			burst:    rlm.config.Burst,
			interval: rlm.config.Interval,
		}
		rlm.limiters[key] = limiter
	}

	return limiter
}

// Allow checks if a request is allowed
func (rlm *RateLimitManager) Allow(key string) bool {
	limiter := rlm.GetLimiter(key)
	return limiter.limiter.Allow()
}

// AllowN checks if N requests are allowed
func (rlm *RateLimitManager) AllowN(key string, n int) bool {
	limiter := rlm.GetLimiter(key)
	return limiter.limiter.AllowN(time.Now(), n)
}

// Wait waits for a request to be allowed
func (rlm *RateLimitManager) Wait(ctx context.Context, key string) error {
	limiter := rlm.GetLimiter(key)
	return limiter.limiter.Wait(ctx)
}

// WaitN waits for N requests to be allowed
func (rlm *RateLimitManager) WaitN(ctx context.Context, key string, n int) error {
	limiter := rlm.GetLimiter(key)
	return limiter.limiter.WaitN(ctx, n)
}

// Reserve reserves a request
func (rlm *RateLimitManager) Reserve(key string) *rate.Reservation {
	limiter := rlm.GetLimiter(key)
	return limiter.limiter.Reserve()
}

// ReserveN reserves N requests
func (rlm *RateLimitManager) ReserveN(key string, n int) *rate.Reservation {
	limiter := rlm.GetLimiter(key)
	return limiter.limiter.ReserveN(time.Now(), n)
}

// HTTPRateLimitMiddleware provides HTTP rate limiting middleware
type HTTPRateLimitMiddleware struct {
	manager *RateLimitManager
}

// NewHTTPRateLimitMiddleware creates a new HTTP rate limit middleware
func NewHTTPRateLimitMiddleware(config RateLimitConfig) *HTTPRateLimitMiddleware {
	return &HTTPRateLimitMiddleware{
		manager: NewRateLimitManager(config),
	}
}

// Middleware returns the HTTP rate limiting middleware
func (m *HTTPRateLimitMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get rate limit key
			key := m.getRateLimitKey(r)

			// Check if request is allowed
			if !m.manager.Allow(key) {
				m.writeRateLimitResponse(w, r)
				return
			}

			// Add rate limit headers
			m.addRateLimitHeaders(w, r, key)

			// Continue to next handler
			next.ServeHTTP(w, r)
		})
	}
}

// getRateLimitKey gets the rate limit key for a request
func (m *HTTPRateLimitMiddleware) getRateLimitKey(r *http.Request) string {
	if m.manager.config.KeyFunc != nil {
		return m.manager.config.KeyFunc(r)
	}

	// Default key function: IP address
	return r.RemoteAddr
}

// writeRateLimitResponse writes a rate limit response
func (m *HTTPRateLimitMiddleware) writeRateLimitResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%.0f", m.manager.config.RequestsPerSecond))
	w.Header().Set("X-RateLimit-Remaining", "0")
	w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(m.manager.config.Interval).Unix()))
	w.WriteHeader(http.StatusTooManyRequests)

	_ = map[string]interface{}{
		"error":   "Rate limit exceeded",
		"code":    "RATE_LIMIT_EXCEEDED",
		"message": "Too many requests",
	}

	fmt.Fprintf(w, `{"error":"Rate limit exceeded","code":"RATE_LIMIT_EXCEEDED","message":"Too many requests"}`)
}

// addRateLimitHeaders adds rate limit headers to the response
func (m *HTTPRateLimitMiddleware) addRateLimitHeaders(w http.ResponseWriter, r *http.Request, key string) {
	limiter := m.manager.GetLimiter(key)

	w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%.0f", m.manager.config.RequestsPerSecond))
	w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", limiter.limiter.Burst()))
	w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(m.manager.config.Interval).Unix()))
}

// GRPCRateLimitInterceptor provides gRPC rate limiting interceptor
type GRPCRateLimitInterceptor struct {
	manager *RateLimitManager
}

// NewGRPCRateLimitInterceptor creates a new gRPC rate limit interceptor
func NewGRPCRateLimitInterceptor(config RateLimitConfig) *GRPCRateLimitInterceptor {
	return &GRPCRateLimitInterceptor{
		manager: NewRateLimitManager(config),
	}
}

// UnaryServerInterceptor returns a unary server interceptor that rate limits requests
func (i *GRPCRateLimitInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Get rate limit key
		key := i.getRateLimitKey(ctx, info.FullMethod)

		// Check if request is allowed
		if !i.manager.Allow(key) {
			return nil, status.Error(codes.ResourceExhausted, "Rate limit exceeded")
		}

		// Continue to next handler
		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a stream server interceptor that rate limits requests
func (i *GRPCRateLimitInterceptor) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Get rate limit key
		key := i.getRateLimitKey(ss.Context(), info.FullMethod)

		// Check if request is allowed
		if !i.manager.Allow(key) {
			return status.Error(codes.ResourceExhausted, "Rate limit exceeded")
		}

		// Continue to next handler
		return handler(srv, ss)
	}
}

// getRateLimitKey gets the rate limit key for a gRPC request
func (i *GRPCRateLimitInterceptor) getRateLimitKey(ctx context.Context, method string) string {
	// Use method name as key
	return method
}

// AdaptiveRateLimiter provides adaptive rate limiting
type AdaptiveRateLimiter struct {
	baseConfig    RateLimitConfig
	currentConfig RateLimitConfig
	manager       *RateLimitManager
	mu            sync.RWMutex
	metrics       *RateLimitMetrics
}

// RateLimitMetrics tracks rate limiting metrics
type RateLimitMetrics struct {
	TotalRequests   int64     `json:"total_requests"`
	AllowedRequests int64     `json:"allowed_requests"`
	BlockedRequests int64     `json:"blocked_requests"`
	CurrentRate     float64   `json:"current_rate"`
	AverageRate     float64   `json:"average_rate"`
	LastUpdated     time.Time `json:"last_updated"`
}

// NewAdaptiveRateLimiter creates a new adaptive rate limiter
func NewAdaptiveRateLimiter(baseConfig RateLimitConfig) *AdaptiveRateLimiter {
	return &AdaptiveRateLimiter{
		baseConfig:    baseConfig,
		currentConfig: baseConfig,
		manager:       NewRateLimitManager(baseConfig),
		metrics:       &RateLimitMetrics{},
	}
}

// Allow checks if a request is allowed with adaptive rate limiting
func (arl *AdaptiveRateLimiter) Allow(key string) bool {
	arl.mu.Lock()
	defer arl.mu.Unlock()

	// Update metrics
	arl.metrics.TotalRequests++
	arl.metrics.LastUpdated = time.Now()

	// Check if request is allowed
	allowed := arl.manager.Allow(key)
	if allowed {
		arl.metrics.AllowedRequests++
	} else {
		arl.metrics.BlockedRequests++
	}

	// Update current rate
	arl.metrics.CurrentRate = float64(arl.metrics.AllowedRequests) / float64(arl.metrics.TotalRequests)

	// Adaptive logic: adjust rate based on metrics
	arl.adaptRate()

	return allowed
}

// adaptRate adjusts the rate limit based on current metrics
func (arl *AdaptiveRateLimiter) adaptRate() {
	// If blocked rate is too high, reduce rate
	if arl.metrics.BlockedRequests > 0 && float64(arl.metrics.BlockedRequests)/float64(arl.metrics.TotalRequests) > 0.1 {
		arl.currentConfig.RequestsPerSecond *= 0.9
		if arl.currentConfig.RequestsPerSecond < 1 {
			arl.currentConfig.RequestsPerSecond = 1
		}
	}

	// If blocked rate is low, increase rate
	if arl.metrics.BlockedRequests == 0 && arl.metrics.TotalRequests > 100 {
		arl.currentConfig.RequestsPerSecond *= 1.1
		if arl.currentConfig.RequestsPerSecond > arl.baseConfig.RequestsPerSecond*2 {
			arl.currentConfig.RequestsPerSecond = arl.baseConfig.RequestsPerSecond * 2
		}
	}

	// Update manager with new config
	arl.manager = NewRateLimitManager(arl.currentConfig)
}

// GetMetrics returns current rate limiting metrics
func (arl *AdaptiveRateLimiter) GetMetrics() *RateLimitMetrics {
	arl.mu.RLock()
	defer arl.mu.RUnlock()

	return &RateLimitMetrics{
		TotalRequests:   arl.metrics.TotalRequests,
		AllowedRequests: arl.metrics.AllowedRequests,
		BlockedRequests: arl.metrics.BlockedRequests,
		CurrentRate:     arl.metrics.CurrentRate,
		AverageRate:     arl.metrics.AverageRate,
		LastUpdated:     arl.metrics.LastUpdated,
	}
}

// Reset resets the rate limiter
func (arl *AdaptiveRateLimiter) Reset() {
	arl.mu.Lock()
	defer arl.mu.Unlock()

	arl.currentConfig = arl.baseConfig
	arl.manager = NewRateLimitManager(arl.currentConfig)
	arl.metrics = &RateLimitMetrics{}
}

// TokenBucketRateLimiter provides token bucket rate limiting
type TokenBucketRateLimiter struct {
	bucket     *rate.Limiter
	capacity   int
	refillRate float64
	lastRefill time.Time
	mu         sync.Mutex
}

// NewTokenBucketRateLimiter creates a new token bucket rate limiter
func NewTokenBucketRateLimiter(capacity int, refillRate float64) *TokenBucketRateLimiter {
	return &TokenBucketRateLimiter{
		bucket:     rate.NewLimiter(rate.Limit(refillRate), capacity),
		capacity:   capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Allow checks if a request is allowed
func (tbrl *TokenBucketRateLimiter) Allow() bool {
	return tbrl.bucket.Allow()
}

// AllowN checks if N requests are allowed
func (tbrl *TokenBucketRateLimiter) AllowN(n int) bool {
	return tbrl.bucket.AllowN(time.Now(), n)
}

// Wait waits for a request to be allowed
func (tbrl *TokenBucketRateLimiter) Wait(ctx context.Context) error {
	return tbrl.bucket.Wait(ctx)
}

// WaitN waits for N requests to be allowed
func (tbrl *TokenBucketRateLimiter) WaitN(ctx context.Context, n int) error {
	return tbrl.bucket.WaitN(ctx, n)
}

// GetCapacity returns the current bucket capacity
func (tbrl *TokenBucketRateLimiter) GetCapacity() int {
	return tbrl.capacity
}

// GetRefillRate returns the refill rate
func (tbrl *TokenBucketRateLimiter) GetRefillRate() float64 {
	return tbrl.refillRate
}

// SetCapacity sets the bucket capacity
func (tbrl *TokenBucketRateLimiter) SetCapacity(capacity int) {
	tbrl.mu.Lock()
	defer tbrl.mu.Unlock()

	tbrl.capacity = capacity
	tbrl.bucket = rate.NewLimiter(rate.Limit(tbrl.refillRate), capacity)
}

// SetRefillRate sets the refill rate
func (tbrl *TokenBucketRateLimiter) SetRefillRate(refillRate float64) {
	tbrl.mu.Lock()
	defer tbrl.mu.Unlock()

	tbrl.refillRate = refillRate
	tbrl.bucket = rate.NewLimiter(rate.Limit(refillRate), tbrl.capacity)
}
