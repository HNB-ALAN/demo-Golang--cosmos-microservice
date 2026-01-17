package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// ConsensusManager manages consensus state and validator operations
type ConsensusManager struct {
	rocksdb *RocksDBManager
	state   *StateManager
	mu      sync.RWMutex
}

// Validator represents a blockchain validator
type Validator struct {
	Address         string               `json:"address"`
	PubKey          string               `json:"pub_key"`
	VotingPower     int64                `json:"voting_power"`
	Commission      string               `json:"commission"`
	Jailed          bool                 `json:"jailed"`
	Status          string               `json:"status"`
	DelegatorShares string               `json:"delegator_shares"`
	UnbondingHeight int64                `json:"unbonding_height"`
	UnbondingTime   time.Time            `json:"unbonding_time"`
	Description     ValidatorDescription `json:"description"`
}

// ValidatorDescription contains validator metadata
type ValidatorDescription struct {
	Moniker         string `json:"moniker"`
	Identity        string `json:"identity"`
	Website         string `json:"website"`
	SecurityContact string `json:"security_contact"`
	Details         string `json:"details"`
}

// ConsensusState represents the current consensus state
type ConsensusState struct {
	Height           int64       `json:"height"`
	Round            int32       `json:"round"`
	Step             int32       `json:"step"`
	Validators       []Validator `json:"validators"`
	Proposer         string      `json:"proposer"`
	LastCommit       string      `json:"last_commit"`
	NextValidators   []Validator `json:"next_validators"`
	TotalVotingPower int64       `json:"total_voting_power"`
	Timestamp        time.Time   `json:"timestamp"`
	AppHash          string      `json:"app_hash"`
	ConsensusHash    string      `json:"consensus_hash"`
	LastResultsHash  string      `json:"last_results_hash"`
	EvidenceHash     string      `json:"evidence_hash"`
}

// ValidatorSet represents a set of validators
type ValidatorSet struct {
	Validators []Validator `json:"validators"`
	Proposer   *Validator  `json:"proposer"`
	Total      int64       `json:"total"`
	Height     int64       `json:"height"`
}

// NewConsensusManager creates a new consensus manager
func NewConsensusManager(rocksdb *RocksDBManager, state *StateManager) *ConsensusManager {
	return &ConsensusManager{
		rocksdb: rocksdb,
		state:   state,
	}
}

// SetValidator stores a validator
func (cm *ConsensusManager) SetValidator(ctx context.Context, validator Validator) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	key := fmt.Sprintf("validator:%s", validator.Address)
	data, err := json.Marshal(validator)
	if err != nil {
		return fmt.Errorf("failed to marshal validator: %w", err)
	}

	err = cm.rocksdb.Set(ctx, []byte(key), data)
	if err != nil {
		return fmt.Errorf("failed to store validator: %w", err)
	}

	return nil
}

// GetValidator retrieves a validator by address
func (cm *ConsensusManager) GetValidator(ctx context.Context, address string) (*Validator, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	key := fmt.Sprintf("validator:%s", address)
	data, err := cm.rocksdb.Get(ctx, []byte(key))
	if err != nil {
		return nil, fmt.Errorf("failed to get validator: %w", err)
	}

	var validator Validator
	if err := json.Unmarshal(data, &validator); err != nil {
		return nil, fmt.Errorf("failed to unmarshal validator: %w", err)
	}

	return &validator, nil
}

// DeleteValidator removes a validator
func (cm *ConsensusManager) DeleteValidator(ctx context.Context, address string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	key := fmt.Sprintf("validator:%s", address)
	err := cm.rocksdb.Delete(ctx, []byte(key))
	if err != nil {
		return fmt.Errorf("failed to delete validator: %w", err)
	}

	return nil
}

// ListValidators returns all validators
func (cm *ConsensusManager) ListValidators(ctx context.Context) ([]Validator, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var validators []Validator
	startKey := []byte("validator:")
	endKey := []byte("validator:~") // End key for range scan

	iter, err := cm.rocksdb.Iterator(ctx, startKey, endKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create iterator: %w", err)
	}
	defer iter.Close()

	for iter.Valid() {
		var validator Validator
		if err := json.Unmarshal(iter.Value(), &validator); err != nil {
			return nil, fmt.Errorf("failed to unmarshal validator: %w", err)
		}

		validators = append(validators, validator)
		iter.Next()
	}

	return validators, nil
}

// SetConsensusState stores the current consensus state
func (cm *ConsensusManager) SetConsensusState(ctx context.Context, state ConsensusState) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal consensus state: %w", err)
	}

	key := fmt.Sprintf("consensus_state:%d", state.Height)
	err = cm.rocksdb.Set(ctx, []byte(key), data)
	if err != nil {
		return fmt.Errorf("failed to store consensus state: %w", err)
	}

	// Also update the current consensus state
	err = cm.rocksdb.Set(ctx, []byte("current_consensus_state"), data)
	if err != nil {
		return fmt.Errorf("failed to store current consensus state: %w", err)
	}

	return nil
}

// GetConsensusState retrieves the current consensus state
func (cm *ConsensusManager) GetConsensusState(ctx context.Context) (*ConsensusState, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	data, err := cm.rocksdb.Get(ctx, []byte("current_consensus_state"))
	if err != nil {
		return nil, fmt.Errorf("failed to get consensus state: %w", err)
	}

	var state ConsensusState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal consensus state: %w", err)
	}

	return &state, nil
}

// GetConsensusStateByHeight retrieves consensus state at a specific height
func (cm *ConsensusManager) GetConsensusStateByHeight(ctx context.Context, height int64) (*ConsensusState, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	key := fmt.Sprintf("consensus_state:%d", height)
	data, err := cm.rocksdb.Get(ctx, []byte(key))
	if err != nil {
		return nil, fmt.Errorf("failed to get consensus state at height %d: %w", height, err)
	}

	var state ConsensusState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal consensus state: %w", err)
	}

	return &state, nil
}

// SetValidatorSet stores a validator set
func (cm *ConsensusManager) SetValidatorSet(ctx context.Context, height int64, validatorSet ValidatorSet) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	data, err := json.Marshal(validatorSet)
	if err != nil {
		return fmt.Errorf("failed to marshal validator set: %w", err)
	}

	key := fmt.Sprintf("validator_set:%d", height)
	err = cm.rocksdb.Set(ctx, []byte(key), data)
	if err != nil {
		return fmt.Errorf("failed to store validator set: %w", err)
	}

	return nil
}

// GetValidatorSet retrieves a validator set by height
func (cm *ConsensusManager) GetValidatorSet(ctx context.Context, height int64) (*ValidatorSet, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	key := fmt.Sprintf("validator_set:%d", height)
	data, err := cm.rocksdb.Get(ctx, []byte(key))
	if err != nil {
		return nil, fmt.Errorf("failed to get validator set at height %d: %w", height, err)
	}

	var validatorSet ValidatorSet
	if err := json.Unmarshal(data, &validatorSet); err != nil {
		return nil, fmt.Errorf("failed to unmarshal validator set: %w", err)
	}

	return &validatorSet, nil
}

// GetCurrentValidatorSet retrieves the current validator set
func (cm *ConsensusManager) GetCurrentValidatorSet(ctx context.Context) (*ValidatorSet, error) {
	// Get the current height
	height, err := cm.state.GetStateHeight(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current height: %w", err)
	}

	return cm.GetValidatorSet(ctx, height)
}

// UpdateValidatorVotingPower updates a validator's voting power
func (cm *ConsensusManager) UpdateValidatorVotingPower(ctx context.Context, address string, votingPower int64) error {
	validator, err := cm.GetValidator(ctx, address)
	if err != nil {
		return fmt.Errorf("failed to get validator: %w", err)
	}

	validator.VotingPower = votingPower
	return cm.SetValidator(ctx, *validator)
}

// JailValidator jails a validator
func (cm *ConsensusManager) JailValidator(ctx context.Context, address string) error {
	validator, err := cm.GetValidator(ctx, address)
	if err != nil {
		return fmt.Errorf("failed to get validator: %w", err)
	}

	validator.Jailed = true
	validator.Status = "jailed"
	return cm.SetValidator(ctx, *validator)
}

// UnjailValidator unjails a validator
func (cm *ConsensusManager) UnjailValidator(ctx context.Context, address string) error {
	validator, err := cm.GetValidator(ctx, address)
	if err != nil {
		return fmt.Errorf("failed to get validator: %w", err)
	}

	validator.Jailed = false
	validator.Status = "bonded"
	return cm.SetValidator(ctx, *validator)
}

// GetTotalVotingPower returns the total voting power
func (cm *ConsensusManager) GetTotalVotingPower(ctx context.Context) (int64, error) {
	validators, err := cm.ListValidators(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to list validators: %w", err)
	}

	var total int64
	for _, validator := range validators {
		if !validator.Jailed {
			total += validator.VotingPower
		}
	}

	return total, nil
}

// GetActiveValidators returns only active (unjailed) validators
func (cm *ConsensusManager) GetActiveValidators(ctx context.Context) ([]Validator, error) {
	validators, err := cm.ListValidators(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list validators: %w", err)
	}

	var activeValidators []Validator
	for _, validator := range validators {
		if !validator.Jailed && validator.Status == "bonded" {
			activeValidators = append(activeValidators, validator)
		}
	}

	return activeValidators, nil
}

// GetConsensusStats returns consensus statistics
func (cm *ConsensusManager) GetConsensusStats(ctx context.Context) (map[string]interface{}, error) {
	validators, err := cm.ListValidators(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list validators: %w", err)
	}

	activeValidators, err := cm.GetActiveValidators(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active validators: %w", err)
	}

	totalVotingPower, err := cm.GetTotalVotingPower(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get total voting power: %w", err)
	}

	return map[string]interface{}{
		"total_validators":    len(validators),
		"active_validators":   len(activeValidators),
		"jailed_validators":   len(validators) - len(activeValidators),
		"total_voting_power":  totalVotingPower,
		"average_commission":  cm.calculateAverageCommission(validators),
		"consensus_threshold": int64(float64(totalVotingPower) * 0.67), // 2/3 threshold
	}, nil
}

// calculateAverageCommission calculates the average commission rate
func (cm *ConsensusManager) calculateAverageCommission(validators []Validator) float64 {
	if len(validators) == 0 {
		return 0.0
	}

	var total float64
	for range validators {
		// Parse commission string to float64
		// This is simplified - in production, use proper decimal parsing
		total += 0.1 // Placeholder
	}

	return total / float64(len(validators))
}
