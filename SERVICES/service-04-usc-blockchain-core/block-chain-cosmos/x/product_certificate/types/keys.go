package types

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Store key prefixes
const (
	// CertificateKeyPrefix is the prefix for certificate keys
	CertificateKeyPrefix = "certificate:"

	// VerificationKeyPrefix is the prefix for verification keys
	VerificationKeyPrefix = "verification:"

	// ParamsKey is the key for module parameters
	ParamsKey = "params"
)

// GetCertificateKey returns the key for a certificate
func GetCertificateKey(id string) []byte {
	return append([]byte(CertificateKeyPrefix), []byte(id)...)
}

// GetVerificationKey returns the key for a verification
func GetVerificationKey(certificateID, verifier string) []byte {
	key := fmt.Sprintf("%s_%s", certificateID, verifier)
	return append([]byte(VerificationKeyPrefix), []byte(key)...)
}

// GetParamsKey returns the key for module parameters
func GetParamsKey() []byte {
	return []byte(ParamsKey)
}

// ParseCertificateKey parses a certificate key to extract the ID
func ParseCertificateKey(key []byte) (string, error) {
	if !bytes.HasPrefix(key, []byte(CertificateKeyPrefix)) {
		return "", fmt.Errorf("invalid certificate key prefix")
	}
	return string(key[len(CertificateKeyPrefix):]), nil
}

// ParseVerificationKey parses a verification key to extract the IDs
func ParseVerificationKey(key []byte) (string, string, error) {
	if !bytes.HasPrefix(key, []byte(VerificationKeyPrefix)) {
		return "", "", fmt.Errorf("invalid verification key prefix")
	}

	keyStr := string(key[len(VerificationKeyPrefix):])
	parts := bytes.Split([]byte(keyStr), []byte("_"))
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid verification key format")
	}

	return string(parts[0]), string(parts[1]), nil
}

// KeyPrefix represents a key prefix
type KeyPrefix []byte

// Key prefixes for different data types
var (
	CertificatePrefix  = KeyPrefix(CertificateKeyPrefix)
	VerificationPrefix = KeyPrefix(VerificationKeyPrefix)
	ParamsPrefix       = KeyPrefix(ParamsKey)
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
