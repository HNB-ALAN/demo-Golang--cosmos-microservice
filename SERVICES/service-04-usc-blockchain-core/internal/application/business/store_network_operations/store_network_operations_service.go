package store_network_operations

import (
	"context"
	"time"

	"service-04/internal/application/repository/store_network_operations"
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

// Service handles store network operations business logic
type Service struct {
	repo              *store_network_operations.Repository
	cosmosApp         *app.USCApp
	blockchainStorage *storage.StateManager
	logger            *logging.Logger
	validator         *validation.Validator
	metrics           *metrics.MetricsService
}

// NewService creates a new store network operations service
func NewService(
	repo *store_network_operations.Repository,
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

// SyncStoreNetworkState syncs external network state
func (s *Service) SyncStoreNetworkState(ctx context.Context, req *proto.SyncStoreNetworkStateRequest) (*proto.SyncStoreNetworkStateResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("sync_store_network_state", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Syncing store network state in business service",
		logging.String("correlation_id", correlationID),
		logging.String("networkId", req.NetworkId),
		logging.String("syncType", req.SyncType))

	// Input validation
	if req.NetworkId == "" {
		s.logger.Error("Network ID is required",
			logging.String("correlation_id", correlationID))
		s.metrics.RecordFailure("sync_store_network_state", "validation_error", map[string]string{
			"network_id": req.NetworkId,
		})
		return nil, status.Errorf(codes.InvalidArgument, "network_id is required")
	}

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.SyncStoreNetworkState(ctx, req)
	if err != nil {
		s.logger.Error("Failed to sync store network state in repository",
			logging.String("correlation_id", correlationID),
			logging.String("network_id", req.NetworkId),
			logging.String("sync_type", req.SyncType),
			logging.Error(err))
		s.metrics.RecordFailure("sync_store_network_state", "repository_error", map[string]string{
			"network_id": req.NetworkId,
			"sync_type":  req.SyncType,
		})
		return nil, status.Errorf(codes.Internal, "failed to sync store network state: %v", err)
	}

	// Record success metrics
	s.logger.Info("Store network state synced successfully",
		logging.String("correlation_id", correlationID),
		logging.String("network_id", req.NetworkId),
		logging.String("sync_type", req.SyncType))
	s.metrics.RecordSuccess("sync_store_network_state", map[string]string{
		"network_id": req.NetworkId,
		"sync_type":  req.SyncType,
	})

	return response, nil
}

// GetStoreNetworkInfo retrieves store network information
func (s *Service) GetStoreNetworkInfo(ctx context.Context, req *proto.GetStoreNetworkInfoRequest) (*proto.GetStoreNetworkInfoResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("get_store_network_info", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Getting store network info in business service",
		logging.String("correlation_id", correlationID),
		logging.String("networkId", req.NetworkId),
		logging.String("infoType", req.InfoType))

	// Input validation
	if req.NetworkId == "" {
		s.logger.Error("Network ID is required",
			logging.String("correlation_id", correlationID))
		s.metrics.RecordFailure("get_store_network_info", "validation_error", map[string]string{
			"network_id": req.NetworkId,
		})
		return nil, status.Errorf(codes.InvalidArgument, "network_id is required")
	}

	// Call repository
	response, err := s.repo.GetStoreNetworkInfo(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get store network info in repository",
			logging.String("correlation_id", correlationID),
			logging.String("network_id", req.NetworkId),
			logging.String("info_type", req.InfoType),
			logging.Error(err))
		s.metrics.RecordFailure("get_store_network_info", "repository_error", map[string]string{
			"network_id": req.NetworkId,
			"info_type":  req.InfoType,
		})
		return nil, status.Errorf(codes.Internal, "failed to get store network info: %v", err)
	}

	// Record success metrics
	s.logger.Info("Store network info retrieved successfully",
		logging.String("correlation_id", correlationID),
		logging.String("network_id", req.NetworkId),
		logging.String("info_type", req.InfoType))
	s.metrics.RecordSuccess("get_store_network_info", map[string]string{
		"network_id": req.NetworkId,
		"info_type":  req.InfoType,
	})

	return response, nil
}

// UpdateStoreBridgeConfig updates bridge configuration
func (s *Service) UpdateStoreBridgeConfig(ctx context.Context, req *proto.UpdateStoreBridgeConfigRequest) (*proto.UpdateStoreBridgeConfigResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("update_store_bridge_config", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Updating store bridge config in business service",
		logging.String("correlation_id", correlationID),
		logging.String("bridge", req.BridgeAddress),
		logging.String("configType", req.ConfigType))

	// Input validation using validator service
	if err := s.validator.ValidateContractAddress(req.BridgeAddress); err != nil {
		s.logger.Error("Bridge address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("bridge_address", req.BridgeAddress),
			logging.Error(err))
		s.metrics.RecordFailure("update_store_bridge_config", "validation_error", map[string]string{
			"bridge_address": req.BridgeAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid bridge_address: %v", err)
	}

	if req.ConfigData == "" {
		s.logger.Error("Config data is required",
			logging.String("correlation_id", correlationID))
		s.metrics.RecordFailure("update_store_bridge_config", "validation_error", map[string]string{
			"config_data": req.ConfigData,
		})
		return nil, status.Errorf(codes.InvalidArgument, "config_data is required")
	}

	// Call repository
	response, err := s.repo.UpdateStoreBridgeConfig(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update store bridge config in repository",
			logging.String("correlation_id", correlationID),
			logging.String("bridge_address", req.BridgeAddress),
			logging.String("config_type", req.ConfigType),
			logging.Error(err))
		s.metrics.RecordFailure("update_store_bridge_config", "repository_error", map[string]string{
			"bridge_address": req.BridgeAddress,
			"config_type":    req.ConfigType,
		})
		return nil, status.Errorf(codes.Internal, "failed to update store bridge config: %v", err)
	}

	// Record success metrics
	s.logger.Info("Store bridge config updated successfully",
		logging.String("correlation_id", correlationID),
		logging.String("bridge_address", req.BridgeAddress),
		logging.String("config_type", req.ConfigType))
	s.metrics.RecordSuccess("update_store_bridge_config", map[string]string{
		"bridge_address": req.BridgeAddress,
		"config_type":    req.ConfigType,
	})

	return response, nil
}
