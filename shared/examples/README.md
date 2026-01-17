# USC Platform Shared Library Examples

This directory contains example services demonstrating how to use the USC Platform shared library components.

## Services Overview

### 1. Basic Service (Port 9090)
- **Purpose**: Demonstrates basic shared library usage
- **Dependencies**: PostgreSQL, Redis
- **Features**: gRPC server, health checks, basic database operations
- **Dockerfile**: `./basic_service/Dockerfile`

### 2. Advanced Service (Port 9091)
- **Purpose**: Demonstrates advanced shared library features
- **Dependencies**: PostgreSQL, Redis, MongoDB, ClickHouse, InfluxDB, Quickwit
- **Features**: Multi-database support, metrics collection, caching, authentication
- **Dockerfile**: `./advanced_service/Dockerfile`

### 3. GraphQL Service (Port 4001)
- **Purpose**: Demonstrates GraphQL federation capabilities
- **Dependencies**: None (standalone)
- **Features**: GraphQL endpoint, federation support, playground
- **Dockerfile**: `./graphql_service/Dockerfile`

### 4. Notification Service (Port 4002)
- **Purpose**: Demonstrates multi-channel notification system
- **Dependencies**: Redis
- **Features**: Email, SMS, Push notifications, webhooks
- **Dockerfile**: `./notification_service/Dockerfile`

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Go 1.22+ (for local development)

### Running All Services

```bash
# Start all services with databases
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down
```

### Running Individual Services

```bash
# Start only databases
docker-compose up -d postgres redis

# Start specific service
docker-compose up -d basic-service

# Build and start specific service
docker-compose up --build basic-service
```

## Service Endpoints

### Basic Service (gRPC)
- **Health Check**: `grpc://localhost:9090/grpc.health.v1.Health/Check`
- **Reflection**: Enabled for gRPC client tools

### Advanced Service (gRPC)
- **Health Check**: `grpc://localhost:9091/grpc.health.v1.Health/Check`
- **Metrics**: `http://localhost:9091/metrics`
- **Reflection**: Enabled for gRPC client tools

### GraphQL Service (HTTP)
- **GraphQL Endpoint**: `http://localhost:4001/graphql`
- **Playground**: `http://localhost:4001/graphql`
- **Health Check**: `http://localhost:4001/health`
- **Federation Info**: `http://localhost:4001/federation`

### Notification Service (HTTP)
- **Health Check**: `http://localhost:4002/health`
- **Send Notification**: `POST http://localhost:4002/notifications/send`
- **Batch Notifications**: `POST http://localhost:4002/notifications/batch`
- **Test Notification**: `POST http://localhost:4002/notifications/test`

## Monitoring

### Prometheus (Port 9092)
- **URL**: `http://localhost:9092`
- **Metrics**: Collected from all services

### Grafana (Port 3000)
- **URL**: `http://localhost:3000`
- **Login**: admin/admin
- **Dashboards**: Pre-configured for service monitoring

## Testing

### Health Checks
```bash
# Check all services health
docker-compose ps

# Test gRPC health
grpcurl -plaintext localhost:9090 grpc.health.v1.Health/Check

# Test HTTP health
curl http://localhost:4001/health
curl http://localhost:4002/health
```

### Load Testing
```bash
# Test GraphQL service
curl -X POST http://localhost:4001/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "{ __schema { types { name } } }"}'

# Test notification service
curl -X POST http://localhost:4002/notifications/test \
  -H "Content-Type: application/json" \
  -d '{"channel": "email", "user_id": "test-user", "type": "welcome"}'
```

## Configuration

### Environment Variables
- `EMAIL_API_KEY`: SendGrid API key for email notifications
- `SMS_API_KEY`: Twilio API key for SMS notifications
- `FIREBASE_PROJECT_ID`: Firebase project ID for push notifications
- `FIREBASE_SERVICE_KEY`: Firebase service account key
- `WEBHOOK_SECRET_KEY`: Secret key for webhook validation

### Database Configuration
All services use the same database instances but with different schemas:
- **PostgreSQL**: `usc_social_media` database
- **MongoDB**: `usc_social_media` database
- **ClickHouse**: `usc_analytics` database
- **InfluxDB**: `metrics` bucket in `usc` organization
- **Quickwit**: `usc_index` index

## Development

### Local Development
```bash
# Run individual services locally
cd basic_service
go run main.go

# With hot reload (requires air)
air -c .air.toml
```

### Building Images
```bash
# Build specific service
docker build -t usc-basic-service ./basic_service

# Build all services
docker-compose build
```

## Troubleshooting

### Common Issues

1. **Port Conflicts**: Ensure ports 3000, 4001, 4002, 5432, 6379, 9090, 9091, 9092, 7280, 27017, 9000, 8086 are available

2. **Database Connection Issues**: Wait for database health checks to pass before starting services

3. **Memory Issues**: Quickwit requires at least 512MB RAM. Adjust memory settings if needed

4. **Permission Issues**: Ensure Docker has proper permissions to access volumes

### Logs
```bash
# View specific service logs
docker-compose logs -f basic-service

# View all logs
docker-compose logs -f

# View database logs
docker-compose logs -f postgres redis
```

### Cleanup
```bash
# Stop and remove containers
docker-compose down

# Remove volumes (WARNING: This will delete all data)
docker-compose down -v

# Remove images
docker-compose down --rmi all
```

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Basic Service  │    │ Advanced Service│    │ GraphQL Service │    │Notification Svc │
│   (gRPC:9090)   │    │  (gRPC:9091)    │    │  (HTTP:4001)    │    │  (HTTP:4002)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │                       │
         └───────────────────────┼───────────────────────┼───────────────────────┘
                                 │                       │
                    ┌─────────────┴─────────────┐       │
                    │      Shared Databases     │       │
                    │  PostgreSQL, Redis, etc.  │       │
                    └───────────────────────────┘       │
                                                         │
                    ┌─────────────────────────────────────┴─────────────────────┐
                    │                Monitoring Stack                          │
                    │         Prometheus (9092) + Grafana (3000)              │
                    └─────────────────────────────────────────────────────────┘
```

Each service demonstrates different aspects of the shared library:
- **Basic Service**: Core functionality (auth, database, logging)
- **Advanced Service**: Full feature set (multi-db, metrics, caching)
- **GraphQL Service**: API gateway and federation
- **Notification Service**: External integrations and messaging