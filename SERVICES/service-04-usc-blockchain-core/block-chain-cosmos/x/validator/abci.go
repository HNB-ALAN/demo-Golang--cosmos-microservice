package validator

import (
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/validator/keeper"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/validator/types"
)

// BeginBlocker is called at the beginning of every block
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	// Log block height
	ctx.Logger().Info(fmt.Sprintf("Validator BeginBlocker: Block %d", ctx.BlockHeight()))

	// Perform any necessary operations at the beginning of each block
	// This could include:
	// - Validator set updates
	// - Staking operations
	// - Delegation updates
	// - Emitting events

	// Example: Emit a block start event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBlockStart,
			sdk.NewAttribute(types.AttributeKeyBlockHeight, fmt.Sprintf("%d", ctx.BlockHeight())),
			sdk.NewAttribute(types.AttributeKeyBlockTime, ctx.BlockTime().String()),
		),
	)
}

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	// Log block height
	ctx.Logger().Info(fmt.Sprintf("Validator EndBlocker: Block %d", ctx.BlockHeight()))

	// Perform any necessary operations at the end of each block
	// This could include:
	// - Validator set updates
	// - Staking rewards distribution
	// - Delegation updates
	// - Emitting events

	// Example: Emit a block end event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBlockEnd,
			sdk.NewAttribute(types.AttributeKeyBlockHeight, fmt.Sprintf("%d", ctx.BlockHeight())),
			sdk.NewAttribute(types.AttributeKeyBlockTime, ctx.BlockTime().String()),
		),
	)

	// Return validator updates (if any)
	return []abci.ValidatorUpdate{}
}

// InitGenesis initializes the genesis state for the validator module
func InitGenesis(ctx sdk.Context, cdc codec.Codec, k keeper.Keeper, genState types.GenesisState) {
	// Validate genesis state
	if err := genState.Validate(); err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to validate genesis state during InitGenesis",
			"module", types.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight(),
			"chain_id", ctx.ChainID())
		panic(fmt.Sprintf("validator: failed to validate genesis state: %s", err.Error()))
	}

	// Set parameters
	if err := k.SetParams(ctx, genState.Params); err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to set parameters during InitGenesis",
			"module", types.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight())
		panic(fmt.Sprintf("validator: failed to set parameters: %s", err.Error()))
	}

	// Set validators
	for _, validator := range genState.Validators {
		if err := k.SetValidator(ctx, validator); err != nil {
			// Log error with detailed context before panic
			ctx.Logger().Error("Failed to set validator during InitGenesis",
				"module", types.ModuleName,
				"validator_address", validator.Address,
				"error", err.Error(),
				"block_height", ctx.BlockHeight())
			panic(fmt.Sprintf("validator: failed to set validator %s: %s", validator.Address, err.Error()))
		}
	}

	// Set delegations
	for _, delegation := range genState.Delegations {
		if err := k.SetDelegation(ctx, delegation); err != nil {
			// Log error with detailed context before panic
			ctx.Logger().Error("Failed to set delegation during InitGenesis",
				"module", types.ModuleName,
				"delegator_address", delegation.DelegatorAddress,
				"validator_address", delegation.ValidatorAddress,
				"error", err.Error(),
				"block_height", ctx.BlockHeight())
			panic(fmt.Sprintf("validator: failed to set delegation from %s to %s: %s", delegation.DelegatorAddress, delegation.ValidatorAddress, err.Error()))
		}
	}
}

// ExportGenesis exports the genesis state for the validator module
// Returns genesis state and error (if any)
func ExportGenesis(ctx sdk.Context, cdc codec.Codec, k keeper.Keeper) (*types.GenesisState, error) {
	// Call keeper's ExportGenesis which now returns error
	genState, err := k.ExportGenesis(ctx)
	if err != nil {
		// Log error and return error instead of panic
		ctx.Logger().Error("Failed to export genesis state",
			"module", types.ModuleName,
			"error", err.Error())
		return nil, fmt.Errorf("failed to export %s genesis state: %w", types.ModuleName, err)
	}

	return genState, nil
}

// ValidateGenesis validates the genesis state for the validator module
func ValidateGenesis(genState types.GenesisState) error {
	return genState.Validate()
}

// GetGenesisState returns the current genesis state for the validator module
func GetGenesisState(ctx sdk.Context, k keeper.Keeper) (*types.GenesisState, error) {
	// Get parameters
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get parameters: %w", err)
	}

	// Get all validators
	validators, err := k.GetAllValidators(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get validators: %w", err)
	}

	// Get all delegations
	delegations, err := k.GetAllDelegations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get delegations: %w", err)
	}

	return &types.GenesisState{
		Validators:  validators,
		Delegations: delegations,
		Params:      params,
	}, nil
}
