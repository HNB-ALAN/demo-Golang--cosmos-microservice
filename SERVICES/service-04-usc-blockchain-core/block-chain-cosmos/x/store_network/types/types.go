package types

import (
	"fmt"
	"time"
)

const (
	// ModuleName defines the module name
	ModuleName = "store_network"

	// RouterKey defines the message route for the store module
	RouterKey = ModuleName

	// QuerierRoute defines the querier route for the store module
	QuerierRoute = ModuleName
)

// Event types
const (
	EventTypeDataStored     = "data_stored"
	EventTypeDataRetrieved  = "data_retrieved"
	EventTypeDataUpdated    = "data_updated"
	EventTypeDataDeleted    = "data_deleted"
	EventTypeStoreCreated   = "store_created"
	EventTypeStoreDeleted   = "store_deleted"
	EventTypeBackupCreated  = "backup_created"
	EventTypeRestoreCreated = "restore_created"
)

// Event attribute keys
const (
	AttributeKeyDataID    = "data_id"
	AttributeKeyDataKey   = "data_key"
	AttributeKeyDataValue = "data_value"
	AttributeKeyDataSize  = "data_size"
	AttributeKeyStoreID   = "store_id"
	AttributeKeyStoreName = "store_name"
	AttributeKeyBackupID  = "backup_id"
	AttributeKeyRestoreID = "restore_id"
	AttributeKeyOperation = "operation"
	AttributeKeyTimestamp = "timestamp"
)

// StoredData represents data stored in the store module
type StoredData struct {
	ID          string            `json:"id"`
	Key         string            `json:"key"`
	Value       []byte            `json:"value"`
	Size        int64             `json:"size"`
	ContentType string            `json:"content_type"`
	Tags        map[string]string `json:"tags"`
	Metadata    map[string]string `json:"metadata"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	ExpiresAt   time.Time         `json:"expires_at"`
	Version     int64             `json:"version"`
}

// Store represents a data store
type Store struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Type        string            `json:"type"` // kv, document, graph, time-series
	Config      map[string]string `json:"config"`
	Tags        map[string]string `json:"tags"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Status      string            `json:"status"` // active, inactive, maintenance
	Size        int64             `json:"size"`
	ItemCount   int64             `json:"item_count"`
}

// Backup represents a data backup
type Backup struct {
	ID          string            `json:"id"`
	StoreID     string            `json:"store_id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Size        int64             `json:"size"`
	ItemCount   int64             `json:"item_count"`
	Status      string            `json:"status"` // pending, in_progress, completed, failed
	CreatedAt   time.Time         `json:"created_at"`
	CompletedAt time.Time         `json:"completed_at"`
	ExpiresAt   time.Time         `json:"expires_at"`
	Tags        map[string]string `json:"tags"`
	Metadata    map[string]string `json:"metadata"`
}

// Restore represents a data restore operation
type Restore struct {
	ID          string            `json:"id"`
	BackupID    string            `json:"backup_id"`
	StoreID     string            `json:"store_id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Status      string            `json:"status"` // pending, in_progress, completed, failed
	CreatedAt   time.Time         `json:"created_at"`
	CompletedAt time.Time         `json:"completed_at"`
	Tags        map[string]string `json:"tags"`
	Metadata    map[string]string `json:"metadata"`
}

// StoreIndex represents an index for a store
type StoreIndex struct {
	ID        string            `json:"id"`
	StoreID   string            `json:"store_id"`
	Name      string            `json:"name"`
	Fields    []string          `json:"fields"`
	Type      string            `json:"type"` // btree, hash, fulltext, spatial
	Config    map[string]string `json:"config"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Status    string            `json:"status"` // active, inactive, building
}

// StoreQuery represents a query for a store
type StoreQuery struct {
	ID        string            `json:"id"`
	StoreID   string            `json:"store_id"`
	Name      string            `json:"name"`
	Query     string            `json:"query"`
	Type      string            `json:"type"` // select, insert, update, delete
	Params    map[string]string `json:"params"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Tags      map[string]string `json:"tags"`
}

// StoreTransaction represents a transaction for a store
type StoreTransaction struct {
	ID        string            `json:"id"`
	StoreID   string            `json:"store_id"`
	Type      string            `json:"type"`   // read, write, read_write
	Status    string            `json:"status"` // pending, committed, rolled_back
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Tags      map[string]string `json:"tags"`
	Metadata  map[string]string `json:"metadata"`
}

// GenesisState represents the genesis state of the store module
type GenesisState struct {
	StoredData   []StoredData       `json:"stored_data"`
	Stores       []Store            `json:"stores"`
	Backups      []Backup           `json:"backups"`
	Restores     []Restore          `json:"restores"`
	StoreIndexes []StoreIndex       `json:"store_indexes"`
	StoreQueries []StoreQuery       `json:"store_queries"`
	Transactions []StoreTransaction `json:"transactions"`
	Params       Params             `json:"params"`
}

// Params represents the parameters for the store module
type Params struct {
	MaxDataSize        int64         `json:"max_data_size"`
	MaxStoreSize       int64         `json:"max_store_size"`
	DefaultRetention   time.Duration `json:"default_retention"`
	BackupInterval     time.Duration `json:"backup_interval"`
	MaxBackups         int64         `json:"max_backups"`
	CompressionEnabled bool          `json:"compression_enabled"`
	EncryptionEnabled  bool          `json:"encryption_enabled"`
}

// DefaultParams returns the default parameters for the store module
func DefaultParams() Params {
	return Params{
		MaxDataSize:        10 * 1024 * 1024,   // 10MB
		MaxStoreSize:       100 * 1024 * 1024,  // 100MB
		DefaultRetention:   7 * 24 * time.Hour, // 7 days
		BackupInterval:     24 * time.Hour,     // 24 hours
		MaxBackups:         30,                 // 30 backups
		CompressionEnabled: true,
		EncryptionEnabled:  true,
	}
}

// Validate validates the parameters
func (p Params) Validate() error {
	if p.MaxDataSize <= 0 {
		return fmt.Errorf("max data size must be positive")
	}
	if p.MaxStoreSize <= 0 {
		return fmt.Errorf("max store size must be positive")
	}
	if p.DefaultRetention <= 0 {
		return fmt.Errorf("default retention must be positive")
	}
	if p.BackupInterval <= 0 {
		return fmt.Errorf("backup interval must be positive")
	}
	if p.MaxBackups <= 0 {
		return fmt.Errorf("max backups must be positive")
	}
	return nil
}

// Validate validates stored data
func (d StoredData) Validate() error {
	if d.ID == "" {
		return fmt.Errorf("data ID cannot be empty")
	}
	if d.Key == "" {
		return fmt.Errorf("data key cannot be empty")
	}
	if d.Size < 0 {
		return fmt.Errorf("data size cannot be negative")
	}
	if d.Version < 0 {
		return fmt.Errorf("data version cannot be negative")
	}
	return nil
}

// Validate validates a store
func (s Store) Validate() error {
	if s.ID == "" {
		return fmt.Errorf("store ID cannot be empty")
	}
	if s.Name == "" {
		return fmt.Errorf("store name cannot be empty")
	}
	if s.Type != "kv" && s.Type != "document" && s.Type != "graph" && s.Type != "time-series" {
		return fmt.Errorf("invalid store type: %s", s.Type)
	}
	if s.Status != "active" && s.Status != "inactive" && s.Status != "maintenance" {
		return fmt.Errorf("invalid store status: %s", s.Status)
	}
	return nil
}

// Validate validates a backup
func (b Backup) Validate() error {
	if b.ID == "" {
		return fmt.Errorf("backup ID cannot be empty")
	}
	if b.StoreID == "" {
		return fmt.Errorf("store ID cannot be empty")
	}
	if b.Name == "" {
		return fmt.Errorf("backup name cannot be empty")
	}
	if b.Status != "pending" && b.Status != "in_progress" && b.Status != "completed" && b.Status != "failed" {
		return fmt.Errorf("invalid backup status: %s", b.Status)
	}
	return nil
}

// Validate validates a restore
func (r Restore) Validate() error {
	if r.ID == "" {
		return fmt.Errorf("restore ID cannot be empty")
	}
	if r.BackupID == "" {
		return fmt.Errorf("backup ID cannot be empty")
	}
	if r.StoreID == "" {
		return fmt.Errorf("store ID cannot be empty")
	}
	if r.Name == "" {
		return fmt.Errorf("restore name cannot be empty")
	}
	if r.Status != "pending" && r.Status != "in_progress" && r.Status != "completed" && r.Status != "failed" {
		return fmt.Errorf("invalid restore status: %s", r.Status)
	}
	return nil
}

// Validate validates a store index
func (i StoreIndex) Validate() error {
	if i.ID == "" {
		return fmt.Errorf("index ID cannot be empty")
	}
	if i.StoreID == "" {
		return fmt.Errorf("store ID cannot be empty")
	}
	if i.Name == "" {
		return fmt.Errorf("index name cannot be empty")
	}
	if i.Type != "btree" && i.Type != "hash" && i.Type != "fulltext" && i.Type != "spatial" {
		return fmt.Errorf("invalid index type: %s", i.Type)
	}
	if i.Status != "active" && i.Status != "inactive" && i.Status != "building" {
		return fmt.Errorf("invalid index status: %s", i.Status)
	}
	return nil
}

// Validate validates a store query
func (q StoreQuery) Validate() error {
	if q.ID == "" {
		return fmt.Errorf("query ID cannot be empty")
	}
	if q.StoreID == "" {
		return fmt.Errorf("store ID cannot be empty")
	}
	if q.Name == "" {
		return fmt.Errorf("query name cannot be empty")
	}
	if q.Query == "" {
		return fmt.Errorf("query cannot be empty")
	}
	if q.Type != "select" && q.Type != "insert" && q.Type != "update" && q.Type != "delete" {
		return fmt.Errorf("invalid query type: %s", q.Type)
	}
	return nil
}

// Validate validates a store transaction
func (t StoreTransaction) Validate() error {
	if t.ID == "" {
		return fmt.Errorf("transaction ID cannot be empty")
	}
	if t.StoreID == "" {
		return fmt.Errorf("store ID cannot be empty")
	}
	if t.Type != "read" && t.Type != "write" && t.Type != "read_write" {
		return fmt.Errorf("invalid transaction type: %s", t.Type)
	}
	if t.Status != "pending" && t.Status != "committed" && t.Status != "rolled_back" {
		return fmt.Errorf("invalid transaction status: %s", t.Status)
	}
	return nil
}
