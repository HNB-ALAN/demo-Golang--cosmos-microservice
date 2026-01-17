// Package cache provides caching utilities for USC platform services.
package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// Cache errors
var (
	ErrCacheMiss   = errors.New("cache miss")
	ErrInvalidType = errors.New("invalid type")
)

// MemoryCache represents an in-memory cache implementation
type MemoryCache struct {
	mu      sync.RWMutex
	items   map[string]*cacheItem
	ttl     time.Duration
	maxSize int
	cleanup *time.Ticker
	stop    chan struct{}
}

// cacheItem represents a cached item
type cacheItem struct {
	value       interface{}
	expiresAt   time.Time
	createdAt   time.Time
	accessCount int64
	lastAccess  time.Time
}

// MemoryConfig represents memory cache configuration
type MemoryConfig struct {
	TTL             time.Duration `mapstructure:"ttl"`
	MaxSize         int           `mapstructure:"max_size"`
	CleanupInterval time.Duration `mapstructure:"cleanup_interval"`
}

// NewMemoryCache creates a new in-memory cache instance
func NewMemoryCache(config MemoryConfig) *MemoryCache {
	cache := &MemoryCache{
		items:   make(map[string]*cacheItem),
		ttl:     config.TTL,
		maxSize: config.MaxSize,
		stop:    make(chan struct{}),
	}

	// Start cleanup goroutine
	if config.CleanupInterval > 0 {
		cache.cleanup = time.NewTicker(config.CleanupInterval)
		go cache.cleanupExpired()
	}

	return cache
}

// Get retrieves a value from cache
func (m *MemoryCache) Get(ctx context.Context, key string) (string, error) {
	m.mu.RLock()
	item, exists := m.items[key]
	if !exists {
		m.mu.RUnlock()
		return "", ErrCacheMiss
	}

	// Check if item is expired
	if time.Now().After(item.expiresAt) {
		m.mu.RUnlock()
		// Item is expired, remove it with write lock
		m.mu.Lock()
		// Double-check if item still exists and is expired
		if item, exists := m.items[key]; exists && time.Now().After(item.expiresAt) {
			delete(m.items, key)
		}
		m.mu.Unlock()
		return "", ErrCacheMiss
	}

	// Update access statistics (safe to do under read lock)
	item.accessCount++
	item.lastAccess = time.Now()
	m.mu.RUnlock()

	// Convert value to string
	if str, ok := item.value.(string); ok {
		return str, nil
	}
	return fmt.Sprintf("%v", item.value), nil
}

// Set stores a value in cache
func (m *MemoryCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Use provided TTL or default TTL
	duration := expiration
	if duration == 0 {
		duration = m.ttl
	}

	// Check if we need to evict items
	if m.maxSize > 0 && len(m.items) >= m.maxSize {
		if _, exists := m.items[key]; !exists {
			m.evictLRU()
		}
	}

	// Create cache item
	item := &cacheItem{
		value:       value,
		expiresAt:   time.Now().Add(duration),
		createdAt:   time.Now(),
		accessCount: 1,
		lastAccess:  time.Now(),
	}

	m.items[key] = item
	return nil
}

// Delete removes a value from cache
func (m *MemoryCache) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.items, key)
	return nil
}

// Exists checks if a key exists in cache
func (m *MemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.items[key]
	if !exists {
		return false, nil
	}

	// Check if item is expired
	if time.Now().After(item.expiresAt) {
		return false, nil
	}

	return true, nil
}

// Expire sets expiration time for a key
func (m *MemoryCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, exists := m.items[key]
	if !exists {
		return ErrCacheMiss
	}

	item.expiresAt = time.Now().Add(ttl)
	return nil
}

// TTL returns the time to live for a key
func (m *MemoryCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.items[key]
	if !exists {
		return 0, ErrCacheMiss
	}

	// Check if item is expired
	if time.Now().After(item.expiresAt) {
		return 0, ErrCacheMiss
	}

	return time.Until(item.expiresAt), nil
}

// Increment increments a numeric value
func (m *MemoryCache) Increment(ctx context.Context, key string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, exists := m.items[key]
	if !exists {
		// Create new item with value 1
		item = &cacheItem{
			value:       int64(1),
			expiresAt:   time.Now().Add(m.ttl),
			createdAt:   time.Now(),
			accessCount: 1,
			lastAccess:  time.Now(),
		}
		m.items[key] = item
		return 1, nil
	}

	// Check if item is expired
	if time.Now().After(item.expiresAt) {
		// Item is expired, create new one
		item = &cacheItem{
			value:       int64(1),
			expiresAt:   time.Now().Add(m.ttl),
			createdAt:   time.Now(),
			accessCount: 1,
			lastAccess:  time.Now(),
		}
		m.items[key] = item
		return 1, nil
	}

	// Try to increment the value
	switch val := item.value.(type) {
	case int:
		newVal := int64(val) + 1
		item.value = newVal
		item.accessCount++
		item.lastAccess = time.Now()
		return newVal, nil
	case int64:
		newVal := val + 1
		item.value = newVal
		item.accessCount++
		item.lastAccess = time.Now()
		return newVal, nil
	case float64:
		newVal := int64(val) + 1
		item.value = newVal
		item.accessCount++
		item.lastAccess = time.Now()
		return newVal, nil
	default:
		return 0, ErrInvalidType
	}
}

// Decrement decrements a numeric value
func (m *MemoryCache) Decrement(ctx context.Context, key string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, exists := m.items[key]
	if !exists {
		// Create new item with value -1
		item = &cacheItem{
			value:       int64(-1),
			expiresAt:   time.Now().Add(m.ttl),
			createdAt:   time.Now(),
			accessCount: 1,
			lastAccess:  time.Now(),
		}
		m.items[key] = item
		return -1, nil
	}

	// Check if item is expired
	if time.Now().After(item.expiresAt) {
		// Item is expired, create new one
		item = &cacheItem{
			value:       int64(-1),
			expiresAt:   time.Now().Add(m.ttl),
			createdAt:   time.Now(),
			accessCount: 1,
			lastAccess:  time.Now(),
		}
		m.items[key] = item
		return -1, nil
	}

	// Try to decrement the value
	switch val := item.value.(type) {
	case int:
		newVal := int64(val) - 1
		item.value = newVal
		item.accessCount++
		item.lastAccess = time.Now()
		return newVal, nil
	case int64:
		newVal := val - 1
		item.value = newVal
		item.accessCount++
		item.lastAccess = time.Now()
		return newVal, nil
	case float64:
		newVal := int64(val) - 1
		item.value = newVal
		item.accessCount++
		item.lastAccess = time.Now()
		return newVal, nil
	default:
		return 0, ErrInvalidType
	}
}

// GetMultiple retrieves multiple values from cache
func (m *MemoryCache) GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, key := range keys {
		item, exists := m.items[key]
		if !exists {
			continue
		}

		// Check if item is expired
		if time.Now().After(item.expiresAt) {
			continue
		}

		result[key] = item.value
		item.accessCount++
		item.lastAccess = time.Now()
	}

	return result, nil
}

// SetMultiple stores multiple values in cache
func (m *MemoryCache) SetMultiple(ctx context.Context, values map[string]interface{}, ttl ...time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Use provided TTL or default TTL
	duration := m.ttl
	if len(ttl) > 0 {
		duration = ttl[0]
	}

	for key, value := range values {
		// Check if we need to evict items
		if m.maxSize > 0 && len(m.items) >= m.maxSize {
			if _, exists := m.items[key]; !exists {
				m.evictLRU()
			}
		}

		item := &cacheItem{
			value:       value,
			expiresAt:   time.Now().Add(duration),
			createdAt:   time.Now(),
			accessCount: 1,
			lastAccess:  time.Now(),
		}

		m.items[key] = item
	}

	return nil
}

// DeleteMultiple removes multiple values from cache
func (m *MemoryCache) DeleteMultiple(ctx context.Context, keys []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, key := range keys {
		delete(m.items, key)
	}

	return nil
}

// Clear clears all items from cache
func (m *MemoryCache) Clear(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items = make(map[string]*cacheItem)
	return nil
}

// Keys returns all keys in cache
func (m *MemoryCache) Keys(ctx context.Context, pattern string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	keys := make([]string, 0)
	now := time.Now()

	for key, item := range m.items {
		// Skip expired items
		if now.After(item.expiresAt) {
			continue
		}

		// Simple pattern matching (in real implementation, use regex)
		if pattern == "" || pattern == "*" || key == pattern {
			keys = append(keys, key)
		}
	}

	return keys, nil
}

// Size returns the number of items in cache
func (m *MemoryCache) Size(ctx context.Context) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Count only non-expired items
	count := int64(0)
	now := time.Now()

	for _, item := range m.items {
		if now.Before(item.expiresAt) {
			count++
		}
	}

	return count, nil
}

// Flush removes all items from cache
func (m *MemoryCache) Flush(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items = make(map[string]*cacheItem)
	return nil
}

// Health checks the health of the memory cache
func (m *MemoryCache) Health(ctx context.Context) error {
	// Memory cache is always healthy if it's running
	return nil
}

// GetWithFallback retrieves a value from cache, with fallback function
func (m *MemoryCache) GetWithFallback(ctx context.Context, key string, fallback func() (interface{}, error), ttl ...time.Duration) (interface{}, error) {
	// Try to get from cache first
	val, err := m.Get(ctx, key)
	if err == nil {
		return val, nil
	}

	// If cache miss, call fallback function
	result, err := fallback()
	if err != nil {
		return nil, err
	}

	// Store result in cache
	var ttlDuration time.Duration
	if len(ttl) > 0 {
		ttlDuration = ttl[0]
	}
	if err := m.Set(ctx, key, result, ttlDuration); err != nil {
		// Log error but don't fail the operation
		// In a real implementation, you would use a logger here
	}

	return result, nil
}

// SetIfNotExists sets a value only if the key doesn't exist
func (m *MemoryCache) SetIfNotExists(ctx context.Context, key string, value interface{}, ttl ...time.Duration) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if key exists and is not expired
	item, exists := m.items[key]
	if exists && time.Now().Before(item.expiresAt) {
		return false, nil
	}

	// Use provided TTL or default TTL
	duration := m.ttl
	if len(ttl) > 0 {
		duration = ttl[0]
	}

	// Create cache item
	item = &cacheItem{
		value:       value,
		expiresAt:   time.Now().Add(duration),
		createdAt:   time.Now(),
		accessCount: 1,
		lastAccess:  time.Now(),
	}

	m.items[key] = item
	return true, nil
}

// GetAndSet atomically gets and sets a value
func (m *MemoryCache) GetAndSet(ctx context.Context, key string, value interface{}, ttl ...time.Duration) (interface{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Use provided TTL or default TTL
	duration := m.ttl
	if len(ttl) > 0 {
		duration = ttl[0]
	}

	// Get old value
	var oldValue interface{}
	item, exists := m.items[key]
	if exists && time.Now().Before(item.expiresAt) {
		oldValue = item.value
	}

	// Set new value
	item = &cacheItem{
		value:       value,
		expiresAt:   time.Now().Add(duration),
		createdAt:   time.Now(),
		accessCount: 1,
		lastAccess:  time.Now(),
	}

	m.items[key] = item
	return oldValue, nil
}

// evictLRU evicts the least recently used item
func (m *MemoryCache) evictLRU() {
	if len(m.items) == 0 {
		return
	}

	var oldestKey string
	var oldestTime time.Time
	now := time.Now()

	for key, item := range m.items {
		// Skip expired items
		if now.After(item.expiresAt) {
			delete(m.items, key)
			return
		}

		if oldestKey == "" || item.lastAccess.Before(oldestTime) {
			oldestKey = key
			oldestTime = item.lastAccess
		}
	}

	if oldestKey != "" {
		delete(m.items, oldestKey)
	}
}

// cleanupExpired removes expired items periodically
func (m *MemoryCache) cleanupExpired() {
	for {
		select {
		case <-m.cleanup.C:
			m.mu.Lock()
			now := time.Now()
			for key, item := range m.items {
				if now.After(item.expiresAt) {
					delete(m.items, key)
				}
			}
			m.mu.Unlock()
		case <-m.stop:
			return
		}
	}
}

// Close stops the cache and cleans up resources
func (m *MemoryCache) Close() error {
	if m.cleanup != nil {
		m.cleanup.Stop()
	}
	close(m.stop)
	return nil
}

// GetStats returns cache statistics
func (m *MemoryCache) GetStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := map[string]interface{}{
		"total_items": len(m.items),
		"max_size":    m.maxSize,
		"default_ttl": m.ttl.String(),
	}

	// Count non-expired items
	nonExpired := 0
	totalAccessCount := int64(0)
	now := time.Now()

	for _, item := range m.items {
		if now.Before(item.expiresAt) {
			nonExpired++
			totalAccessCount += item.accessCount
		}
	}

	stats["non_expired_items"] = nonExpired
	stats["total_access_count"] = totalAccessCount

	return stats
}
