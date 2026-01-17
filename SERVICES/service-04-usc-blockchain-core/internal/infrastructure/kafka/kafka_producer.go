package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/usc-platform/shared/config"
	messaging "github.com/usc-platform/shared/kafka-messaging"
	"github.com/usc-platform/shared/logging"
)

// KafkaProducerManager manages Kafka event production for USC Blockchain Core Service
type KafkaProducerManager struct {
	manager *messaging.KafkaManager
	config  *config.Config
	logger  *logging.Logger
}

// NewKafkaProducerManager creates a new Kafka producer manager
func NewKafkaProducerManager(cfg *config.Config, logger *logging.Logger) (*KafkaProducerManager, error) {
	manager, err := messaging.NewKafkaManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer manager: %w", err)
	}

	return &KafkaProducerManager{
		manager: manager,
		config:  cfg,
		logger:  logger,
	}, nil
}

// GetManager returns the Kafka manager
func (k *KafkaProducerManager) GetManager() *messaging.KafkaManager {
	return k.manager
}

// Close closes the Kafka producer manager
func (k *KafkaProducerManager) Close() error {
	if k.manager != nil {
		return k.manager.Close()
	}
	return nil
}

// HealthCheck performs a health check on the Kafka producer manager
func (k *KafkaProducerManager) HealthCheck(ctx context.Context) error {
	if k.manager == nil {
		return fmt.Errorf("kafka producer manager is not initialized")
	}

	// Simple health check - try to list topics
	_, err := k.manager.ListTopics(ctx)
	if err != nil {
		return fmt.Errorf("kafka producer health check failed: %w", err)
	}

	return nil
}

// SendMessageWithRetry sends a message with retry logic and DLQ support
func (k *KafkaProducerManager) SendMessageWithRetry(ctx context.Context, topic string, key string, value []byte, maxRetries int) error {
	if k.manager == nil {
		return fmt.Errorf("kafka producer manager is not initialized")
	}

	return k.manager.SendMessageWithRetry(ctx, topic, key, value, maxRetries)
}

// SetSchemaRegistry sets the schema registry configuration
func (k *KafkaProducerManager) SetSchemaRegistry(url string, username string, password string, timeout time.Duration) {
	if k.manager == nil {
		return
	}

	config := messaging.SchemaRegistryConfig{
		URL:      url,
		Username: username,
		Password: password,
		Timeout:  timeout,
	}

	k.manager.SetSchemaRegistry(config)
}

// GetSchemaRegistry returns the schema registry configuration
func (k *KafkaProducerManager) GetSchemaRegistry() *messaging.SchemaRegistryConfig {
	if k.manager == nil {
		return nil
	}

	return k.manager.GetSchemaRegistry()
}

// ValidateMessageSchema validates a message against a schema
func (k *KafkaProducerManager) ValidateMessageSchema(topic string, message []byte) error {
	if k.manager == nil {
		return fmt.Errorf("kafka producer manager is not initialized")
	}

	return k.manager.ValidateMessageSchema(topic, message)
}
