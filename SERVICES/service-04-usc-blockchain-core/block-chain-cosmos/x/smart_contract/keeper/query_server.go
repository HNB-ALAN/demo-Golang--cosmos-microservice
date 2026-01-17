package keeper

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	query "github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/protobuf/types/known/timestamppb"

	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/smart_contract/v1/usc/smart_contract/v1"
)

// QueryServer defines the gRPC querier service for the smart contract module using blockchain-proto types
type QueryServer struct {
	Keeper
}

// NewQueryServer creates a new smart contract query server
func NewQueryServer(keeper Keeper) *QueryServer {
	return &QueryServer{Keeper: keeper}
}

// QueryContract handles individual contract queries
func (k QueryServer) QueryContract(ctx context.Context, req *blockchainproto.QueryContractRequest) (*blockchainproto.QueryContractResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Basic validation
	if req.ContractAddress == "" {
		return nil, fmt.Errorf("contract address cannot be empty")
	}

	// Get contract (simplified - using address as ID)
	contract, err := k.Keeper.GetContract(sdkCtx, req.ContractAddress)
	if err != nil {
		return nil, fmt.Errorf("contract not found: %w", err)
	}

	// Convert string status to blockchain-proto enum
	var contractStatus blockchainproto.ContractStatus
	switch string(contract.Status) {
	case "active":
		contractStatus = blockchainproto.ContractStatus_CONTRACT_STATUS_ACTIVE
	case "destroyed":
		contractStatus = blockchainproto.ContractStatus_CONTRACT_STATUS_DESTROYED
	default:
		contractStatus = blockchainproto.ContractStatus_CONTRACT_STATUS_UNSPECIFIED
	}

	// Convert to blockchain-proto SmartContract type
	blockchainContract := &blockchainproto.SmartContract{
		Address:        contract.Address,
		Deployer:       contract.Owner,
		Name:           contract.Name,
		Version:        contract.Version,
		CodeSize:       0, // Placeholder
		CodeHash:       contract.CodeHash,
		Status:         contractStatus,
		DeployedAt:     timestamppb.New(contract.CreatedAt),
		UpdatedAt:      timestamppb.New(contract.UpdatedAt),
		DestroyedAt:    nil,
		ExecutionCount: 0,
		TotalFeesPaid:  &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)},
		GasLimit:       0,
		GasPrice:       &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)},
		Memo:           "",
	}

	return &blockchainproto.QueryContractResponse{
		Contract: blockchainContract,
	}, nil
}

// QueryContracts handles queries for multiple contracts
func (k QueryServer) QueryContracts(ctx context.Context, req *blockchainproto.QueryContractsRequest) (*blockchainproto.QueryContractsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get all contracts
	contracts := k.Keeper.GetAllContracts(sdkCtx)

	// Convert to blockchain-proto SmartContract types
	var blockchainContracts []*blockchainproto.SmartContract
	for _, contract := range contracts {
		// Convert string status to blockchain-proto enum
		var contractStatus blockchainproto.ContractStatus
		switch string(contract.Status) {
		case "active":
			contractStatus = blockchainproto.ContractStatus_CONTRACT_STATUS_ACTIVE
		case "destroyed":
			contractStatus = blockchainproto.ContractStatus_CONTRACT_STATUS_DESTROYED
		default:
			contractStatus = blockchainproto.ContractStatus_CONTRACT_STATUS_UNSPECIFIED
		}

		blockchainContract := &blockchainproto.SmartContract{
			Address:        contract.Address,
			Deployer:       contract.Owner,
			Name:           contract.Name,
			Version:        contract.Version,
			CodeSize:       0, // Placeholder
			CodeHash:       contract.CodeHash,
			Status:         contractStatus,
			DeployedAt:     timestamppb.New(contract.CreatedAt),
			UpdatedAt:      timestamppb.New(contract.UpdatedAt),
			DestroyedAt:    nil,
			ExecutionCount: 0,
			TotalFeesPaid:  &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)},
			GasLimit:       0,
			GasPrice:       &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)},
			Memo:           "",
		}
		blockchainContracts = append(blockchainContracts, blockchainContract)
	}

	// Apply pagination
	pageRes := &query.PageResponse{
		NextKey: nil,
		Total:   uint64(len(blockchainContracts)),
	}

	return &blockchainproto.QueryContractsResponse{
		Contracts:  blockchainContracts,
		Pagination: pageRes,
	}, nil
}

// QueryContractStats handles contract statistics queries
func (k QueryServer) QueryContractStats(ctx context.Context, req *blockchainproto.QueryContractStatsRequest) (*blockchainproto.QueryContractStatsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get all contracts for statistics
	contracts := k.Keeper.GetAllContracts(sdkCtx)

	// Calculate statistics
	totalContracts := len(contracts)
	activeContracts := 0
	destroyedContracts := 0
	totalExecutions := int64(0)

	for _, contract := range contracts {
		if contract.Status == "active" {
			activeContracts++
		} else if contract.Status == "destroyed" {
			destroyedContracts++
		}
		totalExecutions += 0 // Default execution count
	}

	// Create sample statistics
	stats := &blockchainproto.ContractStats{
		TotalContracts:      int64(totalContracts),
		ActiveContracts:     int64(activeContracts),
		DestroyedContracts:  int64(destroyedContracts),
		TotalExecutions:     totalExecutions,
		TotalFeesCollected:  &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)},
		AverageExecutionFee: &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)},
		TotalDeployers:      int64(max(1, totalContracts)),
		MostActiveContract:  "sample_contract",
		LastActivity:        timestamppb.New(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
	}

	return &blockchainproto.QueryContractStatsResponse{
		Stats: stats,
	}, nil
}

// Helper function for max
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
