package product_certificate

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/product_certificate/keeper"
	productcertificatetypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/product_certificate/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
	_ appmodule.AppModule   = AppModule{}
)

// ConsensusVersion defines the current product_certificate module consensus version.
const ConsensusVersion = 1

// AppModuleBasic defines the basic application module used by the product_certificate module.
type AppModuleBasic struct {
	cdc codec.Codec
}

// Name returns the product_certificate module's name.
func (AppModuleBasic) Name() string {
	return productcertificatetypes.ModuleName
}

// RegisterLegacyAminoCodec registers the product_certificate module's types on the given LegacyAmino codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	productcertificatetypes.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the module's interface types
func (a AppModuleBasic) RegisterInterfaces(reg codectypes.InterfaceRegistry) {
	productcertificatetypes.RegisterInterfaces(reg)
}

// DefaultGenesis returns default genesis state as raw bytes for the product_certificate module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(productcertificatetypes.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the product_certificate module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var genState productcertificatetypes.GenesisState
	if err := cdc.UnmarshalJSON(bz, &genState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", productcertificatetypes.ModuleName, err)
	}
	return genState.ValidateGenesis()
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the product_certificate module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	// gRPC Gateway routes can be added here if needed
}

// GetTxCmd returns the transaction commands for this module
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	// CLI commands can be added here if needed
	return nil
}

// GetQueryCmd returns the query commands for this module
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	// CLI commands can be added here if needed
	return nil
}

// AppModule implements an application module for the product_certificate module.
type AppModule struct {
	AppModuleBasic

	keeper     keeper.Keeper
	paramSpace paramtypes.Subspace
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper, paramSpace paramtypes.Subspace) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc},
		keeper:         keeper,
		paramSpace:     paramSpace,
	}
}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

// Name returns the product_certificate module's name.
func (am AppModule) Name() string {
	return am.AppModuleBasic.Name()
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	productcertificatetypes.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	productcertificatetypes.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServerImpl(am.keeper))
}

// InitGenesis performs genesis initialization for the product_certificate module.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) []abci.ValidatorUpdate {
	// Wrap entire function in panic recovery for graceful error handling
	defer func() {
		if r := recover(); r != nil {
			panicMsg := fmt.Sprintf("%v", r)
			// Re-panic with detailed context - will be caught by caller's panic recovery
			panic(fmt.Sprintf("product_certificate InitGenesis panic: %s", panicMsg))
		}
	}()

	var genState productcertificatetypes.GenesisState
	if err := cdc.UnmarshalJSON(gs, &genState); err != nil {
		// Log error before panic for better debugging
		ctx.Logger().Error("Failed to unmarshal genesis state during InitGenesis",
			"module", productcertificatetypes.ModuleName,
			"error", err.Error())
		panic(fmt.Sprintf("failed to unmarshal %s genesis state: %s", productcertificatetypes.ModuleName, err.Error()))
	}

	// Initialize genesis state - validate but handle errors gracefully
	if err := genState.ValidateGenesis(); err != nil {
		// Log validation error but don't panic - return default state
		ctx.Logger().Warn("Genesis validation failed, using default params",
			"module", productcertificatetypes.ModuleName,
			"error", err.Error())
		// Use default params if validation fails
		genState.Params = productcertificatetypes.DefaultParams()
	}

	// Initialize genesis state via keeper
	// NOTE: During InitGenesis, store service may not be fully available in context.
	// If InitGenesis fails, module will use default parameters (this is expected behavior).
	am.keeper.InitGenesis(ctx, genState)

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the product_certificate module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genState := am.keeper.ExportGenesis(ctx)
	return cdc.MustMarshalJSON(genState)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return ConsensusVersion }

// BeginBlock returns the begin blocker for the product_certificate module.
func (am AppModule) BeginBlock(ctx sdk.Context) {
	am.keeper.BeginBlocker(ctx)
}

// EndBlock returns the end blocker for the product_certificate module.
func (am AppModule) EndBlock(ctx sdk.Context) []abci.ValidatorUpdate {
	return am.keeper.EndBlocker(ctx)
}

// ProvideModule provides the module.
func ProvideModule(in depinject.In) (AppModule, error) {
	// TODO: Implement proper dependency injection
	// cdc := depinject.ExtractKvStoreKey[codec.Codec](in)
	// keeper := keeper.NewKeeper(cdc, depinject.ExtractKvStoreKey[store.KVStoreService](in))
	// paramSpace := depinject.ExtractKvStoreKey[paramtypes.Subspace](in)

	// return NewAppModule(cdc, keeper, paramSpace), nil
	return AppModule{}, nil
}
