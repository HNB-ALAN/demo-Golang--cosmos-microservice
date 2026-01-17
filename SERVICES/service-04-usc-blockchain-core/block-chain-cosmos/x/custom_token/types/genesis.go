package types

import (
	"fmt"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesisState returns the default genesis state for the custom_token module
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Tokens:    []CustomToken{},
		Balances:  []TokenBalance{},
		Transfers: []TokenTransfer{},
		Params:    DefaultParams(),
	}
}

// ValidateGenesis validates the genesis state
func (gs GenesisState) ValidateGenesis() error {
	// Validate parameters
	if err := gs.Params.Validate(); err != nil {
		return fmt.Errorf("invalid parameters: %w", err)
	}

	// Validate tokens
	seenIDs := make(map[string]bool)
	for _, token := range gs.Tokens {
		if token.ID == "" {
			return fmt.Errorf("token ID cannot be empty")
		}
		if token.Name == "" {
			return fmt.Errorf("token name cannot be empty")
		}
		if token.Symbol == "" {
			return fmt.Errorf("token symbol cannot be empty")
		}
		if token.Owner == "" {
			return fmt.Errorf("owner cannot be empty")
		}
		if token.Status == "" {
			return fmt.Errorf("status cannot be empty")
		}
		if token.Decimals > 18 {
			return fmt.Errorf("decimals cannot exceed 18")
		}

		// Check for duplicate IDs
		if seenIDs[token.ID] {
			return fmt.Errorf("duplicate token ID: %s", token.ID)
		}
		seenIDs[token.ID] = true
	}

	// Validate balances
	for _, balance := range gs.Balances {
		if balance.TokenID == "" {
			return fmt.Errorf("balance token ID cannot be empty")
		}
		if balance.Owner == "" {
			return fmt.Errorf("balance owner cannot be empty")
		}
		if balance.Amount == "" {
			return fmt.Errorf("balance amount cannot be empty")
		}
		if balance.UpdatedAt <= 0 {
			return fmt.Errorf("balance timestamp must be positive")
		}
	}

	// Validate transfers
	for _, transfer := range gs.Transfers {
		if transfer.ID == "" {
			return fmt.Errorf("transfer ID cannot be empty")
		}
		if transfer.TokenID == "" {
			return fmt.Errorf("transfer token ID cannot be empty")
		}
		if transfer.From == "" {
			return fmt.Errorf("transfer from cannot be empty")
		}
		if transfer.To == "" {
			return fmt.Errorf("transfer to cannot be empty")
		}
		if transfer.Amount == "" {
			return fmt.Errorf("transfer amount cannot be empty")
		}
		if transfer.CreatedAt <= 0 {
			return fmt.Errorf("transfer timestamp must be positive")
		}
	}

	return nil
}

// ExportGenesis exports the genesis state
func ExportGenesis(tokens []CustomToken, balances []TokenBalance, transfers []TokenTransfer, params Params) *GenesisState {
	// Sort tokens by ID for deterministic output
	sort.Slice(tokens, func(i, j int) bool {
		return tokens[i].ID < tokens[j].ID
	})

	// Sort balances by token ID for deterministic output
	sort.Slice(balances, func(i, j int) bool {
		return balances[i].TokenID < balances[j].TokenID
	})

	// Sort transfers by ID for deterministic output
	sort.Slice(transfers, func(i, j int) bool {
		return transfers[i].ID < transfers[j].ID
	})

	return &GenesisState{
		Tokens:    tokens,
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

	// Set tokens
	// for _, token := range gs.Tokens {
	// 	if err := keeper.SetToken(ctx, token); err != nil {
	// 		return fmt.Errorf("failed to set token for ID %s: %w", token.ID, err)
	// 	}
	// }

	// Set balances
	// for _, balance := range gs.Balances {
	// 	if err := keeper.SetBalance(ctx, balance); err != nil {
	// 		return fmt.Errorf("failed to set balance: %w", err)
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

	// Get all tokens
	// tokens, err := keeper.GetAllTokens(ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get tokens: %w", err)
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

	// return ExportGenesis(tokens, balances, transfers, params), nil

	// Return default genesis state for now
	return DefaultGenesisState(), nil
}
