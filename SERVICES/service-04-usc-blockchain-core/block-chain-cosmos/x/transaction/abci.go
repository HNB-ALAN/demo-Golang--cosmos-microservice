package transaction

import (
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/transaction/keeper"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/transaction/types"
)

// BeginBlocker processes transaction timeouts and cleanup
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	k.BeginBlocker(ctx)
}

// EndBlocker processes end of block operations
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	return k.EndBlocker(ctx)
}

// InitGenesis initializes the genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.InitGenesis(ctx, genState)
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return k.ExportGenesis(ctx)
}

// ValidateGenesis validates the genesis state
func ValidateGenesis(genState types.GenesisState) error {
	return types.ValidateGenesis(&genState)
}
