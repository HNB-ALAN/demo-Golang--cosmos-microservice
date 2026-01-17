package types

import (
	"bytes"
	"fmt"
	"strconv"
)

// GetTransactionKey returns the store key for a transaction
func GetTransactionKey(hash string) []byte {
	return []byte(TransactionKeyPrefix + hash)
}

// GetTransactionIDKey returns the store key for a transaction by ID
func GetTransactionIDKey(id string) []byte {
	return []byte(TransactionKeyPrefix + "id:" + id)
}

// GetStatsKey returns the store key for transaction statistics
func GetStatsKey() []byte {
	return []byte(StatsKeyPrefix + "main")
}

// GetTransactionByAddressKey returns the store key for transactions by address
func GetTransactionByAddressKey(address string) []byte {
	return []byte(TransactionKeyPrefix + "address:" + address)
}

// GetTransactionByTypeKey returns the store key for transactions by type
func GetTransactionByTypeKey(txType string) []byte {
	return []byte(TransactionKeyPrefix + "type:" + txType)
}

// GetTransactionByStatusKey returns the store key for transactions by status
func GetTransactionByStatusKey(status string) []byte {
	return []byte(TransactionKeyPrefix + "status:" + status)
}

// GetTransactionByTimeKey returns the store key for transactions by time range
func GetTransactionByTimeKey(timestamp int64) []byte {
	return []byte(TransactionKeyPrefix + "time:" + strconv.FormatInt(timestamp, 10))
}

// ParseTransactionKey parses a transaction key and returns the hash
func ParseTransactionKey(key []byte) string {
	if !bytes.HasPrefix(key, []byte(TransactionKeyPrefix)) {
		return ""
	}
	return string(key[len(TransactionKeyPrefix):])
}

// ParseTransactionIDKey parses a transaction ID key and returns the ID
func ParseTransactionIDKey(key []byte) string {
	if !bytes.HasPrefix(key, []byte(TransactionKeyPrefix+"id:")) {
		return ""
	}
	return string(key[len(TransactionKeyPrefix+"id:"):])
}

// IsValidAddress validates an address format
func IsValidAddress(address string) bool {
	if address == "" {
		return false
	}
	// Basic validation - should be at least 20 characters
	return len(address) >= 20
}

// ValidateAmount validates an amount string
func ValidateAmount(amount string) error {
	if amount == "" {
		return fmt.Errorf("amount cannot be empty")
	}
	// Basic validation - should be numeric
	if _, err := strconv.ParseFloat(amount, 64); err != nil {
		return fmt.Errorf("invalid amount format: %s", amount)
	}
	return nil
}

// ValidateTransactionType validates a transaction type
func ValidateTransactionType(txType string) error {
	validTypes := []string{"transfer", "mint", "burn", "contract_call", "delegate", "undelegate"}
	for _, validType := range validTypes {
		if txType == validType {
			return nil
		}
	}
	return fmt.Errorf("invalid transaction type: %s", txType)
}

// ValidateTransactionStatus validates a transaction status
func ValidateTransactionStatus(status string) error {
	validStatuses := []string{"pending", "validated", "executed", "failed", "cancelled"}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return nil
		}
	}
	return fmt.Errorf("invalid transaction status: %s", status)
}
