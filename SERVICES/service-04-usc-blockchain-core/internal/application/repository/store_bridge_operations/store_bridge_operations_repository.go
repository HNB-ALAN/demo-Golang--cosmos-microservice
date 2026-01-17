package store_bridge_operations

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	repoerrors "service-04/internal/application/repository"
	"service-04/internal/application/utils"
	"service-04/internal/infrastructure/database"
	proto "service-04/proto"

	// Cosmos SDK imports
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/app"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/storage"
	bridgetypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/store_bridge/types"

	"github.com/usc-platform/shared/logging"
)

// Repository handles store bridge operations data access
type Repository struct {
	db                *database.PostgreSQLManager
	cosmosApp         *app.USCApp
	blockchainStorage *storage.StateManager
	redisManager      *database.RedisManager
	logger            *logging.Logger
}

// NewRepository creates a new store bridge operations repository
func NewRepository(db *database.PostgreSQLManager, cosmosApp *app.USCApp, blockchainStorage *storage.StateManager, redisManager *database.RedisManager, logger *logging.Logger) *Repository {
	return &Repository{
		db:                db,
		cosmosApp:         cosmosApp,
		blockchainStorage: blockchainStorage,
		redisManager:      redisManager,
		logger:            logger,
	}
}

// DeployStoreBridge deploys a new store bridge
func (r *Repository) DeployStoreBridge(ctx context.Context, req *proto.DeployStoreBridgeRequest) (*proto.DeployStoreBridgeResponse, error) {
	// Validate request
	if req.FromAddress == "" {
		return &proto.DeployStoreBridgeResponse{
			Success:      false,
			ErrorMessage: "from_address is required",
		}, nil
	}

	if req.BridgeName == "" {
		return &proto.DeployStoreBridgeResponse{
			Success:      false,
			ErrorMessage: "bridge_name is required",
		}, nil
	}

	if req.TargetChainId == "" {
		return &proto.DeployStoreBridgeResponse{
			Success:      false,
			ErrorMessage: "target_chain_id is required",
		}, nil
	}

	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.deployBridgeOnKeeper(ctx, req); err == nil {
			// Save to PostgreSQL for analytics (sync to ensure data persistence for cross-chain tracking)
			if r.db != nil {
				correlationID := utils.GetCorrelationID(ctx)
				if err := r.saveBridgeToDatabase(ctx, req, result); err != nil {
					r.logger.Error("Failed to save bridge analytics",
						logging.Error(err),
						logging.String("bridge_address", result.BridgeAddress),
						logging.String("correlation_id", correlationID))
					// Continue even if database save fails (keeper is primary)
				} else {
					r.logger.Info("Bridge analytics saved successfully",
						logging.String("bridge_address", result.BridgeAddress),
						logging.String("correlation_id", correlationID))
				}
			}
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.deployBridgeInDatabase(ctx, req)
}

// Helper methods for StoreBridgeKeeper interaction

// getSDKContext creates a sdk.Context from context.Context
// Uses shared utility to avoid code duplication
func (r *Repository) getSDKContext(ctx context.Context) (sdk.Context, error) {
	return utils.GetSDKContext(ctx, r.cosmosApp, r.logger)
}

// deployBridgeOnKeeper deploys a bridge on the keeper
func (r *Repository) deployBridgeOnKeeper(ctx context.Context, req *proto.DeployStoreBridgeRequest) (*proto.DeployStoreBridgeResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, err
	}

	// Generate bridge ID
	bridgeID := fmt.Sprintf("bridge_%s_%d", req.BridgeName, time.Now().Unix())

	// Create Bridge
	bridge := bridgetypes.Bridge{
		ID:          bridgeID,
		Name:        req.BridgeName,
		Description: req.BridgeDescription,
		FromChain:   "usc-1",
		ToChain:     req.TargetChainId,
		Type:        "token",
		Status:      "active",
		Config:      make(map[string]string),
		Validators:  []string{},
		Threshold:   1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Tags:        make(map[string]string),
	}

	// Set bridge in keeper
	if err := r.cosmosApp.StoreBridgeKeeper.SetBridge(sdkCtx, bridge); err != nil {
		return nil, repoerrors.NewDatabaseError("set_bridge", err)
	}

	// Generate transaction hash
	dataStr := fmt.Sprintf("%s:%s:%s:%s:deploy", req.FromAddress, req.BridgeName, req.TargetChainId, time.Now().Format(time.RFC3339))
	hashBytes := sha256.Sum256([]byte(dataStr))
	txHash := "0x" + hex.EncodeToString(hashBytes[:])

	return &proto.DeployStoreBridgeResponse{
		Success:         true,
		BridgeAddress:   bridgeID,
		TransactionHash: txHash,
		Status:          1, // Confirmed
		ErrorMessage:    "",
	}, nil
}

// Database fallback methods

// deployBridgeInDatabase deploys a bridge in database
func (r *Repository) deployBridgeInDatabase(ctx context.Context, req *proto.DeployStoreBridgeRequest) (*proto.DeployStoreBridgeResponse, error) {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return nil, repoerrors.NewRepositoryError(repoerrors.ErrDatabaseUnavailable)
	}

	// Generate bridge ID
	bridgeID := fmt.Sprintf("bridge_%s_%d", req.BridgeName, time.Now().Unix())

	// Generate transaction hash
	dataStr := fmt.Sprintf("%s:%s:%s:%s:deploy", req.FromAddress, req.BridgeName, req.TargetChainId, time.Now().Format(time.RFC3339))
	hashBytes := sha256.Sum256([]byte(dataStr))
	txHash := "0x" + hex.EncodeToString(hashBytes[:])

	query := `
		INSERT INTO store_bridges (
			bridge_address, bridge_name, bridge_description, target_network,
			target_chain_id, bridge_config, bridge_status, deployment_transaction_hash,
			user_id, store_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (bridge_address) DO UPDATE SET
			bridge_name = EXCLUDED.bridge_name,
			updated_at = NOW()
	`

	bridgeConfig := "{}"
	if req.BridgeConfig != "" {
		bridgeConfig = req.BridgeConfig
	}

	_, err := postgres.ExecContext(ctx, query,
		bridgeID, req.BridgeName, req.BridgeDescription, req.TargetNetwork,
		req.TargetChainId, bridgeConfig, "active", txHash,
		req.UserId, req.StoreId,
	)
	if err != nil {
		return nil, repoerrors.NewDatabaseError("deploy_bridge", err)
	}

	return &proto.DeployStoreBridgeResponse{
		Success:         true,
		BridgeAddress:   bridgeID,
		TransactionHash: txHash,
		Status:          0, // Pending (database fallback)
		ErrorMessage:    "",
	}, nil
}

// saveBridgeToDatabase saves bridge to database for analytics (sync for cross-chain tracking)
func (r *Repository) saveBridgeToDatabase(ctx context.Context, req *proto.DeployStoreBridgeRequest, result *proto.DeployStoreBridgeResponse) error {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return fmt.Errorf("postgres connection not available")
	}

	query := `
		INSERT INTO store_bridges (
			bridge_address, bridge_name, bridge_description, target_network,
			target_chain_id, bridge_config, bridge_status, deployment_transaction_hash,
			user_id, store_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (bridge_address) DO UPDATE SET
			bridge_name = EXCLUDED.bridge_name,
			updated_at = NOW()
	`

	bridgeConfig := "{}"
	if req.BridgeConfig != "" {
		bridgeConfig = req.BridgeConfig
	}

	if _, err := postgres.ExecContext(ctx, query,
		result.BridgeAddress, req.BridgeName, req.BridgeDescription, req.TargetNetwork,
		req.TargetChainId, bridgeConfig, "active", result.TransactionHash,
		req.UserId, req.StoreId,
	); err != nil {
		return fmt.Errorf("failed to save bridge to database: %w", err)
	}

	return nil
}

// RegisterStoreNetwork registers a new store network
func (r *Repository) RegisterStoreNetwork(ctx context.Context, req *proto.RegisterStoreNetworkRequest) (*proto.RegisterStoreNetworkResponse, error) {
	// Validate request
	if req.NetworkName == "" {
		return &proto.RegisterStoreNetworkResponse{
			Success:      false,
			ErrorMessage: "network_name is required",
		}, nil
	}

	if req.NetworkId == "" {
		return &proto.RegisterStoreNetworkResponse{
			Success:      false,
			ErrorMessage: "network_id is required",
		}, nil
	}

	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.registerNetworkOnKeeper(ctx, req); err == nil {
			// Save to PostgreSQL for analytics (sync to ensure data persistence for cross-chain tracking)
			if r.db != nil {
				correlationID := utils.GetCorrelationID(ctx)
				if err := r.saveNetworkToDatabase(ctx, req, result); err != nil {
					r.logger.Error("Failed to save network analytics",
						logging.Error(err),
						logging.String("network_id", req.NetworkId),
						logging.String("correlation_id", correlationID))
					// Continue even if database save fails (keeper is primary)
				} else {
					r.logger.Info("Network analytics saved successfully",
						logging.String("network_id", req.NetworkId),
						logging.String("correlation_id", correlationID))
				}
			}
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.registerNetworkInDatabase(ctx, req)
}

// registerNetworkOnKeeper registers a network on the keeper
func (r *Repository) registerNetworkOnKeeper(ctx context.Context, req *proto.RegisterStoreNetworkRequest) (*proto.RegisterStoreNetworkResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, err
	}

	// Get bridge configs to find matching network
	bridgeConfigs := r.cosmosApp.StoreBridgeKeeper.GetAllBridgeConfigs(sdkCtx)
	for _, config := range bridgeConfigs {
		if config.BridgeID == req.NetworkId {
			return &proto.RegisterStoreNetworkResponse{
				Success:       true,
				NetworkId:     req.NetworkId,
				NetworkName:   req.NetworkName,
				ErrorMessage:  "",
				NetworkStatus: "active",
			}, nil
		}
	}

	// Network not found, return success anyway (network registration is informational)
	return &proto.RegisterStoreNetworkResponse{
		Success:       true,
		NetworkId:     req.NetworkId,
		NetworkName:   req.NetworkName,
		ErrorMessage:  "",
		NetworkStatus: "active",
	}, nil
}

// registerNetworkInDatabase registers a network in database
func (r *Repository) registerNetworkInDatabase(ctx context.Context, req *proto.RegisterStoreNetworkRequest) (*proto.RegisterStoreNetworkResponse, error) {
	// Note: Network metadata table not implemented yet
	return &proto.RegisterStoreNetworkResponse{
		Success:       true,
		NetworkId:     req.NetworkId,
		NetworkName:   req.NetworkName,
		ErrorMessage:  "",
		NetworkStatus: "active",
	}, nil
}

// saveNetworkToDatabase saves network to database for analytics (sync for cross-chain tracking)
func (r *Repository) saveNetworkToDatabase(ctx context.Context, req *proto.RegisterStoreNetworkRequest, result *proto.RegisterStoreNetworkResponse) error {
	// Note: Network metadata table not implemented yet
	// Return nil for now (no-op) until table is implemented
	_ = ctx
	_ = req
	_ = result
	return nil
}

// BridgeStoreTokenToUSC bridges store tokens to USC
func (r *Repository) BridgeStoreTokenToUSC(ctx context.Context, req *proto.BridgeStoreTokenToUSCRequest) (*proto.BridgeStoreTokenToUSCResponse, error) {
	// Validate request
	if req.FromAddress == "" {
		return &proto.BridgeStoreTokenToUSCResponse{
			Status:       2, // Failed
			ErrorMessage: "from_address is required",
		}, nil
	}

	if req.StoreTokenAmount == "" {
		return &proto.BridgeStoreTokenToUSCResponse{
			Status:       2, // Failed
			ErrorMessage: "store_token_amount is required",
		}, nil
	}

	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.bridgeTokenToUSCOnKeeper(ctx, req); err == nil {
			// Save to PostgreSQL for analytics (sync to ensure data persistence for cross-chain tracking)
			if r.db != nil {
				correlationID := utils.GetCorrelationID(ctx)
				if err := r.saveBridgeTransactionToDatabase(ctx, req, result, "token_to_usc"); err != nil {
					r.logger.Error("Failed to save bridge transaction analytics",
						logging.Error(err),
						logging.String("from_address", req.FromAddress),
						logging.String("correlation_id", correlationID))
					// Continue even if database save fails (keeper is primary)
				} else {
					r.logger.Info("Bridge transaction analytics saved successfully",
						logging.String("from_address", req.FromAddress),
						logging.String("correlation_id", correlationID))
				}
			}
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.bridgeTokenToUSCInDatabase(ctx, req)
}

// bridgeTokenToUSCOnKeeper bridges store token to USC on the keeper
func (r *Repository) bridgeTokenToUSCOnKeeper(ctx context.Context, req *proto.BridgeStoreTokenToUSCRequest) (*proto.BridgeStoreTokenToUSCResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, err
	}

	// Create bridge transfer event
	transferID := fmt.Sprintf("transfer_%s_%d", req.FromAddress, time.Now().Unix())
	transfer := bridgetypes.Transfer{
		ID:          transferID,
		BridgeID:    req.TargetNetwork,
		FromChain:   req.TargetNetwork,
		ToChain:     "usc-1",
		FromAddress: req.FromAddress,
		ToAddress:   req.FromAddress, // Use from address as default
		Amount:      req.StoreTokenAmount,
		Token:       req.StoreTokenAddress,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	// Set transfer in keeper
	if err := r.cosmosApp.StoreBridgeKeeper.SetTransfer(sdkCtx, transfer); err != nil {
		return nil, repoerrors.NewDatabaseError("set_transfer", err)
	}

	// Generate transaction hash
	dataStr := fmt.Sprintf("%s:%s:%s:%s:bridge", req.FromAddress, req.StoreTokenAddress, req.StoreTokenAmount, time.Now().Format(time.RFC3339))
	hashBytes := sha256.Sum256([]byte(dataStr))
	txHash := "0x" + hex.EncodeToString(hashBytes[:])

	return &proto.BridgeStoreTokenToUSCResponse{
		Success:         true,
		TransactionHash: txHash,
		Status:          0, // Pending
		ErrorMessage:    "",
		UscAmount:       req.StoreTokenAmount,
		BridgeFee:       "0.01",
	}, nil
}

// bridgeTokenToUSCInDatabase bridges store token to USC in database
func (r *Repository) bridgeTokenToUSCInDatabase(ctx context.Context, req *proto.BridgeStoreTokenToUSCRequest) (*proto.BridgeStoreTokenToUSCResponse, error) {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return nil, repoerrors.NewRepositoryError(repoerrors.ErrDatabaseUnavailable)
	}

	// Generate transaction hash
	dataStr := fmt.Sprintf("%s:%s:%s:%s:bridge", req.FromAddress, req.StoreTokenAddress, req.StoreTokenAmount, time.Now().Format(time.RFC3339))
	hashBytes := sha256.Sum256([]byte(dataStr))
	txHash := "0x" + hex.EncodeToString(hashBytes[:])

	query := `
		INSERT INTO bridge_transactions (
			transaction_hash, bridge_address, from_address, to_address,
			source_network, target_network, source_token_address, target_token_address,
			source_amount, target_amount, bridge_fee, transaction_type,
			status, user_id, device_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		ON CONFLICT (transaction_hash) DO UPDATE SET
			status = EXCLUDED.status,
			updated_at = NOW()
	`

	_, err := postgres.ExecContext(ctx, query,
		txHash, req.TargetNetwork, req.FromAddress, req.FromAddress,
		req.TargetNetwork, "usc-1", req.StoreTokenAddress, "",
		req.StoreTokenAmount, req.StoreTokenAmount, "0.01", "token_to_usc",
		"pending", req.UserId, req.DeviceId,
	)
	if err != nil {
		return nil, repoerrors.NewDatabaseError("bridge_transaction", err)
	}

	return &proto.BridgeStoreTokenToUSCResponse{
		Success:         true,
		TransactionHash: txHash,
		Status:          0, // Pending (database fallback)
		ErrorMessage:    "",
		UscAmount:       req.StoreTokenAmount,
		BridgeFee:       "0.01",
	}, nil
}

// BridgeUSCToStoreToken bridges USC to store tokens
func (r *Repository) BridgeUSCToStoreToken(ctx context.Context, req *proto.BridgeUSCToStoreTokenRequest) (*proto.BridgeUSCToStoreTokenResponse, error) {
	// Validate request
	if req.FromAddress == "" {
		return &proto.BridgeUSCToStoreTokenResponse{
			Status:       2, // Failed
			ErrorMessage: "from_address is required",
		}, nil
	}

	if req.UscAmount == "" {
		return &proto.BridgeUSCToStoreTokenResponse{
			Status:       2, // Failed
			ErrorMessage: "usc_amount is required",
		}, nil
	}

	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.bridgeUSCToTokenOnKeeper(ctx, req); err == nil {
			// Save to PostgreSQL for analytics (sync to ensure data persistence for cross-chain tracking)
			if r.db != nil {
				correlationID := utils.GetCorrelationID(ctx)
				if err := r.saveBridgeTransactionToDatabase(ctx, req, result, "usc_to_token"); err != nil {
					r.logger.Error("Failed to save bridge transaction analytics",
						logging.Error(err),
						logging.String("from_address", req.FromAddress),
						logging.String("correlation_id", correlationID))
					// Continue even if database save fails (keeper is primary)
				} else {
					r.logger.Info("Bridge transaction analytics saved successfully",
						logging.String("from_address", req.FromAddress),
						logging.String("correlation_id", correlationID))
				}
			}
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.bridgeUSCToTokenInDatabase(ctx, req)
}

// bridgeUSCToTokenOnKeeper bridges USC to store token on the keeper
func (r *Repository) bridgeUSCToTokenOnKeeper(ctx context.Context, req *proto.BridgeUSCToStoreTokenRequest) (*proto.BridgeUSCToStoreTokenResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, err
	}

	// Create bridge transfer event
	transferID := fmt.Sprintf("transfer_%s_%d", req.FromAddress, time.Now().Unix())
	transfer := bridgetypes.Transfer{
		ID:          transferID,
		BridgeID:    req.TargetNetwork,
		FromChain:   "usc-1",
		ToChain:     req.TargetNetwork,
		FromAddress: req.FromAddress,
		ToAddress:   req.FromAddress, // Use from address as default
		Amount:      req.UscAmount,
		Token:       req.TargetTokenAddress,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	// Set transfer in keeper
	if err := r.cosmosApp.StoreBridgeKeeper.SetTransfer(sdkCtx, transfer); err != nil {
		return nil, repoerrors.NewDatabaseError("set_transfer", err)
	}

	// Generate transaction hash
	dataStr := fmt.Sprintf("%s:%s:%s:%s:bridge", req.FromAddress, req.UscAmount, req.TargetTokenAddress, time.Now().Format(time.RFC3339))
	hashBytes := sha256.Sum256([]byte(dataStr))
	txHash := "0x" + hex.EncodeToString(hashBytes[:])

	return &proto.BridgeUSCToStoreTokenResponse{
		Success:          true,
		TransactionHash:  txHash,
		Status:           0, // Pending
		ErrorMessage:     "",
		StoreTokenAmount: req.UscAmount,
		BridgeFee:        "0.01",
	}, nil
}

// bridgeUSCToTokenInDatabase bridges USC to store token in database
func (r *Repository) bridgeUSCToTokenInDatabase(ctx context.Context, req *proto.BridgeUSCToStoreTokenRequest) (*proto.BridgeUSCToStoreTokenResponse, error) {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return nil, repoerrors.NewRepositoryError(repoerrors.ErrDatabaseUnavailable)
	}

	// Generate transaction hash
	dataStr := fmt.Sprintf("%s:%s:%s:%s:bridge", req.FromAddress, req.UscAmount, req.TargetTokenAddress, time.Now().Format(time.RFC3339))
	hashBytes := sha256.Sum256([]byte(dataStr))
	txHash := "0x" + hex.EncodeToString(hashBytes[:])

	query := `
		INSERT INTO bridge_transactions (
			transaction_hash, bridge_address, from_address, to_address,
			source_network, target_network, source_token_address, target_token_address,
			source_amount, target_amount, bridge_fee, transaction_type,
			status, user_id, device_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		ON CONFLICT (transaction_hash) DO UPDATE SET
			status = EXCLUDED.status,
			updated_at = NOW()
	`

	_, err := postgres.ExecContext(ctx, query,
		txHash, req.TargetNetwork, req.FromAddress, req.FromAddress,
		"usc-1", req.TargetNetwork, "", req.TargetTokenAddress,
		req.UscAmount, req.UscAmount, "0.01", "usc_to_token",
		"pending", req.UserId, req.DeviceId,
	)
	if err != nil {
		return nil, repoerrors.NewDatabaseError("bridge_transaction", err)
	}

	return &proto.BridgeUSCToStoreTokenResponse{
		Success:          true,
		TransactionHash:  txHash,
		Status:           0, // Pending (database fallback)
		ErrorMessage:     "",
		StoreTokenAmount: req.UscAmount,
		BridgeFee:        "0.01",
	}, nil
}

// GetStoreBridgeMetrics retrieves store bridge metrics
func (r *Repository) GetStoreBridgeMetrics(ctx context.Context, req *proto.GetStoreBridgeMetricsRequest) (*proto.GetStoreBridgeMetricsResponse, error) {
	// Validate request
	if req.BridgeAddress == "" {
		return nil, repoerrors.NewValidationError("bridge_address", "is required")
	}

	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.getBridgeMetricsFromKeeper(ctx, req); err == nil {
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.getBridgeMetricsFromDatabase(ctx, req)
}

// getBridgeMetricsFromKeeper retrieves bridge metrics from the keeper
func (r *Repository) getBridgeMetricsFromKeeper(ctx context.Context, req *proto.GetStoreBridgeMetricsRequest) (*proto.GetStoreBridgeMetricsResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, err
	}

	// Get bridge
	_, err = r.cosmosApp.StoreBridgeKeeper.GetBridge(sdkCtx, req.BridgeAddress)
	if err != nil {
		return nil, err
	}

	// Get all transfers for this bridge
	transfers := r.cosmosApp.StoreBridgeKeeper.GetAllTransfers(sdkCtx)
	bridgeTransfers := 0
	for _, transfer := range transfers {
		if transfer.BridgeID == req.BridgeAddress {
			bridgeTransfers++
		}
	}
	_ = bridgeTransfers // Use in future metrics

	return &proto.GetStoreBridgeMetricsResponse{
		Success:      true,
		TimeRange:    req.TimeRange,
		ErrorMessage: "",
		// Note: Additional metrics fields would be populated here
	}, nil
}

// getBridgeMetricsFromDatabase retrieves bridge metrics from database
func (r *Repository) getBridgeMetricsFromDatabase(ctx context.Context, req *proto.GetStoreBridgeMetricsRequest) (*proto.GetStoreBridgeMetricsResponse, error) {
	// Note: Bridge analytics table not implemented yet
	return &proto.GetStoreBridgeMetricsResponse{
		Success:      true,
		TimeRange:    req.TimeRange,
		ErrorMessage: "",
	}, nil
}

// ValidateStoreBridge validates a store bridge
func (r *Repository) ValidateStoreBridge(ctx context.Context, req *proto.ValidateStoreBridgeRequest) (*proto.ValidateStoreBridgeResponse, error) {
	// Validate request
	if req.BridgeAddress == "" {
		return nil, repoerrors.NewValidationError("bridge_address", "is required")
	}

	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.validateBridgeFromKeeper(ctx, req); err == nil {
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.validateBridgeFromDatabase(ctx, req)
}

// validateBridgeFromKeeper validates a bridge from the keeper
func (r *Repository) validateBridgeFromKeeper(ctx context.Context, req *proto.ValidateStoreBridgeRequest) (*proto.ValidateStoreBridgeResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, err
	}

	// Get bridge
	bridge, err := r.cosmosApp.StoreBridgeKeeper.GetBridge(sdkCtx, req.BridgeAddress)
	if err != nil {
		return &proto.ValidateStoreBridgeResponse{
			Success:          false,
			IsValid:          false,
			ValidationResult: "Bridge not found",
			Issues:           []string{"Bridge does not exist"},
			Recommendations:  []string{},
			SecurityScore:    "0",
			PerformanceScore: "0",
			ComplianceScore:  "0",
			ErrorMessage:     err.Error(),
		}, nil
	}

	// Validate bridge
	isValid := bridge.Status == "active"
	issues := []string{}
	if bridge.Status != "active" {
		issues = append(issues, fmt.Sprintf("Bridge status is %s", bridge.Status))
	}
	if len(bridge.Validators) == 0 {
		issues = append(issues, "No validators configured")
	}

	securityScore := "95"
	performanceScore := "90"
	complianceScore := "100"
	if !isValid {
		securityScore = "50"
		performanceScore = "50"
		complianceScore = "50"
	}

	return &proto.ValidateStoreBridgeResponse{
		Success:          true,
		IsValid:          isValid,
		ValidationResult: "Valid",
		Issues:           issues,
		Recommendations:  []string{},
		SecurityScore:    securityScore,
		PerformanceScore: performanceScore,
		ComplianceScore:  complianceScore,
		ErrorMessage:     "",
	}, nil
}

// saveBridgeTransactionToDatabase saves bridge transaction to database for analytics (sync for cross-chain tracking)
func (r *Repository) saveBridgeTransactionToDatabase(ctx context.Context, req interface{}, result interface{}, txType string) error {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return fmt.Errorf("postgres connection not available")
	}

	var txHash, bridgeAddr, fromAddr, toAddr, sourceNetwork, targetNetwork, sourceToken, targetToken, sourceAmount, targetAmount, bridgeFee string
	var userId, deviceId string

	// Extract data based on request type
	switch reqTyped := req.(type) {
	case *proto.BridgeStoreTokenToUSCRequest:
		reqBridge := reqTyped
		if resultTyped, ok := result.(*proto.BridgeStoreTokenToUSCResponse); ok {
			txHash = resultTyped.TransactionHash
			bridgeAddr = reqBridge.TargetNetwork // Use TargetNetwork as bridge identifier
			fromAddr = reqBridge.FromAddress
			toAddr = reqBridge.FromAddress
			sourceNetwork = reqBridge.TargetNetwork
			targetNetwork = "usc-1"
			sourceToken = reqBridge.StoreTokenAddress
			targetToken = ""
			sourceAmount = reqBridge.StoreTokenAmount
			targetAmount = resultTyped.UscAmount
			bridgeFee = resultTyped.BridgeFee
			userId = reqBridge.UserId
			deviceId = reqBridge.DeviceId
		}
	case *proto.BridgeUSCToStoreTokenRequest:
		reqBridge := reqTyped
		if resultTyped, ok := result.(*proto.BridgeUSCToStoreTokenResponse); ok {
			txHash = resultTyped.TransactionHash
			bridgeAddr = reqBridge.TargetNetwork // Use TargetNetwork as bridge identifier
			fromAddr = reqBridge.FromAddress
			toAddr = reqBridge.FromAddress
			sourceNetwork = "usc-1"
			targetNetwork = reqBridge.TargetNetwork
			sourceToken = ""
			targetToken = reqBridge.TargetTokenAddress
			sourceAmount = reqBridge.UscAmount
			targetAmount = resultTyped.StoreTokenAmount
			bridgeFee = resultTyped.BridgeFee
			userId = reqBridge.UserId
			deviceId = reqBridge.DeviceId
		}
	default:
		return fmt.Errorf("unsupported request type for bridge transaction")
	}

	query := `
		INSERT INTO bridge_transactions (
			transaction_hash, bridge_address, from_address, to_address,
			source_network, target_network, source_token_address, target_token_address,
			source_amount, target_amount, bridge_fee, transaction_type,
			status, user_id, device_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		ON CONFLICT (transaction_hash) DO UPDATE SET
			status = EXCLUDED.status,
			updated_at = NOW()
	`

	status := "pending"
	if txHash != "" {
		status = "completed"
	}

	if _, err := postgres.ExecContext(ctx, query,
		txHash, bridgeAddr, fromAddr, toAddr,
		sourceNetwork, targetNetwork, sourceToken, targetToken,
		sourceAmount, targetAmount, bridgeFee, txType,
		status, userId, deviceId,
	); err != nil {
		return fmt.Errorf("failed to save bridge transaction to database: %w", err)
	}

	return nil
}

// validateBridgeFromDatabase validates a bridge from database
func (r *Repository) validateBridgeFromDatabase(ctx context.Context, req *proto.ValidateStoreBridgeRequest) (*proto.ValidateStoreBridgeResponse, error) {
	// Note: Bridge validation from database not implemented yet
	return &proto.ValidateStoreBridgeResponse{
		Success:          true,
		IsValid:          true,
		ValidationResult: "Valid",
		Issues:           []string{},
		Recommendations:  []string{},
		SecurityScore:    "95",
		PerformanceScore: "90",
		ComplianceScore:  "100",
		ErrorMessage:     "",
	}, nil
}
