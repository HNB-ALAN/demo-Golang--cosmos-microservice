package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc/reflection"

	"github.com/usc-platform/shared/config"
	"github.com/usc-platform/shared/database"
	grpcShared "github.com/usc-platform/shared/grpc"
	"github.com/usc-platform/shared/health"
	"github.com/usc-platform/shared/logging"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("", "basic-service")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger := logging.NewLogger("basic-service", cfg.Log)
	logger.Info("Starting basic service")

	// Initialize database manager
	dbManager, err := database.NewDatabaseManager(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize database manager", logging.Error(err))
	}
	defer dbManager.Close()

	// Create gRPC server
	grpcServer := grpcShared.NewServer(cfg, logger)

	// Create health service
	healthService := health.NewService("basic-service", "1.0.0")

	// Register database health checks
	healthService.RegisterCheck("postgresql", database.NewPostgreSQLHealthChecker(dbManager.PostgreSQL()))
	healthService.RegisterCheck("redis", database.NewRedisHealthChecker(dbManager.Redis()))

	// Register health service with gRPC server
	grpcServer.RegisterHealthService("basic-service", "1.0.0")

	// Register reflection
	reflection.Register(grpcServer.GetServer())

	// Start server
	logger.Info("Starting gRPC server on :9090")
	if err := grpcServer.Start(); err != nil {
		logger.Fatal("Failed to start gRPC server", logging.Error(err))
	}

	logger.Info("Basic service started successfully")

	// Wait for interrupt signal
	waitForShutdown(logger, grpcServer, dbManager)
}

// waitForShutdown waits for shutdown signal
func waitForShutdown(logger *logging.Logger, grpcServer *grpcShared.Server, dbManager *database.DatabaseManager) {
	// Create a channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for signal
	sig := <-sigChan
	logger.Info("Received shutdown signal", logging.String("signal", sig.String()))

	// Graceful shutdown
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Stop gRPC server
	grpcServer.Stop()

	// Close database connections
	if err := dbManager.Close(); err != nil {
		logger.Error("Failed to close database connections", logging.Error(err))
	}

	logger.Info("Basic service shutdown completed")
}
