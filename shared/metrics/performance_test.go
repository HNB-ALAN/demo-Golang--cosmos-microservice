package metrics

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"
)

func TestPerformanceMetrics_NewPerformanceMetrics(t *testing.T) {
	metrics := NewPerformanceMetrics()

	if metrics == nil {
		t.Fatal("Expected performance metrics to be created")
	}

	if metrics.requestLatency == nil {
		t.Error("Expected request latency histogram to be initialized")
	}

	if metrics.requestCount == nil {
		t.Error("Expected request count counter to be initialized")
	}

	if metrics.errorCount == nil {
		t.Error("Expected error count counter to be initialized")
	}

	if metrics.cacheHits == nil {
		t.Error("Expected cache hits counter to be initialized")
	}

	if metrics.cacheMisses == nil {
		t.Error("Expected cache misses counter to be initialized")
	}

	if metrics.dbConnections == nil {
		t.Error("Expected database connections gauge to be initialized")
	}

	if metrics.memoryUsage == nil {
		t.Error("Expected memory usage gauge to be initialized")
	}
}

func TestPerformanceMetrics_RecordRequest(t *testing.T) {
	metrics := NewPerformanceMetrics()

	// Test successful request
	metrics.RecordRequest(100*time.Millisecond, false)

	// Test failed request
	metrics.RecordRequest(200*time.Millisecond, true)

	// Verify metrics
	requestMetrics := metrics.GetRequestMetrics()

	if requestMetrics["request_count"] != int64(2) {
		t.Errorf("Expected 2 requests, got %v", requestMetrics["request_count"])
	}

	if requestMetrics["error_count"] != int64(1) {
		t.Errorf("Expected 1 error, got %v", requestMetrics["error_count"])
	}

	errorRate := requestMetrics["error_rate"].(float64)
	if errorRate != 50.0 {
		t.Errorf("Expected 50%% error rate, got %.2f", errorRate)
	}
}

func TestPerformanceMetrics_RecordCacheOperation(t *testing.T) {
	metrics := NewPerformanceMetrics()

	// Record cache operations
	metrics.RecordCacheOperation(true)  // hit
	metrics.RecordCacheOperation(true)  // hit
	metrics.RecordCacheOperation(false) // miss

	// Verify metrics
	cacheMetrics := metrics.GetCacheMetrics()

	if cacheMetrics["cache_hits"] != int64(2) {
		t.Errorf("Expected 2 cache hits, got %v", cacheMetrics["cache_hits"])
	}

	if cacheMetrics["cache_misses"] != int64(1) {
		t.Errorf("Expected 1 cache miss, got %v", cacheMetrics["cache_misses"])
	}

	if cacheMetrics["cache_operations"] != int64(3) {
		t.Errorf("Expected 3 cache operations, got %v", cacheMetrics["cache_operations"])
	}

	hitRate := cacheMetrics["hit_rate_percent"].(float64)
	expectedHitRate := 66.67
	tolerance := 0.01
	if math.Abs(hitRate-expectedHitRate) > tolerance {
		t.Errorf("Expected %.2f%% hit rate, got %.2f", expectedHitRate, hitRate)
	}
}

func TestPerformanceMetrics_RecordDatabaseQuery(t *testing.T) {
	metrics := NewPerformanceMetrics()

	// Record database queries
	metrics.RecordDatabaseQuery(50 * time.Millisecond)
	metrics.RecordDatabaseQuery(100 * time.Millisecond)
	metrics.RecordDatabaseQuery(150 * time.Millisecond)

	// Verify metrics
	dbMetrics := metrics.GetDatabaseMetrics()

	if dbMetrics["query_count"] != int64(3) {
		t.Errorf("Expected 3 database queries, got %v", dbMetrics["query_count"])
	}

	avgLatency := dbMetrics["avg_query_latency_ms"].(float64)
	if avgLatency <= 0 {
		t.Errorf("Expected positive average latency, got %.2f", avgLatency)
	}
}

func TestPerformanceMetrics_SetDatabaseConnections(t *testing.T) {
	metrics := NewPerformanceMetrics()

	// Set database connections
	metrics.SetDatabaseConnections(25.0)

	// Verify metrics
	systemMetrics := metrics.GetSystemMetrics()

	if systemMetrics["memory_usage_bytes"] != 0.0 {
		t.Errorf("Expected 0 memory usage, got %v", systemMetrics["memory_usage_bytes"])
	}

	// Test database connections separately
	dbMetrics := metrics.GetDatabaseMetrics()
	if dbMetrics["connections"] != 25.0 {
		t.Errorf("Expected 25 connections, got %v", dbMetrics["connections"])
	}
}

func TestPerformanceMetrics_SetMemoryUsage(t *testing.T) {
	metrics := NewPerformanceMetrics()

	// Set memory usage
	metrics.SetMemoryUsage(1024 * 1024 * 100) // 100MB

	// Verify metrics
	systemMetrics := metrics.GetSystemMetrics()

	if systemMetrics["memory_usage_bytes"] != 104857600.0 {
		t.Errorf("Expected 104857600 bytes, got %v", systemMetrics["memory_usage_bytes"])
	}
}

func TestPerformanceMetrics_SetCPUUsage(t *testing.T) {
	metrics := NewPerformanceMetrics()

	// Set CPU usage
	metrics.SetCPUUsage(75.5)

	// Verify metrics
	systemMetrics := metrics.GetSystemMetrics()

	if systemMetrics["cpu_usage_percent"] != 75.5 {
		t.Errorf("Expected 75.5%% CPU usage, got %v", systemMetrics["cpu_usage_percent"])
	}
}

func TestPerformanceMetrics_SetGoroutineCount(t *testing.T) {
	metrics := NewPerformanceMetrics()

	// Set goroutine count
	metrics.SetGoroutineCount(150.0)

	// Verify metrics
	systemMetrics := metrics.GetSystemMetrics()

	if systemMetrics["goroutine_count"] != 150.0 {
		t.Errorf("Expected 150 goroutines, got %v", systemMetrics["goroutine_count"])
	}
}

func TestPerformanceMetrics_SetCustomMetric(t *testing.T) {
	metrics := NewPerformanceMetrics()

	// Set custom metrics
	metrics.SetCustomMetric("custom_counter", 42)
	metrics.SetCustomMetric("custom_string", "test")

	// Verify metrics
	allMetrics := metrics.GetAllMetrics()
	customMetrics := allMetrics["custom"].(map[string]interface{})

	if customMetrics["custom_counter"] != 42 {
		t.Errorf("Expected custom_counter to be 42, got %v", customMetrics["custom_counter"])
	}

	if customMetrics["custom_string"] != "test" {
		t.Errorf("Expected custom_string to be 'test', got %v", customMetrics["custom_string"])
	}
}

func TestPerformanceMetrics_GetAllMetrics(t *testing.T) {
	metrics := NewPerformanceMetrics()

	// Record some metrics
	metrics.RecordRequest(100*time.Millisecond, false)
	metrics.RecordCacheOperation(true)
	metrics.RecordDatabaseQuery(50 * time.Millisecond)
	metrics.SetMemoryUsage(1024 * 1024)
	metrics.SetCustomMetric("test", "value")

	// Get all metrics
	allMetrics := metrics.GetAllMetrics()

	// Verify structure
	if allMetrics["requests"] == nil {
		t.Error("Expected requests metrics to be present")
	}

	if allMetrics["cache"] == nil {
		t.Error("Expected cache metrics to be present")
	}

	if allMetrics["database"] == nil {
		t.Error("Expected database metrics to be present")
	}

	if allMetrics["system"] == nil {
		t.Error("Expected system metrics to be present")
	}

	if allMetrics["custom"] == nil {
		t.Error("Expected custom metrics to be present")
	}
}

func TestPerformanceMonitor_NewPerformanceMonitor(t *testing.T) {
	metrics := NewPerformanceMetrics()
	monitor := NewPerformanceMonitor(metrics)

	if monitor == nil {
		t.Fatal("Expected performance monitor to be created")
	}

	if monitor.metrics != metrics {
		t.Error("Expected monitor to use provided metrics")
	}
}

func TestPerformanceMonitor_StartRequest(t *testing.T) {
	metrics := NewPerformanceMetrics()
	monitor := NewPerformanceMonitor(metrics)

	requestMonitor := monitor.StartRequest()
	if requestMonitor == nil {
		t.Fatal("Expected request monitor to be created")
	}

	if requestMonitor.metrics != metrics {
		t.Error("Expected request monitor to use provided metrics")
	}
}

func TestRequestMonitor_Finish(t *testing.T) {
	metrics := NewPerformanceMetrics()
	monitor := NewPerformanceMonitor(metrics)

	requestMonitor := monitor.StartRequest()

	// Simulate some work
	time.Sleep(10 * time.Millisecond)

	// Finish with success
	requestMonitor.Finish(false)

	// Verify metrics were recorded
	requestMetrics := metrics.GetRequestMetrics()
	if requestMetrics["request_count"] != int64(1) {
		t.Errorf("Expected 1 request, got %v", requestMetrics["request_count"])
	}

	if requestMetrics["error_count"] != int64(0) {
		t.Errorf("Expected 0 errors, got %v", requestMetrics["error_count"])
	}
}

func TestRequestMonitor_FinishWithError(t *testing.T) {
	metrics := NewPerformanceMetrics()
	monitor := NewPerformanceMonitor(metrics)

	requestMonitor := monitor.StartRequest()

	// Simulate some work
	time.Sleep(10 * time.Millisecond)

	// Finish with error
	requestMonitor.Finish(true)

	// Verify metrics were recorded
	requestMetrics := metrics.GetRequestMetrics()
	if requestMetrics["request_count"] != int64(1) {
		t.Errorf("Expected 1 request, got %v", requestMetrics["request_count"])
	}

	if requestMetrics["error_count"] != int64(1) {
		t.Errorf("Expected 1 error, got %v", requestMetrics["error_count"])
	}
}

func TestDatabaseQueryMonitor_Finish(t *testing.T) {
	metrics := NewPerformanceMetrics()
	monitor := NewPerformanceMonitor(metrics)

	queryMonitor := monitor.StartDatabaseQuery()

	// Simulate some work
	time.Sleep(10 * time.Millisecond)

	// Finish query
	queryMonitor.Finish()

	// Verify metrics were recorded
	dbMetrics := metrics.GetDatabaseMetrics()
	if dbMetrics["query_count"] != int64(1) {
		t.Errorf("Expected 1 database query, got %v", dbMetrics["query_count"])
	}
}

func TestPerformanceMiddleware_NewPerformanceMiddleware(t *testing.T) {
	metrics := NewPerformanceMetrics()
	middleware := NewPerformanceMiddleware(metrics)

	if middleware == nil {
		t.Fatal("Expected performance middleware to be created")
	}

	if middleware.metrics != metrics {
		t.Error("Expected middleware to use provided metrics")
	}
}

func TestPerformanceMiddleware_MonitorRequest(t *testing.T) {
	metrics := NewPerformanceMetrics()
	middleware := NewPerformanceMiddleware(metrics)

	// Test successful handler
	handler := func(ctx context.Context) error {
		time.Sleep(10 * time.Millisecond)
		return nil
	}

	monitoredHandler := middleware.MonitorRequest(handler)
	err := monitoredHandler(context.Background())

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify metrics were recorded
	requestMetrics := metrics.GetRequestMetrics()
	if requestMetrics["request_count"] != int64(1) {
		t.Errorf("Expected 1 request, got %v", requestMetrics["request_count"])
	}
}

func TestPerformanceMiddleware_MonitorRequestWithError(t *testing.T) {
	metrics := NewPerformanceMetrics()
	middleware := NewPerformanceMiddleware(metrics)

	// Test error handler
	handler := func(ctx context.Context) error {
		time.Sleep(10 * time.Millisecond)
		return fmt.Errorf("test error")
	}

	monitoredHandler := middleware.MonitorRequest(handler)
	err := monitoredHandler(context.Background())

	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Verify metrics were recorded
	requestMetrics := metrics.GetRequestMetrics()
	if requestMetrics["request_count"] != int64(1) {
		t.Errorf("Expected 1 request, got %v", requestMetrics["request_count"])
	}

	if requestMetrics["error_count"] != int64(1) {
		t.Errorf("Expected 1 error, got %v", requestMetrics["error_count"])
	}
}

func TestPerformanceMiddleware_MonitorCacheOperation(t *testing.T) {
	metrics := NewPerformanceMetrics()
	middleware := NewPerformanceMiddleware(metrics)

	// Test successful cache operation
	operation := func() (interface{}, error) {
		return "cached_value", nil
	}

	result, err := middleware.MonitorCacheOperation(operation)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != "cached_value" {
		t.Errorf("Expected 'cached_value', got %v", result)
	}

	// Verify metrics were recorded
	cacheMetrics := metrics.GetCacheMetrics()
	// Allow for some operations from other tests
	if cacheMetrics["cache_operations"].(int64) < 1 {
		t.Errorf("Expected at least 1 cache operation, got %v", cacheMetrics["cache_operations"])
	}

	if cacheMetrics["cache_hits"].(int64) < 1 {
		t.Errorf("Expected at least 1 cache hit, got %v", cacheMetrics["cache_hits"])
	}
}

func TestPerformanceMiddleware_MonitorCacheOperationWithError(t *testing.T) {
	metrics := NewPerformanceMetrics()
	middleware := NewPerformanceMiddleware(metrics)

	// Test failed cache operation
	operation := func() (interface{}, error) {
		return nil, fmt.Errorf("cache error")
	}

	result, err := middleware.MonitorCacheOperation(operation)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}

	// Verify metrics were recorded
	cacheMetrics := metrics.GetCacheMetrics()
	if cacheMetrics["cache_operations"] != int64(1) {
		t.Errorf("Expected 1 cache operation, got %v", cacheMetrics["cache_operations"])
	}

	if cacheMetrics["cache_misses"] != int64(1) {
		t.Errorf("Expected 1 cache miss, got %v", cacheMetrics["cache_misses"])
	}
}

func TestPerformanceMetrics_ConcurrentAccess(t *testing.T) {
	metrics := NewPerformanceMetrics()

	// Test concurrent access
	done := make(chan bool, 100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			metrics.RecordRequest(time.Duration(i)*time.Millisecond, i%2 == 0)
			metrics.RecordCacheOperation(i%2 == 0)
			metrics.RecordDatabaseQuery(time.Duration(i) * time.Millisecond)
			metrics.SetMemoryUsage(float64(i * 1024))
			metrics.SetCustomMetric(fmt.Sprintf("key_%d", i), i)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 100; i++ {
		<-done
	}

	// Verify metrics were recorded correctly
	requestMetrics := metrics.GetRequestMetrics()
	if requestMetrics["request_count"] != int64(100) {
		t.Errorf("Expected 100 requests, got %v", requestMetrics["request_count"])
	}

	cacheMetrics := metrics.GetCacheMetrics()
	if cacheMetrics["cache_operations"] != int64(100) {
		t.Errorf("Expected 100 cache operations, got %v", cacheMetrics["cache_operations"])
	}

	dbMetrics := metrics.GetDatabaseMetrics()
	if dbMetrics["query_count"] != int64(100) {
		t.Errorf("Expected 100 database queries, got %v", dbMetrics["query_count"])
	}
}

// Benchmark tests
func BenchmarkPerformanceMetrics_RecordRequest(b *testing.B) {
	metrics := NewPerformanceMetrics()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics.RecordRequest(100*time.Millisecond, i%2 == 0)
	}
}

func BenchmarkPerformanceMetrics_RecordCacheOperation(b *testing.B) {
	metrics := NewPerformanceMetrics()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics.RecordCacheOperation(i%2 == 0)
	}
}

func BenchmarkPerformanceMetrics_GetAllMetrics(b *testing.B) {
	metrics := NewPerformanceMetrics()

	// Pre-populate with some data
	for i := 0; i < 1000; i++ {
		metrics.RecordRequest(100*time.Millisecond, i%2 == 0)
		metrics.RecordCacheOperation(i%2 == 0)
		metrics.RecordDatabaseQuery(50 * time.Millisecond)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics.GetAllMetrics()
	}
}

func BenchmarkPerformanceMetrics_Concurrent(b *testing.B) {
	metrics := NewPerformanceMetrics()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			metrics.RecordRequest(time.Duration(i%100)*time.Millisecond, i%2 == 0)
			metrics.RecordCacheOperation(i%2 == 0)
			metrics.RecordDatabaseQuery(time.Duration(i%50) * time.Millisecond)
			i++
		}
	})
}
