package errors

import (
	"testing"
)

func TestUSCErrorCodes(t *testing.T) {
	t.Run("USC_Blockchain_Error_Codes", func(t *testing.T) {
		// Test USC Token errors
		if ErrCodeUSCInsufficientBalance != "USC_INSUFFICIENT_BALANCE" {
			t.Errorf("Expected 'USC_INSUFFICIENT_BALANCE', got '%s'", ErrCodeUSCInsufficientBalance)
		}

		if ErrCodeUSCInvalidAmount != "USC_INVALID_AMOUNT" {
			t.Errorf("Expected 'USC_INVALID_AMOUNT', got '%s'", ErrCodeUSCInvalidAmount)
		}

		if ErrCodeUSCTransferFailed != "USC_TRANSFER_FAILED" {
			t.Errorf("Expected 'USC_TRANSFER_FAILED', got '%s'", ErrCodeUSCTransferFailed)
		}
	})

	t.Run("Wallet_Error_Codes", func(t *testing.T) {
		if ErrCodeWalletNotFound != "WALLET_NOT_FOUND" {
			t.Errorf("Expected 'WALLET_NOT_FOUND', got '%s'", ErrCodeWalletNotFound)
		}

		if ErrCodeWalletCreationFailed != "WALLET_CREATION_FAILED" {
			t.Errorf("Expected 'WALLET_CREATION_FAILED', got '%s'", ErrCodeWalletCreationFailed)
		}

		if ErrCodeChatMessageEmpty != "CHAT_MESSAGE_EMPTY" {
			t.Errorf("Expected 'CHAT_MESSAGE_EMPTY', got '%s'", ErrCodeChatMessageEmpty)
		}
	})

	t.Run("NFT_Error_Codes", func(t *testing.T) {
		if ErrCodeNFTNotFound != "NFT_NOT_FOUND" {
			t.Errorf("Expected 'NFT_NOT_FOUND', got '%s'", ErrCodeNFTNotFound)
		}

		if ErrCodeNFTMintingFailed != "NFT_MINTING_FAILED" {
			t.Errorf("Expected 'NFT_MINTING_FAILED', got '%s'", ErrCodeNFTMintingFailed)
		}

		if ErrCodeNFTMarketplaceListed != "NFT_MARKETPLACE_LISTED" {
			t.Errorf("Expected 'NFT_MARKETPLACE_LISTED', got '%s'", ErrCodeNFTMarketplaceListed)
		}
	})

	t.Run("Service_Error_Codes", func(t *testing.T) {
		// Gateway errors
		if ErrCodeGatewayTimeout != "GATEWAY_TIMEOUT" {
			t.Errorf("Expected 'GATEWAY_TIMEOUT', got '%s'", ErrCodeGatewayTimeout)
		}

		// Auth errors
		if ErrCodeMFARequired != "MFA_REQUIRED" {
			t.Errorf("Expected 'MFA_REQUIRED', got '%s'", ErrCodeMFARequired)
		}

		// Social errors
		if ErrCodePostNotFound != "POST_NOT_FOUND" {
			t.Errorf("Expected 'POST_NOT_FOUND', got '%s'", ErrCodePostNotFound)
		}

		// Video errors
		if ErrCodeVideoNotFound != "VIDEO_NOT_FOUND" {
			t.Errorf("Expected 'VIDEO_NOT_FOUND', got '%s'", ErrCodeVideoNotFound)
		}

		// Commerce errors
		if ErrCodeProductNotFound != "PRODUCT_NOT_FOUND" {
			t.Errorf("Expected 'PRODUCT_NOT_FOUND', got '%s'", ErrCodeProductNotFound)
		}

		// AI errors
		if ErrCodeAIModelNotFound != "AI_MODEL_NOT_FOUND" {
			t.Errorf("Expected 'AI_MODEL_NOT_FOUND', got '%s'", ErrCodeAIModelNotFound)
		}

		// Kafka errors
		if ErrCodeKafkaConnectionFailed != "KAFKA_CONNECTION_FAILED" {
			t.Errorf("Expected 'KAFKA_CONNECTION_FAILED', got '%s'", ErrCodeKafkaConnectionFailed)
		}
	})
}

func TestUSCErrorCategories(t *testing.T) {
	t.Run("Category_Constants", func(t *testing.T) {
		if CategoryUSCBlockchain != "usc_blockchain" {
			t.Errorf("Expected 'usc_blockchain', got '%s'", CategoryUSCBlockchain)
		}

		if CategoryUSCWallet != "usc_wallet" {
			t.Errorf("Expected 'usc_wallet', got '%s'", CategoryUSCWallet)
		}

		if CategoryUSCNFT != "usc_nft" {
			t.Errorf("Expected 'usc_nft', got '%s'", CategoryUSCNFT)
		}

		if CategorySocial != "social" {
			t.Errorf("Expected 'social', got '%s'", CategorySocial)
		}

		if CategoryVideo != "video" {
			t.Errorf("Expected 'video', got '%s'", CategoryVideo)
		}

		if CategoryCommerce != "commerce" {
			t.Errorf("Expected 'commerce', got '%s'", CategoryCommerce)
		}

		if CategoryAI != "ai" {
			t.Errorf("Expected 'ai', got '%s'", CategoryAI)
		}
	})
}

func TestGetUSCErrorCodesByService(t *testing.T) {
	t.Run("Service_Error_Mapping", func(t *testing.T) {
		serviceErrors := GetUSCErrorCodesByService()

		// Test that all 22 services are included
		expectedServices := []string{
			"service-01-gateway",
			"service-02-auth",
			"service-03-user",
			"service-04-usc-blockchain-core",
			"service-05-usc-wallet",
			"service-06-security",
			"service-07-caching",
			"service-08-monitoring",
			"service-09-social",
			"service-10-usc-bilateral-rewards",
			"service-11-content-management",
			"service-12-video-service",
			"service-13-ai-service",
			"service-14-commerce-service",
			"service-15-notification-service",
			"service-16-search-service",
			"service-17-analytics-service",
			"service-18-moderation-service",
			"service-19-recommendation-service",
			"service-20-advertising-service",
			"service-21-admin-service",
			"service-22-kafka-messaging-service",
		}

		if len(serviceErrors) != len(expectedServices) {
			t.Errorf("Expected %d services, got %d", len(expectedServices), len(serviceErrors))
		}

		for _, service := range expectedServices {
			if _, exists := serviceErrors[service]; !exists {
				t.Errorf("Service '%s' not found in error mapping", service)
			}
		}

		// Test specific service error codes
		gatewayErrors := serviceErrors["service-01-gateway"]
		if len(gatewayErrors) == 0 {
			t.Error("Gateway service should have error codes")
		}

		// Check if specific error is in gateway errors
		found := false
		for _, err := range gatewayErrors {
			if err == ErrCodeGatewayTimeout {
				found = true
				break
			}
		}
		if !found {
			t.Error("ErrCodeGatewayTimeout should be in gateway service errors")
		}

		// Test blockchain service errors
		blockchainErrors := serviceErrors["service-04-usc-blockchain-core"]
		if len(blockchainErrors) == 0 {
			t.Error("Blockchain service should have error codes")
		}

		// Check USC specific errors
		found = false
		for _, err := range blockchainErrors {
			if err == ErrCodeUSCInsufficientBalance {
				found = true
				break
			}
		}
		if !found {
			t.Error("ErrCodeUSCInsufficientBalance should be in blockchain service errors")
		}
	})
}

func TestIsUSCSpecificError(t *testing.T) {
	t.Run("USC_Specific_Detection", func(t *testing.T) {
		// Test USC-specific errors
		uscErrors := []ErrorCode{
			ErrCodeUSCInsufficientBalance,
			ErrCodeUSCInvalidAmount,
			ErrCodeUSCTransferFailed,
			ErrCodeBlockchainConnectionFailed,
			ErrCodeWalletNotFound,
			ErrCodeStakingInsufficientAmount,
			ErrCodeNFTNotFound,
			ErrCodeRewardCalculationFailed,
		}

		for _, err := range uscErrors {
			if !IsUSCSpecificError(err) {
				t.Errorf("Error code '%s' should be detected as USC-specific", err)
			}
		}

		// Test non-USC errors
		nonUSCErrors := []ErrorCode{
			ErrCodeGatewayTimeout,
			ErrCodeMFARequired,
			ErrCodePostNotFound,
			ErrCodeVideoNotFound,
		}

		for _, err := range nonUSCErrors {
			if IsUSCSpecificError(err) {
				t.Errorf("Error code '%s' should not be detected as USC-specific", err)
			}
		}
	})
}

func TestComprehensiveErrorCodeCoverage(t *testing.T) {
	t.Run("All_Services_Have_Errors", func(t *testing.T) {
		serviceErrors := GetUSCErrorCodesByService()

		// Check that each service has at least 3 error codes
		minErrorsPerService := 3
		for serviceName, errors := range serviceErrors {
			if len(errors) < minErrorsPerService {
				t.Errorf("Service '%s' should have at least %d error codes, got %d",
					serviceName, minErrorsPerService, len(errors))
			}
		}
	})

	t.Run("Error_Code_Uniqueness", func(t *testing.T) {
		serviceErrors := GetUSCErrorCodesByService()
		allErrors := make(map[ErrorCode]bool)

		// Collect all error codes
		for _, errors := range serviceErrors {
			for _, err := range errors {
				if allErrors[err] {
					t.Errorf("Duplicate error code found: %s", err)
				}
				allErrors[err] = true
			}
		}

		// Verify we have sufficient error coverage
		if len(allErrors) < 50 {
			t.Errorf("Expected at least 50 unique error codes, got %d", len(allErrors))
		}
	})
}

func BenchmarkGetUSCErrorCodesByService(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GetUSCErrorCodesByService()
	}
}

func BenchmarkIsUSCSpecificError(b *testing.B) {
	testErr := ErrCodeUSCInsufficientBalance

	for i := 0; i < b.N; i++ {
		_ = IsUSCSpecificError(testErr)
	}
}
