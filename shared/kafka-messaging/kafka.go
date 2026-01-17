package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/usc-platform/shared/config"
	"github.com/usc-platform/shared/logging"
)

// KafkaClient represents a Kafka client interface
type KafkaClient interface {
	// Producer methods
	SendMessage(ctx context.Context, topic string, key string, value []byte) error
	SendMessageWithHeaders(ctx context.Context, topic string, key string, value []byte, headers map[string]string) error
	SendBatchMessages(ctx context.Context, topic string, messages []Message) error

	// Consumer methods
	Subscribe(ctx context.Context, topic string, groupID string, handler MessageHandler) error
	SubscribeWithOptions(ctx context.Context, options ConsumerOptions) error

	// Admin methods
	CreateTopic(ctx context.Context, topic string, partitions int, replicationFactor int) error
	DeleteTopic(ctx context.Context, topic string) error
	ListTopics(ctx context.Context) ([]string, error)

	// Health and status
	HealthCheck(ctx context.Context) error
	Close() error
}

// Message represents a Kafka message
type Message struct {
	Key       string
	Value     []byte
	Headers   map[string]string
	Topic     string
	Partition int
	Offset    int64
}

// MessageHandler handles incoming messages
type MessageHandler func(ctx context.Context, message Message) error

// ConsumerOptions contains options for consuming messages
type ConsumerOptions struct {
	Topic     string
	GroupID   string
	Handler   MessageHandler
	Partition int
	Offset    int64
	MinBytes  int
	MaxBytes  int
	MaxWait   time.Duration
}

// ProducerOptions contains options for producing messages
type ProducerOptions struct {
	RequiredAcks int
	Compression  string
	BatchSize    int
	BatchTimeout time.Duration
}

// TopicRoutingConfig contains topic routing configuration
type TopicRoutingConfig struct {
	DefaultPartitions int
	ReplicationFactor int
	RetentionHours    int
	CompressionType   string
	CleanupPolicy     string // "delete" or "compact"
}

// SchemaRegistryConfig contains schema registry configuration
type SchemaRegistryConfig struct {
	URL      string
	Username string
	Password string
	Timeout  time.Duration
}

// DeadLetterQueueConfig contains DLQ configuration
type DeadLetterQueueConfig struct {
	Enabled       bool
	TopicSuffix   string
	MaxRetries    int
	RetryDelay    time.Duration
	RetryBackoff  time.Duration
	MaxRetryDelay time.Duration
}

// EventSourcingConfig contains event sourcing configuration
type EventSourcingConfig struct {
	Enabled           bool
	EventStoreTopic   string
	SnapshotTopic     string
	SnapshotInterval  time.Duration
	MaxEventsPerBatch int
}

// KafkaManager manages Kafka connections and operations
type KafkaManager struct {
	config *config.Config
	logger *logging.Logger

	// Kafka connections
	producer *kafka.Writer
	consumer *kafka.Reader
	conn     *kafka.Conn

	// Configuration
	brokers         []string
	options         ProducerOptions
	topicRouting    TopicRoutingConfig
	schemaRegistry  *SchemaRegistryConfig
	deadLetterQueue DeadLetterQueueConfig
	eventSourcing   EventSourcingConfig

	// Advanced features
	topicCache        map[string]*TopicInfo
	partitionStrategy PartitionStrategy
	eventStore        *EventStore
	dlqHandler        *DeadLetterQueueHandler

	// State management
	connected bool
	mu        sync.RWMutex
}

// TopicInfo contains topic metadata
type TopicInfo struct {
	Name              string
	Partitions        int
	ReplicationFactor int
	Config            map[string]string
	CreatedAt         time.Time
}

// PartitionStrategy defines how messages are partitioned
type PartitionStrategy interface {
	GetPartition(key string, topic string, partitions int) int
}

// RoundRobinPartitionStrategy implements round-robin partitioning
type RoundRobinPartitionStrategy struct {
	counter int64
}

// HashPartitionStrategy implements hash-based partitioning
type HashPartitionStrategy struct{}

// EventStore manages event sourcing
type EventStore struct {
	topic    string
	producer *kafka.Writer
	logger   *logging.Logger
}

// DeadLetterQueueHandler manages failed message handling
type DeadLetterQueueHandler struct {
	config   DeadLetterQueueConfig
	producer *kafka.Writer
	logger   *logging.Logger
}

// NewKafkaManager creates a new Kafka manager
func NewKafkaManager(cfg *config.Config) (*KafkaManager, error) {
	manager := &KafkaManager{
		config:  cfg,
		logger:  logging.NewLogger("kafka", cfg.Log),
		brokers: cfg.GetKafkaBrokers(),
		options: ProducerOptions{
			RequiredAcks: 1,
			Compression:  "snappy",
			BatchSize:    100,
			BatchTimeout: 10 * time.Millisecond,
		},
		topicRouting: TopicRoutingConfig{
			DefaultPartitions: 3,
			ReplicationFactor: 1,
			RetentionHours:    168, // 7 days
			CompressionType:   "snappy",
			CleanupPolicy:     "delete",
		},
		deadLetterQueue: DeadLetterQueueConfig{
			Enabled:       true,
			TopicSuffix:   "-dlq",
			MaxRetries:    3,
			RetryDelay:    5 * time.Second,
			RetryBackoff:  2 * time.Second,
			MaxRetryDelay: 60 * time.Second,
		},
		eventSourcing: EventSourcingConfig{
			Enabled:           false,
			EventStoreTopic:   "event-store",
			SnapshotTopic:     "snapshots",
			SnapshotInterval:  1 * time.Hour,
			MaxEventsPerBatch: 1000,
		},
		topicCache:        make(map[string]*TopicInfo),
		partitionStrategy: &HashPartitionStrategy{},
	}

	// Log configured brokers
	manager.logger.Info("Initializing Kafka manager with brokers",
		logging.String("brokers", fmt.Sprintf("%v", manager.brokers)))

	// Initialize connections
	if err := manager.initializeConnections(); err != nil {
		return nil, fmt.Errorf("failed to initialize kafka connections: %w", err)
	}

	// Initialize advanced features
	if err := manager.initializeAdvancedFeatures(); err != nil {
		return nil, fmt.Errorf("failed to initialize advanced features: %w", err)
	}

	manager.connected = true
	return manager, nil
}

// GetPartition implements PartitionStrategy for RoundRobinPartitionStrategy
func (r *RoundRobinPartitionStrategy) GetPartition(key string, topic string, partitions int) int {
	if partitions <= 0 {
		return 0
	}
	partition := atomic.AddInt64(&r.counter, 1) % int64(partitions)
	return int(partition)
}

// GetPartition implements PartitionStrategy for HashPartitionStrategy
func (h *HashPartitionStrategy) GetPartition(key string, topic string, partitions int) int {
	if partitions <= 0 {
		return 0
	}
	hash := fnv.New32a()
	hash.Write([]byte(key))
	return int(hash.Sum32()) % partitions
}

// initializeConnections initializes Kafka connections
func (m *KafkaManager) initializeConnections() error {
	// Initialize producer
	if err := m.initializeProducer(); err != nil {
		return fmt.Errorf("failed to initialize producer: %w", err)
	}

	// Initialize connection for admin operations
	if err := m.initializeConnection(); err != nil {
		return fmt.Errorf("failed to initialize connection: %w", err)
	}

	return nil
}

// initializeProducer initializes the Kafka producer
func (m *KafkaManager) initializeProducer() error {
	m.producer = &kafka.Writer{
		Addr:         kafka.TCP(m.brokers...),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequiredAcks(m.options.RequiredAcks),
		BatchSize:    m.options.BatchSize,
		BatchTimeout: m.options.BatchTimeout,
		Compression:  kafka.Snappy,
		ErrorLogger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
			m.logger.Error(msg, logging.String("component", "kafka-producer"))
		}),
		Logger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
			m.logger.Info(msg, logging.String("component", "kafka-producer"))
		}),
	}

	return nil
}

// initializeConnection initializes the Kafka connection for admin operations
func (m *KafkaManager) initializeConnection() error {
	conn, err := kafka.Dial("tcp", m.brokers[0])
	if err != nil {
		return fmt.Errorf("failed to dial kafka: %w", err)
	}

	m.conn = conn
	return nil
}

// initializeAdvancedFeatures initializes advanced Kafka features
func (m *KafkaManager) initializeAdvancedFeatures() error {
	// Initialize DLQ handler if enabled
	if m.deadLetterQueue.Enabled {
		m.dlqHandler = &DeadLetterQueueHandler{
			config:   m.deadLetterQueue,
			producer: m.producer,
			logger:   m.logger,
		}
	}

	// Initialize event store if enabled
	if m.eventSourcing.Enabled {
		m.eventStore = &EventStore{
			topic:    m.eventSourcing.EventStoreTopic,
			producer: m.producer,
			logger:   m.logger,
		}
	}

	return nil
}

// SendMessage sends a single message to a topic
func (m *KafkaManager) SendMessage(ctx context.Context, topic string, key string, value []byte) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.connected {
		return fmt.Errorf("kafka manager is not connected")
	}

	message := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
		Time:  time.Now(),
	}

	if err := m.producer.WriteMessages(ctx, message); err != nil {
		m.logger.Error("Producer WriteMessages failed",
			logging.Error(err),
			logging.String("topic", topic),
			logging.String("brokers", fmt.Sprintf("%v", m.brokers)))
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// SendMessageWithHeaders sends a message with headers
func (m *KafkaManager) SendMessageWithHeaders(ctx context.Context, topic string, key string, value []byte, headers map[string]string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.connected {
		return fmt.Errorf("kafka manager is not connected")
	}

	// Convert headers to Kafka headers
	kafkaHeaders := make([]kafka.Header, 0, len(headers))
	for k, v := range headers {
		kafkaHeaders = append(kafkaHeaders, kafka.Header{
			Key:   k,
			Value: []byte(v),
		})
	}

	message := kafka.Message{
		Topic:   topic,
		Key:     []byte(key),
		Value:   value,
		Headers: kafkaHeaders,
		Time:    time.Now(),
	}

	if err := m.producer.WriteMessages(ctx, message); err != nil {
		m.logger.Error("Producer WriteMessages failed",
			logging.Error(err),
			logging.String("topic", topic),
			logging.String("brokers", fmt.Sprintf("%v", m.brokers)))
		return fmt.Errorf("failed to send message with headers: %w", err)
	}

	return nil
}

// SendBatchMessages sends multiple messages in a batch
func (m *KafkaManager) SendBatchMessages(ctx context.Context, topic string, messages []Message) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.connected {
		return fmt.Errorf("kafka manager is not connected")
	}

	kafkaMessages := make([]kafka.Message, len(messages))
	for i, msg := range messages {
		// Convert headers to Kafka headers
		kafkaHeaders := make([]kafka.Header, 0, len(msg.Headers))
		for k, v := range msg.Headers {
			kafkaHeaders = append(kafkaHeaders, kafka.Header{
				Key:   k,
				Value: []byte(v),
			})
		}

		kafkaMessages[i] = kafka.Message{
			Topic:   topic,
			Key:     []byte(msg.Key),
			Value:   msg.Value,
			Headers: kafkaHeaders,
			Time:    time.Now(),
		}
	}

	if err := m.producer.WriteMessages(ctx, kafkaMessages...); err != nil {
		m.logger.Error("Producer WriteMessages failed",
			logging.Error(err),
			logging.String("topic", topic),
			logging.String("brokers", fmt.Sprintf("%v", m.brokers)))
		return fmt.Errorf("failed to send batch messages: %w", err)
	}

	return nil
}

// Subscribe subscribes to a topic and processes messages
func (m *KafkaManager) Subscribe(ctx context.Context, topic string, groupID string, handler MessageHandler) error {
	options := ConsumerOptions{
		Topic:    topic,
		GroupID:  groupID,
		Handler:  handler,
		MinBytes: 1,    // deliver small messages promptly
		MaxBytes: 10e6, // 10MB
		MaxWait:  200 * time.Millisecond,
	}

	return m.SubscribeWithOptions(ctx, options)
}

// SubscribeWithOptions subscribes with custom options
func (m *KafkaManager) SubscribeWithOptions(ctx context.Context, options ConsumerOptions) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.connected {
		return fmt.Errorf("kafka manager is not connected")
	}

	// Create reader: if explicit partition/offset provided, use direct partition reader; otherwise use group consumer
	var reader *kafka.Reader
	if options.Partition > 0 || options.Offset >= 0 {
		startOffset := kafka.FirstOffset
		if options.Offset > 0 {
			startOffset = options.Offset
		}
		reader = kafka.NewReader(kafka.ReaderConfig{
			Brokers:           m.brokers,
			Topic:             options.Topic,
			Partition:         options.Partition,
			MinBytes:          options.MinBytes,
			MaxBytes:          options.MaxBytes,
			MaxWait:           options.MaxWait,
			HeartbeatInterval: 3 * time.Second,
			SessionTimeout:    10 * time.Second,
			RebalanceTimeout:  15 * time.Second,
			JoinGroupBackoff:  1 * time.Second,
			ReadLagInterval:   0,
			CommitInterval:    1 * time.Second,
			ErrorLogger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
				m.logger.Error(msg, logging.String("component", "kafka-consumer"))
			}),
			Logger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
				m.logger.Info(msg, logging.String("component", "kafka-consumer"))
			}),
			StartOffset: startOffset,
		})
	} else {
		reader = kafka.NewReader(kafka.ReaderConfig{
			Brokers:               m.brokers,
			Topic:                 options.Topic,
			GroupID:               options.GroupID,
			MinBytes:              options.MinBytes,
			MaxBytes:              options.MaxBytes,
			MaxWait:               options.MaxWait,
			WatchPartitionChanges: true,
			GroupBalancers:        []kafka.GroupBalancer{kafka.RangeGroupBalancer{}},
			HeartbeatInterval:     3 * time.Second,
			SessionTimeout:        10 * time.Second,
			RebalanceTimeout:      15 * time.Second,
			JoinGroupBackoff:      1 * time.Second,
			ReadLagInterval:       0,
			CommitInterval:        1 * time.Second,
			ErrorLogger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
				m.logger.Error(msg, logging.String("component", "kafka-consumer"))
			}),
			Logger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
				m.logger.Info(msg, logging.String("component", "kafka-consumer"))
			}),
			StartOffset: kafka.FirstOffset,
		})
	}

	// For direct partition reader, SetOffset already applied via StartOffset

	// Start consuming messages
	go func() {
		defer reader.Close()

		for {
			select {
			case <-ctx.Done():
				m.logger.Info("Stopping Kafka consumer", logging.String("topic", options.Topic), logging.String("group", options.GroupID))
				return
			default:
				msg, err := reader.ReadMessage(ctx)
				if err != nil {
					m.logger.Error("Failed to read message", logging.Error(err))
					continue
				}

				// Convert Kafka message to our Message type
				message := Message{
					Key:       string(msg.Key),
					Value:     msg.Value,
					Topic:     msg.Topic,
					Partition: msg.Partition,
					Offset:    msg.Offset,
				}

				// Convert headers
				if len(msg.Headers) > 0 {
					message.Headers = make(map[string]string)
					for _, header := range msg.Headers {
						message.Headers[header.Key] = string(header.Value)
					}
				}

				// Handle the message
				if err := options.Handler(ctx, message); err != nil {
					m.logger.Error("Failed to handle message", logging.Error(err), logging.String("topic", msg.Topic), logging.String("key", string(msg.Key)))
				}
			}
		}
	}()

	return nil
}

// CreateTopic creates a new topic
func (m *KafkaManager) CreateTopic(ctx context.Context, topic string, partitions int, replicationFactor int) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.connected {
		return fmt.Errorf("kafka manager is not connected")
	}

	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
	}

	if err := m.conn.CreateTopics(topicConfig); err != nil {
		return fmt.Errorf("failed to create topic %s: %w", topic, err)
	}

	return nil
}

// DeleteTopic deletes a topic
func (m *KafkaManager) DeleteTopic(ctx context.Context, topic string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.connected {
		return fmt.Errorf("kafka manager is not connected")
	}

	if err := m.conn.DeleteTopics(topic); err != nil {
		return fmt.Errorf("failed to delete topic %s: %w", topic, err)
	}

	return nil
}

// ListTopics lists all topics
func (m *KafkaManager) ListTopics(ctx context.Context) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.connected {
		return nil, fmt.Errorf("kafka manager is not connected")
	}

	partitions, err := m.conn.ReadPartitions()
	if err != nil {
		return nil, fmt.Errorf("failed to read partitions: %w", err)
	}

	topicMap := make(map[string]bool)
	for _, partition := range partitions {
		topicMap[partition.Topic] = true
	}

	topics := make([]string, 0, len(topicMap))
	for topic := range topicMap {
		topics = append(topics, topic)
	}

	return topics, nil
}

// HealthCheck performs a health check on Kafka
func (m *KafkaManager) HealthCheck(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.connected {
		return fmt.Errorf("kafka manager is not connected")
	}

	// Try to list topics as a health check
	_, err := m.ListTopics(ctx)
	if err != nil {
		return fmt.Errorf("kafka health check failed: %w", err)
	}

	return nil
}

// Close closes all Kafka connections
func (m *KafkaManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errors []error

	if m.producer != nil {
		if err := m.producer.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close producer: %w", err))
		}
	}

	if m.consumer != nil {
		if err := m.consumer.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close consumer: %w", err))
		}
	}

	if m.conn != nil {
		if err := m.conn.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close connection: %w", err))
		}
	}

	m.connected = false

	if len(errors) > 0 {
		return fmt.Errorf("errors closing kafka connections: %v", errors)
	}

	return nil
}

// IsConnected returns true if the manager is connected
func (m *KafkaManager) IsConnected() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.connected
}

// SendJSONMessage sends a JSON message
func (m *KafkaManager) SendJSONMessage(ctx context.Context, topic string, key string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	return m.SendMessage(ctx, topic, key, jsonData)
}

// SendJSONMessageWithHeaders sends a JSON message with headers
func (m *KafkaManager) SendJSONMessageWithHeaders(ctx context.Context, topic string, key string, data interface{}, headers map[string]string) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	return m.SendMessageWithHeaders(ctx, topic, key, jsonData, headers)
}

// SetProducerOptions sets producer options
func (m *KafkaManager) SetProducerOptions(options ProducerOptions) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.options = options
}

// GetBrokers returns the list of brokers
func (m *KafkaManager) GetBrokers() []string {
	return m.brokers
}

// AddBroker adds a broker to the list
func (m *KafkaManager) AddBroker(broker string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if broker already exists
	for _, b := range m.brokers {
		if b == broker {
			return
		}
	}

	m.brokers = append(m.brokers, broker)
}

// ===== ADVANCED FEATURES =====

// CreateTopicWithRouting creates a topic with advanced routing configuration
func (m *KafkaManager) CreateTopicWithRouting(ctx context.Context, topic string, routing TopicRoutingConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.connected {
		return fmt.Errorf("kafka manager is not connected")
	}

	// Create topic with custom configuration
	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     routing.DefaultPartitions,
		ReplicationFactor: routing.ReplicationFactor,
		ConfigEntries: []kafka.ConfigEntry{
			{ConfigName: "retention.ms", ConfigValue: fmt.Sprintf("%d", routing.RetentionHours*3600000)},
			{ConfigName: "compression.type", ConfigValue: routing.CompressionType},
			{ConfigName: "cleanup.policy", ConfigValue: routing.CleanupPolicy},
			{ConfigName: "min.insync.replicas", ConfigValue: "1"},
			{ConfigName: "unclean.leader.election.enable", ConfigValue: "false"},
		},
	}

	if err := m.conn.CreateTopics(topicConfig); err != nil {
		return fmt.Errorf("failed to create topic %s with routing: %w", topic, err)
	}

	// Cache topic info
	configMap := make(map[string]string)
	for _, entry := range topicConfig.ConfigEntries {
		configMap[entry.ConfigName] = entry.ConfigValue
	}

	m.topicCache[topic] = &TopicInfo{
		Name:              topic,
		Partitions:        routing.DefaultPartitions,
		ReplicationFactor: routing.ReplicationFactor,
		Config:            configMap,
		CreatedAt:         time.Now(),
	}

	m.logger.Info("Created topic with advanced routing",
		logging.String("topic", topic),
		logging.Int("partitions", routing.DefaultPartitions),
		logging.Int("replication_factor", routing.ReplicationFactor))

	return nil
}

// SendMessageWithPartitioning sends a message with custom partitioning
func (m *KafkaManager) SendMessageWithPartitioning(ctx context.Context, topic string, key string, value []byte, partition int) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.connected {
		return fmt.Errorf("kafka manager is not connected")
	}

	message := kafka.Message{
		Topic:     topic,
		Key:       []byte(key),
		Value:     value,
		Time:      time.Now(),
		Partition: partition,
	}

	if err := m.producer.WriteMessages(ctx, message); err != nil {
		return fmt.Errorf("failed to send message with partitioning: %w", err)
	}

	return nil
}

// SendMessageWithRetry sends a message with retry logic and DLQ support
func (m *KafkaManager) SendMessageWithRetry(ctx context.Context, topic string, key string, value []byte, maxRetries int) error {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		err := m.SendMessage(ctx, topic, key, value)
		if err == nil {
			return nil
		}

		lastErr = err

		if attempt < maxRetries {
			// Exponential backoff
			delay := time.Duration(attempt+1) * time.Second
			time.Sleep(delay)

			m.logger.Warn("Message send failed, retrying",
				logging.String("topic", topic),
				logging.Int("attempt", attempt+1),
				logging.Int("max_retries", maxRetries),
				logging.Error(err))
		}
	}

	// Send to DLQ if enabled and all retries failed
	if m.deadLetterQueue.Enabled && m.dlqHandler != nil {
		dlqTopic := topic + m.deadLetterQueue.TopicSuffix
		dlqErr := m.dlqHandler.SendToDLQ(ctx, dlqTopic, key, value, lastErr)
		if dlqErr != nil {
			m.logger.Error("Failed to send message to DLQ",
				logging.String("topic", topic),
				logging.String("dlq_topic", dlqTopic),
				logging.Error(dlqErr))
		}
	}

	return fmt.Errorf("failed to send message after %d retries: %w", maxRetries, lastErr)
}

// GetTopicInfo returns cached topic information
func (m *KafkaManager) GetTopicInfo(topic string) (*TopicInfo, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	info, exists := m.topicCache[topic]
	return info, exists
}

// SetPartitionStrategy sets the partition strategy
func (m *KafkaManager) SetPartitionStrategy(strategy PartitionStrategy) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.partitionStrategy = strategy
}

// EnableEventSourcing enables event sourcing for the manager
func (m *KafkaManager) EnableEventSourcing(config EventSourcingConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.eventSourcing = config
	m.eventSourcing.Enabled = true

	// Initialize event store
	m.eventStore = &EventStore{
		topic:    config.EventStoreTopic,
		producer: m.producer,
		logger:   m.logger,
	}

	// Create event store topic if it doesn't exist
	ctx := context.Background()
	if err := m.CreateTopic(ctx, config.EventStoreTopic, 3, 1); err != nil {
		return fmt.Errorf("failed to create event store topic: %w", err)
	}

	m.logger.Info("Event sourcing enabled",
		logging.String("event_store_topic", config.EventStoreTopic),
		logging.String("snapshot_topic", config.SnapshotTopic))

	return nil
}

// StoreEvent stores an event in the event store
func (m *KafkaManager) StoreEvent(ctx context.Context, eventType string, aggregateID string, eventData interface{}) error {
	if m.eventStore == nil {
		return fmt.Errorf("event sourcing is not enabled")
	}

	event := map[string]interface{}{
		"event_type":   eventType,
		"aggregate_id": aggregateID,
		"event_data":   eventData,
		"timestamp":    time.Now().Unix(),
		"version":      1,
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	return m.SendMessage(ctx, m.eventStore.topic, aggregateID, eventBytes)
}

// SetSchemaRegistry sets the schema registry configuration
func (m *KafkaManager) SetSchemaRegistry(config SchemaRegistryConfig) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.schemaRegistry = &config
	m.logger.Info("Schema registry configured",
		logging.String("url", config.URL),
		logging.Duration("timeout", config.Timeout))
}

// GetSchemaRegistry returns the schema registry configuration
func (m *KafkaManager) GetSchemaRegistry() *SchemaRegistryConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.schemaRegistry
}

// ValidateMessageSchema validates a message against a schema (placeholder for future implementation)
func (m *KafkaManager) ValidateMessageSchema(topic string, message []byte) error {
	if m.schemaRegistry == nil {
		// No schema registry configured, skip validation
		return nil
	}

	// TODO: Implement actual schema validation when schema registry is integrated
	// This is a placeholder for future schema validation functionality
	m.logger.Debug("Schema validation skipped - not yet implemented",
		logging.String("topic", topic),
		logging.Int("message_size", len(message)))

	return nil
}

// ===== DEAD LETTER QUEUE HANDLER =====

// SendToDLQ sends a failed message to the dead letter queue
func (dlq *DeadLetterQueueHandler) SendToDLQ(ctx context.Context, dlqTopic string, key string, value []byte, originalError error) error {
	dlqMessage := map[string]interface{}{
		"original_topic": "",
		"original_key":   key,
		"original_value": string(value),
		"error":          originalError.Error(),
		"timestamp":      time.Now().Unix(),
		"retry_count":    0,
	}

	dlqBytes, err := json.Marshal(dlqMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal DLQ message: %w", err)
	}

	message := kafka.Message{
		Topic: dlqTopic,
		Key:   []byte(key),
		Value: dlqBytes,
		Time:  time.Now(),
	}

	if err := dlq.producer.WriteMessages(ctx, message); err != nil {
		return fmt.Errorf("failed to send message to DLQ: %w", err)
	}

	dlq.logger.Info("Message sent to DLQ",
		logging.String("dlq_topic", dlqTopic),
		logging.String("key", key),
		logging.Error(originalError))

	return nil
}
