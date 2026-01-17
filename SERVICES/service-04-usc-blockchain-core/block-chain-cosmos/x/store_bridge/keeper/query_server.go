package keeper

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/protobuf/types/known/timestamppb"

	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/store_bridge/v1/usc/store_bridge/v1"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/store_bridge/types"
)

// QueryServer defines the gRPC querier service for the store_bridge module
type QueryServer interface {
	QueryBridge(context.Context, *blockchainproto.QueryBridgeRequest) (*blockchainproto.QueryBridgeResponse, error)
	QueryBridges(context.Context, *blockchainproto.QueryBridgesRequest) (*blockchainproto.QueryBridgesResponse, error)
	QueryBridgeOperations(context.Context, *blockchainproto.QueryBridgeOperationsRequest) (*blockchainproto.QueryBridgeOperationsResponse, error)
	QueryBridgeStats(context.Context, *blockchainproto.QueryBridgeStatsRequest) (*blockchainproto.QueryBridgeStatsResponse, error)
}

// queryServer implements QueryServer
type queryServer struct {
	Keeper
}

// NewQueryServerImpl creates a new store_bridge query server
func NewQueryServerImpl(keeper Keeper) QueryServer {
	return &queryServer{Keeper: keeper}
}

// QueryBridge handles bridge queries by ID
func (k queryServer) QueryBridge(ctx context.Context, req *blockchainproto.QueryBridgeRequest) (*blockchainproto.QueryBridgeResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	bridge, err := k.Keeper.GetBridge(sdkCtx, req.BridgeId)
	if err != nil {
		return nil, fmt.Errorf("bridge not found: %w", err)
	}

	// Convert internal bridge to blockchainproto.Bridge
	blockchainBridge := convertBridgeToProto(bridge)

	return &blockchainproto.QueryBridgeResponse{
		Bridge: blockchainBridge,
	}, nil
}

// QueryBridges handles queries for multiple bridges
func (k queryServer) QueryBridges(ctx context.Context, req *blockchainproto.QueryBridgesRequest) (*blockchainproto.QueryBridgesResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	bridges := k.Keeper.GetAllBridges(sdkCtx)

	// Apply filters
	filteredBridges := []types.Bridge{}
	for _, bridge := range bridges {
		// Filter by bridge_type
		if req.BridgeType != blockchainproto.BridgeType_BRIDGE_TYPE_UNSPECIFIED {
			protoType := convertTypeToProto(bridge.Type)
			if protoType != req.BridgeType {
				continue
			}
		}
		// Filter by status
		if req.Status != blockchainproto.BridgeStatus_BRIDGE_STATUS_UNSPECIFIED {
			protoStatus := convertStatusToProto(bridge.Status)
			if protoStatus != req.Status {
				continue
			}
		}
		filteredBridges = append(filteredBridges, bridge)
	}

	// Convert to proto
	protoBridges := make([]*blockchainproto.Bridge, len(filteredBridges))
	for i, bridge := range filteredBridges {
		protoBridges[i] = convertBridgeToProto(bridge)
	}

	return &blockchainproto.QueryBridgesResponse{
		Bridges: protoBridges,
	}, nil
}

// QueryBridgeOperations handles queries for bridge operations
func (k queryServer) QueryBridgeOperations(ctx context.Context, req *blockchainproto.QueryBridgeOperationsRequest) (*blockchainproto.QueryBridgeOperationsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	transfers := k.Keeper.GetAllTransfers(sdkCtx)

	// Filter by bridge_id if specified
	filteredTransfers := []types.Transfer{}
	for _, transfer := range transfers {
		if req.BridgeId != "" && transfer.BridgeID != req.BridgeId {
			continue
		}
		// Filter by operation_type (using bridge type as fallback)
		if req.OperationType != blockchainproto.OperationType_OPERATION_TYPE_UNSPECIFIED {
			protoType := convertOperationTypeToProto("transfer") // Default to transfer
			if protoType != req.OperationType {
				continue
			}
		}
		// Filter by status
		if req.Status != blockchainproto.OperationStatus_OPERATION_STATUS_UNSPECIFIED {
			protoStatus := convertOperationStatusToProto(transfer.Status)
			if protoStatus != req.Status {
				continue
			}
		}
		filteredTransfers = append(filteredTransfers, transfer)
	}

	// Convert to proto
	protoOperations := make([]*blockchainproto.BridgeOperation, len(filteredTransfers))
	for i, transfer := range filteredTransfers {
		protoOperations[i] = convertTransferToOperation(transfer)
	}

	return &blockchainproto.QueryBridgeOperationsResponse{
		Operations: protoOperations,
	}, nil
}

// QueryBridgeStats handles queries for bridge statistics
func (k queryServer) QueryBridgeStats(ctx context.Context, req *blockchainproto.QueryBridgeStatsRequest) (*blockchainproto.QueryBridgeStatsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	allBridges := k.Keeper.GetAllBridges(sdkCtx)
	allTransfers := k.Keeper.GetAllTransfers(sdkCtx)

	// Calculate statistics
	totalBridges := int64(len(allBridges))
	activeBridges := int64(0)
	totalOperations := int64(len(allTransfers))
	successfulOperations := int64(0)
	failedOperations := int64(0)
	var mostActiveBridge string

	for _, bridge := range allBridges {
		if bridge.Status == "active" {
			activeBridges++
		}
	}

	for _, transfer := range allTransfers {
		if transfer.Status == "completed" {
			successfulOperations++
		} else if transfer.Status == "failed" {
			failedOperations++
		}
	}

	successRate := 0.0
	if totalOperations > 0 {
		successRate = float64(successfulOperations) / float64(totalOperations) * 100
	}

	stats := &blockchainproto.BridgeStats{
		TotalBridges:                totalBridges,
		ActiveBridges:               activeBridges,
		TotalOperations:             totalOperations,
		SuccessfulOperations:        successfulOperations,
		FailedOperations:            failedOperations,
		SuccessRate:                 successRate,
		TotalVolume:                 &sdk.Coin{Denom: "usc", Amount: sdkmath.NewInt(0)}, // Default volume
		AverageOperationTimeSeconds: 0,                                                  // Calculate if needed
		MostActiveBridge:            mostActiveBridge,
		LastActivity:                timestamppb.New(sdkCtx.BlockTime()),
	}

	return &blockchainproto.QueryBridgeStatsResponse{
		Stats: stats,
	}, nil
}

// Helper functions

func convertBridgeToProto(bridge types.Bridge) *blockchainproto.Bridge {
	// Convert source and target chains
	sourceChain := &blockchainproto.ChainIdentifier{
		ChainId:     bridge.FromChain,
		ChainName:   bridge.FromChain,
		ChainType:   "cosmos", // Default
		RpcEndpoint: "",
		ExplorerUrl: "",
		NativeToken: "usc",
		BlockTime:   6, // Default 6 seconds
	}

	targetChain := &blockchainproto.ChainIdentifier{
		ChainId:     bridge.ToChain,
		ChainName:   bridge.ToChain,
		ChainType:   "cosmos", // Default
		RpcEndpoint: "",
		ExplorerUrl: "",
		NativeToken: "usc",
		BlockTime:   6, // Default 6 seconds
	}

	// Convert bridge type
	bridgeType := convertTypeToProto(bridge.Type)

	// Convert status
	status := convertStatusToProto(bridge.Status)

	// Convert config
	config := &blockchainproto.BridgeConfig{
		ConfirmationBlocks:     1, // Default
		TimeoutBlocks:          100,
		EnableAutoFinalization: true,
		EnableAutoRetry:        true,
		MaxRetryAttempts:       3,
		CustomSettings:         bridge.Config,
	}

	// Convert fees (default)
	fees := &blockchainproto.BridgeFees{
		BaseFee:        &sdk.Coin{Denom: "usc", Amount: sdkmath.NewInt(0)},
		FeePercentage:  0.01, // 1%
		MinimumFee:     &sdk.Coin{Denom: "usc", Amount: sdkmath.NewInt(0)},
		MaximumFee:     &sdk.Coin{Denom: "usc", Amount: sdkmath.NewInt(1000000)},
		DynamicPricing: false,
	}

	// Convert security deposit (default)
	securityDeposit := &blockchainproto.SecurityDeposit{
		RequiredDeposit:       &sdk.Coin{Denom: "usc", Amount: sdkmath.NewInt(0)},
		CurrentDeposit:        &sdk.Coin{Denom: "usc", Amount: sdkmath.NewInt(0)},
		DepositDurationBlocks: 1000,
		AutoRefund:            true,
		DepositConditions:     "Standard bridge conditions",
	}

	return &blockchainproto.Bridge{
		Id:                   bridge.ID,
		Creator:              "system", // Default creator
		SourceChain:          sourceChain,
		TargetChain:          targetChain,
		BridgeType:           bridgeType,
		Status:               status,
		Config:               config,
		Fees:                 fees,
		SecurityDeposit:      securityDeposit,
		CreatedAt:            timestamppb.New(bridge.CreatedAt),
		UpdatedAt:            timestamppb.New(bridge.UpdatedAt),
		TotalOperations:      0, // Calculate if needed
		SuccessfulOperations: 0, // Calculate if needed
		FailedOperations:     0, // Calculate if needed
		TotalVolume:          &sdk.Coin{Denom: "usc", Amount: sdkmath.NewInt(0)},
		Memo:                 "",
	}
}

func convertTransferToOperation(transfer types.Transfer) *blockchainproto.BridgeOperation {
	// Convert source and target chains
	sourceChain := &blockchainproto.ChainIdentifier{
		ChainId:     transfer.FromChain,
		ChainName:   transfer.FromChain,
		ChainType:   "cosmos",
		RpcEndpoint: "",
		ExplorerUrl: "",
		NativeToken: "usc",
		BlockTime:   6,
	}

	targetChain := &blockchainproto.ChainIdentifier{
		ChainId:     transfer.ToChain,
		ChainName:   transfer.ToChain,
		ChainType:   "cosmos",
		RpcEndpoint: "",
		ExplorerUrl: "",
		NativeToken: "usc",
		BlockTime:   6,
	}

	// Parse amount
	amount := sdk.Coin{Denom: transfer.Token, Amount: sdkmath.NewInt(0)}
	if transfer.Amount != "" {
		if parsedAmount, err := sdk.ParseCoinNormalized(transfer.Amount); err == nil {
			amount = parsedAmount
		}
	}

	// Convert status
	status := convertOperationStatusToProto(transfer.Status)

	// Convert operation type (default to transfer)
	operationType := convertOperationTypeToProto("transfer")

	completedAt := (*timestamppb.Timestamp)(nil)
	if !transfer.CompletedAt.IsZero() {
		completedAt = timestamppb.New(transfer.CompletedAt)
	}

	return &blockchainproto.BridgeOperation{
		Id:                transfer.ID,
		BridgeId:          transfer.BridgeID,
		OperationType:     operationType,
		Sender:            transfer.FromAddress,
		Recipient:         transfer.ToAddress,
		Amount:            &amount,
		SourceChain:       sourceChain,
		TargetChain:       targetChain,
		Status:            status,
		CreatedAt:         timestamppb.New(transfer.CreatedAt),
		CompletedAt:       completedAt,
		TransactionHash:   transfer.TxHash,
		FinalizationProof: "",
		Memo:              "",
	}
}

func convertTypeToProto(bridgeType string) blockchainproto.BridgeType {
	switch bridgeType {
	case "token", "asset":
		return blockchainproto.BridgeType_BRIDGE_TYPE_ASSET
	case "data":
		return blockchainproto.BridgeType_BRIDGE_TYPE_DATA
	case "message":
		return blockchainproto.BridgeType_BRIDGE_TYPE_MESSAGE
	case "nft":
		return blockchainproto.BridgeType_BRIDGE_TYPE_NFT
	case "contract":
		return blockchainproto.BridgeType_BRIDGE_TYPE_CONTRACT
	default:
		return blockchainproto.BridgeType_BRIDGE_TYPE_UNSPECIFIED
	}
}

func convertStatusToProto(status string) blockchainproto.BridgeStatus {
	switch status {
	case "active":
		return blockchainproto.BridgeStatus_BRIDGE_STATUS_ACTIVE
	case "inactive":
		return blockchainproto.BridgeStatus_BRIDGE_STATUS_INACTIVE
	case "maintenance":
		return blockchainproto.BridgeStatus_BRIDGE_STATUS_MAINTENANCE
	case "error":
		return blockchainproto.BridgeStatus_BRIDGE_STATUS_ERROR
	default:
		return blockchainproto.BridgeStatus_BRIDGE_STATUS_UNSPECIFIED
	}
}

func convertOperationTypeToProto(operationType string) blockchainproto.OperationType {
	switch operationType {
	case "transfer", "token":
		return blockchainproto.OperationType_OPERATION_TYPE_TRANSFER
	case "message":
		return blockchainproto.OperationType_OPERATION_TYPE_MESSAGE
	case "data":
		return blockchainproto.OperationType_OPERATION_TYPE_DATA
	case "nft":
		return blockchainproto.OperationType_OPERATION_TYPE_NFT
	case "contract":
		return blockchainproto.OperationType_OPERATION_TYPE_CONTRACT
	default:
		return blockchainproto.OperationType_OPERATION_TYPE_UNSPECIFIED
	}
}

func convertOperationStatusToProto(status string) blockchainproto.OperationStatus {
	switch status {
	case "pending":
		return blockchainproto.OperationStatus_OPERATION_STATUS_PENDING
	case "completed":
		return blockchainproto.OperationStatus_OPERATION_STATUS_COMPLETED
	case "failed":
		return blockchainproto.OperationStatus_OPERATION_STATUS_FAILED
	case "cancelled":
		return blockchainproto.OperationStatus_OPERATION_STATUS_CANCELLED
	case "timeout":
		return blockchainproto.OperationStatus_OPERATION_STATUS_TIMEOUT
	default:
		return blockchainproto.OperationStatus_OPERATION_STATUS_UNSPECIFIED
	}
}
