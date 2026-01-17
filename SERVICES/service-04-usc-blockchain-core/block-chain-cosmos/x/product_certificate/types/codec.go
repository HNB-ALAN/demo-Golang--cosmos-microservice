package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/product_certificate/v1/usc/product_certificate/v1"
)

// RegisterLegacyAminoCodec registers the product_certificate module's types with the given codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// Register blockchain-proto message types with amino codec
	cdc.RegisterConcrete(&blockchainproto.MsgCreateCertificate{}, "product_certificate/MsgCreateCertificate", nil)
	cdc.RegisterConcrete(&blockchainproto.MsgUpdateCertificate{}, "product_certificate/MsgUpdateCertificate", nil)
	cdc.RegisterConcrete(&blockchainproto.MsgRevokeCertificate{}, "product_certificate/MsgRevokeCertificate", nil)
	cdc.RegisterConcrete(&blockchainproto.MsgVerifyCertificate{}, "product_certificate/MsgVerifyCertificate", nil)
	cdc.RegisterConcrete(&blockchainproto.MsgTransferCertificate{}, "product_certificate/MsgTransferCertificate", nil)

	// Register custom product_certificate types with amino codec
	cdc.RegisterConcrete(&ProductCertificate{}, "product_certificate/ProductCertificate", nil)
	cdc.RegisterConcrete(&CertificateVerification{}, "product_certificate/CertificateVerification", nil)
	cdc.RegisterConcrete(&Params{}, "product_certificate/Params", nil)
	cdc.RegisterConcrete(&GenesisState{}, "product_certificate/GenesisState", nil)
}

// RegisterInterfaces registers the product_certificate module's interfaces with the given interface registry
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	// Register blockchain-proto message types with the interface registry
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&blockchainproto.MsgCreateCertificate{},
		&blockchainproto.MsgUpdateCertificate{},
		&blockchainproto.MsgRevokeCertificate{},
		&blockchainproto.MsgVerifyCertificate{},
		&blockchainproto.MsgTransferCertificate{},
	)
}

// RegisterCodec registers the product_certificate module's types with the given codec
func RegisterCodec(cdc codec.Codec) {
	// Register interfaces
	RegisterInterfaces(cdc.InterfaceRegistry())
}

// Legacy Amino fully removed for this module; Protobuf-only

// RegisterMsgServer registers the product_certificate module's Msg service
func RegisterMsgServer(server interface{}, msgServer interface{}) {
	// The actual registration is done via module.Configurator in module.go
	// This function exists for compatibility but the real registration happens
	// when keeper.NewMsgServerImpl is passed to module.Configurator.MsgServer()
}

// RegisterQueryServer registers the product_certificate module's Query service
func RegisterQueryServer(server interface{}, queryServer interface{}) {
	// The actual registration is done via module.Configurator in module.go
	// This function exists for compatibility but the real registration happens
	// when keeper.NewQueryServerImpl is passed to module.Configurator.QueryServer()
}
