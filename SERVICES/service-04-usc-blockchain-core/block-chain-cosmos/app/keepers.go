package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	blockkeeper "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/block/keeper"
	customtokenkeeper "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/custom_token/keeper"
	monitoringkeeper "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/monitoring/keeper"
	networkkeeper "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/network/keeper"
	nfttokenkeeper "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/nft_token/keeper"
	performancekeeper "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/performance/keeper"
	productcertificatekeeper "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/product_certificate/keeper"
	smartcontractkeeper "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/smart_contract/keeper"
	storebridgekeeper "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/store_bridge/keeper"
	storenetworkkeeper "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/store_network/keeper"
	streamingkeeper "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/streaming/keeper"
	transactionkeeper "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/transaction/keeper"
	usccoinkeeper "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/usc_coin/keeper"
	validatorkeeper "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/validator/keeper"
)

// ============================================================================
// KEEPER INITIALIZATION
// ============================================================================

// initializeKeepers initializes all keepers with proper dependencies
func (app *USCApp) initializeKeepers() error {
	app.Logger().Info("Initializing keepers", "component", "keeper_initialization")

	// Initialize Account Keeper (required by BankKeeper)
	// Cosmos SDK 0.53.4 requires: codec, storeService, accountFactory, maccPerms, addressCodec, bech32Prefix, authority
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		app.appCodec,
		NewKVStoreService(AccountStoreKey),
		authtypes.ProtoBaseAccount,
		app.GetMaccPerms(),
		NewAddressCodec(),
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		authtypes.NewModuleAddress("gov").String(),
	)
	app.Logger().Debug("Account keeper initialized", "component", "keeper_initialization")

	// Initialize Bank Keeper
	app.BankKeeper = bankkeeper.NewBaseKeeper(
		app.appCodec,
		NewKVStoreService(BankStoreKey),
		app.AccountKeeper,
		app.ModuleAccountAddrs(),                   // blocked module accounts
		authtypes.NewModuleAddress("gov").String(), // authority address expected bech32
		app.Logger(),
	)
	app.Logger().Debug("Bank keeper initialized", "component", "keeper_initialization")

	// Create params subspace for USC coin module
	uscCoinSubspace := paramtypes.NewSubspace(
		app.appCodec,
		app.legacyAmino,
		ParamsStoreKey,
		USCCoinStoreKey,
		"usc_coin",
	)

	// Initialize USC Coin keeper with real BankKeeper
	app.USCCoinKeeper = usccoinkeeper.NewKeeper(
		app.appCodec,
		USCCoinStoreKey,
		uscCoinSubspace,
		app.BankKeeper,
	)
	app.Logger().Debug("USC Coin keeper initialized", "component", "keeper_initialization", "note", "with BankKeeper")

	// Initialize Block keeper
	app.BlockKeeper = blockkeeper.NewKeeper(
		app.appCodec,
		BlockStoreKey,
	)
	app.Logger().Debug("Block keeper initialized", "component", "keeper_initialization")

	// Create params subspace for Transaction module
	transactionSubspace := paramtypes.NewSubspace(
		app.appCodec,
		app.legacyAmino,
		ParamsStoreKey,
		TransactionStoreKey,
		"transaction",
	)

	// Initialize Transaction keeper
	app.TransactionKeeper = *transactionkeeper.NewKeeper(
		app.appCodec,
		TransactionStoreKey,
		MemStoreKey, // Use MemStoreKey for transaction mempool
		transactionSubspace,
	)
	app.Logger().Debug("Transaction keeper initialized", "component", "keeper_initialization")

	// Create params subspace for Staking module
	app.stakingSubspace = paramtypes.NewSubspace(
		app.appCodec,
		app.legacyAmino,
		ParamsStoreKey,
		StakingStoreKey,
		"staking",
	)

	// Initialize Staking Keeper (required by Validator module)
	// Cosmos SDK 0.53.4 requires: codec, storeService, accountKeeper, bankKeeper, authority, validatorAddressCodec, consensusAddressCodec
	app.StakingKeeper = stakingkeeper.NewKeeper(
		app.appCodec,
		NewKVStoreService(StakingStoreKey),
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.NewModuleAddress("gov").String(), // authority address
		NewValidatorAddressCodec(),
		NewConsensusAddressCodec(),
	)
	app.Logger().Debug("Staking keeper initialized", "component", "keeper_initialization")

	// Create params subspace for Validator module
	validatorSubspace := paramtypes.NewSubspace(
		app.appCodec,
		app.legacyAmino,
		ParamsStoreKey,
		ValidatorStoreKey,
		"validator",
	)

	// Initialize Validator keeper
	app.ValidatorKeeper = validatorkeeper.NewKeeper(
		app.appCodec,
		ValidatorStoreKey,
		validatorSubspace,
	)
	app.Logger().Debug("Validator keeper initialized", "component", "keeper_initialization")

	// Initialize NFT Token keeper
	app.NFTTokenKeeper = nfttokenkeeper.NewKeeper(
		app.appCodec,
		NFTTokenStoreKey,
	)
	app.Logger().Debug("NFT Token keeper initialized", "component", "keeper_initialization")

	// Initialize Smart Contract keeper
	app.SmartContractKeeper = smartcontractkeeper.NewKeeper(
		app.appCodec,
		SmartContractStoreKey,
	)
	app.Logger().Debug("Smart Contract keeper initialized", "component", "keeper_initialization")

	// Initialize Custom Token keeper
	customTokenSubspace := paramtypes.NewSubspace(
		app.appCodec,
		app.legacyAmino,
		ParamsStoreKey,
		CustomTokenStoreKey,
		"custom_token",
	)
	app.CustomTokenKeeper = customtokenkeeper.NewKeeper(
		app.appCodec,
		CustomTokenStoreKey,
		customTokenSubspace,
	)
	app.Logger().Debug("Custom Token keeper initialized", "component", "keeper_initialization")

	// Initialize Network keeper
	app.NetworkKeeper = networkkeeper.NewKeeper(
		app.appCodec,
		NetworkStoreKey,
	)
	app.Logger().Debug("Network keeper initialized", "component", "keeper_initialization")

	// Initialize Product Certificate keeper
	productCertificateSubspace := paramtypes.NewSubspace(
		app.appCodec,
		app.legacyAmino,
		ParamsStoreKey,
		ProductCertificateStoreKey,
		"product_certificate",
	)
	app.ProductCertificateKeeper = productcertificatekeeper.NewKeeper(
		app.appCodec,
		ProductCertificateStoreKey,
		productCertificateSubspace,
	)
	app.Logger().Debug("Product Certificate keeper initialized", "component", "keeper_initialization")

	// Initialize Store Bridge keeper
	app.StoreBridgeKeeper = storebridgekeeper.NewKeeper(
		app.appCodec,
		StoreBridgeStoreKey,
	)
	app.Logger().Debug("Store Bridge keeper initialized", "component", "keeper_initialization")

	// Initialize Store Network keeper
	app.StoreNetworkKeeper = storenetworkkeeper.NewKeeper(
		app.appCodec,
		StoreNetworkStoreKey,
	)
	app.Logger().Debug("Store Network keeper initialized", "component", "keeper_initialization")

	// Initialize Streaming keeper
	app.StreamingKeeper = streamingkeeper.NewKeeper(
		app.appCodec,
		StreamingStoreKey,
	)
	app.Logger().Debug("Streaming keeper initialized", "component", "keeper_initialization")

	// Initialize Monitoring keeper
	app.MonitoringKeeper = monitoringkeeper.NewKeeper(
		app.appCodec,
		MonitoringStoreKey,
	)
	app.Logger().Debug("Monitoring keeper initialized", "component", "keeper_initialization")

	// Initialize Performance keeper
	app.PerformanceKeeper = performancekeeper.NewKeeper(
		app.appCodec,
		PerformanceStoreKey,
	)
	app.Logger().Debug("Performance keeper initialized", "component", "keeper_initialization")

	app.Logger().Debug("Other core keepers stubbed", "component", "keeper_initialization", "note", "slashing, gov, etc")
	app.Logger().Debug("Other USC keepers stubbed", "component", "keeper_initialization", "note", "will be implemented after core")
	app.Logger().Info("Keepers initialized successfully",
		"component", "keeper_initialization",
		"functional_keepers", "Auth, Bank, USC Coin, Block, Transaction, Staking, Validator, NFT Token, Smart Contract, Custom Token, Network, Product Certificate, Store Bridge, Store Network, Streaming, Monitoring, Performance")
	return nil
}
