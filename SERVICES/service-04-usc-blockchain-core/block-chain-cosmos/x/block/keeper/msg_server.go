package keeper

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/block/v1/usc/block/v1"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/block/types"
)

// MsgServer handles block module messages
type MsgServer struct {
	Keeper
}

// NewMsgServer creates a new block message server
func NewMsgServer(keeper Keeper) *MsgServer {
	return &MsgServer{Keeper: keeper}
}

// CreateBlock handles block creation using blockchain-proto message types
func (k MsgServer) CreateBlock(ctx context.Context, req *blockchainproto.MsgCreateBlock) (*blockchainproto.MsgCreateBlockResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Create block using blockchain-proto types
	block := types.Block{
		ID:           req.Creator,           // Using creator as block ID
		Height:       req.Timestamp.Seconds, // Using timestamp seconds as height for now
		Hash:         req.DataHash,
		PreviousHash: req.PreviousHash,
		Timestamp:    time.Unix(req.Timestamp.Seconds, int64(req.Timestamp.Nanos)),
		Validator:    req.Creator,
		Size:         0, // Will be calculated
		TxCount:      0, // Will be calculated
		GasUsed:      0, // Will be calculated
		GasLimit:     0, // Will be calculated
		Status:       "pending",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Store block
	if err := k.SetBlock(sdkCtx, block); err != nil {
		return nil, fmt.Errorf("failed to create block: %w", err)
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBlockCreated,
			sdk.NewAttribute(types.AttributeKeyBlockID, block.ID),
			sdk.NewAttribute(types.AttributeKeyBlockHeight, fmt.Sprintf("%d", block.Height)),
			sdk.NewAttribute(types.AttributeKeyBlockHash, block.Hash),
			sdk.NewAttribute(types.AttributeKeyBlockValidator, block.Validator),
		),
	)

	return &blockchainproto.MsgCreateBlockResponse{
		Success:         true,
		BlockHash:       block.Hash,
		BlockHeight:     block.Height,
		TransactionHash: block.ID, // Using block ID as transaction hash
	}, nil
}

// Note: UpdateBlock and DeleteBlock handlers removed as they are not part of blockchain-proto interface
// The blockchain-proto interface only includes CreateBlock, ValidateBlock, and FinalizeBlock

// ValidateBlock handles block validation using blockchain-proto message types
func (k MsgServer) ValidateBlock(ctx context.Context, req *blockchainproto.MsgValidateBlock) (*blockchainproto.MsgValidateBlockResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get block by hash
	block, err := k.GetBlockByHash(sdkCtx, req.BlockHash)
	if err != nil {
		return nil, fmt.Errorf("block not found: %w", err)
	}

	// Get module parameters
	params := k.GetParams(sdkCtx)

	// Initialize validation result
	validation := types.BlockValidation{
		BlockID:        block.ID,
		Height:         block.Height,
		Hash:           block.Hash,
		IsValid:        false,
		ValidationTime: time.Now(),
		Validator:      req.Validator,
		Errors:         []string{},
		Warnings:       []string{},
	}

	// Perform actual validation if enabled
	if params.ValidationEnabled {
		// Validate block structure
		if err := block.Validate(); err != nil {
			validation.Errors = append(validation.Errors, fmt.Sprintf("block structure validation failed: %s", err.Error()))
		}

		// Validate block using keeper validation
		if !k.Keeper.ValidateBlock(sdkCtx, block) {
			validation.Errors = append(validation.Errors, "block validation failed: keeper validation returned false")
		}

		// Validate block integrity (hash chain, format, etc.)
		if err := k.ValidateBlockIntegrity(sdkCtx, block.Height); err != nil {
			validation.Errors = append(validation.Errors, fmt.Sprintf("block integrity validation failed: %s", err.Error()))
		}

		// Check block size limits
		if block.Size > params.MaxBlockSize {
			validation.Errors = append(validation.Errors, fmt.Sprintf("block size exceeds maximum: %d > %d", block.Size, params.MaxBlockSize))
		}

		// Check transaction count limits
		if block.TxCount > params.MaxTxCount {
			validation.Errors = append(validation.Errors, fmt.Sprintf("block tx count exceeds maximum: %d > %d", block.TxCount, params.MaxTxCount))
		}

		// Check gas limit
		if block.GasLimit > params.MaxGasLimit {
			validation.Warnings = append(validation.Warnings, fmt.Sprintf("block gas limit exceeds maximum: %d > %d", block.GasLimit, params.MaxGasLimit))
		}

		// Check gas usage
		if block.GasUsed > block.GasLimit {
			validation.Errors = append(validation.Errors, fmt.Sprintf("gas used exceeds gas limit: %d > %d", block.GasUsed, block.GasLimit))
		}

		// Block is valid if no errors
		validation.IsValid = len(validation.Errors) == 0
	} else {
		// Validation disabled - mark as valid but add warning
		validation.IsValid = true
		validation.Warnings = append(validation.Warnings, "block validation is disabled in module parameters")
	}

	// Store validation
	if err := k.SetValidation(sdkCtx, validation); err != nil {
		return nil, fmt.Errorf("failed to store validation: %w", err)
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBlockValidated,
			sdk.NewAttribute(types.AttributeKeyBlockID, block.ID),
			sdk.NewAttribute(types.AttributeKeyBlockHeight, fmt.Sprintf("%d", block.Height)),
			sdk.NewAttribute(types.AttributeKeyBlockValidator, req.Validator),
			sdk.NewAttribute("is_valid", fmt.Sprintf("%t", validation.IsValid)),
			sdk.NewAttribute("error_count", fmt.Sprintf("%d", len(validation.Errors))),
		),
	)

	// Generate validation result message
	validationResult := "valid"
	if !validation.IsValid {
		validationResult = fmt.Sprintf("invalid: %d error(s)", len(validation.Errors))
		if len(validation.Warnings) > 0 {
			validationResult += fmt.Sprintf(", %d warning(s)", len(validation.Warnings))
		}
	} else if len(validation.Warnings) > 0 {
		validationResult = fmt.Sprintf("valid with %d warning(s)", len(validation.Warnings))
	}

	return &blockchainproto.MsgValidateBlockResponse{
		Success:          true,
		IsValid:          validation.IsValid,
		ValidationResult: validationResult,
		TransactionHash:  block.ID,
	}, nil
}

// FinalizeBlock handles block finalization using blockchain-proto message types
func (k MsgServer) FinalizeBlock(ctx context.Context, req *blockchainproto.MsgFinalizeBlock) (*blockchainproto.MsgFinalizeBlockResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get block by hash
	block, err := k.GetBlockByHash(sdkCtx, req.BlockHash)
	if err != nil {
		return nil, fmt.Errorf("block not found: %w", err)
	}

	// Update block status to finalized
	block.Status = "finalized"
	block.UpdatedAt = time.Now()

	// Store updated block
	if err := k.SetBlock(sdkCtx, block); err != nil {
		return nil, fmt.Errorf("failed to finalize block: %w", err)
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBlockFinalized,
			sdk.NewAttribute(types.AttributeKeyBlockID, block.ID),
			sdk.NewAttribute(types.AttributeKeyBlockHeight, fmt.Sprintf("%d", block.Height)),
			sdk.NewAttribute(types.AttributeKeyBlockHash, block.Hash),
			sdk.NewAttribute(types.AttributeKeyBlockFinalizer, req.Finalizer),
		),
	)

	return &blockchainproto.MsgFinalizeBlockResponse{
		Success:            true,
		FinalizationHash:   block.Hash,
		FinalizationHeight: block.Height,
		TransactionHash:    block.ID,
	}, nil
}
