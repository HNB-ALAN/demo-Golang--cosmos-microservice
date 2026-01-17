# 🗄️ **GIẢI THÍCH CÁC LOẠI DATABASE TRONG SERVICE-04**

## 📋 **TỔNG QUAN**

Service-04 sử dụng **4 loại database** với mục đích khác nhau, được tổ chức thành **2 layers**:

1. **APPLICATION LAYER** - Business Logic
2. **BLOCKCHAIN LAYER** - Cosmos SDK & Consensus

---

## 🏗️ **1. APPLICATION LAYER DATABASES**

### **A. PostgreSQL (Business Logic Database)**

**Mục đích:**
- Lưu trữ dữ liệu business (blocks, transactions, smart contracts, NFTs, etc.)
- Analytics và reporting
- Complex queries với JOIN operations

**Cấu hình:**
```yaml
Database: blockchain_db
Host: postgres:5432
Tables: 50+ tables
```

**Sử dụng trong code:**
```go
// cmd/main.go
postgresManager, err := database.NewPostgreSQLManager(
    app.GetConfig(), 
    *app.GetLogger(), 
    cosmosApp, 
    blockchainStorage
)
```

**Ví dụ tables:**
- `blocks` - Block metadata
- `transactions` - Transaction records
- `smart_contracts` - Smart contract data
- `nft_collections` - NFT collections
- `nfts` - Individual NFTs
- `custom_tokens` - Custom token data
- `product_certificates` - Product certificates
- `validators` - Validator information
- `staking` - Staking data

---

### **B. Redis (Cache & Real-time Data)**

**Mục đích:**
- High-speed caching (<1ms response)
- Real-time data storage
- Session management
- Cosmos SDK module cache

**Cấu hình:**
```yaml
Host: redis:6379
Database: 0-15 (multiple databases)
Performance: Sub-millisecond access
```

**Sử dụng trong code:**
```go
// cmd/main.go
redisManager, err := database.NewRedisManager(
    app.GetConfig(), 
    *app.GetLogger()
)
```

**Use cases:**
- Cache block data để giảm database queries
- Real-time metrics và monitoring
- Session storage
- Temporary data (TTL-based)

---

## 🔗 **2. BLOCKCHAIN LAYER DATABASES**

### **A. RocksDB (Business Logic - Blockchain State)**

**Mục đích:**
- High-performance blockchain state storage
- Business logic data (không phải Cosmos SDK state)
- Key-value storage với performance cao

**Cấu hình:**
```yaml
Path: ./data/rocksdb
Backend: RocksDB (CGO required)
Performance: <10ms state access
```

**Sử dụng trong code:**
```go
// cmd/main.go
rocksDBManager, err := storage.NewRocksDBManager(
    storage.DefaultRocksDBConfig()
)
blockchainStorage := storage.NewStateManager(rocksDBManager)
```

**Lưu ý:**
- RocksDB này **KHÔNG phải** cho Cosmos SDK
- Dùng cho business logic layer
- Cần CGO để compile

---

### **B. RocksDB (Cosmos SDK State Database)**

**Mục đích:**
- Cosmos SDK application state storage
- Block data, account states, keeper data
- **RocksDB là backend được khuyến nghị** cho production theo tài liệu Cosmos SDK

**Cấu hình:**
```yaml
Path: ./data/cosmos
Backend: RocksDB (via cometbft-db)
Purpose: Cosmos SDK BaseApp state
Build Tags: rocksdb_legacy rocksdb (required)
CGO: Enabled (required)
```

**Sử dụng trong code:**
```go
// block-chain-cosmos/app/factory.go
cmtDB, err := cmtdb.NewDB(
    "application", 
    cmtdb.RocksDBBackend,  // ← RocksDB (theo tài liệu chính thức)
    dbDir
)
cosmosDatabase := &cosmosDBAdapter{db: cmtDB}
cosmosApp := NewUSCApp(cosmosDatabase)
```

**✅ Cosmos SDK HỖ TRỢ ROCKSDB (theo tài liệu chính thức)**

Theo tài liệu chính thức của Cosmos SDK, **RocksDB được hỗ trợ** và là một trong các backend được khuyến nghị cho production.

**Yêu cầu để dùng RocksDB:**
1. ✅ **Build Tags**: `rocksdb_legacy rocksdb` trong Dockerfile
2. ✅ **CGO Enabled**: `CGO_ENABLED=1` khi build
3. ✅ **librocksdb**: Đã có trong Docker image (rocksdb-dev package)
4. ✅ **ABCIServer**: Cũng dùng RocksDB với build tags tương tự

**Lưu ý:**
- ✅ Cosmos SDK **HỖ TRỢ** RocksDB (theo tài liệu chính thức)
- ✅ Đã enable RocksDB với build tags trong Dockerfile
- ✅ RocksDB được dùng cho cả business logic (RocksDBManager) và Cosmos SDK state
- ✅ Performance tốt hơn GoLevelDB cho production workloads

---

## 📊 **SO SÁNH CÁC DATABASE**

| Database | Layer | Mục đích | Backend | Performance |
|----------|-------|----------|---------|-------------|
| **PostgreSQL** | Application | Business logic, analytics | SQL | <100ms queries |
| **Redis** | Application | Cache, real-time data | In-memory | <1ms access |
| **RocksDB** | Application | Blockchain state (business) | Key-value | <10ms access |
| **RocksDB** | Blockchain | Cosmos SDK state | Key-value | <10ms access |

---

## 🎯 **KHI NÀO DÙNG DATABASE NÀO?**

### **PostgreSQL:**
- ✅ Complex queries với JOIN
- ✅ Analytics và reporting
- ✅ Transactional data (ACID)
- ✅ Relational data (foreign keys)

### **Redis:**
- ✅ Cache để giảm database load
- ✅ Real-time data (sessions, metrics)
- ✅ Temporary data với TTL
- ✅ High-speed lookups

### **RocksDB (Business Logic):**
- ✅ High-performance blockchain state
- ✅ Business logic data storage
- ✅ Key-value operations
- ✅ Sequential writes

### **RocksDB (Cosmos SDK):**
- ✅ Cosmos SDK application state
- ✅ Block data trong Cosmos SDK
- ✅ Keeper state storage
- ✅ Production-ready performance (theo tài liệu Cosmos SDK)

---

## 🔄 **DATA FLOW**

```
┌─────────────────────────────────────────────────────────┐
│                    APPLICATION LAYER                      │
├─────────────────────────────────────────────────────────┤
│                                                           │
│  ┌──────────────┐      ┌──────────────┐                │
│  │  PostgreSQL  │      │    Redis     │                │
│  │ (Business DB)│      │   (Cache)    │                │
│  └──────────────┘      └──────────────┘                │
│         │                      │                         │
│         └──────────┬───────────┘                        │
│                    │                                     │
│         ┌──────────▼──────────┐                         │
│         │   RocksDB Manager   │                         │
│         │  (Business Logic)   │                         │
│         └──────────┬──────────┘                         │
└─────────────────────┼───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│                  BLOCKCHAIN LAYER                        │
├─────────────────────────────────────────────────────────┤
│                                                           │
│         ┌──────────────────────┐                        │
│         │   GoLevelDB         │                        │
│         │ (Cosmos SDK State)  │                        │
│         └──────────────────────┘                        │
│                    │                                     │
│         ┌──────────▼──────────┐                         │
│         │   Cosmos SDK App   │                         │
│         │  (BaseApp, Keepers) │                        │
│         └─────────────────────┘                        │
└─────────────────────────────────────────────────────────┘
```

---

## 📝 **TÓM TẮT**

1. **PostgreSQL**: Business logic, analytics, complex queries
2. **Redis**: Cache, real-time data, sessions
3. **RocksDB**: Business logic blockchain state (Application layer)
4. **RocksDB**: Cosmos SDK application state (Blockchain layer)

**Lưu ý quan trọng:**
- ✅ **Cosmos SDK HỖ TRỢ ROCKSDB** (theo tài liệu chính thức)
- ✅ **Đã enable RocksDB** với build tags `rocksdb_legacy rocksdb` trong Dockerfile
- ✅ RocksDB được dùng cho cả 2 layers (Application và Blockchain)
- ✅ Application layer: RocksDB cho business logic (RocksDBManager)
- ✅ Blockchain layer: RocksDB cho Cosmos SDK state (BaseApp, Keepers)
- ✅ Performance tốt hơn GoLevelDB cho production workloads

**Tài liệu chính thức Cosmos SDK:**
- Cosmos SDK hỗ trợ: **RocksDB, LevelDB, BadgerDB, Pebble**
- RocksDB là backend được khuyến nghị cho production
- RocksDB cần: Build tags `rocksdb`, CGO enabled, librocksdb installed
- ✅ **Đã cấu hình đầy đủ**: Build tags, CGO, và librocksdb đều có trong Dockerfile

