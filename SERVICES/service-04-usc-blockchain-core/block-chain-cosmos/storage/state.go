package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// StateManager manages blockchain state operations
type StateManager struct {
	rocksdb *RocksDBManager
	cache   map[string][]byte
	mu      sync.RWMutex
}

// StateKey represents different types of state keys
type StateKey string

const (
	// Blockchain state keys
	StateKeyGenesis      StateKey = "genesis"
	StateKeyChainID      StateKey = "chain_id"
	StateKeyHeight       StateKey = "height"
	StateKeyAppHash      StateKey = "app_hash"
	StateKeyConsensus    StateKey = "consensus"
	StateKeyValidators   StateKey = "validators"
	StateKeyStaking      StateKey = "staking"
	StateKeyBank         StateKey = "bank"
	StateKeyAuth         StateKey = "auth"
	StateKeyGov          StateKey = "gov"
	StateKeyParams       StateKey = "params"
	StateKeyCrisis       StateKey = "crisis"
	StateKeySlashing     StateKey = "slashing"
	StateKeyMint         StateKey = "mint"
	StateKeyDistribution StateKey = "distribution"

	// USC module state keys
	StateKeyUSC          StateKey = "usc"
	StateKeyUSCBalances  StateKey = "usc_balances"
	StateKeyUSCTransfers StateKey = "usc_transfers"
	StateKeyUSCParams    StateKey = "usc_params"
)

// StateData represents the structure of state data
type StateData struct {
	Key       StateKey  `json:"key"`
	Value     []byte    `json:"value"`
	Height    int64     `json:"height"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Hash      string    `json:"hash"`
}

// NewStateManager creates a new state manager
func NewStateManager(rocksdb *RocksDBManager) *StateManager {
	return &StateManager{
		rocksdb: rocksdb,
		cache:   make(map[string][]byte),
	}
}

// SetState stores state data
func (sm *StateManager) SetState(ctx context.Context, key StateKey, value []byte, height int64) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	stateData := StateData{
		Key:       key,
		Value:     value,
		Height:    height,
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Hash:      fmt.Sprintf("%x", value), // Simplified hash
	}

	data, err := json.Marshal(stateData)
	if err != nil {
		return fmt.Errorf("failed to marshal state data: %w", err)
	}

	// Store in RocksDB
	err = sm.rocksdb.Set(ctx, []byte(string(key)), data)
	if err != nil {
		return fmt.Errorf("failed to store state in RocksDB: %w", err)
	}

	// Update cache
	sm.cache[string(key)] = data

	return nil
}

// GetState retrieves state data
func (sm *StateManager) GetState(ctx context.Context, key StateKey) (*StateData, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// Check cache first
	if data, exists := sm.cache[string(key)]; exists {
		var stateData StateData
		if err := json.Unmarshal(data, &stateData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cached state data: %w", err)
		}
		return &stateData, nil
	}

	// Get from RocksDB
	data, err := sm.rocksdb.Get(ctx, []byte(string(key)))
	if err != nil {
		return nil, fmt.Errorf("failed to get state from RocksDB: %w", err)
	}

	var stateData StateData
	if err := json.Unmarshal(data, &stateData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal state data: %w", err)
	}

	// Update cache
	sm.cache[string(key)] = data

	return &stateData, nil
}

// DeleteState removes state data
func (sm *StateManager) DeleteState(ctx context.Context, key StateKey) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Delete from RocksDB
	err := sm.rocksdb.Delete(ctx, []byte(string(key)))
	if err != nil {
		return fmt.Errorf("failed to delete state from RocksDB: %w", err)
	}

	// Remove from cache
	delete(sm.cache, string(key))

	return nil
}

// HasState checks if state exists
func (sm *StateManager) HasState(ctx context.Context, key StateKey) (bool, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// Check cache first
	if _, exists := sm.cache[string(key)]; exists {
		return true, nil
	}

	// Check RocksDB
	exists, err := sm.rocksdb.Has(ctx, []byte(string(key)))
	if err != nil {
		return false, fmt.Errorf("failed to check state in RocksDB: %w", err)
	}

	return exists, nil
}

// GetStateByHeight retrieves state data at a specific height
func (sm *StateManager) GetStateByHeight(ctx context.Context, key StateKey, height int64) (*StateData, error) {
	// For now, return the latest state
	// In a full implementation, this would query historical state
	return sm.GetState(ctx, key)
}

// ListStates returns all states for a given prefix
func (sm *StateManager) ListStates(ctx context.Context, prefix StateKey) ([]*StateData, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	var states []*StateData
	startKey := []byte(string(prefix))
	endKey := append(startKey, 0xFF) // End key for range scan

	iter, err := sm.rocksdb.Iterator(ctx, startKey, endKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create iterator: %w", err)
	}
	defer iter.Close()

	for iter.Valid() {
		value := iter.Value()

		var stateData StateData
		if err := json.Unmarshal(value, &stateData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal state data: %w", err)
		}

		states = append(states, &stateData)
		iter.Next()
	}

	return states, nil
}

// GetStateHeight returns the current blockchain height
func (sm *StateManager) GetStateHeight(ctx context.Context) (int64, error) {
	state, err := sm.GetState(ctx, StateKeyHeight)
	if err != nil {
		return 0, fmt.Errorf("failed to get blockchain height: %w", err)
	}

	var height int64
	if err := json.Unmarshal(state.Value, &height); err != nil {
		return 0, fmt.Errorf("failed to unmarshal height: %w", err)
	}

	return height, nil
}

// SetStateHeight updates the blockchain height
func (sm *StateManager) SetStateHeight(ctx context.Context, height int64) error {
	heightBytes, err := json.Marshal(height)
	if err != nil {
		return fmt.Errorf("failed to marshal height: %w", err)
	}

	return sm.SetState(ctx, StateKeyHeight, heightBytes, height)
}

// GetChainID returns the chain ID
func (sm *StateManager) GetChainID(ctx context.Context) (string, error) {
	state, err := sm.GetState(ctx, StateKeyChainID)
	if err != nil {
		return "", fmt.Errorf("failed to get chain ID: %w", err)
	}

	var chainID string
	if err := json.Unmarshal(state.Value, &chainID); err != nil {
		return "", fmt.Errorf("failed to unmarshal chain ID: %w", err)
	}

	return chainID, nil
}

// SetChainID sets the chain ID
func (sm *StateManager) SetChainID(ctx context.Context, chainID string) error {
	chainIDBytes, err := json.Marshal(chainID)
	if err != nil {
		return fmt.Errorf("failed to marshal chain ID: %w", err)
	}

	return sm.SetState(ctx, StateKeyChainID, chainIDBytes, 0)
}

// GetAppHash returns the application hash
func (sm *StateManager) GetAppHash(ctx context.Context) (string, error) {
	state, err := sm.GetState(ctx, StateKeyAppHash)
	if err != nil {
		return "", fmt.Errorf("failed to get app hash: %w", err)
	}

	var appHash string
	if err := json.Unmarshal(state.Value, &appHash); err != nil {
		return "", fmt.Errorf("failed to unmarshal app hash: %w", err)
	}

	return appHash, nil
}

// SetAppHash sets the application hash
func (sm *StateManager) SetAppHash(ctx context.Context, appHash string) error {
	appHashBytes, err := json.Marshal(appHash)
	if err != nil {
		return fmt.Errorf("failed to marshal app hash: %w", err)
	}

	height, err := sm.GetStateHeight(ctx)
	if err != nil {
		height = 0
	}

	return sm.SetState(ctx, StateKeyAppHash, appHashBytes, height)
}

// ClearCache clears the state cache
func (sm *StateManager) ClearCache() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.cache = make(map[string][]byte)
}

// GetCacheStats returns cache statistics
func (sm *StateManager) GetCacheStats() map[string]interface{} {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return map[string]interface{}{
		"cache_size": len(sm.cache),
		"cache_keys": func() []string {
			keys := make([]string, 0, len(sm.cache))
			for k := range sm.cache {
				keys = append(keys, k)
			}
			return keys
		}(),
	}
}
