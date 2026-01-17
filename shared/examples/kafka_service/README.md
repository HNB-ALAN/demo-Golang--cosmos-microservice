# Kafka Service Example

This example demonstrates how to use the USC Platform Shared Library's Kafka messaging components to create a comprehensive Kafka service for event streaming and message processing.

## Features

- **Event Publishing**: Single and batch event publishing to Kafka topics
- **Event Consumption**: Real-time event consumption with automatic topic routing
- **Topic Management**: Create, list, and delete Kafka topics
- **Health Monitoring**: Comprehensive health checks for Kafka connectivity
- **Metrics Collection**: Prometheus metrics for monitoring Kafka operations
- **REST API**: HTTP endpoints for Kafka operations
- **Event Processing**: Automatic event routing and processing based on event types

## Components Used

### Kafka Messaging
- `shared/messaging/kafka.go` - Kafka client implementation
- `shared/messaging/kafka_test.go` - Kafka client tests

### Configuration
- `shared/config` - Configuration management with Kafka settings
- `shared/config/defaults.go` - Kafka default configurations

### Health Checks
- `shared/health/kafka.go` - Kafka health check implementations

### Metrics
- `shared/metrics/kafka.go` - Kafka metrics collection

### Core Infrastructure
- `shared/logging` - Structured logging
- `shared/health` - Health check registry

## Quick Start

1. **Install Dependencies**
   ```bash
   go mod tidy
   ```

2. **Configure Environment Variables**
   ```bash
   # Kafka Configuration
   export KAFKA_BROKERS="localhost:9092"
   export KAFKA_CLIENT_ID="kafka-service"
   export KAFKA_GROUP_ID="kafka-service-group"
   
   # Service Configuration
   export SERVER_PORT="9093"
   export LOG_LEVEL="info"
   ```

3. **Start Kafka (if not already running)**
   ```bash
   # Using Docker Compose
   docker-compose up -d kafka
   
   # Or using Kafka directly
   # Start Zookeeper
   bin/zookeeper-server-start.sh config/zookeeper.properties
   
   # Start Kafka
   bin/kafka-server-start.sh config/server.properties
   ```

4. **Run the Service**
   ```bash
   go run main.go
   ```

5. **Test the Service**
   ```bash
   # Health check
   curl http://localhost:9093/health
   
   # List topics
   curl http://localhost:9093/api/v1/topics
   
   # Publish an event
   curl -X POST http://localhost:9093/api/v1/events \
     -H "Content-Type: application/json" \
     -d '{
       "type": "user.created",
       "data": {
         "user_id": "user123",
         "email": "user@example.com",
         "name": "John Doe"
       }
     }'
   ```

## Configuration

The service uses YAML configuration with the following key sections:

### Kafka Configuration
```yaml
kafka:
  brokers: ["kafka-1:9092", "kafka-2:9092", "kafka-3:9092"]
  client_id: "kafka-service"
  group_id: "kafka-service-group"
  security_protocol: "PLAINTEXT"
  compression_type: "snappy"
  batch_size: 16384
  retries: 3
```

### Service Configuration
```yaml
kafka_service:
  producer:
    topics:
      - "user-events"
      - "content-events"
      - "notification-events"
      - "analytics-events"
    batch_size: 100
    batch_timeout: "5s"
    
  consumer:
    topics:
      - "user-events"
      - "content-events"
      - "notification-events"
    group_id: "kafka-service-consumer"
    auto_offset_reset: "latest"
```

## API Endpoints

### Health & Monitoring
- **GET** `/health` - Health check
- **GET** `/metrics` - Prometheus metrics

### Event Management
- **POST** `/api/v1/events` - Publish single event
- **POST** `/api/v1/events/batch` - Publish batch events
- **GET** `/ws/events` - WebSocket for real-time events

### Topic Management
- **GET** `/api/v1/topics` - List all topics
- **POST** `/api/v1/topics` - Create new topic
- **DELETE** `/api/v1/topics/:name` - Delete topic

## Usage Examples

### Publish Single Event
```bash
curl -X POST http://localhost:9093/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{
    "type": "user.created",
    "data": {
      "user_id": "user123",
      "email": "user@example.com",
      "name": "John Doe"
    }
  }'
```

### Publish Batch Events
```bash
curl -X POST http://localhost:9093/api/v1/events/batch \
  -H "Content-Type: application/json" \
  -d '[
    {
      "type": "user.created",
      "data": {"user_id": "user1", "email": "user1@example.com"}
    },
    {
      "type": "content.created",
      "data": {"content_id": "content1", "title": "Sample Content"}
    }
  ]'
```

### Create Topic
```bash
curl -X POST http://localhost:9093/api/v1/topics \
  -H "Content-Type: application/json" \
  -d '{
    "name": "custom-events",
    "partitions": 3,
    "replication_factor": 1
  }'
```

### List Topics
```bash
curl http://localhost:9093/api/v1/topics
```

## Event Types

The service automatically routes events to appropriate topics based on their type:

### User Events → `user-events` topic
- `user.created`
- `user.updated`
- `user.deleted`

### Content Events → `content-events` topic
- `content.created`
- `content.updated`
- `content.deleted`

### Notification Events → `notification-events` topic
- `notification.sent`
- `notification.read`

### Analytics Events → `analytics-events` topic
- `analytics.view`
- `analytics.click`
- `analytics.conversion`

## Event Structure

All events follow this structure:

```json
{
  "id": "evt_1234567890_123456",
  "type": "user.created",
  "timestamp": "2024-01-15T10:30:00Z",
  "data": {
    "user_id": "user123",
    "email": "user@example.com",
    "name": "John Doe"
  },
  "source": "kafka-service"
}
```

## Monitoring and Observability

### Metrics

The service exposes Prometheus metrics for:

- **Producer Metrics**:
  - `kafka_messages_produced_total` - Total messages produced
  - `kafka_messages_produced_bytes_total` - Total bytes produced
  - `kafka_producer_errors_total` - Producer errors
  - `kafka_producer_latency_seconds` - Producer latency

- **Consumer Metrics**:
  - `kafka_messages_consumed_total` - Total messages consumed
  - `kafka_messages_consumed_bytes_total` - Total bytes consumed
  - `kafka_consumer_errors_total` - Consumer errors
  - `kafka_consumer_latency_seconds` - Consumer latency

- **Connection Metrics**:
  - `kafka_connection_status` - Connection status
  - `kafka_connection_errors_total` - Connection errors
  - `kafka_reconnections_total` - Reconnection count

### Health Checks

The service includes comprehensive health checks:

- **Kafka Connection Health**: Verifies Kafka connectivity
- **Topic Health**: Checks if required topics exist
- **Producer Health**: Tests message publishing capability
- **Consumer Health**: Tests message consumption capability

### Logging

Structured JSON logging with the following fields:

- `level` - Log level (info, warn, error)
- `timestamp` - Event timestamp
- `service` - Service name
- `message` - Log message
- `topic` - Kafka topic (when applicable)
- `event_id` - Event ID (when applicable)
- `error` - Error details (when applicable)

## Security Features

### Authentication
- JWT token validation for API endpoints
- API key authentication for Kafka operations

### Data Protection
- Event data encryption in transit
- Secure credential storage
- Input validation and sanitization

### Network Security
- TLS/SSL support for Kafka connections
- Network policies for service communication
- Rate limiting for API endpoints

## Performance Optimization

### Producer Optimization
- Batch message publishing
- Compression (Snappy, Gzip, LZ4)
- Async message sending
- Connection pooling

### Consumer Optimization
- Parallel message processing
- Offset management
- Consumer group coordination
- Backpressure handling

### Resource Management
- Memory-efficient message handling
- Connection reuse
- Garbage collection optimization
- CPU usage monitoring

## Development

### Local Development
```bash
# Start with mock Kafka
export MOCK_KAFKA=true
go run main.go

# Start with local Kafka
docker-compose up -d kafka
go run main.go
```

### Testing
```bash
# Run unit tests
go test ./...

# Run integration tests
go test -tags=integration ./...

# Run with coverage
go test -cover ./...
```

### Debug Mode
```bash
# Enable debug logging
export DEBUG_MODE=true
export LOG_LEVEL=debug
go run main.go
```

## Production Deployment

### Docker Support
```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o kafka-service main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/kafka-service .
CMD ["./kafka-service"]
```

### Kubernetes Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kafka-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: kafka-service
  template:
    metadata:
      labels:
        app: kafka-service
    spec:
      containers:
      - name: kafka-service
        image: kafka-service:latest
        ports:
        - containerPort: 9093
        env:
        - name: KAFKA_BROKERS
          value: "kafka-1:9092,kafka-2:9092,kafka-3:9092"
        - name: SERVER_PORT
          value: "9093"
```

## Best Practices

### Performance
- Use batch processing for high-volume events
- Implement connection pooling
- Use async processing for non-critical events
- Monitor queue depths and processing rates

### Reliability
- Implement retry logic with exponential backoff
- Use dead letter queues for failed events
- Monitor delivery rates and error patterns
- Implement circuit breakers for external dependencies

### Security
- Rotate API keys regularly
- Use environment variables for sensitive data
- Implement rate limiting per client
- Validate all input data

### Monitoring
- Set up comprehensive alerting
- Monitor delivery success rates
- Track topic performance
- Monitor consumer lag and processing metrics

## Troubleshooting

### Common Issues

1. **Kafka Connection Failed**
   - Check Kafka broker addresses
   - Verify network connectivity
   - Check authentication credentials

2. **High Message Loss Rate**
   - Check producer acknowledgments
   - Verify topic replication factor
   - Monitor broker health

3. **Consumer Lag Issues**
   - Check consumer group coordination
   - Monitor processing performance
   - Verify offset management

4. **Topic Creation Failed**
   - Check broker permissions
   - Verify replication factor
   - Check available disk space

### Debug Commands

```bash
# Check service health
curl http://localhost:9093/health

# View metrics
curl http://localhost:9093/metrics

# Test event publishing
curl -X POST http://localhost:9093/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{"type": "test", "data": {"message": "test"}}'

# List topics
curl http://localhost:9093/api/v1/topics
```

## Contributing

1. Follow the coding standards
2. Add comprehensive tests
3. Update documentation
4. Submit pull requests

## License

This example is part of the USC Platform Shared Library and follows the same license terms.
