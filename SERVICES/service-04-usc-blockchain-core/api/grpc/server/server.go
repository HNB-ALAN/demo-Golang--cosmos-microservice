package server

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"service-04/api/grpc/client"
	"service-04/api/grpc/interceptors"
	grpcreflection "service-04/api/grpc/reflection"
	"service-04/internal/application"
	"service-04/internal/infrastructure/auth"
	"service-04/internal/infrastructure/database"
	"service-04/internal/infrastructure/errors"
	"service-04/internal/infrastructure/kafka"
	"service-04/internal/infrastructure/metrics"
	servicemiddleware "service-04/internal/infrastructure/middleware"
	"service-04/internal/infrastructure/validation"
	proto "service-04/proto"

	"github.com/usc-platform/shared/config"
	"github.com/usc-platform/shared/constants"
	sharedgrpc "github.com/usc-platform/shared/grpc"
	"github.com/usc-platform/shared/logging"
	sharedmetrics "github.com/usc-platform/shared/metrics"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/app"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/storage"
)

// Server represents the gRPC server
type Server struct {
	config           *config.Config
	logger           *logging.Logger
	grpcServer       *grpc.Server
	listener         net.Listener
	ctx              context.Context
	postgresManager  *database.PostgreSQLManager
	poolManager      *database.PoolManager
	migrationManager *database.MigrationManager
	authService      *auth.JWTService
	validator        *validation.Validator
	metricsService   *metrics.MetricsService
	errorManager     *errors.ErrorManager
	cancel           context.CancelFunc
	middleware       *servicemiddleware.MiddlewareManager
	// Tier 4 services
	kafkaProducerManager *kafka.KafkaProducerManager
	kafkaConsumerManager *kafka.KafkaConsumerManager
	// gRPC Enhanced services
	grpcClientManager *client.ServiceBlockchainCoreClientManager
	grpcInterceptors  *interceptors.ServiceBlockchainCoreInterceptors
	grpcReflection    *grpcreflection.ServiceBlockchainCoreReflectionService

	// Application Container
	container *application.Container

	// Cosmos SDK Integration
	cosmosApp         *app.USCApp
	blockchainStorage *storage.StateManager
	rocksDBManager    *storage.RocksDBManager
	redisManager      *database.RedisManager

	// Metrics HTTP server shutdown function
	metricsHTTPServerShutdown func(ctx context.Context) error
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config, logger *logging.Logger, postgresManager *database.PostgreSQLManager, poolManager *database.PoolManager, migrationManager *database.MigrationManager, authService *auth.JWTService, validator *validation.Validator, metricsService *metrics.MetricsService, errorManager *errors.ErrorManager, kafkaProducerManager *kafka.KafkaProducerManager, kafkaConsumerManager *kafka.KafkaConsumerManager, grpcClientManager *client.ServiceBlockchainCoreClientManager, grpcInterceptors *interceptors.ServiceBlockchainCoreInterceptors, grpcReflection *grpcreflection.ServiceBlockchainCoreReflectionService, container *application.Container, cosmosApp *app.USCApp, blockchainStorage *storage.StateManager, rocksDBManager *storage.RocksDBManager, redisManager *database.RedisManager) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize gRPC middleware
	grpcMiddleware := servicemiddleware.NewMiddlewareManager(logger)

	return &Server{
		config:           cfg,
		logger:           logger,
		ctx:              ctx,
		cancel:           cancel,
		middleware:       grpcMiddleware,
		postgresManager:  postgresManager,
		poolManager:      poolManager,
		migrationManager: migrationManager,
		authService:      authService,
		validator:        validator,
		metricsService:   metricsService,
		errorManager:     errorManager,
		// Tier 4 services
		kafkaProducerManager: kafkaProducerManager,
		kafkaConsumerManager: kafkaConsumerManager,
		// gRPC Enhanced services
		grpcClientManager: grpcClientManager,
		grpcInterceptors:  grpcInterceptors,
		grpcReflection:    grpcReflection,
		// Application Container
		container: container,
		// Cosmos SDK Integration
		cosmosApp:         cosmosApp,
		blockchainStorage: blockchainStorage,
		rocksDBManager:    rocksDBManager,
		redisManager:      redisManager,
	}
}

// GetGRPCServer returns the gRPC server instance
func (s *Server) GetGRPCServer() *grpc.Server {
	return s.grpcServer
}

// GetConfig returns the server configuration
func (s *Server) GetConfig() *config.Config {
	return s.config
}

// GetLogger returns the server logger
func (s *Server) GetLogger() *logging.Logger {
	return s.logger
}

// GetMiddleware returns the server middleware
func (s *Server) GetMiddleware() *servicemiddleware.MiddlewareManager {
	return s.middleware
}

// GetPostgreSQLManager returns the postgres manager
func (s *Server) GetPostgreSQLManager() *database.PostgreSQLManager {
	return s.postgresManager
}

// GetAuthService returns the auth service
func (s *Server) GetAuthService() *auth.JWTService {
	return s.authService
}

// GetValidator returns the validator
func (s *Server) GetValidator() *validation.Validator {
	return s.validator
}

// GetMetricsService returns the metrics service
func (s *Server) GetMetricsService() *metrics.MetricsService {
	return s.metricsService
}

// GetErrorManager returns the error manager
func (s *Server) GetErrorManager() *errors.ErrorManager {
	return s.errorManager
}

// GetKafkaProducerManager returns the kafka producer manager
func (s *Server) GetKafkaProducerManager() *kafka.KafkaProducerManager {
	return s.kafkaProducerManager
}

// GetKafkaConsumerManager returns the kafka consumer manager
func (s *Server) GetKafkaConsumerManager() *kafka.KafkaConsumerManager {
	return s.kafkaConsumerManager
}

// GetCosmosApp returns the Cosmos SDK app instance
func (s *Server) GetCosmosApp() *app.USCApp {
	return s.cosmosApp
}

// GetBlockchainStorage returns the blockchain storage manager
func (s *Server) GetBlockchainStorage() *storage.StateManager {
	return s.blockchainStorage
}

// GetRocksDBManager returns the RocksDB manager
func (s *Server) GetRocksDBManager() *storage.RocksDBManager {
	return s.rocksDBManager
}

// GetRedisManager returns the Redis manager
func (s *Server) GetRedisManager() *database.RedisManager {
	return s.redisManager
}

// Start starts the gRPC server
func (s *Server) Start() error {
	// Get server address from config
	serverAddr := s.config.GetServerAddress()

	s.logger.Info("Starting gRPC server",
		logging.String("address", serverAddr),
		logging.String("service", "USC Blockchain Core Service"))

	// Initialize Cosmos SDK components
	if s.cosmosApp != nil {
		s.logger.Info("Initializing Cosmos SDK app",
			logging.String("chain_id", "usc-blockchain"),
			logging.String("service", "USC Blockchain Core Service"))
	}

	if s.blockchainStorage != nil {
		s.logger.Info("Initializing blockchain storage",
			logging.String("service", "USC Blockchain Core Service"))
	}

	if s.rocksDBManager != nil {
		s.logger.Info("Initializing RocksDB manager",
			logging.String("service", "USC Blockchain Core Service"))
	}

	if s.redisManager != nil {
		s.logger.Info("Initializing Redis manager",
			logging.String("service", "USC Blockchain Core Service"))
	}

	// Start metrics HTTP server if enabled
	if s.config.Metrics.Enabled {
		metricsAddr := fmt.Sprintf("0.0.0.0:%d", s.config.Metrics.Port)
		if shutdown, err := sharedmetrics.StartHTTPServer(metricsAddr, s.config.Metrics.Path, s.logger); err != nil {
			s.logger.Warn("Failed to start metrics HTTP server",
				logging.String("address", metricsAddr),
				logging.Error(err))
		} else {
			s.metricsHTTPServerShutdown = shutdown
		}
	}

	// Create listener
	listener, err := net.Listen("tcp", serverAddr)
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}
	s.listener = listener

	// Collect interceptors from configured middleware manager
	unaryInterceptors := append([]grpc.UnaryServerInterceptor{}, s.middleware.GetUnaryInterceptors()...)
	streamInterceptors := append([]grpc.StreamServerInterceptor{}, s.middleware.GetStreamInterceptors()...)

	// Get gRPC Enhanced interceptors
	enhancedUnaryInterceptors, enhancedStreamInterceptors := s.grpcInterceptors.GetUnaryServerInterceptors(), s.grpcInterceptors.GetStreamServerInterceptors()
	unaryInterceptors = append(unaryInterceptors, enhancedUnaryInterceptors...)
	streamInterceptors = append(streamInterceptors, enhancedStreamInterceptors...)

	// Create gRPC server with comprehensive middleware stack (chain interceptors)
	var serverOpts []grpc.ServerOption
	if len(unaryInterceptors) > 0 {
		serverOpts = append(serverOpts, grpc.ChainUnaryInterceptor(unaryInterceptors...))
	}
	if len(streamInterceptors) > 0 {
		serverOpts = append(serverOpts, grpc.ChainStreamInterceptor(streamInterceptors...))
	}
	s.grpcServer = grpc.NewServer(serverOpts...)

	// Register gRPC reflection service
	reflection.Register(s.grpcServer)

	// Register health check services
	s.registerHealthServices()

	// Register business services (to be implemented by each service)
	s.registerBusinessServices()

	s.logger.Info("gRPC server started successfully",
		logging.String("address", serverAddr),
		logging.String("service", "USC Blockchain Core Service"))

	// Start serving
	if err := s.grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("gRPC server failed: %w", err)
	}

	return nil
}

// Stop gracefully stops the gRPC server
func (s *Server) Stop() {
	s.logger.Info("Stopping gRPC server...",
		logging.String("service", "USC Blockchain Core Service"))

	// Cancel context
	s.cancel()

	// Graceful stop with timeout
	done := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(done)
	}()

	// Wait for graceful stop or timeout
	select {
	case <-done:
		s.logger.Info("gRPC server stopped gracefully",
			logging.String("service", "USC Blockchain Core Service"))
	case <-time.After(30 * time.Second):
		s.logger.Warn("gRPC server stop timeout, forcing stop",
			logging.String("service", "USC Blockchain Core Service"))
		s.grpcServer.Stop()
	}

	// Close metrics HTTP server
	if s.metricsHTTPServerShutdown != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.metricsHTTPServerShutdown(ctx)
	}

	// Close listener
	if s.listener != nil {
		s.listener.Close()
	}

	// Cleanup Cosmos SDK components
	if s.cosmosApp != nil {
		s.logger.Info("Stopping Cosmos SDK app",
			logging.String("service", "USC Blockchain Core Service"))
		if err := s.cosmosApp.Stop(); err != nil {
			s.logger.Warn("Error stopping Cosmos SDK app",
				logging.Error(err),
				logging.String("service", "USC Blockchain Core Service"))
		}
	}

	if s.rocksDBManager != nil {
		s.logger.Info("Closing RocksDB manager",
			logging.String("service", "USC Blockchain Core Service"))
		if err := s.rocksDBManager.Close(); err != nil {
			s.logger.Warn("Error closing RocksDB manager",
				logging.Error(err),
				logging.String("service", "USC Blockchain Core Service"))
		}
	}

	if s.redisManager != nil {
		s.logger.Info("Closing Redis manager",
			logging.String("service", "USC Blockchain Core Service"))
		if err := s.redisManager.Close(); err != nil {
			s.logger.Warn("Error closing Redis manager",
				logging.Error(err),
				logging.String("service", "USC Blockchain Core Service"))
		}
	}
}

// RegisterServices registers all gRPC services
func (s *Server) RegisterServices() {
	s.logger.Info("Registering gRPC services...",
		logging.String("service", "USC Blockchain Core Service"))

	// Health services are registered in registerHealthServices
	// Business services are registered in registerBusinessServices

	s.logger.Info("All gRPC services registered successfully",
		logging.String("service", "USC Blockchain Core Service"))
}

// registerHealthServices registers health check services using shared health
func (s *Server) registerHealthServices() {
	// Use shared gRPC health service registration
	sharedgrpc.RegisterHealthService(s.grpcServer, "USC Blockchain Core Service", constants.AppVersion)

	s.logger.Info("Health check services registered using shared health",
		logging.String("service", "USC Blockchain Core Service"))
}

// registerBusinessServices registers business-specific services using Container
func (s *Server) registerBusinessServices() {
	if s.container == nil {
		s.logger.Fatal("Container is not initialized")
		return
	}

	// Register all gRPC services using handlers from Container
	proto.RegisterBlockOperationsServiceServer(s.grpcServer, s.container.GetBlockHandlers())
	proto.RegisterTransactionOperationsServiceServer(s.grpcServer, s.container.GetTransactionHandlers())
	proto.RegisterUSCCoinOperationsServiceServer(s.grpcServer, s.container.GetUSCCoinHandlers())
	proto.RegisterSmartContractOperationsServiceServer(s.grpcServer, s.container.GetContractHandlers())
	proto.RegisterNFTTokenOperationsServiceServer(s.grpcServer, s.container.GetNFTHandlers())
	proto.RegisterCustomTokenOperationsServiceServer(s.grpcServer, s.container.GetTokenHandlers())
	proto.RegisterProductCertificateOperationsServiceServer(s.grpcServer, s.container.GetCertificateHandlers())
	proto.RegisterValidatorOperationsServiceServer(s.grpcServer, s.container.GetValidatorHandlers())
	proto.RegisterNetworkOperationsServiceServer(s.grpcServer, s.container.GetNetworkHandlers())
	proto.RegisterStreamingOperationsServiceServer(s.grpcServer, s.container.GetStreamHandlers())
	proto.RegisterStoreBridgeOperationsServiceServer(s.grpcServer, s.container.GetBridgeHandlers())
	proto.RegisterStoreNetworkOperationsServiceServer(s.grpcServer, s.container.GetStoreNetworkHandlers())

	s.logger.Info("Business services registered via Container",
		logging.String("service", "USC Blockchain Core Service"),
		logging.String("registered", "12 gRPC services"))
}
