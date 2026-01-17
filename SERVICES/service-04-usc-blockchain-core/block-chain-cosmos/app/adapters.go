package app

import (
	"context"

	"cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ============================================================================
// STORE ADAPTERS
// ============================================================================

// KVStoreServiceAdapter adapts storetypes.KVStoreKey to store.KVStoreService
// This is needed to bridge between legacy store keys and the new KVStoreService interface
type KVStoreServiceAdapter struct {
	key storetypes.StoreKey
}

// KVStoreAdapter adapts storetypes.KVStore to store.KVStore
type KVStoreAdapter struct {
	store storetypes.KVStore
}

// IteratorAdapter adapts storetypes.Iterator to store.Iterator
type IteratorAdapter struct {
	iter storetypes.Iterator
	err  error
}

// NewKVStoreService creates a new KVStoreService from a StoreKey
func NewKVStoreService(key storetypes.StoreKey) store.KVStoreService {
	return &KVStoreServiceAdapter{key: key}
}

// OpenKVStore opens a KV store using the context
func (k *KVStoreServiceAdapter) OpenKVStore(ctx context.Context) store.KVStore {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return &KVStoreAdapter{store: sdkCtx.KVStore(k.key)}
}

// Get retrieves a value by key
func (k *KVStoreAdapter) Get(key []byte) ([]byte, error) {
	return k.store.Get(key), nil
}

// Has checks if a key exists
func (k *KVStoreAdapter) Has(key []byte) (bool, error) {
	return k.store.Has(key), nil
}

// Set sets a key-value pair
func (k *KVStoreAdapter) Set(key, value []byte) error {
	k.store.Set(key, value)
	return nil
}

// Delete deletes a key
func (k *KVStoreAdapter) Delete(key []byte) error {
	k.store.Delete(key)
	return nil
}

// Iterator creates an iterator over a range of keys
func (k *KVStoreAdapter) Iterator(start, end []byte) (store.Iterator, error) {
	return &IteratorAdapter{iter: k.store.Iterator(start, end)}, nil
}

// ReverseIterator creates a reverse iterator over a range of keys
func (k *KVStoreAdapter) ReverseIterator(start, end []byte) (store.Iterator, error) {
	return &IteratorAdapter{iter: k.store.ReverseIterator(start, end)}, nil
}

// Domain returns the start and end keys of the iterator
func (i *IteratorAdapter) Domain() ([]byte, []byte) {
	return i.iter.Domain()
}

// Valid checks if the iterator is valid
func (i *IteratorAdapter) Valid() bool {
	return i.iter.Valid()
}

// Next advances the iterator
func (i *IteratorAdapter) Next() {
	i.iter.Next()
}

// Key returns the current key
func (i *IteratorAdapter) Key() []byte {
	return i.iter.Key()
}

// Value returns the current value
func (i *IteratorAdapter) Value() []byte {
	return i.iter.Value()
}

// Error returns any error that occurred during iteration
func (i *IteratorAdapter) Error() error {
	return i.err
}

// Close closes the iterator
func (i *IteratorAdapter) Close() error {
	return i.iter.Close()
}

// ============================================================================
// ADDRESS CODEC
// ============================================================================

// AddressCodec implements the address.Codec interface for SDK addresses
type AddressCodec struct{}

// StringToBytes decodes text to bytes
func (ac AddressCodec) StringToBytes(text string) ([]byte, error) {
	addr, err := sdk.AccAddressFromBech32(text)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

// BytesToString encodes bytes to text
func (ac AddressCodec) BytesToString(bz []byte) (string, error) {
	return sdk.AccAddress(bz).String(), nil
}

// NewAddressCodec creates a new AddressCodec for account addresses
func NewAddressCodec() address.Codec {
	return AddressCodec{}
}

// ValidatorAddressCodec implements the address.Codec interface for validator addresses
type ValidatorAddressCodec struct{}

// StringToBytes decodes validator address text to bytes
func (ac ValidatorAddressCodec) StringToBytes(text string) ([]byte, error) {
	addr, err := sdk.ValAddressFromBech32(text)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

// BytesToString encodes validator address bytes to text
func (ac ValidatorAddressCodec) BytesToString(bz []byte) (string, error) {
	return sdk.ValAddress(bz).String(), nil
}

// NewValidatorAddressCodec creates a new AddressCodec for validator addresses
func NewValidatorAddressCodec() address.Codec {
	return ValidatorAddressCodec{}
}

// ConsensusAddressCodec implements the address.Codec interface for consensus addresses
type ConsensusAddressCodec struct{}

// StringToBytes decodes consensus address text to bytes
func (ac ConsensusAddressCodec) StringToBytes(text string) ([]byte, error) {
	addr, err := sdk.ConsAddressFromBech32(text)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

// BytesToString encodes consensus address bytes to text
func (ac ConsensusAddressCodec) BytesToString(bz []byte) (string, error) {
	return sdk.ConsAddress(bz).String(), nil
}

// NewConsensusAddressCodec creates a new AddressCodec for consensus addresses
func NewConsensusAddressCodec() address.Codec {
	return ConsensusAddressCodec{}
}
