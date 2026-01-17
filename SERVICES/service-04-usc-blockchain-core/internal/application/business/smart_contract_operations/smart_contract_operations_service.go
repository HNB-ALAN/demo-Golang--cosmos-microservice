package smart_contract_operations

import (
	"context"
	"time"

	"service-04/internal/application/repository/smart_contract_operations"
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

// Service handles smart contract operations business logic
type Service struct {
	repo              *smart_contract_operations.Repository
	cosmosApp         *app.USCApp
	blockchainStorage *storage.StateManager
	logger            *logging.Logger
	validator         *validation.Validator
	metrics           *metrics.MetricsService
}

// NewService creates a new smart contract operations service
func NewService(
	repo *smart_contract_operations.Repository,
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

// DeployContract deploys a new smart contract
func (s *Service) DeployContract(ctx context.Context, req *proto.DeployContractRequest) (*proto.DeployContractResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("deploy_contract", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Deploying contract in business service",
		logging.String("correlation_id", correlationID),
		logging.String("from", req.FromAddress),
		logging.String("contractName", req.ContractName))

	// Input validation using validator service
	if err := s.validator.ValidateWalletAddress(req.FromAddress); err != nil {
		s.logger.Error("From address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("from_address", req.FromAddress),
			logging.Error(err))
		s.metrics.RecordFailure("deploy_contract", "validation_error", map[string]string{
			"from_address": req.FromAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid from_address: %v", err)
	}

	if req.ContractName == "" {
		s.logger.Error("Contract name is required",
			logging.String("correlation_id", correlationID))
		s.metrics.RecordFailure("deploy_contract", "validation_error", map[string]string{
			"contract_name": req.ContractName,
		})
		return nil, status.Errorf(codes.InvalidArgument, "contract_name is required")
	}

	if req.Bytecode == "" {
		s.logger.Error("Bytecode is required",
			logging.String("correlation_id", correlationID))
		s.metrics.RecordFailure("deploy_contract", "validation_error", map[string]string{})
		return nil, status.Errorf(codes.InvalidArgument, "bytecode is required")
	}

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.DeployContract(ctx, req)
	if err != nil {
		s.logger.Error("Failed to deploy contract in repository",
			logging.String("correlation_id", correlationID),
			logging.String("from_address", req.FromAddress),
			logging.String("contract_name", req.ContractName),
			logging.Error(err))
		s.metrics.RecordFailure("deploy_contract", "repository_error", map[string]string{
			"from_address":  req.FromAddress,
			"contract_name": req.ContractName,
		})
		return nil, status.Errorf(codes.Internal, "failed to deploy contract: %v", err)
	}

	// Record success metrics
	s.logger.Info("Contract deployed successfully",
		logging.String("correlation_id", correlationID),
		logging.String("from_address", req.FromAddress),
		logging.String("contract_name", req.ContractName))
	s.metrics.RecordSuccess("deploy_contract", map[string]string{
		"from_address":  req.FromAddress,
		"contract_name": req.ContractName,
	})

	// Record blockchain-specific metric if contract was deployed
	if response != nil && response.ContractAddress != "" {
		s.metrics.RecordContractDeployed(response.ContractAddress)
	}

	return response, nil
}

// ExecuteContract executes a smart contract function
func (s *Service) ExecuteContract(ctx context.Context, req *proto.ExecuteContractRequest) (*proto.ExecuteContractResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("execute_contract", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Executing contract in business service",
		logging.String("correlation_id", correlationID),
		logging.String("contract", req.ContractAddress),
		logging.String("function", req.FunctionName))

	// Input validation using validator service
	if err := s.validator.ValidateContractAddress(req.ContractAddress); err != nil {
		s.logger.Error("Contract address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("contract_address", req.ContractAddress),
			logging.Error(err))
		s.metrics.RecordFailure("execute_contract", "validation_error", map[string]string{
			"contract_address": req.ContractAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid contract_address: %v", err)
	}

	if req.FunctionName == "" {
		s.logger.Error("Function name is required",
			logging.String("correlation_id", correlationID))
		s.metrics.RecordFailure("execute_contract", "validation_error", map[string]string{
			"function_name": req.FunctionName,
		})
		return nil, status.Errorf(codes.InvalidArgument, "function_name is required")
	}

	if err := s.validator.ValidateWalletAddress(req.FromAddress); err != nil {
		s.logger.Error("From address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("from_address", req.FromAddress),
			logging.Error(err))
		s.metrics.RecordFailure("execute_contract", "validation_error", map[string]string{
			"from_address": req.FromAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid from_address: %v", err)
	}

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.ExecuteContract(ctx, req)
	if err != nil {
		s.logger.Error("Failed to execute contract in repository",
			logging.String("correlation_id", correlationID),
			logging.String("contract_address", req.ContractAddress),
			logging.String("function_name", req.FunctionName),
			logging.Error(err))
		s.metrics.RecordFailure("execute_contract", "repository_error", map[string]string{
			"contract_address": req.ContractAddress,
			"function_name":    req.FunctionName,
		})
		return nil, status.Errorf(codes.Internal, "failed to execute contract: %v", err)
	}

	// Record success metrics
	s.logger.Info("Contract executed successfully",
		logging.String("correlation_id", correlationID),
		logging.String("contract_address", req.ContractAddress),
		logging.String("function_name", req.FunctionName))
	s.metrics.RecordSuccess("execute_contract", map[string]string{
		"contract_address": req.ContractAddress,
		"function_name":    req.FunctionName,
	})

	return response, nil
}

// QueryContract queries a smart contract function
func (s *Service) QueryContract(ctx context.Context, req *proto.QueryContractRequest) (*proto.QueryContractResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("query_contract", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Querying contract in business service",
		logging.String("correlation_id", correlationID),
		logging.String("contract", req.ContractAddress),
		logging.String("function", req.FunctionName))

	// Input validation using validator service
	if err := s.validator.ValidateContractAddress(req.ContractAddress); err != nil {
		s.logger.Error("Contract address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("contract_address", req.ContractAddress),
			logging.Error(err))
		s.metrics.RecordFailure("query_contract", "validation_error", map[string]string{
			"contract_address": req.ContractAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid contract_address: %v", err)
	}

	if req.FunctionName == "" {
		s.logger.Error("Function name is required",
			logging.String("correlation_id", correlationID))
		s.metrics.RecordFailure("query_contract", "validation_error", map[string]string{
			"function_name": req.FunctionName,
		})
		return nil, status.Errorf(codes.InvalidArgument, "function_name is required")
	}

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.QueryContract(ctx, req)
	if err != nil {
		s.logger.Error("Failed to query contract in repository",
			logging.String("correlation_id", correlationID),
			logging.String("contract_address", req.ContractAddress),
			logging.String("function_name", req.FunctionName),
			logging.Error(err))
		s.metrics.RecordFailure("query_contract", "repository_error", map[string]string{
			"contract_address": req.ContractAddress,
			"function_name":    req.FunctionName,
		})
		return nil, status.Errorf(codes.Internal, "failed to query contract: %v", err)
	}

	// Record success metrics
	s.logger.Info("Contract queried successfully",
		logging.String("correlation_id", correlationID),
		logging.String("contract_address", req.ContractAddress),
		logging.String("function_name", req.FunctionName))
	s.metrics.RecordSuccess("query_contract", map[string]string{
		"contract_address": req.ContractAddress,
		"function_name":    req.FunctionName,
	})

	return response, nil
}

// GetContractCode retrieves contract source code
func (s *Service) GetContractCode(ctx context.Context, req *proto.GetContractCodeRequest) (*proto.GetContractCodeResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("get_contract_code", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Getting contract code in business service",
		logging.String("correlation_id", correlationID),
		logging.String("contract", req.ContractAddress))

	// Input validation using validator service
	if err := s.validator.ValidateContractAddress(req.ContractAddress); err != nil {
		s.logger.Error("Contract address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("contract_address", req.ContractAddress),
			logging.Error(err))
		s.metrics.RecordFailure("get_contract_code", "validation_error", map[string]string{
			"contract_address": req.ContractAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid contract_address: %v", err)
	}

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.GetContractCode(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get contract code in repository",
			logging.String("correlation_id", correlationID),
			logging.String("contract_address", req.ContractAddress),
			logging.Error(err))
		s.metrics.RecordFailure("get_contract_code", "repository_error", map[string]string{
			"contract_address": req.ContractAddress,
		})
		return nil, status.Errorf(codes.Internal, "failed to get contract code: %v", err)
	}

	// Record success metrics
	s.logger.Info("Contract code retrieved successfully",
		logging.String("correlation_id", correlationID),
		logging.String("contract_address", req.ContractAddress))
	s.metrics.RecordSuccess("get_contract_code", map[string]string{
		"contract_address": req.ContractAddress,
	})

	return response, nil
}

// GetContractStorage retrieves contract storage
func (s *Service) GetContractStorage(ctx context.Context, req *proto.GetContractStorageRequest) (*proto.GetContractStorageResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("get_contract_storage", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Getting contract storage in business service",
		logging.String("correlation_id", correlationID),
		logging.String("contract", req.ContractAddress),
		logging.String("key", req.StorageKey))

	// Input validation using validator service
	if err := s.validator.ValidateContractAddress(req.ContractAddress); err != nil {
		s.logger.Error("Contract address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("contract_address", req.ContractAddress),
			logging.Error(err))
		s.metrics.RecordFailure("get_contract_storage", "validation_error", map[string]string{
			"contract_address": req.ContractAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid contract_address: %v", err)
	}

	if req.StorageKey == "" {
		s.logger.Error("Storage key is required",
			logging.String("correlation_id", correlationID))
		s.metrics.RecordFailure("get_contract_storage", "validation_error", map[string]string{
			"storage_key": req.StorageKey,
		})
		return nil, status.Errorf(codes.InvalidArgument, "storage_key is required")
	}

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.GetContractStorage(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get contract storage in repository",
			logging.String("correlation_id", correlationID),
			logging.String("contract_address", req.ContractAddress),
			logging.String("storage_key", req.StorageKey),
			logging.Error(err))
		s.metrics.RecordFailure("get_contract_storage", "repository_error", map[string]string{
			"contract_address": req.ContractAddress,
			"storage_key":      req.StorageKey,
		})
		return nil, status.Errorf(codes.Internal, "failed to get contract storage: %v", err)
	}

	// Record success metrics
	s.logger.Info("Contract storage retrieved successfully",
		logging.String("correlation_id", correlationID),
		logging.String("contract_address", req.ContractAddress),
		logging.String("storage_key", req.StorageKey))
	s.metrics.RecordSuccess("get_contract_storage", map[string]string{
		"contract_address": req.ContractAddress,
		"storage_key":      req.StorageKey,
	})

	return response, nil
}
