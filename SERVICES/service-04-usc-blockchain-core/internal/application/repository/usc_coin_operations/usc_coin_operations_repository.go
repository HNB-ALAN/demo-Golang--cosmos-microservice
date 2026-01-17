package usc_coin_operations

import (
	"context"
	"fmt"
	"strings"
	"time"

	repoerrors "service-04/internal/application/repository"
	"service-04/internal/application/utils"
	"service-04/internal/infrastructure/database"
	proto "service-04/proto"

	// Cosmos SDK imports
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/app"
	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/usc_coin/v1/usc/usc_coin/v1"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/storage"
	"google.golang.org/protobuf/types/known/timestamppb"

	"crypto/sha256"
	"encoding/hex"

	"github.com/usc-platform/shared/logging"
)

// Repository handles USC coin operations data access
type Repository struct {
	db                *database.PostgreSQLManager
	cosmosApp         *app.USCApp
	blockchainStorage *storage.StateManager
	redisManager      *database.RedisManager
	logger            *logging.Logger
}

// NewRepository creates a new USC coin operations repository
func NewRepository(db *database.PostgreSQLManager, cosmosApp *app.USCApp, blockchainStorage *storage.StateManager, redisManager *database.RedisManager, logger *logging.Logger) *Repository {
	return &Repository{
		db:                db,
		cosmosApp:         cosmosApp,
		blockchainStorage: blockchainStorage,
		redisManager:      redisManager,
		logger:            logger,
	}
}

// GetUSCBalance retrieves USC balance for an address
func (r *Repository) GetUSCBalance(ctx context.Context, req *proto.GetWalletBalanceRequest) (*proto.GetWalletBalanceResponse, error) {
	// Priority 0: Redis cache
	cacheKey := fmt.Sprintf("usc:balance:%s", req.WalletAddress)
	if r.redisManager != nil {
		if cachedBalance, err := r.redisManager.Get(ctx, cacheKey); err == nil && cachedBalance != "" {
			return &proto.GetWalletBalanceResponse{
				Success:       true,
				WalletAddress: req.WalletAddress,
				Balance:       cachedBalance,
				Currency:      "USC",
			}, nil
		}
	}

	// Priority 1: Keeper (RocksDB)
	var balance string
	var err error
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if balance, err = r.getUSCBalanceFromBlockchain(ctx, req.WalletAddress); err == nil && balance != "" {
			// Cache in Redis (60s TTL)
			if r.redisManager != nil {
				_ = r.redisManager.Set(ctx, cacheKey, balance, 60*time.Second)
			}
			return &proto.GetWalletBalanceResponse{
				Success:       true,
				WalletAddress: req.WalletAddress,
				Balance:       balance,
				Currency:      "USC",
			}, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		balance = "0"
	} else {
		query := `SELECT COALESCE(SUM(amount), 0) as balance FROM usc_coin_analytics WHERE wallet_address = $1 AND status = 'confirmed'`
		if err = postgres.QueryRowContext(ctx, query, req.WalletAddress).Scan(&balance); err != nil {
			balance = "0"
		}
	}

	// Cache in Redis (even if from database)
	if r.redisManager != nil && balance != "" {
		_ = r.redisManager.Set(ctx, cacheKey, balance, 60*time.Second)
	}

	return &proto.GetWalletBalanceResponse{
		Success:       true,
		WalletAddress: req.WalletAddress,
		Balance:       balance,
		Currency:      "USC",
	}, nil
}

// TransferUSC transfers USC between addresses
func (r *Repository) TransferUSC(ctx context.Context, req *proto.TransferUSCBlockchainRequest) (*proto.TransferUSCBlockchainResponse, error) {
	// Priority 1: Transfer on Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.transferUSCOnBlockchain(ctx, req); err == nil {
			// Save to PostgreSQL for analytics (sync to ensure data persistence for financial data)
			if r.db != nil {
				correlationID := utils.GetCorrelationID(ctx)
				if err := r.saveTransferToDatabase(ctx, req, result); err != nil {
					r.logger.Error("Failed to save transfer analytics",
						logging.Error(err),
						logging.String("from_address", req.FromAddress),
						logging.String("to_address", req.ToAddress),
						logging.String("transaction_hash", result.TransactionHash),
						logging.String("correlation_id", correlationID))
					// Continue even if database save fails (keeper is primary)
				} else {
					r.logger.Info("Transfer analytics saved successfully",
						logging.String("from_address", req.FromAddress),
						logging.String("to_address", req.ToAddress),
						logging.String("transaction_hash", result.TransactionHash),
						logging.String("correlation_id", correlationID))
				}
			}
			// Invalidate cache
			r.invalidateBalanceCache(ctx, req.FromAddress, req.ToAddress)
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return nil, fmt.Errorf("database not available")
	}

	query := `INSERT INTO usc_coin_analytics (wallet_address, amount, transaction_type, status, timestamp, from_address, to_address, transaction_hash) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	// Generate transaction hash
	dataStr := fmt.Sprintf("%s:%s:%s:%s", req.FromAddress, req.ToAddress, req.Amount, time.Now().Format(time.RFC3339))
	hashBytes := sha256.Sum256([]byte(dataStr))
	txHash := "0x" + hex.EncodeToString(hashBytes[:])

	// Insert debit and credit transactions
	if _, err := postgres.ExecContext(ctx, query, req.FromAddress, "-"+req.Amount, "transfer_out", "confirmed", time.Now(), req.FromAddress, req.ToAddress, txHash); err != nil {
		return nil, err
	}
	if _, err := postgres.ExecContext(ctx, query, req.ToAddress, req.Amount, "transfer_in", "confirmed", time.Now(), req.FromAddress, req.ToAddress, txHash); err != nil {
		return nil, err
	}

	// Invalidate cache
	r.invalidateBalanceCache(ctx, req.FromAddress, req.ToAddress)

	return &proto.TransferUSCBlockchainResponse{
		Success:         true,
		TransactionHash: txHash,
		Status:          0, // Pending (database fallback)
		ErrorMessage:    "transaction logged in database",
	}, nil
}

// invalidateBalanceCache invalidates balance and supply cache
func (r *Repository) invalidateBalanceCache(ctx context.Context, fromAddress, toAddress string) {
	if r.redisManager == nil {
		return
	}
	_, _ = r.redisManager.Delete(ctx, fmt.Sprintf("usc:balance:%s", fromAddress))
	_, _ = r.redisManager.Delete(ctx, fmt.Sprintf("usc:balance:%s", toAddress))
	_, _ = r.redisManager.Delete(ctx, "usc:supply:total")
}

// GetUSCSupply retrieves total USC supply
func (r *Repository) GetUSCSupply(ctx context.Context) (*proto.GetUSCSupplyResponse, error) {
	// Priority 0: Redis cache
	cacheKey := "usc:supply:total"
	if r.redisManager != nil {
		if cachedSupply, err := r.redisManager.Get(ctx, cacheKey); err == nil && cachedSupply != "" {
			return &proto.GetUSCSupplyResponse{
				Success:     true,
				TotalSupply: cachedSupply,
				Currency:    "USC",
			}, nil
		}
	}

	// Priority 1: Keeper (RocksDB)
	var totalSupply string
	var err error
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if sdkCtx, err := r.getSDKContext(ctx); err == nil {
			if supply, err := r.cosmosApp.USCCoinKeeper.GetTotalSupply(sdkCtx); err == nil && supply != "" {
				totalSupply = supply
				// Cache in Redis (60s TTL)
				if r.redisManager != nil {
					_ = r.redisManager.Set(ctx, cacheKey, totalSupply, 60*time.Second)
				}
				return &proto.GetUSCSupplyResponse{
					Success:     true,
					TotalSupply: totalSupply,
					Currency:    "USC",
				}, nil
			}
		}
	}

	// Priority 2: PostgreSQL (fallback)
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return nil, repoerrors.NewRepositoryError(repoerrors.ErrDatabaseUnavailable)
	}

	query := `SELECT COALESCE(SUM(amount), 0) as total_supply FROM usc_coin_analytics WHERE transaction_type = 'mint' AND status = 'confirmed'`
	if err = postgres.QueryRowContext(ctx, query).Scan(&totalSupply); err != nil {
		return nil, err
	}

	// Cache in Redis (even if from database)
	if r.redisManager != nil && totalSupply != "" {
		_ = r.redisManager.Set(ctx, cacheKey, totalSupply, 60*time.Second)
	}

	return &proto.GetUSCSupplyResponse{
		Success:     true,
		TotalSupply: totalSupply,
		Currency:    "USC",
	}, nil
}

// GetTransactionHistory retrieves transaction history for an address
func (r *Repository) GetTransactionHistory(ctx context.Context, req *proto.GetTransactionHistoryRequest) (*proto.GetTransactionHistoryResponse, error) {
	if req.WalletAddress == "" {
		return &proto.GetTransactionHistoryResponse{Success: false, ErrorMessage: "wallet address is required", ErrorCode: 400}, nil
	}

	limit, offset := utils.NormalizePagination(req.Limit, req.Offset, utils.PaginationConfig{
		DefaultLimit: 50,
		MaxLimit:     100,
	})

	transactions, totalCount, err := r.queryTransactions(ctx, req.WalletAddress, req.StartDate, req.EndDate, req.TransactionType, "", limit, offset)
	if err != nil {
		return &proto.GetTransactionHistoryResponse{Success: false, ErrorMessage: err.Error(), ErrorCode: 500}, nil
	}

	hasMore, nextOffset := r.calculatePagination(int32(len(transactions)), limit, offset, totalCount)

	return &proto.GetTransactionHistoryResponse{
		Success:       true,
		Transactions:  r.convertToHistoryEntries(transactions),
		TotalCount:    totalCount,
		HasMore:       hasMore,
		NextOffset:    nextOffset,
		TotalSent:     "0",
		TotalReceived: "0",
		TotalFees:     "0",
	}, nil
}

// GetUSCTransactions retrieves USC-specific transactions
func (r *Repository) GetUSCTransactions(ctx context.Context, req *proto.GetUSCTransactionsRequest) (*proto.GetUSCTransactionsResponse, error) {
	if req.WalletAddress == "" {
		return &proto.GetUSCTransactionsResponse{Success: false, ErrorMessage: "wallet address is required", ErrorCode: 400}, nil
	}

	limit, offset := utils.NormalizePagination(req.Limit, req.Offset, utils.PaginationConfig{
		DefaultLimit: 50,
		MaxLimit:     100,
	})

	transactions, totalCount, err := r.queryTransactions(ctx, req.WalletAddress, req.StartDate, req.EndDate, r.mapTransactionType(req.TransactionType), req.Status, limit, offset)
	if err != nil {
		return &proto.GetUSCTransactionsResponse{Success: false, ErrorMessage: err.Error(), ErrorCode: 500}, nil
	}

	hasMore, nextOffset := r.calculatePagination(int32(len(transactions)), limit, offset, totalCount)

	return &proto.GetUSCTransactionsResponse{
		Success:       true,
		Transactions:  r.convertToUSCEntries(transactions),
		TotalCount:    totalCount,
		HasMore:       hasMore,
		NextOffset:    nextOffset,
		TotalSent:     "0",
		TotalReceived: "0",
		TotalFees:     "0",
	}, nil
}

// getUSCBalanceFromBlockchain retrieves USC balance from Cosmos SDK blockchain using USCCoinKeeper
func (r *Repository) getUSCBalanceFromBlockchain(ctx context.Context, address string) (string, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return "", repoerrors.NewBlockchainError("get_sdk_context", err)
	}

	balance, err := r.cosmosApp.USCCoinKeeper.GetBalance(sdkCtx, address)
	if err != nil {
		// Balance not found is not an error - return "0"
		if strings.Contains(err.Error(), "balance not found") {
			return "0", nil
		}
		return "", repoerrors.NewDatabaseError("query_balance", err)
	}

	return balance.Amount, nil
}

// getSDKContext creates a sdk.Context from context.Context
// Uses shared utility to avoid code duplication
func (r *Repository) getSDKContext(ctx context.Context) (sdk.Context, error) {
	return utils.GetSDKContext(ctx, r.cosmosApp, r.logger)
}

// transferUSCOnBlockchain transfers USC on Cosmos SDK blockchain
func (r *Repository) transferUSCOnBlockchain(ctx context.Context, req *proto.TransferUSCBlockchainRequest) (*proto.TransferUSCBlockchainResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, repoerrors.NewBlockchainError("get_sdk_context", err)
	}

	// Convert amount string to sdk.Coin (format: "100usc")
	coin, err := sdk.ParseCoinNormalized(req.Amount + "usc")
	if err != nil {
		return nil, repoerrors.NewValidationError("amount", fmt.Sprintf("invalid amount: %s", req.Amount))
	}

	// Execute transfer on blockchain
	resp, err := r.cosmosApp.USCCoinKeeper.TransferUSC(sdkCtx, &blockchainproto.MsgTransferUSC{
		FromAddress: req.FromAddress,
		ToAddress:   req.ToAddress,
		Amount:      &coin,
		GasLimit:    uint64(req.GasLimit),
		Memo:        req.Data,
	})
	if err != nil {
		return nil, repoerrors.WrapRepositoryError(repoerrors.ErrUSCTransferFailed, err)
	}

	return &proto.TransferUSCBlockchainResponse{
		Success:         true,
		TransactionHash: resp.TransactionHash,
		Status:          1, // Confirmed (on blockchain)
		ErrorMessage:    "",
	}, nil
}

// saveTransferToDatabase saves transfer to database for analytics (sync for financial data persistence)
func (r *Repository) saveTransferToDatabase(ctx context.Context, req *proto.TransferUSCBlockchainRequest, resp *proto.TransferUSCBlockchainResponse) error {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return fmt.Errorf("postgres connection not available")
	}
	if resp == nil {
		return fmt.Errorf("response is nil")
	}

	query := `INSERT INTO usc_coin_analytics (wallet_address, amount, transaction_type, status, timestamp, from_address, to_address, transaction_hash) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	// Insert debit transaction (transfer_out)
	if _, err := postgres.ExecContext(ctx, query, req.FromAddress, "-"+req.Amount, "transfer_out", "confirmed", time.Now(), req.FromAddress, req.ToAddress, resp.TransactionHash); err != nil {
		return fmt.Errorf("failed to save transfer_out to analytics: %w", err)
	}

	// Insert credit transaction (transfer_in)
	if _, err := postgres.ExecContext(ctx, query, req.ToAddress, req.Amount, "transfer_in", "confirmed", time.Now(), req.FromAddress, req.ToAddress, resp.TransactionHash); err != nil {
		return fmt.Errorf("failed to save transfer_in to analytics: %w", err)
	}

	return nil
}

// Helper functions for transaction queries

type txRow struct {
	hash        string
	from        string
	to          string
	amount      string
	txType      string
	status      string
	timestamp   time.Time
	blockNumber int64
	blockHash   string
}

// mapTransactionType maps request transaction type to database type
func (r *Repository) mapTransactionType(txType string) string {
	switch txType {
	case "sent":
		return "transfer_out"
	case "received":
		return "transfer_in"
	case "staked":
		return "stake"
	case "unstaked":
		return "unstake"
	default:
		return txType
	}
}

// queryTransactions executes transaction query with filters
func (r *Repository) queryTransactions(ctx context.Context, address string, startDate, endDate *timestamppb.Timestamp, txType, status string, limit, offset int32) ([]txRow, int32, error) {
	// Use helper function to get PostgreSQL connection (reduces duplicate code)
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return nil, 0, fmt.Errorf("database not available")
	}

	// Build query using strings.Builder for better performance
	var queryBuilder strings.Builder
	queryBuilder.WriteString(`SELECT transaction_hash, from_address, to_address, amount, transaction_type, status, timestamp, 
		COALESCE(block_number, 0) as block_number, COALESCE(block_hash, '') as block_hash,
		COUNT(*) OVER() as total_count
		FROM usc_coin_analytics WHERE wallet_address = $1`)
	args := []interface{}{address}
	argIndex := 2

	if startDate != nil {
		queryBuilder.WriteString(fmt.Sprintf(" AND timestamp >= $%d", argIndex))
		args = append(args, startDate.AsTime())
		argIndex++
	}
	if endDate != nil {
		queryBuilder.WriteString(fmt.Sprintf(" AND timestamp <= $%d", argIndex))
		args = append(args, endDate.AsTime())
		argIndex++
	}
	if txType != "" && txType != "all" {
		queryBuilder.WriteString(fmt.Sprintf(" AND transaction_type = $%d", argIndex))
		args = append(args, txType)
		argIndex++
	}
	if status != "" && status != "all" {
		queryBuilder.WriteString(fmt.Sprintf(" AND status = $%d", argIndex))
		args = append(args, status)
		argIndex++
	}

	queryBuilder.WriteString(" ORDER BY timestamp DESC")
	queryBuilder.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1))
	args = append(args, limit, offset)
	query := queryBuilder.String()

	// Execute query
	rows, err := postgres.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, repoerrors.NewDatabaseError("query_transaction_history", err)
	}
	defer rows.Close()

	// Pre-allocate slice with capacity = limit for better performance
	transactions := make([]txRow, 0, limit)
	var totalCount int32
	for rows.Next() {
		var row txRow
		if err := rows.Scan(&row.hash, &row.from, &row.to, &row.amount, &row.txType, &row.status, &row.timestamp, &row.blockNumber, &row.blockHash, &totalCount); err != nil {
			r.logger.Warn("Failed to scan row", logging.Error(err))
			continue
		}
		transactions = append(transactions, row)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		r.logger.Warn("Error during rows iteration", logging.Error(err))
		// Use count from results if iteration error occurs
		if totalCount == 0 {
			totalCount = int32(len(transactions))
		}
	}

	// Fallback to result count if totalCount is 0 (shouldn't happen, but safety check)
	if totalCount == 0 {
		totalCount = int32(len(transactions))
	}

	return transactions, totalCount, nil
}

// convertToHistoryEntries converts txRow to TransactionHistoryEntry
func (r *Repository) convertToHistoryEntries(rows []txRow) []*proto.TransactionHistoryEntry {
	entries := make([]*proto.TransactionHistoryEntry, 0, len(rows))
	for _, row := range rows {
		entries = append(entries, &proto.TransactionHistoryEntry{
			TransactionHash: row.hash,
			FromAddress:     row.from,
			ToAddress:       row.to,
			Amount:          strings.TrimPrefix(row.amount, "-"),
			Type:            r.convertTxType(row.txType),
			Status:          r.convertStatus(row.status),
			BlockNumber:     row.blockNumber,
			BlockHash:       row.blockHash,
			CreatedAt:       timestamppb.New(row.timestamp),
		})
	}
	return entries
}

// convertToUSCEntries converts txRow to USCTransactionEntry
func (r *Repository) convertToUSCEntries(rows []txRow) []*proto.USCTransactionEntry {
	entries := make([]*proto.USCTransactionEntry, 0, len(rows))
	for _, row := range rows {
		entries = append(entries, &proto.USCTransactionEntry{
			TransactionHash: row.hash,
			FromAddress:     row.from,
			ToAddress:       row.to,
			Amount:          strings.TrimPrefix(row.amount, "-"),
			Type:            r.convertTxType(row.txType),
			Status:          r.convertStatus(row.status),
			BlockNumber:     row.blockNumber,
			BlockHash:       row.blockHash,
			CreatedAt:       timestamppb.New(row.timestamp),
		})
	}
	return entries
}

// convertStatus converts status string to int32
func (r *Repository) convertStatus(status string) int32 {
	switch status {
	case "confirmed":
		return 1
	case "failed":
		return 2
	default:
		return 0
	}
}

// convertTxType converts transaction type string to int32
func (r *Repository) convertTxType(txType string) int32 {
	switch txType {
	case "stake":
		return 1
	case "unstake":
		return 2
	case "reward":
		return 3
	case "contract":
		return 4
	default:
		return 0
	}
}

// calculatePagination calculates pagination metadata
func (r *Repository) calculatePagination(returned, limit, offset, total int32) (bool, int32) {
	hasMore := returned == limit && (offset+limit) < total
	nextOffset := offset + limit
	if !hasMore {
		nextOffset = 0
	}
	return hasMore, nextOffset
}
