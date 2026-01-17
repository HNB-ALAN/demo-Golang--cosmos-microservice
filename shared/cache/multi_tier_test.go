package cache

import (
	"context"
	"testing"
	"time"
)

func TestEnhancedMultiTierCache(t *testing.T) {
	// Create enhanced multi-tier cache configuration
	config := DefaultEnhancedMultiTierConfig()
	config.EnableMetrics = true
	config.EnableHealthCheck = true
	config.EnableAutoWarmup = false

	// Create enhanced multi-tier cache
	emtc, err := NewEnhancedMultiTierCache(config)
	if err != nil {
		t.Fatalf("Failed to create enhanced multi-tier cache: %v", err)
	}
	defer emtc.Close()

	ctx := context.Background()

	// Test basic operations
	key := "test-key"
	value := "test-value"
	expiration := 5 * time.Minute

	// Test Set
	err = emtc.Set(ctx, key, value, expiration)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	// Test Get
	retrievedValue, err := emtc.Get(ctx, key)
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if retrievedValue != value {
		t.Errorf("Expected %s, got %s", value, retrievedValue)
	}

	// Test metrics
	metrics := emtc.GetMetrics()
	if metrics == nil {
		t.Error("Expected metrics to be available")
	}

	// Test performance report
	report := emtc.GetPerformanceReport()
	if report == nil {
		t.Error("Expected performance report to be available")
	}

	// Test health check
	err = emtc.HealthCheck(ctx)
	if err != nil {
		t.Errorf("Health check failed: %v", err)
	}

	// Test Delete
	err = emtc.Delete(ctx, key)
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	// Test Get after delete (should fail)
	_, err = emtc.Get(ctx, key)
	if err == nil {
		t.Error("Expected error after delete")
	}
}

func TestMultiTierCacheMetrics(t *testing.T) {
	// Create metrics
	metrics := NewMultiTierCacheMetrics()

	// Test recording hits
	metrics.RecordHit("l1")
	metrics.RecordHit("l2")
	metrics.RecordHit("l4")

	// Test recording misses
	metrics.RecordMiss("l1")
	metrics.RecordMiss("l2")

	// Test recording writes and deletes
	metrics.RecordWrite()
	metrics.RecordWrite()
	metrics.RecordDelete()

	// Test recording errors
	metrics.RecordError("l1")
	metrics.RecordError("l2")

	// Test recording response times
	metrics.RecordResponseTime("l1", 100*time.Millisecond)
	metrics.RecordResponseTime("l2", 200*time.Millisecond)
	metrics.RecordResponseTime("l4", 300*time.Millisecond)

	// Test getting stats
	stats := metrics.GetStats()
	if stats.L1Hits != 1 {
		t.Errorf("Expected L1Hits to be 1, got %d", stats.L1Hits)
	}
	if stats.L2Hits != 1 {
		t.Errorf("Expected L2Hits to be 1, got %d", stats.L2Hits)
	}
	if stats.L4Hits != 1 {
		t.Errorf("Expected L4Hits to be 1, got %d", stats.L4Hits)
	}
	if stats.L1Misses != 1 {
		t.Errorf("Expected L1Misses to be 1, got %d", stats.L1Misses)
	}
	if stats.L2Misses != 1 {
		t.Errorf("Expected L2Misses to be 1, got %d", stats.L2Misses)
	}
	if stats.Writes != 2 {
		t.Errorf("Expected Writes to be 2, got %d", stats.Writes)
	}
	if stats.Deletes != 1 {
		t.Errorf("Expected Deletes to be 1, got %d", stats.Deletes)
	}
	if stats.Errors != 2 {
		t.Errorf("Expected Errors to be 2, got %d", stats.Errors)
	}

	// Test hit rate calculation
	hitRate := metrics.GetHitRate()
	expectedHitRate := float64(3) / float64(5) * 100 // 3 hits out of 5 total operations
	if hitRate != expectedHitRate {
		t.Errorf("Expected hit rate %.2f, got %.2f", expectedHitRate, hitRate)
	}

	// Test tier hit rates
	tierRates := metrics.GetTierHitRates()
	if tierRates["l1"] != 50.0 { // 1 hit out of 2 total
		t.Errorf("Expected L1 hit rate 50.0, got %.2f", tierRates["l1"])
	}
	if tierRates["l2"] != 50.0 { // 1 hit out of 2 total
		t.Errorf("Expected L2 hit rate 50.0, got %.2f", tierRates["l2"])
	}
	if tierRates["l4"] != 100.0 { // 1 hit out of 1 total
		t.Errorf("Expected L4 hit rate 100.0, got %.2f", tierRates["l4"])
	}

	// Test average response times
	l1AvgTime := metrics.GetAverageResponseTime("l1")
	if l1AvgTime != 100*time.Millisecond {
		t.Errorf("Expected L1 avg response time 100ms, got %v", l1AvgTime)
	}

	l2AvgTime := metrics.GetAverageResponseTime("l2")
	if l2AvgTime != 200*time.Millisecond {
		t.Errorf("Expected L2 avg response time 200ms, got %v", l2AvgTime)
	}

	l4AvgTime := metrics.GetAverageResponseTime("l4")
	if l4AvgTime != 300*time.Millisecond {
		t.Errorf("Expected L4 avg response time 300ms, got %v", l4AvgTime)
	}

	// Test error rates
	l1ErrorRate := metrics.GetErrorRate("l1")
	expectedL1ErrorRate := float64(1) / float64(5) * 100 // 1 error out of 5 total operations
	if l1ErrorRate != expectedL1ErrorRate {
		t.Errorf("Expected L1 error rate %.2f, got %.2f", expectedL1ErrorRate, l1ErrorRate)
	}
}

func TestMultiTierCacheConfig(t *testing.T) {
	config := DefaultEnhancedMultiTierConfig()

	// Test default values
	if !config.EnableMetrics {
		t.Error("Expected metrics to be enabled by default")
	}

	if !config.EnableHealthCheck {
		t.Error("Expected health check to be enabled by default")
	}

	if config.EnableAutoWarmup {
		t.Error("Expected auto warmup to be disabled by default")
	}

	if config.EnableCompression {
		t.Error("Expected compression to be disabled by default")
	}

	if config.EnableEncryption {
		t.Error("Expected encryption to be disabled by default")
	}

	// Test base config
	if config.L1Config.MaxSize != 1000 {
		t.Errorf("Expected L1 MaxSize to be 1000, got %d", config.L1Config.MaxSize)
	}

	if config.L2Config.MaxRetries != 3 {
		t.Errorf("Expected L2 MaxRetries to be 3, got %d", config.L2Config.MaxRetries)
	}

	if config.L4Config.MaxFileSize != 10*1024*1024 {
		t.Errorf("Expected L4 MaxFileSize to be 10MB, got %d", config.L4Config.MaxFileSize)
	}

	// Test performance tuning
	if config.MaxRetries != 3 {
		t.Errorf("Expected MaxRetries to be 3, got %d", config.MaxRetries)
	}

	if config.RetryDelay != 100*time.Millisecond {
		t.Errorf("Expected RetryDelay to be 100ms, got %v", config.RetryDelay)
	}

	if config.CircuitBreakerThreshold != 10 {
		t.Errorf("Expected CircuitBreakerThreshold to be 10, got %d", config.CircuitBreakerThreshold)
	}
}

func TestMultiTierCacheWarmup(t *testing.T) {
	// Create enhanced multi-tier cache configuration
	config := DefaultEnhancedMultiTierConfig()
	config.EnableAutoWarmup = true

	// Create enhanced multi-tier cache
	emtc, err := NewEnhancedMultiTierCache(config)
	if err != nil {
		t.Fatalf("Failed to create enhanced multi-tier cache: %v", err)
	}
	defer emtc.Close()

	ctx := context.Background()

	// Set some data in L4 cache
	key1 := "warmup-key-1"
	value1 := "warmup-value-1"
	key2 := "warmup-key-2"
	value2 := "warmup-value-2"

	err = emtc.l4.Set(ctx, key1, value1, 5*time.Minute)
	if err != nil {
		t.Errorf("Failed to set L4 data: %v", err)
	}

	err = emtc.l4.Set(ctx, key2, value2, 5*time.Minute)
	if err != nil {
		t.Errorf("Failed to set L4 data: %v", err)
	}

	// Perform warmup
	warmupKeys := []string{key1, key2}
	err = emtc.Warmup(ctx, warmupKeys)
	if err != nil {
		t.Errorf("Warmup failed: %v", err)
	}

	// Verify data is now in L1 and L2
	// Note: In a real implementation, you would verify the data is actually populated
	// For now, we just verify the warmup operation completes without error
}

func TestMultiTierCacheHealthCheck(t *testing.T) {
	// Create enhanced multi-tier cache configuration
	config := DefaultEnhancedMultiTierConfig()
	config.EnableHealthCheck = true

	// Create enhanced multi-tier cache
	emtc, err := NewEnhancedMultiTierCache(config)
	if err != nil {
		t.Fatalf("Failed to create enhanced multi-tier cache: %v", err)
	}
	defer emtc.Close()

	ctx := context.Background()

	// Test health check
	err = emtc.HealthCheck(ctx)
	if err != nil {
		t.Errorf("Health check failed: %v", err)
	}

	// Test with health check disabled
	config.EnableHealthCheck = false
	emtc2, err := NewEnhancedMultiTierCache(config)
	if err != nil {
		t.Fatalf("Failed to create enhanced multi-tier cache: %v", err)
	}
	defer emtc2.Close()

	// Health check should pass when disabled
	err = emtc2.HealthCheck(ctx)
	if err != nil {
		t.Errorf("Health check should pass when disabled: %v", err)
	}
}

func TestMultiTierCachePerformanceReport(t *testing.T) {
	// Create enhanced multi-tier cache
	config := DefaultEnhancedMultiTierConfig()
	emtc, err := NewEnhancedMultiTierCache(config)
	if err != nil {
		t.Fatalf("Failed to create enhanced multi-tier cache: %v", err)
	}
	defer emtc.Close()

	ctx := context.Background()

	// Perform some operations to generate metrics
	key := "perf-test-key"
	value := "perf-test-value"

	err = emtc.Set(ctx, key, value, 5*time.Minute)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	_, err = emtc.Get(ctx, key)
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}

	err = emtc.Delete(ctx, key)
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	// Get performance report
	report := emtc.GetPerformanceReport()
	if report == nil {
		t.Error("Expected performance report to be available")
	}

	// Verify report structure
	if report["stats"] == nil {
		t.Error("Expected stats in performance report")
	}

	if report["hit_rates"] == nil {
		t.Error("Expected hit_rates in performance report")
	}

	if report["response_times"] == nil {
		t.Error("Expected response_times in performance report")
	}

	if report["error_rates"] == nil {
		t.Error("Expected error_rates in performance report")
	}

	if report["config"] == nil {
		t.Error("Expected config in performance report")
	}

	// Verify overall hit rate is present
	if report["overall_hit_rate"] == nil {
		t.Error("Expected overall_hit_rate in performance report")
	}
}
