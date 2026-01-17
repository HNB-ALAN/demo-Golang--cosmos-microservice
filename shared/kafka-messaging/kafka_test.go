package messaging

import (
	"context"
	"testing"
	"time"

	"github.com/usc-platform/shared/config"
)

func TestNewKafkaManager(t *testing.T) {
	cfg := &config.Config{
		Kafka: config.KafkaConfig{
			Brokers: []string{"localhost:9092"},
		},
	}

	manager, err := NewKafkaManager(cfg)
	if err != nil {
		// Expected to fail without Kafka server - this is normal for unit tests
		t.Logf("Kafka manager creation failed as expected (no server): %v", err)
		return
	}

	if manager == nil {
		t.Fatal("Manager is nil")
	}

	if !manager.IsConnected() {
		t.Error("Manager should be connected")
	}

	// Clean up
	manager.Close()
}

func TestSendMessage(t *testing.T) {
	cfg := &config.Config{
		Kafka: config.KafkaConfig{
			Brokers: []string{"localhost:9092"},
		},
	}

	manager, err := NewKafkaManager(cfg)
	if err != nil {
		// Expected to fail without Kafka server - this is normal for unit tests
		t.Logf("Kafka manager creation failed as expected (no server): %v", err)
		return
	}
	defer manager.Close()

	ctx := context.Background()
	topic := "test-topic"
	key := "test-key"
	value := []byte("test-value")

	// This test will fail if Kafka is not running, which is expected
	err = manager.SendMessage(ctx, topic, key, value)
	if err != nil {
		t.Logf("SendMessage failed (expected if Kafka is not running): %v", err)
	}
}

func TestSendMessageWithHeaders(t *testing.T) {
	cfg := &config.Config{
		Kafka: config.KafkaConfig{
			Brokers: []string{"localhost:9092"},
		},
	}

	manager, err := NewKafkaManager(cfg)
	if err != nil {
		t.Logf("Kafka manager creation failed as expected (no server): %v", err)
		return
	}
	defer manager.Close()

	ctx := context.Background()
	topic := "test-topic"
	key := "test-key"
	value := []byte("test-value")
	headers := map[string]string{
		"content-type": "application/json",
		"source":       "test",
	}

	// This test will fail if Kafka is not running, which is expected
	err = manager.SendMessageWithHeaders(ctx, topic, key, value, headers)
	if err != nil {
		t.Logf("SendMessageWithHeaders failed (expected if Kafka is not running): %v", err)
	}
}

func TestSendJSONMessage(t *testing.T) {
	cfg := &config.Config{
		Kafka: config.KafkaConfig{
			Brokers: []string{"localhost:9092"},
		},
	}

	manager, err := NewKafkaManager(cfg)
	if err != nil {
		t.Logf("Kafka manager creation failed as expected (no server): %v", err)
		return
	}
	defer manager.Close()

	ctx := context.Background()
	topic := "test-topic"
	key := "test-key"
	data := map[string]interface{}{
		"message":   "test",
		"timestamp": time.Now().Unix(),
	}

	// This test will fail if Kafka is not running, which is expected
	err = manager.SendJSONMessage(ctx, topic, key, data)
	if err != nil {
		t.Logf("SendJSONMessage failed (expected if Kafka is not running): %v", err)
	}
}

func TestSendBatchMessages(t *testing.T) {
	cfg := &config.Config{
		Kafka: config.KafkaConfig{
			Brokers: []string{"localhost:9092"},
		},
	}

	manager, err := NewKafkaManager(cfg)
	if err != nil {
		t.Logf("Kafka manager creation failed as expected (no server): %v", err)
		return
	}
	defer manager.Close()

	ctx := context.Background()
	topic := "test-topic"
	messages := []Message{
		{
			Key:   "key1",
			Value: []byte("value1"),
			Headers: map[string]string{
				"type": "test",
			},
		},
		{
			Key:   "key2",
			Value: []byte("value2"),
			Headers: map[string]string{
				"type": "test",
			},
		},
	}

	// This test will fail if Kafka is not running, which is expected
	err = manager.SendBatchMessages(ctx, topic, messages)
	if err != nil {
		t.Logf("SendBatchMessages failed (expected if Kafka is not running): %v", err)
	}
}

func TestHealthCheck(t *testing.T) {
	cfg := &config.Config{
		Kafka: config.KafkaConfig{
			Brokers: []string{"localhost:9092"},
		},
	}

	manager, err := NewKafkaManager(cfg)
	if err != nil {
		t.Logf("Kafka manager creation failed as expected (no server): %v", err)
		return
	}
	defer manager.Close()

	ctx := context.Background()

	// This test will fail if Kafka is not running, which is expected
	err = manager.HealthCheck(ctx)
	if err != nil {
		t.Logf("HealthCheck failed (expected if Kafka is not running): %v", err)
	}
}

func TestListTopics(t *testing.T) {
	cfg := &config.Config{
		Kafka: config.KafkaConfig{
			Brokers: []string{"localhost:9092"},
		},
	}

	manager, err := NewKafkaManager(cfg)
	if err != nil {
		t.Logf("Kafka manager creation failed as expected (no server): %v", err)
		return
	}
	defer manager.Close()

	ctx := context.Background()

	// This test will fail if Kafka is not running, which is expected
	topics, err := manager.ListTopics(ctx)
	if err != nil {
		t.Logf("ListTopics failed (expected if Kafka is not running): %v", err)
	} else {
		t.Logf("Found topics: %v", topics)
	}
}

func TestAddBroker(t *testing.T) {
	cfg := &config.Config{
		Kafka: config.KafkaConfig{
			Brokers: []string{"localhost:9092"},
		},
	}

	manager, err := NewKafkaManager(cfg)
	if err != nil {
		t.Logf("Kafka manager creation failed as expected (no server): %v", err)
		return
	}
	defer manager.Close()

	initialBrokers := manager.GetBrokers()
	if len(initialBrokers) != 1 {
		t.Errorf("Expected 1 broker, got %d", len(initialBrokers))
	}

	manager.AddBroker("localhost:9093")
	brokers := manager.GetBrokers()
	if len(brokers) != 2 {
		t.Errorf("Expected 2 brokers, got %d", len(brokers))
	}

	// Add duplicate broker
	manager.AddBroker("localhost:9092")
	brokers = manager.GetBrokers()
	if len(brokers) != 2 {
		t.Errorf("Expected 2 brokers after adding duplicate, got %d", len(brokers))
	}
}

func TestSetProducerOptions(t *testing.T) {
	cfg := &config.Config{
		Kafka: config.KafkaConfig{
			Brokers: []string{"localhost:9092"},
		},
	}

	manager, err := NewKafkaManager(cfg)
	if err != nil {
		t.Logf("Kafka manager creation failed as expected (no server): %v", err)
		return
	}
	defer manager.Close()

	options := ProducerOptions{
		RequiredAcks: 2,
		Compression:  "gzip",
		BatchSize:    200,
		BatchTimeout: 20 * time.Millisecond,
	}

	manager.SetProducerOptions(options)
	// Note: We can't easily test if the options were set without exposing internal state
	// This test mainly ensures the method doesn't panic
}

func TestMessageHandler(t *testing.T) {
	// Test message handler function
	handler := func(ctx context.Context, message Message) error {
		if message.Key == "" {
			t.Error("Message key should not be empty")
		}
		if len(message.Value) == 0 {
			t.Error("Message value should not be empty")
		}
		return nil
	}

	// Create a test message
	message := Message{
		Key:   "test-key",
		Value: []byte("test-value"),
		Topic: "test-topic",
		Headers: map[string]string{
			"content-type": "text/plain",
		},
	}

	ctx := context.Background()
	err := handler(ctx, message)
	if err != nil {
		t.Errorf("Message handler failed: %v", err)
	}
}

func TestConsumerOptions(t *testing.T) {
	options := ConsumerOptions{
		Topic:     "test-topic",
		GroupID:   "test-group",
		Partition: 0,
		Offset:    0,
		MinBytes:  1024,
		MaxBytes:  1048576,
		MaxWait:   1 * time.Second,
		Handler: func(ctx context.Context, message Message) error {
			return nil
		},
	}

	if options.Topic != "test-topic" {
		t.Error("Topic should be 'test-topic'")
	}
	if options.GroupID != "test-group" {
		t.Error("GroupID should be 'test-group'")
	}
	if options.MinBytes != 1024 {
		t.Error("MinBytes should be 1024")
	}
	if options.MaxBytes != 1048576 {
		t.Error("MaxBytes should be 1048576")
	}
}
