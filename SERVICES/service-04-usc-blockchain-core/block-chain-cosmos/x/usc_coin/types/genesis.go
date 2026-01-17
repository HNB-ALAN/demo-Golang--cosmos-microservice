package types

import (
	"fmt"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesisState returns the default genesis state for the USC module
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Balances:  []Balance{},
		Transfers: []Transfer{},
		Params:    DefaultParams(),
	}
}

// ValidateGenesis validates the genesis state
func (gs GenesisState) ValidateGenesis() error {
	// Validate parameters
	if err := gs.Params.Validate(); err != nil {
		return fmt.Errorf("invalid parameters: %w", err)
	}

	// Validate balances
	seenAddresses := make(map[string]bool)
	for _, balance := range gs.Balances {
		if balance.Address == "" {
			return fmt.Errorf("balance address cannot be empty")
		}
		if balance.Amount == "" {
			return fmt.Errorf("balance amount cannot be empty")
		}
		if balance.Denom == "" {
			return fmt.Errorf("balance denomination cannot be empty")
		}

		// Check for duplicate addresses
		if seenAddresses[balance.Address] {
			return fmt.Errorf("duplicate balance address: %s", balance.Address)
		}
		seenAddresses[balance.Address] = true
	}

	// Validate transfers
	for _, transfer := range gs.Transfers {
		if transfer.FromAddress == "" {
			return fmt.Errorf("transfer from address cannot be empty")
		}
		if transfer.ToAddress == "" {
			return fmt.Errorf("transfer to address cannot be empty")
		}
		if transfer.Amount == "" {
			return fmt.Errorf("transfer amount cannot be empty")
		}
		if transfer.Denom == "" {
			return fmt.Errorf("transfer denomination cannot be empty")
		}
		if transfer.Timestamp <= 0 {
			return fmt.Errorf("transfer timestamp must be positive")
		}
	}

	return nil
}

// ExportGenesis exports the genesis state
func ExportGenesis(balances []Balance, transfers []Transfer, params Params) *GenesisState {
	// Sort balances by address for deterministic output
	sort.Slice(balances, func(i, j int) bool {
		return balances[i].Address < balances[j].Address
	})

	// Sort transfers by timestamp for deterministic output
	sort.Slice(transfers, func(i, j int) bool {
		return transfers[i].Timestamp < transfers[j].Timestamp
	})

	return &GenesisState{
		Balances:  balances,
		Transfers: transfers,
		Params:    params,
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

	// Set balances
	// for _, balance := range gs.Balances {
	// 	if err := keeper.SetBalance(ctx, balance.Address, balance); err != nil {
	// 		return fmt.Errorf("failed to set balance for address %s: %w", balance.Address, err)
	// 	}
	// }

	// Set transfers
	// for _, transfer := range gs.Transfers {
	// 	if err := keeper.SetTransfer(ctx, transfer); err != nil {
	// 		return fmt.Errorf("failed to set transfer: %w", err)
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

	// Get all balances
	// balances, err := keeper.GetAllBalances(ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get balances: %w", err)
	// }

	// Get all transfers
	// transfers, err := keeper.GetAllTransfers(ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get transfers: %w", err)
	// }

	// return ExportGenesis(balances, transfers, params), nil

	// Return default genesis state for now
	return DefaultGenesisState(), nil
}
