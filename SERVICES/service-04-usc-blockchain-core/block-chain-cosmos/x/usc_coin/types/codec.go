package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/usc_coin/v1/usc/usc_coin/v1"
)

// RegisterLegacyAminoCodec registers the USC module's types with the given codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// Register blockchain-proto message types with amino codec
	cdc.RegisterConcrete(&blockchainproto.MsgTransferUSC{}, "usc/coin/MsgTransferUSC", nil)
	cdc.RegisterConcrete(&blockchainproto.MsgMintUSC{}, "usc/coin/MsgMintUSC", nil)
	cdc.RegisterConcrete(&blockchainproto.MsgBurnUSC{}, "usc/coin/MsgBurnUSC", nil)

	// Register custom USC types with amino codec
	cdc.RegisterConcrete(&Balance{}, "usc/Balance", nil)
	cdc.RegisterConcrete(&Transfer{}, "usc/Transfer", nil)
	cdc.RegisterConcrete(&Params{}, "usc/Params", nil)
	cdc.RegisterConcrete(&GenesisState{}, "usc/GenesisState", nil)
}

// RegisterInterfaces registers the USC module's interfaces with the given interface registry
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	// Register blockchain-proto message types with the interface registry
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&blockchainproto.MsgTransferUSC{},
		&blockchainproto.MsgMintUSC{},
		&blockchainproto.MsgBurnUSC{},
	)
}

// RegisterCodec registers the USC module's types with the given codec
func RegisterCodec(cdc codec.Codec) {
	// Register interfaces
	RegisterInterfaces(cdc.InterfaceRegistry())
}

// Legacy Amino fully removed for this module; Protobuf-only
