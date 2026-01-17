# Notification Service Example

This example demonstrates how to use the USC Platform Shared Library's Notification Channels components to create a comprehensive notification service.

## Features

- **Multi-Channel Support**: Email, SMS, Push, In-app, and Webhook notifications
- **Provider Management**: Multiple providers per channel with fallback support
- **Delivery Tracking**: Real-time delivery status and retry logic
- **Batch Processing**: High-throughput notification delivery
- **Template System**: Dynamic notification templates
- **Queue Integration**: Asynchronous notification processing

## Components Used

### Notification Channels
- `shared/notifications/channels.go` - Multi-channel notification system

### Core Infrastructure
- `shared/config` - Configuration management
- `shared/logging` - Structured logging
- `shared/middleware` - HTTP middleware (CORS, rate limiting)

## Quick Start

1. **Install Dependencies**
   ```bash
   go mod tidy
   ```

2. **Configure Environment Variables**
   ```bash
   # Email Configuration
   export EMAIL_PROVIDER="sendgrid"
   export EMAIL_API_KEY="your_sendgrid_api_key"
   export FROM_EMAIL="noreply@uscplatform.com"
   export FROM_NAME="USC Platform"

   # SMS Configuration
   export SMS_PROVIDER="twilio"
   export SMS_API_KEY="your_twilio_api_key"
   export FROM_PHONE="+1234567890"

   # Push Configuration
   export FIREBASE_PROJECT_ID="your_firebase_project_id"
   export FIREBASE_SERVICE_KEY="your_firebase_service_key"

   # Service Configuration
   export PORT="4002"
   ```

3. **Run the Service**
   ```bash
   go run main.go
   ```

4. **Test the Service**
   ```bash
   # Send a test notification
   curl -X POST http://localhost:4002/notifications/test \
     -H "Content-Type: application/json" \
     -d '{
       "channel": "email",
       "user_id": "user123",
       "type": "welcome"
     }'
   ```

## Configuration

The service uses YAML configuration with the following key sections:

### Channel Configuration
```yaml
channels:
  default_channel: "email"
  retry_attempts: 3
  retry_delay: "5s"
  
  email:
    provider: "sendgrid"
    api_key: ""
    from_email: "noreply@uscplatform.com"
    rate_limit_per_minute: 1000
```

### Provider Configuration
```yaml
push:
  firebase:
    project_id: ""
    service_key: ""
  apns:
    key_id: ""
    team_id: ""
    bundle_id: ""
```

## API Endpoints

### Notification Endpoints
- **POST** `/notifications/send` - Send single notification
- **POST** `/notifications/batch` - Send batch notifications
- **GET** `/notifications/delivery/:id` - Get delivery status
- **POST** `/notifications/retry/:id` - Retry failed notification
- **POST** `/notifications/test` - Send test notification

### Service Endpoints
- **GET** `/health` - Health check
- **GET** `/metrics` - Prometheus metrics

## Usage Examples

### Send Single Notification
```bash
curl -X POST http://localhost:4002/notifications/send \
  -H "Content-Type: application/json" \
  -d '{
    "id": "msg_123",
    "user_id": "user_456",
    "channel": "email",
    "type": "welcome",
    "title": "Welcome to USC Platform",
    "content": "Thank you for joining our platform!",
    "data": {
      "user_name": "John Doe",
      "activation_link": "https://uscplatform.com/activate"
    },
    "priority": "high"
  }'
```

### Send Batch Notifications
```bash
curl -X POST http://localhost:4002/notifications/batch \
  -H "Content-Type: application/json" \
  -d '[
    {
      "user_id": "user_1",
      "channel": "email",
      "type": "newsletter",
      "title": "Weekly Newsletter",
      "content": "Check out our latest updates!"
    },
    {
      "user_id": "user_2",
      "channel": "sms",
      "type": "alert",
      "title": "Security Alert",
      "content": "Unusual activity detected on your account"
    }
  ]'
```

### Get Delivery Status
```bash
curl http://localhost:4002/notifications/delivery/delivery_123
```

### Retry Failed Notification
```bash
curl -X POST http://localhost:4002/notifications/retry/delivery_123
```

## Supported Channels

### Email Notifications
- **Providers**: SendGrid, AWS SES, SMTP
- **Features**: HTML/Text templates, attachments, tracking
- **Rate Limits**: Configurable per provider

### SMS Notifications
- **Providers**: Twilio, AWS SNS, SMS Gateway
- **Features**: Unicode support, delivery receipts
- **Rate Limits**: Provider-specific limits

### Push Notifications
- **Providers**: Firebase, Apple Push Notification Service (APNs), Web Push
- **Features**: Rich media, deep linking, badge updates
- **Platforms**: iOS, Android, Web

### In-App Notifications
- **Storage**: Database-backed notifications
- **Features**: Real-time delivery, read receipts, persistence
- **TTL**: Configurable expiration

### Webhook Notifications
- **Features**: Custom endpoints, retry logic, signature verification
- **Security**: HMAC signature validation
- **Reliability**: Configurable retry policies

## Template System

### Email Templates
```html
<!-- templates/email/welcome.html -->
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
</head>
<body>
    <h1>Welcome {{.User.Name}}!</h1>
    <p>{{.Content}}</p>
    <a href="{{.Data.ActivationLink}}">Activate Account</a>
</body>
</html>
```

### SMS Templates
```text
<!-- templates/sms/alert.txt -->
{{.Title}}: {{.Content}}
Reply STOP to unsubscribe.
```

### Push Templates
```json
{
  "title": "{{.Title}}",
  "body": "{{.Content}}",
  "data": {
    "type": "{{.Type}}",
    "user_id": "{{.UserID}}"
  }
}
```

## Queue Integration

### Redis Queue
```yaml
queue:
  provider: "redis"
  redis:
    url: "redis://localhost:6379"
    db: 0
    pool_size: 10
```

### Kafka Queue
```yaml
queue:
  provider: "kafka"
  kafka:
    brokers: ["localhost:9092"]
    topic: "notifications"
    group_id: "notification-service"
```

## Batch Processing

### Configuration
```yaml
batch:
  enabled: true
  max_batch_size: 100
  batch_timeout: "5s"
  max_workers: 10
```

### Usage
```go
// Send batch notifications
messages := []*notifications.NotificationMessage{
    {UserID: "user1", Channel: "email", Title: "Newsletter"},
    {UserID: "user2", Channel: "sms", Title: "Alert"},
}

deliveries, err := notificationService.SendBatchNotifications(ctx, messages)
```

## Monitoring and Observability

### Metrics
- Notification delivery rates
- Channel performance metrics
- Error rates and retry counts
- Queue processing metrics

### Health Checks
- Provider connectivity
- Queue health
- Database connectivity
- Service status

### Logging
- Structured JSON logging
- Request/response logging
- Error tracking
- Performance metrics

## Security Features

### Authentication
- API key authentication
- JWT token validation
- Rate limiting per user

### Data Protection
- Sensitive data encryption
- PII data masking
- Secure template rendering

### Provider Security
- API key rotation
- Secure credential storage
- Provider-specific security policies

## Development

### Local Development
```bash
# Start with mock providers
export MOCK_PROVIDERS=true
go run main.go
```

### Testing
```bash
# Run tests
go test ./...

# Run integration tests
go test -tags=integration ./...
```

### Debug Mode
```bash
# Enable debug logging
export DEBUG_MODE=true
go run main.go
```

## Production Deployment

### Docker Support
```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o notification-service main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/notification-service .
CMD ["./notification-service"]
```

### Kubernetes Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: notification-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: notification-service
  template:
    metadata:
      labels:
        app: notification-service
    spec:
      containers:
      - name: notification-service
        image: notification-service:latest
        ports:
        - containerPort: 4002
        env:
        - name: EMAIL_API_KEY
          valueFrom:
            secretKeyRef:
              name: notification-secrets
              key: email-api-key
        - name: SMS_API_KEY
          valueFrom:
            secretKeyRef:
              name: notification-secrets
              key: sms-api-key
```

## Best Practices

### Performance
- Use batch processing for high-volume notifications
- Implement connection pooling
- Use async processing for non-critical notifications
- Monitor queue depths and processing rates

### Reliability
- Implement retry logic with exponential backoff
- Use dead letter queues for failed notifications
- Monitor delivery rates and error patterns
- Implement circuit breakers for external providers

### Security
- Rotate API keys regularly
- Use environment variables for sensitive data
- Implement rate limiting per user
- Validate all input data

### Monitoring
- Set up comprehensive alerting
- Monitor delivery success rates
- Track provider performance
- Monitor queue processing metrics

## Troubleshooting

### Common Issues

1. **Provider Authentication Failed**
   - Check API key configuration
   - Verify provider credentials
   - Check network connectivity

2. **High Delivery Failure Rate**
   - Check provider rate limits
   - Verify message format
   - Check provider status

3. **Queue Processing Delays**
   - Monitor queue depths
   - Check worker performance
   - Verify database connectivity

### Debug Commands

```bash
# Check service health
curl http://localhost:4002/health

# View metrics
curl http://localhost:9092/metrics

# Test notification
curl -X POST http://localhost:4002/notifications/test \
  -H "Content-Type: application/json" \
  -d '{"channel": "email", "user_id": "test"}'
```

## Contributing

1. Follow the coding standards
2. Add comprehensive tests
3. Update documentation
4. Submit pull requests

## License

This example is part of the USC Platform Shared Library and follows the same license terms.
