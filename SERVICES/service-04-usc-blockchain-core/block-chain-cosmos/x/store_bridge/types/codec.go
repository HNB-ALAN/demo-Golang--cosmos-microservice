package types

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

// RegisterCodec registers the bridge module's types with the legacy amino codec
func RegisterCodec(cdc codec.Codec) {
	// Register interfaces
	RegisterInterfaces(cdc.InterfaceRegistry())
}

// RegisterInterfaces registers the bridge module interfaces
func RegisterInterfaces(registry types.InterfaceRegistry) {
	// TODO: Register bridge types when proto messages are implemented
	// registry.RegisterImplementations(
	//     (*codec.ProtoMarshaler)(nil),
	//     &Bridge{},
	//     &Transfer{},
	//     &Validator{},
	//     &BridgeConfig{},
	//     &BridgeFee{},
	//     &BridgeLimit{},
	//     &BridgeEvent{},
	//     &Params{},
	// )
}

// RegisterLegacyAminoCodec registers the bridge module types with the legacy amino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// Register bridge types with amino codec
	cdc.RegisterConcrete(&Bridge{}, "bridge/Bridge", nil)
	cdc.RegisterConcrete(&Transfer{}, "bridge/Transfer", nil)
	cdc.RegisterConcrete(&Validator{}, "bridge/Validator", nil)
	cdc.RegisterConcrete(&BridgeConfig{}, "bridge/BridgeConfig", nil)
	cdc.RegisterConcrete(&BridgeFee{}, "bridge/BridgeFee", nil)
	cdc.RegisterConcrete(&BridgeLimit{}, "bridge/BridgeLimit", nil)
	cdc.RegisterConcrete(&BridgeEvent{}, "bridge/BridgeEvent", nil)
	cdc.RegisterConcrete(&Params{}, "bridge/Params", nil)
}

// NewQueryClient creates a new query client for the bridge module
func NewQueryClient(clientCtx client.Context) interface{} {
	// TODO: Implement query client creation
	return nil
}

// RegisterMsgServer registers the bridge module's Msg service
func RegisterMsgServer(server grpc.ServiceRegistrar, msgServer interface{}) {
	// The actual registration is done via module.Configurator in module.go
	// This function exists for compatibility but the real registration happens
	// when keeper.NewMsgServerImpl is passed to module.Configurator.MsgServer()
}

// RegisterQueryServer registers the bridge module's Query service
func RegisterQueryServer(server grpc.ServiceRegistrar, queryServer interface{}) {
	// The actual registration is done via module.Configurator in module.go
	// This function exists for compatibility but the real registration happens
	// when keeper.NewQueryServerImpl is passed to module.Configurator.QueryServer()
}

// RegisterQueryHandlerClient registers the bridge module's Query service with the gRPC Gateway
func RegisterQueryHandlerClient(ctx context.Context, mux *runtime.ServeMux, client interface{}) error {
	// TODO: Implement query handler client registration
	return nil
}
