package network_operations

import (
	"context"
	"time"

	"service-04/internal/application/repository/network_operations"
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

// Service handles network operations business logic
type Service struct {
	repo              *network_operations.Repository
	cosmosApp         *app.USCApp
	blockchainStorage *storage.StateManager
	logger            *logging.Logger
	validator         *validation.Validator
	metrics           *metrics.MetricsService
}

// NewService creates a new network operations service
func NewService(
	repo *network_operations.Repository,
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

// GetNetworkInfo retrieves network information
func (s *Service) GetNetworkInfo(ctx context.Context) (*proto.GetNetworkInfoResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("get_network_info", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Getting network info in business service",
		logging.String("correlation_id", correlationID))

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.GetNetworkInfo(ctx)
	if err != nil {
		s.logger.Error("Failed to get network info in repository",
			logging.String("correlation_id", correlationID),
			logging.Error(err))
		s.metrics.RecordFailure("get_network_info", "repository_error", map[string]string{})
		return nil, status.Errorf(codes.Internal, "failed to get network info: %v", err)
	}

	// Record success metrics
	s.logger.Info("Network info retrieved successfully",
		logging.String("correlation_id", correlationID))
	s.metrics.RecordSuccess("get_network_info", map[string]string{})

	return response, nil
}

// GetChainInfo retrieves chain information
func (s *Service) GetChainInfo(ctx context.Context) (*proto.GetChainInfoResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("get_chain_info", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Getting chain info in business service",
		logging.String("correlation_id", correlationID))

	// Call repository
	response, err := s.repo.GetChainInfo(ctx)
	if err != nil {
		s.logger.Error("Failed to get chain info in repository",
			logging.String("correlation_id", correlationID),
			logging.Error(err))
		s.metrics.RecordFailure("get_chain_info", "repository_error", map[string]string{})
		return nil, status.Errorf(codes.Internal, "failed to get chain info: %v", err)
	}

	// Record success metrics
	s.logger.Info("Chain info retrieved successfully",
		logging.String("correlation_id", correlationID))
	s.metrics.RecordSuccess("get_chain_info", map[string]string{})

	return response, nil
}

// GetPeers retrieves list of peers
func (s *Service) GetPeers(ctx context.Context, req *proto.GetPeersRequest) (*proto.GetPeersResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("get_peers", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Getting peers in business service",
		logging.String("correlation_id", correlationID),
		logging.Int32("limit", req.Limit),
		logging.Int32("offset", req.Offset))

	// Normalize pagination
	// Use helper function to normalize pagination (reduces duplicate code)
	req.Limit, req.Offset = utils.NormalizePaginationWithDefaults(req.Limit, req.Offset)

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.GetPeers(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get peers in repository",
			logging.String("correlation_id", correlationID),
			logging.Error(err))
		s.metrics.RecordFailure("get_peers", "repository_error", map[string]string{})
		return nil, status.Errorf(codes.Internal, "failed to get peers: %v", err)
	}

	// Record success metrics
	s.logger.Info("Peers retrieved successfully",
		logging.String("correlation_id", correlationID))
	s.metrics.RecordSuccess("get_peers", map[string]string{})

	return response, nil
}

// GetNetworkStats retrieves network statistics
func (s *Service) GetNetworkStats(ctx context.Context, req *proto.GetNetworkStatsRequest) (*proto.GetNetworkStatsResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("get_network_stats", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Getting network stats in business service",
		logging.String("correlation_id", correlationID),
		logging.String("timeRange", req.TimeRange),
		logging.String("metricType", req.MetricType))

	// Call repository
	response, err := s.repo.GetNetworkStats(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get network stats in repository",
			logging.String("correlation_id", correlationID),
			logging.String("time_range", req.TimeRange),
			logging.String("metric_type", req.MetricType),
			logging.Error(err))
		s.metrics.RecordFailure("get_network_stats", "repository_error", map[string]string{
			"time_range":  req.TimeRange,
			"metric_type": req.MetricType,
		})
		return nil, status.Errorf(codes.Internal, "failed to get network stats: %v", err)
	}

	// Record success metrics
	s.logger.Info("Network stats retrieved successfully",
		logging.String("correlation_id", correlationID),
		logging.String("time_range", req.TimeRange),
		logging.String("metric_type", req.MetricType))
	s.metrics.RecordSuccess("get_network_stats", map[string]string{
		"time_range":  req.TimeRange,
		"metric_type": req.MetricType,
	})

	return response, nil
}

