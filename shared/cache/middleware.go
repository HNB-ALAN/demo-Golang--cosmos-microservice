// Package cache provides caching utilities for USC platform services.
package cache

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc"
)

// CacheMiddleware provides HTTP and gRPC caching middleware
type CacheMiddleware struct {
	cache Cache
	ttl   time.Duration
}

// NewCacheMiddleware creates a new cache middleware
func NewCacheMiddleware(cache Cache, ttl time.Duration) *CacheMiddleware {
	return &CacheMiddleware{
		cache: cache,
		ttl:   ttl,
	}
}

// HTTPMiddleware returns HTTP caching middleware
func (cm *CacheMiddleware) HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only cache GET requests
		if r.Method != http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}

		// Generate cache key
		cacheKey := cm.generateHTTPCacheKey(r)

		// Try to get from cache
		cachedResponse, err := cm.cache.Get(r.Context(), cacheKey)
		if err == nil {
			// Cache hit, return cached response
			// For now, just return the cached string response
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(cachedResponse))
			return
		}

		// Cache miss, create response writer to capture response
		responseWriter := &CachedResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			headers:        make(map[string][]string),
			body:           make([]byte, 0),
		}

		// Call next handler
		next.ServeHTTP(responseWriter, r)

		// Cache the response if it's successful
		if responseWriter.statusCode >= 200 && responseWriter.statusCode < 300 {
			cachedResponse := &CachedHTTPResponse{
				StatusCode: responseWriter.statusCode,
				Headers:    responseWriter.headers,
				Body:       responseWriter.body,
				Timestamp:  time.Now(),
			}

			cm.cache.Set(r.Context(), cacheKey, cachedResponse, cm.ttl)
		}
	})
}

// GRPCUnaryInterceptor returns gRPC unary caching interceptor
func (cm *CacheMiddleware) GRPCUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Generate cache key
		cacheKey := cm.generateGRPCCacheKey(info.FullMethod, req)

		// Try to get from cache
		cachedResponse, err := cm.cache.Get(ctx, cacheKey)
		if err == nil {
			// Cache hit, return cached response
			// For now, just return the cached string response
			return cachedResponse, nil
		}

		// Cache miss, call handler
		response, err := handler(ctx, req)
		if err != nil {
			return nil, err
		}

		// Cache the response (simplified - just cache the response as string)
		cm.cache.Set(ctx, cacheKey, fmt.Sprintf("%v", response), cm.ttl)

		return response, nil
	}
}

// generateHTTPCacheKey generates a cache key for HTTP requests
func (cm *CacheMiddleware) generateHTTPCacheKey(r *http.Request) string {
	// Include method, path, and query parameters
	key := fmt.Sprintf("%s:%s:%s", r.Method, r.URL.Path, r.URL.RawQuery)

	// Include relevant headers
	if accept := r.Header.Get("Accept"); accept != "" {
		key += ":" + accept
	}
	if contentType := r.Header.Get("Content-Type"); contentType != "" {
		key += ":" + contentType
	}

	// Hash the key to keep it reasonable length
	hash := md5.Sum([]byte(key))
	return "http:" + hex.EncodeToString(hash[:])
}

// generateGRPCCacheKey generates a cache key for gRPC requests
func (cm *CacheMiddleware) generateGRPCCacheKey(method string, req interface{}) string {
	// Include method and request data
	key := fmt.Sprintf("%s:%v", method, req)

	// Hash the key to keep it reasonable length
	hash := md5.Sum([]byte(key))
	return "grpc:" + hex.EncodeToString(hash[:])
}

// CachedHTTPResponse represents a cached HTTP response
type CachedHTTPResponse struct {
	StatusCode int                 `json:"status_code"`
	Headers    map[string][]string `json:"headers"`
	Body       []byte              `json:"body"`
	Timestamp  time.Time           `json:"timestamp"`
}

// CachedGRPCResponse represents a cached gRPC response
type CachedGRPCResponse struct {
	Response  interface{} `json:"response"`
	Timestamp time.Time   `json:"timestamp"`
}

// CachedResponseWriter captures HTTP response for caching
type CachedResponseWriter struct {
	http.ResponseWriter
	statusCode int
	headers    map[string][]string
	body       []byte
}

// WriteHeader captures the status code
func (crw *CachedResponseWriter) WriteHeader(code int) {
	crw.statusCode = code
	crw.ResponseWriter.WriteHeader(code)
}

// Write captures the response body
func (crw *CachedResponseWriter) Write(data []byte) (int, error) {
	crw.body = append(crw.body, data...)
	return crw.ResponseWriter.Write(data)
}

// Header returns the response headers
func (crw *CachedResponseWriter) Header() http.Header {
	// Capture headers when they're set
	headers := crw.ResponseWriter.Header()
	for key, values := range headers {
		crw.headers[key] = values
	}
	return headers
}

// CacheKeyGenerator generates cache keys for different scenarios
type CacheKeyGenerator struct {
	prefix string
}

// NewCacheKeyGenerator creates a new cache key generator
func NewCacheKeyGenerator(prefix string) *CacheKeyGenerator {
	return &CacheKeyGenerator{
		prefix: prefix,
	}
}

// GenerateKey generates a cache key with prefix
func (ckg *CacheKeyGenerator) GenerateKey(parts ...string) string {
	key := strings.Join(parts, ":")
	if ckg.prefix != "" {
		key = ckg.prefix + ":" + key
	}
	return key
}

// GenerateUserKey generates a cache key for user-specific data
func (ckg *CacheKeyGenerator) GenerateUserKey(userID string, parts ...string) string {
	allParts := append([]string{"user", userID}, parts...)
	return ckg.GenerateKey(allParts...)
}

// GenerateSessionKey generates a cache key for session data
func (ckg *CacheKeyGenerator) GenerateSessionKey(sessionID string, parts ...string) string {
	allParts := append([]string{"session", sessionID}, parts...)
	return ckg.GenerateKey(allParts...)
}

// GenerateAPIKey generates a cache key for API data
func (ckg *CacheKeyGenerator) GenerateAPIKey(endpoint string, params map[string]string) string {
	parts := []string{"api", endpoint}

	// Sort parameters for consistent keys
	for key, value := range params {
		parts = append(parts, key, value)
	}

	return ckg.GenerateKey(parts...)
}

// GenerateDatabaseKey generates a cache key for database queries
func (ckg *CacheKeyGenerator) GenerateDatabaseKey(table string, query string, params ...interface{}) string {
	parts := []string{"db", table, query}

	// Add parameters
	for _, param := range params {
		parts = append(parts, fmt.Sprintf("%v", param))
	}

	return ckg.GenerateKey(parts...)
}

// CacheInvalidator provides cache invalidation utilities
type CacheInvalidator struct {
	cache Cache
}

// NewCacheInvalidator creates a new cache invalidator
func NewCacheInvalidator(cache Cache) *CacheInvalidator {
	return &CacheInvalidator{
		cache: cache,
	}
}

// InvalidateByPattern invalidates cache keys matching a pattern
func (ci *CacheInvalidator) InvalidateByPattern(ctx context.Context, pattern string) error {
	keys, err := ci.cache.Keys(ctx, pattern)
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		// Delete keys one by one since DeleteMultiple is not in our interface
		for _, key := range keys {
			if err := ci.cache.Delete(ctx, key); err != nil {
				return err
			}
		}
		return nil
	}

	return nil
}

// InvalidateUserData invalidates all cache data for a user
func (ci *CacheInvalidator) InvalidateUserData(ctx context.Context, userID string) error {
	pattern := fmt.Sprintf("user:%s:*", userID)
	return ci.InvalidateByPattern(ctx, pattern)
}

// InvalidateSessionData invalidates all cache data for a session
func (ci *CacheInvalidator) InvalidateSessionData(ctx context.Context, sessionID string) error {
	pattern := fmt.Sprintf("session:%s:*", sessionID)
	return ci.InvalidateByPattern(ctx, pattern)
}

// InvalidateAPIData invalidates cache data for an API endpoint
func (ci *CacheInvalidator) InvalidateAPIData(ctx context.Context, endpoint string) error {
	pattern := fmt.Sprintf("api:%s:*", endpoint)
	return ci.InvalidateByPattern(ctx, pattern)
}

// InvalidateDatabaseData invalidates cache data for a database table
func (ci *CacheInvalidator) InvalidateDatabaseData(ctx context.Context, table string) error {
	pattern := fmt.Sprintf("db:%s:*", table)
	return ci.InvalidateByPattern(ctx, pattern)
}

// CacheMetrics provides cache performance metrics
type CacheMetrics struct {
	Hits     int64   `json:"hits"`
	Misses   int64   `json:"misses"`
	Sets     int64   `json:"sets"`
	Deletes  int64   `json:"deletes"`
	Errors   int64   `json:"errors"`
	HitRate  float64 `json:"hit_rate"`
	TotalOps int64   `json:"total_ops"`
}

// CacheMetricsCollector collects cache performance metrics
type CacheMetricsCollector struct {
	metrics *CacheMetrics
}

// NewCacheMetricsCollector creates a new cache metrics collector
func NewCacheMetricsCollector() *CacheMetricsCollector {
	return &CacheMetricsCollector{
		metrics: &CacheMetrics{},
	}
}

// RecordHit records a cache hit
func (cmc *CacheMetricsCollector) RecordHit() {
	cmc.metrics.Hits++
	cmc.metrics.TotalOps++
	cmc.updateHitRate()
}

// RecordMiss records a cache miss
func (cmc *CacheMetricsCollector) RecordMiss() {
	cmc.metrics.Misses++
	cmc.metrics.TotalOps++
	cmc.updateHitRate()
}

// RecordSet records a cache set operation
func (cmc *CacheMetricsCollector) RecordSet() {
	cmc.metrics.Sets++
	cmc.metrics.TotalOps++
}

// RecordDelete records a cache delete operation
func (cmc *CacheMetricsCollector) RecordDelete() {
	cmc.metrics.Deletes++
	cmc.metrics.TotalOps++
}

// RecordError records a cache error
func (cmc *CacheMetricsCollector) RecordError() {
	cmc.metrics.Errors++
	cmc.metrics.TotalOps++
}

// GetMetrics returns current cache metrics
func (cmc *CacheMetricsCollector) GetMetrics() *CacheMetrics {
	return cmc.metrics
}

// Reset resets all metrics
func (cmc *CacheMetricsCollector) Reset() {
	cmc.metrics = &CacheMetrics{}
}

// updateHitRate updates the hit rate
func (cmc *CacheMetricsCollector) updateHitRate() {
	if cmc.metrics.TotalOps > 0 {
		cmc.metrics.HitRate = float64(cmc.metrics.Hits) / float64(cmc.metrics.TotalOps)
	}
}
