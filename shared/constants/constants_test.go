package constants

import (
	"testing"
)

func TestConfigConstants(t *testing.T) {
	// Test server configuration constants
	if ConfigServerHost == "" {
		t.Error("ConfigServerHost should not be empty")
	}

	if ConfigServerPort == "" {
		t.Error("ConfigServerPort should not be empty")
	}

	if ConfigServerTimeout == "" {
		t.Error("ConfigServerTimeout should not be empty")
	}

	// Test database configuration constants
	if ConfigDatabaseHost == "" {
		t.Error("ConfigDatabaseHost should not be empty")
	}

	if ConfigDatabasePort == "" {
		t.Error("ConfigDatabasePort should not be empty")
	}

	if ConfigDatabaseName == "" {
		t.Error("ConfigDatabaseName should not be empty")
	}

	// Test Redis configuration constants
	if ConfigRedisHost == "" {
		t.Error("ConfigRedisHost should not be empty")
	}

	if ConfigRedisPort == "" {
		t.Error("ConfigRedisPort should not be empty")
	}

	if ConfigRedisPassword == "" {
		t.Error("ConfigRedisPassword should not be empty")
	}
}

// TestErrorConstants removed - error codes moved to errors package

func TestCleanServiceConstants(t *testing.T) {
	// Test service name constants (clean version - matches actual services.go)
	if ServiceGateway == "" {
		t.Error("ServiceGateway should not be empty")
	}

	if ServiceAuth == "" {
		t.Error("ServiceAuth should not be empty")
	}

	if ServiceUser == "" {
		t.Error("ServiceUser should not be empty")
	}

	if ServiceBlockchainCore == "" {
		t.Error("ServiceBlockchainCore should not be empty")
	}

	if ServiceWallet == "" {
		t.Error("ServiceWallet should not be empty")
	}

	if ServiceSecurity == "" {
		t.Error("ServiceSecurity should not be empty")
	}

	if ServiceCaching == "" {
		t.Error("ServiceCaching should not be empty")
	}

	if ServiceMonitoring == "" {
		t.Error("ServiceMonitoring should not be empty")
	}

	if ServiceSocial == "" {
		t.Error("ServiceSocial should not be empty")
	}

	if ServiceBilateralRewards == "" {
		t.Error("ServiceBilateralRewards should not be empty")
	}

	if ServiceContentManagement == "" {
		t.Error("ServiceContentManagement should not be empty")
	}

	if ServiceVideo == "" {
		t.Error("ServiceVideo should not be empty")
	}

	if ServiceAI == "" {
		t.Error("ServiceAI should not be empty")
	}

	if ServiceCommerce == "" {
		t.Error("ServiceCommerce should not be empty")
	}

	if ServiceNotification == "" {
		t.Error("ServiceNotification should not be empty")
	}

	if ServiceSearch == "" {
		t.Error("ServiceSearch should not be empty")
	}

	if ServiceAnalytics == "" {
		t.Error("ServiceAnalytics should not be empty")
	}

	if ServiceModeration == "" {
		t.Error("ServiceModeration should not be empty")
	}

	if ServiceRecommendation == "" {
		t.Error("ServiceRecommendation should not be empty")
	}

	if ServiceAdvertising == "" {
		t.Error("ServiceAdvertising should not be empty")
	}

	if ServiceAdmin == "" {
		t.Error("ServiceAdmin should not be empty")
	}

	if ServiceKafkaMessaging == "" {
		t.Error("ServiceKafkaMessaging should not be empty")
	}
}

func TestCleanPortConstants(t *testing.T) {
	// Test port constants (clean version - matches actual ports.go)
	if PortGateway == 0 {
		t.Error("PortGateway should not be zero")
	}

	if PortAuth == 0 {
		t.Error("PortAuth should not be zero")
	}

	if PortUser == 0 {
		t.Error("PortUser should not be zero")
	}

	if PortBlockchainCore == 0 {
		t.Error("PortBlockchainCore should not be zero")
	}

	if PortWallet == 0 {
		t.Error("PortWallet should not be zero")
	}

	if PortSecurity == 0 {
		t.Error("PortSecurity should not be zero")
	}

	if PortCaching == 0 {
		t.Error("PortCaching should not be zero")
	}

	if PortMonitoring == 0 {
		t.Error("PortMonitoring should not be zero")
	}

	if PortSocial == 0 {
		t.Error("PortSocial should not be zero")
	}

	if PortBilateralRewards == 0 {
		t.Error("PortBilateralRewards should not be zero")
	}

	if PortContentManagement == 0 {
		t.Error("PortContentManagement should not be zero")
	}

	if PortVideo == 0 {
		t.Error("PortVideo should not be zero")
	}

	if PortAI == 0 {
		t.Error("PortAI should not be zero")
	}

	if PortCommerce == 0 {
		t.Error("PortCommerce should not be zero")
	}

	if PortNotification == 0 {
		t.Error("PortNotification should not be zero")
	}

	if PortSearch == 0 {
		t.Error("PortSearch should not be zero")
	}

	if PortAnalytics == 0 {
		t.Error("PortAnalytics should not be zero")
	}

	if PortModeration == 0 {
		t.Error("PortModeration should not be zero")
	}

	if PortRecommendation == 0 {
		t.Error("PortRecommendation should not be zero")
	}

	if PortAdvertising == 0 {
		t.Error("PortAdvertising should not be zero")
	}

	if PortAdmin == 0 {
		t.Error("PortAdmin should not be zero")
	}

	if PortKafkaMessaging == 0 {
		t.Error("PortKafkaMessaging should not be zero")
	}
}

func TestCategoryConstants(t *testing.T) {
	// Test category constants
	if CategoryCore == "" {
		t.Error("CategoryCore should not be empty")
	}

	if CategoryUser == "" {
		t.Error("CategoryUser should not be empty")
	}

	if CategoryBlockchain == "" {
		t.Error("CategoryBlockchain should not be empty")
	}

	if CategoryTransaction == "" {
		t.Error("CategoryTransaction should not be empty")
	}

	if CategoryNotification == "" {
		t.Error("CategoryNotification should not be empty")
	}
}
