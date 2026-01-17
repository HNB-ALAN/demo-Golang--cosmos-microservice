package types

import (
	"fmt"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesisState returns the default genesis state for the validator module
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Validators:  []Validator{},
		Delegations: []Delegation{},
		Params:      DefaultParams(),
	}
}

// ValidateGenesis validates the genesis state
func (gs GenesisState) ValidateGenesis() error {
	// Validate parameters
	if err := gs.Params.Validate(); err != nil {
		return fmt.Errorf("invalid parameters: %w", err)
	}

	// Validate validators
	seenAddresses := make(map[string]bool)
	for _, validator := range gs.Validators {
		if validator.Address == "" {
			return fmt.Errorf("validator address cannot be empty")
		}
		if validator.PubKey == "" {
			return fmt.Errorf("validator pub key cannot be empty")
		}
		if validator.Description == "" {
			return fmt.Errorf("validator description cannot be empty")
		}
		if validator.Commission == "" {
			return fmt.Errorf("validator commission cannot be empty")
		}

		// Check for duplicate addresses
		if seenAddresses[validator.Address] {
			return fmt.Errorf("duplicate validator address: %s", validator.Address)
		}
		seenAddresses[validator.Address] = true
	}

	// Validate delegations
	for _, delegation := range gs.Delegations {
		if delegation.DelegatorAddress == "" {
			return fmt.Errorf("delegation delegator address cannot be empty")
		}
		if delegation.ValidatorAddress == "" {
			return fmt.Errorf("delegation validator address cannot be empty")
		}
		if delegation.Amount == "" {
			return fmt.Errorf("delegation amount cannot be empty")
		}
		if delegation.CreatedAt <= 0 {
			return fmt.Errorf("delegation timestamp must be positive")
		}
	}

	return nil
}

// ExportGenesis exports the genesis state
func ExportGenesis(validators []Validator, delegations []Delegation, params Params) *GenesisState {
	// Sort validators by address for deterministic output
	sort.Slice(validators, func(i, j int) bool {
		return validators[i].Address < validators[j].Address
	})

	// Sort delegations by delegator address for deterministic output
	sort.Slice(delegations, func(i, j int) bool {
		return delegations[i].DelegatorAddress < delegations[j].DelegatorAddress
	})

	return &GenesisState{
		Validators:  validators,
		Delegations: delegations,
		Params:      params,
	}
}

// InitGenesis initializes the genesis state
func InitGenesis(ctx sdk.Context, keeper interface{}, gs GenesisState) error {
	// Validate genesis state
	if err := gs.ValidateGenesis(); err != nil {
		return fmt.Errorf("invalid genesis state: %w", err)
	}

	// TODO: Implement keeper operations when keeper interface is defined
	// Set parameters
	// if err := keeper.SetParams(ctx, gs.Params); err != nil {
	// 	return fmt.Errorf("failed to set parameters: %w", err)
	// }

	// Set validators
	// for _, validator := range gs.Validators {
	// 	if err := keeper.SetValidator(ctx, validator); err != nil {
	// 		return fmt.Errorf("failed to set validator for address %s: %w", validator.Address, err)
	// 	}
	// }

	// Set delegations
	// for _, delegation := range gs.Delegations {
	// 	if err := keeper.SetDelegation(ctx, delegation); err != nil {
	// 		return fmt.Errorf("failed to set delegation: %w", err)
	// 	}
	// }

	return nil
}

// GetGenesisState returns the current genesis state
func GetGenesisState(ctx sdk.Context, keeper interface{}) (*GenesisState, error) {
	// TODO: Implement keeper operations when keeper interface is defined
	// Get parameters
	// params, err := keeper.GetParams(ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get parameters: %w", err)
	// }

	// Get all validators
	// validators, err := keeper.GetAllValidators(ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get validators: %w", err)
	// }

	// Get all delegations
	// delegations, err := keeper.GetAllDelegations(ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get delegations: %w", err)
	// }

	// return ExportGenesis(validators, delegations, params), nil

	// Return default genesis state for now
	return DefaultGenesisState(), nil
}
