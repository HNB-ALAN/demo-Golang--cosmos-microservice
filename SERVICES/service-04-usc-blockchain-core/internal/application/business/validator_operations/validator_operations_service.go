package validator_operations

import (
	"context"
	"time"

	"service-04/internal/application/repository/validator_operations"
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

// Service handles validator operations business logic
type Service struct {
	repo              *validator_operations.Repository
	cosmosApp         *app.USCApp
	blockchainStorage *storage.StateManager
	logger            *logging.Logger
	validator         *validation.Validator
	metrics           *metrics.MetricsService
}

// NewService creates a new validator operations service
func NewService(
	repo *validator_operations.Repository,
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

// RegisterValidator registers a new validator
func (s *Service) RegisterValidator(ctx context.Context, req *proto.RegisterValidatorRequest) (*proto.RegisterValidatorResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("register_validator", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Registering validator",
		logging.String("correlation_id", correlationID),
		logging.String("validator_address", req.ValidatorAddress),
		logging.String("service", "validator_operations"))

	// Input validation using validator service
	if err := s.validator.ValidateValidatorAddress(req.ValidatorAddress); err != nil {
		s.logger.Error("Validator address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("validator_address", req.ValidatorAddress),
			logging.Error(err))
		s.metrics.RecordFailure("register_validator", "validation_error", map[string]string{
			"validator_address": req.ValidatorAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid validator_address: %v", err)
	}

	if req.ValidatorPublicKey == "" {
		s.logger.Error("Validator public key is required",
			logging.String("correlation_id", correlationID))
		s.metrics.RecordFailure("register_validator", "validation_error", map[string]string{
			"validator_public_key": req.ValidatorPublicKey,
		})
		return nil, status.Errorf(codes.InvalidArgument, "validator_public_key is required")
	}

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.RegisterValidator(ctx, req)
	if err != nil {
		s.logger.Error("Failed to register validator in repository",
			logging.String("correlation_id", correlationID),
			logging.String("validator_address", req.ValidatorAddress),
			logging.Error(err))
		s.metrics.RecordFailure("register_validator", "repository_error", map[string]string{
			"validator_address": req.ValidatorAddress,
		})
		return nil, status.Errorf(codes.Internal, "failed to register validator: %v", err)
	}

	// Record success metrics
	s.logger.Info("Validator registered successfully",
		logging.String("correlation_id", correlationID),
		logging.String("validator_address", req.ValidatorAddress))
	s.metrics.RecordSuccess("register_validator", map[string]string{
		"validator_address": req.ValidatorAddress,
	})

	// Record blockchain-specific metric if validator was registered
	if response != nil && response.Success {
		s.metrics.RecordValidatorRegistered(req.ValidatorAddress)
	}

	return response, nil
}

// GetValidators gets list of validators
func (s *Service) GetValidators(ctx context.Context, req *proto.GetValidatorsRequest) (*proto.GetValidatorsResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("get_validators", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Getting validators",
		logging.String("correlation_id", correlationID),
		logging.String("service", "validator_operations"))

	// Business logic validation
	// Use helper function to normalize pagination (reduces duplicate code)
	req.Limit, req.Offset = utils.NormalizePaginationWithDefaults(req.Limit, req.Offset)

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.GetValidators(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get validators in repository",
			logging.String("correlation_id", correlationID),
			logging.Error(err))
		s.metrics.RecordFailure("get_validators", "repository_error", map[string]string{})
		return nil, status.Errorf(codes.Internal, "failed to get validators: %v", err)
	}

	// Record success metrics
	s.logger.Info("Validators retrieved successfully",
		logging.String("correlation_id", correlationID))
	s.metrics.RecordSuccess("get_validators", map[string]string{})

	return response, nil
}

// GetValidatorStatus gets validator status
func (s *Service) GetValidatorStatus(ctx context.Context, req *proto.GetValidatorStatusRequest) (*proto.GetValidatorStatusResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("get_validator_status", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Getting validator status",
		logging.String("correlation_id", correlationID),
		logging.String("validator_address", req.ValidatorAddress),
		logging.String("service", "validator_operations"))

	// Input validation using validator service
	if err := s.validator.ValidateValidatorAddress(req.ValidatorAddress); err != nil {
		s.logger.Error("Validator address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("validator_address", req.ValidatorAddress),
			logging.Error(err))
		s.metrics.RecordFailure("get_validator_status", "validation_error", map[string]string{
			"validator_address": req.ValidatorAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid validator_address: %v", err)
	}

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.GetValidatorStatus(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get validator status in repository",
			logging.String("correlation_id", correlationID),
			logging.String("validator_address", req.ValidatorAddress),
			logging.Error(err))
		s.metrics.RecordFailure("get_validator_status", "repository_error", map[string]string{
			"validator_address": req.ValidatorAddress,
		})
		return nil, status.Errorf(codes.Internal, "failed to get validator status: %v", err)
	}

	// Record success metrics
	s.logger.Info("Validator status retrieved successfully",
		logging.String("correlation_id", correlationID),
		logging.String("validator_address", req.ValidatorAddress))
	s.metrics.RecordSuccess("get_validator_status", map[string]string{
		"validator_address": req.ValidatorAddress,
	})

	return response, nil
}


// StakeUSC stakes USC tokens
func (s *Service) StakeUSC(ctx context.Context, req *proto.StakeUSCRequest) (*proto.StakeUSCResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("stake_usc", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Staking USC",
		logging.String("correlation_id", correlationID),
		logging.String("validator_address", req.ValidatorAddress),
		logging.String("stake_amount", req.StakeAmount),
		logging.String("service", "validator_operations"))

	// Input validation using validator service
	if err := s.validator.ValidateValidatorAddress(req.ValidatorAddress); err != nil {
		s.logger.Error("Validator address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("validator_address", req.ValidatorAddress),
			logging.Error(err))
		s.metrics.RecordFailure("stake_usc", "validation_error", map[string]string{
			"validator_address": req.ValidatorAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid validator_address: %v", err)
	}

	if err := s.validator.ValidateAmount(req.StakeAmount); err != nil {
		s.logger.Error("Stake amount validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("stake_amount", req.StakeAmount),
			logging.Error(err))
		s.metrics.RecordFailure("stake_usc", "validation_error", map[string]string{
			"stake_amount": req.StakeAmount,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid stake_amount: %v", err)
	}

	// Delegate to repository
	response, err := s.repo.StakeUSC(ctx, req)
	if err != nil {
		s.logger.Error("Failed to stake USC in repository",
			logging.String("correlation_id", correlationID),
			logging.String("validator_address", req.ValidatorAddress),
			logging.String("stake_amount", req.StakeAmount),
			logging.Error(err))
		s.metrics.RecordFailure("stake_usc", "repository_error", map[string]string{
			"validator_address": req.ValidatorAddress,
			"stake_amount":      req.StakeAmount,
		})
		return nil, status.Errorf(codes.Internal, "failed to stake USC: %v", err)
	}

	// Record success metrics
	s.logger.Info("USC staked successfully",
		logging.String("correlation_id", correlationID),
		logging.String("validator_address", req.ValidatorAddress),
		logging.String("stake_amount", req.StakeAmount))
	s.metrics.RecordSuccess("stake_usc", map[string]string{
		"validator_address": req.ValidatorAddress,
		"stake_amount":      req.StakeAmount,
	})

	return response, nil
}

// UnstakeUSC unstakes USC tokens
func (s *Service) UnstakeUSC(ctx context.Context, req *proto.UnstakeUSCRequest) (*proto.UnstakeUSCResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("unstake_usc", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Unstaking USC",
		logging.String("correlation_id", correlationID),
		logging.String("validator_address", req.ValidatorAddress),
		logging.String("unstake_amount", req.UnstakeAmount),
		logging.String("service", "validator_operations"))

	// Input validation using validator service
	if err := s.validator.ValidateValidatorAddress(req.ValidatorAddress); err != nil {
		s.logger.Error("Validator address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("validator_address", req.ValidatorAddress),
			logging.Error(err))
		s.metrics.RecordFailure("unstake_usc", "validation_error", map[string]string{
			"validator_address": req.ValidatorAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid validator_address: %v", err)
	}

	if err := s.validator.ValidateAmount(req.UnstakeAmount); err != nil {
		s.logger.Error("Unstake amount validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("unstake_amount", req.UnstakeAmount),
			logging.Error(err))
		s.metrics.RecordFailure("unstake_usc", "validation_error", map[string]string{
			"unstake_amount": req.UnstakeAmount,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid unstake_amount: %v", err)
	}

	// Delegate to repository
	response, err := s.repo.UnstakeUSC(ctx, req)
	if err != nil {
		s.logger.Error("Failed to unstake USC in repository",
			logging.String("correlation_id", correlationID),
			logging.String("validator_address", req.ValidatorAddress),
			logging.String("unstake_amount", req.UnstakeAmount),
			logging.Error(err))
		s.metrics.RecordFailure("unstake_usc", "repository_error", map[string]string{
			"validator_address": req.ValidatorAddress,
			"unstake_amount":    req.UnstakeAmount,
		})
		return nil, status.Errorf(codes.Internal, "failed to unstake USC: %v", err)
	}

	// Record success metrics
	s.logger.Info("USC unstaked successfully",
		logging.String("correlation_id", correlationID),
		logging.String("validator_address", req.ValidatorAddress),
		logging.String("unstake_amount", req.UnstakeAmount))
	s.metrics.RecordSuccess("unstake_usc", map[string]string{
		"validator_address": req.ValidatorAddress,
		"unstake_amount":    req.UnstakeAmount,
	})

	return response, nil
}
