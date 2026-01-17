package validator_operations

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
	validatortypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/validator/types"

	"github.com/usc-platform/shared/logging"
)

// Repository handles validator operations data access
type Repository struct {
	db                *database.PostgreSQLManager
	cosmosApp         *app.USCApp
	blockchainStorage *storage.StateManager
	redisManager      *database.RedisManager
	logger            *logging.Logger
}

// NewRepository creates a new validator operations repository
func NewRepository(db *database.PostgreSQLManager, cosmosApp *app.USCApp, blockchainStorage *storage.StateManager, redisManager *database.RedisManager, logger *logging.Logger) *Repository {
	return &Repository{
		db:                db,
		cosmosApp:         cosmosApp,
		blockchainStorage: blockchainStorage,
		redisManager:      redisManager,
		logger:            logger,
	}
}

// RegisterValidator registers a new validator
func (r *Repository) RegisterValidator(ctx context.Context, req *proto.RegisterValidatorRequest) (*proto.RegisterValidatorResponse, error) {
	// Priority 1: Register on Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.registerValidatorOnKeeper(ctx, req); err == nil {
			// Save to PostgreSQL for analytics (sync to ensure data persistence)
			if r.db != nil {
				if err := r.saveValidatorToDatabase(ctx, req, result); err != nil {
					r.logger.Error("Failed to save validator to database",
						logging.String("validator_address", req.ValidatorAddress),
						logging.Error(err))
					// Continue even if database save fails (keeper is primary)
				} else {
					r.logger.Info("Validator saved to database successfully",
						logging.String("validator_address", req.ValidatorAddress))
				}
			}
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.registerValidatorInDatabase(ctx, req)
}

// GetValidatorStatus retrieves validator status
func (r *Repository) GetValidatorStatus(ctx context.Context, req *proto.GetValidatorStatusRequest) (*proto.GetValidatorStatusResponse, error) {
	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if status, err := r.getValidatorStatusFromKeeper(ctx, req.ValidatorAddress); err == nil {
			return status, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.getValidatorStatusFromDatabase(ctx, req)
}

// StakeUSC stakes USC tokens for a validator
func (r *Repository) StakeUSC(ctx context.Context, req *proto.StakeUSCRequest) (*proto.StakeUSCResponse, error) {
	// Priority 1: Stake on Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.stakeUSCOnKeeper(ctx, req); err == nil {
			// Save to PostgreSQL for analytics (sync to ensure data persistence)
			if r.db != nil {
				if err := r.saveStakingToDatabase(ctx, req, result, "stake"); err != nil {
					r.logger.Error("Failed to save staking to database",
						logging.String("delegator_address", req.DelegatorAddress),
						logging.String("validator_address", req.ValidatorAddress),
						logging.Error(err))
					// Continue even if database save fails (keeper is primary)
				} else {
					r.logger.Info("Staking saved to database successfully",
						logging.String("delegator_address", req.DelegatorAddress),
						logging.String("validator_address", req.ValidatorAddress))
				}
			}
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.stakeUSCInDatabase(ctx, req)
}

// UnstakeUSC unstakes USC tokens from a validator
func (r *Repository) UnstakeUSC(ctx context.Context, req *proto.UnstakeUSCRequest) (*proto.UnstakeUSCResponse, error) {
	// Priority 1: Unstake on Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if result, err := r.unstakeUSCOnKeeper(ctx, req); err == nil {
			// Save to PostgreSQL for analytics (sync to ensure data persistence)
			if r.db != nil {
				if err := r.saveStakingToDatabase(ctx, req, result, "unstake"); err != nil {
					r.logger.Error("Failed to save unstaking to database",
						logging.String("delegator_address", req.DelegatorAddress),
						logging.String("validator_address", req.ValidatorAddress),
						logging.Error(err))
					// Continue even if database save fails (keeper is primary)
				} else {
					r.logger.Info("Unstaking saved to database successfully",
						logging.String("delegator_address", req.DelegatorAddress),
						logging.String("validator_address", req.ValidatorAddress))
				}
			}
			return result, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.unstakeUSCInDatabase(ctx, req)
}

// GetValidators retrieves list of validators
func (r *Repository) GetValidators(ctx context.Context, req *proto.GetValidatorsRequest) (*proto.GetValidatorsResponse, error) {
	limit, offset := utils.NormalizePagination(req.Limit, req.Offset, utils.PaginationConfig{
		DefaultLimit: 100,
		MaxLimit:     1000,
	})

	// Priority 1: Keeper (RocksDB)
	if utils.IsCosmosAppAvailable(r.cosmosApp) {
		if validators, err := r.getValidatorsFromKeeper(ctx, limit, offset); err == nil {
			return validators, nil
		}
	}

	// Priority 2: PostgreSQL (fallback)
	return r.getValidatorsFromDatabase(ctx, req)
}

// Helper methods for ValidatorKeeper interaction

// getSDKContext creates a sdk.Context from context.Context
// Uses shared utility to avoid code duplication
func (r *Repository) getSDKContext(ctx context.Context) (sdk.Context, error) {
	return utils.GetSDKContext(ctx, r.cosmosApp, r.logger)
}

// getValidatorsFromKeeper retrieves validators from ValidatorKeeper
func (r *Repository) getValidatorsFromKeeper(ctx context.Context, limit, offset int32) (*proto.GetValidatorsResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, err
	}

	allValidators, err := r.cosmosApp.ValidatorKeeper.GetAllValidators(sdkCtx)
	if err != nil {
		return nil, repoerrors.NewDatabaseError("get_validators", err)
	}

	// Apply pagination
	start := int(offset)
	end := start + int(limit)
	if end > len(allValidators) {
		end = len(allValidators)
	}

	// Pre-allocate slice with capacity = (end - start) for better performance
	validators := make([]*proto.ValidatorInfo, 0, end-start)
	for i := start; i < end; i++ {
		validators = append(validators, r.convertValidatorToProto(&allValidators[i]))
	}

	return &proto.GetValidatorsResponse{
		Validators: validators,
		TotalCount: int32(len(allValidators)),
		HasMore:    end < len(allValidators),
	}, nil
}

// getValidatorStatusFromKeeper retrieves validator status from ValidatorKeeper
func (r *Repository) getValidatorStatusFromKeeper(ctx context.Context, address string) (*proto.GetValidatorStatusResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, err
	}

	validator, err := r.cosmosApp.ValidatorKeeper.GetValidator(sdkCtx, address)
	if err != nil {
		return &proto.GetValidatorStatusResponse{
			Success:      false,
			ErrorMessage: fmt.Sprintf("validator not found: %v", err),
		}, nil
	}

	return &proto.GetValidatorStatusResponse{
		Success:   true,
		Validator: r.convertValidatorToProto(&validator),
	}, nil
}

// convertValidatorToProto converts validatortypes.Validator to proto.ValidatorInfo
func (r *Repository) convertValidatorToProto(v *validatortypes.Validator) *proto.ValidatorInfo {
	return &proto.ValidatorInfo{
		ValidatorId:      v.Address, // Use address as ID
		ValidatorAddress: v.Address,
		ValidatorName:    v.Description,
		Description:      v.Description,
		CommissionRate:   v.Commission,
		Status:           v.Status,
		StakeAmount:      fmt.Sprintf("%d", v.Power), // Convert power to string
	}
}

// registerValidatorOnKeeper registers a validator on Cosmos SDK blockchain
func (r *Repository) registerValidatorOnKeeper(ctx context.Context, req *proto.RegisterValidatorRequest) (*proto.RegisterValidatorResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, repoerrors.NewBlockchainError("get_sdk_context", err)
	}

	validator := validatortypes.NewValidator(
		req.ValidatorAddress,
		req.ValidatorPublicKey,
		req.ValidatorName,
		req.CommissionRate,
	)

	if err := r.cosmosApp.ValidatorKeeper.CreateValidator(sdkCtx, validator); err != nil {
		return nil, repoerrors.WrapRepositoryError(repoerrors.ErrValidatorRegistrationFailed, err)
	}

	return &proto.RegisterValidatorResponse{
		Success:          true,
		ValidatorAddress: req.ValidatorAddress,
		ErrorMessage:     "",
	}, nil
}

// registerValidatorInDatabase registers a validator in database (fallback)
func (r *Repository) registerValidatorInDatabase(ctx context.Context, req *proto.RegisterValidatorRequest) (*proto.RegisterValidatorResponse, error) {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return nil, repoerrors.NewRepositoryError(repoerrors.ErrDatabaseUnavailable)
	}

	query := `
		INSERT INTO usc_validator_analytics (
			validator_address, validator_name, validator_public_key, commission_rate,
			status, registered_at, last_updated
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (validator_address) DO UPDATE SET
			validator_name = EXCLUDED.validator_name,
			validator_public_key = EXCLUDED.validator_public_key,
			commission_rate = EXCLUDED.commission_rate,
			status = EXCLUDED.status,
			last_updated = EXCLUDED.last_updated
	`

	now := time.Now()
	_, err := postgres.ExecContext(ctx, query,
		req.ValidatorAddress,
		req.ValidatorName,
		req.ValidatorPublicKey,
		req.CommissionRate,
		"active",
		now,
		now,
	)

	if err != nil {
		return nil, repoerrors.NewDatabaseError("register_validator", err)
	}

	// Save to main validators table
	mainQuery := `
		INSERT INTO validators (
			validator_id, validator_address, validator_name, validator_public_key,
			commission_rate, status, user_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (validator_address) DO UPDATE SET
			validator_name = EXCLUDED.validator_name,
			validator_public_key = EXCLUDED.validator_public_key,
			commission_rate = EXCLUDED.commission_rate,
			status = EXCLUDED.status,
			updated_at = NOW()
	`

	if _, err := postgres.ExecContext(ctx, mainQuery,
		req.ValidatorAddress,
		req.ValidatorAddress,
		req.ValidatorName,
		req.ValidatorPublicKey,
		req.CommissionRate,
		"active",
		req.UserId,
	); err != nil {
		r.logger.Error("Failed to save validator to main table",
			logging.Error(err),
			logging.String("validator_address", req.ValidatorAddress),
			logging.String("validator_name", req.ValidatorName))
		// Continue even if database save fails (keeper is primary)
	}

	return &proto.RegisterValidatorResponse{
		Success:          true,
		ValidatorAddress: req.ValidatorAddress,
		ErrorMessage:     "",
	}, nil
}

// saveValidatorToDatabase saves validator to database for analytics
func (r *Repository) saveValidatorToDatabase(ctx context.Context, req *proto.RegisterValidatorRequest, resp *proto.RegisterValidatorResponse) error {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil || resp == nil {
		return fmt.Errorf("postgres connection not available or response is nil")
	}

	now := time.Now()

	// Save to analytics table
	analyticsQuery := `
		INSERT INTO usc_validator_analytics (
			validator_address, validator_name, validator_public_key, commission_rate,
			status, registered_at, last_updated
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (validator_address) DO UPDATE SET
			validator_name = EXCLUDED.validator_name,
			validator_public_key = EXCLUDED.validator_public_key,
			commission_rate = EXCLUDED.commission_rate,
			status = EXCLUDED.status,
			last_updated = EXCLUDED.last_updated
	`

	if _, err := postgres.ExecContext(ctx, analyticsQuery,
		req.ValidatorAddress,
		req.ValidatorName,
		req.ValidatorPublicKey,
		req.CommissionRate,
		"active",
		now,
		now,
	); err != nil {
		return fmt.Errorf("failed to save to analytics table: %w", err)
	}

	// Save to main validators table
	mainQuery := `
		INSERT INTO validators (
			validator_id, validator_address, validator_name, validator_public_key,
			commission_rate, status, user_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (validator_address) DO UPDATE SET
			validator_name = EXCLUDED.validator_name,
			validator_public_key = EXCLUDED.validator_public_key,
			commission_rate = EXCLUDED.commission_rate,
			status = EXCLUDED.status,
			updated_at = NOW()
	`

	if _, err := postgres.ExecContext(ctx, mainQuery,
		req.ValidatorAddress,
		req.ValidatorAddress,
		req.ValidatorName,
		req.ValidatorPublicKey,
		req.CommissionRate,
		"active",
		req.UserId,
	); err != nil {
		return fmt.Errorf("failed to save to main table: %w", err)
	}

	return nil
}

// stakeUSCOnKeeper stakes USC on Cosmos SDK blockchain
func (r *Repository) stakeUSCOnKeeper(ctx context.Context, req *proto.StakeUSCRequest) (*proto.StakeUSCResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, repoerrors.NewBlockchainError("get_sdk_context", err)
	}

	delegation := validatortypes.NewDelegation(
		req.DelegatorAddress,
		req.ValidatorAddress,
		req.StakeAmount,
	)

	if err := r.cosmosApp.ValidatorKeeper.Delegate(sdkCtx, delegation); err != nil {
		return nil, repoerrors.WrapRepositoryError(repoerrors.ErrStakingFailed, err)
	}

	txHash := blocktypes.CalculateTransactionHash(sdkCtx, req.DelegatorAddress, req.ValidatorAddress, req.StakeAmount, "stake_usc", "", "")

	return &proto.StakeUSCResponse{
		Success:         true,
		TransactionHash: txHash,
		Status:          1, // Confirmed (on blockchain)
		ErrorMessage:    "",
	}, nil
}

// stakeUSCInDatabase stakes USC in database (fallback)
func (r *Repository) stakeUSCInDatabase(ctx context.Context, req *proto.StakeUSCRequest) (*proto.StakeUSCResponse, error) {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return nil, repoerrors.NewRepositoryError(repoerrors.ErrDatabaseUnavailable)
	}

	query := `
		INSERT INTO usc_staking_analytics (
			validator_address, staker_address, amount, transaction_type,
			status, timestamp, transaction_hash
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	// Generate real hash for database analytics
	dataStr := fmt.Sprintf("%s:%s:%s:%s:stake", req.DelegatorAddress, req.ValidatorAddress, req.StakeAmount, time.Now().Format(time.RFC3339))
	hashBytes := sha256.Sum256([]byte(dataStr))
	txHash := "0x" + hex.EncodeToString(hashBytes[:])

	_, err := postgres.ExecContext(ctx, query,
		req.ValidatorAddress,
		req.DelegatorAddress,
		req.StakeAmount,
		"stake",
		"confirmed",
		time.Now(),
		txHash,
	)

	if err != nil {
		return nil, repoerrors.NewDatabaseError("stake_usc", err)
	}

	return &proto.StakeUSCResponse{
		Success:         true,
		TransactionHash: txHash,
		Status:          0, // 0=Pending
		ErrorMessage:    "",
	}, nil
}

// unstakeUSCOnKeeper unstakes USC on Cosmos SDK blockchain
func (r *Repository) unstakeUSCOnKeeper(ctx context.Context, req *proto.UnstakeUSCRequest) (*proto.UnstakeUSCResponse, error) {
	sdkCtx, err := r.getSDKContext(ctx)
	if err != nil {
		return nil, repoerrors.NewBlockchainError("get_sdk_context", err)
	}

	// Check if delegation exists
	if _, err = r.cosmosApp.ValidatorKeeper.GetDelegation(sdkCtx, req.DelegatorAddress, req.ValidatorAddress); err != nil {
		return nil, repoerrors.NewNotFoundError("delegation", "")
	}

	if err := r.cosmosApp.ValidatorKeeper.Undelegate(sdkCtx, req.DelegatorAddress, req.ValidatorAddress); err != nil {
		return nil, repoerrors.WrapRepositoryError(repoerrors.ErrUnstakingFailed, err)
	}

	txHash := blocktypes.CalculateTransactionHash(sdkCtx, req.DelegatorAddress, req.ValidatorAddress, req.UnstakeAmount, "unstake_usc", "", "")

	return &proto.UnstakeUSCResponse{
		Success:         true,
		TransactionHash: txHash,
		Status:          1, // Confirmed (on blockchain)
		ErrorMessage:    "",
	}, nil
}

// unstakeUSCInDatabase unstakes USC in database (fallback)
func (r *Repository) unstakeUSCInDatabase(ctx context.Context, req *proto.UnstakeUSCRequest) (*proto.UnstakeUSCResponse, error) {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return nil, repoerrors.NewRepositoryError(repoerrors.ErrDatabaseUnavailable)
	}

	query := `
		INSERT INTO usc_staking_analytics (
			validator_address, staker_address, amount, transaction_type,
			status, timestamp, transaction_hash
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	// Generate real hash for database analytics
	dataStr := fmt.Sprintf("%s:%s:%s:%s:unstake", req.DelegatorAddress, req.ValidatorAddress, req.UnstakeAmount, time.Now().Format(time.RFC3339))
	hashBytes := sha256.Sum256([]byte(dataStr))
	txHash := "0x" + hex.EncodeToString(hashBytes[:])

	_, err := postgres.ExecContext(ctx, query,
		req.ValidatorAddress,
		req.DelegatorAddress,
		req.UnstakeAmount,
		"unstake",
		"confirmed",
		time.Now(),
		txHash,
	)

	if err != nil {
		return nil, repoerrors.NewDatabaseError("unstake_usc", err)
	}

	return &proto.UnstakeUSCResponse{
		Success:         true,
		TransactionHash: txHash,
		Status:          0, // 0=Pending
		ErrorMessage:    "",
	}, nil
}

// saveStakingToDatabase saves staking transaction to database for analytics
func (r *Repository) saveStakingToDatabase(ctx context.Context, req interface{}, resp interface{}, txType string) error {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return fmt.Errorf("postgres connection not available")
	}

	var validatorAddress, delegatorAddress, amount, txHash, userId string

	switch reqType := req.(type) {
	case *proto.StakeUSCRequest:
		validatorAddress = reqType.ValidatorAddress
		delegatorAddress = reqType.DelegatorAddress
		amount = reqType.StakeAmount
		userId = reqType.UserId
		if respType, ok := resp.(*proto.StakeUSCResponse); ok && respType != nil {
			txHash = respType.TransactionHash
		}
	case *proto.UnstakeUSCRequest:
		validatorAddress = reqType.ValidatorAddress
		delegatorAddress = reqType.DelegatorAddress
		amount = reqType.UnstakeAmount
		userId = reqType.UserId
		if respType, ok := resp.(*proto.UnstakeUSCResponse); ok && respType != nil {
			txHash = respType.TransactionHash
		}
	default:
		return fmt.Errorf("unknown request type: %T", req)
	}

	if txHash == "" {
		return fmt.Errorf("transaction hash is empty")
	}

	now := time.Now()

	// Save to analytics table
	analyticsQuery := `
		INSERT INTO usc_staking_analytics (
			validator_address, staker_address, amount, transaction_type,
			status, timestamp, transaction_hash
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (transaction_hash) DO UPDATE SET
			amount = EXCLUDED.amount,
			status = EXCLUDED.status,
			timestamp = EXCLUDED.timestamp
	`

	if _, err := postgres.ExecContext(ctx, analyticsQuery,
		validatorAddress,
		delegatorAddress,
		amount,
		txType,
		"confirmed",
		now,
		txHash,
	); err != nil {
		return fmt.Errorf("failed to save to analytics table: %w", err)
	}

	// Save to main staking table
	mainQuery := `
		INSERT INTO staking (
			delegator_address, validator_address, stake_amount, stake_type,
			transaction_hash, user_id
		) VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (delegator_address, validator_address, transaction_hash) DO UPDATE SET
			stake_amount = EXCLUDED.stake_amount,
			stake_type = EXCLUDED.stake_type,
			updated_at = NOW()
	`

	if _, err := postgres.ExecContext(ctx, mainQuery,
		delegatorAddress,
		validatorAddress,
		amount,
		txType,
		txHash,
		userId,
	); err != nil {
		return fmt.Errorf("failed to save to main table: %w", err)
	}

	return nil
}

// Database fallback methods

// getValidatorsFromDatabase retrieves validators from database
func (r *Repository) getValidatorsFromDatabase(ctx context.Context, req *proto.GetValidatorsRequest) (*proto.GetValidatorsResponse, error) {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return &proto.GetValidatorsResponse{
			Validators: []*proto.ValidatorInfo{},
			TotalCount: 0,
			HasMore:    false,
		}, nil
	}

	limit, offset := utils.NormalizePagination(req.Limit, req.Offset, utils.PaginationConfig{
		DefaultLimit: 100,
		MaxLimit:     1000,
	})

	query := `
		SELECT validator_address, validator_name, commission_rate, status,
		       COUNT(*) OVER() as total_count
		FROM usc_validator_analytics
		ORDER BY registered_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := postgres.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return &proto.GetValidatorsResponse{
			Validators: []*proto.ValidatorInfo{},
			TotalCount: 0,
			HasMore:    false,
		}, nil
	}
	defer rows.Close()

	validators := make([]*proto.ValidatorInfo, 0, limit)
	var totalCount int32
	for rows.Next() {
		var v proto.ValidatorInfo
		if err := rows.Scan(&v.ValidatorAddress, &v.ValidatorName, &v.CommissionRate, &v.Status, &totalCount); err != nil {
			continue
		}
		v.ValidatorId = v.ValidatorAddress
		validators = append(validators, &v)
	}

	if err := rows.Err(); err != nil {
		if totalCount == 0 {
			totalCount = int32(len(validators))
		}
	}

	if totalCount == 0 {
		totalCount = int32(len(validators))
	}

	hasMore := int32(len(validators)) == limit && (offset+limit) < totalCount

	return &proto.GetValidatorsResponse{
		Validators: validators,
		TotalCount: totalCount,
		HasMore:    hasMore,
	}, nil
}

// getValidatorStatusFromDatabase retrieves validator status from database
func (r *Repository) getValidatorStatusFromDatabase(ctx context.Context, req *proto.GetValidatorStatusRequest) (*proto.GetValidatorStatusResponse, error) {
	postgres := repoerrors.GetPostgresConnection(r.db)
	if postgres == nil {
		return &proto.GetValidatorStatusResponse{
			Success:      false,
			ErrorMessage: "database not available",
		}, nil
	}

	query := `SELECT validator_address, validator_name, commission_rate, status FROM usc_validator_analytics WHERE validator_address = $1 LIMIT 1`

	var v proto.ValidatorInfo
	err := postgres.QueryRowContext(ctx, query, req.ValidatorAddress).Scan(
		&v.ValidatorAddress, &v.ValidatorName, &v.CommissionRate, &v.Status,
	)
	if err != nil {
		return &proto.GetValidatorStatusResponse{
			Success:      false,
			ErrorMessage: "validator not found",
		}, nil
	}

	v.ValidatorId = v.ValidatorAddress
	return &proto.GetValidatorStatusResponse{
		Success:   true,
		Validator: &v,
	}, nil
}
