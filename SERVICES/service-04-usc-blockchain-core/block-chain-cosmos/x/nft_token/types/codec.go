package types

import (
    "context"

    "github.com/cosmos/cosmos-sdk/client"
    "github.com/cosmos/cosmos-sdk/codec"
    codectypes "github.com/cosmos/cosmos-sdk/codec/types"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/grpc-ecosystem/grpc-gateway/runtime"
    "google.golang.org/grpc"

    blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/nft_token/v1/usc/nft_token/v1"
)

// RegisterCodec registers the nft module's types with the legacy amino codec
func RegisterCodec(cdc codec.Codec) {
    RegisterInterfaces(cdc.InterfaceRegistry())
}

// RegisterInterfaces registers the nft module interfaces
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
    // Register blockchain-proto message types
    registry.RegisterImplementations((*sdk.Msg)(nil),
        &blockchainproto.MsgMintNFT{},
        &blockchainproto.MsgTransferNFT{},
        &blockchainproto.MsgBurnNFT{},
        &blockchainproto.MsgUpdateNFT{},
        &blockchainproto.MsgCreateCollection{},
    )
}

// RegisterLegacyAminoCodec registers the nft module types with the legacy amino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
    cdc.RegisterConcrete(&NFT{}, "nft/NFT", nil)
    cdc.RegisterConcrete(&Collection{}, "nft/Collection", nil)
    cdc.RegisterConcrete(&Params{}, "nft/Params", nil)
}

// NewQueryClient creates a new query client for the nft module
func NewQueryClient(clientCtx client.Context) interface{} {
    return nil
}

// RegisterMsgServer registers the nft module's Msg service
func RegisterMsgServer(server grpc.ServiceRegistrar, msgServer interface{}) {
    // TODO: Implement message server registration
}

// RegisterQueryServer registers the nft module's Query service
func RegisterQueryServer(server grpc.ServiceRegistrar, queryServer interface{}) {
    // TODO: Implement query server registration
}

// RegisterQueryHandlerClient registers the nft module's Query service with the gRPC Gateway
func RegisterQueryHandlerClient(ctx context.Context, mux *runtime.ServeMux, client interface{}) error {
    return nil
}
