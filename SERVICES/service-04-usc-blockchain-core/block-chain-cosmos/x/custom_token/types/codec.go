package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc"

	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/custom_token/v1/usc/custom_token/v1"
)

// RegisterLegacyAminoCodec registers the custom_token module's types with the given codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&CustomToken{}, "custom_token/CustomToken", nil)
	cdc.RegisterConcrete(&TokenBalance{}, "custom_token/TokenBalance", nil)
	cdc.RegisterConcrete(&TokenTransfer{}, "custom_token/TokenTransfer", nil)
	cdc.RegisterConcrete(&Params{}, "custom_token/Params", nil)
	cdc.RegisterConcrete(&GenesisState{}, "custom_token/GenesisState", nil)
}

// RegisterInterfaces registers the custom_token module's interfaces with the given interface registry
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&blockchainproto.MsgCreateToken{},
		&blockchainproto.MsgMintToken{},
		&blockchainproto.MsgBurnToken{},
		&blockchainproto.MsgTransferToken{},
		&blockchainproto.MsgUpdateToken{},
	)
}

// RegisterCodec registers the custom_token module's types with the given codec
func RegisterCodec(cdc codec.Codec) {
	// Register interfaces
	RegisterInterfaces(cdc.InterfaceRegistry())
}

// Legacy Amino fully removed for this module; Protobuf-only

// RegisterMsgServer registers the custom token module's Msg service
func RegisterMsgServer(server grpc.ServiceRegistrar, msgServer interface{}) {
	// The actual registration is done via module.Configurator in module.go
	// This function exists for compatibility but the real registration happens
	// when keeper.NewMsgServerImpl is passed to module.Configurator.MsgServer()
}

// RegisterQueryServer registers the custom token module's Query service
func RegisterQueryServer(server grpc.ServiceRegistrar, queryServer interface{}) {
	// The actual registration is done via module.Configurator in module.go
	// This function exists for compatibility but the real registration happens
	// when keeper.NewQueryServerImpl is passed to module.Configurator.QueryServer()
}
