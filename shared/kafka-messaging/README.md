# Kafka Messaging Package

This package provides comprehensive Kafka messaging support for the USC Platform shared library, enabling reliable event streaming and message processing across all microservices.

## Features

### 🚀 **Producer Features**
- **Single Message Publishing**: Send individual messages to Kafka topics
- **Batch Message Publishing**: Efficiently send multiple messages in batches
- **Message Headers**: Support for custom message headers
- **JSON Serialization**: Built-in JSON message serialization
- **Configurable Options**: Compression, batching, retries, and timeouts
- **Metrics Integration**: Automatic metrics collection for monitoring

### 📥 **Consumer Features**
- **Topic Subscription**: Subscribe to Kafka topics with message handlers
- **Consumer Groups**: Support for consumer group coordination
- **Offset Management**: Automatic and manual offset management
- **Configurable Polling**: Customizable polling intervals and batch sizes
- **Error Handling**: Robust error handling with retry logic
- **Graceful Shutdown**: Clean shutdown with proper resource cleanup

### 🛠️ **Admin Features**
- **Topic Management**: Create, delete, and list Kafka topics
- **Partition Configuration**: Configure partitions and replication factors
- **Health Monitoring**: Comprehensive health checks for Kafka connectivity
- **Connection Management**: Automatic connection pooling and reconnection

### 📊 **Monitoring & Observability**
- **Health Checks**: Multiple health check types for different scenarios
- **Prometheus Metrics**: Comprehensive metrics collection
- **Structured Logging**: Detailed logging with context
- **Error Tracking**: Categorized error tracking and reporting

## Quick Start

### 1. Configuration

Add Kafka configuration to your service config:

```yaml
kafka:
  brokers: ["kafka-1:9092", "kafka-2:9092", "kafka-3:9092"]
  client_id: "your-service"
  group_id: "your-service-group"
  security_protocol: "PLAINTEXT"
  compression_type: "snappy"
  batch_size: 16384
  retries: 3
```

### 2. Initialize Kafka Manager

```go
import (
    "github.com/usc-platform/shared/config"
    "github.com/usc-platform/shared/messaging"
)

// Load configuration
cfg, err := config.LoadConfig("", "your-service")
if err != nil {
    log.Fatal(err)
}

// Create Kafka manager
kafkaManager, err := messaging.NewManager(cfg)
if err != nil {
    log.Fatal(err)
}
defer kafkaManager.Close()
```

### 3. Publishing Messages

```go
// Send single message
err := kafkaManager.SendMessage(ctx, "user-events", "user123", []byte("user data"))

// Send message with headers
headers := map[string]string{
    "content-type": "application/json",
    "source": "user-service",
}
err := kafkaManager.SendMessageWithHeaders(ctx, "user-events", "user123", data, headers)

// Send JSON message
userData := map[string]interface{}{
    "id": "user123",
    "name": "John Doe",
    "email": "john@example.com",
}
err := kafkaManager.SendJSONMessage(ctx, "user-events", "user123", userData)

// Send batch messages
messages := []messaging.Message{
    {Key: "key1", Value: []byte("value1")},
    {Key: "key2", Value: []byte("value2")},
}
err := kafkaManager.SendBatchMessages(ctx, "user-events", messages)
```

### 4. Consuming Messages

```go
// Define message handler
handler := func(ctx context.Context, message messaging.Message) error {
    // Process the message
    fmt.Printf("Received message: %s\n", string(message.Value))
    return nil
}

// Subscribe to topic
err := kafkaManager.Subscribe(ctx, "user-events", "user-service-group", handler)

// Subscribe with custom options
options := messaging.ConsumerOptions{
    Topic:     "user-events",
    GroupID:   "user-service-group",
    Handler:   handler,
    MinBytes:  1024,
    MaxBytes:  1048576,
    MaxWait:   1 * time.Second,
}
err := kafkaManager.SubscribeWithOptions(ctx, options)
```

### 5. Topic Management

```go
// Create topic
err := kafkaManager.CreateTopic(ctx, "new-topic", 3, 1)

// List topics
topics, err := kafkaManager.ListTopics(ctx)

// Delete topic
err := kafkaManager.DeleteTopic(ctx, "old-topic")
```

### 6. Health Checks

```go
import "github.com/usc-platform/shared/health"

// Create health checker
kafkaHealthChecker := health.NewKafkaHealthChecker("kafka", "Kafka connection health check", kafkaManager)

// Register with health service
healthService := health.NewService("my-service", "1.0.0")
healthService.RegisterCheck("kafka", kafkaHealthChecker)
```

### 7. Metrics Integration

```go
import "github.com/usc-platform/shared/metrics"

// Create metrics
kafkaMetrics := metrics.NewKafkaMetrics("my-service")

// Use metrics wrapper
producerWrapper := metrics.NewKafkaProducerWrapper(kafkaManager, kafkaMetrics)
consumerWrapper := metrics.NewKafkaConsumerWrapper(kafkaManager, kafkaMetrics)
```

## Configuration Options

### Kafka Configuration

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `brokers` | []string | ["localhost:9092"] | List of Kafka broker addresses |
| `client_id` | string | service name | Client identifier |
| `group_id` | string | service name + "-group" | Consumer group ID |
| `security_protocol` | string | "PLAINTEXT" | Security protocol (PLAINTEXT, SSL, SASL_PLAINTEXT, SASL_SSL) |
| `sasl_mechanism` | string | "PLAIN" | SASL mechanism |
| `sasl_username` | string | "" | SASL username |
| `sasl_password` | string | "" | SASL password |
| `ssl_ca_file` | string | "" | SSL CA certificate file |
| `ssl_cert_file` | string | "" | SSL client certificate file |
| `ssl_key_file` | string | "" | SSL client key file |
| `session_timeout` | int | 30000 | Session timeout in milliseconds |
| `heartbeat_interval` | int | 3000 | Heartbeat interval in milliseconds |
| `max_poll_records` | int | 500 | Maximum records per poll |
| `auto_offset_reset` | string | "latest" | Auto offset reset strategy |
| `enable_auto_commit` | bool | true | Enable automatic offset commits |
| `compression_type` | string | "snappy" | Compression type (none, gzip, snappy, lz4) |
| `batch_size` | int | 16384 | Batch size in bytes |
| `linger_ms` | int | 5 | Linger time in milliseconds |
| `retries` | int | 3 | Number of retries |
| `request_timeout` | int | 30000 | Request timeout in milliseconds |

### Producer Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `RequiredAcks` | int | 1 | Required acknowledgments |
| `Compression` | string | "snappy" | Compression type |
| `BatchSize` | int | 100 | Batch size |
| `BatchTimeout` | time.Duration | 10ms | Batch timeout |

### Consumer Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `Topic` | string | "" | Topic name |
| `GroupID` | string | "" | Consumer group ID |
| `Partition` | int | -1 | Partition number (-1 for all) |
| `Offset` | int64 | -1 | Starting offset (-1 for latest) |
| `MinBytes` | int | 10240 | Minimum bytes to fetch |
| `MaxBytes` | int | 10485760 | Maximum bytes to fetch |
| `MaxWait` | time.Duration | 1s | Maximum wait time |

## Message Structure

```go
type Message struct {
    Key     string            // Message key
    Value   []byte            // Message value
    Headers map[string]string // Message headers
    Topic   string            // Topic name
}
```

## Error Handling

The package provides comprehensive error handling:

- **Connection Errors**: Automatic reconnection with exponential backoff
- **Publishing Errors**: Retry logic with configurable retry count
- **Consumer Errors**: Error handling in message handlers
- **Health Check Errors**: Detailed error reporting for monitoring

## Best Practices

### Performance
- Use batch publishing for high-volume messages
- Configure appropriate batch sizes and timeouts
- Use compression for large messages
- Monitor producer and consumer lag

### Reliability
- Implement proper error handling in message handlers
- Use dead letter queues for failed messages
- Monitor consumer group coordination
- Set appropriate timeouts and retry policies

### Security
- Use SSL/TLS for production environments
- Implement SASL authentication when needed
- Secure credential storage
- Network security policies

### Monitoring
- Set up comprehensive alerting
- Monitor message throughput and latency
- Track consumer lag and processing rates
- Monitor connection health and errors

## Examples

See the complete example in `shared/examples/kafka_service/` for a full-featured Kafka service implementation.

## Dependencies

- `github.com/segmentio/kafka-go` - Kafka client library
- `github.com/prometheus/client_golang` - Metrics collection
- `go.uber.org/zap` - Structured logging

## License

This package is part of the USC Platform shared library and follows the same license terms.
