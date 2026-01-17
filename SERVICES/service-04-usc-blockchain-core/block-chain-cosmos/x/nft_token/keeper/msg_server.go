package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/nft_token/v1/usc/nft_token/v1"
	blocktypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/block/types"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/nft_token/types"
)

// MsgServer defines the gRPC message server for the NFT module
type MsgServer struct {
	Keeper
}

// NewMsgServer creates a new NFT message server
func NewMsgServer(keeper Keeper) *MsgServer {
	return &MsgServer{Keeper: keeper}
}

// MintNFT handles NFT minting messages
func (k MsgServer) MintNFT(ctx context.Context, msg *blockchainproto.MsgMintNFT) (*blockchainproto.MsgMintNFTResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Basic validation
	if msg.Minter == "" {
		return nil, fmt.Errorf("minter cannot be empty")
	}
	if msg.Recipient == "" {
		return nil, fmt.Errorf("recipient cannot be empty")
	}
	if msg.NftName == "" {
		return nil, fmt.Errorf("NFT name cannot be empty")
	}
	if msg.CollectionId == "" {
		return nil, fmt.Errorf("collection ID cannot be empty")
	}

	// Create NFT ID (simplified - using timestamp + minter)
	nftID := fmt.Sprintf("nft_%d_%s", sdkCtx.BlockTime().Unix(), msg.Minter)

	// Create NFT object
	nft := types.NFT{
		ID:           nftID,
		CollectionID: msg.CollectionId,
		Owner:        msg.Recipient,
		TokenURI:     msg.NftImageUri,
		Name:         msg.NftName,
		Description:  msg.NftDescription,
		Image:        msg.NftImageUri,
		Attributes:   make(map[string]string),
		CreatedAt:    sdkCtx.BlockTime(),
		UpdatedAt:    sdkCtx.BlockTime(),
		Metadata:     make(map[string]string),
	}

	// Set the NFT
	if err := k.SetNFT(sdkCtx, nft); err != nil {
		return nil, fmt.Errorf("failed to mint NFT: %w", err)
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeNFTCreated,
			sdk.NewAttribute(types.AttributeKeyNFTID, nft.ID),
			sdk.NewAttribute(types.AttributeKeyCollectionID, nft.CollectionID),
			sdk.NewAttribute(types.AttributeKeyOwner, nft.Owner),
		),
	)

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Minter, msg.Recipient, "", "mint_nft", fmt.Sprintf("%s:%s", msg.NftName, msg.CollectionId), "")

	return &blockchainproto.MsgMintNFTResponse{
		Success:         true,
		NftId:           nft.ID,
		TransactionHash: txHash,
	}, nil
}

// TransferNFT handles NFT transfer messages
func (k MsgServer) TransferNFT(ctx context.Context, msg *blockchainproto.MsgTransferNFT) (*blockchainproto.MsgTransferNFTResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Basic validation
	if msg.NftId == "" {
		return nil, fmt.Errorf("NFT ID cannot be empty")
	}
	if msg.Sender == "" {
		return nil, fmt.Errorf("sender address cannot be empty")
	}
	if msg.Recipient == "" {
		return nil, fmt.Errorf("recipient address cannot be empty")
	}

	// Get existing NFT
	nft, err := k.GetNFT(sdkCtx, msg.NftId)
	if err != nil {
		return nil, fmt.Errorf("NFT not found: %w", err)
	}

	// Update NFT owner
	nft.Owner = msg.Recipient
	nft.UpdatedAt = sdkCtx.BlockTime()

	// Validate updated NFT
	if err := nft.Validate(); err != nil {
		return nil, fmt.Errorf("invalid NFT update: %w", err)
	}

	// Set the updated NFT
	if err := k.SetNFT(sdkCtx, nft); err != nil {
		return nil, fmt.Errorf("failed to transfer NFT: %w", err)
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeNFTTransferred,
			sdk.NewAttribute(types.AttributeKeyNFTID, msg.NftId),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyRecipient, msg.Recipient),
		),
	)

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Sender, msg.Recipient, "", "transfer_nft", msg.NftId, "")

	return &blockchainproto.MsgTransferNFTResponse{
		Success:         true,
		TransactionHash: txHash,
	}, nil
}

// UpdateNFT handles NFT update messages
func (k MsgServer) UpdateNFT(ctx context.Context, msg *blockchainproto.MsgUpdateNFT) (*blockchainproto.MsgUpdateNFTResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Basic validation
	if msg.NftId == "" {
		return nil, fmt.Errorf("NFT ID cannot be empty")
	}

	// Get existing NFT
	nft, err := k.GetNFT(sdkCtx, msg.NftId)
	if err != nil {
		return nil, fmt.Errorf("NFT not found: %w", err)
	}

	// Update NFT fields
	if msg.NewName != "" {
		nft.Name = msg.NewName
	}
	if msg.NewDescription != "" {
		nft.Description = msg.NewDescription
	}
	if msg.NewImageUri != "" {
		nft.Image = msg.NewImageUri
		nft.TokenURI = msg.NewImageUri
	}
	nft.UpdatedAt = sdkCtx.BlockTime()

	// Validate updated NFT
	if err := nft.Validate(); err != nil {
		return nil, fmt.Errorf("invalid NFT update: %w", err)
	}

	// Set the updated NFT
	if err := k.SetNFT(sdkCtx, nft); err != nil {
		return nil, fmt.Errorf("failed to update NFT: %w", err)
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeNFTUpdated,
			sdk.NewAttribute(types.AttributeKeyNFTID, msg.NftId),
			sdk.NewAttribute(types.AttributeKeyOwner, nft.Owner),
		),
	)

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, nft.Owner, "", "", "update_nft", msg.NftId, "")

	return &blockchainproto.MsgUpdateNFTResponse{
		Success:         true,
		TransactionHash: txHash,
	}, nil
}

// BurnNFT handles NFT burn messages
func (k MsgServer) BurnNFT(ctx context.Context, msg *blockchainproto.MsgBurnNFT) (*blockchainproto.MsgBurnNFTResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Basic validation
	if msg.NftId == "" {
		return nil, fmt.Errorf("NFT ID cannot be empty")
	}
	if msg.Burner == "" {
		return nil, fmt.Errorf("burner cannot be empty")
	}

	// Get existing NFT
	nft, err := k.GetNFT(sdkCtx, msg.NftId)
	if err != nil {
		return nil, fmt.Errorf("NFT not found: %w", err)
	}

	// Validate ownership - only owner can burn
	if nft.Owner != msg.Burner {
		return nil, fmt.Errorf("unauthorized: only NFT owner can burn, owner is %s, burner is %s", nft.Owner, msg.Burner)
	}

	// Store collection ID before deletion for collection update
	collectionID := nft.CollectionID

	// Delete NFT from storage (burning)
	if err := k.DeleteNFT(sdkCtx, msg.NftId); err != nil {
		return nil, fmt.Errorf("failed to burn NFT: %w", err)
	}

	// Update collection if it exists (log for tracking)
	if collectionID != "" {
		_, err := k.GetCollection(sdkCtx, collectionID)
		if err == nil {
			// Collection exists - NFT has been removed from it
			// Note: Collection NFT count would need to be tracked separately if needed
			sdkCtx.Logger().Debug("NFT burned from collection", "collection_id", collectionID, "nft_id", msg.NftId)
		}
	}

	// Emit event with detailed burn information
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeNFTBurned,
			sdk.NewAttribute(types.AttributeKeyNFTID, msg.NftId),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Burner),
			sdk.NewAttribute("collection_id", collectionID),
			sdk.NewAttribute("token_uri", nft.TokenURI),
		),
	)

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Burner, "", "", "burn_nft", msg.NftId, "")

	return &blockchainproto.MsgBurnNFTResponse{
		Success:         true,
		TransactionHash: txHash,
	}, nil
}

// CreateCollection handles collection creation messages
func (k MsgServer) CreateCollection(ctx context.Context, msg *blockchainproto.MsgCreateCollection) (*blockchainproto.MsgCreateCollectionResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Basic validation
	if msg.Creator == "" {
		return nil, fmt.Errorf("creator cannot be empty")
	}
	if msg.CollectionName == "" {
		return nil, fmt.Errorf("collection name cannot be empty")
	}
	if msg.CollectionSymbol == "" {
		return nil, fmt.Errorf("collection symbol cannot be empty")
	}

	// Create collection ID (simplified - using timestamp + creator)
	collectionID := fmt.Sprintf("collection_%d_%s", sdkCtx.BlockTime().Unix(), msg.Creator)

	// Create collection object
	collection := types.Collection{
		ID:            collectionID,
		Name:          msg.CollectionName,
		Description:   msg.CollectionDescription,
		Symbol:        msg.CollectionSymbol,
		Image:         msg.CollectionImageUri,
		Owner:         msg.Creator,
		MaxSupply:     10000, // Default max supply
		CurrentSupply: 0,
		CreatedAt:     sdkCtx.BlockTime(),
		UpdatedAt:     sdkCtx.BlockTime(),
		Metadata:      make(map[string]string),
	}

	// Validate the collection
	if err := collection.Validate(); err != nil {
		return nil, fmt.Errorf("invalid collection: %w", err)
	}

	// Set the collection
	if err := k.SetCollection(sdkCtx, collection); err != nil {
		return nil, fmt.Errorf("failed to create collection: %w", err)
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCollectionCreated,
			sdk.NewAttribute(types.AttributeKeyCollectionID, collection.ID),
			sdk.NewAttribute(types.AttributeKeyOwner, collection.Owner),
		),
	)

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Creator, "", "", "create_collection", collection.ID, "")
	creationHash := blocktypes.CalculateHashFromString(fmt.Sprintf("create_%s", collection.ID))

	return &blockchainproto.MsgCreateCollectionResponse{
		Success:         true,
		CollectionId:    collection.ID,
		CreationHash:    creationHash,
		TransactionHash: txHash,
	}, nil
}
