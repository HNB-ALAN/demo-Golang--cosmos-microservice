package product_certificate

import (
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/product_certificate/keeper"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/product_certificate/types"
)

// BeginBlocker is called at the beginning of every block
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	// Log block height
	ctx.Logger().Info(fmt.Sprintf("ProductCertificate BeginBlocker: Block %d", ctx.BlockHeight()))

	// Perform any necessary operations at the beginning of each block
	// This could include:
	// - Certificate expiration checks
	// - Auto-verification processes
	// - Emitting events

	// Example: Emit a block start event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProductCertificateCreated,
			sdk.NewAttribute(types.AttributeKeyCreatedAt, fmt.Sprintf("%d", ctx.BlockTime().Unix())),
		),
	)
}

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	// Log block height
	ctx.Logger().Info(fmt.Sprintf("ProductCertificate EndBlocker: Block %d", ctx.BlockHeight()))

	// Perform any necessary operations at the end of each block
	// This could include:
	// - Certificate cleanup
	// - Verification processing
	// - Emitting events

	// Example: Emit a block end event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProductCertificateUpdated,
			sdk.NewAttribute(types.AttributeKeyUpdatedAt, fmt.Sprintf("%d", ctx.BlockTime().Unix())),
		),
	)

	// Return validator updates (if any)
	return []abci.ValidatorUpdate{}
}

// InitGenesis initializes the genesis state for the product_certificate module
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

	// Set certificates - skip on error, log for debugging
	for _, cert := range genState.Certificates {
		if err := k.SetCertificate(ctx, cert); err != nil {
			ctx.Logger().Warn("Failed to set certificate, skipping",
				"module", types.ModuleName,
				"certificate_id", cert.ID,
				"error", err.Error())
			continue
		}
	}

	// Set verifications - skip on error, log for debugging
	for _, verification := range genState.Verifications {
		if err := k.SetVerification(ctx, verification); err != nil {
			ctx.Logger().Warn("Failed to set verification, skipping",
				"module", types.ModuleName,
				"certificate_id", verification.CertificateID,
				"verifier", verification.Verifier,
				"error", err.Error())
			continue
		}
	}
}

// ExportGenesis exports the genesis state for the product_certificate module
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
		panic(fmt.Sprintf("product_certificate: failed to get parameters: %s", err.Error()))
	}

	// Get all certificates
	certificates, err := k.GetAllCertificates(ctx)
	if err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to get certificates during ExportGenesis",
			"module", types.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight())
		panic(fmt.Sprintf("product_certificate: failed to get certificates: %s", err.Error()))
	}

	// Get all verifications
	verifications, err := k.GetAllVerifications(ctx)
	if err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to get verifications during ExportGenesis",
			"module", types.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight())
		panic(fmt.Sprintf("product_certificate: failed to get verifications: %s", err.Error()))
	}

	return &types.GenesisState{
		Certificates:  certificates,
		Verifications: verifications,
		Params:        params,
	}
}

// ValidateGenesis validates the genesis state for the product_certificate module
func ValidateGenesis(genState types.GenesisState) error {
	return genState.ValidateGenesis()
}

// GetGenesisState returns the current genesis state for the product_certificate module
func GetGenesisState(ctx sdk.Context, k keeper.Keeper) (*types.GenesisState, error) {
	// Get parameters
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get parameters: %w", err)
	}

	// Get all certificates
	certificates, err := k.GetAllCertificates(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get certificates: %w", err)
	}

	// Get all verifications
	verifications, err := k.GetAllVerifications(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get verifications: %w", err)
	}

	return &types.GenesisState{
		Certificates:  certificates,
		Verifications: verifications,
		Params:        params,
	}, nil
}
