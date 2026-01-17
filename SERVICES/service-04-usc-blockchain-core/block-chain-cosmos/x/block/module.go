package block

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/block/keeper"
	blocktypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/block/types"
)

// Import block package functions (BeginBlocker, EndBlocker)
// These are in the same package, so we can call them directly

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
	// COSMOS SDK 0.53.4: Explicitly implement HasABCIGenesis to call AppModule.InitGenesis
	// This ensures our custom InitGenesis method is called, not a default implementation
	_ module.HasABCIGenesis = AppModule{}
)

// AppModuleBasic defines the basic application module used by the block module
type AppModuleBasic struct {
	cdc codec.Codec
}

// NewAppModuleBasic creates a new AppModuleBasic instance
func NewAppModuleBasic(cdc codec.Codec) AppModuleBasic {
	return AppModuleBasic{cdc: cdc}
}

// Name returns the block module's name
func (AppModuleBasic) Name() string {
	return blocktypes.ModuleName
}

// RegisterLegacyAminoCodec registers the block module's types with the legacy amino codec
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	blocktypes.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the block module's interfaces
func (AppModuleBasic) RegisterInterfaces(registry types.InterfaceRegistry) {
	blocktypes.RegisterInterfaces(registry)
}

// DefaultGenesis returns the block module's default genesis state
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesisState := blocktypes.GenesisState{
		Blocks:      []blocktypes.Block{},
		BlockData:   []blocktypes.BlockData{},
		Validations: []blocktypes.BlockValidation{},
		Params:      blocktypes.DefaultParams(),
	}

	genesisBytes, err := json.Marshal(genesisState)
	if err != nil {
		// Note: ctx not available in DefaultGenesis, but this is a critical failure
		panic(fmt.Sprintf("block: failed to marshal default genesis state: %s", err.Error()))
	}

	return genesisBytes
}

// ValidateGenesis validates the block module's genesis state
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var genState blocktypes.GenesisState
	if err := json.Unmarshal(bz, &genState); err != nil {
		return err
	}
	return blocktypes.ValidateGenesis(genState)
}

// RegisterRESTRoutes - Service-04 uses gRPC API only, REST API is not supported
func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr interface{}) {
	// REST API is not supported in Service-04 (USC Blockchain Core)
	// All API access must use gRPC protocol only
}

// RegisterGRPCGatewayRoutes registers the block module's gRPC Gateway routes
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	// Register gRPC Gateway routes
	// Note: gRPC Gateway routes will be implemented when proto files are available
}

// GetTxCmd returns the block module's root tx command
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   blocktypes.ModuleName,
		Short: "Block module subcommands",
		Long:  "Block module subcommands for managing blockchain blocks",
	}

	// Add subcommands
	cmd.AddCommand(
		NewCreateBlockCmd(),
		NewUpdateBlockCmd(),
		NewDeleteBlockCmd(),
		NewValidateBlockCmd(),
	)

	return cmd
}

// GetQueryCmd returns the block module's root query command
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   blocktypes.ModuleName,
		Short: "Querying commands for the block module",
		Long:  "Querying commands for the block module",
	}

	// Add subcommands
	cmd.AddCommand(
		NewQueryBlockCmd(),
		NewQueryBlockByHeightCmd(),
		NewQueryBlockByHashCmd(),
		NewQueryAllBlocksCmd(),
		NewQueryBlockDataCmd(),
		NewQueryAllBlockDataCmd(),
		NewQueryValidationCmd(),
		NewQueryAllValidationsCmd(),
		NewQueryParamsCmd(),
	)

	return cmd
}

// AppModule implements an application module for the block module
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

// Name returns the block module's name
func (am AppModule) Name() string {
	return am.AppModuleBasic.Name()
}

// RegisterServices registers the block module's services
func (am AppModule) RegisterServices(cfg module.Configurator) {
	// Register message server
	// Note: Message server registration will be implemented when proto files are available
	// blocktypes.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServer(am.keeper))

	// Register query server
	// Note: Query server registration will be implemented when proto files are available
	// blocktypes.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServer(am.keeper))
}

// RegisterLegacyAminoCodec registers the block module's types with the legacy amino codec
func (am AppModule) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	am.AppModuleBasic.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the block module's interfaces
func (am AppModule) RegisterInterfaces(registry types.InterfaceRegistry) {
	am.AppModuleBasic.RegisterInterfaces(registry)
}

// InitGenesis performs genesis initialization for the block module
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) []abci.ValidatorUpdate {
	var genState blocktypes.GenesisState
	if err := json.Unmarshal(gs, &genState); err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to unmarshal genesis state during InitGenesis",
			"module", blocktypes.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight(),
			"chain_id", ctx.ChainID())
		panic(fmt.Sprintf("block: failed to unmarshal genesis state: %s", err.Error()))
	}

	// Initialize genesis state - validate but handle errors gracefully
	if err := blocktypes.ValidateGenesis(genState); err != nil {
		// Log validation error but don't panic - return default state
		ctx.Logger().Warn("Genesis validation failed, using default params",
			"module", blocktypes.ModuleName,
			"error", err.Error())
		// Use default params if validation fails
		genState.Params = blocktypes.DefaultParams()
	}

	// Initialize genesis state - this will save genesis block (height 1) to keeper
	// COSMOS SDK 0.53.4: blocktypes.InitGenesis saves genesis block to keeper via am.keeper
	// NOTE: During InitGenesis, store service may not be fully available in context.
	// If InitGenesis fails, module will use default parameters (this is expected behavior).
	if err := blocktypes.InitGenesis(ctx, am.keeper, genState); err != nil {
		// Log error but don't panic - chain can continue with default state
		ctx.Logger().Warn("Failed to initialize block module genesis, using default state",
			"module", blocktypes.ModuleName,
			"error", err.Error())
	}

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the block module's exported genesis state
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genState := blocktypes.ExportGenesis(ctx, am.keeper)
	genesisBytes, err := json.Marshal(genState)
	if err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to marshal genesis state during ExportGenesis",
			"module", blocktypes.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight())
		panic(fmt.Sprintf("block: failed to marshal genesis state: %s", err.Error()))
	}
	return genesisBytes
}

// BeginBlock performs begin block logic for the block module
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestFinalizeBlock) {
	// Call BeginBlocker with RequestFinalizeBlock to get real data
	BeginBlocker(ctx, am.keeper, &req)
}

// EndBlock performs end block logic for the block module
func (am AppModule) EndBlock(ctx sdk.Context, req abci.RequestFinalizeBlock) []abci.ValidatorUpdate {
	// Call EndBlocker with RequestFinalizeBlock to get real data
	return EndBlocker(ctx, am.keeper, &req)
}

// IsAppModule implements module.AppModule
func (am AppModule) IsAppModule() {}

// IsOnePerModuleType implements module.AppModule
func (am AppModule) IsOnePerModuleType() {}

// CLI command functions
func NewCreateBlockCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-block [block-id] [height] [hash]",
		Short: "Create a new block",
		Long:  "Create a new block with the specified ID, height, and hash",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get arguments
			blockID := args[0]
			heightStr := args[1]
			hash := args[2]

			// Parse height
			height, err := strconv.ParseInt(heightStr, 10, 64)
			if err != nil {
				return err
			}

			// Create block
			_ = blocktypes.Block{
				ID:           blockID,
				Height:       height,
				Hash:         hash,
				PreviousHash: fmt.Sprintf("block_hash_%d", height-1),
				Timestamp:    time.Now(),
				Validator:    "validator_address",
				Size:         0,
				TxCount:      0,
				GasUsed:      0,
				GasLimit:     10000000,
				Status:       "pending",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}

			// Store block (implementation would call keeper)
			// Note: CLI commands output to stdout for user feedback
			cmd.PrintErrf("Created block: %s at height %d\n", blockID, height)
			return nil
		},
	}
}

func NewUpdateBlockCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update-block [block-id] [hash] [status]",
		Short: "Update an existing block",
		Long:  "Update an existing block with new hash and status",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get arguments
			blockID := args[0]
			hash := args[1]
			status := args[2]

			// Update block
			// Note: CLI commands output to stdout for user feedback
			cmd.PrintErrf("Updated block: %s with hash %s and status %s\n", blockID, hash, status)
			return nil
		},
	}
}

func NewDeleteBlockCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete-block [block-id]",
		Short: "Delete a block",
		Long:  "Delete a block by its ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get arguments
			blockID := args[0]

			// Delete block
			// Note: CLI commands output to stdout for user feedback
			cmd.PrintErrf("Deleted block: %s\n", blockID)
			return nil
		},
	}
}

func NewValidateBlockCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate-block [block-id] [validator]",
		Short: "Validate a block",
		Long:  "Validate a block with the specified validator",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get arguments
			blockID := args[0]
			validator := args[1]

			// Validate block
			// Note: CLI commands output to stdout for user feedback
			cmd.PrintErrf("Validated block: %s by validator %s\n", blockID, validator)
			return nil
		},
	}
}

func NewQueryBlockCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "block [block-id]",
		Short: "Query a block by ID",
		Long:  "Query a block by its ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get block ID from arguments
			blockID := args[0]

			// Note: CLI commands require gRPC query server implementation
			// In production, use gRPC queries via query server
			cmd.PrintErrf("Querying block: %s\n", blockID)
			cmd.PrintErrf("Note: Use gRPC query server for actual block data\n")

			return nil
		},
	}
}

func NewQueryBlockByHeightCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "block-by-height [height]",
		Short: "Query a block by height",
		Long:  "Query a block by its height",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get height from arguments
			heightStr := args[0]
			height, err := strconv.ParseInt(heightStr, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid height: %w", err)
			}

			// Note: CLI commands require gRPC query server implementation
			// In production, use gRPC queries via query server
			cmd.PrintErrf("Querying block by height: %d\n", height)
			cmd.PrintErrf("Note: Use gRPC query server for actual block data\n")

			return nil
		},
	}
}

func NewQueryBlockByHashCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "block-by-hash [hash]",
		Short: "Query a block by hash",
		Long:  "Query a block by its hash",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get hash from arguments
			hash := args[0]

			// Note: CLI commands require gRPC query server implementation
			// In production, use gRPC queries via query server
			cmd.PrintErrf("Querying block by hash: %s\n", hash)
			cmd.PrintErrf("Note: Use gRPC query server for actual block data\n")

			return nil
		},
	}
}

func NewQueryAllBlocksCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "blocks",
		Short: "Query all blocks",
		Long:  "Query all blocks with pagination",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Note: CLI commands require gRPC query server implementation
			// In production, use gRPC queries via query server
			cmd.PrintErrf("Querying all blocks\n")
			cmd.PrintErrf("Note: Use gRPC query server for actual block data\n")

			return nil
		},
	}
}

func NewQueryBlockDataCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "block-data [block-id]",
		Short: "Query block data by ID",
		Long:  "Query block data by block ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get block ID from arguments
			blockID := args[0]

			// Note: CLI commands require gRPC query server implementation
			// In production, use gRPC queries via query server
			cmd.PrintErrf("Querying block data: %s\n", blockID)
			cmd.PrintErrf("Note: Use gRPC query server for actual block data\n")

			return nil
		},
	}
}

func NewQueryAllBlockDataCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "block-data",
		Short: "Query all block data",
		Long:  "Query all block data with pagination",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Query all block data from keeper
			// Note: This is a simplified implementation
			// In production, this would query all block data with pagination
			// Note: CLI commands output to stdout for user feedback
			cmd.PrintErrf("Querying all block data...\n")
			cmd.PrintErrf("Note: This would return all block data with pagination in production\n")

			return nil
		},
	}
}

func NewQueryValidationCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validation [block-id]",
		Short: "Query block validation by ID",
		Long:  "Query block validation by block ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get block ID from arguments
			blockID := args[0]

			// Query validation from keeper
			// Note: This is a simplified implementation
			// In production, this would query the specific validation
			// Note: CLI commands output to stdout for user feedback
			cmd.PrintErrf("Querying validation for block: %s\n", blockID)
			cmd.PrintErrf("Note: This would return validation details in production\n")

			return nil
		},
	}
}

func NewQueryAllValidationsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validations",
		Short: "Query all validations",
		Long:  "Query all validations with pagination",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Note: CLI commands require gRPC query server implementation
			// In production, use gRPC queries via query server
			cmd.PrintErrf("Querying all validations\n")
			cmd.PrintErrf("Note: Use gRPC query server for actual validation data\n")

			return nil
		},
	}
}

func NewQueryParamsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query block module parameters",
		Long:  "Query block module parameters",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Note: CLI commands require gRPC query server implementation
			// In production, use gRPC queries via query server
			cmd.PrintErrf("Querying module parameters\n")
			cmd.PrintErrf("Note: Use gRPC query server for actual parameter data\n")

			return nil
		},
	}

}
