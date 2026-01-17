package keeper

import (
	"context"
	"fmt"
	"strconv"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	query "github.com/cosmos/cosmos-sdk/types/query"
	blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/validator/v1/usc/validator/v1"
)

// QueryServer defines the query server interface using blockchain-proto types
type QueryServer interface {
	QueryValidator(context.Context, *blockchainproto.QueryValidatorRequest) (*blockchainproto.QueryValidatorResponse, error)
	QueryValidators(context.Context, *blockchainproto.QueryValidatorsRequest) (*blockchainproto.QueryValidatorsResponse, error)
	QueryValidatorDelegations(context.Context, *blockchainproto.QueryValidatorDelegationsRequest) (*blockchainproto.QueryValidatorDelegationsResponse, error)
	QueryValidatorStats(context.Context, *blockchainproto.QueryValidatorStatsRequest) (*blockchainproto.QueryValidatorStatsResponse, error)
}

// queryServer implements the QueryServer interface
type queryServer struct {
	Keeper
}

// NewQueryServerImpl returns an implementation of the QueryServer interface
func NewQueryServerImpl(keeper Keeper) QueryServer {
	return &queryServer{Keeper: keeper}
}

// QueryValidator returns a specific validator
func (k queryServer) QueryValidator(ctx context.Context, req *blockchainproto.QueryValidatorRequest) (*blockchainproto.QueryValidatorResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("empty request")
	}

	if req.ValidatorAddress == "" {
		return nil, fmt.Errorf("validator address cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	validator, err := k.GetValidator(sdkCtx, req.ValidatorAddress)
	if err != nil {
		return nil, fmt.Errorf("validator not found: %w", err)
	}

	// Convert string types to blockchain-proto enums
	var validatorStatus blockchainproto.ValidatorStatus
	switch validator.Status {
	case "active":
		validatorStatus = blockchainproto.ValidatorStatus_VALIDATOR_STATUS_ACTIVE
	case "inactive":
		validatorStatus = blockchainproto.ValidatorStatus_VALIDATOR_STATUS_INACTIVE
	case "jailed":
		validatorStatus = blockchainproto.ValidatorStatus_VALIDATOR_STATUS_JAILED
	default:
		validatorStatus = blockchainproto.ValidatorStatus_VALIDATOR_STATUS_UNSPECIFIED
	}

	// Parse commission rate
	commissionRate, _ := strconv.ParseFloat(validator.Commission, 64)

	// Convert to blockchain-proto Validator type
	blockchainValidator := &blockchainproto.Validator{
		Address:         validator.Address,
		Name:            validator.Address, // Using address as name for now
		Description:     &blockchainproto.ValidatorDescription{Details: validator.Description},
		Status:          validatorStatus,
		CommissionRate:  commissionRate,
		TotalStake:      &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)}, // Placeholder
		SelfStake:       &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)}, // Placeholder
		DelegatedStake:  &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)}, // Placeholder
		DelegationCount: 0,                                               // Placeholder
		CreatedAt:       nil,                                             // Placeholder
		UpdatedAt:       nil,                                             // Placeholder
		JailedAt:        nil,
		UnjailedAt:      nil,
		Memo:            "",
	}

	return &blockchainproto.QueryValidatorResponse{
		Validator: blockchainValidator,
	}, nil
}

// QueryValidators returns all validators with pagination
func (k queryServer) QueryValidators(ctx context.Context, req *blockchainproto.QueryValidatorsRequest) (*blockchainproto.QueryValidatorsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get all validators
	validators, err := k.GetAllValidators(sdkCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to get validators: %w", err)
	}

	// Convert to blockchain-proto Validator types
	var blockchainValidators []*blockchainproto.Validator
	for _, validator := range validators {
		// Convert string types to blockchain-proto enums
		var validatorStatus blockchainproto.ValidatorStatus
		switch validator.Status {
		case "active":
			validatorStatus = blockchainproto.ValidatorStatus_VALIDATOR_STATUS_ACTIVE
		case "inactive":
			validatorStatus = blockchainproto.ValidatorStatus_VALIDATOR_STATUS_INACTIVE
		case "jailed":
			validatorStatus = blockchainproto.ValidatorStatus_VALIDATOR_STATUS_JAILED
		default:
			validatorStatus = blockchainproto.ValidatorStatus_VALIDATOR_STATUS_UNSPECIFIED
		}

		// Parse commission rate
		commissionRate, _ := strconv.ParseFloat(validator.Commission, 64)

		blockchainValidator := &blockchainproto.Validator{
			Address:         validator.Address,
			Name:            validator.Address, // Using address as name for now
			Description:     &blockchainproto.ValidatorDescription{Details: validator.Description},
			Status:          validatorStatus,
			CommissionRate:  commissionRate,
			TotalStake:      &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)}, // Placeholder
			SelfStake:       &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)}, // Placeholder
			DelegatedStake:  &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)}, // Placeholder
			DelegationCount: 0,                                               // Placeholder
			CreatedAt:       nil,                                             // Placeholder
			UpdatedAt:       nil,                                             // Placeholder
			JailedAt:        nil,
			UnjailedAt:      nil,
			Memo:            "",
		}
		blockchainValidators = append(blockchainValidators, blockchainValidator)
	}

	// Apply pagination
	pageRes := &query.PageResponse{
		NextKey: nil,
		Total:   uint64(len(blockchainValidators)),
	}

	return &blockchainproto.QueryValidatorsResponse{
		Validators: blockchainValidators,
		Pagination: pageRes,
	}, nil
}

// QueryValidatorDelegations returns delegations for a specific validator
func (k queryServer) QueryValidatorDelegations(ctx context.Context, req *blockchainproto.QueryValidatorDelegationsRequest) (*blockchainproto.QueryValidatorDelegationsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("empty request")
	}

	if req.ValidatorAddress == "" {
		return nil, fmt.Errorf("validator address cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get all delegations for the validator
	delegations, err := k.GetAllDelegations(sdkCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to get delegations: %w", err)
	}

	// Filter delegations for the specific validator
	var validatorDelegations []*blockchainproto.Delegation
	for _, delegation := range delegations {
		if delegation.ValidatorAddress == req.ValidatorAddress {
			// Parse amount string to Coin
			amount, err := sdk.ParseCoinNormalized(delegation.Amount)
			if err != nil {
				amount = sdk.NewCoin("usc", math.NewInt(0))
			}

			blockchainDelegation := &blockchainproto.Delegation{
				Id:               "delegation_" + delegation.DelegatorAddress + "_" + delegation.ValidatorAddress,
				Delegator:        delegation.DelegatorAddress,
				ValidatorAddress: delegation.ValidatorAddress,
				DelegationAmount: &amount,
				DelegationTime:   nil, // Placeholder
				Status:           blockchainproto.DelegationStatus_DELEGATION_STATUS_ACTIVE,
				UndelegationTime: nil,
				RewardsEarned:    &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)},
				Memo:             "",
			}
			validatorDelegations = append(validatorDelegations, blockchainDelegation)
		}
	}

	// Apply pagination
	pageRes := &query.PageResponse{
		NextKey: nil,
		Total:   uint64(len(validatorDelegations)),
	}

	return &blockchainproto.QueryValidatorDelegationsResponse{
		Delegations: validatorDelegations,
		Pagination:  pageRes,
	}, nil
}

// QueryValidatorStats returns statistics for a specific validator
func (k queryServer) QueryValidatorStats(ctx context.Context, req *blockchainproto.QueryValidatorStatsRequest) (*blockchainproto.QueryValidatorStatsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("empty request")
	}

	if req.ValidatorAddress == "" {
		return nil, fmt.Errorf("validator address cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get validator
	validator, err := k.GetValidator(sdkCtx, req.ValidatorAddress)
	if err != nil {
		return nil, fmt.Errorf("validator not found: %w", err)
	}

	// Get delegations for the validator
	delegations, err := k.GetAllDelegations(sdkCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to get delegations: %w", err)
	}

	// Calculate stats
	var totalDelegations int64
	var totalAmount int64
	for _, delegation := range delegations {
		if delegation.ValidatorAddress == req.ValidatorAddress {
			totalDelegations++
			// Parse amount (assuming it's a string representation of int64)
			// totalAmount += parseAmount(delegation.Amount)
		}
	}

	// Parse commission rate
	commissionRate, _ := strconv.ParseFloat(validator.Commission, 64)

	// Create stats response
	stats := &blockchainproto.ValidatorStats{
		TotalValidators:        1, // This validator only
		ActiveValidators:       1, // Assuming active
		InactiveValidators:     0,
		JailedValidators:       0,
		TotalStake:             &sdk.Coin{Denom: "usc", Amount: math.NewInt(totalAmount)},
		AverageCommission:      commissionRate,
		TotalDelegations:       &sdk.Coin{Denom: "usc", Amount: math.NewInt(totalDelegations)},
		AverageDelegation:      &sdk.Coin{Denom: "usc", Amount: math.NewInt(0)},
		MostDelegatedValidator: validator.Address,
	}

	return &blockchainproto.QueryValidatorStatsResponse{
		Stats: stats,
	}, nil
}
