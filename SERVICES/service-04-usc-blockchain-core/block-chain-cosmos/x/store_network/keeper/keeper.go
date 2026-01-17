package keeper

import (
	"encoding/json"
	"fmt"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	networktypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/store_network/types"
)

// Keeper manages the store module state
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
}

// NewKeeper creates a new store keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
	}
}

// GetStoredData returns stored data by its ID
func (k Keeper) GetStoredData(ctx sdk.Context, id string) (networktypes.StoredData, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(networktypes.StoredDataKey(id))
	if bz == nil {
		return networktypes.StoredData{}, fmt.Errorf("stored data with ID %s not found", id)
	}

	var data networktypes.StoredData
	if err := json.Unmarshal(bz, &data); err != nil {
		return networktypes.StoredData{}, fmt.Errorf("failed to unmarshal stored data: %w", err)
	}

	return data, nil
}

// SetStoredData sets stored data
func (k Keeper) SetStoredData(ctx sdk.Context, data networktypes.StoredData) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal stored data: %w", err)
	}
	store.Set(networktypes.StoredDataKey(data.ID), bz)
	return nil
}

// GetAllStoredData returns all stored data
func (k Keeper) GetAllStoredData(ctx sdk.Context) []networktypes.StoredData {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, networktypes.StoredDataKeyPrefix)
	defer iterator.Close()

	var dataList []networktypes.StoredData
	for ; iterator.Valid(); iterator.Next() {
		var data networktypes.StoredData
		if err := json.Unmarshal(iterator.Value(), &data); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		dataList = append(dataList, data)
	}
	return dataList
}

// GetStore returns a store by its ID
func (k Keeper) GetStore(ctx sdk.Context, id string) (networktypes.Store, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(networktypes.StoreKey(id))
	if bz == nil {
		return networktypes.Store{}, fmt.Errorf("store with ID %s not found", id)
	}

	var storeData networktypes.Store
	if err := json.Unmarshal(bz, &storeData); err != nil {
		return networktypes.Store{}, fmt.Errorf("failed to unmarshal store: %w", err)
	}

	return storeData, nil
}

// SetStore sets a store
func (k Keeper) SetStore(ctx sdk.Context, storeData networktypes.Store) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(storeData)
	if err != nil {
		return fmt.Errorf("failed to marshal store: %w", err)
	}
	store.Set(networktypes.StoreKey(storeData.ID), bz)
	return nil
}

// GetAllStores returns all stores
func (k Keeper) GetAllStores(ctx sdk.Context) []networktypes.Store {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, networktypes.StoreKeyPrefix)
	defer iterator.Close()

	var stores []networktypes.Store
	for ; iterator.Valid(); iterator.Next() {
		var storeData networktypes.Store
		if err := json.Unmarshal(iterator.Value(), &storeData); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		stores = append(stores, storeData)
	}
	return stores
}

// GetBackup returns a backup by its ID
func (k Keeper) GetBackup(ctx sdk.Context, id string) (networktypes.Backup, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(networktypes.BackupKey(id))
	if bz == nil {
		return networktypes.Backup{}, fmt.Errorf("backup with ID %s not found", id)
	}

	var backup networktypes.Backup
	if err := json.Unmarshal(bz, &backup); err != nil {
		return networktypes.Backup{}, fmt.Errorf("failed to unmarshal backup: %w", err)
	}

	return backup, nil
}

// SetBackup sets a backup
func (k Keeper) SetBackup(ctx sdk.Context, backup networktypes.Backup) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(backup)
	if err != nil {
		return fmt.Errorf("failed to marshal backup: %w", err)
	}
	store.Set(networktypes.BackupKey(backup.ID), bz)
	return nil
}

// GetAllBackups returns all backups
func (k Keeper) GetAllBackups(ctx sdk.Context) []networktypes.Backup {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, networktypes.BackupKeyPrefix)
	defer iterator.Close()

	var backups []networktypes.Backup
	for ; iterator.Valid(); iterator.Next() {
		var backup networktypes.Backup
		if err := json.Unmarshal(iterator.Value(), &backup); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		backups = append(backups, backup)
	}
	return backups
}

// GetRestore returns a restore by its ID
func (k Keeper) GetRestore(ctx sdk.Context, id string) (networktypes.Restore, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(networktypes.RestoreKey(id))
	if bz == nil {
		return networktypes.Restore{}, fmt.Errorf("restore with ID %s not found", id)
	}

	var restore networktypes.Restore
	if err := json.Unmarshal(bz, &restore); err != nil {
		return networktypes.Restore{}, fmt.Errorf("failed to unmarshal restore: %w", err)
	}

	return restore, nil
}

// SetRestore sets a restore
func (k Keeper) SetRestore(ctx sdk.Context, restore networktypes.Restore) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(restore)
	if err != nil {
		return fmt.Errorf("failed to marshal restore: %w", err)
	}
	store.Set(networktypes.RestoreKey(restore.ID), bz)
	return nil
}

// GetAllRestores returns all restores
func (k Keeper) GetAllRestores(ctx sdk.Context) []networktypes.Restore {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, networktypes.RestoreKeyPrefix)
	defer iterator.Close()

	var restores []networktypes.Restore
	for ; iterator.Valid(); iterator.Next() {
		var restore networktypes.Restore
		if err := json.Unmarshal(iterator.Value(), &restore); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		restores = append(restores, restore)
	}
	return restores
}

// GetStoreIndex returns a store index by its ID
func (k Keeper) GetStoreIndex(ctx sdk.Context, id string) (networktypes.StoreIndex, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(networktypes.StoreIndexKey(id))
	if bz == nil {
		return networktypes.StoreIndex{}, fmt.Errorf("store index with ID %s not found", id)
	}

	var index networktypes.StoreIndex
	if err := json.Unmarshal(bz, &index); err != nil {
		return networktypes.StoreIndex{}, fmt.Errorf("failed to unmarshal store index: %w", err)
	}

	return index, nil
}

// SetStoreIndex sets a store index
func (k Keeper) SetStoreIndex(ctx sdk.Context, index networktypes.StoreIndex) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(index)
	if err != nil {
		return fmt.Errorf("failed to marshal store index: %w", err)
	}
	store.Set(networktypes.StoreIndexKey(index.ID), bz)
	return nil
}

// GetAllStoreIndexes returns all store indexes
func (k Keeper) GetAllStoreIndexes(ctx sdk.Context) []networktypes.StoreIndex {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, networktypes.StoreIndexKeyPrefix)
	defer iterator.Close()

	var indexes []networktypes.StoreIndex
	for ; iterator.Valid(); iterator.Next() {
		var index networktypes.StoreIndex
		if err := json.Unmarshal(iterator.Value(), &index); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		indexes = append(indexes, index)
	}
	return indexes
}

// GetStoreQuery returns a store query by its ID
func (k Keeper) GetStoreQuery(ctx sdk.Context, id string) (networktypes.StoreQuery, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(networktypes.StoreQueryKey(id))
	if bz == nil {
		return networktypes.StoreQuery{}, fmt.Errorf("store query with ID %s not found", id)
	}

	var query networktypes.StoreQuery
	if err := json.Unmarshal(bz, &query); err != nil {
		return networktypes.StoreQuery{}, fmt.Errorf("failed to unmarshal store query: %w", err)
	}

	return query, nil
}

// SetStoreQuery sets a store query
func (k Keeper) SetStoreQuery(ctx sdk.Context, query networktypes.StoreQuery) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(query)
	if err != nil {
		return fmt.Errorf("failed to marshal store query: %w", err)
	}
	store.Set(networktypes.StoreQueryKey(query.ID), bz)
	return nil
}

// GetAllStoreQueries returns all store queries
func (k Keeper) GetAllStoreQueries(ctx sdk.Context) []networktypes.StoreQuery {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, networktypes.StoreQueryKeyPrefix)
	defer iterator.Close()

	var queries []networktypes.StoreQuery
	for ; iterator.Valid(); iterator.Next() {
		var query networktypes.StoreQuery
		if err := json.Unmarshal(iterator.Value(), &query); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		queries = append(queries, query)
	}
	return queries
}

// GetStoreTransaction returns a store transaction by its ID
func (k Keeper) GetStoreTransaction(ctx sdk.Context, id string) (networktypes.StoreTransaction, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(networktypes.TransactionKey(id))
	if bz == nil {
		return networktypes.StoreTransaction{}, fmt.Errorf("store transaction with ID %s not found", id)
	}

	var transaction networktypes.StoreTransaction
	if err := json.Unmarshal(bz, &transaction); err != nil {
		return networktypes.StoreTransaction{}, fmt.Errorf("failed to unmarshal store transaction: %w", err)
	}

	return transaction, nil
}

// SetStoreTransaction sets a store transaction
func (k Keeper) SetStoreTransaction(ctx sdk.Context, transaction networktypes.StoreTransaction) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to marshal store transaction: %w", err)
	}
	store.Set(networktypes.TransactionKey(transaction.ID), bz)
	return nil
}

// GetAllStoreTransactions returns all store transactions
func (k Keeper) GetAllStoreTransactions(ctx sdk.Context) []networktypes.StoreTransaction {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, networktypes.TransactionKeyPrefix)
	defer iterator.Close()

	var transactions []networktypes.StoreTransaction
	for ; iterator.Valid(); iterator.Next() {
		var transaction networktypes.StoreTransaction
		if err := json.Unmarshal(iterator.Value(), &transaction); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		transactions = append(transactions, transaction)
	}
	return transactions
}

// GetParams returns the store module's parameters
func (k Keeper) GetParams(ctx sdk.Context) networktypes.Params {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(networktypes.ParamsKey)
	if bz == nil {
		return networktypes.DefaultParams()
	}

	var params networktypes.Params
	if err := json.Unmarshal(bz, &params); err != nil {
		return networktypes.DefaultParams()
	}

	return params
}

// SetParams sets the store module's parameters
func (k Keeper) SetParams(ctx sdk.Context, params networktypes.Params) {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(params)
	if err != nil {
		// Handle error appropriately
		return
	}
	store.Set(networktypes.ParamsKey, bz)
}
