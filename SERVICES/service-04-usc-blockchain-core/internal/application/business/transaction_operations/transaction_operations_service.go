package transaction_operations

import (
	"context"
	"fmt"
	"time"

	"service-04/internal/application/repository/transaction_operations"
	"service-04/internal/application/utils"
	"service-04/internal/infrastructure/metrics"
	"service-04/internal/infrastructure/validation"
	proto "service-04/proto"

	"github.com/usc-platform/shared/logging"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/app"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service handles transaction operations business logic
type Service struct {
	repo              *transaction_operations.Repository
	cosmosApp         *app.USCApp
	blockchainStorage *storage.StateManager
	logger            *logging.Logger
	validator         *validation.Validator
	metrics           *metrics.MetricsService
}

// NewService creates a new transaction operations service
func NewService(
	repo *transaction_operations.Repository,
	cosmosApp *app.USCApp,
	blockchainStorage *storage.StateManager,
	logger *logging.Logger,
	validator *validation.Validator,
	metricsService *metrics.MetricsService,
) *Service {
	return &Service{
		repo:              repo,
		cosmosApp:         cosmosApp,
		blockchainStorage: blockchainStorage,
		logger:            logger,
		validator:         validator,
		metrics:           metricsService,
	}
}

// SubmitTransaction submits a new transaction
func (s *Service) SubmitTransaction(ctx context.Context, req *proto.SubmitTransactionRequest) (*proto.SubmitTransactionResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("submit_transaction", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Submitting transaction in business service",
		logging.String("correlation_id", correlationID),
		logging.String("from", req.FromAddress),
		logging.String("to", req.ToAddress))

	// Input validation using validator service
	if err := s.validator.ValidateWalletAddress(req.FromAddress); err != nil {
		s.logger.Error("From address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("from_address", req.FromAddress),
			logging.Error(err))
		s.metrics.RecordFailure("submit_transaction", "validation_error", map[string]string{
			"from_address": req.FromAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid from_address: %v", err)
	}

	if err := s.validator.ValidateWalletAddress(req.ToAddress); err != nil {
		s.logger.Error("To address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("to_address", req.ToAddress),
			logging.Error(err))
		s.metrics.RecordFailure("submit_transaction", "validation_error", map[string]string{
			"to_address": req.ToAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid to_address: %v", err)
	}

	if err := s.validator.ValidateAmount(req.Amount); err != nil {
		s.logger.Error("Amount validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("amount", req.Amount),
			logging.Error(err))
		s.metrics.RecordFailure("submit_transaction", "validation_error", map[string]string{
			"amount": req.Amount,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid amount: %v", err)
	}

	// Validate gas limit if provided
	if req.GasLimit > 0 {
		if err := s.validator.ValidateGasLimit(req.GasLimit); err != nil {
			s.logger.Error("Gas limit validation failed",
				logging.String("correlation_id", correlationID),
				logging.Int64("gas_limit", req.GasLimit),
				logging.Error(err))
			s.metrics.RecordFailure("submit_transaction", "validation_error", map[string]string{
				"gas_limit": fmt.Sprintf("%d", req.GasLimit),
			})
			return nil, status.Errorf(codes.InvalidArgument, "invalid gas_limit: %v", err)
		}
	}

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.SubmitTransaction(ctx, req)
	if err != nil {
		s.logger.Error("Failed to submit transaction in repository",
			logging.String("correlation_id", correlationID),
			logging.String("from_address", req.FromAddress),
			logging.String("to_address", req.ToAddress),
			logging.Error(err))
		s.metrics.RecordFailure("submit_transaction", "repository_error", map[string]string{
			"from_address": req.FromAddress,
			"to_address":   req.ToAddress,
		})
		return nil, status.Errorf(codes.Internal, "failed to submit transaction: %v", err)
	}

	// Record success metrics
	s.logger.Info("Transaction submitted successfully",
		logging.String("correlation_id", correlationID),
		logging.String("from_address", req.FromAddress),
		logging.String("to_address", req.ToAddress))
	s.metrics.RecordSuccess("submit_transaction", map[string]string{
		"from_address": req.FromAddress,
		"to_address":   req.ToAddress,
	})

	// Record blockchain-specific metric if transaction was submitted
	if response != nil && response.TransactionHash != "" {
		s.metrics.RecordTransactionSubmitted(response.TransactionHash, req.FromAddress, req.ToAddress)
	}

	return response, nil
}

// GetTransaction retrieves a transaction by hash
func (s *Service) GetTransaction(ctx context.Context, req *proto.GetTransactionRequest) (*proto.GetTransactionResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("get_transaction", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Getting transaction in business service",
		logging.String("correlation_id", correlationID),
		logging.String("hash", req.TransactionHash))

	// Input validation using validator service
	if err := s.validator.ValidateTransactionHash(req.TransactionHash); err != nil {
		s.logger.Error("Transaction hash validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("transaction_hash", req.TransactionHash),
			logging.Error(err))
		s.metrics.RecordFailure("get_transaction", "validation_error", map[string]string{
			"transaction_hash": req.TransactionHash,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid transaction_hash: %v", err)
	}

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.GetTransaction(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get transaction in repository",
			logging.String("correlation_id", correlationID),
			logging.String("transaction_hash", req.TransactionHash),
			logging.Error(err))
		s.metrics.RecordFailure("get_transaction", "repository_error", map[string]string{
			"transaction_hash": req.TransactionHash,
		})
		return nil, status.Errorf(codes.Internal, "failed to get transaction: %v", err)
	}

	// Record success metrics
	s.logger.Info("Transaction retrieved successfully",
		logging.String("correlation_id", correlationID),
		logging.String("transaction_hash", req.TransactionHash))
	s.metrics.RecordSuccess("get_transaction", map[string]string{
		"transaction_hash": req.TransactionHash,
	})

	return response, nil
}

// GetTransactionStatus retrieves transaction status
func (s *Service) GetTransactionStatus(ctx context.Context, req *proto.GetTransactionStatusRequest) (*proto.GetTransactionStatusResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("get_transaction_status", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Getting transaction status in business service",
		logging.String("correlation_id", correlationID),
		logging.String("hash", req.TransactionHash))

	// Input validation using validator service
	if err := s.validator.ValidateTransactionHash(req.TransactionHash); err != nil {
		s.logger.Error("Transaction hash validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("transaction_hash", req.TransactionHash),
			logging.Error(err))
		s.metrics.RecordFailure("get_transaction_status", "validation_error", map[string]string{
			"transaction_hash": req.TransactionHash,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid transaction_hash: %v", err)
	}

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.GetTransactionStatus(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get transaction status in repository",
			logging.String("correlation_id", correlationID),
			logging.String("transaction_hash", req.TransactionHash),
			logging.Error(err))
		s.metrics.RecordFailure("get_transaction_status", "repository_error", map[string]string{
			"transaction_hash": req.TransactionHash,
		})
		return nil, status.Errorf(codes.Internal, "failed to get transaction status: %v", err)
	}

	// Record success metrics
	s.logger.Info("Transaction status retrieved successfully",
		logging.String("correlation_id", correlationID),
		logging.String("transaction_hash", req.TransactionHash))
	s.metrics.RecordSuccess("get_transaction_status", map[string]string{
		"transaction_hash": req.TransactionHash,
	})

	return response, nil
}

// GetPendingTransactions retrieves pending transactions
func (s *Service) GetPendingTransactions(ctx context.Context, req *proto.GetPendingTransactionsRequest) (*proto.GetPendingTransactionsResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("get_pending_transactions", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Getting pending transactions in business service",
		logging.String("correlation_id", correlationID))

	// Business logic validation
	// Use helper function to normalize pagination (reduces duplicate code)
	req.Limit, _ = utils.NormalizePaginationWithDefaults(req.Limit, 0)

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.GetPendingTransactions(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get pending transactions in repository",
			logging.String("correlation_id", correlationID),
			logging.Error(err))
		s.metrics.RecordFailure("get_pending_transactions", "repository_error", map[string]string{})
		return nil, status.Errorf(codes.Internal, "failed to get pending transactions: %v", err)
	}

	// Record success metrics
	s.logger.Info("Pending transactions retrieved successfully",
		logging.String("correlation_id", correlationID))
	s.metrics.RecordSuccess("get_pending_transactions", map[string]string{})

	return response, nil
}

// EstimateTransactionFee estimates transaction fee
func (s *Service) EstimateTransactionFee(ctx context.Context, req *proto.EstimateTransactionFeeRequest) (*proto.EstimateTransactionFeeResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("estimate_transaction_fee", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Estimating transaction fee in business service",
		logging.String("correlation_id", correlationID),
		logging.String("from", req.FromAddress),
		logging.String("to", req.ToAddress))

	// Input validation using validator service
	if err := s.validator.ValidateWalletAddress(req.FromAddress); err != nil {
		s.logger.Error("From address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("from_address", req.FromAddress),
			logging.Error(err))
		s.metrics.RecordFailure("estimate_transaction_fee", "validation_error", map[string]string{
			"from_address": req.FromAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid from_address: %v", err)
	}

	// Validate gas limit if provided
	if req.GasLimit > 0 {
		if err := s.validator.ValidateGasLimit(req.GasLimit); err != nil {
			s.logger.Error("Gas limit validation failed",
				logging.String("correlation_id", correlationID),
				logging.Int64("gas_limit", req.GasLimit),
				logging.Error(err))
			s.metrics.RecordFailure("estimate_transaction_fee", "validation_error", map[string]string{
				"gas_limit": fmt.Sprintf("%d", req.GasLimit),
			})
			return nil, status.Errorf(codes.InvalidArgument, "invalid gas_limit: %v", err)
		}
	} else {
		req.GasLimit = 21000 // Default gas limit
	}

	// Call repository
	response, err := s.repo.EstimateTransactionFee(ctx, req)
	if err != nil {
		s.logger.Error("Failed to estimate transaction fee in repository",
			logging.String("correlation_id", correlationID),
			logging.String("from_address", req.FromAddress),
			logging.Error(err))
		s.metrics.RecordFailure("estimate_transaction_fee", "repository_error", map[string]string{
			"from_address": req.FromAddress,
		})
		return nil, status.Errorf(codes.Internal, "failed to estimate transaction fee: %v", err)
	}

	// Record success metrics
	s.logger.Info("Transaction fee estimated successfully",
		logging.String("correlation_id", correlationID),
		logging.String("from_address", req.FromAddress))
	s.metrics.RecordSuccess("estimate_transaction_fee", map[string]string{
		"from_address": req.FromAddress,
	})

	return response, nil
}
