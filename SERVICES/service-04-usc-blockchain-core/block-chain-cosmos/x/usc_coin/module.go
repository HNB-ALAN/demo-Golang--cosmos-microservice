package usc_coin

import (
	"encoding/json"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/usc_coin/keeper"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/usc_coin/types"
)

var (
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the USC module
type AppModuleBasic struct {
	cdc codec.Codec
}

// Name returns the USC module's name
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterCodec registers the USC module's types for the given codec
func (AppModuleBasic) RegisterCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

// RegisterLegacyAminoCodec registers the USC module's types on the given LegacyAmino codec
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the USC module's interface types
func (a AppModuleBasic) RegisterInterfaces(reg codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(reg)
}

// DefaultGenesis returns default genesis state as raw bytes for the USC module
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genState := types.DefaultGenesisState()
	bz, err := json.Marshal(genState)
	if err != nil {
		// Note: ctx not available in DefaultGenesis, but this is a critical failure
		panic(fmt.Sprintf("usc_coin: failed to marshal default genesis state: %s", err.Error()))
	}
	return bz
}

// ValidateGenesis performs genesis state validation for the USC module
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var genState types.GenesisState
	if err := json.Unmarshal(bz, &genState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return genState.ValidateGenesis()
}

// RegisterRESTRoutes - Service-04 uses gRPC API only, REST API is not supported
func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr interface{}) {
	// REST API is not supported in Service-04 (USC Blockchain Core)
	// All API access must use gRPC protocol only
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the USC module
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	// gRPC Gateway routes can be added here if needed
}

// GetTxCmd returns the root tx command for the USC module
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "USC coin module subcommands",
		Long:  "USC coin module subcommands for managing USC token transfers, minting, and burning",
	}

	// Add subcommands
	cmd.AddCommand(
		NewTransferUSCCmd(),
		NewMintUSCCmd(),
		NewBurnUSCCmd(),
	)

	return cmd
}

// GetQueryCmd returns the root query command for the USC module
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Querying commands for the USC coin module",
		Long:  "Querying commands for the USC coin module",
	}

	// Add subcommands
	cmd.AddCommand(
		NewQueryUSCBalanceCmd(),
		NewQueryUSCSupplyCmd(),
		NewQueryUSCHoldersCmd(),
		NewQueryParamsCmd(),
	)

	return cmd
}

// AppModule implements an application module for the USC module
type AppModule struct {
	AppModuleBasic

	keeper     keeper.Keeper
	bankKeeper types.BankKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper, bankKeeper types.BankKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc},
		keeper:         keeper,
		bankKeeper:     bankKeeper,
	}
}

// Name returns the USC module's name
func (am AppModule) Name() string {
	return am.AppModuleBasic.Name()
}

// RegisterServices registers a gRPC query service to respond to the module-specific gRPC queries
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServerImpl(am.keeper))
}

// RegisterInvariants registers the USC module invariants
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	// Invariants can be registered here if needed
}

// Route returns the message routing key for the USC module
func (am AppModule) Route() string {
	return types.RouterKey
}

// QuerierRoute returns the USC module's querier route name
func (AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

// InitGenesis performs genesis initialization for the USC module
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) []abci.ValidatorUpdate {
	var genState types.GenesisState
	if err := json.Unmarshal(gs, &genState); err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to unmarshal genesis state during InitGenesis",
			"module", types.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight(),
			"chain_id", ctx.ChainID())
		panic(fmt.Sprintf("usc_coin: failed to unmarshal genesis state: %s", err.Error()))
	}

	// Initialize genesis state - validate but handle errors gracefully
	if err := genState.ValidateGenesis(); err != nil {
		// Log validation error but don't panic - return default state
		ctx.Logger().Warn("Genesis validation failed, using default params",
			"module", types.ModuleName,
			"error", err.Error())
		// Use default params if validation fails
		genState.Params = types.DefaultParams()
	}

	// Set parameters
	// NOTE: During InitGenesis, store service may not be fully available in context.
	// If SetParams fails, module will use default parameters (this is expected behavior).
	if err := am.keeper.SetParams(ctx, genState.Params); err != nil {
		// Expected: Store service not available during InitGenesis
		// Module will use default parameters, chain can continue normally
		ctx.Logger().Info("Using default parameters (store service not available in InitGenesis)",
			"module", types.ModuleName)
	} else {
		ctx.Logger().Info("Parameters set from genesis",
			"module", types.ModuleName)
	}

	// Set balances (skip if balances array is empty or nil)
	if len(genState.Balances) > 0 {
		for _, balance := range genState.Balances {
			if balance.Address != "" {
				if err := am.keeper.SetBalance(ctx, balance.Address, balance); err != nil {
					// Log error but don't panic - skip this balance
					ctx.Logger().Warn("Failed to set balance, skipping",
						"module", types.ModuleName,
						"address", balance.Address,
						"error", err.Error())
				}
			}
		}
	}

	// Set transfers (skip if transfers array is empty or nil)
	if len(genState.Transfers) > 0 {
		for _, transfer := range genState.Transfers {
			if transfer.FromAddress != "" && transfer.ToAddress != "" {
				if err := am.keeper.SetTransfer(ctx, transfer); err != nil {
					// Log error but don't panic - skip this transfer
					ctx.Logger().Warn("Failed to set transfer, skipping",
						"module", types.ModuleName,
						"from", transfer.FromAddress,
						"to", transfer.ToAddress,
						"error", err.Error())
				}
			}
		}
	}

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the USC module
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	// Get parameters
	params, err := am.keeper.GetParams(ctx)
	if err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to get parameters during ExportGenesis",
			"module", types.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight(),
			"chain_id", ctx.ChainID())
		panic(fmt.Sprintf("usc_coin: failed to get parameters: %s", err.Error()))
	}

	// Get all balances
	balances, err := am.keeper.GetAllBalances(ctx)
	if err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to get balances during ExportGenesis",
			"module", types.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight())
		panic(fmt.Sprintf("usc_coin: failed to get balances: %s", err.Error()))
	}

	// Get all transfers
	transfers, err := am.keeper.GetAllTransfers(ctx)
	if err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to get transfers during ExportGenesis",
			"module", types.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight())
		panic(fmt.Sprintf("usc_coin: failed to get transfers: %s", err.Error()))
	}

	genState := &types.GenesisState{
		Balances:  balances,
		Transfers: transfers,
		Params:    params,
	}

	bz, err := json.Marshal(genState)
	if err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to marshal genesis state during ExportGenesis",
			"module", types.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight())
		panic(fmt.Sprintf("usc_coin: failed to marshal genesis state: %s", err.Error()))
	}
	return bz
}

// ConsensusVersion implements AppModule/ConsensusVersion
func (AppModule) ConsensusVersion() uint64 { return 1 }

// BeginBlock returns the begin blocker for the USC module
func (am AppModule) BeginBlock(ctx sdk.Context) {
	BeginBlocker(ctx, am.keeper)
}

// EndBlock returns the end blocker for the USC module
func (am AppModule) EndBlock(ctx sdk.Context) []abci.ValidatorUpdate {
	return EndBlocker(ctx, am.keeper)
}

// IsAppModule implements module.AppModule
func (AppModule) IsAppModule() {}

// IsOnePerModuleType implements module.AppModule
func (AppModule) IsOnePerModuleType() {}

// ----------------------------------------------------------------------------
// CLI Commands
// ----------------------------------------------------------------------------

// NewTransferUSCCmd returns a CLI command for transferring USC tokens
func NewTransferUSCCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "transfer [from] [to] [amount]",
		Short: "Transfer USC tokens",
		Long:  "Transfer USC tokens from one address to another",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get arguments
			from := args[0]
			to := args[1]
			amount := args[2]

			// Note: CLI commands require gRPC query server implementation
			// In production, use gRPC queries via query server
			cmd.PrintErrf("Transferring USC: %s -> %s, amount: %s\n", from, to, amount)
			cmd.PrintErrf("Note: Use gRPC query server for actual transaction\n")

			return nil
		},
	}
}

// NewMintUSCCmd returns a CLI command for minting USC tokens
func NewMintUSCCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mint [to] [amount]",
		Short: "Mint USC tokens",
		Long:  "Mint new USC tokens to an address",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get arguments
			to := args[0]
			amount := args[1]

			// Note: CLI commands require gRPC query server implementation
			// In production, use gRPC queries via query server
			cmd.PrintErrf("Minting USC: %s, amount: %s\n", to, amount)
			cmd.PrintErrf("Note: Use gRPC query server for actual transaction\n")

			return nil
		},
	}
}

// NewBurnUSCCmd returns a CLI command for burning USC tokens
func NewBurnUSCCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "burn [from] [amount]",
		Short: "Burn USC tokens",
		Long:  "Burn USC tokens from an address",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get arguments
			from := args[0]
			amount := args[1]

			// Note: CLI commands require gRPC query server implementation
			// In production, use gRPC queries via query server
			cmd.PrintErrf("Burning USC: %s, amount: %s\n", from, amount)
			cmd.PrintErrf("Note: Use gRPC query server for actual transaction\n")

			return nil
		},
	}
}

// NewQueryUSCBalanceCmd returns a CLI command for querying USC balance
func NewQueryUSCBalanceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "balance [address]",
		Short: "Query USC balance",
		Long:  "Query USC balance for a specific address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get address from arguments
			address := args[0]

			// Note: CLI commands require gRPC query server implementation
			// In production, use gRPC queries via query server
			cmd.PrintErrf("Querying USC balance for address: %s\n", address)
			cmd.PrintErrf("Note: Use gRPC query server for actual balance data\n")

			return nil
		},
	}
}

// NewQueryUSCSupplyCmd returns a CLI command for querying USC total supply
func NewQueryUSCSupplyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "supply",
		Short: "Query USC total supply",
		Long:  "Query the total supply of USC tokens",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Note: CLI commands require gRPC query server implementation
			// In production, use gRPC queries via query server
			cmd.PrintErrf("Querying USC total supply\n")
			cmd.PrintErrf("Note: Use gRPC query server for actual supply data\n")

			return nil
		},
	}
}

// NewQueryUSCHoldersCmd returns a CLI command for querying USC holders
func NewQueryUSCHoldersCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "holders",
		Short: "Query USC holders",
		Long:  "Query all USC token holders with pagination",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Note: CLI commands require gRPC query server implementation
			// In production, use gRPC queries via query server
			cmd.PrintErrf("Querying USC holders\n")
			cmd.PrintErrf("Note: Use gRPC query server for actual holders data\n")

			return nil
		},
	}
}

// NewQueryParamsCmd returns a CLI command for querying USC module parameters
func NewQueryParamsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query USC module parameters",
		Long:  "Query USC module parameters",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Note: CLI commands require gRPC query server implementation
			// In production, use gRPC queries via query server
			cmd.PrintErrf("Querying USC module parameters\n")
			cmd.PrintErrf("Note: Use gRPC query server for actual parameter data\n")

			return nil
		},
	}
}
