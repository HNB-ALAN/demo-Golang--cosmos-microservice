package usc_coin

import (
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/usc_coin/keeper"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/usc_coin/types"
)

// BeginBlocker is called at the beginning of every block
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	// Log block height
	ctx.Logger().Info(fmt.Sprintf("USC BeginBlocker: Block %d", ctx.BlockHeight()))

	// Perform any necessary operations at the beginning of each block
	// This could include:
	// - Updating validator sets
	// - Processing pending transactions
	// - Updating module state
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
	ctx.Logger().Info(fmt.Sprintf("USC EndBlocker: Block %d", ctx.BlockHeight()))

	// Perform any necessary operations at the end of each block
	// This could include:
	// - Finalizing transactions
	// - Updating module state
	// - Processing rewards
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

// InitGenesis initializes the genesis state for the USC module
func InitGenesis(ctx sdk.Context, cdc codec.Codec, k keeper.Keeper, genState types.GenesisState) {
	// Validate genesis state
	if err := genState.ValidateGenesis(); err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to validate genesis state during InitGenesis",
			"module", types.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight(),
			"chain_id", ctx.ChainID())
		panic(fmt.Sprintf("usc_coin: failed to validate genesis state: %s", err.Error()))
	}

	// Set parameters
	if err := k.SetParams(ctx, genState.Params); err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to set parameters during InitGenesis",
			"module", types.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight())
		panic(fmt.Sprintf("usc_coin: failed to set parameters: %s", err.Error()))
	}

	// Set balances
	for _, balance := range genState.Balances {
		if err := k.SetBalance(ctx, balance.Address, balance); err != nil {
			// Log error with detailed context before panic
			ctx.Logger().Error("Failed to set balance during InitGenesis",
				"module", types.ModuleName,
				"address", balance.Address,
				"amount", balance.Amount,
				"error", err.Error(),
				"block_height", ctx.BlockHeight())
			panic(fmt.Sprintf("usc_coin: failed to set balance for address %s: %s", balance.Address, err.Error()))
		}
	}

	// Set transfers
	for _, transfer := range genState.Transfers {
		if err := k.SetTransfer(ctx, transfer); err != nil {
			// Log error with detailed context before panic
			ctx.Logger().Error("Failed to set transfer during InitGenesis",
				"module", types.ModuleName,
				"from_address", transfer.FromAddress,
				"to_address", transfer.ToAddress,
				"amount", transfer.Amount,
				"error", err.Error(),
				"block_height", ctx.BlockHeight())
			panic(fmt.Sprintf("usc_coin: failed to set transfer from %s to %s: %s", transfer.FromAddress, transfer.ToAddress, err.Error()))
		}
	}

	// Set initial total supply
	if len(genState.Balances) > 0 {
		totalSupply := "0"
		for _, balance := range genState.Balances {
			// Add balance amount to total supply
			// This is a simplified calculation
			totalSupply = balance.Amount // For now, just use the last balance amount
		}
		if err := k.SetTotalSupply(ctx, totalSupply); err != nil {
			// Log error with detailed context before panic
			ctx.Logger().Error("Failed to set total supply during InitGenesis",
				"module", types.ModuleName,
				"total_supply", totalSupply,
				"error", err.Error(),
				"block_height", ctx.BlockHeight())
			panic(fmt.Sprintf("usc_coin: failed to set total supply %s: %s", totalSupply, err.Error()))
		}
	}
}

// ExportGenesis exports the genesis state for the USC module
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
		panic(fmt.Sprintf("usc_coin: failed to get parameters: %s", err.Error()))
	}

	// Get all balances
	balances, err := k.GetAllBalances(ctx)
	if err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to get balances during ExportGenesis",
			"module", types.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight())
		panic(fmt.Sprintf("usc_coin: failed to get balances: %s", err.Error()))
	}

	// Get all transfers
	transfers, err := k.GetAllTransfers(ctx)
	if err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to get transfers during ExportGenesis",
			"module", types.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight())
		panic(fmt.Sprintf("usc_coin: failed to get transfers: %s", err.Error()))
	}

	return &types.GenesisState{
		Balances:  balances,
		Transfers: transfers,
		Params:    params,
	}
}

// ValidateGenesis validates the genesis state for the USC module
func ValidateGenesis(genState types.GenesisState) error {
	return genState.ValidateGenesis()
}

// GetGenesisState returns the current genesis state for the USC module
func GetGenesisState(ctx sdk.Context, k keeper.Keeper) (*types.GenesisState, error) {
	// Get parameters
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get parameters: %w", err)
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
		Balances:  balances,
		Transfers: transfers,
		Params:    params,
	}, nil
}
