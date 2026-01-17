package notifications

import (
	"context"
	"fmt"
	"testing"
	"time"

	sharedConfig "github.com/usc-platform/shared/config"
	"github.com/usc-platform/shared/logging"
)

func TestNotificationChannelService_NewNotificationChannelService(t *testing.T) {
	logger := logging.NewLogger("test", sharedConfig.LogConfig{})
	config := &ChannelConfig{
		EmailConfig: &EmailConfig{
			Provider:     "smtp",
			APIKey:       "test-key",
			FromEmail:    "test@example.com",
			FromName:     "Test Sender",
			TemplatePath: "./templates",
			RateLimit:    100,
		},
		SMSConfig: &SMSConfig{
			Provider:  "twilio",
			APIKey:    "test-sms-key",
			FromPhone: "+1234567890",
			RateLimit: 50,
		},
		PushConfig: &PushConfig{
			FirebaseConfig: &FirebaseConfig{
				ProjectID: "test-project",
			},
		},
		InAppConfig: &InAppConfig{
			DatabaseURL: "postgres://localhost/test",
			TableName:   "notifications",
			TTL:         24,
		},
		WebhookConfig: &WebhookConfig{
			Timeout:    30 * time.Second,
			RetryCount: 3,
			RetryDelay: 5 * time.Second,
			SecretKey:  "test-webhook-secret",
		},
		DefaultChannel: "email",
		RetryAttempts:  3,
		RetryDelay:     5 * time.Second,
	}

	service := NewNotificationChannelService(logger, config)

	if service == nil {
		t.Fatal("Expected notification channel service to be created")
	}

	if service.logger != logger {
		t.Error("Expected service to use provided logger")
	}

	if service.config != config {
		t.Error("Expected service to use provided config")
	}
}

func TestNotificationChannelService_SendNotification(t *testing.T) {
	logger := logging.NewLogger("test", sharedConfig.LogConfig{})
	config := &ChannelConfig{
		EmailConfig: &EmailConfig{
			Provider:     "smtp",
			APIKey:       "test-key",
			FromEmail:    "test@example.com",
			FromName:     "Test Sender",
			TemplatePath: "./templates",
			RateLimit:    100,
		},
		DefaultChannel: "email",
		RetryAttempts:  3,
		RetryDelay:     5 * time.Second,
	}

	service := NewNotificationChannelService(logger, config)

	message := &NotificationMessage{
		ID:        "test-message-1",
		UserID:    "user-123",
		Channel:   "email",
		Type:      "welcome",
		Title:     "Welcome!",
		Content:   "Welcome to our platform!",
		Recipient: "user@example.com",
		Subject:   "Welcome to USC Platform",
		Data:      make(map[string]interface{}),
		Priority:  "normal",
		CreatedAt: time.Now(),
	}

	ctx := context.Background()
	delivery, err := service.SendNotification(ctx, message)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if delivery == nil {
		t.Fatal("Expected delivery to be returned")
	}

	if delivery.MessageID != message.ID {
		t.Errorf("Expected message ID %s, got %s", message.ID, delivery.MessageID)
	}

	if delivery.Channel != message.Channel {
		t.Errorf("Expected channel %s, got %s", message.Channel, delivery.Channel)
	}
}

func TestNotificationChannelService_SendNotificationWithInvalidChannel(t *testing.T) {
	logger := logging.NewLogger("test", sharedConfig.LogConfig{})
	config := &ChannelConfig{
		EmailConfig: &EmailConfig{
			Provider:     "smtp",
			APIKey:       "test-key",
			FromEmail:    "test@example.com",
			FromName:     "Test Sender",
			TemplatePath: "./templates",
			RateLimit:    100,
		},
		DefaultChannel: "email",
		RetryAttempts:  3,
		RetryDelay:     5 * time.Second,
	}

	service := NewNotificationChannelService(logger, config)

	message := &NotificationMessage{
		ID:        "test-message-1",
		UserID:    "user-123",
		Channel:   "invalid-channel", // Invalid channel
		Type:      "welcome",
		Title:     "Welcome!",
		Content:   "Welcome to our platform!",
		Recipient: "user@example.com",
		Subject:   "Welcome to USC Platform",
		Data:      make(map[string]interface{}),
		Priority:  "normal",
		CreatedAt: time.Now(),
	}

	ctx := context.Background()
	_, err := service.SendNotification(ctx, message)

	if err == nil {
		t.Error("Expected error for invalid channel, got nil")
	}
}

func TestNotificationChannelService_SendBatchNotifications(t *testing.T) {
	logger := logging.NewLogger("test", sharedConfig.LogConfig{})
	config := &ChannelConfig{
		EmailConfig: &EmailConfig{
			Provider:     "smtp",
			APIKey:       "test-key",
			FromEmail:    "test@example.com",
			FromName:     "Test Sender",
			TemplatePath: "./templates",
			RateLimit:    100,
		},
		DefaultChannel: "email",
		RetryAttempts:  3,
		RetryDelay:     5 * time.Second,
	}

	service := NewNotificationChannelService(logger, config)

	messages := []*NotificationMessage{
		{
			ID:        "test-message-1",
			UserID:    "user-123",
			Channel:   "email",
			Type:      "welcome",
			Title:     "Welcome!",
			Content:   "Welcome to our platform!",
			Recipient: "user1@example.com",
			Subject:   "Welcome to USC Platform",
			Data:      make(map[string]interface{}),
			Priority:  "normal",
			CreatedAt: time.Now(),
		},
		{
			ID:        "test-message-2",
			UserID:    "user-456",
			Channel:   "email",
			Type:      "welcome",
			Title:     "Welcome!",
			Content:   "Welcome to our platform!",
			Recipient: "user2@example.com",
			Subject:   "Welcome to USC Platform",
			Data:      make(map[string]interface{}),
			Priority:  "normal",
			CreatedAt: time.Now(),
		},
	}

	ctx := context.Background()
	deliveries, err := service.SendBatchNotifications(ctx, messages)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(deliveries) != len(messages) {
		t.Errorf("Expected %d deliveries, got %d", len(messages), len(deliveries))
	}

	for i, delivery := range deliveries {
		if delivery.MessageID != messages[i].ID {
			t.Errorf("Expected message ID %s, got %s", messages[i].ID, delivery.MessageID)
		}
	}
}

func TestNotificationChannelService_GetDeliveryStatus(t *testing.T) {
	logger := logging.NewLogger("test", sharedConfig.LogConfig{})
	config := &ChannelConfig{
		EmailConfig: &EmailConfig{
			Provider:     "smtp",
			APIKey:       "test-key",
			FromEmail:    "test@example.com",
			FromName:     "Test Sender",
			TemplatePath: "./templates",
			RateLimit:    100,
		},
		DefaultChannel: "email",
		RetryAttempts:  3,
		RetryDelay:     5 * time.Second,
	}

	service := NewNotificationChannelService(logger, config)

	ctx := context.Background()
	status, err := service.GetDeliveryStatus(ctx, "test-delivery-id")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if status == nil {
		t.Fatal("Expected delivery status to be returned")
	}

	if status.ID != "test-delivery-id" {
		t.Errorf("Expected delivery ID %s, got %s", "test-delivery-id", status.ID)
	}
}

func TestNotificationChannelService_RetryFailedNotification(t *testing.T) {
	logger := logging.NewLogger("test", sharedConfig.LogConfig{})
	config := &ChannelConfig{
		EmailConfig: &EmailConfig{
			Provider:     "smtp",
			APIKey:       "test-key",
			FromEmail:    "test@example.com",
			FromName:     "Test Sender",
			TemplatePath: "./templates",
			RateLimit:    100,
		},
		DefaultChannel: "email",
		RetryAttempts:  3,
		RetryDelay:     5 * time.Second,
	}

	service := NewNotificationChannelService(logger, config)

	ctx := context.Background()
	delivery, err := service.RetryFailedNotification(ctx, "test-delivery-id")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if delivery == nil {
		t.Fatal("Expected delivery to be returned")
	}

	if delivery.ID != "test-delivery-id" {
		t.Errorf("Expected delivery ID %s, got %s", "test-delivery-id", delivery.ID)
	}
}

func TestEmailProvider_Send(t *testing.T) {
	config := &EmailConfig{
		Provider:     "smtp",
		APIKey:       "test-key",
		FromEmail:    "test@example.com",
		FromName:     "Test Sender",
		TemplatePath: "./templates",
		RateLimit:    100,
	}

	logger := logging.NewLogger("test", sharedConfig.LogConfig{})
	provider := &EmailProvider{
		config: config,
		logger: logger,
	}

	message := &NotificationMessage{
		ID:        "test-message-1",
		UserID:    "user-123",
		Channel:   "email",
		Type:      "welcome",
		Title:     "Welcome!",
		Content:   "Welcome to our platform!",
		Recipient: "user@example.com",
		Subject:   "Welcome to USC Platform",
		Data:      make(map[string]interface{}),
		Priority:  "normal",
		CreatedAt: time.Now(),
	}

	ctx := context.Background()
	delivery, err := provider.Send(ctx, message)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if delivery == nil {
		t.Fatal("Expected delivery to be returned")
	}

	if delivery.MessageID != message.ID {
		t.Errorf("Expected message ID %s, got %s", message.ID, delivery.MessageID)
	}

	if delivery.Channel != "email" {
		t.Errorf("Expected channel 'email', got %s", delivery.Channel)
	}
}

func TestEmailProvider_SendWithInvalidEmail(t *testing.T) {
	config := &EmailConfig{
		Provider:     "smtp",
		APIKey:       "test-key",
		FromEmail:    "test@example.com",
		FromName:     "Test Sender",
		TemplatePath: "./templates",
		RateLimit:    100,
	}

	logger := logging.NewLogger("test", sharedConfig.LogConfig{})
	provider := &EmailProvider{
		config: config,
		logger: logger,
	}

	message := &NotificationMessage{
		ID:        "test-message-1",
		UserID:    "user-123",
		Channel:   "email",
		Type:      "welcome",
		Title:     "Welcome!",
		Content:   "Welcome to our platform!",
		Recipient: "invalid-email", // Invalid email
		Subject:   "Welcome to USC Platform",
		Data:      make(map[string]interface{}),
		Priority:  "normal",
		CreatedAt: time.Now(),
	}

	ctx := context.Background()
	_, err := provider.Send(ctx, message)

	if err == nil {
		t.Error("Expected error for invalid email, got nil")
	}
}

func TestEmailProvider_SendWithEmptyRecipient(t *testing.T) {
	config := &EmailConfig{
		Provider:     "smtp",
		APIKey:       "test-key",
		FromEmail:    "test@example.com",
		FromName:     "Test Sender",
		TemplatePath: "./templates",
		RateLimit:    100,
	}

	logger := logging.NewLogger("test", sharedConfig.LogConfig{})
	provider := &EmailProvider{
		config: config,
		logger: logger,
	}

	message := &NotificationMessage{
		ID:        "test-message-1",
		UserID:    "user-123",
		Channel:   "email",
		Type:      "welcome",
		Title:     "Welcome!",
		Content:   "Welcome to our platform!",
		Recipient: "", // Empty recipient
		Subject:   "Welcome to USC Platform",
		Data:      make(map[string]interface{}),
		Priority:  "normal",
		CreatedAt: time.Now(),
	}

	ctx := context.Background()
	_, err := provider.Send(ctx, message)

	if err == nil {
		t.Error("Expected error for empty recipient, got nil")
	}
}

func TestEmailProvider_SendWithEmptySubject(t *testing.T) {
	config := &EmailConfig{
		Provider:     "smtp",
		APIKey:       "test-key",
		FromEmail:    "test@example.com",
		FromName:     "Test Sender",
		TemplatePath: "./templates",
		RateLimit:    100,
	}

	logger := logging.NewLogger("test", sharedConfig.LogConfig{})
	provider := &EmailProvider{
		config: config,
		logger: logger,
	}

	message := &NotificationMessage{
		ID:        "test-message-1",
		UserID:    "user-123",
		Channel:   "email",
		Type:      "welcome",
		Title:     "Welcome!",
		Content:   "Welcome to our platform!",
		Recipient: "user@example.com",
		Subject:   "", // Empty subject
		Data:      make(map[string]interface{}),
		Priority:  "normal",
		CreatedAt: time.Now(),
	}

	ctx := context.Background()
	_, err := provider.Send(ctx, message)

	if err == nil {
		t.Error("Expected error for empty subject, got nil")
	}
}

func TestEmailProvider_SendWithEmptyContent(t *testing.T) {
	config := &EmailConfig{
		Provider:     "smtp",
		APIKey:       "test-key",
		FromEmail:    "test@example.com",
		FromName:     "Test Sender",
		TemplatePath: "./templates",
		RateLimit:    100,
	}

	logger := logging.NewLogger("test", sharedConfig.LogConfig{})
	provider := &EmailProvider{
		config: config,
		logger: logger,
	}

	message := &NotificationMessage{
		ID:        "test-message-1",
		UserID:    "user-123",
		Channel:   "email",
		Type:      "welcome",
		Title:     "Welcome!",
		Content:   "", // Empty content
		Recipient: "user@example.com",
		Subject:   "Welcome to USC Platform",
		Data:      make(map[string]interface{}),
		Priority:  "normal",
		CreatedAt: time.Now(),
	}

	ctx := context.Background()
	_, err := provider.Send(ctx, message)

	if err == nil {
		t.Error("Expected error for empty content, got nil")
	}
}

func TestEmailProvider_Validate(t *testing.T) {
	config := &EmailConfig{
		Provider:     "smtp",
		APIKey:       "test-key",
		FromEmail:    "test@example.com",
		FromName:     "Test Sender",
		TemplatePath: "./templates",
		RateLimit:    100,
	}

	logger := logging.NewLogger("test", sharedConfig.LogConfig{})
	provider := &EmailProvider{
		config: config,
		logger: logger,
	}

	message := &NotificationMessage{
		ID:        "test-message-1",
		UserID:    "user-123",
		Channel:   "email",
		Type:      "welcome",
		Title:     "Welcome!",
		Content:   "Welcome to our platform!",
		Recipient: "user@example.com",
		Subject:   "Welcome to USC Platform",
		Data:      make(map[string]interface{}),
		Priority:  "normal",
		CreatedAt: time.Now(),
	}

	err := provider.Validate(message)
	if err != nil {
		t.Errorf("Expected no error for valid message, got %v", err)
	}
}

func TestEmailProvider_GetName(t *testing.T) {
	config := &EmailConfig{
		Provider:     "smtp",
		APIKey:       "test-key",
		FromEmail:    "test@example.com",
		FromName:     "Test Sender",
		TemplatePath: "./templates",
		RateLimit:    100,
	}

	logger := logging.NewLogger("test", sharedConfig.LogConfig{})
	provider := &EmailProvider{
		config: config,
		logger: logger,
	}

	name := provider.GetName()
	if name != "email" {
		t.Errorf("Expected name 'email', got %s", name)
	}
}

func TestEmailProvider_IsEnabled(t *testing.T) {
	config := &EmailConfig{
		Provider:     "smtp",
		APIKey:       "test-key",
		FromEmail:    "test@example.com",
		FromName:     "Test Sender",
		TemplatePath: "./templates",
		RateLimit:    100,
	}

	logger := logging.NewLogger("test", sharedConfig.LogConfig{})
	provider := &EmailProvider{
		config: config,
		logger: logger,
	}

	enabled := provider.IsEnabled()
	if !enabled {
		t.Error("Expected provider to be enabled")
	}
}

func TestEmailProvider_IsEnabledWithEmptyAPIKey(t *testing.T) {
	config := &EmailConfig{
		Provider:     "smtp",
		APIKey:       "", // Empty API key
		FromEmail:    "test@example.com",
		FromName:     "Test Sender",
		TemplatePath: "./templates",
		RateLimit:    100,
	}

	logger := logging.NewLogger("test", sharedConfig.LogConfig{})
	provider := &EmailProvider{
		config: config,
		logger: logger,
	}

	enabled := provider.IsEnabled()
	if enabled {
		t.Error("Expected provider to be disabled with empty API key")
	}
}

func TestEmailProvider_validateEmailMessage(t *testing.T) {
	config := &EmailConfig{
		Provider:     "smtp",
		APIKey:       "test-key",
		FromEmail:    "test@example.com",
		FromName:     "Test Sender",
		TemplatePath: "./templates",
		RateLimit:    100,
	}

	logger := logging.NewLogger("test", sharedConfig.LogConfig{})
	provider := &EmailProvider{
		config: config,
		logger: logger,
	}

	// Test with nil message
	err := provider.validateEmailMessage(nil)
	if err == nil {
		t.Error("Expected error for nil message, got nil")
	}

	// Test with valid message
	message := &NotificationMessage{
		ID:        "test-message-1",
		UserID:    "user-123",
		Channel:   "email",
		Type:      "welcome",
		Title:     "Welcome!",
		Content:   "Welcome to our platform!",
		Recipient: "user@example.com",
		Subject:   "Welcome to USC Platform",
		Data:      make(map[string]interface{}),
		Priority:  "normal",
		CreatedAt: time.Now(),
	}

	err = provider.validateEmailMessage(message)
	if err != nil {
		t.Errorf("Expected no error for valid message, got %v", err)
	}
}

func TestEmailProvider_isValidEmail(t *testing.T) {
	config := &EmailConfig{
		Provider:     "smtp",
		APIKey:       "test-key",
		FromEmail:    "test@example.com",
		FromName:     "Test Sender",
		TemplatePath: "./templates",
		RateLimit:    100,
	}

	logger := logging.NewLogger("test", sharedConfig.LogConfig{})
	provider := &EmailProvider{
		config: config,
		logger: logger,
	}

	// Test valid emails
	validEmails := []string{
		"user@example.com",
		"test.user@domain.co.uk",
		"user+tag@example.org",
	}

	for _, email := range validEmails {
		if !provider.isValidEmail(email) {
			t.Errorf("Expected email %s to be valid", email)
		}
	}

	// Test invalid emails
	invalidEmails := []string{
		"invalid-email",
		"user@",
		"user.example.com",
		"",
	}

	for _, email := range invalidEmails {
		if provider.isValidEmail(email) {
			t.Errorf("Expected email %s to be invalid", email)
		}
	}

	// Test edge case - email with @ and . but no local part
	if provider.isValidEmail("@example.com") {
		t.Errorf("Expected email @example.com to be invalid")
	}
}

func TestEmailProvider_sendViaSendGrid(t *testing.T) {
	config := &EmailConfig{
		Provider:     "sendgrid",
		APIKey:       "test-key",
		FromEmail:    "test@example.com",
		FromName:     "Test Sender",
		TemplatePath: "./templates",
		RateLimit:    100,
	}

	logger := logging.NewLogger("test", sharedConfig.LogConfig{})
	provider := &EmailProvider{
		config: config,
		logger: logger,
	}

	message := &NotificationMessage{
		ID:        "test-message-1",
		UserID:    "user-123",
		Channel:   "email",
		Type:      "welcome",
		Title:     "Welcome!",
		Content:   "Welcome to our platform!",
		Recipient: "user@example.com",
		Subject:   "Welcome to USC Platform",
		Data:      make(map[string]interface{}),
		Priority:  "normal",
		CreatedAt: time.Now(),
	}

	delivery, err := provider.sendViaSendGrid(message)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if delivery == nil {
		t.Fatal("Expected delivery to be returned")
	}

	if delivery.Channel != "email" {
		t.Errorf("Expected channel 'email', got %s", delivery.Channel)
	}
}

func TestEmailProvider_sendViaSES(t *testing.T) {
	config := &EmailConfig{
		Provider:     "ses",
		APIKey:       "test-key",
		FromEmail:    "test@example.com",
		FromName:     "Test Sender",
		TemplatePath: "./templates",
		RateLimit:    100,
	}

	logger := logging.NewLogger("test", sharedConfig.LogConfig{})
	provider := &EmailProvider{
		config: config,
		logger: logger,
	}

	message := &NotificationMessage{
		ID:        "test-message-1",
		UserID:    "user-123",
		Channel:   "email",
		Type:      "welcome",
		Title:     "Welcome!",
		Content:   "Welcome to our platform!",
		Recipient: "user@example.com",
		Subject:   "Welcome to USC Platform",
		Data:      make(map[string]interface{}),
		Priority:  "normal",
		CreatedAt: time.Now(),
	}

	delivery, err := provider.sendViaSES(message)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if delivery == nil {
		t.Fatal("Expected delivery to be returned")
	}

	if delivery.Channel != "email" {
		t.Errorf("Expected channel 'email', got %s", delivery.Channel)
	}
}

func TestEmailProvider_sendViaSMTP(t *testing.T) {
	config := &EmailConfig{
		Provider:     "smtp",
		APIKey:       "test-key",
		FromEmail:    "test@example.com",
		FromName:     "Test Sender",
		TemplatePath: "./templates",
		RateLimit:    100,
	}

	logger := logging.NewLogger("test", sharedConfig.LogConfig{})
	provider := &EmailProvider{
		config: config,
		logger: logger,
	}

	message := &NotificationMessage{
		ID:        "test-message-1",
		UserID:    "user-123",
		Channel:   "email",
		Type:      "welcome",
		Title:     "Welcome!",
		Content:   "Welcome to our platform!",
		Recipient: "user@example.com",
		Subject:   "Welcome to USC Platform",
		Data:      make(map[string]interface{}),
		Priority:  "normal",
		CreatedAt: time.Now(),
	}

	delivery, err := provider.sendViaSMTP(message)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if delivery == nil {
		t.Fatal("Expected delivery to be returned")
	}

	if delivery.Channel != "email" {
		t.Errorf("Expected channel 'email', got %s", delivery.Channel)
	}
}

func TestNotificationChannelService_ConcurrentAccess(t *testing.T) {
	logger := logging.NewLogger("test", sharedConfig.LogConfig{})
	config := &ChannelConfig{
		EmailConfig: &EmailConfig{
			Provider:     "smtp",
			APIKey:       "test-key",
			FromEmail:    "test@example.com",
			FromName:     "Test Sender",
			TemplatePath: "./templates",
			RateLimit:    100,
		},
		DefaultChannel: "email",
		RetryAttempts:  3,
		RetryDelay:     5 * time.Second,
	}

	service := NewNotificationChannelService(logger, config)

	// Test concurrent access
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			message := &NotificationMessage{
				ID:        fmt.Sprintf("test-message-%d", i),
				UserID:    fmt.Sprintf("user-%d", i),
				Channel:   "email",
				Type:      "welcome",
				Title:     "Welcome!",
				Content:   "Welcome to our platform!",
				Recipient: fmt.Sprintf("user%d@example.com", i),
				Subject:   "Welcome to USC Platform",
				Data:      make(map[string]interface{}),
				Priority:  "normal",
				CreatedAt: time.Now(),
			}

			ctx := context.Background()
			_, err := service.SendNotification(ctx, message)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// Benchmark tests
func BenchmarkNotificationChannelService_SendNotification(b *testing.B) {
	logger := logging.NewLogger("test", sharedConfig.LogConfig{})
	config := &ChannelConfig{
		EmailConfig: &EmailConfig{
			Provider:     "smtp",
			APIKey:       "test-key",
			FromEmail:    "test@example.com",
			FromName:     "Test Sender",
			TemplatePath: "./templates",
			RateLimit:    100,
		},
		DefaultChannel: "email",
		RetryAttempts:  3,
		RetryDelay:     5 * time.Second,
	}

	service := NewNotificationChannelService(logger, config)

	message := &NotificationMessage{
		ID:        "test-message-1",
		UserID:    "user-123",
		Channel:   "email",
		Type:      "welcome",
		Title:     "Welcome!",
		Content:   "Welcome to our platform!",
		Recipient: "user@example.com",
		Subject:   "Welcome to USC Platform",
		Data:      make(map[string]interface{}),
		Priority:  "normal",
		CreatedAt: time.Now(),
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.SendNotification(ctx, message)
	}
}

func BenchmarkEmailProvider_Send(b *testing.B) {
	config := &EmailConfig{
		Provider:     "smtp",
		APIKey:       "test-key",
		FromEmail:    "test@example.com",
		FromName:     "Test Sender",
		TemplatePath: "./templates",
		RateLimit:    100,
	}

	logger := logging.NewLogger("test", sharedConfig.LogConfig{})
	provider := &EmailProvider{
		config: config,
		logger: logger,
	}

	message := &NotificationMessage{
		ID:        "test-message-1",
		UserID:    "user-123",
		Channel:   "email",
		Type:      "welcome",
		Title:     "Welcome!",
		Content:   "Welcome to our platform!",
		Recipient: "user@example.com",
		Subject:   "Welcome to USC Platform",
		Data:      make(map[string]interface{}),
		Priority:  "normal",
		CreatedAt: time.Now(),
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		provider.Send(ctx, message)
	}
}

func BenchmarkNotificationChannelService_Concurrent(b *testing.B) {
	logger := logging.NewLogger("test", sharedConfig.LogConfig{})
	config := &ChannelConfig{
		EmailConfig: &EmailConfig{
			Provider:     "smtp",
			APIKey:       "test-key",
			FromEmail:    "test@example.com",
			FromName:     "Test Sender",
			TemplatePath: "./templates",
			RateLimit:    100,
		},
		DefaultChannel: "email",
		RetryAttempts:  3,
		RetryDelay:     5 * time.Second,
	}

	service := NewNotificationChannelService(logger, config)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			message := &NotificationMessage{
				ID:        fmt.Sprintf("test-message-%d", i),
				UserID:    fmt.Sprintf("user-%d", i),
				Channel:   "email",
				Type:      "welcome",
				Title:     "Welcome!",
				Content:   "Welcome to our platform!",
				Recipient: fmt.Sprintf("user%d@example.com", i),
				Subject:   "Welcome to USC Platform",
				Data:      make(map[string]interface{}),
				Priority:  "normal",
				CreatedAt: time.Now(),
			}

			ctx := context.Background()
			service.SendNotification(ctx, message)
			i++
		}
	})
}
