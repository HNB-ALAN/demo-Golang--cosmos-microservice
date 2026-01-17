package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/usc-platform/shared/config"
	"github.com/usc-platform/shared/health"
	messaging "github.com/usc-platform/shared/kafka-messaging"
	"github.com/usc-platform/shared/logging"
	"github.com/usc-platform/shared/metrics"
)

// Event represents a sample event structure
type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
	Source    string                 `json:"source"`
}

// KafkaService represents the Kafka service
type KafkaService struct {
	config  *config.Config
	logger  *logging.Logger
	kafka   messaging.KafkaClient
	metrics *metrics.KafkaMetrics
	health  *health.Registry
	server  *gin.Engine
}

// NewKafkaService creates a new Kafka service
func NewKafkaService(cfg *config.Config) (*KafkaService, error) {
	logger := logging.NewLogger("kafka-service", cfg.Log)

	// Initialize Kafka client
	kafkaClient, err := messaging.NewKafkaManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Kafka client: %w", err)
	}

	// Initialize metrics
	kafkaMetrics := metrics.NewKafkaMetrics(cfg.Service.Name)

	// Initialize health registry
	healthRegistry := health.NewRegistry()

	// Register Kafka service
	kafkaService := healthRegistry.RegisterService("kafka", "1.0.0")

	// Add Kafka health checks
	kafkaHealthChecker := health.NewKafkaHealthChecker("kafka", "Kafka connection health check", kafkaClient)
	kafkaService.RegisterCheck("kafka", kafkaHealthChecker)

	// Initialize Gin server
	gin.SetMode(gin.ReleaseMode)
	server := gin.New()
	server.Use(gin.Logger(), gin.Recovery())

	service := &KafkaService{
		config:  cfg,
		logger:  logger,
		kafka:   kafkaClient,
		metrics: kafkaMetrics,
		health:  healthRegistry,
		server:  server,
	}

	// Setup routes
	service.setupRoutes()

	return service, nil
}

// setupRoutes sets up HTTP routes
func (ks *KafkaService) setupRoutes() {
	// Health check endpoint
	ks.server.GET("/health", ks.healthCheckHandler)

	// Metrics endpoint
	ks.server.GET("/metrics", ks.metricsHandler)

	// Kafka endpoints
	api := ks.server.Group("/api/v1")
	{
		api.POST("/events", ks.publishEventHandler)
		api.POST("/events/batch", ks.publishBatchEventHandler)
		api.GET("/topics", ks.listTopicsHandler)
		api.POST("/topics", ks.createTopicHandler)
		api.DELETE("/topics/:name", ks.deleteTopicHandler)
	}

	// WebSocket endpoint for real-time events
	ks.server.GET("/ws/events", ks.websocketHandler)
}

// healthCheckHandler handles health check requests
func (ks *KafkaService) healthCheckHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	isHealthy := ks.health.IsHealthy(ctx)
	if !isHealthy {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now(),
	})
}

// metricsHandler handles metrics requests
func (ks *KafkaService) metricsHandler(c *gin.Context) {
	// In a real implementation, you would expose Prometheus metrics here
	c.JSON(http.StatusOK, gin.H{
		"message": "Metrics endpoint - implement Prometheus metrics exposure",
	})
}

// publishEventHandler handles single event publishing
func (ks *KafkaService) publishEventHandler(c *gin.Context) {
	var event Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid event format",
		})
		return
	}

	// Set event metadata
	event.ID = generateEventID()
	event.Timestamp = time.Now()
	event.Source = ks.config.Service.Name

	// Determine topic based on event type
	topic := ks.getTopicForEventType(event.Type)
	if topic == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unknown event type",
		})
		return
	}

	// Serialize event
	eventData, err := json.Marshal(event)
	if err != nil {
		ks.logger.Error("Failed to marshal event", logging.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to serialize event",
		})
		return
	}

	// Publish event
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	err = ks.kafka.SendMessage(ctx, topic, event.ID, eventData)
	if err != nil {
		ks.logger.Error("Failed to publish event", logging.Error(err), logging.String("topic", topic), logging.String("event_id", event.ID))
		ks.metrics.RecordProducerError(topic, "publish_failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to publish event",
		})
		return
	}

	ks.logger.Info("Event published successfully", logging.String("topic", topic), logging.String("event_id", event.ID), logging.String("event_type", event.Type))
	ks.metrics.RecordMessageProduced(topic, len(eventData), 0) // Latency would be measured in real implementation

	c.JSON(http.StatusOK, gin.H{
		"message":  "Event published successfully",
		"event_id": event.ID,
		"topic":    topic,
	})
}

// publishBatchEventHandler handles batch event publishing
func (ks *KafkaService) publishBatchEventHandler(c *gin.Context) {
	var events []Event
	if err := c.ShouldBindJSON(&events); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid events format",
		})
		return
	}

	if len(events) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No events provided",
		})
		return
	}

	// Group events by topic
	eventsByTopic := make(map[string][]messaging.Message)

	for _, event := range events {
		event.ID = generateEventID()
		event.Timestamp = time.Now()
		event.Source = ks.config.Service.Name

		topic := ks.getTopicForEventType(event.Type)
		if topic == "" {
			ks.logger.Warn("Skipping event with unknown type", logging.String("event_type", event.Type))
			continue
		}

		eventData, err := json.Marshal(event)
		if err != nil {
			ks.logger.Error("Failed to marshal event", logging.Error(err))
			continue
		}

		message := messaging.Message{
			Key:   event.ID,
			Value: eventData,
			Headers: map[string]string{
				"event_type": event.Type,
				"source":     event.Source,
			},
		}

		eventsByTopic[topic] = append(eventsByTopic[topic], message)
	}

	// Publish events by topic
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	var publishedCount int
	var errors []string

	for topic, messages := range eventsByTopic {
		err := ks.kafka.SendBatchMessages(ctx, topic, messages)
		if err != nil {
			ks.logger.Error("Failed to publish batch events", logging.Error(err), logging.String("topic", topic), logging.Int("count", len(messages)))
			ks.metrics.RecordProducerError(topic, "batch_publish_failed")
			errors = append(errors, fmt.Sprintf("Failed to publish to topic %s: %v", topic, err))
		} else {
			publishedCount += len(messages)
			ks.logger.Info("Batch events published successfully", logging.String("topic", topic), logging.Int("count", len(messages)))

			// Calculate total size for metrics
			totalSize := 0
			for _, msg := range messages {
				totalSize += len(msg.Value)
			}
			ks.metrics.RecordMessageProducedBatch(topic, len(messages), totalSize, 0)
		}
	}

	response := gin.H{
		"message":         "Batch events processed",
		"published_count": publishedCount,
		"total_count":     len(events),
	}

	if len(errors) > 0 {
		response["errors"] = errors
		c.JSON(http.StatusPartialContent, response)
	} else {
		c.JSON(http.StatusOK, response)
	}
}

// listTopicsHandler handles topic listing requests
func (ks *KafkaService) listTopicsHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	topics, err := ks.kafka.ListTopics(ctx)
	if err != nil {
		ks.logger.Error("Failed to list topics", logging.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list topics",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"topics": topics,
		"count":  len(topics),
	})
}

// createTopicHandler handles topic creation requests
func (ks *KafkaService) createTopicHandler(c *gin.Context) {
	var request struct {
		Name              string `json:"name" binding:"required"`
		Partitions        int    `json:"partitions"`
		ReplicationFactor int    `json:"replication_factor"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Set defaults
	if request.Partitions == 0 {
		request.Partitions = 1
	}
	if request.ReplicationFactor == 0 {
		request.ReplicationFactor = 1
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	err := ks.kafka.CreateTopic(ctx, request.Name, request.Partitions, request.ReplicationFactor)
	if err != nil {
		ks.logger.Error("Failed to create topic", logging.Error(err), logging.String("topic", request.Name))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create topic",
		})
		return
	}

	ks.logger.Info("Topic created successfully", logging.String("topic", request.Name), logging.Int("partitions", request.Partitions), logging.Int("replication_factor", request.ReplicationFactor))

	c.JSON(http.StatusCreated, gin.H{
		"message":            "Topic created successfully",
		"topic":              request.Name,
		"partitions":         request.Partitions,
		"replication_factor": request.ReplicationFactor,
	})
}

// deleteTopicHandler handles topic deletion requests
func (ks *KafkaService) deleteTopicHandler(c *gin.Context) {
	topicName := c.Param("name")
	if topicName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Topic name is required",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	err := ks.kafka.DeleteTopic(ctx, topicName)
	if err != nil {
		ks.logger.Error("Failed to delete topic", logging.Error(err), logging.String("topic", topicName))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete topic",
		})
		return
	}

	ks.logger.Info("Topic deleted successfully", logging.String("topic", topicName))

	c.JSON(http.StatusOK, gin.H{
		"message": "Topic deleted successfully",
		"topic":   topicName,
	})
}

// websocketHandler handles WebSocket connections for real-time events
func (ks *KafkaService) websocketHandler(c *gin.Context) {
	// In a real implementation, you would set up WebSocket connection here
	// and stream events from Kafka to connected clients
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "WebSocket endpoint not implemented",
	})
}

// getTopicForEventType returns the appropriate topic for an event type
func (ks *KafkaService) getTopicForEventType(eventType string) string {
	switch eventType {
	case "user.created", "user.updated", "user.deleted":
		return "user-events"
	case "content.created", "content.updated", "content.deleted":
		return "content-events"
	case "notification.sent", "notification.read":
		return "notification-events"
	case "analytics.view", "analytics.click", "analytics.conversion":
		return "analytics-events"
	default:
		return ""
	}
}

// generateEventID generates a unique event ID
func generateEventID() string {
	return fmt.Sprintf("evt_%d_%d", time.Now().Unix(), time.Now().UnixNano()%1000000)
}

// Start starts the Kafka service
func (ks *KafkaService) Start() error {
	ks.logger.Info("Starting Kafka service", logging.String("port", ks.config.Server.Port))

	// Start HTTP server
	go func() {
		if err := ks.server.Run(":" + ks.config.Server.Port); err != nil {
			ks.logger.Error("Failed to start HTTP server", logging.Error(err))
		}
	}()

	// Start event consumer
	go ks.startEventConsumer()

	ks.logger.Info("Kafka service started successfully")
	return nil
}

// startEventConsumer starts consuming events from Kafka
func (ks *KafkaService) startEventConsumer() {
	ctx := context.Background()

	// Consumer for user events
	go ks.consumeEvents(ctx, "user-events", "user-events-consumer")

	// Consumer for content events
	go ks.consumeEvents(ctx, "content-events", "content-events-consumer")

	// Consumer for notification events
	go ks.consumeEvents(ctx, "notification-events", "notification-events-consumer")

	// Consumer for analytics events
	go ks.consumeEvents(ctx, "analytics-events", "analytics-events-consumer")
}

// consumeEvents consumes events from a specific topic
func (ks *KafkaService) consumeEvents(ctx context.Context, topic, groupID string) {
	handler := func(ctx context.Context, message messaging.Message) error {
		var event Event
		if err := json.Unmarshal(message.Value, &event); err != nil {
			ks.logger.Error("Failed to unmarshal event", logging.Error(err), logging.String("topic", topic))
			return err
		}

		ks.logger.Info("Processing event",
			logging.String("topic", topic),
			logging.String("event_id", event.ID),
			logging.String("event_type", event.Type),
			logging.String("source", event.Source))

		// Process the event based on its type
		switch event.Type {
		case "user.created":
			ks.processUserCreatedEvent(event)
		case "user.updated":
			ks.processUserUpdatedEvent(event)
		case "content.created":
			ks.processContentCreatedEvent(event)
		case "notification.sent":
			ks.processNotificationSentEvent(event)
		case "analytics.view":
			ks.processAnalyticsViewEvent(event)
		default:
			ks.logger.Info("Unknown event type, skipping", logging.String("event_type", event.Type))
		}

		return nil
	}

	err := ks.kafka.Subscribe(ctx, topic, groupID, handler)
	if err != nil {
		ks.logger.Error("Failed to subscribe to topic", logging.Error(err), logging.String("topic", topic), logging.String("group_id", groupID))
	}
}

// Event processing methods
func (ks *KafkaService) processUserCreatedEvent(event Event) {
	ks.logger.Info("Processing user created event", logging.String("event_id", event.ID), logging.String("user_id", fmt.Sprintf("%v", event.Data["user_id"])))
	// Implement user creation logic here
}

func (ks *KafkaService) processUserUpdatedEvent(event Event) {
	ks.logger.Info("Processing user updated event", logging.String("event_id", event.ID), logging.String("user_id", fmt.Sprintf("%v", event.Data["user_id"])))
	// Implement user update logic here
}

func (ks *KafkaService) processContentCreatedEvent(event Event) {
	ks.logger.Info("Processing content created event", logging.String("event_id", event.ID), logging.String("content_id", fmt.Sprintf("%v", event.Data["content_id"])))
	// Implement content creation logic here
}

func (ks *KafkaService) processNotificationSentEvent(event Event) {
	ks.logger.Info("Processing notification sent event", logging.String("event_id", event.ID), logging.String("notification_id", fmt.Sprintf("%v", event.Data["notification_id"])))
	// Implement notification processing logic here
}

func (ks *KafkaService) processAnalyticsViewEvent(event Event) {
	ks.logger.Info("Processing analytics view event", logging.String("event_id", event.ID), logging.String("content_id", fmt.Sprintf("%v", event.Data["content_id"])))
	// Implement analytics processing logic here
}

// Stop stops the Kafka service
func (ks *KafkaService) Stop() error {
	ks.logger.Info("Stopping Kafka service")

	// Close Kafka client
	if err := ks.kafka.Close(); err != nil {
		ks.logger.Error("Failed to close Kafka client", logging.Error(err))
	}

	ks.logger.Info("Kafka service stopped")
	return nil
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("", "kafka-service")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create Kafka service
	service, err := NewKafkaService(cfg)
	if err != nil {
		log.Fatalf("Failed to create Kafka service: %v", err)
	}

	// Start service
	if err := service.Start(); err != nil {
		log.Fatalf("Failed to start Kafka service: %v", err)
	}

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Stop service
	if err := service.Stop(); err != nil {
		log.Printf("Error stopping service: %v", err)
	}

	log.Println("Kafka service stopped")
}
