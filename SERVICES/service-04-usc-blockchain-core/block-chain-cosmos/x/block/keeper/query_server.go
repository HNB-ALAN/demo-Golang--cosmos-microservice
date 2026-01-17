package keeper

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	query "github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/protobuf/types/known/timestamppb"

	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/block/v1/usc/block/v1"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/block/types"
)

// QueryServer handles block module queries
type QueryServer struct {
	Keeper
}

// NewQueryServer creates a new block query server
func NewQueryServer(keeper Keeper) *QueryServer {
	return &QueryServer{Keeper: keeper}
}

// QueryBlock handles block queries using blockchain-proto query types
func (k QueryServer) QueryBlock(ctx context.Context, req *blockchainproto.QueryBlockRequest) (*blockchainproto.QueryBlockResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	var block types.Block
	var err error

	// Query by hash if provided
	if req.BlockHash != "" {
		block, err = k.Keeper.GetBlockByHash(sdkCtx, req.BlockHash)
	} else if req.BlockHeight > 0 {
		// Query by height if provided
		block, err = k.Keeper.GetBlockByHeight(sdkCtx, req.BlockHeight)
	} else {
		return nil, fmt.Errorf("either block_hash or block_height must be provided")
	}

	if err != nil {
		return nil, fmt.Errorf("block not found: %w", err)
	}

	// Convert string status to blockchain-proto enum
	var blockStatus blockchainproto.BlockStatus
	switch block.Status {
	case "pending":
		blockStatus = blockchainproto.BlockStatus_BLOCK_STATUS_PENDING
	case "validated":
		blockStatus = blockchainproto.BlockStatus_BLOCK_STATUS_VALIDATED
	case "finalized":
		blockStatus = blockchainproto.BlockStatus_BLOCK_STATUS_FINALIZED
	default:
		blockStatus = blockchainproto.BlockStatus_BLOCK_STATUS_UNSPECIFIED
	}

	// Convert to blockchain-proto Block type
	blockchainBlock := &blockchainproto.Block{
		Hash:             block.Hash,
		Height:           block.Height,
		PreviousHash:     block.PreviousHash,
		MerkleRoot:       block.Hash, // Using hash as merkle root for now
		Timestamp:        timestamppb.New(block.Timestamp),
		Creator:          block.Validator,
		Validator:        block.Validator,
		Finalizer:        block.Validator,
		Status:           blockStatus,
		DataHash:         block.Hash,
		TransactionCount: block.TxCount,
		GasUsed:          block.GasUsed,
		GasLimit:         block.GasLimit,
		Memo:             "",
		CreatedAt:        timestamppb.New(block.CreatedAt),
		ValidatedAt:      timestamppb.New(block.UpdatedAt),
		FinalizedAt:      timestamppb.New(block.UpdatedAt),
	}

	return &blockchainproto.QueryBlockResponse{
		Block: blockchainBlock,
	}, nil
}

// QueryBlocks handles multiple block queries using blockchain-proto query types
func (k QueryServer) QueryBlocks(ctx context.Context, req *blockchainproto.QueryBlocksRequest) (*blockchainproto.QueryBlocksResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get all blocks
	allBlocks := k.Keeper.GetAllBlocks(sdkCtx)

	// Apply filters
	var filteredBlocks []types.Block
	for _, block := range allBlocks {
		// Filter by height range
		if req.StartHeight > 0 && block.Height < req.StartHeight {
			continue
		}
		if req.EndHeight > 0 && block.Height > req.EndHeight {
			continue
		}
		// Filter by creator
		if req.Creator != "" && block.Validator != req.Creator {
			continue
		}
		// Filter by status
		if req.Status != blockchainproto.BlockStatus_BLOCK_STATUS_UNSPECIFIED {
			// Convert string status to enum for comparison
			var blockStatus blockchainproto.BlockStatus
			switch block.Status {
			case "pending":
				blockStatus = blockchainproto.BlockStatus_BLOCK_STATUS_PENDING
			case "validated":
				blockStatus = blockchainproto.BlockStatus_BLOCK_STATUS_VALIDATED
			case "finalized":
				blockStatus = blockchainproto.BlockStatus_BLOCK_STATUS_FINALIZED
			default:
				blockStatus = blockchainproto.BlockStatus_BLOCK_STATUS_UNSPECIFIED
			}
			if blockStatus != req.Status {
				continue
			}
		}
		filteredBlocks = append(filteredBlocks, block)
	}

	// Convert to blockchain-proto Block types
	var blockchainBlocks []*blockchainproto.Block
	for _, block := range filteredBlocks {
		// Convert string status to blockchain-proto enum
		var blockStatus blockchainproto.BlockStatus
		switch block.Status {
		case "pending":
			blockStatus = blockchainproto.BlockStatus_BLOCK_STATUS_PENDING
		case "validated":
			blockStatus = blockchainproto.BlockStatus_BLOCK_STATUS_VALIDATED
		case "finalized":
			blockStatus = blockchainproto.BlockStatus_BLOCK_STATUS_FINALIZED
		default:
			blockStatus = blockchainproto.BlockStatus_BLOCK_STATUS_UNSPECIFIED
		}

		blockchainBlock := &blockchainproto.Block{
			Hash:             block.Hash,
			Height:           block.Height,
			PreviousHash:     block.PreviousHash,
			MerkleRoot:       block.Hash,
			Timestamp:        timestamppb.New(block.Timestamp),
			Creator:          block.Validator,
			Validator:        block.Validator,
			Finalizer:        block.Validator,
			Status:           blockStatus,
			DataHash:         block.Hash,
			TransactionCount: block.TxCount,
			GasUsed:          block.GasUsed,
			GasLimit:         block.GasLimit,
			Memo:             "",
			CreatedAt:        timestamppb.New(block.CreatedAt),
			ValidatedAt:      timestamppb.New(block.UpdatedAt),
			FinalizedAt:      timestamppb.New(block.UpdatedAt),
		}
		blockchainBlocks = append(blockchainBlocks, blockchainBlock)
	}

	// Apply pagination
	var start, end int64
	if req.Pagination != nil {
		if req.Pagination.Offset > 0 {
			start = int64(req.Pagination.Offset)
		}
		if req.Pagination.Limit > 0 {
			end = start + int64(req.Pagination.Limit)
		} else {
			end = int64(len(blockchainBlocks))
		}
	} else {
		start = 0
		end = int64(len(blockchainBlocks))
	}

	if end > int64(len(blockchainBlocks)) {
		end = int64(len(blockchainBlocks))
	}

	if start >= int64(len(blockchainBlocks)) {
		blockchainBlocks = []*blockchainproto.Block{}
	} else {
		blockchainBlocks = blockchainBlocks[start:end]
	}

	// Create pagination response
	var pagination *query.PageResponse
	if req.Pagination != nil {
		// Calculate next key for pagination
		var nextKey []byte
		if end < int64(len(filteredBlocks)) {
			// There are more blocks, set next key to the height of the last block in this page
			if len(blockchainBlocks) > 0 {
				lastBlock := blockchainBlocks[len(blockchainBlocks)-1]
				// Encode height as next key for continuation
				nextKey = []byte(fmt.Sprintf("%d", lastBlock.Height+1))
			}
		}
		
		pagination = &query.PageResponse{
			NextKey: nextKey,
			Total:   uint64(len(filteredBlocks)),
		}
	}

	return &blockchainproto.QueryBlocksResponse{
		Blocks:     blockchainBlocks,
		Pagination: pagination,
	}, nil
}

// QueryBlockStats handles block statistics queries using blockchain-proto query types
func (k QueryServer) QueryBlockStats(ctx context.Context, req *blockchainproto.QueryBlockStatsRequest) (*blockchainproto.QueryBlockStatsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get all blocks
	allBlocks := k.Keeper.GetAllBlocks(sdkCtx)

	// Filter by height range
	var filteredBlocks []types.Block
	for _, block := range allBlocks {
		if req.StartHeight > 0 && block.Height < req.StartHeight {
			continue
		}
		if req.EndHeight > 0 && block.Height > req.EndHeight {
			continue
		}
		filteredBlocks = append(filteredBlocks, block)
	}

	// Calculate statistics
	var totalBlocks, pendingBlocks, validatedBlocks, finalizedBlocks int64
	var totalTransactions, totalGasUsed int64
	var currentHeight int64
	var lastBlockTime string

	for _, block := range filteredBlocks {
		totalBlocks++
		switch block.Status {
		case "pending":
			pendingBlocks++
		case "validated":
			validatedBlocks++
		case "finalized":
			finalizedBlocks++
		}
		totalTransactions += block.TxCount
		totalGasUsed += block.GasUsed
		if block.Height > currentHeight {
			currentHeight = block.Height
			lastBlockTime = block.Timestamp.Format(time.RFC3339)
		}
	}

	// Calculate averages
	var averageGasPerBlock int64
	if totalBlocks > 0 {
		averageGasPerBlock = totalGasUsed / totalBlocks
	}

	// Parse last block time
	var lastBlockTimeProto *timestamppb.Timestamp
	if lastBlockTime != "" {
		if t, err := time.Parse(time.RFC3339, lastBlockTime); err == nil {
			lastBlockTimeProto = timestamppb.New(t)
		}
	}

	stats := &blockchainproto.BlockStats{
		TotalBlocks:        totalBlocks,
		PendingBlocks:      pendingBlocks,
		ValidatedBlocks:    validatedBlocks,
		FinalizedBlocks:    finalizedBlocks,
		TotalTransactions:  totalTransactions,
		TotalGasUsed:       totalGasUsed,
		AverageGasPerBlock: averageGasPerBlock,
		CurrentHeight:      currentHeight,
		LastBlockTime:      lastBlockTimeProto,
	}

	return &blockchainproto.QueryBlockStatsResponse{
		Stats: stats,
	}, nil
}

// Note: Custom query handlers removed as they are replaced by blockchain-proto query handlers
// The blockchain-proto interface provides QueryBlock, QueryBlocks, and QueryBlockStats
