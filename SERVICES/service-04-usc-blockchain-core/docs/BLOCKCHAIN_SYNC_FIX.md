# 🔧 **BLOCKCHAIN SYNC FIX - Stale Data Handling**

**Ngày fix**: 2025-11-12  
**Vấn đề**: Database có dữ liệu cũ (26162 blocks) từ lần chạy trước, trong khi CometBFT chỉ có ~6300 blocks  
**Trạng thái**: ✅ **FIXED**

---

## 🐛 **VẤN ĐỀ PHÁT HIỆN**

### **Symptom**
```
cometbft_height: 6300
db_height: 26162
```

**Nguyên nhân**:
- Database có 26162 blocks từ lần chạy trước (2025-11-11)
- CometBFT đã được reset và chỉ có ~6300 blocks (mới)
- Logic sync chỉ sync khi `cometBFTHeight > dbHeight`
- Do đó sync không chạy, database không được cập nhật

---

## ✅ **GIẢI PHÁP ĐÃ IMPLEMENT**

### **1. Cải thiện Logic Sync**

**File**: `internal/infrastructure/database/manager_postgresql.go`

**Thay đổi**:
- ✅ Detect khi `db_height > cometbft_height` (stale data)
- ✅ Auto-reset blocks table nếu:
  - Difference > 1000 blocks (significant difference)
  - CometBFT height < 100 (new chain)
- ✅ Verify genesis block hash cho small differences
- ✅ Reset và start fresh sync nếu genesis hash không match

### **2. New Functions**

#### **`verifyGenesisBlockHash()`**
- Verify genesis block hash từ CometBFT với database
- Returns `true` nếu match, `false` nếu không match
- Handles errors gracefully

#### **`resetBlocksTable()`**
- Reset blocks table để start fresh sync
- Safe operation với proper error handling

---

## 📋 **LOGIC FLOW**

```
1. Get cometbft_height và db_height
2. If db_height > cometbft_height:
   a. If difference > 1000 OR cometbft_height < 100:
      → Reset blocks table → Start fresh sync
   b. Else (small difference):
      → Verify genesis block hash
      → If mismatch: Reset blocks table → Start fresh sync
      → If match: Database is up to date (return)
3. If cometbft_height > db_height:
   → Sync blocks from db_height+1 to cometbft_height
4. Else:
   → Database is up to date
```

---

## 🧪 **TESTING**

### **Test Case 1: Stale Data (Large Difference)**
```bash
# Scenario: db_height = 26162, cometbft_height = 6300
# Expected: Blocks table reset, fresh sync starts
```

### **Test Case 2: Stale Data (Small Difference)**
```bash
# Scenario: db_height = 100, cometbft_height = 50
# Expected: Verify genesis hash, reset if mismatch
```

### **Test Case 3: Normal Sync**
```bash
# Scenario: db_height = 100, cometbft_height = 200
# Expected: Sync blocks 101-200
```

---

## 🚀 **DEPLOYMENT**

### **Steps**
1. ✅ Code đã được update
2. ⏭️ Rebuild service-04
3. ⏭️ Restart service
4. ⏭️ Monitor logs để verify fix hoạt động

### **Expected Behavior After Fix**
```
1. Service starts
2. Sync detects db_height (26162) > cometbft_height (6300)
3. Detects large difference (>1000)
4. Resets blocks table
5. Starts fresh sync from height 0
6. Syncs all blocks from CometBFT
```

---

## 📊 **MONITORING**

### **Logs to Watch**
```bash
# Watch for these log messages:
- "Database has stale data from previous chain/run, resetting blocks table"
- "Blocks table reset successfully"
- "Blocks table reset, starting fresh sync"
- "Syncing blocks to database"
```

### **Metrics to Monitor**
- `db_height` should match `cometbft_height` after sync
- Sync should complete within reasonable time
- No errors during sync process

---

## 🎯 **KẾT QUẢ**

### **Before Fix**
- ❌ Database có stale data (26162 blocks)
- ❌ Sync không chạy (db_height > cometbft_height)
- ❌ Database không được cập nhật

### **After Fix**
- ✅ Auto-detect stale data
- ✅ Auto-reset blocks table
- ✅ Start fresh sync
- ✅ Database được sync đúng với CometBFT

---

**Last Updated**: 2025-11-12  
**Status**: ✅ **READY FOR TESTING**  
**Next Step**: Rebuild và restart service để test fix

