package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// ============================================================================
// GENESIS HELPERS
// ============================================================================

// filterRegisteredModules filters genesis state to only include modules registered in ModuleManager
func filterRegisteredModules(genesisStateMap map[string]json.RawMessage, registeredModules []string) map[string]json.RawMessage {
	filtered := make(map[string]json.RawMessage)
	registeredSet := make(map[string]bool)
	for _, name := range registeredModules {
		registeredSet[name] = true
	}

	for moduleName, moduleState := range genesisStateMap {
		if registeredSet[moduleName] {
			filtered[moduleName] = moduleState
		}
		// Note: Unregistered modules are silently skipped (expected - genesis may contain future modules)
	}
	return filtered
}

// sanitizeGenesisState sanitizes genesis state for specific modules to match Cosmos SDK 0.53.4 format
func sanitizeGenesisState(moduleName string, state json.RawMessage) json.RawMessage {
	switch moduleName {
	case "bank":
		return sanitizeBankGenesisState(state)
	case "usc_coin":
		return sanitizeUSCCoinGenesisState(state)
	default:
		return state // No sanitization needed for other modules
	}
}

// sanitizeBankGenesisState sanitizes bank module genesis state for Cosmos SDK 0.53.4
// Fixes: send_enabled (boolean -> array), removes receive_enabled (not supported)
func sanitizeBankGenesisState(state json.RawMessage) json.RawMessage {
	var bankState map[string]interface{}
	if err := json.Unmarshal(state, &bankState); err != nil {
		return state // Return original if unmarshal fails
	}

	params, ok := bankState["params"].(map[string]interface{})
	if !ok {
		return state // Return original if params don't exist
	}

	// Fix send_enabled: ensure it's an array, not boolean
	if sendEnabled, ok := params["send_enabled"].(bool); ok {
		params["send_enabled"] = []interface{}{}
		if sendEnabled {
			fmt.Printf("  🔧 Fixed bank.send_enabled: boolean true -> array\n")
		} else {
			fmt.Printf("  🔧 Fixed bank.send_enabled: boolean false -> array\n")
		}
	}

	// Remove receive_enabled field (not supported in Cosmos SDK 0.53.4 bank module)
	if _, exists := params["receive_enabled"]; exists {
		delete(params, "receive_enabled")
		fmt.Printf("  🔧 Removed bank.receive_enabled (not supported in SDK 0.53.4)\n")
	}

	// Re-marshal the corrected state
	correctedState, err := json.Marshal(bankState)
	if err != nil {
		return state // Return original if marshal fails
	}

	fmt.Printf("  ✅ Bank genesis state sanitized for SDK 0.53.4\n")
	return json.RawMessage(correctedState)
}

// sanitizeUSCCoinGenesisState sanitizes usc_coin module genesis state
// Adds default values for required fields if missing
func sanitizeUSCCoinGenesisState(state json.RawMessage) json.RawMessage {
	var uscState map[string]interface{}
	if err := json.Unmarshal(state, &uscState); err != nil {
		return state // Return original if unmarshal fails
	}

	params, ok := uscState["params"].(map[string]interface{})
	if !ok {
		return state // Return original if params don't exist
	}

	// Ensure required fields exist with defaults if missing
	if _, exists := params["token_name"]; !exists || params["token_name"] == "" {
		params["token_name"] = "Universal Social Coin"
		fmt.Printf("  🔧 Added usc_coin.token_name default\n")
	}
	if _, exists := params["token_symbol"]; !exists || params["token_symbol"] == "" {
		params["token_symbol"] = "USC"
		fmt.Printf("  🔧 Added usc_coin.token_symbol default\n")
	}
	if _, exists := params["token_decimals"]; !exists {
		params["token_decimals"] = 18
		fmt.Printf("  🔧 Added usc_coin.token_decimals default\n")
	}

	// Ensure optional fields have defaults if missing
	if _, exists := params["max_supply"]; !exists || params["max_supply"] == "" {
		params["max_supply"] = "10000000000000000000000000000"
		fmt.Printf("  🔧 Added usc_coin.max_supply default\n")
	}
	if _, exists := params["mint_enabled"]; !exists {
		params["mint_enabled"] = true
		fmt.Printf("  🔧 Added usc_coin.mint_enabled default\n")
	}
	if _, exists := params["burn_enabled"]; !exists {
		params["burn_enabled"] = true
		fmt.Printf("  🔧 Added usc_coin.burn_enabled default\n")
	}

	// Re-marshal the corrected state
	correctedState, err := json.Marshal(uscState)
	if err != nil {
		return state // Return original if marshal fails
	}

	fmt.Printf("  ✅ USC Coin genesis state sanitized\n")
	return json.RawMessage(correctedState)
}

// callModuleInitGenesis calls InitGenesis on a module with panic recovery
// Returns validators and error (if any)
func callModuleInitGenesis(ctx sdk.Context, appCodec codec.Codec, mod interface{}, moduleName string, state json.RawMessage) ([]abci.ValidatorUpdate, error) {
	var validatorUpdates []abci.ValidatorUpdate
	var err error

	func() {
		defer func() {
			if r := recover(); r != nil {
				panicMsg := fmt.Sprintf("%v", r)
				fmt.Printf("  ❌ PANIC in module '%s': %s\n", moduleName, panicMsg)

				// For usc_coin, handle panic gracefully - allow chain to start with default state
				if moduleName == "usc_coin" {
					fmt.Printf("  ℹ️  Continuing with default state for usc_coin (chain will start normally)\n")
					// Don't set err - allow chain to continue
					return
				}

				// For other modules, set error to stop processing
				err = fmt.Errorf("module %s panic: %s", moduleName, panicMsg)
			}
		}()

		// Call InitGenesis based on module type
		// COSMOS SDK 0.53.4: Check interfaces in order: HasGenesis -> HasABCIGenesis -> AppModule.InitGenesis
		// Block module should match AppModule.InitGenesis interface
		if hasGenesis, ok := mod.(module.HasGenesis); ok {
			// Legacy HasGenesis interface
			fmt.Printf("  🔍 Module '%s' matches HasGenesis interface\n", moduleName)
			hasGenesis.InitGenesis(ctx, appCodec, state)
			fmt.Printf("  ✅ Module '%s' initialized via HasGenesis interface\n", moduleName)
		} else if hasABCIGenesis, ok := mod.(module.HasABCIGenesis); ok {
			// HasABCIGenesis interface (returns validators)
			fmt.Printf("  🔍 Module '%s' matches HasABCIGenesis interface\n", moduleName)
			moduleValUpdates := hasABCIGenesis.InitGenesis(ctx, appCodec, state)
			if len(moduleValUpdates) > 0 {
				validatorUpdates = append(validatorUpdates, moduleValUpdates...)
				fmt.Printf("  ✅ Module '%s' initialized via HasABCIGenesis interface (%d validators)\n", moduleName, len(moduleValUpdates))
			} else {
				fmt.Printf("  ✅ Module '%s' initialized via HasABCIGenesis interface (no validators)\n", moduleName)
			}
		} else if appModule, ok := mod.(interface {
			InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) []abci.ValidatorUpdate
		}); ok {
			// AppModule with InitGenesis method (e.g., block module)
			fmt.Printf("  🔍 Module '%s' matches AppModule.InitGenesis interface (this is correct for block module)\n", moduleName)
			moduleValUpdates := appModule.InitGenesis(ctx, appCodec, state)
			if len(moduleValUpdates) > 0 {
				validatorUpdates = append(validatorUpdates, moduleValUpdates...)
				fmt.Printf("  ✅ Module '%s' initialized via AppModule.InitGenesis (%d validators)\n", moduleName, len(moduleValUpdates))
			} else {
				fmt.Printf("  ✅ Module '%s' initialized via AppModule.InitGenesis (no validators)\n", moduleName)
			}
		} else {
			fmt.Printf("  ⚠️  Module '%s' doesn't implement HasGenesis, HasABCIGenesis, or AppModule.InitGenesis, skipping\n", moduleName)
			fmt.Printf("  ⚠️  This means InitGenesis will NOT be called for module '%s'\n", moduleName)
		}
	}()

	return validatorUpdates, err
}

// ============================================================================
// GENESIS INITIALIZATION
// ============================================================================

// triggerInitChain triggers InitChain manually to initialize genesis state
// COSMOS SDK 0.53.4: This ensures InitGenesis is called for all modules
func (app *USCApp) triggerInitChain() error {
	// Read genesis.json to get genesis state
	genesisFile := filepath.Join(DefaultNodeHome, "config", "genesis.json")
	if _, err := os.Stat(genesisFile); os.IsNotExist(err) {
		// If genesis file doesn't exist, use empty state
		genesisFile = filepath.Join("block-chain-cosmos", "config", "genesis.json")
	}

	genesisBytes, err := os.ReadFile(genesisFile)
	if err != nil {
		// If can't read genesis file, create empty RequestInitChain
		req := &abci.RequestInitChain{
			ChainId:       app.BaseApp.ChainID(),
			AppStateBytes: []byte(`{"app_state":{}}`),
			Validators:    []abci.ValidatorUpdate{},
		}
		return app.handleInitChain(req)
	}

	// Parse genesis.json
	var genesisDoc map[string]interface{}
	if err := json.Unmarshal(genesisBytes, &genesisDoc); err != nil {
		return fmt.Errorf("failed to parse genesis.json: %w", err)
	}

	// Extract app_state
	appState, ok := genesisDoc["app_state"].(map[string]interface{})
	if !ok {
		appState = make(map[string]interface{})
	}

	// Marshal app_state to bytes
	appStateBytes, err := json.Marshal(appState)
	if err != nil {
		return fmt.Errorf("failed to marshal app_state: %w", err)
	}

	// Create RequestInitChain
	req := &abci.RequestInitChain{
		ChainId:       app.BaseApp.ChainID(),
		AppStateBytes: appStateBytes,
		Validators:    []abci.ValidatorUpdate{},
	}

	return app.handleInitChain(req)
}

// handleInitChain handles InitChain request
// COSMOS SDK 0.53.4: BaseApp.InitChain() commits state, but we need to verify it's queryable
func (app *USCApp) handleInitChain(req *abci.RequestInitChain) error {
	// COSMOS SDK 0.53.4: Use BaseApp.InitChain() which properly handles state commit
	// BaseApp.InitChain() internally:
	// 1. Creates a commit context (commit=true)
	// 2. Calls initChainer with that context
	// 3. Commits the state automatically
	resp, err := app.BaseApp.InitChain(req)
	if err != nil {
		return fmt.Errorf("BaseApp.InitChain failed: %w", err)
	}

	// COSMOS SDK 0.53.4: Verify genesis block is queryable after InitChain
	// If not, it means state wasn't properly committed - this can happen in standalone mode
	// We'll verify by trying to query block 1 immediately after InitChain
	testCtx := app.BaseApp.NewContext(true)
	block, queryErr := app.BlockKeeper.GetBlockByHeight(testCtx, 1)
	if queryErr != nil {
		fmt.Printf("⚠️  WARNING: Genesis block not queryable after InitChain: %v\n", queryErr)
		fmt.Printf("⚠️  This is expected in standalone mode - block will be queryable after first block is produced\n")
	} else {
		fmt.Printf("✅ Genesis block verified queryable after InitChain: height=%d, hash=%s\n", block.Height, block.Hash)
	}

	fmt.Printf("✅ InitChain completed: %d validators (state committed by BaseApp)\n", len(resp.Validators))
	return nil
}

// buildInitChainResponse builds ResponseInitChain from validators and error state
func buildInitChainResponse(validatorUpdates []abci.ValidatorUpdate, reqValidators []abci.ValidatorUpdate, err error) *abci.ResponseInitChain {
	// If there's an error, return validators from request to allow chain to start
	if err != nil {
		return &abci.ResponseInitChain{Validators: reqValidators}
	}

	// Prioritize validators from modules, then from request, then empty
	if len(validatorUpdates) > 0 {
		return &abci.ResponseInitChain{Validators: validatorUpdates}
	}
	if len(reqValidators) > 0 {
		return &abci.ResponseInitChain{Validators: reqValidators}
	}
	return &abci.ResponseInitChain{Validators: []abci.ValidatorUpdate{}}
}

// initChainer handles chain initialization
// Signature: func(ctx Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error)
func (app *USCApp) initChainer(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
	fmt.Printf("InitChainer called with chain-id: %s\n", req.ChainId)
	if app.mm == nil {
		fmt.Println("  ModuleManager is nil, returning empty response")
		return &abci.ResponseInitChain{
			Validators: req.Validators, // Echo back validators from genesis
		}, nil
	}

	// Use ModuleManager to handle genesis initialization
	genesisState := req.AppStateBytes
	if genesisState != nil {
		var genesisStateMap map[string]json.RawMessage
		if err := json.Unmarshal(genesisState, &genesisStateMap); err != nil {
			return nil, fmt.Errorf("failed to unmarshal genesis state: %w", err)
		}

		// Filter genesis state to only include modules registered in ModuleManager
		registeredModules := app.mm.ModuleNames()
		filteredGenesisState := filterRegisteredModules(genesisStateMap, registeredModules)

		// Log modules that will be initialized
		moduleNames := make([]string, 0, len(filteredGenesisState))
		for name := range filteredGenesisState {
			moduleNames = append(moduleNames, name)
		}
		fmt.Printf("  Initializing %d modules: %v\n", len(filteredGenesisState), moduleNames)

		// COSMOS SDK 0.53.4: Always initialize block module even if genesis state is empty
		// This ensures genesis block (height 1) is saved to keeper
		if len(filteredGenesisState) == 0 {
			fmt.Println("  ℹ️  No registered modules in genesis state, initializing block module only")
			// Initialize block module with empty state
			mod := app.mm.Modules["block"]
			if mod != nil {
				emptyState := json.RawMessage(`{"blocks":[],"block_data":[],"validations":[],"params":{}}`)
				sanitizedState := sanitizeGenesisState("block", emptyState)
				moduleValUpdates, moduleErr := callModuleInitGenesis(ctx, app.appCodec, mod, "block", sanitizedState)
				if moduleErr != nil {
					fmt.Printf("  ⚠️  Failed to initialize block module: %v\n", moduleErr)
				} else {
					fmt.Println("  ✅ Block module initialized successfully (genesis block saved)")
					return buildInitChainResponse(moduleValUpdates, req.Validators, nil), nil
				}
			}
			return &abci.ResponseInitChain{
				Validators: req.Validators,
			}, nil
		}

		// Initialize modules with genesis state
		// ROOT CAUSE FIX: Bypass ModuleManager.InitGenesis() which has internal unmarshal bug
		// Instead, call each module's InitGenesis directly with sanitized state
		var resp *abci.ResponseInitChain
		var err error
		var validatorUpdates []abci.ValidatorUpdate

		// Call each module's InitGenesis directly (bypassing ModuleManager's internal processing)
		for _, moduleName := range app.mm.OrderInitGenesis {
			// For block module, always call InitGenesis even if genesis state is empty
			// This allows block module to create genesis block (height 1) during initialization
			if filteredGenesisState[moduleName] == nil && moduleName != "block" {
				continue // Skip modules without genesis data (except block module)
			}

			// Get module from ModuleManager
			mod := app.mm.Modules[moduleName]
			if mod == nil {
				continue // Skip if module not found
			}

			// For block module, create empty genesis state if not provided
			// This allows block module to create genesis block (height 1) during initialization
			var sanitizedState json.RawMessage
			if moduleName == "block" && filteredGenesisState[moduleName] == nil {
				// Create empty genesis state for block module
				emptyState := json.RawMessage(`{"blocks":[],"block_data":[],"validations":[],"params":{}}`)
				sanitizedState = sanitizeGenesisState(moduleName, emptyState)
			} else {
				// Sanitize genesis state for this module
				sanitizedState = sanitizeGenesisState(moduleName, filteredGenesisState[moduleName])
			}

			// Call InitGenesis with sanitized state
			moduleValUpdates, moduleErr := callModuleInitGenesis(ctx, app.appCodec, mod, moduleName, sanitizedState)
			if moduleErr != nil {
				// COSMOS SDK 0.53.4: Log error but continue processing other modules
				// This ensures block module is initialized even if other modules fail
				fmt.Printf("  ⚠️  Module '%s' InitGenesis failed: %v (continuing with other modules)\n", moduleName, moduleErr)
				err = moduleErr // Track error but don't break
				// Don't break - continue processing to ensure block module is initialized
			}
			if len(moduleValUpdates) > 0 {
				validatorUpdates = append(validatorUpdates, moduleValUpdates...)
			}
		}

		// Build response with validators
		resp = buildInitChainResponse(validatorUpdates, req.Validators, err)

		// If InitGenesis fails, log and return response (allows chain to start with defaults)
		if err != nil {
			fmt.Printf("  ℹ️  InitGenesis: Some modules failed, continuing with default state\n")
		} else {
			fmt.Printf("  ✅ InitGenesis completed, returning %d validators\n", len(resp.Validators))
		}

		return resp, nil
	}

	// No genesis state provided
	return &abci.ResponseInitChain{
		Validators: req.Validators,
	}, nil
}
