package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewRateLimitManager(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerSecond: 10.0,
		Burst:             5,
		Interval:          time.Second,
	}

	manager := NewRateLimitManager(config)

	if manager == nil {
		t.Error("Expected rate limit manager, got nil")
		return
	}

	if manager.config.RequestsPerSecond != config.RequestsPerSecond {
		t.Errorf("Expected %f requests per second, got %f", config.RequestsPerSecond, manager.config.RequestsPerSecond)
	}
}

func TestRateLimitManager_GetLimiter(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerSecond: 10.0,
		Burst:             5,
		Interval:          time.Second,
	}

	manager := NewRateLimitManager(config)

	// Test getting a new limiter
	limiter := manager.GetLimiter("test-key")
	if limiter == nil {
		t.Error("Expected limiter, got nil")
	}

	// Test getting the same limiter
	limiter2 := manager.GetLimiter("test-key")
	if limiter != limiter2 {
		t.Error("Expected same limiter instance")
	}

	// Test getting a different limiter
	limiter3 := manager.GetLimiter("different-key")
	if limiter == limiter3 {
		t.Error("Expected different limiter instance")
	}
}

func TestRateLimitManager_Allow(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerSecond: 2.0,
		Burst:             1,
		Interval:          time.Second,
	}

	manager := NewRateLimitManager(config)

	// First request should be allowed
	if !manager.Allow("test-key") {
		t.Error("Expected first request to be allowed")
	}

	// Second request should be rate limited
	if manager.Allow("test-key") {
		t.Error("Expected second request to be rate limited")
	}

	// Wait for rate limit to reset
	time.Sleep(600 * time.Millisecond)

	// Third request should be allowed again
	if !manager.Allow("test-key") {
		t.Error("Expected third request to be allowed after reset")
	}
}

func TestHTTPRateLimitMiddleware(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerSecond: 2.0,
		Burst:             1,
		Interval:          time.Second,
		KeyFunc: func(r *http.Request) string {
			return "test-key"
		},
	}

	middleware := NewHTTPRateLimitMiddleware(config)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Test first request
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	middleware.Middleware()(handler).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	// Test second request (should be rate limited)
	req2 := httptest.NewRequest("GET", "/test", nil)
	rr2 := httptest.NewRecorder()
	middleware.Middleware()(handler).ServeHTTP(rr2, req2)

	if rr2.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status %d, got %d", http.StatusTooManyRequests, rr2.Code)
	}
}

func TestRateLimitManager_AllowN(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerSecond: 10.0,
		Burst:             5,
		Interval:          time.Second,
	}

	manager := NewRateLimitManager(config)

	// Test allowing multiple requests
	if !manager.AllowN("test-key", 3) {
		t.Error("Expected 3 requests to be allowed")
	}

	// Test exceeding burst limit
	if manager.AllowN("test-key", 5) {
		t.Error("Expected 5 requests to be rate limited")
	}
}

func TestRateLimitManager_Wait(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerSecond: 1.0,
		Burst:             1,
		Interval:          time.Second,
	}

	manager := NewRateLimitManager(config)

	// First request should be allowed immediately
	ctx := context.Background()
	err := manager.Wait(ctx, "test-key")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Second request should wait
	ctx2, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	err = manager.Wait(ctx2, "test-key")
	if err == nil {
		t.Error("Expected timeout error")
	}
}
