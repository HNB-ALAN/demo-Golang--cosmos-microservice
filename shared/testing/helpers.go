// Package testing provides testing utilities for USC platform services.
package testing

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// HTTPTestHelper provides HTTP testing utilities
type HTTPTestHelper struct {
	server *httptest.Server
	client *http.Client
}

// NewHTTPTestHelper creates a new HTTP test helper
func NewHTTPTestHelper(handler http.Handler) *HTTPTestHelper {
	server := httptest.NewServer(handler)
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	return &HTTPTestHelper{
		server: server,
		client: client,
	}
}

// Close closes the test server
func (h *HTTPTestHelper) Close() {
	h.server.Close()
}

// Get performs a GET request
func (h *HTTPTestHelper) Get(path string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("GET", h.server.URL+path, nil)
	if err != nil {
		return nil, err
	}

	// Add headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return h.client.Do(req)
}

// Post performs a POST request
func (h *HTTPTestHelper) Post(path string, body interface{}, headers map[string]string) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest("POST", h.server.URL+path, bodyReader)
	if err != nil {
		return nil, err
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")

	// Add headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return h.client.Do(req)
}

// Put performs a PUT request
func (h *HTTPTestHelper) Put(path string, body interface{}, headers map[string]string) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest("PUT", h.server.URL+path, bodyReader)
	if err != nil {
		return nil, err
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")

	// Add headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return h.client.Do(req)
}

// Delete performs a DELETE request
func (h *HTTPTestHelper) Delete(path string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", h.server.URL+path, nil)
	if err != nil {
		return nil, err
	}

	// Add headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return h.client.Do(req)
}

// ParseResponse parses a JSON response
func (h *HTTPTestHelper) ParseResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}

// AssertStatusCode asserts the response status code
func (h *HTTPTestHelper) AssertStatusCode(resp *http.Response, expected int) error {
	if resp.StatusCode != expected {
		return fmt.Errorf("expected status code %d, got %d", expected, resp.StatusCode)
	}
	return nil
}

// AssertHeader asserts a response header
func (h *HTTPTestHelper) AssertHeader(resp *http.Response, name, expected string) error {
	actual := resp.Header.Get(name)
	if actual != expected {
		return fmt.Errorf("expected header %s to be %s, got %s", name, expected, actual)
	}
	return nil
}

// GRPCTestHelper provides gRPC testing utilities
type GRPCTestHelper struct {
	server *grpc.Server
	conn   *grpc.ClientConn
}

// NewGRPCTestHelper creates a new gRPC test helper
func NewGRPCTestHelper(server *grpc.Server) *GRPCTestHelper {
	return &GRPCTestHelper{
		server: server,
	}
}

// Start starts the gRPC server
func (h *GRPCTestHelper) Start() error {
	// In a real implementation, you would start the server
	// For testing, we'll just return nil
	return nil
}

// Stop stops the gRPC server
func (h *GRPCTestHelper) Stop() {
	if h.conn != nil {
		h.conn.Close()
	}
	if h.server != nil {
		h.server.Stop()
	}
}

// Connect connects to the gRPC server
func (h *GRPCTestHelper) Connect(address string) error {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	h.conn = conn
	return nil
}

// GetConnection returns the gRPC connection
func (h *GRPCTestHelper) GetConnection() *grpc.ClientConn {
	return h.conn
}

// AssertGRPCError asserts a gRPC error
func (h *GRPCTestHelper) AssertGRPCError(err error, expectedCode codes.Code) error {
	if err == nil {
		return fmt.Errorf("expected gRPC error, got nil")
	}

	st, ok := status.FromError(err)
	if !ok {
		return fmt.Errorf("expected gRPC status error, got %v", err)
	}

	if st.Code() != expectedCode {
		return fmt.Errorf("expected gRPC error code %v, got %v", expectedCode, st.Code())
	}

	return nil
}

// TestContext provides test context utilities
type TestContext struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// NewTestContext creates a new test context
func NewTestContext(timeout time.Duration) *TestContext {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return &TestContext{
		ctx:    ctx,
		cancel: cancel,
	}
}

// GetContext returns the test context
func (tc *TestContext) GetContext() context.Context {
	return tc.ctx
}

// Cancel cancels the test context
func (tc *TestContext) Cancel() {
	tc.cancel()
}

// Cleanup cleans up the test context
func (tc *TestContext) Cleanup() {
	tc.cancel()
}

// TestAssertion provides test assertion utilities
type TestAssertion struct{}

// NewTestAssertion creates a new test assertion
func NewTestAssertion() *TestAssertion {
	return &TestAssertion{}
}

// AssertEqual asserts that two values are equal
func (ta *TestAssertion) AssertEqual(actual, expected interface{}) error {
	if actual != expected {
		return fmt.Errorf("expected %v, got %v", expected, actual)
	}
	return nil
}

// AssertNotEqual asserts that two values are not equal
func (ta *TestAssertion) AssertNotEqual(actual, expected interface{}) error {
	if actual == expected {
		return fmt.Errorf("expected %v to not equal %v", actual, expected)
	}
	return nil
}

// AssertTrue asserts that a value is true
func (ta *TestAssertion) AssertTrue(actual bool) error {
	if !actual {
		return fmt.Errorf("expected true, got false")
	}
	return nil
}

// AssertFalse asserts that a value is false
func (ta *TestAssertion) AssertFalse(actual bool) error {
	if actual {
		return fmt.Errorf("expected false, got true")
	}
	return nil
}

// AssertNil asserts that a value is nil
func (ta *TestAssertion) AssertNil(actual interface{}) error {
	if actual != nil {
		return fmt.Errorf("expected nil, got %v", actual)
	}
	return nil
}

// AssertNotNil asserts that a value is not nil
func (ta *TestAssertion) AssertNotNil(actual interface{}) error {
	if actual == nil {
		return fmt.Errorf("expected not nil, got nil")
	}
	return nil
}

// AssertContains asserts that a string contains a substring
func (ta *TestAssertion) AssertContains(str, substr string) error {
	if !strings.Contains(str, substr) {
		return fmt.Errorf("expected %s to contain %s", str, substr)
	}
	return nil
}

// AssertNotContains asserts that a string does not contain a substring
func (ta *TestAssertion) AssertNotContains(str, substr string) error {
	if strings.Contains(str, substr) {
		return fmt.Errorf("expected %s to not contain %s", str, substr)
	}
	return nil
}

// TestTimer provides test timing utilities
type TestTimer struct {
	start time.Time
}

// NewTestTimer creates a new test timer
func NewTestTimer() *TestTimer {
	return &TestTimer{
		start: time.Now(),
	}
}

// Elapsed returns the elapsed time
func (tt *TestTimer) Elapsed() time.Duration {
	return time.Since(tt.start)
}

// Reset resets the timer
func (tt *TestTimer) Reset() {
	tt.start = time.Now()
}

// TestLogger provides test logging utilities
type TestLogger struct {
	logs []string
}

// NewTestLogger creates a new test logger
func NewTestLogger() *TestLogger {
	return &TestLogger{
		logs: make([]string, 0),
	}
}

// Log logs a message
func (tl *TestLogger) Log(message string) {
	tl.logs = append(tl.logs, message)
}

// Logf logs a formatted message
func (tl *TestLogger) Logf(format string, args ...interface{}) {
	tl.logs = append(tl.logs, fmt.Sprintf(format, args...))
}

// GetLogs returns all logged messages
func (tl *TestLogger) GetLogs() []string {
	return tl.logs
}

// Clear clears all logs
func (tl *TestLogger) Clear() {
	tl.logs = make([]string, 0)
}

// TestCleanup provides test cleanup utilities
type TestCleanup struct {
	cleanupFuncs []func()
}

// NewTestCleanup creates a new test cleanup
func NewTestCleanup() *TestCleanup {
	return &TestCleanup{
		cleanupFuncs: make([]func(), 0),
	}
}

// Add adds a cleanup function
func (tc *TestCleanup) Add(cleanup func()) {
	tc.cleanupFuncs = append(tc.cleanupFuncs, cleanup)
}

// Run runs all cleanup functions
func (tc *TestCleanup) Run() {
	for _, cleanup := range tc.cleanupFuncs {
		cleanup()
	}
}

// TestData provides test data utilities
type TestData struct {
	data map[string]interface{}
}

// NewTestData creates a new test data
func NewTestData() *TestData {
	return &TestData{
		data: make(map[string]interface{}),
	}
}

// Set sets a value
func (td *TestData) Set(key string, value interface{}) {
	td.data[key] = value
}

// Get gets a value
func (td *TestData) Get(key string) (interface{}, bool) {
	value, exists := td.data[key]
	return value, exists
}

// GetString gets a string value
func (td *TestData) GetString(key string) (string, bool) {
	value, exists := td.data[key]
	if !exists {
		return "", false
	}
	str, ok := value.(string)
	return str, ok
}

// GetInt gets an int value
func (td *TestData) GetInt(key string) (int, bool) {
	value, exists := td.data[key]
	if !exists {
		return 0, false
	}
	i, ok := value.(int)
	return i, ok
}

// GetBool gets a bool value
func (td *TestData) GetBool(key string) (bool, bool) {
	value, exists := td.data[key]
	if !exists {
		return false, false
	}
	b, ok := value.(bool)
	return b, ok
}

// Clear clears all data
func (td *TestData) Clear() {
	td.data = make(map[string]interface{})
}

// TestEnvironment provides test environment utilities
type TestEnvironment struct {
	env map[string]string
}

// NewTestEnvironment creates a new test environment
func NewTestEnvironment() *TestEnvironment {
	return &TestEnvironment{
		env: make(map[string]string),
	}
}

// Set sets an environment variable
func (te *TestEnvironment) Set(key, value string) {
	te.env[key] = value
}

// Get gets an environment variable
func (te *TestEnvironment) Get(key string) (string, bool) {
	value, exists := te.env[key]
	return value, exists
}

// Clear clears all environment variables
func (te *TestEnvironment) Clear() {
	te.env = make(map[string]string)
}

// TestConfigManager provides test configuration utilities
type TestConfigManager struct {
	config map[string]interface{}
}

// NewTestConfigManager creates a new test config manager
func NewTestConfigManager() *TestConfigManager {
	return &TestConfigManager{
		config: make(map[string]interface{}),
	}
}

// Set sets a config value
func (tc *TestConfigManager) Set(key string, value interface{}) {
	tc.config[key] = value
}

// Get gets a config value
func (tc *TestConfigManager) Get(key string) (interface{}, bool) {
	value, exists := tc.config[key]
	return value, exists
}

// GetString gets a string config value
func (tc *TestConfigManager) GetString(key string) (string, bool) {
	value, exists := tc.config[key]
	if !exists {
		return "", false
	}
	str, ok := value.(string)
	return str, ok
}

// GetInt gets an int config value
func (tc *TestConfigManager) GetInt(key string) (int, bool) {
	value, exists := tc.config[key]
	if !exists {
		return 0, false
	}
	i, ok := value.(int)
	return i, ok
}

// GetBool gets a bool config value
func (tc *TestConfigManager) GetBool(key string) (bool, bool) {
	value, exists := tc.config[key]
	if !exists {
		return false, false
	}
	b, ok := value.(bool)
	return b, ok
}

// Clear clears all config values
func (tc *TestConfigManager) Clear() {
	tc.config = make(map[string]interface{})
}
