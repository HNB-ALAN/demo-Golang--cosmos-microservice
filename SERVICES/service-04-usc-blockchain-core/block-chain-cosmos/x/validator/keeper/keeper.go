package keeper

import (
	"encoding/json"
	"fmt"
	"strconv"

	storetypes "cosmossdk.io/store/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/validator/types"
)

// Keeper manages the validator module state
type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	paramSpace paramtypes.Subspace
}

// NewKeeper creates a new validator keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	paramSpace paramtypes.Subspace,
) Keeper {
	// Set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		paramSpace: paramSpace,
	}
}

// GetParams returns the current parameters
// COSMOS SDK 0.53.4: Handle nil paramSpace or uninitialized params gracefully
func (k Keeper) GetParams(ctx sdk.Context) (types.Params, error) {
	// Check if paramSpace has key table
	if !k.paramSpace.HasKeyTable() {
		// Return default params if paramSpace is not initialized
		return types.DefaultParams(), nil
	}

	// Try to get params from paramSpace
	// Use recover to handle panic if params are not set
	var params types.Params
	func() {
		defer func() {
			if r := recover(); r != nil {
				// Params not set, use default params
				params = types.DefaultParams()
			}
		}()
		k.paramSpace.GetParamSet(ctx, &params)
	}()

	// If params are zero/empty, return default params
	if params.MinDelegation == "" {
		return types.DefaultParams(), nil
	}

	return params, nil
}

// SetParams sets the parameters
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}
	k.paramSpace.SetParamSet(ctx, &params)
	return nil
}

// GetValidator returns a validator by address
func (k Keeper) GetValidator(ctx sdk.Context, address string) (types.Validator, error) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetValidatorKey(address)

	if !store.Has(key) {
		return types.Validator{}, fmt.Errorf("validator not found: %s", address)
	}

	var validator types.Validator
	if err := json.Unmarshal(store.Get(key), &validator); err != nil {
		return types.Validator{}, err
	}

	return validator, nil
}

// SetValidator sets a validator
func (k Keeper) SetValidator(ctx sdk.Context, validator types.Validator) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetValidatorKey(validator.Address)

	bz, err := json.Marshal(&validator)
	if err != nil {
		return err
	}

	store.Set(key, bz)
	return nil
}

// GetAllValidators returns all validators
func (k Keeper) GetAllValidators(ctx sdk.Context) ([]types.Validator, error) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte(types.ValidatorKeyPrefix))
	defer iterator.Close()

	var validators []types.Validator
	for ; iterator.Valid(); iterator.Next() {
		var validator types.Validator
		if err := json.Unmarshal(iterator.Value(), &validator); err != nil {
			return nil, err
		}
		validators = append(validators, validator)
	}

	return validators, nil
}

// GetDelegation returns a delegation
func (k Keeper) GetDelegation(ctx sdk.Context, delegatorAddress, validatorAddress string) (types.Delegation, error) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetDelegationKey(delegatorAddress, validatorAddress)

	if !store.Has(key) {
		return types.Delegation{}, fmt.Errorf("delegation not found: %s -> %s", delegatorAddress, validatorAddress)
	}

	var delegation types.Delegation
	if err := json.Unmarshal(store.Get(key), &delegation); err != nil {
		return types.Delegation{}, err
	}

	return delegation, nil
}

// SetDelegation sets a delegation
func (k Keeper) SetDelegation(ctx sdk.Context, delegation types.Delegation) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetDelegationKey(delegation.DelegatorAddress, delegation.ValidatorAddress)

	bz, err := json.Marshal(&delegation)
	if err != nil {
		return err
	}

	store.Set(key, bz)
	return nil
}

// GetAllDelegations returns all delegations
func (k Keeper) GetAllDelegations(ctx sdk.Context) ([]types.Delegation, error) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte(types.DelegationKeyPrefix))
	defer iterator.Close()

	var delegations []types.Delegation
	for ; iterator.Valid(); iterator.Next() {
		var delegation types.Delegation
		if err := json.Unmarshal(iterator.Value(), &delegation); err != nil {
			return nil, err
		}
		delegations = append(delegations, delegation)
	}

	return delegations, nil
}

// CreateValidator creates a new validator
func (k Keeper) CreateValidator(ctx sdk.Context, validator types.Validator) error {
	// Check if validator already exists
	if _, err := k.GetValidator(ctx, validator.Address); err == nil {
		return fmt.Errorf("validator already exists: %s", validator.Address)
	}

	// Validate validator
	if validator.Address == "" {
		return fmt.Errorf("validator address cannot be empty")
	}
	if validator.PubKey == "" {
		return fmt.Errorf("validator pub key cannot be empty")
	}
	if validator.Description == "" {
		return fmt.Errorf("validator description cannot be empty")
	}
	if validator.Commission == "" {
		return fmt.Errorf("validator commission cannot be empty")
	}

	// Set validator
	return k.SetValidator(ctx, validator)
}

// UpdateValidator updates an existing validator
func (k Keeper) UpdateValidator(ctx sdk.Context, validator types.Validator) error {
	// Check if validator exists
	if _, err := k.GetValidator(ctx, validator.Address); err != nil {
		return fmt.Errorf("validator not found: %s", validator.Address)
	}

	// Set validator
	return k.SetValidator(ctx, validator)
}

// RemoveValidator removes a validator
func (k Keeper) RemoveValidator(ctx sdk.Context, address string) error {
	// Check if validator exists
	if _, err := k.GetValidator(ctx, address); err != nil {
		return fmt.Errorf("validator not found: %s", address)
	}

	store := ctx.KVStore(k.storeKey)
	key := types.GetValidatorKey(address)
	store.Delete(key)

	return nil
}

// Delegate creates a new delegation
func (k Keeper) Delegate(ctx sdk.Context, delegation types.Delegation) error {
	// Validate delegation
	if delegation.DelegatorAddress == "" {
		return fmt.Errorf("delegator address cannot be empty")
	}
	if delegation.ValidatorAddress == "" {
		return fmt.Errorf("validator address cannot be empty")
	}
	if delegation.Amount == "" {
		return fmt.Errorf("delegation amount cannot be empty")
	}

	// Check if validator exists
	if _, err := k.GetValidator(ctx, delegation.ValidatorAddress); err != nil {
		return fmt.Errorf("validator not found: %s", delegation.ValidatorAddress)
	}

	// Validate amount
	amount, err := strconv.ParseInt(delegation.Amount, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid delegation amount: %w", err)
	}

	// Get parameters
	params, err := k.GetParams(ctx)
	if err != nil {
		return fmt.Errorf("failed to get parameters: %w", err)
	}

	// Check minimum delegation
	minDelegation, err := strconv.ParseInt(params.MinDelegation, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid min delegation: %w", err)
	}

	if amount < minDelegation {
		return fmt.Errorf("delegation amount below minimum: %d < %d", amount, minDelegation)
	}

	// Set delegation
	return k.SetDelegation(ctx, delegation)
}

// Undelegate removes a delegation
func (k Keeper) Undelegate(ctx sdk.Context, delegatorAddress, validatorAddress string) error {
	// Check if delegation exists
	if _, err := k.GetDelegation(ctx, delegatorAddress, validatorAddress); err != nil {
		return fmt.Errorf("delegation not found: %s -> %s", delegatorAddress, validatorAddress)
	}

	store := ctx.KVStore(k.storeKey)
	key := types.GetDelegationKey(delegatorAddress, validatorAddress)
	store.Delete(key)

	return nil
}

// InitGenesis initializes the genesis state
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) error {
	// Set parameters
	if err := k.SetParams(ctx, genState.Params); err != nil {
		return fmt.Errorf("failed to set parameters: %w", err)
	}

	// Set validators
	for _, validator := range genState.Validators {
		if err := k.SetValidator(ctx, validator); err != nil {
			return fmt.Errorf("failed to set validator: %w", err)
		}
	}

	// Set delegations
	for _, delegation := range genState.Delegations {
		if err := k.SetDelegation(ctx, delegation); err != nil {
			return fmt.Errorf("failed to set delegation: %w", err)
		}
	}

	return nil
}

// ExportGenesis exports the genesis state
// Returns genesis state and error (if any)
func (k Keeper) ExportGenesis(ctx sdk.Context) (*types.GenesisState, error) {
	// Get parameters
	params, err := k.GetParams(ctx)
	if err != nil {
		// Log error and return error instead of panic
		ctx.Logger().Error("Failed to get parameters during ExportGenesis",
			"module", types.ModuleName,
			"error", err.Error())
		return nil, fmt.Errorf("failed to get parameters: %w", err)
	}

	// Get all validators
	validators, err := k.GetAllValidators(ctx)
	if err != nil {
		// Log error and return error instead of panic
		ctx.Logger().Error("Failed to get validators during ExportGenesis",
			"module", types.ModuleName,
			"error", err.Error())
		return nil, fmt.Errorf("failed to get validators: %w", err)
	}

	// Get all delegations
	delegations, err := k.GetAllDelegations(ctx)
	if err != nil {
		// Log error and return error instead of panic
		ctx.Logger().Error("Failed to get delegations during ExportGenesis",
			"module", types.ModuleName,
			"error", err.Error())
		return nil, fmt.Errorf("failed to get delegations: %w", err)
	}

	return &types.GenesisState{
		Validators:  validators,
		Delegations: delegations,
		Params:      params,
	}, nil
}

// BeginBlocker is called at the beginning of every block
func (k Keeper) BeginBlocker(ctx sdk.Context) {
	// Begin block logic can be added here if needed
}

// EndBlocker is called at the end of every block
func (k Keeper) EndBlocker(ctx sdk.Context) []abci.ValidatorUpdate {
	// End block logic can be added here if needed
	return []abci.ValidatorUpdate{}
}
