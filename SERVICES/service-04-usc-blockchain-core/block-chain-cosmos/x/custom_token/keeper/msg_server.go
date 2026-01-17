package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/custom_token/v1/usc/custom_token/v1"

	blocktypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/block/types"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/custom_token/types"
)

// MsgServer defines the interface for the custom_token module's message server
type MsgServer interface {
	CreateToken(context.Context, *blockchainproto.MsgCreateToken) (*blockchainproto.MsgCreateTokenResponse, error)
	UpdateToken(context.Context, *blockchainproto.MsgUpdateToken) (*blockchainproto.MsgUpdateTokenResponse, error)
	MintToken(context.Context, *blockchainproto.MsgMintToken) (*blockchainproto.MsgMintTokenResponse, error)
	BurnToken(context.Context, *blockchainproto.MsgBurnToken) (*blockchainproto.MsgBurnTokenResponse, error)
	TransferToken(context.Context, *blockchainproto.MsgTransferToken) (*blockchainproto.MsgTransferTokenResponse, error)
}

// msgServer implements the MsgServer interface
type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
func NewMsgServerImpl(keeper Keeper) MsgServer {
	return &msgServer{Keeper: keeper}
}

// CreateToken handles token creation
func (k msgServer) CreateToken(ctx context.Context, msg *blockchainproto.MsgCreateToken) (*blockchainproto.MsgCreateTokenResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if msg.Creator == "" {
		return nil, fmt.Errorf("creator cannot be empty")
	}
	if msg.TokenName == "" || msg.TokenSymbol == "" {
		return nil, fmt.Errorf("token name/symbol cannot be empty")
	}

	// Create internal token model (mapping from proto)
	token := types.CustomToken{
		ID:          fmt.Sprintf("token_%d", sdkCtx.BlockHeight()),
		Name:        msg.TokenName,
		Symbol:      msg.TokenSymbol,
		Decimals:    uint8(msg.Decimals),
		TotalSupply: msg.InitialSupply.String(),
		Owner:       msg.Creator,
		Status:      "active",
		Metadata:    msg.TokenMetadata.String(),
		CreatedAt:   sdkCtx.BlockTime().Unix(),
		UpdatedAt:   sdkCtx.BlockTime().Unix(),
	}

	if err := k.Keeper.CreateToken(sdkCtx, token); err != nil {
		return nil, err
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTokenCreated,
			sdk.NewAttribute(types.AttributeKeyTokenID, token.ID),
			sdk.NewAttribute(types.AttributeKeyTokenName, token.Name),
			sdk.NewAttribute(types.AttributeKeyTokenSymbol, token.Symbol),
			sdk.NewAttribute(types.AttributeKeyOwner, token.Owner),
		),
	)

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Creator, "", msg.InitialSupply.String(), "create_token", msg.TokenMetadata.String(), "")

	return &blockchainproto.MsgCreateTokenResponse{
		Success:         true,
		TokenId:         token.ID,
		TokenAddress:    "",
		CreationHash:    blocktypes.CalculateHashFromString(fmt.Sprintf("create_%s", token.ID)),
		TransactionHash: txHash,
	}, nil
}

// UpdateToken handles token updates
func (k msgServer) UpdateToken(ctx context.Context, msg *blockchainproto.MsgUpdateToken) (*blockchainproto.MsgUpdateTokenResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get existing token
	existingToken, exists := k.Keeper.GetToken(sdkCtx, msg.TokenId)
	if !exists {
		return nil, fmt.Errorf("token with ID %s does not exist", msg.TokenId)
	}

	// Update token fields
	if msg.NewName != "" {
		existingToken.Name = msg.NewName
	}
	if msg.NewMetadata != nil {
		existingToken.Metadata = msg.NewMetadata.String()
	}
	if msg.NewMaxSupply != nil {
		existingToken.TotalSupply = msg.NewMaxSupply.String()
	}
	existingToken.UpdatedAt = sdkCtx.BlockTime().Unix()

	if err := k.Keeper.UpdateToken(sdkCtx, existingToken); err != nil {
		return nil, err
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTokenUpdated,
			sdk.NewAttribute(types.AttributeKeyTokenID, msg.TokenId),
			sdk.NewAttribute(types.AttributeKeyTokenName, existingToken.Name),
			sdk.NewAttribute(types.AttributeKeyTokenSymbol, existingToken.Symbol),
			sdk.NewAttribute(types.AttributeKeyOwner, existingToken.Owner),
		),
	)

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, existingToken.Owner, "", "", "update_token", existingToken.Metadata, "")

	return &blockchainproto.MsgUpdateTokenResponse{
		Success:         true,
		UpdateHash:      blocktypes.CalculateHashFromString(fmt.Sprintf("update_%s", existingToken.ID)),
		TransactionHash: txHash,
	}, nil
}

// MintToken handles token minting
func (k msgServer) MintToken(ctx context.Context, msg *blockchainproto.MsgMintToken) (*blockchainproto.MsgMintTokenResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if err := k.Keeper.MintToken(sdkCtx, msg.TokenId, msg.Recipient, msg.MintAmount.String()); err != nil {
		return nil, err
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTokenMinted,
			sdk.NewAttribute(types.AttributeKeyTokenID, msg.TokenId),
			sdk.NewAttribute(types.AttributeKeyTo, msg.Recipient),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.MintAmount.String()),
		),
	)

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Recipient, "", msg.MintAmount.String(), "mint_token", msg.TokenId, "")

	return &blockchainproto.MsgMintTokenResponse{
		Success:         true,
		MintingHash:     blocktypes.CalculateHashFromString(fmt.Sprintf("mint_%s", msg.TokenId)),
		TransactionHash: txHash,
	}, nil
}

// BurnToken handles token burning
func (k msgServer) BurnToken(ctx context.Context, msg *blockchainproto.MsgBurnToken) (*blockchainproto.MsgBurnTokenResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if err := k.Keeper.BurnToken(sdkCtx, msg.TokenId, msg.Burner, msg.BurnAmount.String()); err != nil {
		return nil, err
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTokenBurned,
			sdk.NewAttribute(types.AttributeKeyTokenID, msg.TokenId),
			sdk.NewAttribute(types.AttributeKeyFrom, msg.Burner),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.BurnAmount.String()),
		),
	)

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Burner, "", msg.BurnAmount.String(), "burn_token", msg.TokenId, "")

	return &blockchainproto.MsgBurnTokenResponse{
		Success:         true,
		BurnHash:        blocktypes.CalculateHashFromString(fmt.Sprintf("burn_%s", msg.TokenId)),
		TransactionHash: txHash,
	}, nil
}

// TransferToken handles token transfers
func (k msgServer) TransferToken(ctx context.Context, msg *blockchainproto.MsgTransferToken) (*blockchainproto.MsgTransferTokenResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if err := k.Keeper.TransferToken(sdkCtx, msg.TokenId, msg.Sender, msg.Recipient, msg.TransferAmount.String()); err != nil {
		return nil, err
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTokenTransferred,
			sdk.NewAttribute(types.AttributeKeyTokenID, msg.TokenId),
			sdk.NewAttribute(types.AttributeKeyFrom, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyTo, msg.Recipient),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.TransferAmount.String()),
		),
	)

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Sender, msg.Recipient, msg.TransferAmount.String(), "transfer_token", msg.TokenId, "")

	return &blockchainproto.MsgTransferTokenResponse{
		Success:         true,
		TransferHash:    blocktypes.CalculateHashFromString(fmt.Sprintf("transfer_%s", msg.TokenId)),
		TransactionHash: txHash,
	}, nil
}

// Note: legacy custom message structs and validators removed in favor of blockchain-proto types
