package block_operations

import (
	"context"
	"fmt"
	"time"

	"service-04/internal/application/repository/block_operations"
	"service-04/internal/application/utils"
	"service-04/internal/infrastructure/metrics"
	"service-04/internal/infrastructure/validation"
	proto "service-04/proto"

	"github.com/usc-platform/shared/logging"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/app"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Service handles block operations business logic
type Service struct {
	repo              *block_operations.Repository
	cosmosApp         *app.USCApp
	blockchainStorage *storage.StateManager
	logger            *logging.Logger
	validator         *validation.Validator
	metrics           *metrics.MetricsService
}

// NewService creates a new block operations service
func NewService(
	repo *block_operations.Repository,
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

// ProduceBlock creates a new block
func (s *Service) ProduceBlock(ctx context.Context, req *proto.ProduceBlockRequest) (*proto.ProduceBlockResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("produce_block", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Producing block in business service",
		logging.String("correlation_id", correlationID),
		logging.String("validator", req.ValidatorId))

	// Input validation using validator service
	if err := s.validator.ValidateValidatorId(req.ValidatorId); err != nil {
		s.logger.Error("Validator ID validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("validator_id", req.ValidatorId),
			logging.Error(err))
		s.metrics.RecordFailure("produce_block", "validation_error", map[string]string{
			"validator_id": req.ValidatorId,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid validator_id: %v", err)
	}

	// Validate timestamp
	if req.Timestamp <= 0 {
		req.Timestamp = time.Now().Unix()
	}

	// Delegate to repository (repository will handle blockchain or database)
	response, err := s.repo.ProduceBlock(ctx, req)
	if err != nil {
		s.logger.Error("Failed to produce block in repository",
			logging.String("correlation_id", correlationID),
			logging.String("validator_id", req.ValidatorId),
			logging.Error(err))
		s.metrics.RecordFailure("produce_block", "repository_error", map[string]string{
			"validator_id": req.ValidatorId,
		})
		return nil, status.Errorf(codes.Internal, "failed to produce block: %v", err)
	}

	// Record success metrics
	s.logger.Info("Block produced successfully",
		logging.String("correlation_id", correlationID),
		logging.String("validator_id", req.ValidatorId))
	s.metrics.RecordSuccess("produce_block", map[string]string{
		"validator_id": req.ValidatorId,
	})

	// Record blockchain-specific metric if block was created
	if response != nil && response.Success && response.BlockHash != "" {
		s.metrics.RecordBlockCreated(response.BlockNumber, response.BlockHash)
	}

	return response, nil
}

// ValidateBlock validates a block
func (s *Service) ValidateBlock(ctx context.Context, req *proto.ValidateBlockRequest) (*proto.ValidateBlockResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("validate_block", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Validating block in business service",
		logging.String("correlation_id", correlationID),
		logging.String("block_hash", req.BlockHash),
		logging.Int64("block_number", req.BlockNumber))

	// Input validation using validator service
	if req.BlockHash != "" {
		if err := s.validator.ValidateBlockHash(req.BlockHash); err != nil {
			s.logger.Error("Block hash validation failed",
				logging.String("correlation_id", correlationID),
				logging.String("block_hash", req.BlockHash),
				logging.Error(err))
			s.metrics.RecordFailure("validate_block", "validation_error", map[string]string{
				"block_hash": req.BlockHash,
			})
			return nil, status.Errorf(codes.InvalidArgument, "invalid block_hash: %v", err)
		}
	} else if req.BlockNumber > 0 {
		if err := s.validator.ValidateBlockNumber(req.BlockNumber); err != nil {
			s.logger.Error("Block number validation failed",
				logging.String("correlation_id", correlationID),
				logging.Int64("block_number", req.BlockNumber),
				logging.Error(err))
			s.metrics.RecordFailure("validate_block", "validation_error", map[string]string{
				"block_number": fmt.Sprintf("%d", req.BlockNumber),
			})
			return nil, status.Errorf(codes.InvalidArgument, "invalid block_number: %v", err)
		}
	} else {
		s.logger.Error("Block hash or block number is required",
			logging.String("correlation_id", correlationID))
		s.metrics.RecordFailure("validate_block", "validation_error", map[string]string{})
		return nil, status.Errorf(codes.InvalidArgument, "block_hash or block_number is required")
	}

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.ValidateBlock(ctx, req)
	if err != nil {
		s.logger.Error("Failed to validate block in repository",
			logging.String("correlation_id", correlationID),
			logging.String("block_hash", req.BlockHash),
			logging.Int64("block_number", req.BlockNumber),
			logging.Error(err))
		s.metrics.RecordFailure("validate_block", "repository_error", map[string]string{
			"block_hash":   req.BlockHash,
			"block_number": fmt.Sprintf("%d", req.BlockNumber),
		})
		return nil, status.Errorf(codes.Internal, "failed to validate block: %v", err)
	}

	// Record success metrics
	s.logger.Info("Block validation completed",
		logging.String("correlation_id", correlationID),
		logging.Bool("valid", response.Valid))
	s.metrics.RecordSuccess("validate_block", map[string]string{
		"valid": fmt.Sprintf("%v", response.Valid),
	})

	return response, nil
}

// GetBlock retrieves a block by number
func (s *Service) GetBlock(ctx context.Context, req *proto.GetBlockRequest) (*proto.GetBlockResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("get_block", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Getting block in business service",
		logging.String("correlation_id", correlationID),
		logging.Int64("block_number", req.BlockNumber))

	// Input validation using validator service
	if err := s.validator.ValidateBlockNumber(req.BlockNumber); err != nil {
		s.logger.Error("Block number validation failed",
			logging.String("correlation_id", correlationID),
			logging.Int64("block_number", req.BlockNumber),
			logging.Error(err))
		s.metrics.RecordFailure("get_block", "validation_error", map[string]string{
			"block_number": fmt.Sprintf("%d", req.BlockNumber),
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid block_number: %v", err)
	}

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.GetBlock(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get block in repository",
			logging.String("correlation_id", correlationID),
			logging.Int64("block_number", req.BlockNumber),
			logging.Error(err))
		s.metrics.RecordFailure("get_block", "repository_error", map[string]string{
			"block_number": fmt.Sprintf("%d", req.BlockNumber),
		})
		return nil, status.Errorf(codes.Internal, "failed to get block: %v", err)
	}

	// Record success metrics
	s.logger.Info("Block retrieved successfully",
		logging.String("correlation_id", correlationID),
		logging.Int64("block_number", req.BlockNumber))
	s.metrics.RecordSuccess("get_block", map[string]string{
		"block_number": fmt.Sprintf("%d", req.BlockNumber),
	})

	return response, nil
}

// GetBlockByHash retrieves a block by hash
func (s *Service) GetBlockByHash(ctx context.Context, req *proto.GetBlockByHashRequest) (*proto.GetBlockResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("get_block_by_hash", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Getting block by hash in business service",
		logging.String("correlation_id", correlationID),
		logging.String("block_hash", req.BlockHash))

	// Input validation using validator service
	if err := s.validator.ValidateBlockHash(req.BlockHash); err != nil {
		s.logger.Error("Block hash validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("block_hash", req.BlockHash),
			logging.Error(err))
		s.metrics.RecordFailure("get_block_by_hash", "validation_error", map[string]string{
			"block_hash": req.BlockHash,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid block_hash: %v", err)
	}

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.GetBlockByHash(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get block by hash in repository",
			logging.String("correlation_id", correlationID),
			logging.String("block_hash", req.BlockHash),
			logging.Error(err))
		s.metrics.RecordFailure("get_block_by_hash", "repository_error", map[string]string{
			"block_hash": req.BlockHash,
		})
		return nil, status.Errorf(codes.Internal, "failed to get block by hash: %v", err)
	}

	// Record success metrics
	s.logger.Info("Block retrieved by hash successfully",
		logging.String("correlation_id", correlationID),
		logging.String("block_hash", req.BlockHash))
	s.metrics.RecordSuccess("get_block_by_hash", map[string]string{
		"block_hash": req.BlockHash,
	})

	return response, nil
}

// GetLatestBlock retrieves the latest block
func (s *Service) GetLatestBlock(ctx context.Context, req *emptypb.Empty) (*proto.GetBlockResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("get_latest_block", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Getting latest block in business service",
		logging.String("correlation_id", correlationID))

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.GetLatestBlock(ctx)
	if err != nil {
		s.logger.Error("Failed to get latest block in repository",
			logging.String("correlation_id", correlationID),
			logging.Error(err))
		s.metrics.RecordFailure("get_latest_block", "repository_error", map[string]string{})
		return nil, status.Errorf(codes.Internal, "failed to get latest block: %v", err)
	}

	// Record success metrics
	s.logger.Info("Latest block retrieved successfully",
		logging.String("correlation_id", correlationID))
	s.metrics.RecordSuccess("get_latest_block", map[string]string{})

	return response, nil
}

// GetBlockRange retrieves a range of blocks
func (s *Service) GetBlockRange(ctx context.Context, req *proto.GetBlockRangeRequest) (*proto.GetBlockRangeResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("get_block_range", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Getting block range in business service",
		logging.String("correlation_id", correlationID),
		logging.Int64("start_block", req.StartBlock),
		logging.Int64("end_block", req.EndBlock))

	// Input validation using validator service
	if err := s.validator.ValidateBlockNumber(req.StartBlock); err != nil {
		s.logger.Error("Start block number validation failed",
			logging.String("correlation_id", correlationID),
			logging.Int64("start_block", req.StartBlock),
			logging.Error(err))
		s.metrics.RecordFailure("get_block_range", "validation_error", map[string]string{
			"start_block": fmt.Sprintf("%d", req.StartBlock),
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid start_block: %v", err)
	}
	if err := s.validator.ValidateBlockNumber(req.EndBlock); err != nil {
		s.logger.Error("End block number validation failed",
			logging.String("correlation_id", correlationID),
			logging.Int64("end_block", req.EndBlock),
			logging.Error(err))
		s.metrics.RecordFailure("get_block_range", "validation_error", map[string]string{
			"end_block": fmt.Sprintf("%d", req.EndBlock),
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid end_block: %v", err)
	}
	if req.StartBlock > req.EndBlock {
		s.logger.Error("Start block greater than end block",
			logging.String("correlation_id", correlationID),
			logging.Int64("start_block", req.StartBlock),
			logging.Int64("end_block", req.EndBlock))
		s.metrics.RecordFailure("get_block_range", "validation_error", map[string]string{
			"start_block": fmt.Sprintf("%d", req.StartBlock),
			"end_block":   fmt.Sprintf("%d", req.EndBlock),
		})
		return nil, status.Errorf(codes.InvalidArgument, "start_block must be less than or equal to end_block")
	}

	// Validate limit
	// Use helper function to normalize pagination (reduces duplicate code)
	// Note: Block range uses different defaults (50 default, 100 max)
	req.Limit, req.Offset = utils.NormalizePagination(req.Limit, req.Offset, utils.PaginationConfig{
		DefaultLimit:  50,
		MaxLimit:      100,
		DefaultOffset: 0,
	})

	// Delegate to repository
	response, err := s.repo.GetBlockRange(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get block range in repository",
			logging.String("correlation_id", correlationID),
			logging.Int64("start_block", req.StartBlock),
			logging.Int64("end_block", req.EndBlock),
			logging.Error(err))
		s.metrics.RecordFailure("get_block_range", "repository_error", map[string]string{
			"start_block": fmt.Sprintf("%d", req.StartBlock),
			"end_block":   fmt.Sprintf("%d", req.EndBlock),
		})
		return nil, status.Errorf(codes.Internal, "failed to get block range: %v", err)
	}

	// Record success metrics
	s.logger.Info("Block range retrieved successfully",
		logging.String("correlation_id", correlationID),
		logging.Int64("start_block", req.StartBlock),
		logging.Int64("end_block", req.EndBlock))
	s.metrics.RecordSuccess("get_block_range", map[string]string{
		"start_block": fmt.Sprintf("%d", req.StartBlock),
		"end_block":   fmt.Sprintf("%d", req.EndBlock),
	})

	return response, nil
}
