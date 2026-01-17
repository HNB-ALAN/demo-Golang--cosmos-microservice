package types

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

// RegisterCodec registers the store module's types with the legacy amino codec
func RegisterCodec(cdc codec.Codec) {
	// Register interfaces
	RegisterInterfaces(cdc.InterfaceRegistry())
}

// RegisterInterfaces registers the store module interfaces
func RegisterInterfaces(registry types.InterfaceRegistry) {
	// TODO: Register store types when proto messages are implemented
	// registry.RegisterImplementations(
	//     (*codec.ProtoMarshaler)(nil),
	//     &StoredData{},
	//     &Store{},
	//     &Backup{},
	//     &Restore{},
	//     &StoreIndex{},
	//     &StoreQuery{},
	//     &StoreTransaction{},
	//     &Params{},
	// )
}

// RegisterLegacyAminoCodec registers the store module types with the legacy amino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// Register store types with amino codec
	cdc.RegisterConcrete(&StoredData{}, "store/StoredData", nil)
	cdc.RegisterConcrete(&Store{}, "store/Store", nil)
	cdc.RegisterConcrete(&Backup{}, "store/Backup", nil)
	cdc.RegisterConcrete(&Restore{}, "store/Restore", nil)
	cdc.RegisterConcrete(&StoreIndex{}, "store/StoreIndex", nil)
	cdc.RegisterConcrete(&StoreQuery{}, "store/StoreQuery", nil)
	cdc.RegisterConcrete(&StoreTransaction{}, "store/StoreTransaction", nil)
	cdc.RegisterConcrete(&Params{}, "store/Params", nil)
}

// NewQueryClient creates a new query client for the store module
func NewQueryClient(clientCtx client.Context) interface{} {
	// TODO: Implement query client creation
	return nil
}

// RegisterMsgServer registers the store module's Msg service
func RegisterMsgServer(server grpc.ServiceRegistrar, msgServer interface{}) {
	// The actual registration is done via module.Configurator in module.go
	// This function exists for compatibility but the real registration happens
	// when keeper.NewMsgServerImpl is passed to module.Configurator.MsgServer()
}

// RegisterQueryServer registers the store module's Query service
func RegisterQueryServer(server grpc.ServiceRegistrar, queryServer interface{}) {
	// The actual registration is done via module.Configurator in module.go
	// This function exists for compatibility but the real registration happens
	// when keeper.NewQueryServerImpl is passed to module.Configurator.QueryServer()
}

// RegisterQueryHandlerClient registers the store module's Query service with the gRPC Gateway
func RegisterQueryHandlerClient(ctx context.Context, mux *runtime.ServeMux, client interface{}) error {
	// TODO: Implement query handler client registration
	return nil
}
