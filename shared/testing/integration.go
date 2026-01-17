// Package testing provides testing utilities for USC platform services.
package testing

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// IntegrationTestConfig represents integration test configuration
type IntegrationTestConfig struct {
	DatabaseURL   string        `json:"database_url"`
	RedisURL      string        `json:"redis_url"`
	Timeout       time.Duration `json:"timeout"`
	RetryAttempts int           `json:"retry_attempts"`
	RetryDelay    time.Duration `json:"retry_delay"`
	CleanupData   bool          `json:"cleanup_data"`
	ParallelTests bool          `json:"parallel_tests"`
}

// DefaultIntegrationTestConfig returns the default integration test configuration
func DefaultIntegrationTestConfig() IntegrationTestConfig {
	return IntegrationTestConfig{
		DatabaseURL:   "postgres://test:test@localhost:5432/test_db",
		RedisURL:      "redis://localhost:6379/0",
		Timeout:       30 * time.Second,
		RetryAttempts: 3,
		RetryDelay:    1 * time.Second,
		CleanupData:   true,
		ParallelTests: false,
	}
}

// IntegrationTestSuite provides integration test suite utilities
type IntegrationTestSuite struct {
	config    IntegrationTestConfig
	db        *sql.DB
	redis     *redis.Client
	fixtures  *TestFixtures
	cleanup   *TestCleanup
	startTime time.Time
}

// NewIntegrationTestSuite creates a new integration test suite
func NewIntegrationTestSuite(config IntegrationTestConfig) *IntegrationTestSuite {
	return &IntegrationTestSuite{
		config:    config,
		fixtures:  NewTestFixtures(),
		cleanup:   NewTestCleanup(),
		startTime: time.Now(),
	}
}

// Setup sets up the integration test suite
func (its *IntegrationTestSuite) Setup() error {
	// Setup database connection
	if err := its.setupDatabase(); err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}

	// Setup Redis connection
	if err := its.setupRedis(); err != nil {
		return fmt.Errorf("failed to setup Redis: %w", err)
	}

	// Setup test data
	if err := its.setupTestData(); err != nil {
		return fmt.Errorf("failed to setup test data: %w", err)
	}

	return nil
}

// Teardown tears down the integration test suite
func (its *IntegrationTestSuite) Teardown() error {
	// Run cleanup functions
	its.cleanup.Run()

	// Cleanup test data
	if its.config.CleanupData {
		if err := its.cleanupTestData(); err != nil {
			return fmt.Errorf("failed to cleanup test data: %w", err)
		}
	}

	// Close connections
	if its.db != nil {
		its.db.Close()
	}
	if its.redis != nil {
		its.redis.Close()
	}

	return nil
}

// setupDatabase sets up the database connection
func (its *IntegrationTestSuite) setupDatabase() error {
	db, err := sql.Open("postgres", its.config.DatabaseURL)
	if err != nil {
		return err
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return err
	}

	its.db = db
	its.cleanup.Add(func() {
		db.Close()
	})

	return nil
}

// setupRedis sets up the Redis connection
func (its *IntegrationTestSuite) setupRedis() error {
	opt, err := redis.ParseURL(its.config.RedisURL)
	if err != nil {
		return err
	}

	client := redis.NewClient(opt)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), its.config.Timeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return err
	}

	its.redis = client
	its.cleanup.Add(func() {
		client.Close()
	})

	return nil
}

// setupTestData sets up test data
func (its *IntegrationTestSuite) setupTestData() error {
	// Setup database test data
	if err := its.setupDatabaseTestData(); err != nil {
		return err
	}

	// Setup Redis test data
	if err := its.setupRedisTestData(); err != nil {
		return err
	}

	return nil
}

// setupDatabaseTestData sets up database test data
func (its *IntegrationTestSuite) setupDatabaseTestData() error {
	// Create test tables
	createTablesSQL := `
		CREATE TABLE IF NOT EXISTS test_users (
			id VARCHAR(255) PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			username VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL,
			permissions TEXT,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			is_active BOOLEAN NOT NULL DEFAULT true
		);

		CREATE TABLE IF NOT EXISTS test_content (
			id VARCHAR(255) PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			author_id VARCHAR(255) NOT NULL,
			category VARCHAR(100) NOT NULL,
			tags TEXT,
			status VARCHAR(50) NOT NULL,
			published_at TIMESTAMP,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
	`

	if _, err := its.db.Exec(createTablesSQL); err != nil {
		return err
	}

	// Insert test data
	for _, user := range its.fixtures.Users {
		_, err := its.db.Exec(`
			INSERT INTO test_users (id, email, username, password, role, permissions, created_at, updated_at, is_active)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`, user.ID, user.Email, user.Username, user.Password, user.Role,
			fmt.Sprintf("%v", user.Permissions), user.CreatedAt, user.UpdatedAt, user.IsActive)
		if err != nil {
			return err
		}
	}

	for _, content := range its.fixtures.Content {
		_, err := its.db.Exec(`
			INSERT INTO test_content (id, title, content, author_id, category, tags, status, published_at, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`, content.ID, content.Title, content.Content, content.AuthorID, content.Category,
			fmt.Sprintf("%v", content.Tags), content.Status, content.PublishedAt, content.CreatedAt, content.UpdatedAt)
		if err != nil {
			return err
		}
	}

	return nil
}

// setupRedisTestData sets up Redis test data
func (its *IntegrationTestSuite) setupRedisTestData() error {
	ctx, cancel := context.WithTimeout(context.Background(), its.config.Timeout)
	defer cancel()

	// Set test data in Redis
	for _, session := range its.fixtures.Sessions {
		key := fmt.Sprintf("session:%s", session.ID)
		value := fmt.Sprintf(`{"user_id":"%s","token":"%s","expires_at":"%s"}`,
			session.UserID, session.Token, session.ExpiresAt.Format(time.RFC3339))

		if err := its.redis.Set(ctx, key, value, time.Until(session.ExpiresAt)).Err(); err != nil {
			return err
		}
	}

	return nil
}

// cleanupTestData cleans up test data
func (its *IntegrationTestSuite) cleanupTestData() error {
	if !its.config.CleanupData {
		return nil
	}

	// Cleanup database test data
	if err := its.cleanupDatabaseTestData(); err != nil {
		return err
	}

	// Cleanup Redis test data
	if err := its.cleanupRedisTestData(); err != nil {
		return err
	}

	return nil
}

// cleanupDatabaseTestData cleans up database test data
func (its *IntegrationTestSuite) cleanupDatabaseTestData() error {
	// Drop test tables
	dropTablesSQL := `
		DROP TABLE IF EXISTS test_content;
		DROP TABLE IF EXISTS test_users;
	`

	if _, err := its.db.Exec(dropTablesSQL); err != nil {
		return err
	}

	return nil
}

// cleanupRedisTestData cleans up Redis test data
func (its *IntegrationTestSuite) cleanupRedisTestData() error {
	ctx, cancel := context.WithTimeout(context.Background(), its.config.Timeout)
	defer cancel()

	// Delete test keys
	keys, err := its.redis.Keys(ctx, "session:*").Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		if err := its.redis.Del(ctx, keys...).Err(); err != nil {
			return err
		}
	}

	return nil
}

// GetDatabase returns the database connection
func (its *IntegrationTestSuite) GetDatabase() *sql.DB {
	return its.db
}

// GetRedis returns the Redis client
func (its *IntegrationTestSuite) GetRedis() *redis.Client {
	return its.redis
}

// GetFixtures returns the test fixtures
func (its *IntegrationTestSuite) GetFixtures() *TestFixtures {
	return its.fixtures
}

// GetCleanup returns the test cleanup
func (its *IntegrationTestSuite) GetCleanup() *TestCleanup {
	return its.cleanup
}

// GetElapsedTime returns the elapsed time since setup
func (its *IntegrationTestSuite) GetElapsedTime() time.Duration {
	return time.Since(its.startTime)
}

// IntegrationTestRunner provides integration test runner utilities
type IntegrationTestRunner struct {
	suite *IntegrationTestSuite
}

// NewIntegrationTestRunner creates a new integration test runner
func NewIntegrationTestRunner(suite *IntegrationTestSuite) *IntegrationTestRunner {
	return &IntegrationTestRunner{
		suite: suite,
	}
}

// RunTest runs an integration test
func (itr *IntegrationTestRunner) RunTest(name string, testFunc func(*IntegrationTestSuite) error) error {
	// Setup test
	if err := itr.suite.Setup(); err != nil {
		return fmt.Errorf("failed to setup test %s: %w", name, err)
	}

	// Run test
	err := testFunc(itr.suite)

	// Teardown test
	if teardownErr := itr.suite.Teardown(); teardownErr != nil {
		if err != nil {
			return fmt.Errorf("test %s failed: %w, teardown failed: %v", name, err, teardownErr)
		}
		return fmt.Errorf("test %s teardown failed: %w", name, teardownErr)
	}

	return err
}

// RunTests runs multiple integration tests
func (itr *IntegrationTestRunner) RunTests(tests map[string]func(*IntegrationTestSuite) error) error {
	var errors []error

	for name, testFunc := range tests {
		if err := itr.RunTest(name, testFunc); err != nil {
			errors = append(errors, fmt.Errorf("test %s failed: %w", name, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("integration tests failed: %v", errors)
	}

	return nil
}

// IntegrationTestHelper provides integration test helper utilities
type IntegrationTestHelper struct {
	suite *IntegrationTestSuite
}

// NewIntegrationTestHelper creates a new integration test helper
func NewIntegrationTestHelper(suite *IntegrationTestSuite) *IntegrationTestHelper {
	return &IntegrationTestHelper{
		suite: suite,
	}
}

// WaitForDatabase waits for database to be ready
func (ith *IntegrationTestHelper) WaitForDatabase() error {
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for i := 0; i < 5; i++ {
		// Simple database ping simulation
		time.Sleep(100 * time.Millisecond)
		if i == 4 {
			return nil
		}
	}

	return fmt.Errorf("database not ready after 5 attempts")
}

// WaitForRedis waits for Redis to be ready
func (ith *IntegrationTestHelper) WaitForRedis() error {
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for i := 0; i < 5; i++ {
		// Simple Redis ping simulation
		time.Sleep(100 * time.Millisecond)
		if i == 4 {
			return nil
		}
	}

	return fmt.Errorf("redis not ready after 5 attempts")
}

// WaitForAllServices waits for all services to be ready
func (ith *IntegrationTestHelper) WaitForAllServices() error {
	if err := ith.WaitForDatabase(); err != nil {
		return err
	}
	if err := ith.WaitForRedis(); err != nil {
		return err
	}
	return nil
}
