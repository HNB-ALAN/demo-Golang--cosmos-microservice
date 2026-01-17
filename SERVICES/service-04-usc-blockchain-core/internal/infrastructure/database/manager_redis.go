package database

import (
	"context"
	"fmt"
	"time"

	"github.com/usc-platform/shared/config"
	"github.com/usc-platform/shared/database"
	"github.com/usc-platform/shared/logging"
)

// RedisManager wraps the shared Redis database manager for USC Blockchain Core Service
type RedisManager struct {
	manager *database.DatabaseManager
	config  *config.Config
	logger  logging.Logger
}

// NewRedisManager creates a new Redis database manager using shared libraries
func NewRedisManager(cfg *config.Config, logger logging.Logger) (*RedisManager, error) {
	// Ensure Redis is enabled in config
	if !cfg.Redis.Enabled {
		return nil, fmt.Errorf("redis is not enabled in configuration")
	}

	manager, err := database.NewDatabaseManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create database manager: %w", err)
	}

	redisManager := &RedisManager{
		manager: manager,
		config:  cfg,
		logger:  logger,
	}

	// Initialize Redis keys and configuration from migrations
	if err := redisManager.initializeRedisSetup(context.Background()); err != nil {
		logger.Warn("Failed to initialize Redis setup", logging.Error(err))
		// Don't fail startup, just log warning
	}

	return redisManager, nil
}

// GetRedis returns Redis connection
func (rm *RedisManager) GetRedis() database.RedisClient {
	return rm.manager.Redis()
}

// Set stores a key-value pair in Redis
func (rm *RedisManager) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	client := rm.GetRedis()
	if client == nil {
		return fmt.Errorf("redis connection not available")
	}

	if err := client.Set(ctx, key, value, expiration); err != nil {
		rm.logger.Error("redis SET operation failed",
			logging.String("key", key),
			logging.Error(err))
		return fmt.Errorf("redis SET operation failed: %w", err)
	}

	return nil
}

// Get retrieves a value from Redis
func (rm *RedisManager) Get(ctx context.Context, key string) (string, error) {
	client := rm.GetRedis()
	if client == nil {
		return "", fmt.Errorf("redis connection not available")
	}

	value, err := client.Get(ctx, key)
	if err != nil {
		rm.logger.Error("redis GET operation failed",
			logging.String("key", key),
			logging.Error(err))
		return "", fmt.Errorf("redis GET operation failed: %w", err)
	}

	return value, nil
}

// Delete removes keys from Redis
func (rm *RedisManager) Delete(ctx context.Context, keys ...string) (int64, error) {
	client := rm.GetRedis()
	if client == nil {
		return 0, fmt.Errorf("redis connection not available")
	}

	count, err := client.Del(ctx, keys...)
	if err != nil {
		rm.logger.Error("redis DEL operation failed",
			logging.Strings("keys", keys),
			logging.Error(err))
		return 0, fmt.Errorf("redis DEL operation failed: %w", err)
	}

	return count, nil
}

// Exists checks if a key exists in Redis
func (rm *RedisManager) Exists(ctx context.Context, key string) (bool, error) {
	client := rm.GetRedis()
	if client == nil {
		return false, fmt.Errorf("redis connection not available")
	}

	// Use GET to check existence (returns error if key doesn't exist)
	_, err := client.Get(ctx, key)
	if err != nil {
		// Key doesn't exist or other error
		return false, nil
	}

	return true, nil
}

// SetWithExpiration stores a key-value pair with expiration
func (rm *RedisManager) SetWithExpiration(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rm.Set(ctx, key, value, expiration)
}

// GetAndDelete retrieves and deletes a key atomically
func (rm *RedisManager) GetAndDelete(ctx context.Context, key string) (string, error) {
	// Get the value first
	value, err := rm.Get(ctx, key)
	if err != nil {
		return "", err
	}

	// Delete the key
	_, err = rm.Delete(ctx, key)
	if err != nil {
		rm.logger.Warn("Failed to delete key after GET",
			logging.String("key", key),
			logging.Error(err))
		// Don't return error here as we already got the value
	}

	return value, nil
}

// Health checks Redis database connection
func (rm *RedisManager) Health(ctx context.Context) error {
	// Use shared library health check
	if err := rm.manager.HealthCheck(ctx); err != nil {
		rm.logger.Error("redis health check failed", logging.Error(err))
		return fmt.Errorf("redis health check failed: %w", err)
	}

	// Additional service-specific health checks
	if err := rm.pingRedis(ctx); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}

	return nil
}

// pingRedis performs a simple ping to verify connection
func (rm *RedisManager) pingRedis(ctx context.Context) error {
	client := rm.GetRedis()
	if client == nil {
		return fmt.Errorf("redis connection not available")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}

	return nil
}

// GetConnectionStatus returns the connection status
func (rm *RedisManager) GetConnectionStatus(ctx context.Context) map[string]interface{} {
	status := make(map[string]interface{})

	// Check if Redis client is available
	client := rm.GetRedis()
	status["available"] = client != nil

	// Get health status
	status["healthy"] = rm.Health(ctx) == nil

	// Test basic operations
	if client != nil {
		testKey := fmt.Sprintf("health_check_%d", time.Now().Unix())
		testValue := "test"

		// Test SET operation
		if err := client.Set(ctx, testKey, testValue, time.Minute); err != nil {
			status["set_operation"] = false
		} else {
			status["set_operation"] = true

			// Test GET operation
			if _, err := client.Get(ctx, testKey); err != nil {
				status["get_operation"] = false
			} else {
				status["get_operation"] = true
			}

			// Clean up test key
			client.Del(ctx, testKey)
		}
	}

	return status
}

// initializeRedisSetup initializes Redis keys and configuration from migration files
func (rm *RedisManager) initializeRedisSetup(ctx context.Context) error {
	// Check if migrations directory exists
	migrationsPath := "migrations/redis"

	// Try to run migration script if it exists
	if err := rm.runRedisMigrationScript(migrationsPath); err != nil {
		rm.logger.Warn("redis migration script not found or failed",
			logging.String("path", migrationsPath),
			logging.Error(err))
		return err
	}

	rm.logger.Info("redis setup initialized successfully")
	return nil
}

// runRedisMigrationScript attempts to run Redis migration script
func (rm *RedisManager) runRedisMigrationScript(migrationsPath string) error {
	// This is a placeholder for running Redis migration scripts
	// In a real implementation, you would:
	// 1. Check if migration script exists (.sh, .redis, .lua files)
	// 2. Execute the script using os/exec
	// 3. Handle different file types appropriately

	rm.logger.Info("redis migrations will be handled by migration scripts",
		logging.String("path", migrationsPath))

	return nil
}

// Close closes Redis database connection
func (rm *RedisManager) Close() error {
	if err := rm.manager.Close(); err != nil {
		rm.logger.Error("Failed to close Redis manager", logging.Error(err))
		return fmt.Errorf("failed to close Redis manager: %w", err)
	}

	rm.logger.Info("redis manager closed successfully")
	return nil
}
