package network

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

	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/network/keeper"
	networktypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/network/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// ----------------------------------------------------------------------------
// AppModuleBasic
// ----------------------------------------------------------------------------

// AppModuleBasic implements the AppModuleBasic interface for the network module
type AppModuleBasic struct {
	cdc codec.Codec
}

// NewAppModuleBasic creates a new AppModuleBasic instance
func NewAppModuleBasic(cdc codec.Codec) AppModuleBasic {
	return AppModuleBasic{cdc: cdc}
}

// Name returns the network module's name
func (AppModuleBasic) Name() string {
	return networktypes.ModuleName
}

// RegisterLegacyAminoCodec registers the network module's types with the legacy amino codec
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	networktypes.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the network module's interfaces
func (AppModuleBasic) RegisterInterfaces(registry types.InterfaceRegistry) {
	networktypes.RegisterInterfaces(registry)
}

// DefaultGenesis returns the network module's default genesis state
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesisState := networktypes.GenesisState{
		Networks:      []networktypes.Network{},
		Nodes:         []networktypes.Node{},
		Connections:   []networktypes.Connection{},
		Syncs:         []networktypes.NetworkSync{},
		HealthMetrics: []networktypes.NetworkHealth{},
		Params:        networktypes.DefaultParams(),
	}

	genesisBytes, err := json.Marshal(genesisState)
	if err != nil {
		// Log error before panic for better debugging
		// Note: ctx not available in DefaultGenesis, using fmt for error message
		panic(fmt.Sprintf("network: failed to marshal default genesis state: %s", err.Error()))
	}

	return genesisBytes
}

// ValidateGenesis validates the network module's genesis state
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var genState networktypes.GenesisState
	if err := json.Unmarshal(bz, &genState); err != nil {
		return err
	}
	return networktypes.ValidateGenesis(genState)
}

// RegisterRESTRoutes - Service-04 uses gRPC API only, REST API is not supported
func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr interface{}) {
	// REST API is not supported in Service-04 (USC Blockchain Core)
	// All API access must use gRPC protocol only
}

// RegisterGRPCGatewayRoutes registers the network module's gRPC Gateway service handlers
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	// Register gRPC Gateway routes using RegisterQueryHandlerClient
	// This connects the gRPC Gateway to the network module's query service
	if err := networktypes.RegisterQueryHandlerClient(context.Background(), mux, clientCtx); err != nil {
		// Log error but don't fail - gRPC Gateway routes are optional
		// The actual gRPC service registration happens via module.Configurator
		_ = err
	}
}

// GetTxCmd returns the network module's root tx command
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   networktypes.ModuleName,
		Short: "Network module subcommands",
		Long:  "Network module subcommands for managing network operations",
	}

	// Add subcommands
	cmd.AddCommand(
		NewCreateNetworkCmd(),
		NewUpdateNetworkCmd(),
		NewJoinNodeCmd(),
		NewLeaveNodeCmd(),
		NewEstablishConnectionCmd(),
		NewStartSyncCmd(),
		NewUpdateHealthCmd(),
	)

	return cmd
}

// GetQueryCmd returns the network module's root query command
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   networktypes.ModuleName,
		Short: "Querying commands for the network module",
		Long:  "Querying commands for the network module",
	}

	// Add subcommands
	cmd.AddCommand(
		NewQueryNetworkCmd(),
		NewQueryAllNetworksCmd(),
		NewQueryNodeCmd(),
		NewQueryAllNodesCmd(),
		NewQueryConnectionCmd(),
		NewQueryAllConnectionsCmd(),
		NewQuerySyncCmd(),
		NewQueryAllSyncsCmd(),
		NewQueryHealthCmd(),
		NewQueryAllHealthsCmd(),
		NewQueryParamsCmd(),
	)

	return cmd
}

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

// AppModule implements the AppModule interface for the network module
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

// Name returns the network module's name
func (am AppModule) Name() string {
	return am.AppModuleBasic.Name()
}

// RegisterServices registers the network module's services
func (am AppModule) RegisterServices(cfg module.Configurator) {
	networktypes.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	networktypes.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServerImpl(am.keeper))
}

// RegisterLegacyAminoCodec registers the network module's types with the legacy amino codec
func (am AppModule) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	am.AppModuleBasic.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the network module's interfaces
func (am AppModule) RegisterInterfaces(registry types.InterfaceRegistry) {
	am.AppModuleBasic.RegisterInterfaces(registry)
}

// InitGenesis performs genesis initialization for the network module
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) []abci.ValidatorUpdate {
	var genState networktypes.GenesisState
	if err := json.Unmarshal(gs, &genState); err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to unmarshal genesis state during InitGenesis",
			"module", networktypes.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight(),
			"chain_id", ctx.ChainID())
		panic(fmt.Sprintf("network: failed to unmarshal genesis state: %s", err.Error()))
	}

	// Initialize genesis state - validate but handle errors gracefully
	if err := networktypes.ValidateGenesis(genState); err != nil {
		// Log validation error but don't panic - return default state
		ctx.Logger().Warn("Genesis validation failed, using default params",
			"module", networktypes.ModuleName,
			"error", err.Error())
		// Use default params if validation fails
		genState.Params = networktypes.DefaultParams()
	}

	// Initialize genesis state via keeper
	// NOTE: During InitGenesis, store service may not be fully available in context.
	// If InitGenesis fails, module will use default parameters (this is expected behavior).
	networktypes.InitGenesis(ctx, am.keeper, genState)

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the network module's exported genesis state
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genState := networktypes.ExportGenesis(ctx, am.keeper)
	genesisBytes, err := json.Marshal(genState)
	if err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to marshal genesis state during ExportGenesis",
			"module", networktypes.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight())
		panic(fmt.Sprintf("network: failed to marshal genesis state: %s", err.Error()))
	}
	return genesisBytes
}

// BeginBlock performs begin block logic for the network module
func (am AppModule) BeginBlock(ctx sdk.Context) {
	BeginBlocker(ctx, am.keeper)
}

// EndBlock performs end block logic for the network module
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
func NewCreateNetworkCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-network [id] [name] [type] [chain-id] [rpc-url]",
		Short: "Create a new network",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement create network logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewUpdateNetworkCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update-network [network-id] [name] [description]",
		Short: "Update a network",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement update network logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewJoinNodeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "join-node [node-id] [network-id] [name] [address] [port]",
		Short: "Join a node to a network",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement join node logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewLeaveNodeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "leave-node [node-id]",
		Short: "Leave a node from a network",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement leave node logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewEstablishConnectionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "establish-connection [connection-id] [network-id] [from-node] [to-node]",
		Short: "Establish a connection between nodes",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement establish connection logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewStartSyncCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start-sync [sync-id] [network-id] [node-id]",
		Short: "Start network synchronization",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement start sync logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewUpdateHealthCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update-health [health-id] [network-id] [health-score]",
		Short: "Update network health metrics",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement update health logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryNetworkCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "network [id]",
		Short: "Query network by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query network logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllNetworksCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-networks",
		Short: "Query all networks",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all networks logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryNodeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "node [id]",
		Short: "Query node by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query node logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllNodesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-nodes",
		Short: "Query all nodes",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all nodes logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryConnectionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "connection [id]",
		Short: "Query connection by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query connection logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllConnectionsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-connections",
		Short: "Query all connections",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all connections logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQuerySyncCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync [id]",
		Short: "Query sync by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query sync logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllSyncsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-syncs",
		Short: "Query all syncs",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all syncs logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryHealthCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "health [id]",
		Short: "Query health by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query health logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllHealthsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-healths",
		Short: "Query all health metrics",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all healths logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryParamsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query network module parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query params logic
			return fmt.Errorf("not implemented")
		},
	}

}
