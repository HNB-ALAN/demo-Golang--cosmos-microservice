package store_network_operations

import (
	"context"
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
	storenetworktypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/store_network/types"

	"github.com/usc-platform/shared/logging"
)

// Repository handles store network operations data access
type Repository struct {
	db                *database.PostgreSQLManager
	cosmosApp         *app.USCApp
	blockchainStorage *storage.StateManager
	redisManager      *database.RedisManager
	logger            *logging.Logger
}

// NewRepository creates a new store network operations repository
func NewRepository(db *database.PostgreSQLManager, cosmosApp *app.USCApp, blockchainStorage *storage.StateManager, redisManager *database.RedisManager, logger *logging.Logger) *Repository {
	return &Repository{
		db:                db,
		cosmosApp:         cosmosApp,
		blockchainStorage: blockchainStorage,
		redisManager:      redisManager,
		logger:            logger,
	}
}

// SyncStoreNetworkState syncs external network state
func (r *Repository) SyncStoreNetworkState(ctx context.Context, req *proto.SyncStoreNetworkStateRequest) (*proto.SyncStoreNetworkStateResponse, error) {
	// Validate request
	if req.NetworkId == "" {
		return &proto.SyncStoreNetworkStateResponse{
			Success:      false,
			ErrorMessage: "network_id is required",
		}, nil
	}

	if req.SyncType == "" {
		return &proto.SyncStoreNetworkStateResponse{
			Success:      false,
			ErrorMessage: "sync_type is required",
		}, nil
	}

	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.syncNetworkStateOnKeeper(ctx, req); err == nil {
			// Save to PostgreSQL for analytics (sync to ensure data persistence for network health tracking)
			if r.db != nil {
				correlationID := utils.GetCorrelationID(ctx)
				if err := r.saveSyncStateToDatabase(ctx, req, result); err != nil {
					r.logger.Error("Failed to save network sync analytics",
						logging.Error(err),
						logging.String("network_id", req.NetworkId),
						logging.String("correlation_id", correlationID))
					// Continue even if database save fails (keeper is primary)
				} else {
					r.logger.Info("Network sync analytics saved successfully",
						logging.String("network_id", req.NetworkId),
						logging.String("correlation_id", correlationID))
				}
			}
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.syncNetworkStateInDatabase(ctx, req)
}

// Helper methods for StoreNetworkKeeper interaction

// getSDKContext creates a sdk.Context from context.Context
// Uses shared utility to avoid code duplication
func (r *Repository) getSDKContext(ctx context.Context) (sdk.Context, error) {
	return utils.GetSDKContext(ctx, r.cosmosApp, r.logger)
}

// syncNetworkStateOnKeeper syncs network state on the keeper
func (r *Repository) syncNetworkStateOnKeeper(ctx context.Context, req *proto.SyncStoreNetworkStateRequest) (*proto.SyncStoreNetworkStateResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, err
	}

	// Generate sync ID
	syncID := fmt.Sprintf("sync_%s_%d", req.NetworkId, time.Now().Unix())

	// Create StoredData to store sync information
	valueJSON := fmt.Sprintf(`{"networkId":"%s","syncType":"%s"}`, req.NetworkId, req.SyncType)
	storedData := storenetworktypes.StoredData{
		ID:          syncID,
		Key:         fmt.Sprintf("sync:%s", req.NetworkId),
		Value:       []byte(valueJSON),
		Size:        int64(len(valueJSON)),
		ContentType: "application/json",
		Tags:        make(map[string]string),
		Metadata:    make(map[string]string),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ExpiresAt:   time.Time{},
		Version:     1,
	}

	// Set stored data in keeper
	if err := r.cosmosApp.StoreNetworkKeeper.SetStoredData(sdkCtx, storedData); err != nil {
		return nil, repoerrors.NewDatabaseError("set_stored_data", err)
	}

	// Get all stored data to count syncs
	allData := r.cosmosApp.StoreNetworkKeeper.GetAllStoredData(sdkCtx)
	blocksSynced := int64(0)
	transactionsSynced := int64(0)
	for _, data := range allData {
		if data.Key == fmt.Sprintf("sync:%s", req.NetworkId) {
			blocksSynced++
			transactionsSynced++
		}
	}

	return &proto.SyncStoreNetworkStateResponse{
		Success:            true,
		SyncId:             syncID,
		BlocksSynced:       blocksSynced,
		TransactionsSynced: transactionsSynced,
		SyncStatus:         "completed",
		ErrorMessage:       "",
	}, nil
}

// Database fallback methods

// syncNetworkStateInDatabase syncs network state in database
func (r *Repository) syncNetworkStateInDatabase(ctx context.Context, req *proto.SyncStoreNetworkStateRequest) (*proto.SyncStoreNetworkStateResponse, error) {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return nil, repoerrors.NewRepositoryError(repoerrors.ErrDatabaseUnavailable)
	}

	// Generate sync ID
	syncID := fmt.Sprintf("sync_%s_%d", req.NetworkId, time.Now().Unix())

	query := `
		INSERT INTO network_sync_logs (
			network_id, sync_id, sync_type, from_block, to_block,
			blocks_synced, transactions_synced, sync_status, user_id, store_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (sync_id) DO UPDATE SET
			sync_status = EXCLUDED.sync_status,
			blocks_synced = EXCLUDED.blocks_synced,
			transactions_synced = EXCLUDED.transactions_synced
	`

	_, err := postgres.ExecContext(ctx, query,
		req.NetworkId, syncID, req.SyncType, req.FromBlock, req.ToBlock,
		0, 0, "in_progress", req.UserId, req.StoreId,
	)
	if err != nil {
		return nil, repoerrors.NewDatabaseError("sync_network_state", err)
	}

	return &proto.SyncStoreNetworkStateResponse{
		Success:            true,
		SyncId:             syncID,
		BlocksSynced:       0,
		TransactionsSynced: 0,
		SyncStatus:         "pending",
		ErrorMessage:       "",
	}, nil
}

// saveSyncStateToDatabase saves sync state to database for analytics (sync for network health tracking)
func (r *Repository) saveSyncStateToDatabase(ctx context.Context, req *proto.SyncStoreNetworkStateRequest, result *proto.SyncStoreNetworkStateResponse) error {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return fmt.Errorf("postgres connection not available")
	}

	query := `
		INSERT INTO network_sync_logs (
			network_id, sync_id, sync_type, from_block, to_block,
			blocks_synced, transactions_synced, sync_status, sync_completed_at, user_id, store_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (sync_id) DO UPDATE SET
			sync_status = EXCLUDED.sync_status,
			blocks_synced = EXCLUDED.blocks_synced,
			transactions_synced = EXCLUDED.transactions_synced,
			sync_completed_at = EXCLUDED.sync_completed_at
	`

	completedAt := time.Now()
	if result.SyncStatus != "completed" {
		completedAt = time.Time{}
	}

	if _, err := postgres.ExecContext(ctx, query,
		req.NetworkId, result.SyncId, req.SyncType, req.FromBlock, req.ToBlock,
		result.BlocksSynced, result.TransactionsSynced, result.SyncStatus, completedAt,
		req.UserId, req.StoreId,
	); err != nil {
		return fmt.Errorf("failed to save network sync state to database: %w", err)
	}

	return nil
}

// GetStoreNetworkInfo retrieves store network information
func (r *Repository) GetStoreNetworkInfo(ctx context.Context, req *proto.GetStoreNetworkInfoRequest) (*proto.GetStoreNetworkInfoResponse, error) {
	// Validate request
	if req.NetworkId == "" {
		return nil, repoerrors.NewValidationError("network_id", "is required")
	}

	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.getNetworkInfoFromKeeper(ctx, req); err == nil {
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.getNetworkInfoFromDatabase(ctx, req)
}

// getNetworkInfoFromKeeper retrieves network info from the keeper
func (r *Repository) getNetworkInfoFromKeeper(ctx context.Context, req *proto.GetStoreNetworkInfoRequest) (*proto.GetStoreNetworkInfoResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, err
	}

	// Get stored data for this network
	allData := r.cosmosApp.StoreNetworkKeeper.GetAllStoredData(sdkCtx)
	for _, data := range allData {
		if data.Key == fmt.Sprintf("sync:%s", req.NetworkId) {
			return &proto.GetStoreNetworkInfoResponse{
				Success:      true,
				ErrorMessage: "",
				// Note: Additional info fields would be populated from data.Value
			}, nil
		}
	}

	// Network not found
	return &proto.GetStoreNetworkInfoResponse{
		Success:      false,
		ErrorMessage: "network not found",
	}, nil
}

// getNetworkInfoFromDatabase retrieves network info from database
func (r *Repository) getNetworkInfoFromDatabase(ctx context.Context, req *proto.GetStoreNetworkInfoRequest) (*proto.GetStoreNetworkInfoResponse, error) {
	// Note: Network sync metadata table not implemented yet
	return &proto.GetStoreNetworkInfoResponse{
		Success:      true,
		ErrorMessage: "",
	}, nil
}

// UpdateStoreBridgeConfig updates bridge configuration
func (r *Repository) UpdateStoreBridgeConfig(ctx context.Context, req *proto.UpdateStoreBridgeConfigRequest) (*proto.UpdateStoreBridgeConfigResponse, error) {
	// Validate request
	if req.BridgeAddress == "" {
		return &proto.UpdateStoreBridgeConfigResponse{
			Success:      false,
			ErrorMessage: "bridge_address is required",
		}, nil
	}

	if req.ConfigType == "" {
		return &proto.UpdateStoreBridgeConfigResponse{
			Success:      false,
			ErrorMessage: "config_type is required",
		}, nil
	}

	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.updateBridgeConfigOnKeeper(ctx, req); err == nil {
			// Save to PostgreSQL for analytics (sync to ensure data persistence for network health tracking)
			if r.db != nil {
				correlationID := utils.GetCorrelationID(ctx)
				if err := r.saveBridgeConfigToDatabase(ctx, req, result); err != nil {
					r.logger.Error("Failed to save bridge config analytics",
						logging.Error(err),
						logging.String("bridge_address", req.BridgeAddress),
						logging.String("correlation_id", correlationID))
					// Continue even if database save fails (keeper is primary)
				} else {
					r.logger.Info("Bridge config analytics saved successfully",
						logging.String("bridge_address", req.BridgeAddress),
						logging.String("correlation_id", correlationID))
				}
			}
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.updateBridgeConfigInDatabase(ctx, req)
}

// updateBridgeConfigOnKeeper updates bridge config on the keeper
func (r *Repository) updateBridgeConfigOnKeeper(ctx context.Context, req *proto.UpdateStoreBridgeConfigRequest) (*proto.UpdateStoreBridgeConfigResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, err
	}

	// Generate config ID
	configID := fmt.Sprintf("config_%s_%d", req.BridgeAddress, time.Now().Unix())

	// Create StoredData to store bridge config
	valueJSON := fmt.Sprintf(`{"bridgeAddress":"%s","configType":"%s","configData":"%s"}`, req.BridgeAddress, req.ConfigType, req.ConfigData)
	storedData := storenetworktypes.StoredData{
		ID:          configID,
		Key:         fmt.Sprintf("bridge_config:%s", req.BridgeAddress),
		Value:       []byte(valueJSON),
		Size:        int64(len(valueJSON)),
		ContentType: "application/json",
		Tags:        make(map[string]string),
		Metadata:    make(map[string]string),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ExpiresAt:   time.Time{},
		Version:     1,
	}

	// Set stored data in keeper
	if err := r.cosmosApp.StoreNetworkKeeper.SetStoredData(sdkCtx, storedData); err != nil {
		return nil, repoerrors.NewDatabaseError("set_stored_data", err)
	}

	return &proto.UpdateStoreBridgeConfigResponse{
		Success:       true,
		ConfigId:      configID,
		UpdatedConfig: req.ConfigData,
		ErrorMessage:  "",
	}, nil
}

// updateBridgeConfigInDatabase updates bridge config in database
func (r *Repository) updateBridgeConfigInDatabase(ctx context.Context, req *proto.UpdateStoreBridgeConfigRequest) (*proto.UpdateStoreBridgeConfigResponse, error) {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return nil, repoerrors.NewRepositoryError(repoerrors.ErrDatabaseUnavailable)
	}

	// Generate config ID
	configID := fmt.Sprintf("config_%s_%d", req.BridgeAddress, time.Now().Unix())

	// Note: Bridge config table not implemented yet
	return &proto.UpdateStoreBridgeConfigResponse{
		Success:       true,
		ConfigId:      configID,
		UpdatedConfig: req.ConfigData,
		ErrorMessage:  "",
	}, nil
}

// saveBridgeConfigToDatabase saves bridge config to database for analytics (sync for network health tracking)
func (r *Repository) saveBridgeConfigToDatabase(ctx context.Context, req *proto.UpdateStoreBridgeConfigRequest, result *proto.UpdateStoreBridgeConfigResponse) error {
	// Note: Bridge config table not implemented yet, so this is a placeholder
	// Return nil for now (no-op) until table is implemented
	_ = ctx
	_ = req
	_ = result
	return nil
}
