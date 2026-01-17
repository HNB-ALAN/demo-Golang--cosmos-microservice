package app

import (
	"fmt"
	"os"

	cmtdb "github.com/cometbft/cometbft-db"
	cosmosdb "github.com/cosmos/cosmos-db"
)

// NewUSCAppWithRocksDB creates and initializes a USC application with RocksDB database
// This is a convenience function that handles database creation and app initialization
// Note: Uses RocksDB for Cosmos SDK (requires build tags "rocksdb_legacy rocksdb")
// RocksDB is also used for business logic via RocksDBManager
// dbDir: Directory path for the RocksDB database (default: "./data/cosmos")
func NewUSCAppWithRocksDB(dbDir string) (*USCApp, cosmosdb.DB, error) {
	// Default directory if not provided
	if dbDir == "" {
		dbDir = "./data/cosmos"
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, nil, fmt.Errorf("failed to create cosmos database directory: %w", err)
	}

	// Open RocksDB database for Cosmos SDK
	// Note: Requires build tags "rocksdb_legacy rocksdb" and CGO enabled
	// RocksDB is the recommended backend for production according to Cosmos SDK docs
	cmtDB, err := cmtdb.NewDB("application", cmtdb.RocksDBBackend, dbDir)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open Cosmos SDK RocksDB: %w", err)
	}

	// Wrap in cosmos-db adapter
	cosmosDatabase := &cosmosDBAdapter{db: cmtDB}

	// Create USC app
	uscApp := NewUSCApp(cosmosDatabase)

	// Initialize the app
	if err := uscApp.Initialize(); err != nil {
		cosmosDatabase.Close()
		return nil, nil, fmt.Errorf("failed to initialize Cosmos SDK app: %w", err)
	}

	return uscApp, cosmosDatabase, nil
}

// cosmosDBAdapter adapts cometbft-db to cosmos-db interface
type cosmosDBAdapter struct {
	db cmtdb.DB
}

func (a *cosmosDBAdapter) Get(key []byte) ([]byte, error) {
	return a.db.Get(key)
}

func (a *cosmosDBAdapter) Has(key []byte) (bool, error) {
	return a.db.Has(key)
}

func (a *cosmosDBAdapter) Set(key, value []byte) error {
	return a.db.Set(key, value)
}

func (a *cosmosDBAdapter) SetSync(key, value []byte) error {
	return a.db.SetSync(key, value)
}

func (a *cosmosDBAdapter) Delete(key []byte) error {
	return a.db.Delete(key)
}

func (a *cosmosDBAdapter) DeleteSync(key []byte) error {
	return a.db.DeleteSync(key)
}

func (a *cosmosDBAdapter) Iterator(start, end []byte) (cosmosdb.Iterator, error) {
	iter, err := a.db.Iterator(start, end)
	if err != nil {
		return nil, err
	}
	return &iteratorAdapter{iter: iter}, nil
}

func (a *cosmosDBAdapter) ReverseIterator(start, end []byte) (cosmosdb.Iterator, error) {
	iter, err := a.db.ReverseIterator(start, end)
	if err != nil {
		return nil, err
	}
	return &iteratorAdapter{iter: iter}, nil
}

func (a *cosmosDBAdapter) Close() error {
	return a.db.Close()
}

func (a *cosmosDBAdapter) NewBatch() cosmosdb.Batch {
	return &batchAdapter{batch: a.db.NewBatch()}
}

func (a *cosmosDBAdapter) NewBatchWithSize(size int) cosmosdb.Batch {
	return &batchAdapter{batch: a.db.NewBatch()}
}

func (a *cosmosDBAdapter) Print() error {
	return a.db.Print()
}

func (a *cosmosDBAdapter) Stats() map[string]string {
	return a.db.Stats()
}

// iteratorAdapter adapts cometbft-db Iterator to cosmos-db Iterator
type iteratorAdapter struct {
	iter cmtdb.Iterator
}

func (i *iteratorAdapter) Domain() (start, end []byte) {
	return i.iter.Domain()
}

func (i *iteratorAdapter) Valid() bool {
	return i.iter.Valid()
}

func (i *iteratorAdapter) Next() {
	i.iter.Next()
}

func (i *iteratorAdapter) Key() []byte {
	return i.iter.Key()
}

func (i *iteratorAdapter) Value() []byte {
	return i.iter.Value()
}

func (i *iteratorAdapter) Error() error {
	return i.iter.Error()
}

func (i *iteratorAdapter) Close() error {
	i.iter.Close()
	return nil
}

// batchAdapter adapts cometbft-db Batch to cosmos-db Batch
type batchAdapter struct {
	batch cmtdb.Batch
}

func (b *batchAdapter) Set(key, value []byte) error {
	return b.batch.Set(key, value)
}

func (b *batchAdapter) Delete(key []byte) error {
	return b.batch.Delete(key)
}

func (b *batchAdapter) Write() error {
	return b.batch.Write()
}

func (b *batchAdapter) WriteSync() error {
	return b.batch.WriteSync()
}

func (b *batchAdapter) Close() error {
	return b.batch.Close()
}

func (b *batchAdapter) GetByteSize() (int, error) {
	// cometbft-db Batch doesn't provide size info
	return 0, nil
}
