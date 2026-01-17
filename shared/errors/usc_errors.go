// Package errors provides USC-specific error codes for the USC platform services.
package errors

// USC-Specific Error Codes
const (
	// =====================
	// USC BLOCKCHAIN ERRORS
	// =====================

	// USC Token Errors
	ErrCodeUSCInsufficientBalance ErrorCode = "USC_INSUFFICIENT_BALANCE"
	ErrCodeUSCInvalidAmount       ErrorCode = "USC_INVALID_AMOUNT"
	ErrCodeUSCTransferFailed      ErrorCode = "USC_TRANSFER_FAILED"
	ErrCodeUSCInvalidAddress      ErrorCode = "USC_INVALID_ADDRESS"
	ErrCodeUSCTransactionPending  ErrorCode = "USC_TRANSACTION_PENDING"
	ErrCodeUSCTransactionFailed   ErrorCode = "USC_TRANSACTION_FAILED"
	ErrCodeUSCTransactionTimeout  ErrorCode = "USC_TRANSACTION_TIMEOUT"
	ErrCodeUSCGasEstimationFailed ErrorCode = "USC_GAS_ESTIMATION_FAILED"
	ErrCodeUSCNonceMismatch       ErrorCode = "USC_NONCE_MISMATCH"

	// USC Blockchain Core Errors
	ErrCodeBlockchainConnectionFailed ErrorCode = "BLOCKCHAIN_CONNECTION_FAILED"
	ErrCodeBlockchainSyncFailed       ErrorCode = "BLOCKCHAIN_SYNC_FAILED"
	ErrCodeBlockValidationFailed      ErrorCode = "BLOCK_VALIDATION_FAILED"
	ErrCodeConsensusFailure           ErrorCode = "CONSENSUS_FAILURE"
	ErrCodeInvalidBlockHeight         ErrorCode = "INVALID_BLOCK_HEIGHT"
	ErrCodeChainReorganization        ErrorCode = "CHAIN_REORGANIZATION"

	// USC Wallet Errors
	ErrCodeWalletNotFound            ErrorCode = "WALLET_NOT_FOUND"
	ErrCodeWalletCreationFailed      ErrorCode = "WALLET_CREATION_FAILED"
	ErrCodeWalletUnlockFailed        ErrorCode = "WALLET_UNLOCK_FAILED"
	ErrCodeWalletBackupFailed        ErrorCode = "WALLET_BACKUP_FAILED"
	ErrCodeWalletRestoreFailed       ErrorCode = "WALLET_RESTORE_FAILED"
	ErrCodeInvalidMnemonic           ErrorCode = "INVALID_MNEMONIC"
	ErrCodeInvalidPrivateKey         ErrorCode = "INVALID_PRIVATE_KEY"
	ErrCodeWalletChatMessageTooLong  ErrorCode = "WALLET_CHAT_MESSAGE_TOO_LONG"
	ErrCodeWalletChatHistoryExceeded ErrorCode = "WALLET_CHAT_HISTORY_EXCEEDED"

	// USC Staking Errors
	ErrCodeStakingInsufficientAmount      ErrorCode = "STAKING_INSUFFICIENT_AMOUNT"
	ErrCodeStakingPeriodInvalid           ErrorCode = "STAKING_PERIOD_INVALID"
	ErrCodeStakingAlreadyActive           ErrorCode = "STAKING_ALREADY_ACTIVE"
	ErrCodeStakingNotActive               ErrorCode = "STAKING_NOT_ACTIVE"
	ErrCodeStakingUnlockNotReady          ErrorCode = "STAKING_UNLOCK_NOT_READY"
	ErrCodeStakingRewardCalculationFailed ErrorCode = "STAKING_REWARD_CALCULATION_FAILED"

	// USC NFT Errors
	ErrCodeNFTNotFound             ErrorCode = "NFT_NOT_FOUND"
	ErrCodeNFTMintingFailed        ErrorCode = "NFT_MINTING_FAILED"
	ErrCodeNFTTransferFailed       ErrorCode = "NFT_TRANSFER_FAILED"
	ErrCodeNFTInvalidMetadata      ErrorCode = "NFT_INVALID_METADATA"
	ErrCodeNFTMaxSupplyExceeded    ErrorCode = "NFT_MAX_SUPPLY_EXCEEDED"
	ErrCodeNFTInvalidRoyalty       ErrorCode = "NFT_INVALID_ROYALTY"
	ErrCodeNFTMarketplaceListed    ErrorCode = "NFT_MARKETPLACE_LISTED"
	ErrCodeNFTMarketplaceNotListed ErrorCode = "NFT_MARKETPLACE_NOT_LISTED"

	// ============================
	// SERVICE-SPECIFIC ERROR CODES
	// ============================

	// Service-01-Gateway Errors
	ErrCodeGatewayTimeout               ErrorCode = "GATEWAY_TIMEOUT"
	ErrCodeGraphQLParsingFailed         ErrorCode = "GRAPHQL_PARSING_FAILED"
	ErrCodeGraphQLValidationFailed      ErrorCode = "GRAPHQL_VALIDATION_FAILED"
	ErrCodeFederationServiceUnavailable ErrorCode = "FEDERATION_SERVICE_UNAVAILABLE"
	ErrCodeQueryComplexityExceeded      ErrorCode = "QUERY_COMPLEXITY_EXCEEDED"
	ErrCodeSubscriptionLimitExceeded    ErrorCode = "SUBSCRIPTION_LIMIT_EXCEEDED"

	// Service-02-Auth Errors
	ErrCodeMFARequired               ErrorCode = "MFA_REQUIRED"
	ErrCodeMFAInvalid                ErrorCode = "MFA_INVALID"
	ErrCodeMFASetupRequired          ErrorCode = "MFA_SETUP_REQUIRED"
	ErrCodePasswordTooWeak           ErrorCode = "PASSWORD_TOO_WEAK"
	ErrCodePasswordExpired           ErrorCode = "PASSWORD_EXPIRED"
	ErrCodeSessionExpired            ErrorCode = "SESSION_EXPIRED"
	ErrCodeOAuthProviderError        ErrorCode = "OAUTH_PROVIDER_ERROR"
	ErrCodeEmailVerificationRequired ErrorCode = "EMAIL_VERIFICATION_REQUIRED"

	// Service-03-User Errors
	ErrCodeUserProfileIncomplete    ErrorCode = "USER_PROFILE_INCOMPLETE"
	ErrCodeUserRelationshipExists   ErrorCode = "USER_RELATIONSHIP_EXISTS"
	ErrCodeUserRelationshipNotFound ErrorCode = "USER_RELATIONSHIP_NOT_FOUND"
	ErrCodeUserPreferencesInvalid   ErrorCode = "USER_PREFERENCES_INVALID"
	ErrCodeUserSuspended            ErrorCode = "USER_SUSPENDED"
	ErrCodeUserDeactivated          ErrorCode = "USER_DEACTIVATED"

	// Service-05-Wallet Chat Errors
	ErrCodeChatMessageEmpty         ErrorCode = "CHAT_MESSAGE_EMPTY"
	ErrCodeChatParticipantNotFound  ErrorCode = "CHAT_PARTICIPANT_NOT_FOUND"
	ErrCodeChatConversationNotFound ErrorCode = "CHAT_CONVERSATION_NOT_FOUND"
	ErrCodeChatMessageNotFound      ErrorCode = "CHAT_MESSAGE_NOT_FOUND"
	ErrCodeChatEncryptionFailed     ErrorCode = "CHAT_ENCRYPTION_FAILED"
	ErrCodeChatDecryptionFailed     ErrorCode = "CHAT_DECRYPTION_FAILED"

	// Service-06-Security Errors
	ErrCodeSecurityThreatDetected  ErrorCode = "SECURITY_THREAT_DETECTED"
	ErrCodeSecurityRuleViolation   ErrorCode = "SECURITY_RULE_VIOLATION"
	ErrCodeSecurityScanFailed      ErrorCode = "SECURITY_SCAN_FAILED"
	ErrCodeFraudDetected           ErrorCode = "FRAUD_DETECTED"
	ErrCodeSuspiciousActivity      ErrorCode = "SUSPICIOUS_ACTIVITY"
	ErrCodeSecurityPolicyViolation ErrorCode = "SECURITY_POLICY_VIOLATION"
	ErrCodeIPBlocked               ErrorCode = "IP_BLOCKED"
	ErrCodeGeoLocationBlocked      ErrorCode = "GEO_LOCATION_BLOCKED"

	// Service-07-Caching Errors
	ErrCodeCacheConnectionFailed      ErrorCode = "CACHE_CONNECTION_FAILED"
	ErrCodeCacheKeyNotFound           ErrorCode = "CACHE_KEY_NOT_FOUND"
	ErrCodeCacheSerializationFailed   ErrorCode = "CACHE_SERIALIZATION_FAILED"
	ErrCodeCacheDeserializationFailed ErrorCode = "CACHE_DESERIALIZATION_FAILED"
	ErrCodeCacheEvictionFailed        ErrorCode = "CACHE_EVICTION_FAILED"
	ErrCodeCacheMemoryExhausted       ErrorCode = "CACHE_MEMORY_EXHAUSTED"

	// Service-08-Monitoring Errors
	ErrCodeMetricCollectionFailed    ErrorCode = "METRIC_COLLECTION_FAILED"
	ErrCodeAlertConfigurationInvalid ErrorCode = "ALERT_CONFIGURATION_INVALID"
	ErrCodeMonitoringServiceDown     ErrorCode = "MONITORING_SERVICE_DOWN"
	ErrCodeHealthCheckFailed         ErrorCode = "HEALTH_CHECK_FAILED"
	ErrCodeMetricThresholdExceeded   ErrorCode = "METRIC_THRESHOLD_EXCEEDED"

	// Service-09-Social Errors
	ErrCodePostNotFound            ErrorCode = "POST_NOT_FOUND"
	ErrCodePostCreationFailed      ErrorCode = "POST_CREATION_FAILED"
	ErrCodePostContentInvalid      ErrorCode = "POST_CONTENT_INVALID"
	ErrCodeSocialInteractionFailed ErrorCode = "SOCIAL_INTERACTION_FAILED"
	ErrCodeFeedGenerationFailed    ErrorCode = "FEED_GENERATION_FAILED"
	ErrCodeContentModerationFailed ErrorCode = "CONTENT_MODERATION_FAILED"
	ErrCodeSocialLimitExceeded     ErrorCode = "SOCIAL_LIMIT_EXCEEDED"

	// Service-10-USC-Bilateral-Rewards Errors
	ErrCodeRewardCalculationFailed  ErrorCode = "REWARD_CALCULATION_FAILED"
	ErrCodeRewardDistributionFailed ErrorCode = "REWARD_DISTRIBUTION_FAILED"
	ErrCodeRewardAlreadyClaimed     ErrorCode = "REWARD_ALREADY_CLAIMED"
	ErrCodeRewardExpired            ErrorCode = "REWARD_EXPIRED"
	ErrCodeRewardEligibilityFailed  ErrorCode = "REWARD_ELIGIBILITY_FAILED"
	ErrCodeRewardPoolExhausted      ErrorCode = "REWARD_POOL_EXHAUSTED"

	// Service-11-Content-Management Errors
	// Note: ErrCodeFileUploadFailed and ErrCodeFileSizeExceeded moved to domain.go
	ErrCodeFileFormatUnsupported   ErrorCode = "FILE_FORMAT_UNSUPPORTED"
	ErrCodeContentProcessingFailed ErrorCode = "CONTENT_PROCESSING_FAILED"
	ErrCodeCDNDeploymentFailed     ErrorCode = "CDN_DEPLOYMENT_FAILED"
	ErrCodeMediaTranscodingFailed  ErrorCode = "MEDIA_TRANSCODING_FAILED"

	// Service-12-Video-Service Errors
	ErrCodeVideoNotFound             ErrorCode = "VIDEO_NOT_FOUND"
	ErrCodeVideoUploadFailed         ErrorCode = "VIDEO_UPLOAD_FAILED"
	ErrCodeVideoProcessingFailed     ErrorCode = "VIDEO_PROCESSING_FAILED"
	ErrCodeVideoStreamingFailed      ErrorCode = "VIDEO_STREAMING_FAILED"
	ErrCodeLiveStreamNotFound        ErrorCode = "LIVE_STREAM_NOT_FOUND"
	ErrCodeLiveStreamStartFailed     ErrorCode = "LIVE_STREAM_START_FAILED"
	ErrCodeLiveStreamStopFailed      ErrorCode = "LIVE_STREAM_STOP_FAILED"
	ErrCodeVideoQualityUnsupported   ErrorCode = "VIDEO_QUALITY_UNSUPPORTED"
	ErrCodeVideoChatConnectionFailed ErrorCode = "VIDEO_CHAT_CONNECTION_FAILED"

	// Service-13-AI-Service Errors
	ErrCodeAIModelNotFound           ErrorCode = "AI_MODEL_NOT_FOUND"
	ErrCodeAIInferenceFailed         ErrorCode = "AI_INFERENCE_FAILED"
	ErrCodeAIModelLoadingFailed      ErrorCode = "AI_MODEL_LOADING_FAILED"
	ErrCodeAIFeatureExtractionFailed ErrorCode = "AI_FEATURE_EXTRACTION_FAILED"
	ErrCodeAIRecommendationFailed    ErrorCode = "AI_RECOMMENDATION_FAILED"
	ErrCodeAIFraudDetectionFailed    ErrorCode = "AI_FRAUD_DETECTION_FAILED"
	ErrCodeAIContentAnalysisFailed   ErrorCode = "AI_CONTENT_ANALYSIS_FAILED"
	ErrCodeAIModelTrainingFailed     ErrorCode = "AI_MODEL_TRAINING_FAILED"

	// Service-14-Commerce-Service Errors
	ErrCodeProductNotFound           ErrorCode = "PRODUCT_NOT_FOUND"
	ErrCodeOrderCreationFailed       ErrorCode = "ORDER_CREATION_FAILED"
	ErrCodeOrderNotFound             ErrorCode = "ORDER_NOT_FOUND"
	ErrCodePaymentProcessingFailed   ErrorCode = "PAYMENT_PROCESSING_FAILED"
	ErrCodeInventoryInsufficient     ErrorCode = "INVENTORY_INSUFFICIENT"
	ErrCodePricingCalculationFailed  ErrorCode = "PRICING_CALCULATION_FAILED"
	ErrCodeShippingCalculationFailed ErrorCode = "SHIPPING_CALCULATION_FAILED"
	ErrCodeMarketplaceRuleViolation  ErrorCode = "MARKETPLACE_RULE_VIOLATION"

	// Service-15-Notification-Service Errors
	ErrCodeNotificationDeliveryFailed     ErrorCode = "NOTIFICATION_DELIVERY_FAILED"
	ErrCodeNotificationTemplateNotFound   ErrorCode = "NOTIFICATION_TEMPLATE_NOT_FOUND"
	ErrCodeNotificationChannelUnavailable ErrorCode = "NOTIFICATION_CHANNEL_UNAVAILABLE"
	ErrCodeNotificationPreferencesInvalid ErrorCode = "NOTIFICATION_PREFERENCES_INVALID"
	ErrCodeNotificationRateLimitExceeded  ErrorCode = "NOTIFICATION_RATE_LIMIT_EXCEEDED"
	ErrCodePushNotificationFailed         ErrorCode = "PUSH_NOTIFICATION_FAILED"
	ErrCodeEmailNotificationFailed        ErrorCode = "EMAIL_NOTIFICATION_FAILED"
	ErrCodeSMSNotificationFailed          ErrorCode = "SMS_NOTIFICATION_FAILED"

	// Service-16-Search-Service Errors
	ErrCodeSearchIndexingFailed    ErrorCode = "SEARCH_INDEXING_FAILED"
	ErrCodeSearchQueryInvalid      ErrorCode = "SEARCH_QUERY_INVALID"
	ErrCodeSearchTimeout           ErrorCode = "SEARCH_TIMEOUT"
	ErrCodeSearchIndexNotFound     ErrorCode = "SEARCH_INDEX_NOT_FOUND"
	ErrCodeSearchReindexingFailed  ErrorCode = "SEARCH_REINDEXING_FAILED"
	ErrCodeSearchFilterInvalid     ErrorCode = "SEARCH_FILTER_INVALID"
	ErrCodeSearchAggregationFailed ErrorCode = "SEARCH_AGGREGATION_FAILED"

	// Service-17-Analytics-Service Errors
	ErrCodeAnalyticsDataProcessingFailed   ErrorCode = "ANALYTICS_DATA_PROCESSING_FAILED"
	ErrCodeAnalyticsQueryInvalid           ErrorCode = "ANALYTICS_QUERY_INVALID"
	ErrCodeAnalyticsReportGenerationFailed ErrorCode = "ANALYTICS_REPORT_GENERATION_FAILED"
	ErrCodeAnalyticsDataSourceUnavailable  ErrorCode = "ANALYTICS_DATA_SOURCE_UNAVAILABLE"
	ErrCodeAnalyticsMetricNotFound         ErrorCode = "ANALYTICS_METRIC_NOT_FOUND"
	ErrCodeAnalyticsAggregationFailed      ErrorCode = "ANALYTICS_AGGREGATION_FAILED"
	ErrCodeAnalyticsPredictionFailed       ErrorCode = "ANALYTICS_PREDICTION_FAILED"

	// Service-18-Moderation-Service Errors
	ErrCodeModerationRuleNotFound   ErrorCode = "MODERATION_RULE_NOT_FOUND"
	ErrCodeContentFlaggingFailed    ErrorCode = "CONTENT_FLAGGING_FAILED"
	ErrCodeModerationActionFailed   ErrorCode = "MODERATION_ACTION_FAILED"
	ErrCodeAutoModerationFailed     ErrorCode = "AUTO_MODERATION_FAILED"
	ErrCodeModerationReviewRequired ErrorCode = "MODERATION_REVIEW_REQUIRED"
	ErrCodeContentViolationDetected ErrorCode = "CONTENT_VIOLATION_DETECTED"

	// Service-19-Recommendation-Service Errors
	ErrCodeRecommendationGenerationFailed ErrorCode = "RECOMMENDATION_GENERATION_FAILED"
	ErrCodeRecommendationModelNotFound    ErrorCode = "RECOMMENDATION_MODEL_NOT_FOUND"
	ErrCodeUserPreferenceNotFound         ErrorCode = "USER_PREFERENCE_NOT_FOUND"
	ErrCodeRecommendationCacheMiss        ErrorCode = "RECOMMENDATION_CACHE_MISS"
	ErrCodeRecommendationFilteringFailed  ErrorCode = "RECOMMENDATION_FILTERING_FAILED"

	// Service-20-Advertising-Service Errors
	ErrCodeAdCampaignNotFound ErrorCode = "AD_CAMPAIGN_NOT_FOUND"
	ErrCodeAdBudgetExhausted  ErrorCode = "AD_BUDGET_EXHAUSTED"
	ErrCodeAdTargetingInvalid ErrorCode = "AD_TARGETING_INVALID"
	ErrCodeAdCreativeInvalid  ErrorCode = "AD_CREATIVE_INVALID"
	ErrCodeAdServingFailed    ErrorCode = "AD_SERVING_FAILED"
	ErrCodeAdBiddingFailed    ErrorCode = "AD_BIDDING_FAILED"
	ErrCodeAdAnalyticsFailed  ErrorCode = "AD_ANALYTICS_FAILED"

	// Service-21-Admin-Service Errors
	ErrCodeAdminPermissionDenied      ErrorCode = "ADMIN_PERMISSION_DENIED"
	ErrCodeSystemConfigurationInvalid ErrorCode = "SYSTEM_CONFIGURATION_INVALID"
	ErrCodeAdminActionFailed          ErrorCode = "ADMIN_ACTION_FAILED"
	ErrCodeSystemMaintenanceMode      ErrorCode = "SYSTEM_MAINTENANCE_MODE"
	ErrCodeAdminAuditLogFailed        ErrorCode = "ADMIN_AUDIT_LOG_FAILED"

	// Service-22-Kafka-Messaging-Service Errors
	ErrCodeKafkaConnectionFailed    ErrorCode = "KAFKA_CONNECTION_FAILED"
	ErrCodeKafkaTopicNotFound       ErrorCode = "KAFKA_TOPIC_NOT_FOUND"
	ErrCodeKafkaProducerFailed      ErrorCode = "KAFKA_PRODUCER_FAILED"
	ErrCodeKafkaConsumerFailed      ErrorCode = "KAFKA_CONSUMER_FAILED"
	ErrCodeKafkaPartitionError      ErrorCode = "KAFKA_PARTITION_ERROR"
	ErrCodeKafkaOffsetCommitFailed  ErrorCode = "KAFKA_OFFSET_COMMIT_FAILED"
	ErrCodeKafkaRebalanceFailed     ErrorCode = "KAFKA_REBALANCE_FAILED"
	ErrCodeKafkaMessageTooLarge     ErrorCode = "KAFKA_MESSAGE_TOO_LARGE"
	ErrCodeKafkaTopicCreationFailed ErrorCode = "KAFKA_TOPIC_CREATION_FAILED"
)

// GetUSCErrorCodesByService returns error codes grouped by service
func GetUSCErrorCodesByService() map[string][]ErrorCode {
	return map[string][]ErrorCode{
		"service-01-gateway": {
			ErrCodeGatewayTimeout,
			ErrCodeGraphQLParsingFailed,
			ErrCodeGraphQLValidationFailed,
			ErrCodeFederationServiceUnavailable,
			ErrCodeQueryComplexityExceeded,
			ErrCodeSubscriptionLimitExceeded,
		},
		"service-02-auth": {
			ErrCodeMFARequired,
			ErrCodeMFAInvalid,
			ErrCodeMFASetupRequired,
			ErrCodePasswordTooWeak,
			ErrCodePasswordExpired,
			ErrCodeSessionExpired,
			ErrCodeOAuthProviderError,
			ErrCodeEmailVerificationRequired,
		},
		"service-03-user": {
			ErrCodeUserProfileIncomplete,
			ErrCodeUserRelationshipExists,
			ErrCodeUserRelationshipNotFound,
			ErrCodeUserPreferencesInvalid,
			ErrCodeUserSuspended,
			ErrCodeUserDeactivated,
		},
		"service-04-usc-blockchain-core": {
			ErrCodeUSCInsufficientBalance,
			ErrCodeUSCInvalidAmount,
			ErrCodeUSCTransferFailed,
			ErrCodeUSCInvalidAddress,
			ErrCodeBlockchainConnectionFailed,
			ErrCodeBlockchainSyncFailed,
			ErrCodeBlockValidationFailed,
			ErrCodeConsensusFailure,
		},
		"service-05-usc-wallet": {
			ErrCodeWalletNotFound,
			ErrCodeWalletCreationFailed,
			ErrCodeWalletUnlockFailed,
			ErrCodeChatMessageEmpty,
			ErrCodeChatParticipantNotFound,
			ErrCodeChatEncryptionFailed,
			ErrCodeUSCTransactionPending,
			ErrCodeUSCTransactionFailed,
		},
		"service-06-security": {
			ErrCodeSecurityThreatDetected,
			ErrCodeSecurityRuleViolation,
			ErrCodeFraudDetected,
			ErrCodeSuspiciousActivity,
			ErrCodeIPBlocked,
			ErrCodeGeoLocationBlocked,
		},
		"service-07-caching": {
			ErrCodeCacheConnectionFailed,
			ErrCodeCacheKeyNotFound,
			ErrCodeCacheSerializationFailed,
			ErrCodeCacheMemoryExhausted,
		},
		"service-08-monitoring": {
			ErrCodeMetricCollectionFailed,
			ErrCodeHealthCheckFailed,
			ErrCodeMonitoringServiceDown,
			ErrCodeMetricThresholdExceeded,
		},
		"service-09-social": {
			ErrCodePostNotFound,
			ErrCodePostCreationFailed,
			ErrCodeSocialInteractionFailed,
			ErrCodeFeedGenerationFailed,
			ErrCodeSocialLimitExceeded,
		},
		"service-10-usc-bilateral-rewards": {
			ErrCodeRewardCalculationFailed,
			ErrCodeRewardDistributionFailed,
			ErrCodeRewardAlreadyClaimed,
			ErrCodeRewardExpired,
			ErrCodeRewardPoolExhausted,
		},
		"service-11-content-management": {
			ErrCodeFileUploadFailed,
			ErrCodeFileFormatUnsupported,
			ErrCodeFileSizeExceeded,
			ErrCodeContentProcessingFailed,
			ErrCodeCDNDeploymentFailed,
		},
		"service-12-video-service": {
			ErrCodeVideoNotFound,
			ErrCodeVideoUploadFailed,
			ErrCodeVideoStreamingFailed,
			ErrCodeLiveStreamNotFound,
			ErrCodeVideoChatConnectionFailed,
		},
		"service-13-ai-service": {
			ErrCodeAIModelNotFound,
			ErrCodeAIInferenceFailed,
			ErrCodeAIRecommendationFailed,
			ErrCodeAIFraudDetectionFailed,
			ErrCodeAIContentAnalysisFailed,
		},
		"service-14-commerce-service": {
			ErrCodeProductNotFound,
			ErrCodeOrderCreationFailed,
			ErrCodePaymentProcessingFailed,
			ErrCodeInventoryInsufficient,
			ErrCodeMarketplaceRuleViolation,
		},
		"service-15-notification-service": {
			ErrCodeNotificationDeliveryFailed,
			ErrCodeNotificationTemplateNotFound,
			ErrCodePushNotificationFailed,
			ErrCodeEmailNotificationFailed,
			ErrCodeSMSNotificationFailed,
		},
		"service-16-search-service": {
			ErrCodeSearchIndexingFailed,
			ErrCodeSearchQueryInvalid,
			ErrCodeSearchTimeout,
			ErrCodeSearchReindexingFailed,
		},
		"service-17-analytics-service": {
			ErrCodeAnalyticsDataProcessingFailed,
			ErrCodeAnalyticsReportGenerationFailed,
			ErrCodeAnalyticsMetricNotFound,
			ErrCodeAnalyticsPredictionFailed,
		},
		"service-18-moderation-service": {
			ErrCodeModerationRuleNotFound,
			ErrCodeContentFlaggingFailed,
			ErrCodeModerationActionFailed,
			ErrCodeContentViolationDetected,
		},
		"service-19-recommendation-service": {
			ErrCodeRecommendationGenerationFailed,
			ErrCodeRecommendationModelNotFound,
			ErrCodeUserPreferenceNotFound,
			ErrCodeRecommendationFilteringFailed,
		},
		"service-20-advertising-service": {
			ErrCodeAdCampaignNotFound,
			ErrCodeAdBudgetExhausted,
			ErrCodeAdTargetingInvalid,
			ErrCodeAdServingFailed,
		},
		"service-21-admin-service": {
			ErrCodeAdminPermissionDenied,
			ErrCodeSystemConfigurationInvalid,
			ErrCodeAdminActionFailed,
			ErrCodeSystemMaintenanceMode,
		},
		"service-22-kafka-messaging-service": {
			ErrCodeKafkaConnectionFailed,
			ErrCodeKafkaTopicNotFound,
			ErrCodeKafkaProducerFailed,
			ErrCodeKafkaConsumerFailed,
			ErrCodeKafkaMessageTooLarge,
		},
	}
}

// IsUSCSpecificError checks if an error code is USC-specific
func IsUSCSpecificError(code ErrorCode) bool {
	uscErrors := []ErrorCode{
		ErrCodeUSCInsufficientBalance,
		ErrCodeUSCInvalidAmount,
		ErrCodeUSCTransferFailed,
		ErrCodeUSCInvalidAddress,
		ErrCodeBlockchainConnectionFailed,
		ErrCodeWalletNotFound,
		ErrCodeWalletCreationFailed,
		ErrCodeStakingInsufficientAmount,
		ErrCodeNFTNotFound,
		ErrCodeNFTMintingFailed,
		ErrCodeRewardCalculationFailed,
	}

	for _, uscErr := range uscErrors {
		if code == uscErr {
			return true
		}
	}
	return false
}
