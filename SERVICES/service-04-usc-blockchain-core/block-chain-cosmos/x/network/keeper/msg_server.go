package keeper

import (
	"context"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/network/v1/usc/network/v1"
	blocktypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/block/types"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/network/types"
)

// MsgServer defines the gRPC message server for the network module
type MsgServer interface {
	JoinNetwork(context.Context, *blockchainproto.MsgJoinNetwork) (*blockchainproto.MsgJoinNetworkResponse, error)
	LeaveNetwork(context.Context, *blockchainproto.MsgLeaveNetwork) (*blockchainproto.MsgLeaveNetworkResponse, error)
	UpdateNetwork(context.Context, *blockchainproto.MsgUpdateNetwork) (*blockchainproto.MsgUpdateNetworkResponse, error)
	SyncNetwork(context.Context, *blockchainproto.MsgSyncNetwork) (*blockchainproto.MsgSyncNetworkResponse, error)
}

// msgServer implements the MsgServer interface
type msgServer struct {
	Keeper
}

// NewMsgServerImpl creates a new Network message server
func NewMsgServerImpl(keeper Keeper) MsgServer {
	return &msgServer{Keeper: keeper}
}

// JoinNetwork handles network joining messages
func (k msgServer) JoinNetwork(ctx context.Context, msg *blockchainproto.MsgJoinNetwork) (*blockchainproto.MsgJoinNetworkResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Create node info
	nodeInfo := types.Node{
		ID:        msg.NodeId,
		Address:   msg.NodeAddress,
		Name:      msg.NodeId,
		NetworkID: msg.NetworkId,
		Status:    types.NodeStatusOnline,
		Metadata:  map[string]string{"owner": msg.Joiner},
	}

	// Add node to network
	if err := k.SetNode(sdkCtx, nodeInfo); err != nil {
		return nil, fmt.Errorf("failed to join network: %w", err)
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeNodeJoined,
			sdk.NewAttribute(types.AttributeKeyNodeID, msg.NodeId),
			sdk.NewAttribute(types.AttributeKeyNetworkID, msg.NetworkId),
			sdk.NewAttribute(types.AttributeKeyModule, types.ModuleName),
		),
	)

	return &blockchainproto.MsgJoinNetworkResponse{
		Success:             true,
		ConnectionId:        msg.NodeId,
		NetworkMembershipId: msg.NetworkId,
	}, nil
}

// UpdateNetwork handles network update messages
func (k msgServer) UpdateNetwork(ctx context.Context, msg *blockchainproto.MsgUpdateNetwork) (*blockchainproto.MsgUpdateNetworkResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get existing network
	network, err := k.GetNetwork(sdkCtx, msg.NetworkId)
	if err != nil {
		return nil, fmt.Errorf("network not found: %w", err)
	}

	// Update network fields based on consensus params and topology
	network.UpdatedAt = sdkCtx.BlockTime()

	// Validate updated network
	if err := network.Validate(); err != nil {
		return nil, fmt.Errorf("invalid network update: %w", err)
	}

	// Set the updated network
	if err := k.SetNetwork(sdkCtx, network); err != nil {
		return nil, fmt.Errorf("failed to update network: %w", err)
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeNetworkUpdated,
			sdk.NewAttribute(types.AttributeKeyNetworkID, msg.NetworkId),
			sdk.NewAttribute(types.AttributeKeyModule, types.ModuleName),
		),
	)

	return &blockchainproto.MsgUpdateNetworkResponse{
		Success:    true,
		UpdateHash: msg.NetworkId,
	}, nil
}

// LeaveNetwork handles network leaving messages
func (k msgServer) LeaveNetwork(ctx context.Context, msg *blockchainproto.MsgLeaveNetwork) (*blockchainproto.MsgLeaveNetworkResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Remove node from network - just mark as offline
	node, err := k.GetNode(sdkCtx, msg.NodeId)
	if err != nil {
		return nil, fmt.Errorf("node not found: %w", err)
	}
	node.Status = types.NodeStatusOffline
	if err := k.SetNode(sdkCtx, node); err != nil {
		return nil, fmt.Errorf("failed to leave network: %w", err)
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeNodeLeft,
			sdk.NewAttribute(types.AttributeKeyNodeID, msg.NodeId),
			sdk.NewAttribute(types.AttributeKeyNetworkID, msg.NetworkId),
			sdk.NewAttribute(types.AttributeKeyModule, types.ModuleName),
		),
	)

	return &blockchainproto.MsgLeaveNetworkResponse{
		Success:   true,
		LeaveHash: msg.NodeId,
	}, nil
}

// SyncNetwork handles network synchronization messages
func (k msgServer) SyncNetwork(ctx context.Context, msg *blockchainproto.MsgSyncNetwork) (*blockchainproto.MsgSyncNetworkResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Extract sync parameters from sync_params map
	networkID := msg.SyncParams["network_id"]
	if networkID == "" {
		networkID = "default" // Use default network if not specified
	}

	// Get or create node
	node, err := k.GetNode(sdkCtx, msg.NodeId)
	if err != nil {
		// Node doesn't exist, create a new one for sync
		node = types.Node{
			ID:        msg.NodeId,
			NetworkID: networkID,
			Status:    types.NodeStatusConnecting,
			LastSeen:  sdkCtx.BlockTime(),
			CreatedAt: sdkCtx.BlockTime(),
			UpdatedAt: sdkCtx.BlockTime(),
			Metadata:  make(map[string]string),
		}
		if err := k.SetNode(sdkCtx, node); err != nil {
			return nil, fmt.Errorf("failed to create node for sync: %w", err)
		}
	}

	// Get current block height for sync progress calculation
	currentHeight := sdkCtx.BlockHeight()
	startHeight := currentHeight
	if startHeightStr := msg.SyncParams["start_height"]; startHeightStr != "" {
		if parsed, err := strconv.ParseInt(startHeightStr, 10, 64); err == nil && parsed > 0 {
			startHeight = parsed
		}
	}
	endHeight := currentHeight
	if endHeightStr := msg.SyncParams["end_height"]; endHeightStr != "" {
		if parsed, err := strconv.ParseInt(endHeightStr, 10, 64); err == nil && parsed > 0 {
			endHeight = parsed
		}
	}

	// Calculate sync progress (0-100)
	var progress int64
	if endHeight > startHeight {
		progress = ((currentHeight - startHeight) * 100) / (endHeight - startHeight)
		if progress < 0 {
			progress = 0
		}
		if progress > 100 {
			progress = 100
		}
	} else {
		progress = 100 // Already synced
	}

	// Determine sync status based on sync type and progress
	syncStatus := types.NetworkSync{
		ID:            fmt.Sprintf("sync_%s_%d", msg.NodeId, sdkCtx.BlockTime().Unix()),
		NetworkID:     networkID,
		NodeID:        msg.NodeId,
		StartHeight:   startHeight,
		EndHeight:     endHeight,
		CurrentHeight: currentHeight,
		StartedAt:     sdkCtx.BlockTime(),
		Metadata:      make(map[string]string),
	}

	// Update sync status based on progress and sync type
	if progress == 100 {
		syncStatus.Status = "synced"
		completedAt := sdkCtx.BlockTime()
		syncStatus.CompletedAt = &completedAt
		node.Status = types.NodeStatusOnline
	} else if progress > 0 {
		syncStatus.Status = "syncing"
		node.Status = types.NodeStatusConnecting
	} else {
		syncStatus.Status = "pending"
		node.Status = types.NodeStatusConnecting
	}
	syncStatus.Progress = progress

	// Store sync type in metadata
	syncStatus.Metadata["sync_type"] = msg.SyncType.String()
	syncStatus.Metadata["syncer"] = msg.Syncer

	// Update node last seen and status
	node.LastSeen = sdkCtx.BlockTime()
	node.UpdatedAt = sdkCtx.BlockTime()
	if err := k.SetNode(sdkCtx, node); err != nil {
		return nil, fmt.Errorf("failed to update node: %w", err)
	}

	// Store sync record
	if err := k.SetSync(sdkCtx, syncStatus); err != nil {
		return nil, fmt.Errorf("failed to store sync record: %w", err)
	}

	// Update network topology by checking connections
	connections := k.GetAllConnections(sdkCtx)
	activeConnections := 0
	for _, conn := range connections {
		if conn.NetworkID == networkID && conn.Status == types.ConnectionStatusActive {
			activeConnections++
		}
	}

	// Calculate transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Syncer, msg.NodeId, "", "sync_network", networkID, msg.Memo)

	// Emit event with detailed sync information
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeNetworkSync,
			sdk.NewAttribute(types.AttributeKeyNodeID, msg.NodeId),
			sdk.NewAttribute(types.AttributeKeyNetworkID, networkID),
			sdk.NewAttribute("sync_id", syncStatus.ID),
			sdk.NewAttribute("sync_status", syncStatus.Status),
			sdk.NewAttribute("sync_type", msg.SyncType.String()),
			sdk.NewAttribute("progress", fmt.Sprintf("%d", progress)),
			sdk.NewAttribute("current_height", fmt.Sprintf("%d", currentHeight)),
			sdk.NewAttribute("active_connections", fmt.Sprintf("%d", activeConnections)),
			sdk.NewAttribute(types.AttributeKeyModule, types.ModuleName),
		),
	)

	// Convert sync status to proto enum
	var protoSyncStatus blockchainproto.SyncStatus
	switch syncStatus.Status {
	case "synced":
		protoSyncStatus = blockchainproto.SyncStatus_SYNC_STATUS_COMPLETED
	case "syncing":
		protoSyncStatus = blockchainproto.SyncStatus_SYNC_STATUS_IN_PROGRESS
	case "pending":
		protoSyncStatus = blockchainproto.SyncStatus_SYNC_STATUS_PENDING
	default:
		protoSyncStatus = blockchainproto.SyncStatus_SYNC_STATUS_FAILED
	}

	return &blockchainproto.MsgSyncNetworkResponse{
		Success:         true,
		SyncId:          syncStatus.ID,
		SyncStatus:      protoSyncStatus,
		TransactionHash: txHash,
	}, nil
}
