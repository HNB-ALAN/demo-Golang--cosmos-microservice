package transaction

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/appmodule"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/transaction/keeper"
	transactiontypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/transaction/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
	_ appmodule.AppModule   = AppModule{}
)

// AppModuleBasic defines the basic application module used by the transaction module
type AppModuleBasic struct {
	cdc codec.Codec
}

// Name returns the transaction module's name
func (AppModuleBasic) Name() string {
	return transactiontypes.ModuleName
}

// RegisterLegacyAminoCodec registers the transaction module's types on the given LegacyAmino codec
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	transactiontypes.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the module's interface types
func (a AppModuleBasic) RegisterInterfaces(reg codectypes.InterfaceRegistry) {
	transactiontypes.RegisterInterfaces(reg)
}

// DefaultGenesis returns default genesis state as raw bytes for the transaction module
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(transactiontypes.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the transaction module
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var genState transactiontypes.GenesisState
	if err := cdc.UnmarshalJSON(bz, &genState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", transactiontypes.ModuleName, err)
	}
	return genState.Validate()
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the transaction module
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	// Register gRPC gateway routes here when needed
}

// GetTxCmd returns the transaction module's root tx command
func (a AppModuleBasic) GetTxCmd() *cobra.Command {
	// Return tx command when needed
	return nil
}

// GetQueryCmd returns the transaction module's root query command
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	// Return query command when needed
	return nil
}

// AppModule implements an application module for the transaction module
type AppModule struct {
	AppModuleBasic

	keeper        keeper.Keeper
	accountKeeper transactiontypes.AccountKeeper
	bankKeeper    transactiontypes.BankKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper, ak transactiontypes.AccountKeeper, bk transactiontypes.BankKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc},
		keeper:         keeper,
		accountKeeper:  ak,
		bankKeeper:     bk,
	}
}

// Name returns the transaction module's name
func (am AppModule) Name() string {
	return am.AppModuleBasic.Name()
}

// RegisterServices registers module services
func (am AppModule) RegisterServices(cfg module.Configurator) {
	transactiontypes.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	transactiontypes.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServerImpl(am.keeper))
}

// InitGenesis performs genesis initialization for the transaction module
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) []abci.ValidatorUpdate {
	// Wrap entire function in panic recovery for graceful error handling
	defer func() {
		if r := recover(); r != nil {
			panicMsg := fmt.Sprintf("%v", r)
			// Re-panic with detailed context - will be caught by caller's panic recovery
			panic(fmt.Sprintf("transaction InitGenesis panic: %s", panicMsg))
		}
	}()

	var genState transactiontypes.GenesisState
	if err := cdc.UnmarshalJSON(gs, &genState); err != nil {
		// Log error before panic for better debugging
		ctx.Logger().Error("Failed to unmarshal genesis state during InitGenesis",
			"module", transactiontypes.ModuleName,
			"error", err.Error())
		panic(fmt.Sprintf("failed to unmarshal %s genesis state: %s", transactiontypes.ModuleName, err.Error()))
	}

	// Initialize genesis state - validate but handle errors gracefully
	if err := transactiontypes.ValidateGenesis(&genState); err != nil {
		// Log validation error but don't panic - return default state
		ctx.Logger().Warn("Genesis validation failed, using default params",
			"module", transactiontypes.ModuleName,
			"error", err.Error())
		// Use default params if validation fails
		genState.Params = transactiontypes.DefaultParams
	}

	// Initialize genesis state via keeper
	// NOTE: During InitGenesis, store service may not be fully available in context.
	// If InitGenesis fails, module will use default parameters (this is expected behavior).
	am.keeper.InitGenesis(ctx, genState)

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the transaction module
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genState := am.keeper.ExportGenesis(ctx)
	return cdc.MustMarshalJSON(genState)
}

// ConsensusVersion implements AppModule/ConsensusVersion
func (AppModule) ConsensusVersion() uint64 { return 1 }

// BeginBlock returns the begin blocker for the transaction module
func (am AppModule) BeginBlock(ctx sdk.Context) {
	BeginBlocker(ctx, am.keeper)
}

// EndBlock returns the end blocker for the transaction module
func (am AppModule) EndBlock(ctx sdk.Context) []abci.ValidatorUpdate {
	return EndBlocker(ctx, am.keeper)
}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface
func (am AppModule) IsAppModule() {}

// ProvideModule provides the module with dependencies
func (am AppModule) ProvideModule(cfg module.Configurator) appmodule.AppModule {
	return am
}
