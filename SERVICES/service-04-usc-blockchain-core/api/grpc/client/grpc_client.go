package client

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/usc-platform/shared/config"
	sharedgrpc "github.com/usc-platform/shared/grpc"
	"github.com/usc-platform/shared/logging"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/app"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/storage"
)

// ServiceBlockchainCoreClientManager manages gRPC client connections for USC Blockchain Core Service
type ServiceBlockchainCoreClientManager struct {
	config  *config.Config
	logger  *logging.Logger
	factory *sharedgrpc.ClientFactory
	pool    *sharedgrpc.ConnectionPool
	clients map[string]*grpc.ClientConn
	mu      sync.RWMutex
	// Cosmos SDK Integration
	cosmosApp         *app.USCApp
	blockchainStorage *storage.StateManager
	rocksDBManager    *storage.RocksDBManager
}

// NewServiceBlockchainCoreClientManager creates a new client manager for USC Blockchain Core Service
func NewServiceBlockchainCoreClientManager(cfg *config.Config, logger *logging.Logger, cosmosApp *app.USCApp, blockchainStorage *storage.StateManager, rocksDBManager *storage.RocksDBManager) *ServiceBlockchainCoreClientManager {
	// Create client factory
	factory := sharedgrpc.NewClientFactory(cfg, logger)

	// Create connection pool with default config
	poolConfig := sharedgrpc.DefaultConnectionPoolConfig()
	pool := sharedgrpc.NewConnectionPool(poolConfig, factory)

	return &ServiceBlockchainCoreClientManager{
		config:  cfg,
		logger:  logger,
		factory: factory,
		pool:    pool,
		clients: make(map[string]*grpc.ClientConn),
		// Cosmos SDK Integration
		cosmosApp:         cosmosApp,
		blockchainStorage: blockchainStorage,
		rocksDBManager:    rocksDBManager,
	}
}

// GetClient gets or creates a client connection for USC Blockchain Core Service
func (m *ServiceBlockchainCoreClientManager) GetClient(name, address string) (*grpc.ClientConn, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if client already exists and is healthy
	if conn, exists := m.clients[name]; exists {
		if m.isConnectionHealthy(conn) {
			m.logger.Debug("Using existing gRPC client connection",
				logging.String("service", "USC Blockchain Core Service"),
				logging.String("name", name),
				logging.String("address", address),
			)
			return conn, nil
		}
		// Remove unhealthy connection
		conn.Close()
		delete(m.clients, name)
	}

	// Create new connection using connection pool
	conn, err := m.pool.GetConnection(name, address)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection for %s: %w", name, err)
	}

	m.clients[name] = conn
	m.logger.Info("Created new gRPC client connection",
		logging.String("service", "USC Blockchain Core Service"),
		logging.String("name", name),
		logging.String("address", address),
	)

	return conn, nil
}

// GetClientWithTimeout gets or creates a client connection with timeout
func (m *ServiceBlockchainCoreClientManager) GetClientWithTimeout(name, address string, timeout time.Duration) (*grpc.ClientConn, error) {
	// Create connection with timeout
	conn, err := m.factory.CreateClientWithTimeout(address, timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection with timeout for %s: %w", name, err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.clients[name] = conn
	m.logger.Info("Created gRPC client connection with timeout",
		logging.String("service", "USC Blockchain Core Service"),
		logging.String("name", name),
		logging.String("address", address),
		logging.Duration("timeout", timeout),
	)

	return conn, nil
}

// isConnectionHealthy checks if a connection is healthy
func (m *ServiceBlockchainCoreClientManager) isConnectionHealthy(conn *grpc.ClientConn) bool {
	state := conn.GetState()
	return state.String() == "READY"
}

// HealthCheck performs health check on all clients
func (m *ServiceBlockchainCoreClientManager) HealthCheck(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var errors []error
	for name, conn := range m.clients {
		if !m.isConnectionHealthy(conn) {
			errors = append(errors, fmt.Errorf("client %s is not healthy", name))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("health check failed for USC Blockchain Core Service: %v", errors)
	}

	m.logger.Debug("All gRPC clients healthy",
		logging.String("service", "USC Blockchain Core Service"),
		logging.Int("client_count", len(m.clients)),
	)

	return nil
}

// GetClientNames returns the names of all clients
func (m *ServiceBlockchainCoreClientManager) GetClientNames() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.clients))
	for name := range m.clients {
		names = append(names, name)
	}
	return names
}

// GetMetrics returns client metrics
func (m *ServiceBlockchainCoreClientManager) GetMetrics() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	metrics := make(map[string]interface{})
	metrics["service"] = "USC Blockchain Core Service"
	metrics["total_clients"] = len(m.clients)
	metrics["healthy_clients"] = m.getHealthyClientCount()

	// Add client details
	clientDetails := make(map[string]interface{})
	for name, conn := range m.clients {
		clientDetails[name] = map[string]interface{}{
			"healthy": m.isConnectionHealthy(conn),
			"state":   conn.GetState().String(),
		}
	}
	metrics["clients"] = clientDetails

	return metrics
}

// getHealthyClientCount returns the number of healthy clients
func (m *ServiceBlockchainCoreClientManager) getHealthyClientCount() int {
	count := 0
	for _, conn := range m.clients {
		if m.isConnectionHealthy(conn) {
			count++
		}
	}
	return count
}

// CleanupIdleConnections cleans up idle connections
func (m *ServiceBlockchainCoreClientManager) CleanupIdleConnections() {
	m.pool.CleanupIdleConnections()
	m.logger.Debug("Cleaned up idle gRPC connections",
		logging.String("service", "USC Blockchain Core Service"),
	)
}

// CloseClient closes a specific client connection
func (m *ServiceBlockchainCoreClientManager) CloseClient(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if conn, exists := m.clients[name]; exists {
		err := conn.Close()
		delete(m.clients, name)
		m.logger.Info("Closed gRPC client connection",
			logging.String("service", "USC Blockchain Core Service"),
			logging.String("name", name),
		)
		return err
	}

	return fmt.Errorf("client %s not found", name)
}

// CloseAll closes all client connections
func (m *ServiceBlockchainCoreClientManager) CloseAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errors []error

	// Close all clients
	for name, conn := range m.clients {
		if err := conn.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close client %s: %w", name, err))
		}
	}

	// Close connection pool
	if err := m.pool.Close(); err != nil {
		errors = append(errors, fmt.Errorf("failed to close connection pool: %w", err))
	}

	m.clients = make(map[string]*grpc.ClientConn)

	if len(errors) > 0 {
		return fmt.Errorf("errors closing USC Blockchain Core Service clients: %v", errors)
	}

	m.logger.Info("All gRPC client connections closed",
		logging.String("service", "USC Blockchain Core Service"),
	)

	return nil
}

// GetCosmosApp returns the Cosmos SDK app instance
func (m *ServiceBlockchainCoreClientManager) GetCosmosApp() *app.USCApp {
	return m.cosmosApp
}

// GetBlockchainStorage returns the blockchain storage manager
func (m *ServiceBlockchainCoreClientManager) GetBlockchainStorage() *storage.StateManager {
	return m.blockchainStorage
}

// GetRocksDBManager returns the RocksDB manager
func (m *ServiceBlockchainCoreClientManager) GetRocksDBManager() *storage.RocksDBManager {
	return m.rocksDBManager
}

// ServiceBlockchainCoreServiceClients represents all service clients for USC Blockchain Core Service
type ServiceBlockchainCoreServiceClients struct {
	manager *ServiceBlockchainCoreClientManager
}

// NewServiceBlockchainCoreServiceClients creates a new service clients instance
func NewServiceBlockchainCoreServiceClients(cfg *config.Config, logger *logging.Logger, cosmosApp *app.USCApp, blockchainStorage *storage.StateManager, rocksDBManager *storage.RocksDBManager) *ServiceBlockchainCoreServiceClients {
	manager := NewServiceBlockchainCoreClientManager(cfg, logger, cosmosApp, blockchainStorage, rocksDBManager)

	return &ServiceBlockchainCoreServiceClients{
		manager: manager,
	}
}

// GetServiceClient gets a client connection to a specific service
func (c *ServiceBlockchainCoreServiceClients) GetServiceClient(serviceName, address string) (*grpc.ClientConn, error) {
	return c.manager.GetClient(serviceName, address)
}

// HealthCheck performs health check on all service clients
func (c *ServiceBlockchainCoreServiceClients) HealthCheck(ctx context.Context) error {
	return c.manager.HealthCheck(ctx)
}

// GetMetrics returns metrics for all service clients
func (c *ServiceBlockchainCoreServiceClients) GetMetrics() map[string]interface{} {
	return c.manager.GetMetrics()
}

// Close closes all service clients
func (c *ServiceBlockchainCoreServiceClients) Close() error {
	return c.manager.CloseAll()
}
