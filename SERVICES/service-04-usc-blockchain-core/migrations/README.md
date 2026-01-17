# 🗄️ USC Blockchain Core Service - Database Migrations

**⚠️ LƯU Ý: Đây chỉ là ví dụ template. Sau khi hoàn chỉnh service thì nó sẽ biến đổi theo cho khớp với yêu cầu chức năng từng service để tối ưu hoàn chỉnh.**

## 📁 **Directory Structure**

```
migrations/
├── README.md                           # This file
├── postgresql/                         # APPLICATION LAYER - PostgreSQL business logic migrations
│   ├── 001_create_blockchain_tables.up.sql
│   └── 001_create_blockchain_tables.down.sql
├── redis/                              # APPLICATION LAYER - Redis cache & real-time data migrations
│   ├── 001_initial_redis_schema.up.sh
│   ├── 001_initial_redis_schema.down.sh
│   ├── 002_cosmos_modules_cache.up.sh
│   └── 002_cosmos_modules_cache.down.sh
└── init-database.sh                    # Database initialization script
```

## 🎯 **Migration Strategy**

### **Up/Down Pattern**
- **`.up.*` files**: Apply migrations forward
- **`.down.*` files**: Rollback migrations backward
- **Version numbering**: Sequential versioning (001, 002, 003...)
- **Atomic operations**: Each migration is atomic and reversible

### **Two-Layer Architecture**

#### **APPLICATION LAYER** (`migrations/`)
- **Purpose**: Business logic, user data, analytics, caching
- **PostgreSQL**: Business tables, user transactions, NFT metadata
- **Redis**: Application cache, real-time data, Cosmos SDK module cache
- **Responsibility**: Application-specific data management

#### **BLOCKCHAIN LAYER** (RocksDB - Managed by Cosmos SDK)
- **Purpose**: Blockchain state storage (managed by Cosmos SDK)
- **RocksDB**: Blockchain state storage (managed automatically by Cosmos SDK)
- **Note**: PostgreSQL migrations for blockchain layer have been removed (not used)
- **Responsibility**: Blockchain state is managed by Cosmos SDK, not via migration scripts

### **Database Types**

#### 1. **PostgreSQL (Application Layer)** - Business Logic
- **Purpose**: Complex queries, analytics, reporting
- **Content**: User transactions, NFT metadata, business analytics
- **Format**: SQL DDL statements
- **Features**: ACID compliance, complex joins, indexing

#### 2. **RocksDB** - Blockchain State Storage
- **Purpose**: High-performance blockchain data storage
- **Content**: Blocks, transactions, account states, smart contracts
- **Format**: YAML configuration with RocksDB-specific settings
- **Performance**: Optimized for read/write operations

#### 3. **Redis (Application Layer)** - Application Cache & Real-time Data
- **Purpose**: High-speed caching and real-time data
- **Content**: Application cache, Cosmos SDK module cache, performance metrics
- **Format**: Shell scripts with Redis commands and error handling
- **Performance**: Sub-millisecond response times

## 🚀 **Running Migrations**

### **Prerequisites**
```bash
# Install required tools
sudo apt-get update
sudo apt-get install -y postgresql-client redis-tools

# Set environment variables for APPLICATION LAYER
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=blockchain_db
export DB_USER=postgres
export DB_PASSWORD=password

# Set Redis configuration
export REDIS_HOST=localhost
export REDIS_PORT=6379
export REDIS_PASSWORD=
export REDIS_DB=0
```

### **Running All Migrations**
```bash
# Run the initialization script for application layer
# Note: Blockchain layer uses RocksDB (managed by Cosmos SDK, no migration scripts needed)
./migrations/init-database.sh
```

### **Application Layer Migrations**
```bash
# Apply application layer PostgreSQL migration
cd migrations/postgresql
psql -h localhost -p 5432 -U postgres -d blockchain_db -f 001_create_blockchain_tables.up.sql

# Apply application layer Redis migration
cd migrations/redis
chmod +x 001_initial_redis_schema.up.sh
./001_initial_redis_schema.up.sh

# Apply Cosmos SDK modules Redis migration
chmod +x 002_cosmos_modules_cache.up.sh
./002_cosmos_modules_cache.up.sh

# Rollback application layer migrations
psql -h localhost -p 5432 -U postgres -d blockchain_db -f 001_create_blockchain_tables.down.sql
chmod +x 001_initial_redis_schema.down.sh
./001_initial_redis_schema.down.sh
chmod +x 002_cosmos_modules_cache.down.sh
./002_cosmos_modules_cache.down.sh
```

### **Blockchain Layer (RocksDB)**
```bash
# Note: Blockchain layer uses RocksDB which is managed by Cosmos SDK
# No migration scripts needed - RocksDB state is managed automatically by Cosmos SDK
# RocksDB data is stored in ./data/cosmos (for Cosmos SDK) and ./data/rocksdb (for business logic)
```

## 📊 **Migration Status Tracking**

### **Version Control**
- Each migration has a unique version number
- Versions are applied sequentially
- Down migrations reverse the exact changes made in up migrations

### **Migration Table (PostgreSQL)**
```sql
CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMP DEFAULT NOW()
);
```

### **Migration State (Redis)**
```yaml
# Stored in Redis with key: "migrations:applied"
migrations:
  - version: "001"
    name: "initial_redis_schema"
    applied_at: "2024-01-01T00:00:00Z"
    database_type: "redis"
```

## 🔧 **Migration Development**

### **Creating New Migrations**

#### **PostgreSQL Migration**
```bash
# Create new migration
cd migrations/postgresql
touch 002_add_nft_storage.up.sql
touch 002_add_nft_storage.down.sql
```

#### **Redis Migration**
```bash
# Create new migration
cd migrations/redis
touch 002_add_nft_cache.up.sh
touch 002_add_nft_cache.down.sh
chmod +x 002_add_nft_cache.up.sh
chmod +x 002_add_nft_cache.down.sh
```

### **Migration Naming Convention**
```
{version}_{description}.{direction}.{extension}
```

**Examples:**
- `001_initial_redis_schema.up.sh`
- `001_initial_redis_schema.down.sh`
- `002_add_nft_storage.up.sql`
- `002_add_nft_storage.down.sql`

### **Migration Content Guidelines**

#### **PostgreSQL SQL**
- Use descriptive table and column names
- Include proper foreign key constraints
- Create appropriate indexes for performance
- Use ENUMs for fixed value sets
- Include proper timestamps and audit fields

#### **Redis Shell Scripts**
- Include connection testing and error handling
- Use environment variables for configuration
- Provide user feedback with progress messages
- Document key patterns and TTL values
- Test Redis connection before executing commands

## 🧪 **Testing Migrations**

### **Test Environment Setup**
```bash
# Create test databases
createdb usc_blockchain_test
redis-server --port 6380 --dbnum 0
```

### **Migration Testing**
```bash
# Test up migration
./test-migration.sh up 001

# Test down migration
./test-migration.sh down 001

# Verify data integrity
./verify-migration.sh 001
```

### **Rollback Testing**
```bash
# Apply migration
./migrate.sh up 001

# Verify data exists
./verify-data.sh

# Rollback migration
./migrate.sh down 001

# Verify data is removed
./verify-rollback.sh
```

## 📈 **Performance Considerations**

### **PostgreSQL**
- **Indexing**: Create indexes on frequently queried columns
- **Partitioning**: Use table partitioning for large tables
- **Connection Pooling**: Optimize connection management
- **Query Optimization**: Use EXPLAIN ANALYZE for slow queries

### **Redis**
- **Memory Management**: Set appropriate maxmemory policies
- **TTL Values**: Balance cache freshness with memory usage
- **Database Separation**: Use separate databases for different data types
- **Persistence**: Configure RDB snapshots for backup

## 🚨 **Troubleshooting**

### **Common Issues**

#### **Migration Fails**
```bash
# Check migration status
./migrate.sh status

# View migration logs
./migrate.sh logs

# Force migration version
./migrate.sh force VERSION
```

#### **Data Corruption**
```bash
# Verify data integrity
./verify-data.sh

# Check database consistency
./check-consistency.sh

# Restore from backup if needed
./restore-backup.sh BACKUP_FILE
```

#### **Performance Issues**
```bash
# Monitor database performance
./monitor-performance.sh

# Analyze slow queries
./analyze-queries.sh

# Optimize indexes
./optimize-indexes.sh
```

### **Recovery Procedures**

#### **Partial Migration Failure**
1. Identify the failed step
2. Manually complete or rollback the step
3. Verify data consistency
4. Continue with remaining migrations

#### **Complete Rollback**
1. Stop all services
2. Apply down migrations in reverse order
3. Verify clean state
4. Restart services

## 📚 **Additional Resources**

### **Documentation**
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Redis Documentation](https://redis.io/documentation)
- [RocksDB Documentation](https://rocksdb.org/docs/) - Note: RocksDB is managed by Cosmos SDK, not via migration scripts

### **Tools**
- [Golang Migrate](https://github.com/golang-migrate/migrate)
- [Redis CLI](https://redis.io/topics/rediscli)

### **Best Practices**
- Always backup before running migrations
- Test migrations in development environment first
- Use transaction wrapping for complex migrations
- Monitor migration performance and impact
- Document all schema changes and their purpose

## 🤝 **Contributing**

### **Migration Review Process**
1. Create migration files following naming conventions
2. Include comprehensive up/down operations
3. Add appropriate tests and verification scripts
4. Document performance implications
5. Submit for review and testing

### **Migration Standards**
- **Reversibility**: All migrations must be reversible
- **Performance**: Consider impact on production systems
- **Documentation**: Clear description of changes and purpose
- **Testing**: Comprehensive testing in staging environment
- **Rollback**: Verified rollback procedures

---

**Last Updated**: 2024-01-01  
**Version**: 1.0.0  
**Maintainer**: USC Blockchain Core Team
