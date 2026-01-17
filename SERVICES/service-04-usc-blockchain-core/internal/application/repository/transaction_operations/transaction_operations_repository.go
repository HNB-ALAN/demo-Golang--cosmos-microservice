package transaction_operations

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
	transactiontypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/transaction/types"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/usc-platform/shared/logging"
)

// Repository handles transaction operations data access
type Repository struct {
	db                *database.PostgreSQLManager
	cosmosApp         *app.USCApp
	blockchainStorage *storage.StateManager
	redisManager      *database.RedisManager
	logger            *logging.Logger
}

// NewRepository creates a new transaction operations repository
func NewRepository(db *database.PostgreSQLManager, cosmosApp *app.USCApp, blockchainStorage *storage.StateManager, redisManager *database.RedisManager, logger *logging.Logger) *Repository {
	return &Repository{
		db:                db,
		cosmosApp:         cosmosApp,
		blockchainStorage: blockchainStorage,
		redisManager:      redisManager,
		logger:            logger,
	}
}

// SubmitTransaction submits a new transaction to the mempool
func (r *Repository) SubmitTransaction(ctx context.Context, req *proto.SubmitTransactionRequest) (*proto.SubmitTransactionResponse, error) {
	if req.FromAddress == "" || req.ToAddress == "" || req.Amount == "" {
		return &proto.SubmitTransactionResponse{
			Status:       2, // Failed
			ErrorMessage: "from_address, to_address, and amount are required",
		}, nil
	}

	// Priority 1: Submit on Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if txHash, err := r.submitTransactionOnBlockchain(ctx, req); err == nil && txHash != "" {
			// Save to PostgreSQL for analytics (async with error handling)
			go func() {
				if r.db != nil {
					// Use background context với timeout for async operation
					bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()
					
					// Preserve correlation ID from original context
					correlationID := utils.GetCorrelationID(ctx)
					
					if err := r.saveTransactionToDatabase(bgCtx, txHash, req); err != nil {
						r.logger.Error("Failed to save transaction analytics (async)",
							logging.Error(err),
							logging.String("transaction_hash", txHash),
							logging.String("correlation_id", correlationID))
						// Continue even if database save fails (keeper is primary, analytics only)
					}
				}
			}()
			return &proto.SubmitTransactionResponse{
				TransactionHash: txHash,
				Status:          0, // Pending
				SubmittedAt:     timestamppb.New(time.Now()),
			}, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.submitTransactionToDatabase(ctx, req)
}

// GetTransaction retrieves a transaction by hash
func (r *Repository) GetTransaction(ctx context.Context, req *proto.GetTransactionRequest) (*proto.GetTransactionResponse, error) {
	if req.TransactionHash == "" {
		return nil, repoerrors.NewValidationError("transaction_hash", "is required")
	}

	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if tx, err := r.getTransactionFromKeeper(ctx, req.TransactionHash); err == nil && tx != nil {
			return tx, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.getTransactionFromDatabase(ctx, req.TransactionHash)
}

// GetTransactionStatus retrieves transaction status
func (r *Repository) GetTransactionStatus(ctx context.Context, req *proto.GetTransactionStatusRequest) (*proto.GetTransactionStatusResponse, error) {
	if req.TransactionHash == "" {
		return &proto.GetTransactionStatusResponse{
			Status:       0,
			ErrorMessage: "transaction_hash is required",
		}, nil
	}

	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if status, err := r.getTransactionStatusFromKeeper(ctx, req.TransactionHash); err == nil {
			return status, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.getTransactionStatusFromDatabase(ctx, req.TransactionHash)
}

// GetPendingTransactions retrieves pending transactions
func (r *Repository) GetPendingTransactions(ctx context.Context, req *proto.GetPendingTransactionsRequest) (*proto.GetPendingTransactionsResponse, error) {
	limit, _ := utils.NormalizePagination(req.Limit, 0, utils.PaginationConfig{
		DefaultLimit: 100,
		MaxLimit:     1000,
	})

	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if txs, err := r.getPendingTransactionsFromKeeper(ctx, limit, req.Address); err == nil {
			return txs, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.getPendingTransactionsFromDatabase(ctx, req)
}

// EstimateTransactionFee estimates transaction fee
func (r *Repository) EstimateTransactionFee(ctx context.Context, req *proto.EstimateTransactionFeeRequest) (*proto.EstimateTransactionFeeResponse, error) {
	r.logger.Info("Estimating transaction fee in repository",
		logging.String("from", req.FromAddress),
		logging.String("to", req.ToAddress))

	// Normalize gas limit
	gasLimit := req.GasLimit
	if gasLimit <= 0 {
		gasLimit = 21000 // Default gas limit for simple transfer
	}

	// Simple fee estimation: gas_limit * gas_price
	// Default gas price: 1 USC per gas unit
	defaultGasPrice := "1"

	// Calculate total fee (simplified)
	// In production, this would query current gas prices from the network
	totalFee := fmt.Sprintf("%d", gasLimit) // Assuming 1 USC per gas unit

	return &proto.EstimateTransactionFeeResponse{
		GasPrice:                  defaultGasPrice,
		GasLimit:                  gasLimit,
		TotalFee:                  totalFee,
		UscFee:                    totalFee,
		EstimatedConfirmationTime: 5, // 5 seconds (simplified)
	}, nil
}

// Helper methods for TransactionKeeper interaction

// getSDKContext creates a sdk.Context from context.Context
// Uses shared utility to avoid code duplication
func (r *Repository) getSDKContext(ctx context.Context) (sdk.Context, error) {
	return utils.GetSDKContext(ctx, r.cosmosApp, r.logger)
}

// submitTransactionOnBlockchain submits a transaction using TransactionKeeper
func (r *Repository) submitTransactionOnBlockchain(ctx context.Context, req *proto.SubmitTransactionRequest) (string, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return "", repoerrors.NewBlockchainError("get_sdk_context", err)
	}

	// Generate real transaction hash using blocktypes helper
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, req.FromAddress, req.ToAddress, req.Amount, "transfer", req.Data, "")
	txID := txHash

	// Create transaction using TransactionKeeper
	transaction, err := r.cosmosApp.TransactionKeeper.CreateTransaction(
		sdkCtx,
		txHash,
		txID,
		req.FromAddress,
		req.ToAddress,
		req.Amount,
		"transfer", // Default transaction type
		req.Data,
		"", // memo not available
	)
	if err != nil {
		return "", repoerrors.WrapRepositoryError(repoerrors.ErrTransactionSubmissionFailed, err)
	}

	return transaction.Hash, nil
}

// getTransactionFromKeeper retrieves a transaction from TransactionKeeper
func (r *Repository) getTransactionFromKeeper(ctx context.Context, hash string) (*proto.GetTransactionResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, err
	}

	transaction, exists := r.cosmosApp.TransactionKeeper.GetTransaction(sdkCtx, hash)
	if !exists {
		return nil, repoerrors.NewNotFoundError("transaction", "")
	}

	return r.convertTransactionToProto(&transaction), nil
}

// getTransactionStatusFromKeeper retrieves transaction status from TransactionKeeper
func (r *Repository) getTransactionStatusFromKeeper(ctx context.Context, hash string) (*proto.GetTransactionStatusResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, err
	}

	transaction, exists := r.cosmosApp.TransactionKeeper.GetTransaction(sdkCtx, hash)
	if !exists {
		return &proto.GetTransactionStatusResponse{
			TransactionHash: hash,
			Status:          0, // Pending
			ErrorMessage:    "transaction not found",
		}, nil
	}

	// Convert status string to int32
	var status int32
	switch transaction.Status {
	case "pending":
		status = 0
	case "validated":
		status = 1
	case "executed":
		status = 1 // Confirmed
	case "failed":
		status = 2 // Failed
	case "cancelled":
		status = 2 // Failed
	default:
		status = 0
	}

	return &proto.GetTransactionStatusResponse{
		TransactionHash: hash,
		Status:          status,
		ErrorMessage:    transaction.ErrorMessage,
	}, nil
}

// getPendingTransactionsFromKeeper retrieves pending transactions from TransactionKeeper
func (r *Repository) getPendingTransactionsFromKeeper(ctx context.Context, limit int32, addressFilter string) (*proto.GetPendingTransactionsResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, err
	}

	allTransactions := r.cosmosApp.TransactionKeeper.GetAllTransactions(sdkCtx)

	// Filter pending transactions
	// Pre-allocate slice with capacity = limit for better performance (worst case: all filtered transactions fit in limit)
	pending := make([]transactiontypes.Transaction, 0, limit)
	for _, tx := range allTransactions {
		if tx.Status == "pending" {
			// Apply address filter if provided
			if addressFilter == "" || tx.FromAddress == addressFilter || tx.ToAddress == addressFilter {
				pending = append(pending, tx)
				if int32(len(pending)) >= limit {
					break
				}
			}
		}
	}

	// Convert to proto response
	transactions := make([]*proto.PendingTransaction, 0, len(pending))
	for _, tx := range pending {
		transactions = append(transactions, r.convertTransactionToPendingTransaction(&tx))
	}

	return &proto.GetPendingTransactionsResponse{
		Transactions: transactions,
		TotalCount:   int32(len(pending)),
		HasMore:      int32(len(pending)) >= limit,
	}, nil
}

// Database helper methods

// saveTransactionToDatabase saves a transaction to database for analytics (async)
func (r *Repository) saveTransactionToDatabase(ctx context.Context, txHash string, req *proto.SubmitTransactionRequest) error {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return fmt.Errorf("postgres connection not available")
	}
	if txHash == "" {
		return fmt.Errorf("transaction hash is empty")
	}

	query := `
		INSERT INTO usc_transaction_analytics (
			transaction_hash, from_address, to_address, amount, gas_price, gas_limit,
			timestamp, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (transaction_hash) DO UPDATE SET
			status = EXCLUDED.status
	`

	if _, err := postgres.ExecContext(ctx, query,
		txHash,
		req.FromAddress,
		req.ToAddress,
		req.Amount,
		req.GasPrice,
		req.GasLimit,
		time.Now(),
		"pending",
	); err != nil {
		return fmt.Errorf("failed to save transaction to database for analytics: %w", err)
	}

	return nil
}

// submitTransactionToDatabase submits a transaction to the database (fallback)
func (r *Repository) submitTransactionToDatabase(ctx context.Context, req *proto.SubmitTransactionRequest) (*proto.SubmitTransactionResponse, error) {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return &proto.SubmitTransactionResponse{
			Status:       2, // Failed
			ErrorMessage: "database not available",
		}, nil
	}

	query := `
		INSERT INTO usc_transaction_analytics (
			transaction_hash, from_address, to_address, amount, gas_price, gas_limit,
			timestamp, status, block_number, transaction_index
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (transaction_hash) DO UPDATE SET
			status = EXCLUDED.status,
			block_number = EXCLUDED.block_number,
			transaction_index = EXCLUDED.transaction_index
	`

	// Generate real transaction hash for database analytics
	// Use a simplified hash calculation since we don't have sdkCtx in database fallback
	// This is acceptable for analytics purposes
	dataStr := fmt.Sprintf("%s:%s:%s:%s", req.FromAddress, req.ToAddress, req.Amount, time.Now().Format(time.RFC3339))
	hashBytes := sha256.Sum256([]byte(dataStr))
	txHash := "0x" + hex.EncodeToString(hashBytes[:])

	_, err := postgres.ExecContext(ctx, query,
		txHash,
		req.FromAddress,
		req.ToAddress,
		req.Amount,
		req.GasPrice,
		req.GasLimit,
		time.Now(),
		"pending",
		nil, // block_number (null for pending)
		nil, // transaction_index (null for pending)
	)

	if err != nil {
		return &proto.SubmitTransactionResponse{
			Status:       2, // Failed
			ErrorMessage: fmt.Sprintf("failed to submit transaction: %v", err),
		}, nil
	}

	return &proto.SubmitTransactionResponse{
		TransactionHash: txHash,
		Status:          0, // Pending
		SubmittedAt:     timestamppb.New(time.Now()),
	}, nil
}

// getTransactionFromDatabase retrieves a transaction from the database
func (r *Repository) getTransactionFromDatabase(ctx context.Context, hash string) (*proto.GetTransactionResponse, error) {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return &proto.GetTransactionResponse{}, fmt.Errorf("database not available")
	}

	query := `
		SELECT transaction_hash, from_address, to_address, amount, gas_price, gas_limit,
			gas_used, timestamp, status, block_number, block_hash, transaction_index
		FROM usc_transaction_analytics
		WHERE transaction_hash = $1
	`

	var tx proto.GetTransactionResponse
	var createdAt time.Time
	var statusStr string
	var gasUsed sql.NullInt64
	var blockNumber sql.NullInt64
	var blockHash sql.NullString
	var txIndex sql.NullInt64

	err := postgres.QueryRowContext(ctx, query, hash).Scan(
		&tx.TransactionHash,
		&tx.FromAddress,
		&tx.ToAddress,
		&tx.Amount,
		&tx.GasPrice,
		&tx.GasLimit,
		&gasUsed,
		&createdAt,
		&statusStr,
		&blockNumber,
		&blockHash,
		&txIndex,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &proto.GetTransactionResponse{}, fmt.Errorf("transaction not found: %s", hash)
		}
		return &proto.GetTransactionResponse{}, fmt.Errorf("failed to query transaction: %w", err)
	}

	// Handle NULL values
	if gasUsed.Valid {
		tx.GasUsed = gasUsed.Int64
	}
	if blockNumber.Valid {
		tx.BlockNumber = blockNumber.Int64
	}
	if blockHash.Valid {
		tx.BlockHash = blockHash.String
	}

	if txIndex.Valid {
		tx.TransactionIndex = txIndex.Int64
	}

	// Convert status string to int32
	switch statusStr {
	case "pending":
		tx.Status = 0
	case "confirmed":
		tx.Status = 1
	case "failed":
		tx.Status = 2
	default:
		tx.Status = 0
	}

	tx.CreatedAt = timestamppb.New(createdAt)
	if tx.Status == 1 {
		tx.ConfirmedAt = timestamppb.New(createdAt)
	}

	return &tx, nil
}

// getTransactionStatusFromDatabase retrieves transaction status from the database
func (r *Repository) getTransactionStatusFromDatabase(ctx context.Context, hash string) (*proto.GetTransactionStatusResponse, error) {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return &proto.GetTransactionStatusResponse{}, fmt.Errorf("database not available")
	}

	query := `
		SELECT status, block_number, block_hash, timestamp
		FROM usc_transaction_analytics
		WHERE transaction_hash = $1
	`

	var statusStr string
	var blockNumber *int64
	var blockHash *string
	var createdAt *time.Time

	err := postgres.QueryRowContext(ctx, query, hash).Scan(
		&statusStr,
		&blockNumber,
		&blockHash,
		&createdAt,
	)

	if err != nil {
		return &proto.GetTransactionStatusResponse{
			TransactionHash: hash,
			Status:          0, // Pending
			ErrorMessage:    "transaction not found",
		}, nil
	}

	// Convert status string to int32
	var status int32
	switch statusStr {
	case "pending":
		status = 0
	case "confirmed":
		status = 1
	case "failed":
		status = 2
	default:
		status = 0
	}

	resp := &proto.GetTransactionStatusResponse{
		TransactionHash: hash,
		Status:          status,
	}

	if blockNumber != nil {
		resp.BlockNumber = *blockNumber
	}
	if blockHash != nil {
		resp.BlockHash = *blockHash
	}
	if createdAt != nil && status == 1 {
		resp.ConfirmedAt = timestamppb.New(*createdAt)
	}

	return resp, nil
}

// getPendingTransactionsFromDatabase retrieves pending transactions from the database
func (r *Repository) getPendingTransactionsFromDatabase(ctx context.Context, req *proto.GetPendingTransactionsRequest) (*proto.GetPendingTransactionsResponse, error) {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return nil, repoerrors.NewDatabaseError("get_pending_transactions", fmt.Errorf("database not available"))
	}

	limit, _ := utils.NormalizePagination(req.Limit, 0, utils.PaginationConfig{
		DefaultLimit: 100,
		MaxLimit:     1000,
	})

	query := `
		SELECT transaction_hash, from_address, to_address, amount, gas_price, gas_limit, timestamp
		FROM usc_transaction_analytics
		WHERE status = 'pending'
		AND ($1 = '' OR from_address = $1 OR to_address = $1)
		ORDER BY timestamp DESC
		LIMIT $2
	`

	rows, err := postgres.QueryContext(ctx, query, req.Address, limit)
	if err != nil {
		return nil, repoerrors.NewDatabaseError("query_pending_transactions", err)
	}
	defer rows.Close()

	// Pre-allocate slice with capacity = limit for better performance
	transactions := make([]*proto.PendingTransaction, 0, limit)
	for rows.Next() {
		var tx proto.PendingTransaction
		var submittedAt time.Time

		err := rows.Scan(
			&tx.TransactionHash,
			&tx.FromAddress,
			&tx.ToAddress,
			&tx.Amount,
			&tx.GasPrice,
			&tx.GasLimit,
			&submittedAt,
		)
		if err != nil {
			r.logger.Warn("Failed to scan pending transaction row", logging.Error(err))
			continue
		}

		tx.SubmittedAt = timestamppb.New(submittedAt)
		transactions = append(transactions, &tx)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		r.logger.Warn("Error during rows iteration", logging.Error(err))
		// Continue with partial results if iteration error occurs
	}

	return &proto.GetPendingTransactionsResponse{
		Transactions: transactions,
		TotalCount:   int32(len(transactions)),
		HasMore:      int32(len(transactions)) >= limit,
	}, nil
}

// Conversion helper methods

// convertTransactionToProto converts transactiontypes.Transaction to proto.GetTransactionResponse
func (r *Repository) convertTransactionToProto(tx *transactiontypes.Transaction) *proto.GetTransactionResponse {
	// Convert status string to int32
	var status int32
	switch tx.Status {
	case "pending":
		status = 0
	case "validated":
		status = 1
	case "executed":
		status = 1 // Confirmed
	case "failed":
		status = 2
	case "cancelled":
		status = 2
	default:
		status = 0
	}

	resp := &proto.GetTransactionResponse{
		TransactionHash: tx.Hash,
		FromAddress:     tx.FromAddress,
		ToAddress:       tx.ToAddress,
		Amount:          tx.Amount,
		GasLimit:        tx.GasLimit,
		GasUsed:         tx.GasUsed,
		Data:            tx.Data,
		Status:          status,
	}

	if tx.CreatedAt > 0 {
		resp.CreatedAt = timestamppb.New(time.Unix(tx.CreatedAt, 0))
	}
	if tx.ExecutedAt > 0 {
		resp.ConfirmedAt = timestamppb.New(time.Unix(tx.ExecutedAt, 0))
	}

	return resp
}

// convertTransactionToPendingTransaction converts transactiontypes.Transaction to proto.PendingTransaction
func (r *Repository) convertTransactionToPendingTransaction(tx *transactiontypes.Transaction) *proto.PendingTransaction {
	ptx := &proto.PendingTransaction{
		TransactionHash: tx.Hash,
		FromAddress:     tx.FromAddress,
		ToAddress:       tx.ToAddress,
		Amount:          tx.Amount,
		GasLimit:        tx.GasLimit,
	}

	if tx.CreatedAt > 0 {
		ptx.SubmittedAt = timestamppb.New(time.Unix(tx.CreatedAt, 0))
	}

	return ptx
}
