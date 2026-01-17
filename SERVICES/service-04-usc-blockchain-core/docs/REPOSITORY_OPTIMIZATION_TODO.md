# Repository Layer Optimization TODO - 12 Groups

## 📋 Tổng Quan

TODO list để tối ưu Repository layer theo 12 groups, focus vào:
1. **Silent Failures**: Fix `_, _ = postgres.ExecContext` (ignore errors)
2. **Async Saves**: Review và convert critical ones to sync

**Ngày tạo**: $(date +%Y-%m-%d)

---

## 🎯 Optimization Strategy

### Priority Levels
- **HIGH**: Silent failures (should fix)
- **MEDIUM**: Critical async saves (should review)
- **LOW**: Non-critical async saves (optional)

### Pattern to Apply
```go
// BEFORE (Silent Failure)
_, _ = postgres.ExecContext(ctx, query, ...)

// AFTER (With Error Handling)
if _, err := postgres.ExecContext(ctx, query, ...); err != nil {
    r.logger.Error("Failed to save to database",
        logging.Error(err),
        logging.String("operation", "save_analytics"))
    // Continue even if database save fails (keeper is primary)
}
```

---

## 📊 TODO List by Groups

### ✅ GROUP 1: Block Operations
**File**: `block_operations_repository.go`

**Issues Found**:
- ⚠️ **Async Save**: `saveBlockToDatabase` (line ~59) - `go func()` async save

**TODO**:
- [ ] **OPTIONAL**: Review `saveBlockToDatabase` async save
  - Current: `go func()` async save
  - Decision: Keep async (analytics only, non-critical) OR convert to sync
  - Location: `ProduceBlock` method

**Status**: ✅ **LOW PRIORITY** (analytics only)

---

### ✅ GROUP 2: Transaction Operations
**File**: `transaction_operations_repository.go`

**Issues Found**:
- ⚠️ **Async Save**: `saveTransactionToDatabase` (line ~61) - `go func()` async save

**TODO**:
- [ ] **OPTIONAL**: Review `saveTransactionToDatabase` async save
  - Current: `go func()` async save
  - Decision: Keep async (analytics only, non-critical) OR convert to sync
  - Location: `SubmitTransaction` method

**Status**: ✅ **LOW PRIORITY** (analytics only)

---

### ⚠️ GROUP 3: USC Coin Operations
**File**: `usc_coin_operations_repository.go`

**Issues Found**:
- 🔴 **Silent Failure**: 2 instances (lines ~348-349) - `_, _ = postgres.ExecContext`
- ⚠️ **Async Save**: `saveTransferToDatabase` (line ~110) - `go func()` async save

**TODO**:
- [ ] **HIGH**: Fix silent failures in `saveTransferToDatabase`
  - Location: Lines ~348-349
  - Fix: Add error handling và logging
  - Impact: Better error tracking cho transfer operations
- [ ] **OPTIONAL**: Review `saveTransferToDatabase` async save
  - Current: `go func()` async save
  - Decision: Keep async (analytics only) OR convert to sync

**Status**: ⚠️ **HIGH PRIORITY** (silent failures)

---

### ⚠️ GROUP 4: Smart Contract Operations
**File**: `smart_contract_operations_repository.go`

**Issues Found**:
- 🔴 **Silent Failure**: 1 instance (line ~521) - `_, _ = postgres.ExecContext`
- ⚠️ **Async Save**: `saveContractExecutionToDatabase` (line ~77) - `go func()` async save

**TODO**:
- [ ] **HIGH**: Fix silent failure in contract execution save
  - Location: Line ~521
  - Fix: Add error handling và logging
  - Impact: Better error tracking cho contract execution analytics
- [ ] **OPTIONAL**: Review `saveContractExecutionToDatabase` async save
  - Current: `go func()` async save
  - Decision: Keep async (analytics only) OR convert to sync

**Status**: ⚠️ **HIGH PRIORITY** (silent failure)

---

### ✅ GROUP 5: NFT Token Operations
**File**: `nft_token_operations_repository.go`

**Issues Found**:
- ✅ **No Issues**: Already fixed in previous optimization

**TODO**:
- [x] ✅ **COMPLETED**: All async saves converted to sync với error handling

**Status**: ✅ **NO ACTION NEEDED**

---

### ✅ GROUP 6: Custom Token Operations
**File**: `custom_token_operations_repository.go`

**Issues Found**:
- ✅ **No Issues**: No silent failures or async saves found

**TODO**:
- [x] ✅ **COMPLETED**: No issues found

**Status**: ✅ **NO ACTION NEEDED**

---

### ✅ GROUP 7: Product Certificate Operations
**File**: `product_certificate_operations_repository.go`

**Issues Found**:
- ✅ **No Issues**: Already fixed in previous optimization

**TODO**:
- [x] ✅ **COMPLETED**: All async saves converted to sync với error handling

**Status**: ✅ **NO ACTION NEEDED**

---

### ⚠️ GROUP 8: Validator Operations
**File**: `validator_operations_repository.go`

**Issues Found**:
- 🔴 **Silent Failure**: 1 instance (line ~303) - `_, _ = postgres.ExecContext`

**TODO**:
- [ ] **HIGH**: Fix silent failure in validator save
  - Location: Line ~303
  - Fix: Add error handling và logging
  - Impact: Better error tracking cho validator analytics
  - Note: Already has sync saves for `saveValidatorToDatabase` và `saveStakingToDatabase` (fixed previously)

**Status**: ⚠️ **HIGH PRIORITY** (silent failure)

---

### ✅ GROUP 9: Network Operations
**File**: `network_operations_repository.go`

**Issues Found**:
- ✅ **No Issues**: No silent failures or async saves found

**TODO**:
- [x] ✅ **COMPLETED**: No issues found

**Status**: ✅ **NO ACTION NEEDED**

---

### ✅ GROUP 10: Streaming Operations
**File**: `streaming_operations_repository.go`

**Issues Found**:
- ✅ **No Issues**: No async saves (streaming operations don't need database saves)

**TODO**:
- [x] ✅ **COMPLETED**: No issues found

**Status**: ✅ **NO ACTION NEEDED**

---

### ⚠️ GROUP 11: Store Bridge Operations
**File**: `store_bridge_operations_repository.go`

**Issues Found**:
- 🔴 **Silent Failure**: 2 instances (lines ~213, ~744) - `_, _ = postgres.ExecContext`
- ⚠️ **Async Saves**: 4 instances (lines ~71, ~240, ~325, ~445) - `go func()` async saves

**TODO**:
- [ ] **HIGH**: Fix silent failures in bridge operations
  - Location: Lines ~213, ~744
  - Fix: Add error handling và logging
  - Impact: Better error tracking cho bridge operations
- [ ] **OPTIONAL**: Review 4 async saves
  - Locations: `DeployStoreBridge`, `RegisterStoreNetwork`, `BridgeStoreTokenToUSC`, `BridgeUSCToStoreToken`
  - Decision: Keep async (analytics only) OR convert critical ones to sync

**Status**: ⚠️ **HIGH PRIORITY** (silent failures)

---

### ⚠️ GROUP 12: Store Network Operations
**File**: `store_network_operations_repository.go`

**Issues Found**:
- 🔴 **Silent Failure**: 1 instance (line ~201) - `_, _ = postgres.ExecContext`
- ⚠️ **Async Saves**: 2 instances (lines ~63, ~282) - `go func()` async saves

**TODO**:
- [ ] **HIGH**: Fix silent failure in network sync save
  - Location: Line ~201
  - Fix: Add error handling và logging
  - Impact: Better error tracking cho network sync analytics
- [ ] **OPTIONAL**: Review 2 async saves
  - Locations: `SyncStoreNetworkState`, `UpdateStoreBridgeConfig`
  - Decision: Keep async (analytics only) OR convert to sync

**Status**: ⚠️ **HIGH PRIORITY** (silent failure)

---

## 📊 Summary Statistics

### Issues by Priority

| Priority | Count | Groups |
|----------|-------|--------|
| **HIGH** (Silent Failures) | **7** | Groups 3, 4, 8, 11, 12 |
| **MEDIUM** (Critical Async) | **0** | - |
| **LOW** (Optional Async) | **9** | Groups 1, 2, 3, 4, 11, 12 |

### Groups Status

| Status | Count | Groups |
|--------|-------|--------|
| ✅ **NO ACTION** | **6** | Groups 5, 6, 7, 9, 10 |
| ⚠️ **NEEDS FIX** | **5** | Groups 3, 4, 8, 11, 12 |
| ✅ **OPTIONAL** | **1** | Group 1, 2 |

---

## 🎯 Action Plan

### Phase 1: Fix Silent Failures (HIGH PRIORITY)
**Effort**: ~2-3 hours
**Groups**: 3, 4, 8, 11, 12

1. **GROUP 3**: Fix 2 silent failures in `usc_coin_operations_repository.go`
2. **GROUP 4**: Fix 1 silent failure in `smart_contract_operations_repository.go`
3. **GROUP 8**: Fix 1 silent failure in `validator_operations_repository.go`
4. **GROUP 11**: Fix 2 silent failures in `store_bridge_operations_repository.go`
5. **GROUP 12**: Fix 1 silent failure in `store_network_operations_repository.go`

### Phase 2: Review Async Saves (OPTIONAL)
**Effort**: ~3-4 hours
**Groups**: 1, 2, 3, 4, 11, 12

1. Review each async save
2. Decide: sync vs async based on use case
3. Convert critical ones to sync

---

## 📝 Implementation Pattern

### Fix Silent Failure Pattern
```go
// BEFORE
_, _ = postgres.ExecContext(ctx, query, args...)

// AFTER
if _, err := postgres.ExecContext(ctx, query, args...); err != nil {
    r.logger.Error("Failed to save to database",
        logging.Error(err),
        logging.String("operation", "save_analytics"),
        logging.String("table", "table_name"))
    // Continue even if database save fails (keeper is primary)
}
```

### Convert Async to Sync Pattern
```go
// BEFORE
go func() {
    if r.db != nil {
        r.saveToDatabase(context.Background(), data)
    }
}()

// AFTER
if r.db != nil {
    if err := r.saveToDatabase(ctx, data); err != nil {
        r.logger.Error("Failed to save to database",
            logging.Error(err),
            logging.String("operation", "save_analytics"))
        // Continue even if database save fails (keeper is primary)
    }
}
```

---

## ✅ Verification Checklist

After completing optimizations:

- [ ] All silent failures fixed (7 instances)
- [ ] Error handling added to all database saves
- [ ] Logging added for all errors
- [ ] Async saves reviewed và converted if needed
- [ ] Linter: 0 errors
- [ ] Tests: All pass
- [ ] Code review: No regressions

---

## 🎯 Priority Order

1. **HIGH**: Fix silent failures (Groups 3, 4, 8, 11, 12)
2. **OPTIONAL**: Review async saves (Groups 1, 2, 3, 4, 11, 12)

---

**Status**: 📋 **TODO CREATED - READY FOR IMPLEMENTATION**

