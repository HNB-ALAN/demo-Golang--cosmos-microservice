package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	messaging "github.com/usc-platform/shared/kafka-messaging"
)

// KafkaMetrics represents Kafka-specific metrics
type KafkaMetrics struct {
	// Producer metrics
	MessagesProducedTotal prometheus.Counter
	MessagesProducedBytes prometheus.Counter
	ProducerErrorsTotal   prometheus.Counter
	ProducerLatency       prometheus.Histogram
	ProducerBatchSize     prometheus.Histogram

	// Consumer metrics
	MessagesConsumedTotal prometheus.Counter
	MessagesConsumedBytes prometheus.Counter
	ConsumerErrorsTotal   prometheus.Counter
	ConsumerLag           prometheus.Gauge
	ConsumerLatency       prometheus.Histogram
	ConsumerOffset        prometheus.Gauge

	// Connection metrics
	ConnectionStatus      prometheus.Gauge
	ConnectionErrorsTotal prometheus.Counter
	ReconnectionTotal     prometheus.Counter

	// Topic metrics
	TopicPartitions        prometheus.Gauge
	TopicReplicationFactor prometheus.Gauge
	TopicSize              prometheus.Gauge

	// Registry for custom metrics
	registry    prometheus.Registerer
	serviceName string
}

// NewKafkaMetrics creates a new Kafka metrics instance
func NewKafkaMetrics(serviceName string) *KafkaMetrics {
	labels := prometheus.Labels{"service": serviceName}

	return &KafkaMetrics{
		// Producer metrics
		MessagesProducedTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name:        "kafka_messages_produced_total",
			Help:        "Total number of messages produced to Kafka",
			ConstLabels: labels,
		}),
		MessagesProducedBytes: promauto.NewCounter(prometheus.CounterOpts{
			Name:        "kafka_messages_produced_bytes_total",
			Help:        "Total bytes of messages produced to Kafka",
			ConstLabels: labels,
		}),
		ProducerErrorsTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name:        "kafka_producer_errors_total",
			Help:        "Total number of producer errors",
			ConstLabels: labels,
		}),
		ProducerLatency: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:        "kafka_producer_latency_seconds",
			Help:        "Producer message latency in seconds",
			ConstLabels: labels,
			Buckets:     prometheus.ExponentialBuckets(0.001, 2, 15), // 1ms to 32s
		}),
		ProducerBatchSize: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:        "kafka_producer_batch_size",
			Help:        "Producer batch size",
			ConstLabels: labels,
			Buckets:     prometheus.ExponentialBuckets(1, 2, 12), // 1 to 4096
		}),

		// Consumer metrics
		MessagesConsumedTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name:        "kafka_messages_consumed_total",
			Help:        "Total number of messages consumed from Kafka",
			ConstLabels: labels,
		}),
		MessagesConsumedBytes: promauto.NewCounter(prometheus.CounterOpts{
			Name:        "kafka_messages_consumed_bytes_total",
			Help:        "Total bytes of messages consumed from Kafka",
			ConstLabels: labels,
		}),
		ConsumerErrorsTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name:        "kafka_consumer_errors_total",
			Help:        "Total number of consumer errors",
			ConstLabels: labels,
		}),
		ConsumerLag: promauto.NewGauge(prometheus.GaugeOpts{
			Name:        "kafka_consumer_lag",
			Help:        "Consumer lag in messages",
			ConstLabels: labels,
		}),
		ConsumerLatency: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:        "kafka_consumer_latency_seconds",
			Help:        "Consumer message processing latency in seconds",
			ConstLabels: labels,
			Buckets:     prometheus.ExponentialBuckets(0.001, 2, 15), // 1ms to 32s
		}),
		ConsumerOffset: promauto.NewGauge(prometheus.GaugeOpts{
			Name:        "kafka_consumer_offset",
			Help:        "Current consumer offset",
			ConstLabels: labels,
		}),

		// Connection metrics
		ConnectionStatus: promauto.NewGauge(prometheus.GaugeOpts{
			Name:        "kafka_connection_status",
			Help:        "Kafka connection status (1=connected, 0=disconnected)",
			ConstLabels: labels,
		}),
		ConnectionErrorsTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name:        "kafka_connection_errors_total",
			Help:        "Total number of connection errors",
			ConstLabels: labels,
		}),
		ReconnectionTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name:        "kafka_reconnections_total",
			Help:        "Total number of reconnections",
			ConstLabels: labels,
		}),

		// Topic metrics
		TopicPartitions: promauto.NewGauge(prometheus.GaugeOpts{
			Name:        "kafka_topic_partitions",
			Help:        "Number of partitions in topic",
			ConstLabels: labels,
		}),
		TopicReplicationFactor: promauto.NewGauge(prometheus.GaugeOpts{
			Name:        "kafka_topic_replication_factor",
			Help:        "Replication factor of topic",
			ConstLabels: labels,
		}),
		TopicSize: promauto.NewGauge(prometheus.GaugeOpts{
			Name:        "kafka_topic_size_bytes",
			Help:        "Size of topic in bytes",
			ConstLabels: labels,
		}),

		registry:    prometheus.DefaultRegisterer,
		serviceName: serviceName,
	}
}

// RecordMessageProduced records a produced message
func (km *KafkaMetrics) RecordMessageProduced(topic string, size int, latency time.Duration) {
	km.MessagesProducedTotal.Inc()
	km.MessagesProducedBytes.Add(float64(size))
	km.ProducerLatency.Observe(latency.Seconds())
}

// RecordMessageProducedBatch records a batch of produced messages
func (km *KafkaMetrics) RecordMessageProducedBatch(topic string, batchSize int, totalSize int, latency time.Duration) {
	km.MessagesProducedTotal.Add(float64(batchSize))
	km.MessagesProducedBytes.Add(float64(totalSize))
	km.ProducerLatency.Observe(latency.Seconds())
	km.ProducerBatchSize.Observe(float64(batchSize))
}

// RecordProducerError records a producer error
func (km *KafkaMetrics) RecordProducerError(topic string, errorType string) {
	km.ProducerErrorsTotal.Inc()
}

// RecordMessageConsumed records a consumed message
func (km *KafkaMetrics) RecordMessageConsumed(topic string, size int, latency time.Duration) {
	km.MessagesConsumedTotal.Inc()
	km.MessagesConsumedBytes.Add(float64(size))
	km.ConsumerLatency.Observe(latency.Seconds())
}

// RecordConsumerError records a consumer error
func (km *KafkaMetrics) RecordConsumerError(topic string, errorType string) {
	km.ConsumerErrorsTotal.Inc()
}

// SetConsumerLag sets the consumer lag
func (km *KafkaMetrics) SetConsumerLag(topic string, lag int64) {
	km.ConsumerLag.Set(float64(lag))
}

// SetConsumerOffset sets the consumer offset
func (km *KafkaMetrics) SetConsumerOffset(topic string, offset int64) {
	km.ConsumerOffset.Set(float64(offset))
}

// SetConnectionStatus sets the connection status
func (km *KafkaMetrics) SetConnectionStatus(connected bool) {
	if connected {
		km.ConnectionStatus.Set(1)
	} else {
		km.ConnectionStatus.Set(0)
	}
}

// RecordConnectionError records a connection error
func (km *KafkaMetrics) RecordConnectionError(errorType string) {
	km.ConnectionErrorsTotal.Inc()
}

// RecordReconnection records a reconnection
func (km *KafkaMetrics) RecordReconnection() {
	km.ReconnectionTotal.Inc()
}

// SetTopicPartitions sets the number of partitions for a topic
func (km *KafkaMetrics) SetTopicPartitions(topic string, partitions int) {
	km.TopicPartitions.Set(float64(partitions))
}

// SetTopicReplicationFactor sets the replication factor for a topic
func (km *KafkaMetrics) SetTopicReplicationFactor(topic string, replicationFactor int) {
	km.TopicReplicationFactor.Set(float64(replicationFactor))
}

// SetTopicSize sets the size of a topic
func (km *KafkaMetrics) SetTopicSize(topic string, size int64) {
	km.TopicSize.Set(float64(size))
}

// KafkaMetricsCollector represents a Kafka metrics collector
type KafkaMetricsCollector struct {
	metrics *KafkaMetrics
	client  messaging.KafkaClient
}

// NewKafkaMetricsCollector creates a new Kafka metrics collector
func NewKafkaMetricsCollector(serviceName string, client messaging.KafkaClient) *KafkaMetricsCollector {
	return &KafkaMetricsCollector{
		metrics: NewKafkaMetrics(serviceName),
		client:  client,
	}
}

// CollectMetrics collects Kafka metrics
func (kmc *KafkaMetricsCollector) CollectMetrics(ctx context.Context) error {
	if kmc.client == nil {
		return nil
	}

	// Check connection status
	if err := kmc.client.HealthCheck(ctx); err != nil {
		kmc.metrics.SetConnectionStatus(false)
		kmc.metrics.RecordConnectionError("health_check_failed")
		return err
	}

	kmc.metrics.SetConnectionStatus(true)

	// Collect topic information
	topics, err := kmc.client.ListTopics(ctx)
	if err != nil {
		kmc.metrics.RecordConnectionError("list_topics_failed")
		return err
	}

	// Update topic metrics (simplified - in production you'd get more detailed info)
	for _, topic := range topics {
		kmc.metrics.SetTopicPartitions(topic, 1)        // Default to 1 partition
		kmc.metrics.SetTopicReplicationFactor(topic, 1) // Default to 1 replica
		kmc.metrics.SetTopicSize(topic, 0)              // Default to 0 size
	}

	return nil
}

// GetMetrics returns the metrics instance
func (kmc *KafkaMetricsCollector) GetMetrics() *KafkaMetrics {
	return kmc.metrics
}

// KafkaProducerWrapper wraps a Kafka client with metrics
type KafkaProducerWrapper struct {
	client  messaging.KafkaClient
	metrics *KafkaMetrics
}

// NewKafkaProducerWrapper creates a new Kafka producer wrapper with metrics
func NewKafkaProducerWrapper(client messaging.KafkaClient, metrics *KafkaMetrics) *KafkaProducerWrapper {
	return &KafkaProducerWrapper{
		client:  client,
		metrics: metrics,
	}
}

// SendMessage sends a message with metrics
func (kpw *KafkaProducerWrapper) SendMessage(ctx context.Context, topic string, key string, value []byte) error {
	start := time.Now()

	err := kpw.client.SendMessage(ctx, topic, key, value)

	latency := time.Since(start)

	if err != nil {
		kpw.metrics.RecordProducerError(topic, "send_failed")
		return err
	}

	kpw.metrics.RecordMessageProduced(topic, len(value), latency)
	return nil
}

// SendMessageWithHeaders sends a message with headers and metrics
func (kpw *KafkaProducerWrapper) SendMessageWithHeaders(ctx context.Context, topic string, key string, value []byte, headers map[string]string) error {
	start := time.Now()

	err := kpw.client.SendMessageWithHeaders(ctx, topic, key, value, headers)

	latency := time.Since(start)

	if err != nil {
		kpw.metrics.RecordProducerError(topic, "send_with_headers_failed")
		return err
	}

	kpw.metrics.RecordMessageProduced(topic, len(value), latency)
	return nil
}

// SendBatchMessages sends batch messages with metrics
func (kpw *KafkaProducerWrapper) SendBatchMessages(ctx context.Context, topic string, messages []messaging.Message) error {
	start := time.Now()

	err := kpw.client.SendBatchMessages(ctx, topic, messages)

	latency := time.Since(start)

	if err != nil {
		kpw.metrics.RecordProducerError(topic, "batch_send_failed")
		return err
	}

	// Calculate total size
	totalSize := 0
	for _, msg := range messages {
		totalSize += len(msg.Value)
	}

	kpw.metrics.RecordMessageProducedBatch(topic, len(messages), totalSize, latency)
	return nil
}

// KafkaConsumerWrapper wraps a Kafka consumer with metrics
type KafkaConsumerWrapper struct {
	client  messaging.KafkaClient
	metrics *KafkaMetrics
}

// NewKafkaConsumerWrapper creates a new Kafka consumer wrapper with metrics
func NewKafkaConsumerWrapper(client messaging.KafkaClient, metrics *KafkaMetrics) *KafkaConsumerWrapper {
	return &KafkaConsumerWrapper{
		client:  client,
		metrics: metrics,
	}
}

// Subscribe subscribes to a topic with metrics
func (kcw *KafkaConsumerWrapper) Subscribe(ctx context.Context, topic string, groupID string, handler messaging.MessageHandler) error {
	// Wrap the handler with metrics
	wrappedHandler := func(ctx context.Context, message messaging.Message) error {
		start := time.Now()

		err := handler(ctx, message)

		latency := time.Since(start)

		if err != nil {
			kcw.metrics.RecordConsumerError(topic, "handler_failed")
			return err
		}

		kcw.metrics.RecordMessageConsumed(topic, len(message.Value), latency)
		return nil
	}

	return kcw.client.Subscribe(ctx, topic, groupID, wrappedHandler)
}

// SubscribeWithOptions subscribes with options and metrics
func (kcw *KafkaConsumerWrapper) SubscribeWithOptions(ctx context.Context, options messaging.ConsumerOptions) error {
	// Wrap the handler with metrics
	wrappedHandler := func(ctx context.Context, message messaging.Message) error {
		start := time.Now()

		err := options.Handler(ctx, message)

		latency := time.Since(start)

		if err != nil {
			kcw.metrics.RecordConsumerError(options.Topic, "handler_failed")
			return err
		}

		kcw.metrics.RecordMessageConsumed(options.Topic, len(message.Value), latency)
		return nil
	}

	// Create new options with wrapped handler
	newOptions := options
	newOptions.Handler = wrappedHandler

	return kcw.client.SubscribeWithOptions(ctx, newOptions)
}
