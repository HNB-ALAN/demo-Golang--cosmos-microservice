package types

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Store key prefixes
const (
	// TokenKeyPrefix is the prefix for token keys
	TokenKeyPrefix = "token:"

	// BalanceKeyPrefix is the prefix for balance keys
	BalanceKeyPrefix = "balance:"

	// TransferKeyPrefix is the prefix for transfer keys
	TransferKeyPrefix = "transfer:"

	// ParamsKey is the key for module parameters
	ParamsKey = "params"
)

// GetTokenKey returns the key for a token
func GetTokenKey(id string) []byte {
	return append([]byte(TokenKeyPrefix), []byte(id)...)
}

// GetBalanceKey returns the key for a balance
func GetBalanceKey(tokenID, owner string) []byte {
	key := fmt.Sprintf("%s_%s", tokenID, owner)
	return append([]byte(BalanceKeyPrefix), []byte(key)...)
}

// GetTransferKey returns the key for a transfer
func GetTransferKey(id string) []byte {
	return append([]byte(TransferKeyPrefix), []byte(id)...)
}

// GetParamsKey returns the key for module parameters
func GetParamsKey() []byte {
	return []byte(ParamsKey)
}

// ParseTokenKey parses a token key to extract the ID
func ParseTokenKey(key []byte) (string, error) {
	if !bytes.HasPrefix(key, []byte(TokenKeyPrefix)) {
		return "", fmt.Errorf("invalid token key prefix")
	}
	return string(key[len(TokenKeyPrefix):]), nil
}

// ParseBalanceKey parses a balance key to extract the IDs
func ParseBalanceKey(key []byte) (string, string, error) {
	if !bytes.HasPrefix(key, []byte(BalanceKeyPrefix)) {
		return "", "", fmt.Errorf("invalid balance key prefix")
	}

	keyStr := string(key[len(BalanceKeyPrefix):])
	parts := bytes.Split([]byte(keyStr), []byte("_"))
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid balance key format")
	}

	return string(parts[0]), string(parts[1]), nil
}

// ParseTransferKey parses a transfer key to extract the ID
func ParseTransferKey(key []byte) (string, error) {
	if !bytes.HasPrefix(key, []byte(TransferKeyPrefix)) {
		return "", fmt.Errorf("invalid transfer key prefix")
	}
	return string(key[len(TransferKeyPrefix):]), nil
}

// KeyPrefix represents a key prefix
type KeyPrefix []byte

// Key prefixes for different data types
var (
	TokenPrefix    = KeyPrefix(TokenKeyPrefix)
	BalancePrefix  = KeyPrefix(BalanceKeyPrefix)
	TransferPrefix = KeyPrefix(TransferKeyPrefix)
	ParamsPrefix   = KeyPrefix(ParamsKey)
)

// String returns the string representation of the key prefix
func (kp KeyPrefix) String() string {
	return string(kp)
}

// Bytes returns the byte representation of the key prefix
func (kp KeyPrefix) Bytes() []byte {
	return []byte(kp)
}

// IsValidAddress checks if an address is valid
func IsValidAddress(address string) bool {
	_, err := sdk.AccAddressFromBech32(address)
	return err == nil
}

// ValidateAmount validates an amount
func ValidateAmount(amount string) error {
	if amount == "" {
		return fmt.Errorf("amount cannot be empty")
	}
	if len(amount) > 128 {
		return fmt.Errorf("amount too long")
	}
	return nil
}
