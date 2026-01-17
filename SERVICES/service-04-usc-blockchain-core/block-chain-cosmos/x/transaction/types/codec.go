package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/transaction/v1/usc/transaction/v1"
	"google.golang.org/grpc"
)

// RegisterLegacyAminoCodec registers the transaction module's types on the given LegacyAmino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// Register blockchain-proto message types with amino codec
	cdc.RegisterConcrete(&blockchainproto.MsgCreateTransaction{}, "transaction/MsgCreateTransaction", nil)
	cdc.RegisterConcrete(&blockchainproto.MsgValidateTransaction{}, "transaction/MsgValidateTransaction", nil)
	cdc.RegisterConcrete(&blockchainproto.MsgExecuteTransaction{}, "transaction/MsgExecuteTransaction", nil)
	cdc.RegisterConcrete(&blockchainproto.MsgCancelTransaction{}, "transaction/MsgCancelTransaction", nil)

	// Register custom transaction types with amino codec
	cdc.RegisterConcrete(&Transaction{}, "transaction/Transaction", nil)
	cdc.RegisterConcrete(&TransactionStats{}, "transaction/TransactionStats", nil)
	cdc.RegisterConcrete(&GenesisState{}, "transaction/GenesisState", nil)
	cdc.RegisterConcrete(&Params{}, "transaction/Params", nil)
}

// RegisterInterfaces registers the transaction module's interface types
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	// Register blockchain-proto message types with the interface registry
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&blockchainproto.MsgCreateTransaction{},
		&blockchainproto.MsgValidateTransaction{},
		&blockchainproto.MsgExecuteTransaction{},
		&blockchainproto.MsgCancelTransaction{},
	)
}

// RegisterCodec registers the transaction module's types with the given codec
func RegisterCodec(cdc *codec.LegacyAmino) {}

// NewCodec creates a new codec for the transaction module
func NewCodec() *codec.LegacyAmino { return nil }

var ()

func init() {}

// RegisterMsgServer registers the transaction module's Msg service
func RegisterMsgServer(server grpc.ServiceRegistrar, msgServer interface{}) {
	// The actual registration is done via module.Configurator in module.go
	// This function exists for compatibility but the real registration happens
	// when keeper.NewMsgServerImpl is passed to module.Configurator.MsgServer()
}

// RegisterQueryServer registers the transaction module's Query service
func RegisterQueryServer(server grpc.ServiceRegistrar, queryServer interface{}) {
	// The actual registration is done via module.Configurator in module.go
	// This function exists for compatibility but the real registration happens
	// when keeper.NewQueryServerImpl is passed to module.Configurator.QueryServer()
}
