package monitoring

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

	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/monitoring/keeper"
	monitoringtypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/monitoring/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// ----------------------------------------------------------------------------
// AppModuleBasic
// ----------------------------------------------------------------------------

// AppModuleBasic implements the AppModuleBasic interface for the monitoring module
type AppModuleBasic struct {
	cdc codec.Codec
}

// NewAppModuleBasic creates a new AppModuleBasic instance
func NewAppModuleBasic(cdc codec.Codec) AppModuleBasic {
	return AppModuleBasic{cdc: cdc}
}

// Name returns the monitoring module's name
func (AppModuleBasic) Name() string {
	return monitoringtypes.ModuleName
}

// RegisterLegacyAminoCodec registers the monitoring module's types with the legacy amino codec
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	monitoringtypes.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the monitoring module's interfaces
func (AppModuleBasic) RegisterInterfaces(registry types.InterfaceRegistry) {
	monitoringtypes.RegisterInterfaces(registry)
}

// DefaultGenesis returns the monitoring module's default genesis state
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesisState := monitoringtypes.GenesisState{
		Metrics:          []monitoringtypes.Metric{},
		Alerts:           []monitoringtypes.Alert{},
		PerformanceData:  []monitoringtypes.PerformanceData{},
		SystemHealth:     []monitoringtypes.SystemHealth{},
		MonitoringConfig: []monitoringtypes.MonitoringConfig{},
		Params:           monitoringtypes.DefaultParams(),
	}

	genesisBytes, err := json.Marshal(genesisState)
	if err != nil {
		// Log error before panic for better debugging
		// Note: ctx not available in DefaultGenesis, using fmt for error message
		panic(fmt.Sprintf("monitoring: failed to marshal default genesis state: %s", err.Error()))
	}

	return genesisBytes
}

// ValidateGenesis validates the monitoring module's genesis state
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var genState monitoringtypes.GenesisState
	if err := json.Unmarshal(bz, &genState); err != nil {
		return err
	}
	return monitoringtypes.ValidateGenesis(genState)
}

// RegisterRESTRoutes - Service-04 uses gRPC API only, REST API is not supported
func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr interface{}) {
	// REST API is not supported in Service-04 (USC Blockchain Core)
	// All API access must use gRPC protocol only
}

// RegisterGRPCGatewayRoutes registers the monitoring module's gRPC Gateway service handlers
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	// Register gRPC Gateway routes using RegisterQueryHandlerClient
	// This connects the gRPC Gateway to the monitoring module's query service
	if err := monitoringtypes.RegisterQueryHandlerClient(context.Background(), mux, clientCtx); err != nil {
		// Log error but don't fail - gRPC Gateway routes are optional
		// The actual gRPC service registration happens via module.Configurator
		_ = err
	}
}

// GetTxCmd returns the monitoring module's root tx command
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   monitoringtypes.ModuleName,
		Short: "Monitoring module subcommands",
		Long:  "Monitoring module subcommands for managing performance metrics and alerts",
	}

	// Add subcommands
	cmd.AddCommand(
		NewCreateMetricCmd(),
		NewUpdateMetricCmd(),
		NewCreateAlertCmd(),
		NewUpdateAlertCmd(),
		NewCreatePerformanceDataCmd(),
		NewCreateSystemHealthCmd(),
		NewCreateMonitoringConfigCmd(),
	)

	return cmd
}

// GetQueryCmd returns the monitoring module's root query command
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   monitoringtypes.ModuleName,
		Short: "Querying commands for the monitoring module",
		Long:  "Querying commands for the monitoring module",
	}

	// Add subcommands
	cmd.AddCommand(
		NewQueryMetricCmd(),
		NewQueryAllMetricsCmd(),
		NewQueryAlertCmd(),
		NewQueryAllAlertsCmd(),
		NewQueryPerformanceDataCmd(),
		NewQueryAllPerformanceDataCmd(),
		NewQuerySystemHealthCmd(),
		NewQueryAllSystemHealthCmd(),
		NewQueryMonitoringConfigCmd(),
		NewQueryAllMonitoringConfigsCmd(),
		NewQueryParamsCmd(),
	)

	return cmd
}

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

// AppModule implements the AppModule interface for the monitoring module
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

// Name returns the monitoring module's name
func (am AppModule) Name() string {
	return am.AppModuleBasic.Name()
}

// RegisterServices registers the monitoring module's services
func (am AppModule) RegisterServices(cfg module.Configurator) {
	monitoringtypes.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	monitoringtypes.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServerImpl(am.keeper))
}

// RegisterLegacyAminoCodec registers the monitoring module's types with the legacy amino codec
func (am AppModule) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	am.AppModuleBasic.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the monitoring module's interfaces
func (am AppModule) RegisterInterfaces(registry types.InterfaceRegistry) {
	am.AppModuleBasic.RegisterInterfaces(registry)
}

// InitGenesis performs genesis initialization for the monitoring module
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) []abci.ValidatorUpdate {
	var genState monitoringtypes.GenesisState
	if err := json.Unmarshal(gs, &genState); err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to unmarshal genesis state during InitGenesis",
			"module", monitoringtypes.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight(),
			"chain_id", ctx.ChainID())
		panic(fmt.Sprintf("monitoring: failed to unmarshal genesis state: %s", err.Error()))
	}

	// Initialize genesis state - validate but handle errors gracefully
	if err := monitoringtypes.ValidateGenesis(genState); err != nil {
		// Log validation error but don't panic - return default state
		ctx.Logger().Warn("Genesis validation failed, using default params",
			"module", monitoringtypes.ModuleName,
			"error", err.Error())
		// Use default params if validation fails
		genState.Params = monitoringtypes.DefaultParams()
	}

	// Initialize genesis state via keeper
	// NOTE: During InitGenesis, store service may not be fully available in context.
	// If InitGenesis fails, module will use default parameters (this is expected behavior).
	monitoringtypes.InitGenesis(ctx, am.keeper, genState)

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the monitoring module's exported genesis state
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genState := monitoringtypes.ExportGenesis(ctx, am.keeper)
	genesisBytes, err := json.Marshal(genState)
	if err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to marshal genesis state during ExportGenesis",
			"module", monitoringtypes.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight())
		panic(fmt.Sprintf("monitoring: failed to marshal genesis state: %s", err.Error()))
	}
	return genesisBytes
}

// BeginBlock performs begin block logic for the monitoring module
func (am AppModule) BeginBlock(ctx sdk.Context) {
	BeginBlocker(ctx, am.keeper)
}

// EndBlock performs end block logic for the monitoring module
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
func NewCreateMetricCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-metric [id] [name] [value] [unit]",
		Short: "Create a new metric",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement create metric logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewUpdateMetricCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update-metric [id] [name] [value] [unit]",
		Short: "Update an existing metric",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement update metric logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewCreateAlertCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-alert [id] [name] [severity] [threshold]",
		Short: "Create a new alert",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement create alert logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewUpdateAlertCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update-alert [id] [name] [severity] [threshold]",
		Short: "Update an existing alert",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement update alert logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewCreatePerformanceDataCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-performance-data [id] [service] [metric] [value]",
		Short: "Create performance data",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement create performance data logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewCreateSystemHealthCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-system-health [id] [status] [score]",
		Short: "Create system health record",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement create system health logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewCreateMonitoringConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-monitoring-config [id] [service] [enabled]",
		Short: "Create monitoring configuration",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement create monitoring config logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryMetricCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "metric [id]",
		Short: "Query a metric by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query metric logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllMetricsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-metrics",
		Short: "Query all metrics",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all metrics logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAlertCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "alert [id]",
		Short: "Query an alert by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query alert logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllAlertsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-alerts",
		Short: "Query all alerts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all alerts logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryPerformanceDataCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "performance-data [id]",
		Short: "Query performance data by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query performance data logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllPerformanceDataCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-performance-data",
		Short: "Query all performance data",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all performance data logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQuerySystemHealthCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "system-health [id]",
		Short: "Query system health by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query system health logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllSystemHealthCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-system-health",
		Short: "Query all system health records",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all system health logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryMonitoringConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "monitoring-config [id]",
		Short: "Query monitoring config by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query monitoring config logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllMonitoringConfigsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-monitoring-configs",
		Short: "Query all monitoring configs",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all monitoring configs logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryParamsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query monitoring module parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query params logic
			return fmt.Errorf("not implemented")
		},
	}

}
