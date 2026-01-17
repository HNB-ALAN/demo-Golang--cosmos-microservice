package app

// ============================================================================
// SECTION 1: DECLARATIONS
// ============================================================================
// This section contains all package-level declarations including:
// - Imports (standard library, Cosmos SDK, CometBFT, USC modules)
// - Constants (application name)
// - Variables (store keys, default paths)
// - Types (USCApp struct, ConsensusParamStore)
// - Module Basics (exported module definitions)
// - Init function (package initialization)

import (
	// Standard library imports
	"fmt"
	"os"
	"path/filepath"

	// Cosmos SDK core imports
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	cosmosdb "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/rs/zerolog"

	// Core modules
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	// Evidence and Upgrade
	"cosmossdk.io/x/evidence"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	"cosmossdk.io/x/upgrade"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"

	// USC modules - All 14 modules
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/block"
	blockkeeper "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/block/keeper"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/custom_token"
	customtokenkeeper "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/custom_token/keeper"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/monitoring"
	monitoringkeeper "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/monitoring/keeper"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/network"
	networkkeeper "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/network/keeper"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/nft_token"
	nfttokenkeeper "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/nft_token/keeper"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/performance"
	performancekeeper "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/performance/keeper"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/product_certificate"
	productcertificatekeeper "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/product_certificate/keeper"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/smart_contract"
	smartcontractkeeper "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/smart_contract/keeper"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/store_bridge"
	storebridgekeeper "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/store_bridge/keeper"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/store_network"
	storenetworkkeeper "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/store_network/keeper"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/streaming"
	streamingkeeper "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/streaming/keeper"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/transaction"
	transactionkeeper "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/transaction/keeper"
	usccoin "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/usc_coin"
	usccoinkeeper "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/usc_coin/keeper"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/validator"
	validatorkeeper "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/validator/keeper"
)

// ============================================================================
// SECTION 1.1: CONSTANTS
// ============================================================================

// Name is the application name used throughout the Cosmos SDK
const (
	Name = "usc" // Application identifier for USC blockchain
)

// ============================================================================
// SECTION 1.2: VARIABLES - DEFAULT PATHS
// ============================================================================

// DefaultNodeHome is the default directory where node data is stored
// Initialized in init() function to use user's home directory
var (
	DefaultNodeHome string // Default: ~/.uscd
)

// ============================================================================
// SECTION 1.3: VARIABLES - STORE KEYS
// ============================================================================
// Store keys are used to identify different KV stores in the application.
// Each module has its own store key for data isolation.

var (
	// Core store keys - Base application stores
	// Note: MainStoreKey is reserved for future use (currently not mounted)
	MainStoreKey = storetypes.NewKVStoreKey("main")
	MemStoreKey  = storetypes.NewMemoryStoreKey("mem_capability") // Memory store for transient data

	// Cosmos SDK store keys - Standard Cosmos SDK modules
	// These store keys are used by Cosmos SDK modules for state management
	AccountStoreKey   = storetypes.NewKVStoreKey("account")      // Account and authentication data
	BankStoreKey      = storetypes.NewKVStoreKey("bank")         // Token balances and transfers
	StakingStoreKey   = storetypes.NewKVStoreKey("staking")      // Validator staking data (not yet initialized)
	SlashingStoreKey  = storetypes.NewKVStoreKey("slashing")     // Validator slashing data (not yet initialized)
	MintStoreKey      = storetypes.NewKVStoreKey("mint")         // Token minting data (not yet initialized)
	DistrStoreKey     = storetypes.NewKVStoreKey("distribution") // Distribution rewards (not yet initialized)
	GovStoreKey       = storetypes.NewKVStoreKey("gov")          // Governance proposals (not yet initialized)
	ConsensusStoreKey = storetypes.NewKVStoreKey("consensus")    // Consensus parameters (not yet initialized)
	UpgradeStoreKey   = storetypes.NewKVStoreKey("upgrade")      // Upgrade coordination (not yet initialized)
	EvidenceStoreKey  = storetypes.NewKVStoreKey("evidence")     // Evidence handling (not yet initialized)

	// USC custom store keys - All 14 USC-specific modules
	// These store keys are used by USC custom modules for specialized functionality
	// Note: Only USCCoinStoreKey is currently in use; others reserved for future implementation
	USCCoinStoreKey            = storetypes.NewKVStoreKey("usc_coin")            // ✅ USC token management (initialized)
	BlockStoreKey              = storetypes.NewKVStoreKey("block")               // Block metadata (not yet initialized)
	CustomTokenStoreKey        = storetypes.NewKVStoreKey("custom_token")        // Custom token creation (not yet initialized)
	MonitoringStoreKey         = storetypes.NewKVStoreKey("monitoring")          // System monitoring (not yet initialized)
	NetworkStoreKey            = storetypes.NewKVStoreKey("network")             // Network management (not yet initialized)
	NFTTokenStoreKey           = storetypes.NewKVStoreKey("nft_token")           // NFT token management (not yet initialized)
	PerformanceStoreKey        = storetypes.NewKVStoreKey("performance")         // Performance metrics (not yet initialized)
	ProductCertificateStoreKey = storetypes.NewKVStoreKey("product_certificate") // Product certificates (not yet initialized)
	SmartContractStoreKey      = storetypes.NewKVStoreKey("smart_contract")      // Smart contract execution (not yet initialized)
	StoreBridgeStoreKey        = storetypes.NewKVStoreKey("store_bridge")        // Store bridging (not yet initialized)
	StoreNetworkStoreKey       = storetypes.NewKVStoreKey("store_network")       // Store network (not yet initialized)
	StreamingStoreKey          = storetypes.NewKVStoreKey("streaming")           // Streaming data (not yet initialized)
	TransactionStoreKey        = storetypes.NewKVStoreKey("transaction")         // Transaction records (not yet initialized)
	ValidatorStoreKey          = storetypes.NewKVStoreKey("validator")           // Validator management (not yet initialized)

	// Params subspace key - Used for parameter management across modules
	ParamsStoreKey = storetypes.NewKVStoreKey("params") // Module parameters and configuration
)

// ============================================================================
// SECTION 1.4: TYPES - APPLICATION STRUCTURE
// ============================================================================

// USCApp represents the USC blockchain application
// It embeds Cosmos SDK's BaseApp and adds USC-specific functionality
// Contains all keepers, module managers, and application state
type USCApp struct {
	// BaseApp is the core Cosmos SDK application instance
	// Handles ABCI calls, transaction routing, and state management
	*baseapp.BaseApp

	// Codec and encoding - Used for serialization/deserialization
	legacyAmino       *codec.LegacyAmino      // Legacy Amino codec for backward compatibility
	appCodec          codec.Codec             // Primary application codec (protobuf-based)
	txConfig          client.TxConfig         // Transaction encoding/decoding configuration
	interfaceRegistry types.InterfaceRegistry // Type registry for interface resolution

	// Module management
	sm *module.SimulationManager // Simulation manager for testing
	mm *module.Manager           // Module manager - orchestrates all modules

	// Core Cosmos SDK keepers - State management for standard modules
	AccountKeeper authkeeper.AccountKeeper // ✅ Initialized - Account management (authentication, nonces)
	BankKeeper    bankkeeper.Keeper        // ✅ Initialized - Token balances, transfers, supply

	// Additional Cosmos SDK keepers (declared but not initialized yet)
	// These will be implemented after core modules are stable
	StakingKeeper   *stakingkeeper.Keeper  // Validator staking, delegation
	stakingSubspace paramtypes.Subspace    // Staking module parameter subspace
	SlashingKeeper  slashingkeeper.Keeper  // Validator slashing penalties
	MintKeeper      mintkeeper.Keeper      // Token minting
	DistrKeeper     distrkeeper.Keeper     // Distribution rewards
	GovKeeper       govkeeper.Keeper       // Governance proposals and voting
	ConsensusKeeper consensuskeeper.Keeper // Consensus parameters
	UpgradeKeeper   upgradekeeper.Keeper   // Upgrade coordination
	EvidenceKeeper  evidencekeeper.Keeper  // Evidence handling

	// USC custom keepers - All 14 USC-specific modules
	// Note: Only USCCoinKeeper is currently initialized
	// Other USC keepers are declared but will be implemented progressively
	USCCoinKeeper            usccoinkeeper.Keeper            // ✅ Initialized - USC token operations
	BlockKeeper              blockkeeper.Keeper              // Block metadata management
	CustomTokenKeeper        customtokenkeeper.Keeper        // Custom token creation
	MonitoringKeeper         monitoringkeeper.Keeper         // System monitoring
	NetworkKeeper            networkkeeper.Keeper            // Network management
	NFTTokenKeeper           nfttokenkeeper.Keeper           // NFT token management
	PerformanceKeeper        performancekeeper.Keeper        // Performance metrics
	ProductCertificateKeeper productcertificatekeeper.Keeper // Product certificates
	SmartContractKeeper      smartcontractkeeper.Keeper      // Smart contract execution
	StoreBridgeKeeper        storebridgekeeper.Keeper        // Store bridging
	StoreNetworkKeeper       storenetworkkeeper.Keeper       // Store network
	StreamingKeeper          streamingkeeper.Keeper          // Streaming data
	TransactionKeeper        transactionkeeper.Keeper        // Transaction records
	ValidatorKeeper          validatorkeeper.Keeper          // Validator management

	// Database reference - Used for state persistence
	// RocksDB backend for high-performance key-value storage
	db cosmosdb.DB
}

// NOTE: EncodingConfig type moved to encoding.go

// NOTE: Store Adapters and Address Codec moved to adapters.go

// ============================================================================
// SECTION 1.8: MODULE ACCOUNT PERMISSIONS
// ============================================================================

// maccPerms defines module account permissions used by core modules
var maccPerms = map[string][]string{
	authtypes.FeeCollectorName:     nil, // fee collector (no permissions)
	authtypes.Minter:               {authtypes.Minter},
	authtypes.Burner:               {authtypes.Burner},
	stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
	stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
}

// moduleAccountAddrs defines which module accounts are blocked from receiving user funds
var moduleAccountAddrs = map[string]bool{
	authtypes.FeeCollectorName:     true,
	authtypes.Minter:               true,
	authtypes.Burner:               true,
	stakingtypes.BondedPoolName:    true,
	stakingtypes.NotBondedPoolName: true,
}

// ============================================================================
// SECTION 1.9: MODULE BASICS (Exported)
// ============================================================================

// ModuleBasics defines the module BasicManager - All 14 modules
var ModuleBasics = module.NewBasicManager(
	// Core Cosmos SDK modules
	auth.AppModuleBasic{},
	bank.AppModuleBasic{},
	staking.AppModuleBasic{},
	mint.AppModuleBasic{},
	distribution.AppModuleBasic{},
	gov.AppModuleBasic{},
	slashing.AppModuleBasic{},
	consensus.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},

	// USC custom modules - All 14 modules
	usccoin.AppModuleBasic{},
	block.AppModuleBasic{},
	custom_token.AppModuleBasic{},
	monitoring.AppModuleBasic{},
	network.AppModuleBasic{},
	nft_token.AppModuleBasic{},
	performance.AppModuleBasic{},
	product_certificate.AppModuleBasic{},
	smart_contract.AppModuleBasic{},
	store_bridge.AppModuleBasic{},
	store_network.AppModuleBasic{},
	streaming.AppModuleBasic{},
	transaction.AppModuleBasic{},
	validator.AppModuleBasic{},
)

// ============================================================================
// SECTION 1.10: INIT FUNCTION
// ============================================================================

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		// Log error with context before panic
		// Note: Logger not available in init(), so use fmt for error message
		panic(fmt.Sprintf("failed to get user home directory: %v (required for DefaultNodeHome)", err))
	}
	DefaultNodeHome = filepath.Join(userHomeDir, ".uscd")
}

// ============================================================================
// SECTION 2: APP CONSTRUCTION
// ============================================================================

// NOTE: Encoding Configuration and readGenesisChainID moved to encoding.go

// ============================================================================
// SECTION 2.3: APPLICATION CONSTRUCTOR
// ============================================================================

// NewUSCApp creates a new USC application
func NewUSCApp(db cosmosdb.DB) *USCApp {
	// Ensure Bech32 prefixes are set
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount("usc", "uscpub")
	cfg.SetBech32PrefixForValidator("uscvaloper", "uscvaloperpub")
	cfg.SetBech32PrefixForConsensusNode("uscvalcons", "uscvalconspub")
	cfg.Seal()

	encodingConfig := MakeEncodingConfig()

	// Create structured logger for production
	// Use log.NewLogger with proper configuration for production observability
	// Default to info level, JSON format for structured logging
	// Note: cosmossdk.io/log.NewLogger signature: NewLogger(w io.Writer, options ...Option) Logger
	// Use LevelOption, OutputJSONOption, and ColorOption for structured logging
	logger := log.NewLogger(os.Stdout,
		log.LevelOption(zerolog.InfoLevel), // Info level for production
		log.OutputJSONOption(),             // JSON format for structured logging
		log.ColorOption(false),             // Disable color in production
	)

	// Read chain-id from genesis.json to set in BaseApp
	// This ensures BaseApp.chainID is set before LoadLatestVersion()
	// BaseApp.InitChain() validates req.ChainId against app.chainID
	genesisChainID := readGenesisChainID()

	// Create BaseApp with database and proper parameters for Cosmos SDK 0.53.4
	// Set chain-id via SetChainID() option if found in genesis
	var baseAppOptions []func(*baseapp.BaseApp)
	if genesisChainID != "" {
		baseAppOptions = append(baseAppOptions, baseapp.SetChainID(genesisChainID))
	}
	bApp := baseapp.NewBaseApp(
		Name,
		logger,
		db,
		encodingConfig.TxConfig.TxDecoder(),
		baseAppOptions..., // Pass SetChainID option if chain-id found
	)
	bApp.SetVersion("1.0.0")
	// Ensure routers/registries are initialized
	bApp.SetInterfaceRegistry(encodingConfig.InterfaceRegistry)
	bApp.SetTxEncoder(encodingConfig.TxConfig.TxEncoder())
	if bApp.MsgServiceRouter() == nil {
		bApp.SetMsgServiceRouter(baseapp.NewMsgServiceRouter())
	}
	if bApp.GRPCQueryRouter() == nil {
		bApp.SetGRPCQueryRouter(baseapp.NewGRPCQueryRouter())
	}

	app := &USCApp{
		BaseApp:           bApp,
		legacyAmino:       encodingConfig.Amino,
		appCodec:          encodingConfig.Codec,
		txConfig:          encodingConfig.TxConfig,
		interfaceRegistry: encodingConfig.InterfaceRegistry,
		db:                db, // Store database reference for checking if empty
	}

	// Mount stores after BaseApp is created with database
	app.mountStores()

	// Set ParamStore for BaseApp to store consensus params during InitChain
	// BaseApp requires ParamStore interface to store consensus parameters from genesis
	// Create a simple ParamStore adapter using params subspace
	paramSubspace := paramtypes.NewSubspace(
		app.appCodec,
		app.legacyAmino,
		ParamsStoreKey,
		ParamsStoreKey,
		"baseapp",
	)
	paramStore := NewConsensusParamStore(paramSubspace, app.db, ParamsStoreKey)
	bApp.SetParamStore(paramStore)
	// Log using BaseApp logger (structured logging)
	bApp.Logger().Info("ParamStore set for BaseApp", "component", "consensus_params")

	app.sm = module.NewSimulationManager()
	app.sm.RegisterStoreDecoders()

	return app
}

// ============================================================================
// SECTION 2.4: STORE MOUNTING
// ============================================================================

// mountStores mounts all KV and Memory stores for core and USC modules
func (app *USCApp) mountStores() {
	kvStores := map[string]*storetypes.KVStoreKey{
		// Core
		"account":      AccountStoreKey,
		"bank":         BankStoreKey,
		"staking":      StakingStoreKey,
		"slashing":     SlashingStoreKey,
		"mint":         MintStoreKey,
		"distribution": DistrStoreKey,
		"gov":          GovStoreKey,
		"consensus":    ConsensusStoreKey,
		"upgrade":      UpgradeStoreKey,
		"evidence":     EvidenceStoreKey,
		"params":       ParamsStoreKey,
		// USC modules
		"usc_coin":            USCCoinStoreKey,
		"block":               BlockStoreKey,
		"custom_token":        CustomTokenStoreKey,
		"monitoring":          MonitoringStoreKey,
		"network":             NetworkStoreKey,
		"nft_token":           NFTTokenStoreKey,
		"performance":         PerformanceStoreKey,
		"product_certificate": ProductCertificateStoreKey,
		"smart_contract":      SmartContractStoreKey,
		"store_bridge":        StoreBridgeStoreKey,
		"store_network":       StoreNetworkStoreKey,
		"streaming":           StreamingStoreKey,
		"transaction":         TransactionStoreKey,
		"validator":           ValidatorStoreKey,
	}

	app.MountKVStores(kvStores)
	app.MountMemoryStores(map[string]*storetypes.MemoryStoreKey{
		"mem_capability": MemStoreKey,
	})
}

// ============================================================================
// SECTION 2.5: APPLICATION INITIALIZATION
// ============================================================================

// GetName returns the application name
func (app *USCApp) GetName() string {
	return Name
}

// Initialize initializes the application
func (app *USCApp) Initialize() error {
	// Use BaseApp logger for structured logging
	app.BaseApp.Logger().Info("Initializing USC blockchain application",
		"app_name", app.GetName(),
		"component", "app_initialization")

	if err := app.initializeKeepers(); err != nil {
		return fmt.Errorf("failed to initialize keepers: %w", err)
	}

	if err := app.registerModules(); err != nil {
		return fmt.Errorf("failed to register modules: %w", err)
	}

	// Load latest version as per Cosmos SDK standard initialization flow
	// BaseApp.chainID was already set via SetChainID() option in NewUSCApp()
	// when genesis.json was read. This ensures InitChain validation will pass:
	// req.ChainId == app.chainID
	if err := app.LoadLatestVersion(); err != nil {
		return fmt.Errorf("failed to load latest version: %w", err)
	}

	height := app.BaseApp.LastBlockHeight()
	if height == 0 {
		app.BaseApp.Logger().Info("Database was empty - InitChain will initialize from genesis",
			"height", height,
			"component", "genesis_initialization")
		// COSMOS SDK 0.53.4: Trigger InitChain directly (not in goroutine) to ensure state is committed
		// This ensures InitGenesis is called for all modules including block module
		// InitChain will save genesis block (height 1) to keeper via InitGenesis
		// Note: Calling directly (not in goroutine) allows BaseApp to commit state properly
		if err := app.triggerInitChain(); err != nil {
			app.BaseApp.Logger().Warn("Failed to trigger InitChain",
				"error", err.Error(),
				"component", "genesis_initialization",
				"note", "Block 1 will be saved when EndBlocker is called for the first block")
		} else {
			app.BaseApp.Logger().Info("InitChain triggered successfully - genesis block created",
				"height", 1,
				"component", "genesis_initialization")
		}
	} else {
		app.BaseApp.Logger().Info("Loaded existing database state",
			"height", height,
			"component", "app_initialization")

		// Check if block 1 exists in keeper, if not, trigger InitChain to create it
		// This handles the case where database has height=1 but block 1 was not properly saved
		// Use proper error handling instead of panic/recover
		if app.BlockKeeper.StoreKey() == nil {
			app.BaseApp.Logger().Warn("BlockKeeper.StoreKey() is nil, triggering InitChain to create block 1",
				"component", "genesis_initialization")
			if err := app.triggerInitChain(); err != nil {
				app.BaseApp.Logger().Warn("Failed to trigger InitChain",
					"error", err.Error(),
					"component", "genesis_initialization")
			} else {
				app.BaseApp.Logger().Info("InitChain triggered successfully - genesis block created",
					"height", 1,
					"component", "genesis_initialization")
			}
		} else {
			testCtx := app.BaseApp.NewContext(true)
			block, err := app.BlockKeeper.GetBlockByHeight(testCtx, 1)
			if err != nil {
				app.BaseApp.Logger().Warn("Database has height but block 1 not found in keeper, triggering InitChain",
					"height", height,
					"error", err.Error(),
					"component", "genesis_initialization")
				if err := app.triggerInitChain(); err != nil {
					app.BaseApp.Logger().Warn("Failed to trigger InitChain",
						"error", err.Error(),
						"component", "genesis_initialization")
				} else {
					app.BaseApp.Logger().Info("InitChain triggered successfully - genesis block created",
						"height", 1,
						"component", "genesis_initialization")
				}
			} else {
				app.BaseApp.Logger().Info("Block 1 found in keeper",
					"height", block.Height,
					"hash", block.Hash,
					"component", "genesis_initialization")
			}
		}
	}

	app.BaseApp.Logger().Info("USC blockchain application initialized successfully",
		"app_name", app.GetName(),
		"component", "app_initialization")
	return nil
}

// NOTE: Consensus Param Store moved to consensus_params.go
// NOTE: Keeper Initialization moved to keepers.go
// NOTE: Module Registration moved to module_registration.go

// NOTE: Genesis Helpers and Initialization moved to genesis.go
// NOTE: Block Handlers moved to blocks.go
// NOTE: Utility Methods moved to utils.go
