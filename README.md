# 🚀 Demo Golang Cosmos Microservice Platform

A comprehensive **22-microservice architecture** demonstration built with **Go**, **Cosmos Blockchain**, **CometBFT**, and **gRPC**. This project showcases enterprise-grade microservices patterns with complete infrastructure, production-ready patterns, and blockchain integration.

## 📋 Overview

This repository demonstrates a **production-ready microservices platform** that combines:

- **22 Business Microservices** - Gateway, Auth, User Management, Wallet, Blockchain Core, AI, Commerce, Analytics, and more
- **Cosmos Blockchain Integration** - USC custom blockchain with Cosmos SDK and CometBFT consensus
- **Complete Infrastructure** - PostgreSQL, Redis, ClickHouse, InfluxDB, Quickwit, Kafka, Prometheus, Grafana
- **Shared Library** - 100% code reuse across all services for configuration, database, gRPC, logging, authentication
- **Enterprise Patterns** - Service discovery, health checks, metrics, distributed logging, circuit breakers
- **Production-Ready** - Docker Compose orchestration, Kubernetes manifests, comprehensive security

**Key Statistics:**
- **Language Composition**: Go (64%), HTML (31.8%), Shell (3.2%), PLpgSQL (0.8%), Dockerfile (0.2%)
- **Total Services**: 22 microservices + 10+ infrastructure services
- **Architecture**: Event-driven, API gateway, blockchain-enabled
- **Deployment**: Docker Compose (development), Kubernetes-ready (production)

---

## 🎯 Key Features

### 🔧 Shared Library (100% Code Reuse)

The foundation for all 22 microservices with:

- ✅ **Configuration Management**: Environment-based config with validation
- ✅ **Multi-Database Support**:
  - Direct: PostgreSQL, Redis, ClickHouse, Qdrant, MinIO
  - Via gRPC: Quickwit (Service-16), InfluxDB (Service-17)
- ✅ **gRPC Framework**: Server setup, client management, interceptors, middleware
- ✅ **Security Layer**: JWT validation, password hashing (Argon2), encryption, OAuth2
- ✅ **Advanced Caching**: Multi-tier cache with Redis and in-memory backends
- ✅ **Notifications**: Multi-channel support (Email, SMS, Push, In-App, Webhooks)
- ✅ **Observability**:
  - Structured JSON logging with correlation IDs
  - Distributed tracing (OpenTelemetry compatible)
  - Prometheus metrics (HTTP, Kafka, Performance)
- ✅ **Middleware Stack**: Rate limiting, circuit breaker, retry (exponential backoff), timeout, CORS
- ✅ **Health Checks**: Service, database, and Kafka health monitoring
- ✅ **Kafka Integration**: Production-grade Kafka producer/consumer (KRaft mode)
- ✅ **Error Handling**: Unified error types with USC-specific error codes
- ✅ **Blockchain Support**: Cosmos SDK integration, USC chain constants

---

## 🏗️ Architecture

```
┌────────────────────────────────────────────────────────────────────────┐
│                     API Gateway (Service 01)                            │
│              (gRPC + GraphQL + REST + WebSocket)                        │
└────────────┬────────────────────────────────────┬──────────────────────┘
             │                                    │
     ┌───────▼────────┐              ┌────────────▼──────────┐
     │ Business Logic │              │  Blockchain Services  │
     │  Services      │              │  (Service 04)         │
     │  (02-22)       │              │  + CometBFT Node      │
     └───────┬────────┘              └────────────┬──────────┘
             │                                    │
     ┌───────▼──────────────────────────────────▼────────────┐
     │          Shared Library (USC Platform)                 │
     │  - Config, Database, gRPC, Auth, Logging, Metrics     │
     └───────┬──────────────────────────────────┬────────────┘
             │                                  │
     ┌───────▼──────────────────┐  ┌───────────▼──────────┐
     │   Database Layer          │  │  Message Queue       │
     │ - PostgreSQL (metadata)   │  │  - Kafka             │
     │ - Redis (cache/sessions)  │  │  - Event streaming   │
     │ - ClickHouse (analytics)  │  │  - KRaft mode        │
     │ - InfluxDB (metrics)      │  │                      │
     │ - Quickwit (search/logs)  │  │                      │
     │ - Qdrant (vectors)        │  │                      │
     │ - MinIO (object store)    │  │                      │
     └───────────────────────────┘  └──────────────────────┘
```

### Service Communication Patterns

1. **Synchronous**: gRPC calls between microservices
2. **Asynchronous**: Kafka events for cross-service communication
3. **Blockchain**: Cosmos SDK chain interactions via ABCI protocol
4. **Cache**: Redis for distributed caching and sessions
5. **Search**: Quickwit for full-text search on logs and documents

---

## 🛠️ Technology Stack

### Core Technologies
| Component | Technology | Version | Purpose |
|-----------|-----------|---------|---------|
| **Language** | Go | 1.23.4+ | Microservice implementation |
| **RPC Framework** | gRPC | 1.75.1 | Service-to-service communication |
| **Blockchain** | Cosmos SDK | Latest | Custom blockchain functionality |
| **Consensus** | CometBFT | 0.38.19 | Byzantine fault-tolerant consensus |
| **Protocol** | Protocol Buffers | 1.36.6 | Message serialization |
| **API Gateway** | Custom gRPC | - | Request routing, GraphQL, REST |

### Data Layer
| Database | Version | Purpose |
|----------|---------|---------|
| **PostgreSQL** | 16+ | Primary metadata storage |
| **Redis** | 7.4+ | Caching and session management |
| **ClickHouse** | 24.8+ | Analytics and data warehouse |
| **InfluxDB** | 2.7+ | Time-series metrics storage |
| **Quickwit** | Latest | Full-text search and logging |
| **Qdrant** | Latest | Vector embeddings & similarity search |
| **MinIO** | Latest | S3-compatible object storage |

### Infrastructure & Monitoring
| Service | Version | Purpose |
|---------|---------|---------|
| **Kafka** | 3.7.0 | Event streaming (KRaft mode) |
| **Prometheus** | 2.48.0 | Metrics collection |
| **Grafana** | 10.2.0 | Visualization and dashboards |
| **TorchServe** | Latest | ML model serving (GPU enabled) |
| **Docker Compose** | Latest | Local orchestration |

### Development & DevOps
- **Logging**: Uber Zap (structured JSON logging)
- **Metrics**: Prometheus client
- **Configuration**: Viper (environment variables + YAML)
- **Testing**: Testify, mocking frameworks
- **CI/CD**: GitHub Actions ready
- **Dependency Management**: Go modules

---

## 📁 Repository Structure

```
demo-Golang--cosmos-microservice/
├── README.md                           # This file
├── .env.example                        # Environment variables template
├── .gitignore                          # Git ignore rules
│
├── SERVICES/                           # All 22 microservices
│   ├── docker-compose.yml              # Complete infrastructure setup
│   ├── credentials/                    # Service credentials (BigQuery, etc)
│   ├── quickwit-config/                # Quickwit configuration
│   │
│   ├── service-01-gateway/             # 🔌 API Gateway
│   │   ├── handlers/                   # API endpoints
│   │   ├── middleware/                 # Request processing
│   │   ├── proto/                      # gRPC definitions
│   │   └── main.go
│   │
│   ├── service-02-auth/                # 🔐 Authentication
│   ├── service-03-user/                # 👤 User Management
│   ├── service-04-usc-blockchain-core/ # ⛓️ Blockchain Core
│   ├── service-04-cometbft/            # 🤝 CometBFT Node
│   ├── service-05-usc-wallet/          # 💰 Wallet Management
│   ├── service-06-security/            # 🛡️ Security Service
│   ├── service-07-caching/             # ⚡ Caching Service
│   ├── service-08-monitoring/          # 📊 Monitoring
│   ├── service-09-social/              # 👥 Social Features
│   ├── service-10-usc-bilateral-rewards/ # 🎁 Rewards
│   ├── service-11-content-management/  # 📝 Content
│   ├── service-12-video-service/       # 🎬 Video Service
│   ├── service-13-ai-service/          # 🤖 AI Service
│   ├── service-14-commerce-service/    # 🛒 Commerce
│   ├── service-15-notification-service/ # 📬 Notifications
│   ├── service-16-search-service/      # 🔍 Search (Quickwit)
│   ├── service-17-analytics-service/   # 📈 Analytics (BigQuery)
│   ├── service-18-moderation-service/  # 🚨 Moderation
│   ├── service-19-recommendation-service/ # 💡 Recommendations
│   └── service-20-advertising-service/ # 📢 Advertising
│
└── shared/                             # Shared Library (100% code reuse)
    ├── README.md                       # Shared library documentation
    ├── SECURITY.md                     # Security guidelines
    ├── SHARING-FILE.md                 # Library structure
    ├── go.mod & go.sum                 # Dependencies
    │
    ├── config/                         # Configuration Management
    ├── database/                       # Multi-Database Support
    ├── grpc/                           # gRPC Utilities
    ├── auth/                           # Authentication & Authorization
    ├── logging/                        # Structured Logging (Zap)
    ├── metrics/                        # Prometheus Metrics
    ├── health/                         # Health Checks
    ├── cache/                          # Caching Layer
    ├── errors/                         # Error Handling
    ├── validation/                     # Input Validation
    ├── middleware/                     # Common Middleware
    ├── utils/                          # Utility Functions
    ├── constants/                      # Application Constants
    ├── testing/                        # Testing Utilities
    ├── monitoring/                     # Monitoring Setup
    ├── notifications/                  # Notification Channels
    ├── kafka-messaging/                # Kafka Integration
    ├── graphql/                        # GraphQL Federation
    ├── proto/                          # Protocol Buffers
    ├── docs/                           # Documentation
    ├── examples/                       # Usage Examples
    │   ├── basic_service/              # Simple service example
    │   ├── advanced_service/           # Full-featured example
    │   ├── graphql_service/            # GraphQL integration
    │   ├── kafka_service/              # Kafka messaging
    │   ├── notification_service/       # Multi-channel notifications
    │   └── k8s/                        # Kubernetes deployment examples
    └── scripts/                        # Build and utility scripts
```

---

## 🚀 Quick Start

### Prerequisites

```bash
- Go 1.23.4+
- Docker & Docker Compose (latest)
- PostgreSQL 16+ (Docker provides it)
- Redis 7.4+ (Docker provides it)
- Git
```

### 1. Clone Repository

```bash
git clone https://github.com/HNB-ALAN/demo-Golang--cosmos-microservice.git
cd demo-Golang--cosmos-microservice
```

### 2. Setup Environment

```bash
# Copy environment template
cp .env.example .env

# Edit configuration with your values
nano .env

# Required variables to set:
# - POSTGRES_PASSWORD (secure password)
# - REDIS_PASSWORD (secure password)
# - JWT_SECRET (min 32 characters)
# - JWT_SECRET_BLOCKCHAIN (min 32 characters)
```

### 3. Start Infrastructure

```bash
cd SERVICES

# Start all infrastructure services
docker-compose up -d

# Verify all services are healthy
docker-compose ps

# Check specific service health
docker-compose logs -f postgres
docker-compose logs -f redis-master-1
docker-compose logs -f kafka-broker-1
```

### 4. Verify Services

```bash
# Wait for all services to be ready (~30 seconds)
sleep 30

# Check Prometheus
curl http://localhost:9091/-/healthy

# Check Redis
redis-cli -h localhost -p 6379 ping

# Check Kafka
docker exec kafka-broker-1 kafka-broker-api-versions.sh --bootstrap-server localhost:9092
```

### 5. Access Services

| Service | URL | Credentials |
|---------|-----|-------------|
| GraphQL Gateway | http://localhost:4000/graphql | - |
| gRPC Gateway | localhost:8001 | - |
| Prometheus | http://localhost:9091 | - |
| Grafana | http://localhost:3000 | admin/admin |
| CometBFT RPC | http://localhost:26657 | - |
| MinIO Console | http://localhost:9031 | minioadmin/minioadmin |
| Quickwit | http://localhost:7280 | - |

---

## 🏗️ Core Components

### Shared Library - 100% Code Reuse

All 22 services use the shared library for consistent patterns:

**Configuration Management**
```go
cfg, err := config.LoadConfig("", "service-name")
if err != nil {
    log.Fatal(err)
}
```

**Multi-Database Support with Retry Logic**
```go
dbManager, err := database.NewManager(cfg)
if err != nil {
    logger.Fatal("Failed to initialize database", logging.Error(err))
}

postgres := dbManager.PostgreSQL()
redis := dbManager.Redis()
clickhouse := dbManager.ClickHouse()
```

**gRPC Server Setup**
```go
grpcServer := grpc.NewServer(cfg, logger)
grpc.RegisterHealthService(grpcServer, "service-name", "1.0.0")
grpc.RegisterMetrics(grpcServer)

// Start server
grpcServer.Start(cfg.Server.Port)
```

**Structured Logging with Correlation IDs**
```go
logger := logging.NewLogger("service-name", cfg.Log)
logger.Info("Service started", 
    logging.String("port", cfg.Server.Port),
    logging.String("correlationID", ctx.Value("correlation-id")),
)
```

**Authentication & Authorization**
```go
token, err := auth.GenerateJWT(userID, cfg.JWT.Secret)
if err != nil {
    return err
}

claims, err := auth.ValidateJWT(token, cfg.JWT.Secret)
if err != nil {
    return err
}
```

**Prometheus Metrics**
```go
metrics.Init("service-name")
metrics.IncrementCounter("requests_total", map[string]string{
    "method": "GetUser",
    "status": "200",
})
metrics.RecordDuration("request_duration_ms", duration)
```

**Health Checks**
```go
healthService := health.NewService("service-name", "1.0.0")
healthService.RegisterCheck("database", dbManager.HealthCheck)
healthService.RegisterCheck("redis", cacheManager.HealthCheck)
healthService.RegisterCheck("kafka", kafkaClient.HealthCheck)

// Expose health endpoint
healthService.Serve(":8080")
```

---

## 📊 Infrastructure Services

### Databases
- **PostgreSQL 16**: Metadata storage, user data, transactions
- **Redis 7.4**: Caching, session management, rate limiting state
- **ClickHouse 24.8**: Analytics, time-series data, reporting
- **InfluxDB 2.7**: Performance metrics, monitoring data
- **Quickwit**: Full-text search, log aggregation
- **Qdrant**: Vector embeddings, similarity search
- **MinIO**: S3-compatible object storage, file uploads

### Message Queue
- **Kafka 3.7.0**: Event streaming, KRaft mode (no Zookeeper), high throughput

### Monitoring & Observability
- **Prometheus 2.48.0**: Metrics collection and storage
- **Grafana 10.2.0**: Visualization, dashboards, alerting (admin/admin)

### Blockchain
- **CometBFT 0.38.19**: Byzantine Fault Tolerant consensus engine
- **Cosmos SDK**: Blockchain framework and utilities

### Machine Learning
- **TorchServe**: GPU-enabled model serving for AI Service-13

---

## ⚙️ Configuration

### Environment Variables (.env)

```bash
# Database Configuration
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=secure_postgres_password
POSTGRES_DB=usc_social_media

# Redis Configuration
REDIS_HOST=redis-master-1
REDIS_PORT=6379
REDIS_PASSWORD=secure_redis_password
REDIS_CLUSTER_NODES=redis-master-1:6379,redis-master-2:6380,redis-master-3:6381

# JWT Configuration (min 32 chars)
JWT_SECRET=your_super_secure_jwt_secret_key_must_be_at_least_32_characters_long
JWT_SECRET_BLOCKCHAIN=blockchain_jwt_secret_key_min_32_chars
JWT_EXPIRATION=24h

# Kafka Configuration
KAFKA_BROKERS=kafka-broker-1:9092,kafka-broker-2:9093,kafka-broker-3:9094
KAFKA_SCHEMA_REGISTRY=http://schema-registry:8081

# ClickHouse Configuration
CLICKHOUSE_HOST=clickhouse
CLICKHOUSE_PORT=9000
CLICKHOUSE_HTTP_PORT=8123
CLICKHOUSE_USER=default
CLICKHOUSE_PASSWORD=

# InfluxDB Configuration
INFLUX_HOST=influxdb
INFLUX_PORT=8086
INFLUX_TOKEN=influx_token
INFLUX_ADMIN_TOKEN=influx_admin_token
INFLUX_ORG=usc-platform
INFLUX_BUCKET=metrics

# AWS SES (optional, for email notifications)
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_REGION=us-east-1
SES_FROM_EMAIL=no-reply@usc-platform.com

# BigQuery (optional, for Analytics Service-17)
BIGQUERY_PROJECT_ID=your_project_id
BIGQUERY_DATASET=usc_analytics_warehouse
GOOGLE_APPLICATION_CREDENTIALS=/app/credentials/service-account-key.json

# Service Configuration
SERVICE_ENVIRONMENT=development
SERVICE_LOG_LEVEL=info
RUST_LOG=info
```

### Configuration Files
- `shared/config/` - Shared library default configurations
- `SERVICES/service-XX/config/` - Service-specific configurations
- `SERVICES/.env` - Runtime environment variables

---

## 🚀 Deployment

### Docker Compose (Development/Local)

```bash
cd SERVICES

# Start all services
docker-compose up -d

# View all running services
docker-compose ps

# View logs for a specific service
docker-compose logs -f service-01-gateway

# Rebuild a specific service
docker-compose up -d --build service-02-auth

# Stop all services
docker-compose down

# Stop all services and remove volumes (clean slate)
docker-compose down -v
```

### Kubernetes (Production)

Production Kubernetes manifests are available in `shared/examples/k8s/`:

```bash
cd shared/examples/k8s

# Deploy to Kubernetes
kubectl apply -f namespace.yaml
kubectl apply -f configmap.yaml
kubectl apply -f secrets.yaml
kubectl apply -f services/
kubectl apply -f deployments/

# Monitor deployment
kubectl get pods -n usc-platform
kubectl logs -f deployment/service-01-gateway -n usc-platform
```

---

## 🔒 Security

### Best Practices
✅ All secrets in environment variables (`.env` file)  
✅ JWT-based authentication with configurable expiration  
✅ Password hashing with Argon2  
✅ Rate limiting and circuit breakers  
✅ Input validation on all endpoints  
✅ SQL injection protection via parameterized queries  
✅ TLS/SSL enabled for production  
✅ Regular dependency updates with `go get -u`  
✅ `.env` file in `.gitignore` (never commit secrets)

### Security Checklist
- [ ] Change all default passwords in `.env`
- [ ] Set JWT_SECRET to unique 32+ character value
- [ ] Enable TLS for Kafka and Redis in production
- [ ] Review `shared/SECURITY.md` for guidelines
- [ ] Audit database user permissions
- [ ] Enable Prometheus authentication
- [ ] Implement rate limiting per service
- [ ] Setup API key rotation schedule

See `shared/SECURITY.md` for detailed security guidelines.

---

## 👨‍💻 Development

### Build a New Service

```go
package main

import (
    "github.com/HNB-ALAN/demo-Golang--cosmos-microservice/shared/config"
    "github.com/HNB-ALAN/demo-Golang--cosmos-microservice/shared/database"
    "github.com/HNB-ALAN/demo-Golang--cosmos-microservice/shared/logging"
)

func main() {
    cfg, err := config.LoadConfig("", "my-service")
    if err != nil {
        panic(err)
    }

    logger := logging.NewLogger("my-service", cfg.Log)
    
    dbManager, err := database.NewManager(cfg)
    if err != nil {
        logger.Fatal("Database initialization failed", logging.Error(err))
    }

    logger.Info("Service started successfully")
}
```

### Run Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run tests for a specific package
go test ./shared/config -v
```

### Code Standards

- Follow `golangci-lint` configuration
- Use `gofmt` for code formatting
- Implement error handling with `pkg/errors`
- Add tests for all exported functions
- Use context for cancellation
- Follow Go idioms and best practices

---

## 📈 Monitoring & Metrics

### Prometheus Endpoints

```bash
# Service metrics (X = service number)
http://localhost:900X/metrics

# Gateway metrics
http://localhost:9001/metrics

# Prometheus dashboard
http://localhost:9091
```

### Grafana Dashboards

```bash
# Access Grafana
http://localhost:3000
# Login: admin / admin

# Dashboards included:
# - Service Overview
# - Request Latency
# - Error Rates
# - Database Performance
# - Kafka Throughput
# - Redis Hit Rates
```

### Health Checks

```bash
# Service health endpoint
curl http://localhost:800X/health

# Gateway health
curl http://localhost:8001/health

# Response includes:
# - Service status
# - Database connectivity
# - Kafka cluster status
# - Cache availability
```

---

## 🔄 Resilience Features

### Automatic Retry Logic
All database connections include automatic retry with exponential backoff:
- **Max Retries**: 5 attempts
- **Base Delay**: 1 second
- **Exponential Backoff**: 1s, 2s, 4s, 8s, 16s
- **Timeout**: 30 seconds per attempt

### Health Checks
- **Concurrent Checks**: All dependencies checked simultaneously
- **Timeout Protection**: 5-second timeout per health check
- **Graceful Degradation**: Services continue with partial failures

### Connection Pooling
- **Dynamic Pool Sizing**: Automatic adjustment based on load
- **Health Monitoring**: Continuous pool health verification
- **Connection Reuse**: Efficient resource utilization

---

## 📚 Examples

### Basic Service
```bash
cd shared/examples/basic_service
go run main.go
```

### Advanced Service with Database
```bash
cd shared/examples/advanced_service
go run main.go
```

### GraphQL Service
```bash
cd shared/examples/graphql_service
go run main.go
```

### Kafka Messaging Service
```bash
cd shared/examples/kafka_service
go run main.go
```

### Multi-Channel Notification Service
```bash
cd shared/examples/notification_service
go run main.go
```

---

## 📞 Support & Documentation

- **Shared Library**: `shared/README.md`
- **Security Guide**: `shared/SECURITY.md`
- **Library Structure**: `shared/SHARING-FILE.md`
- **Examples**: `shared/examples/`
- **Documentation**: `shared/docs/`
- **Issues**: [GitHub Issues](https://github.com/HNB-ALAN/demo-Golang--cosmos-microservice/issues)

---

## 🎯 Key Stats

- **22 Microservices** + 10+ Infrastructure Services
- **7 Database Engines** with automatic retry logic and health checks
- **100% Code Reuse** via shared library
- **Full Blockchain Integration** with Cosmos SDK and CometBFT
- **Production-Ready** patterns, configurations, and security
- **Event-Driven Architecture** with Kafka (KRaft mode)
- **Observable** with Prometheus, Grafana, and structured logging

---

## 🤝 Contributing

This is a demo implementation showcasing USC Platform architecture. Contributions are welcome:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Areas for contribution:
- Code improvements and optimizations
- Documentation enhancements
- Additional service examples
- Performance optimizations
- Test coverage improvements

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🎉 Welcome to USC Platform - Production-Ready Microservices Demo!

**Last Updated**: 2026-07-09  
**Version**: 1.1.0  
**Language Composition**: Go (64%) | HTML (31.8%) | Shell (3.2%) | PLpgSQL (0.8%) | Dockerfile (0.2%)

For more information, see `shared/README.md` and explore `shared/examples/` for implementation patterns.

---

### Quick Links
- 📖 [Shared Library Guide](shared/README.md)
- 🔒 [Security Guidelines](shared/SECURITY.md)
- 🏗️ [Architecture Details](shared/SHARING-FILE.md)
- 📚 [Examples](shared/examples/)
- 🐳 [Docker Documentation](SERVICES/docker-compose.yml)
- ⛓️ [Blockchain Setup](SERVICES/service-04-usc-blockchain-core/)
