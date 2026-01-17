package database

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/usc-platform/shared/config"
)

// InfluxDBHealthChecker implements health checking for InfluxDB
type InfluxDBHealthChecker struct {
	client InfluxDBClient
}

// NewInfluxDBHealthChecker creates a new InfluxDB health checker
func NewInfluxDBHealthChecker(client InfluxDBClient) *InfluxDBHealthChecker {
	return &InfluxDBHealthChecker{client: client}
}

// Check performs a health check on InfluxDB
func (h *InfluxDBHealthChecker) Check(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := h.client.Ping(ctx); err != nil {
		return fmt.Errorf("InfluxDB ping failed: %w", err)
	}

	return nil
}

// Name returns the name of the health checker
func (h *InfluxDBHealthChecker) Name() string {
	return "influxdb"
}

// Description returns the description of the health checker
func (h *InfluxDBHealthChecker) Description() string {
	return "InfluxDB database health check"
}

// InfluxDBConnection represents an InfluxDB connection
type InfluxDBConnection struct {
	client   influxdb2.Client
	writeAPI api.WriteAPI
	queryAPI api.QueryAPI
}

// NewInfluxDBConnection creates a new InfluxDB connection
func NewInfluxDBConnection(cfg *config.Config) (*InfluxDBConnection, error) {
	url := cfg.GetInfluxDBURL()

	client := influxdb2.NewClient(url, cfg.InfluxDB.Token)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to ping InfluxDB: %w", err)
	}

	writeAPI := client.WriteAPI(cfg.InfluxDB.Org, cfg.InfluxDB.Bucket)
	queryAPI := client.QueryAPI(cfg.InfluxDB.Org)

	return &InfluxDBConnection{
		client:   client,
		writeAPI: writeAPI,
		queryAPI: queryAPI,
	}, nil
}

// Client returns the underlying InfluxDB client
func (i *InfluxDBConnection) Client() influxdb2.Client {
	return i.client
}

// WriteAPI returns the write API
func (i *InfluxDBConnection) WriteAPI() api.WriteAPI {
	return i.writeAPI
}

// QueryAPI returns the query API
func (i *InfluxDBConnection) QueryAPI() api.QueryAPI {
	return i.queryAPI
}

// Ping tests the connection
func (i *InfluxDBConnection) Ping(ctx context.Context) error {
	_, err := i.client.Ping(ctx)
	return err
}

// Write writes a point to InfluxDB
func (i *InfluxDBConnection) Write(ctx context.Context, bucket string, point interface{}) error {
	writeAPI := i.client.WriteAPI("", bucket)
	writeAPI.WritePoint(point.(*write.Point))
	writeAPI.Flush()
	return nil
}

// Query executes a query
func (i *InfluxDBConnection) Query(ctx context.Context, query string) (*QueryResult, error) {
	result, err := i.queryAPI.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	return &QueryResult{result: result}, nil
}

// Close closes the connection
func (i *InfluxDBConnection) Close() error {
	i.client.Close()
	return nil
}

// HealthCheck performs a health check
func (i *InfluxDBConnection) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := i.client.Ping(ctx); err != nil {
		return fmt.Errorf("InfluxDB ping failed: %w", err)
	}

	return nil
}

// Point represents an InfluxDB data point
type Point interface {
	AddTag(key, value string)
	AddField(key string, value interface{})
	SetTime(t time.Time)
}

// QueryResult represents the result of an InfluxDB query
type QueryResult struct {
	result *api.QueryTableResult
}

// Next returns the next row from the query result
func (q *QueryResult) Next() bool {
	return q.result.Next()
}

// Record returns the current record
func (q *QueryResult) Record() interface{} {
	return q.result.Record()
}

// Close closes the query result
func (q *QueryResult) Close() error {
	q.result.Close()
	return nil
}

// Err returns any error from the query
func (q *QueryResult) Err() error {
	return q.result.Err()
}

// initializeInfluxDB initializes InfluxDB connection
func (m *DatabaseManager) initializeInfluxDB() error {
	url := m.config.GetInfluxDBURL()

	client := influxdb2.NewClient(url, m.config.InfluxDB.Token)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx); err != nil {
		client.Close()
		return fmt.Errorf("failed to ping InfluxDB: %w", err)
	}

	writeAPI := client.WriteAPI(m.config.InfluxDB.Org, m.config.InfluxDB.Bucket)
	queryAPI := client.QueryAPI(m.config.InfluxDB.Org)

	m.influxdb = &InfluxDBConnection{
		client:   client,
		writeAPI: writeAPI,
		queryAPI: queryAPI,
	}

	return nil
}

// IsInfluxDBError checks if an error is InfluxDB-specific
func IsInfluxDBError(err error) bool {
	if err == nil {
		return false
	}

	// Check for common InfluxDB error patterns
	errStr := err.Error()
	influxdbErrors := []string{
		"InfluxDB",
		"influxdb",
		"influx",
		"bucket",
		"measurement",
		"field",
		"tag",
		"line protocol",
	}

	for _, pattern := range influxdbErrors {
		if contains(errStr, pattern) {
			return true
		}
	}

	return false
}
