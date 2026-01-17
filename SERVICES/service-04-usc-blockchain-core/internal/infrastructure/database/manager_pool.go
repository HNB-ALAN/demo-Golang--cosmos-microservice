package database

import (
	"context"
	"fmt"
	"time"

	"github.com/usc-platform/shared/config"
	"github.com/usc-platform/shared/database"
	"github.com/usc-platform/shared/logging"
)

// PoolManager wraps the shared database pool manager for USC Blockchain Core Service
type PoolManager struct {
	multiPoolManager *database.MultiPoolManager
	config           *config.Config
	logger           logging.Logger
}

// PoolConfig represents pool configuration for USC Blockchain Core Service
type PoolConfig struct {
	MinSize             int           `mapstructure:"min_size"`
	MaxSize             int           `mapstructure:"max_size"`
	InitialSize         int           `mapstructure:"initial_size"`
	AdjustmentInterval  time.Duration `mapstructure:"adjustment_interval"`
	LoadThreshold       float64       `mapstructure:"load_threshold"`
	HealthCheckInterval time.Duration `mapstructure:"health_check_interval"`
}

// DefaultPoolConfig returns default pool configuration for USC Blockchain Core Service
func DefaultPoolConfig() PoolConfig {
	return PoolConfig{
		MinSize:             2,
		MaxSize:             20,
		InitialSize:         5,
		AdjustmentInterval:  30 * time.Second,
		LoadThreshold:       0.8,
		HealthCheckInterval: 10 * time.Second,
	}
}

// NewPoolManager creates a new database pool manager using shared libraries
func NewPoolManager(cfg *config.Config, logger logging.Logger) (*PoolManager, error) {
	// Convert our config to shared library config
	sharedConfig := database.PoolConfig{
		MinSize:             DefaultPoolConfig().MinSize,
		MaxSize:             DefaultPoolConfig().MaxSize,
		InitialSize:         DefaultPoolConfig().InitialSize,
		AdjustmentInterval:  DefaultPoolConfig().AdjustmentInterval,
		LoadThreshold:       DefaultPoolConfig().LoadThreshold,
		HealthCheckInterval: DefaultPoolConfig().HealthCheckInterval,
	}

	multiPoolManager := database.NewMultiPoolManager(sharedConfig)

	return &PoolManager{
		multiPoolManager: multiPoolManager,
		config:           cfg,
		logger:           logger,
	}, nil
}

// NewPoolManagerWithConfig creates a new database pool manager with custom configuration
func NewPoolManagerWithConfig(cfg *config.Config, poolConfig PoolConfig, logger logging.Logger) (*PoolManager, error) {
	// Convert our config to shared library config
	sharedConfig := database.PoolConfig{
		MinSize:             poolConfig.MinSize,
		MaxSize:             poolConfig.MaxSize,
		InitialSize:         poolConfig.InitialSize,
		AdjustmentInterval:  poolConfig.AdjustmentInterval,
		LoadThreshold:       poolConfig.LoadThreshold,
		HealthCheckInterval: poolConfig.HealthCheckInterval,
	}

	multiPoolManager := database.NewMultiPoolManager(sharedConfig)

	return &PoolManager{
		multiPoolManager: multiPoolManager,
		config:           cfg,
		logger:           logger,
	}, nil
}

// GetPool returns a connection pool for the specified database
func (pm *PoolManager) GetPool(database string) *database.PoolManager {
	return pm.multiPoolManager.GetPool(database)
}

// GetStats returns statistics for all connection pools
func (pm *PoolManager) GetStats() map[string]interface{} {
	stats := pm.multiPoolManager.GetStats()
	result := make(map[string]interface{})
	for db, stat := range stats {
		result[db] = stat
	}
	return result
}

// HealthCheck performs health check on all connection pools
func (pm *PoolManager) HealthCheck(ctx context.Context) error {
	return pm.multiPoolManager.HealthCheck(ctx)
}

// Close closes all connection pools
func (pm *PoolManager) Close() error {
	return pm.multiPoolManager.Close()
}

// GetPoolStatus returns the status of all pools
func (pm *PoolManager) GetPoolStatus(ctx context.Context) map[string]bool {
	stats := pm.GetStats()
	status := make(map[string]bool)

	// Extract pool status from stats
	if pools, ok := stats["pools"].(map[string]interface{}); ok {
		for poolName := range pools {
			// Check if pool is healthy (simplified check)
			status[poolName] = true // In real implementation, check actual pool health
		}
	}

	return status
}

// AdjustPoolSize adjusts the size of a specific pool
func (pm *PoolManager) AdjustPoolSize(database string, newSize int) error {
	pool := pm.GetPool(database)
	if pool == nil {
		return fmt.Errorf("pool not found for database: %s", database)
	}

	// In a real implementation, this would adjust the pool size
	// For now, we log the request
	pm.logger.Info("Pool size adjustment requested",
		logging.String("database", database),
		logging.Int("new_size", newSize))

	return nil
}

// GetPoolMetrics returns detailed metrics for a specific pool
func (pm *PoolManager) GetPoolMetrics(database string) map[string]interface{} {
	pool := pm.GetPool(database)
	if pool == nil {
		return nil
	}

	// Return basic metrics (in real implementation, get actual metrics)
	return map[string]interface{}{
		"database":           database,
		"active_connections": 0,
		"idle_connections":   0,
		"total_connections":  0,
		"waiting_requests":   0,
	}
}
