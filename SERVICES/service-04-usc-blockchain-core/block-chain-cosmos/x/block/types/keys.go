package types

import (
	"fmt"
	"strconv"
)

// Key prefixes for different data types
const (
	BlockKeyPrefix       = "block:"
	BlockHeightKeyPrefix = "block:height:"
	BlockHashKeyPrefix   = "block:hash:"
	BlockDataKeyPrefix   = "block_data:"
	ValidationKeyPrefix  = "validation:"
	ParamsKeyPrefix      = "params:"
)

// Event types
const (
	EventTypeBlockValidation   = "block_validation"
	EventTypeBlockStateUpdate  = "block_state_update"
	EventTypeBlockEvent        = "block_event"
	EventTypeBlockFinalization = "block_finalization"
	EventTypeBlockCleanup      = "block_cleanup"
	EventTypeBlockEnd          = "block_end"
)

// Additional attribute keys
const (
	AttributeKeyHeight    = "height"
	AttributeKeyHash      = "hash"
	AttributeKeyValidator = "validator"
	AttributeKeyStatus    = "status"
)

// GetBlockKey returns the key for a block
func GetBlockKey(blockID string) []byte {
	return []byte(fmt.Sprintf("%s%s", BlockKeyPrefix, blockID))
}

// GetBlockHeightKey returns the key for a block by height
func GetBlockHeightKey(height int64) []byte {
	return []byte(fmt.Sprintf("%s%d", BlockHeightKeyPrefix, height))
}

// GetBlockHashKey returns the key for a block by hash
func GetBlockHashKey(hash string) []byte {
	return []byte(fmt.Sprintf("%s%s", BlockHashKeyPrefix, hash))
}

// GetBlockDataKey returns the key for block data
func GetBlockDataKey(blockID string) []byte {
	return []byte(fmt.Sprintf("%s%s", BlockDataKeyPrefix, blockID))
}

// GetBlockDataHeightKey returns the key for block data by height
func GetBlockDataHeightKey(height int64) []byte {
	return []byte(fmt.Sprintf("%s%d", BlockDataKeyPrefix, height))
}

// GetValidationKey returns the key for block validation
func GetValidationKey(blockID string) []byte {
	return []byte(fmt.Sprintf("%s%s", ValidationKeyPrefix, blockID))
}

// GetValidationHeightKey returns the key for validation by height
func GetValidationHeightKey(height int64) []byte {
	return []byte(fmt.Sprintf("%s%d", ValidationKeyPrefix, height))
}

// GetParamsKey returns the key for module parameters
func GetParamsKey() []byte {
	return []byte(ParamsKeyPrefix)
}

// ParseBlockKey parses a block key to extract block ID
func ParseBlockKey(key []byte) (string, error) {
	keyStr := string(key)
	if len(keyStr) <= len(BlockKeyPrefix) {
		return "", fmt.Errorf("invalid block key")
	}
	return keyStr[len(BlockKeyPrefix):], nil
}

// ParseBlockDataKey parses a block data key to extract block ID
func ParseBlockDataKey(key []byte) (string, error) {
	keyStr := string(key)
	if len(keyStr) <= len(BlockDataKeyPrefix) {
		return "", fmt.Errorf("invalid block data key")
	}
	return keyStr[len(BlockDataKeyPrefix):], nil
}

// ParseValidationKey parses a validation key to extract block ID
func ParseValidationKey(key []byte) (string, error) {
	keyStr := string(key)
	if len(keyStr) <= len(ValidationKeyPrefix) {
		return "", fmt.Errorf("invalid validation key")
	}
	return keyStr[len(ValidationKeyPrefix):], nil
}

// GetBlockHeightFromKey extracts height from a height-based key
func GetBlockHeightFromKey(key []byte) (int64, error) {
	keyStr := string(key)
	if len(keyStr) <= len(BlockHeightKeyPrefix) {
		return 0, fmt.Errorf("invalid block height key")
	}
	heightStr := keyStr[len(BlockHeightKeyPrefix):]
	return strconv.ParseInt(heightStr, 10, 64)
}

// GetBlockDataHeightFromKey extracts height from a block data height key
func GetBlockDataHeightFromKey(key []byte) (int64, error) {
	keyStr := string(key)
	if len(keyStr) <= len(BlockDataKeyPrefix) {
		return 0, fmt.Errorf("invalid block data height key")
	}
	heightStr := keyStr[len(BlockDataKeyPrefix):]
	return strconv.ParseInt(heightStr, 10, 64)
}

// GetValidationHeightFromKey extracts height from a validation height key
func GetValidationHeightFromKey(key []byte) (int64, error) {
	keyStr := string(key)
	if len(keyStr) <= len(ValidationKeyPrefix) {
		return 0, fmt.Errorf("invalid validation height key")
	}
	heightStr := keyStr[len(ValidationKeyPrefix):]
	return strconv.ParseInt(heightStr, 10, 64)
}
