# 🏗️ USC SHARED LIBRARY - COMPREHENSIVE VERSION

## 📋 Tổng quan

Shared library này được thiết kế dựa trên phân tích chi tiết 21 microservices trong USC platform. Tất cả components đều được extract từ patterns thực tế của các services.

## 🎯 Mục tiêu

- ✅ **100% code reuse** - Tất cả code chung được extract
- ✅ **Consistency** - Đảm bảo tất cả services hoạt động giống nhau
- ✅ **Maintainability** - Dễ dàng update và maintain
- ✅ **Performance** - Tối ưu hóa cho production
- ✅ **Clean Architecture** - Tuân thủ SOLID principles

## 📁 Cấu trúc thư mục đầy đủ

```
shared/
├── README.md                           # Tài liệu này
├── go.mod                              # Go module definition
├── go.sum                              # Go dependencies checksums
├── config/                             # Configuration management
│   ├── config.go                       # Main config struct và loading (200 lines)
│   ├── defaults.go                     # Default values cho tất cả services (150 lines)
│   ├── validation.go                   # Config validation logic (100 lines)
│   └── env.go                          # Environment variables handling (80 lines)
├── database/                           # Database connection management
│   ├── manager.go                      # Database manager chính (250 lines)
│   ├── postgresql.go                   # PostgreSQL connection (120 lines)
│   ├── redis.go                        # Redis connection (100 lines)
│   ├── clickhouse.go                   # ClickHouse connection (100 lines)
│   ├── influxdb.go                     # InfluxDB connection (100 lines)
│   ├── quickwit.go                     # Quickwit connection (100 lines)
│   └── health.go                       # Database health checks (80 lines)
├── grpc/                               # gRPC utilities
│   ├── server.go                       # gRPC server setup (200 lines)
│   ├── client.go                       # gRPC client factory (150 lines)
│   ├── interceptors.go                 # Common interceptors (180 lines)
│   ├── health.go                       # Health service implementation (120 lines)
│   ├── reflection.go                   # Reflection setup (60 lines)
│   └── middleware.go                   # gRPC middleware (100 lines)
├── health/                             # Health checking system
│   ├── service.go                      # Health service implementation (150 lines)
│   ├── checker.go                      # Health checkers (120 lines)
│   ├── status.go                       # Health status definitions (80 lines)
│   └── registry.go                     # Health check registry (100 lines)
├── proto/                              # Common protobuf definitions
│   ├── health.proto                    # Health service proto (100 lines)
│   ├── common.proto                    # Common messages (150 lines)
│   ├── error.proto                     # Error definitions (80 lines)
│   └── timestamp.proto                 # Timestamp utilities (60 lines)
├── logging/                            # Logging framework
│   ├── logger.go                       # Logger implementation (200 lines)
│   ├── structured.go                   # Structured logging (150 lines)
│   ├── fields.go                       # Log fields definitions (100 lines)
│   ├── middleware.go                   # Logging middleware (120 lines)
│   └── formatters.go                   # Log formatters (100 lines)
├── metrics/                            # Metrics collection
│   ├── prometheus.go                   # Prometheus metrics (150 lines)
│   ├── custom.go                       # Custom metrics (120 lines)
│   ├── middleware.go                   # Metrics middleware (100 lines)
│   └── collectors.go                   # Metric collectors (120 lines)
├── auth/                               # Authentication & Authorization
│   ├── jwt.go                          # JWT utilities (150 lines)
│   ├── middleware.go                   # Auth middleware (180 lines)
│   ├── permissions.go                  # Permission system (120 lines)
│   ├── tokens.go                       # Token management (100 lines)
│   └── validation.go                   # Auth validation (80 lines)
├── cache/                              # Caching layer
│   ├── redis.go                        # Redis cache (150 lines)
│   ├── memory.go                       # Memory cache (120 lines)
│   ├── patterns.go                     # Cache patterns (100 lines)
│   └── middleware.go                   # Cache middleware (80 lines)
├── validation/                         # Input validation
│   ├── validator.go                    # Validation interface (80 lines)
│   ├── struct.go                       # Struct validation (120 lines)
│   ├── business.go                     # Business validation (150 lines)
│   └── errors.go                       # Validation errors (60 lines)
├── errors/                             # Error handling
│   ├── domain.go                       # Domain errors (100 lines)
│   ├── grpc.go                         # gRPC errors (120 lines)
│   ├── http.go                         # HTTP errors (80 lines)
│   ├── codes.go                        # Error codes (60 lines)
│   └── handlers.go                     # Error handlers (100 lines)
├── utils/                              # Utility functions
│   ├── time.go                         # Time utilities (80 lines)
│   ├── string.go                       # String utilities (120 lines)
│   ├── json.go                         # JSON utilities (100 lines)
│   ├── conversion.go                   # Data conversion (100 lines)
│   ├── crypto.go                       # Cryptographic utilities (120 lines)
│   ├── uuid.go                         # UUID utilities (60 lines)
│   └── http.go                         # HTTP utilities (80 lines)
├── constants/                          # Application constants
│   ├── app.go                          # App constants (100 lines)
│   ├── services.go                     # Service constants (120 lines)
│   ├── errors.go                       # Error constants (80 lines)
│   ├── config.go                       # Config constants (60 lines)
│   └── ports.go                        # Port constants (40 lines)
├── middleware/                         # Common middleware
│   ├── rate_limit.go                   # Rate limiting (120 lines)
│   ├── circuit_breaker.go              # Circuit breaker (150 lines)
│   ├── recovery.go                     # Panic recovery (80 lines)
│   ├── cors.go                         # CORS handling (60 lines)
│   ├── timeout.go                      # Request timeout (60 lines)
│   └── compression.go                  # Compression middleware (80 lines)
├── testing/                            # Testing utilities
│   ├── mocks.go                        # Mock implementations (150 lines)
│   ├── fixtures.go                     # Test fixtures (100 lines)
│   ├── helpers.go                      # Test helpers (120 lines)
│   └── integration.go                  # Integration test utilities (100 lines)
└── examples/                           # Usage examples
    ├── basic_service/                  # Basic service example
    │   ├── main.go                     # Example main.go (80 lines)
    │   ├── config.yaml                 # Example config
    │   └── README.md                   # Example documentation
    ├── advanced_service/               # Advanced service example
    │   ├── main.go                     # Advanced main.go (120 lines)
    │   ├── config.yaml                 # Advanced config
    │   └── README.md                   # Advanced documentation
    └── migration_guide.md              # Migration guide từ individual services
```

## 🚀 Sử dụng cơ bản

### 1. Configuration Management
```go
import "github.com/usc-platform/shared/config"

// Load configuration với service name
cfg, err := config.LoadConfig("", "service-name")
if err != nil {
    log.Fatalf("Failed to load config: %v", err)
}

// Sử dụng config
serverAddr := cfg.GetServerAddress()
dbDSN := cfg.GetDatabaseDSN()
redisAddr := cfg.GetRedisAddress()
```

### 2. Database Management
```go
import "github.com/usc-platform/shared/database"

// Initialize database manager
dbManager, err := database.NewManager(cfg)
if err != nil {
    log.Fatalf("Failed to initialize database: %v", err)
}

// Get connections
postgres := dbManager.PostgreSQL()
redis := dbManager.Redis()
clickhouse := dbManager.ClickHouse()
influxdb := dbManager.InfluxDB()
quickwit := dbManager.Quickwit()
```

### 3. gRPC Server Setup
```go
import "github.com/usc-platform/shared/grpc"
import "github.com/usc-platform/shared/logging"

// Initialize logger
logger := logging.NewLogger("service-name", cfg.Log)

// Create gRPC server với interceptors
grpcServer := grpc.NewServer(cfg, logger)

// Register health service
grpc.RegisterHealthService(grpcServer, "service-name", "1.0.0")

// Register reflection
grpc.RegisterReflection(grpcServer)

// Start server
grpcServer.Start()
```

### 4. Health Checking
```go
import "github.com/usc-platform/shared/health"

// Create health service
healthService := health.NewService("service-name", "1.0.0")

// Register health checks
healthService.RegisterCheck("database", dbManager.HealthCheck)
healthService.RegisterCheck("redis", redis.HealthCheck)

// Get health status
status := healthService.GetStatus()
```

### 5. Logging
```go
import "github.com/usc-platform/shared/logging"

// Initialize logger (using zap like most services)
logger := logging.NewLogger("service-name", cfg.Log)

// Structured logging
logger.Info("Service started", 
    logging.String("port", cfg.Server.Port),
    logging.String("version", "1.0.0"),
    logging.String("environment", cfg.Service.Environment),
)
```

### 6. Metrics
```go
import "github.com/usc-platform/shared/metrics"

// Initialize metrics
metrics.Init("service-name")

// Record metrics
metrics.IncrementCounter("requests_total", map[string]string{
    "service": "service-name",
    "method": "GetUser",
})
```

## 📦 Dependencies

### Core Dependencies (Production-Ready Versions)
```go
require (
    github.com/spf13/viper v1.21.0                    # Configuration management (latest stable)
    go.uber.org/zap v1.27.0                          # Logging (latest stable, used by most services)
    google.golang.org/grpc v1.75.1                   # gRPC framework (latest stable)
    google.golang.org/protobuf v1.36.6               # Protocol buffers (latest stable)
    github.com/go-redis/redis/v8 v8.11.5             # Redis client (stable for Redis 7+)
    github.com/lib/pq v1.10.9                        # PostgreSQL driver (compatible with PostgreSQL 16+)
    github.com/prometheus/client_golang v1.23.0      # Prometheus metrics (latest stable)
    github.com/google/uuid v1.6.0                    # UUID generation (latest stable)
    github.com/stretchr/testify v1.10.0              # Testing (latest stable)
)
```

### Database Dependencies (Production-Ready Versions)
```go
require (
    github.com/ClickHouse/clickhouse-go/v2 v2.40.1   # ClickHouse driver (compatible with ClickHouse 24+)
    github.com/influxdata/influxdb-client-go/v2 v2.14.0 # InfluxDB client (compatible with InfluxDB 2.7+)
    github.com/aws/aws-sdk-go v1.55.8                # AWS SDK (latest stable for cloud storage)
    cloud.google.com/go/storage v1.30.1              # Google Cloud Storage (latest stable)
    # Quickwit uses standard net/http (no external dependency needed)
)
```

### Development Dependencies (Production-Ready)
```go
require (
    github.com/golangci/golangci-lint v1.60.0         # Linting (latest stable)
    github.com/go-playground/validator/v10 v10.22.0   # Validation (latest stable)
    github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.2 # gRPC middleware (latest stable)
    github.com/grpc-ecosystem/grpc-gateway/v2 v2.25.0 # gRPC gateway (latest stable)
    github.com/golang-migrate/migrate/v4 v4.18.0      # Database migrations (latest stable)
    github.com/swaggo/swag v1.16.4                    # API documentation (latest stable)
)
```

## 🔧 Yêu cầu hệ thống (Production-Ready)

- **Go**: 1.24.4+ (latest stable, based on all 21 services)
- **PostgreSQL**: 16+ (latest stable with LTS support)
- **Redis**: 7.4+ (latest stable with performance improvements)
- **ClickHouse**: 24.8+ (latest stable with improved analytics)
- **InfluxDB**: 2.7+ (latest stable with better performance)
- **Quickwit**: Latest (uses standard HTTP, no version dependency)

## 📝 Ghi chú quan trọng (Production-Ready)

- **File size**: Tất cả files đều < 250 lines
- **Architecture**: Clean Architecture với SOLID principles
- **Environment**: Hỗ trợ đầy đủ environment variables
- **Logging**: Structured logging với JSON format (Zap)
- **Health**: Comprehensive health checking system
- **Metrics**: Prometheus metrics integration
- **Testing**: Comprehensive test coverage
- **Security**: Latest stable versions với security patches
- **Performance**: Optimized cho production workloads
- **Compatibility**: Tương thích với tất cả 21 services
- **LTS Support**: Sử dụng các phiên bản có hỗ trợ lâu dài

## 🚀 Migration từ individual services

### Bước 1: Update go.mod
```go
// Thêm shared library với Go 1.24.4 (latest stable)
go 1.23.4

require github.com/usc-platform/shared v1.0.0
```

### Bước 2: Update main.go
```go
// Thay thế individual imports
import (
    "github.com/usc-platform/shared/config"
    "github.com/usc-platform/shared/database"
    "github.com/usc-platform/shared/grpc"
    "github.com/usc-platform/shared/logging"
)
```

### Bước 3: Update config
```go
// Sử dụng shared config
cfg, err := config.LoadConfig("", "service-name")
```

### Bước 4: Update database
```go
// Sử dụng shared database manager
dbManager, err := database.NewManager(cfg)
```

### Bước 5: Update gRPC
```go
// Sử dụng shared gRPC server
grpcServer := grpc.NewServer(cfg, logger)
```

## 📞 Hỗ trợ

- **Documentation**: Inline Go documentation
- **Examples**: Usage examples trong examples/ directory
- **Testing**: Comprehensive test coverage
- **Migration**: Migration guide trong examples/migration_guide.md

---

**🎯 Mục tiêu: Tạo shared library hoàn chỉnh cho 21 microservices, đảm bảo 100% code reuse và consistency!**
