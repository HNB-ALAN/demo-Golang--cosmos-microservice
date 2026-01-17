package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	block "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/block"
	customtoken "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/custom_token"
	monitoring "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/monitoring"
	network "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/network"
	nfttoken "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/nft_token"
	performance "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/performance"
	productcertificate "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/product_certificate"
	smartcontract "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/smart_contract"
	storebridge "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/store_bridge"
	storenetwork "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/store_network"
	streaming "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/streaming"
	transaction "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/transaction"
	usccoin "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/usc_coin"
	validator "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/validator"
)

// ============================================================================
// MODULE REGISTRATION
// ============================================================================

// registerModules registers all modules and sets up handlers
func (app *USCApp) registerModules() error {
	app.Logger().Info("Registering modules", "component", "module_registration")

	// Core modules
	authModule := auth.NewAppModule(app.appCodec, app.AccountKeeper, nil, nil)
	bankModule := bank.NewAppModule(app.appCodec, app.BankKeeper, app.AccountKeeper, nil)

	// Staking module (required by Validator module)
	// staking.NewAppModule requires: codec, keeper, accountKeeper, bankKeeper, subspace
	stakingModule := staking.NewAppModule(
		app.appCodec,
		app.StakingKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		app.stakingSubspace,
	)

	// USC modules
	uscCoinModule := usccoin.NewAppModule(app.appCodec, app.USCCoinKeeper, app.BankKeeper)
	blockModule := block.NewAppModule(app.appCodec, app.BlockKeeper)

	// Transaction module - AccountKeeper and BankKeeper are stored but not actively used in keeper
	// Pass nil (similar to how other modules handle unused dependencies)
	transactionModule := transaction.NewAppModule(app.appCodec, app.TransactionKeeper, nil, nil)

	// Validator module - requires paramSpace
	validatorSubspace := paramtypes.NewSubspace(
		app.appCodec,
		app.legacyAmino,
		ParamsStoreKey,
		ValidatorStoreKey,
		"validator",
	)
	validatorModule := validator.NewAppModule(app.appCodec, app.ValidatorKeeper, validatorSubspace)

	// NFT Token module
	nftTokenModule := nfttoken.NewAppModule(app.appCodec, app.NFTTokenKeeper)

	// Smart Contract module
	smartContractModule := smartcontract.NewAppModule(app.appCodec, app.SmartContractKeeper)

	// Custom Token module
	customTokenSubspace := paramtypes.NewSubspace(
		app.appCodec,
		app.legacyAmino,
		ParamsStoreKey,
		CustomTokenStoreKey,
		"custom_token",
	)
	customTokenModule := customtoken.NewAppModule(app.appCodec, app.CustomTokenKeeper, customTokenSubspace)

	// Network module
	networkModule := network.NewAppModule(app.appCodec, app.NetworkKeeper)

	// Product Certificate module
	productCertificateSubspace := paramtypes.NewSubspace(
		app.appCodec,
		app.legacyAmino,
		ParamsStoreKey,
		ProductCertificateStoreKey,
		"product_certificate",
	)
	productCertificateModule := productcertificate.NewAppModule(app.appCodec, app.ProductCertificateKeeper, productCertificateSubspace)

	// Store Bridge module
	storeBridgeModule := storebridge.NewAppModule(app.appCodec, app.StoreBridgeKeeper)

	// Store Network module
	storeNetworkModule := storenetwork.NewAppModule(app.appCodec, app.StoreNetworkKeeper)

	// Streaming module
	streamingModule := streaming.NewAppModule(app.appCodec, app.StreamingKeeper)

	// Monitoring module
	monitoringModule := monitoring.NewAppModule(app.appCodec, app.MonitoringKeeper)

	// Performance module
	performanceModule := performance.NewAppModule(app.appCodec, app.PerformanceKeeper)

	// Create ModuleManager and add modules
	app.mm = module.NewManager(
		authModule,
		bankModule,
		stakingModule,
		uscCoinModule,
		blockModule,
		transactionModule,
		validatorModule,
		nftTokenModule,
		smartContractModule,
		customTokenModule,
		networkModule,
		productCertificateModule,
		storeBridgeModule,
		storeNetworkModule,
		streamingModule,
		monitoringModule,
		performanceModule,
	)

	// Set module initialization order for InitGenesis
	// COSMOS SDK 0.53.4: Order matters - dependencies must be initialized before dependents
	// See module_order.go for detailed dependency documentation
	app.mm.SetOrderInitGenesis(InitGenesisOrder...)

	// Set block processing order
	// COSMOS SDK 0.53.4: Order matters for state updates and finalization
	app.mm.SetOrderBeginBlockers(BeginBlockersOrder...)
	app.mm.SetOrderEndBlockers(EndBlockersOrder...)

	// Register module services (gRPC)
	app.mm.RegisterServices(module.NewConfigurator(app.appCodec, app.BaseApp.MsgServiceRouter(), app.BaseApp.GRPCQueryRouter()))

	// Set up handlers (required by Cosmos SDK 0.53.4)
	app.SetInitChainer(app.initChainer)
	app.SetBeginBlocker(app.beginBlocker)
	app.SetEndBlocker(app.endBlocker)

	// Set AnteHandler for transaction validation with full fee and signature verification
	app.SetAnteHandler(NewAnteHandler(app.AccountKeeper, app.BankKeeper))

	app.Logger().Info("Modules registered successfully", "component", "module_registration")
	return nil
}

// ============================================================================
// TRANSACTION HANDLER
// ============================================================================

// NewAnteHandler creates a full AnteHandler with fee deduction and basic transaction validation
// This follows the standard Cosmos SDK pattern for transaction preprocessing
func NewAnteHandler(
	accountKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		// Setup decorator - sets up gas meter and panic recovery
		ante.NewSetUpContextDecorator(),

		// Extension options decorator - validates extension options
		ante.NewExtensionOptionsDecorator(nil),

		// Validate basic decorator - validates transaction structure
		ante.NewValidateBasicDecorator(),

		// Tx timeout height decorator - validates timeout height
		ante.NewTxTimeoutHeightDecorator(),

		// Validate memo decorator - validates memo size
		ante.NewValidateMemoDecorator(accountKeeper),

		// Consume gas for tx size - charges gas based on transaction size
		ante.NewConsumeGasForTxSizeDecorator(accountKeeper),

		// Deduct fee decorator - deducts transaction fees from sender
		ante.NewDeductFeeDecorator(accountKeeper, bankKeeper, nil, nil),

		// Set pub key decorator - sets public key for account
		ante.NewSetPubKeyDecorator(accountKeeper),

		// Increment sequence decorator - increments account sequence (nonce)
		ante.NewIncrementSequenceDecorator(accountKeeper),
	)
}
