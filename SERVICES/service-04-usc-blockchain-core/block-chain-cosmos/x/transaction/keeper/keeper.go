package keeper

import (
	"fmt"
	"time"

	storetypes "cosmossdk.io/store/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/transaction/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
	}
}

func (k Keeper) Logger(ctx sdk.Context) interface{} {
	return fmt.Sprintf("module: x/%s", types.ModuleName)
}

// GetParams returns the parameters for the transaction module
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	var params types.Params
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the parameters for the transaction module
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// GetTransaction returns a transaction by hash
func (k Keeper) GetTransaction(ctx sdk.Context, hash string) (types.Transaction, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetTransactionKey(hash))
	if bz == nil {
		return types.Transaction{}, false
	}

	var transaction types.Transaction
	if err := k.cdc.Unmarshal(bz, &transaction); err != nil {
		ctx.Logger().Error("Failed to unmarshal transaction",
			"error", err,
			"key", string(types.GetTransactionKey(hash)))
		return types.Transaction{}, false
	}
	return transaction, true
}

// SetTransaction sets a transaction
func (k Keeper) SetTransaction(ctx sdk.Context, transaction types.Transaction) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&transaction)
	store.Set(types.GetTransactionKey(transaction.Hash), bz)
	store.Set(types.GetTransactionIDKey(transaction.ID), []byte(transaction.Hash))
}

// GetAllTransactions returns all transactions
// COSMOS SDK 0.53.4: Handle unmarshal errors gracefully (skip corrupted data)
func (k Keeper) GetAllTransactions(ctx sdk.Context) []types.Transaction {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte(types.TransactionKeyPrefix))
	defer iterator.Close()

	var transactions []types.Transaction
	for ; iterator.Valid(); iterator.Next() {
		var transaction types.Transaction
		// Use Unmarshal instead of MustUnmarshal to handle errors gracefully
		if err := k.cdc.Unmarshal(iterator.Value(), &transaction); err != nil {
			// Skip corrupted transactions (e.g., old format without protobuf tags)
			ctx.Logger().Error("Failed to unmarshal transaction, skipping",
				"error", err,
				"key", string(iterator.Key()))
			continue
		}
		transactions = append(transactions, transaction)
	}
	return transactions
}

// GetTransactionStats returns transaction statistics
func (k Keeper) GetTransactionStats(ctx sdk.Context) types.TransactionStats {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetStatsKey())
	if bz == nil {
		return types.NewTransactionStats()
	}

	var stats types.TransactionStats
	if err := k.cdc.Unmarshal(bz, &stats); err != nil {
		ctx.Logger().Error("Failed to unmarshal transaction stats",
			"error", err)
		return types.NewTransactionStats()
	}
	return stats
}

// SetTransactionStats sets transaction statistics
func (k Keeper) SetTransactionStats(ctx sdk.Context, stats types.TransactionStats) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&stats)
	store.Set(types.GetStatsKey(), bz)
}

// CreateTransaction creates a new transaction
func (k Keeper) CreateTransaction(ctx sdk.Context, hash, id, from, to, amount, txType, data, memo string) (types.Transaction, error) {
	// Validate inputs
	if !types.IsValidAddress(from) {
		return types.Transaction{}, fmt.Errorf("invalid from address: %s", from)
	}
	if !types.IsValidAddress(to) {
		return types.Transaction{}, fmt.Errorf("invalid to address: %s", to)
	}
	if err := types.ValidateAmount(amount); err != nil {
		return types.Transaction{}, err
	}
	if err := types.ValidateTransactionType(txType); err != nil {
		return types.Transaction{}, err
	}

	// Check if transaction already exists
	if _, exists := k.GetTransaction(ctx, hash); exists {
		return types.Transaction{}, fmt.Errorf("transaction with hash %s already exists", hash)
	}

	// Create transaction
	transaction := types.NewTransaction(hash, id, from, to, amount, txType, data, memo)
	k.SetTransaction(ctx, transaction)

	// Update statistics
	stats := k.GetTransactionStats(ctx)
	stats.TotalTransactions++
	stats.PendingTransactions++
	stats.LastTransactionTime = time.Now().Format(time.RFC3339)
	k.SetTransactionStats(ctx, stats)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateTransaction,
			sdk.NewAttribute(types.AttributeKeyTransactionHash, hash),
			sdk.NewAttribute(types.AttributeKeyTransactionID, id),
			sdk.NewAttribute(types.AttributeKeyFromAddress, from),
			sdk.NewAttribute(types.AttributeKeyToAddress, to),
			sdk.NewAttribute(types.AttributeKeyAmount, amount),
			sdk.NewAttribute(types.AttributeKeyTransactionType, txType),
			sdk.NewAttribute(types.AttributeKeyStatus, "pending"),
		),
	)

	return transaction, nil
}

// ValidateTransaction validates a transaction
func (k Keeper) ValidateTransaction(ctx sdk.Context, hash, validator, proof string) error {
	transaction, exists := k.GetTransaction(ctx, hash)
	if !exists {
		return fmt.Errorf("transaction with hash %s not found", hash)
	}

	if transaction.Status != "pending" {
		return fmt.Errorf("transaction %s is not pending", hash)
	}

	// Update transaction status
	transaction.Status = "validated"
	transaction.ValidatedAt = time.Now().Unix()
	transaction.ValidationProof = proof
	k.SetTransaction(ctx, transaction)

	// Update statistics
	stats := k.GetTransactionStats(ctx)
	stats.PendingTransactions--
	stats.ValidatedTransactions++
	k.SetTransactionStats(ctx, stats)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeValidateTransaction,
			sdk.NewAttribute(types.AttributeKeyTransactionHash, hash),
			sdk.NewAttribute(types.AttributeKeyValidator, validator),
			sdk.NewAttribute(types.AttributeKeyStatus, "validated"),
		),
	)

	return nil
}

// ExecuteTransaction executes a transaction
func (k Keeper) ExecuteTransaction(ctx sdk.Context, hash, executor, proof string) error {
	transaction, exists := k.GetTransaction(ctx, hash)
	if !exists {
		return fmt.Errorf("transaction with hash %s not found", hash)
	}

	if transaction.Status != "validated" {
		return fmt.Errorf("transaction %s is not validated", hash)
	}

	// Update transaction status
	transaction.Status = "executed"
	transaction.ExecutedAt = time.Now().Unix()
	transaction.ExecutionProof = proof
	k.SetTransaction(ctx, transaction)

	// Update statistics
	stats := k.GetTransactionStats(ctx)
	stats.ValidatedTransactions--
	stats.ExecutedTransactions++

	// Update success rate
	if stats.TotalTransactions > 0 {
		successRate := float64(stats.ExecutedTransactions) / float64(stats.TotalTransactions) * 100
		stats.SuccessRate = fmt.Sprintf("%.2f", successRate)
	}

	k.SetTransactionStats(ctx, stats)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeExecuteTransaction,
			sdk.NewAttribute(types.AttributeKeyTransactionHash, hash),
			sdk.NewAttribute(types.AttributeKeyExecutor, executor),
			sdk.NewAttribute(types.AttributeKeyStatus, "executed"),
		),
	)

	return nil
}

// CancelTransaction cancels a transaction
func (k Keeper) CancelTransaction(ctx sdk.Context, hash, canceller, reason string) error {
	transaction, exists := k.GetTransaction(ctx, hash)
	if !exists {
		return fmt.Errorf("transaction with hash %s not found", hash)
	}

	if transaction.Status == "executed" || transaction.Status == "cancelled" {
		return fmt.Errorf("transaction %s cannot be cancelled", hash)
	}

	// Update transaction status
	transaction.Status = "cancelled"
	transaction.ErrorMessage = reason
	k.SetTransaction(ctx, transaction)

	// Update statistics
	stats := k.GetTransactionStats(ctx)
	if transaction.Status == "pending" {
		stats.PendingTransactions--
	} else if transaction.Status == "validated" {
		stats.ValidatedTransactions--
	}
	stats.CancelledTransactions++
	k.SetTransactionStats(ctx, stats)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCancelTransaction,
			sdk.NewAttribute(types.AttributeKeyTransactionHash, hash),
			sdk.NewAttribute(types.AttributeKeyCanceller, canceller),
			sdk.NewAttribute(types.AttributeKeyStatus, "cancelled"),
		),
	)

	return nil
}

// GetTransactionsByAddress returns transactions for a specific address
func (k Keeper) GetTransactionsByAddress(ctx sdk.Context, address string) []types.Transaction {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte(types.TransactionKeyPrefix))
	defer iterator.Close()

	var transactions []types.Transaction
	for ; iterator.Valid(); iterator.Next() {
		var transaction types.Transaction
		if err := k.cdc.Unmarshal(iterator.Value(), &transaction); err != nil {
			ctx.Logger().Error("Failed to unmarshal transaction, skipping",
				"error", err,
				"key", string(iterator.Key()))
			continue
		}
		if transaction.FromAddress == address || transaction.ToAddress == address {
			transactions = append(transactions, transaction)
		}
	}
	return transactions
}

// GetTransactionsByType returns transactions of a specific type
func (k Keeper) GetTransactionsByType(ctx sdk.Context, txType string) []types.Transaction {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte(types.TransactionKeyPrefix))
	defer iterator.Close()

	var transactions []types.Transaction
	for ; iterator.Valid(); iterator.Next() {
		var transaction types.Transaction
		if err := k.cdc.Unmarshal(iterator.Value(), &transaction); err != nil {
			ctx.Logger().Error("Failed to unmarshal transaction, skipping",
				"error", err,
				"key", string(iterator.Key()))
			continue
		}
		if transaction.TransactionType == txType {
			transactions = append(transactions, transaction)
		}
	}
	return transactions
}

// GetTransactionsByStatus returns transactions with a specific status
func (k Keeper) GetTransactionsByStatus(ctx sdk.Context, status string) []types.Transaction {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte(types.TransactionKeyPrefix))
	defer iterator.Close()

	var transactions []types.Transaction
	for ; iterator.Valid(); iterator.Next() {
		var transaction types.Transaction
		if err := k.cdc.Unmarshal(iterator.Value(), &transaction); err != nil {
			ctx.Logger().Error("Failed to unmarshal transaction, skipping",
				"error", err,
				"key", string(iterator.Key()))
			continue
		}
		if transaction.Status == status {
			transactions = append(transactions, transaction)
		}
	}
	return transactions
}

// InitGenesis initializes the genesis state
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	// Wrap in panic recovery for graceful error handling
	defer func() {
		if r := recover(); r != nil {
			panicMsg := fmt.Sprintf("%v", r)
			panic(fmt.Sprintf("transaction InitGenesis panic: %s", panicMsg))
		}
	}()

	k.SetParams(ctx, genState.Params)
	k.SetTransactionStats(ctx, genState.TransactionStats)

	for _, transaction := range genState.Transactions {
		k.SetTransaction(ctx, transaction)
	}
}

// ExportGenesis exports the genesis state
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	transactions := k.GetAllTransactions(ctx)
	stats := k.GetTransactionStats(ctx)
	params := k.GetParams(ctx)

	return &types.GenesisState{
		Transactions:     transactions,
		TransactionStats: stats,
		Params:           params,
	}
}

// BeginBlocker processes transaction timeouts
func (k Keeper) BeginBlocker(ctx sdk.Context) {
	params := k.GetParams(ctx)
	timeout := time.Duration(params.TransactionTimeout) * time.Second
	cutoff := time.Now().Add(-timeout).Unix()

	// Process pending transactions that have timed out
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, []byte(types.TransactionKeyPrefix))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var transaction types.Transaction
		if err := k.cdc.Unmarshal(iterator.Value(), &transaction); err != nil {
			ctx.Logger().Error("Failed to unmarshal transaction, skipping",
				"error", err,
				"key", string(iterator.Key()))
			continue
		}

		if transaction.Status == "pending" && transaction.CreatedAt < cutoff {
			// Cancel timed out transaction
			k.CancelTransaction(ctx, transaction.Hash, "system", "transaction timeout")
		}
	}
}

// EndBlocker processes end of block operations
func (k Keeper) EndBlocker(ctx sdk.Context) []abci.ValidatorUpdate {
	// Update current height in statistics
	stats := k.GetTransactionStats(ctx)
	stats.CurrentHeight = ctx.BlockHeight()
	k.SetTransactionStats(ctx, stats)

	return []abci.ValidatorUpdate{}
}
