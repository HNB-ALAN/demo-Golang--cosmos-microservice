// Package testing provides testing utilities for USC platform services.
package testing

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MockCache provides a mock cache implementation
type MockCache struct {
	data map[string]interface{}
	ttl  map[string]time.Time
}

// NewMockCache creates a new mock cache
func NewMockCache() *MockCache {
	return &MockCache{
		data: make(map[string]interface{}),
		ttl:  make(map[string]time.Time),
	}
}

// Get retrieves a value from the mock cache
func (mc *MockCache) Get(ctx context.Context, key string) (interface{}, error) {
	// Check if key exists
	value, exists := mc.data[key]
	if !exists {
		return nil, errors.New("key not found")
	}

	// Check if TTL has expired
	if expiry, hasTTL := mc.ttl[key]; hasTTL && time.Now().After(expiry) {
		delete(mc.data, key)
		delete(mc.ttl, key)
		return nil, errors.New("key expired")
	}

	return value, nil
}

// Set stores a value in the mock cache
func (mc *MockCache) Set(ctx context.Context, key string, value interface{}, ttl ...time.Duration) error {
	mc.data[key] = value

	if len(ttl) > 0 {
		mc.ttl[key] = time.Now().Add(ttl[0])
	}

	return nil
}

// Delete removes a value from the mock cache
func (mc *MockCache) Delete(ctx context.Context, key string) error {
	delete(mc.data, key)
	delete(mc.ttl, key)
	return nil
}

// Exists checks if a key exists in the mock cache
func (mc *MockCache) Exists(ctx context.Context, key string) (bool, error) {
	_, exists := mc.data[key]
	if !exists {
		return false, nil
	}

	// Check if TTL has expired
	if expiry, hasTTL := mc.ttl[key]; hasTTL && time.Now().After(expiry) {
		delete(mc.data, key)
		delete(mc.ttl, key)
		return false, nil
	}

	return true, nil
}

// Clear clears all data from the mock cache
func (mc *MockCache) Clear() {
	mc.data = make(map[string]interface{})
	mc.ttl = make(map[string]time.Time)
}

// MockDatabase provides a mock database implementation
type MockDatabase struct {
	tables map[string]map[string]interface{}
}

// NewMockDatabase creates a new mock database
func NewMockDatabase() *MockDatabase {
	return &MockDatabase{
		tables: make(map[string]map[string]interface{}),
	}
}

// Insert inserts a document into the mock database
func (md *MockDatabase) Insert(table string, id string, document interface{}) error {
	if md.tables[table] == nil {
		md.tables[table] = make(map[string]interface{})
	}
	md.tables[table][id] = document
	return nil
}

// Find finds a document in the mock database
func (md *MockDatabase) Find(table string, id string) (interface{}, error) {
	if md.tables[table] == nil {
		return nil, errors.New("table not found")
	}

	document, exists := md.tables[table][id]
	if !exists {
		return nil, errors.New("document not found")
	}

	return document, nil
}

// Update updates a document in the mock database
func (md *MockDatabase) Update(table string, id string, document interface{}) error {
	if md.tables[table] == nil {
		return errors.New("table not found")
	}

	if _, exists := md.tables[table][id]; !exists {
		return errors.New("document not found")
	}

	md.tables[table][id] = document
	return nil
}

// Delete deletes a document from the mock database
func (md *MockDatabase) Delete(table string, id string) error {
	if md.tables[table] == nil {
		return errors.New("table not found")
	}

	if _, exists := md.tables[table][id]; !exists {
		return errors.New("document not found")
	}

	delete(md.tables[table], id)
	return nil
}

// List lists all documents in a table
func (md *MockDatabase) List(table string) ([]interface{}, error) {
	if md.tables[table] == nil {
		return nil, errors.New("table not found")
	}

	documents := make([]interface{}, 0, len(md.tables[table]))
	for _, doc := range md.tables[table] {
		documents = append(documents, doc)
	}

	return documents, nil
}

// Clear clears all data from the mock database
func (md *MockDatabase) Clear() {
	md.tables = make(map[string]map[string]interface{})
}

// MockRedisClient provides a mock Redis client
type MockRedisClient struct {
	data map[string]string
	ttl  map[string]time.Time
}

// NewMockRedisClient creates a new mock Redis client
func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		data: make(map[string]string),
		ttl:  make(map[string]time.Time),
	}
}

// Get retrieves a value from the mock Redis
func (mrc *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	cmd := redis.NewStringCmd(ctx, "get", key)

	value, exists := mrc.data[key]
	if !exists {
		cmd.SetErr(redis.Nil)
		return cmd
	}

	// Check if TTL has expired
	if expiry, hasTTL := mrc.ttl[key]; hasTTL && time.Now().After(expiry) {
		delete(mrc.data, key)
		delete(mrc.ttl, key)
		cmd.SetErr(redis.Nil)
		return cmd
	}

	cmd.SetVal(value)
	return cmd
}

// Set stores a value in the mock Redis
func (mrc *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(ctx, "set", key, value)

	mrc.data[key] = value.(string)
	if expiration > 0 {
		mrc.ttl[key] = time.Now().Add(expiration)
	}

	cmd.SetVal("OK")
	return cmd
}

// Del deletes a key from the mock Redis
func (mrc *MockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	cmd := redis.NewIntCmd(ctx, "del")

	count := int64(0)
	for _, key := range keys {
		if _, exists := mrc.data[key]; exists {
			delete(mrc.data, key)
			delete(mrc.ttl, key)
			count++
		}
	}

	cmd.SetVal(count)
	return cmd
}

// Exists checks if a key exists in the mock Redis
func (mrc *MockRedisClient) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	cmd := redis.NewIntCmd(ctx, "exists")

	count := int64(0)
	for _, key := range keys {
		if _, exists := mrc.data[key]; exists {
			// Check if TTL has expired
			if expiry, hasTTL := mrc.ttl[key]; hasTTL && time.Now().After(expiry) {
				delete(mrc.data, key)
				delete(mrc.ttl, key)
			} else {
				count++
			}
		}
	}

	cmd.SetVal(count)
	return cmd
}

// Clear clears all data from the mock Redis
func (mrc *MockRedisClient) Clear() {
	mrc.data = make(map[string]string)
	mrc.ttl = make(map[string]time.Time)
}

// MockGRPCServer provides a mock gRPC server
type MockGRPCServer struct {
	services map[string]interface{}
}

// NewMockGRPCServer creates a new mock gRPC server
func NewMockGRPCServer() *MockGRPCServer {
	return &MockGRPCServer{
		services: make(map[string]interface{}),
	}
}

// RegisterService registers a service with the mock server
func (mgs *MockGRPCServer) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	mgs.services[desc.ServiceName] = impl
}

// GetService returns a registered service
func (mgs *MockGRPCServer) GetService(name string) (interface{}, bool) {
	service, exists := mgs.services[name]
	return service, exists
}

// Clear clears all services from the mock server
func (mgs *MockGRPCServer) Clear() {
	mgs.services = make(map[string]interface{})
}

// MockGRPCClient provides a mock gRPC client
type MockGRPCClient struct {
	responses map[string]interface{}
	errors    map[string]error
}

// NewMockGRPCClient creates a new mock gRPC client
func NewMockGRPCClient() *MockGRPCClient {
	return &MockGRPCClient{
		responses: make(map[string]interface{}),
		errors:    make(map[string]error),
	}
}

// SetResponse sets a response for a method
func (mgc *MockGRPCClient) SetResponse(method string, response interface{}) {
	mgc.responses[method] = response
}

// SetError sets an error for a method
func (mgc *MockGRPCClient) SetError(method string, err error) {
	mgc.errors[method] = err
}

// Invoke invokes a method on the mock client
func (mgc *MockGRPCClient) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	// Check if there's an error for this method
	if err, exists := mgc.errors[method]; exists {
		return err
	}

	// Check if there's a response for this method
	if _, exists := mgc.responses[method]; exists {
		// Copy response to reply
		// This is a simplified implementation
		// In a real implementation, you would use reflection or type assertion
		return nil
	}

	return status.Error(codes.Unimplemented, "method not implemented")
}

// Clear clears all responses and errors from the mock client
func (mgc *MockGRPCClient) Clear() {
	mgc.responses = make(map[string]interface{})
	mgc.errors = make(map[string]error)
}

// MockLogger provides a mock logger
type MockLogger struct {
	logs []LogEntry
}

// LogEntry represents a log entry
type LogEntry struct {
	Level   string
	Message string
	Fields  map[string]interface{}
	Time    time.Time
}

// NewMockLogger creates a new mock logger
func NewMockLogger() *MockLogger {
	return &MockLogger{
		logs: make([]LogEntry, 0),
	}
}

// Log logs a message
func (ml *MockLogger) Log(level string, message string, fields map[string]interface{}) {
	ml.logs = append(ml.logs, LogEntry{
		Level:   level,
		Message: message,
		Fields:  fields,
		Time:    time.Now(),
	})
}

// GetLogs returns all logged messages
func (ml *MockLogger) GetLogs() []LogEntry {
	return ml.logs
}

// Clear clears all logs
func (ml *MockLogger) Clear() {
	ml.logs = make([]LogEntry, 0)
}

// MockMetrics provides a mock metrics collector
type MockMetrics struct {
	counters   map[string]int64
	gauges     map[string]float64
	histograms map[string][]float64
}

// NewMockMetrics creates a new mock metrics collector
func NewMockMetrics() *MockMetrics {
	return &MockMetrics{
		counters:   make(map[string]int64),
		gauges:     make(map[string]float64),
		histograms: make(map[string][]float64),
	}
}

// IncrementCounter increments a counter
func (mm *MockMetrics) IncrementCounter(name string, labels map[string]string) {
	mm.counters[name]++
}

// SetGauge sets a gauge value
func (mm *MockMetrics) SetGauge(name string, value float64, labels map[string]string) {
	mm.gauges[name] = value
}

// ObserveHistogram observes a histogram value
func (mm *MockMetrics) ObserveHistogram(name string, value float64, labels map[string]string) {
	mm.histograms[name] = append(mm.histograms[name], value)
}

// GetCounter returns a counter value
func (mm *MockMetrics) GetCounter(name string) int64 {
	return mm.counters[name]
}

// GetGauge returns a gauge value
func (mm *MockMetrics) GetGauge(name string) float64 {
	return mm.gauges[name]
}

// GetHistogram returns a histogram values
func (mm *MockMetrics) GetHistogram(name string) []float64 {
	return mm.histograms[name]
}

// Clear clears all metrics
func (mm *MockMetrics) Clear() {
	mm.counters = make(map[string]int64)
	mm.gauges = make(map[string]float64)
	mm.histograms = make(map[string][]float64)
}
