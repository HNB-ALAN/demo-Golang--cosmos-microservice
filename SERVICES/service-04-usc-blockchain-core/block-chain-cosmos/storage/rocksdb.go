//go:build cgo
// +build cgo

package storage

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/linxGnu/grocksdb"
)

// RocksDBManager manages RocksDB database operations for the USC blockchain
type RocksDBManager struct {
	db     *grocksdb.DB
	ro     *grocksdb.ReadOptions
	wo     *grocksdb.WriteOptions
	mu     sync.RWMutex
	closed bool
}

// RocksDBConfig contains configuration for RocksDB
type RocksDBConfig struct {
	DataPath        string        `json:"data_path"`
	MaxOpenFiles    int           `json:"max_open_files"`
	WriteBufferSize int           `json:"write_buffer_size"`
	MaxWriteBuffer  int           `json:"max_write_buffer"`
	BlockSize       int           `json:"block_size"`
	CacheSize       int           `json:"cache_size"`
	Compression     string        `json:"compression"`
	SyncWrites      bool          `json:"sync_writes"`
	Timeout         time.Duration `json:"timeout"`
}

// DefaultRocksDBConfig returns the default RocksDB configuration
func DefaultRocksDBConfig() RocksDBConfig {
	return RocksDBConfig{
		DataPath:        "./data/rocksdb",
		MaxOpenFiles:    1000,
		WriteBufferSize: 64 * 1024 * 1024, // 64MB
		MaxWriteBuffer:  3,
		BlockSize:       4 * 1024,          // 4KB
		CacheSize:       128 * 1024 * 1024, // 128MB
		Compression:     "snappy",
		SyncWrites:      true,
		Timeout:         30 * time.Second,
	}
}

// NewRocksDBManager creates a new RocksDB manager
func NewRocksDBManager(config RocksDBConfig) (*RocksDBManager, error) {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(config.DataPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Configure RocksDB options
	opts := grocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(true)
	opts.SetMaxOpenFiles(config.MaxOpenFiles)
	opts.SetWriteBufferSize(uint64(config.WriteBufferSize))
	opts.SetMaxWriteBufferNumber(config.MaxWriteBuffer)

	// Set compression - use default compression for now
	// opts.SetCompression(grocksdb.SnappyCompression)

	// Open database
	db, err := grocksdb.OpenDb(opts, config.DataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open RocksDB: %w", err)
	}

	// Create read and write options
	ro := grocksdb.NewDefaultReadOptions()
	wo := grocksdb.NewDefaultWriteOptions()
	wo.SetSync(config.SyncWrites)

	return &RocksDBManager{
		db: db,
		ro: ro,
		wo: wo,
	}, nil
}

// Get retrieves a value by key
func (r *RocksDBManager) Get(ctx context.Context, key []byte) ([]byte, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.closed {
		return nil, fmt.Errorf("database is closed")
	}

	value, err := r.db.Get(r.ro, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get value: %w", err)
	}
	defer value.Free()

	if value.Data() == nil {
		return nil, fmt.Errorf("key not found: %x", key)
	}

	// Copy data to avoid memory issues
	result := make([]byte, len(value.Data()))
	copy(result, value.Data())
	return result, nil
}

// Set stores a key-value pair
func (r *RocksDBManager) Set(ctx context.Context, key, value []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return fmt.Errorf("database is closed")
	}

	err := r.db.Put(r.wo, key, value)
	if err != nil {
		return fmt.Errorf("failed to set value: %w", err)
	}

	return nil
}

// Delete removes a key
func (r *RocksDBManager) Delete(ctx context.Context, key []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return fmt.Errorf("database is closed")
	}

	err := r.db.Delete(r.wo, key)
	if err != nil {
		return fmt.Errorf("failed to delete key: %w", err)
	}

	return nil
}

// Has checks if a key exists
func (r *RocksDBManager) Has(ctx context.Context, key []byte) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.closed {
		return false, fmt.Errorf("database is closed")
	}

	value, err := r.db.Get(r.ro, key)
	if err != nil {
		return false, fmt.Errorf("failed to check key existence: %w", err)
	}
	defer value.Free()

	return value.Data() != nil, nil
}

// Iterator creates an iterator for key range scanning
func (r *RocksDBManager) Iterator(ctx context.Context, start, end []byte) (*RocksDBIterator, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.closed {
		return nil, fmt.Errorf("database is closed")
	}

	iter := r.db.NewIterator(r.ro)
	if start != nil {
		iter.Seek(start)
	} else {
		iter.SeekToFirst()
	}

	return &RocksDBIterator{
		iter: iter,
		end:  end,
	}, nil
}

// Batch performs batch operations
func (r *RocksDBManager) Batch(ctx context.Context, operations []BatchOperation) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return fmt.Errorf("database is closed")
	}

	// Create write batch
	wb := grocksdb.NewWriteBatch()
	defer wb.Destroy()

	// Add operations to batch
	for _, op := range operations {
		switch op.Type {
		case BatchOpSet:
			wb.Put(op.Key, op.Value)
		case BatchOpDelete:
			wb.Delete(op.Key)
		default:
			return fmt.Errorf("unknown batch operation type: %d", op.Type)
		}
	}

	// Execute batch
	err := r.db.Write(r.wo, wb)
	if err != nil {
		return fmt.Errorf("failed to execute batch: %w", err)
	}

	return nil
}

// Close closes the database
func (r *RocksDBManager) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return nil
	}

	r.closed = true

	// Close RocksDB resources
	if r.ro != nil {
		r.ro.Destroy()
	}
	if r.wo != nil {
		r.wo.Destroy()
	}
	if r.db != nil {
		r.db.Close()
	}

	return nil
}

// Stats returns database statistics
func (r *RocksDBManager) Stats() map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.closed {
		return nil
	}

	stats := make(map[string]interface{})
	stats["is_closed"] = r.closed

	// Get RocksDB statistics
	if r.db != nil {
		stats["db_stats"] = r.db.GetProperty("rocksdb.stats")
		stats["num_keys"] = r.db.GetProperty("rocksdb.num-keys")
		stats["db_size"] = r.db.GetProperty("rocksdb.estimate-num-keys")
	}

	return stats
}

// BatchOperation represents a single batch operation
type BatchOperation struct {
	Type  BatchOpType `json:"type"`
	Key   []byte      `json:"key"`
	Value []byte      `json:"value,omitempty"`
}

// BatchOpType represents the type of batch operation
type BatchOpType int

const (
	BatchOpSet BatchOpType = iota
	BatchOpDelete
)

// RocksDBIterator provides iteration over key-value pairs
type RocksDBIterator struct {
	iter *grocksdb.Iterator
	end  []byte
}

// Valid returns true if the iterator is valid
func (i *RocksDBIterator) Valid() bool {
	if i.iter == nil {
		return false
	}

	if !i.iter.Valid() {
		return false
	}

	// Check if we've reached the end key
	if i.end != nil {
		key := i.iter.Key()
		defer key.Free()
		if key.Data() != nil && string(key.Data()) >= string(i.end) {
			return false
		}
	}

	return true
}

// Key returns the current key
func (i *RocksDBIterator) Key() []byte {
	if !i.Valid() {
		return nil
	}

	key := i.iter.Key()
	defer key.Free()

	if key.Data() == nil {
		return nil
	}

	// Copy data to avoid memory issues
	result := make([]byte, len(key.Data()))
	copy(result, key.Data())
	return result
}

// Value returns the current value
func (i *RocksDBIterator) Value() []byte {
	if !i.Valid() {
		return nil
	}

	value := i.iter.Value()
	defer value.Free()

	if value.Data() == nil {
		return nil
	}

	// Copy data to avoid memory issues
	result := make([]byte, len(value.Data()))
	copy(result, value.Data())
	return result
}

// Next moves to the next key-value pair
func (i *RocksDBIterator) Next() {
	if i.iter != nil && i.iter.Valid() {
		i.iter.Next()
	}
}

// Close closes the iterator
func (i *RocksDBIterator) Close() {
	if i.iter != nil {
		i.iter.Close()
		i.iter = nil
	}
}
