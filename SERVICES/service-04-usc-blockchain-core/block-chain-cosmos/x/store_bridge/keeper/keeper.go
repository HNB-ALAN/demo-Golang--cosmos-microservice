package keeper

import (
	"encoding/json"
	"fmt"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	bridgetypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/store_bridge/types"
)

// Keeper manages the bridge module state
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
}

// NewKeeper creates a new bridge keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
	}
}

// GetBridge returns a bridge by its ID
func (k Keeper) GetBridge(ctx sdk.Context, id string) (bridgetypes.Bridge, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(bridgetypes.BridgeKey(id))
	if bz == nil {
		return bridgetypes.Bridge{}, fmt.Errorf("bridge with ID %s not found", id)
	}

	var bridge bridgetypes.Bridge
	if err := json.Unmarshal(bz, &bridge); err != nil {
		return bridgetypes.Bridge{}, fmt.Errorf("failed to unmarshal bridge: %w", err)
	}

	return bridge, nil
}

// SetBridge sets a bridge
func (k Keeper) SetBridge(ctx sdk.Context, bridge bridgetypes.Bridge) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(bridge)
	if err != nil {
		return fmt.Errorf("failed to marshal bridge: %w", err)
	}
	store.Set(bridgetypes.BridgeKey(bridge.ID), bz)
	return nil
}

// GetAllBridges returns all bridges
func (k Keeper) GetAllBridges(ctx sdk.Context) []bridgetypes.Bridge {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, bridgetypes.BridgeKeyPrefix)
	defer iterator.Close()

	var bridges []bridgetypes.Bridge
	for ; iterator.Valid(); iterator.Next() {
		var bridge bridgetypes.Bridge
		if err := json.Unmarshal(iterator.Value(), &bridge); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		bridges = append(bridges, bridge)
	}
	return bridges
}

// GetTransfer returns a transfer by its ID
func (k Keeper) GetTransfer(ctx sdk.Context, id string) (bridgetypes.Transfer, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(bridgetypes.TransferKey(id))
	if bz == nil {
		return bridgetypes.Transfer{}, fmt.Errorf("transfer with ID %s not found", id)
	}

	var transfer bridgetypes.Transfer
	if err := json.Unmarshal(bz, &transfer); err != nil {
		return bridgetypes.Transfer{}, fmt.Errorf("failed to unmarshal transfer: %w", err)
	}

	return transfer, nil
}

// SetTransfer sets a transfer
func (k Keeper) SetTransfer(ctx sdk.Context, transfer bridgetypes.Transfer) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(transfer)
	if err != nil {
		return fmt.Errorf("failed to marshal transfer: %w", err)
	}
	store.Set(bridgetypes.TransferKey(transfer.ID), bz)
	return nil
}

// GetAllTransfers returns all transfers
func (k Keeper) GetAllTransfers(ctx sdk.Context) []bridgetypes.Transfer {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, bridgetypes.TransferKeyPrefix)
	defer iterator.Close()

	var transfers []bridgetypes.Transfer
	for ; iterator.Valid(); iterator.Next() {
		var transfer bridgetypes.Transfer
		if err := json.Unmarshal(iterator.Value(), &transfer); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		transfers = append(transfers, transfer)
	}
	return transfers
}

// GetValidator returns a validator by its ID
func (k Keeper) GetValidator(ctx sdk.Context, id string) (bridgetypes.Validator, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(bridgetypes.ValidatorKey(id))
	if bz == nil {
		return bridgetypes.Validator{}, fmt.Errorf("validator with ID %s not found", id)
	}

	var validator bridgetypes.Validator
	if err := json.Unmarshal(bz, &validator); err != nil {
		return bridgetypes.Validator{}, fmt.Errorf("failed to unmarshal validator: %w", err)
	}

	return validator, nil
}

// SetValidator sets a validator
func (k Keeper) SetValidator(ctx sdk.Context, validator bridgetypes.Validator) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(validator)
	if err != nil {
		return fmt.Errorf("failed to marshal validator: %w", err)
	}
	store.Set(bridgetypes.ValidatorKey(validator.ID), bz)
	return nil
}

// GetAllValidators returns all validators
func (k Keeper) GetAllValidators(ctx sdk.Context) []bridgetypes.Validator {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, bridgetypes.ValidatorKeyPrefix)
	defer iterator.Close()

	var validators []bridgetypes.Validator
	for ; iterator.Valid(); iterator.Next() {
		var validator bridgetypes.Validator
		if err := json.Unmarshal(iterator.Value(), &validator); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		validators = append(validators, validator)
	}
	return validators
}

// GetBridgeConfig returns a bridge config by its ID
func (k Keeper) GetBridgeConfig(ctx sdk.Context, id string) (bridgetypes.BridgeConfig, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(bridgetypes.ConfigKey(id))
	if bz == nil {
		return bridgetypes.BridgeConfig{}, fmt.Errorf("bridge config with ID %s not found", id)
	}

	var config bridgetypes.BridgeConfig
	if err := json.Unmarshal(bz, &config); err != nil {
		return bridgetypes.BridgeConfig{}, fmt.Errorf("failed to unmarshal bridge config: %w", err)
	}

	return config, nil
}

// SetBridgeConfig sets a bridge config
func (k Keeper) SetBridgeConfig(ctx sdk.Context, config bridgetypes.BridgeConfig) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal bridge config: %w", err)
	}
	store.Set(bridgetypes.ConfigKey(config.ID), bz)
	return nil
}

// GetAllBridgeConfigs returns all bridge configs
func (k Keeper) GetAllBridgeConfigs(ctx sdk.Context) []bridgetypes.BridgeConfig {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, bridgetypes.ConfigKeyPrefix)
	defer iterator.Close()

	var configs []bridgetypes.BridgeConfig
	for ; iterator.Valid(); iterator.Next() {
		var config bridgetypes.BridgeConfig
		if err := json.Unmarshal(iterator.Value(), &config); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		configs = append(configs, config)
	}
	return configs
}

// GetBridgeFee returns a bridge fee by its ID
func (k Keeper) GetBridgeFee(ctx sdk.Context, id string) (bridgetypes.BridgeFee, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(bridgetypes.FeeKey(id))
	if bz == nil {
		return bridgetypes.BridgeFee{}, fmt.Errorf("bridge fee with ID %s not found", id)
	}

	var fee bridgetypes.BridgeFee
	if err := json.Unmarshal(bz, &fee); err != nil {
		return bridgetypes.BridgeFee{}, fmt.Errorf("failed to unmarshal bridge fee: %w", err)
	}

	return fee, nil
}

// SetBridgeFee sets a bridge fee
func (k Keeper) SetBridgeFee(ctx sdk.Context, fee bridgetypes.BridgeFee) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(fee)
	if err != nil {
		return fmt.Errorf("failed to marshal bridge fee: %w", err)
	}
	store.Set(bridgetypes.FeeKey(fee.ID), bz)
	return nil
}

// GetAllBridgeFees returns all bridge fees
func (k Keeper) GetAllBridgeFees(ctx sdk.Context) []bridgetypes.BridgeFee {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, bridgetypes.FeeKeyPrefix)
	defer iterator.Close()

	var fees []bridgetypes.BridgeFee
	for ; iterator.Valid(); iterator.Next() {
		var fee bridgetypes.BridgeFee
		if err := json.Unmarshal(iterator.Value(), &fee); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		fees = append(fees, fee)
	}
	return fees
}

// GetBridgeLimit returns a bridge limit by its ID
func (k Keeper) GetBridgeLimit(ctx sdk.Context, id string) (bridgetypes.BridgeLimit, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(bridgetypes.LimitKey(id))
	if bz == nil {
		return bridgetypes.BridgeLimit{}, fmt.Errorf("bridge limit with ID %s not found", id)
	}

	var limit bridgetypes.BridgeLimit
	if err := json.Unmarshal(bz, &limit); err != nil {
		return bridgetypes.BridgeLimit{}, fmt.Errorf("failed to unmarshal bridge limit: %w", err)
	}

	return limit, nil
}

// SetBridgeLimit sets a bridge limit
func (k Keeper) SetBridgeLimit(ctx sdk.Context, limit bridgetypes.BridgeLimit) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(limit)
	if err != nil {
		return fmt.Errorf("failed to marshal bridge limit: %w", err)
	}
	store.Set(bridgetypes.LimitKey(limit.ID), bz)
	return nil
}

// GetAllBridgeLimits returns all bridge limits
func (k Keeper) GetAllBridgeLimits(ctx sdk.Context) []bridgetypes.BridgeLimit {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, bridgetypes.LimitKeyPrefix)
	defer iterator.Close()

	var limits []bridgetypes.BridgeLimit
	for ; iterator.Valid(); iterator.Next() {
		var limit bridgetypes.BridgeLimit
		if err := json.Unmarshal(iterator.Value(), &limit); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		limits = append(limits, limit)
	}
	return limits
}

// GetBridgeEvent returns a bridge event by its ID
func (k Keeper) GetBridgeEvent(ctx sdk.Context, id string) (bridgetypes.BridgeEvent, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(bridgetypes.EventKey(id))
	if bz == nil {
		return bridgetypes.BridgeEvent{}, fmt.Errorf("bridge event with ID %s not found", id)
	}

	var event bridgetypes.BridgeEvent
	if err := json.Unmarshal(bz, &event); err != nil {
		return bridgetypes.BridgeEvent{}, fmt.Errorf("failed to unmarshal bridge event: %w", err)
	}

	return event, nil
}

// SetBridgeEvent sets a bridge event
func (k Keeper) SetBridgeEvent(ctx sdk.Context, event bridgetypes.BridgeEvent) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal bridge event: %w", err)
	}
	store.Set(bridgetypes.EventKey(event.ID), bz)
	return nil
}

// GetAllBridgeEvents returns all bridge events
func (k Keeper) GetAllBridgeEvents(ctx sdk.Context) []bridgetypes.BridgeEvent {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, bridgetypes.EventKeyPrefix)
	defer iterator.Close()

	var events []bridgetypes.BridgeEvent
	for ; iterator.Valid(); iterator.Next() {
		var event bridgetypes.BridgeEvent
		if err := json.Unmarshal(iterator.Value(), &event); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		events = append(events, event)
	}
	return events
}

// GetParams returns the bridge module's parameters
func (k Keeper) GetParams(ctx sdk.Context) bridgetypes.Params {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(bridgetypes.ParamsKey)
	if bz == nil {
		return bridgetypes.DefaultParams()
	}

	var params bridgetypes.Params
	if err := json.Unmarshal(bz, &params); err != nil {
		return bridgetypes.DefaultParams()
	}

	return params
}

// SetParams sets the bridge module's parameters
func (k Keeper) SetParams(ctx sdk.Context, params bridgetypes.Params) {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(params)
	if err != nil {
		// Handle error appropriately
		return
	}
	store.Set(bridgetypes.ParamsKey, bz)
}
