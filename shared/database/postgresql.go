package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/usc-platform/shared/config"
)

// PostgreSQLHealthChecker implements health checking for PostgreSQL
type PostgreSQLHealthChecker struct {
	db *sql.DB
}

// NewPostgreSQLHealthChecker creates a new PostgreSQL health checker
func NewPostgreSQLHealthChecker(db *sql.DB) *PostgreSQLHealthChecker {
	return &PostgreSQLHealthChecker{db: db}
}

// Check performs a health check on PostgreSQL
func (h *PostgreSQLHealthChecker) Check(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		return fmt.Errorf("PostgreSQL ping failed: %w", err)
	}

	// Check if we can execute a simple query
	var result int
	if err := h.db.QueryRowContext(ctx, "SELECT 1").Scan(&result); err != nil {
		return fmt.Errorf("PostgreSQL query failed: %w", err)
	}

	return nil
}

// Name returns the name of the health checker
func (h *PostgreSQLHealthChecker) Name() string {
	return "postgresql"
}

// Description returns the description of the health checker
func (h *PostgreSQLHealthChecker) Description() string {
	return "PostgreSQL database health check"
}

// initializePostgreSQL initializes PostgreSQL connection
func (m *DatabaseManager) initializePostgreSQL() error {
	dsn := m.config.GetDatabaseDSN()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open PostgreSQL connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(m.config.Database.MaxConns)
	db.SetMaxIdleConns(m.config.Database.MinConns)
	db.SetConnMaxLifetime(time.Hour)
	db.SetConnMaxIdleTime(time.Minute * 30)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	m.postgres = db
	return nil
}

// PostgreSQLConnection represents a PostgreSQL connection with additional methods
type PostgreSQLConnection struct {
	db *sql.DB
}

// NewPostgreSQLConnection creates a new PostgreSQL connection
func NewPostgreSQLConnection(cfg *config.Config) (*PostgreSQLConnection, error) {
	dsn := cfg.GetDatabaseDSN()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open PostgreSQL connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.Database.MaxConns)
	db.SetMaxIdleConns(cfg.Database.MinConns)
	db.SetConnMaxLifetime(time.Hour)
	db.SetConnMaxIdleTime(time.Minute * 30)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	return &PostgreSQLConnection{db: db}, nil
}

// DB returns the underlying sql.DB
func (p *PostgreSQLConnection) DB() *sql.DB {
	return p.db
}

// Ping tests the connection
func (p *PostgreSQLConnection) Ping(ctx context.Context) error {
	return p.db.PingContext(ctx)
}

// Query executes a query and returns rows
func (p *PostgreSQLConnection) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return p.db.QueryContext(ctx, query, args...)
}

// QueryRow executes a query and returns a single row
func (p *PostgreSQLConnection) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return p.db.QueryRowContext(ctx, query, args...)
}

// Exec executes a query without returning rows
func (p *PostgreSQLConnection) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return p.db.ExecContext(ctx, query, args...)
}

// Begin starts a transaction
func (p *PostgreSQLConnection) Begin(ctx context.Context) (*sql.Tx, error) {
	return p.db.BeginTx(ctx, nil)
}

// BeginTx starts a transaction with options
func (p *PostgreSQLConnection) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return p.db.BeginTx(ctx, opts)
}

// Close closes the connection
func (p *PostgreSQLConnection) Close() error {
	return p.db.Close()
}

// Stats returns connection pool statistics
func (p *PostgreSQLConnection) Stats() sql.DBStats {
	return p.db.Stats()
}

// HealthCheck performs a health check
func (p *PostgreSQLConnection) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := p.db.PingContext(ctx); err != nil {
		return fmt.Errorf("PostgreSQL ping failed: %w", err)
	}

	// Check if we can execute a simple query
	var result int
	if err := p.db.QueryRowContext(ctx, "SELECT 1").Scan(&result); err != nil {
		return fmt.Errorf("PostgreSQL query failed: %w", err)
	}

	return nil
}

// IsPostgreSQLError checks if an error is a PostgreSQL-specific error
func IsPostgreSQLError(err error) bool {
	if err == nil {
		return false
	}

	// Check for PostgreSQL-specific error types
	if _, ok := err.(*pq.Error); ok {
		return true
	}

	// Check for common PostgreSQL error patterns
	errStr := err.Error()
	postgresErrors := []string{
		"pq:",
		"PostgreSQL",
		"postgres",
		"relation",
		"column",
		"constraint",
		"duplicate key",
		"foreign key",
		"unique constraint",
	}

	for _, pattern := range postgresErrors {
		if contains(errStr, pattern) {
			return true
		}
	}

	return false
}

// GetPostgreSQLErrorCode extracts the PostgreSQL error code
func GetPostgreSQLErrorCode(err error) string {
	if pqErr, ok := err.(*pq.Error); ok {
		return string(pqErr.Code)
	}
	return ""
}

// IsUniqueViolation checks if the error is a unique constraint violation
func IsUniqueViolation(err error) bool {
	return GetPostgreSQLErrorCode(err) == "23505"
}

// IsForeignKeyViolation checks if the error is a foreign key constraint violation
func IsForeignKeyViolation(err error) bool {
	return GetPostgreSQLErrorCode(err) == "23503"
}

// IsNotNullViolation checks if the error is a not null constraint violation
func IsNotNullViolation(err error) bool {
	return GetPostgreSQLErrorCode(err) == "23502"
}

// IsCheckViolation checks if the error is a check constraint violation
func IsCheckViolation(err error) bool {
	return GetPostgreSQLErrorCode(err) == "23514"
}
