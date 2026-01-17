package types

import (
	"fmt"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesisState returns the default genesis state for the product_certificate module
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Certificates:  []ProductCertificate{},
		Verifications: []CertificateVerification{},
		Params:        DefaultParams(),
	}
}

// ValidateGenesis validates the genesis state
func (gs GenesisState) ValidateGenesis() error {
	// Validate parameters
	if err := gs.Params.Validate(); err != nil {
		return fmt.Errorf("invalid parameters: %w", err)
	}

	// Validate certificates
	seenIDs := make(map[string]bool)
	for _, cert := range gs.Certificates {
		if cert.ID == "" {
			return fmt.Errorf("certificate ID cannot be empty")
		}
		if cert.ProductID == "" {
			return fmt.Errorf("product ID cannot be empty")
		}
		if cert.Owner == "" {
			return fmt.Errorf("owner cannot be empty")
		}
		if cert.Status == "" {
			return fmt.Errorf("status cannot be empty")
		}

		// Check for duplicate IDs
		if seenIDs[cert.ID] {
			return fmt.Errorf("duplicate certificate ID: %s", cert.ID)
		}
		seenIDs[cert.ID] = true
	}

	// Validate verifications
	for _, verification := range gs.Verifications {
		if verification.CertificateID == "" {
			return fmt.Errorf("verification certificate ID cannot be empty")
		}
		if verification.Verifier == "" {
			return fmt.Errorf("verification verifier cannot be empty")
		}
		if verification.Status == "" {
			return fmt.Errorf("verification status cannot be empty")
		}
		if verification.VerifiedAt <= 0 {
			return fmt.Errorf("verification timestamp must be positive")
		}
	}

	return nil
}

// ExportGenesis exports the genesis state
func ExportGenesis(certificates []ProductCertificate, verifications []CertificateVerification, params Params) *GenesisState {
	// Sort certificates by ID for deterministic output
	sort.Slice(certificates, func(i, j int) bool {
		return certificates[i].ID < certificates[j].ID
	})

	// Sort verifications by certificate ID for deterministic output
	sort.Slice(verifications, func(i, j int) bool {
		return verifications[i].CertificateID < verifications[j].CertificateID
	})

	return &GenesisState{
		Certificates:  certificates,
		Verifications: verifications,
		Params:        params,
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

	// Set certificates
	// for _, cert := range gs.Certificates {
	// 	if err := keeper.SetCertificate(ctx, cert); err != nil {
	// 		return fmt.Errorf("failed to set certificate for ID %s: %w", cert.ID, err)
	// 	}
	// }

	// Set verifications
	// for _, verification := range gs.Verifications {
	// 	if err := keeper.SetVerification(ctx, verification); err != nil {
	// 		return fmt.Errorf("failed to set verification: %w", err)
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

	// Get all certificates
	// certificates, err := keeper.GetAllCertificates(ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get certificates: %w", err)
	// }

	// Get all verifications
	// verifications, err := keeper.GetAllVerifications(ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get verifications: %w", err)
	// }

	// return ExportGenesis(certificates, verifications, params), nil

	// Return default genesis state for now
	return DefaultGenesisState(), nil
}
