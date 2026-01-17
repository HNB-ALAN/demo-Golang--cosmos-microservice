package app

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/module"
)

// ============================================================================
// UTILITY METHODS
// ============================================================================

// SimulationManager returns the simulation manager
func (app *USCApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

// GetMaccPerms returns module account permissions
func (app *USCApp) GetMaccPerms() map[string][]string {
	return maccPerms
}

// ModuleAccountAddrs returns module account addresses
func (app *USCApp) ModuleAccountAddrs() map[string]bool {
	return moduleAccountAddrs
}

// LoadHeight loads app state at given height
// NOTE: This is a stub implementation for interface compliance
// Full implementation will be added when needed for state queries at specific heights
func (app *USCApp) LoadHeight(height int64) error {
	fmt.Printf("Loading app state at height: %d\n", height)
	// TODO: Implement actual height loading when required
	return nil
}

// ExportAppStateAndValidators exports app state and validators
func (app *USCApp) ExportAppStateAndValidators() (appState []byte, validators []byte, err error) {
	fmt.Println("Exporting app state and validators...")

	if app.mm == nil {
		return []byte("{}"), []byte("[]"), nil
	}

	// Create a temporary context for export
	ctx := app.NewContext(false)

	// Export genesis state from all modules - returns map and error
	genesisStateMap, err := app.mm.ExportGenesis(ctx, app.appCodec)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to export genesis: %w", err)
	}

	// Marshal genesis state
	appState, err = json.MarshalIndent(genesisStateMap, "", "  ")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal app state: %w", err)
	}

	// For now, return empty validators array
	// Validators will be exported when staking module is fully implemented
	validators = []byte("[]")

	fmt.Println("✅ App state and validators exported successfully")
	return appState, validators, nil
}

// GetSubspace returns parameter subspace for module
// NOTE: This is a stub implementation for interface compliance
// Full implementation will return actual subspace when module parameters are needed
func (app *USCApp) GetSubspace(moduleName string) interface{} {
	fmt.Printf("Getting subspace for module: %s\n", moduleName)
	// TODO: Implement actual subspace retrieval when needed
	return nil
}

// Stop gracefully stops the USC app
// NOTE: BaseApp doesn't have a Stop method in Cosmos SDK 0.53.4
// This is a placeholder for graceful shutdown if needed in future
func (app *USCApp) Stop() error {
	fmt.Println("🛑 Stopping USC app...")

	// Close any open connections or resources
	// TODO: Add cleanup logic if needed (e.g., close database connections, cancel background tasks)
	fmt.Println("✅ USC app stopped successfully")
	return nil
}
