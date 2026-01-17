# Basic Service Example

This is a basic example of how to use the USC shared library to create a simple microservice.

## Features

- Configuration management
- Database connections (PostgreSQL, Redis, MongoDB, ClickHouse, InfluxDB, Quickwit)
- gRPC server setup
- Health checking
- Structured logging
- Metrics collection

## Usage

1. **Start the service:**
   ```bash
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

## Configuration

The service uses a YAML configuration file (`config.yaml`) with the following sections:

- **service**: Service metadata (name, version, environment)
- **server**: Server configuration (host, port, gRPC settings)
- **database**: Database connection settings
- **redis**: Redis connection settings
- **mongodb**: MongoDB connection settings
- **clickhouse**: ClickHouse connection settings
- **influxdb**: InfluxDB connection settings
- **quickwit**: Quickwit connection settings
- **log**: Logging configuration
- **metrics**: Metrics configuration
- **auth**: Authentication settings
- **cache**: Cache settings

## Environment Variables

You can override any configuration value using environment variables:

```bash
export USC_SERVER_PORT=8080
export USC_DATABASE_HOST=localhost
export USC_REDIS_HOST=localhost
export USC_LOG_LEVEL=debug
```

## Health Checks

The service includes health checks for:

- PostgreSQL database
- Redis cache
- Service itself

## Metrics

The service exposes Prometheus metrics at `/metrics` endpoint.

## Logging

The service uses structured logging with JSON format by default.

## Database Connections

The service automatically connects to all configured databases:

- PostgreSQL (main database)
- Redis (caching)
- MongoDB (document storage)
- ClickHouse (analytics)
- InfluxDB (time series)
- Quickwit (search)

## gRPC Server

The service starts a gRPC server with:

- Health service
- Reflection enabled
- Interceptors for logging and metrics
- Graceful shutdown

## Error Handling

The service includes comprehensive error handling and logging.

## Security

The service includes basic security features:

- JWT authentication
- Rate limiting
- Input validation
- Error sanitization

## Monitoring

The service includes monitoring capabilities:

- Health checks
- Metrics collection
- Structured logging
- Performance monitoring

## Development

To run in development mode:

```bash
export USC_SERVICE_ENVIRONMENT=development
export USC_LOG_LEVEL=debug
go run main.go
```

## Production

To run in production mode:

```bash
export USC_SERVICE_ENVIRONMENT=production
export USC_LOG_LEVEL=info
export USC_LOG_FORMAT=json
go run main.go
```

## Docker

To run with Docker:

```bash
docker build -t basic-service .
docker run -p 8080:8080 -p 8081:8081 -p 9090:9090 basic-service
```

## Kubernetes

To deploy to Kubernetes:

```bash
kubectl apply -f k8s/
```

## Testing

To run tests:

```bash
go test ./...
```

## Building

To build the service:

```bash
go build -o basic-service main.go
```

## Dependencies

The service depends on:

- USC shared library
- gRPC
- PostgreSQL driver
- Redis client
- MongoDB driver
- ClickHouse driver
- InfluxDB client
- Quickwit client (uses standard HTTP)
- Prometheus client
- Zap logger
- Viper configuration
