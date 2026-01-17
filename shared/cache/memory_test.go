package cache

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestMemoryCache_NewMemoryCache(t *testing.T) {
	config := MemoryConfig{
		TTL:             5 * time.Minute,
		MaxSize:         100,
		CleanupInterval: 1 * time.Minute,
	}

	cache := NewMemoryCache(config)

	if cache == nil {
		t.Fatal("Expected cache to be created")
	}

	if cache.ttl != config.TTL {
		t.Errorf("Expected TTL %v, got %v", config.TTL, cache.ttl)
	}

	if cache.maxSize != config.MaxSize {
		t.Errorf("Expected MaxSize %d, got %d", config.MaxSize, cache.maxSize)
	}
}

func TestMemoryCache_SetAndGet(t *testing.T) {
	config := MemoryConfig{
		TTL:             5 * time.Minute,
		MaxSize:         100,
		CleanupInterval: 1 * time.Minute,
	}

	cache := NewMemoryCache(config)
	ctx := context.Background()

	// Test setting and getting a value
	key := "test-key"
	value := "test-value"

	err := cache.Set(ctx, key, value, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	retrieved, err := cache.Get(ctx, key)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if retrieved != value {
		t.Errorf("Expected %s, got %s", value, retrieved)
	}
}

func TestMemoryCache_GetNonExistent(t *testing.T) {
	config := MemoryConfig{
		TTL:             5 * time.Minute,
		MaxSize:         100,
		CleanupInterval: 1 * time.Minute,
	}

	cache := NewMemoryCache(config)
	ctx := context.Background()

	_, err := cache.Get(ctx, "non-existent")
	if err != ErrCacheMiss {
		t.Errorf("Expected ErrCacheMiss, got %v", err)
	}
}

func TestMemoryCache_Expiration(t *testing.T) {
	config := MemoryConfig{
		TTL:             5 * time.Minute,
		MaxSize:         100,
		CleanupInterval: 1 * time.Minute,
	}

	cache := NewMemoryCache(config)
	ctx := context.Background()

	key := "expired-key"
	value := "expired-value"

	// Set with very short expiration
	err := cache.Set(ctx, key, value, 1*time.Millisecond)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Wait for expiration
	time.Sleep(10 * time.Millisecond)

	_, err = cache.Get(ctx, key)
	if err != ErrCacheMiss {
		t.Errorf("Expected ErrCacheMiss due to expiration, got %v", err)
	}
}

func TestMemoryCache_Delete(t *testing.T) {
	config := MemoryConfig{
		TTL:             5 * time.Minute,
		MaxSize:         100,
		CleanupInterval: 1 * time.Minute,
	}

	cache := NewMemoryCache(config)
	ctx := context.Background()

	key := "delete-key"
	value := "delete-value"

	// Set value
	err := cache.Set(ctx, key, value, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify it exists
	_, err = cache.Get(ctx, key)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Delete it
	err = cache.Delete(ctx, key)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify it's gone
	_, err = cache.Get(ctx, key)
	if err != ErrCacheMiss {
		t.Errorf("Expected ErrCacheMiss after delete, got %v", err)
	}
}

func TestMemoryCache_Exists(t *testing.T) {
	config := MemoryConfig{
		TTL:             5 * time.Minute,
		MaxSize:         100,
		CleanupInterval: 1 * time.Minute,
	}

	cache := NewMemoryCache(config)
	ctx := context.Background()

	key := "exists-key"
	value := "exists-value"

	// Test non-existent key
	exists, err := cache.Exists(ctx, key)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if exists {
		t.Error("Expected key to not exist")
	}

	// Set value
	err = cache.Set(ctx, key, value, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Test existing key
	exists, err = cache.Exists(ctx, key)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !exists {
		t.Error("Expected key to exist")
	}
}

func TestMemoryCache_Expire(t *testing.T) {
	config := MemoryConfig{
		TTL:             5 * time.Minute,
		MaxSize:         100,
		CleanupInterval: 1 * time.Minute,
	}

	cache := NewMemoryCache(config)
	ctx := context.Background()

	key := "expire-key"
	value := "expire-value"

	// Set value
	err := cache.Set(ctx, key, value, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Set expiration
	err = cache.Expire(ctx, key, 1*time.Millisecond)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Wait for expiration
	time.Sleep(10 * time.Millisecond)

	_, err = cache.Get(ctx, key)
	if err != ErrCacheMiss {
		t.Errorf("Expected ErrCacheMiss after expiration, got %v", err)
	}
}

func TestMemoryCache_TTL(t *testing.T) {
	config := MemoryConfig{
		TTL:             5 * time.Minute,
		MaxSize:         100,
		CleanupInterval: 1 * time.Minute,
	}

	cache := NewMemoryCache(config)
	ctx := context.Background()

	key := "ttl-key"
	value := "ttl-value"

	// Set value with TTL
	err := cache.Set(ctx, key, value, 1*time.Second)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check TTL
	ttl, err := cache.TTL(ctx, key)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if ttl <= 0 || ttl > 1*time.Second {
		t.Errorf("Expected TTL between 0 and 1s, got %v", ttl)
	}
}

func TestMemoryCache_Increment(t *testing.T) {
	config := MemoryConfig{
		TTL:             5 * time.Minute,
		MaxSize:         100,
		CleanupInterval: 1 * time.Minute,
	}

	cache := NewMemoryCache(config)
	ctx := context.Background()

	key := "increment-key"

	// Test increment on non-existent key
	val, err := cache.Increment(ctx, key)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if val != 1 {
		t.Errorf("Expected 1, got %d", val)
	}

	// Test increment on existing key
	val, err = cache.Increment(ctx, key)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if val != 2 {
		t.Errorf("Expected 2, got %d", val)
	}
}

func TestMemoryCache_Decrement(t *testing.T) {
	config := MemoryConfig{
		TTL:             5 * time.Minute,
		MaxSize:         100,
		CleanupInterval: 1 * time.Minute,
	}

	cache := NewMemoryCache(config)
	ctx := context.Background()

	key := "decrement-key"

	// Test decrement on non-existent key
	val, err := cache.Decrement(ctx, key)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if val != -1 {
		t.Errorf("Expected -1, got %d", val)
	}

	// Set a numeric value first
	err = cache.Set(ctx, key, 5, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Test decrement on existing key
	val, err = cache.Decrement(ctx, key)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if val != 4 {
		t.Errorf("Expected 4, got %d", val)
	}
}

func TestMemoryCache_Keys(t *testing.T) {
	config := MemoryConfig{
		TTL:             5 * time.Minute,
		MaxSize:         100,
		CleanupInterval: 1 * time.Minute,
	}

	cache := NewMemoryCache(config)
	ctx := context.Background()

	// Set multiple keys
	keys := []string{"key1", "key2", "key3"}
	for _, key := range keys {
		err := cache.Set(ctx, key, "value", 0)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	}

	// Get all keys
	retrievedKeys, err := cache.Keys(ctx, "*")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(retrievedKeys) != len(keys) {
		t.Errorf("Expected %d keys, got %d", len(keys), len(retrievedKeys))
	}
}

func TestMemoryCache_Flush(t *testing.T) {
	config := MemoryConfig{
		TTL:             5 * time.Minute,
		MaxSize:         100,
		CleanupInterval: 1 * time.Minute,
	}

	cache := NewMemoryCache(config)
	ctx := context.Background()

	// Set some values
	err := cache.Set(ctx, "key1", "value1", 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = cache.Set(ctx, "key2", "value2", 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Flush cache
	err = cache.Flush(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify keys are gone
	_, err = cache.Get(ctx, "key1")
	if err != ErrCacheMiss {
		t.Errorf("Expected ErrCacheMiss after flush, got %v", err)
	}

	_, err = cache.Get(ctx, "key2")
	if err != ErrCacheMiss {
		t.Errorf("Expected ErrCacheMiss after flush, got %v", err)
	}
}

func TestMemoryCache_Health(t *testing.T) {
	config := MemoryConfig{
		TTL:             5 * time.Minute,
		MaxSize:         100,
		CleanupInterval: 1 * time.Minute,
	}

	cache := NewMemoryCache(config)
	ctx := context.Background()

	err := cache.Health(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestMemoryCache_ConcurrentAccess(t *testing.T) {
	config := MemoryConfig{
		TTL:             5 * time.Minute,
		MaxSize:         100,
		CleanupInterval: 1 * time.Minute,
	}

	cache := NewMemoryCache(config)
	ctx := context.Background()

	// Test concurrent writes
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			key := fmt.Sprintf("concurrent-key-%d", i)
			value := fmt.Sprintf("concurrent-value-%d", i)
			err := cache.Set(ctx, key, value, 0)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all values were set
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("concurrent-key-%d", i)
		expectedValue := fmt.Sprintf("concurrent-value-%d", i)
		value, err := cache.Get(ctx, key)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if value != expectedValue {
			t.Errorf("Expected %s, got %s", expectedValue, value)
		}
	}
}

func TestMemoryCache_MaxSize(t *testing.T) {
	config := MemoryConfig{
		TTL:             5 * time.Minute,
		MaxSize:         3,
		CleanupInterval: 1 * time.Minute,
	}

	cache := NewMemoryCache(config)
	ctx := context.Background()

	// Fill cache to max size
	for i := 0; i < 3; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		err := cache.Set(ctx, key, value, 0)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	}

	// Add one more item (should trigger eviction)
	err := cache.Set(ctx, "key-3", "value-3", 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify cache size is still within limit
	keys, err := cache.Keys(ctx, "*")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(keys) > config.MaxSize {
		t.Errorf("Expected cache size <= %d, got %d", config.MaxSize, len(keys))
	}
}

// Benchmark tests
func BenchmarkMemoryCache_Set(b *testing.B) {
	config := MemoryConfig{
		TTL:             5 * time.Minute,
		MaxSize:         10000,
		CleanupInterval: 1 * time.Minute,
	}

	cache := NewMemoryCache(config)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("bench-key-%d", i)
		value := fmt.Sprintf("bench-value-%d", i)
		cache.Set(ctx, key, value, 0)
	}
}

func BenchmarkMemoryCache_Get(b *testing.B) {
	config := MemoryConfig{
		TTL:             5 * time.Minute,
		MaxSize:         10000,
		CleanupInterval: 1 * time.Minute,
	}

	cache := NewMemoryCache(config)
	ctx := context.Background()

	// Pre-populate cache
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("bench-key-%d", i)
		value := fmt.Sprintf("bench-value-%d", i)
		cache.Set(ctx, key, value, 0)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("bench-key-%d", i%1000)
		cache.Get(ctx, key)
	}
}

func BenchmarkMemoryCache_Concurrent(b *testing.B) {
	config := MemoryConfig{
		TTL:             5 * time.Minute,
		MaxSize:         10000,
		CleanupInterval: 1 * time.Minute,
	}

	cache := NewMemoryCache(config)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("concurrent-bench-key-%d", i)
			value := fmt.Sprintf("concurrent-bench-value-%d", i)
			cache.Set(ctx, key, value, 0)
			cache.Get(ctx, key)
			i++
		}
	})
}
