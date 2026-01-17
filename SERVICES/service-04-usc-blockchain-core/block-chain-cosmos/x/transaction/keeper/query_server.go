package keeper

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	query "github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	sdk "github.com/cosmos/cosmos-sdk/types"
	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/transaction/v1/usc/transaction/v1"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/transaction/types"
)

// QueryServer defines the query server interface using blockchain-proto types
type QueryServer interface {
	QueryTransaction(context.Context, *blockchainproto.QueryTransactionRequest) (*blockchainproto.QueryTransactionResponse, error)
	QueryTransactions(context.Context, *blockchainproto.QueryTransactionsRequest) (*blockchainproto.QueryTransactionsResponse, error)
	QueryTransactionStats(context.Context, *blockchainproto.QueryTransactionStatsRequest) (*blockchainproto.QueryTransactionStatsResponse, error)
}

// queryServer implements the QueryServer interface
type queryServer struct {
	Keeper
}

// NewQueryServerImpl returns an implementation of the QueryServer interface
func NewQueryServerImpl(keeper Keeper) QueryServer {
	return &queryServer{Keeper: keeper}
}

// QueryTransaction returns a specific transaction using blockchain-proto types
func (k queryServer) QueryTransaction(ctx context.Context, req *blockchainproto.QueryTransactionRequest) (*blockchainproto.QueryTransactionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	var transaction types.Transaction
	var found bool

	if req.TransactionHash != "" {
		transaction, found = k.GetTransaction(sdkCtx, req.TransactionHash)
	} else if req.TransactionId != "" {
		// Find transaction by ID
		store := sdkCtx.KVStore(k.storeKey)
		hashBytes := store.Get(types.GetTransactionIDKey(req.TransactionId))
		if hashBytes != nil {
			transaction, found = k.GetTransaction(sdkCtx, string(hashBytes))
		}
	} else {
		return nil, status.Error(codes.InvalidArgument, "transaction hash or ID must be provided")
	}

	if !found {
		return nil, status.Error(codes.NotFound, "transaction not found")
	}

	// Convert to blockchain-proto Transaction type
	// Convert string amount to Coin
	amount, err := sdk.ParseCoinNormalized(transaction.Amount)
	if err != nil {
		// If parsing fails, create a default coin
		amount = sdk.NewCoin("usc", math.NewInt(0))
	}

	// Convert string types to blockchain-proto enums
	var txType blockchainproto.TransactionType
	switch transaction.TransactionType {
	case "transfer":
		txType = blockchainproto.TransactionType_TRANSACTION_TYPE_TRANSFER
	case "mint":
		txType = blockchainproto.TransactionType_TRANSACTION_TYPE_MINT
	case "burn":
		txType = blockchainproto.TransactionType_TRANSACTION_TYPE_BURN
	default:
		txType = blockchainproto.TransactionType_TRANSACTION_TYPE_UNSPECIFIED
	}

	var txStatus blockchainproto.TransactionStatus
	switch transaction.Status {
	case "pending":
		txStatus = blockchainproto.TransactionStatus_TRANSACTION_STATUS_PENDING
	case "validated":
		txStatus = blockchainproto.TransactionStatus_TRANSACTION_STATUS_VALIDATED
	case "executed":
		txStatus = blockchainproto.TransactionStatus_TRANSACTION_STATUS_EXECUTED
	case "failed":
		txStatus = blockchainproto.TransactionStatus_TRANSACTION_STATUS_FAILED
	case "cancelled":
		txStatus = blockchainproto.TransactionStatus_TRANSACTION_STATUS_CANCELLED
	default:
		txStatus = blockchainproto.TransactionStatus_TRANSACTION_STATUS_UNSPECIFIED
	}

	// Convert timestamp
	createdAt := time.Unix(transaction.CreatedAt, 0)
	createdAtProto := timestamppb.New(createdAt)

	blockchainTransaction := &blockchainproto.Transaction{
		Id:              transaction.ID,
		Hash:            transaction.Hash,
		FromAddress:     transaction.FromAddress,
		ToAddress:       transaction.ToAddress,
		Amount:          &amount,
		TransactionType: txType,
		Data:            &blockchainproto.TransactionData{},
		Memo:            transaction.Memo,
		Status:          txStatus,
		CreatedAt:       createdAtProto,
	}

	return &blockchainproto.QueryTransactionResponse{
		Transaction: blockchainTransaction,
	}, nil
}

// QueryTransactions returns multiple transactions with pagination using blockchain-proto types
func (k queryServer) QueryTransactions(ctx context.Context, req *blockchainproto.QueryTransactionsRequest) (*blockchainproto.QueryTransactionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	store := sdkCtx.KVStore(k.storeKey)
	transactionStore := prefix.NewStore(store, []byte(types.TransactionKeyPrefix))

	var transactions []types.Transaction
	pageRes, err := query.Paginate(transactionStore, req.Pagination, func(key []byte, value []byte) error {
		var transaction types.Transaction
		if err := k.cdc.Unmarshal(value, &transaction); err != nil {
			return err
		}

		// Apply filters
		if req.FromAddress != "" && transaction.FromAddress != req.FromAddress {
			return nil
		}
		if req.ToAddress != "" && transaction.ToAddress != req.ToAddress {
			return nil
		}
		if req.TransactionType != blockchainproto.TransactionType_TRANSACTION_TYPE_UNSPECIFIED {
			// Convert string to enum for comparison
			var txType blockchainproto.TransactionType
			switch transaction.TransactionType {
			case "transfer":
				txType = blockchainproto.TransactionType_TRANSACTION_TYPE_TRANSFER
			case "mint":
				txType = blockchainproto.TransactionType_TRANSACTION_TYPE_MINT
			case "burn":
				txType = blockchainproto.TransactionType_TRANSACTION_TYPE_BURN
			default:
				txType = blockchainproto.TransactionType_TRANSACTION_TYPE_UNSPECIFIED
			}
			if txType != req.TransactionType {
				return nil
			}
		}
		if req.Status != blockchainproto.TransactionStatus_TRANSACTION_STATUS_UNSPECIFIED {
			// Convert string to enum for comparison
			var txStatus blockchainproto.TransactionStatus
			switch transaction.Status {
			case "pending":
				txStatus = blockchainproto.TransactionStatus_TRANSACTION_STATUS_PENDING
			case "validated":
				txStatus = blockchainproto.TransactionStatus_TRANSACTION_STATUS_VALIDATED
			case "executed":
				txStatus = blockchainproto.TransactionStatus_TRANSACTION_STATUS_EXECUTED
			case "failed":
				txStatus = blockchainproto.TransactionStatus_TRANSACTION_STATUS_FAILED
			case "cancelled":
				txStatus = blockchainproto.TransactionStatus_TRANSACTION_STATUS_CANCELLED
			default:
				txStatus = blockchainproto.TransactionStatus_TRANSACTION_STATUS_UNSPECIFIED
			}
			if txStatus != req.Status {
				return nil
			}
		}
		if req.StartTime != nil && transaction.CreatedAt < req.StartTime.Seconds {
			return nil
		}
		if req.EndTime != nil && transaction.CreatedAt > req.EndTime.Seconds {
			return nil
		}

		transactions = append(transactions, transaction)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert to blockchain-proto Transaction types
	var blockchainTransactions []*blockchainproto.Transaction
	for _, transaction := range transactions {
		// Convert string amount to Coin
		amount, err := sdk.ParseCoinNormalized(transaction.Amount)
		if err != nil {
			// If parsing fails, create a default coin
			amount = sdk.NewCoin("usc", math.NewInt(0))
		}

		// Convert string types to blockchain-proto enums
		var txType blockchainproto.TransactionType
		switch transaction.TransactionType {
		case "transfer":
			txType = blockchainproto.TransactionType_TRANSACTION_TYPE_TRANSFER
		case "mint":
			txType = blockchainproto.TransactionType_TRANSACTION_TYPE_MINT
		case "burn":
			txType = blockchainproto.TransactionType_TRANSACTION_TYPE_BURN
		default:
			txType = blockchainproto.TransactionType_TRANSACTION_TYPE_UNSPECIFIED
		}

		var txStatus blockchainproto.TransactionStatus
		switch transaction.Status {
		case "pending":
			txStatus = blockchainproto.TransactionStatus_TRANSACTION_STATUS_PENDING
		case "validated":
			txStatus = blockchainproto.TransactionStatus_TRANSACTION_STATUS_VALIDATED
		case "executed":
			txStatus = blockchainproto.TransactionStatus_TRANSACTION_STATUS_EXECUTED
		case "failed":
			txStatus = blockchainproto.TransactionStatus_TRANSACTION_STATUS_FAILED
		case "cancelled":
			txStatus = blockchainproto.TransactionStatus_TRANSACTION_STATUS_CANCELLED
		default:
			txStatus = blockchainproto.TransactionStatus_TRANSACTION_STATUS_UNSPECIFIED
		}

		// Convert timestamp
		createdAt := time.Unix(transaction.CreatedAt, 0)
		createdAtProto := timestamppb.New(createdAt)

		blockchainTransaction := &blockchainproto.Transaction{
			Id:              transaction.ID,
			Hash:            transaction.Hash,
			FromAddress:     transaction.FromAddress,
			ToAddress:       transaction.ToAddress,
			Amount:          &amount,
			TransactionType: txType,
			Data:            &blockchainproto.TransactionData{},
			Memo:            transaction.Memo,
			Status:          txStatus,
			CreatedAt:       createdAtProto,
		}
		blockchainTransactions = append(blockchainTransactions, blockchainTransaction)
	}

	return &blockchainproto.QueryTransactionsResponse{
		Transactions: blockchainTransactions,
		Pagination:   pageRes,
	}, nil
}

// QueryTransactionStats returns transaction statistics using blockchain-proto types
func (k queryServer) QueryTransactionStats(ctx context.Context, req *blockchainproto.QueryTransactionStatsRequest) (*blockchainproto.QueryTransactionStatsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	stats := k.GetTransactionStats(sdkCtx)

	// Apply time filter if provided
	if req.StartTime != nil || req.EndTime != nil {
		// Filter transactions by time range
		allTransactions := k.GetAllTransactions(sdkCtx)
		var filteredTransactions []types.Transaction

		for _, tx := range allTransactions {
			if req.StartTime != nil && tx.CreatedAt < req.StartTime.Seconds {
				continue
			}
			if req.EndTime != nil && tx.CreatedAt > req.EndTime.Seconds {
				continue
			}
			if req.TransactionType != blockchainproto.TransactionType_TRANSACTION_TYPE_UNSPECIFIED {
				// Convert string to enum for comparison
				var txType blockchainproto.TransactionType
				switch tx.TransactionType {
				case "transfer":
					txType = blockchainproto.TransactionType_TRANSACTION_TYPE_TRANSFER
				case "mint":
					txType = blockchainproto.TransactionType_TRANSACTION_TYPE_MINT
				case "burn":
					txType = blockchainproto.TransactionType_TRANSACTION_TYPE_BURN
				default:
					txType = blockchainproto.TransactionType_TRANSACTION_TYPE_UNSPECIFIED
				}
				if txType != req.TransactionType {
					continue
				}
			}
			filteredTransactions = append(filteredTransactions, tx)
		}

		// Recalculate stats for filtered transactions
		stats = types.TransactionStats{
			TotalTransactions:       int64(len(filteredTransactions)),
			PendingTransactions:     0,
			ValidatedTransactions:   0,
			ExecutedTransactions:    0,
			FailedTransactions:      0,
			CancelledTransactions:   0,
			TotalVolume:             "0",
			AverageTransactionValue: "0",
			SuccessRate:             "0.00",
			AverageExecutionTime:    0,
			CurrentHeight:           sdkCtx.BlockHeight(),
			LastTransactionTime:     time.Now().Format(time.RFC3339),
		}

		var totalVolume float64
		var executedCount int64

		for _, tx := range filteredTransactions {
			if amount, err := strconv.ParseFloat(tx.Amount, 64); err == nil {
				totalVolume += amount
			}

			switch tx.Status {
			case "pending":
				stats.PendingTransactions++
			case "validated":
				stats.ValidatedTransactions++
			case "executed":
				stats.ExecutedTransactions++
				executedCount++
			case "failed":
				stats.FailedTransactions++
			case "cancelled":
				stats.CancelledTransactions++
			}
		}

		stats.TotalVolume = fmt.Sprintf("%.2f", totalVolume)
		if len(filteredTransactions) > 0 {
			stats.AverageTransactionValue = fmt.Sprintf("%.2f", totalVolume/float64(len(filteredTransactions)))
		}
		if stats.TotalTransactions > 0 {
			successRate := float64(executedCount) / float64(stats.TotalTransactions) * 100
			stats.SuccessRate = fmt.Sprintf("%.2f", successRate)
		}
	}

	// Convert to blockchain-proto TransactionStats type
	// Parse success rate from string to float64
	successRate, _ := strconv.ParseFloat(stats.SuccessRate, 64)

	// Parse last transaction time
	var lastTransactionTime *timestamppb.Timestamp
	if stats.LastTransactionTime != "" {
		if t, err := time.Parse(time.RFC3339, stats.LastTransactionTime); err == nil {
			lastTransactionTime = timestamppb.New(t)
		}
	}

	// Convert string amounts to Coin types
	var totalVolume *sdk.Coin
	if stats.TotalVolume != "" {
		if coin, err := sdk.ParseCoinNormalized(stats.TotalVolume); err == nil {
			totalVolume = &coin
		}
	}

	var avgTransactionValue *sdk.Coin
	if stats.AverageTransactionValue != "" {
		if coin, err := sdk.ParseCoinNormalized(stats.AverageTransactionValue); err == nil {
			avgTransactionValue = &coin
		}
	}

	blockchainStats := &blockchainproto.TransactionStats{
		TotalTransactions:       stats.TotalTransactions,
		PendingTransactions:     stats.PendingTransactions,
		ValidatedTransactions:   stats.ValidatedTransactions,
		ExecutedTransactions:    stats.ExecutedTransactions,
		FailedTransactions:      stats.FailedTransactions,
		CancelledTransactions:   stats.CancelledTransactions,
		TotalVolume:             totalVolume,
		AverageTransactionValue: avgTransactionValue,
		SuccessRate:             successRate,
		CurrentHeight:           stats.CurrentHeight,
		LastTransactionTime:     lastTransactionTime,
	}

	return &blockchainproto.QueryTransactionStatsResponse{
		Stats: blockchainStats,
	}, nil
}

// Note: Custom query types removed as they are replaced by blockchain-proto query types
// The blockchain-proto interface provides QueryTransaction, QueryTransactions, and QueryTransactionStats
