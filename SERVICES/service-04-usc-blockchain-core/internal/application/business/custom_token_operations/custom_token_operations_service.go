package custom_token_operations

import (
	"context"
	"time"

	"service-04/internal/application/repository/custom_token_operations"
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

// Service handles custom token operations business logic
type Service struct {
	repo              *custom_token_operations.Repository
	cosmosApp         *app.USCApp
	blockchainStorage *storage.StateManager
	logger            *logging.Logger
	validator         *validation.Validator
	metrics           *metrics.MetricsService
}

// NewService creates a new custom token operations service
func NewService(
	repo *custom_token_operations.Repository,
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

// CreateBlockchainToken creates a new custom token
func (s *Service) CreateBlockchainToken(ctx context.Context, req *proto.CreateBlockchainTokenRequest) (*proto.CreateBlockchainTokenResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("create_blockchain_token", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Creating blockchain token in business service",
		logging.String("correlation_id", correlationID),
		logging.String("from", req.FromAddress),
		logging.String("name", req.TokenName))

	// Input validation using validator service
	if err := s.validator.ValidateWalletAddress(req.FromAddress); err != nil {
		s.logger.Error("From address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("from_address", req.FromAddress),
			logging.Error(err))
		s.metrics.RecordFailure("create_blockchain_token", "validation_error", map[string]string{
			"from_address": req.FromAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid from_address: %v", err)
	}

	if req.TokenName == "" {
		s.logger.Error("Token name is required",
			logging.String("correlation_id", correlationID))
		s.metrics.RecordFailure("create_blockchain_token", "validation_error", map[string]string{
			"token_name": req.TokenName,
		})
		return nil, status.Errorf(codes.InvalidArgument, "token_name is required")
	}

	if req.TokenSymbol == "" {
		s.logger.Error("Token symbol is required",
			logging.String("correlation_id", correlationID))
		s.metrics.RecordFailure("create_blockchain_token", "validation_error", map[string]string{
			"token_symbol": req.TokenSymbol,
		})
		return nil, status.Errorf(codes.InvalidArgument, "token_symbol is required")
	}

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.CreateBlockchainToken(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create blockchain token in repository",
			logging.String("correlation_id", correlationID),
			logging.String("from_address", req.FromAddress),
			logging.String("token_name", req.TokenName),
			logging.Error(err))
		s.metrics.RecordFailure("create_blockchain_token", "repository_error", map[string]string{
			"from_address": req.FromAddress,
			"token_name":   req.TokenName,
		})
		return nil, status.Errorf(codes.Internal, "failed to create blockchain token: %v", err)
	}

	// Record success metrics
	s.logger.Info("Blockchain token created successfully",
		logging.String("correlation_id", correlationID),
		logging.String("from_address", req.FromAddress),
		logging.String("token_name", req.TokenName))
	s.metrics.RecordSuccess("create_blockchain_token", map[string]string{
		"from_address": req.FromAddress,
		"token_name":   req.TokenName,
	})

	// Record blockchain-specific metric if token was created
	if response != nil && response.ContractAddress != "" {
		s.metrics.RecordTokenCreated(response.ContractAddress, req.TokenName)
	}

	return response, nil
}

// MintTokens mints custom tokens
func (s *Service) MintTokens(ctx context.Context, req *proto.MintTokensRequest) (*proto.MintTokensResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("mint_tokens", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Minting tokens in business service",
		logging.String("correlation_id", correlationID),
		logging.String("contract", req.ContractAddress),
		logging.String("to", req.ToAddress),
		logging.String("amount", req.Amount))

	// Input validation using validator service
	if err := s.validator.ValidateContractAddress(req.ContractAddress); err != nil {
		s.logger.Error("Contract address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("contract_address", req.ContractAddress),
			logging.Error(err))
		s.metrics.RecordFailure("mint_tokens", "validation_error", map[string]string{
			"contract_address": req.ContractAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid contract_address: %v", err)
	}

	if err := s.validator.ValidateWalletAddress(req.ToAddress); err != nil {
		s.logger.Error("To address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("to_address", req.ToAddress),
			logging.Error(err))
		s.metrics.RecordFailure("mint_tokens", "validation_error", map[string]string{
			"to_address": req.ToAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid to_address: %v", err)
	}

	if err := s.validator.ValidateAmount(req.Amount); err != nil {
		s.logger.Error("Amount validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("amount", req.Amount),
			logging.Error(err))
		s.metrics.RecordFailure("mint_tokens", "validation_error", map[string]string{
			"amount": req.Amount,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid amount: %v", err)
	}

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.MintTokens(ctx, req)
	if err != nil {
		s.logger.Error("Failed to mint tokens in repository",
			logging.String("correlation_id", correlationID),
			logging.String("contract_address", req.ContractAddress),
			logging.String("to_address", req.ToAddress),
			logging.Error(err))
		s.metrics.RecordFailure("mint_tokens", "repository_error", map[string]string{
			"contract_address": req.ContractAddress,
			"to_address":       req.ToAddress,
		})
		return nil, status.Errorf(codes.Internal, "failed to mint tokens: %v", err)
	}

	// Record success metrics
	s.logger.Info("Tokens minted successfully",
		logging.String("correlation_id", correlationID),
		logging.String("contract_address", req.ContractAddress),
		logging.String("to_address", req.ToAddress))
	s.metrics.RecordSuccess("mint_tokens", map[string]string{
		"contract_address": req.ContractAddress,
		"to_address":       req.ToAddress,
	})

	return response, nil
}

// GetTokenBalance retrieves token balance for an address
func (s *Service) GetTokenBalance(ctx context.Context, req *proto.GetTokenBalanceRequest) (*proto.GetTokenBalanceResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("get_token_balance", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Getting token balance in business service",
		logging.String("correlation_id", correlationID),
		logging.String("contract", req.ContractAddress),
		logging.String("wallet", req.WalletAddress))

	// Input validation using validator service
	if err := s.validator.ValidateContractAddress(req.ContractAddress); err != nil {
		s.logger.Error("Contract address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("contract_address", req.ContractAddress),
			logging.Error(err))
		s.metrics.RecordFailure("get_token_balance", "validation_error", map[string]string{
			"contract_address": req.ContractAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid contract_address: %v", err)
	}

	if err := s.validator.ValidateWalletAddress(req.WalletAddress); err != nil {
		s.logger.Error("Wallet address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("wallet_address", req.WalletAddress),
			logging.Error(err))
		s.metrics.RecordFailure("get_token_balance", "validation_error", map[string]string{
			"wallet_address": req.WalletAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid wallet_address: %v", err)
	}

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.GetTokenBalance(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get token balance in repository",
			logging.String("correlation_id", correlationID),
			logging.String("contract_address", req.ContractAddress),
			logging.String("wallet_address", req.WalletAddress),
			logging.Error(err))
		s.metrics.RecordFailure("get_token_balance", "repository_error", map[string]string{
			"contract_address": req.ContractAddress,
			"wallet_address":   req.WalletAddress,
		})
		return nil, status.Errorf(codes.Internal, "failed to get token balance: %v", err)
	}

	// Record success metrics
	s.logger.Info("Token balance retrieved successfully",
		logging.String("correlation_id", correlationID),
		logging.String("contract_address", req.ContractAddress),
		logging.String("wallet_address", req.WalletAddress))
	s.metrics.RecordSuccess("get_token_balance", map[string]string{
		"contract_address": req.ContractAddress,
		"wallet_address":   req.WalletAddress,
	})

	return response, nil
}

// GetTokenInfo retrieves token information
func (s *Service) GetTokenInfo(ctx context.Context, req *proto.GetTokenInfoRequest) (*proto.GetTokenInfoResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("get_token_info", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Getting token info in business service",
		logging.String("correlation_id", correlationID),
		logging.String("contract", req.ContractAddress))

	// Input validation using validator service
	if err := s.validator.ValidateContractAddress(req.ContractAddress); err != nil {
		s.logger.Error("Contract address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("contract_address", req.ContractAddress),
			logging.Error(err))
		s.metrics.RecordFailure("get_token_info", "validation_error", map[string]string{
			"contract_address": req.ContractAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid contract_address: %v", err)
	}

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.GetTokenInfo(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get token info in repository",
			logging.String("correlation_id", correlationID),
			logging.String("contract_address", req.ContractAddress),
			logging.Error(err))
		s.metrics.RecordFailure("get_token_info", "repository_error", map[string]string{
			"contract_address": req.ContractAddress,
		})
		return nil, status.Errorf(codes.Internal, "failed to get token info: %v", err)
	}

	// Record success metrics
	s.logger.Info("Token info retrieved successfully",
		logging.String("correlation_id", correlationID),
		logging.String("contract_address", req.ContractAddress))
	s.metrics.RecordSuccess("get_token_info", map[string]string{
		"contract_address": req.ContractAddress,
	})

	return response, nil
}

// BurnTokens burns custom tokens
func (s *Service) BurnTokens(ctx context.Context, req *proto.BurnTokensRequest) (*proto.BurnTokensResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("burn_tokens", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Burning tokens in business service",
		logging.String("correlation_id", correlationID),
		logging.String("contract", req.ContractAddress),
		logging.String("from", req.FromAddress),
		logging.String("amount", req.Amount))

	// Input validation using validator service
	if err := s.validator.ValidateContractAddress(req.ContractAddress); err != nil {
		s.logger.Error("Contract address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("contract_address", req.ContractAddress),
			logging.Error(err))
		s.metrics.RecordFailure("burn_tokens", "validation_error", map[string]string{
			"contract_address": req.ContractAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid contract_address: %v", err)
	}

	if err := s.validator.ValidateWalletAddress(req.FromAddress); err != nil {
		s.logger.Error("From address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("from_address", req.FromAddress),
			logging.Error(err))
		s.metrics.RecordFailure("burn_tokens", "validation_error", map[string]string{
			"from_address": req.FromAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid from_address: %v", err)
	}

	if err := s.validator.ValidateAmount(req.Amount); err != nil {
		s.logger.Error("Amount validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("amount", req.Amount),
			logging.Error(err))
		s.metrics.RecordFailure("burn_tokens", "validation_error", map[string]string{
			"amount": req.Amount,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid amount: %v", err)
	}

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.BurnTokens(ctx, req)
	if err != nil {
		s.logger.Error("Failed to burn tokens in repository",
			logging.String("correlation_id", correlationID),
			logging.String("contract_address", req.ContractAddress),
			logging.String("from_address", req.FromAddress),
			logging.Error(err))
		s.metrics.RecordFailure("burn_tokens", "repository_error", map[string]string{
			"contract_address": req.ContractAddress,
			"from_address":     req.FromAddress,
		})
		return nil, status.Errorf(codes.Internal, "failed to burn tokens: %v", err)
	}

	// Record success metrics
	s.logger.Info("Tokens burned successfully",
		logging.String("correlation_id", correlationID),
		logging.String("contract_address", req.ContractAddress),
		logging.String("from_address", req.FromAddress))
	s.metrics.RecordSuccess("burn_tokens", map[string]string{
		"contract_address": req.ContractAddress,
		"from_address":     req.FromAddress,
	})

	return response, nil
}
