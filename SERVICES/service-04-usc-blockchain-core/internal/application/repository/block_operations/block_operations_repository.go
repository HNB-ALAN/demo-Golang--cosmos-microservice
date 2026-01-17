package block_operations

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	repoerrors "service-04/internal/application/repository"
	"service-04/internal/application/utils"
	"service-04/internal/infrastructure/database"
	proto "service-04/proto"

	// Cosmos SDK imports
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/app"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/storage"
	blocktypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/block/types"
	"google.golang.org/protobuf/types/known/timestamppb"

	sherrors "github.com/usc-platform/shared/errors"
	"github.com/usc-platform/shared/logging"
)

// Repository handles block operations data access
type Repository struct {
	db                *database.PostgreSQLManager
	cosmosApp         *app.USCApp
	blockchainStorage *storage.StateManager
	redisManager      *database.RedisManager
	logger            *logging.Logger
}

// NewRepository creates a new block operations repository
func NewRepository(db *database.PostgreSQLManager, cosmosApp *app.USCApp, blockchainStorage *storage.StateManager, redisManager *database.RedisManager, logger *logging.Logger) *Repository {
	return &Repository{
		db:                db,
		cosmosApp:         cosmosApp,
		blockchainStorage: blockchainStorage,
		redisManager:      redisManager,
		logger:            logger,
	}
}

// ProduceBlock creates a new block
func (r *Repository) ProduceBlock(ctx context.Context, req *proto.ProduceBlockRequest) (*proto.ProduceBlockResponse, error) {
	r.logger.Info("Producing block in repository",
		logging.String("validator", req.ValidatorId),
		logging.Int("transaction_count", len(req.TransactionHashes)))

	// Priority 1: Produce block on Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		result, err := r.produceBlockOnKeeper(ctx, req)
		if err == nil {
			// Save to PostgreSQL for analytics (async with error handling)
			go func() {
				if r.db != nil {
					// Use background context với timeout for async operation
					bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()

					// Preserve correlation ID from original context
					correlationID := utils.GetCorrelationID(ctx)

					if err := r.saveBlockToDatabase(bgCtx, result); err != nil {
						r.logger.Error("Failed to save block analytics (async)",
							logging.Error(err),
							logging.String("block_hash", result.BlockHash),
							logging.Int64("block_number", result.BlockNumber),
							logging.String("correlation_id", correlationID))
						// Continue even if database save fails (keeper is primary, analytics only)
					}
				}
			}()
			return result, nil
		}
		r.logger.Warn("Failed to produce block on keeper, falling back to database",
			logging.Error(err),
			logging.String("validator", req.ValidatorId))
	}

	// Fallback to database
	return r.produceBlockInDatabase(ctx, req)
}

// ValidateBlock validates a block
func (r *Repository) ValidateBlock(ctx context.Context, req *proto.ValidateBlockRequest) (*proto.ValidateBlockResponse, error) {
	r.logger.Info("Validating block in repository",
		logging.String("block_hash", req.BlockHash),
		logging.Int64("block_number", req.BlockNumber))

	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		var isValid bool
		var err error
		if req.BlockHash != "" {
			isValid, err = r.validateBlockByHash(ctx, req.BlockHash)
		} else if req.BlockNumber > 0 {
			isValid, err = r.validateBlockByNumber(ctx, req.BlockNumber)
		}
		if err == nil {
			return &proto.ValidateBlockResponse{
				Valid:            isValid,
				ValidationResult: r.getValidationResult(isValid),
				ValidatedAt:      timestamppb.New(time.Now()),
			}, nil
		}
	}

	// Fallback to database validation
	return r.validateBlockFromDatabase(ctx, req)
}

// GetBlock retrieves a block by number
func (r *Repository) GetBlock(ctx context.Context, req *proto.GetBlockRequest) (*proto.GetBlockResponse, error) {
	if req.BlockNumber <= 0 {
		return nil, repoerrors.NewValidationError("block_number", "must be greater than 0")
	}

	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if block, err := r.getBlockFromKeeper(ctx, req.BlockNumber); err == nil && block != nil {
			return block, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	block, err := r.getBlockFromDatabase(ctx, req.BlockNumber)
	if err != nil {
		if de, ok := err.(*sherrors.DomainError); ok && string(de.Code) == fmt.Sprintf("REPO_%d", repoerrors.ErrBlockNotFound) {
			return &proto.GetBlockResponse{}, nil
		}
		return nil, err
	}
	return block, nil
}

// GetBlockByHash retrieves a block by hash
func (r *Repository) GetBlockByHash(ctx context.Context, req *proto.GetBlockByHashRequest) (*proto.GetBlockResponse, error) {
	if req.BlockHash == "" {
		return nil, repoerrors.NewValidationError("block_hash", "is required")
	}

	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if block, err := r.getBlockByHashFromKeeper(ctx, req.BlockHash); err == nil && block != nil {
			return block, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	block, err := r.getBlockByHashFromDatabase(ctx, req.BlockHash)
	if err != nil {
		if de, ok := err.(*sherrors.DomainError); ok && string(de.Code) == fmt.Sprintf("REPO_%d", repoerrors.ErrBlockNotFound) {
			return &proto.GetBlockResponse{}, nil
		}
		return nil, err
	}
	return block, nil
}

// GetLatestBlock retrieves the latest block
func (r *Repository) GetLatestBlock(ctx context.Context) (*proto.GetBlockResponse, error) {
	// Priority 1: Keeper (RocksDB - blockchain state)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if block, err := r.getLatestBlockFromKeeper(ctx); err == nil && block != nil {
			return block, nil
		}
	}

	// Priority 2: PostgreSQL (analytics - fallback)
	if resp, err := r.getLatestBlockFromDatabase(ctx); err == nil && resp != nil {
		return resp, nil
	}

	return &proto.GetBlockResponse{}, repoerrors.NewNotFoundError("block", "latest block not found")
}

// GetBlockRange retrieves a range of blocks
func (r *Repository) GetBlockRange(ctx context.Context, req *proto.GetBlockRangeRequest) (*proto.GetBlockRangeResponse, error) {
	if req.StartBlock <= 0 || req.EndBlock <= 0 {
		return &proto.GetBlockRangeResponse{}, fmt.Errorf("start_block and end_block must be greater than 0")
	}
	if req.StartBlock > req.EndBlock {
		return &proto.GetBlockRangeResponse{}, fmt.Errorf("start_block must be less than or equal to end_block")
	}

	limit, offset := utils.NormalizePagination(req.Limit, req.Offset, utils.PaginationConfig{
		DefaultLimit: 50,
		MaxLimit:     100,
	})

	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.getBlockRangeFromKeeper(ctx, req.StartBlock, req.EndBlock, limit, offset, req.IncludeTransactions); err == nil {
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.getBlockRangeFromDatabase(ctx, int32(req.StartBlock), int32(req.EndBlock), limit, offset, req.IncludeTransactions)
}

// Helper methods for BlockKeeper integration

// getSDKContext creates a sdk.Context from context.Context
// Uses shared utility to avoid code duplication
func (r *Repository) getSDKContext(ctx context.Context) (sdk.Context, error) {
	return utils.GetSDKContext(ctx, r.cosmosApp, r.logger)
}

// getBlockFromKeeper gets a block by number from BlockKeeper
func (r *Repository) getBlockFromKeeper(ctx context.Context, blockNumber int64) (*proto.GetBlockResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		r.logger.Debug("Failed to get SDK context",
			logging.Error(err),
			logging.Int64("block_number", blockNumber))
		return nil, repoerrors.NewBlockchainError("get_sdk_context", err)
	}

	block, err := r.cosmosApp.BlockKeeper.GetBlockByHeight(sdkCtx, blockNumber)
	if err != nil {
		// COSMOS SDK 0.53.4: If CheckTx context is used, blocks saved in InitGenesis may not be visible
		// This is because CheckTx context reads from cache, not committed state
		// Return error to trigger fallback to database (acceptable for queries)
		if sdkCtx.IsCheckTx() {
			r.logger.Debug("Block not found in keeper (CheckTx context) - will fallback to database",
				logging.Int64("block_number", blockNumber),
				logging.Bool("is_check_tx", true))
			// Return error to trigger fallback to database in GetBlock()
			return nil, repoerrors.WrapRepositoryError(repoerrors.ErrBlockNotFound, err,
				fmt.Sprintf("block_number=%d (CheckTx context)", blockNumber))
		}

		// Log detailed error info for debugging
		r.logger.Info("Block not found in keeper - detailed debug",
			logging.Error(err),
			logging.Int64("block_number", blockNumber),
			logging.Int64("current_height", sdkCtx.BlockHeight()),
			logging.String("error_message", err.Error()),
			logging.Bool("is_check_tx", sdkCtx.IsCheckTx()))

		// Return error to trigger fallback to database
		return nil, repoerrors.WrapRepositoryError(repoerrors.ErrBlockNotFound, err,
			fmt.Sprintf("block_number=%d, current_height=%d", blockNumber, sdkCtx.BlockHeight()))
	}

	// Check if block is zero value (empty struct)
	if block.Height == 0 && block.Hash == "" {
		r.logger.Debug("Block is zero value from keeper",
			logging.Int64("block_number", blockNumber))
		return nil, repoerrors.NewNotFoundError("block", fmt.Sprintf("block_number=%d", blockNumber))
	}

	r.logger.Debug("Block retrieved from keeper",
		logging.Int64("block_number", blockNumber),
		logging.String("block_hash", block.Hash))
	return r.convertBlockToProto(block), nil
}

// getBlockByHashFromKeeper gets a block by hash from BlockKeeper
func (r *Repository) getBlockByHashFromKeeper(ctx context.Context, blockHash string) (*proto.GetBlockResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		r.logger.Debug("Failed to get SDK context",
			logging.Error(err),
			logging.String("block_hash", blockHash))
		return nil, repoerrors.NewBlockchainError("get_sdk_context", err)
	}

	block, err := r.cosmosApp.BlockKeeper.GetBlockByHash(sdkCtx, blockHash)
	if err != nil {
		r.logger.Debug("Block not found in keeper by hash",
			logging.Error(err),
			logging.String("block_hash", blockHash),
			logging.Int64("current_height", sdkCtx.BlockHeight()))
		return nil, repoerrors.WrapRepositoryError(repoerrors.ErrBlockNotFound, err,
			fmt.Sprintf("block_hash=%s, current_height=%d", blockHash, sdkCtx.BlockHeight()))
	}

	// Check if block is zero value (empty struct)
	if block.Height == 0 && block.Hash == "" {
		r.logger.Debug("Block is zero value from keeper by hash",
			logging.String("block_hash", blockHash))
		return nil, repoerrors.NewNotFoundError("block", fmt.Sprintf("block_hash=%s", blockHash))
	}

	r.logger.Debug("Block retrieved from keeper by hash",
		logging.String("block_hash", blockHash),
		logging.Int64("block_number", block.Height))
	return r.convertBlockToProto(block), nil
}

// getLatestBlockFromKeeper gets the latest block from BlockKeeper
func (r *Repository) getLatestBlockFromKeeper(ctx context.Context) (*proto.GetBlockResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		r.logger.Debug("Failed to get SDK context for latest block",
			logging.Error(err))
		return nil, repoerrors.NewBlockchainError("get_sdk_context", err)
	}

	// Get current block height
	height := sdkCtx.BlockHeight()
	if height <= 0 {
		r.logger.Debug("No blocks available in keeper",
			logging.Int64("height", height))
		return nil, repoerrors.NewNotFoundError("block", fmt.Sprintf("height=%d", height))
	}

	block, err := r.cosmosApp.BlockKeeper.GetBlockByHeight(sdkCtx, height)
	if err != nil {
		r.logger.Debug("Latest block not found in keeper",
			logging.Error(err),
			logging.Int64("height", height))
		return nil, repoerrors.WrapRepositoryError(repoerrors.ErrBlockNotFound, err,
			fmt.Sprintf("height=%d", height))
	}

	// Check if block is zero value (empty struct)
	if block.Height == 0 && block.Hash == "" {
		r.logger.Debug("Latest block is zero value from keeper",
			logging.Int64("height", height))
		return nil, repoerrors.NewNotFoundError("block", fmt.Sprintf("height=%d", height))
	}

	r.logger.Debug("Latest block retrieved from keeper",
		logging.Int64("height", height),
		logging.String("block_hash", block.Hash))
	return r.convertBlockToProto(block), nil
}

// validateBlockByHash validates a block by hash using BlockKeeper
func (r *Repository) validateBlockByHash(ctx context.Context, blockHash string) (bool, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return false, err
	}

	block, err := r.cosmosApp.BlockKeeper.GetBlockByHash(sdkCtx, blockHash)
	if err != nil {
		return false, err
	}

	return r.cosmosApp.BlockKeeper.ValidateBlock(sdkCtx, block), nil
}

// validateBlockByNumber validates a block by number using BlockKeeper
func (r *Repository) validateBlockByNumber(ctx context.Context, blockNumber int64) (bool, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return false, err
	}

	block, err := r.cosmosApp.BlockKeeper.GetBlockByHeight(sdkCtx, blockNumber)
	if err != nil {
		return false, err
	}

	return r.cosmosApp.BlockKeeper.ValidateBlock(sdkCtx, block), nil
}

// convertBlockToProto converts blocktypes.Block to proto.GetBlockResponse
func (r *Repository) convertBlockToProto(block blocktypes.Block) *proto.GetBlockResponse {
	return &proto.GetBlockResponse{
		BlockHash:         block.Hash,
		BlockNumber:       block.Height,
		PreviousBlockHash: block.PreviousHash,
		MerkleRoot:        block.Hash, // Use block hash as merkle root (as per query_server.go)
		Timestamp:         block.Timestamp.Unix(),
		ValidatorAddress:  block.Validator,
		GasUsed:           block.GasUsed,
		GasLimit:          block.GasLimit,
		CreatedAt:         timestamppb.New(block.CreatedAt),
	}
}

// getValidationResult returns validation result message
func (r *Repository) getValidationResult(isValid bool) string {
	if isValid {
		return "Block validation successful"
	}
	return "Block validation failed"
}

// produceBlockOnKeeper produces a block on Cosmos SDK blockchain
func (r *Repository) produceBlockOnKeeper(ctx context.Context, req *proto.ProduceBlockRequest) (*proto.ProduceBlockResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, repoerrors.NewBlockchainError("get_sdk_context", err)
	}

	// Determine next block height: use GetAllBlocks to find highest block
	allBlocks := r.cosmosApp.BlockKeeper.GetAllBlocks(sdkCtx)
	var highestHeight int64
	for _, block := range allBlocks {
		if block.Height > highestHeight {
			highestHeight = block.Height
		}
	}
	newHeight := highestHeight + 1

	// Generate block hash from header data
	header := sdkCtx.BlockHeader()
	data := fmt.Sprintf("%s:%d:%s:%s:%d", header.ChainID, newHeight, hex.EncodeToString(header.AppHash), hex.EncodeToString(header.DataHash), len(req.TransactionHashes))
	hashBytes := sha256.Sum256([]byte(data))
	blockHash := hex.EncodeToString(hashBytes[:])

	// Get previous block hash
	var previousHash string
	if highestHeight > 0 {
		if prevBlock, err := r.cosmosApp.BlockKeeper.GetBlockByHeight(sdkCtx, highestHeight); err == nil {
			previousHash = prevBlock.Hash
		}
	}

	// Create and store block
	block := blocktypes.Block{
		ID:           fmt.Sprintf("block_%d", newHeight),
		Height:       newHeight,
		Hash:         blockHash,
		PreviousHash: previousHash,
		Timestamp:    time.Unix(req.Timestamp, 0),
		Validator:    req.ValidatorId,
		Size:         int64(len(req.TransactionHashes) * 100),
		TxCount:      int64(len(req.TransactionHashes)),
		GasUsed:      21000,
		GasLimit:     21000000,
		Status:       "pending",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := r.cosmosApp.BlockKeeper.SetBlock(sdkCtx, block); err != nil {
		return nil, fmt.Errorf("failed to store block in keeper: %w", err)
	}

	r.logger.Info("Block produced on keeper",
		logging.Int64("block_number", newHeight),
		logging.String("block_hash", blockHash),
		logging.String("validator", req.ValidatorId))

	return &proto.ProduceBlockResponse{
		Success:     true,
		BlockHash:   blockHash,
		BlockNumber: newHeight,
	}, nil
}

// produceBlockInDatabase produces a block in database (fallback)
func (r *Repository) produceBlockInDatabase(ctx context.Context, req *proto.ProduceBlockRequest) (*proto.ProduceBlockResponse, error) {
	// Use helper function to get PostgreSQL connection (reduces duplicate code)
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return nil, fmt.Errorf("database not available")
	}

	query := `
		INSERT INTO usc_block_analytics (
			block_number, block_hash, validator_address, timestamp, 
			transaction_count, total_usc_transferred, gas_used, gas_limit,
			block_size_bytes, is_finalized
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (block_number) DO UPDATE SET
			block_hash = EXCLUDED.block_hash,
			validator_address = EXCLUDED.validator_address,
			timestamp = EXCLUDED.timestamp,
			transaction_count = EXCLUDED.transaction_count,
			total_usc_transferred = EXCLUDED.total_usc_transferred,
			gas_used = EXCLUDED.gas_used,
			gas_limit = EXCLUDED.gas_limit,
			block_size_bytes = EXCLUDED.block_size_bytes,
			is_finalized = EXCLUDED.is_finalized
	`

	blockNumber := time.Now().Unix()
	// Generate real hash for database analytics
	dataStr := fmt.Sprintf("%s:%d:%s:%s", req.ValidatorId, blockNumber, fmt.Sprintf("%v", req.TransactionHashes), time.Now().Format(time.RFC3339))
	hashBytes := sha256.Sum256([]byte(dataStr))
	blockHash := "0x" + hex.EncodeToString(hashBytes[:])

	_, err := postgres.ExecContext(ctx, query,
		blockNumber,
		blockHash,
		req.ValidatorId,
		time.Now(),
		len(req.TransactionHashes),
		"0",
		21000,
		21000000,
		1024,
		true,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to produce block in database: %w", err)
	}

	return &proto.ProduceBlockResponse{
		Success:     true,
		BlockHash:   blockHash,
		BlockNumber: blockNumber,
	}, nil
}

// saveBlockToDatabase saves a block to database for analytics (async with error return)
func (r *Repository) saveBlockToDatabase(ctx context.Context, resp *proto.ProduceBlockResponse) error {
	// Use helper function to get PostgreSQL connection (reduces duplicate code)
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return fmt.Errorf("postgres connection not available")
	}
	if resp == nil {
		return fmt.Errorf("response is nil")
	}

	query := `
		INSERT INTO usc_block_analytics (
			block_number, block_hash, validator_address, timestamp, 
			transaction_count, total_usc_transferred, gas_used, gas_limit,
			block_size_bytes, is_finalized
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (block_number) DO UPDATE SET
			block_hash = EXCLUDED.block_hash,
			validator_address = EXCLUDED.validator_address,
			timestamp = EXCLUDED.timestamp
	`

	if _, err := postgres.ExecContext(ctx, query,
		resp.BlockNumber,
		resp.BlockHash,
		"", // validator not in response
		time.Now(),
		0,   // tx count
		"0", // total USC
		21000,
		21000000,
		1024,
		true,
	); err != nil {
		return fmt.Errorf("failed to save block to database for analytics: %w", err)
	}

	return nil
}

// getBlockRangeFromKeeper retrieves a range of blocks from BlockKeeper
func (r *Repository) getBlockRangeFromKeeper(ctx context.Context, startBlock, endBlock int64, limit, offset int32, includeTransactions bool) (*proto.GetBlockRangeResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, repoerrors.NewBlockchainError("get_sdk_context", err)
	}

	// Get all blocks from keeper
	allBlocks := r.cosmosApp.BlockKeeper.GetAllBlocks(sdkCtx)

	// Filter by range
	// Pre-allocate slice with estimated capacity (worst case: all blocks in range)
	// Use min(endBlock-startBlock+1, len(allBlocks)) as capacity estimate
	estimatedCapacity := int(endBlock - startBlock + 1)
	if estimatedCapacity > len(allBlocks) {
		estimatedCapacity = len(allBlocks)
	}
	filteredBlocks := make([]blocktypes.Block, 0, estimatedCapacity)
	for _, block := range allBlocks {
		if block.Height >= startBlock && block.Height <= endBlock {
			filteredBlocks = append(filteredBlocks, block)
		}
	}

	// Apply pagination
	totalCount := int32(len(filteredBlocks))
	startIdx := int(offset)
	endIdx := startIdx + int(limit)
	if endIdx > len(filteredBlocks) {
		endIdx = len(filteredBlocks)
	}

	if startIdx >= len(filteredBlocks) {
		return &proto.GetBlockRangeResponse{
			Blocks:     []*proto.GetBlockResponse{},
			TotalCount: totalCount,
			HasMore:    false,
			NextOffset: int64(offset),
		}, nil
	}

	// Convert to proto
	blocks := make([]*proto.GetBlockResponse, 0, endIdx-startIdx)
	for i := startIdx; i < endIdx; i++ {
		blocks = append(blocks, r.convertBlockToProto(filteredBlocks[i]))
	}

	hasMore := endIdx < len(filteredBlocks)
	nextOffset := offset + int32(len(blocks))

	return &proto.GetBlockRangeResponse{
		Blocks:     blocks,
		TotalCount: totalCount,
		HasMore:    hasMore,
		NextOffset: int64(nextOffset),
	}, nil
}

// Database helper methods

// getBlockFromDatabase retrieves a block by number from database
func (r *Repository) getBlockFromDatabase(ctx context.Context, blockNumber int64) (*proto.GetBlockResponse, error) {
	// Use helper function to get PostgreSQL connection (reduces duplicate code)
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return &proto.GetBlockResponse{}, fmt.Errorf("database not available")
	}

	query := `
		SELECT block_number, block_hash, validator_address, timestamp,
		       transaction_count, gas_used, gas_limit, block_size_bytes, is_finalized
		FROM usc_block_analytics
		WHERE block_number = $1
		LIMIT 1
	`

	var blockNum, txCount, gasUsed, gasLimit, blockSize int64
	var blockHash, validator string
	var timestamp time.Time
	var isFinalized bool

	err := postgres.QueryRowContext(ctx, query, blockNumber).Scan(
		&blockNum, &blockHash, &validator, &timestamp,
		&txCount, &gasUsed, &gasLimit, &blockSize, &isFinalized,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Debug("Block not found in database",
				logging.Int64("block_number", blockNumber))
			return &proto.GetBlockResponse{}, fmt.Errorf("block not found: block_number=%d", blockNumber)
		}
		r.logger.Warn("Failed to query block from database",
			logging.Error(err),
			logging.Int64("block_number", blockNumber))
		return &proto.GetBlockResponse{}, fmt.Errorf("failed to query block: %w", err)
	}

	return &proto.GetBlockResponse{
		BlockNumber:      blockNum,
		BlockHash:        blockHash,
		ValidatorAddress: validator,
		Timestamp:        timestamp.Unix(),
		GasUsed:          gasUsed,
		GasLimit:         gasLimit,
		CreatedAt:        timestamppb.New(timestamp),
	}, nil
}

// getBlockByHashFromDatabase retrieves a block by hash from database
func (r *Repository) getBlockByHashFromDatabase(ctx context.Context, blockHash string) (*proto.GetBlockResponse, error) {
	// Use helper function to get PostgreSQL connection (reduces duplicate code)
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return &proto.GetBlockResponse{}, fmt.Errorf("database not available")
	}

	query := `
		SELECT block_number, block_hash, validator_address, timestamp,
		       transaction_count, gas_used, gas_limit, block_size_bytes, is_finalized
		FROM usc_block_analytics
		WHERE block_hash = $1
		LIMIT 1
	`

	var blockNum, txCount, gasUsed, gasLimit, blockSize int64
	var hash, validator string
	var timestamp time.Time
	var isFinalized bool

	err := postgres.QueryRowContext(ctx, query, blockHash).Scan(
		&blockNum, &hash, &validator, &timestamp,
		&txCount, &gasUsed, &gasLimit, &blockSize, &isFinalized,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Debug("Block not found in database",
				logging.String("block_hash", blockHash))
			return &proto.GetBlockResponse{}, fmt.Errorf("block not found: block_hash=%s", blockHash)
		}
		r.logger.Warn("Failed to query block from database",
			logging.Error(err),
			logging.String("block_hash", blockHash))
		return &proto.GetBlockResponse{}, fmt.Errorf("failed to query block: %w", err)
	}

	return &proto.GetBlockResponse{
		BlockNumber:      blockNum,
		BlockHash:        hash,
		ValidatorAddress: validator,
		Timestamp:        timestamp.Unix(),
		GasUsed:          gasUsed,
		GasLimit:         gasLimit,
		CreatedAt:        timestamppb.New(timestamp),
	}, nil
}

// getLatestBlockFromDatabase retrieves the latest block from database
func (r *Repository) getLatestBlockFromDatabase(ctx context.Context) (*proto.GetBlockResponse, error) {
	// Use helper function to get PostgreSQL connection (reduces duplicate code)
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return &proto.GetBlockResponse{}, fmt.Errorf("database not available")
	}

	query := `
		SELECT block_number, block_hash, validator_address, timestamp,
		       transaction_count, gas_used, gas_limit, block_size_bytes, is_finalized
		FROM usc_block_analytics
		ORDER BY block_number DESC
		LIMIT 1
	`

	var blockNum, txCount, gasUsed, gasLimit, blockSize int64
	var blockHash, validator string
	var timestamp time.Time
	var isFinalized bool

	err := postgres.QueryRowContext(ctx, query).Scan(
		&blockNum, &blockHash, &validator, &timestamp,
		&txCount, &gasUsed, &gasLimit, &blockSize, &isFinalized,
	)
	if err != nil {
		return &proto.GetBlockResponse{}, fmt.Errorf("no blocks found: %w", err)
	}

	return &proto.GetBlockResponse{
		BlockNumber:      blockNum,
		BlockHash:        blockHash,
		ValidatorAddress: validator,
		Timestamp:        timestamp.Unix(),
		GasUsed:          gasUsed,
		GasLimit:         gasLimit,
		CreatedAt:        timestamppb.New(timestamp),
	}, nil
}

// getBlockRangeFromDatabase retrieves a range of blocks from database
func (r *Repository) getBlockRangeFromDatabase(ctx context.Context, startBlock, endBlock, limit, offset int32, includeTransactions bool) (*proto.GetBlockRangeResponse, error) {
	// Use helper function to get PostgreSQL connection (reduces duplicate code)
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return &proto.GetBlockRangeResponse{}, fmt.Errorf("database not available")
	}

	// Optimized query: Use window function to get total count in single query (eliminates N+1 query)
	query := `
		SELECT block_number, block_hash, validator_address, timestamp,
		       transaction_count, gas_used, gas_limit, block_size_bytes, is_finalized,
		       COUNT(*) OVER() as total_count
		FROM usc_block_analytics
		WHERE block_number >= $1 AND block_number <= $2
		ORDER BY block_number DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := postgres.QueryContext(ctx, query, startBlock, endBlock, limit, offset)
	if err != nil {
		return &proto.GetBlockRangeResponse{}, fmt.Errorf("failed to query blocks: %w", err)
	}
	defer rows.Close()

	// Pre-allocate slice with capacity = limit for better performance
	blocks := make([]*proto.GetBlockResponse, 0, limit)
	var totalCount int32
	for rows.Next() {
		var blockNum, txCount, gasUsed, gasLimit, blockSize int64
		var blockHash, validator string
		var timestamp time.Time
		var isFinalized bool

		if err := rows.Scan(
			&blockNum, &blockHash, &validator, &timestamp,
			&txCount, &gasUsed, &gasLimit, &blockSize, &isFinalized,
			&totalCount,
		); err != nil {
			r.logger.Warn("Failed to scan block row", logging.Error(err))
			continue
		}

		blocks = append(blocks, &proto.GetBlockResponse{
			BlockNumber:      blockNum,
			BlockHash:        blockHash,
			ValidatorAddress: validator,
			Timestamp:        timestamp.Unix(),
			GasUsed:          gasUsed,
			GasLimit:         gasLimit,
			CreatedAt:        timestamppb.New(timestamp),
		})
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		r.logger.Warn("Error during rows iteration", logging.Error(err))
		// Use count from results if iteration error occurs
		if totalCount == 0 {
			totalCount = int32(len(blocks))
		}
	}

	// Fallback to result count if totalCount is 0 (shouldn't happen, but safety check)
	if totalCount == 0 {
		totalCount = int32(len(blocks))
	}

	hasMore := int32(len(blocks)) == limit && (int32(offset)+limit) < totalCount
	nextOffset := int32(offset) + limit
	if !hasMore {
		nextOffset = 0
	}

	return &proto.GetBlockRangeResponse{
		Blocks:     blocks,
		TotalCount: totalCount,
		HasMore:    hasMore,
		NextOffset: int64(nextOffset),
	}, nil
}

// validateBlockFromDatabase validates a block from database
func (r *Repository) validateBlockFromDatabase(ctx context.Context, req *proto.ValidateBlockRequest) (*proto.ValidateBlockResponse, error) {
	// Use helper function to get PostgreSQL connection (reduces duplicate code)
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return &proto.ValidateBlockResponse{
			Valid:            false,
			ValidationResult: "database not available",
		}, nil
	}

	var exists bool
	var err error

	if req.BlockHash != "" {
		err = postgres.QueryRowContext(ctx,
			"SELECT EXISTS(SELECT 1 FROM usc_block_analytics WHERE block_hash = $1)",
			req.BlockHash).Scan(&exists)
	} else if req.BlockNumber > 0 {
		err = postgres.QueryRowContext(ctx,
			"SELECT EXISTS(SELECT 1 FROM usc_block_analytics WHERE block_number = $1)",
			req.BlockNumber).Scan(&exists)
	} else {
		return &proto.ValidateBlockResponse{
			Valid:            false,
			ValidationResult: "block_hash or block_number is required",
		}, nil
	}

	if err != nil {
		r.logger.Error("Failed to validate block",
			logging.Error(err))
		return &proto.ValidateBlockResponse{
			Valid:            false,
			ValidationResult: fmt.Sprintf("validation error: %v", err),
		}, nil
	}

	if !exists {
		return &proto.ValidateBlockResponse{
			Valid:            false,
			ValidationResult: "Block not found",
			ValidatedAt:      timestamppb.New(time.Now()),
		}, nil
	}

	return &proto.ValidateBlockResponse{
		Valid:            true,
		ValidationResult: "Block validation successful",
		ValidatedAt:      timestamppb.New(time.Now()),
	}, nil
}
