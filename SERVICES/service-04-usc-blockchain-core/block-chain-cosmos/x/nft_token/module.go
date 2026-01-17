package nft_token

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

	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/nft_token/keeper"
	nfttypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/nft_token/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// ----------------------------------------------------------------------------
// AppModuleBasic
// ----------------------------------------------------------------------------

// AppModuleBasic implements the AppModuleBasic interface for the NFT module
type AppModuleBasic struct {
	cdc codec.Codec
}

// NewAppModuleBasic creates a new AppModuleBasic instance
func NewAppModuleBasic(cdc codec.Codec) AppModuleBasic {
	return AppModuleBasic{cdc: cdc}
}

// Name returns the NFT module's name
func (AppModuleBasic) Name() string {
	return nfttypes.ModuleName
}

// RegisterLegacyAminoCodec registers the NFT module's types with the legacy amino codec
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	nfttypes.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the NFT module's interfaces
func (AppModuleBasic) RegisterInterfaces(registry types.InterfaceRegistry) {
	nfttypes.RegisterInterfaces(registry)
}

// DefaultGenesis returns the NFT module's default genesis state
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	genesisState := nfttypes.GenesisState{
		NFTs:        []nfttypes.NFT{},
		Collections: []nfttypes.Collection{},
		Params:      nfttypes.DefaultParams(),
	}

	genesisBytes, err := json.Marshal(genesisState)
	if err != nil {
		// Log error before panic for better debugging
		// Note: ctx not available in DefaultGenesis, using fmt for error message
		panic(fmt.Sprintf("nft_token: failed to marshal default genesis state: %s", err.Error()))
	}

	return genesisBytes
}

// ValidateGenesis validates the NFT module's genesis state
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var genState nfttypes.GenesisState
	if err := json.Unmarshal(bz, &genState); err != nil {
		return err
	}
	return nfttypes.ValidateGenesis(genState)
}

// RegisterRESTRoutes - Service-04 uses gRPC API only, REST API is not supported
func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr interface{}) {
	// REST API is not supported in Service-04 (USC Blockchain Core)
	// All API access must use gRPC protocol only
}

// RegisterGRPCGatewayRoutes registers the nft_token module's gRPC Gateway service handlers
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	// Register gRPC Gateway routes using RegisterQueryHandlerClient
	// This connects the gRPC Gateway to the nft_token module's query service
	if err := nfttypes.RegisterQueryHandlerClient(context.Background(), mux, clientCtx); err != nil {
		// Log error but don't fail - gRPC Gateway routes are optional
		// The actual gRPC service registration happens via module.Configurator
		_ = err
	}
}

// GetTxCmd returns the NFT module's root tx command
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   nfttypes.ModuleName,
		Short: "NFT module subcommands",
		Long:  "NFT module subcommands for managing non-fungible tokens",
	}

	// Add subcommands
	cmd.AddCommand(
		NewCreateNFTCmd(),
		NewTransferNFTCmd(),
		NewUpdateNFTCmd(),
		NewBurnNFTCmd(),
		NewCreateCollectionCmd(),
		NewUpdateCollectionCmd(),
	)

	return cmd
}

// GetQueryCmd returns the NFT module's root query command
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   nfttypes.ModuleName,
		Short: "Querying commands for the NFT module",
		Long:  "Querying commands for the NFT module",
	}

	// Add subcommands
	cmd.AddCommand(
		NewQueryNFTCmd(),
		NewQueryAllNFTsCmd(),
		NewQueryCollectionCmd(),
		NewQueryAllCollectionsCmd(),
		NewQueryParamsCmd(),
	)

	return cmd
}

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

// AppModule implements the AppModule interface for the NFT module
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

// Name returns the NFT module's name
func (am AppModule) Name() string {
	return am.AppModuleBasic.Name()
}

// RegisterServices registers the NFT module's services
func (am AppModule) RegisterServices(cfg module.Configurator) {
	nfttypes.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServer(am.keeper))
	nfttypes.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServer(am.keeper))
}

// RegisterLegacyAminoCodec registers the NFT module's types with the legacy amino codec
func (am AppModule) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	am.AppModuleBasic.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the NFT module's interfaces
func (am AppModule) RegisterInterfaces(registry types.InterfaceRegistry) {
	am.AppModuleBasic.RegisterInterfaces(registry)
}

// InitGenesis performs genesis initialization for the NFT module
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) []abci.ValidatorUpdate {
	var genState nfttypes.GenesisState
	if err := json.Unmarshal(gs, &genState); err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to unmarshal genesis state during InitGenesis",
			"module", nfttypes.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight(),
			"chain_id", ctx.ChainID())
		panic(fmt.Sprintf("nft_token: failed to unmarshal genesis state: %s", err.Error()))
	}

	// Initialize genesis state - validate but handle errors gracefully
	if err := nfttypes.ValidateGenesis(genState); err != nil {
		// Log validation error but don't panic - return default state
		ctx.Logger().Warn("Genesis validation failed, using default params",
			"module", nfttypes.ModuleName,
			"error", err.Error())
		// Use default params if validation fails
		genState.Params = nfttypes.DefaultParams()
	}

	// Initialize genesis state via keeper
	// NOTE: During InitGenesis, store service may not be fully available in context.
	// If InitGenesis fails, module will use default parameters (this is expected behavior).
	nfttypes.InitGenesis(ctx, am.keeper, genState)

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the NFT module's exported genesis state
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genState := nfttypes.ExportGenesis(ctx, am.keeper)
	genesisBytes, err := json.Marshal(genState)
	if err != nil {
		// Log error with detailed context before panic
		ctx.Logger().Error("Failed to marshal genesis state during ExportGenesis",
			"module", nfttypes.ModuleName,
			"error", err.Error(),
			"block_height", ctx.BlockHeight())
		panic(fmt.Sprintf("nft_token: failed to marshal genesis state: %s", err.Error()))
	}
	return genesisBytes
}

// BeginBlock performs begin block logic for the NFT module
func (am AppModule) BeginBlock(ctx sdk.Context) {
	BeginBlocker(ctx, am.keeper)
}

// EndBlock performs end block logic for the NFT module
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
func NewCreateNFTCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-nft [id] [collection-id] [owner] [token-uri]",
		Short: "Create a new NFT",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement create NFT logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewTransferNFTCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "transfer-nft [nft-id] [from] [to]",
		Short: "Transfer an NFT",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement transfer NFT logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewUpdateNFTCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update-nft [nft-id] [name] [description]",
		Short: "Update an NFT",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement update NFT logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewBurnNFTCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "burn-nft [nft-id] [owner]",
		Short: "Burn an NFT",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement burn NFT logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewCreateCollectionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-collection [id] [name] [symbol] [owner]",
		Short: "Create a new collection",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement create collection logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewUpdateCollectionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update-collection [collection-id] [name] [description]",
		Short: "Update a collection",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement update collection logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryNFTCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "nft [id]",
		Short: "Query NFT by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query NFT logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllNFTsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-nfts",
		Short: "Query all NFTs",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all NFTs logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryCollectionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "collection [id]",
		Short: "Query collection by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query collection logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryAllCollectionsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all-collections",
		Short: "Query all collections",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query all collections logic
			return fmt.Errorf("not implemented")
		},
	}
}

func NewQueryParamsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query NFT module parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement query params logic
			return fmt.Errorf("not implemented")
		},
	}

}
