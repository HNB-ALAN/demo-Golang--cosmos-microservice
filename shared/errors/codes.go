// Package errors provides error handling utilities for USC platform services.
package errors

// ErrorCodeRegistry maintains a registry of all error codes
type ErrorCodeRegistry struct {
	codes map[ErrorCode]ErrorCodeInfo
}

// ErrorCodeInfo contains information about an error code
type ErrorCodeInfo struct {
	Code        ErrorCode     `json:"code"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Severity    ErrorSeverity `json:"severity"`
	Category    ErrorCategory `json:"category"`
	Retryable   bool          `json:"retryable"`
	HTTPCode    int           `json:"http_code"`
	GRPCCode    string        `json:"grpc_code"`
}

// NewErrorCodeRegistry creates a new error code registry
func NewErrorCodeRegistry() *ErrorCodeRegistry {
	registry := &ErrorCodeRegistry{
		codes: make(map[ErrorCode]ErrorCodeInfo),
	}

	// Register all error codes
	registry.registerDefaultCodes()

	return registry
}

// RegisterCode registers a new error code
func (r *ErrorCodeRegistry) RegisterCode(info ErrorCodeInfo) {
	r.codes[info.Code] = info
}

// GetCodeInfo retrieves information about an error code
func (r *ErrorCodeRegistry) GetCodeInfo(code ErrorCode) (ErrorCodeInfo, bool) {
	info, exists := r.codes[code]
	return info, exists
}

// GetAllCodes returns all registered error codes
func (r *ErrorCodeRegistry) GetCodes() map[ErrorCode]ErrorCodeInfo {
	return r.codes
}

// GetCodesByCategory returns error codes filtered by category
func (r *ErrorCodeRegistry) GetCodesByCategory(category ErrorCategory) map[ErrorCode]ErrorCodeInfo {
	result := make(map[ErrorCode]ErrorCodeInfo)
	for code, info := range r.codes {
		if info.Category == category {
			result[code] = info
		}
	}
	return result
}

// GetCodesBySeverity returns error codes filtered by severity
func (r *ErrorCodeRegistry) GetCodesBySeverity(severity ErrorSeverity) map[ErrorCode]ErrorCodeInfo {
	result := make(map[ErrorCode]ErrorCodeInfo)
	for code, info := range r.codes {
		if info.Severity == severity {
			result[code] = info
		}
	}
	return result
}

// GetRetryableCodes returns all retryable error codes
func (r *ErrorCodeRegistry) GetRetryableCodes() map[ErrorCode]ErrorCodeInfo {
	result := make(map[ErrorCode]ErrorCodeInfo)
	for code, info := range r.codes {
		if info.Retryable {
			result[code] = info
		}
	}
	return result
}

// registerDefaultCodes registers all default error codes
func (r *ErrorCodeRegistry) registerDefaultCodes() {
	// General errors
	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeInternal,
		Name:        "Internal Error",
		Description: "An internal server error occurred",
		Severity:    SeverityHigh,
		Category:    CategorySystem,
		Retryable:   false,
		HTTPCode:    500,
		GRPCCode:    "INTERNAL",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeInvalidInput,
		Name:        "Invalid Input",
		Description: "The provided input is invalid",
		Severity:    SeverityMedium,
		Category:    CategoryValidation,
		Retryable:   false,
		HTTPCode:    400,
		GRPCCode:    "INVALID_ARGUMENT",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeNotFound,
		Name:        "Not Found",
		Description: "The requested resource was not found",
		Severity:    SeverityMedium,
		Category:    CategoryBusiness,
		Retryable:   false,
		HTTPCode:    404,
		GRPCCode:    "NOT_FOUND",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeUnauthorized,
		Name:        "Unauthorized",
		Description: "Authentication is required",
		Severity:    SeverityMedium,
		Category:    CategoryAuth,
		Retryable:   false,
		HTTPCode:    401,
		GRPCCode:    "UNAUTHENTICATED",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeForbidden,
		Name:        "Forbidden",
		Description: "Access to the resource is forbidden",
		Severity:    SeverityMedium,
		Category:    CategoryAuth,
		Retryable:   false,
		HTTPCode:    403,
		GRPCCode:    "PERMISSION_DENIED",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeConflict,
		Name:        "Conflict",
		Description: "The request conflicts with the current state",
		Severity:    SeverityMedium,
		Category:    CategoryBusiness,
		Retryable:   false,
		HTTPCode:    409,
		GRPCCode:    "ALREADY_EXISTS",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeTimeout,
		Name:        "Timeout",
		Description: "The request timed out",
		Severity:    SeverityMedium,
		Category:    CategorySystem,
		Retryable:   true,
		HTTPCode:    408,
		GRPCCode:    "DEADLINE_EXCEEDED",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeRateLimited,
		Name:        "Rate Limited",
		Description: "Too many requests",
		Severity:    SeverityMedium,
		Category:    CategorySystem,
		Retryable:   true,
		HTTPCode:    429,
		GRPCCode:    "RESOURCE_EXHAUSTED",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeServiceUnavailable,
		Name:        "Service Unavailable",
		Description: "The service is temporarily unavailable",
		Severity:    SeverityHigh,
		Category:    CategorySystem,
		Retryable:   true,
		HTTPCode:    503,
		GRPCCode:    "UNAVAILABLE",
	})

	// Database errors
	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeDatabaseConnection,
		Name:        "Database Connection Error",
		Description: "Failed to connect to the database",
		Severity:    SeverityHigh,
		Category:    CategoryDatabase,
		Retryable:   true,
		HTTPCode:    503,
		GRPCCode:    "UNAVAILABLE",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeDatabaseQuery,
		Name:        "Database Query Error",
		Description: "Failed to execute database query",
		Severity:    SeverityHigh,
		Category:    CategoryDatabase,
		Retryable:   true,
		HTTPCode:    500,
		GRPCCode:    "INTERNAL",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeDatabaseTransaction,
		Name:        "Database Transaction Error",
		Description: "Failed to execute database transaction",
		Severity:    SeverityHigh,
		Category:    CategoryDatabase,
		Retryable:   true,
		HTTPCode:    500,
		GRPCCode:    "INTERNAL",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeDatabaseConstraint,
		Name:        "Database Constraint Error",
		Description: "Database constraint violation",
		Severity:    SeverityMedium,
		Category:    CategoryDatabase,
		Retryable:   false,
		HTTPCode:    400,
		GRPCCode:    "INVALID_ARGUMENT",
	})

	// Authentication errors
	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeInvalidCredentials,
		Name:        "Invalid Credentials",
		Description: "The provided credentials are invalid",
		Severity:    SeverityMedium,
		Category:    CategoryAuth,
		Retryable:   false,
		HTTPCode:    401,
		GRPCCode:    "UNAUTHENTICATED",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeTokenExpired,
		Name:        "Token Expired",
		Description: "The authentication token has expired",
		Severity:    SeverityMedium,
		Category:    CategoryAuth,
		Retryable:   false,
		HTTPCode:    401,
		GRPCCode:    "UNAUTHENTICATED",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeTokenInvalid,
		Name:        "Token Invalid",
		Description: "The authentication token is invalid",
		Severity:    SeverityMedium,
		Category:    CategoryAuth,
		Retryable:   false,
		HTTPCode:    401,
		GRPCCode:    "UNAUTHENTICATED",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeAccountLocked,
		Name:        "Account Locked",
		Description: "The account is locked",
		Severity:    SeverityMedium,
		Category:    CategoryAuth,
		Retryable:   false,
		HTTPCode:    403,
		GRPCCode:    "PERMISSION_DENIED",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeAccountDisabled,
		Name:        "Account Disabled",
		Description: "The account is disabled",
		Severity:    SeverityMedium,
		Category:    CategoryAuth,
		Retryable:   false,
		HTTPCode:    403,
		GRPCCode:    "PERMISSION_DENIED",
	})

	// Validation errors
	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeValidationFailed,
		Name:        "Validation Failed",
		Description: "Input validation failed",
		Severity:    SeverityLow,
		Category:    CategoryValidation,
		Retryable:   false,
		HTTPCode:    400,
		GRPCCode:    "INVALID_ARGUMENT",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeRequiredField,
		Name:        "Required Field",
		Description: "A required field is missing",
		Severity:    SeverityLow,
		Category:    CategoryValidation,
		Retryable:   false,
		HTTPCode:    400,
		GRPCCode:    "INVALID_ARGUMENT",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeInvalidFormat,
		Name:        "Invalid Format",
		Description: "The provided format is invalid",
		Severity:    SeverityLow,
		Category:    CategoryValidation,
		Retryable:   false,
		HTTPCode:    400,
		GRPCCode:    "INVALID_ARGUMENT",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeOutOfRange,
		Name:        "Out of Range",
		Description: "The value is out of the allowed range",
		Severity:    SeverityLow,
		Category:    CategoryValidation,
		Retryable:   false,
		HTTPCode:    400,
		GRPCCode:    "OUT_OF_RANGE",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeDuplicateValue,
		Name:        "Duplicate Value",
		Description: "The value already exists",
		Severity:    SeverityMedium,
		Category:    CategoryValidation,
		Retryable:   false,
		HTTPCode:    409,
		GRPCCode:    "ALREADY_EXISTS",
	})

	// Business logic errors
	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeInsufficientFunds,
		Name:        "Insufficient Funds",
		Description: "Insufficient funds for the operation",
		Severity:    SeverityMedium,
		Category:    CategoryBusiness,
		Retryable:   false,
		HTTPCode:    400,
		GRPCCode:    "FAILED_PRECONDITION",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeResourceExhausted,
		Name:        "Resource Exhausted",
		Description: "System resources are exhausted",
		Severity:    SeverityHigh,
		Category:    CategorySystem,
		Retryable:   true,
		HTTPCode:    429,
		GRPCCode:    "RESOURCE_EXHAUSTED",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeOperationNotAllowed,
		Name:        "Operation Not Allowed",
		Description: "The operation is not allowed",
		Severity:    SeverityMedium,
		Category:    CategoryBusiness,
		Retryable:   false,
		HTTPCode:    403,
		GRPCCode:    "PERMISSION_DENIED",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeQuotaExceeded,
		Name:        "Quota Exceeded",
		Description: "The quota has been exceeded",
		Severity:    SeverityMedium,
		Category:    CategoryBusiness,
		Retryable:   true,
		HTTPCode:    429,
		GRPCCode:    "RESOURCE_EXHAUSTED",
	})

	// Register USC-specific error codes
	r.registerUSCErrorCodes()
}

// registerUSCErrorCodes registers all USC-specific error codes
func (r *ErrorCodeRegistry) registerUSCErrorCodes() {
	// USC Blockchain Core Errors
	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeUSCInsufficientBalance,
		Name:        "USC Insufficient Balance",
		Description: "Insufficient USC token balance for transaction",
		Severity:    SeverityMedium,
		Category:    CategoryUSCBlockchain,
		Retryable:   false,
		HTTPCode:    400,
		GRPCCode:    "FAILED_PRECONDITION",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeUSCInvalidAmount,
		Name:        "USC Invalid Amount",
		Description: "Invalid USC token amount specified",
		Severity:    SeverityMedium,
		Category:    CategoryUSCBlockchain,
		Retryable:   false,
		HTTPCode:    400,
		GRPCCode:    "INVALID_ARGUMENT",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeUSCTransferFailed,
		Name:        "USC Transfer Failed",
		Description: "USC token transfer operation failed",
		Severity:    SeverityHigh,
		Category:    CategoryUSCBlockchain,
		Retryable:   true,
		HTTPCode:    500,
		GRPCCode:    "INTERNAL",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeBlockchainConnectionFailed,
		Name:        "Blockchain Connection Failed",
		Description: "Failed to connect to USC blockchain network",
		Severity:    SeverityHigh,
		Category:    CategoryUSCBlockchain,
		Retryable:   true,
		HTTPCode:    503,
		GRPCCode:    "UNAVAILABLE",
	})

	// USC Wallet Errors
	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeWalletNotFound,
		Name:        "Wallet Not Found",
		Description: "USC wallet not found for user",
		Severity:    SeverityMedium,
		Category:    CategoryUSCWallet,
		Retryable:   false,
		HTTPCode:    404,
		GRPCCode:    "NOT_FOUND",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeWalletCreationFailed,
		Name:        "Wallet Creation Failed",
		Description: "Failed to create USC wallet",
		Severity:    SeverityHigh,
		Category:    CategoryUSCWallet,
		Retryable:   true,
		HTTPCode:    500,
		GRPCCode:    "INTERNAL",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeChatMessageEmpty,
		Name:        "Chat Message Empty",
		Description: "Wallet chat message cannot be empty",
		Severity:    SeverityLow,
		Category:    CategoryUSCWallet,
		Retryable:   false,
		HTTPCode:    400,
		GRPCCode:    "INVALID_ARGUMENT",
	})

	// USC NFT Errors
	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeNFTNotFound,
		Name:        "NFT Not Found",
		Description: "USC NFT not found",
		Severity:    SeverityMedium,
		Category:    CategoryUSCNFT,
		Retryable:   false,
		HTTPCode:    404,
		GRPCCode:    "NOT_FOUND",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeNFTMintingFailed,
		Name:        "NFT Minting Failed",
		Description: "Failed to mint USC NFT",
		Severity:    SeverityHigh,
		Category:    CategoryUSCNFT,
		Retryable:   true,
		HTTPCode:    500,
		GRPCCode:    "INTERNAL",
	})

	// USC Staking Errors
	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeStakingInsufficientAmount,
		Name:        "Staking Insufficient Amount",
		Description: "Insufficient amount for USC staking",
		Severity:    SeverityMedium,
		Category:    CategoryUSCStaking,
		Retryable:   false,
		HTTPCode:    400,
		GRPCCode:    "FAILED_PRECONDITION",
	})

	// USC Rewards Errors
	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeRewardCalculationFailed,
		Name:        "Reward Calculation Failed",
		Description: "Failed to calculate USC bilateral rewards",
		Severity:    SeverityHigh,
		Category:    CategoryUSCRewards,
		Retryable:   true,
		HTTPCode:    500,
		GRPCCode:    "INTERNAL",
	})

	// Service-Specific Errors
	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeGatewayTimeout,
		Name:        "Gateway Timeout",
		Description: "Gateway request timeout",
		Severity:    SeverityMedium,
		Category:    CategoryGateway,
		Retryable:   true,
		HTTPCode:    504,
		GRPCCode:    "DEADLINE_EXCEEDED",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeMFARequired,
		Name:        "MFA Required",
		Description: "Multi-factor authentication is required",
		Severity:    SeverityMedium,
		Category:    CategoryAuth,
		Retryable:   false,
		HTTPCode:    401,
		GRPCCode:    "UNAUTHENTICATED",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodePostNotFound,
		Name:        "Post Not Found",
		Description: "Social media post not found",
		Severity:    SeverityMedium,
		Category:    CategorySocial,
		Retryable:   false,
		HTTPCode:    404,
		GRPCCode:    "NOT_FOUND",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeVideoNotFound,
		Name:        "Video Not Found",
		Description: "Video content not found",
		Severity:    SeverityMedium,
		Category:    CategoryVideo,
		Retryable:   false,
		HTTPCode:    404,
		GRPCCode:    "NOT_FOUND",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeProductNotFound,
		Name:        "Product Not Found",
		Description: "Commerce product not found",
		Severity:    SeverityMedium,
		Category:    CategoryCommerce,
		Retryable:   false,
		HTTPCode:    404,
		GRPCCode:    "NOT_FOUND",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeAIModelNotFound,
		Name:        "AI Model Not Found",
		Description: "AI model not found",
		Severity:    SeverityMedium,
		Category:    CategoryAI,
		Retryable:   false,
		HTTPCode:    404,
		GRPCCode:    "NOT_FOUND",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeSearchIndexingFailed,
		Name:        "Search Indexing Failed",
		Description: "Failed to index content for search",
		Severity:    SeverityMedium,
		Category:    CategorySearch,
		Retryable:   true,
		HTTPCode:    500,
		GRPCCode:    "INTERNAL",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeAnalyticsDataProcessingFailed,
		Name:        "Analytics Data Processing Failed",
		Description: "Failed to process analytics data",
		Severity:    SeverityMedium,
		Category:    CategoryAnalytics,
		Retryable:   true,
		HTTPCode:    500,
		GRPCCode:    "INTERNAL",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeContentFlaggingFailed,
		Name:        "Content Flagging Failed",
		Description: "Failed to flag content for moderation",
		Severity:    SeverityMedium,
		Category:    CategoryModeration,
		Retryable:   true,
		HTTPCode:    500,
		GRPCCode:    "INTERNAL",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeAdCampaignNotFound,
		Name:        "Ad Campaign Not Found",
		Description: "Advertising campaign not found",
		Severity:    SeverityMedium,
		Category:    CategoryAdvertising,
		Retryable:   false,
		HTTPCode:    404,
		GRPCCode:    "NOT_FOUND",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeKafkaConnectionFailed,
		Name:        "Kafka Connection Failed",
		Description: "Failed to connect to Kafka messaging service",
		Severity:    SeverityHigh,
		Category:    CategoryMessaging,
		Retryable:   true,
		HTTPCode:    503,
		GRPCCode:    "UNAVAILABLE",
	})

	r.RegisterCode(ErrorCodeInfo{
		Code:        ErrCodeAdminPermissionDenied,
		Name:        "Admin Permission Denied",
		Description: "Admin permission required for this operation",
		Severity:    SeverityMedium,
		Category:    CategoryAdmin,
		Retryable:   false,
		HTTPCode:    403,
		GRPCCode:    "PERMISSION_DENIED",
	})
}

// Global error code registry instance
var globalErrorCodeRegistry = NewErrorCodeRegistry()

// GetGlobalErrorCodeRegistry returns the global error code registry
func GetGlobalErrorCodeRegistry() *ErrorCodeRegistry {
	return globalErrorCodeRegistry
}

// GetErrorCodeInfo retrieves information about an error code from the global registry
func GetErrorCodeInfo(code ErrorCode) (ErrorCodeInfo, bool) {
	return globalErrorCodeRegistry.GetCodeInfo(code)
}

// GetErrorCodesByCategory returns error codes filtered by category from the global registry
func GetErrorCodesByCategory(category ErrorCategory) map[ErrorCode]ErrorCodeInfo {
	return globalErrorCodeRegistry.GetCodesByCategory(category)
}

// GetErrorCodesBySeverity returns error codes filtered by severity from the global registry
func GetErrorCodesBySeverity(severity ErrorSeverity) map[ErrorCode]ErrorCodeInfo {
	return globalErrorCodeRegistry.GetCodesBySeverity(severity)
}

// GetRetryableErrorCodes returns all retryable error codes from the global registry
func GetRetryableErrorCodes() map[ErrorCode]ErrorCodeInfo {
	return globalErrorCodeRegistry.GetRetryableCodes()
}

// ErrorCodeValidator validates error codes
type ErrorCodeValidator struct {
	registry *ErrorCodeRegistry
}

// NewErrorCodeValidator creates a new error code validator
func NewErrorCodeValidator(registry *ErrorCodeRegistry) *ErrorCodeValidator {
	return &ErrorCodeValidator{
		registry: registry,
	}
}

// ValidateErrorCode validates if an error code is registered
func (v *ErrorCodeValidator) ValidateErrorCode(code ErrorCode) bool {
	_, exists := v.registry.GetCodeInfo(code)
	return exists
}

// ValidateDomainError validates if a domain error has a valid error code
func (v *ErrorCodeValidator) ValidateDomainError(err *DomainError) bool {
	return v.ValidateErrorCode(err.Code)
}

// GetErrorCodeStatistics returns statistics about error codes
func (v *ErrorCodeValidator) GetErrorCodeStatistics() map[string]interface{} {
	stats := make(map[string]interface{})

	allCodes := v.registry.GetCodes()
	stats["total_codes"] = len(allCodes)

	// Count by category
	categoryCount := make(map[ErrorCategory]int)
	for _, info := range allCodes {
		categoryCount[info.Category]++
	}
	stats["by_category"] = categoryCount

	// Count by severity
	severityCount := make(map[ErrorSeverity]int)
	for _, info := range allCodes {
		severityCount[info.Severity]++
	}
	stats["by_severity"] = severityCount

	// Count retryable codes
	retryableCount := 0
	for _, info := range allCodes {
		if info.Retryable {
			retryableCount++
		}
	}
	stats["retryable_codes"] = retryableCount

	return stats
}
