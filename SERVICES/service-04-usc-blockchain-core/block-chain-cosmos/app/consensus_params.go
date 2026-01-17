package app

import (
	"context"
	"fmt"

	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cosmosdb "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// ============================================================================
// CONSENSUS PARAM STORE
// ============================================================================

// ConsensusParamStore implements baseapp.ParamStore interface for storing consensus params
type ConsensusParamStore struct {
	subspace paramtypes.Subspace
	db       cosmosdb.DB
	key      []byte
}

// NewConsensusParamStore creates a new ConsensusParamStore
func NewConsensusParamStore(subspace paramtypes.Subspace, db cosmosdb.DB, storeKey storetypes.StoreKey) baseapp.ParamStore {
	return &ConsensusParamStore{
		subspace: subspace,
		db:       db,
		key:      []byte("consensus_params"), // Key for storing consensus params
	}
}

// Get retrieves consensus params from database
func (ps *ConsensusParamStore) Get(ctx context.Context) (cmtproto.ConsensusParams, error) {
	var params cmtproto.ConsensusParams

	// Try to get from database
	data, err := ps.db.Get(ps.key)
	if err != nil || data == nil {
		// Return default consensus params if not found
		return cmtproto.ConsensusParams{
			Block: &cmtproto.BlockParams{
				MaxBytes: 1048576, // 1MB
				MaxGas:   10000000,
			},
			Evidence: &cmtproto.EvidenceParams{
				MaxAgeNumBlocks: 100000,
				MaxAgeDuration:  172800000000000, // 2 days in nanoseconds
			},
			Validator: &cmtproto.ValidatorParams{
				PubKeyTypes: []string{"ed25519"},
			},
			Abci: &cmtproto.ABCIParams{
				VoteExtensionsEnableHeight: 0,
			},
		}, nil
	}

	// Unmarshal consensus params
	if err := params.Unmarshal(data); err != nil {
		return cmtproto.ConsensusParams{}, fmt.Errorf("failed to unmarshal consensus params: %w", err)
	}

	return params, nil
}

// Set stores consensus params to database
func (ps *ConsensusParamStore) Set(ctx context.Context, params cmtproto.ConsensusParams) error {
	data, err := params.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal consensus params: %w", err)
	}

	return ps.db.Set(ps.key, data)
}

// Has checks if consensus params exist in database
func (ps *ConsensusParamStore) Has(ctx context.Context) (bool, error) {
	return ps.db.Has(ps.key)
}
