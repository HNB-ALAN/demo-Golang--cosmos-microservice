package usc_coin_operations

import (
	"context"
	"time"

	"service-04/internal/application/repository/usc_coin_operations"
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

// Service handles USC coin operations business logic
type Service struct {
	repo              *usc_coin_operations.Repository
	cosmosApp         *app.USCApp
	blockchainStorage *storage.StateManager
	logger            *logging.Logger
	validator         *validation.Validator
	metrics           *metrics.MetricsService
}

// NewService creates a new USC coin operations service
func NewService(
	repo *usc_coin_operations.Repository,
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

// GetUSCBalance retrieves USC balance for an address
func (s *Service) GetUSCBalance(ctx context.Context, req *proto.GetWalletBalanceRequest) (*proto.GetWalletBalanceResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("get_usc_balance", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Getting USC balance in business service",
		logging.String("correlation_id", correlationID),
		logging.String("address", req.WalletAddress))

	// Input validation using validator service
	if err := s.validator.ValidateWalletAddress(req.WalletAddress); err != nil {
		s.logger.Error("Wallet address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("wallet_address", req.WalletAddress),
			logging.Error(err))
		s.metrics.RecordFailure("get_usc_balance", "validation_error", map[string]string{
			"wallet_address": req.WalletAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid wallet_address: %v", err)
	}

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.GetUSCBalance(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get USC balance in repository",
			logging.String("correlation_id", correlationID),
			logging.String("wallet_address", req.WalletAddress),
			logging.Error(err))
		s.metrics.RecordFailure("get_usc_balance", "repository_error", map[string]string{
			"wallet_address": req.WalletAddress,
		})
		return nil, status.Errorf(codes.Internal, "failed to get USC balance: %v", err)
	}

	// Record success metrics
	s.logger.Info("USC balance retrieved successfully",
		logging.String("correlation_id", correlationID),
		logging.String("wallet_address", req.WalletAddress))
	s.metrics.RecordSuccess("get_usc_balance", map[string]string{
		"wallet_address": req.WalletAddress,
	})

	return response, nil
}

// TransferUSC transfers USC between addresses
func (s *Service) TransferUSC(ctx context.Context, req *proto.TransferUSCBlockchainRequest) (*proto.TransferUSCBlockchainResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("transfer_usc", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Transferring USC in business service",
		logging.String("correlation_id", correlationID),
		logging.String("from", req.FromAddress),
		logging.String("to", req.ToAddress),
		logging.String("amount", req.Amount))

	// Input validation using validator service
	if err := s.validator.ValidateWalletAddress(req.FromAddress); err != nil {
		s.logger.Error("From address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("from_address", req.FromAddress),
			logging.Error(err))
		s.metrics.RecordFailure("transfer_usc", "validation_error", map[string]string{
			"from_address": req.FromAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid from_address: %v", err)
	}

	if err := s.validator.ValidateWalletAddress(req.ToAddress); err != nil {
		s.logger.Error("To address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("to_address", req.ToAddress),
			logging.Error(err))
		s.metrics.RecordFailure("transfer_usc", "validation_error", map[string]string{
			"to_address": req.ToAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid to_address: %v", err)
	}

	if err := s.validator.ValidateAmount(req.Amount); err != nil {
		s.logger.Error("Amount validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("amount", req.Amount),
			logging.Error(err))
		s.metrics.RecordFailure("transfer_usc", "validation_error", map[string]string{
			"amount": req.Amount,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid amount: %v", err)
	}

	// Business logic validation: from and to addresses must be different
	if req.FromAddress == req.ToAddress {
		s.logger.Error("Cannot transfer to self",
			logging.String("correlation_id", correlationID),
			logging.String("address", req.FromAddress))
		s.metrics.RecordFailure("transfer_usc", "validation_error", map[string]string{
			"from_address": req.FromAddress,
			"to_address":   req.ToAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "cannot transfer to self")
	}

	// Transfer USC - repository will handle blockchain submission or database fallback
	// Business service only validates and delegates to repository
	response, err := s.repo.TransferUSC(ctx, req)
	if err != nil {
		s.logger.Error("Failed to transfer USC in repository",
			logging.String("correlation_id", correlationID),
			logging.String("from_address", req.FromAddress),
			logging.String("to_address", req.ToAddress),
			logging.Error(err))
		s.metrics.RecordFailure("transfer_usc", "repository_error", map[string]string{
			"from_address": req.FromAddress,
			"to_address":   req.ToAddress,
		})
		return nil, status.Errorf(codes.Internal, "failed to transfer USC: %v", err)
	}

	// Record success metrics
	s.logger.Info("USC transfer completed successfully",
		logging.String("correlation_id", correlationID),
		logging.String("from_address", req.FromAddress),
		logging.String("to_address", req.ToAddress),
		logging.String("amount", req.Amount))
	s.metrics.RecordSuccess("transfer_usc", map[string]string{
		"from_address": req.FromAddress,
		"to_address":   req.ToAddress,
		"amount":       req.Amount,
	})

	// Record blockchain-specific metric if transfer was successful
	if response != nil && response.Success {
		s.metrics.RecordUSCTransfer(req.Amount, req.FromAddress, req.ToAddress)
	}

	return response, nil
}

// GetUSCSupply retrieves total USC supply
func (s *Service) GetUSCSupply(ctx context.Context) (*proto.GetUSCSupplyResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("get_usc_supply", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Getting USC supply in business service",
		logging.String("correlation_id", correlationID))

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.GetUSCSupply(ctx)
	if err != nil {
		s.logger.Error("Failed to get USC supply in repository",
			logging.String("correlation_id", correlationID),
			logging.Error(err))
		s.metrics.RecordFailure("get_usc_supply", "repository_error", map[string]string{})
		return nil, status.Errorf(codes.Internal, "failed to get USC supply: %v", err)
	}

	// Record success metrics
	s.logger.Info("USC supply retrieved successfully",
		logging.String("correlation_id", correlationID))
	s.metrics.RecordSuccess("get_usc_supply", map[string]string{})

	return response, nil
}

// GetTransactionHistory retrieves transaction history for an address
func (s *Service) GetTransactionHistory(ctx context.Context, req *proto.GetTransactionHistoryRequest) (*proto.GetTransactionHistoryResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("get_transaction_history", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Getting transaction history in business service",
		logging.String("correlation_id", correlationID),
		logging.String("address", req.WalletAddress))

	// Input validation using validator service
	if err := s.validator.ValidateWalletAddress(req.WalletAddress); err != nil {
		s.logger.Error("Wallet address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("wallet_address", req.WalletAddress),
			logging.Error(err))
		s.metrics.RecordFailure("get_transaction_history", "validation_error", map[string]string{
			"wallet_address": req.WalletAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid wallet_address: %v", err)
	}

	// Call repository
	response, err := s.repo.GetTransactionHistory(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get transaction history in repository",
			logging.String("correlation_id", correlationID),
			logging.String("wallet_address", req.WalletAddress),
			logging.Error(err))
		s.metrics.RecordFailure("get_transaction_history", "repository_error", map[string]string{
			"wallet_address": req.WalletAddress,
		})
		return nil, status.Errorf(codes.Internal, "failed to get transaction history: %v", err)
	}

	// Record success metrics
	s.logger.Info("Transaction history retrieved successfully",
		logging.String("correlation_id", correlationID),
		logging.String("wallet_address", req.WalletAddress))
	s.metrics.RecordSuccess("get_transaction_history", map[string]string{
		"wallet_address": req.WalletAddress,
	})

	return response, nil
}

// GetUSCTransactions retrieves USC-specific transactions
func (s *Service) GetUSCTransactions(ctx context.Context, req *proto.GetUSCTransactionsRequest) (*proto.GetUSCTransactionsResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("get_usc_transactions", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Getting USC transactions in business service",
		logging.String("correlation_id", correlationID),
		logging.String("address", req.WalletAddress))

	// Input validation using validator service
	if err := s.validator.ValidateWalletAddress(req.WalletAddress); err != nil {
		s.logger.Error("Wallet address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("wallet_address", req.WalletAddress),
			logging.Error(err))
		s.metrics.RecordFailure("get_usc_transactions", "validation_error", map[string]string{
			"wallet_address": req.WalletAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid wallet_address: %v", err)
	}

	// Call repository
	response, err := s.repo.GetUSCTransactions(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get USC transactions in repository",
			logging.String("correlation_id", correlationID),
			logging.String("wallet_address", req.WalletAddress),
			logging.Error(err))
		s.metrics.RecordFailure("get_usc_transactions", "repository_error", map[string]string{
			"wallet_address": req.WalletAddress,
		})
		return nil, status.Errorf(codes.Internal, "failed to get USC transactions: %v", err)
	}

	// Record success metrics
	s.logger.Info("USC transactions retrieved successfully",
		logging.String("correlation_id", correlationID),
		logging.String("wallet_address", req.WalletAddress))
	s.metrics.RecordSuccess("get_usc_transactions", map[string]string{
		"wallet_address": req.WalletAddress,
	})

	return response, nil
}
