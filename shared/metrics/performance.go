package metrics

import (
	"context"
	"sync"
	"time"
)

// PerformanceMetrics provides performance monitoring capabilities
type PerformanceMetrics struct {
	mu sync.RWMutex

	// Request metrics
	requestLatency *Histogram
	requestCount   *Counter
	errorCount     *Counter

	// Cache metrics
	cacheHits       *Counter
	cacheMisses     *Counter
	cacheOperations *Counter

	// Database metrics
	dbConnections  *Gauge
	dbQueryLatency *Histogram
	dbQueryCount   *Counter

	// System metrics
	memoryUsage    *Gauge
	cpuUsage       *Gauge
	goroutineCount *Gauge

	// Custom metrics
	customMetrics map[string]interface{}
}

// NewPerformanceMetrics creates a new performance metrics instance
func NewPerformanceMetrics() *PerformanceMetrics {
	return &PerformanceMetrics{
		requestLatency:  NewHistogram("request_latency_ms"),
		requestCount:    NewCounter("request_count"),
		errorCount:      NewCounter("error_count"),
		cacheHits:       NewCounter("cache_hits"),
		cacheMisses:     NewCounter("cache_misses"),
		cacheOperations: NewCounter("cache_operations"),
		dbConnections:   NewGauge("db_connections"),
		dbQueryLatency:  NewHistogram("db_query_latency_ms"),
		dbQueryCount:    NewCounter("db_query_count"),
		memoryUsage:     NewGauge("memory_usage_bytes"),
		cpuUsage:        NewGauge("cpu_usage_percent"),
		goroutineCount:  NewGauge("goroutine_count"),
		customMetrics:   make(map[string]interface{}),
	}
}

// RecordRequest records a request metric
func (p *PerformanceMetrics) RecordRequest(latency time.Duration, isError bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.requestCount.Increment()
	p.requestLatency.Observe(float64(latency.Milliseconds()))

	if isError {
		p.errorCount.Increment()
	}
}

// RecordCacheOperation records a cache operation
func (p *PerformanceMetrics) RecordCacheOperation(isHit bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.cacheOperations.Increment()
	if isHit {
		p.cacheHits.Increment()
	} else {
		p.cacheMisses.Increment()
	}
}

// RecordDatabaseQuery records a database query
func (p *PerformanceMetrics) RecordDatabaseQuery(latency time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.dbQueryCount.Increment()
	p.dbQueryLatency.Observe(float64(latency.Milliseconds()))
}

// SetDatabaseConnections sets the current database connection count
func (p *PerformanceMetrics) SetDatabaseConnections(count float64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.dbConnections.Set(count)
}

// SetMemoryUsage sets the current memory usage
func (p *PerformanceMetrics) SetMemoryUsage(bytes float64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.memoryUsage.Set(bytes)
}

// SetCPUUsage sets the current CPU usage
func (p *PerformanceMetrics) SetCPUUsage(percent float64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.cpuUsage.Set(percent)
}

// SetGoroutineCount sets the current goroutine count
func (p *PerformanceMetrics) SetGoroutineCount(count float64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.goroutineCount.Set(count)
}

// SetCustomMetric sets a custom metric
func (p *PerformanceMetrics) SetCustomMetric(name string, value interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.customMetrics[name] = value
}

// GetRequestMetrics returns request-related metrics
func (p *PerformanceMetrics) GetRequestMetrics() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return map[string]interface{}{
		"request_count":   p.requestCount.Get(),
		"error_count":     p.errorCount.Get(),
		"error_rate":      p.calculateErrorRate(),
		"avg_latency_ms":  p.calculateAverageLatency(),
		"latency_buckets": p.requestLatency.GetBuckets(),
	}
}

// GetCacheMetrics returns cache-related metrics
func (p *PerformanceMetrics) GetCacheMetrics() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	hits := p.cacheHits.Get()
	misses := p.cacheMisses.Get()
	total := hits + misses

	hitRate := float64(0)
	if total > 0 {
		hitRate = float64(hits) / float64(total) * 100
	}

	return map[string]interface{}{
		"cache_hits":       hits,
		"cache_misses":     misses,
		"cache_operations": p.cacheOperations.Get(),
		"hit_rate_percent": hitRate,
	}
}

// GetDatabaseMetrics returns database-related metrics
func (p *PerformanceMetrics) GetDatabaseMetrics() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return map[string]interface{}{
		"connections":           p.dbConnections.Get(),
		"query_count":           p.dbQueryCount.Get(),
		"avg_query_latency_ms":  p.calculateAverageQueryLatency(),
		"query_latency_buckets": p.dbQueryLatency.GetBuckets(),
	}
}

// GetSystemMetrics returns system-related metrics
func (p *PerformanceMetrics) GetSystemMetrics() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return map[string]interface{}{
		"memory_usage_bytes": p.memoryUsage.Get(),
		"cpu_usage_percent":  p.cpuUsage.Get(),
		"goroutine_count":    p.goroutineCount.Get(),
	}
}

// GetAllMetrics returns all performance metrics
func (p *PerformanceMetrics) GetAllMetrics() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return map[string]interface{}{
		"requests": p.GetRequestMetrics(),
		"cache":    p.GetCacheMetrics(),
		"database": p.GetDatabaseMetrics(),
		"system":   p.GetSystemMetrics(),
		"custom":   p.customMetrics,
	}
}

// calculateErrorRate calculates the error rate percentage
func (p *PerformanceMetrics) calculateErrorRate() float64 {
	requests := p.requestCount.Get()
	errors := p.errorCount.Get()

	if requests == 0 {
		return 0
	}

	return float64(errors) / float64(requests) * 100
}

// calculateAverageLatency calculates the average request latency
func (p *PerformanceMetrics) calculateAverageLatency() float64 {
	// This is a simplified calculation
	// In a real implementation, you'd want to track the sum of latencies
	buckets := p.requestLatency.GetBuckets()
	total := int64(0)
	count := int64(0)

	for bucket, bucketCount := range buckets {
		// Parse bucket range and use midpoint
		midpoint := p.getBucketMidpoint(bucket)
		total += int64(midpoint * float64(bucketCount))
		count += bucketCount
	}

	if count == 0 {
		return 0
	}

	return float64(total) / float64(count)
}

// calculateAverageQueryLatency calculates the average database query latency
func (p *PerformanceMetrics) calculateAverageQueryLatency() float64 {
	// Similar to calculateAverageLatency but for database queries
	buckets := p.dbQueryLatency.GetBuckets()
	total := int64(0)
	count := int64(0)

	for bucket, bucketCount := range buckets {
		midpoint := p.getBucketMidpoint(bucket)
		total += int64(midpoint * float64(bucketCount))
		count += bucketCount
	}

	if count == 0 {
		return 0
	}

	return float64(total) / float64(count)
}

// getBucketMidpoint returns the midpoint of a bucket range
func (p *PerformanceMetrics) getBucketMidpoint(bucket string) float64 {
	switch bucket {
	case "0-1":
		return 0.5
	case "1-5":
		return 3.0
	case "5-10":
		return 7.5
	case "10-50":
		return 30.0
	case "50-100":
		return 75.0
	case "100+":
		return 150.0 // Assumed average for 100+ range
	default:
		return 0
	}
}

// PerformanceMonitor provides monitoring capabilities with context
type PerformanceMonitor struct {
	metrics *PerformanceMetrics
	start   time.Time
}

// NewPerformanceMonitor creates a new performance monitor
func NewPerformanceMonitor(metrics *PerformanceMetrics) *PerformanceMonitor {
	return &PerformanceMonitor{
		metrics: metrics,
		start:   time.Now(),
	}
}

// StartRequest starts monitoring a request
func (m *PerformanceMonitor) StartRequest() *RequestMonitor {
	return &RequestMonitor{
		metrics: m.metrics,
		start:   time.Now(),
	}
}

// StartDatabaseQuery starts monitoring a database query
func (m *PerformanceMonitor) StartDatabaseQuery() *DatabaseQueryMonitor {
	return &DatabaseQueryMonitor{
		metrics: m.metrics,
		start:   time.Now(),
	}
}

// RequestMonitor monitors individual requests
type RequestMonitor struct {
	metrics *PerformanceMetrics
	start   time.Time
}

// Finish finishes monitoring a request
func (r *RequestMonitor) Finish(isError bool) {
	latency := time.Since(r.start)
	r.metrics.RecordRequest(latency, isError)
}

// DatabaseQueryMonitor monitors database queries
type DatabaseQueryMonitor struct {
	metrics *PerformanceMetrics
	start   time.Time
}

// Finish finishes monitoring a database query
func (d *DatabaseQueryMonitor) Finish() {
	latency := time.Since(d.start)
	d.metrics.RecordDatabaseQuery(latency)
}

// PerformanceMiddleware provides middleware for automatic performance monitoring
type PerformanceMiddleware struct {
	metrics *PerformanceMetrics
}

// NewPerformanceMiddleware creates a new performance middleware
func NewPerformanceMiddleware(metrics *PerformanceMetrics) *PerformanceMiddleware {
	return &PerformanceMiddleware{
		metrics: metrics,
	}
}

// MonitorRequest monitors a request with the given handler
func (p *PerformanceMiddleware) MonitorRequest(handler func(context.Context) error) func(context.Context) error {
	return func(ctx context.Context) error {
		start := time.Now()
		err := handler(ctx)
		latency := time.Since(start)

		p.metrics.RecordRequest(latency, err != nil)
		return err
	}
}

// MonitorCacheOperation monitors a cache operation
func (p *PerformanceMiddleware) MonitorCacheOperation(operation func() (interface{}, error)) (interface{}, error) {
	p.metrics.RecordCacheOperation(false) // Assume miss initially

	result, err := operation()
	if err == nil {
		p.metrics.RecordCacheOperation(true) // Update to hit if successful
	}

	return result, err
}
