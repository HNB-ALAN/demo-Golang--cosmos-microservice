package application

import (
	"context"
	"os"
	"time"

	// Cosmos SDK imports
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/app"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/storage"

	"service-04/internal/infrastructure/database"

	"github.com/usc-platform/shared/config"
	"github.com/usc-platform/shared/constants"
	"github.com/usc-platform/shared/logging"
	"github.com/usc-platform/shared/utils"
)

// App represents the application instance
type App struct {
	config            *config.Config
	logger            *logging.Logger
	ctx               context.Context
	cancel            context.CancelFunc
	startTime         time.Time
	cosmosApp         *app.USCApp
	rocksDBManager    *storage.RocksDBManager
	blockchainStorage *storage.StateManager
	redisManager      *database.RedisManager
}

// NewApp creates a new application instance
func NewApp() (*App, error) {
	// Create context
	ctx, cancel := context.WithCancel(context.Background())

	// Load configuration
	cfg, err := config.LoadConfig("configs/config.yaml", constants.ServiceBlockchainCore)
	if err != nil {
		cancel()
		return nil, err
	}

	// Override secrets from environment variables (security best practice)
	overrideConfigFromEnv(cfg)

	// Create logger
	logger := logging.NewLogger(constants.ServiceBlockchainCore, cfg.Log)

	// Log application startup with timestamp
	timeUtils := utils.NewTimeUtils()
	startTime := timeUtils.FormatTimeISO(timeUtils.UTCNow())
	logger.Info("Application starting",
		logging.String("service", constants.ServiceBlockchainCore),
		logging.String("version", constants.AppVersion),
		logging.String("start_time", startTime))

	return &App{
		config:    cfg,
		logger:    logger,
		ctx:       ctx,
		cancel:    cancel,
		startTime: time.Now(),
	}, nil
}

// Run initializes and runs the application
func (a *App) Run() error {
	a.logger.Info("Application running",
		logging.String("service", constants.ServiceBlockchainCore))

	// Application initialization logic here
	// Database connections are handled in main.go
	// External services initialization can be added here

	return nil
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown() {
	a.logger.Info("Application shutting down",
		logging.String("service", constants.ServiceBlockchainCore))

	// Graceful shutdown logic here
	// Database connections are closed in main.go
	// Stop background tasks, cleanup resources, etc.

	// Cancel context
	a.cancel()
}

// GetConfig returns the application configuration
func (a *App) GetConfig() *config.Config {
	return a.config
}

// GetLogger returns the application logger
func (a *App) GetLogger() *logging.Logger {
	return a.logger
}

// GetContext returns the application context
func (a *App) GetContext() context.Context {
	return a.ctx
}

// IsRunning checks if the application is running
func (a *App) IsRunning() bool {
	select {
	case <-a.ctx.Done():
		return false
	default:
		return true
	}
}

// WaitForShutdown waits for the application to be shut down
func (a *App) WaitForShutdown() {
	<-a.ctx.Done()
}

// GetUptime returns the application uptime
func (a *App) GetUptime() time.Duration {
	return time.Since(a.startTime)
}

// GetCosmosApp returns the Cosmos SDK app instance
func (a *App) GetCosmosApp() *app.USCApp {
	return a.cosmosApp
}

// GetRocksDBManager returns the RocksDB manager instance
func (a *App) GetRocksDBManager() *storage.RocksDBManager {
	return a.rocksDBManager
}

// GetBlockchainStorage returns the blockchain storage manager instance
func (a *App) GetBlockchainStorage() *storage.StateManager {
	return a.blockchainStorage
}

// SetCosmosApp sets the Cosmos SDK app instance
func (a *App) SetCosmosApp(cosmosApp *app.USCApp) {
	a.cosmosApp = cosmosApp
}

// SetRocksDBManager sets the RocksDB manager instance
func (a *App) SetRocksDBManager(rocksDBManager *storage.RocksDBManager) {
	a.rocksDBManager = rocksDBManager
}

// SetBlockchainStorage sets the blockchain storage manager instance
func (a *App) SetBlockchainStorage(blockchainStorage *storage.StateManager) {
	a.blockchainStorage = blockchainStorage
}

// GetRedisManager returns the Redis manager instance
func (a *App) GetRedisManager() *database.RedisManager {
	return a.redisManager
}

// SetRedisManager sets the Redis manager instance
func (a *App) SetRedisManager(redisManager *database.RedisManager) {
	a.redisManager = redisManager
}

// overrideConfigFromEnv overrides config secrets from environment variables
// This ensures secrets are not hardcoded in config.yaml (security best practice)
func overrideConfigFromEnv(cfg *config.Config) {
	// JWT Secret
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		cfg.Auth.JWTSecret = jwtSecret
	}

	// Database passwords
	if dbPassword := os.Getenv("POSTGRES_PASSWORD"); dbPassword != "" {
		cfg.Database.Password = dbPassword
	}

	// Redis password
	if redisPassword := os.Getenv("REDIS_PASSWORD"); redisPassword != "" {
		cfg.Redis.Password = redisPassword
	}

	// ClickHouse password
	if clickhousePassword := os.Getenv("CLICKHOUSE_PASSWORD"); clickhousePassword != "" {
		cfg.ClickHouse.Password = clickhousePassword
	}

	// InfluxDB token
	if influxToken := os.Getenv("INFLUXDB_TOKEN"); influxToken != "" {
		cfg.InfluxDB.Token = influxToken
	}
}
