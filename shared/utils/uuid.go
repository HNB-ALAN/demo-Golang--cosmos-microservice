// Package utils provides utility functions for USC platform services.
package utils

import (
	"crypto/rand"
	"fmt"
	"regexp"
	"strings"
)

// UUIDUtils provides UUID utility functions
type UUIDUtils struct{}

// NewUUIDUtils creates a new UUID utils instance
func NewUUIDUtils() *UUIDUtils {
	return &UUIDUtils{}
}

// GenerateUUID generates a UUID v4
func (uu *UUIDUtils) GenerateUUID() (string, error) {
	// Generate random bytes
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Set version (4) and variant bits
	bytes[6] = (bytes[6] & 0x0f) | 0x40 // Version 4
	bytes[8] = (bytes[8] & 0x3f) | 0x80 // Variant bits

	// Format as UUID
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:16]), nil
}

// GenerateUUIDWithoutHyphens generates a UUID v4 without hyphens
func (uu *UUIDUtils) GenerateUUIDWithoutHyphens() (string, error) {
	uuid, err := uu.GenerateUUID()
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(uuid, "-", ""), nil
}

// GenerateShortUUID generates a short UUID (8 characters)
func (uu *UUIDUtils) GenerateShortUUID() (string, error) {
	// Generate random bytes
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Format as short UUID
	return fmt.Sprintf("%x", bytes), nil
}

// GenerateLongUUID generates a long UUID (32 characters without hyphens)
func (uu *UUIDUtils) GenerateLongUUID() (string, error) {
	// Generate random bytes
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Format as long UUID
	return fmt.Sprintf("%x", bytes), nil
}

// ValidateUUID validates a UUID
func (uu *UUIDUtils) ValidateUUID(uuid string) bool {
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

// ValidateUUIDWithoutHyphens validates a UUID without hyphens
func (uu *UUIDUtils) ValidateUUIDWithoutHyphens(uuid string) bool {
	if len(uuid) != 32 {
		return false
	}

	// Check if all characters are valid hex
	for _, char := range uuid {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
			return false
		}
	}

	return true
}

// NormalizeUUID normalizes a UUID to standard format
func (uu *UUIDUtils) NormalizeUUID(uuid string) string {
	// Remove hyphens and convert to lowercase
	normalized := strings.ToLower(strings.ReplaceAll(uuid, "-", ""))

	// Add hyphens in correct positions
	if len(normalized) == 32 {
		return fmt.Sprintf("%s-%s-%s-%s-%s",
			normalized[0:8], normalized[8:12], normalized[12:16],
			normalized[16:20], normalized[20:32])
	}

	return uuid
}

// ExtractUUIDFromString extracts UUID from a string
func (uu *UUIDUtils) ExtractUUIDFromString(text string) []string {
	var uuids []string

	// Pattern for UUID with hyphens
	pattern := `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`

	// Find all matches
	re := regexp.MustCompile(pattern)
	matches := re.FindAllString(text, -1)

	for _, match := range matches {
		if uu.ValidateUUID(match) {
			uuids = append(uuids, match)
		}
	}

	return uuids
}

// ExtractUUIDWithoutHyphensFromString extracts UUID without hyphens from a string
func (uu *UUIDUtils) ExtractUUIDWithoutHyphensFromString(text string) []string {
	var uuids []string

	// Pattern for UUID without hyphens
	pattern := `[0-9a-fA-F]{32}`

	// Find all matches
	re := regexp.MustCompile(pattern)
	matches := re.FindAllString(text, -1)

	for _, match := range matches {
		if uu.ValidateUUIDWithoutHyphens(match) {
			uuids = append(uuids, match)
		}
	}

	return uuids
}

// IsValidUUID checks if a string is a valid UUID
func (uu *UUIDUtils) IsValidUUID(uuid string) bool {
	return uu.ValidateUUID(uuid) || uu.ValidateUUIDWithoutHyphens(uuid)
}

// GenerateUUIDWithPrefix generates a UUID with a prefix
func (uu *UUIDUtils) GenerateUUIDWithPrefix(prefix string) (string, error) {
	uuid, err := uu.GenerateUUID()
	if err != nil {
		return "", err
	}
	return prefix + "_" + uuid, nil
}

// GenerateUUIDWithSuffix generates a UUID with a suffix
func (uu *UUIDUtils) GenerateUUIDWithSuffix(suffix string) (string, error) {
	uuid, err := uu.GenerateUUID()
	if err != nil {
		return "", err
	}
	return uuid + "_" + suffix, nil
}

// GenerateUUIDWithSeparator generates a UUID with a custom separator
func (uu *UUIDUtils) GenerateUUIDWithSeparator(separator string) (string, error) {
	uuid, err := uu.GenerateUUID()
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(uuid, "-", separator), nil
}

// GenerateMultipleUUIDs generates multiple UUIDs
func (uu *UUIDUtils) GenerateMultipleUUIDs(count int) ([]string, error) {
	var uuids []string

	for i := 0; i < count; i++ {
		uuid, err := uu.GenerateUUID()
		if err != nil {
			return nil, err
		}
		uuids = append(uuids, uuid)
	}

	return uuids, nil
}

// GenerateUniqueUUIDs generates unique UUIDs (no duplicates)
func (uu *UUIDUtils) GenerateUniqueUUIDs(count int) ([]string, error) {
	uuids := make(map[string]bool)
	var result []string

	for len(result) < count {
		uuid, err := uu.GenerateUUID()
		if err != nil {
			return nil, err
		}

		if !uuids[uuid] {
			uuids[uuid] = true
			result = append(result, uuid)
		}
	}

	return result, nil
}

// CompareUUIDs compares two UUIDs
func (uu *UUIDUtils) CompareUUIDs(uuid1, uuid2 string) bool {
	// Normalize both UUIDs
	normalized1 := uu.NormalizeUUID(uuid1)
	normalized2 := uu.NormalizeUUID(uuid2)

	return normalized1 == normalized2
}

// GetUUIDVersion gets the version of a UUID
func (uu *UUIDUtils) GetUUIDVersion(uuid string) (int, error) {
	if !uu.ValidateUUID(uuid) {
		return 0, fmt.Errorf("invalid UUID format")
	}

	// Remove hyphens
	clean := strings.ReplaceAll(uuid, "-", "")

	// Get version from 13th character
	versionChar := clean[12]

	switch versionChar {
	case '1':
		return 1, nil
	case '2':
		return 2, nil
	case '3':
		return 3, nil
	case '4':
		return 4, nil
	case '5':
		return 5, nil
	default:
		return 0, fmt.Errorf("unknown UUID version")
	}
}

// GetUUIDVariant gets the variant of a UUID
func (uu *UUIDUtils) GetUUIDVariant(uuid string) (string, error) {
	if !uu.ValidateUUID(uuid) {
		return "", fmt.Errorf("invalid UUID format")
	}

	// Remove hyphens
	clean := strings.ReplaceAll(uuid, "-", "")

	// Get variant from 17th character
	variantChar := clean[16]

	switch {
	case variantChar >= '0' && variantChar <= '7':
		return "NCS", nil
	case variantChar >= '8' && variantChar <= 'b':
		return "RFC4122", nil
	case variantChar >= 'c' && variantChar <= 'd':
		return "Microsoft", nil
	case variantChar >= 'e' && variantChar <= 'f':
		return "Future", nil
	default:
		return "Unknown", nil
	}
}
