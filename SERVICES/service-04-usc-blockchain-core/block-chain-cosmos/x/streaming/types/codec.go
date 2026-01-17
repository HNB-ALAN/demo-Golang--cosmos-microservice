package types

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

// RegisterCodec registers the streaming module's types with the legacy amino codec
func RegisterCodec(cdc codec.Codec) {
	RegisterInterfaces(cdc.InterfaceRegistry())
}

// RegisterInterfaces registers the streaming module interfaces
func RegisterInterfaces(registry types.InterfaceRegistry) {
	// TODO: Register streaming types when proto messages are implemented
}

// RegisterLegacyAminoCodec registers the streaming module types with the legacy amino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&Stream{}, "streaming/Stream", nil)
	cdc.RegisterConcrete(&StreamViewer{}, "streaming/StreamViewer", nil)
	cdc.RegisterConcrete(&StreamQualityMetrics{}, "streaming/StreamQualityMetrics", nil)
	cdc.RegisterConcrete(&StreamAnalytics{}, "streaming/StreamAnalytics", nil)
	cdc.RegisterConcrete(&StreamChat{}, "streaming/StreamChat", nil)
	cdc.RegisterConcrete(&StreamDonation{}, "streaming/StreamDonation", nil)
	cdc.RegisterConcrete(&StreamModeration{}, "streaming/StreamModeration", nil)
	cdc.RegisterConcrete(&StreamEvent{}, "streaming/StreamEvent", nil)
	cdc.RegisterConcrete(&StreamQuality{}, "streaming/StreamQuality", nil)
	cdc.RegisterConcrete(&Params{}, "streaming/Params", nil)
}

// NewQueryClient creates a new query client for the streaming module
func NewQueryClient(clientCtx client.Context) interface{} {
	// TODO: Implement query client creation
	return nil
}

// RegisterMsgServer registers the streaming module's Msg service
func RegisterMsgServer(server grpc.ServiceRegistrar, msgServer interface{}) {
	// The actual registration is done via module.Configurator in module.go
	// This function exists for compatibility but the real registration happens
	// when keeper.NewMsgServer is passed to module.Configurator.MsgServer()
}

// RegisterQueryServer registers the streaming module's Query service
func RegisterQueryServer(server grpc.ServiceRegistrar, queryServer interface{}) {
	// The actual registration is done via module.Configurator in module.go
	// This function exists for compatibility but the real registration happens
	// when keeper.NewQueryServer is passed to module.Configurator.QueryServer()
}

// RegisterQueryHandlerClient registers the streaming module's Query service with the gRPC Gateway
func RegisterQueryHandlerClient(ctx context.Context, mux *runtime.ServeMux, client interface{}) error {
	// Note: This function is called from RegisterGRPCGatewayRoutes in module.go
	// The actual gRPC Gateway route registration would require generated code from proto service definitions
	// For now, this is a placeholder that allows the module to compile
	// When proto service definitions are available, this should register routes like:
	// return streamingproto.RegisterQueryHandlerClient(ctx, mux, client)
	
	// Since streaming proto doesn't have a Query service definition yet,
	// we return nil to allow compilation. Routes will be registered when proto is updated.
	return nil
}
