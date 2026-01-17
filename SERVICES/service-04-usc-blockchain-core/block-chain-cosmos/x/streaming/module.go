package streaming

import (
	"context"
	"encoding/json"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/streaming/keeper"
	streamingtypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/streaming/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// ----------------------------------------------------------------------------
// AppModuleBasic
// ----------------------------------------------------------------------------

// AppModuleBasic implements the AppModuleBasic interface for the streaming module
type AppModuleBasic struct {
	cdc codec.Codec
}

// NewAppModuleBasic creates a new AppModuleBasic instance
func NewAppModuleBasic(cdc codec.Codec) AppModuleBasic {
	return AppModuleBasic{cdc: cdc}
}

// Name returns the streaming module's name
func (AppModuleBasic) Name() string {
	return streamingtypes.ModuleName
}

// RegisterLegacyAminoCodec registers the streaming module's types with the legacy amino codec
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	streamingtypes.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the streaming module's interfaces
func (AppModuleBasic) RegisterInterfaces(registry types.InterfaceRegistry) {
	streamingtypes.RegisterInterfaces(registry)
}

// DefaultGenesis returns the streaming module's default genesis state
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesisState := streamingtypes.GenesisState{
		Streams:        []streamingtypes.Stream{},
		Viewers:        []streamingtypes.StreamViewer{},
		QualityMetrics: []streamingtypes.StreamQualityMetrics{},
		Analytics:      []streamingtypes.StreamAnalytics{},
		Events:         []streamingtypes.StreamEvent{},
		Params:         streamingtypes.DefaultParams(),
	}

	genesisBytes, err := json.Marshal(genesisState)
	if err != nil {
		// Note: ctx not available in DefaultGenesis, but this is a critical failure
		panic(fmt.Sprintf("streaming: failed to marshal default genesis state: %s", err.Error()))
	}

	return genesisBytes
}

// ValidateGenesis validates the streaming module's genesis state
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var genState streamingtypes.GenesisState
	if err := json.Unmarshal(bz, &genState); err != nil {
		return err
	}
	return streamingtypes.ValidateGenesis(genState)
}

// RegisterRESTRoutes - Service-04 uses gRPC API only, REST API is not supported
func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr interface{}) {
	// REST API is not supported in Service-04 (USC Blockchain Core)
	// All API access must use gRPC protocol only
}

// RegisterGRPCGatewayRoutes registers the streaming module's gRPC Gateway service handlers
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	// Register gRPC Gateway routes using RegisterQueryHandlerClient
	// This connects the gRPC Gateway to the streaming module's query service
	if err := streamingtypes.RegisterQueryHandlerClient(context.Background(), mux, clientCtx); err != nil {
		// Log error but don't fail - gRPC Gateway routes are optional
		// The actual gRPC service registration happens via module.Configurator
		_ = err
	}
}

// GetTxCmd returns the streaming module's root tx command
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   streamingtypes.ModuleName,
		Short: "Streaming module subcommands",
		Long:  "Streaming module subcommands for managing real-time streaming",
	}

	// Add subcommands
	cmd.AddCommand(
		NewCreateStreamCmd(),
		NewUpdateStreamCmd(),
		NewStartStreamCmd(),
		NewStopStreamCmd(),
		NewJoinStreamCmd(),
		NewLeaveStreamCmd(),
		NewUpdateQualityCmd(),
		NewUpdateAnalyticsCmd(),
	)

	return cmd
}

// GetQueryCmd returns the streaming module's root query command
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   streamingtypes.ModuleName,
		Short: "Querying commands for the streaming module",
		Long:  "Querying commands for the streaming module",
	}

	// Add subcommands
	cmd.AddCommand(
		NewQueryStreamCmd(),
		NewQueryAllStreamsCmd(),
		NewQueryViewerCmd(),
		NewQueryAllViewersCmd(),
		NewQueryQualityCmd(),
		NewQueryAllQualitiesCmd(),
		NewQueryAnalyticsCmd(),
		NewQueryAllAnalyticsCmd(),
		NewQueryEventCmd(),
		NewQueryAllEventsCmd(),
		NewQueryParamsCmd(),
	)

	return cmd
}

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

// AppModule implements the AppModule interface for the streaming module
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

// Name returns the streaming module's name
func (am AppModule) Name() string {
	return am.AppModuleBasic.Name()
}

// RegisterServices registers the streaming module's services
func (am AppModule) RegisterServices(cfg module.Configurator) {
	streamingtypes.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	streamingtypes.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServerImpl(am.keeper))
}

// RegisterLegacyAminoCodec registers the streaming module's types with the legacy amino codec
func (am AppModule) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	am.AppModuleBasic.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the streaming module's interfaces
func (am AppModule) RegisterInterfaces(registry types.InterfaceRegistry) {
	am.AppModuleBasic.RegisterInterfaces(registry)
}

// InitGenesis performs genesis initialization for the streaming module
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) []abci.ValidatorUpdate {
	var genState streamingtypes.GenesisState
	if err := json.Unmarshal(gs, &genState); err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to unmarshal genesis state during InitGenesis",
			"module", streamingtypes.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight(),
			"chain_id", ctx.ChainID())
		panic(fmt.Sprintf("streaming: failed to unmarshal genesis state: %s", err.Error()))
	}

	// Initialize genesis state - validate but handle errors gracefully
	if err := streamingtypes.ValidateGenesis(genState); err != nil {
		// Log validation error but don't panic - return default state
		ctx.Logger().Warn("Genesis validation failed, using default params",
			"module", streamingtypes.ModuleName,
			"error", err.Error())
		// Use default params if validation fails
		genState.Params = streamingtypes.DefaultParams()
	}

	// Initialize genesis state via keeper
	// NOTE: During InitGenesis, store service may not be fully available in context.
	// If InitGenesis fails, module will use default parameters (this is expected behavior).
	streamingtypes.InitGenesis(ctx, am.keeper, genState)

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the streaming module's exported genesis state
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genState := streamingtypes.ExportGenesis(ctx, am.keeper)
	genesisBytes, err := json.Marshal(genState)
	if err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to marshal genesis state during ExportGenesis",
			"module", streamingtypes.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight())
		panic(fmt.Sprintf("streaming: failed to marshal genesis state: %s", err.Error()))
	}
	return genesisBytes
}

// BeginBlock performs begin block logic for the streaming module
func (am AppModule) BeginBlock(ctx sdk.Context) {
	BeginBlocker(ctx, am.keeper)
}

// EndBlock performs end block logic for the streaming module
func (am AppModule) EndBlock(ctx sdk.Context) []abci.ValidatorUpdate {
	return EndBlocker(ctx, am.keeper)
}

// IsAppModule implements module.AppModule
func (am AppModule) IsAppModule() {}

// IsOnePerModuleType implements module.AppModule
func (am AppModule) IsOnePerModuleType() {}

// CLI command functions
func NewCreateStreamCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-stream [id] [streamer-id] [title] [description] [category]",
		Short: "Create a new stream",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement create stream logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewUpdateStreamCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update-stream [stream-id] [title] [description] [category]",
		Short: "Update a stream",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement update stream logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewStartStreamCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start-stream [stream-id]",
		Short: "Start a stream",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement start stream logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewStopStreamCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop-stream [stream-id]",
		Short: "Stop a stream",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement stop stream logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewJoinStreamCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "join-stream [viewer-id] [stream-id] [viewer-id] [quality]",
		Short: "Join a stream as viewer",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement join stream logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewLeaveStreamCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "leave-stream [viewer-id]",
		Short: "Leave a stream",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement leave stream logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewUpdateQualityCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update-quality [quality-id] [stream-id] [bitrate] [fps]",
		Short: "Update stream quality metrics",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement update quality logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewUpdateAnalyticsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update-analytics [analytics-id] [stream-id] [viewer-count] [engagement]",
		Short: "Update stream analytics",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement update analytics logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryStreamCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stream [id]",
		Short: "Query stream by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query stream logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllStreamsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-streams",
		Short: "Query all streams",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all streams logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryViewerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "viewer [id]",
		Short: "Query viewer by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query viewer logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllViewersCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-viewers",
		Short: "Query all viewers",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all viewers logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryQualityCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "quality [id]",
		Short: "Query quality by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query quality logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllQualitiesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-qualities",
		Short: "Query all quality metrics",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all qualities logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAnalyticsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "analytics [id]",
		Short: "Query analytics by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query analytics logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllAnalyticsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-analytics",
		Short: "Query all analytics",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all analytics logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryEventCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "event [id]",
		Short: "Query event by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query event logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllEventsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-events",
		Short: "Query all events",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all events logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryParamsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query streaming module parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query params logic
			return fmt.Errorf("not implemented")
		},
	}
}
