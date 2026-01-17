# GraphQL Service Example

This example demonstrates how to use the USC Platform Shared Library's GraphQL Federation components to create a GraphQL service.

## Features

- **GraphQL Federation**: Service registration and schema federation
- **Query Middleware**: Complexity analysis, timeout handling, metrics collection
- **Performance Monitoring**: Query tracing and performance metrics
- **Security**: Rate limiting, CORS, and security middleware
- **Health Checks**: Service health monitoring and reporting

## Components Used

### GraphQL Federation
- `shared/graphql/federation.go` - Federation service management
- `shared/graphql/middleware.go` - Query middleware and optimization

### Core Infrastructure
- `shared/config` - Configuration management
- `shared/logging` - Structured logging
- `shared/middleware` - HTTP middleware (CORS, rate limiting)
- `shared/health` - Health check system

## Quick Start

1. **Install Dependencies**
   ```bash
   go mod tidy
   ```

2. **Configure Environment**
   ```bash
   export GATEWAY_URL="http://localhost:4000"
   export SERVICE_URL="http://localhost:4001"
   export PORT="4001"
   ```

3. **Run the Service**
   ```bash
   go run main.go
   ```

4. **Access GraphQL Playground**
   ```
   http://localhost:4001/graphql
   ```

## Configuration

The service uses YAML configuration with the following key sections:

### Federation Configuration
```yaml
federation:
  gateway_url: "http://localhost:4000"
  service_url: "http://localhost:4001"
  service_name: "graphql-service"
  introspection: true
  playground: true
```

### Middleware Configuration
```yaml
middleware:
  max_query_depth: 10
  max_query_complexity: 1000
  query_timeout: "30s"
  enable_tracing: true
  enable_metrics: true
```

## API Endpoints

### GraphQL Endpoint
- **POST** `/graphql` - GraphQL query execution
- **GET** `/graphql` - GraphQL Playground

### Service Endpoints
- **GET** `/health` - Health check
- **GET** `/federation` - Federation information
- **GET** `/metrics` - Prometheus metrics

## GraphQL Federation

### Service Registration
The service automatically registers with the federation gateway:

```go
serviceInfo := &graphql.ServiceInfo{
    Name:    "graphql-service",
    Version: "1.0.0",
    URL:     "http://localhost:4001",
    Health: &graphql.ServiceHealth{
        Status: "healthy",
        Timestamp: time.Now(),
    },
}

federationService.RegisterService(ctx, serviceInfo)
```

### Query Execution
Federated queries are executed through the federation service:

```go
request := &graphql.FederationRequest{
    Query: `
        query GetUsers {
            users {
                id
                name
                email
            }
        }
    `,
    Variables: map[string]interface{}{},
}

response, err := federationService.ExecuteFederatedQuery(ctx, request)
```

## Middleware Features

### Query Complexity Analysis
- Analyzes query depth and complexity
- Enforces limits to prevent resource exhaustion
- Provides detailed complexity metrics

### Query Timeout
- Configurable query timeout
- Prevents long-running queries from blocking the service
- Graceful timeout handling

### Performance Metrics
- Query execution time tracking
- Complexity and depth metrics
- Error rate monitoring
- Custom metrics collection

### Query Tracing
- Distributed tracing support
- Request flow analysis
- Performance bottleneck identification

## Security Features

### Rate Limiting
- Configurable rate limits per minute
- Burst handling
- IP-based rate limiting

### CORS Support
- Configurable CORS policies
- Origin validation
- Method and header restrictions

### Query Validation
- Query syntax validation
- Schema validation
- Security rule enforcement

## Monitoring and Observability

### Health Checks
- Service health monitoring
- Dependency health checks
- Health status reporting

### Metrics Collection
- Prometheus metrics integration
- Custom business metrics
- Performance metrics

### Logging
- Structured JSON logging
- Request/response logging
- Error tracking and reporting

## Development

### Hot Reload
- Schema hot reloading
- Configuration hot reloading
- Development mode optimizations

### Debug Mode
- Verbose logging
- Query debugging
- Performance profiling

### Mock Data
- Development mock data
- Test data generation
- Schema mocking

## Production Deployment

### Docker Support
```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o graphql-service main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/graphql-service .
CMD ["./graphql-service"]
```

### Kubernetes Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: graphql-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: graphql-service
  template:
    metadata:
      labels:
        app: graphql-service
    spec:
      containers:
      - name: graphql-service
        image: graphql-service:latest
        ports:
        - containerPort: 4001
        env:
        - name: GATEWAY_URL
          value: "http://gateway-service:4000"
        - name: SERVICE_URL
          value: "http://graphql-service:4001"
```

## Best Practices

### Query Optimization
- Use query complexity analysis
- Implement query caching
- Optimize database queries
- Use data loaders

### Security
- Implement proper authentication
- Use rate limiting
- Validate all inputs
- Monitor for suspicious queries

### Performance
- Monitor query performance
- Use connection pooling
- Implement caching strategies
- Optimize schema design

### Monitoring
- Set up comprehensive monitoring
- Use distributed tracing
- Monitor error rates
- Track business metrics

## Troubleshooting

### Common Issues

1. **Service Registration Failed**
   - Check gateway URL configuration
   - Verify network connectivity
   - Check service health status

2. **Query Timeout**
   - Increase timeout configuration
   - Optimize query complexity
   - Check database performance

3. **High Memory Usage**
   - Monitor query complexity
   - Implement query limits
   - Use connection pooling

### Debug Commands

```bash
# Check service health
curl http://localhost:4001/health

# Get federation info
curl http://localhost:4001/federation

# View metrics
curl http://localhost:9091/metrics
```

## Contributing

1. Follow the coding standards
2. Add comprehensive tests
3. Update documentation
4. Submit pull requests

## License

This example is part of the USC Platform Shared Library and follows the same license terms.
