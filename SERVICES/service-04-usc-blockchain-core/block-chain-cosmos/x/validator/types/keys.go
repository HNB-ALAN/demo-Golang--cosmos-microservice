package types

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Store key prefixes
const (
	// ValidatorKeyPrefix is the prefix for validator keys
	ValidatorKeyPrefix = "validator:"

	// DelegationKeyPrefix is the prefix for delegation keys
	DelegationKeyPrefix = "delegation:"

	// ParamsKey is the key for module parameters
	ParamsKey = "params"
)

// GetValidatorKey returns the key for a validator
func GetValidatorKey(address string) []byte {
	return append([]byte(ValidatorKeyPrefix), []byte(address)...)
}

// GetDelegationKey returns the key for a delegation
func GetDelegationKey(delegatorAddress, validatorAddress string) []byte {
	key := fmt.Sprintf("%s_%s", delegatorAddress, validatorAddress)
	return append([]byte(DelegationKeyPrefix), []byte(key)...)
}

// GetParamsKey returns the key for module parameters
func GetParamsKey() []byte {
	return []byte(ParamsKey)
}

// ParseValidatorKey parses a validator key to extract the address
func ParseValidatorKey(key []byte) (string, error) {
	if !bytes.HasPrefix(key, []byte(ValidatorKeyPrefix)) {
		return "", fmt.Errorf("invalid validator key prefix")
	}
	return string(key[len(ValidatorKeyPrefix):]), nil
}

// ParseDelegationKey parses a delegation key to extract the addresses
func ParseDelegationKey(key []byte) (string, string, error) {
	if !bytes.HasPrefix(key, []byte(DelegationKeyPrefix)) {
		return "", "", fmt.Errorf("invalid delegation key prefix")
	}

	keyStr := string(key[len(DelegationKeyPrefix):])
	parts := bytes.Split([]byte(keyStr), []byte("_"))
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid delegation key format")
	}

	return string(parts[0]), string(parts[1]), nil
}

// KeyPrefix represents a key prefix
type KeyPrefix []byte

// Key prefixes for different data types
var (
	ValidatorPrefix  = KeyPrefix(ValidatorKeyPrefix)
	DelegationPrefix = KeyPrefix(DelegationKeyPrefix)
	ParamsPrefix     = KeyPrefix(ParamsKey)
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
