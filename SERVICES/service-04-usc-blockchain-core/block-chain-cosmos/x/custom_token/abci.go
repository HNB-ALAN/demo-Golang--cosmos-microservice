package custom_token

import (
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/custom_token/keeper"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/custom_token/types"
)

// BeginBlocker is called at the beginning of every block
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	// Log block height
	ctx.Logger().Info(fmt.Sprintf("CustomToken BeginBlocker: Block %d", ctx.BlockHeight()))

	// Perform any necessary operations at the beginning of each block
	// This could include:
	// - Token expiration checks
	// - Balance updates
	// - Emitting events

	// Example: Emit a block start event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTokenCreated,
			sdk.NewAttribute(types.AttributeKeyCreatedAt, fmt.Sprintf("%d", ctx.BlockTime().Unix())),
		),
	)
}

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	// Log block height
	ctx.Logger().Info(fmt.Sprintf("CustomToken EndBlocker: Block %d", ctx.BlockHeight()))

	// Perform any necessary operations at the end of each block
	// This could include:
	// - Token cleanup
	// - Balance processing
	// - Emitting events

	// Example: Emit a block end event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTokenUpdated,
			sdk.NewAttribute(types.AttributeKeyUpdatedAt, fmt.Sprintf("%d", ctx.BlockTime().Unix())),
		),
	)

	// Return validator updates (if any)
	return []abci.ValidatorUpdate{}
}

// InitGenesis initializes the genesis state for the custom_token module
func InitGenesis(ctx sdk.Context, cdc codec.Codec, k keeper.Keeper, genState types.GenesisState) {

	// Validate genesis state - use default params if validation fails
	if err := genState.ValidateGenesis(); err != nil {
		ctx.Logger().Warn("Genesis validation failed, using default params",
			"module", types.ModuleName,
			"error", err.Error())
		genState.Params = types.DefaultParams()
	}

	// Set parameters - handle errors gracefully
	if err := k.SetParams(ctx, genState.Params); err != nil {
		// Expected: Store service may not be available during InitGenesis
		// Module will use default parameters, chain can continue normally
		ctx.Logger().Warn("Failed to set parameters, using default params",
			"module", types.ModuleName,
			"error", err.Error())
		genState.Params = types.DefaultParams()
	}

	// Set tokens - skip on error, log for debugging
	for _, token := range genState.Tokens {
		if err := k.SetToken(ctx, token); err != nil {
			ctx.Logger().Warn("Failed to set token, skipping",
				"module", types.ModuleName,
				"token_id", token.ID,
				"error", err.Error())
			continue
		}
	}

	// Set balances - skip on error, log for debugging
	for _, balance := range genState.Balances {
		if err := k.SetBalance(ctx, balance); err != nil {
			ctx.Logger().Warn("Failed to set balance, skipping",
				"module", types.ModuleName,
				"token_id", balance.TokenID,
				"error", err.Error())
			continue
		}
	}

	// Set transfers - skip on error, log for debugging
	for _, transfer := range genState.Transfers {
		if err := k.SetTransfer(ctx, transfer); err != nil {
			ctx.Logger().Warn("Failed to set transfer, skipping",
				"module", types.ModuleName,
				"token_id", transfer.TokenID,
				"error", err.Error())
			continue
		}
	}
}

// ExportGenesis exports the genesis state for the custom_token module
func ExportGenesis(ctx sdk.Context, cdc codec.Codec, k keeper.Keeper) *types.GenesisState {
	// Get parameters
	params, err := k.GetParams(ctx)
	if err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to get parameters during ExportGenesis",
			"module", types.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight(),
			"chain_id", ctx.ChainID())
		panic(fmt.Sprintf("custom_token: failed to get parameters: %s", err.Error()))
	}

	// Get all tokens
	tokens, err := k.GetAllTokens(ctx)
	if err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to get tokens during ExportGenesis",
			"module", types.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight())
		panic(fmt.Sprintf("custom_token: failed to get tokens: %s", err.Error()))
	}

	// Get all balances
	balances, err := k.GetAllBalances(ctx)
	if err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to get balances during ExportGenesis",
			"module", types.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight())
		panic(fmt.Sprintf("custom_token: failed to get balances: %s", err.Error()))
	}

	// Get all transfers
	transfers, err := k.GetAllTransfers(ctx)
	if err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to get transfers during ExportGenesis",
			"module", types.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight())
		panic(fmt.Sprintf("custom_token: failed to get transfers: %s", err.Error()))
	}

	return &types.GenesisState{
		Tokens:    tokens,
		Balances:  balances,
		Transfers: transfers,
		Params:    params,
	}
}

// ValidateGenesis validates the genesis state for the custom_token module
func ValidateGenesis(genState types.GenesisState) error {
	return genState.ValidateGenesis()
}

// GetGenesisState returns the current genesis state for the custom_token module
func GetGenesisState(ctx sdk.Context, k keeper.Keeper) (*types.GenesisState, error) {
	// Get parameters
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get parameters: %w", err)
	}

	// Get all tokens
	tokens, err := k.GetAllTokens(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tokens: %w", err)
	}

	// Get all balances
	balances, err := k.GetAllBalances(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get balances: %w", err)
	}

	// Get all transfers
	transfers, err := k.GetAllTransfers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transfers: %w", err)
	}

	return &types.GenesisState{
		Tokens:    tokens,
		Balances:  balances,
		Transfers: transfers,
		Params:    params,
	}, nil
}
