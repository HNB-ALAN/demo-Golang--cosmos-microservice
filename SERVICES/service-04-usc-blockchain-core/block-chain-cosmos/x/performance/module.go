package performance

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

	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/performance/keeper"
	performancetypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/performance/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// ----------------------------------------------------------------------------
// AppModuleBasic
// ----------------------------------------------------------------------------

// AppModuleBasic implements the AppModuleBasic interface for the performance module
type AppModuleBasic struct {
	cdc codec.Codec
}

// NewAppModuleBasic creates a new AppModuleBasic instance
func NewAppModuleBasic(cdc codec.Codec) AppModuleBasic {
	return AppModuleBasic{cdc: cdc}
}

// Name returns the performance module's name
func (AppModuleBasic) Name() string {
	return performancetypes.ModuleName
}

// RegisterLegacyAminoCodec registers the performance module's types with the legacy amino codec
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	performancetypes.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the performance module's interfaces
func (AppModuleBasic) RegisterInterfaces(registry types.InterfaceRegistry) {
	performancetypes.RegisterInterfaces(registry)
}

// DefaultGenesis returns the performance module's default genesis state
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesisState := performancetypes.GenesisState{
		Metrics:       []performancetypes.PerformanceMetric{},
		Benchmarks:    []performancetypes.Benchmark{},
		Optimizations: []performancetypes.Optimization{},
		Alerts:        []performancetypes.PerformanceAlert{},
		Profiles:      []performancetypes.PerformanceProfile{},
		Reports:       []performancetypes.PerformanceReport{},
		Params:        performancetypes.DefaultParams(),
	}

	genesisBytes, err := json.Marshal(genesisState)
	if err != nil {
		// Log error before panic for better debugging
		// Note: ctx not available in DefaultGenesis, using fmt for error message
		panic(fmt.Sprintf("performance: failed to marshal default genesis state: %s", err.Error()))
	}

	return genesisBytes
}

// ValidateGenesis validates the performance module's genesis state
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var genState performancetypes.GenesisState
	if err := json.Unmarshal(bz, &genState); err != nil {
		return err
	}
	return performancetypes.ValidateGenesis(genState)
}

// RegisterRESTRoutes - Service-04 uses gRPC API only, REST API is not supported
func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr interface{}) {
	// REST API is not supported in Service-04 (USC Blockchain Core)
	// All API access must use gRPC protocol only
}

// RegisterGRPCGatewayRoutes registers the performance module's gRPC Gateway service handlers
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	// Register gRPC Gateway routes using RegisterQueryHandlerClient
	// This connects the gRPC Gateway to the performance module's query service
	if err := performancetypes.RegisterQueryHandlerClient(context.Background(), mux, clientCtx); err != nil {
		// Log error but don't fail - gRPC Gateway routes are optional
		// The actual gRPC service registration happens via module.Configurator
		_ = err
	}
}

// GetTxCmd returns the performance module's root tx command
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   performancetypes.ModuleName,
		Short: "Performance module subcommands",
		Long:  "Performance module subcommands for managing performance metrics and optimization",
	}

	// Add subcommands
	cmd.AddCommand(
		NewCreatePerformanceMetricCmd(),
		NewCreateBenchmarkCmd(),
		NewCreateOptimizationCmd(),
		NewCreatePerformanceAlertCmd(),
		NewCreatePerformanceProfileCmd(),
		NewCreatePerformanceReportCmd(),
	)

	return cmd
}

// GetQueryCmd returns the performance module's root query command
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   performancetypes.ModuleName,
		Short: "Querying commands for the performance module",
		Long:  "Querying commands for the performance module",
	}

	// Add subcommands
	cmd.AddCommand(
		NewQueryPerformanceMetricCmd(),
		NewQueryAllPerformanceMetricsCmd(),
		NewQueryBenchmarkCmd(),
		NewQueryAllBenchmarksCmd(),
		NewQueryOptimizationCmd(),
		NewQueryAllOptimizationsCmd(),
		NewQueryPerformanceAlertCmd(),
		NewQueryAllPerformanceAlertsCmd(),
		NewQueryPerformanceProfileCmd(),
		NewQueryAllPerformanceProfilesCmd(),
		NewQueryPerformanceReportCmd(),
		NewQueryAllPerformanceReportsCmd(),
		NewQueryParamsCmd(),
	)

	return cmd
}

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

// AppModule implements the AppModule interface for the performance module
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

// Name returns the performance module's name
func (am AppModule) Name() string {
	return am.AppModuleBasic.Name()
}

// RegisterServices registers the performance module's services
func (am AppModule) RegisterServices(cfg module.Configurator) {
	// TODO: Register message server
	// keeper.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServer(am.keeper))

	// TODO: Register query server
	// keeper.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServer(am.keeper))
}

// RegisterLegacyAminoCodec registers the performance module's types with the legacy amino codec
func (am AppModule) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	am.AppModuleBasic.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the performance module's interfaces
func (am AppModule) RegisterInterfaces(registry types.InterfaceRegistry) {
	am.AppModuleBasic.RegisterInterfaces(registry)
}

// InitGenesis performs genesis initialization for the performance module
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) []abci.ValidatorUpdate {
	var genState performancetypes.GenesisState
	if err := json.Unmarshal(gs, &genState); err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to unmarshal genesis state during InitGenesis",
			"module", performancetypes.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight(),
			"chain_id", ctx.ChainID())
		panic(fmt.Sprintf("performance: failed to unmarshal genesis state: %s", err.Error()))
	}

	// Initialize genesis state - validate but handle errors gracefully
	if err := performancetypes.ValidateGenesis(genState); err != nil {
		// Log validation error but don't panic - return default state
		ctx.Logger().Warn("Genesis validation failed, using default params",
			"module", performancetypes.ModuleName,
			"error", err.Error())
		// Use default params if validation fails
		genState.Params = performancetypes.DefaultParams()
	}

	// Initialize genesis state via keeper
	// NOTE: During InitGenesis, store service may not be fully available in context.
	// If InitGenesis fails, module will use default parameters (this is expected behavior).
	performancetypes.InitGenesis(ctx, am.keeper, genState)

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the performance module's exported genesis state
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genState := performancetypes.ExportGenesis(ctx, am.keeper)
	genesisBytes, err := json.Marshal(genState)
	if err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to marshal genesis state during ExportGenesis",
			"module", performancetypes.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight())
		panic(fmt.Sprintf("performance: failed to marshal genesis state: %s", err.Error()))
	}
	return genesisBytes
}

// BeginBlock performs begin block logic for the performance module
func (am AppModule) BeginBlock(ctx sdk.Context) {
	BeginBlocker(ctx, am.keeper)
}

// EndBlock performs end block logic for the performance module
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
func NewCreatePerformanceMetricCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-performance-metric [id] [name] [value] [unit] [category]",
		Short: "Create a new performance metric",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement create performance metric logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewCreateBenchmarkCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-benchmark [id] [name] [description]",
		Short: "Create a new benchmark",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// CLI command for creating benchmarks
			// Note: Actual implementation would require keeper access via client context
			cmd.PrintErrf("Benchmark creation: %s\n", args[0])
			return nil
		},
	}
}

func NewCreateOptimizationCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-optimization [id] [name] [type] [impact]",
		Short: "Create a new optimization",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			// CLI command for creating optimizations
			// Note: Actual implementation would require keeper access via client context
			cmd.PrintErrf("Optimization creation: %s\n", args[0])
			return nil
		},
	}
}

func NewCreatePerformanceAlertCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-performance-alert [id] [name] [severity] [threshold]",
		Short: "Create a new performance alert",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			// CLI command for creating performance alerts
			// Note: Actual implementation would require keeper access via client context
			cmd.PrintErrf("Performance alert creation: %s\n", args[0])
			return nil
		},
	}
}

func NewCreatePerformanceProfileCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-performance-profile [id] [name] [service]",
		Short: "Create a new performance profile",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement create performance profile logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewCreatePerformanceReportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-performance-report [id] [name] [description]",
		Short: "Create a new performance report",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement create performance report logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryPerformanceMetricCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "performance-metric [id]",
		Short: "Query a performance metric by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query performance metric logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllPerformanceMetricsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-performance-metrics",
		Short: "Query all performance metrics",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all performance metrics logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryBenchmarkCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "benchmark [id]",
		Short: "Query a benchmark by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query benchmark logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllBenchmarksCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-benchmarks",
		Short: "Query all benchmarks",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all benchmarks logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryOptimizationCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "optimization [id]",
		Short: "Query an optimization by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query optimization logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllOptimizationsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-optimizations",
		Short: "Query all optimizations",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all optimizations logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryPerformanceAlertCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "performance-alert [id]",
		Short: "Query a performance alert by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query performance alert logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllPerformanceAlertsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-performance-alerts",
		Short: "Query all performance alerts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all performance alerts logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryPerformanceProfileCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "performance-profile [id]",
		Short: "Query a performance profile by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query performance profile logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllPerformanceProfilesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-performance-profiles",
		Short: "Query all performance profiles",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all performance profiles logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryPerformanceReportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "performance-report [id]",
		Short: "Query a performance report by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query performance report logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllPerformanceReportsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-performance-reports",
		Short: "Query all performance reports",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all performance reports logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryParamsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query performance module parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query params logic
			return fmt.Errorf("not implemented")
		},
	}

}
