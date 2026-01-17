# Advanced Service Example

This is an advanced example of how to use the USC shared library to create a production-ready microservice with comprehensive features.

## Features

- **Configuration Management**: Advanced configuration with environment variables and secrets
- **Database Connections**: All database types (PostgreSQL, Redis, MongoDB, ClickHouse, InfluxDB, Quickwit)
- **gRPC Server**: Production-ready gRPC server with interceptors and middleware
- **Health Checking**: Comprehensive health checks for all components
- **Structured Logging**: Advanced logging with structured fields and context
- **Metrics Collection**: Prometheus metrics with custom collectors
- **Authentication**: JWT-based authentication with refresh tokens
- **Caching**: Multi-level caching with Redis and in-memory
- **Validation**: Input validation with custom rules
- **Error Handling**: Comprehensive error handling and recovery
- **Background Services**: Automated maintenance and monitoring
- **Graceful Shutdown**: Proper shutdown handling with cleanup
- **Security**: Rate limiting, CORS, CSRF protection
- **Performance**: Compression, caching, optimization
- **Monitoring**: Health checks, metrics, alerting
- **Deployment**: Kubernetes, Docker, production-ready

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Load Balancer │    │   API Gateway   │    │   Service Mesh  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │ Advanced Service │
                    └─────────────────┘
                                 │
         ┌───────────────────────┼───────────────────────┐
         │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   PostgreSQL    │    │     Redis       │    │    MongoDB      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   ClickHouse    │    │    InfluxDB     │    │ Quickwit        │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Usage

### Development

1. **Start the service:**
   ```bash
   export USC_SERVICE_ENVIRONMENT=development
   export USC_LOG_LEVEL=debug
   go run main.go
   ```

2. **Check health:**
   ```bash
   curl http://localhost:8081/health
   ```

3. **View metrics:**
   ```bash
   curl http://localhost:9090/metrics
   ```

### Production

1. **Set environment variables:**
   ```bash
   export USC_SERVICE_ENVIRONMENT=production
   export USC_LOG_LEVEL=info
   export USC_LOG_FORMAT=json
   export POSTGRES_PASSWORD=your-secure-password
   export REDIS_PASSWORD=your-redis-password
   export JWT_SECRET=your-jwt-secret
   ```

2. **Start the service:**
   ```bash
   go run main.go
   ```

### Docker

1. **Build the image:**
   ```bash
   docker build -t advanced-service .
   ```

2. **Run the container:**
   ```bash
   docker run -p 8080:8080 -p 8081:8081 -p 9090:9090 \
     -e POSTGRES_PASSWORD=your-password \
     -e REDIS_PASSWORD=your-password \
     -e JWT_SECRET=your-secret \
     advanced-service
   ```

### Kubernetes

1. **Deploy to Kubernetes:**
   ```bash
   kubectl apply -f k8s/
   ```

2. **Check deployment:**
   ```bash
   kubectl get pods -l app=advanced-service
   kubectl get services -l app=advanced-service
   ```

## Configuration

### Environment Variables

The service supports extensive environment variable configuration:

```bash
# Service configuration
export USC_SERVICE_NAME=advanced-service
export USC_SERVICE_VERSION=1.0.0
export USC_SERVICE_ENVIRONMENT=production

# Server configuration
export USC_SERVER_HOST=0.0.0.0
export USC_SERVER_PORT=8080

# Database configuration
export USC_POSTGRES_HOST=postgres-cluster.internal
export USC_POSTGRES_PORT=5432
export USC_POSTGRES_USER=usc_user
export USC_POSTGRES_PASSWORD=your-password
export USC_POSTGRES_DB=usc_social_media

# Redis configuration
export USC_REDIS_HOST=redis-cluster.internal
export USC_REDIS_PORT=6379
export USC_REDIS_PASSWORD=your-password

# Logging configuration
export USC_LOG_LEVEL=info
export USC_LOG_FORMAT=json
export USC_LOG_OUTPUT=stdout

# Metrics configuration
export USC_METRICS_ENABLED=true
export USC_METRICS_PORT=9090

# Authentication configuration
export USC_JWT_SECRET=your-jwt-secret
export USC_JWT_EXPIRY=24h
export USC_REFRESH_EXPIRY=168h
```

### Configuration File

The service uses a comprehensive YAML configuration file with the following sections:

- **service**: Service metadata and environment settings
- **server**: Server configuration with gRPC and HTTP settings
- **database**: Database connection settings for all supported databases
- **log**: Logging configuration with rotation and formatting
- **metrics**: Metrics collection and Prometheus settings
- **auth**: Authentication and authorization settings
- **cache**: Caching configuration and TTL settings
- **rate_limiting**: Rate limiting configuration
- **circuit_breaker**: Circuit breaker settings
- **retry**: Retry logic configuration
- **timeout**: Timeout settings for different operations
- **monitoring**: Health checks and monitoring settings
- **security**: Security features and CORS settings
- **performance**: Performance optimization settings
- **backup**: Backup configuration
- **alerting**: Alerting and notification settings
- **feature_flags**: Feature flag configuration
- **external_services**: External service integration settings
- **kubernetes**: Kubernetes deployment configuration
- **docker**: Docker configuration
- **deployment**: Deployment and scaling settings

## Health Checks

The service includes comprehensive health checks for:

- **PostgreSQL**: Database connectivity and query execution
- **Redis**: Cache connectivity and operations
- **MongoDB**: Document database connectivity
- **ClickHouse**: Analytics database connectivity
- **InfluxDB**: Time series database connectivity
- **Quickwit**: Search engine connectivity
- **Cache**: Cache system health
- **Auth**: Authentication service health
- **Validator**: Input validation service health

### Health Check Endpoints

- **Basic Health**: `GET /health`
- **Detailed Health**: `GET /health/detailed`
- **Readiness**: `GET /ready`
- **Liveness**: `GET /live`

## Metrics

The service exposes comprehensive Prometheus metrics:

### HTTP Metrics
- `http_requests_total`: Total HTTP requests
- `http_request_duration_seconds`: HTTP request duration
- `http_request_size_bytes`: HTTP request size
- `http_response_size_bytes`: HTTP response size

### gRPC Metrics
- `grpc_requests_total`: Total gRPC requests
- `grpc_request_duration_seconds`: gRPC request duration
- `grpc_streams_total`: Total gRPC streams
- `grpc_stream_duration_seconds`: gRPC stream duration

### Database Metrics
- `database_queries_total`: Total database queries
- `database_query_duration_seconds`: Database query duration
- `database_connections`: Database connection count

### Cache Metrics
- `cache_hits_total`: Total cache hits
- `cache_misses_total`: Total cache misses
- `cache_operations_total`: Total cache operations
- `cache_operation_duration_seconds`: Cache operation duration

### Business Metrics
- `business_events_total`: Total business events
- `user_actions_total`: Total user actions
- `errors_total`: Total errors

### System Metrics
- `memory_alloc_bytes`: Memory allocation
- `memory_total_alloc_bytes`: Total memory allocation
- `memory_sys_bytes`: System memory
- `goroutines_count`: Goroutine count
- `gc_pause_total_ns`: GC pause time

## Logging

The service uses structured logging with the following features:

- **JSON Format**: Machine-readable log format
- **Structured Fields**: Consistent field naming and types
- **Context Propagation**: Request context throughout the call chain
- **Log Levels**: Debug, Info, Warn, Error, Fatal, Panic
- **Log Rotation**: Automatic log file rotation
- **Compression**: Log file compression
- **Performance**: High-performance logging with minimal overhead

### Log Fields

- **Service Fields**: `service`, `version`, `environment`
- **Request Fields**: `request_id`, `method`, `path`, `status_code`
- **Timing Fields**: `duration`, `start_time`, `end_time`
- **Database Fields**: `database`, `table`, `query_type`, `query_duration`
- **User Fields**: `user_id`, `user_email`, `user_role`
- **Error Fields**: `error_code`, `error_message`, `error_type`
- **Environment Fields**: `environment`, `region`, `instance`

## Security

The service includes comprehensive security features:

### Authentication
- **JWT Tokens**: Secure token-based authentication
- **Refresh Tokens**: Long-lived refresh tokens
- **Token Validation**: Comprehensive token validation
- **Session Management**: Secure session handling

### Authorization
- **Role-Based Access Control**: RBAC implementation
- **Permission System**: Fine-grained permissions
- **Resource Protection**: Resource-level access control

### Rate Limiting
- **Request Rate Limiting**: Per-endpoint rate limiting
- **User Rate Limiting**: Per-user rate limiting
- **IP Rate Limiting**: Per-IP rate limiting
- **Burst Protection**: Burst request protection

### Input Validation
- **Request Validation**: Comprehensive request validation
- **Data Sanitization**: Input data sanitization
- **SQL Injection Protection**: SQL injection prevention
- **XSS Protection**: Cross-site scripting prevention

### CORS
- **Cross-Origin Resource Sharing**: CORS configuration
- **Origin Validation**: Origin validation
- **Method Restrictions**: HTTP method restrictions
- **Header Restrictions**: Header restrictions

### CSRF Protection
- **CSRF Tokens**: CSRF token validation
- **SameSite Cookies**: SameSite cookie configuration
- **Origin Validation**: Origin validation

## Performance

The service includes performance optimization features:

### Caching
- **Multi-Level Caching**: Redis + in-memory caching
- **Cache Warming**: Automatic cache warming
- **Cache Invalidation**: Smart cache invalidation
- **Cache Compression**: Cache data compression

### Compression
- **Response Compression**: Gzip compression
- **Request Compression**: Request body compression
- **Configurable Levels**: Compression level configuration

### Connection Pooling
- **Database Pooling**: Database connection pooling
- **Redis Pooling**: Redis connection pooling
- **HTTP Pooling**: HTTP client connection pooling

### Background Services
- **Cache Warming**: Periodic cache warming
- **Database Maintenance**: Automated database maintenance
- **Metrics Collection**: System metrics collection
- **Health Monitoring**: Continuous health monitoring

## Monitoring

The service includes comprehensive monitoring:

### Health Monitoring
- **Health Checks**: Continuous health monitoring
- **Dependency Checks**: External dependency monitoring
- **Performance Monitoring**: Performance metrics collection
- **Error Monitoring**: Error tracking and alerting

### Metrics Collection
- **System Metrics**: CPU, memory, disk, network
- **Application Metrics**: Business and technical metrics
- **Database Metrics**: Database performance metrics
- **Cache Metrics**: Cache performance metrics

### Alerting
- **Webhook Alerts**: Webhook-based alerting
- **Slack Integration**: Slack channel notifications
- **Email Alerts**: Email notifications
- **Threshold Monitoring**: Configurable thresholds

### Tracing
- **Distributed Tracing**: Request tracing across services
- **Performance Tracing**: Performance bottleneck identification
- **Error Tracing**: Error root cause analysis

## Deployment

### Docker

```dockerfile
FROM golang:1.24.4-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o advanced-service main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/advanced-service .
COPY --from=builder /app/config.yaml .
CMD ["./advanced-service"]
```

### Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: advanced-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: advanced-service
  template:
    metadata:
      labels:
        app: advanced-service
    spec:
      containers:
      - name: advanced-service
        image: advanced-service:latest
        ports:
        - containerPort: 8080
        - containerPort: 8081
        - containerPort: 9090
        env:
        - name: USC_SERVICE_ENVIRONMENT
          value: "production"
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: postgres-secret
              key: password
        resources:
          requests:
            cpu: 100m
            memory: 256Mi
          limits:
            cpu: 500m
            memory: 512Mi
        livenessProbe:
          httpGet:
            path: /health
            port: 8081
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8081
          initialDelaySeconds: 5
          periodSeconds: 5
```

### Helm

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: advanced-service-config
data:
  config.yaml: |
    service:
      name: "advanced-service"
      version: "1.0.0"
      environment: "production"
    server:
      host: "0.0.0.0"
      port: "8080"
    # ... rest of configuration
```

## Testing

### Unit Tests

```bash
go test ./...
```

### Integration Tests

```bash
go test -tags=integration ./...
```

### Load Tests

```bash
go test -tags=load ./...
```

### End-to-End Tests

```bash
go test -tags=e2e ./...
```

## Development

### Prerequisites

- Go 1.24.4+
- Docker
- Kubernetes (for deployment)
- PostgreSQL 16+
- Redis 7.4+
- MongoDB 7.0+
- ClickHouse 24.8+
- InfluxDB 2.7+
- Quickwit (latest)

### Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/usc-platform/shared.git
   cd shared/examples/advanced_service
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Set up databases:**
   ```bash
   docker-compose up -d postgres redis mongodb clickhouse influxdb quickwit
   ```

4. **Run the service:**
   ```bash
   go run main.go
   ```

### Development Tools

- **Linting**: `golangci-lint run`
- **Formatting**: `go fmt ./...`
- **Vulnerability Scanning**: `govulncheck ./...`
- **Dependency Updates**: `go get -u ./...`

## Troubleshooting

### Common Issues

1. **Database Connection Issues**
   - Check database credentials
   - Verify network connectivity
   - Check firewall rules

2. **Memory Issues**
   - Monitor memory usage
   - Check for memory leaks
   - Adjust connection pool sizes

3. **Performance Issues**
   - Check database query performance
   - Monitor cache hit rates
   - Review log levels

4. **Health Check Failures**
   - Check service dependencies
   - Verify configuration
   - Review error logs

### Debugging

1. **Enable Debug Logging:**
   ```bash
   export USC_LOG_LEVEL=debug
   ```

2. **Enable Profiling:**
   ```bash
   export USC_ENABLE_PROFILING=true
   ```

3. **Check Metrics:**
   ```bash
   curl http://localhost:9090/metrics
   ```

4. **Review Health Status:**
   ```bash
   curl http://localhost:8081/health/detailed
   ```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License.

## Support

For support and questions:

- **Documentation**: [https://docs.usc-platform.com](https://docs.usc-platform.com)
- **Issues**: [https://github.com/usc-platform/shared/issues](https://github.com/usc-platform/shared/issues)
- **Discussions**: [https://github.com/usc-platform/shared/discussions](https://github.com/usc-platform/shared/discussions)
- **Email**: support@usc-platform.com