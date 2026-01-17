package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/smart_contract/v1/usc/smart_contract/v1"

	blocktypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/block/types"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/smart_contract/types"
)

// MsgServer defines the gRPC message server for the smart contract module using blockchain-proto types
type MsgServer struct {
	Keeper
}

// NewMsgServer creates a new smart contract message server
func NewMsgServer(keeper Keeper) *MsgServer {
	return &MsgServer{Keeper: keeper}
}

// DeployContract handles contract deployment messages
func (k MsgServer) DeployContract(ctx context.Context, msg *blockchainproto.MsgDeployContract) (*blockchainproto.MsgDeployContractResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Basic validation
	if msg.Deployer == "" {
		return nil, fmt.Errorf("deployer cannot be empty")
	}
	if msg.ContractCode == "" {
		return nil, fmt.Errorf("contract code cannot be empty")
	}
	if msg.ContractName == "" {
		return nil, fmt.Errorf("contract name cannot be empty")
	}

	// Generate contract ID and address
	contractID := fmt.Sprintf("contract_%s_%d", msg.ContractName, sdkCtx.BlockHeight())

	// Calculate code hash from contract code
	codeHashBytes := sha256.Sum256([]byte(msg.ContractCode))
	codeHash := hex.EncodeToString(codeHashBytes[:])

	// Get deployer nonce (number of contracts deployed by this deployer)
	// For now, use block height as nonce (will be improved with proper nonce tracking)
	nonce := uint64(sdkCtx.BlockHeight())
	contractAddress := blocktypes.CalculateContractAddress(sdkCtx, msg.Deployer, codeHash, nonce)

	// Calculate real transaction hash
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Deployer, contractAddress, "", "deploy", msg.ContractCode, msg.Memo)

	// Create contract deployment record
	deployment := types.ContractDeployment{
		ID:          contractID,
		ContractID:  contractID,
		Deployer:    msg.Deployer,
		Network:     "usc_network",
		Address:     contractAddress,
		TxHash:      txHash,
		BlockNumber: uint64(sdkCtx.BlockHeight()),
		GasUsed:     1000000, // Default gas used
		DeployedAt:  sdkCtx.BlockTime(),
		Metadata:    map[string]string{"version": msg.ContractVersion},
	}

	// Set the deployment
	if err := k.SetDeployment(sdkCtx, deployment); err != nil {
		return nil, fmt.Errorf("failed to record deployment: %w", err)
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeContractDeployed,
			sdk.NewAttribute(types.AttributeKeyContractID, contractID),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Deployer),
		),
	)

	return &blockchainproto.MsgDeployContractResponse{
		Success:         true,
		ContractAddress: contractAddress,
		DeploymentHash:  codeHash, // Use code hash as deployment hash
		TransactionHash: txHash,
	}, nil
}

// UpdateContract handles contract update messages
func (k MsgServer) UpdateContract(ctx context.Context, msg *blockchainproto.MsgUpdateContract) (*blockchainproto.MsgUpdateContractResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Basic validation
	if msg.ContractAddress == "" {
		return nil, fmt.Errorf("contract address cannot be empty")
	}
	if msg.Updater == "" {
		return nil, fmt.Errorf("updater cannot be empty")
	}
	if msg.NewCode == "" {
		return nil, fmt.Errorf("new code cannot be empty")
	}

	// Find contract by address (simplified - using address as ID)
	contract, err := k.GetContract(sdkCtx, msg.ContractAddress)
	if err != nil {
		return nil, fmt.Errorf("contract not found: %w", err)
	}

	// Update contract fields
	contract.Version = msg.NewVersion
	contract.UpdatedAt = sdkCtx.BlockTime()
	contract.Metadata = map[string]string{"version": msg.NewVersion}

	// Set the updated contract
	if err := k.SetContract(sdkCtx, contract); err != nil {
		return nil, fmt.Errorf("failed to update contract: %w", err)
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeContractUpdated,
			sdk.NewAttribute(types.AttributeKeyContractID, contract.ID),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Updater),
		),
	)

	// Calculate real transaction hash for update
	updateTxHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Updater, msg.ContractAddress, "", "update", msg.NewCode, msg.Memo)

	return &blockchainproto.MsgUpdateContractResponse{
		Success:         true,
		UpdateHash:      blocktypes.CalculateHashFromString(fmt.Sprintf("update_%s_%d", contract.ID, sdkCtx.BlockHeight())),
		TransactionHash: updateTxHash,
	}, nil
}

// ExecuteContract handles contract execution messages
func (k MsgServer) ExecuteContract(ctx context.Context, msg *blockchainproto.MsgExecuteContract) (*blockchainproto.MsgExecuteContractResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Basic validation
	if msg.ContractAddress == "" {
		return nil, fmt.Errorf("contract address cannot be empty")
	}
	if msg.Executor == "" {
		return nil, fmt.Errorf("executor cannot be empty")
	}
	if msg.MethodName == "" {
		return nil, fmt.Errorf("method name cannot be empty")
	}

	// Get existing contract (simplified - using address as ID)
	contract, err := k.GetContract(sdkCtx, msg.ContractAddress)
	if err != nil {
		return nil, fmt.Errorf("contract not found: %w", err)
	}

	// Create execution record
	executionID := fmt.Sprintf("exec_%s_%d", contract.ID, sdkCtx.BlockHeight())
	execution := types.ContractExecution{
		ID:         executionID,
		ContractID: contract.ID,
		Executor:   msg.Executor,
		Method:     msg.MethodName,
		Input:      []byte(fmt.Sprintf("%v", msg.MethodParams)),
		Output:     []byte{}, // Will be set after execution
		GasUsed:    100000,   // Default gas used
		GasLimit:   1000000,  // Default gas limit
		Status:     "pending",
		ExecutedAt: sdkCtx.BlockTime(),
		Metadata:   map[string]string{"memo": msg.Memo},
	}

	// TODO: Implement actual contract execution logic
	// This would typically involve:
	// - Validating execution permissions
	// - Executing the contract method
	// - Processing the result
	// - Updating gas usage

	execution.Status = "success"
	execution.Output = []byte("execution_result")

	// Set the execution
	if err := k.SetExecution(sdkCtx, execution); err != nil {
		return nil, fmt.Errorf("failed to record execution: %w", err)
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeContractExecuted,
			sdk.NewAttribute(types.AttributeKeyContractID, contract.ID),
			sdk.NewAttribute(types.AttributeKeyExecutor, msg.Executor),
		),
	)

	// Calculate real transaction hash for execution
	execTxHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Executor, msg.ContractAddress, "", "execute", fmt.Sprintf("%v", msg.MethodParams), msg.Memo)

	return &blockchainproto.MsgExecuteContractResponse{
		Success:         true,
		ReturnValue:     string(execution.Output),
		Logs:            []string{"execution_completed"},
		ExecutionHash:   executionID,
		TransactionHash: execTxHash,
	}, nil
}

// DestroyContract handles contract destruction messages
func (k MsgServer) DestroyContract(ctx context.Context, msg *blockchainproto.MsgDestroyContract) (*blockchainproto.MsgDestroyContractResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Basic validation
	if msg.ContractAddress == "" {
		return nil, fmt.Errorf("contract address cannot be empty")
	}
	if msg.Destroyer == "" {
		return nil, fmt.Errorf("destroyer cannot be empty")
	}

	// Get existing contract (simplified - using address as ID)
	contract, err := k.GetContract(sdkCtx, msg.ContractAddress)
	if err != nil {
		return nil, fmt.Errorf("contract not found: %w", err)
	}

	// Mark contract as destroyed
	contract.Status = "destroyed"
	contract.UpdatedAt = sdkCtx.BlockTime()

	// Set the updated contract
	if err := k.SetContract(sdkCtx, contract); err != nil {
		return nil, fmt.Errorf("failed to update contract: %w", err)
	}

	// Emit event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeContractDestroyed,
			sdk.NewAttribute(types.AttributeKeyContractID, contract.ID),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Destroyer),
		),
	)

	// Calculate real transaction hash for destruction
	destroyTxHash := blocktypes.CalculateTransactionHash(sdkCtx, msg.Destroyer, msg.ContractAddress, "", "destroy", "", msg.Memo)

	return &blockchainproto.MsgDestroyContractResponse{
		Success:         true,
		DestructionHash: blocktypes.CalculateHashFromString(fmt.Sprintf("destroy_%s_%d", contract.ID, sdkCtx.BlockHeight())),
		TransactionHash: destroyTxHash,
	}, nil
}
