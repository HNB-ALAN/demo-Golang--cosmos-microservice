package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/validator/v1/usc/validator/v1"
	"google.golang.org/grpc"
)

// RegisterLegacyAminoCodec registers the validator module's types with the given codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&Validator{}, "validator/Validator", nil)
	cdc.RegisterConcrete(&Delegation{}, "validator/Delegation", nil)
	cdc.RegisterConcrete(&Params{}, "validator/Params", nil)
	cdc.RegisterConcrete(&GenesisState{}, "validator/GenesisState", nil)

	// Register blockchain-proto message types
	cdc.RegisterConcrete(&blockchainproto.MsgCreateValidator{}, "validator/MsgCreateValidator", nil)
	cdc.RegisterConcrete(&blockchainproto.MsgUpdateValidator{}, "validator/MsgUpdateValidator", nil)
	cdc.RegisterConcrete(&blockchainproto.MsgDelegateValidator{}, "validator/MsgDelegateValidator", nil)
	cdc.RegisterConcrete(&blockchainproto.MsgUndelegateValidator{}, "validator/MsgUndelegateValidator", nil)
	cdc.RegisterConcrete(&blockchainproto.MsgSlashValidator{}, "validator/MsgSlashValidator", nil)
}

// RegisterInterfaces registers the validator module's interfaces with the given interface registry
func RegisterInterfaces(registry types.InterfaceRegistry) {
	// Register blockchain-proto message types as implementations of sdk.Msg
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&blockchainproto.MsgCreateValidator{},
		&blockchainproto.MsgUpdateValidator{},
		&blockchainproto.MsgDelegateValidator{},
		&blockchainproto.MsgUndelegateValidator{},
		&blockchainproto.MsgSlashValidator{},
	)
}

// RegisterCodec registers the validator module's types with the given codec
func RegisterCodec(cdc codec.Codec) {
	// Register interfaces
	RegisterInterfaces(cdc.InterfaceRegistry())
}

// Legacy Amino fully removed for this module; Protobuf-only

// RegisterMsgServer registers the validator module's Msg service
func RegisterMsgServer(server grpc.ServiceRegistrar, msgServer interface{}) {
	// The actual registration is done via module.Configurator in module.go
	// This function exists for compatibility but the real registration happens
	// when keeper.NewMsgServerImpl is passed to module.Configurator.MsgServer()
}

// RegisterQueryServer registers the validator module's Query service
func RegisterQueryServer(server grpc.ServiceRegistrar, queryServer interface{}) {
	// The actual registration is done via module.Configurator in module.go
	// This function exists for compatibility but the real registration happens
	// when keeper.NewQueryServerImpl is passed to module.Configurator.QueryServer()
}
