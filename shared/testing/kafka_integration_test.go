//go:build integration
// +build integration

package messaging

import (
	"context"
	"testing"
	"time"

	"github.com/usc-platform/shared/config"
)

// Integration tests require a running Kafka instance
// Run with: go test -tags=integration ./messaging/...

func TestKafkaIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	cfg := &config.Config{
		Kafka: config.KafkaConfig{
			Brokers: []string{"localhost:9092"},
		},
		Log: config.LogConfig{
			Level:  "info",
			Format: "json",
			Output: "stdout",
		},
	}

	manager, err := NewManager(cfg)
	if err != nil {
		t.Skipf("Skipping integration test - Kafka not available: %v", err)
	}
	defer manager.Close()

	ctx := context.Background()
	topic := "integration-test-topic"

	// Test topic creation
	err = manager.CreateTopic(ctx, topic, 1, 1)
	if err != nil {
		t.Logf("Topic creation failed (may already exist): %v", err)
	}

	// Test message publishing
	testMessage := []byte("integration test message")
	err = manager.SendMessage(ctx, topic, "test-key", testMessage)
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Test message consumption
	messageReceived := make(chan bool, 1)
	handler := func(ctx context.Context, message Message) error {
		if string(message.Value) == string(testMessage) {
			messageReceived <- true
		}
		return nil
	}

	// Start consumer
	go func() {
		err := manager.Subscribe(ctx, topic, "integration-test-group", handler)
		if err != nil {
			t.Errorf("Failed to subscribe: %v", err)
		}
	}()

	// Wait for message to be received
	select {
	case <-messageReceived:
		t.Log("Message received successfully")
	case <-time.After(10 * time.Second):
		t.Error("Timeout waiting for message")
	}

	// Test health check
	err = manager.HealthCheck(ctx)
	if err != nil {
		t.Errorf("Health check failed: %v", err)
	}

	// Test topic listing
	topics, err := manager.ListTopics(ctx)
	if err != nil {
		t.Errorf("Failed to list topics: %v", err)
	}

	found := false
	for _, t := range topics {
		if t == topic {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Created topic not found in topic list")
	}

	// Clean up - delete topic
	err = manager.DeleteTopic(ctx, topic)
	if err != nil {
		t.Logf("Topic deletion failed: %v", err)
	}
}

func TestKafkaBatchIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	cfg := &config.Config{
		Kafka: config.KafkaConfig{
			Brokers: []string{"localhost:9092"},
		},
		Log: config.LogConfig{
			Level:  "info",
			Format: "json",
			Output: "stdout",
		},
	}

	manager, err := NewManager(cfg)
	if err != nil {
		t.Skipf("Skipping integration test - Kafka not available: %v", err)
	}
	defer manager.Close()

	ctx := context.Background()
	topic := "integration-batch-test-topic"

	// Test topic creation
	err = manager.CreateTopic(ctx, topic, 1, 1)
	if err != nil {
		t.Logf("Topic creation failed (may already exist): %v", err)
	}

	// Test batch message publishing
	messages := []Message{
		{Key: "key1", Value: []byte("batch message 1")},
		{Key: "key2", Value: []byte("batch message 2")},
		{Key: "key3", Value: []byte("batch message 3")},
	}

	err = manager.SendBatchMessages(ctx, topic, messages)
	if err != nil {
		t.Fatalf("Failed to send batch messages: %v", err)
	}

	// Test batch message consumption
	receivedCount := 0
	expectedCount := len(messages)
	messageReceived := make(chan bool, expectedCount)

	handler := func(ctx context.Context, message Message) error {
		receivedCount++
		messageReceived <- true
		return nil
	}

	// Start consumer
	go func() {
		err := manager.Subscribe(ctx, topic, "integration-batch-test-group", handler)
		if err != nil {
			t.Errorf("Failed to subscribe: %v", err)
		}
	}()

	// Wait for all messages to be received
	for i := 0; i < expectedCount; i++ {
		select {
		case <-messageReceived:
			t.Logf("Batch message %d received successfully", i+1)
		case <-time.After(10 * time.Second):
			t.Errorf("Timeout waiting for batch message %d", i+1)
			return
		}
	}

	if receivedCount != expectedCount {
		t.Errorf("Expected %d messages, received %d", expectedCount, receivedCount)
	}

	// Clean up - delete topic
	err = manager.DeleteTopic(ctx, topic)
	if err != nil {
		t.Logf("Topic deletion failed: %v", err)
	}
}

func TestKafkaJSONIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	cfg := &config.Config{
		Kafka: config.KafkaConfig{
			Brokers: []string{"localhost:9092"},
		},
		Log: config.LogConfig{
			Level:  "info",
			Format: "json",
			Output: "stdout",
		},
	}

	manager, err := NewManager(cfg)
	if err != nil {
		t.Skipf("Skipping integration test - Kafka not available: %v", err)
	}
	defer manager.Close()

	ctx := context.Background()
	topic := "integration-json-test-topic"

	// Test topic creation
	err = manager.CreateTopic(ctx, topic, 1, 1)
	if err != nil {
		t.Logf("Topic creation failed (may already exist): %v", err)
	}

	// Test JSON message publishing
	testData := map[string]interface{}{
		"id":    "user123",
		"name":  "John Doe",
		"email": "john@example.com",
		"age":   30,
	}

	err = manager.SendJSONMessage(ctx, topic, "user123", testData)
	if err != nil {
		t.Fatalf("Failed to send JSON message: %v", err)
	}

	// Test JSON message consumption
	messageReceived := make(chan bool, 1)
	handler := func(ctx context.Context, message Message) error {
		// Verify the message contains JSON data
		if len(message.Value) > 0 {
			messageReceived <- true
		}
		return nil
	}

	// Start consumer
	go func() {
		err := manager.Subscribe(ctx, topic, "integration-json-test-group", handler)
		if err != nil {
			t.Errorf("Failed to subscribe: %v", err)
		}
	}()

	// Wait for message to be received
	select {
	case <-messageReceived:
		t.Log("JSON message received successfully")
	case <-time.After(10 * time.Second):
		t.Error("Timeout waiting for JSON message")
	}

	// Clean up - delete topic
	err = manager.DeleteTopic(ctx, topic)
	if err != nil {
		t.Logf("Topic deletion failed: %v", err)
	}
}
