package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/usc-platform/shared/config"
	messaging "github.com/usc-platform/shared/kafka-messaging"
	"github.com/usc-platform/shared/logging"
)

// EventType represents the type of event
type EventType string

const (
	// User Events
	EventTypeUserCreated EventType = "user.created"
	EventTypeUserUpdated EventType = "user.updated"
	EventTypeUserDeleted EventType = "user.deleted"
	EventTypeUserLogin   EventType = "user.login"
	EventTypeUserLogout  EventType = "user.logout"

	// Content Events
	EventTypeContentCreated EventType = "content.created"
	EventTypeContentUpdated EventType = "content.updated"
	EventTypeContentDeleted EventType = "content.deleted"
	EventTypeContentLiked   EventType = "content.liked"
	EventTypeContentShared  EventType = "content.shared"

	// Transaction Events
	EventTypeTransactionCreated   EventType = "transaction.created"
	EventTypeTransactionCompleted EventType = "transaction.completed"
	EventTypeTransactionFailed    EventType = "transaction.failed"

	// System Events
	EventTypeSystemStartup  EventType = "system.startup"
	EventTypeSystemShutdown EventType = "system.shutdown"
	EventTypeSystemError    EventType = "system.error"
)

// Event represents a domain event
type Event struct {
	ID          string                 `json:"id"`
	Type        EventType              `json:"type"`
	Source      string                 `json:"source"`
	Data        map[string]interface{} `json:"data"`
	Metadata    map[string]interface{} `json:"metadata"`
	Timestamp   time.Time              `json:"timestamp"`
	Version     int                    `json:"version"`
	AggregateID string                 `json:"aggregate_id"`
}

// EventHandler handles domain events
type EventHandler func(ctx context.Context, event Event) error

// KafkaEventManager manages Kafka event publishing and consumption for USC Blockchain Core Service
type KafkaEventManager struct {
	manager *messaging.KafkaManager
	config  *config.Config
	logger  *logging.Logger
}

// NewKafkaEventManager creates a new Kafka event manager
func NewKafkaEventManager(cfg *config.Config, logger *logging.Logger) (*KafkaEventManager, error) {
	manager, err := messaging.NewKafkaManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka event manager: %w", err)
	}

	return &KafkaEventManager{
		manager: manager,
		config:  cfg,
		logger:  logger,
	}, nil
}

// GetManager returns the Kafka manager
func (k *KafkaEventManager) GetManager() *messaging.KafkaManager {
	return k.manager
}

// Close closes the Kafka event manager
func (k *KafkaEventManager) Close() error {
	if k.manager != nil {
		return k.manager.Close()
	}
	return nil
}

// HealthCheck performs a health check on the Kafka event manager
func (k *KafkaEventManager) HealthCheck(ctx context.Context) error {
	if k.manager == nil {
		return fmt.Errorf("kafka event manager is not initialized")
	}

	// Simple health check - try to list topics
	_, err := k.manager.ListTopics(ctx)
	if err != nil {
		return fmt.Errorf("kafka event health check failed: %w", err)
	}

	return nil
}

// PublishEvent publishes a domain event to Kafka
func (k *KafkaEventManager) PublishEvent(ctx context.Context, event Event) error {
	if k.manager == nil {
		return fmt.Errorf("kafka event manager is not initialized")
	}

	// Generate event ID if not provided
	if event.ID == "" {
		event.ID = fmt.Sprintf("%s-%d", event.Type, time.Now().UnixNano())
	}

	// Set timestamp if not provided
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// Set version if not provided
	if event.Version == 0 {
		event.Version = 1
	}

	// Marshal event to JSON
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Determine topic based on event type
	topic := k.getTopicForEventType(event.Type)

	// Publish event
	err = k.manager.SendMessage(ctx, topic, event.ID, eventData)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	k.logger.Info("Event published",
		logging.String("event_id", event.ID),
		logging.String("event_type", string(event.Type)),
		logging.String("topic", topic),
		logging.String("source", event.Source))

	return nil
}

// PublishEventWithRetry publishes a domain event with retry logic
func (k *KafkaEventManager) PublishEventWithRetry(ctx context.Context, event Event, maxRetries int) error {
	if k.manager == nil {
		return fmt.Errorf("kafka event manager is not initialized")
	}

	// Generate event ID if not provided
	if event.ID == "" {
		event.ID = fmt.Sprintf("%s-%d", event.Type, time.Now().UnixNano())
	}

	// Set timestamp if not provided
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// Set version if not provided
	if event.Version == 0 {
		event.Version = 1
	}

	// Marshal event to JSON
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Determine topic based on event type
	topic := k.getTopicForEventType(event.Type)

	// Publish event with retry
	err = k.manager.SendMessageWithRetry(ctx, topic, event.ID, eventData, maxRetries)
	if err != nil {
		return fmt.Errorf("failed to publish event with retry: %w", err)
	}

	k.logger.Info("Event published with retry",
		logging.String("event_id", event.ID),
		logging.String("event_type", string(event.Type)),
		logging.String("topic", topic),
		logging.String("source", event.Source),
		logging.Int("max_retries", maxRetries))

	return nil
}

// SubscribeToEvents subscribes to events of a specific type
func (k *KafkaEventManager) SubscribeToEvents(ctx context.Context, eventType EventType, handler EventHandler) error {
	if k.manager == nil {
		return fmt.Errorf("kafka event manager is not initialized")
	}

	topic := k.getTopicForEventType(eventType)
	groupID := fmt.Sprintf("%s-%s-consumer", k.config.Service.Name, eventType)

	// Create message handler
	messageHandler := func(ctx context.Context, message messaging.Message) error {
		var event Event
		if err := json.Unmarshal(message.Value, &event); err != nil {
			k.logger.Error("Failed to unmarshal event",
				logging.Error(err),
				logging.String("topic", topic),
				logging.String("key", message.Key))
			return err
		}

		// Call the event handler
		if err := handler(ctx, event); err != nil {
			k.logger.Error("Failed to handle event",
				logging.Error(err),
				logging.String("event_id", event.ID),
				logging.String("event_type", string(event.Type)))
			return err
		}

		k.logger.Debug("Event handled successfully",
			logging.String("event_id", event.ID),
			logging.String("event_type", string(event.Type)))

		return nil
	}

	// Subscribe to topic
	err := k.manager.Subscribe(ctx, topic, groupID, messageHandler)
	if err != nil {
		return fmt.Errorf("failed to subscribe to events: %w", err)
	}

	k.logger.Info("Subscribed to events",
		logging.String("event_type", string(eventType)),
		logging.String("topic", topic),
		logging.String("group_id", groupID))

	return nil
}

// SubscribeToAllEvents subscribes to all events for this service
func (k *KafkaEventManager) SubscribeToAllEvents(ctx context.Context, handler EventHandler) error {
	if k.manager == nil {
		return fmt.Errorf("kafka event manager is not initialized")
	}

	// Subscribe to all relevant event types for this service
	eventTypes := k.getRelevantEventTypes()

	for _, eventType := range eventTypes {
		if err := k.SubscribeToEvents(ctx, eventType, handler); err != nil {
			k.logger.Error("Failed to subscribe to event type",
				logging.Error(err),
				logging.String("event_type", string(eventType)))
			// Continue with other event types
		}
	}

	return nil
}

// getTopicForEventType returns the Kafka topic for a given event type
func (k *KafkaEventManager) getTopicForEventType(eventType EventType) string {
	// Map event types to topics
	switch {
	case eventType == EventTypeUserCreated || eventType == EventTypeUserUpdated || eventType == EventTypeUserDeleted:
		return "user-events"
	case eventType == EventTypeUserLogin || eventType == EventTypeUserLogout:
		return "auth-events"
	case eventType == EventTypeContentCreated || eventType == EventTypeContentUpdated || eventType == EventTypeContentDeleted:
		return "content-events"
	case eventType == EventTypeContentLiked || eventType == EventTypeContentShared:
		return "social-events"
	case eventType == EventTypeTransactionCreated || eventType == EventTypeTransactionCompleted || eventType == EventTypeTransactionFailed:
		return "transaction-events"
	case eventType == EventTypeSystemStartup || eventType == EventTypeSystemShutdown || eventType == EventTypeSystemError:
		return "system-events"
	default:
		return "general-events"
	}
}

// getRelevantEventTypes returns the event types relevant to this service
func (k *KafkaEventManager) getRelevantEventTypes() []EventType {
	// This should be customized based on the specific service
	// For now, return all event types
	return []EventType{
		EventTypeUserCreated,
		EventTypeUserUpdated,
		EventTypeUserDeleted,
		EventTypeUserLogin,
		EventTypeUserLogout,
		EventTypeContentCreated,
		EventTypeContentUpdated,
		EventTypeContentDeleted,
		EventTypeContentLiked,
		EventTypeContentShared,
		EventTypeTransactionCreated,
		EventTypeTransactionCompleted,
		EventTypeTransactionFailed,
		EventTypeSystemStartup,
		EventTypeSystemShutdown,
		EventTypeSystemError,
	}
}

// CreateEvent creates a new domain event
func CreateEvent(eventType EventType, source string, data map[string]interface{}) Event {
	return Event{
		ID:        fmt.Sprintf("%s-%d", eventType, time.Now().UnixNano()),
		Type:      eventType,
		Source:    source,
		Data:      data,
		Metadata:  make(map[string]interface{}),
		Timestamp: time.Now(),
		Version:   1,
	}
}

// CreateEventWithAggregate creates a new domain event with aggregate ID
func CreateEventWithAggregate(eventType EventType, source string, aggregateID string, data map[string]interface{}) Event {
	event := CreateEvent(eventType, source, data)
	event.AggregateID = aggregateID
	return event
}
