package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesisState returns the default genesis state for the NFT module
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		NFTs:        []NFT{},
		Collections: []Collection{},
		Params:      DefaultParams(),
	}
}

// InitGenesis initializes the nft module's genesis state
func InitGenesis(ctx sdk.Context, k interface{}, genState GenesisState) {
	// TODO: Implement genesis initialization
}

// ExportGenesis returns the nft module's exported genesis state
func ExportGenesis(ctx sdk.Context, k interface{}) GenesisState {
	return GenesisState{
		NFTs:        []NFT{},
		Collections: []Collection{},
		Params:      DefaultParams(),
	}
}

// ValidateGenesis validates the nft module's genesis state
func ValidateGenesis(genState GenesisState) error {
	if err := genState.Params.Validate(); err != nil {
		return fmt.Errorf("invalid parameters: %w", err)
	}

	for _, nft := range genState.NFTs {
		if err := nft.Validate(); err != nil {
			return fmt.Errorf("invalid NFT: %w", err)
		}
	}

	for _, collection := range genState.Collections {
		if err := collection.Validate(); err != nil {
			return fmt.Errorf("invalid collection: %w", err)
		}
	}

	return nil
}
