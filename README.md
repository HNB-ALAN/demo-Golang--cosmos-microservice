# 🚀 Demo Golang Cosmos Microservice Platform

A comprehensive **22-microservice architecture** demonstration built with **Go**, **Cosmos Blockchain**, **CometBFT**, and **gRPC**. This project showcases enterprise-grade microservices patterns with complete infrastructure, featuring a unified shared library, advanced database support, and production-ready deployment configurations.

---

## 📋 Table of Contents

- [Project Overview](#-project-overview)
- [Architecture](#-architecture)
- [Technology Stack](#-technology-stack)
- [Repository Structure](#-repository-structure)
- [Quick Start](#-quick-start)
- [Core Components](#-core-components)
- [Services Overview](#-services-overview)
- [Infrastructure Services](#-infrastructure-services)
- [Configuration](#-configuration)
- [Development Guide](#-development-guide)
- [Deployment](#-deployment)
- [Security](#-security)
- [Contributing](#-contributing)

---

## 🎯 Project Overview

This is a **production-ready demonstration** of a modern microservices platform that combines:

- **21+ Business Microservices** - Gateway, Auth, User Management, Wallet, Blockchain Core, etc.
- **Cosmos Blockchain Integration** - USC custom blockchain with Cosmos SDK and CometBFT consensus
- **Complete Infrastructure** - PostgreSQL, Redis, ClickHouse, InfluxDB, Quickwit, Kafka, Prometheus, Grafana
- **Enterprise Patterns** - Service discovery, health checks, metrics, distributed logging, circuit breakers
- **Shared Library** - 100% code reuse across all services for configuration, database, gRPC, logging, authentication

**Key Statistics:**
- Language Composition: Go (64%), HTML (31.8%), Shell (3.2%), PLpgSQL (0.8%), Dockerfile (0.2%)
- Total Services: 22 microservices + 10+ infrastructure services
- Architecture: Event-driven, API gateway, blockchain-enabled
- Deployment: Docker Compose, Kubernetes-ready

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                     API Gateway (Service 01)                     │
│              (gRPC + GraphQL + REST)                             │
└────────────┬────────────────────────────────────┬────────────────┘
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
     │ - ClickHouse (analytics)  │  │                      │
     │ - InfluxDB (metrics)      │  │                      │
     │ - Quickwit (search/logs)  │  │                      │
     │ - Qdrant (vectors)        │  │                      │
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
| **API Gateway** | Custom gRPC | - | Request routing and GraphQL |

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
| **Docker Compose** | Latest | Orchestration |

### Development & DevOps
- **Logging**: Uber Zap (structured JSON logging)
- **Metrics**: Prometheus client
- **Configuration**: Viper (environment variables + YAML)
- **Testing**: Testify, mocking frameworks
- **CI/CD**: GitHub Actions ready

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
│   ├── credentials/                    # Service credentials
│   ├── quickwit-config/                # Quickwit configuration
│   │
│   ├── service-01-gateway/             # API Gateway
│   ├── service-02-auth/                # Authentication
│   ├── service-03-user/                # User Management
│   ├── service-04-usc-blockchain-core/ # Blockchain Core
│   ├── service-04-cometbft/            # CometBFT Node
│   ├── service-05-usc-wallet/          # Wallet Management
│   ├── service-06-security/            # Security Service
│   ├── service-07-caching/             # Caching Service
│   ├── service-08-monitoring/          # Monitoring
│   ├── service-09-social/              # Social Features
│   ├── service-10-usc-bilateral-rewards/ # Rewards
│   ├── service-11-content-management/  # Content
│   ├── service-12-video-service/       # Video Service
│   ├── service-13-ai-service/          # AI Service
│   ├── service-14-commerce-service/    # Commerce
│   ├── service-15-notification-service/ # Notifications
│   ├── service-16-search-service/      # Search
│   ├── service-17-analytics-service/   # Analytics
│   ├── service-18-moderation-service/  # Moderation
│   ├── service-19-recommendation-service/ # Recommendations
│   └── service-20-advertising-service/ # Advertising
│
└── shared/                             # Shared Library (100% code reuse)
    ├── README.md                       # Shared library docs
    ├── SECURITY.md                     # Security guidelines
    ├── SHARING-FILE.md                 # Library structure
    ├── go.mod & go.sum                 # Dependencies
    │
    ├── config/                         # Configuration Management
    ├── database/                       # Multi-Database Support
    ├── grpc/                           # gRPC Utilities
    ├── auth/                           # Authentication & Authorization
    ├── logging/                        # Structured Logging
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
    │   ├── basic_service/
    │   ├── advanced_service/
    │   ├── graphql_service/
    │   ├── kafka_service/
    │   ├── notification_service/
    │   └── k8s/                        # Kubernetes examples
    └── scripts/                        # Build and utility scripts
```

---

## 🚀 Quick Start

### Prerequisites

```bash
- Go 1.23.4+
- Docker & Docker Compose
- PostgreSQL 16+ (Docker provides it)
- Redis 7.4+ (Docker provides it)
```

### 1. Clone Repository

```bash
git clone https://github.com/HNB-ALAN/demo-Golang--cosmos-microservice.git
cd demo-Golang--cosmos-microservice
```

### 2. Setup Environment

```bash
cp .env.example .env
nano .env
# Edit these key variables:
# POSTGRES_PASSWORD, REDIS_PASSWORD, JWT_SECRET (min 32 chars), etc.
```

### 3. Start Infrastructure

```bash
cd SERVICES
docker-compose up -d
docker-compose ps  # Verify all services are healthy
```

### 4. Access Services

| Service | URL |
|---------|-----|
| GraphQL Gateway | http://localhost:4000 |
| gRPC Gateway | localhost:8001 |
| Prometheus | http://localhost:9091 |
| Grafana | http://localhost:3000 (admin/admin) |
| CometBFT RPC | http://localhost:26657 |
| MinIO | http://localhost:9031 |

---

## 🏗️ Core Components

### Shared Library - 100% Code Reuse

All 22 services use the shared library for:

**Configuration Management**
```go
cfg, err := config.LoadConfig("", "service-name")
```

**Multi-Database Support**
```go
dbManager, err := database.NewManager(cfg)
postgres := dbManager.PostgreSQL()
redis := dbManager.Redis()
clickhouse := dbManager.ClickHouse()
```

**gRPC Server Setup**
```go
grpcServer := grpc.NewServer(cfg, logger)
grpc.RegisterHealthService(grpcServer, "service-name", "1.0.0")
```

**Structured Logging**
```go
logger := logging.NewLogger("service-name", cfg.Log)
logger.Info("Service started", logging.String("port", cfg.Server.Port))
```

**Authentication**
```go
token := auth.GenerateJWT(userID, cfg.JWT.Secret)
```

**Metrics**
```go
metrics.IncrementCounter("requests_total", map[string]string{"method": "GetUser"})
```

---

## 📊 Infrastructure Services

### Databases
- **PostgreSQL 16**: Metadata storage
- **Redis 7.4**: Caching and sessions
- **ClickHouse 24.8**: Analytics
- **InfluxDB 2.7**: Time-series metrics
- **Quickwit**: Full-text search
- **Qdrant**: Vector similarity search
- **MinIO**: S3-compatible storage

### Message Queue
- **Kafka 3.7.0**: Event streaming (KRaft mode)

### Monitoring
- **Prometheus 2.48.0**: Metrics collection
- **Grafana 10.2.0**: Dashboards (port 3000)

### Blockchain
- **CometBFT 0.38.19**: Consensus engine

### ML
- **TorchServe**: GPU-enabled model serving

---

## ⚙️ Configuration

### Environment Variables (.env)

```bash
# Database
POSTGRES_PASSWORD=secure_password
REDIS_PASSWORD=secure_password

# JWT (min 32 chars)
JWT_SECRET=your_super_secure_jwt_secret_key_must_be_at_least_32_characters_long
JWT_SECRET_BLOCKCHAIN=blockchain_jwt_secret_key_min_32_chars

# InfluxDB
INFLUX_TOKEN=token
INFLUX_ADMIN_TOKEN=admin_token

# AWS SES (optional)
AWS_ACCESS_KEY_ID=key_id
AWS_SECRET_ACCESS_KEY=secret_key

# BigQuery (optional)
BIGQUERY_PROJECT_ID=project_id
```

---

## 🚀 Deployment

### Docker Compose (Development)

```bash
cd SERVICES
docker-compose up -d
docker-compose ps
docker-compose logs -f service-01-gateway
```

### Kubernetes (Production)

See `shared/examples/k8s/` for complete K8s manifests.

---

## 🔒 Security

### Best Practices
✅ All secrets in environment variables  
✅ JWT-based authentication  
✅ Rate limiting and circuit breakers  
✅ Input validation  
✅ SQL injection protection  
✅ Regular dependency updates  

See `shared/SECURITY.md` for detailed guidelines.

---

## 👨‍💻 Development

### Build a Service

```go
package main

import "github.com/usc-platform/shared/config"

func main() {
    cfg, _ := config.LoadConfig("", "my-service")
    // Build your service...
}
```

### Run Tests

```bash
go test ./...
go test -coverprofile=coverage.out ./...
```

---

## 📞 Support

- **Shared Library**: `shared/README.md`
- **Security**: `shared/SECURITY.md`
- **Examples**: `shared/examples/`
- **Docs**: `shared/docs/`

---

## 🎯 Key Stats

- **22 Microservices** + 10+ Infrastructure Services
- **6 Database Engines** with automatic retry logic
- **100% Code Reuse** via shared library
- **Full Blockchain Integration** with Cosmos/CometBFT
- **Production-Ready** patterns and configurations

---

**🎉 Welcome to the USC Platform - Production-Ready Microservices Demo!**

For more information, see `shared/README.md` and explore `shared/examples/`.
