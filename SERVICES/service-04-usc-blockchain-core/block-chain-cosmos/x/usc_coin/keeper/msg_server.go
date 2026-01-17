package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/usc_coin/v1/usc/usc_coin/v1"
	blocktypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/block/types"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/usc_coin/types"
)

// MsgServer defines the message server interface using blockchain-proto types
type MsgServer interface {
	TransferUSC(context.Context, *blockchainproto.MsgTransferUSC) (*blockchainproto.MsgTransferUSCResponse, error)
	MintUSC(context.Context, *blockchainproto.MsgMintUSC) (*blockchainproto.MsgMintUSCResponse, error)
	BurnUSC(context.Context, *blockchainproto.MsgBurnUSC) (*blockchainproto.MsgBurnUSCResponse, error)
}

// msgServer implements the MsgServer interface
type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
func NewMsgServerImpl(keeper Keeper) MsgServer {
	return &msgServer{Keeper: keeper}
}

// TransferUSC handles USC transfer messages using blockchain-proto types
func (k msgServer) TransferUSC(ctx context.Context, msg *blockchainproto.MsgTransferUSC) (*blockchainproto.MsgTransferUSCResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Check if sender and receiver are different
	if msg.FromAddress == msg.ToAddress {
		return nil, fmt.Errorf("cannot transfer to self")
	}

	// Validate addresses
	if !types.IsValidAddress(msg.FromAddress) {
		return nil, fmt.Errorf("invalid sender address")
	}

	if !types.IsValidAddress(msg.ToAddress) {
		return nil, fmt.Errorf("invalid receiver address")
	}

	// Validate amount
	if msg.Amount == nil || msg.Amount.Amount.IsZero() {
		return nil, fmt.Errorf("amount cannot be zero")
	}

	// Perform transfer using keeper method
	_, err := k.TransferUSC(sdkCtx, msg)
	if err != nil {
		return nil, fmt.Errorf("transfer failed: %w", err)
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTransfer,
			sdk.NewAttribute(types.AttributeKeyFrom, msg.FromAddress),
			sdk.NewAttribute(types.AttributeKeyTo, msg.ToAddress),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.Amount.String()),
		),
	)

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.FromAddress, msg.ToAddress, msg.Amount.Amount.String(), "transfer_usc", "", "")

	return &blockchainproto.MsgTransferUSCResponse{
		Success:         true,
		TransactionHash: txHash,
	}, nil
}

// MintUSC handles USC mint messages using blockchain-proto types
func (k msgServer) MintUSC(ctx context.Context, msg *blockchainproto.MsgMintUSC) (*blockchainproto.MsgMintUSCResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Validate address
	if !types.IsValidAddress(msg.Minter) {
		return nil, fmt.Errorf("invalid minter address")
	}

	// Validate amount
	if msg.Amount == nil || msg.Amount.Amount.IsZero() {
		return nil, fmt.Errorf("amount cannot be zero")
	}

	// Perform mint using keeper method
	_, err := k.MintUSC(sdkCtx, msg)
	if err != nil {
		return nil, fmt.Errorf("mint failed: %w", err)
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeyTo, msg.Minter),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.Amount.String()),
		),
	)

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Minter, "", msg.Amount.Amount.String(), "mint_usc", "", "")

	return &blockchainproto.MsgMintUSCResponse{
		Success:         true,
		TransactionHash: txHash,
	}, nil
}

// BurnUSC handles USC burn messages using blockchain-proto types
func (k msgServer) BurnUSC(ctx context.Context, msg *blockchainproto.MsgBurnUSC) (*blockchainproto.MsgBurnUSCResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Validate address
	if !types.IsValidAddress(msg.Burner) {
		return nil, fmt.Errorf("invalid burner address")
	}

	// Validate amount
	if msg.Amount == nil || msg.Amount.Amount.IsZero() {
		return nil, fmt.Errorf("amount cannot be zero")
	}

	// Perform burn using keeper method
	_, err := k.BurnUSC(sdkCtx, msg)
	if err != nil {
		return nil, fmt.Errorf("burn failed: %w", err)
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBurn,
			sdk.NewAttribute(types.AttributeKeyFrom, msg.Burner),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.Amount.String()),
		),
	)

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Burner, "", msg.Amount.Amount.String(), "burn_usc", "", "")

	return &blockchainproto.MsgBurnUSCResponse{
		Success:         true,
		TransactionHash: txHash,
	}, nil
}

// Note: Custom message types removed as they are replaced by blockchain-proto message types
// The blockchain-proto interface provides MsgTransferUSC, MsgMintUSC, and MsgBurnUSC
