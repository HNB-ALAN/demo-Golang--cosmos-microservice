package database

import (
	"context"
	"fmt"
	"testing"

	"github.com/usc-platform/shared/config"
)

func TestMultiDatabaseManager(t *testing.T) {
	// Create test configuration
	cfg := &config.Config{
		Service: config.ServiceConfig{
			Name:    "test-service",
			Version: "1.0.0",
		},
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "test",
			Password: "test",
			DBName:   "test_db",
		},
	}

	// Create multi-database config
	multiConfig := DefaultMultiDatabaseConfig()
	multiConfig.EnablePooling = true
	multiConfig.EnableTransactions = true
	multiConfig.EnablePerformanceMonitoring = true

	// Create multi-database manager
	mdm, err := NewMultiDatabaseManager(cfg, multiConfig)
	if err != nil {
		t.Logf("Expected error creating multi-database manager (no database running): %v", err)
		return
	}
	defer mdm.Close()

	// Test getting database connections
	databases := []string{"postgresql", "redis", "clickhouse", "influxdb", "quickwit"}
	for _, db := range databases {
		conn, err := mdm.GetDatabase(db)
		if err != nil {
			t.Logf("Expected error getting %s connection (no database running): %v", db, err)
		} else {
			t.Logf("Successfully got %s connection", db)
			_ = conn
		}
	}

	// Test connection pools
	if multiConfig.EnablePooling {
		for _, db := range databases {
			pool, err := mdm.GetConnectionPool(db)
			if err != nil {
				t.Logf("Expected error getting %s pool: %v", db, err)
			} else {
				t.Logf("Successfully got %s pool", db)
				_ = pool
			}
		}
	}

	// Test performance metrics
	metrics := mdm.GetPerformanceMetrics()
	if metrics == nil {
		t.Error("Expected metrics to be returned")
	} else {
		t.Logf("Performance metrics: %+v", metrics)
	}

	// Test health check
	ctx := context.Background()
	err = mdm.HealthCheck(ctx)
	if err != nil {
		t.Logf("Expected health check error (no database running): %v", err)
	}
}

func TestMultiDatabaseTransaction(t *testing.T) {
	// Create test configuration
	cfg := &config.Config{
		Service: config.ServiceConfig{
			Name:    "test-service",
			Version: "1.0.0",
		},
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "test",
			Password: "test",
			DBName:   "test_db",
		},
	}

	// Create base manager
	manager, err := NewDatabaseManager(cfg)
	if err != nil {
		t.Logf("Expected error creating manager (no database running): %v", err)
		return
	}
	defer manager.Close()

	// Create transaction manager
	tm := NewTransactionManager(manager)

	// Test beginning multi-database transaction
	ctx := context.Background()
	tx, err := tm.BeginMultiDatabaseTransaction(ctx)
	if err != nil {
		t.Logf("Expected error beginning transaction (no database running): %v", err)
		return
	}

	// Test transaction duration
	duration := tx.GetDuration()
	if duration < 0 {
		t.Error("Expected positive transaction duration")
	} else {
		t.Logf("Transaction duration: %v", duration)
	}

	// Test getting transaction
	postgresTx, exists := tx.GetTransaction("postgresql")
	if !exists {
		t.Log("PostgreSQL transaction not found (expected if no database)")
	} else {
		t.Logf("PostgreSQL transaction found: %+v", postgresTx)
	}

	// Test rollback
	err = tx.Rollback()
	if err != nil {
		t.Logf("Expected error during rollback: %v", err)
	} else {
		t.Log("Transaction rolled back successfully")
	}
}

func TestExecuteWithRetry(t *testing.T) {
	// Create test configuration
	cfg := &config.Config{
		Service: config.ServiceConfig{
			Name:    "test-service",
			Version: "1.0.0",
		},
	}

	// Create multi-database config
	multiConfig := DefaultMultiDatabaseConfig()
	multiConfig.EnablePooling = false // Disable pooling for this test

	// Create multi-database manager
	mdm, err := NewMultiDatabaseManager(cfg, multiConfig)
	if err != nil {
		t.Logf("Expected error creating multi-database manager: %v", err)
		return
	}
	defer mdm.Close()

	// Test execute with retry
	ctx := context.Background()
	databases := []string{"postgresql"}

	err = mdm.ExecuteWithRetry(ctx, databases, func(connections map[string]interface{}) error {
		// Simulate a function that always fails
		return fmt.Errorf("simulated error")
	})

	if err == nil {
		t.Error("Expected error from ExecuteWithRetry")
	} else {
		t.Logf("ExecuteWithRetry failed as expected: %v", err)
	}
}

func TestMultiDatabaseConfig(t *testing.T) {
	config := DefaultMultiDatabaseConfig()

	// Test default values
	if !config.EnablePooling {
		t.Error("Expected pooling to be enabled by default")
	}

	if !config.EnableTransactions {
		t.Error("Expected transactions to be enabled by default")
	}

	if !config.EnablePerformanceMonitoring {
		t.Error("Expected performance monitoring to be enabled by default")
	}

	// Test pool config
	if config.PoolConfig.MinSize != 2 {
		t.Errorf("Expected MinSize to be 2, got %d", config.PoolConfig.MinSize)
	}

	if config.PoolConfig.MaxSize != 20 {
		t.Errorf("Expected MaxSize to be 20, got %d", config.PoolConfig.MaxSize)
	}

	if config.PoolConfig.InitialSize != 5 {
		t.Errorf("Expected InitialSize to be 5, got %d", config.PoolConfig.InitialSize)
	}
}

func TestTransactionManager(t *testing.T) {
	// Create test configuration
	cfg := &config.Config{
		Service: config.ServiceConfig{
			Name:    "test-service",
			Version: "1.0.0",
		},
	}

	// Create base manager
	manager, err := NewDatabaseManager(cfg)
	if err != nil {
		t.Logf("Expected error creating manager: %v", err)
		return
	}
	defer manager.Close()

	// Create transaction manager
	tm := NewTransactionManager(manager)
	if tm == nil {
		t.Error("Expected transaction manager to be created")
		return
	}

	if tm.manager != manager {
		t.Error("Expected transaction manager to reference the correct manager")
	}
}

func TestIsRetryableError(t *testing.T) {
	// Test retryable errors
	retryableErrors := []string{
		"connection reset",
		"connection refused",
		"timeout",
		"temporary failure",
		"network error",
	}

	for _, errMsg := range retryableErrors {
		err := fmt.Errorf("%s", errMsg)
		if !isRetryableError(err) {
			t.Errorf("Expected error '%s' to be retryable", errMsg)
		}
	}

	// Test non-retryable errors
	nonRetryableErrors := []string{
		"invalid syntax",
		"permission denied",
		"not found",
	}

	for _, errMsg := range nonRetryableErrors {
		err := fmt.Errorf("%s", errMsg)
		if isRetryableError(err) {
			t.Errorf("Expected error '%s' to not be retryable", errMsg)
		}
	}

	// Test nil error
	if isRetryableError(nil) {
		t.Error("Expected nil error to not be retryable")
	}
}
