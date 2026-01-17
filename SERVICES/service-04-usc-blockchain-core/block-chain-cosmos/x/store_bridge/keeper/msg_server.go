package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/store_bridge/v1/usc/store_bridge/v1"
	blocktypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/block/types"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/store_bridge/types"
)

// MsgServer defines the gRPC message server for the store_bridge module
type MsgServer interface {
	CreateBridge(context.Context, *blockchainproto.MsgCreateBridge) (*blockchainproto.MsgCreateBridgeResponse, error)
	TransferAssets(context.Context, *blockchainproto.MsgTransferAssets) (*blockchainproto.MsgTransferAssetsResponse, error)
	FinalizeBridge(context.Context, *blockchainproto.MsgFinalizeBridge) (*blockchainproto.MsgFinalizeBridgeResponse, error)
	CancelBridge(context.Context, *blockchainproto.MsgCancelBridge) (*blockchainproto.MsgCancelBridgeResponse, error)
}

// msgServer implements MsgServer
type msgServer struct {
	Keeper
}

// NewMsgServerImpl creates a new store_bridge message server
func NewMsgServerImpl(keeper Keeper) MsgServer {
	return &msgServer{Keeper: keeper}
}

// CreateBridge handles bridge creation messages
func (k msgServer) CreateBridge(ctx context.Context, msg *blockchainproto.MsgCreateBridge) (*blockchainproto.MsgCreateBridgeResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Create bridge from proto message
	bridgeID := fmt.Sprintf("bridge-%s-%s", msg.SourceChain.ChainId, msg.TargetChain.ChainId)
	bridge := types.Bridge{
		ID:          bridgeID,
		Name:        fmt.Sprintf("Bridge %s -> %s", msg.SourceChain.ChainName, msg.TargetChain.ChainName),
		Description: fmt.Sprintf("Bridge between %s and %s", msg.SourceChain.ChainName, msg.TargetChain.ChainName),
		FromChain:   msg.SourceChain.ChainId,
		ToChain:     msg.TargetChain.ChainId,
		Type:        msg.BridgeType.String(),
		Status:      "active",
		Config: map[string]string{
			"confirmation_blocks": fmt.Sprintf("%d", msg.BridgeConfig.ConfirmationBlocks),
			"timeout_blocks":      fmt.Sprintf("%d", msg.BridgeConfig.TimeoutBlocks),
			"auto_finalization":   fmt.Sprintf("%t", msg.BridgeConfig.EnableAutoFinalization),
			"auto_retry":          fmt.Sprintf("%t", msg.BridgeConfig.EnableAutoRetry),
			"max_retry_attempts":  fmt.Sprintf("%d", msg.BridgeConfig.MaxRetryAttempts),
		},
		Validators: []string{},
		Threshold:  1,
		CreatedAt:  sdkCtx.BlockTime(),
		UpdatedAt:  sdkCtx.BlockTime(),
		Tags:       msg.BridgeConfig.CustomSettings,
	}

	// Validate the bridge
	if err := bridge.Validate(); err != nil {
		return nil, fmt.Errorf("invalid bridge: %w", err)
	}

	// Set the bridge
	if err := k.SetBridge(sdkCtx, bridge); err != nil {
		return nil, fmt.Errorf("failed to create bridge: %w", err)
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBridgeCreated,
			sdk.NewAttribute(types.AttributeKeyBridgeID, bridgeID),
			sdk.NewAttribute(types.AttributeKeyBridgeName, bridge.Name),
			sdk.NewAttribute(types.AttributeKeyFromChain, msg.SourceChain.ChainId),
			sdk.NewAttribute(types.AttributeKeyToChain, msg.TargetChain.ChainId),
		),
	)

	// Calculate real transaction hash
	creator := msg.SourceChain.ChainId // Use source chain as creator
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, creator, msg.TargetChain.ChainId, "", "create_bridge", bridgeID, "")
	creationHash := blocktypes.CalculateHashFromString(fmt.Sprintf("create-%s", bridgeID))

	return &blockchainproto.MsgCreateBridgeResponse{
		Success:         true,
		BridgeId:        bridgeID,
		CreationHash:    creationHash,
		TransactionHash: txHash,
	}, nil
}

// TransferAssets handles asset transfer messages
func (k msgServer) TransferAssets(ctx context.Context, msg *blockchainproto.MsgTransferAssets) (*blockchainproto.MsgTransferAssetsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get existing bridge
	_, err := k.GetBridge(sdkCtx, msg.BridgeId)
	if err != nil {
		return nil, fmt.Errorf("bridge not found: %w", err)
	}

	// Create transfer from proto message
	transferID := fmt.Sprintf("transfer-%s-%d", msg.BridgeId, sdkCtx.BlockHeight())

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Sender, msg.Recipient, msg.Amount.String(), "transfer_assets", transferID, "")

	transfer := types.Transfer{
		ID:          transferID,
		BridgeID:    msg.BridgeId,
		FromChain:   msg.SourceChain.ChainId,
		ToChain:     msg.TargetChain.ChainId,
		FromAddress: msg.Sender,
		ToAddress:   msg.Recipient,
		Amount:      msg.Amount.String(),
		Token:       msg.Amount.Denom,
		Status:      "pending",
		TxHash:      txHash,
		BlockHeight: sdkCtx.BlockHeight(),
		CreatedAt:   sdkCtx.BlockTime(),
		Metadata:    msg.TransferData,
		Tags:        map[string]string{},
	}

	// Validate the transfer
	if err := transfer.Validate(); err != nil {
		return nil, fmt.Errorf("invalid transfer: %w", err)
	}

	// Set the transfer
	if err := k.SetTransfer(sdkCtx, transfer); err != nil {
		return nil, fmt.Errorf("failed to initiate transfer: %w", err)
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTransferInitiated,
			sdk.NewAttribute(types.AttributeKeyTransferID, transferID),
			sdk.NewAttribute(types.AttributeKeyFromChain, msg.SourceChain.ChainId),
			sdk.NewAttribute(types.AttributeKeyToChain, msg.TargetChain.ChainId),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
		),
	)

	// Use the transaction hash already calculated above
	transferHash := blocktypes.CalculateHashFromString(fmt.Sprintf("transfer-%s", transferID))

	return &blockchainproto.MsgTransferAssetsResponse{
		Success:         true,
		TransferId:      transferID,
		TransferHash:    transferHash,
		TransactionHash: txHash,
	}, nil
}

// FinalizeBridge handles bridge finalization messages
func (k msgServer) FinalizeBridge(ctx context.Context, msg *blockchainproto.MsgFinalizeBridge) (*blockchainproto.MsgFinalizeBridgeResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get existing transfer
	transfer, err := k.GetTransfer(sdkCtx, msg.OperationId)
	if err != nil {
		return nil, fmt.Errorf("transfer not found: %w", err)
	}

	// Update transfer status
	transfer.Status = "completed"
	transfer.CompletedAt = sdkCtx.BlockTime()
	transfer.TxHash = msg.FinalizationProof

	// Set the updated transfer
	if err := k.SetTransfer(sdkCtx, transfer); err != nil {
		return nil, fmt.Errorf("failed to finalize transfer: %w", err)
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTransferCompleted,
			sdk.NewAttribute(types.AttributeKeyTransferID, msg.OperationId),
			sdk.NewAttribute(types.AttributeKeyStatus, "completed"),
		),
	)

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, transfer.FromAddress, transfer.ToAddress, transfer.Amount, "finalize_bridge", msg.OperationId, "")
	finalizationHash := blocktypes.CalculateHashFromString(fmt.Sprintf("finalize-%s", msg.OperationId))

	return &blockchainproto.MsgFinalizeBridgeResponse{
		Success:          true,
		FinalizationHash: finalizationHash,
		TransactionHash:  txHash,
	}, nil
}

// CancelBridge handles bridge cancellation messages
func (k msgServer) CancelBridge(ctx context.Context, msg *blockchainproto.MsgCancelBridge) (*blockchainproto.MsgCancelBridgeResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get existing transfer
	transfer, err := k.GetTransfer(sdkCtx, msg.OperationId)
	if err != nil {
		return nil, fmt.Errorf("transfer not found: %w", err)
	}

	// Update transfer status
	transfer.Status = "failed"
	transfer.FailedAt = sdkCtx.BlockTime()
	transfer.FailureReason = msg.CancelReason

	// Set the updated transfer
	if err := k.SetTransfer(sdkCtx, transfer); err != nil {
		return nil, fmt.Errorf("failed to cancel transfer: %w", err)
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTransferFailed,
			sdk.NewAttribute(types.AttributeKeyTransferID, msg.OperationId),
			sdk.NewAttribute(types.AttributeKeyStatus, "cancelled"),
		),
	)

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, transfer.FromAddress, transfer.ToAddress, transfer.Amount, "cancel_bridge", msg.OperationId, "")
	cancellationHash := blocktypes.CalculateHashFromString(fmt.Sprintf("cancel-%s", msg.OperationId))

	return &blockchainproto.MsgCancelBridgeResponse{
		Success:          true,
		CancellationHash: cancellationHash,
		TransactionHash:  txHash,
	}, nil
}
