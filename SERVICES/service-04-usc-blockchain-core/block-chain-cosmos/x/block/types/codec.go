package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/block/v1/usc/block/v1"
)

// RegisterCodec registers the block module types with the given codec
func RegisterCodec(cdc codec.Codec) {
	// Register interfaces - check if InterfaceRegistry method exists
	if registry := cdc.InterfaceRegistry(); registry != nil {
		RegisterInterfaces(registry)
	}
}

// RegisterCodecWithInterfaceRegistry registers the block module types with the given codec and interface registry
func RegisterCodecWithInterfaceRegistry(cdc codec.Codec, registry codectypes.InterfaceRegistry) {
	// Register interfaces
	RegisterInterfaces(registry)
}

// RegisterInterfaces registers the block module interfaces
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	// Register blockchain-proto message types with the interface registry
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&blockchainproto.MsgCreateBlock{},
		&blockchainproto.MsgValidateBlock{},
		&blockchainproto.MsgFinalizeBlock{},
	)
}

// RegisterLegacyAminoCodec registers the block module types with the legacy amino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// Register blockchain-proto message types with amino codec
	cdc.RegisterConcrete(&blockchainproto.MsgCreateBlock{}, "usc/block/MsgCreateBlock", nil)
	cdc.RegisterConcrete(&blockchainproto.MsgValidateBlock{}, "usc/block/MsgValidateBlock", nil)
	cdc.RegisterConcrete(&blockchainproto.MsgFinalizeBlock{}, "usc/block/MsgFinalizeBlock", nil)

	// Register custom block types with amino codec
	cdc.RegisterConcrete(&Block{}, "block/Block", nil)
	cdc.RegisterConcrete(&BlockData{}, "block/BlockData", nil)
	cdc.RegisterConcrete(&BlockValidation{}, "block/BlockValidation", nil)
	cdc.RegisterConcrete(&BlockHeader{}, "block/BlockHeader", nil)
	cdc.RegisterConcrete(&Params{}, "block/Params", nil)
}

// RegisterMsgServer registers the block module message server
func RegisterMsgServer(server interface{}) {
	// Register message server with the given server interface
	// This function will be called during module initialization
	// to register the block module's message handlers
	// Implementation will be added when proto files are available
}

// RegisterQueryServer registers the block module query server
func RegisterQueryServer(server interface{}) {
	// Register query server with the given server interface
	// This function will be called during module initialization
	// to register the block module's query handlers
	// Implementation will be added when proto files are available
}

// RegisterQueryHandlerClient registers the block module query handler client
func RegisterQueryHandlerClient(client interface{}) {
	// Query handler client registration will be implemented when proto files are available
	// This is a placeholder for future proto-based query handler client registration
	// Implementation will be added when proto files are available
}

// NewQueryClient creates a new block module query client
func NewQueryClient(client interface{}) interface{} {
	// Query client creation will be implemented when proto files are available
	// This is a placeholder for future proto-based query client creation
	// Implementation will be added when proto files are available
	return nil
}
