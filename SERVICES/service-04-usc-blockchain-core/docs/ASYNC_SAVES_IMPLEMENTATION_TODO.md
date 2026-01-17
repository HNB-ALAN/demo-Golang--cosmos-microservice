# Async Saves Implementation TODO

## 📋 Tổng Quan

TODO list chi tiết để optimize async saves trong Repository layer theo 12 groups.

**Ngày tạo**: $(date +%Y-%m-%d)

---

## 🎯 Implementation Strategy

### Pattern 1: Improve Async (Keep async + add error handling)
**Groups**: 1, 2
**Reason**: High volume, non-critical, industry standard

### Pattern 2: Convert to Sync (Critical operations)
**Groups**: 3, 4 (ExecuteContract), 11, 12
**Reason**: Critical operations cần data persistence guarantee

---

## 📋 TODO List

### ✅ TODO 1: GROUP 1 - Block Operations
**Status**: ⏳ PENDING
**Priority**: MEDIUM
**Type**: Improve Async

**File**: `block_operations_repository.go`
**Method**: `ProduceBlock`
**Current**: Async save không có error handling

**Tasks**:
- [ ] Add error handling cho `saveBlockToDatabase`
- [ ] Add context với timeout (5 seconds)
- [ ] Add correlation ID logging
- [ ] Add retry logic (optional)

**Implementation**:
```go
// BEFORE
go func() {
    if r.db != nil {
        r.saveBlockToDatabase(context.Background(), result)
    }
}()

// AFTER
go func() {
    if r.db != nil {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        
        if err := r.saveBlockToDatabase(ctx, result); err != nil {
            r.logger.Error("Failed to save block analytics (async)",
                logging.Error(err),
                logging.String("block_hash", result.BlockHash),
                logging.String("correlation_id", utils.GetCorrelationID(ctx)))
            // Optional: Retry logic hoặc dead letter queue
        }
    }
}()
```

---

### ✅ TODO 2: GROUP 2 - Transaction Operations
**Status**: ⏳ PENDING
**Priority**: MEDIUM
**Type**: Improve Async

**File**: `transaction_operations_repository.go`
**Method**: `SubmitTransaction`
**Current**: Async save không có error handling

**Tasks**:
- [ ] Add error handling cho `saveTransactionToDatabase`
- [ ] Add context với timeout (5 seconds)
- [ ] Add correlation ID logging
- [ ] Add retry logic (optional)

**Implementation**:
```go
// BEFORE
go func() {
    if r.db != nil {
        r.saveTransactionToDatabase(context.Background(), txHash, req)
    }
}()

// AFTER
go func() {
    if r.db != nil {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        
        if err := r.saveTransactionToDatabase(ctx, txHash, req); err != nil {
            r.logger.Error("Failed to save transaction analytics (async)",
                logging.Error(err),
                logging.String("transaction_hash", txHash),
                logging.String("correlation_id", utils.GetCorrelationID(ctx)))
            // Optional: Retry logic hoặc dead letter queue
        }
    }
}()
```

---

### ⚠️ TODO 3: GROUP 3 - USC Coin Operations
**Status**: ⏳ PENDING
**Priority**: HIGH
**Type**: Convert to Sync

**File**: `usc_coin_operations_repository.go`
**Method**: `TransferUSC`
**Current**: Async save
**Reason**: Financial data, audit trail important

**Tasks**:
- [ ] Convert `saveTransferToDatabase` từ async → sync
- [ ] Add error handling và logging
- [ ] Update function signature to return error
- [ ] Add correlation ID logging

**Implementation**:
```go
// BEFORE
go func() {
    if r.db != nil {
        r.saveTransferToDatabase(context.Background(), req, result)
    }
}()

// AFTER
if r.db != nil {
    if err := r.saveTransferToDatabase(ctx, req, result); err != nil {
        r.logger.Error("Failed to save transfer analytics",
            logging.Error(err),
            logging.String("from_address", req.FromAddress),
            logging.String("to_address", req.ToAddress),
            logging.String("correlation_id", utils.GetCorrelationID(ctx)))
        // Continue even if database save fails (keeper is primary)
    }
}
```

**Function Update**:
```go
// Update function signature
func (r *Repository) saveTransferToDatabase(ctx context.Context, req *proto.TransferUSCBlockchainRequest, resp *proto.TransferUSCBlockchainResponse) error {
    // ... existing code ...
    if _, err := postgres.ExecContext(ctx, query, ...); err != nil {
        return fmt.Errorf("failed to save transfer analytics: %w", err)
    }
    return nil
}
```

---

### ⚠️ TODO 4: GROUP 4 - Smart Contract Operations
**Status**: ⏳ PENDING
**Priority**: HIGH
**Type**: Convert to Sync

**File**: `smart_contract_operations_repository.go`
**Method**: `ExecuteContract`
**Current**: Async save
**Reason**: Debugging important

**Tasks**:
- [ ] Convert `saveContractExecutionToDatabase` từ async → sync
- [ ] Add error handling và logging
- [ ] Update function signature to return error
- [ ] Add correlation ID logging

**Implementation**:
```go
// BEFORE
go func() {
    if r.db != nil {
        r.saveContractExecutionToDatabase(context.Background(), req, result)
    }
}()

// AFTER
if r.db != nil {
    if err := r.saveContractExecutionToDatabase(ctx, req, result); err != nil {
        r.logger.Error("Failed to save contract execution analytics",
            logging.Error(err),
            logging.String("contract_address", req.ContractAddress),
            logging.String("correlation_id", utils.GetCorrelationID(ctx)))
        // Continue even if database save fails (keeper is primary)
    }
}
```

**Function Update**:
```go
// Update function signature
func (r *Repository) saveContractExecutionToDatabase(ctx context.Context, req *proto.ExecuteContractRequest, resp *proto.ExecuteContractResponse) error {
    // ... existing code ...
    if _, err := postgres.ExecContext(ctx, query, ...); err != nil {
        return fmt.Errorf("failed to save contract execution analytics: %w", err)
    }
    return nil
}
```

---

### ⚠️ TODO 5: GROUP 11 - Store Bridge Operations
**Status**: ⏳ PENDING
**Priority**: HIGH
**Type**: Convert to Sync

**File**: `store_bridge_operations_repository.go`
**Methods**: 4 instances
- `DeployStoreBridge`
- `RegisterStoreNetwork`
- `BridgeStoreTokenToUSC`
- `BridgeUSCToStoreToken`

**Current**: Async saves
**Reason**: Cross-chain tracking important

**Tasks**:
- [ ] Convert `saveBridgeToDatabase` từ async → sync
- [ ] Convert `saveNetworkToDatabase` từ async → sync
- [ ] Convert `saveBridgeTransactionToDatabase` từ async → sync (2 instances)
- [ ] Add error handling và logging cho tất cả
- [ ] Update function signatures to return error

**Implementation Pattern**:
```go
// BEFORE (for each method)
go func() {
    if r.db != nil {
        r.saveBridgeToDatabase(context.Background(), req, result)
    }
}()

// AFTER (for each method)
if r.db != nil {
    if err := r.saveBridgeToDatabase(ctx, req, result); err != nil {
        r.logger.Error("Failed to save bridge analytics",
            logging.Error(err),
            logging.String("bridge_address", result.BridgeAddress),
            logging.String("correlation_id", utils.GetCorrelationID(ctx)))
        // Continue even if database save fails (keeper is primary)
    }
}
```

---

### ⚠️ TODO 6: GROUP 12 - Store Network Operations
**Status**: ⏳ PENDING
**Priority**: HIGH
**Type**: Convert to Sync

**File**: `store_network_operations_repository.go`
**Methods**: 2 instances
- `SyncStoreNetworkState`
- `UpdateStoreBridgeConfig`

**Current**: Async saves
**Reason**: Network health tracking important

**Tasks**:
- [ ] Convert `saveSyncStateToDatabase` từ async → sync
- [ ] Convert `saveBridgeConfigToDatabase` từ async → sync
- [ ] Add error handling và logging cho cả 2
- [ ] Update function signatures to return error

**Implementation Pattern**:
```go
// BEFORE (for each method)
go func() {
    if r.db != nil {
        r.saveSyncStateToDatabase(context.Background(), req, result)
    }
}()

// AFTER (for each method)
if r.db != nil {
    if err := r.saveSyncStateToDatabase(ctx, req, result); err != nil {
        r.logger.Error("Failed to save network sync analytics",
            logging.Error(err),
            logging.String("network_id", req.NetworkId),
            logging.String("correlation_id", utils.GetCorrelationID(ctx)))
        // Continue even if database save fails (keeper is primary)
    }
}
```

---

## 📊 Summary

### Tasks by Priority

| Priority | Count | Groups |
|----------|-------|--------|
| **HIGH** | 4 | Groups 3, 4, 11, 12 |
| **MEDIUM** | 2 | Groups 1, 2 |

### Tasks by Type

| Type | Count | Groups |
|------|-------|--------|
| **Improve Async** | 2 | Groups 1, 2 |
| **Convert to Sync** | 4 | Groups 3, 4, 11, 12 |

### Total Instances

- **Improve Async**: 2 instances
- **Convert to Sync**: 7 instances (1 + 1 + 4 + 2)
- **Total**: 9 instances

---

## 🎯 Implementation Order

### Phase 1: Convert Critical to Sync (HIGH PRIORITY)
1. GROUP 3: USC Coin Operations
2. GROUP 4: Smart Contract Execution
3. GROUP 11: Store Bridge Operations
4. GROUP 12: Store Network Operations

### Phase 2: Improve Async Error Handling (MEDIUM PRIORITY)
1. GROUP 1: Block Operations
2. GROUP 2: Transaction Operations

---

## ✅ Verification Checklist

After completing each TODO:

- [ ] Function signature updated (if converting to sync)
- [ ] Error handling added
- [ ] Logging added với correlation ID
- [ ] Context với timeout (for async)
- [ ] Linter: 0 errors
- [ ] Tests: All pass
- [ ] Code review: No regressions

---

**Status**: 📋 **TODO CREATED - READY FOR IMPLEMENTATION**

