package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"service-04/api/grpc/client"
	"service-04/api/grpc/interceptors"
	grpcreflection "service-04/api/grpc/reflection"
	"service-04/api/grpc/server"
	"service-04/internal/application"
	"service-04/internal/infrastructure/auth"
	"service-04/internal/infrastructure/database"
	"service-04/internal/infrastructure/errors"
	"service-04/internal/infrastructure/kafka"
	"service-04/internal/infrastructure/metrics"
	"service-04/internal/infrastructure/validation"

	// Cosmos SDK imports
	cosmosapp "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/app"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/storage"

	"github.com/usc-platform/shared/constants"
	"github.com/usc-platform/shared/logging"
)

// main is the entry point of the application
func main() {
	// Create application instance
	app, err := application.NewApp()
	if err != nil {
		// Use shared error handling
		log.Fatalf("Failed to create application: %v", err)
	}

	// Initialize RocksDB for business logic (production - requires CGO)
	rocksDBManager, err := storage.NewRocksDBManager(storage.DefaultRocksDBConfig())
	if err != nil {
		app.GetLogger().Fatal("Failed to initialize RocksDB manager (CGO required)",
			logging.Error(err))
	}
	app.GetLogger().Info("RocksDB initialized successfully", logging.String("path", storage.DefaultRocksDBConfig().DataPath))

	blockchainStorage := storage.NewStateManager(rocksDBManager)

	// Initialize Cosmos SDK app using factory function from block-chain-cosmos
	// This keeps Cosmos SDK initialization logic in the block-chain-cosmos module
	cosmosDBDir := "./data/cosmos" // Default directory for Cosmos SDK database
	cosmosApp, cosmosDatabase, err := cosmosapp.NewUSCAppWithRocksDB(cosmosDBDir)
	if err != nil {
		app.GetLogger().Fatal("Failed to initialize Cosmos SDK app",
			logging.Error(err))
	}
	app.GetLogger().Info("Cosmos SDK app initialized successfully",
		logging.String("db_dir", cosmosDBDir))

	// cosmosDatabase will be closed in setupGracefulShutdown

	// Initialize database managers
	postgresManager, err := database.NewPostgreSQLManager(app.GetConfig(), *app.GetLogger(), cosmosApp, blockchainStorage)
	if err != nil {
		app.GetLogger().Fatal("Failed to create postgres manager",
			logging.Error(err))
	}

	// Initialize Redis manager
	redisManager, err := database.NewRedisManager(app.GetConfig(), *app.GetLogger())
	if err != nil {
		app.GetLogger().Fatal("Failed to create redis manager",
			logging.Error(err))
	}

	// Initialize pool manager
	poolManager, err := database.NewPoolManager(app.GetConfig(), *app.GetLogger())
	if err != nil {
		app.GetLogger().Fatal("Failed to create pool manager",
			logging.Error(err))
	}

	// Initialize migration manager
	migrationManager, err := database.NewMigrationManager(postgresManager.GetManager(), app.GetConfig(), *app.GetLogger())
	if err != nil {
		app.GetLogger().Fatal("Failed to create migration manager",
			logging.Error(err))
	}

	// Initialize Tier 3 services
	authService, err := auth.NewJWTService(app.GetConfig(), *app.GetLogger(), redisManager)
	if err != nil {
		app.GetLogger().Fatal("Failed to create auth service",
			logging.Error(err))
	}

	validator := validation.NewValidator(*app.GetLogger())

	metricsService, err := metrics.NewMetricsService(*app.GetLogger())
	if err != nil {
		app.GetLogger().Fatal("Failed to create metrics service",
			logging.Error(err))
	}

	// Initialize error manager
	errorManager := errors.NewErrorManager()

	// Initialize Tier 4 services
	kafkaProducerManager, err := kafka.NewKafkaProducerManager(app.GetConfig(), app.GetLogger())
	if err != nil {
		app.GetLogger().Fatal("Failed to create kafka producer manager",
			logging.Error(err))
	}

	kafkaConsumerManager, err := kafka.NewKafkaConsumerManager(app.GetConfig(), app.GetLogger())
	if err != nil {
		app.GetLogger().Fatal("Failed to create kafka consumer manager",
			logging.Error(err))
	}

	// Cosmos SDK components already initialized above

	app.GetLogger().Info("Cosmos SDK blockchain components initialized successfully")

	// Run database migrations - Hybrid approach for development and production
	ctx := context.Background()

	// Check if we're running in Docker (init-database.sh already ran)
	// vs Development mode (need to run migrations)
	if os.Getenv("DOCKER_CONTAINER") == "true" {
		app.GetLogger().Info("Running in Docker - migrations handled by init-database.sh")
	} else {
		app.GetLogger().Info("Running in Development mode - running migrations")
		if err := migrationManager.MigrateUp(ctx); err != nil {
			app.GetLogger().Fatal("Failed to run database migrations",
				logging.Error(err))
		}
		app.GetLogger().Info("Database migrations completed successfully")
	}

	// Initialize Phase 3 Container (Dependency Injection)
	container := application.NewContainer(
		app.GetConfig(),
		app.GetLogger(),
		postgresManager,
		poolManager,
		migrationManager,
		redisManager,
		authService,
		validator,
		metricsService,
		kafkaProducerManager,
		kafkaConsumerManager,
		cosmosApp,
		rocksDBManager,
		blockchainStorage,
	)
	if err := container.Initialize(app.GetContext()); err != nil {
		app.GetLogger().Fatal("Failed to initialize container",
			logging.Error(err))
	}

	// Initialize gRPC enhanced services
	grpcClientManager := client.NewServiceBlockchainCoreClientManager(app.GetConfig(), app.GetLogger(), cosmosApp, blockchainStorage, rocksDBManager)
	grpcInterceptors := interceptors.NewServiceBlockchainCoreInterceptors(app.GetLogger(), nil, cosmosApp, blockchainStorage, rocksDBManager)
	grpcReflection := grpcreflection.NewServiceBlockchainCoreReflectionService(app.GetLogger(), nil)

	// Create server instance with all services
	grpcServer := server.NewServer(app.GetConfig(), app.GetLogger(), postgresManager, poolManager, migrationManager, authService, validator, metricsService, errorManager, kafkaProducerManager, kafkaConsumerManager, grpcClientManager, grpcInterceptors, grpcReflection, container, cosmosApp, blockchainStorage, rocksDBManager, redisManager)

	// Registration is handled inside server.Start()

	// Reflection will be registered by the server when the gRPC server is created

	// Setup graceful shutdown
	setupGracefulShutdown(grpcServer, app, container, postgresManager, poolManager, migrationManager, authService, metricsService, errorManager, kafkaProducerManager, kafkaConsumerManager, grpcClientManager, grpcInterceptors, grpcReflection, cosmosApp, cosmosDatabase, blockchainStorage)

	// Initialize the application
	if err := app.Run(); err != nil {
		// Use shared error handling with structured logging
		app.GetLogger().Fatal("Failed to run application",
			logging.Error(err))
	}

	// Start background blockchain sync job BEFORE blocking server start
	go func() {
		ticker := time.NewTicker(30 * time.Second) // Sync every 30 seconds
		defer ticker.Stop()

		// Initial sync (wait a bit for server to be ready)
		time.Sleep(2 * time.Second)
		ctx := context.Background()
		if err := postgresManager.SyncWithBlockchain(ctx); err != nil {
			app.GetLogger().Error("Initial blockchain sync failed",
				logging.Error(err))
		} else {
			app.GetLogger().Info("Initial blockchain sync completed successfully")
		}

		// Periodic sync
		for range ticker.C {
			ctx := context.Background()
			if err := postgresManager.SyncWithBlockchain(ctx); err != nil {
				app.GetLogger().Error("Periodic blockchain sync failed",
					logging.Error(err))
			}
		}
	}()

	// Start the server (blocking call)
	if err := grpcServer.Start(); err != nil {
		// Use shared error handling
		app.GetLogger().Fatal("Failed to start server",
			logging.Error(err))
	}
}

// (Removed) register* helper functions; registration is centralized in server.Start()

// setupGracefulShutdown sets up graceful shutdown handling
func setupGracefulShutdown(grpcServer *server.Server, app *application.App, container *application.Container, postgresManager *database.PostgreSQLManager, poolManager *database.PoolManager, migrationManager *database.MigrationManager, authService *auth.JWTService, metricsService *metrics.MetricsService, errorManager *errors.ErrorManager, kafkaProducerManager *kafka.KafkaProducerManager, kafkaConsumerManager *kafka.KafkaConsumerManager, grpcClientManager *client.ServiceBlockchainCoreClientManager, grpcInterceptors *interceptors.ServiceBlockchainCoreInterceptors, grpcReflection *grpcreflection.ServiceBlockchainCoreReflectionService, cosmosApp *cosmosapp.USCApp, cosmosDatabase interface{ Close() error }, blockchainStorage *storage.StateManager) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		app.GetLogger().Info("Received shutdown signal, shutting down gracefully...",
			logging.String("service", constants.ServiceBlockchainCore))

		// Stop the server
		grpcServer.Stop()

		// Use container's shutdown method for proper cleanup
		if err := container.Shutdown(context.Background()); err != nil {
			app.GetLogger().Error("Failed to shutdown container",
				logging.Error(err))
		}

		// Close additional services not managed by container
		if err := poolManager.Close(); err != nil {
			app.GetLogger().Error("Failed to close pool manager",
				logging.Error(err))
		}

		// Close auth service
		if err := authService.Close(); err != nil {
			app.GetLogger().Error("Failed to close auth service",
				logging.Error(err))
		}

		// Close metrics service
		if err := metricsService.Close(); err != nil {
			app.GetLogger().Error("Failed to close metrics service",
				logging.Error(err))
		}

		// Close Tier 4 services
		if kafkaProducerManager != nil {
			if err := kafkaProducerManager.Close(); err != nil {
				app.GetLogger().Error("Failed to close kafka producer manager",
					logging.Error(err))
			}
		}
		if kafkaConsumerManager != nil {
			if err := kafkaConsumerManager.Close(); err != nil {
				app.GetLogger().Error("Failed to close kafka consumer manager",
					logging.Error(err))
			}
		}

		// Close gRPC Enhanced components
		if grpcClientManager != nil {
			if err := grpcClientManager.CloseAll(); err != nil {
				app.GetLogger().Error("Failed to close gRPC client manager",
					logging.Error(err))
			}
		}

		if grpcReflection != nil {
			if err := grpcReflection.Stop(context.Background()); err != nil {
				app.GetLogger().Error("Failed to stop gRPC reflection service",
					logging.Error(err))
			}
		}

		// Close Cosmos SDK blockchain components
		if cosmosApp != nil {
			if err := cosmosApp.Stop(); err != nil {
				app.GetLogger().Error("Failed to stop Cosmos SDK app",
					logging.Error(err))
			}
		}

		// Close Cosmos SDK database
		if cosmosDatabase != nil {
			if err := cosmosDatabase.Close(); err != nil {
				app.GetLogger().Error("Failed to close Cosmos SDK database",
					logging.Error(err))
			}
		}

		// Shutdown Phase 3 Container
		if container != nil {
			if err := container.Shutdown(context.Background()); err != nil {
				app.GetLogger().Error("Failed to shutdown container",
					logging.Error(err))
			}
		}

		// Shutdown the application
		app.Shutdown()

		os.Exit(0)
	}()
}
