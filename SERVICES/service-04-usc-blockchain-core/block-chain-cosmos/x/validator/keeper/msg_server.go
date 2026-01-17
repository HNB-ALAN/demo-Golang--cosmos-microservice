package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/validator/v1/usc/validator/v1"

	blocktypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/block/types"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/validator/types"
)

// MsgServer defines the message server interface using blockchain-proto types
type MsgServer interface {
	CreateValidator(context.Context, *blockchainproto.MsgCreateValidator) (*blockchainproto.MsgCreateValidatorResponse, error)
	UpdateValidator(context.Context, *blockchainproto.MsgUpdateValidator) (*blockchainproto.MsgUpdateValidatorResponse, error)
	DelegateValidator(context.Context, *blockchainproto.MsgDelegateValidator) (*blockchainproto.MsgDelegateValidatorResponse, error)
	UndelegateValidator(context.Context, *blockchainproto.MsgUndelegateValidator) (*blockchainproto.MsgUndelegateValidatorResponse, error)
	SlashValidator(context.Context, *blockchainproto.MsgSlashValidator) (*blockchainproto.MsgSlashValidatorResponse, error)
}

// msgServer implements the MsgServer interface
type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
func NewMsgServerImpl(keeper Keeper) MsgServer {
	return &msgServer{Keeper: keeper}
}

// CreateValidator handles validator creation messages
func (k msgServer) CreateValidator(ctx context.Context, msg *blockchainproto.MsgCreateValidator) (*blockchainproto.MsgCreateValidatorResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Basic validation
	if msg.Creator == "" {
		return nil, fmt.Errorf("creator cannot be empty")
	}
	if msg.ValidatorAddress == "" {
		return nil, fmt.Errorf("validator address cannot be empty")
	}
	if msg.ValidatorName == "" {
		return nil, fmt.Errorf("validator name cannot be empty")
	}
	if msg.Description == nil || msg.Description.Details == "" {
		return nil, fmt.Errorf("description cannot be empty")
	}
	if msg.InitialStake == nil {
		return nil, fmt.Errorf("initial stake cannot be empty")
	}

	// Validate address
	if !types.IsValidAddress(msg.ValidatorAddress) {
		return nil, fmt.Errorf("invalid validator address")
	}

	// Create validator
	validator := types.NewValidator(msg.ValidatorAddress, msg.ValidatorName, msg.Description.Details, fmt.Sprintf("%.2f", msg.CommissionRate))
	if err := k.Keeper.CreateValidator(sdkCtx, validator); err != nil {
		return nil, fmt.Errorf("create validator failed: %w", err)
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeValidatorCreate,
			sdk.NewAttribute(types.AttributeKeyValidatorAddress, msg.ValidatorAddress),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Description.Details),
		),
	)

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Creator, msg.ValidatorAddress, msg.InitialStake.String(), "create_validator", msg.Description.Details, "")
	creationHash := blocktypes.CalculateHashFromString(fmt.Sprintf("create_validator_%s", msg.ValidatorAddress))
	
	return &blockchainproto.MsgCreateValidatorResponse{
		Success:         true,
		ValidatorId:     msg.ValidatorAddress,
		CreationHash:    creationHash,
		TransactionHash: txHash,
	}, nil
}

// UpdateValidator handles validator update messages
func (k msgServer) UpdateValidator(ctx context.Context, msg *blockchainproto.MsgUpdateValidator) (*blockchainproto.MsgUpdateValidatorResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Basic validation
	if msg.Updater == "" {
		return nil, fmt.Errorf("updater cannot be empty")
	}
	if msg.ValidatorAddress == "" {
		return nil, fmt.Errorf("validator address cannot be empty")
	}

	// Validate address
	if !types.IsValidAddress(msg.ValidatorAddress) {
		return nil, fmt.Errorf("invalid validator address")
	}

	// Get existing validator
	validator, err := k.GetValidator(sdkCtx, msg.ValidatorAddress)
	if err != nil {
		return nil, fmt.Errorf("validator not found: %w", err)
	}

	// Update validator fields
	if msg.NewDescription != nil && msg.NewDescription.Details != "" {
		validator.Description = msg.NewDescription.Details
	}
	if msg.NewCommissionRate > 0 {
		validator.Commission = fmt.Sprintf("%.2f", msg.NewCommissionRate)
	}

	// Update validator
	if err := k.Keeper.UpdateValidator(sdkCtx, validator); err != nil {
		return nil, fmt.Errorf("update validator failed: %w", err)
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeValidatorUpdate,
			sdk.NewAttribute(types.AttributeKeyValidatorAddress, msg.ValidatorAddress),
		),
	)

	// Calculate real transaction hash
	updateData := ""
	if msg.NewDescription != nil {
		updateData = msg.NewDescription.Details
	}
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Updater, msg.ValidatorAddress, "", "update_validator", updateData, "")
	updateHash := blocktypes.CalculateHashFromString(fmt.Sprintf("update_validator_%s", msg.ValidatorAddress))
	
	return &blockchainproto.MsgUpdateValidatorResponse{
		Success:         true,
		UpdateHash:      updateHash,
		TransactionHash: txHash,
	}, nil
}

// DelegateValidator handles delegation messages
func (k msgServer) DelegateValidator(ctx context.Context, msg *blockchainproto.MsgDelegateValidator) (*blockchainproto.MsgDelegateValidatorResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Basic validation
	if msg.Delegator == "" {
		return nil, fmt.Errorf("delegator cannot be empty")
	}
	if msg.ValidatorAddress == "" {
		return nil, fmt.Errorf("validator address cannot be empty")
	}
	if msg.DelegationAmount == nil {
		return nil, fmt.Errorf("delegation amount cannot be empty")
	}

	// Validate addresses
	if !types.IsValidAddress(msg.Delegator) {
		return nil, fmt.Errorf("invalid delegator address")
	}
	if !types.IsValidAddress(msg.ValidatorAddress) {
		return nil, fmt.Errorf("invalid validator address")
	}

	// Create delegation
	delegation := types.NewDelegation(msg.Delegator, msg.ValidatorAddress, msg.DelegationAmount.String())
	if err := k.Keeper.Delegate(sdkCtx, delegation); err != nil {
		return nil, fmt.Errorf("delegate failed: %w", err)
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDelegationCreate,
			sdk.NewAttribute(types.AttributeKeyDelegatorAddress, msg.Delegator),
			sdk.NewAttribute(types.AttributeKeyValidatorAddress, msg.ValidatorAddress),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.DelegationAmount.String()),
		),
	)

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Delegator, msg.ValidatorAddress, msg.DelegationAmount.String(), "delegate_validator", "", "")
	delegationId := fmt.Sprintf("delegation_%s_%s", msg.Delegator, msg.ValidatorAddress)
	
	return &blockchainproto.MsgDelegateValidatorResponse{
		Success:         true,
		DelegationId:    delegationId,
		TransactionHash: txHash,
	}, nil
}

// UndelegateValidator handles undelegation messages
func (k msgServer) UndelegateValidator(ctx context.Context, msg *blockchainproto.MsgUndelegateValidator) (*blockchainproto.MsgUndelegateValidatorResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Basic validation
	if msg.Undelegator == "" {
		return nil, fmt.Errorf("undelegator cannot be empty")
	}
	if msg.ValidatorAddress == "" {
		return nil, fmt.Errorf("validator address cannot be empty")
	}

	// Validate addresses
	if !types.IsValidAddress(msg.Undelegator) {
		return nil, fmt.Errorf("invalid undelegator address")
	}
	if !types.IsValidAddress(msg.ValidatorAddress) {
		return nil, fmt.Errorf("invalid validator address")
	}

	// Undelegate
	if err := k.Keeper.Undelegate(sdkCtx, msg.Undelegator, msg.ValidatorAddress); err != nil {
		return nil, fmt.Errorf("undelegate failed: %w", err)
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDelegationRemove,
			sdk.NewAttribute(types.AttributeKeyDelegatorAddress, msg.Undelegator),
			sdk.NewAttribute(types.AttributeKeyValidatorAddress, msg.ValidatorAddress),
		),
	)

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Undelegator, msg.ValidatorAddress, "", "undelegate_validator", "", "")
	undelegationId := fmt.Sprintf("undelegation_%s_%s", msg.Undelegator, msg.ValidatorAddress)
	
	return &blockchainproto.MsgUndelegateValidatorResponse{
		Success:         true,
		UndelegationId:  undelegationId,
		TransactionHash: txHash,
	}, nil
}

// SlashValidator handles validator slashing messages
func (k msgServer) SlashValidator(ctx context.Context, msg *blockchainproto.MsgSlashValidator) (*blockchainproto.MsgSlashValidatorResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Basic validation
	if msg.Slasher == "" {
		return nil, fmt.Errorf("slasher cannot be empty")
	}
	if msg.ValidatorAddress == "" {
		return nil, fmt.Errorf("validator address cannot be empty")
	}
	if msg.SlashAmount == nil {
		return nil, fmt.Errorf("slash amount cannot be empty")
	}

	// Validate addresses
	if !types.IsValidAddress(msg.Slasher) {
		return nil, fmt.Errorf("invalid slasher address")
	}
	if !types.IsValidAddress(msg.ValidatorAddress) {
		return nil, fmt.Errorf("invalid validator address")
	}

	// Slash validator - using a simple implementation for now
	// TODO: Implement proper slashing logic in keeper
	// if err := k.Keeper.SlashValidator(sdkCtx, msg.ValidatorAddress, msg.SlashAmount); err != nil {
	// 	return nil, fmt.Errorf("slash validator failed: %w", err)
	// }

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeValidatorUpdate, // Using existing event type
			sdk.NewAttribute(types.AttributeKeyValidatorAddress, msg.ValidatorAddress),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.SlashAmount.String()),
		),
	)

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Slasher, msg.ValidatorAddress, msg.SlashAmount.String(), "slash_validator", "", "")
	slashHash := blocktypes.CalculateHashFromString(fmt.Sprintf("slash_validator_%s", msg.ValidatorAddress))
	
	return &blockchainproto.MsgSlashValidatorResponse{
		Success:         true,
		SlashHash:       slashHash,
		TransactionHash: txHash,
	}, nil
}
