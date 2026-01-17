package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/product_certificate/v1/usc/product_certificate/v1"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/product_certificate/types"
)

// MsgServer defines the interface for the product_certificate module's message server
type MsgServer interface {
	CreateCertificate(context.Context, *blockchainproto.MsgCreateCertificate) (*blockchainproto.MsgCreateCertificateResponse, error)
	UpdateCertificate(context.Context, *blockchainproto.MsgUpdateCertificate) (*blockchainproto.MsgUpdateCertificateResponse, error)
	RevokeCertificate(context.Context, *blockchainproto.MsgRevokeCertificate) (*blockchainproto.MsgRevokeCertificateResponse, error)
	VerifyCertificate(context.Context, *blockchainproto.MsgVerifyCertificate) (*blockchainproto.MsgVerifyCertificateResponse, error)
	TransferCertificate(context.Context, *blockchainproto.MsgTransferCertificate) (*blockchainproto.MsgTransferCertificateResponse, error)
}

// msgServer implements the MsgServer interface
type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
func NewMsgServerImpl(keeper Keeper) MsgServer {
	return &msgServer{Keeper: keeper}
}

// CreateCertificate handles certificate creation
func (k msgServer) CreateCertificate(ctx context.Context, msg *blockchainproto.MsgCreateCertificate) (*blockchainproto.MsgCreateCertificateResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Create certificate
	cert := types.NewProductCertificate(msg.ProductId, msg.ProductId, msg.Creator, "active", msg.ProductName)

	if err := k.Keeper.CreateCertificate(sdkCtx, cert); err != nil {
		return nil, err
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProductCertificateCreated,
			sdk.NewAttribute(types.AttributeKeyCertificateID, msg.ProductId),
			sdk.NewAttribute(types.AttributeKeyProductID, msg.ProductId),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyStatus, "active"),
		),
	)

	return &blockchainproto.MsgCreateCertificateResponse{}, nil
}

// UpdateCertificate handles certificate updates
func (k msgServer) UpdateCertificate(ctx context.Context, msg *blockchainproto.MsgUpdateCertificate) (*blockchainproto.MsgUpdateCertificateResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get existing certificate
	existingCert, exists := k.Keeper.GetCertificate(sdkCtx, msg.CertificateId)
	if !exists {
		return nil, fmt.Errorf("certificate with ID %s does not exist", msg.CertificateId)
	}

	// Update certificate fields
	existingCert.ProductID = msg.CertificateId
	existingCert.Owner = msg.Updater
	existingCert.Status = "active"
	existingCert.Metadata = msg.NewProductName
	existingCert.UpdatedAt = sdkCtx.BlockTime().Unix()

	if err := k.Keeper.UpdateCertificate(sdkCtx, existingCert); err != nil {
		return nil, err
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProductCertificateUpdated,
			sdk.NewAttribute(types.AttributeKeyCertificateID, msg.CertificateId),
			sdk.NewAttribute(types.AttributeKeyProductID, msg.CertificateId),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Updater),
			sdk.NewAttribute(types.AttributeKeyStatus, "active"),
		),
	)

	return &blockchainproto.MsgUpdateCertificateResponse{}, nil
}

// RevokeCertificate handles certificate revocation
func (k msgServer) RevokeCertificate(ctx context.Context, msg *blockchainproto.MsgRevokeCertificate) (*blockchainproto.MsgRevokeCertificateResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if err := k.Keeper.DeleteCertificate(sdkCtx, msg.CertificateId); err != nil {
		return nil, err
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProductCertificateDeleted,
			sdk.NewAttribute(types.AttributeKeyCertificateID, msg.CertificateId),
		),
	)

	return &blockchainproto.MsgRevokeCertificateResponse{}, nil
}

// VerifyCertificate handles certificate verification
func (k msgServer) VerifyCertificate(ctx context.Context, msg *blockchainproto.MsgVerifyCertificate) (*blockchainproto.MsgVerifyCertificateResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if err := k.Keeper.VerifyCertificate(sdkCtx, msg.CertificateId, msg.Verifier, "verified", msg.VerificationProof); err != nil {
		return nil, err
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProductCertificateVerified,
			sdk.NewAttribute(types.AttributeKeyCertificateID, msg.CertificateId),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Verifier),
			sdk.NewAttribute(types.AttributeKeyStatus, "verified"),
		),
	)

	return &blockchainproto.MsgVerifyCertificateResponse{}, nil
}

// TransferCertificate handles certificate transfer
func (k msgServer) TransferCertificate(ctx context.Context, msg *blockchainproto.MsgTransferCertificate) (*blockchainproto.MsgTransferCertificateResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get existing certificate
	existingCert, exists := k.Keeper.GetCertificate(sdkCtx, msg.CertificateId)
	if !exists {
		return nil, fmt.Errorf("certificate with ID %s does not exist", msg.CertificateId)
	}

	// Update certificate owner
	existingCert.Owner = msg.Recipient
	existingCert.UpdatedAt = sdkCtx.BlockTime().Unix()

	if err := k.Keeper.UpdateCertificate(sdkCtx, existingCert); err != nil {
		return nil, err
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"product_certificate_transferred",
			sdk.NewAttribute(types.AttributeKeyCertificateID, msg.CertificateId),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Recipient),
		),
	)

	return &blockchainproto.MsgTransferCertificateResponse{}, nil
}
