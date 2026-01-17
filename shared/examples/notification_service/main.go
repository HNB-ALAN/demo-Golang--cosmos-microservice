package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/usc-platform/shared/config"
	"github.com/usc-platform/shared/logging"
	"github.com/usc-platform/shared/middleware"
	"github.com/usc-platform/shared/notifications"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("", "notification-service")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger := logging.NewLogger("notification-service", cfg.Log)
	logger.Info("Starting notification service")

	// Initialize notification channel service
	channelConfig := &notifications.ChannelConfig{
		EmailConfig: &notifications.EmailConfig{
			Provider:     getEnvOrDefault("EMAIL_PROVIDER", "sendgrid"),
			APIKey:       getEnvOrDefault("EMAIL_API_KEY", ""),
			FromEmail:    getEnvOrDefault("FROM_EMAIL", "noreply@uscplatform.com"),
			FromName:     getEnvOrDefault("FROM_NAME", "USC Platform"),
			TemplatePath: getEnvOrDefault("EMAIL_TEMPLATE_PATH", "./templates"),
			RateLimit:    1000,
		},
		SMSConfig: &notifications.SMSConfig{
			Provider:  getEnvOrDefault("SMS_PROVIDER", "twilio"),
			APIKey:    getEnvOrDefault("SMS_API_KEY", ""),
			FromPhone: getEnvOrDefault("FROM_PHONE", "+1234567890"),
			RateLimit: 100,
		},
		PushConfig: &notifications.PushConfig{
			FirebaseConfig: &notifications.FirebaseConfig{
				ProjectID:     getEnvOrDefault("FIREBASE_PROJECT_ID", ""),
				ServiceKey:    getEnvOrDefault("FIREBASE_SERVICE_KEY", ""),
				DatabaseURL:   getEnvOrDefault("FIREBASE_DATABASE_URL", ""),
				StorageBucket: getEnvOrDefault("FIREBASE_STORAGE_BUCKET", ""),
			},
			APNsConfig: &notifications.APNsConfig{
				KeyID:       getEnvOrDefault("APNS_KEY_ID", ""),
				TeamID:      getEnvOrDefault("APNS_TEAM_ID", ""),
				BundleID:    getEnvOrDefault("APNS_BUNDLE_ID", ""),
				PrivateKey:  getEnvOrDefault("APNS_PRIVATE_KEY", ""),
				Environment: getEnvOrDefault("APNS_ENVIRONMENT", "sandbox"),
			},
			WebPushConfig: &notifications.WebPushConfig{
				VAPIDPublicKey:  getEnvOrDefault("VAPID_PUBLIC_KEY", ""),
				VAPIDPrivateKey: getEnvOrDefault("VAPID_PRIVATE_KEY", ""),
				VAPIDSubject:    getEnvOrDefault("VAPID_SUBJECT", "mailto:admin@uscplatform.com"),
			},
		},
		InAppConfig: &notifications.InAppConfig{
			DatabaseURL: getEnvOrDefault("INAPP_DATABASE_URL", ""),
			TableName:   "notifications",
			TTL:         24,
		},
		WebhookConfig: &notifications.WebhookConfig{
			Timeout:    30 * time.Second,
			RetryCount: 3,
			RetryDelay: 5 * time.Second,
			SecretKey:  getEnvOrDefault("WEBHOOK_SECRET_KEY", ""),
		},
		DefaultChannel: "email",
		RetryAttempts:  3,
		RetryDelay:     5 * time.Second,
	}

	notificationService := notifications.NewNotificationChannelService(logger, channelConfig)

	// Create Gin router
	router := gin.Default()

	// Add CORS middleware
	corsConfig := middleware.DefaultCORSConfig()
	corsMiddleware := middleware.NewCORSMiddleware(corsConfig)
	router.Use(gin.WrapH(corsMiddleware.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This will be handled by the CORS middleware
	}))))

	// Add rate limiting middleware
	rateLimitConfig := middleware.RateLimitConfig{
		RequestsPerSecond: 16.67, // 1000 requests per minute = 16.67 per second
		Burst:             100,
	}
	rateLimitMiddleware := middleware.NewHTTPRateLimitMiddleware(rateLimitConfig)
	router.Use(gin.WrapH(rateLimitMiddleware.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This will be handled by the rate limit middleware
	}))))

	// Add request logging middleware
	router.Use(gin.Logger())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "timestamp": time.Now()})
	})

	// Send single notification
	router.POST("/notifications/send", func(c *gin.Context) {
		var message notifications.NotificationMessage
		if err := c.ShouldBindJSON(&message); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Set message ID if not provided
		if message.ID == "" {
			message.ID = generateMessageID()
		}

		// Set creation timestamp
		message.CreatedAt = time.Now()

		// Send notification
		delivery, err := notificationService.SendNotification(context.Background(), &message)
		if err != nil {
			logger.Error("Failed to send notification",
				logging.String("message_id", message.ID),
				logging.String("channel", message.Channel),
				logging.Error(err),
			)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"message_id":  message.ID,
			"delivery_id": delivery.ID,
			"status":      delivery.Status,
			"channel":     delivery.Channel,
		})
	})

	// Send batch notifications
	router.POST("/notifications/batch", func(c *gin.Context) {
		var messages []notifications.NotificationMessage
		if err := c.ShouldBindJSON(&messages); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Set message IDs and timestamps
		for i := range messages {
			if messages[i].ID == "" {
				messages[i].ID = generateMessageID()
			}
			messages[i].CreatedAt = time.Now()
		}

		// Convert to pointers
		messagePointers := make([]*notifications.NotificationMessage, len(messages))
		for i := range messages {
			messagePointers[i] = &messages[i]
		}

		// Send batch notifications
		deliveries, err := notificationService.SendBatchNotifications(context.Background(), messagePointers)
		if err != nil {
			logger.Error("Failed to send batch notifications",
				logging.Int("count", len(messages)),
				logging.Error(err),
			)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"total_count":   len(messages),
			"success_count": len(deliveries),
			"deliveries":    deliveries,
		})
	})

	// Get delivery status
	router.GET("/notifications/delivery/:id", func(c *gin.Context) {
		deliveryID := c.Param("id")

		delivery, err := notificationService.GetDeliveryStatus(context.Background(), deliveryID)
		if err != nil {
			logger.Error("Failed to get delivery status",
				logging.String("delivery_id", deliveryID),
				logging.Error(err),
			)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, delivery)
	})

	// Retry failed notification
	router.POST("/notifications/retry/:id", func(c *gin.Context) {
		deliveryID := c.Param("id")

		delivery, err := notificationService.RetryFailedNotification(context.Background(), deliveryID)
		if err != nil {
			logger.Error("Failed to retry notification",
				logging.String("delivery_id", deliveryID),
				logging.Error(err),
			)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, delivery)
	})

	// Test notification endpoint
	router.POST("/notifications/test", func(c *gin.Context) {
		var request struct {
			Channel string `json:"channel" binding:"required"`
			UserID  string `json:"user_id" binding:"required"`
			Type    string `json:"type"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Create test message
		message := &notifications.NotificationMessage{
			ID:      generateMessageID(),
			UserID:  request.UserID,
			Channel: request.Channel,
			Type:    getOrDefault(request.Type, "test"),
			Title:   "Test Notification",
			Content: "This is a test notification from USC Platform",
			Data: map[string]interface{}{
				"test":      true,
				"timestamp": time.Now(),
			},
			Priority:  "normal",
			CreatedAt: time.Now(),
		}

		// Send test notification
		delivery, err := notificationService.SendNotification(context.Background(), message)
		if err != nil {
			logger.Error("Failed to send test notification",
				logging.String("channel", request.Channel),
				logging.Error(err),
			)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"message":  "Test notification sent successfully",
			"delivery": delivery,
		})
	})

	// Start server
	port := getEnvOrDefault("PORT", "4002")
	logger.Info("Starting notification service", logging.String("port", port))

	if err := router.Run(":" + port); err != nil {
		logger.Fatal("Failed to start notification service", logging.Error(err))
	}

	logger.Info("Notification service started successfully")
}

// generateMessageID generates a unique message ID
func generateMessageID() string {
	return time.Now().Format("20060102150405") + "-" + "notification"
}

// getEnvOrDefault returns environment variable value or default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getOrDefault returns value or default
func getOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
