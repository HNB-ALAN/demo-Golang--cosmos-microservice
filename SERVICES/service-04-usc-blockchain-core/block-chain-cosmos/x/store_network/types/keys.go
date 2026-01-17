package types

import (
	"bytes"
	"fmt"
	"time"
)

// Store key prefixes
var (
	StoredDataKeyPrefix  = []byte{0x01}
	StoreKeyPrefix       = []byte{0x02}
	BackupKeyPrefix      = []byte{0x03}
	RestoreKeyPrefix     = []byte{0x04}
	StoreIndexKeyPrefix  = []byte{0x05}
	StoreQueryKeyPrefix  = []byte{0x06}
	TransactionKeyPrefix = []byte{0x07}
	ParamsKey            = []byte{0x08}
)

// StoredDataKey returns the key for stored data
func StoredDataKey(id string) []byte {
	return append(StoredDataKeyPrefix, []byte(id)...)
}

// StoreKey returns the key for a store
func StoreKey(id string) []byte {
	return append(StoreKeyPrefix, []byte(id)...)
}

// BackupKey returns the key for a backup
func BackupKey(id string) []byte {
	return append(BackupKeyPrefix, []byte(id)...)
}

// RestoreKey returns the key for a restore
func RestoreKey(id string) []byte {
	return append(RestoreKeyPrefix, []byte(id)...)
}

// StoreIndexKey returns the key for a store index
func StoreIndexKey(id string) []byte {
	return append(StoreIndexKeyPrefix, []byte(id)...)
}

// StoreQueryKey returns the key for a store query
func StoreQueryKey(id string) []byte {
	return append(StoreQueryKeyPrefix, []byte(id)...)
}

// TransactionKey returns the key for a transaction
func TransactionKey(id string) []byte {
	return append(TransactionKeyPrefix, []byte(id)...)
}

// StoredDataByKeyKey returns the key for stored data by key
func StoredDataByKeyKey(dataKey string) []byte {
	return append(StoredDataKeyPrefix, []byte("key:"+dataKey)...)
}

// StoredDataByStoreKey returns the key for stored data by store
func StoredDataByStoreKey(storeID, dataID string) []byte {
	return append(append(StoredDataKeyPrefix, []byte("store:"+storeID)...), []byte(dataID)...)
}

// StoredDataByTypeKey returns the key for stored data by content type
func StoredDataByTypeKey(contentType, dataID string) []byte {
	return append(append(StoredDataKeyPrefix, []byte("type:"+contentType)...), []byte(dataID)...)
}

// StoredDataByTagKey returns the key for stored data by tag
func StoredDataByTagKey(tagKey, tagValue, dataID string) []byte {
	return append(append(StoredDataKeyPrefix, []byte("tag:"+tagKey+":"+tagValue)...), []byte(dataID)...)
}

// StoredDataByTimestampKey returns the key for stored data by timestamp
func StoredDataByTimestampKey(timestamp time.Time, dataID string) []byte {
	return append(append(StoredDataKeyPrefix, []byte("time:"+timestamp.Format(time.RFC3339))...), []byte(dataID)...)
}

// StoredDataBySizeKey returns the key for stored data by size
func StoredDataBySizeKey(minSize, maxSize int64, dataID string) []byte {
	sizeRange := fmt.Sprintf("%d-%d", minSize, maxSize)
	return append(append(StoredDataKeyPrefix, []byte("size:"+sizeRange)...), []byte(dataID)...)
}

// StoreByNameKey returns the key for stores by name
func StoreByNameKey(name string) []byte {
	return append(StoreKeyPrefix, []byte("name:"+name)...)
}

// StoreByTypeKey returns the key for stores by type
func StoreByTypeKey(storeType, storeID string) []byte {
	return append(append(StoreKeyPrefix, []byte("type:"+storeType)...), []byte(storeID)...)
}

// StoreByStatusKey returns the key for stores by status
func StoreByStatusKey(status, storeID string) []byte {
	return append(append(StoreKeyPrefix, []byte("status:"+status)...), []byte(storeID)...)
}

// StoreByCreatedTimeKey returns the key for stores by creation time
func StoreByCreatedTimeKey(createdAt time.Time, storeID string) []byte {
	return append(append(StoreKeyPrefix, []byte("created:"+createdAt.Format(time.RFC3339))...), []byte(storeID)...)
}

// StoreBySizeKey returns the key for stores by size
func StoreBySizeKey(minSize, maxSize int64, storeID string) []byte {
	sizeRange := fmt.Sprintf("%d-%d", minSize, maxSize)
	return append(append(StoreKeyPrefix, []byte("size:"+sizeRange)...), []byte(storeID)...)
}

// BackupByStoreKey returns the key for backups by store
func BackupByStoreKey(storeID, backupID string) []byte {
	return append(append(BackupKeyPrefix, []byte("store:"+storeID)...), []byte(backupID)...)
}

// BackupByStatusKey returns the key for backups by status
func BackupByStatusKey(status, backupID string) []byte {
	return append(append(BackupKeyPrefix, []byte("status:"+status)...), []byte(backupID)...)
}

// BackupByCreatedTimeKey returns the key for backups by creation time
func BackupByCreatedTimeKey(createdAt time.Time, backupID string) []byte {
	return append(append(BackupKeyPrefix, []byte("created:"+createdAt.Format(time.RFC3339))...), []byte(backupID)...)
}

// BackupBySizeKey returns the key for backups by size
func BackupBySizeKey(minSize, maxSize int64, backupID string) []byte {
	sizeRange := fmt.Sprintf("%d-%d", minSize, maxSize)
	return append(append(BackupKeyPrefix, []byte("size:"+sizeRange)...), []byte(backupID)...)
}

// RestoreByBackupKey returns the key for restores by backup
func RestoreByBackupKey(backupID, restoreID string) []byte {
	return append(append(RestoreKeyPrefix, []byte("backup:"+backupID)...), []byte(restoreID)...)
}

// RestoreByStoreKey returns the key for restores by store
func RestoreByStoreKey(storeID, restoreID string) []byte {
	return append(append(RestoreKeyPrefix, []byte("store:"+storeID)...), []byte(restoreID)...)
}

// RestoreByStatusKey returns the key for restores by status
func RestoreByStatusKey(status, restoreID string) []byte {
	return append(append(RestoreKeyPrefix, []byte("status:"+status)...), []byte(restoreID)...)
}

// RestoreByCreatedTimeKey returns the key for restores by creation time
func RestoreByCreatedTimeKey(createdAt time.Time, restoreID string) []byte {
	return append(append(RestoreKeyPrefix, []byte("created:"+createdAt.Format(time.RFC3339))...), []byte(restoreID)...)
}

// StoreIndexByStoreKey returns the key for indexes by store
func StoreIndexByStoreKey(storeID, indexID string) []byte {
	return append(append(StoreIndexKeyPrefix, []byte("store:"+storeID)...), []byte(indexID)...)
}

// StoreIndexByTypeKey returns the key for indexes by type
func StoreIndexByTypeKey(indexType, indexID string) []byte {
	return append(append(StoreIndexKeyPrefix, []byte("type:"+indexType)...), []byte(indexID)...)
}

// StoreIndexByStatusKey returns the key for indexes by status
func StoreIndexByStatusKey(status, indexID string) []byte {
	return append(append(StoreIndexKeyPrefix, []byte("status:"+status)...), []byte(indexID)...)
}

// StoreQueryByStoreKey returns the key for queries by store
func StoreQueryByStoreKey(storeID, queryID string) []byte {
	return append(append(StoreQueryKeyPrefix, []byte("store:"+storeID)...), []byte(queryID)...)
}

// StoreQueryByTypeKey returns the key for queries by type
func StoreQueryByTypeKey(queryType, queryID string) []byte {
	return append(append(StoreQueryKeyPrefix, []byte("type:"+queryType)...), []byte(queryID)...)
}

// StoreQueryByCreatedTimeKey returns the key for queries by creation time
func StoreQueryByCreatedTimeKey(createdAt time.Time, queryID string) []byte {
	return append(append(StoreQueryKeyPrefix, []byte("created:"+createdAt.Format(time.RFC3339))...), []byte(queryID)...)
}

// TransactionByStoreKey returns the key for transactions by store
func TransactionByStoreKey(storeID, transactionID string) []byte {
	return append(append(TransactionKeyPrefix, []byte("store:"+storeID)...), []byte(transactionID)...)
}

// TransactionByTypeKey returns the key for transactions by type
func TransactionByTypeKey(transactionType, transactionID string) []byte {
	return append(append(TransactionKeyPrefix, []byte("type:"+transactionType)...), []byte(transactionID)...)
}

// TransactionByStatusKey returns the key for transactions by status
func TransactionByStatusKey(status, transactionID string) []byte {
	return append(append(TransactionKeyPrefix, []byte("status:"+status)...), []byte(transactionID)...)
}

// TransactionByCreatedTimeKey returns the key for transactions by creation time
func TransactionByCreatedTimeKey(createdAt time.Time, transactionID string) []byte {
	return append(append(TransactionKeyPrefix, []byte("created:"+createdAt.Format(time.RFC3339))...), []byte(transactionID)...)
}

// GetStoredDataIDFromKey extracts stored data ID from key
func GetStoredDataIDFromKey(key []byte) string {
	if !bytes.HasPrefix(key, StoredDataKeyPrefix) {
		return ""
	}
	return string(key[len(StoredDataKeyPrefix):])
}

// GetStoreIDFromKey extracts store ID from key
func GetStoreIDFromKey(key []byte) string {
	if !bytes.HasPrefix(key, StoreKeyPrefix) {
		return ""
	}
	return string(key[len(StoreKeyPrefix):])
}

// GetBackupIDFromKey extracts backup ID from key
func GetBackupIDFromKey(key []byte) string {
	if !bytes.HasPrefix(key, BackupKeyPrefix) {
		return ""
	}
	return string(key[len(BackupKeyPrefix):])
}

// GetRestoreIDFromKey extracts restore ID from key
func GetRestoreIDFromKey(key []byte) string {
	if !bytes.HasPrefix(key, RestoreKeyPrefix) {
		return ""
	}
	return string(key[len(RestoreKeyPrefix):])
}

// GetStoreIndexIDFromKey extracts store index ID from key
func GetStoreIndexIDFromKey(key []byte) string {
	if !bytes.HasPrefix(key, StoreIndexKeyPrefix) {
		return ""
	}
	return string(key[len(StoreIndexKeyPrefix):])
}

// GetStoreQueryIDFromKey extracts store query ID from key
func GetStoreQueryIDFromKey(key []byte) string {
	if !bytes.HasPrefix(key, StoreQueryKeyPrefix) {
		return ""
	}
	return string(key[len(StoreQueryKeyPrefix):])
}

// GetTransactionIDFromKey extracts transaction ID from key
func GetTransactionIDFromKey(key []byte) string {
	if !bytes.HasPrefix(key, TransactionKeyPrefix) {
		return ""
	}
	return string(key[len(TransactionKeyPrefix):])
}
