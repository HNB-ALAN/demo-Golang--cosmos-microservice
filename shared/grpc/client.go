package grpc

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"github.com/usc-platform/shared/config"
	"github.com/usc-platform/shared/logging"
	"github.com/usc-platform/shared/metrics"
)

// ClientFactory creates gRPC clients
type ClientFactory struct {
	config *config.Config
	logger *logging.Logger
}

// NewClientFactory creates a new client factory
func NewClientFactory(cfg *config.Config, logger *logging.Logger) *ClientFactory {
	return &ClientFactory{
		config: cfg,
		logger: logger,
	}
}

// CreateClient creates a gRPC client connection
func (f *ClientFactory) CreateClient(address string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             3 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithUnaryInterceptor(UnaryClientInterceptor(f.logger)),
		grpc.WithStreamInterceptor(StreamClientInterceptor(f.logger)),
	}

	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", address, err)
	}

	f.logger.Info("gRPC client connected",
		logging.String("address", address),
	)

	return conn, nil
}

// CreateClientWithTimeout creates a gRPC client connection with timeout
func (f *ClientFactory) CreateClientWithTimeout(address string, timeout time.Duration) (*grpc.ClientConn, error) {
	_, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             3 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithUnaryInterceptor(UnaryClientInterceptor(f.logger)),
		grpc.WithStreamInterceptor(StreamClientInterceptor(f.logger)),
	}

	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", address, err)
	}

	f.logger.Info("gRPC client connected with timeout",
		logging.String("address", address),
		logging.Duration("timeout", timeout),
	)

	return conn, nil
}

// UnaryClientInterceptor creates a unary client interceptor
func UnaryClientInterceptor(logger *logging.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()

		logger.Debug("gRPC client request started",
			logging.String("method", method),
		)

		err := invoker(ctx, method, req, reply, cc, opts...)

		duration := time.Since(start)

		if err != nil {
			logger.Error("gRPC client request failed",
				logging.String("method", method),
				logging.Duration("duration", duration),
				logging.Error(err),
			)
		} else {
			logger.Debug("gRPC client request completed",
				logging.String("method", method),
				logging.Duration("duration", duration),
			)
		}

		return err
	}
}

// StreamClientInterceptor creates a stream client interceptor
func StreamClientInterceptor(logger *logging.Logger) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		start := time.Now()

		logger.Debug("gRPC client stream started",
			logging.String("method", method),
		)

		stream, err := streamer(ctx, desc, cc, method, opts...)

		duration := time.Since(start)

		if err != nil {
			logger.Error("gRPC client stream failed",
				logging.String("method", method),
				logging.Duration("duration", duration),
				logging.Error(err),
			)
		} else {
			logger.Debug("gRPC client stream completed",
				logging.String("method", method),
				logging.Duration("duration", duration),
			)
		}

		return stream, err
	}
}

// ConnectionPoolConfig represents connection pool configuration
type ConnectionPoolConfig struct {
	MaxConnections    int           `mapstructure:"max_connections"`
	MinConnections    int           `mapstructure:"min_connections"`
	MaxIdleTime       time.Duration `mapstructure:"max_idle_time"`
	ConnectionTimeout time.Duration `mapstructure:"connection_timeout"`
	RetryAttempts     int           `mapstructure:"retry_attempts"`
	RetryDelay        time.Duration `mapstructure:"retry_delay"`
}

// DefaultConnectionPoolConfig returns default connection pool configuration
func DefaultConnectionPoolConfig() ConnectionPoolConfig {
	return ConnectionPoolConfig{
		MaxConnections:    10,
		MinConnections:    2,
		MaxIdleTime:       5 * time.Minute,
		ConnectionTimeout: 30 * time.Second,
		RetryAttempts:     3,
		RetryDelay:        1 * time.Second,
	}
}

// ConnectionPool manages a pool of gRPC connections
type ConnectionPool struct {
	config      ConnectionPoolConfig
	factory     *ClientFactory
	connections map[string]*grpc.ClientConn
	lastUsed    map[string]time.Time
	mu          sync.RWMutex
}

// NewConnectionPool creates a new connection pool
func NewConnectionPool(config ConnectionPoolConfig, factory *ClientFactory) *ConnectionPool {
	return &ConnectionPool{
		config:      config,
		factory:     factory,
		connections: make(map[string]*grpc.ClientConn),
		lastUsed:    make(map[string]time.Time),
	}
}

// GetConnection gets a connection from the pool or creates a new one
func (p *ConnectionPool) GetConnection(name, address string) (*grpc.ClientConn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Check if connection exists and is healthy
	if conn, exists := p.connections[name]; exists {
		// Check if connection is still healthy
		if p.isConnectionHealthy(conn) {
			p.lastUsed[name] = time.Now()
			return conn, nil
		}
		// Remove unhealthy connection
		conn.Close()
		delete(p.connections, name)
		delete(p.lastUsed, name)
	}

	// Create new connection
	conn, err := p.factory.CreateClientWithTimeout(address, p.config.ConnectionTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection for %s: %w", name, err)
	}

	p.connections[name] = conn
	p.lastUsed[name] = time.Now()

	return conn, nil
}

// isConnectionHealthy checks if a connection is healthy
func (p *ConnectionPool) isConnectionHealthy(conn *grpc.ClientConn) bool {
	state := conn.GetState()
	return state.String() == "READY"
}

// CleanupIdleConnections removes idle connections
func (p *ConnectionPool) CleanupIdleConnections() {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	for name, lastUsed := range p.lastUsed {
		if now.Sub(lastUsed) > p.config.MaxIdleTime {
			if conn, exists := p.connections[name]; exists {
				conn.Close()
				delete(p.connections, name)
				delete(p.lastUsed, name)
			}
		}
	}
}

// Close closes all connections in the pool
func (p *ConnectionPool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	var errors []error
	for name, conn := range p.connections {
		if err := conn.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close connection %s: %w", name, err))
		}
	}

	p.connections = make(map[string]*grpc.ClientConn)
	p.lastUsed = make(map[string]time.Time)

	if len(errors) > 0 {
		return fmt.Errorf("errors closing connections: %v", errors)
	}

	return nil
}

// ClientManager manages gRPC client connections
type ClientManager struct {
	factory *ClientFactory
	clients map[string]*grpc.ClientConn
	pool    *ConnectionPool
}

// NewClientManager creates a new client manager
func NewClientManager(factory *ClientFactory) *ClientManager {
	poolConfig := DefaultConnectionPoolConfig()
	pool := NewConnectionPool(poolConfig, factory)

	return &ClientManager{
		factory: factory,
		clients: make(map[string]*grpc.ClientConn),
		pool:    pool,
	}
}

// NewClientManagerWithPool creates a new client manager with custom pool config
func NewClientManagerWithPool(factory *ClientFactory, poolConfig ConnectionPoolConfig) *ClientManager {
	pool := NewConnectionPool(poolConfig, factory)

	return &ClientManager{
		factory: factory,
		clients: make(map[string]*grpc.ClientConn),
		pool:    pool,
	}
}

// GetClient gets or creates a client connection using connection pool
func (m *ClientManager) GetClient(name, address string) (*grpc.ClientConn, error) {
	// Use connection pool for better performance
	return m.pool.GetConnection(name, address)
}

// GetClientLegacy gets or creates a client connection (legacy method)
func (m *ClientManager) GetClientLegacy(name, address string) (*grpc.ClientConn, error) {
	if conn, exists := m.clients[name]; exists {
		return conn, nil
	}

	conn, err := m.factory.CreateClient(address)
	if err != nil {
		return nil, err
	}

	m.clients[name] = conn
	return conn, nil
}

// CloseClient closes a client connection
func (m *ClientManager) CloseClient(name string) error {
	if conn, exists := m.clients[name]; exists {
		err := conn.Close()
		delete(m.clients, name)
		return err
	}

	return fmt.Errorf("client %s not found", name)
}

// CloseAll closes all client connections
func (m *ClientManager) CloseAll() error {
	var errors []error

	// Close pool connections
	if err := m.pool.Close(); err != nil {
		errors = append(errors, fmt.Errorf("failed to close connection pool: %w", err))
	}

	// Close legacy connections
	for name, conn := range m.clients {
		if err := conn.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close client %s: %w", name, err))
		}
	}

	m.clients = make(map[string]*grpc.ClientConn)

	if len(errors) > 0 {
		return fmt.Errorf("errors closing clients: %v", errors)
	}

	return nil
}

// CleanupIdleConnections cleans up idle connections in the pool
func (m *ClientManager) CleanupIdleConnections() {
	m.pool.CleanupIdleConnections()
}

// GetClientNames returns the names of all clients
func (m *ClientManager) GetClientNames() []string {
	names := make([]string, 0, len(m.clients))
	for name := range m.clients {
		names = append(names, name)
	}
	return names
}

// EnhancedGRPCClient represents an enhanced gRPC client with advanced features
type EnhancedGRPCClient struct {
	config           *config.Config
	logger           *logging.Logger
	pool             *ConnectionPool
	retryConfig      RetryConfig
	circuitBreaker   *CircuitBreaker
	loadBalancer     *LoadBalancer
	metricsCollector *metrics.PrometheusMetricsCollector
	clients          map[string]*grpc.ClientConn
	mu               sync.RWMutex
}

// EnhancedClientConfig represents configuration for enhanced gRPC client
type EnhancedClientConfig struct {
	ConnectionPool       ConnectionPoolConfig
	Retry                RetryConfig
	CircuitBreaker       CircuitBreakerConfig
	LoadBalancer         LoadBalancerConfig
	Addresses            []string
	EnableMetrics        bool
	EnableRetry          bool
	EnableCircuitBreaker bool
	EnableLoadBalancer   bool
}

// DefaultEnhancedClientConfig returns default enhanced client configuration
func DefaultEnhancedClientConfig() EnhancedClientConfig {
	return EnhancedClientConfig{
		ConnectionPool:       DefaultConnectionPoolConfig(),
		Retry:                DefaultRetryConfig(),
		CircuitBreaker:       DefaultCircuitBreakerConfig(),
		LoadBalancer:         DefaultLoadBalancerConfig(),
		Addresses:            []string{"localhost:9090"},
		EnableMetrics:        true,
		EnableRetry:          true,
		EnableCircuitBreaker: true,
		EnableLoadBalancer:   true,
	}
}

// NewEnhancedGRPCClient creates a new enhanced gRPC client
func NewEnhancedGRPCClient(cfg *config.Config, logger *logging.Logger, clientConfig EnhancedClientConfig) *EnhancedGRPCClient {
	// Create client factory
	factory := NewClientFactory(cfg, logger)

	// Create connection pool
	pool := NewConnectionPool(clientConfig.ConnectionPool, factory)

	// Create circuit breaker
	var circuitBreaker *CircuitBreaker
	if clientConfig.EnableCircuitBreaker {
		circuitBreaker = NewCircuitBreaker(clientConfig.CircuitBreaker)
	}

	// Create load balancer
	var loadBalancer *LoadBalancer
	if clientConfig.EnableLoadBalancer {
		loadBalancer = NewLoadBalancer(clientConfig.LoadBalancer, clientConfig.Addresses)
	}

	// Create metrics collector
	var metricsCollector *metrics.PrometheusMetricsCollector
	if clientConfig.EnableMetrics {
		metricsCollector = metrics.NewPrometheusMetricsCollector("usc", "grpc_client")
	}

	return &EnhancedGRPCClient{
		config:           cfg,
		logger:           logger,
		pool:             pool,
		retryConfig:      clientConfig.Retry,
		circuitBreaker:   circuitBreaker,
		loadBalancer:     loadBalancer,
		metricsCollector: metricsCollector,
		clients:          make(map[string]*grpc.ClientConn),
	}
}

// CreateEnhancedClient creates a gRPC client connection with enhanced features
func (c *EnhancedGRPCClient) CreateEnhancedClient(name, address string) (*grpc.ClientConn, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if client already exists
	if conn, exists := c.clients[name]; exists {
		if c.isConnectionHealthy(conn) {
			return conn, nil
		}
		// Remove unhealthy connection
		conn.Close()
		delete(c.clients, name)
	}

	// Create interceptors
	var interceptors []grpc.UnaryClientInterceptor

	// Add basic logging interceptor
	interceptors = append(interceptors, UnaryClientInterceptor(c.logger))

	// Add retry interceptor
	if c.retryConfig.MaxAttempts > 0 {
		interceptors = append(interceptors, UnaryClientRetryInterceptor(c.retryConfig, c.logger))
	}

	// Add circuit breaker interceptor
	if c.circuitBreaker != nil {
		interceptors = append(interceptors, CircuitBreakerInterceptor(c.circuitBreaker, c.logger))
	}

	// Add load balancing interceptor
	if c.loadBalancer != nil {
		interceptors = append(interceptors, LoadBalancingInterceptor(c.loadBalancer, c.logger))
	}

	// Add metrics interceptor
	if c.metricsCollector != nil {
		interceptors = append(interceptors, MetricsInterceptor(c.metricsCollector, c.logger))
	}

	// Create client options
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             3 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithChainUnaryInterceptor(interceptors...),
	}

	// Create connection
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create enhanced client for %s: %w", name, err)
	}

	c.clients[name] = conn
	c.logger.Info("Enhanced gRPC client created",
		logging.String("name", name),
		logging.String("address", address),
	)

	return conn, nil
}

// GetClient gets a client connection by name
func (c *EnhancedGRPCClient) GetClient(name string) (*grpc.ClientConn, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	conn, exists := c.clients[name]
	return conn, exists
}

// GetClientWithLoadBalancing gets a client connection using load balancing
func (c *EnhancedGRPCClient) GetClientWithLoadBalancing(name string) (*grpc.ClientConn, error) {
	if c.loadBalancer == nil {
		return nil, fmt.Errorf("load balancer not configured")
	}

	// Get next address from load balancer
	address := c.loadBalancer.GetNextAddress()
	if address == "" {
		return nil, fmt.Errorf("no healthy addresses available")
	}

	// Create or get client for this address
	return c.CreateEnhancedClient(name, address)
}

// isConnectionHealthy checks if a connection is healthy
func (c *EnhancedGRPCClient) isConnectionHealthy(conn *grpc.ClientConn) bool {
	state := conn.GetState()
	return state.String() == "READY"
}

// HealthCheck performs health check on all clients
func (c *EnhancedGRPCClient) HealthCheck(ctx context.Context) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var errors []error
	for name, conn := range c.clients {
		if !c.isConnectionHealthy(conn) {
			errors = append(errors, fmt.Errorf("client %s is not healthy", name))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("health check failed: %v", errors)
	}

	return nil
}

// GetMetrics returns client metrics
func (c *EnhancedGRPCClient) GetMetrics() map[string]interface{} {
	metrics := make(map[string]interface{})

	// Connection pool metrics
	metrics["connection_pool"] = map[string]interface{}{
		"total_connections":   len(c.clients),
		"healthy_connections": c.getHealthyConnectionCount(),
	}

	// Circuit breaker metrics
	if c.circuitBreaker != nil {
		metrics["circuit_breaker"] = map[string]interface{}{
			"state": c.circuitBreaker.GetState(),
		}
	}

	// Load balancer metrics
	if c.loadBalancer != nil {
		metrics["load_balancer"] = map[string]interface{}{
			"total_addresses":   len(c.loadBalancer.addresses),
			"healthy_addresses": c.getHealthyAddressCount(),
		}
	}

	return metrics
}

// getHealthyConnectionCount returns the number of healthy connections
func (c *EnhancedGRPCClient) getHealthyConnectionCount() int {
	count := 0
	for _, conn := range c.clients {
		if c.isConnectionHealthy(conn) {
			count++
		}
	}
	return count
}

// getHealthyAddressCount returns the number of healthy addresses
func (c *EnhancedGRPCClient) getHealthyAddressCount() int {
	if c.loadBalancer == nil {
		return 0
	}

	count := 0
	for _, address := range c.loadBalancer.addresses {
		if c.loadBalancer.IsHealthy(address) {
			count++
		}
	}
	return count
}

// Close closes all client connections
func (c *EnhancedGRPCClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var errors []error

	// Close all clients
	for name, conn := range c.clients {
		if err := conn.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close client %s: %w", name, err))
		}
	}

	// Close connection pool
	if err := c.pool.Close(); err != nil {
		errors = append(errors, fmt.Errorf("failed to close connection pool: %w", err))
	}

	c.clients = make(map[string]*grpc.ClientConn)

	if len(errors) > 0 {
		return fmt.Errorf("errors closing enhanced client: %v", errors)
	}

	c.logger.Info("Enhanced gRPC client closed")
	return nil
}

// CleanupIdleConnections cleans up idle connections
func (c *EnhancedGRPCClient) CleanupIdleConnections() {
	c.pool.CleanupIdleConnections()
}
