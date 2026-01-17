package types

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/network/v1/usc/network/v1"
)

// RegisterCodec registers the network module's types with the legacy amino codec
func RegisterCodec(cdc codec.Codec) {
	RegisterInterfaces(cdc.InterfaceRegistry())
}

// RegisterInterfaces registers the network module interfaces
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	// Register blockchain-proto message types
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&blockchainproto.MsgUpdateNetwork{},
		&blockchainproto.MsgJoinNetwork{},
		&blockchainproto.MsgLeaveNetwork{},
		&blockchainproto.MsgSyncNetwork{},
	)
}

// RegisterLegacyAminoCodec registers the network module types with the legacy amino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// Register blockchain-proto message types with amino codec
	cdc.RegisterConcrete(&blockchainproto.MsgUpdateNetwork{}, "network/MsgUpdateNetwork", nil)
	cdc.RegisterConcrete(&blockchainproto.MsgJoinNetwork{}, "network/MsgJoinNetwork", nil)
	cdc.RegisterConcrete(&blockchainproto.MsgLeaveNetwork{}, "network/MsgLeaveNetwork", nil)
	cdc.RegisterConcrete(&blockchainproto.MsgSyncNetwork{}, "network/MsgSyncNetwork", nil)

	// Register custom network types with amino codec
	cdc.RegisterConcrete(&Network{}, "network/Network", nil)
	cdc.RegisterConcrete(&Node{}, "network/Node", nil)
	cdc.RegisterConcrete(&Connection{}, "network/Connection", nil)
	cdc.RegisterConcrete(&NetworkSync{}, "network/NetworkSync", nil)
	cdc.RegisterConcrete(&NetworkHealth{}, "network/NetworkHealth", nil)
	cdc.RegisterConcrete(&Params{}, "network/Params", nil)
}

// NewQueryClient creates a new query client for the network module
func NewQueryClient(clientCtx client.Context) interface{} {
	// TODO: Implement query client creation
	return nil
}

// RegisterMsgServer registers the network module's Msg service
func RegisterMsgServer(server grpc.ServiceRegistrar, msgServer interface{}) {
	// The actual registration is done via module.Configurator in module.go
	// This function exists for compatibility but the real registration happens
	// when keeper.NewMsgServerImpl is passed to module.Configurator.MsgServer()
}

// RegisterQueryServer registers the network module's Query service
func RegisterQueryServer(server grpc.ServiceRegistrar, queryServer interface{}) {
	// The actual registration is done via module.Configurator in module.go
	// This function exists for compatibility but the real registration happens
	// when keeper.NewQueryServerImpl is passed to module.Configurator.QueryServer()
}

// RegisterQueryHandlerClient registers the network module's Query service with the gRPC Gateway
func RegisterQueryHandlerClient(ctx context.Context, mux *runtime.ServeMux, client interface{}) error {
	// TODO: Implement query handler client registration
	return nil
}
