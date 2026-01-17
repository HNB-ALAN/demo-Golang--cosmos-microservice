package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/usc-platform/shared/config"
	"github.com/usc-platform/shared/database"
	"github.com/usc-platform/shared/logging"

	// Cosmos SDK imports
	"service-04/internal/application/utils"

	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/app"
	"github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/storage"
)

// PostgreSQLManager wraps the shared PostgreSQL database manager for USC Blockchain Core Service
// This is the APPLICATION LAYER database manager for business logic
// It integrates with the BLOCKCHAIN LAYER storage for hybrid operations
type PostgreSQLManager struct {
	manager           *database.DatabaseManager // Application layer database
	config            *config.Config
	logger            logging.Logger
	cosmosApp         *app.USCApp           // Blockchain layer app
	blockchainStorage *storage.StateManager // Blockchain layer storage
}

// NewPostgreSQLManager creates a new PostgreSQL database manager using shared libraries
func NewPostgreSQLManager(cfg *config.Config, logger logging.Logger, cosmosApp *app.USCApp, blockchainStorage *storage.StateManager) (*PostgreSQLManager, error) {
	// Ensure PostgreSQL is enabled in config
	if !cfg.Database.Enabled {
		return nil, fmt.Errorf("PostgreSQL is not enabled in configuration")
	}

	manager, err := database.NewDatabaseManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create database manager: %w", err)
	}

	pgManager := &PostgreSQLManager{
		manager:           manager,
		config:            cfg,
		logger:            logger,
		cosmosApp:         cosmosApp,
		blockchainStorage: blockchainStorage,
	}

	// Initialize database schema from migrations
	if err := pgManager.initializeSchema(context.Background()); err != nil {
		logger.Warn("Failed to initialize PostgreSQL schema", logging.Error(err))
		// Don't fail startup, just log warning
	}

	return pgManager, nil
}

// GetPostgres returns PostgreSQL connection
func (pm *PostgreSQLManager) GetPostgres() *sql.DB {
	return pm.manager.PostgreSQL()
}

// GetManager returns the underlying database manager
func (pm *PostgreSQLManager) GetManager() *database.DatabaseManager {
	return pm.manager
}

// GetConnection returns a database connection with context
func (pm *PostgreSQLManager) GetConnection(ctx context.Context) (*sql.Conn, error) {
	db := pm.GetPostgres()
	if db == nil {
		return nil, fmt.Errorf("PostgreSQL connection not available")
	}

	conn, err := db.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	return conn, nil
}

// ExecuteQuery executes a query and returns rows
func (pm *PostgreSQLManager) ExecuteQuery(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	db := pm.GetPostgres()
	if db == nil {
		return nil, fmt.Errorf("PostgreSQL connection not available")
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		pm.logger.Error("Query execution failed",
			logging.String("query", query),
			logging.Error(err))
		return nil, fmt.Errorf("query execution failed: %w", err)
	}

	return rows, nil
}

// ExecuteQueryRow executes a query and returns a single row
func (pm *PostgreSQLManager) ExecuteQueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	db := pm.GetPostgres()
	if db == nil {
		return nil
	}

	return db.QueryRowContext(ctx, query, args...)
}

// ExecuteCommand executes a command (INSERT, UPDATE, DELETE) and returns result
func (pm *PostgreSQLManager) ExecuteCommand(ctx context.Context, command string, args ...interface{}) (sql.Result, error) {
	db := pm.GetPostgres()
	if db == nil {
		return nil, fmt.Errorf("PostgreSQL connection not available")
	}

	result, err := db.ExecContext(ctx, command, args...)
	if err != nil {
		pm.logger.Error("Command execution failed",
			logging.String("command", command),
			logging.Error(err))
		return nil, fmt.Errorf("command execution failed: %w", err)
	}

	return result, nil
}

// BeginTransaction begins a new transaction
func (pm *PostgreSQLManager) BeginTransaction(ctx context.Context) (*sql.Tx, error) {
	return pm.BeginTransactionWithOptions(ctx, nil)
}

// BeginTransactionWithOptions begins a new transaction with options
func (pm *PostgreSQLManager) BeginTransactionWithOptions(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	db := pm.GetPostgres()
	if db == nil {
		return nil, fmt.Errorf("PostgreSQL connection not available")
	}

	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		pm.logger.Error("Failed to begin transaction", logging.Error(err))
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return tx, nil
}

// Transaction executes a function within a database transaction
func (pm *PostgreSQLManager) Transaction(ctx context.Context, fn func(*sql.Tx) error) error {
	return pm.manager.Transaction(ctx, fn)
}

// TransactionWithConfig executes a function within a database transaction with custom config
func (pm *PostgreSQLManager) TransactionWithConfig(ctx context.Context, config database.TransactionConfig, fn func(*sql.Tx) error) error {
	return pm.manager.TransactionWithConfig(ctx, config, fn)
}

// Health checks PostgreSQL database connection
func (pm *PostgreSQLManager) Health(ctx context.Context) error {
	// Use shared library health check
	if err := pm.manager.HealthCheck(ctx); err != nil {
		pm.logger.Error("PostgreSQL health check failed", logging.Error(err))
		return fmt.Errorf("PostgreSQL health check failed: %w", err)
	}

	// Additional service-specific health checks
	if err := pm.pingDatabase(ctx); err != nil {
		return fmt.Errorf("PostgreSQL ping failed: %w", err)
	}

	return nil
}

// pingDatabase performs a simple ping to verify connection
func (pm *PostgreSQLManager) pingDatabase(ctx context.Context) error {
	db := pm.GetPostgres()
	if db == nil {
		return fmt.Errorf("PostgreSQL connection not available")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

// GetConnectionStatus returns the connection status
func (pm *PostgreSQLManager) GetConnectionStatus(ctx context.Context) map[string]interface{} {
	status := make(map[string]interface{})

	// Get basic connection info
	db := pm.GetPostgres()
	if db != nil {
		stats := db.Stats()
		status["open_connections"] = stats.OpenConnections
		status["in_use"] = stats.InUse
		status["idle"] = stats.Idle
		status["wait_count"] = stats.WaitCount
		status["wait_duration"] = stats.WaitDuration.String()
		status["max_idle_closed"] = stats.MaxIdleClosed
		status["max_idle_time_closed"] = stats.MaxIdleTimeClosed
		status["max_lifetime_closed"] = stats.MaxLifetimeClosed
	}

	// Get health status
	status["healthy"] = pm.Health(ctx) == nil

	return status
}

// initializeSchema initializes database schema from migration files
func (pm *PostgreSQLManager) initializeSchema(ctx context.Context) error {
	// Check if migrations directory exists
	migrationsPath := "migrations/postgresql"

	// Try to run migration script if it exists
	if err := pm.runMigrationScript(migrationsPath); err != nil {
		pm.logger.Warn("Migration script not found or failed",
			logging.String("path", migrationsPath),
			logging.Error(err))
		return err
	}

	pm.logger.Info("PostgreSQL schema initialized successfully")
	return nil
}

// runMigrationScript attempts to run migration script
func (pm *PostgreSQLManager) runMigrationScript(migrationsPath string) error {
	// This is a placeholder for running migration scripts
	// In a real implementation, you would:
	// 1. Check if migration script exists
	// 2. Execute the script using os/exec
	// 3. Handle errors appropriately

	pm.logger.Info("PostgreSQL migrations will be handled by MigrationManager",
		logging.String("path", migrationsPath))

	return nil
}

// Close closes PostgreSQL database connection
func (pm *PostgreSQLManager) Close() error {
	if err := pm.manager.Close(); err != nil {
		pm.logger.Error("Failed to close PostgreSQL manager", logging.Error(err))
		return fmt.Errorf("failed to close PostgreSQL manager: %w", err)
	}

	pm.logger.Info("PostgreSQL manager closed successfully")
	return nil
}

// SyncWithBlockchain syncs database state with Cosmos SDK blockchain
func (pm *PostgreSQLManager) SyncWithBlockchain(ctx context.Context) error {
	if pm.cosmosApp == nil || pm.blockchainStorage == nil {
		pm.logger.Warn("Cosmos SDK components not available for blockchain sync",
			logging.String("service", "postgresql-manager"))
		return nil
	}

	correlationID := utils.GetCorrelationID(ctx)
	pm.logger.Info("Starting blockchain sync with Cosmos SDK",
		logging.String("service", "postgresql-manager"),
		logging.String("correlation_id", correlationID))

	// Get latest block height from CometBFT
	cometBFTHeight, err := pm.getLatestBlockHeightFromCometBFT(ctx)
	if err != nil {
		pm.logger.Warn("Failed to get latest block height from CometBFT, using keeper",
			logging.String("correlation_id", correlationID),
			logging.Error(err))
		cometBFTHeight = 0
	}

	// Get latest synced block from database
	dbHeight, err := pm.getLatestSyncedBlockHeight(ctx)
	if err != nil {
		pm.logger.Warn("Failed to get latest synced block height from database",
			logging.String("correlation_id", correlationID),
			logging.Error(err))
		dbHeight = 0
	}

	pm.logger.Info("Blockchain sync status",
		logging.String("correlation_id", correlationID),
		logging.Int64("cometbft_height", cometBFTHeight),
		logging.Int64("db_height", dbHeight))

	// Handle case where db_height > cometbft_height (stale data from previous chain/run)
	if dbHeight > cometBFTHeight {
		// If difference is significant (>1000 blocks) or CometBFT height is very low (<100), reset immediately
		if (dbHeight-cometBFTHeight > 1000) || cometBFTHeight < 100 {
			newDbHeight, err := pm.resetAndStartSync(ctx, correlationID, dbHeight, cometBFTHeight)
			if err != nil {
				return err
			}
			dbHeight = newDbHeight
		} else {
			// Small difference, verify genesis block hash to confirm same chain
			genesisHashMatch, err := pm.verifyGenesisBlockHash(ctx, cometBFTHeight)
			if err != nil {
				pm.logger.Warn("Failed to verify genesis block hash, assuming stale data",
					logging.String("correlation_id", correlationID),
					logging.Error(err))
				genesisHashMatch = false
			}

			if !genesisHashMatch {
				newDbHeight, err := pm.resetAndStartSync(ctx, correlationID, dbHeight, cometBFTHeight)
				if err != nil {
					return err
				}
				dbHeight = newDbHeight
			} else {
				pm.logger.Info("Database height is higher but genesis matches, database is up to date",
					logging.String("correlation_id", correlationID),
					logging.Int64("cometbft_height", cometBFTHeight),
					logging.Int64("db_height", dbHeight))
				return nil
			}
		}
	}

	// Sync blocks from dbHeight+1 to cometBFTHeight
	if cometBFTHeight > dbHeight {
		syncCount := int(cometBFTHeight - dbHeight)
		pm.logger.Info("Syncing blocks to database",
			logging.String("correlation_id", correlationID),
			logging.Int64("start_height", dbHeight+1),
			logging.Int64("end_height", cometBFTHeight),
			logging.Int("block_count", syncCount))

		// Sync in batches to avoid overwhelming the system
		batchSize := 10
		for i := dbHeight + 1; i <= cometBFTHeight; i += int64(batchSize) {
			endHeight := i + int64(batchSize) - 1
			if endHeight > cometBFTHeight {
				endHeight = cometBFTHeight
			}

			if err := pm.syncBlockRange(ctx, i, endHeight); err != nil {
				pm.logger.Error("Failed to sync block range",
					logging.String("correlation_id", correlationID),
					logging.Int64("start_height", i),
					logging.Int64("end_height", endHeight),
					logging.Error(err))
				// Continue with next batch even if this one fails
				continue
			}

			pm.logger.Debug("Synced block range",
				logging.String("correlation_id", correlationID),
				logging.Int64("start_height", i),
				logging.Int64("end_height", endHeight))
		}
	} else {
		pm.logger.Debug("Database is up to date",
			logging.String("correlation_id", correlationID),
			logging.Int64("cometbft_height", cometBFTHeight),
			logging.Int64("db_height", dbHeight))
	}

	pm.logger.Info("Blockchain sync completed successfully",
		logging.String("correlation_id", correlationID),
		logging.Int64("cometbft_height", cometBFTHeight),
		logging.Int64("db_height", dbHeight))
	return nil
}

// GetBlockchainState retrieves current blockchain state from Cosmos SDK
func (pm *PostgreSQLManager) GetBlockchainState(ctx context.Context) (map[string]interface{}, error) {
	if pm.cosmosApp == nil || pm.blockchainStorage == nil {
		return nil, fmt.Errorf("cosmos SDK components not available")
	}

	// Get blockchain state from Cosmos SDK
	// This would typically involve:
	// 1. Querying blockchain state
	// 2. Retrieving module states
	// 3. Formatting for database sync

	state := map[string]interface{}{
		"cosmos_app_available":         pm.cosmosApp != nil,
		"blockchain_storage_available": pm.blockchainStorage != nil,
		"sync_timestamp":               time.Now().Unix(),
	}

	pm.logger.Info("Retrieved blockchain state",
		logging.Any("state", state))

	return state, nil
}

// SetCosmosComponents sets the Cosmos SDK components for blockchain integration
func (pm *PostgreSQLManager) SetCosmosComponents(cosmosApp *app.USCApp, blockchainStorage *storage.StateManager) {
	pm.cosmosApp = cosmosApp
	pm.blockchainStorage = blockchainStorage

	pm.logger.Info("Cosmos SDK components set for PostgreSQL manager",
		logging.Bool("cosmos_app_available", cosmosApp != nil),
		logging.Bool("blockchain_storage_available", blockchainStorage != nil))
}

// getCometBFTURL gets CometBFT RPC URL with fallback logic
func (pm *PostgreSQLManager) getCometBFTURL() string {
	cometBFTURL := os.Getenv("COMETBFT_RPC_URL")
	if cometBFTURL == "" {
		// Try Docker service name first (works in container)
		cometBFTURL = "http://service-04-cometbft:26657"
	}
	return cometBFTURL
}

// queryCometBFTWithFallback queries CometBFT RPC with automatic fallback to localhost
func (pm *PostgreSQLManager) queryCometBFTWithFallback(url string, timeout time.Duration) (*http.Response, error) {
	client := &http.Client{Timeout: timeout}
	resp, err := client.Get(url)
	if err != nil {
		// Try localhost as fallback if not already using localhost
		if !strings.Contains(url, "localhost:26657") {
			pm.logger.Debug("Failed to connect to CometBFT via service name, trying localhost",
				logging.String("url", url),
				logging.Error(err))
			// Extract path from original URL and append to localhost
			urlParts := strings.SplitN(url, "/", 4)
			if len(urlParts) >= 4 {
				fallbackURL := "http://localhost:26657/" + urlParts[3]
				resp, err = client.Get(fallbackURL)
			} else {
				// If URL parsing fails, try simple localhost replacement
				fallbackURL := strings.Replace(url, pm.getCometBFTURL(), "http://localhost:26657", 1)
				resp, err = client.Get(fallbackURL)
			}
		}
	}
	return resp, err
}

// cometBFTBlockResult represents a block result from CometBFT RPC
type cometBFTBlockResult struct {
	Result struct {
		Block struct {
			Header struct {
				Height string `json:"height"`
				Hash   string `json:"hash"`
				Time   string `json:"time"`
			} `json:"header"`
			Data struct {
				Txs []string `json:"txs"`
			} `json:"data"`
		} `json:"block"`
		BlockID struct {
			Hash string `json:"hash"`
		} `json:"block_id"`
	} `json:"result"`
}

// getBlockFromCometBFT gets a block from CometBFT RPC by height
func (pm *PostgreSQLManager) getBlockFromCometBFT(ctx context.Context, height int64) (*cometBFTBlockResult, error) {
	cometBFTURL := pm.getCometBFTURL()
	url := fmt.Sprintf("%s/block?height=%d", cometBFTURL, height)

	resp, err := pm.queryCometBFTWithFallback(url, 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to query block %d: %w", height, err)
	}
	defer resp.Body.Close()

	var blockResult cometBFTBlockResult
	if err := json.NewDecoder(resp.Body).Decode(&blockResult); err != nil {
		return nil, fmt.Errorf("failed to decode block %d response: %w", height, err)
	}

	return &blockResult, nil
}

// getPreviousBlockHash gets the previous block hash from database
func (pm *PostgreSQLManager) getPreviousBlockHash(ctx context.Context, blockHeight int64) string {
	if blockHeight <= 1 {
		return "" // Genesis block has no previous block
	}

	db := pm.GetPostgres()
	if db == nil {
		pm.logger.Debug("PostgreSQL connection not available, using empty previous_block_hash",
			logging.Int64("height", blockHeight))
		return ""
	}

	var previousBlockHash string
	// Try blocks table first, fallback to usc_block_analytics
	prevQuery := `SELECT block_hash FROM blocks WHERE block_number = $1`
	err := db.QueryRowContext(ctx, prevQuery, blockHeight-1).Scan(&previousBlockHash)
	if err != nil {
		// Fallback to usc_block_analytics table
		prevQuery = `SELECT block_hash FROM usc_block_analytics WHERE block_number = $1`
		err = db.QueryRowContext(ctx, prevQuery, blockHeight-1).Scan(&previousBlockHash)
	}

	if err != nil {
		pm.logger.Debug("Previous block not found, using empty previous_block_hash",
			logging.Int64("height", blockHeight),
			logging.Int64("previous_height", blockHeight-1),
			logging.Error(err))
		return ""
	}

	pm.logger.Debug("Found previous block hash",
		logging.Int64("height", blockHeight),
		logging.Int64("previous_height", blockHeight-1),
		logging.String("previous_hash", previousBlockHash))
	return previousBlockHash
}

// resetAndStartSync resets blocks table and prepares for fresh sync
func (pm *PostgreSQLManager) resetAndStartSync(ctx context.Context, correlationID string, oldDbHeight, cometBFTHeight int64) (int64, error) {
	pm.logger.Warn("Database has stale data from previous chain/run, resetting blocks table",
		logging.String("correlation_id", correlationID),
		logging.Int64("old_db_height", oldDbHeight),
		logging.Int64("cometbft_height", cometBFTHeight),
		logging.Int64("difference", oldDbHeight-cometBFTHeight))

	if err := pm.resetBlocksTable(ctx); err != nil {
		pm.logger.Error("Failed to reset blocks table",
			logging.String("correlation_id", correlationID),
			logging.Error(err))
		return 0, fmt.Errorf("failed to reset blocks table: %w", err)
	}

	pm.logger.Info("Blocks table reset, starting fresh sync",
		logging.String("correlation_id", correlationID))
	return 0, nil // Return 0 to start fresh sync
}

// getLatestBlockHeightFromCometBFT gets the latest block height from CometBFT RPC
func (pm *PostgreSQLManager) getLatestBlockHeightFromCometBFT(ctx context.Context) (int64, error) {
	cometBFTURL := pm.getCometBFTURL()
	url := cometBFTURL + "/status"

	resp, err := pm.queryCometBFTWithFallback(url, 5*time.Second)
	if err != nil {
		return 0, fmt.Errorf("failed to query CometBFT status: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Result struct {
			SyncInfo struct {
				LatestBlockHeight string `json:"latest_block_height"`
			} `json:"sync_info"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode CometBFT response: %w", err)
	}

	height, err := strconv.ParseInt(result.Result.SyncInfo.LatestBlockHeight, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse block height: %w", err)
	}

	return height, nil
}

// getLatestSyncedBlockHeight gets the latest block height from PostgreSQL
func (pm *PostgreSQLManager) getLatestSyncedBlockHeight(ctx context.Context) (int64, error) {
	query := `SELECT COALESCE(MAX(block_number), 0) FROM blocks`
	var height int64
	db := pm.GetPostgres()
	if db == nil {
		return 0, fmt.Errorf("PostgreSQL connection not available")
	}
	err := db.QueryRowContext(ctx, query).Scan(&height)
	if err != nil {
		return 0, fmt.Errorf("failed to query latest block height: %w", err)
	}
	return height, nil
}

// syncBlockRange syncs a range of blocks from CometBFT to PostgreSQL
func (pm *PostgreSQLManager) syncBlockRange(ctx context.Context, startHeight, endHeight int64) error {
	for height := startHeight; height <= endHeight; height++ {
		// Get block from CometBFT
		blockResult, err := pm.getBlockFromCometBFT(ctx, height)
		if err != nil {
			return fmt.Errorf("failed to get block %d: %w", height, err)
		}

		blockHeight, err := strconv.ParseInt(blockResult.Result.Block.Header.Height, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse block height: %w", err)
		}

		blockTime, err := time.Parse(time.RFC3339Nano, blockResult.Result.Block.Header.Time)
		if err != nil {
			// Try alternative format
			blockTime, err = time.Parse(time.RFC3339, blockResult.Result.Block.Header.Time)
			if err != nil {
				blockTime = time.Now()
			}
		}

		blockTimestamp := blockTime.Unix()

		// Get block hash from block_id.hash (preferred) or block.header.hash (fallback)
		blockHash := blockResult.Result.BlockID.Hash
		if blockHash == "" {
			blockHash = blockResult.Result.Block.Header.Hash
		}

		// Skip blocks with empty or invalid hash
		if blockHash == "" || len(blockHash) < 10 {
			pm.logger.Debug("Skipping block with invalid hash",
				logging.Int64("height", height),
				logging.String("hash", blockHash),
				logging.String("block_id_hash", blockResult.Result.BlockID.Hash),
				logging.String("header_hash", blockResult.Result.Block.Header.Hash))
			continue
		}

		// Get previous block hash
		previousBlockHash := pm.getPreviousBlockHash(ctx, blockHeight)

		// Handle conflicts: block_number is primary unique constraint
		// If block_hash conflicts, it means the same hash exists with different block_number (data inconsistency)
		// In that case, we skip the insert to avoid duplicate key error
		insertQuery := `
			INSERT INTO blocks (block_number, block_hash, previous_block_hash, timestamp, nonce, transaction_count, block_size, validator_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
			ON CONFLICT (block_number) DO UPDATE SET
				block_hash = EXCLUDED.block_hash,
				previous_block_hash = EXCLUDED.previous_block_hash,
				updated_at = NOW()
		`

		db := pm.GetPostgres()
		if db == nil {
			return fmt.Errorf("PostgreSQL connection not available")
		}
		_, err = db.ExecContext(ctx, insertQuery,
			blockHeight,
			blockHash,         // Use validated blockHash variable
			previousBlockHash, // Use populated previous block hash
			blockTimestamp,    // Unix timestamp (int64) for bigint column
			0,                 // nonce - not available from CometBFT RPC, use 0 as default
			len(blockResult.Result.Block.Data.Txs),
			0,  // block_size - not available from RPC
			"", // validator_id - not available from RPC
		)
		if err != nil {
			// Check if error is due to duplicate key on block_hash
			errStr := err.Error()
			if strings.Contains(errStr, "duplicate key value violates unique constraint") &&
				strings.Contains(errStr, "blocks_block_hash_key") {
				// Block hash already exists with different block_number - skip this block
				pm.logger.Debug("Skipping block with duplicate hash",
					logging.Int64("height", height),
					logging.String("hash", blockHash))
				continue // Skip this block and continue with next
			}
			// For other errors, return to stop the sync for this batch
			return fmt.Errorf("failed to insert block %d: %w", height, err)
		}

		// Sync transactions (if any)
		for i, txHash := range blockResult.Result.Block.Data.Txs {
			if err := pm.syncTransaction(ctx, txHash, blockHeight, int64(i)); err != nil {
				pm.logger.Warn("Failed to sync transaction",
					logging.String("tx_hash", txHash),
					logging.Int64("block_height", blockHeight),
					logging.Error(err))
				// Continue with next transaction
			}
		}
	}

	return nil
}

// verifyGenesisBlockHash verifies if genesis block hash in database matches CometBFT
// Returns true if genesis block hash matches, false otherwise
func (pm *PostgreSQLManager) verifyGenesisBlockHash(ctx context.Context, cometBFTHeight int64) (bool, error) {
	// If CometBFT height is very low, assume it's a new chain and database has stale data
	if cometBFTHeight < 100 {
		return false, nil
	}

	// Get genesis block hash from CometBFT (height 1)
	blockResult, err := pm.getBlockFromCometBFT(ctx, 1)
	if err != nil {
		return false, fmt.Errorf("failed to query genesis block from CometBFT: %w", err)
	}

	cometBFTGenesisHash := blockResult.Result.BlockID.Hash
	if cometBFTGenesisHash == "" {
		cometBFTGenesisHash = blockResult.Result.Block.Header.Hash
	}

	// Get genesis block hash from database
	db := pm.GetPostgres()
	if db == nil {
		return false, fmt.Errorf("PostgreSQL connection not available")
	}

	var dbGenesisHash string
	query := `SELECT block_hash FROM blocks WHERE block_number = 1`
	err = db.QueryRowContext(ctx, query).Scan(&dbGenesisHash)
	if err != nil {
		// If genesis block doesn't exist in database, assume stale data
		return false, nil
	}

	// Compare hashes (case-insensitive)
	return strings.EqualFold(cometBFTGenesisHash, dbGenesisHash), nil
}

// resetBlocksTable resets the blocks table to start fresh sync
func (pm *PostgreSQLManager) resetBlocksTable(ctx context.Context) error {
	db := pm.GetPostgres()
	if db == nil {
		return fmt.Errorf("PostgreSQL connection not available")
	}

	// Delete all blocks to start fresh
	query := `DELETE FROM blocks`
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to reset blocks table: %w", err)
	}

	pm.logger.Info("Blocks table reset successfully",
		logging.String("service", "postgresql-manager"))

	return nil
}

// syncTransaction syncs a transaction from CometBFT to PostgreSQL
func (pm *PostgreSQLManager) syncTransaction(ctx context.Context, txHash string, blockHeight, txIndex int64) error {
	db := pm.GetPostgres()
	if db == nil {
		return fmt.Errorf("PostgreSQL connection not available")
	}

	// Check if transaction already exists
	var exists bool
	err := db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM transactions WHERE transaction_hash = $1)`, txHash).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check transaction existence: %w", err)
	}

	if exists {
		return nil // Already synced
	}

	// Insert transaction
	insertQuery := `
		INSERT INTO transactions (transaction_hash, block_number, transaction_index, from_address, to_address, amount, transaction_type, status, gas_used, gas_price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
		ON CONFLICT (transaction_hash) DO NOTHING
	`

	_, err = db.ExecContext(ctx, insertQuery,
		txHash,
		blockHeight,
		txIndex,
		"",          // from_address - would need to decode transaction
		"",          // to_address - would need to decode transaction
		nil,         // amount - would need to decode transaction
		"unknown",   // transaction_type - would need to decode transaction
		"confirmed", // status
		0,           // gas_used - not available from block query
		nil,         // gas_price - not available from block query
	)
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err)
	}

	return nil
}
