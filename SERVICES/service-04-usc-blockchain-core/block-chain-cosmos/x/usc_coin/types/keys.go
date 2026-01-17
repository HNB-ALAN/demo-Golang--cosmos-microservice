package types

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Store key prefixes
const (
	// BalanceKeyPrefix is the prefix for balance keys
	BalanceKeyPrefix = "balance:"

	// TransferKeyPrefix is the prefix for transfer keys
	TransferKeyPrefix = "transfer:"

	// SupplyKey is the key for total supply
	SupplyKey = "supply"

	// ParamsKey is the key for module parameters
	ParamsKey = "params"
)

// GetBalanceKey returns the key for a balance
func GetBalanceKey(address string) []byte {
	return append([]byte(BalanceKeyPrefix), []byte(address)...)
}

// GetTransferKey returns the key for a transfer
func GetTransferKey(txHash string) []byte {
	return append([]byte(TransferKeyPrefix), []byte(txHash)...)
}

// GetSupplyKey returns the key for total supply
func GetSupplyKey() []byte {
	return []byte(SupplyKey)
}

// GetParamsKey returns the key for module parameters
func GetParamsKey() []byte {
	return []byte(ParamsKey)
}

// ParseBalanceKey parses a balance key to extract the address
func ParseBalanceKey(key []byte) (string, error) {
	if !bytes.HasPrefix(key, []byte(BalanceKeyPrefix)) {
		return "", fmt.Errorf("invalid balance key prefix")
	}
	return string(key[len(BalanceKeyPrefix):]), nil
}

// ParseTransferKey parses a transfer key to extract the transaction hash
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
	BalancePrefix  = KeyPrefix(BalanceKeyPrefix)
	TransferPrefix = KeyPrefix(TransferKeyPrefix)
	SupplyPrefix   = KeyPrefix(SupplyKey)
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

// ValidateDenom validates a denomination
func ValidateDenom(denom string) error {
	if denom == "" {
		return fmt.Errorf("denomination cannot be empty")
	}
	if len(denom) > 128 {
		return fmt.Errorf("denomination too long")
	}
	return nil
}
