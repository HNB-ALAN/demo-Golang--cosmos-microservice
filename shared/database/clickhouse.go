package database

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/usc-platform/shared/config"
)

// ClickHouseHealthChecker implements health checking for ClickHouse
type ClickHouseHealthChecker struct {
	client ClickHouseClient
}

// NewClickHouseHealthChecker creates a new ClickHouse health checker
func NewClickHouseHealthChecker(client ClickHouseClient) *ClickHouseHealthChecker {
	return &ClickHouseHealthChecker{client: client}
}

// Check performs a health check on ClickHouse
func (h *ClickHouseHealthChecker) Check(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := h.client.Ping(ctx); err != nil {
		return fmt.Errorf("ClickHouse ping failed: %w", err)
	}

	return nil
}

// Name returns the name of the health checker
func (h *ClickHouseHealthChecker) Name() string {
	return "clickhouse"
}

// Description returns the description of the health checker
func (h *ClickHouseHealthChecker) Description() string {
	return "ClickHouse database health check"
}

// ClickHouseConnection represents a ClickHouse connection
type ClickHouseConnection struct {
	conn clickhouse.Conn
}

// NewClickHouseConnection creates a new ClickHouse connection
func NewClickHouseConnection(cfg *config.Config) (*ClickHouseConnection, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", cfg.ClickHouse.Host, cfg.ClickHouse.Port)},
		Auth: clickhouse.Auth{
			Database: cfg.ClickHouse.DBName,
			Username: cfg.ClickHouse.User,
			Password: cfg.ClickHouse.Password,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: 10 * time.Second,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to ClickHouse: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := conn.Ping(ctx); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to ping ClickHouse: %w", err)
	}

	return &ClickHouseConnection{conn: conn}, nil
}

// Conn returns the underlying ClickHouse connection
func (c *ClickHouseConnection) Conn() clickhouse.Conn {
	return c.conn
}

// Ping tests the connection
func (c *ClickHouseConnection) Ping(ctx context.Context) error {
	return c.conn.Ping(ctx)
}

// Query executes a query and returns rows
func (c *ClickHouseConnection) Query(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	return c.conn.Query(ctx, query, args...)
}

// Exec executes a query without returning rows
func (c *ClickHouseConnection) Exec(ctx context.Context, query string, args ...interface{}) error {
	return c.conn.Exec(ctx, query, args...)
}

// Close closes the connection
func (c *ClickHouseConnection) Close() error {
	return c.conn.Close()
}

// HealthCheck performs a health check
func (c *ClickHouseConnection) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := c.conn.Ping(ctx); err != nil {
		return fmt.Errorf("ClickHouse ping failed: %w", err)
	}

	return nil
}

// Rows interface for ClickHouse query results
type Rows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close() error
	Err() error
	Columns() []string
}

// initializeClickHouse initializes ClickHouse connection
func (m *DatabaseManager) initializeClickHouse() error {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", m.config.ClickHouse.Host, m.config.ClickHouse.Port)},
		Auth: clickhouse.Auth{
			Database: m.config.ClickHouse.DBName,
			Username: m.config.ClickHouse.User,
			Password: m.config.ClickHouse.Password,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: 10 * time.Second,
	})

	if err != nil {
		return fmt.Errorf("failed to connect to ClickHouse: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := conn.Ping(ctx); err != nil {
		conn.Close()
		return fmt.Errorf("failed to ping ClickHouse: %w", err)
	}

	m.clickhouse = &ClickHouseConnection{conn: conn}
	return nil
}

// IsClickHouseError checks if an error is ClickHouse-specific
func IsClickHouseError(err error) bool {
	if err == nil {
		return false
	}

	// Check for common ClickHouse error patterns
	errStr := err.Error()
	clickhouseErrors := []string{
		"ClickHouse",
		"clickhouse",
		"Code:",
		"DB::Exception",
		"Table",
		"Column",
		"Database",
	}

	for _, pattern := range clickhouseErrors {
		if contains(errStr, pattern) {
			return true
		}
	}

	return false
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

// containsSubstring checks if a string contains a substring
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
