package types

import (
	"bytes"
	"fmt"
	"time"
)

// Store key prefixes
var (
	BridgeKeyPrefix    = []byte{0x01}
	TransferKeyPrefix  = []byte{0x02}
	ValidatorKeyPrefix = []byte{0x03}
	ConfigKeyPrefix    = []byte{0x04}
	FeeKeyPrefix       = []byte{0x05}
	LimitKeyPrefix     = []byte{0x06}
	EventKeyPrefix     = []byte{0x07}
	ParamsKey          = []byte{0x08}
)

// BridgeKey returns the key for a bridge
func BridgeKey(id string) []byte {
	return append(BridgeKeyPrefix, []byte(id)...)
}

// TransferKey returns the key for a transfer
func TransferKey(id string) []byte {
	return append(TransferKeyPrefix, []byte(id)...)
}

// ValidatorKey returns the key for a validator
func ValidatorKey(id string) []byte {
	return append(ValidatorKeyPrefix, []byte(id)...)
}

// ConfigKey returns the key for a bridge config
func ConfigKey(id string) []byte {
	return append(ConfigKeyPrefix, []byte(id)...)
}

// FeeKey returns the key for a bridge fee
func FeeKey(id string) []byte {
	return append(FeeKeyPrefix, []byte(id)...)
}

// LimitKey returns the key for a bridge limit
func LimitKey(id string) []byte {
	return append(LimitKeyPrefix, []byte(id)...)
}

// EventKey returns the key for a bridge event
func EventKey(id string) []byte {
	return append(EventKeyPrefix, []byte(id)...)
}

// BridgeByNameKey returns the key for bridges by name
func BridgeByNameKey(name string) []byte {
	return append(BridgeKeyPrefix, []byte("name:"+name)...)
}

// BridgeByTypeKey returns the key for bridges by type
func BridgeByTypeKey(bridgeType, bridgeID string) []byte {
	return append(append(BridgeKeyPrefix, []byte("type:"+bridgeType)...), []byte(bridgeID)...)
}

// BridgeByStatusKey returns the key for bridges by status
func BridgeByStatusKey(status, bridgeID string) []byte {
	return append(append(BridgeKeyPrefix, []byte("status:"+status)...), []byte(bridgeID)...)
}

// BridgeByChainKey returns the key for bridges by chain
func BridgeByChainKey(fromChain, toChain, bridgeID string) []byte {
	chainPair := fmt.Sprintf("%s-%s", fromChain, toChain)
	return append(append(BridgeKeyPrefix, []byte("chain:"+chainPair)...), []byte(bridgeID)...)
}

// TransferByBridgeKey returns the key for transfers by bridge
func TransferByBridgeKey(bridgeID, transferID string) []byte {
	return append(append(TransferKeyPrefix, []byte("bridge:"+bridgeID)...), []byte(transferID)...)
}

// TransferByStatusKey returns the key for transfers by status
func TransferByStatusKey(status, transferID string) []byte {
	return append(append(TransferKeyPrefix, []byte("status:"+status)...), []byte(transferID)...)
}

// TransferByChainKey returns the key for transfers by chain
func TransferByChainKey(fromChain, toChain, transferID string) []byte {
	chainPair := fmt.Sprintf("%s-%s", fromChain, toChain)
	return append(append(TransferKeyPrefix, []byte("chain:"+chainPair)...), []byte(transferID)...)
}

// TransferByAddressKey returns the key for transfers by address
func TransferByAddressKey(address, transferID string) []byte {
	return append(append(TransferKeyPrefix, []byte("address:"+address)...), []byte(transferID)...)
}

// TransferByTokenKey returns the key for transfers by token
func TransferByTokenKey(token, transferID string) []byte {
	return append(append(TransferKeyPrefix, []byte("token:"+token)...), []byte(transferID)...)
}

// TransferByCreatedTimeKey returns the key for transfers by creation time
func TransferByCreatedTimeKey(createdAt time.Time, transferID string) []byte {
	return append(append(TransferKeyPrefix, []byte("created:"+createdAt.Format(time.RFC3339))...), []byte(transferID)...)
}

// TransferByAmountKey returns the key for transfers by amount range
func TransferByAmountKey(minAmount, maxAmount, transferID string) []byte {
	amountRange := fmt.Sprintf("%s-%s", minAmount, maxAmount)
	return append(append(TransferKeyPrefix, []byte("amount:"+amountRange)...), []byte(transferID)...)
}

// ValidatorByAddressKey returns the key for validators by address
func ValidatorByAddressKey(address string) []byte {
	return append(ValidatorKeyPrefix, []byte("address:"+address)...)
}

// ValidatorByStatusKey returns the key for validators by status
func ValidatorByStatusKey(status, validatorID string) []byte {
	return append(append(ValidatorKeyPrefix, []byte("status:"+status)...), []byte(validatorID)...)
}

// ValidatorByStakeKey returns the key for validators by stake
func ValidatorByStakeKey(minStake, maxStake, validatorID string) []byte {
	stakeRange := fmt.Sprintf("%s-%s", minStake, maxStake)
	return append(append(ValidatorKeyPrefix, []byte("stake:"+stakeRange)...), []byte(validatorID)...)
}

// ValidatorByCreatedTimeKey returns the key for validators by creation time
func ValidatorByCreatedTimeKey(createdAt time.Time, validatorID string) []byte {
	return append(append(ValidatorKeyPrefix, []byte("created:"+createdAt.Format(time.RFC3339))...), []byte(validatorID)...)
}

// ConfigByBridgeKey returns the key for configs by bridge
func ConfigByBridgeKey(bridgeID, configID string) []byte {
	return append(append(ConfigKeyPrefix, []byte("bridge:"+bridgeID)...), []byte(configID)...)
}

// ConfigByCreatedTimeKey returns the key for configs by creation time
func ConfigByCreatedTimeKey(createdAt time.Time, configID string) []byte {
	return append(append(ConfigKeyPrefix, []byte("created:"+createdAt.Format(time.RFC3339))...), []byte(configID)...)
}

// FeeByBridgeKey returns the key for fees by bridge
func FeeByBridgeKey(bridgeID, feeID string) []byte {
	return append(append(FeeKeyPrefix, []byte("bridge:"+bridgeID)...), []byte(feeID)...)
}

// FeeByTokenKey returns the key for fees by token
func FeeByTokenKey(token, feeID string) []byte {
	return append(append(FeeKeyPrefix, []byte("token:"+token)...), []byte(feeID)...)
}

// FeeByCreatedTimeKey returns the key for fees by creation time
func FeeByCreatedTimeKey(createdAt time.Time, feeID string) []byte {
	return append(append(FeeKeyPrefix, []byte("created:"+createdAt.Format(time.RFC3339))...), []byte(feeID)...)
}

// LimitByBridgeKey returns the key for limits by bridge
func LimitByBridgeKey(bridgeID, limitID string) []byte {
	return append(append(LimitKeyPrefix, []byte("bridge:"+bridgeID)...), []byte(limitID)...)
}

// LimitByTokenKey returns the key for limits by token
func LimitByTokenKey(token, limitID string) []byte {
	return append(append(LimitKeyPrefix, []byte("token:"+token)...), []byte(limitID)...)
}

// LimitByCreatedTimeKey returns the key for limits by creation time
func LimitByCreatedTimeKey(createdAt time.Time, limitID string) []byte {
	return append(append(LimitKeyPrefix, []byte("created:"+createdAt.Format(time.RFC3339))...), []byte(limitID)...)
}

// EventByBridgeKey returns the key for events by bridge
func EventByBridgeKey(bridgeID, eventID string) []byte {
	return append(append(EventKeyPrefix, []byte("bridge:"+bridgeID)...), []byte(eventID)...)
}

// EventByTypeKey returns the key for events by type
func EventByTypeKey(eventType, eventID string) []byte {
	return append(append(EventKeyPrefix, []byte("type:"+eventType)...), []byte(eventID)...)
}

// EventByCreatedTimeKey returns the key for events by creation time
func EventByCreatedTimeKey(createdAt time.Time, eventID string) []byte {
	return append(append(EventKeyPrefix, []byte("created:"+createdAt.Format(time.RFC3339))...), []byte(eventID)...)
}

// GetBridgeIDFromKey extracts bridge ID from key
func GetBridgeIDFromKey(key []byte) string {
	if !bytes.HasPrefix(key, BridgeKeyPrefix) {
		return ""
	}
	return string(key[len(BridgeKeyPrefix):])
}

// GetTransferIDFromKey extracts transfer ID from key
func GetTransferIDFromKey(key []byte) string {
	if !bytes.HasPrefix(key, TransferKeyPrefix) {
		return ""
	}
	return string(key[len(TransferKeyPrefix):])
}

// GetValidatorIDFromKey extracts validator ID from key
func GetValidatorIDFromKey(key []byte) string {
	if !bytes.HasPrefix(key, ValidatorKeyPrefix) {
		return ""
	}
	return string(key[len(ValidatorKeyPrefix):])
}

// GetConfigIDFromKey extracts config ID from key
func GetConfigIDFromKey(key []byte) string {
	if !bytes.HasPrefix(key, ConfigKeyPrefix) {
		return ""
	}
	return string(key[len(ConfigKeyPrefix):])
}

// GetFeeIDFromKey extracts fee ID from key
func GetFeeIDFromKey(key []byte) string {
	if !bytes.HasPrefix(key, FeeKeyPrefix) {
		return ""
	}
	return string(key[len(FeeKeyPrefix):])
}

// GetLimitIDFromKey extracts limit ID from key
func GetLimitIDFromKey(key []byte) string {
	if !bytes.HasPrefix(key, LimitKeyPrefix) {
		return ""
	}
	return string(key[len(LimitKeyPrefix):])
}

// GetEventIDFromKey extracts event ID from key
func GetEventIDFromKey(key []byte) string {
	if !bytes.HasPrefix(key, EventKeyPrefix) {
		return ""
	}
	return string(key[len(EventKeyPrefix):])
}
