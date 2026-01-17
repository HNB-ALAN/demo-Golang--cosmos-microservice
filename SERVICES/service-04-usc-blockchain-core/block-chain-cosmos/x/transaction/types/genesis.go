package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
)

// DefaultGenesisState returns the default genesis state
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Transactions:     []Transaction{},
		TransactionStats: NewTransactionStats(),
		Params:           DefaultParams,
	}
}

// ValidateGenesis validates the genesis state
func ValidateGenesis(gs *GenesisState) error {
	// Validate transactions
	for _, tx := range gs.Transactions {
		if tx.Hash == "" {
			return fmt.Errorf("transaction hash cannot be empty")
		}
		if tx.FromAddress == "" {
			return fmt.Errorf("transaction from address cannot be empty")
		}
		if tx.ToAddress == "" {
			return fmt.Errorf("transaction to address cannot be empty")
		}
		if tx.Amount == "" {
			return fmt.Errorf("transaction amount cannot be empty")
		}
		if err := ValidateTransactionType(tx.TransactionType); err != nil {
			return err
		}
		if err := ValidateTransactionStatus(tx.Status); err != nil {
			return err
		}
	}

	// Validate parameters
	if err := gs.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}

	return nil
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx interface{}, k interface{}) *GenesisState {
	// This will be implemented by the keeper
	return DefaultGenesis()
}

// InitGenesis initializes the genesis state
func InitGenesis(ctx interface{}, k interface{}, genState GenesisState) {
	// This will be implemented by the keeper
}

// GetGenesisState returns the genesis state
func GetGenesisState(cdc codec.Codec, data []byte) (GenesisState, error) {
	var genState GenesisState
	if err := cdc.UnmarshalJSON(data, &genState); err != nil {
		return GenesisState{}, err
	}
	return genState, nil
}
