# Migration Guide: From Individual Services to Shared Library

This guide helps you migrate your existing microservices to use the USC shared library, reducing code duplication and improving consistency across your 21 services.

## Overview

The USC shared library provides common functionality that was previously duplicated across all 21 services:

- Configuration management
- Database connections
- gRPC server setup
- Health checking
- Logging
- Metrics collection
- Authentication
- Caching
- Validation
- Error handling

## Migration Steps

### Step 1: Update go.mod

**Before:**
```go
module service-01-gateway

go 1.23.4

require (
    github.com/spf13/viper v1.21.0
    go.uber.org/zap v1.27.0
    google.golang.org/grpc v1.75.1
    github.com/go-redis/redis/v8 v8.11.5
    go.mongodb.org/mongo-driver v1.17.4
    github.com/lib/pq v1.10.9
    // ... many more dependencies
)
```

**After:**
```go
module service-01-gateway

go 1.23.4

require (
    github.com/usc-platform/shared v1.0.0
)

replace github.com/usc-platform/shared => ../shared
```

### Step 2: Update main.go

**Before:**
```go
package main

import (
    "log"
    "net"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    "google.golang.org/grpc/health/grpc_health_v1"
    "service-01-gateway/config"
    "service-01-gateway/proto"
)

func main() {
    cfg, err := config.LoadConfig("")
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.Port))
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    s := grpc.NewServer()
    
    // Register health service
    grpc_health_v1.RegisterHealthServer(s, &health.Server{})
    
    // Register reflection
    reflection.Register(s)
    
    // Register your service
    proto.RegisterGatewayServiceServer(s, &gatewayServer{})
    
    log.Printf("Server listening at %v", lis.Addr())
    if err := s.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
```

**After:**
```go
package main

import (
    "log"
    "github.com/usc-platform/shared/config"
    "github.com/usc-platform/shared/database"
    "github.com/usc-platform/shared/grpc"
    "github.com/usc-platform/shared/logging"
    "github.com/usc-platform/shared/health"
    "service-01-gateway/proto"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig("", "service-01-gateway")
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Initialize logger
    logger := logging.NewLogger("service-01-gateway", cfg.Log)
    logger.Info("Starting service-01-gateway")

    // Initialize database manager
    dbManager, err := database.NewManager(cfg)
    if err != nil {
        logger.Fatal("Failed to initialize database manager", logging.Error(err))
    }
    defer dbManager.Close()

    // Create gRPC server
    grpcServer := grpc.NewServer(cfg, logger)

    // Register health service
    healthService := grpc.RegisterHealthService(grpcServer, "service-01-gateway", "1.0.0")
    
    // Register database health checks
    healthService.RegisterCheck("postgresql", &database.PostgreSQLHealthChecker{db: dbManager.PostgreSQL()})
    healthService.RegisterCheck("redis", &database.RedisHealthChecker{client: dbManager.Redis()})

    // Register reflection
    grpc.RegisterReflection(grpcServer)

    // Register your service
    grpcServer.RegisterService(&proto.GatewayService_ServiceDesc, &gatewayServer{})

    // Start server
    if err := grpcServer.Start(); err != nil {
        logger.Fatal("Failed to start gRPC server", logging.Error(err))
    }

    logger.Info("Service-01-gateway started successfully")
}
```

### Step 3: Update Configuration

**Before:**
```go
// config/config.go
type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
    Log      LogConfig      `mapstructure:"log"`
}

func LoadConfig(configPath string) (*Config, error) {
    // ... custom implementation
}
```

**After:**
```go
// Remove config/config.go entirely
// Use shared config instead
import "github.com/usc-platform/shared/config"

cfg, err := config.LoadConfig("", "service-name")
```

### Step 4: Update Database Connections

**Before:**
```go
// database/manager.go
type Manager struct {
    postgres *sql.DB
    redis    *redis.Client
    mongodb  *mongo.Client
}

func NewManager(cfg *Config) (*Manager, error) {
    // ... custom implementation
}
```

**After:**
```go
// Remove database/manager.go entirely
// Use shared database manager instead
import "github.com/usc-platform/shared/database"

dbManager, err := database.NewManager(cfg)
postgres := dbManager.PostgreSQL()
redis := dbManager.Redis()
mongodb := dbManager.MongoDB()
```

### Step 5: Update Logging

**Before:**
```go
// logging/logger.go
type Logger struct {
    zap *zap.Logger
}

func NewLogger(cfg LogConfig) *Logger {
    // ... custom implementation
}
```

**After:**
```go
// Remove logging/logger.go entirely
// Use shared logger instead
import "github.com/usc-platform/shared/logging"

logger := logging.NewLogger("service-name", cfg.Log)
logger.Info("Service started", logging.String("version", "1.0.0"))
```

### Step 6: Update Health Checks

**Before:**
```go
// health/service.go
type Service struct {
    name    string
    version string
    checks  map[string]HealthChecker
}

func NewService(name, version string) *Service {
    // ... custom implementation
}
```

**After:**
```go
// Remove health/service.go entirely
// Use shared health service instead
import "github.com/usc-platform/shared/health"

healthService := health.NewService("service-name", "1.0.0")
healthService.RegisterCheck("database", dbManager.HealthCheck)
```

### Step 7: Update Metrics

**Before:**
```go
// metrics/collector.go
type Collector struct {
    requestsTotal prometheus.Counter
    requestDuration prometheus.Histogram
}

func NewCollector() *Collector {
    // ... custom implementation
}
```

**After:**
```go
// Remove metrics/collector.go entirely
// Use shared metrics collector instead
import "github.com/usc-platform/shared/metrics"

collector := metrics.NewMetricsCollector("service_name", "api")
collector.RecordHTTPRequest("GET", "/users", "200", duration, requestSize, responseSize)
```

### Step 8: Update Authentication

**Before:**
```go
// auth/service.go
type Service struct {
    jwtSecret string
    jwtExpiry time.Duration
}

func NewService(cfg AuthConfig) *Service {
    // ... custom implementation
}
```

**After:**
```go
// Remove auth/service.go entirely
// Use shared auth service instead
import "github.com/usc-platform/shared/auth"

authService := auth.NewService(cfg.Auth)
token, err := authService.GenerateToken(userID, claims)
```

### Step 9: Update Caching

**Before:**
```go
// cache/manager.go
type Manager struct {
    redis *redis.Client
    memory map[string]interface{}
}

func NewManager(cfg CacheConfig) *Manager {
    // ... custom implementation
}
```

**After:**
```go
// Remove cache/manager.go entirely
// Use shared cache manager instead
import "github.com/usc-platform/shared/cache"

cacheManager := cache.NewManager(cfg)
cacheManager.Set("key", "value", time.Hour)
```

### Step 10: Update Validation

**Before:**
```go
// validation/validator.go
type Validator struct {
    rules map[string]ValidationRule
}

func NewValidator() *Validator {
    // ... custom implementation
}
```

**After:**
```go
// Remove validation/validator.go entirely
// Use shared validator instead
import "github.com/usc-platform/shared/validation"

validator := validation.NewValidator()
err := validator.ValidateStruct(user)
```

### Step 11: Update Error Handling

**Before:**
```go
// errors/handler.go
type Handler struct {
    logger *zap.Logger
}

func NewHandler(logger *zap.Logger) *Handler {
    // ... custom implementation
}
```

**After:**
```go
// Remove errors/handler.go entirely
// Use shared error handling instead
import "github.com/usc-platform/shared/errors"

err := errors.NewValidationError("invalid input", "field", "value")
```

### Step 12: Update Constants

**Before:**
```go
// constants/app.go
const (
    DefaultPort = 8080
    DefaultTimeout = 30
    MaxRetries = 3
)
```

**After:**
```go
// Remove constants/app.go entirely
// Use shared constants instead
import "github.com/usc-platform/shared/constants"

port := constants.DefaultHTTPPort
timeout := constants.DefaultTimeout
```

### Step 13: Update Utils

**Before:**
```go
// utils/string.go
func IsEmpty(s string) bool {
    return len(strings.TrimSpace(s)) == 0
}
```

**After:**
```go
// Remove utils/string.go entirely
// Use shared utils instead
import "github.com/usc-platform/shared/utils"

isEmpty := utils.IsEmpty(s)
```

### Step 14: Update Middleware

**Before:**
```go
// middleware/auth.go
func AuthMiddleware() grpc.UnaryServerInterceptor {
    // ... custom implementation
}
```

**After:**
```go
// Remove middleware/auth.go entirely
// Use shared middleware instead
import "github.com/usc-platform/shared/grpc"

interceptor := grpc.AuthInterceptor(logger)
```

### Step 15: Update Testing

**Before:**
```go
// testing/mocks.go
type MockDatabase struct {
    // ... custom implementation
}
```

**After:**
```go
// Remove testing/mocks.go entirely
// Use shared testing utilities instead
import "github.com/usc-platform/shared/testing"

mockDB := testing.NewMockDatabase()
```

## Configuration Migration

### Before (Individual Service Config)

```yaml
# config.yaml
server:
  host: "0.0.0.0"
  port: "8080"

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "password"
  dbname: "service_db"

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

log:
  level: "info"
  format: "json"
  output: "stdout"
```

### After (Shared Library Config)

```yaml
# config.yaml
service:
  name: "service-01-gateway"
  version: "1.0.0"
  environment: "development"
  region: "us-east-1"
  instance: "default"

server:
  host: "0.0.0.0"
  port: "8080"
  grpc:
    max_recv_msg_size: 4194304
    max_send_msg_size: 4194304
    keep_alive: true
  http:
    read_timeout: 30
    write_timeout: 30
    idle_timeout: 120

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "password"
  dbname: "usc_social_media"
  sslmode: "disable"
  max_conns: 25
  min_conns: 5

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
  pool_size: 10

mongodb:
  host: "localhost"
  port: 27017
  user: ""
  password: ""
  dbname: "usc_social_media"
  auth_db: "admin"

clickhouse:
  host: "localhost"
  port: 9000
  user: "default"
  password: ""
  dbname: "usc_analytics"
  secure: false

influxdb:
  host: "localhost"
  port: 8086
  token: ""
  org: "usc"
  bucket: "metrics"

quickwit:
  host: "localhost"
  port: 7280
  user: ""
  password: ""
  index: "usc_index"

log:
  level: "info"
  format: "json"
  output: "stdout"
  filename: ""
  max_size: 100
  max_age: 30
  compress: true

metrics:
  enabled: true
  port: 9090
  path: "/metrics"

auth:
  jwt_secret: "your-super-secret-jwt-key-here"
  jwt_expiry: "24h"
  refresh_expiry: "168h"
  issuer: "usc-platform"

cache:
  ttl: 3600
  max_size: 1000
  enabled: true
```

## Benefits of Migration

### Code Reduction

- **Before**: ~2000 lines of duplicated code per service
- **After**: ~100 lines of service-specific code
- **Reduction**: 95% code reduction

### Consistency

- **Before**: 21 different implementations
- **After**: 1 shared implementation
- **Result**: 100% consistency across services

### Maintenance

- **Before**: Update 21 services for bug fixes
- **After**: Update 1 shared library
- **Result**: 95% reduction in maintenance effort

### Testing

- **Before**: Test 21 different implementations
- **After**: Test 1 shared implementation
- **Result**: 95% reduction in testing effort

### Performance

- **Before**: 21 different performance characteristics
- **After**: 1 optimized implementation
- **Result**: Consistent performance across services

### Security

- **Before**: 21 different security implementations
- **After**: 1 secure implementation
- **Result**: Consistent security across services

## Migration Checklist

### Pre-Migration

- [ ] Backup existing code
- [ ] Document current functionality
- [ ] Identify service-specific code
- [ ] Plan migration timeline
- [ ] Set up testing environment

### During Migration

- [ ] Update go.mod
- [ ] Update main.go
- [ ] Remove duplicate packages
- [ ] Update configuration
- [ ] Update imports
- [ ] Test functionality
- [ ] Update documentation

### Post-Migration

- [ ] Verify all functionality works
- [ ] Run integration tests
- [ ] Performance testing
- [ ] Security testing
- [ ] Update deployment scripts
- [ ] Update monitoring
- [ ] Update documentation

## Common Issues and Solutions

### Issue 1: Import Conflicts

**Problem**: Import conflicts between old and new packages

**Solution**: Remove old packages completely before adding shared library

### Issue 2: Configuration Differences

**Problem**: Configuration structure differences

**Solution**: Update configuration files to match shared library structure

### Issue 3: API Differences

**Problem**: API differences between old and new implementations

**Solution**: Update service code to use shared library APIs

### Issue 4: Database Schema Differences

**Problem**: Database schema differences

**Solution**: Update database schemas to match shared library expectations

### Issue 5: Logging Format Differences

**Problem**: Logging format differences

**Solution**: Update log parsing to handle new format

## Testing Migration

### Unit Tests

```bash
# Before migration
go test ./...

# After migration
go test ./...
```

### Integration Tests

```bash
# Before migration
go test -tags=integration ./...

# After migration
go test -tags=integration ./...
```

### Load Tests

```bash
# Before migration
go test -tags=load ./...

# After migration
go test -tags=load ./...
```

### End-to-End Tests

```bash
# Before migration
go test -tags=e2e ./...

# After migration
go test -tags=e2e ./...
```

## Rollback Plan

If migration fails, you can rollback by:

1. **Restore backup code**
2. **Revert go.mod changes**
3. **Restore original configuration**
4. **Restart services**
5. **Verify functionality**

## Support

For migration support:

- **Documentation**: [https://docs.usc-platform.com](https://docs.usc-platform.com)
- **Issues**: [https://github.com/usc-platform/shared/issues](https://github.com/usc-platform/shared/issues)
- **Discussions**: [https://github.com/usc-platform/shared/discussions](https://github.com/usc-platform/shared/discussions)
- **Email**: support@usc-platform.com

## Conclusion

Migrating to the USC shared library provides significant benefits:

- **95% code reduction**
- **100% consistency**
- **95% maintenance reduction**
- **95% testing reduction**
- **Consistent performance**
- **Consistent security**

The migration process is straightforward and well-documented, with comprehensive examples and support available.
