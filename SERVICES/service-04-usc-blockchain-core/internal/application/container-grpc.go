package application

import (
	"context"

	// Business layer imports
	blockbiz "service-04/internal/application/business/block_operations"
	tokenbiz "service-04/internal/application/business/custom_token_operations"
	netbiz "service-04/internal/application/business/network_operations"
	nftbiz "service-04/internal/application/business/nft_token_operations"
	certbiz "service-04/internal/application/business/product_certificate_operations"
	contractbiz "service-04/internal/application/business/smart_contract_operations"
	bridgebiz "service-04/internal/application/business/store_bridge_operations"
	storenetbiz "service-04/internal/application/business/store_network_operations"
	streambiz "service-04/internal/application/business/streaming_operations"
	txbiz "service-04/internal/application/business/transaction_operations"
	uscbiz "service-04/internal/application/business/usc_coin_operations"
	valbiz "service-04/internal/application/business/validator_operations"

	// Handlers layer imports
	blockhandlers "service-04/internal/application/handlers/block_operations"
	tokenhandlers "service-04/internal/application/handlers/custom_token_operations"
	nethandlers "service-04/internal/application/handlers/network_operations"
	nfthandlers "service-04/internal/application/handlers/nft_token_operations"
	certhandlers "service-04/internal/application/handlers/product_certificate_operations"
	contracthandlers "service-04/internal/application/handlers/smart_contract_operations"
	bridgehandlers "service-04/internal/application/handlers/store_bridge_operations"
	storenethandlers "service-04/internal/application/handlers/store_network_operations"
	streamhandlers "service-04/internal/application/handlers/streaming_operations"
	txhandlers "service-04/internal/application/handlers/transaction_operations"
	uschandlers "service-04/internal/application/handlers/usc_coin_operations"
	valhandlers "service-04/internal/application/handlers/validator_operations"

	// Repository layer imports
	blockrepo "service-04/internal/application/repository/block_operations"
	tokenrepo "service-04/internal/application/repository/custom_token_operations"
	netrepo "service-04/internal/application/repository/network_operations"
	nftrepo "service-04/internal/application/repository/nft_token_operations"
	certrepo "service-04/internal/application/repository/product_certificate_operations"
	contractrepo "service-04/internal/application/repository/smart_contract_operations"
	bridgerepo "service-04/internal/application/repository/store_bridge_operations"
	storenetrepo "service-04/internal/application/repository/store_network_operations"
	streamrepo "service-04/internal/application/repository/streaming_operations"
	txrepo "service-04/internal/application/repository/transaction_operations"
	uscrepo "service-04/internal/application/repository/usc_coin_operations"
	valrepo "service-04/internal/application/repository/validator_operations"

	"service-04/api/api"
	"service-04/internal/infrastructure/auth"
	"service-04/internal/infrastructure/database"
	"service-04/internal/infrastructure/kafka"
	"service-04/internal/infrastructure/metrics"
	"service-04/internal/infrastructure/validation"

	// Cosmos SDK imports
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/app"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/storage"

	"github.com/usc-platform/shared/config"
	"github.com/usc-platform/shared/constants"
	"github.com/usc-platform/shared/logging"
)

// Container manages all application dependencies for USC Blockchain Core Service
type Container struct {
	// Configuration
	config *config.Config
	logger *logging.Logger

	// Infrastructure components
	db            *database.PostgreSQLManager
	pool          *database.PoolManager
	migrations    *database.MigrationManager
	redisManager  *database.RedisManager
	auth          *auth.JWTService
	validator     *validation.Validator
	metrics       *metrics.MetricsService
	kafkaProducer *kafka.KafkaProducerManager
	kafkaConsumer *kafka.KafkaConsumerManager

	// Cosmos SDK blockchain components
	cosmosApp         *app.USCApp
	rocksDBManager    *storage.RocksDBManager
	blockchainStorage *storage.StateManager

	// Block Operations domain
	BlockRepository *blockrepo.Repository
	BlockService    *blockbiz.Service
	BlockHandlers   *blockhandlers.Handlers

	// Transaction Operations domain
	TransactionRepository *txrepo.Repository
	TransactionService    *txbiz.Service
	TransactionHandlers   *txhandlers.Handlers

	// USC Coin Operations domain
	USCCoinRepository *uscrepo.Repository
	USCCoinService    *uscbiz.Service
	USCCoinHandlers   *uschandlers.Handlers

	// Smart Contract Operations domain
	ContractRepository *contractrepo.Repository
	ContractService    *contractbiz.Service
	ContractHandlers   *contracthandlers.Handlers

	// NFT Token Operations domain
	NFTRepository *nftrepo.Repository
	NFTService    *nftbiz.Service
	NFTHandlers   *nfthandlers.Handlers

	// Custom Token Operations domain
	TokenRepository *tokenrepo.Repository
	TokenService    *tokenbiz.Service
	TokenHandlers   *tokenhandlers.Handlers

	// Product Certificate Operations domain
	CertificateRepository *certrepo.Repository
	CertificateService    *certbiz.Service
	CertificateHandlers   *certhandlers.Handlers

	// Validator Operations domain
	ValidatorRepository *valrepo.Repository
	ValidatorService    *valbiz.Service
	ValidatorHandlers   *valhandlers.Handlers

	// Network Operations domain
	NetworkRepository *netrepo.Repository
	NetworkService    *netbiz.Service
	NetworkHandlers   *nethandlers.Handlers

	// Streaming Operations domain
	StreamRepository *streamrepo.Repository
	StreamService    *streambiz.Service
	StreamHandlers   *streamhandlers.Handlers

	// Store Bridge Operations domain
	BridgeRepository *bridgerepo.Repository
	BridgeService    *bridgebiz.Service
	BridgeHandlers   *bridgehandlers.Handlers

	// Store Network Operations domain
	StoreNetworkRepository *storenetrepo.Repository
	StoreNetworkService    *storenetbiz.Service
	StoreNetworkHandlers   *storenethandlers.Handlers

	// API endpoints
	grpcEndpoints *api.GRPCEndpoints
}

// NewContainer creates a new dependency injection container
func NewContainer(
	cfg *config.Config,
	logger *logging.Logger,
	db *database.PostgreSQLManager,
	pool *database.PoolManager,
	migrations *database.MigrationManager,
	redisManager *database.RedisManager,
	auth *auth.JWTService,
	validator *validation.Validator,
	metrics *metrics.MetricsService,
	kafkaProducer *kafka.KafkaProducerManager,
	kafkaConsumer *kafka.KafkaConsumerManager,
	cosmosApp *app.USCApp,
	rocksDBManager *storage.RocksDBManager,
	blockchainStorage *storage.StateManager,
) *Container {
	return &Container{
		config:            cfg,
		logger:            logger,
		db:                db,
		pool:              pool,
		migrations:        migrations,
		redisManager:      redisManager,
		auth:              auth,
		validator:         validator,
		metrics:           metrics,
		kafkaProducer:     kafkaProducer,
		kafkaConsumer:     kafkaConsumer,
		cosmosApp:         cosmosApp,
		rocksDBManager:    rocksDBManager,
		blockchainStorage: blockchainStorage,
	}
}

// Initialize initializes all dependencies in the correct order
func (c *Container) Initialize(ctx context.Context) error {
	c.logger.Info("Initializing USC Blockchain Core Service container",
		logging.String("service", constants.ServiceBlockchainCore))

	// Initialize repository layer
	if err := c.initializeRepositories(); err != nil {
		return err
	}

	// Initialize business layer
	if err := c.initializeBusiness(); err != nil {
		return err
	}

	// Initialize handlers layer
	if err := c.initializeHandlers(); err != nil {
		return err
	}

	// Initialize API endpoints
	if err := c.initializeAPIEndpoints(); err != nil {
		return err
	}

	c.logger.Info("USC Blockchain Core Service container initialized successfully",
		logging.String("service", constants.ServiceBlockchainCore))

	return nil
}

// initializeRepositories initializes all repository layers
func (c *Container) initializeRepositories() error {
	c.logger.Info("Initializing repository layer for USC Blockchain Core Service")

	// Block Operations
	c.BlockRepository = blockrepo.NewRepository(c.db, c.cosmosApp, c.blockchainStorage, c.redisManager, c.logger)

	// Transaction Operations
	c.TransactionRepository = txrepo.NewRepository(c.db, c.cosmosApp, c.blockchainStorage, c.redisManager, c.logger)

	// USC Coin Operations
	c.USCCoinRepository = uscrepo.NewRepository(c.db, c.cosmosApp, c.blockchainStorage, c.redisManager, c.logger)

	// Smart Contract Operations
	c.ContractRepository = contractrepo.NewRepository(c.db, c.cosmosApp, c.blockchainStorage, c.redisManager, c.logger)

	// NFT Token Operations
	c.NFTRepository = nftrepo.NewRepository(c.db, c.cosmosApp, c.blockchainStorage, c.redisManager, c.logger)

	// Custom Token Operations
	c.TokenRepository = tokenrepo.NewRepository(c.db, c.cosmosApp, c.blockchainStorage, c.redisManager, c.logger)

	// Product Certificate Operations
	c.CertificateRepository = certrepo.NewRepository(c.db, c.cosmosApp, c.blockchainStorage, c.redisManager, c.logger)

	// Validator Operations
	c.ValidatorRepository = valrepo.NewRepository(c.db, c.cosmosApp, c.blockchainStorage, c.redisManager, c.logger)

	// Network Operations
	c.NetworkRepository = netrepo.NewRepository(c.db, c.cosmosApp, c.blockchainStorage, c.redisManager, c.logger)

	// Streaming Operations
	c.StreamRepository = streamrepo.NewRepository(c.db, c.cosmosApp, c.blockchainStorage, c.redisManager, c.logger)

	// Store Bridge Operations
	c.BridgeRepository = bridgerepo.NewRepository(c.db, c.cosmosApp, c.blockchainStorage, c.redisManager, c.logger)

	// Store Network Operations
	c.StoreNetworkRepository = storenetrepo.NewRepository(c.db, c.cosmosApp, c.blockchainStorage, c.redisManager, c.logger)

	return nil
}

// initializeBusiness initializes all business layers
func (c *Container) initializeBusiness() error {
	c.logger.Info("Initializing business layer for USC Blockchain Core Service")

	// Block Operations
	c.BlockService = blockbiz.NewService(c.BlockRepository, c.cosmosApp, c.blockchainStorage, c.logger, c.validator, c.metrics)

	// Transaction Operations
	c.TransactionService = txbiz.NewService(c.TransactionRepository, c.cosmosApp, c.blockchainStorage, c.logger, c.validator, c.metrics)

	// USC Coin Operations
	c.USCCoinService = uscbiz.NewService(c.USCCoinRepository, c.cosmosApp, c.blockchainStorage, c.logger, c.validator, c.metrics)

	// Smart Contract Operations
	c.ContractService = contractbiz.NewService(c.ContractRepository, c.cosmosApp, c.blockchainStorage, c.logger, c.validator, c.metrics)

	// NFT Token Operations
	c.NFTService = nftbiz.NewService(c.NFTRepository, c.cosmosApp, c.blockchainStorage, c.logger, c.validator, c.metrics)

	// Custom Token Operations
	c.TokenService = tokenbiz.NewService(c.TokenRepository, c.cosmosApp, c.blockchainStorage, c.logger, c.validator, c.metrics)

	// Product Certificate Operations
	c.CertificateService = certbiz.NewService(c.CertificateRepository, c.cosmosApp, c.blockchainStorage, c.logger, c.validator, c.metrics)

	// Validator Operations
	c.ValidatorService = valbiz.NewService(c.ValidatorRepository, c.cosmosApp, c.blockchainStorage, c.logger, c.validator, c.metrics)

	// Network Operations
	c.NetworkService = netbiz.NewService(c.NetworkRepository, c.cosmosApp, c.blockchainStorage, c.logger, c.validator, c.metrics)

	// Streaming Operations
	c.StreamService = streambiz.NewService(c.StreamRepository, c.cosmosApp, c.blockchainStorage, c.logger, c.validator, c.metrics)

	// Store Bridge Operations
	c.BridgeService = bridgebiz.NewService(c.BridgeRepository, c.cosmosApp, c.blockchainStorage, c.logger, c.validator, c.metrics)

	// Store Network Operations
	c.StoreNetworkService = storenetbiz.NewService(c.StoreNetworkRepository, c.cosmosApp, c.blockchainStorage, c.logger, c.validator, c.metrics)

	return nil
}

// initializeHandlers initializes all handlers layers
func (c *Container) initializeHandlers() error {
	c.logger.Info("Initializing handlers layer for USC Blockchain Core Service")

	// Block Operations
	c.BlockHandlers = blockhandlers.NewHandlers(c.BlockService, c.logger)

	// Transaction Operations
	c.TransactionHandlers = txhandlers.NewHandlers(c.TransactionService, c.logger)

	// USC Coin Operations
	c.USCCoinHandlers = uschandlers.NewHandlers(c.USCCoinService, c.logger)

	// Smart Contract Operations
	c.ContractHandlers = contracthandlers.NewHandlers(c.ContractService, c.logger)

	// NFT Token Operations
	c.NFTHandlers = nfthandlers.NewHandlers(c.NFTService, c.logger)

	// Custom Token Operations
	c.TokenHandlers = tokenhandlers.NewHandlers(c.TokenService, c.logger)

	// Product Certificate Operations
	c.CertificateHandlers = certhandlers.NewHandlers(c.CertificateService, c.logger)

	// Validator Operations
	c.ValidatorHandlers = valhandlers.NewHandlers(c.ValidatorService, c.logger)

	// Network Operations
	c.NetworkHandlers = nethandlers.NewHandlers(c.NetworkService, c.logger)

	// Streaming Operations
	c.StreamHandlers = streamhandlers.NewHandlers(c.StreamService, c.logger)

	// Store Bridge Operations
	c.BridgeHandlers = bridgehandlers.NewHandlers(c.BridgeService, c.logger)

	// Store Network Operations
	c.StoreNetworkHandlers = storenethandlers.NewHandlers(c.StoreNetworkService, c.logger)

	return nil
}

// initializeAPIEndpoints initializes API endpoints
func (c *Container) initializeAPIEndpoints() error {
	c.logger.Info("Initializing API endpoints for USC Blockchain Core Service")

	// Initialize gRPC endpoints (for all services)
	c.grpcEndpoints = api.NewGRPCEndpoints(c.logger)

	return nil
}

// Getter methods for all handlers
func (c *Container) GetBlockHandlers() *blockhandlers.Handlers {
	return c.BlockHandlers
}

func (c *Container) GetTransactionHandlers() *txhandlers.Handlers {
	return c.TransactionHandlers
}

func (c *Container) GetUSCCoinHandlers() *uschandlers.Handlers {
	return c.USCCoinHandlers
}

func (c *Container) GetContractHandlers() *contracthandlers.Handlers {
	return c.ContractHandlers
}

func (c *Container) GetNFTHandlers() *nfthandlers.Handlers {
	return c.NFTHandlers
}

func (c *Container) GetTokenHandlers() *tokenhandlers.Handlers {
	return c.TokenHandlers
}

func (c *Container) GetCertificateHandlers() *certhandlers.Handlers {
	return c.CertificateHandlers
}

func (c *Container) GetValidatorHandlers() *valhandlers.Handlers {
	return c.ValidatorHandlers
}

func (c *Container) GetNetworkHandlers() *nethandlers.Handlers {
	return c.NetworkHandlers
}

func (c *Container) GetStreamHandlers() *streamhandlers.Handlers {
	return c.StreamHandlers
}

func (c *Container) GetBridgeHandlers() *bridgehandlers.Handlers {
	return c.BridgeHandlers
}

func (c *Container) GetStoreNetworkHandlers() *storenethandlers.Handlers {
	return c.StoreNetworkHandlers
}

// GetGRPCEndpoints returns the gRPC endpoints
func (c *Container) GetGRPCEndpoints() *api.GRPCEndpoints {
	return c.grpcEndpoints
}

// Infrastructure component getters
func (c *Container) GetDatabase() *database.PostgreSQLManager {
	return c.db
}

func (c *Container) GetPool() *database.PoolManager {
	return c.pool
}

func (c *Container) GetMigrations() *database.MigrationManager {
	return c.migrations
}

func (c *Container) GetAuth() *auth.JWTService {
	return c.auth
}

func (c *Container) GetValidator() *validation.Validator {
	return c.validator
}

func (c *Container) GetMetrics() *metrics.MetricsService {
	return c.metrics
}

func (c *Container) GetKafkaProducer() *kafka.KafkaProducerManager {
	return c.kafkaProducer
}

func (c *Container) GetKafkaConsumer() *kafka.KafkaConsumerManager {
	return c.kafkaConsumer
}

// Cosmos SDK component getters
func (c *Container) GetCosmosApp() *app.USCApp {
	return c.cosmosApp
}

func (c *Container) GetRocksDBManager() *storage.RocksDBManager {
	return c.rocksDBManager
}

func (c *Container) GetBlockchainStorage() *storage.StateManager {
	return c.blockchainStorage
}

// Shutdown gracefully shuts down all dependencies
func (c *Container) Shutdown(ctx context.Context) error {
	c.logger.Info("Shutting down USC Blockchain Core Service container",
		logging.String("service", constants.ServiceBlockchainCore))

	// Shutdown API endpoints
	if c.grpcEndpoints != nil {
		c.grpcEndpoints.Stop()
	}

	// Close database connections
	if c.db != nil {
		if err := c.db.Close(); err != nil {
			c.logger.Error("Error closing database connection", logging.Error(err))
		}
	}

	if c.pool != nil {
		if err := c.pool.Close(); err != nil {
			c.logger.Error("Error closing pool manager", logging.Error(err))
		}
	}

	// Close Kafka connections
	if c.kafkaProducer != nil {
		if err := c.kafkaProducer.Close(); err != nil {
			c.logger.Error("Error closing kafka producer", logging.Error(err))
		}
	}

	if c.kafkaConsumer != nil {
		if err := c.kafkaConsumer.Close(); err != nil {
			c.logger.Error("Error closing kafka consumer", logging.Error(err))
		}
	}

	// Close auth service
	if c.auth != nil {
		if err := c.auth.Close(); err != nil {
			c.logger.Error("Error closing auth service", logging.Error(err))
		}
	}

	// Close metrics service
	if c.metrics != nil {
		if err := c.metrics.Close(); err != nil {
			c.logger.Error("Error closing metrics service", logging.Error(err))
		}
	}

	// Close Cosmos SDK blockchain components
	if c.cosmosApp != nil {
		if err := c.cosmosApp.Stop(); err != nil {
			c.logger.Error("Error stopping Cosmos SDK app", logging.Error(err))
		}
	}

	if c.rocksDBManager != nil {
		if err := c.rocksDBManager.Close(); err != nil {
			c.logger.Error("Error closing RocksDB manager", logging.Error(err))
		}
	}

	// Note: StateManager doesn't need explicit closing
	if c.blockchainStorage != nil {
		c.logger.Info("Blockchain storage state manager available")
	}

	c.logger.Info("USC Blockchain Core Service container shutdown completed",
		logging.String("service", constants.ServiceBlockchainCore))

	return nil
}
