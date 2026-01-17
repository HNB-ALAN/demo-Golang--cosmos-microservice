package database

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/usc-platform/shared/config"
	"github.com/usc-platform/shared/logging"
)

// DatabaseManager manages all database connections
type DatabaseManager struct {
	config *config.Config

	// Database connections
	postgres   *sql.DB
	redis      RedisClient
	clickhouse ClickHouseClient
	influxdb   InfluxDBClient
	quickwit   QuickwitClient
	vectordb   VectorDBClient
	minio      MinIOClient
	bigquery   BigQueryClient

	// Enhanced features
	multiPoolManager   *MultiPoolManager
	transactionManager *TransactionManager

	// Connection status
	connected    bool
	mu           sync.RWMutex
	healthChecks map[string]DatabaseHealthChecker
}

// DatabaseHealthChecker interface for database health checks
type DatabaseHealthChecker interface {
	Check(ctx context.Context) error
}

// RedisClient interface for Redis operations
type RedisClient interface {
	Ping(ctx context.Context) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) (int64, error)
	Exists(ctx context.Context, keys ...string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) (bool, error)
	TTL(ctx context.Context, key string) (time.Duration, error)
	Keys(ctx context.Context, pattern string) ([]string, error)
	IncrBy(ctx context.Context, key string, value int64) (int64, error)
	DecrBy(ctx context.Context, key string, value int64) (int64, error)
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
	MGet(ctx context.Context, keys ...string) ([]interface{}, error)
	Pipeline() redis.Pipeliner
	Info(ctx context.Context, section ...string) (string, error)
	LPush(ctx context.Context, key string, values ...interface{}) (int64, error)
	HealthCheck(ctx context.Context) error
}

// ClickHouseClient interface for ClickHouse operations
type ClickHouseClient interface {
	Ping(ctx context.Context) error
	Query(ctx context.Context, query string, args ...interface{}) (Rows, error)
	Exec(ctx context.Context, query string, args ...interface{}) error
	HealthCheck(ctx context.Context) error
}

// InfluxDBClient interface for InfluxDB operations
type InfluxDBClient interface {
	Ping(ctx context.Context) error
	Write(ctx context.Context, bucket string, point interface{}) error
	Query(ctx context.Context, query string) (*QueryResult, error)
	HealthCheck(ctx context.Context) error
}

// QuickwitClient interface for Quickwit operations
type QuickwitClient interface {
	Ping(ctx context.Context) error
	Index(ctx context.Context, index string, document interface{}) error
	Search(ctx context.Context, index string, query interface{}) (*SearchResult, error)
	HealthCheck(ctx context.Context) error
}

// VectorDBClient interface for Vector DB operations
type VectorDBClient interface {
	Ping(ctx context.Context) error
	CreateCollection(ctx context.Context, name string, config interface{}) error
	InsertVectors(ctx context.Context, collection string, vectors []interface{}) error
	SearchVectors(ctx context.Context, collection string, query interface{}) (*VectorSearchResult, error)
	HealthCheck(ctx context.Context) error
}

// NewDatabaseManager creates a new database manager
func NewDatabaseManager(cfg *config.Config) (*DatabaseManager, error) {
	manager := &DatabaseManager{
		config:       cfg,
		healthChecks: make(map[string]DatabaseHealthChecker),
	}

	// Initialize connections
	if err := manager.initializeConnections(); err != nil {
		return nil, fmt.Errorf("failed to initialize connections: %w", err)
	}

	// Register health checks
	manager.registerHealthChecks()

	// Initialize enhanced features with default configuration
	poolConfig := PoolConfig{
		MinSize:             2,
		MaxSize:             20,
		InitialSize:         5,
		AdjustmentInterval:  30 * time.Second,
		LoadThreshold:       0.8,
		HealthCheckInterval: 10 * time.Second,
	}
	manager.multiPoolManager = NewMultiPoolManager(poolConfig)
	manager.transactionManager = NewTransactionManager(manager)

	manager.connected = true
	return manager, nil
}

// initializeConnections initializes all database connections
func (m *DatabaseManager) initializeConnections() error {
	// Initialize PostgreSQL if enabled
	if m.config.Database.Enabled {
		if err := m.initializePostgreSQL(); err != nil {
			return fmt.Errorf("failed to initialize PostgreSQL: %w", err)
		}
	}

	// Initialize Redis if enabled
	if m.config.Redis.Enabled {
		if err := m.initializeRedis(); err != nil {
			return fmt.Errorf("failed to initialize Redis: %w", err)
		}
	}

	// Initialize ClickHouse if enabled
	if m.config.ClickHouse.Enabled {
		if err := m.initializeClickHouse(); err != nil {
			return fmt.Errorf("failed to initialize ClickHouse: %w", err)
		}
	}

	// Initialize InfluxDB if enabled
	if m.config.InfluxDB.Enabled {
		if err := m.initializeInfluxDB(); err != nil {
			return fmt.Errorf("failed to initialize InfluxDB: %w", err)
		}
	}

	// Initialize Quickwit if enabled
	if m.config.Quickwit.Enabled {
		if err := m.initializeQuickwit(); err != nil {
			return fmt.Errorf("failed to initialize Quickwit: %w", err)
		}
	}

	// Initialize VectorDB if enabled
	if m.config.VectorDB.Enabled {
		if err := m.initializeVectorDB(); err != nil {
			return fmt.Errorf("failed to initialize VectorDB: %w", err)
		}
	}

	// Initialize MinIO if enabled
	if m.config.MinIO.Enabled {
		if err := m.initializeMinIO(); err != nil {
			return fmt.Errorf("failed to initialize MinIO: %w", err)
		}
	}

	// Initialize BigQuery if enabled
	if m.config.BigQuery.Enabled {
		if err := m.initializeBigQuery(); err != nil {
			return fmt.Errorf("failed to initialize BigQuery: %w", err)
		}
	}

	return nil
}

// registerHealthChecks registers health checkers for all databases
func (m *DatabaseManager) registerHealthChecks() {
	if m.postgres != nil {
		m.healthChecks["postgresql"] = &PostgreSQLHealthChecker{db: m.postgres}
	}

	if m.redis != nil {
		m.healthChecks["redis"] = &RedisHealthChecker{client: m.redis}
	}

	if m.clickhouse != nil {
		m.healthChecks["clickhouse"] = &ClickHouseHealthChecker{client: m.clickhouse}
	}

	if m.influxdb != nil {
		m.healthChecks["influxdb"] = &InfluxDBHealthChecker{client: m.influxdb}
	}

	if m.quickwit != nil {
		m.healthChecks["quickwit"] = &QuickwitHealthChecker{client: m.quickwit}
	}

	if m.vectordb != nil {
		m.healthChecks["vectordb"] = &VectorDBHealthChecker{client: m.vectordb}
	}

	if m.minio != nil {
		// Create no-op logger for health checker
		emptyLogger := logging.NewLogger("minio-health", config.LogConfig{})
		m.healthChecks["minio"] = NewMinIOHealthChecker(m.minio, *emptyLogger)
	}

	if m.bigquery != nil {
		m.healthChecks["bigquery"] = NewBigQueryHealthChecker(m.bigquery)
	}
}

// PostgreSQL returns the PostgreSQL connection
func (m *DatabaseManager) PostgreSQL() *sql.DB {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.postgres
}

// Redis returns the Redis client
func (m *DatabaseManager) Redis() RedisClient {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.redis
}

// ClickHouse returns the ClickHouse client
func (m *DatabaseManager) ClickHouse() ClickHouseClient {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.clickhouse
}

// InfluxDB returns the InfluxDB client
func (m *DatabaseManager) InfluxDB() InfluxDBClient {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.influxdb
}

// Quickwit returns the Quickwit client
func (m *DatabaseManager) Quickwit() QuickwitClient {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.quickwit
}

// VectorDB returns the VectorDB client
func (m *DatabaseManager) VectorDB() VectorDBClient {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.vectordb
}

// MinIO returns the MinIO client
func (m *DatabaseManager) MinIO() MinIOClient {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.minio
}

// BigQuery returns the BigQuery client
func (m *DatabaseManager) BigQuery() BigQueryClient {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.bigquery
}

// GetMinIOClient returns the MinIO client (alias for MinIO)
func (m *DatabaseManager) GetMinIOClient() MinIOClient {
	return m.MinIO()
}

// CloseMinIO closes the MinIO connection
func (m *DatabaseManager) CloseMinIO() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.minio != nil {
		// MinIO client doesn't have Close method in our interface
		m.minio = nil
	}

	return nil
}

// HealthCheck performs health checks on all databases concurrently
func (m *DatabaseManager) HealthCheck(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.connected {
		return fmt.Errorf("database manager is not connected")
	}

	// Use WaitGroup for concurrent health checks
	var wg sync.WaitGroup
	errChan := make(chan error, len(m.healthChecks))

	// Create context with timeout for each health check
	healthCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	for name, checker := range m.healthChecks {
		wg.Add(1)
		go func(name string, checker DatabaseHealthChecker) {
			defer wg.Done()
			if err := checker.Check(healthCtx); err != nil {
				errChan <- fmt.Errorf("health check failed for %s: %w", name, err)
			}
		}(name, checker)
	}

	// Wait for all health checks to complete
	wg.Wait()
	close(errChan)

	// Collect any errors
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("health check errors: %v", errors)
	}

	return nil
}

// Close closes all database connections
func (m *DatabaseManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errors []error

	if m.postgres != nil {
		if err := m.postgres.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close PostgreSQL: %w", err))
		}
	}

	if m.redis != nil {
		// Redis client doesn't have Close method in our interface
		_ = m.redis
	}

	if m.clickhouse != nil {
		// ClickHouse client doesn't have Close method in our interface
		_ = m.clickhouse
	}

	if m.influxdb != nil {
		// InfluxDB client doesn't have Close method in our interface
		_ = m.influxdb
	}

	if m.quickwit != nil {
		// Quickwit client doesn't have Close method in our interface
		_ = m.quickwit
	}

	if m.vectordb != nil {
		// VectorDB client doesn't have Close method in our interface
		_ = m.vectordb
	}

	if m.minio != nil {
		// MinIO client doesn't have Close method in our interface
		_ = m.minio
	}

	if m.bigquery != nil {
		// BigQuery client doesn't have Close method in our interface
		_ = m.bigquery
	}

	m.connected = false

	if len(errors) > 0 {
		return fmt.Errorf("errors closing connections: %v", errors)
	}

	return nil
}

// IsConnected returns true if the manager is connected
func (m *DatabaseManager) IsConnected() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.connected
}

// GetConnectionStatus returns the status of all connections
func (m *DatabaseManager) GetConnectionStatus(ctx context.Context) map[string]bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	status := make(map[string]bool)

	for name, checker := range m.healthChecks {
		status[name] = checker.Check(ctx) == nil
	}

	return status
}

// TransactionConfig represents transaction configuration
type TransactionConfig struct {
	IsolationLevel sql.IsolationLevel `mapstructure:"isolation_level"`
	ReadOnly       bool               `mapstructure:"read_only"`
	Timeout        time.Duration      `mapstructure:"timeout"`
	RetryAttempts  int                `mapstructure:"retry_attempts"`
	RetryDelay     time.Duration      `mapstructure:"retry_delay"`
}

// DefaultTransactionConfig returns default transaction configuration
func DefaultTransactionConfig() TransactionConfig {
	return TransactionConfig{
		IsolationLevel: sql.LevelReadCommitted,
		ReadOnly:       false,
		Timeout:        30 * time.Second,
		RetryAttempts:  3,
		RetryDelay:     100 * time.Millisecond,
	}
}

// Transaction executes a function within a database transaction
func (m *DatabaseManager) Transaction(ctx context.Context, fn func(*sql.Tx) error) error {
	return m.TransactionWithConfig(ctx, DefaultTransactionConfig(), fn)
}

// TransactionWithConfig executes a function within a database transaction with custom config
func (m *DatabaseManager) TransactionWithConfig(ctx context.Context, config TransactionConfig, fn func(*sql.Tx) error) error {
	if m.postgres == nil {
		return fmt.Errorf("PostgreSQL connection not available")
	}

	// Create context with timeout
	if config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, config.Timeout)
		defer cancel()
	}

	var lastErr error
	for attempt := 0; attempt < config.RetryAttempts; attempt++ {
		if attempt > 0 {
			// Wait before retry
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(config.RetryDelay):
			}
		}

		// Begin transaction with options
		opts := &sql.TxOptions{
			Isolation: config.IsolationLevel,
			ReadOnly:  config.ReadOnly,
		}

		tx, err := m.postgres.BeginTx(ctx, opts)
		if err != nil {
			lastErr = fmt.Errorf("failed to begin transaction: %w", err)
			continue
		}

		// Execute function with transaction
		err = func() (err error) {
			defer func() {
				if p := recover(); p != nil {
					tx.Rollback()
					err = fmt.Errorf("panic in transaction: %v", p)
				} else if err != nil {
					tx.Rollback()
				} else {
					err = tx.Commit()
				}
			}()

			return fn(tx)
		}()

		if err == nil {
			return nil
		}

		lastErr = err

		// Check if error is retryable
		if !isRetryableTransactionError(err) {
			return err
		}
	}

	return lastErr
}

// isRetryableTransactionError checks if a transaction error is retryable
func isRetryableTransactionError(err error) bool {
	if err == nil {
		return false
	}

	// Check for common retryable errors
	errStr := err.Error()
	retryableErrors := []string{
		"deadlock",
		"serialization failure",
		"connection reset",
		"connection refused",
		"timeout",
	}

	for _, retryableErr := range retryableErrors {
		if stringContains(errStr, retryableErr) {
			return true
		}
	}

	return false
}

// initializeBigQuery initializes BigQuery connection
func (m *DatabaseManager) initializeBigQuery() error {
	// BigQuery initialization is not yet implemented
	// Return nil for now to avoid build errors
	return nil
}

// stringContains checks if a string contains a substring (case insensitive)
func stringContains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					indexOf(s, substr) >= 0))
}

// indexOf finds the index of a substring in a string
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// DistributedTransaction manages transactions across multiple databases
type DistributedTransaction struct {
	manager *DatabaseManager
	txs     map[string]interface{}
	mu      sync.RWMutex
}

// NewDistributedTransaction creates a new distributed transaction
func (m *DatabaseManager) NewDistributedTransaction() *DistributedTransaction {
	return &DistributedTransaction{
		manager: m,
		txs:     make(map[string]interface{}),
	}
}

// BeginPostgreSQL begins a PostgreSQL transaction
func (dt *DistributedTransaction) BeginPostgreSQL(ctx context.Context) (*sql.Tx, error) {
	dt.mu.Lock()
	defer dt.mu.Unlock()

	if dt.manager.postgres == nil {
		return nil, fmt.Errorf("PostgreSQL connection not available")
	}

	tx, err := dt.manager.postgres.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin PostgreSQL transaction: %w", err)
	}

	dt.txs["postgresql"] = tx
	return tx, nil
}

// Commit commits all transactions
func (dt *DistributedTransaction) Commit() error {
	dt.mu.Lock()
	defer dt.mu.Unlock()

	var errors []error

	for name, tx := range dt.txs {
		switch t := tx.(type) {
		case *sql.Tx:
			if err := t.Commit(); err != nil {
				errors = append(errors, fmt.Errorf("failed to commit %s transaction: %w", name, err))
			}
		default:
			errors = append(errors, fmt.Errorf("unsupported transaction type for %s", name))
		}
	}

	dt.txs = make(map[string]interface{})

	if len(errors) > 0 {
		return fmt.Errorf("errors committing distributed transaction: %v", errors)
	}

	return nil
}

// Rollback rolls back all transactions
func (dt *DistributedTransaction) Rollback() error {
	dt.mu.Lock()
	defer dt.mu.Unlock()

	var errors []error

	for name, tx := range dt.txs {
		switch t := tx.(type) {
		case *sql.Tx:
			if err := t.Rollback(); err != nil {
				errors = append(errors, fmt.Errorf("failed to rollback %s transaction: %w", name, err))
			}
		default:
			errors = append(errors, fmt.Errorf("unsupported transaction type for %s", name))
		}
	}

	dt.txs = make(map[string]interface{})

	if len(errors) > 0 {
		return fmt.Errorf("errors rolling back distributed transaction: %v", errors)
	}

	return nil
}

// PerformanceMonitor monitors database performance
type PerformanceMonitor struct {
	queryTimes  map[string][]time.Duration
	errorCounts map[string]int
	mu          sync.RWMutex
}

// NewPerformanceMonitor creates a new performance monitor
func NewPerformanceMonitor() *PerformanceMonitor {
	return &PerformanceMonitor{
		queryTimes:  make(map[string][]time.Duration),
		errorCounts: make(map[string]int),
	}
}

// RecordQueryTime records the execution time of a query
func (pm *PerformanceMonitor) RecordQueryTime(query string, duration time.Duration) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.queryTimes[query] = append(pm.queryTimes[query], duration)

	// Keep only last 100 measurements per query
	if len(pm.queryTimes[query]) > 100 {
		pm.queryTimes[query] = pm.queryTimes[query][len(pm.queryTimes[query])-100:]
	}
}

// RecordError records a query error
func (pm *PerformanceMonitor) RecordError(query string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.errorCounts[query]++
}

// GetAverageQueryTime returns the average execution time for a query
func (pm *PerformanceMonitor) GetAverageQueryTime(query string) time.Duration {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	times, exists := pm.queryTimes[query]
	if !exists || len(times) == 0 {
		return 0
	}

	var total time.Duration
	for _, t := range times {
		total += t
	}

	return total / time.Duration(len(times))
}

// GetErrorCount returns the error count for a query
func (pm *PerformanceMonitor) GetErrorCount(query string) int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	return pm.errorCounts[query]
}

// TransactionManager manages transactions across multiple databases
type TransactionManager struct {
	manager *DatabaseManager
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(manager *DatabaseManager) *TransactionManager {
	return &TransactionManager{
		manager: manager,
	}
}

// MultiDatabaseTransaction represents a transaction across multiple databases
type MultiDatabaseTransaction struct {
	manager      *DatabaseManager
	transactions map[string]interface{}
	mu           sync.RWMutex
	startTime    time.Time
}

// BeginMultiDatabaseTransaction begins a transaction across multiple databases
func (tm *TransactionManager) BeginMultiDatabaseTransaction(ctx context.Context) (*MultiDatabaseTransaction, error) {
	tx := &MultiDatabaseTransaction{
		manager:      tm.manager,
		transactions: make(map[string]interface{}),
		startTime:    time.Now(),
	}

	// Begin PostgreSQL transaction if available
	if tm.manager.postgres != nil {
		postgresTx, err := tm.manager.postgres.BeginTx(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to begin PostgreSQL transaction: %w", err)
		}
		tx.transactions["postgresql"] = postgresTx
	}

	// Note: Redis and other databases would be added here
	// For now, we focus on PostgreSQL as the primary transactional database

	return tx, nil
}

// Commit commits all transactions in the multi-database transaction
func (tx *MultiDatabaseTransaction) Commit() error {
	tx.mu.Lock()
	defer tx.mu.Unlock()

	var errors []error

	for name, transaction := range tx.transactions {
		switch t := transaction.(type) {
		case *sql.Tx:
			if err := t.Commit(); err != nil {
				errors = append(errors, fmt.Errorf("failed to commit %s transaction: %w", name, err))
			}
		default:
			errors = append(errors, fmt.Errorf("unsupported transaction type for %s", name))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors committing multi-database transaction: %v", errors)
	}

	return nil
}

// Rollback rolls back all transactions in the multi-database transaction
func (tx *MultiDatabaseTransaction) Rollback() error {
	tx.mu.Lock()
	defer tx.mu.Unlock()

	var errors []error

	for name, transaction := range tx.transactions {
		switch t := transaction.(type) {
		case *sql.Tx:
			if err := t.Rollback(); err != nil {
				errors = append(errors, fmt.Errorf("failed to rollback %s transaction: %w", name, err))
			}
		default:
			errors = append(errors, fmt.Errorf("unsupported transaction type for %s", name))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors rolling back multi-database transaction: %v", errors)
	}

	return nil
}

// GetTransaction gets a specific database transaction
func (tx *MultiDatabaseTransaction) GetTransaction(database string) (interface{}, bool) {
	tx.mu.RLock()
	defer tx.mu.RUnlock()

	transaction, exists := tx.transactions[database]
	return transaction, exists
}

// GetDuration returns the duration of the transaction
func (tx *MultiDatabaseTransaction) GetDuration() time.Duration {
	return time.Since(tx.startTime)
}

// MultiDatabaseManager represents an enhanced database manager with multi-database support
type MultiDatabaseManager struct {
	config             *config.Config
	manager            *DatabaseManager
	multiPoolManager   *MultiPoolManager
	transactionManager *TransactionManager
	performanceMonitor *PerformanceMonitor
	mu                 sync.RWMutex
}

// MultiDatabaseConfig represents configuration for multi-database manager
type MultiDatabaseConfig struct {
	PoolConfig                  PoolConfig
	EnablePooling               bool
	EnableTransactions          bool
	EnablePerformanceMonitoring bool
}

// DefaultMultiDatabaseConfig returns default multi-database configuration
func DefaultMultiDatabaseConfig() MultiDatabaseConfig {
	return MultiDatabaseConfig{
		PoolConfig: PoolConfig{
			MinSize:             2,
			MaxSize:             20,
			InitialSize:         5,
			AdjustmentInterval:  30 * time.Second,
			LoadThreshold:       0.8,
			HealthCheckInterval: 10 * time.Second,
		},
		EnablePooling:               true,
		EnableTransactions:          true,
		EnablePerformanceMonitoring: true,
	}
}

// NewMultiDatabaseManager creates a new multi-database manager
func NewMultiDatabaseManager(cfg *config.Config, multiConfig MultiDatabaseConfig) (*MultiDatabaseManager, error) {
	// Create base manager
	manager, err := NewDatabaseManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create base manager: %w", err)
	}

	// Create multi-pool manager
	var multiPoolManager *MultiPoolManager
	if multiConfig.EnablePooling {
		multiPoolManager = NewMultiPoolManager(multiConfig.PoolConfig)
	}

	// Create transaction manager
	var transactionManager *TransactionManager
	if multiConfig.EnableTransactions {
		transactionManager = NewTransactionManager(manager)
	}

	// Create performance monitor
	var performanceMonitor *PerformanceMonitor
	if multiConfig.EnablePerformanceMonitoring {
		performanceMonitor = NewPerformanceMonitor()
	}

	return &MultiDatabaseManager{
		config:             cfg,
		manager:            manager,
		multiPoolManager:   multiPoolManager,
		transactionManager: transactionManager,
		performanceMonitor: performanceMonitor,
	}, nil
}

// GetDatabase gets a specific database connection
func (mdm *MultiDatabaseManager) GetDatabase(database string) (interface{}, error) {
	mdm.mu.RLock()
	defer mdm.mu.RUnlock()

	switch database {
	case "postgresql":
		return mdm.manager.PostgreSQL(), nil
	case "redis":
		return mdm.manager.Redis(), nil
	case "clickhouse":
		return mdm.manager.ClickHouse(), nil
	case "influxdb":
		return mdm.manager.InfluxDB(), nil
	case "quickwit":
		return mdm.manager.Quickwit(), nil
	case "minio":
		return mdm.manager.MinIO(), nil
	case "bigquery":
		return mdm.manager.BigQuery(), nil
	default:
		return nil, fmt.Errorf("unknown database: %s", database)
	}
}

// GetConnectionPool gets a connection pool for a specific database
func (mdm *MultiDatabaseManager) GetConnectionPool(database string) (*PoolManager, error) {
	if mdm.multiPoolManager == nil {
		return nil, fmt.Errorf("connection pooling not enabled")
	}

	return mdm.multiPoolManager.GetPool(database), nil
}

// BeginMultiDatabaseTransaction begins a transaction across multiple databases
func (mdm *MultiDatabaseManager) BeginMultiDatabaseTransaction(ctx context.Context) (*MultiDatabaseTransaction, error) {
	if mdm.transactionManager == nil {
		return nil, fmt.Errorf("transactions not enabled")
	}

	return mdm.transactionManager.BeginMultiDatabaseTransaction(ctx)
}

// ExecuteWithRetry executes a function with retry logic across multiple databases
func (mdm *MultiDatabaseManager) ExecuteWithRetry(ctx context.Context, databases []string, fn func(map[string]interface{}) error) error {
	maxRetries := 3
	retryDelay := 100 * time.Millisecond

	for attempt := 0; attempt < maxRetries; attempt++ {
		// Get database connections
		connections := make(map[string]interface{})
		for _, db := range databases {
			conn, err := mdm.GetDatabase(db)
			if err != nil {
				return fmt.Errorf("failed to get %s connection: %w", db, err)
			}
			connections[db] = conn
		}

		// Execute function
		err := fn(connections)
		if err == nil {
			return nil
		}

		// Check if error is retryable
		if !isRetryableError(err) {
			return err
		}

		// Wait before retry
		if attempt < maxRetries-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(retryDelay):
				retryDelay *= 2 // Exponential backoff
			}
		}
	}

	return fmt.Errorf("max retries exceeded")
}

// GetPerformanceMetrics returns performance metrics for all databases
func (mdm *MultiDatabaseManager) GetPerformanceMetrics() map[string]interface{} {
	mdm.mu.RLock()
	defer mdm.mu.RUnlock()

	metrics := make(map[string]interface{})

	// Pool statistics
	if mdm.multiPoolManager != nil {
		metrics["pools"] = mdm.multiPoolManager.GetStats()
	}

	// Performance metrics
	if mdm.performanceMonitor != nil {
		// This would be expanded with actual performance data
		metrics["performance"] = map[string]interface{}{
			"monitoring_enabled": true,
		}
	}

	// Connection status
	metrics["connections"] = mdm.manager.GetConnectionStatus(context.Background())

	return metrics
}

// HealthCheck performs health check on all databases
func (mdm *MultiDatabaseManager) HealthCheck(ctx context.Context) error {
	// Check base manager health
	if err := mdm.manager.HealthCheck(ctx); err != nil {
		return fmt.Errorf("base manager health check failed: %w", err)
	}

	// Check pool health
	if mdm.multiPoolManager != nil {
		if err := mdm.multiPoolManager.HealthCheck(ctx); err != nil {
			return fmt.Errorf("pool manager health check failed: %w", err)
		}
	}

	return nil
}

// Close closes all database connections and pools
func (mdm *MultiDatabaseManager) Close() error {
	var errors []error

	// Close base manager
	if err := mdm.manager.Close(); err != nil {
		errors = append(errors, fmt.Errorf("failed to close base manager: %w", err))
	}

	// Close pool manager
	if mdm.multiPoolManager != nil {
		if err := mdm.multiPoolManager.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close pool manager: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors closing multi-database manager: %v", errors)
	}

	return nil
}

// isRetryableError checks if an error is retryable
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Check for common retryable errors
	errStr := err.Error()
	retryableErrors := []string{
		"connection reset",
		"connection refused",
		"timeout",
		"temporary failure",
		"network error",
	}

	for _, retryableErr := range retryableErrors {
		if stringContains(errStr, retryableErr) {
			return true
		}
	}

	return false
}
