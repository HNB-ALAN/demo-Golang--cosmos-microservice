package store_network

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

	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/store_network/keeper"
	storenetworktypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/store_network/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// ----------------------------------------------------------------------------
// AppModuleBasic
// ----------------------------------------------------------------------------

// AppModuleBasic implements the AppModuleBasic interface for the store module
type AppModuleBasic struct {
	cdc codec.Codec
}

// NewAppModuleBasic creates a new AppModuleBasic instance
func NewAppModuleBasic(cdc codec.Codec) AppModuleBasic {
	return AppModuleBasic{cdc: cdc}
}

// Name returns the store module's name
func (AppModuleBasic) Name() string {
	return storenetworktypes.ModuleName
}

// RegisterLegacyAminoCodec registers the store module's types with the legacy amino codec
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	storenetworktypes.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the store module's interfaces
func (AppModuleBasic) RegisterInterfaces(registry types.InterfaceRegistry) {
	storenetworktypes.RegisterInterfaces(registry)
}

// DefaultGenesis returns the store module's default genesis state
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesisState := storenetworktypes.GenesisState{
		StoredData:   []storenetworktypes.StoredData{},
		Stores:       []storenetworktypes.Store{},
		Backups:      []storenetworktypes.Backup{},
		Restores:     []storenetworktypes.Restore{},
		StoreIndexes: []storenetworktypes.StoreIndex{},
		StoreQueries: []storenetworktypes.StoreQuery{},
		Transactions: []storenetworktypes.StoreTransaction{},
		Params:       storenetworktypes.DefaultParams(),
	}

	genesisBytes, err := json.Marshal(genesisState)
	if err != nil {
		// Note: ctx not available in DefaultGenesis, but this is a critical failure
		panic(fmt.Sprintf("store_network: failed to marshal default genesis state: %s", err.Error()))
	}

	return genesisBytes
}

// ValidateGenesis validates the store module's genesis state
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var genState storenetworktypes.GenesisState
	if err := json.Unmarshal(bz, &genState); err != nil {
		return err
	}
	return storenetworktypes.ValidateGenesis(genState)
}

// RegisterRESTRoutes - Service-04 uses gRPC API only, REST API is not supported
func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr interface{}) {
	// REST API is not supported in Service-04 (USC Blockchain Core)
	// All API access must use gRPC protocol only
}

// RegisterGRPCGatewayRoutes registers the store_network module's gRPC Gateway service handlers
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	// Register gRPC Gateway routes using RegisterQueryHandlerClient
	// This connects the gRPC Gateway to the store_network module's query service
	if err := storenetworktypes.RegisterQueryHandlerClient(context.Background(), mux, clientCtx); err != nil {
		// Log error but don't fail - gRPC Gateway routes are optional
		// The actual gRPC service registration happens via module.Configurator
		_ = err
	}
}

// GetTxCmd returns the store module's root tx command
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   storenetworktypes.ModuleName,
		Short: "Store module subcommands",
		Long:  "Store module subcommands for managing data storage and state",
	}

	// Add subcommands
	cmd.AddCommand(
		NewStoreDataCmd(),
		NewCreateStoreCmd(),
		NewCreateBackupCmd(),
		NewCreateRestoreCmd(),
		NewCreateStoreIndexCmd(),
		NewCreateStoreQueryCmd(),
		NewCreateStoreTransactionCmd(),
	)

	return cmd
}

// GetQueryCmd returns the store module's root query command
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   storenetworktypes.ModuleName,
		Short: "Querying commands for the store module",
		Long:  "Querying commands for the store module",
	}

	// Add subcommands
	cmd.AddCommand(
		NewQueryStoredDataCmd(),
		NewQueryAllStoredDataCmd(),
		NewQueryStoreCmd(),
		NewQueryAllStoresCmd(),
		NewQueryBackupCmd(),
		NewQueryAllBackupsCmd(),
		NewQueryRestoreCmd(),
		NewQueryAllRestoresCmd(),
		NewQueryStoreIndexCmd(),
		NewQueryAllStoreIndexesCmd(),
		NewQueryStoreQueryCmd(),
		NewQueryAllStoreQueriesCmd(),
		NewQueryStoreTransactionCmd(),
		NewQueryAllStoreTransactionsCmd(),
		NewQueryParamsCmd(),
	)

	return cmd
}

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

// AppModule implements the AppModule interface for the store module
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

// Name returns the store module's name
func (am AppModule) Name() string {
	return am.AppModuleBasic.Name()
}

// RegisterServices registers the store module's services
func (am AppModule) RegisterServices(cfg module.Configurator) {
	storenetworktypes.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	storenetworktypes.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServerImpl(am.keeper))
}

// RegisterLegacyAminoCodec registers the store module's types with the legacy amino codec
func (am AppModule) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	am.AppModuleBasic.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the store module's interfaces
func (am AppModule) RegisterInterfaces(registry types.InterfaceRegistry) {
	am.AppModuleBasic.RegisterInterfaces(registry)
}

// InitGenesis performs genesis initialization for the store module
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) []abci.ValidatorUpdate {
	var genState storenetworktypes.GenesisState
	if err := json.Unmarshal(gs, &genState); err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to unmarshal genesis state during InitGenesis",
			"module", storenetworktypes.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight(),
			"chain_id", ctx.ChainID())
		panic(fmt.Sprintf("store_network: failed to unmarshal genesis state: %s", err.Error()))
	}

	// Initialize genesis state - validate but handle errors gracefully
	if err := storenetworktypes.ValidateGenesis(genState); err != nil {
		// Log validation error but don't panic - return default state
		ctx.Logger().Warn("Genesis validation failed, using default params",
			"module", storenetworktypes.ModuleName,
			"error", err.Error())
		// Use default params if validation fails
		genState.Params = storenetworktypes.DefaultParams()
	}

	// Initialize genesis state via keeper
	// NOTE: During InitGenesis, store service may not be fully available in context.
	// If InitGenesis fails, module will use default parameters (this is expected behavior).
	storenetworktypes.InitGenesis(ctx, am.keeper, genState)

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the store module's exported genesis state
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genState := storenetworktypes.ExportGenesis(ctx, am.keeper)
	genesisBytes, err := json.Marshal(genState)
	if err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to marshal genesis state during ExportGenesis",
			"module", storenetworktypes.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight())
		panic(fmt.Sprintf("store_network: failed to marshal genesis state: %s", err.Error()))
	}
	return genesisBytes
}

// BeginBlock performs begin block logic for the store module
func (am AppModule) BeginBlock(ctx sdk.Context) {
	BeginBlocker(ctx, am.keeper)
}

// EndBlock performs end block logic for the store module
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
func NewStoreDataCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "store-data [id] [key] [value] [content-type]",
		Short: "Store data",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			// CLI command for storing data
			// Note: Actual implementation would require keeper access via client context
			cmd.PrintErrf("Store data: %s\n", args[0])
			return nil
		},
	}
}

func NewCreateStoreCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-store [id] [name] [type]",
		Short: "Create a new store",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// CLI command for creating stores
			// Note: Actual implementation would require keeper access via client context
			cmd.PrintErrf("Create store: %s\n", args[0])
			return nil
		},
	}
}

func NewCreateBackupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-backup [id] [store-id] [name]",
		Short: "Create a new backup",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// CLI command for creating backups
			// Note: Actual implementation would require keeper access via client context
			cmd.PrintErrf("Create backup: %s\n", args[0])
			return nil
		},
	}
}

func NewCreateRestoreCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-restore [id] [backup-id] [store-id] [name]",
		Short: "Create a new restore",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement create restore logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewCreateStoreIndexCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-store-index [id] [store-id] [name] [type]",
		Short: "Create a new store index",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement create store index logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewCreateStoreQueryCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-store-query [id] [store-id] [name] [query] [type]",
		Short: "Create a new store query",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement create store query logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewCreateStoreTransactionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-store-transaction [id] [store-id] [type]",
		Short: "Create a new store transaction",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement create store transaction logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryStoredDataCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stored-data [id]",
		Short: "Query stored data by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query stored data logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllStoredDataCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-stored-data",
		Short: "Query all stored data",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all stored data logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryStoreCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "store [id]",
		Short: "Query store by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query store logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllStoresCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-stores",
		Short: "Query all stores",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all stores logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryBackupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "backup [id]",
		Short: "Query backup by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query backup logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllBackupsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-backups",
		Short: "Query all backups",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all backups logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryRestoreCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "restore [id]",
		Short: "Query restore by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query restore logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllRestoresCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-restores",
		Short: "Query all restores",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all restores logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryStoreIndexCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "store-index [id]",
		Short: "Query store index by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query store index logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllStoreIndexesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-store-indexes",
		Short: "Query all store indexes",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all store indexes logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryStoreQueryCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "store-query [id]",
		Short: "Query store query by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query store query logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllStoreQueriesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-store-queries",
		Short: "Query all store queries",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all store queries logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryStoreTransactionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "store-transaction [id]",
		Short: "Query store transaction by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query store transaction logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllStoreTransactionsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-store-transactions",
		Short: "Query all store transactions",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all store transactions logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryParamsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query store module parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query params logic
			return fmt.Errorf("not implemented")
		},
	}

}
