package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/transaction/v1/usc/transaction/v1"
	blocktypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/block/types"
)

// MsgServer defines the message server interface using blockchain-proto types
type MsgServer interface {
	CreateTransaction(context.Context, *blockchainproto.MsgCreateTransaction) (*blockchainproto.MsgCreateTransactionResponse, error)
	ValidateTransaction(context.Context, *blockchainproto.MsgValidateTransaction) (*blockchainproto.MsgValidateTransactionResponse, error)
	ExecuteTransaction(context.Context, *blockchainproto.MsgExecuteTransaction) (*blockchainproto.MsgExecuteTransactionResponse, error)
	CancelTransaction(context.Context, *blockchainproto.MsgCancelTransaction) (*blockchainproto.MsgCancelTransactionResponse, error)
}

// msgServer implements the MsgServer interface
type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
func NewMsgServerImpl(keeper Keeper) MsgServer {
	return &msgServer{Keeper: keeper}
}

// CreateTransaction handles transaction creation messages using blockchain-proto types
func (k msgServer) CreateTransaction(ctx context.Context, msg *blockchainproto.MsgCreateTransaction) (*blockchainproto.MsgCreateTransactionResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Basic validation
	if msg.Creator == "" {
		return nil, fmt.Errorf("creator cannot be empty")
	}
	if msg.FromAddress == "" {
		return nil, fmt.Errorf("from address cannot be empty")
	}
	if msg.ToAddress == "" {
		return nil, fmt.Errorf("to address cannot be empty")
	}

	// Generate real transaction hash using blocktypes helper
	hash := blocktypes.CalculateTransactionHash(sdkCtx, msg.FromAddress, msg.ToAddress, msg.Amount.String(), msg.TransactionType.String(), msg.Data.String(), msg.Memo)

	// Create transaction using keeper method
	transaction, err := k.Keeper.CreateTransaction(sdkCtx, hash, hash, msg.FromAddress, msg.ToAddress, msg.Amount.String(), msg.TransactionType.String(), msg.Data.String(), msg.Memo)
	if err != nil {
		return nil, err
	}

	return &blockchainproto.MsgCreateTransactionResponse{
		Success:         true,
		TransactionHash: transaction.Hash,
		TransactionId:   transaction.ID,
	}, nil
}

// ValidateTransaction handles transaction validation messages using blockchain-proto types
func (k msgServer) ValidateTransaction(ctx context.Context, msg *blockchainproto.MsgValidateTransaction) (*blockchainproto.MsgValidateTransactionResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Basic validation
	if msg.Validator == "" {
		return nil, fmt.Errorf("validator cannot be empty")
	}
	if msg.TransactionHash == "" {
		return nil, fmt.Errorf("transaction hash cannot be empty")
	}

	// Validate transaction
	err := k.Keeper.ValidateTransaction(sdkCtx, msg.TransactionHash, msg.Validator, msg.ValidationProof)
	if err != nil {
		return &blockchainproto.MsgValidateTransactionResponse{
			Success:          false,
			IsValid:          false,
			ValidationResult: err.Error(),
		}, nil
	}

	return &blockchainproto.MsgValidateTransactionResponse{
		Success:          true,
		IsValid:          true,
		ValidationResult: "transaction validated successfully",
	}, nil
}

// ExecuteTransaction handles transaction execution messages using blockchain-proto types
func (k msgServer) ExecuteTransaction(ctx context.Context, msg *blockchainproto.MsgExecuteTransaction) (*blockchainproto.MsgExecuteTransactionResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Basic validation
	if msg.Executor == "" {
		return nil, fmt.Errorf("executor cannot be empty")
	}
	if msg.TransactionHash == "" {
		return nil, fmt.Errorf("transaction hash cannot be empty")
	}

	// Execute transaction
	err := k.Keeper.ExecuteTransaction(sdkCtx, msg.TransactionHash, msg.Executor, msg.ExecutionProof)
	if err != nil {
		return &blockchainproto.MsgExecuteTransactionResponse{
			Success:         false,
			ExecutionResult: err.Error(),
		}, nil
	}

	// Generate real execution hash
	executionHash := blocktypes.CalculateHashFromString(fmt.Sprintf("execute_%s_%d", msg.TransactionHash, sdkCtx.BlockHeight()))

	return &blockchainproto.MsgExecuteTransactionResponse{
		Success:         true,
		ExecutionResult: "transaction executed successfully",
		ExecutionHash:   executionHash,
	}, nil
}

// CancelTransaction handles transaction cancellation messages using blockchain-proto types
func (k msgServer) CancelTransaction(ctx context.Context, msg *blockchainproto.MsgCancelTransaction) (*blockchainproto.MsgCancelTransactionResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Basic validation
	if msg.Canceller == "" {
		return nil, fmt.Errorf("canceller cannot be empty")
	}
	if msg.TransactionHash == "" {
		return nil, fmt.Errorf("transaction hash cannot be empty")
	}

	// Cancel transaction
	err := k.Keeper.CancelTransaction(sdkCtx, msg.TransactionHash, msg.Canceller, msg.CancelReason)
	if err != nil {
		return &blockchainproto.MsgCancelTransactionResponse{
			Success: false,
		}, err
	}

	return &blockchainproto.MsgCancelTransactionResponse{
		Success: true,
	}, nil
}

// Note: Custom message types removed as they are replaced by blockchain-proto message types
// The blockchain-proto interface provides MsgCreateTransaction, MsgValidateTransaction, MsgExecuteTransaction, and MsgCancelTransaction
