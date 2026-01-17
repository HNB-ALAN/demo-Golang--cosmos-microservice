# 🏗️ USC Platform Shared Library

A comprehensive shared library for the USC Platform microservices architecture, providing common infrastructure components, utilities, and configurations for all 21 microservices.

## 🚀 **Quick Start**

### **1. Installation**
```bash
go get github.com/usc-platform/shared
```

### **2. Environment Setup**
```bash
# Copy environment template
cp examples/k8s/env.example .env

# Edit with your values
nano .env

# Set environment variables
source .env
```

### **3. Basic Usage**
```go
package main

import (
    "github.com/usc-platform/shared/config"
    "github.com/usc-platform/shared/database"
    "github.com/usc-platform/shared/logging"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig("", "your-service")
    if err != nil {
        log.Fatal(err)
    }

    // Initialize logger
    logger := logging.NewLogger("your-service", cfg.Log)

    // Initialize database
    dbManager, err := database.NewManager(cfg)
    if err != nil {
        logger.Fatal("Failed to initialize database", logging.Error(err))
    }

    logger.Info("Service started successfully")
}
```

## 🔒 **Security Requirements**

### **⚠️ CRITICAL: Environment Variables**

**NEVER** use hardcoded secrets! Always use environment variables:

```bash
# Required environment variables
export JWT_SECRET="your-super-secure-jwt-secret-key-must-be-at-least-32-characters-long"
export POSTGRES_PASSWORD="your-secure-postgres-password"
export REDIS_PASSWORD="your-secure-redis-password"
```

### **Security Checklist**
- [ ] All hardcoded secrets removed
- [ ] Environment variables properly set
- [ ] JWT secret is at least 32 characters
- [ ] Database passwords are strong and unique
- [ ] SSL/TLS enabled for production
- [ ] `.env` file not committed to Git

📖 **See [SECURITY.md](SECURITY.md) for detailed security guide**

## 📁 **Library Structure**

```
shared/
├── auth/                    # Authentication & Authorization
├── cache/                   # Caching System
├── config/                  # Configuration Management
├── constants/               # Constants
├── database/                # Database Connections
├── errors/                  # Error Handling
├── graphql/                 # GraphQL Federation
├── grpc/                    # gRPC Middleware
├── health/                  # Health Checks
├── kafka-messaging/         # Kafka Messaging
├── logging/                 # Logging System
├── metrics/                 # Metrics Collection
├── middleware/              # HTTP Middleware
├── monitoring/              # Monitoring Setup
├── notifications/           # Notification Channels
├── proto/                   # Protocol Buffers
├── testing/                 # Testing Framework
├── utils/                   # Common Utilities
├── validation/              # Validation Framework
└── examples/                # Usage Examples
```

## 🛠️ **Core Components**

### **Configuration Management**
```go
import "github.com/usc-platform/shared/config"

cfg, err := config.LoadConfig("", "service-name")
```

### **Database Management with Retry Logic**
```go
import "github.com/usc-platform/shared/database"

// Automatic retry with exponential backoff for all database connections
dbManager, err := database.NewManager(cfg)
postgres := dbManager.PostgreSQL()
redis := dbManager.Redis()
quickwit := dbManager.Quickwit()
```

### **Logging**
```go
import "github.com/usc-platform/shared/logging"

logger := logging.NewLogger("service-name", cfg.Log)
logger.Info("Service started", logging.String("port", cfg.Server.Port))
```

### **Health Checks**
```go
import "github.com/usc-platform/shared/health"

healthService := health.NewService("service-name", "1.0.0")
healthService.RegisterCheck("database", dbManager.HealthCheck)
```

### **gRPC Server**
```go
import "github.com/usc-platform/shared/grpc"

grpcServer := grpc.NewServer(cfg, logger)
grpc.RegisterHealthService(grpcServer, "service-name", "1.0.0")
```

## 📊 **Supported Databases**

- **PostgreSQL** - Primary database (metadata only, **NOT for vector operations**)
- **Redis** - Caching and sessions
- **ClickHouse** - Analytics
- **InfluxDB** - Time series data
- **Quickwit** - Search and logging
- **Vector Database** (Qdrant/Weaviate/Pinecone) - Vector embeddings and similarity search

## 🔄 **Resilience Features**

### **Automatic Retry Logic**
All database connections include automatic retry with exponential backoff:
- **Max Retries**: 5 attempts
- **Base Delay**: 2 seconds
- **Exponential Backoff**: 2s, 4s, 8s, 16s, 32s
- **Timeout**: 10 seconds per attempt

### **Health Checks**
- **Concurrent Health Checks**: All databases checked simultaneously
- **Timeout Protection**: 5-second timeout per health check
- **Graceful Degradation**: Services continue operating despite partial failures

### **Connection Pooling**
- **Multi-tier Pooling**: L1 (Memory) + L2 (Redis) + L3 (Database)
- **Dynamic Pool Sizing**: Automatic adjustment based on load
- **Health Monitoring**: Continuous pool health monitoring

## 🔧 **Configuration**

### **Environment Variables**
See [examples/k8s/env.example](examples/k8s/env.example) for all available configuration options.

### **Configuration Files**
```yaml
# config.yaml
service:
  name: "your-service"
  version: "1.0.0"
  environment: "production"

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "${POSTGRES_PASSWORD}"  # Use environment variable
  dbname: "usc_social_media"
```

## 🧪 **Testing**

### **Run Tests**
```bash
go test ./...
```

### **Test Coverage**
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### **Testing Utilities**
```go
import "github.com/usc-platform/shared/testing"

// HTTP testing
helper := testing.NewHTTPTestHelper(handler)
resp, err := helper.Get("/api/users", nil)

// gRPC testing
grpcHelper := testing.NewGRPCTestHelper(server)
err := grpcHelper.AssertGRPCError(grpcErr, codes.Internal)
```

## 📈 **Monitoring & Metrics**

### **Prometheus Metrics**
```go
import "github.com/usc-platform/shared/metrics"

metrics.Init("service-name")
metrics.IncrementCounter("requests_total", map[string]string{
    "service": "service-name",
    "method": "GetUser",
})
```

### **Health Checks**
```bash
curl http://localhost:8080/health
```

## 🚀 **Deployment**

### **Docker**
```dockerfile
FROM golang:1.24.4-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

### **Kubernetes**
See [examples/k8s/](examples/k8s/) for complete Kubernetes deployment configurations.

## 📚 **Examples**

### **Basic Service**
```bash
cd examples/basic_service
go run main.go
```

### **Advanced Service**
```bash
cd examples/advanced_service
go run main.go
```

### **GraphQL Service**
```bash
cd examples/graphql_service
go run main.go
```

### **Kafka Service**
```bash
cd examples/kafka_service
go run main.go
```

### **Notification Service**
```bash
cd examples/notification_service
go run main.go
```

## 🔄 **Migration Guide**

If you're migrating from individual service configurations to the shared library, see [examples/migration_guide.md](examples/migration_guide.md).

## 📋 **Requirements**

- **Go**: 1.24.4+
- **PostgreSQL**: 16+
- **Redis**: 7.4+
- **ClickHouse**: 24.8+
- **InfluxDB**: 2.7+
- **Quickwit**: Latest (uses standard HTTP, no version dependency)

## 🤝 **Contributing**

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 **License**

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 **Support**

- **Documentation**: [Wiki](https://wiki.usc-platform.com)
- **Security Issues**: [SECURITY.md](SECURITY.md)
- **Resilience Features**: [RESILIENCE-IMPROVEMENTS.md](RESILIENCE-IMPROVEMENTS.md)
- **Environment Setup**: [examples/k8s/ENVIRONMENT-SETUP.md](examples/k8s/ENVIRONMENT-SETUP.md)
- **Issues**: [GitHub Issues](https://github.com/usc-platform/shared/issues)

---

**🎯 The USC Platform Shared Library provides enterprise-grade infrastructure components for all 21 microservices with security, performance, and reliability built-in.**
