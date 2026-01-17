package app

// ============================================================================
// MODULE INITIALIZATION ORDER
// ============================================================================
// This file defines the order in which modules are initialized during genesis
// and processed during block execution. Order matters due to dependencies.

// Module initialization dependencies:
// - auth: No dependencies (must be first)
// - bank: Depends on auth (must be after auth)
// - usc_coin: Depends on bank (must be after bank)
// - block: No dependencies, but should be early to ensure genesis block is saved
// - staking: Depends on bank (must be after bank)
// - transaction: Depends on auth, bank (must be after auth, bank)
// - validator: Depends on staking (must be after staking)
// - Other modules: No strict dependencies, can be in any order
// - Observability modules (monitoring, performance): Should be last as they depend on other modules

// InitGenesisOrder defines the order for module initialization during genesis
// COSMOS SDK 0.53.4: Order matters - dependencies must be initialized before dependents
var InitGenesisOrder = []string{
	"auth",                // First: No dependencies
	"bank",                // After auth: Depends on auth
	"usc_coin",            // After bank: Depends on bank
	"block",               // Early: No dependencies, but should be initialized early to ensure genesis block is saved
	"staking",             // After bank: Depends on bank
	"transaction",         // After auth, bank: Depends on auth, bank
	"validator",           // After staking: Depends on staking
	"nft_token",           // No strict dependencies
	"smart_contract",      // No strict dependencies
	"custom_token",        // No strict dependencies
	"network",             // No strict dependencies
	"product_certificate", // No strict dependencies
	"store_bridge",        // No strict dependencies
	"store_network",       // No strict dependencies
	"streaming",           // No strict dependencies
	"monitoring",          // Last: Observability module
	"performance",         // Last: Observability module
}

// BeginBlockersOrder defines the order for BeginBlock processing
// COSMOS SDK 0.53.4: Order matters for state updates
var BeginBlockersOrder = []string{
	"auth",                // First: Account updates
	"bank",                // After auth: Balance updates
	"staking",             // After bank: Staking operations
	"usc_coin",            // After bank: USC coin operations
	"block",               // Block processing
	"transaction",         // Transaction processing
	"validator",           // Validator operations
	"nft_token",           // NFT operations
	"smart_contract",      // Smart contract operations
	"custom_token",        // Custom token operations
	"network",             // Network operations
	"product_certificate", // Product certificate operations
	"store_bridge",        // Store bridge operations
	"store_network",       // Store network operations
	"streaming",           // Streaming operations
	"monitoring",          // Monitoring operations
	"performance",         // Performance operations
}

// EndBlockersOrder defines the order for EndBlock processing
// COSMOS SDK 0.53.4: Order matters for state finalization
var EndBlockersOrder = []string{
	"auth",                // First: Account finalization
	"bank",                // After auth: Balance finalization
	"staking",             // After bank: Staking finalization
	"usc_coin",            // After bank: USC coin finalization
	"block",               // Block finalization (creates new blocks)
	"transaction",         // Transaction finalization
	"validator",           // Validator finalization
	"nft_token",           // NFT finalization
	"smart_contract",      // Smart contract finalization
	"custom_token",        // Custom token finalization
	"network",             // Network finalization
	"product_certificate", // Product certificate finalization
	"store_bridge",        // Store bridge finalization
	"store_network",       // Store network finalization
	"streaming",           // Streaming finalization
	"monitoring",          // Monitoring finalization
	"performance",         // Performance finalization
}
