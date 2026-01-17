// Package utils provides utility functions for USC platform services.
package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

// CryptoUtils provides cryptographic utility functions
type CryptoUtils struct{}

// NewCryptoUtils creates a new crypto utils instance
func NewCryptoUtils() *CryptoUtils {
	return &CryptoUtils{}
}

// HashMD5 generates MD5 hash
func (cu *CryptoUtils) HashMD5(data string) string {
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

// HashSHA1 generates SHA1 hash
func (cu *CryptoUtils) HashSHA1(data string) string {
	hash := sha1.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

// HashSHA256 generates SHA256 hash
func (cu *CryptoUtils) HashSHA256(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// HashSHA512 generates SHA512 hash
func (cu *CryptoUtils) HashSHA512(data string) string {
	hash := sha512.Sum512([]byte(data))
	return hex.EncodeToString(hash[:])
}

// HashSHA256WithSalt generates SHA256 hash with salt
func (cu *CryptoUtils) HashSHA256WithSalt(data, salt string) string {
	hash := sha256.Sum256([]byte(data + salt))
	return hex.EncodeToString(hash[:])
}

// HashSHA512WithSalt generates SHA512 hash with salt
func (cu *CryptoUtils) HashSHA512WithSalt(data, salt string) string {
	hash := sha512.Sum512([]byte(data + salt))
	return hex.EncodeToString(hash[:])
}

// GenerateSalt generates a random salt
func (cu *CryptoUtils) GenerateSalt(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateRandomBytes generates random bytes
func (cu *CryptoUtils) GenerateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	return bytes, nil
}

// GenerateRandomString generates a random string
func (cu *CryptoUtils) GenerateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}
	return string(bytes), nil
}

// GenerateRandomHex generates a random hex string
func (cu *CryptoUtils) GenerateRandomHex(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// EncryptAES encrypts data using AES
func (cu *CryptoUtils) EncryptAES(data, key string) (string, error) {
	// Convert key to 32 bytes for AES-256
	keyBytes := []byte(key)
	if len(keyBytes) != 32 {
		return "", fmt.Errorf("key must be 32 bytes long")
	}

	// Create cipher
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt data
	ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)

	// Encode to base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptAES decrypts data using AES
func (cu *CryptoUtils) DecryptAES(encryptedData, key string) (string, error) {
	// Convert key to 32 bytes for AES-256
	keyBytes := []byte(key)
	if len(keyBytes) != 32 {
		return "", fmt.Errorf("key must be 32 bytes long")
	}

	// Decode from base64
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	// Create cipher
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Extract nonce
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt data
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// GenerateHMAC generates HMAC
func (cu *CryptoUtils) GenerateHMAC(data, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// VerifyHMAC verifies HMAC
func (cu *CryptoUtils) VerifyHMAC(data, key, expectedHMAC string) bool {
	actualHMAC := cu.GenerateHMAC(data, key)
	return hmac.Equal([]byte(actualHMAC), []byte(expectedHMAC))
}

// EncodeBase64 encodes data to base64
func (cu *CryptoUtils) EncodeBase64(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// DecodeBase64 decodes data from base64
func (cu *CryptoUtils) DecodeBase64(encodedData string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// EncodeBase64URL encodes data to base64 URL
func (cu *CryptoUtils) EncodeBase64URL(data string) string {
	return base64.URLEncoding.EncodeToString([]byte(data))
}

// DecodeBase64URL decodes data from base64 URL
func (cu *CryptoUtils) DecodeBase64URL(encodedData string) (string, error) {
	decoded, err := base64.URLEncoding.DecodeString(encodedData)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// EncodeHex encodes data to hex
func (cu *CryptoUtils) EncodeHex(data string) string {
	return hex.EncodeToString([]byte(data))
}

// DecodeHex decodes data from hex
func (cu *CryptoUtils) DecodeHex(encodedData string) (string, error) {
	decoded, err := hex.DecodeString(encodedData)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// XOR encrypts/decrypts data using XOR
func (cu *CryptoUtils) XOR(data, key string) string {
	result := make([]byte, len(data))
	keyBytes := []byte(key)

	for i := 0; i < len(data); i++ {
		result[i] = data[i] ^ keyBytes[i%len(keyBytes)]
	}

	return string(result)
}

// CaesarCipher encrypts/decrypts data using Caesar cipher
func (cu *CryptoUtils) CaesarCipher(data string, shift int) string {
	result := make([]byte, len(data))

	for i := 0; i < len(data); i++ {
		char := data[i]
		if char >= 'a' && char <= 'z' {
			result[i] = byte((int(char-'a')+shift)%26 + 'a')
		} else if char >= 'A' && char <= 'Z' {
			result[i] = byte((int(char-'A')+shift)%26 + 'A')
		} else {
			result[i] = char
		}
	}

	return string(result)
}

// VigenereCipher encrypts/decrypts data using Vigenère cipher
func (cu *CryptoUtils) VigenereCipher(data, key string) string {
	result := make([]byte, len(data))
	keyBytes := []byte(key)
	keyIndex := 0

	for i := 0; i < len(data); i++ {
		char := data[i]
		if char >= 'a' && char <= 'z' {
			shift := int(keyBytes[keyIndex%len(keyBytes)] - 'a')
			result[i] = byte((int(char-'a')+shift)%26 + 'a')
			keyIndex++
		} else if char >= 'A' && char <= 'Z' {
			shift := int(keyBytes[keyIndex%len(keyBytes)] - 'A')
			result[i] = byte((int(char-'A')+shift)%26 + 'A')
			keyIndex++
		} else {
			result[i] = char
		}
	}

	return string(result)
}

// PasswordHash hashes a password with salt
func (cu *CryptoUtils) PasswordHash(password string) (string, error) {
	// Generate salt
	salt, err := cu.GenerateSalt(16)
	if err != nil {
		return "", err
	}

	// Hash password with salt
	hash := cu.HashSHA256WithSalt(password, salt)

	// Return salt and hash combined
	return salt + ":" + hash, nil
}

// PasswordVerify verifies a password against a hash
func (cu *CryptoUtils) PasswordVerify(password, hash string) bool {
	// Split salt and hash
	parts := splitString(hash, ":")
	if len(parts) != 2 {
		return false
	}

	salt := parts[0]
	expectedHash := parts[1]

	// Hash password with salt
	actualHash := cu.HashSHA256WithSalt(password, salt)

	// Compare hashes
	return actualHash == expectedHash
}

// GenerateAPIKey generates an API key
func (cu *CryptoUtils) GenerateAPIKey() (string, error) {
	// Generate random bytes
	bytes, err := cu.GenerateRandomBytes(32)
	if err != nil {
		return "", err
	}

	// Encode to base64 URL
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// GenerateToken generates a token
func (cu *CryptoUtils) GenerateToken() (string, error) {
	// Generate random bytes
	bytes, err := cu.GenerateRandomBytes(16)
	if err != nil {
		return "", err
	}

	// Encode to hex
	return hex.EncodeToString(bytes), nil
}

// GenerateUUID generates a UUID v4
func (cu *CryptoUtils) GenerateUUID() (string, error) {
	// Generate random bytes
	bytes, err := cu.GenerateRandomBytes(16)
	if err != nil {
		return "", err
	}

	// Set version (4) and variant bits
	bytes[6] = (bytes[6] & 0x0f) | 0x40 // Version 4
	bytes[8] = (bytes[8] & 0x3f) | 0x80 // Variant bits

	// Format as UUID
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:16]), nil
}

// ValidateUUID validates a UUID
func (cu *CryptoUtils) ValidateUUID(uuid string) bool {
	if len(uuid) != 36 {
		return false
	}

	// Check format
	if uuid[8] != '-' || uuid[13] != '-' || uuid[18] != '-' || uuid[23] != '-' {
		return false
	}

	// Check if all characters are valid hex
	for i, char := range uuid {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if char != '-' {
				return false
			}
		} else {
			if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
				return false
			}
		}
	}

	return true
}

// splitString splits a string by delimiter
func splitString(s, delimiter string) []string {
	var result []string
	start := 0

	for i := 0; i < len(s); i++ {
		if i+len(delimiter) <= len(s) && s[i:i+len(delimiter)] == delimiter {
			result = append(result, s[start:i])
			start = i + len(delimiter)
			i += len(delimiter) - 1
		}
	}

	result = append(result, s[start:])
	return result
}
