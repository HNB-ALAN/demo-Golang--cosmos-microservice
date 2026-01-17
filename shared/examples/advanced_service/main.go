package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"google.golang.org/grpc/reflection"

	"github.com/usc-platform/shared/auth"
	"github.com/usc-platform/shared/cache"
	"github.com/usc-platform/shared/config"
	"github.com/usc-platform/shared/database"
	grpcShared "github.com/usc-platform/shared/grpc"
	"github.com/usc-platform/shared/health"
	"github.com/usc-platform/shared/logging"
	"github.com/usc-platform/shared/metrics"
	"github.com/usc-platform/shared/validation"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("", "advanced-service")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger := logging.NewLogger("advanced-service", cfg.Log)
	logger.Info("Starting advanced service")

	// Initialize database manager
	dbManager, err := database.NewDatabaseManager(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize database manager", logging.Error(err))
	}
	defer dbManager.Close()

	// Initialize metrics collector
	metricsCollector := metrics.NewMetricsCollector()
	defer metricsCollector.ClearMetrics()

	// Initialize cache
	cacheManager := cache.NewMemoryCache(cache.MemoryConfig{
		MaxSize: 1000,
		TTL:     time.Hour,
	})

	// Initialize auth service
	authService, err := auth.NewJWTService(auth.Config{
		SecretKey:     getEnvOrDefault("JWT_SECRET_KEY", generateSecureSecret()),
		Issuer:        getEnvOrDefault("JWT_ISSUER", "advanced-service"),
		AccessExpiry:  time.Hour,
		RefreshExpiry: time.Hour * 24 * 7,
	})
	if err != nil {
		logger.Fatal("Failed to initialize auth service", logging.Error(err))
	}

	// Initialize validator
	validator := validation.NewValidator()

	// Create gRPC server with advanced configuration
	grpcServer := grpcShared.NewServer(cfg, logger)

	// Create health service
	healthService := health.NewService("advanced-service", "1.0.0")

	// Register comprehensive health checks
	healthService.RegisterCheck("postgresql", database.NewPostgreSQLHealthChecker(dbManager.PostgreSQL()))
	healthService.RegisterCheck("redis", database.NewRedisHealthChecker(dbManager.Redis()))
	// MongoDB support removed
	healthService.RegisterCheck("clickhouse", database.NewClickHouseHealthChecker(dbManager.ClickHouse()))
	healthService.RegisterCheck("influxdb", database.NewInfluxDBHealthChecker(dbManager.InfluxDB()))
	healthService.RegisterCheck("quickwit", database.NewQuickwitHealthChecker(dbManager.Quickwit()))

	// Register custom health checks
	healthService.RegisterCheck("cache", &CacheHealthChecker{cache: cacheManager})
	healthService.RegisterCheck("auth", &AuthHealthChecker{auth: authService})
	healthService.RegisterCheck("validator", &ValidatorHealthChecker{validator: validator})

	// Register health service with gRPC server
	grpcServer.RegisterHealthService("advanced-service", "1.0.0")

	// Register reflection
	reflection.Register(grpcServer.GetServer())

	// Start server
	logger.Info("Starting advanced gRPC server on :9090")
	if err := grpcServer.Start(); err != nil {
		logger.Fatal("Failed to start gRPC server", logging.Error(err))
	}

	// Start background services
	startBackgroundServices(cfg, logger, dbManager, metricsCollector, cacheManager)

	// Wait for interrupt signal
	waitForShutdown(logger, grpcServer, dbManager, cacheManager)
}

// startBackgroundServices starts background services
func startBackgroundServices(cfg *config.Config, logger *logging.Logger, dbManager *database.DatabaseManager, metricsCollector *metrics.MetricsCollector, cacheManager cache.Cache) {
	// Start metrics collection
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			collectSystemMetrics(metricsCollector)
		}
	}()

	// Start cache warming
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			warmCache(cacheManager, logger)
		}
	}()

	// Start database maintenance
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			performDatabaseMaintenance(dbManager, logger)
		}
	}()

	logger.Info("Background services started")
}

// collectSystemMetrics collects system metrics
func collectSystemMetrics(collector *metrics.MetricsCollector) {
	// Collect memory usage
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	collector.RecordMetric("memory_alloc_bytes", m.Alloc)
	collector.RecordMetric("memory_total_alloc_bytes", m.TotalAlloc)
	collector.RecordMetric("memory_sys_bytes", m.Sys)
	collector.RecordMetric("goroutines_count", runtime.NumGoroutine())
}

// warmCache warms the cache with frequently accessed data
func warmCache(cacheManager cache.Cache, logger *logging.Logger) {
	// Implement cache warming logic
	logger.Debug("Warming cache")
}

// performDatabaseMaintenance performs database maintenance tasks
func performDatabaseMaintenance(dbManager *database.DatabaseManager, logger *logging.Logger) {
	// Implement database maintenance logic
	logger.Debug("Performing database maintenance")
}

// waitForShutdown waits for shutdown signal
func waitForShutdown(logger *logging.Logger, grpcServer *grpcShared.Server, dbManager *database.DatabaseManager, cacheManager cache.Cache) {
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

	// Close cache connections
	// cacheManager.Close() // MemoryCache doesn't have Close method

	logger.Info("Advanced service shutdown completed")
}

// Custom health checkers
type CacheHealthChecker struct {
	cache cache.Cache
}

func (c *CacheHealthChecker) Check(ctx context.Context) error {
	return c.cache.Health(ctx)
}

func (c *CacheHealthChecker) Name() string {
	return "cache"
}

func (c *CacheHealthChecker) Description() string {
	return "Cache health check"
}

type AuthHealthChecker struct {
	auth *auth.JWTService
}

func (a *AuthHealthChecker) Check(ctx context.Context) error {
	// Simple health check - just verify the service is initialized
	if a.auth == nil {
		return fmt.Errorf("auth service not initialized")
	}
	return nil
}

func (a *AuthHealthChecker) Name() string {
	return "auth"
}

func (a *AuthHealthChecker) Description() string {
	return "Authentication service health check"
}

type ValidatorHealthChecker struct {
	validator *validation.Validator
}

func (v *ValidatorHealthChecker) Check(ctx context.Context) error {
	// Simple health check - just verify the validator is initialized
	if v.validator == nil {
		return fmt.Errorf("validator not initialized")
	}
	return nil
}

func (v *ValidatorHealthChecker) Name() string {
	return "validator"
}

func (v *ValidatorHealthChecker) Description() string {
	return "Validation service health check"
}

// getEnvOrDefault returns environment variable value or default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// generateSecureSecret generates a cryptographically secure secret
func generateSecureSecret() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based secret if crypto/rand fails
		return fmt.Sprintf("fallback-secret-%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}
