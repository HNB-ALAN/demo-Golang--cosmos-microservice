package types

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

// RegisterCodec registers the monitoring module's types with the legacy amino codec
func RegisterCodec(cdc codec.Codec) {
	// Register interfaces
	RegisterInterfaces(cdc.InterfaceRegistry())
}

// RegisterInterfaces registers the monitoring module interfaces
func RegisterInterfaces(registry types.InterfaceRegistry) {
	// TODO: Register monitoring types when proto messages are implemented
	// registry.RegisterImplementations(
	//     (*codec.ProtoMarshaler)(nil),
	//     &Metric{},
	//     &Alert{},
	//     &PerformanceData{},
	//     &SystemHealth{},
	//     &MonitoringConfig{},
	//     &Params{},
	// )
}

// RegisterLegacyAminoCodec registers the monitoring module types with the legacy amino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// Register monitoring types with amino codec
	cdc.RegisterConcrete(&Metric{}, "monitoring/Metric", nil)
	cdc.RegisterConcrete(&Alert{}, "monitoring/Alert", nil)
	cdc.RegisterConcrete(&PerformanceData{}, "monitoring/PerformanceData", nil)
	cdc.RegisterConcrete(&SystemHealth{}, "monitoring/SystemHealth", nil)
	cdc.RegisterConcrete(&MonitoringConfig{}, "monitoring/MonitoringConfig", nil)
	cdc.RegisterConcrete(&Params{}, "monitoring/Params", nil)
}

// NewQueryClient creates a new query client for the monitoring module
func NewQueryClient(clientCtx client.Context) interface{} {
	// TODO: Implement query client creation
	return nil
}

// RegisterMsgServer registers the monitoring module's Msg service
func RegisterMsgServer(server grpc.ServiceRegistrar, msgServer interface{}) {
	// The actual registration is done via module.Configurator in module.go
	// This function exists for compatibility but the real registration happens
	// when keeper.NewMsgServerImpl is passed to module.Configurator.MsgServer()
}

// RegisterQueryServer registers the monitoring module's Query service
func RegisterQueryServer(server grpc.ServiceRegistrar, queryServer interface{}) {
	// The actual registration is done via module.Configurator in module.go
	// This function exists for compatibility but the real registration happens
	// when keeper.NewQueryServerImpl is passed to module.Configurator.QueryServer()
}

// RegisterQueryHandlerClient registers the monitoring module's Query service with the gRPC Gateway
func RegisterQueryHandlerClient(ctx context.Context, mux *runtime.ServeMux, client interface{}) error {
	// TODO: Implement query handler client registration
	return nil
}
