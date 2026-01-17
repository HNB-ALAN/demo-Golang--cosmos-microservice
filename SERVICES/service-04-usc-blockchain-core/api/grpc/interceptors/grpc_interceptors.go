package interceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	sharedgrpc "github.com/usc-platform/shared/grpc"
	"github.com/usc-platform/shared/logging"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/app"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/storage"
)

// ServiceBlockchainCoreInterceptors manages gRPC interceptors for USC Blockchain Core Service
type ServiceBlockchainCoreInterceptors struct {
	logger *logging.Logger
	config *ServiceBlockchainCoreInterceptorConfig
	// Cosmos SDK Integration
	cosmosApp         *app.USCApp
	blockchainStorage *storage.StateManager
	rocksDBManager    *storage.RocksDBManager
}

// ServiceBlockchainCoreInterceptorConfig represents interceptor configuration for USC Blockchain Core Service
type ServiceBlockchainCoreInterceptorConfig struct {
	EnableLogging        bool          `mapstructure:"enable_logging"`
	EnableRecovery       bool          `mapstructure:"enable_recovery"`
	EnableAuth           bool          `mapstructure:"enable_auth"`
	EnableRateLimit      bool          `mapstructure:"enable_rate_limit"`
	EnableTimeout        bool          `mapstructure:"enable_timeout"`
	RequestTimeout       time.Duration `mapstructure:"request_timeout"`
	MaxRequestsPerMin    int           `mapstructure:"max_requests_per_min"`
	RetryConfig          sharedgrpc.RetryConfig
	CircuitBreakerConfig sharedgrpc.CircuitBreakerConfig
}

// DefaultServiceBlockchainCoreInterceptorConfig returns default interceptor configuration
func DefaultServiceBlockchainCoreInterceptorConfig() *ServiceBlockchainCoreInterceptorConfig {
	return &ServiceBlockchainCoreInterceptorConfig{
		EnableLogging:        true,
		EnableRecovery:       true,
		EnableAuth:           false, // Set to true for services that need auth
		EnableRateLimit:      false, // Set to true for services that need rate limiting
		EnableTimeout:        true,
		RequestTimeout:       30 * time.Second,
		MaxRequestsPerMin:    1000,
		RetryConfig:          sharedgrpc.DefaultRetryConfig(),
		CircuitBreakerConfig: sharedgrpc.DefaultCircuitBreakerConfig(),
	}
}

// NewServiceBlockchainCoreInterceptors creates a new interceptor manager for USC Blockchain Core Service
func NewServiceBlockchainCoreInterceptors(logger *logging.Logger, config *ServiceBlockchainCoreInterceptorConfig, cosmosApp *app.USCApp, blockchainStorage *storage.StateManager, rocksDBManager *storage.RocksDBManager) *ServiceBlockchainCoreInterceptors {
	if config == nil {
		config = DefaultServiceBlockchainCoreInterceptorConfig()
	}

	return &ServiceBlockchainCoreInterceptors{
		logger: logger,
		config: config,
		// Cosmos SDK Integration
		cosmosApp:         cosmosApp,
		blockchainStorage: blockchainStorage,
		rocksDBManager:    rocksDBManager,
	}
}

// GetUnaryServerInterceptors returns unary server interceptors for USC Blockchain Core Service
func (i *ServiceBlockchainCoreInterceptors) GetUnaryServerInterceptors() []grpc.UnaryServerInterceptor {
	var interceptors []grpc.UnaryServerInterceptor

	// Recovery interceptor (should be first)
	if i.config.EnableRecovery {
		interceptors = append(interceptors, sharedgrpc.RecoveryInterceptor(i.logger))
	}

	// Timeout interceptor
	if i.config.EnableTimeout {
		interceptors = append(interceptors, sharedgrpc.TimeoutInterceptor(i.config.RequestTimeout))
	}

	// Auth interceptor
	if i.config.EnableAuth {
		interceptors = append(interceptors, i.createAuthInterceptor())
	}

	// Rate limit interceptor
	if i.config.EnableRateLimit {
		interceptors = append(interceptors, i.createRateLimitInterceptor())
	}

	// Logging interceptor (should be last)
	if i.config.EnableLogging {
		interceptors = append(interceptors, sharedgrpc.UnaryServerInterceptor(i.logger))
	}

	return interceptors
}

// GetStreamServerInterceptors returns stream server interceptors for USC Blockchain Core Service
func (i *ServiceBlockchainCoreInterceptors) GetStreamServerInterceptors() []grpc.StreamServerInterceptor {
	var interceptors []grpc.StreamServerInterceptor

	// Logging interceptor
	if i.config.EnableLogging {
		interceptors = append(interceptors, sharedgrpc.StreamServerInterceptor(i.logger))
	}

	return interceptors
}

// GetUnaryClientInterceptors returns unary client interceptors for USC Blockchain Core Service
func (i *ServiceBlockchainCoreInterceptors) GetUnaryClientInterceptors() []grpc.UnaryClientInterceptor {
	var interceptors []grpc.UnaryClientInterceptor

	// Basic logging interceptor
	interceptors = append(interceptors, sharedgrpc.UnaryClientInterceptor(i.logger))

	// Retry interceptor
	interceptors = append(interceptors, sharedgrpc.UnaryClientRetryInterceptor(i.config.RetryConfig, i.logger))

	return interceptors
}

// GetStreamClientInterceptors returns stream client interceptors for USC Blockchain Core Service
func (i *ServiceBlockchainCoreInterceptors) GetStreamClientInterceptors() []grpc.StreamClientInterceptor {
	var interceptors []grpc.StreamClientInterceptor

	// Basic logging interceptor
	interceptors = append(interceptors, sharedgrpc.StreamClientInterceptor(i.logger))

	return interceptors
}

// createAuthInterceptor creates a custom auth interceptor for USC Blockchain Core Service
func (i *ServiceBlockchainCoreInterceptors) createAuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Skip auth for health checks
		if info.FullMethod == "/grpc.health.v1.Health/Check" ||
			info.FullMethod == "/grpc.health.v1.Health/Watch" {
			return handler(ctx, req)
		}

		// Check if request comes from internal service (Gateway, other services)
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			// Allow requests from internal services (service-to-service auth)
			if serviceName := md.Get("x-service-name"); len(serviceName) > 0 &&
				(serviceName[0] == "service-01-gateway" || serviceName[0] == "gateway") {
				i.logger.Debug("Allowing internal service request without auth",
					logging.String("method", info.FullMethod),
					logging.String("source_service", serviceName[0]))
				return handler(ctx, req)
			}
		}

		// Custom auth logic for USC Blockchain Core Service
		// This is a placeholder - implement actual auth logic based on service requirements

		i.logger.Debug("USC Blockchain Core Service auth check",
			logging.String("method", info.FullMethod),
		)

		// For now, allow all requests
		// In production, implement proper JWT validation, API key validation, etc.

		return handler(ctx, req)
	}
}

// createRateLimitInterceptor creates a custom rate limit interceptor for USC Blockchain Core Service
func (i *ServiceBlockchainCoreInterceptors) createRateLimitInterceptor() grpc.UnaryServerInterceptor {
	// Rate limiting is configured in middleware config (rate_limit_per_sec: 100)
	// Gateway service handles rate limiting for external requests
	// Service-04 rate limiting is primarily for internal protection
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Rate limiting logic:
		// - Gateway service handles rate limiting for external requests
		// - Service-04 applies additional rate limiting for internal protection
		// - Config: middleware.rate_limit_per_sec (default: 100 req/sec)
		// - Implementation: Use shared library rate limiter or Redis-based limiter

		i.logger.Debug("USC Blockchain Core Service rate limit check",
			logging.String("method", info.FullMethod),
			logging.String("note", "Rate limiting handled by Gateway and middleware config"))

		// Allow request (rate limiting enforced at Gateway level)
		// Additional rate limiting can be added here using Redis or in-memory counters if needed
		return handler(ctx, req)
	}
}

// GetCosmosApp returns the Cosmos SDK app instance
func (i *ServiceBlockchainCoreInterceptors) GetCosmosApp() *app.USCApp {
	return i.cosmosApp
}

// GetBlockchainStorage returns the blockchain storage manager
func (i *ServiceBlockchainCoreInterceptors) GetBlockchainStorage() *storage.StateManager {
	return i.blockchainStorage
}

// GetRocksDBManager returns the RocksDB manager
func (i *ServiceBlockchainCoreInterceptors) GetRocksDBManager() *storage.RocksDBManager {
	return i.rocksDBManager
}

// ServiceBlockchainCoreEnhancedInterceptors provides enhanced interceptor functionality for USC Blockchain Core Service
type ServiceBlockchainCoreEnhancedInterceptors struct {
	*ServiceBlockchainCoreInterceptors
	circuitBreaker *sharedgrpc.CircuitBreaker
	// loadBalancer   *sharedgrpc.LoadBalancer // TODO: Implement load balancer functionality
}

// NewServiceBlockchainCoreEnhancedInterceptors creates enhanced interceptors for USC Blockchain Core Service
func NewServiceBlockchainCoreEnhancedInterceptors(logger *logging.Logger, config *ServiceBlockchainCoreInterceptorConfig, cosmosApp *app.USCApp, blockchainStorage *storage.StateManager, rocksDBManager *storage.RocksDBManager) *ServiceBlockchainCoreEnhancedInterceptors {
	base := NewServiceBlockchainCoreInterceptors(logger, config, cosmosApp, blockchainStorage, rocksDBManager)

	// Create circuit breaker
	circuitBreaker := sharedgrpc.NewCircuitBreaker(config.CircuitBreakerConfig)

	// Create load balancer (if needed)
	// addresses := []string{"localhost:8002", "localhost:8003"} // Example addresses
	// loadBalancer := grpc.NewLoadBalancer(grpc.DefaultLoadBalancerConfig(), addresses)

	return &ServiceBlockchainCoreEnhancedInterceptors{
		ServiceBlockchainCoreInterceptors: base,
		circuitBreaker:                    circuitBreaker,
		// loadBalancer:   loadBalancer,
	}
}

// GetEnhancedUnaryClientInterceptors returns enhanced unary client interceptors
func (ei *ServiceBlockchainCoreEnhancedInterceptors) GetEnhancedUnaryClientInterceptors() []grpc.UnaryClientInterceptor {
	var interceptors []grpc.UnaryClientInterceptor

	// Basic logging
	interceptors = append(interceptors, sharedgrpc.UnaryClientInterceptor(ei.logger))

	// Circuit breaker
	if ei.circuitBreaker != nil {
		interceptors = append(interceptors, sharedgrpc.CircuitBreakerInterceptor(ei.circuitBreaker, ei.logger))
	}

	// Retry
	interceptors = append(interceptors, sharedgrpc.UnaryClientRetryInterceptor(ei.config.RetryConfig, ei.logger))

	return interceptors
}

// GetCircuitBreakerState returns the current circuit breaker state
func (ei *ServiceBlockchainCoreEnhancedInterceptors) GetCircuitBreakerState() sharedgrpc.CircuitBreakerState {
	if ei.circuitBreaker == nil {
		return sharedgrpc.StateClosed
	}
	return ei.circuitBreaker.GetState()
}

// ServiceBlockchainCoreInterceptorManager manages all interceptors for USC Blockchain Core Service
type ServiceBlockchainCoreInterceptorManager struct {
	interceptors *ServiceBlockchainCoreInterceptors
	enhanced     *ServiceBlockchainCoreEnhancedInterceptors
	logger       *logging.Logger
}

// NewServiceBlockchainCoreInterceptorManager creates a new interceptor manager
func NewServiceBlockchainCoreInterceptorManager(logger *logging.Logger, config *ServiceBlockchainCoreInterceptorConfig, cosmosApp *app.USCApp, blockchainStorage *storage.StateManager, rocksDBManager *storage.RocksDBManager) *ServiceBlockchainCoreInterceptorManager {
	return &ServiceBlockchainCoreInterceptorManager{
		interceptors: NewServiceBlockchainCoreInterceptors(logger, config, cosmosApp, blockchainStorage, rocksDBManager),
		enhanced:     NewServiceBlockchainCoreEnhancedInterceptors(logger, config, cosmosApp, blockchainStorage, rocksDBManager),
		logger:       logger,
	}
}

// GetServerInterceptors returns server interceptors
func (m *ServiceBlockchainCoreInterceptorManager) GetServerInterceptors() ([]grpc.UnaryServerInterceptor, []grpc.StreamServerInterceptor) {
	return m.interceptors.GetUnaryServerInterceptors(), m.interceptors.GetStreamServerInterceptors()
}

// GetClientInterceptors returns client interceptors
func (m *ServiceBlockchainCoreInterceptorManager) GetClientInterceptors() ([]grpc.UnaryClientInterceptor, []grpc.StreamClientInterceptor) {
	return m.interceptors.GetUnaryClientInterceptors(), m.interceptors.GetStreamClientInterceptors()
}

// GetEnhancedClientInterceptors returns enhanced client interceptors
func (m *ServiceBlockchainCoreInterceptorManager) GetEnhancedClientInterceptors() []grpc.UnaryClientInterceptor {
	return m.enhanced.GetEnhancedUnaryClientInterceptors()
}

// GetMetrics returns interceptor metrics
func (m *ServiceBlockchainCoreInterceptorManager) GetMetrics() map[string]interface{} {
	metrics := make(map[string]interface{})
	metrics["service"] = "USC Blockchain Core Service"
	metrics["circuit_breaker_state"] = m.enhanced.GetCircuitBreakerState()
	metrics["config"] = map[string]interface{}{
		"enable_logging":    m.interceptors.config.EnableLogging,
		"enable_recovery":   m.interceptors.config.EnableRecovery,
		"enable_auth":       m.interceptors.config.EnableAuth,
		"enable_rate_limit": m.interceptors.config.EnableRateLimit,
		"enable_timeout":    m.interceptors.config.EnableTimeout,
		"request_timeout":   m.interceptors.config.RequestTimeout.String(),
	}
	return metrics
}
