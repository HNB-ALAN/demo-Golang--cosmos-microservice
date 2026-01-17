package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"

	"github.com/usc-platform/shared/config"
	"github.com/usc-platform/shared/graphql"
	"github.com/usc-platform/shared/logging"
	"github.com/usc-platform/shared/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("", "graphql-service")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger := logging.NewLogger("graphql-service", cfg.Log)
	logger.Info("Starting GraphQL service")

	// Initialize GraphQL federation service
	federationConfig := &graphql.FederationConfig{
		GatewayURL:    getEnvOrDefault("GATEWAY_URL", "http://localhost:4000"),
		ServiceURL:    getEnvOrDefault("SERVICE_URL", "http://localhost:4001"),
		ServiceName:   "graphql-service",
		SchemaPath:    "./schema.graphql",
		Introspection: true,
		Playground:    true,
		Extensions: map[string]string{
			"version": "1.0.0",
		},
	}

	federationService := graphql.NewFederationService(logger, federationConfig)

	// Initialize GraphQL middleware
	// middlewareConfig := &graphql.MiddlewareConfig{
	// 	MaxQueryDepth:      10,
	// 	MaxQueryComplexity: 1000,
	// 	QueryTimeout:       30 * time.Second,
	// 	EnableTracing:      true,
	// 	EnableMetrics:      true,
	// 	RateLimitPerMinute: 1000,
	// }

	// graphqlMiddleware := graphql.NewGraphQLMiddleware(logger, middlewareConfig)

	// Create GraphQL handler
	// TODO: Implement actual GraphQL schema
	// graphqlHandler := handler.NewDefaultServer(createExecutableSchema())

	// For now, create a simple handler
	graphqlHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "GraphQL service is running"}`))
	})

	// Create Gin router
	router := gin.Default()

	// Add CORS middleware
	corsConfig := middleware.DefaultCORSConfig()
	corsMiddleware := middleware.NewCORSMiddleware(corsConfig)
	router.Use(gin.WrapH(corsMiddleware.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This will be handled by the CORS middleware
	}))))

	// Add rate limiting middleware
	rateLimitConfig := middleware.RateLimitConfig{
		RequestsPerSecond: 16.67, // 1000 requests per minute = 16.67 per second
		Burst:             100,
	}
	rateLimitMiddleware := middleware.NewHTTPRateLimitMiddleware(rateLimitConfig)
	router.Use(gin.WrapH(rateLimitMiddleware.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This will be handled by the rate limit middleware
	}))))

	// GraphQL endpoint
	router.POST("/graphql", gin.WrapH(graphqlHandler))
	router.GET("/graphql", gin.WrapH(playground.Handler("GraphQL Playground", "/graphql")))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		if err := federationService.HealthCheck(context.Background()); err != nil {
			c.JSON(500, gin.H{"status": "unhealthy", "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Federation info endpoint
	router.GET("/federation", func(c *gin.Context) {
		info, err := federationService.GetFederationInfo(context.Background())
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, info)
	})

	// Register service with federation gateway
	serviceInfo := &graphql.ServiceInfo{
		Name:    "graphql-service",
		Version: "1.0.0",
		URL:     getEnvOrDefault("SERVICE_URL", "http://localhost:4001"),
		Health: &graphql.ServiceHealth{
			Status:    "healthy",
			Timestamp: time.Now(),
			Latency:   10,
			Errors:    0,
		},
	}

	if err := federationService.RegisterService(context.Background(), serviceInfo); err != nil {
		logger.Error("Failed to register service with federation gateway", logging.Error(err))
	}

	// Start server
	port := getEnvOrDefault("PORT", "4001")
	logger.Info("Starting GraphQL server", logging.String("port", port))

	if err := router.Run(":" + port); err != nil {
		logger.Fatal("Failed to start GraphQL server", logging.Error(err))
	}

	logger.Info("GraphQL service started successfully")
}

// createExecutableSchema creates the GraphQL executable schema
// func createExecutableSchema() interface{} {
// 	// TODO: Implement actual GraphQL schema
// 	// This would typically use gqlgen to generate the schema
// 	return nil
// }

// getEnvOrDefault returns environment variable value or default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
