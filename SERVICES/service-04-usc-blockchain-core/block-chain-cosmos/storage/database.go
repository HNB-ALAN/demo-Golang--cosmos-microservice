package storage

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// DatabaseManager manages blockchain-specific database connections
// This is ONLY for blockchain layer - PostgreSQL is handled by application layer
type DatabaseManager struct {
	rocksdb   *RocksDBManager
	state     *StateManager
	consensus *ConsensusManager
	mu        sync.RWMutex
	closed    bool
}

// DatabaseConfig contains configuration for blockchain database connections
type DatabaseConfig struct {
	// RocksDB configuration
	RocksDB RocksDBConfig `json:"rocksdb"`

	// Note: PostgreSQL is handled by application layer, not blockchain layer

	// Connection settings
	MaxConnections     int           `json:"max_connections"`
	MaxIdleConns       int           `json:"max_idle_conns"`
	ConnMaxLifetime    time.Duration `json:"conn_max_lifetime"`
	ConnMaxIdleTime    time.Duration `json:"conn_max_idle_time"`
	QueryTimeout       time.Duration `json:"query_timeout"`
	TransactionTimeout time.Duration `json:"transaction_timeout"`
}

// DefaultDatabaseConfig returns the default database configuration
func DefaultDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		RocksDB: DefaultRocksDBConfig(),
		// Note: PostgreSQL is handled by application layer
		MaxConnections:     100,
		MaxIdleConns:       10,
		ConnMaxLifetime:    30 * time.Minute,
		ConnMaxIdleTime:    5 * time.Minute,
		QueryTimeout:       30 * time.Second,
		TransactionTimeout: 60 * time.Second,
	}
}

// NewDatabaseManager creates a new blockchain database manager
func NewDatabaseManager(config DatabaseConfig) (*DatabaseManager, error) {
	// Initialize RocksDB
	rocksdb, err := NewRocksDBManager(config.RocksDB)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize RocksDB: %w", err)
	}

	// Initialize state manager
	state := NewStateManager(rocksdb)

	// Initialize consensus manager
	consensus := NewConsensusManager(rocksdb, state)

	manager := &DatabaseManager{
		rocksdb:   rocksdb,
		state:     state,
		consensus: consensus,
	}

	return manager, nil
}

// GetRocksDB returns the RocksDB manager
func (dm *DatabaseManager) GetRocksDB() *RocksDBManager {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	return dm.rocksdb
}

// GetState returns the state manager
func (dm *DatabaseManager) GetState() *StateManager {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	return dm.state
}

// GetConsensus returns the consensus manager
func (dm *DatabaseManager) GetConsensus() *ConsensusManager {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	return dm.consensus
}

// Health checks the health of blockchain databases
func (dm *DatabaseManager) Health(ctx context.Context) map[string]interface{} {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	health := make(map[string]interface{})

	// Check RocksDB health
	if dm.rocksdb != nil {
		health["rocksdb"] = map[string]interface{}{
			"status": "healthy",
			"type":   "blockchain_storage",
		}
	} else {
		health["rocksdb"] = map[string]interface{}{
			"status": "unhealthy",
			"error":  "rocksdb not initialized",
		}
	}

	// Check state manager health
	if dm.state != nil {
		health["state"] = map[string]interface{}{
			"status": "healthy",
			"type":   "blockchain_state",
		}
	} else {
		health["state"] = map[string]interface{}{
			"status": "unhealthy",
			"error":  "state manager not initialized",
		}
	}

	// Check consensus manager health
	if dm.consensus != nil {
		health["consensus"] = map[string]interface{}{
			"status": "healthy",
			"type":   "blockchain_consensus",
		}
	} else {
		health["consensus"] = map[string]interface{}{
			"status": "unhealthy",
			"error":  "consensus manager not initialized",
		}
	}

	return health
}

// Stats returns statistics for blockchain databases
func (dm *DatabaseManager) Stats() map[string]interface{} {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	stats := make(map[string]interface{})

	// RocksDB stats
	if dm.rocksdb != nil {
		stats["rocksdb"] = map[string]interface{}{
			"type":   "blockchain_storage",
			"status": "active",
		}
	}

	// State manager stats
	if dm.state != nil {
		stats["state"] = map[string]interface{}{
			"type":   "blockchain_state",
			"status": "active",
		}
	}

	// Consensus manager stats
	if dm.consensus != nil {
		stats["consensus"] = map[string]interface{}{
			"type":   "blockchain_consensus",
			"status": "active",
		}
	}

	return stats
}

// Close closes all blockchain database connections
func (dm *DatabaseManager) Close() error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if dm.closed {
		return nil
	}

	var err error

	// Close RocksDB
	if dm.rocksdb != nil {
		if closeErr := dm.rocksdb.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close RocksDB: %w", closeErr)
		}
	}

	dm.closed = true
	return err
}
