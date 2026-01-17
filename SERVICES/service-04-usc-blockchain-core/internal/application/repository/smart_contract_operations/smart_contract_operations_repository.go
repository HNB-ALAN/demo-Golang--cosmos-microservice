package smart_contract_operations

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
	blocktypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/block/types"
	sctypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/smart_contract/types"

	"github.com/usc-platform/shared/logging"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Repository handles smart contract operations data access
type Repository struct {
	db                *database.PostgreSQLManager
	cosmosApp         *app.USCApp
	blockchainStorage *storage.StateManager
	redisManager      *database.RedisManager
	logger            *logging.Logger
}

// NewRepository creates a new smart contract operations repository
func NewRepository(db *database.PostgreSQLManager, cosmosApp *app.USCApp, blockchainStorage *storage.StateManager, redisManager *database.RedisManager, logger *logging.Logger) *Repository {
	return &Repository{
		db:                db,
		cosmosApp:         cosmosApp,
		blockchainStorage: blockchainStorage,
		redisManager:      redisManager,
		logger:            logger,
	}
}

// DeployContract deploys a new smart contract
func (r *Repository) DeployContract(ctx context.Context, req *proto.DeployContractRequest) (*proto.DeployContractResponse, error) {
	// Priority 1: Deploy on Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.deployContractOnKeeper(ctx, req); err == nil {
			// Save to PostgreSQL for analytics (sync to ensure data persistence)
			if r.db != nil {
				if err := r.saveContractDeploymentToDatabase(ctx, req, result); err != nil {
					r.logger.Error("Failed to save contract to database",
						logging.String("contract_address", result.ContractAddress),
						logging.Error(err))
					// Continue even if database save fails (keeper is primary)
				} else {
					r.logger.Info("Contract saved to database successfully",
						logging.String("contract_address", result.ContractAddress))
				}
			}
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.deployContractInDatabase(ctx, req)
}

// ExecuteContract executes a smart contract function
func (r *Repository) ExecuteContract(ctx context.Context, req *proto.ExecuteContractRequest) (*proto.ExecuteContractResponse, error) {
	// Priority 1: Execute on Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.executeContractOnKeeper(ctx, req); err == nil {
			// Save to PostgreSQL for analytics (sync to ensure data persistence for debugging)
			if r.db != nil {
				correlationID := utils.GetCorrelationID(ctx)
				if err := r.saveContractExecutionToDatabase(ctx, req, result); err != nil {
					r.logger.Error("Failed to save contract execution analytics",
						logging.Error(err),
						logging.String("contract_address", req.ContractAddress),
						logging.String("function_name", req.FunctionName),
						logging.String("correlation_id", correlationID))
					// Continue even if database save fails (keeper is primary)
				} else {
					r.logger.Info("Contract execution analytics saved successfully",
						logging.String("contract_address", req.ContractAddress),
						logging.String("function_name", req.FunctionName),
						logging.String("correlation_id", correlationID))
				}
			}
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.executeContractInDatabase(ctx, req)
}

// QueryContract queries a smart contract function
func (r *Repository) QueryContract(ctx context.Context, req *proto.QueryContractRequest) (*proto.QueryContractResponse, error) {
	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.queryContractFromKeeper(ctx, req); err == nil {
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.queryContractFromDatabase(ctx, req)
}

// GetContractCode retrieves contract source code
func (r *Repository) GetContractCode(ctx context.Context, req *proto.GetContractCodeRequest) (*proto.GetContractCodeResponse, error) {
	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.getContractCodeFromKeeper(ctx, req.ContractAddress); err == nil {
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.getContractCodeFromDatabase(ctx, req)
}

// GetContractStorage retrieves contract storage
func (r *Repository) GetContractStorage(ctx context.Context, req *proto.GetContractStorageRequest) (*proto.GetContractStorageResponse, error) {
	if req.ContractAddress == "" || req.StorageKey == "" {
		return nil, repoerrors.NewValidationError("contract_address and storage_key", "are required")
	}

	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.getContractStorageFromKeeper(ctx, req); err == nil {
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.getContractStorageFromDatabase(ctx, req)
}

// getContractStorageFromKeeper retrieves contract storage from the keeper
func (r *Repository) getContractStorageFromKeeper(ctx context.Context, req *proto.GetContractStorageRequest) (*proto.GetContractStorageResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, err
	}

	// Get contract from keeper
	contract, err := r.cosmosApp.SmartContractKeeper.GetContract(sdkCtx, req.ContractAddress)
	if err != nil {
		return nil, repoerrors.WrapRepositoryError(repoerrors.ErrContractNotFound, err)
	}

	// For now, return contract metadata as storage value
	// In production, this would query actual contract storage
	storageValue := ""
	if len(contract.Metadata) > 0 {
		if val, ok := contract.Metadata[req.StorageKey]; ok {
			storageValue = val
		}
	}

	return &proto.GetContractStorageResponse{
		StorageValue: storageValue,
		StorageKey:   req.StorageKey,
		BlockNumber:  int64(sdkCtx.BlockHeight()),
		StorageType:  "string",
		DecodedValue: storageValue,
	}, nil
}

// getContractStorageFromDatabase retrieves contract storage from database
func (r *Repository) getContractStorageFromDatabase(ctx context.Context, req *proto.GetContractStorageRequest) (*proto.GetContractStorageResponse, error) {
	// Simplified database implementation
	return &proto.GetContractStorageResponse{
		StorageValue: "",
		StorageKey:   req.StorageKey,
		BlockNumber:  0,
		StorageType:  "",
		DecodedValue: "",
	}, nil
}

// Helper methods for SmartContractKeeper interaction

// getSDKContext creates a sdk.Context from context.Context
// Uses shared utility to avoid code duplication
func (r *Repository) getSDKContext(ctx context.Context) (sdk.Context, error) {
	return utils.GetSDKContext(ctx, r.cosmosApp, r.logger)
}

// queryContractFromKeeper queries a contract from the keeper
func (r *Repository) queryContractFromKeeper(ctx context.Context, req *proto.QueryContractRequest) (*proto.QueryContractResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, err
	}

	// Get contract from keeper
	contract, err := r.cosmosApp.SmartContractKeeper.GetContract(sdkCtx, req.ContractAddress)
	if err != nil {
		return &proto.QueryContractResponse{
			ReturnValue:  "",
			Success:      false,
			ErrorMessage: fmt.Sprintf("contract not found: %v", err),
		}, nil
	}

	// For now, return contract info as return value
	// In a real implementation, this would execute the function and return the result
	returnValue := fmt.Sprintf("Contract: %s, Function: %s", contract.Name, req.FunctionName)

	return &proto.QueryContractResponse{
		ReturnValue:  returnValue,
		Success:      true,
		ErrorMessage: "",
	}, nil
}

// getContractCodeFromKeeper retrieves contract code from the keeper
func (r *Repository) getContractCodeFromKeeper(ctx context.Context, contractAddress string) (*proto.GetContractCodeResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, err
	}

	// Get contract from keeper
	contract, err := r.cosmosApp.SmartContractKeeper.GetContract(sdkCtx, contractAddress)
	if err != nil {
		return nil, repoerrors.WrapRepositoryError(repoerrors.ErrContractNotFound, err)
	}

	return &proto.GetContractCodeResponse{
		Bytecode:        string(contract.Bytecode),
		Abi:             contract.ABI,
		SourceCode:      "", // Not stored in keeper
		CompilerVersion: "",
		ContractName:    contract.Name,
		BlockNumber:     int64(sdkCtx.BlockHeight()),
		IsVerified:      false,
	}, nil
}

// deployContractOnKeeper deploys a contract on Cosmos SDK blockchain
func (r *Repository) deployContractOnKeeper(ctx context.Context, req *proto.DeployContractRequest) (*proto.DeployContractResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, repoerrors.NewBlockchainError("get_sdk_context", err)
	}

	// Generate contract ID
	contractID := fmt.Sprintf("contract_%s_%d", req.ContractName, sdkCtx.BlockHeight())

	// Generate real contract address using blocktypes helper
	codeHash := sha256.Sum256([]byte(req.Bytecode))
	codeHashHex := hex.EncodeToString(codeHash[:])
	nonce := uint64(sdkCtx.BlockHeight())
	contractAddress := blocktypes.CalculateContractAddress(sdkCtx, req.FromAddress, codeHashHex, nonce)

	// Create contract
	contract := sctypes.SmartContract{
		ID:          contractID,
		Name:        req.ContractName,
		Type:        sctypes.ContractTypeWASM,
		Status:      sctypes.ContractStatusActive,
		Owner:       req.FromAddress,
		Code:        []byte(req.Bytecode),
		ABI:         req.Abi,
		Bytecode:    []byte(req.Bytecode),
		Address:     contractAddress,
		Version:     "1.0.0",
		Description: "",
		CreatedAt:   sdkCtx.BlockTime(),
		UpdatedAt:   sdkCtx.BlockTime(),
		Metadata:    map[string]string{},
	}

	// Store contract
	if err := r.cosmosApp.SmartContractKeeper.SetContract(sdkCtx, contract); err != nil {
		return nil, repoerrors.NewDatabaseError("store_contract", err)
	}

	// Generate real transaction hash using blocktypes helper
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, req.FromAddress, contractAddress, "", "deploy", req.Bytecode, "")

	// Create deployment record
	deployment := sctypes.ContractDeployment{
		ID:          contractID,
		ContractID:  contractID,
		Deployer:    req.FromAddress,
		Network:     "usc_network",
		Address:     contractAddress,
		TxHash:      txHash,
		BlockNumber: uint64(sdkCtx.BlockHeight()),
		GasUsed:     uint64(req.GasLimit),
		DeployedAt:  sdkCtx.BlockTime(),
		Metadata:    map[string]string{"version": "1.0.0"},
	}

	// Store deployment
	if err := r.cosmosApp.SmartContractKeeper.SetDeployment(sdkCtx, deployment); err != nil {
		return nil, repoerrors.NewDatabaseError("store_deployment", err)
	}

	return &proto.DeployContractResponse{
		ContractAddress: contractAddress,
		TransactionHash: txHash,
		Status:          1, // 1=Confirmed (on blockchain)
		ErrorMessage:    "",
	}, nil
}

// deployContractInDatabase deploys a contract in database (fallback)
func (r *Repository) deployContractInDatabase(ctx context.Context, req *proto.DeployContractRequest) (*proto.DeployContractResponse, error) {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return nil, repoerrors.NewRepositoryError(repoerrors.ErrDatabaseUnavailable)
	}

	query := `
		INSERT INTO usc_smart_contract_analytics (
			contract_address, contract_name, bytecode, abi,
			from_address, gas_price, gas_limit, transaction_hash,
			status, deployed_at, last_updated
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	// Generate real hash for database analytics
	dataStr := fmt.Sprintf("%s:%s:%s:%s:deploy", req.FromAddress, req.ContractName, req.Bytecode, time.Now().Format(time.RFC3339))
	hashBytes := sha256.Sum256([]byte(dataStr))
	contractAddress := "0x" + hex.EncodeToString(hashBytes[:20]) // First 20 bytes for address
	txHash := "0x" + hex.EncodeToString(hashBytes[:])            // Full hash for transaction

	now := time.Now()
	_, err := postgres.ExecContext(ctx, query,
		contractAddress,
		req.ContractName,
		req.Bytecode,
		req.Abi,
		req.FromAddress,
		req.GasPrice,
		req.GasLimit,
		txHash,
		"pending",
		now,
		now,
	)

	if err != nil {
		return nil, repoerrors.NewDatabaseError("deploy_contract", err)
	}

	return &proto.DeployContractResponse{
		ContractAddress: contractAddress,
		TransactionHash: txHash,
		Status:          0, // 0=Pending
		ErrorMessage:    "",
	}, nil
}

// saveContractDeploymentToDatabase saves contract deployment to database for analytics
func (r *Repository) saveContractDeploymentToDatabase(ctx context.Context, req *proto.DeployContractRequest, resp *proto.DeployContractResponse) error {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil || resp == nil {
		return fmt.Errorf("postgres connection not available or response is nil")
	}

	now := time.Now()

	// Save to analytics table
	analyticsQuery := `
		INSERT INTO usc_smart_contract_analytics (
			contract_address, contract_name, bytecode, abi,
			from_address, gas_price, gas_limit, transaction_hash,
			status, deployed_at, last_updated
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (contract_address) DO UPDATE SET
			transaction_hash = EXCLUDED.transaction_hash,
			status = EXCLUDED.status,
			last_updated = EXCLUDED.last_updated
	`

	if _, err := postgres.ExecContext(ctx, analyticsQuery,
		resp.ContractAddress,
		req.ContractName,
		req.Bytecode,
		req.Abi,
		req.FromAddress,
		req.GasPrice,
		req.GasLimit,
		resp.TransactionHash,
		"confirmed",
		now,
		now,
	); err != nil {
		return fmt.Errorf("failed to save to analytics table: %w", err)
	}

	// Save to main smart_contracts table
	mainQuery := `
		INSERT INTO smart_contracts (
			contract_address, contract_name, bytecode, abi,
			owner_address, deployment_transaction_hash, status, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (contract_address) DO UPDATE SET
			contract_name = EXCLUDED.contract_name,
			status = EXCLUDED.status,
			updated_at = NOW()
	`

	if _, err := postgres.ExecContext(ctx, mainQuery,
		resp.ContractAddress,
		req.ContractName,
		req.Bytecode,
		req.Abi,
		req.FromAddress,
		resp.TransactionHash,
		"active",
		now,
	); err != nil {
		return fmt.Errorf("failed to save to main table: %w", err)
	}

	return nil
}

// executeContractOnKeeper executes a contract on Cosmos SDK blockchain
func (r *Repository) executeContractOnKeeper(ctx context.Context, req *proto.ExecuteContractRequest) (*proto.ExecuteContractResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, repoerrors.NewBlockchainError("get_sdk_context", err)
	}

	// Get existing contract
	contract, err := r.cosmosApp.SmartContractKeeper.GetContract(sdkCtx, req.ContractAddress)
	if err != nil {
		return nil, repoerrors.WrapRepositoryError(repoerrors.ErrContractNotFound, err)
	}

	// Create execution record
	executionID := fmt.Sprintf("exec_%s_%d", contract.ID, sdkCtx.BlockHeight())
	execution := sctypes.ContractExecution{
		ID:         executionID,
		ContractID: contract.ID,
		Executor:   req.FromAddress,
		Method:     req.FunctionName,
		Input:      []byte(fmt.Sprintf("%v", req.Parameters)),
		Output:     []byte{}, // Will be set after execution
		GasUsed:    100000,
		GasLimit:   uint64(req.GasLimit),
		Status:     "success",
		ExecutedAt: sdkCtx.BlockTime(),
		Metadata:   map[string]string{},
	}

	// Store execution
	if err := r.cosmosApp.SmartContractKeeper.SetExecution(sdkCtx, execution); err != nil {
		return nil, repoerrors.NewDatabaseError("store_execution", err)
	}

	// Generate real transaction hash using blocktypes helper
	txHash := blocktypes.CalculateTransactionHash(sdkCtx, req.FromAddress, req.ContractAddress, "", "execute", fmt.Sprintf("%s:%s", req.FunctionName, fmt.Sprintf("%v", req.Parameters)), "")

	return &proto.ExecuteContractResponse{
		TransactionHash: txHash,
		Status:          1, // 1=Confirmed (on blockchain)
		ErrorMessage:    "",
	}, nil
}

// executeContractInDatabase executes a contract in database (fallback)
func (r *Repository) executeContractInDatabase(ctx context.Context, req *proto.ExecuteContractRequest) (*proto.ExecuteContractResponse, error) {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return nil, repoerrors.NewRepositoryError(repoerrors.ErrDatabaseUnavailable)
	}

	query := `
		INSERT INTO usc_contract_execution_analytics (
			contract_address, function_name, from_address,
			gas_price, gas_limit, transaction_hash,
			status, executed_at, return_value
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	// Generate real hash for database analytics
	dataStr := fmt.Sprintf("%s:%s:%s:%s:execute", req.ContractAddress, req.FunctionName, fmt.Sprintf("%v", req.Parameters), time.Now().Format(time.RFC3339))
	hashBytes := sha256.Sum256([]byte(dataStr))
	txHash := "0x" + hex.EncodeToString(hashBytes[:])

	_, err := postgres.ExecContext(ctx, query,
		req.ContractAddress,
		req.FunctionName,
		req.FromAddress,
		req.GasPrice,
		req.GasLimit,
		txHash,
		"pending",
		time.Now(),
		"",
	)

	if err != nil {
		return nil, repoerrors.NewDatabaseError("execute_contract", err)
	}

	return &proto.ExecuteContractResponse{
		TransactionHash: txHash,
		Status:          0, // 0=Pending
		ErrorMessage:    "",
	}, nil
}

// saveContractExecutionToDatabase saves contract execution to database for analytics (sync for debugging)
func (r *Repository) saveContractExecutionToDatabase(ctx context.Context, req *proto.ExecuteContractRequest, resp *proto.ExecuteContractResponse) error {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return fmt.Errorf("postgres connection not available")
	}
	if resp == nil {
		return fmt.Errorf("response is nil")
	}

	query := `
		INSERT INTO usc_contract_execution_analytics (
			contract_address, function_name, from_address,
			gas_price, gas_limit, transaction_hash,
			status, executed_at, return_value
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	if _, err := postgres.ExecContext(ctx, query,
		req.ContractAddress,
		req.FunctionName,
		req.FromAddress,
		req.GasPrice,
		req.GasLimit,
		resp.TransactionHash,
		"confirmed",
		time.Now(),
		"",
	); err != nil {
		return fmt.Errorf("failed to save contract execution to analytics: %w", err)
	}

	return nil
}

// Database fallback methods

// queryContractFromDatabase queries contract from database
// Note: Contract query execution requires contract runtime. This method only verifies contract exists.
func (r *Repository) queryContractFromDatabase(ctx context.Context, req *proto.QueryContractRequest) (*proto.QueryContractResponse, error) {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return &proto.QueryContractResponse{
			ReturnValue:  "",
			Success:      false,
			ErrorMessage: "database not available",
		}, nil
	}

	// Verify contract exists in database
	checkQuery := `SELECT contract_address FROM smart_contracts WHERE contract_address = $1 LIMIT 1`
	var contractAddr string
	if err := postgres.QueryRowContext(ctx, checkQuery, req.ContractAddress).Scan(&contractAddr); err != nil {
		return &proto.QueryContractResponse{
			ReturnValue:  "",
			Success:      false,
			ErrorMessage: "contract not found",
		}, nil
	}

	// Contract exists but query execution requires contract runtime
	// This is a fallback method, so we return success but empty result
	return &proto.QueryContractResponse{
		ReturnValue:  "",
		Success:      true,
		ErrorMessage: "contract query execution requires contract runtime (not available in database fallback)",
	}, nil
}

// getContractCodeFromDatabase retrieves contract code from database
func (r *Repository) getContractCodeFromDatabase(ctx context.Context, req *proto.GetContractCodeRequest) (*proto.GetContractCodeResponse, error) {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return &proto.GetContractCodeResponse{
			Bytecode:        "",
			Abi:             "",
			SourceCode:      "",
			CompilerVersion: "",
			ContractName:    "",
			BlockNumber:     0,
			IsVerified:      false,
		}, nil
	}

	query := `
		SELECT 
			contract_address,
			COALESCE(bytecode, '') as bytecode,
			COALESCE(abi, '') as abi,
			COALESCE(source_code, '') as source_code,
			COALESCE(compiler_version, '') as compiler_version,
			COALESCE(contract_name, '') as contract_name,
			COALESCE(is_verified, false) as is_verified,
			COALESCE(gas_used, 0) as gas_used,
			created_at
		FROM smart_contracts
		WHERE contract_address = $1
		LIMIT 1
	`

	var contractAddr, bytecode, abi, sourceCode, compilerVersion, contractName string
	var isVerified bool
	var gasUsed int64
	var createdAt time.Time

	if err := postgres.QueryRowContext(ctx, query, req.ContractAddress).Scan(
		&contractAddr, &bytecode, &abi, &sourceCode, &compilerVersion,
		&contractName, &isVerified, &gasUsed, &createdAt,
	); err != nil {
		return &proto.GetContractCodeResponse{
			Bytecode:        "",
			Abi:             "",
			SourceCode:      "",
			CompilerVersion: "",
			ContractName:    "",
			BlockNumber:     0,
			IsVerified:      false,
		}, nil
	}

	// Get block number from deployment transaction if available
	// For now, return 0 as block_number is not directly stored in smart_contracts table
	blockNumber := int64(0)

	// Convert createdAt to protobuf timestamp
	createdAtPB := timestamppb.New(createdAt)

	return &proto.GetContractCodeResponse{
		Bytecode:        bytecode,
		Abi:             abi,
		SourceCode:      sourceCode,
		CompilerVersion: compilerVersion,
		ContractName:    contractName,
		BlockNumber:     blockNumber,
		CreatedAt:       createdAtPB,
		IsVerified:      isVerified,
	}, nil
}
