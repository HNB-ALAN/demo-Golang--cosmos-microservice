package store_bridge_operations

import (
	"context"
	"time"

	"service-04/internal/application/repository/store_bridge_operations"
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

// Service handles store bridge operations business logic
type Service struct {
	repo              *store_bridge_operations.Repository
	cosmosApp         *app.USCApp
	blockchainStorage *storage.StateManager
	logger            *logging.Logger
	validator         *validation.Validator
	metrics           *metrics.MetricsService
}

// NewService creates a new store bridge operations service
func NewService(
	repo *store_bridge_operations.Repository,
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

// DeployStoreBridge deploys a new store bridge
func (s *Service) DeployStoreBridge(ctx context.Context, req *proto.DeployStoreBridgeRequest) (*proto.DeployStoreBridgeResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("deploy_store_bridge", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Deploying store bridge in business service",
		logging.String("correlation_id", correlationID),
		logging.String("from", req.FromAddress),
		logging.String("bridgeName", req.BridgeName))

	// Input validation using validator service
	if err := s.validator.ValidateWalletAddress(req.FromAddress); err != nil {
		s.logger.Error("From address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("from_address", req.FromAddress),
			logging.Error(err))
		s.metrics.RecordFailure("deploy_store_bridge", "validation_error", map[string]string{
			"from_address": req.FromAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid from_address: %v", err)
	}

	if req.BridgeName == "" {
		s.logger.Error("Bridge name is required",
			logging.String("correlation_id", correlationID))
		s.metrics.RecordFailure("deploy_store_bridge", "validation_error", map[string]string{
			"bridge_name": req.BridgeName,
		})
		return nil, status.Errorf(codes.InvalidArgument, "bridge_name is required")
	}

	// Delegate to repository (repository handles Keeper → Database fallback)
	response, err := s.repo.DeployStoreBridge(ctx, req)
	if err != nil {
		s.logger.Error("Failed to deploy store bridge in repository",
			logging.String("correlation_id", correlationID),
			logging.String("from_address", req.FromAddress),
			logging.String("bridge_name", req.BridgeName),
			logging.Error(err))
		s.metrics.RecordFailure("deploy_store_bridge", "repository_error", map[string]string{
			"from_address": req.FromAddress,
			"bridge_name":  req.BridgeName,
		})
		return nil, status.Errorf(codes.Internal, "failed to deploy store bridge: %v", err)
	}

	// Record success metrics
	s.logger.Info("Store bridge deployed successfully",
		logging.String("correlation_id", correlationID),
		logging.String("from_address", req.FromAddress),
		logging.String("bridge_name", req.BridgeName))
	s.metrics.RecordSuccess("deploy_store_bridge", map[string]string{
		"from_address": req.FromAddress,
		"bridge_name":  req.BridgeName,
	})

	// Record blockchain-specific metric if bridge was deployed
	if response != nil && response.BridgeAddress != "" {
		s.metrics.RecordContractDeployed(response.BridgeAddress)
	}

	return response, nil
}

// RegisterStoreNetwork registers a new store network
func (s *Service) RegisterStoreNetwork(ctx context.Context, req *proto.RegisterStoreNetworkRequest) (*proto.RegisterStoreNetworkResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("register_store_network", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Registering store network in business service",
		logging.String("correlation_id", correlationID),
		logging.String("networkName", req.NetworkName),
		logging.String("networkId", req.NetworkId))

	// Input validation
	if req.NetworkName == "" {
		s.logger.Error("Network name is required",
			logging.String("correlation_id", correlationID))
		s.metrics.RecordFailure("register_store_network", "validation_error", map[string]string{
			"network_name": req.NetworkName,
		})
		return nil, status.Errorf(codes.InvalidArgument, "network_name is required")
	}

	if req.NetworkId == "" {
		s.logger.Error("Network ID is required",
			logging.String("correlation_id", correlationID))
		s.metrics.RecordFailure("register_store_network", "validation_error", map[string]string{
			"network_id": req.NetworkId,
		})
		return nil, status.Errorf(codes.InvalidArgument, "network_id is required")
	}

	// Call repository
	response, err := s.repo.RegisterStoreNetwork(ctx, req)
	if err != nil {
		s.logger.Error("Failed to register store network in repository",
			logging.String("correlation_id", correlationID),
			logging.String("network_name", req.NetworkName),
			logging.String("network_id", req.NetworkId),
			logging.Error(err))
		s.metrics.RecordFailure("register_store_network", "repository_error", map[string]string{
			"network_name": req.NetworkName,
			"network_id":   req.NetworkId,
		})
		return nil, status.Errorf(codes.Internal, "failed to register store network: %v", err)
	}

	// Record success metrics
	s.logger.Info("Store network registered successfully",
		logging.String("correlation_id", correlationID),
		logging.String("network_name", req.NetworkName),
		logging.String("network_id", req.NetworkId))
	s.metrics.RecordSuccess("register_store_network", map[string]string{
		"network_name": req.NetworkName,
		"network_id":   req.NetworkId,
	})

	return response, nil
}

// BridgeStoreTokenToUSC bridges store tokens to USC
func (s *Service) BridgeStoreTokenToUSC(ctx context.Context, req *proto.BridgeStoreTokenToUSCRequest) (*proto.BridgeStoreTokenToUSCResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("bridge_store_token_to_usc", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Bridging store token to USC in business service",
		logging.String("correlation_id", correlationID),
		logging.String("from", req.FromAddress),
		logging.String("amount", req.StoreTokenAmount))

	// Input validation using validator service
	if err := s.validator.ValidateWalletAddress(req.FromAddress); err != nil {
		s.logger.Error("From address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("from_address", req.FromAddress),
			logging.Error(err))
		s.metrics.RecordFailure("bridge_store_token_to_usc", "validation_error", map[string]string{
			"from_address": req.FromAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid from_address: %v", err)
	}

	if err := s.validator.ValidateAmount(req.StoreTokenAmount); err != nil {
		s.logger.Error("Store token amount validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("store_token_amount", req.StoreTokenAmount),
			logging.Error(err))
		s.metrics.RecordFailure("bridge_store_token_to_usc", "validation_error", map[string]string{
			"store_token_amount": req.StoreTokenAmount,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid store_token_amount: %v", err)
	}

	// Call repository
	response, err := s.repo.BridgeStoreTokenToUSC(ctx, req)
	if err != nil {
		s.logger.Error("Failed to bridge store token to USC in repository",
			logging.String("correlation_id", correlationID),
			logging.String("from_address", req.FromAddress),
			logging.String("store_token_amount", req.StoreTokenAmount),
			logging.Error(err))
		s.metrics.RecordFailure("bridge_store_token_to_usc", "repository_error", map[string]string{
			"from_address":       req.FromAddress,
			"store_token_amount": req.StoreTokenAmount,
		})
		return nil, status.Errorf(codes.Internal, "failed to bridge store token to USC: %v", err)
	}

	// Record success metrics
	s.logger.Info("Store token bridged to USC successfully",
		logging.String("correlation_id", correlationID),
		logging.String("from_address", req.FromAddress),
		logging.String("store_token_amount", req.StoreTokenAmount))
	s.metrics.RecordSuccess("bridge_store_token_to_usc", map[string]string{
		"from_address":       req.FromAddress,
		"store_token_amount": req.StoreTokenAmount,
	})

	return response, nil
}

// BridgeUSCToStoreToken bridges USC to store tokens
func (s *Service) BridgeUSCToStoreToken(ctx context.Context, req *proto.BridgeUSCToStoreTokenRequest) (*proto.BridgeUSCToStoreTokenResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("bridge_usc_to_store_token", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Bridging USC to store token in business service",
		logging.String("correlation_id", correlationID),
		logging.String("from", req.FromAddress),
		logging.String("amount", req.UscAmount))

	// Input validation using validator service
	if err := s.validator.ValidateWalletAddress(req.FromAddress); err != nil {
		s.logger.Error("From address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("from_address", req.FromAddress),
			logging.Error(err))
		s.metrics.RecordFailure("bridge_usc_to_store_token", "validation_error", map[string]string{
			"from_address": req.FromAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid from_address: %v", err)
	}

	if err := s.validator.ValidateAmount(req.UscAmount); err != nil {
		s.logger.Error("USC amount validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("usc_amount", req.UscAmount),
			logging.Error(err))
		s.metrics.RecordFailure("bridge_usc_to_store_token", "validation_error", map[string]string{
			"usc_amount": req.UscAmount,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid usc_amount: %v", err)
	}

	// Call repository
	response, err := s.repo.BridgeUSCToStoreToken(ctx, req)
	if err != nil {
		s.logger.Error("Failed to bridge USC to store token in repository",
			logging.String("correlation_id", correlationID),
			logging.String("from_address", req.FromAddress),
			logging.String("usc_amount", req.UscAmount),
			logging.Error(err))
		s.metrics.RecordFailure("bridge_usc_to_store_token", "repository_error", map[string]string{
			"from_address": req.FromAddress,
			"usc_amount":   req.UscAmount,
		})
		return nil, status.Errorf(codes.Internal, "failed to bridge USC to store token: %v", err)
	}

	// Record success metrics
	s.logger.Info("USC bridged to store token successfully",
		logging.String("correlation_id", correlationID),
		logging.String("from_address", req.FromAddress),
		logging.String("usc_amount", req.UscAmount))
	s.metrics.RecordSuccess("bridge_usc_to_store_token", map[string]string{
		"from_address": req.FromAddress,
		"usc_amount":   req.UscAmount,
	})

	return response, nil
}

// GetStoreBridgeMetrics retrieves store bridge metrics
func (s *Service) GetStoreBridgeMetrics(ctx context.Context, req *proto.GetStoreBridgeMetricsRequest) (*proto.GetStoreBridgeMetricsResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("get_store_bridge_metrics", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Getting store bridge metrics in business service",
		logging.String("correlation_id", correlationID),
		logging.String("bridge", req.BridgeAddress),
		logging.String("timeRange", req.TimeRange))

	// Input validation using validator service
	if err := s.validator.ValidateContractAddress(req.BridgeAddress); err != nil {
		s.logger.Error("Bridge address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("bridge_address", req.BridgeAddress),
			logging.Error(err))
		s.metrics.RecordFailure("get_store_bridge_metrics", "validation_error", map[string]string{
			"bridge_address": req.BridgeAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid bridge_address: %v", err)
	}

	// Call repository
	response, err := s.repo.GetStoreBridgeMetrics(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get store bridge metrics in repository",
			logging.String("correlation_id", correlationID),
			logging.String("bridge_address", req.BridgeAddress),
			logging.Error(err))
		s.metrics.RecordFailure("get_store_bridge_metrics", "repository_error", map[string]string{
			"bridge_address": req.BridgeAddress,
		})
		return nil, status.Errorf(codes.Internal, "failed to get store bridge metrics: %v", err)
	}

	// Record success metrics
	s.logger.Info("Store bridge metrics retrieved successfully",
		logging.String("correlation_id", correlationID),
		logging.String("bridge_address", req.BridgeAddress))
	s.metrics.RecordSuccess("get_store_bridge_metrics", map[string]string{
		"bridge_address": req.BridgeAddress,
	})

	return response, nil
}

// ValidateStoreBridge validates a store bridge
func (s *Service) ValidateStoreBridge(ctx context.Context, req *proto.ValidateStoreBridgeRequest) (*proto.ValidateStoreBridgeResponse, error) {
	start := time.Now()
	defer func() {
		s.metrics.RecordDuration("validate_store_bridge", time.Since(start))
	}()

	correlationID := utils.GetCorrelationID(ctx)
	s.logger.Info("Validating store bridge in business service",
		logging.String("correlation_id", correlationID),
		logging.String("bridge", req.BridgeAddress),
		logging.String("type", req.ValidationType))

	// Input validation using validator service
	if err := s.validator.ValidateContractAddress(req.BridgeAddress); err != nil {
		s.logger.Error("Bridge address validation failed",
			logging.String("correlation_id", correlationID),
			logging.String("bridge_address", req.BridgeAddress),
			logging.Error(err))
		s.metrics.RecordFailure("validate_store_bridge", "validation_error", map[string]string{
			"bridge_address": req.BridgeAddress,
		})
		return nil, status.Errorf(codes.InvalidArgument, "invalid bridge_address: %v", err)
	}

	// Call repository
	response, err := s.repo.ValidateStoreBridge(ctx, req)
	if err != nil {
		s.logger.Error("Failed to validate store bridge in repository",
			logging.String("correlation_id", correlationID),
			logging.String("bridge_address", req.BridgeAddress),
			logging.Error(err))
		s.metrics.RecordFailure("validate_store_bridge", "repository_error", map[string]string{
			"bridge_address": req.BridgeAddress,
		})
		return nil, status.Errorf(codes.Internal, "failed to validate store bridge: %v", err)
	}

	// Record success metrics
	s.logger.Info("Store bridge validated successfully",
		logging.String("correlation_id", correlationID),
		logging.String("bridge_address", req.BridgeAddress))
	s.metrics.RecordSuccess("validate_store_bridge", map[string]string{
		"bridge_address": req.BridgeAddress,
	})

	return response, nil
}
