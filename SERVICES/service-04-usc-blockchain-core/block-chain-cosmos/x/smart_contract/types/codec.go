package types

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

// RegisterCodec registers the contract module's types with the legacy amino codec
func RegisterCodec(cdc codec.Codec) {
	RegisterInterfaces(cdc.InterfaceRegistry())
}

// RegisterInterfaces registers the contract module interfaces
func RegisterInterfaces(registry types.InterfaceRegistry) {
	// TODO: Register contract types when proto messages are implemented
}

// RegisterLegacyAminoCodec registers the contract module types with the legacy amino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&SmartContract{}, "contract/SmartContract", nil)
	cdc.RegisterConcrete(&ContractExecution{}, "contract/ContractExecution", nil)
	cdc.RegisterConcrete(&ContractDeployment{}, "contract/ContractDeployment", nil)
	cdc.RegisterConcrete(&ContractUpgrade{}, "contract/ContractUpgrade", nil)
	cdc.RegisterConcrete(&ContractMigration{}, "contract/ContractMigration", nil)
	cdc.RegisterConcrete(&Params{}, "contract/Params", nil)
}

// NewQueryClient creates a new query client for the contract module
func NewQueryClient(clientCtx client.Context) interface{} {
	// TODO: Implement query client creation
	return nil
}

// RegisterMsgServer registers the contract module's Msg service
func RegisterMsgServer(server grpc.ServiceRegistrar, msgServer interface{}) {
	// TODO: Implement message server registration
}

// RegisterQueryServer registers the contract module's Query service
func RegisterQueryServer(server grpc.ServiceRegistrar, queryServer interface{}) {
	// TODO: Implement query server registration
}

// RegisterQueryHandlerClient registers the contract module's Query service with the gRPC Gateway
func RegisterQueryHandlerClient(ctx context.Context, mux *runtime.ServeMux, client interface{}) error {
	// TODO: Implement query handler client registration
	return nil
}
