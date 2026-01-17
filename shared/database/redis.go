package database

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/usc-platform/shared/config"
)

// RedisHealthChecker implements health checking for Redis
type RedisHealthChecker struct {
	client RedisClient
}

// NewRedisHealthChecker creates a new Redis health checker
func NewRedisHealthChecker(client RedisClient) *RedisHealthChecker {
	return &RedisHealthChecker{client: client}
}

// Check performs a health check on Redis
func (h *RedisHealthChecker) Check(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := h.client.Ping(ctx); err != nil {
		return fmt.Errorf("Redis ping failed: %w", err)
	}

	return nil
}

// Name returns the name of the health checker
func (h *RedisHealthChecker) Name() string {
	return "redis"
}

// Description returns the description of the health checker
func (h *RedisHealthChecker) Description() string {
	return "Redis database health check"
}

// RedisConnection represents a Redis connection
type RedisConnection struct {
	client *redis.Client
}

// RedisWrapper wraps redis.Client to implement our RedisClient interface
type RedisWrapper struct {
	client *redis.Client
}

// Ping tests the connection
func (r *RedisWrapper) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// Set sets a key-value pair
func (r *RedisWrapper) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get gets a value by key
func (r *RedisWrapper) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Del deletes keys
func (r *RedisWrapper) Del(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Del(ctx, keys...).Result()
}

// HealthCheck performs a health check
func (r *RedisWrapper) HealthCheck(ctx context.Context) error {
	return r.Ping(ctx)
}

// Exists checks if keys exist
func (r *RedisWrapper) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Exists(ctx, keys...).Result()
}

// Expire sets expiration for a key
func (r *RedisWrapper) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return r.client.Expire(ctx, key, expiration).Result()
}

// TTL returns the time to live for a key
func (r *RedisWrapper) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

// Keys returns all keys matching pattern
func (r *RedisWrapper) Keys(ctx context.Context, pattern string) ([]string, error) {
	return r.client.Keys(ctx, pattern).Result()
}

// IncrBy atomically increments a counter
func (r *RedisWrapper) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return r.client.IncrBy(ctx, key, value).Result()
}

// DecrBy atomically decrements a counter
func (r *RedisWrapper) DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	return r.client.DecrBy(ctx, key, value).Result()
}

// SetNX sets a key only if it doesn't exist
func (r *RedisWrapper) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, value, expiration).Result()
}

// MGet retrieves multiple values
func (r *RedisWrapper) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	return r.client.MGet(ctx, keys...).Result()
}

// Pipeline creates a pipeline
func (r *RedisWrapper) Pipeline() redis.Pipeliner {
	return r.client.Pipeline()
}

// Info returns Redis server information
func (r *RedisWrapper) Info(ctx context.Context, section ...string) (string, error) {
	return r.client.Info(ctx, section...).Result()
}

// LPush pushes values to the left of a list
func (r *RedisWrapper) LPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return r.client.LPush(ctx, key, values...).Result()
}

// NewRedisConnection creates a new Redis connection
func NewRedisConnection(cfg *config.Config) (*RedisConnection, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddress(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	return &RedisConnection{client: client}, nil
}

// Client returns the underlying Redis client
func (r *RedisConnection) Client() *redis.Client {
	return r.client
}

// Ping tests the connection
func (r *RedisConnection) Ping(ctx context.Context) *redis.StatusCmd {
	return r.client.Ping(ctx)
}

// Set sets a key-value pair
func (r *RedisConnection) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.client.Set(ctx, key, value, expiration)
}

// Get gets a value by key
func (r *RedisConnection) Get(ctx context.Context, key string) *redis.StringCmd {
	return r.client.Get(ctx, key)
}

// Del deletes keys
func (r *RedisConnection) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return r.client.Del(ctx, keys...)
}

// Exists checks if keys exist
func (r *RedisConnection) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	return r.client.Exists(ctx, keys...)
}

// Expire sets expiration for a key
func (r *RedisConnection) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return r.client.Expire(ctx, key, expiration)
}

// TTL returns the time to live for a key
func (r *RedisConnection) TTL(ctx context.Context, key string) *redis.DurationCmd {
	return r.client.TTL(ctx, key)
}

// Keys returns all keys matching pattern
func (r *RedisConnection) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {
	return r.client.Keys(ctx, pattern)
}

// FlushDB flushes the current database
func (r *RedisConnection) FlushDB(ctx context.Context) *redis.StatusCmd {
	return r.client.FlushDB(ctx)
}

// FlushAll flushes all databases
func (r *RedisConnection) FlushAll(ctx context.Context) *redis.StatusCmd {
	return r.client.FlushAll(ctx)
}

// Info returns Redis server information
func (r *RedisConnection) Info(ctx context.Context, section ...string) *redis.StringCmd {
	return r.client.Info(ctx, section...)
}

// Close closes the connection
func (r *RedisConnection) Close() error {
	return r.client.Close()
}

// HealthCheck performs a health check
func (r *RedisConnection) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := r.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis ping failed: %w", err)
	}

	return nil
}

// Stats returns connection pool statistics
func (r *RedisConnection) Stats() *redis.PoolStats {
	return r.client.PoolStats()
}

// initializeRedis initializes Redis connection
func (m *DatabaseManager) initializeRedis() error {
	client := redis.NewClient(&redis.Options{
		Addr:     m.config.GetRedisAddress(),
		Password: m.config.Redis.Password,
		DB:       m.config.Redis.DB,
		PoolSize: m.config.Redis.PoolSize,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return fmt.Errorf("failed to ping Redis: %w", err)
	}

	m.redis = &RedisWrapper{client: client}
	return nil
}

// RedisCache provides caching functionality
type RedisCache struct {
	client *redis.Client
	prefix string
	ttl    time.Duration
}

// NewRedisCache creates a new Redis cache
func NewRedisCache(client *redis.Client, prefix string, ttl time.Duration) *RedisCache {
	return &RedisCache{
		client: client,
		prefix: prefix,
		ttl:    ttl,
	}
}

// Get retrieves a value from cache
func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	result := c.client.Get(ctx, c.prefix+key)
	if err := result.Err(); err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("key not found")
		}
		return "", err
	}
	return result.Val(), nil
}

// Set stores a value in cache
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}) error {
	return c.client.Set(ctx, c.prefix+key, value, c.ttl).Err()
}

// Delete removes a value from cache
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, c.prefix+key).Err()
}

// Exists checks if a key exists in cache
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	result := c.client.Exists(ctx, c.prefix+key)
	return result.Val() > 0, result.Err()
}

// Clear removes all keys with the prefix
func (c *RedisCache) Clear(ctx context.Context) error {
	pattern := c.prefix + "*"
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return c.client.Del(ctx, keys...).Err()
	}

	return nil
}

// IsRedisError checks if an error is Redis-specific
func IsRedisError(err error) bool {
	if err == nil {
		return false
	}

	// Check for Redis-specific error types
	if _, ok := err.(redis.Error); ok {
		return true
	}

	// Check for common Redis error patterns
	errStr := err.Error()
	redisErrors := []string{
		"redis:",
		"Redis",
		"NOAUTH",
		"WRONGTYPE",
		"ERR",
		"MOVED",
		"ASK",
	}

	for _, pattern := range redisErrors {
		if contains(errStr, pattern) {
			return true
		}
	}

	return false
}

// IsRedisNil checks if the error is Redis nil (key not found)
func IsRedisNil(err error) bool {
	return err == redis.Nil
}
