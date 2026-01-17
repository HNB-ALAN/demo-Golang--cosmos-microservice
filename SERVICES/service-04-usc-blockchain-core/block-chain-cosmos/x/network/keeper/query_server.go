package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/protobuf/types/known/timestamppb"

	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/network/v1/usc/network/v1"
)

// QueryServer defines the gRPC querier service for the network module
type QueryServer interface {
	QueryNetwork(context.Context, *blockchainproto.QueryNetworkRequest) (*blockchainproto.QueryNetworkResponse, error)
	QueryNetworks(context.Context, *blockchainproto.QueryNetworksRequest) (*blockchainproto.QueryNetworksResponse, error)
	QueryNetworkStats(context.Context, *blockchainproto.QueryNetworkStatsRequest) (*blockchainproto.QueryNetworkStatsResponse, error)
}

// queryServer implements the QueryServer interface
type queryServer struct {
	Keeper
}

// NewQueryServerImpl creates a new Network query server
func NewQueryServerImpl(keeper Keeper) QueryServer {
	return &queryServer{Keeper: keeper}
}

// QueryNetwork handles network queries by ID
func (k queryServer) QueryNetwork(ctx context.Context, req *blockchainproto.QueryNetworkRequest) (*blockchainproto.QueryNetworkResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	network, err := k.Keeper.GetNetwork(sdkCtx, req.NetworkId)
	if err != nil {
		return nil, fmt.Errorf("network not found: %w", err)
	}

	// Convert internal network to blockchain-proto format
	blockchainNetwork := &blockchainproto.Network{
		Id:          network.ID,
		Name:        network.Name,
		Description: network.Description,
		Status:      blockchainproto.NodeStatus_NODE_STATUS_ACTIVE,
		CreatedAt:   timestamppb.New(network.CreatedAt),
		UpdatedAt:   timestamppb.New(network.UpdatedAt),
	}

	return &blockchainproto.QueryNetworkResponse{Network: blockchainNetwork}, nil
}

// QueryNetworks handles queries for all networks
func (k queryServer) QueryNetworks(ctx context.Context, req *blockchainproto.QueryNetworksRequest) (*blockchainproto.QueryNetworksResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	networks := k.Keeper.GetAllNetworks(sdkCtx)

	// Convert internal networks to blockchain-proto format
	var blockchainNetworks []*blockchainproto.Network
	for _, network := range networks {
		blockchainNetwork := &blockchainproto.Network{
			Id:          network.ID,
			Name:        network.Name,
			Description: network.Description,
			Status:      blockchainproto.NodeStatus_NODE_STATUS_ACTIVE,
			CreatedAt:   timestamppb.New(network.CreatedAt),
			UpdatedAt:   timestamppb.New(network.UpdatedAt),
		}
		blockchainNetworks = append(blockchainNetworks, blockchainNetwork)
	}

	return &blockchainproto.QueryNetworksResponse{Networks: blockchainNetworks}, nil
}

// QueryNetworkStats handles network statistics queries
func (k queryServer) QueryNetworkStats(ctx context.Context, req *blockchainproto.QueryNetworkStatsRequest) (*blockchainproto.QueryNetworkStatsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get network
	network, err := k.Keeper.GetNetwork(sdkCtx, req.NetworkId)
	if err != nil {
		return nil, fmt.Errorf("network not found: %w", err)
	}

	// Get all nodes for this network
	nodes := k.Keeper.GetAllNodes(sdkCtx)
	networkNodes := 0
	for _, node := range nodes {
		if node.NetworkID == req.NetworkId {
			networkNodes++
		}
	}

	// Calculate statistics
	stats := &blockchainproto.NetworkStats{
		TotalNetworks:          1,
		ActiveNetworks:         1,
		TotalNodes:             int64(networkNodes),
		ActiveNodes:            int64(networkNodes), // Assume all nodes are active
		AverageNodesPerNetwork: int64(networkNodes),
		TotalSyncOperations:    0,
		AverageSyncTimeSeconds: 0,
		MostActiveNetwork:      network.ID,
		LastNetworkActivity:    timestamppb.New(network.UpdatedAt),
	}

	return &blockchainproto.QueryNetworkStatsResponse{Stats: stats}, nil
}
