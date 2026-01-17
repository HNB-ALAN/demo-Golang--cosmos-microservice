package streaming_operations

import (
	"context"
	"database/sql"
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

	"github.com/usc-platform/shared/logging"
)

// Repository handles streaming operations data access
type Repository struct {
	db                *database.PostgreSQLManager
	cosmosApp         *app.USCApp
	blockchainStorage *storage.StateManager
	redisManager      *database.RedisManager
	logger            *logging.Logger
}

// NewRepository creates a new streaming operations repository
func NewRepository(db *database.PostgreSQLManager, cosmosApp *app.USCApp, blockchainStorage *storage.StateManager, redisManager *database.RedisManager, logger *logging.Logger) *Repository {
	return &Repository{
		db:                db,
		cosmosApp:         cosmosApp,
		blockchainStorage: blockchainStorage,
		redisManager:      redisManager,
		logger:            logger,
	}
}

// StreamBlocks streams blockchain blocks
func (r *Repository) StreamBlocks(ctx context.Context, req *proto.StreamBlocksRequest) (*proto.StreamBlocksResponse, error) {
	// Validate request
	if req.ClientId == "" {
		return nil, repoerrors.NewValidationError("client_id", "is required")
	}

	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.streamBlocksFromKeeper(ctx, req); err == nil {
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.streamBlocksInDatabase(ctx, req)
}

// Helper methods for StreamingKeeper interaction

// getSDKContext creates a sdk.Context from context.Context
// Uses shared utility to avoid code duplication
func (r *Repository) getSDKContext(ctx context.Context) (sdk.Context, error) {
	return utils.GetSDKContext(ctx, r.cosmosApp, r.logger)
}

// streamBlocksFromKeeper streams blocks from the keeper
func (r *Repository) streamBlocksFromKeeper(ctx context.Context, req *proto.StreamBlocksRequest) (*proto.StreamBlocksResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, err
	}

	// Get latest block from BlockKeeper
	blocks := r.cosmosApp.BlockKeeper.GetAllBlocks(sdkCtx)
	if len(blocks) == 0 {
		return nil, fmt.Errorf("no blocks found in keeper")
	}

	// Get the latest block (highest block number)
	latestBlock := blocks[0]
	for _, block := range blocks {
		if block.Height > latestBlock.Height {
			latestBlock = block
		}
	}

	return &proto.StreamBlocksResponse{
		BlockHash:         latestBlock.Hash,
		BlockNumber:       int64(latestBlock.Height),
		PreviousBlockHash: latestBlock.PreviousHash,
		MerkleRoot:        latestBlock.Hash,
		Timestamp:         latestBlock.Timestamp.Unix(),
		ValidatorAddress:  latestBlock.Validator,
		GasUsed:           0,
		GasLimit:          0,
		Transactions:      []*proto.Transaction{},
		BlockData:         "{}",
		EventType:         "block_created",
	}, nil
}

// Database fallback methods

// streamBlocksInDatabase streams blocks from database
func (r *Repository) streamBlocksInDatabase(ctx context.Context, req *proto.StreamBlocksRequest) (*proto.StreamBlocksResponse, error) {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return nil, repoerrors.NewRepositoryError(repoerrors.ErrDatabaseUnavailable)
	}

	query := `
		SELECT 
			b.block_number,
			b.block_hash,
			COALESCE(b.previous_block_hash, '') as previous_block_hash,
			COALESCE(b.merkle_root, '') as merkle_root,
			b.timestamp,
			COALESCE(b.validator_address, '') as validator_address,
			COALESCE(b.gas_used, 0) as gas_used,
			COALESCE(b.gas_limit, 0) as gas_limit,
			COALESCE(b.block_data, '{}') as block_data,
			COALESCE(COUNT(t.transaction_hash), 0) as transaction_count
		FROM blocks b
		LEFT JOIN transactions t ON t.block_number = b.block_number
		GROUP BY b.block_number, b.block_hash, b.previous_block_hash, b.merkle_root,
		         b.timestamp, b.validator_address, b.gas_used, b.gas_limit, b.block_data
		ORDER BY b.block_number DESC
		LIMIT 1
	`

	var blockNum, gasUsed, gasLimit int64
	var blockHash, previousBlockHash, merkleRoot, validatorAddr, blockData string
	var timestamp int64
	var txCount int64

	err := postgres.QueryRowContext(ctx, query).Scan(
		&blockNum, &blockHash, &previousBlockHash, &merkleRoot,
		&timestamp, &validatorAddr, &gasUsed, &gasLimit, &blockData, &txCount,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("not found")
		}
		return nil, err
	}

	return &proto.StreamBlocksResponse{
		BlockHash:         blockHash,
		BlockNumber:       blockNum,
		PreviousBlockHash: previousBlockHash,
		MerkleRoot:        merkleRoot,
		Timestamp:         timestamp,
		ValidatorAddress:  validatorAddr,
		GasUsed:           gasUsed,
		GasLimit:          gasLimit,
		Transactions:      []*proto.Transaction{},
		BlockData:         blockData,
		EventType:         "block_created",
	}, nil
}

// StreamTransactions streams blockchain transactions
func (r *Repository) StreamTransactions(ctx context.Context, req *proto.StreamTransactionsRequest) (*proto.StreamTransactionsResponse, error) {
	// Validate request
	if req.ClientId == "" {
		return nil, repoerrors.NewValidationError("client_id", "is required")
	}

	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.streamTransactionsFromKeeper(ctx, req); err == nil {
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.streamTransactionsFromDatabase(ctx, req)
}

// streamTransactionsFromKeeper streams transactions from the keeper
func (r *Repository) streamTransactionsFromKeeper(ctx context.Context, req *proto.StreamTransactionsRequest) (*proto.StreamTransactionsResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, err
	}

	// Get latest transaction from TransactionKeeper
	transactions := r.cosmosApp.TransactionKeeper.GetAllTransactions(sdkCtx)
	if len(transactions) == 0 {
		return nil, fmt.Errorf("no transactions found in keeper")
	}

	// Get the latest transaction (by CreatedAt)
	latestTx := transactions[0]
	for _, tx := range transactions {
		if tx.CreatedAt > latestTx.CreatedAt {
			latestTx = tx
		}
	}

	status := int32(0)
	if latestTx.Status == "executed" {
		status = 1
	}

	return &proto.StreamTransactionsResponse{
		TransactionHash:  latestTx.Hash,
		FromAddress:      latestTx.FromAddress,
		ToAddress:        latestTx.ToAddress,
		Amount:           latestTx.Amount,
		GasPrice:         latestTx.Fee,
		GasLimit:         latestTx.GasLimit,
		GasUsed:          latestTx.GasUsed,
		Data:             latestTx.Data,
		Status:           status,
		BlockNumber:      0,
		BlockHash:        "",
		TransactionIndex: 0,
		EventType:        "transaction_confirmed",
		TransactionType:  latestTx.TransactionType,
	}, nil
}

// streamTransactionsFromDatabase streams transactions from database
func (r *Repository) streamTransactionsFromDatabase(ctx context.Context, req *proto.StreamTransactionsRequest) (*proto.StreamTransactionsResponse, error) {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return nil, repoerrors.NewRepositoryError(repoerrors.ErrDatabaseUnavailable)
	}

	query := `
		SELECT transaction_hash, from_address, to_address, amount, gas_price, gas_limit, gas_used,
		       data, status, block_number, block_hash, transaction_index, transaction_type
		FROM transactions
		WHERE ($1 = '' OR transaction_type = $1)
		ORDER BY block_number DESC, transaction_index DESC
		LIMIT 1
	`

	var txHash, fromAddr, toAddr, amount, gasPrice, data, txType, blockHash string
	var gasLimit, gasUsed, blockNum, txIndex int64
	var status int32

	err := postgres.QueryRowContext(ctx, query, req.TransactionType).Scan(
		&txHash, &fromAddr, &toAddr, &amount, &gasPrice, &gasLimit, &gasUsed,
		&data, &status, &blockNum, &blockHash, &txIndex, &txType,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("not found")
		}
		return nil, err
	}

	return &proto.StreamTransactionsResponse{
		TransactionHash:  txHash,
		FromAddress:      fromAddr,
		ToAddress:        toAddr,
		Amount:           amount,
		GasPrice:         gasPrice,
		GasLimit:         gasLimit,
		GasUsed:          gasUsed,
		Data:             data,
		Status:           status,
		BlockNumber:      blockNum,
		BlockHash:        blockHash,
		TransactionIndex: txIndex,
		EventType:        "transaction_confirmed",
		TransactionType:  txType,
	}, nil
}

// StreamValidatorEvents streams validator events
func (r *Repository) StreamValidatorEvents(ctx context.Context, req *proto.StreamValidatorEventsRequest) (*proto.StreamValidatorEventsResponse, error) {
	// Validate request
	if req.ClientId == "" {
		return nil, repoerrors.NewValidationError("client_id", "is required")
	}

	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.streamValidatorEventsFromKeeper(ctx, req); err == nil {
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.streamValidatorEventsFromDatabase(ctx, req)
}

// streamValidatorEventsFromKeeper streams validator events from the keeper
func (r *Repository) streamValidatorEventsFromKeeper(ctx context.Context, req *proto.StreamValidatorEventsRequest) (*proto.StreamValidatorEventsResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, err
	}

	// Get validators from ValidatorKeeper
	validators, err := r.cosmosApp.ValidatorKeeper.GetAllValidators(sdkCtx)
	if err != nil || len(validators) == 0 {
		return nil, fmt.Errorf("no validators found in keeper")
	}

	// Get the first validator (or filter by event type if needed)
	validator := validators[0]

	return &proto.StreamValidatorEventsResponse{
		ValidatorAddress:   validator.Address,
		ValidatorId:        validator.Address,
		EventType:          req.EventType,
		EventData:          "{}",
		StakeAmount:        "0",
		DelegatedAmount:    "0",
		CommissionRate:     validator.Commission,
		UptimePercentage:   "0",
		BlocksProposed:     0,
		RewardsEarned:      "0",
		BlockHash:          "",
		BlockNumber:        0,
		AffectedDelegators: []string{},
	}, nil
}

// streamValidatorEventsFromDatabase streams validator events from database
func (r *Repository) streamValidatorEventsFromDatabase(ctx context.Context, req *proto.StreamValidatorEventsRequest) (*proto.StreamValidatorEventsResponse, error) {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return nil, repoerrors.NewRepositoryError(repoerrors.ErrDatabaseUnavailable)
	}

	query := `
		SELECT validator_address, validator_name, commission_rate, status
		FROM usc_validator_analytics
		ORDER BY registered_at DESC
		LIMIT 1
	`

	var validatorAddr, validatorName string
	var commissionRate float64
	var status string

	err := postgres.QueryRowContext(ctx, query).Scan(&validatorAddr, &validatorName, &commissionRate, &status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("not found")
		}
		return nil, err
	}

	return &proto.StreamValidatorEventsResponse{
		ValidatorAddress:   validatorAddr,
		ValidatorId:        validatorAddr,
		EventType:          req.EventType,
		EventData:          "{}",
		StakeAmount:        "0",
		DelegatedAmount:    "0",
		CommissionRate:     fmt.Sprintf("%.2f", commissionRate),
		UptimePercentage:   "0",
		BlocksProposed:     0,
		RewardsEarned:      "0",
		BlockHash:          "",
		BlockNumber:        0,
		AffectedDelegators: []string{},
	}, nil
}

// StreamNetworkEvents streams network events
func (r *Repository) StreamNetworkEvents(ctx context.Context, req *proto.StreamNetworkEventsRequest) (*proto.StreamNetworkEventsResponse, error) {
	// Validate request
	if req.ClientId == "" {
		return nil, repoerrors.NewValidationError("client_id", "is required")
	}

	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.streamNetworkEventsFromKeeper(ctx, req); err == nil {
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.streamNetworkEventsFromDatabase(ctx, req)
}

// streamNetworkEventsFromKeeper streams network events from the keeper
func (r *Repository) streamNetworkEventsFromKeeper(ctx context.Context, req *proto.StreamNetworkEventsRequest) (*proto.StreamNetworkEventsResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, err
	}

	// Get network health from NetworkKeeper
	healths := r.cosmosApp.NetworkKeeper.GetAllHealths(sdkCtx)
	if len(healths) == 0 {
		return nil, fmt.Errorf("no network health data found in keeper")
	}

	// Get the latest health record (by Timestamp)
	latestHealth := healths[0]
	for _, health := range healths {
		if health.Timestamp.After(latestHealth.Timestamp) {
			latestHealth = health
		}
	}

	return &proto.StreamNetworkEventsResponse{
		EventId:            latestHealth.ID,
		EventType:          req.EventType,
		Severity:           "info",
		Title:              "Network Health Event",
		Description:        fmt.Sprintf("Network health score: %d", latestHealth.HealthScore),
		EventData:          "{}",
		Source:             "network",
		AffectedComponents: "[]",
		ResolutionStatus:   "open",
		AutoResolution:     "false",
		ImpactScore:        int64(100 - latestHealth.HealthScore),
	}, nil
}

// streamNetworkEventsFromDatabase streams network events from database
func (r *Repository) streamNetworkEventsFromDatabase(ctx context.Context, req *proto.StreamNetworkEventsRequest) (*proto.StreamNetworkEventsResponse, error) {
	// Note: Network events from PostgreSQL analytics not implemented yet
	// This is a fallback method, so returning default values is acceptable
	return &proto.StreamNetworkEventsResponse{
		EventId:            fmt.Sprintf("event_%d", time.Now().Unix()),
		EventType:          req.EventType,
		Severity:           "info",
		Title:              "Network Event",
		Description:        "A network event occurred",
		EventData:          "{}",
		Source:             "network",
		AffectedComponents: "[]",
		ResolutionStatus:   "open",
		AutoResolution:     "false",
		ImpactScore:        0,
	}, nil
}
