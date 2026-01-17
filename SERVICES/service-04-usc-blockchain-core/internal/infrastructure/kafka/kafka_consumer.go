package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/usc-platform/shared/config"
	messaging "github.com/usc-platform/shared/kafka-messaging"
	"github.com/usc-platform/shared/logging"
)

// KafkaConsumerManager manages Kafka event consumption for USC Blockchain Core Service
type KafkaConsumerManager struct {
	manager *messaging.KafkaManager
	config  *config.Config
	logger  *logging.Logger
}

// NewKafkaConsumerManager creates a new Kafka consumer manager
func NewKafkaConsumerManager(cfg *config.Config, logger *logging.Logger) (*KafkaConsumerManager, error) {
	manager, err := messaging.NewKafkaManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer manager: %w", err)
	}

	return &KafkaConsumerManager{
		manager: manager,
		config:  cfg,
		logger:  logger,
	}, nil
}

// GetManager returns the Kafka manager
func (k *KafkaConsumerManager) GetManager() *messaging.KafkaManager {
	return k.manager
}

// Close closes the Kafka consumer manager
func (k *KafkaConsumerManager) Close() error {
	if k.manager != nil {
		return k.manager.Close()
	}
	return nil
}

// HealthCheck performs a health check on the Kafka consumer manager
func (k *KafkaConsumerManager) HealthCheck(ctx context.Context) error {
	if k.manager == nil {
		return fmt.Errorf("kafka consumer manager is not initialized")
	}

	// Simple health check - try to list topics
	_, err := k.manager.ListTopics(ctx)
	if err != nil {
		return fmt.Errorf("kafka consumer health check failed: %w", err)
	}

	return nil
}

// SetSchemaRegistry sets the schema registry configuration
func (k *KafkaConsumerManager) SetSchemaRegistry(url string, username string, password string, timeout time.Duration) {
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
func (k *KafkaConsumerManager) GetSchemaRegistry() *messaging.SchemaRegistryConfig {
	if k.manager == nil {
		return nil
	}

	return k.manager.GetSchemaRegistry()
}

// ValidateMessageSchema validates a message against a schema
func (k *KafkaConsumerManager) ValidateMessageSchema(topic string, message []byte) error {
	if k.manager == nil {
		return fmt.Errorf("kafka consumer manager is not initialized")
	}

	return k.manager.ValidateMessageSchema(topic, message)
}
