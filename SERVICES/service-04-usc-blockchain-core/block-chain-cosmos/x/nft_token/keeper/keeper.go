package keeper

import (
	"encoding/json"
	"fmt"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/nft_token/types"
)

// Keeper manages the nft module state
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
}

// NewKeeper creates a new NFT keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
	}
}

// GetNFT returns an NFT by its ID
func (k Keeper) GetNFT(ctx sdk.Context, id string) (types.NFT, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NFTKey(id))
	if bz == nil {
		return types.NFT{}, fmt.Errorf("NFT with ID %s not found", id)
	}

	var nft types.NFT
	if err := json.Unmarshal(bz, &nft); err != nil {
		return types.NFT{}, fmt.Errorf("failed to unmarshal NFT: %w", err)
	}

	return nft, nil
}

// SetNFT sets an NFT
func (k Keeper) SetNFT(ctx sdk.Context, nft types.NFT) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(nft)
	if err != nil {
		return fmt.Errorf("failed to marshal NFT: %w", err)
	}
	store.Set(types.NFTKey(nft.ID), bz)
	return nil
}

// GetAllNFTs returns all NFTs
func (k Keeper) GetAllNFTs(ctx sdk.Context) []types.NFT {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.NFTKeyPrefix)
	defer iterator.Close()

	var nfts []types.NFT
	for ; iterator.Valid(); iterator.Next() {
		var nft types.NFT
		if err := json.Unmarshal(iterator.Value(), &nft); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		nfts = append(nfts, nft)
	}
	return nfts
}

// GetCollection returns a collection by its ID
func (k Keeper) GetCollection(ctx sdk.Context, id string) (types.Collection, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.CollectionKey(id))
	if bz == nil {
		return types.Collection{}, fmt.Errorf("collection with ID %s not found", id)
	}

	var collection types.Collection
	if err := json.Unmarshal(bz, &collection); err != nil {
		return types.Collection{}, fmt.Errorf("failed to unmarshal collection: %w", err)
	}

	return collection, nil
}

// SetCollection sets a collection
func (k Keeper) SetCollection(ctx sdk.Context, collection types.Collection) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(collection)
	if err != nil {
		return fmt.Errorf("failed to marshal collection: %w", err)
	}
	store.Set(types.CollectionKey(collection.ID), bz)
	return nil
}

// GetAllCollections returns all collections
func (k Keeper) GetAllCollections(ctx sdk.Context) []types.Collection {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.CollectionKeyPrefix)
	defer iterator.Close()

	var collections []types.Collection
	for ; iterator.Valid(); iterator.Next() {
		var collection types.Collection
		if err := json.Unmarshal(iterator.Value(), &collection); err != nil {
			// Log error or handle it as appropriate
			continue
		}
		collections = append(collections, collection)
	}
	return collections
}

// GetParams returns the NFT module's parameters
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return types.DefaultParams()
	}

	var params types.Params
	if err := json.Unmarshal(bz, &params); err != nil {
		return types.DefaultParams()
	}

	return params
}

// SetParams sets the NFT module's parameters
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(params)
	if err != nil {
		// Handle error appropriately
		return
	}
	store.Set(types.ParamsKey, bz)
}

// DeleteNFT deletes an NFT by its ID (used for burning)
func (k Keeper) DeleteNFT(ctx sdk.Context, id string) error {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.NFTKey(id))
	return nil
}
