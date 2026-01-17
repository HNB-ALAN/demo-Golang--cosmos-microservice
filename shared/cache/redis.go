// Package cache provides caching utilities for USC platform services.
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Cache defines the interface for cache operations
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	TTL(ctx context.Context, key string) (time.Duration, error)
	Increment(ctx context.Context, key string) (int64, error)
	Decrement(ctx context.Context, key string) (int64, error)
	Keys(ctx context.Context, pattern string) ([]string, error)
	Flush(ctx context.Context) error
	Health(ctx context.Context) error
}

// RedisCache represents a Redis-based cache implementation
type RedisCache struct {
	client *redis.Client
	prefix string
	ttl    time.Duration
}

// RedisConfig represents Redis cache configuration
type RedisConfig struct {
	Host     string        `mapstructure:"host"`
	Port     int           `mapstructure:"port"`
	Password string        `mapstructure:"password"`
	DB       int           `mapstructure:"db"`
	PoolSize int           `mapstructure:"pool_size"`
	Prefix   string        `mapstructure:"prefix"`
	TTL      time.Duration `mapstructure:"ttl"`
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(client *redis.Client, config RedisConfig) *RedisCache {
	return &RedisCache{
		client: client,
		prefix: config.Prefix,
		ttl:    config.TTL,
	}
}

// Get retrieves a value from cache
func (r *RedisCache) Get(ctx context.Context, key string) (interface{}, error) {
	fullKey := r.buildKey(key)

	val, err := r.client.Get(ctx, fullKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrCacheMiss
		}
		return nil, fmt.Errorf("failed to get from cache: %w", err)
	}

	// Try to unmarshal as JSON first
	var result interface{}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		// If JSON unmarshaling fails, return as string
		return val, nil
	}

	return result, nil
}

// Set stores a value in cache
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl ...time.Duration) error {
	fullKey := r.buildKey(key)

	// Marshal value to JSON
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	// Use provided TTL or default TTL
	duration := r.ttl
	if len(ttl) > 0 {
		duration = ttl[0]
	}

	err = r.client.Set(ctx, fullKey, data, duration).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// Delete removes a value from cache
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	fullKey := r.buildKey(key)

	err := r.client.Del(ctx, fullKey).Err()
	if err != nil {
		return fmt.Errorf("failed to delete from cache: %w", err)
	}

	return nil
}

// Exists checks if a key exists in cache
func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	fullKey := r.buildKey(key)

	count, err := r.client.Exists(ctx, fullKey).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence: %w", err)
	}

	return count > 0, nil
}

// Expire sets expiration time for a key
func (r *RedisCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	fullKey := r.buildKey(key)

	err := r.client.Expire(ctx, fullKey, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set expiration: %w", err)
	}

	return nil
}

// TTL returns the time to live for a key
func (r *RedisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	fullKey := r.buildKey(key)

	ttl, err := r.client.TTL(ctx, fullKey).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get TTL: %w", err)
	}

	return ttl, nil
}

// Increment increments a numeric value
func (r *RedisCache) Increment(ctx context.Context, key string, delta int64) (int64, error) {
	fullKey := r.buildKey(key)

	val, err := r.client.IncrBy(ctx, fullKey, delta).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment: %w", err)
	}

	return val, nil
}

// Decrement decrements a numeric value
func (r *RedisCache) Decrement(ctx context.Context, key string, delta int64) (int64, error) {
	fullKey := r.buildKey(key)

	val, err := r.client.DecrBy(ctx, fullKey, delta).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to decrement: %w", err)
	}

	return val, nil
}

// GetMultiple retrieves multiple values from cache
func (r *RedisCache) GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, error) {
	if len(keys) == 0 {
		return make(map[string]interface{}), nil
	}

	fullKeys := make([]string, len(keys))
	for i, key := range keys {
		fullKeys[i] = r.buildKey(key)
	}

	vals, err := r.client.MGet(ctx, fullKeys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get multiple values: %w", err)
	}

	result := make(map[string]interface{})
	for i, val := range vals {
		if val != nil {
			// Try to unmarshal as JSON first
			var parsed interface{}
			if err := json.Unmarshal([]byte(val.(string)), &parsed); err != nil {
				// If JSON unmarshaling fails, use as string
				result[keys[i]] = val
			} else {
				result[keys[i]] = parsed
			}
		}
	}

	return result, nil
}

// SetMultiple stores multiple values in cache
func (r *RedisCache) SetMultiple(ctx context.Context, values map[string]interface{}, ttl ...time.Duration) error {
	if len(values) == 0 {
		return nil
	}

	pipe := r.client.Pipeline()

	// Use provided TTL or default TTL
	duration := r.ttl
	if len(ttl) > 0 {
		duration = ttl[0]
	}

	for key, value := range values {
		fullKey := r.buildKey(key)

		// Marshal value to JSON
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
		}

		pipe.Set(ctx, fullKey, data, duration)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to set multiple values: %w", err)
	}

	return nil
}

// DeleteMultiple removes multiple values from cache
func (r *RedisCache) DeleteMultiple(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	fullKeys := make([]string, len(keys))
	for i, key := range keys {
		fullKeys[i] = r.buildKey(key)
	}

	err := r.client.Del(ctx, fullKeys...).Err()
	if err != nil {
		return fmt.Errorf("failed to delete multiple values: %w", err)
	}

	return nil
}

// Clear clears all keys with the cache prefix
func (r *RedisCache) Clear(ctx context.Context) error {
	pattern := r.prefix + "*"

	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get keys: %w", err)
	}

	if len(keys) > 0 {
		err = r.client.Del(ctx, keys...).Err()
		if err != nil {
			return fmt.Errorf("failed to clear cache: %w", err)
		}
	}

	return nil
}

// Keys returns all keys matching a pattern
func (r *RedisCache) Keys(ctx context.Context, pattern string) ([]string, error) {
	fullPattern := r.buildKey(pattern)

	keys, err := r.client.Keys(ctx, fullPattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys: %w", err)
	}

	// Remove prefix from keys
	result := make([]string, len(keys))
	for i, key := range keys {
		result[i] = r.removePrefix(key)
	}

	return result, nil
}

// Size returns the number of keys in cache
func (r *RedisCache) Size(ctx context.Context) (int64, error) {
	pattern := r.prefix + "*"

	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get keys: %w", err)
	}

	return int64(len(keys)), nil
}

// Health checks the health of the Redis connection
func (r *RedisCache) Health(ctx context.Context) error {
	_, err := r.client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("redis health check failed: %w", err)
	}
	return nil
}

// buildKey builds a full key with prefix
func (r *RedisCache) buildKey(key string) string {
	if r.prefix == "" {
		return key
	}
	return r.prefix + ":" + key
}

// removePrefix removes the prefix from a key
func (r *RedisCache) removePrefix(key string) string {
	if r.prefix == "" {
		return key
	}
	prefix := r.prefix + ":"
	if len(key) > len(prefix) && key[:len(prefix)] == prefix {
		return key[len(prefix):]
	}
	return key
}

// GetWithFallback retrieves a value from cache, with fallback function
func (r *RedisCache) GetWithFallback(ctx context.Context, key string, fallback func() (interface{}, error), ttl ...time.Duration) (interface{}, error) {
	// Try to get from cache first
	val, err := r.Get(ctx, key)
	if err == nil {
		return val, nil
	}

	// If cache miss, call fallback function
	result, err := fallback()
	if err != nil {
		return nil, err
	}

	// Store result in cache
	if err := r.Set(ctx, key, result, ttl...); err != nil {
		// Log error but don't fail the operation
		// In a real implementation, you would use a logger here
	}

	return result, nil
}

// SetIfNotExists sets a value only if the key doesn't exist
func (r *RedisCache) SetIfNotExists(ctx context.Context, key string, value interface{}, ttl ...time.Duration) (bool, error) {
	fullKey := r.buildKey(key)

	// Marshal value to JSON
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}

	// Use provided TTL or default TTL
	duration := r.ttl
	if len(ttl) > 0 {
		duration = ttl[0]
	}

	result, err := r.client.SetNX(ctx, fullKey, data, duration).Result()
	if err != nil {
		return false, fmt.Errorf("failed to set if not exists: %w", err)
	}

	return result, nil
}

// GetAndSet atomically gets and sets a value
func (r *RedisCache) GetAndSet(ctx context.Context, key string, value interface{}, ttl ...time.Duration) (interface{}, error) {
	fullKey := r.buildKey(key)

	// Marshal value to JSON
	data, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal value: %w", err)
	}

	// Use provided TTL or default TTL
	duration := r.ttl
	if len(ttl) > 0 {
		duration = ttl[0]
	}

	// Use GETSET command
	oldVal, err := r.client.GetSet(ctx, fullKey, data).Result()
	if err != nil {
		if err == redis.Nil {
			// Key didn't exist, set expiration
			r.client.Expire(ctx, fullKey, duration)
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get and set: %w", err)
	}

	// Set expiration for new value
	r.client.Expire(ctx, fullKey, duration)

	// Try to unmarshal old value
	var result interface{}
	if err := json.Unmarshal([]byte(oldVal), &result); err != nil {
		return oldVal, nil
	}

	return result, nil
}
