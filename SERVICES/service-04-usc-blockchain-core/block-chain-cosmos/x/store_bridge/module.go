package store_bridge

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/spf13/cobra"

	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/store_bridge/keeper"
	storebridgetypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/store_bridge/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// ----------------------------------------------------------------------------
// AppModuleBasic
// ----------------------------------------------------------------------------

// AppModuleBasic implements the AppModuleBasic interface for the bridge module
type AppModuleBasic struct {
	cdc codec.Codec
}

// NewAppModuleBasic creates a new AppModuleBasic instance
func NewAppModuleBasic(cdc codec.Codec) AppModuleBasic {
	return AppModuleBasic{cdc: cdc}
}

// Name returns the bridge module's name
func (AppModuleBasic) Name() string {
	return storebridgetypes.ModuleName
}

// RegisterLegacyAminoCodec registers the bridge module's types with the legacy amino codec
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	storebridgetypes.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the bridge module's interfaces
func (AppModuleBasic) RegisterInterfaces(registry types.InterfaceRegistry) {
	storebridgetypes.RegisterInterfaces(registry)
}

// DefaultGenesis returns the bridge module's default genesis state
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesisState := storebridgetypes.GenesisState{
		Bridges:    []storebridgetypes.Bridge{},
		Transfers:  []storebridgetypes.Transfer{},
		Validators: []storebridgetypes.Validator{},
		Configs:    []storebridgetypes.BridgeConfig{},
		Fees:       []storebridgetypes.BridgeFee{},
		Limits:     []storebridgetypes.BridgeLimit{},
		Events:     []storebridgetypes.BridgeEvent{},
		Params:     storebridgetypes.DefaultParams(),
	}

	genesisBytes, err := json.Marshal(genesisState)
	if err != nil {
		// Note: ctx not available in DefaultGenesis, but this is a critical failure
		panic(fmt.Sprintf("store_bridge: failed to marshal default genesis state: %s", err.Error()))
	}

	return genesisBytes
}

// ValidateGenesis validates the bridge module's genesis state
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var genState storebridgetypes.GenesisState
	if err := json.Unmarshal(bz, &genState); err != nil {
		return err
	}
	return storebridgetypes.ValidateGenesis(genState)
}

// RegisterRESTRoutes - Service-04 uses gRPC API only, REST API is not supported
func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr interface{}) {
	// REST API is not supported in Service-04 (USC Blockchain Core)
	// All API access must use gRPC protocol only
}

// RegisterGRPCGatewayRoutes registers the store_bridge module's gRPC Gateway service handlers
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	// Register gRPC Gateway routes using RegisterQueryHandlerClient
	// This connects the gRPC Gateway to the store_bridge module's query service
	if err := storebridgetypes.RegisterQueryHandlerClient(context.Background(), mux, clientCtx); err != nil {
		// Log error but don't fail - gRPC Gateway routes are optional
		// The actual gRPC service registration happens via module.Configurator
		_ = err
	}
}

// GetTxCmd returns the bridge module's root tx command
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   storebridgetypes.ModuleName,
		Short: "Bridge module subcommands",
		Long:  "Bridge module subcommands for managing cross-chain bridge operations",
	}

	// Add subcommands
	cmd.AddCommand(
		NewCreateBridgeCmd(),
		NewUpdateBridgeCmd(),
		NewInitiateTransferCmd(),
		NewCompleteTransferCmd(),
		NewAddValidatorCmd(),
		NewRemoveValidatorCmd(),
	)

	return cmd
}

// GetQueryCmd returns the bridge module's root query command
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   storebridgetypes.ModuleName,
		Short: "Querying commands for the bridge module",
		Long:  "Querying commands for the bridge module",
	}

	// Add subcommands
	cmd.AddCommand(
		NewQueryBridgeCmd(),
		NewQueryAllBridgesCmd(),
		NewQueryTransferCmd(),
		NewQueryAllTransfersCmd(),
		NewQueryValidatorCmd(),
		NewQueryAllValidatorsCmd(),
		NewQueryBridgeConfigCmd(),
		NewQueryAllBridgeConfigsCmd(),
		NewQueryBridgeFeeCmd(),
		NewQueryAllBridgeFeesCmd(),
		NewQueryBridgeLimitCmd(),
		NewQueryAllBridgeLimitsCmd(),
		NewQueryBridgeEventCmd(),
		NewQueryAllBridgeEventsCmd(),
		NewQueryParamsCmd(),
	)

	return cmd
}

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

// AppModule implements the AppModule interface for the bridge module
type AppModule struct {
	AppModuleBasic
	keeper keeper.Keeper
}

// NewAppModule creates a new AppModule instance
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(cdc),
		keeper:         keeper,
	}
}

// Name returns the bridge module's name
func (am AppModule) Name() string {
	return am.AppModuleBasic.Name()
}

// RegisterServices registers the bridge module's services
func (am AppModule) RegisterServices(cfg module.Configurator) {
	storebridgetypes.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	storebridgetypes.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServerImpl(am.keeper))
}

// RegisterLegacyAminoCodec registers the bridge module's types with the legacy amino codec
func (am AppModule) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	am.AppModuleBasic.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the bridge module's interfaces
func (am AppModule) RegisterInterfaces(registry types.InterfaceRegistry) {
	am.AppModuleBasic.RegisterInterfaces(registry)
}

// InitGenesis performs genesis initialization for the bridge module
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) []abci.ValidatorUpdate {
	var genState storebridgetypes.GenesisState
	if err := json.Unmarshal(gs, &genState); err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to unmarshal genesis state during InitGenesis",
			"module", storebridgetypes.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight(),
			"chain_id", ctx.ChainID())
		panic(fmt.Sprintf("store_bridge: failed to unmarshal genesis state: %s", err.Error()))
	}

	// Initialize genesis state - validate but handle errors gracefully
	if err := storebridgetypes.ValidateGenesis(genState); err != nil {
		// Log validation error but don't panic - return default state
		ctx.Logger().Warn("Genesis validation failed, using default params",
			"module", storebridgetypes.ModuleName,
			"error", err.Error())
		// Use default params if validation fails
		genState.Params = storebridgetypes.DefaultParams()
	}

	// Initialize genesis state via keeper
	// NOTE: During InitGenesis, store service may not be fully available in context.
	// If InitGenesis fails, module will use default parameters (this is expected behavior).
	storebridgetypes.InitGenesis(ctx, am.keeper, genState)

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the bridge module's exported genesis state
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genState := storebridgetypes.ExportGenesis(ctx, am.keeper)
	genesisBytes, err := json.Marshal(genState)
	if err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to marshal genesis state during ExportGenesis",
			"module", storebridgetypes.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight())
		panic(fmt.Sprintf("store_bridge: failed to marshal genesis state: %s", err.Error()))
	}
	return genesisBytes
}

// BeginBlock performs begin block logic for the bridge module
func (am AppModule) BeginBlock(ctx sdk.Context) {
	BeginBlocker(ctx, am.keeper)
}

// EndBlock performs end block logic for the bridge module
func (am AppModule) EndBlock(ctx sdk.Context) []abci.ValidatorUpdate {
	return EndBlocker(ctx, am.keeper)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes
func (am AppModule) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	am.AppModuleBasic.RegisterGRPCGatewayRoutes(clientCtx, mux)
}

// IsAppModule implements module.AppModule
func (am AppModule) IsAppModule() {}

// IsOnePerModuleType implements module.AppModule
func (am AppModule) IsOnePerModuleType() {}

// REST API handlers

// CLI command functions
func NewCreateBridgeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-bridge [id] [name] [from-chain] [to-chain] [type]",
		Short: "Create a new bridge",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement create bridge logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewUpdateBridgeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update-bridge [id] [name] [status]",
		Short: "Update an existing bridge",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement update bridge logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewInitiateTransferCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "initiate-transfer [id] [bridge-id] [from-chain] [to-chain] [from-address] [to-address] [amount] [token]",
		Short: "Initiate a cross-chain transfer",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement initiate transfer logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewCompleteTransferCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "complete-transfer [transfer-id]",
		Short: "Complete a cross-chain transfer",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement complete transfer logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewAddValidatorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add-validator [id] [address] [name] [stake]",
		Short: "Add a new bridge validator",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement add validator logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewRemoveValidatorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove-validator [validator-id]",
		Short: "Remove a bridge validator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement remove validator logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryBridgeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bridge [id]",
		Short: "Query bridge by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query bridge logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllBridgesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-bridges",
		Short: "Query all bridges",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all bridges logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryTransferCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "transfer [id]",
		Short: "Query transfer by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query transfer logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllTransfersCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-transfers",
		Short: "Query all transfers",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all transfers logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryValidatorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validator [id]",
		Short: "Query validator by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query validator logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllValidatorsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-validators",
		Short: "Query all validators",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all validators logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryBridgeConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bridge-config [id]",
		Short: "Query bridge config by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query bridge config logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllBridgeConfigsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-bridge-configs",
		Short: "Query all bridge configs",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all bridge configs logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryBridgeFeeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bridge-fee [id]",
		Short: "Query bridge fee by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query bridge fee logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllBridgeFeesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-bridge-fees",
		Short: "Query all bridge fees",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all bridge fees logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryBridgeLimitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bridge-limit [id]",
		Short: "Query bridge limit by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query bridge limit logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllBridgeLimitsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-bridge-limits",
		Short: "Query all bridge limits",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all bridge limits logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryBridgeEventCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bridge-event [id]",
		Short: "Query bridge event by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query bridge event logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllBridgeEventsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-bridge-events",
		Short: "Query all bridge events",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all bridge events logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryParamsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query bridge module parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query params logic
			return fmt.Errorf("not implemented")
		},
	}

}
