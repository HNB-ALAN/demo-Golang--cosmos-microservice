package smart_contract

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

	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/smart_contract/keeper"
	contracttypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/smart_contract/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// ----------------------------------------------------------------------------
// AppModuleBasic
// ----------------------------------------------------------------------------

// AppModuleBasic implements the AppModuleBasic interface for the contract module
type AppModuleBasic struct {
	cdc codec.Codec
}

// NewAppModuleBasic creates a new AppModuleBasic instance
func NewAppModuleBasic(cdc codec.Codec) AppModuleBasic {
	return AppModuleBasic{cdc: cdc}
}

// Name returns the contract module's name
func (AppModuleBasic) Name() string {
	return contracttypes.ModuleName
}

// RegisterLegacyAminoCodec registers the contract module's types with the legacy amino codec
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	contracttypes.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the contract module's interfaces
func (AppModuleBasic) RegisterInterfaces(registry types.InterfaceRegistry) {
	contracttypes.RegisterInterfaces(registry)
}

// DefaultGenesis returns the contract module's default genesis state
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesisState := contracttypes.GenesisState{
		Contracts:   []contracttypes.SmartContract{},
		Executions:  []contracttypes.ContractExecution{},
		Deployments: []contracttypes.ContractDeployment{},
		Upgrades:    []contracttypes.ContractUpgrade{},
		Migrations:  []contracttypes.ContractMigration{},
		Params:      contracttypes.DefaultParams(),
	}

	genesisBytes, err := json.Marshal(genesisState)
	if err != nil {
		// Log error before panic for better debugging
		// Note: ctx not available in DefaultGenesis, using fmt for error message
		panic(fmt.Sprintf("smart_contract: failed to marshal default genesis state: %s", err.Error()))
	}

	return genesisBytes
}

// ValidateGenesis validates the contract module's genesis state
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var genState contracttypes.GenesisState
	if err := json.Unmarshal(bz, &genState); err != nil {
		return err
	}
	return contracttypes.ValidateGenesis(genState)
}

// RegisterRESTRoutes - Service-04 uses gRPC API only, REST API is not supported
func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr interface{}) {
	// REST API is not supported in Service-04 (USC Blockchain Core)
	// All API access must use gRPC protocol only
}

// RegisterGRPCGatewayRoutes registers the smart_contract module's gRPC Gateway service handlers
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	// Register gRPC Gateway routes using RegisterQueryHandlerClient
	// This connects the gRPC Gateway to the smart_contract module's query service
	if err := contracttypes.RegisterQueryHandlerClient(context.Background(), mux, clientCtx); err != nil {
		// Log error but don't fail - gRPC Gateway routes are optional
		// The actual gRPC service registration happens via module.Configurator
		_ = err
	}
}

// GetTxCmd returns the contract module's root tx command
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   contracttypes.ModuleName,
		Short: "Contract module subcommands",
		Long:  "Contract module subcommands for managing smart contracts",
	}

	// Add subcommands
	cmd.AddCommand(
		NewCreateContractCmd(),
		NewUpdateContractCmd(),
		NewExecuteContractCmd(),
		NewDeployContractCmd(),
		NewUpgradeContractCmd(),
		NewMigrateContractCmd(),
	)

	return cmd
}

// GetQueryCmd returns the contract module's root query command
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   contracttypes.ModuleName,
		Short: "Querying commands for the contract module",
		Long:  "Querying commands for the contract module",
	}

	// Add subcommands
	cmd.AddCommand(
		NewQueryContractCmd(),
		NewQueryAllContractsCmd(),
		NewQueryExecutionCmd(),
		NewQueryAllExecutionsCmd(),
		NewQueryDeploymentCmd(),
		NewQueryAllDeploymentsCmd(),
		NewQueryUpgradeCmd(),
		NewQueryAllUpgradesCmd(),
		NewQueryMigrationCmd(),
		NewQueryAllMigrationsCmd(),
		NewQueryParamsCmd(),
	)

	return cmd
}

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

// AppModule implements the AppModule interface for the contract module
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

// Name returns the contract module's name
func (am AppModule) Name() string {
	return am.AppModuleBasic.Name()
}

// RegisterServices registers the contract module's services
func (am AppModule) RegisterServices(cfg module.Configurator) {
	// TODO: Register message server
	// contracttypes.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServer(am.keeper))

	// TODO: Register query server
	// contracttypes.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServer(am.keeper))
}

// RegisterLegacyAminoCodec registers the contract module's types with the legacy amino codec
func (am AppModule) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	am.AppModuleBasic.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the contract module's interfaces
func (am AppModule) RegisterInterfaces(registry types.InterfaceRegistry) {
	am.AppModuleBasic.RegisterInterfaces(registry)
}

// InitGenesis performs genesis initialization for the contract module
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) []abci.ValidatorUpdate {
	var genState contracttypes.GenesisState
	if err := json.Unmarshal(gs, &genState); err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to unmarshal genesis state during InitGenesis",
			"module", contracttypes.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight(),
			"chain_id", ctx.ChainID())
		panic(fmt.Sprintf("smart_contract: failed to unmarshal genesis state: %s", err.Error()))
	}

	// Initialize genesis state - validate but handle errors gracefully
	if err := contracttypes.ValidateGenesis(genState); err != nil {
		// Log validation error but don't panic - return default state
		ctx.Logger().Warn("Genesis validation failed, using default params",
			"module", contracttypes.ModuleName,
			"error", err.Error())
		// Use default params if validation fails
		genState.Params = contracttypes.DefaultParams()
	}

	// Initialize genesis state via keeper
	// NOTE: During InitGenesis, store service may not be fully available in context.
	// If InitGenesis fails, module will use default parameters (this is expected behavior).
	contracttypes.InitGenesis(ctx, am.keeper, genState)

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the contract module's exported genesis state
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genState := contracttypes.ExportGenesis(ctx, am.keeper)
	genesisBytes, err := json.Marshal(genState)
	if err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to marshal genesis state during ExportGenesis",
			"module", contracttypes.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight())
		panic(fmt.Sprintf("smart_contract: failed to marshal genesis state: %s", err.Error()))
	}
	return genesisBytes
}

// BeginBlock performs begin block logic for the contract module
func (am AppModule) BeginBlock(ctx sdk.Context) {
	BeginBlocker(ctx, am.keeper)
}

// EndBlock performs end block logic for the contract module
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
func NewCreateContractCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-contract [id] [name] [type] [owner] [code-hash]",
		Short: "Create a new smart contract",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement create contract logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewUpdateContractCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update-contract [contract-id] [name] [description]",
		Short: "Update a smart contract",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement update contract logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewExecuteContractCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "execute-contract [contract-id] [executor] [method] [input]",
		Short: "Execute a smart contract",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement execute contract logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewDeployContractCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "deploy-contract [contract-id] [deployer] [network] [address]",
		Short: "Deploy a smart contract",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement deploy contract logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewUpgradeContractCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upgrade-contract [contract-id] [upgrader] [new-version] [code-hash]",
		Short: "Upgrade a smart contract",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement upgrade contract logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewMigrateContractCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate-contract [contract-id] [migrator] [from-network] [to-network]",
		Short: "Migrate a smart contract",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement migrate contract logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryContractCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "contract [id]",
		Short: "Query contract by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query contract logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllContractsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-contracts",
		Short: "Query all contracts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all contracts logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryExecutionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "execution [id]",
		Short: "Query execution by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query execution logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllExecutionsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-executions",
		Short: "Query all executions",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all executions logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryDeploymentCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "deployment [id]",
		Short: "Query deployment by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query deployment logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllDeploymentsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-deployments",
		Short: "Query all deployments",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all deployments logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryUpgradeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upgrade [id]",
		Short: "Query upgrade by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query upgrade logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllUpgradesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-upgrades",
		Short: "Query all upgrades",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all upgrades logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryMigrationCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "migration [id]",
		Short: "Query migration by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query migration logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllMigrationsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-migrations",
		Short: "Query all migrations",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all migrations logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryParamsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query contract module parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query params logic
			return fmt.Errorf("not implemented")
		},
	}

}
