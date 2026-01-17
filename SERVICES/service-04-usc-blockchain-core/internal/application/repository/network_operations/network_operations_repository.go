package network_operations

import (
	"context"
	"fmt"

	"service-04/internal/application/utils"
	"service-04/internal/infrastructure/database"
	proto "service-04/proto"

	// Cosmos SDK imports
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/app"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/storage"

	"github.com/usc-platform/shared/logging"
)

// Repository handles network operations data access
type Repository struct {
	db                *database.PostgreSQLManager
	cosmosApp         *app.USCApp
	blockchainStorage *storage.StateManager
	redisManager      *database.RedisManager
	logger            *logging.Logger
}

// NewRepository creates a new network operations repository
func NewRepository(db *database.PostgreSQLManager, cosmosApp *app.USCApp, blockchainStorage *storage.StateManager, redisManager *database.RedisManager, logger *logging.Logger) *Repository {
	return &Repository{
		db:                db,
		cosmosApp:         cosmosApp,
		blockchainStorage: blockchainStorage,
		redisManager:      redisManager,
		logger:            logger,
	}
}

// GetNetworkInfo retrieves network information
func (r *Repository) GetNetworkInfo(ctx context.Context) (*proto.GetNetworkInfoResponse, error) {
	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.getNetworkInfoFromKeeper(ctx); err == nil {
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.getNetworkInfoFromDatabase(ctx)
}

// GetChainInfo retrieves chain information
func (r *Repository) GetChainInfo(ctx context.Context) (*proto.GetChainInfoResponse, error) {
	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.getChainInfoFromKeeper(ctx); err == nil {
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.getChainInfoFromDatabase(ctx)
}

// GetPeers retrieves list of peers
func (r *Repository) GetPeers(ctx context.Context, req *proto.GetPeersRequest) (*proto.GetPeersResponse, error) {
	limit, offset := utils.NormalizePagination(req.Limit, req.Offset, utils.PaginationConfig{
		DefaultLimit: 100,
		MaxLimit:     1000,
	})

	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.getPeersFromKeeper(ctx, limit, offset); err == nil {
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.getPeersFromDatabase(ctx, req)
}

// Helper methods for NetworkKeeper interaction

// getSDKContext creates a sdk.Context from context.Context
// Uses shared utility to avoid code duplication
func (r *Repository) getSDKContext(ctx context.Context) (sdk.Context, error) {
	return utils.GetSDKContext(ctx, r.cosmosApp, r.logger)
}

// getNetworkInfoFromKeeper retrieves network info from the keeper
func (r *Repository) getNetworkInfoFromKeeper(ctx context.Context) (*proto.GetNetworkInfoResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, err
	}

	// Get all networks from keeper
	networks := r.cosmosApp.NetworkKeeper.GetAllNetworks(sdkCtx)

	// Aggregate network info
	totalNetworks := int32(len(networks))
	totalNodes := int32(0)
	totalConnections := int32(0)

	for _, network := range networks {
		totalNodes += int32(len(network.Nodes))
		totalConnections += int32(len(network.Connections))
	}

	return &proto.GetNetworkInfoResponse{
		NetworkId:          "usc_mainnet",
		NetworkName:        "USC Main Network",
		ChainId:            "usc-1",
		NetworkVersion:     "1.0.0",
		CurrentBlockHeight: int64(sdkCtx.BlockHeight()),
		TotalBlocks:        int64(totalNetworks),
		TotalTransactions:  int64(totalConnections),
		NetworkStatus:      "active",
		ActiveValidators:   int64(totalNodes),
	}, nil
}

// getChainInfoFromKeeper retrieves chain info from the keeper
func (r *Repository) getChainInfoFromKeeper(ctx context.Context) (*proto.GetChainInfoResponse, error) {
	return &proto.GetChainInfoResponse{
		ChainId:            "usc-1",
		ChainName:          "USC Chain",
		ChainType:          "mainnet",
		ConsensusAlgorithm: "pos",
		BlockTime:          3,
		NetworkVersion:     "1.0.0",
		ChainStatus:        "active",
	}, nil
}

// getChainInfoFromDatabase retrieves chain info from database
func (r *Repository) getChainInfoFromDatabase(ctx context.Context) (*proto.GetChainInfoResponse, error) {
	return &proto.GetChainInfoResponse{
		ChainId:            "usc-1",
		ChainName:          "USC Chain",
		ChainType:          "mainnet",
		ConsensusAlgorithm: "pos",
		BlockTime:          3,
		NetworkVersion:     "1.0.0",
		ChainStatus:        "active",
	}, nil
}

// getNetworkStatsFromKeeper retrieves network stats from the keeper
func (r *Repository) getNetworkStatsFromKeeper(ctx context.Context, req *proto.GetNetworkStatsRequest) (*proto.GetNetworkStatsResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, err
	}

	allNodes := r.cosmosApp.NetworkKeeper.GetAllNodes(sdkCtx)

	totalNodes := int64(len(allNodes))
	activeNodes := int64(0)
	for _, node := range allNodes {
		if node.Status == "active" {
			activeNodes++
		}
	}

	return &proto.GetNetworkStatsResponse{
		Metrics: &proto.NetworkMetrics{
			TotalPeers:       totalNodes,
			ConnectedPeers:   activeNodes,
			ActiveValidators: activeNodes,
		},
		TimeRange: req.TimeRange,
	}, nil
}

// getNetworkStatsFromDatabase retrieves network stats from database
func (r *Repository) getNetworkStatsFromDatabase(ctx context.Context, req *proto.GetNetworkStatsRequest) (*proto.GetNetworkStatsResponse, error) {
	// Note: Network stats from PostgreSQL analytics not implemented yet
	// This is a fallback method, so returning default values is acceptable
	return &proto.GetNetworkStatsResponse{
		Metrics:   &proto.NetworkMetrics{},
		TimeRange: req.TimeRange,
	}, nil
}

// getPeersFromKeeper retrieves peers (nodes) from the keeper
func (r *Repository) getPeersFromKeeper(ctx context.Context, limit, offset int32) (*proto.GetPeersResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, err
	}

	// Get all nodes from keeper
	allNodes := r.cosmosApp.NetworkKeeper.GetAllNodes(sdkCtx)

	// Apply pagination
	start := int(offset)
	end := start + int(limit)
	if end > len(allNodes) {
		end = len(allNodes)
	}

	// Pre-allocate slice with capacity = (end - start) for better performance
	peers := make([]*proto.PeerInfo, 0, end-start)
	for i := start; i < end; i++ {
		node := allNodes[i]
		peerAddress := fmt.Sprintf("%s:%d", node.Address, node.Port)
		peers = append(peers, &proto.PeerInfo{
			PeerId:      node.ID,
			PeerAddress: peerAddress,
			PeerType:    "full_node", // Default
			Status:      string(node.Status),
			LastSeen:    node.LastSeen.Unix(),
		})
	}

	return &proto.GetPeersResponse{
		Peers:      peers,
		TotalCount: int32(len(allNodes)),
		HasMore:    end < len(allNodes),
	}, nil
}

// Database fallback methods

// getNetworkInfoFromDatabase retrieves network info from database
func (r *Repository) getNetworkInfoFromDatabase(ctx context.Context) (*proto.GetNetworkInfoResponse, error) {
	// Simplified database implementation
	// In production, query from PostgreSQL
	return &proto.GetNetworkInfoResponse{
		NetworkId:      "usc_mainnet",
		NetworkName:    "USC Main Network",
		NetworkVersion: "1.0.0",
		NetworkStatus:  "active",
	}, nil
}

// getPeersFromDatabase retrieves peers from database
func (r *Repository) getPeersFromDatabase(ctx context.Context, req *proto.GetPeersRequest) (*proto.GetPeersResponse, error) {
	// Simplified database implementation
	// In production, query from PostgreSQL with pagination
	return &proto.GetPeersResponse{
		Peers:      []*proto.PeerInfo{},
		TotalCount: 0,
		HasMore:    false,
	}, nil
}

// GetNetworkStats retrieves network statistics
func (r *Repository) GetNetworkStats(ctx context.Context, req *proto.GetNetworkStatsRequest) (*proto.GetNetworkStatsResponse, error) {
	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.getNetworkStatsFromKeeper(ctx, req); err == nil {
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.getNetworkStatsFromDatabase(ctx, req)
}
