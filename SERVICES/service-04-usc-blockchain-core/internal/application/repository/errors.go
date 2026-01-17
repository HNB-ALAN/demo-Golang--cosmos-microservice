package repository

import (
	"database/sql"
	"fmt"

	"service-04/internal/infrastructure/database"

	"github.com/usc-platform/shared/errors"
)

// Repository error codes (2000-2999 range for repository layer)
const (
	// Block operations errors (2000-2099)
	ErrBlockNotFound = 2000 + iota
	ErrBlockInvalid
	ErrBlockValidationFailed
	ErrBlockProductionFailed

	// Transaction operations errors (2100-2199)
	ErrTransactionNotFound
	ErrTransactionInvalid
	ErrTransactionSubmissionFailed
	ErrTransactionValidationFailed

	// USC Coin operations errors (2200-2299)
	ErrUSCBalanceNotFound
	ErrUSCTransferFailed
	ErrUSCInsufficientBalance
	ErrUSCInvalidAmount

	// Smart Contract operations errors (2300-2399)
	ErrContractNotFound
	ErrContractDeploymentFailed
	ErrContractExecutionFailed
	ErrContractQueryFailed

	// NFT Token operations errors (2400-2499)
	ErrNFTNotFound
	ErrNFTMintFailed
	ErrNFTTransferFailed
	ErrNFTBurnFailed

	// Custom Token operations errors (2500-2599)
	ErrTokenNotFound
	ErrTokenCreationFailed
	ErrTokenMintFailed
	ErrTokenBurnFailed

	// Product Certificate operations errors (2600-2699)
	ErrCertificateNotFound
	ErrCertificateCreationFailed
	ErrCertificateVerificationFailed

	// Validator operations errors (2700-2799)
	ErrValidatorNotFound
	ErrValidatorRegistrationFailed
	ErrStakingFailed
	ErrUnstakingFailed

	// Network operations errors (2800-2899)
	ErrNetworkInfoNotFound
	ErrNetworkSyncFailed

	// Streaming operations errors (2900-2999)
	ErrStreamNotFound
	ErrStreamCreationFailed

	// Store Bridge operations errors (3000-3099)
	ErrBridgeNotFound
	ErrBridgeDeploymentFailed
	ErrBridgeValidationFailed

	// Store Network operations errors (3100-3199)
	ErrStoreNetworkNotFound
	ErrStoreNetworkSyncFailed

	// Common repository errors (3200-3299)
	ErrDatabaseUnavailable
	ErrBlockchainUnavailable
	ErrInvalidRequest
	ErrContextCancelled
)

// Error messages map
var errorMessages = map[int]string{
	// Block operations
	ErrBlockNotFound:         "block not found",
	ErrBlockInvalid:          "invalid block",
	ErrBlockValidationFailed: "block validation failed",
	ErrBlockProductionFailed: "block production failed",

	// Transaction operations
	ErrTransactionNotFound:         "transaction not found",
	ErrTransactionInvalid:          "invalid transaction",
	ErrTransactionSubmissionFailed: "transaction submission failed",
	ErrTransactionValidationFailed: "transaction validation failed",

	// USC Coin operations
	ErrUSCBalanceNotFound:     "USC balance not found",
	ErrUSCTransferFailed:      "USC transfer failed",
	ErrUSCInsufficientBalance: "insufficient USC balance",
	ErrUSCInvalidAmount:       "invalid USC amount",

	// Smart Contract operations
	ErrContractNotFound:         "contract not found",
	ErrContractDeploymentFailed: "contract deployment failed",
	ErrContractExecutionFailed:  "contract execution failed",
	ErrContractQueryFailed:      "contract query failed",

	// NFT Token operations
	ErrNFTNotFound:       "NFT not found",
	ErrNFTMintFailed:     "NFT mint failed",
	ErrNFTTransferFailed: "NFT transfer failed",
	ErrNFTBurnFailed:     "NFT burn failed",

	// Custom Token operations
	ErrTokenNotFound:       "token not found",
	ErrTokenCreationFailed: "token creation failed",
	ErrTokenMintFailed:     "token mint failed",
	ErrTokenBurnFailed:     "token burn failed",

	// Product Certificate operations
	ErrCertificateNotFound:           "certificate not found",
	ErrCertificateCreationFailed:     "certificate creation failed",
	ErrCertificateVerificationFailed: "certificate verification failed",

	// Validator operations
	ErrValidatorNotFound:           "validator not found",
	ErrValidatorRegistrationFailed: "validator registration failed",
	ErrStakingFailed:               "staking failed",
	ErrUnstakingFailed:             "unstaking failed",

	// Network operations
	ErrNetworkInfoNotFound: "network info not found",
	ErrNetworkSyncFailed:   "network sync failed",

	// Streaming operations
	ErrStreamNotFound:       "stream not found",
	ErrStreamCreationFailed: "stream creation failed",

	// Store Bridge operations
	ErrBridgeNotFound:         "bridge not found",
	ErrBridgeDeploymentFailed: "bridge deployment failed",
	ErrBridgeValidationFailed: "bridge validation failed",

	// Store Network operations
	ErrStoreNetworkNotFound:   "store network not found",
	ErrStoreNetworkSyncFailed: "store network sync failed",

	// Common repository errors
	ErrDatabaseUnavailable:   "database not available",
	ErrBlockchainUnavailable: "blockchain not available",
	ErrInvalidRequest:        "invalid request",
	ErrContextCancelled:      "context cancelled",
}

// NewRepositoryError creates a new repository error
func NewRepositoryError(code int, details ...interface{}) error {
	message, exists := errorMessages[code]
	if !exists {
		message = "unknown repository error"
	}

	if len(details) > 0 {
		message = fmt.Sprintf("%s: %v", message, details[0])
	}

	return &errors.DomainError{
		Code:      errors.ErrorCode(fmt.Sprintf("REPO_%d", code)),
		Message:   message,
		Details:   fmt.Sprintf("%v", details),
		Severity:  errors.SeverityMedium,
		Category:  errors.CategorySystem,
		Retryable: isRetryableError(code),
	}
}

// WrapRepositoryError wraps an existing error with repository error context
func WrapRepositoryError(code int, err error, details ...interface{}) error {
	if err == nil {
		return nil
	}

	message, exists := errorMessages[code]
	if !exists {
		message = "unknown repository error"
	}

	if len(details) > 0 {
		message = fmt.Sprintf("%s: %v", message, details[0])
	}

	return &errors.DomainError{
		Code:      errors.ErrorCode(fmt.Sprintf("REPO_%d", code)),
		Message:   fmt.Sprintf("%s: %s", message, err.Error()),
		Details:   fmt.Sprintf("%v", details),
		Severity:  errors.SeverityMedium,
		Category:  errors.CategorySystem,
		Retryable: isRetryableError(code),
	}
}

// isRetryableError determines if an error is retryable
func isRetryableError(code int) bool {
	retryableCodes := map[int]bool{
		ErrDatabaseUnavailable:   true,
		ErrBlockchainUnavailable: true,
		ErrContextCancelled:      false,
		ErrInvalidRequest:        false,
		ErrBlockNotFound:         false,
		ErrTransactionNotFound:   false,
		ErrContractNotFound:      false,
		ErrNFTNotFound:           false,
		ErrTokenNotFound:         false,
		ErrCertificateNotFound:   false,
		ErrValidatorNotFound:     false,
		ErrStreamNotFound:        false,
		ErrBridgeNotFound:        false,
		ErrStoreNetworkNotFound:  false,
	}

	retryable, exists := retryableCodes[code]
	if !exists {
		// Default: operation errors are retryable, validation errors are not
		return code%100 < 50 // Operations (0-49) are retryable, validations (50-99) are not
	}

	return retryable
}

// Helper functions for common error patterns

// NewNotFoundError creates a not found error
func NewNotFoundError(resourceType string, identifier string) error {
	return NewRepositoryError(ErrBlockNotFound, fmt.Sprintf("%s not found: %s", resourceType, identifier))
}

// NewValidationError creates a validation error
func NewValidationError(field string, reason string) error {
	return NewRepositoryError(ErrInvalidRequest, fmt.Sprintf("validation failed for field '%s': %s", field, reason))
}

// NewDatabaseError creates a database error
func NewDatabaseError(operation string, err error) error {
	return WrapRepositoryError(ErrDatabaseUnavailable, err, fmt.Sprintf("database operation failed: %s", operation))
}

// NewBlockchainError creates a blockchain error
func NewBlockchainError(operation string, err error) error {
	return WrapRepositoryError(ErrBlockchainUnavailable, err, fmt.Sprintf("blockchain operation failed: %s", operation))
}

// Database helper functions

// GetPostgresConnection safely retrieves PostgreSQL connection from database manager
// Returns (connection, nil) if available, (nil, nil) if database manager is nil or connection unavailable
// This helper reduces duplicate code in database fallback methods
func GetPostgresConnection(db *database.PostgreSQLManager) *sql.DB {
	if db == nil {
		return nil
	}
	return db.GetPostgres()
}

// IsPostgresAvailable checks if PostgreSQL connection is available
// Returns true if database manager and connection are both available
func IsPostgresAvailable(db *database.PostgreSQLManager) bool {
	return GetPostgresConnection(db) != nil
}
