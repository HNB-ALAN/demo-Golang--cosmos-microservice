package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	query "github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/protobuf/types/known/timestamppb"

	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/nft_token/v1/usc/nft_token/v1"
)

// QueryServer defines the gRPC querier service for the NFT module
type QueryServer struct {
	Keeper
}

// NewQueryServer creates a new NFT query server
func NewQueryServer(keeper Keeper) *QueryServer {
	return &QueryServer{Keeper: keeper}
}

// QueryNFT handles NFT queries by ID
func (k QueryServer) QueryNFT(ctx context.Context, req *blockchainproto.QueryNFTRequest) (*blockchainproto.QueryNFTResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	nft, err := k.Keeper.GetNFT(sdkCtx, req.NftId)
	if err != nil {
		return nil, fmt.Errorf("NFT not found: %w", err)
	}

	// Convert to blockchain-proto NFT type
	blockchainNFT := &blockchainproto.NFT{
		Id:                nft.ID,
		TokenId:           nft.ID, // Use same ID for token ID
		Name:              nft.Name,
		Description:       nft.Description,
		ImageUri:          nft.Image,
		Metadata:          nil, // TODO: map attributes to NFTMetadata
		Owner:             nft.Owner,
		Creator:           nft.Owner, // Use owner as creator for now
		CollectionId:      nft.CollectionID,
		Status:            blockchainproto.NFTStatus_NFT_STATUS_ACTIVE,
		RoyaltyPercentage: 0,
		TransferFee:       &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)},
		MintedAt:          timestamppb.New(nft.CreatedAt),
		LastTransferredAt: timestamppb.New(nft.UpdatedAt),
		BurnedAt:          nil,
		TransferCount:     0,
		CurrentValue:      &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)},
		Memo:              "",
	}

	return &blockchainproto.QueryNFTResponse{
		Nft: blockchainNFT,
	}, nil
}

// QueryNFTs handles queries for all NFTs
func (k QueryServer) QueryNFTs(ctx context.Context, req *blockchainproto.QueryNFTsRequest) (*blockchainproto.QueryNFTsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	nfts := k.Keeper.GetAllNFTs(sdkCtx)

	// Convert to blockchain-proto NFT types
	var blockchainNFTs []*blockchainproto.NFT
	for _, nft := range nfts {
		blockchainNFT := &blockchainproto.NFT{
			Id:                nft.ID,
			TokenId:           nft.ID, // Use same ID for token ID
			Name:              nft.Name,
			Description:       nft.Description,
			ImageUri:          nft.Image,
			Metadata:          nil,
			Owner:             nft.Owner,
			Creator:           nft.Owner, // Use owner as creator for now
			CollectionId:      nft.CollectionID,
			Status:            blockchainproto.NFTStatus_NFT_STATUS_ACTIVE,
			RoyaltyPercentage: 0,
			TransferFee:       &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)},
			MintedAt:          timestamppb.New(nft.CreatedAt),
			LastTransferredAt: timestamppb.New(nft.UpdatedAt),
			BurnedAt:          nil,
			TransferCount:     0,
			CurrentValue:      &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)},
			Memo:              "",
		}
		blockchainNFTs = append(blockchainNFTs, blockchainNFT)
	}

	// Apply pagination
	pageRes := &query.PageResponse{
		NextKey: nil,
		Total:   uint64(len(blockchainNFTs)),
	}

	return &blockchainproto.QueryNFTsResponse{
		Nfts:       blockchainNFTs,
		Pagination: pageRes,
	}, nil
}

// QueryCollection handles collection queries by ID
func (k QueryServer) QueryCollection(ctx context.Context, req *blockchainproto.QueryCollectionRequest) (*blockchainproto.QueryCollectionResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	collection, err := k.Keeper.GetCollection(sdkCtx, req.CollectionId)
	if err != nil {
		return nil, fmt.Errorf("collection not found: %w", err)
	}

	// Convert to blockchain-proto NFTCollection type
	blockchainCollection := &blockchainproto.NFTCollection{
		Id:          collection.ID,
		Name:        collection.Name,
		Description: collection.Description,
		ImageUri:    collection.Image,
		Metadata:    nil,
		Symbol:      collection.Symbol,
		Creator:     collection.Owner,
		Status:      blockchainproto.CollectionStatus_COLLECTION_STATUS_ACTIVE,
		CreatedAt:   timestamppb.New(collection.CreatedAt),
		UpdatedAt:   timestamppb.New(collection.UpdatedAt),
		TotalNfts:   collection.CurrentSupply,
		ActiveNfts:  collection.CurrentSupply,
		BurnedNfts:  0,
		TotalVolume: &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)},
		Memo:        "",
	}

	return &blockchainproto.QueryCollectionResponse{
		Collection: blockchainCollection,
	}, nil
}

// QueryCollections handles queries for all collections
func (k QueryServer) QueryCollections(ctx context.Context, req *blockchainproto.QueryCollectionsRequest) (*blockchainproto.QueryCollectionsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	collections := k.Keeper.GetAllCollections(sdkCtx)

	// Convert to blockchain-proto NFTCollection types
	var blockchainCollections []*blockchainproto.NFTCollection
	for _, collection := range collections {
		blockchainCollection := &blockchainproto.NFTCollection{
			Id:          collection.ID,
			Name:        collection.Name,
			Description: collection.Description,
			ImageUri:    collection.Image,
			Metadata:    nil,
			Symbol:      collection.Symbol,
			Creator:     collection.Owner,
			Status:      blockchainproto.CollectionStatus_COLLECTION_STATUS_ACTIVE,
			CreatedAt:   timestamppb.New(collection.CreatedAt),
			UpdatedAt:   timestamppb.New(collection.UpdatedAt),
			TotalNfts:   collection.CurrentSupply,
			ActiveNfts:  collection.CurrentSupply,
			BurnedNfts:  0,
			TotalVolume: &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)},
			Memo:        "",
		}
		blockchainCollections = append(blockchainCollections, blockchainCollection)
	}

	// Apply pagination
	pageRes := &query.PageResponse{
		NextKey: nil,
		Total:   uint64(len(blockchainCollections)),
	}

	return &blockchainproto.QueryCollectionsResponse{
		Collections: blockchainCollections,
		Pagination:  pageRes,
	}, nil
}

// QueryNFTStats handles NFT statistics queries
func (k QueryServer) QueryNFTStats(ctx context.Context, req *blockchainproto.QueryNFTStatsRequest) (*blockchainproto.QueryNFTStatsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get all NFTs and collections for statistics
	nfts := k.Keeper.GetAllNFTs(sdkCtx)
	collections := k.Keeper.GetAllCollections(sdkCtx)

	// Calculate statistics
	totalNFTs := int64(len(nfts))
	totalCollections := int64(len(collections))

	// Calculate total supply across all collections
	var totalSupply int64
	for _, collection := range collections {
		totalSupply += collection.CurrentSupply
	}

	// Create stats response
	stats := &blockchainproto.NFTStats{
		TotalNfts:         totalNFTs,
		ActiveNfts:        totalNFTs,
		BurnedNfts:        0,
		TotalCollections:  totalCollections,
		ActiveCollections: totalCollections,
		TotalOwners:       totalNFTs, // Simplified - assume each NFT has unique owner
		TotalVolume:       &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)},
		AverageNftValue:   &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)},
		MostValuableNft:   "",
		LastActivity:      timestamppb.New(sdkCtx.BlockTime()),
	}

	return &blockchainproto.QueryNFTStatsResponse{
		Stats: stats,
	}, nil
}
