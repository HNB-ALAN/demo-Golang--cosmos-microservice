#!/bin/bash

# Script to check data in all databases
set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=== 🔍 KIỂM TRA DỮ LIỆU TRONG CÁC DATABASE ===${NC}\n"

# 1. PostgreSQL
echo -e "${YELLOW}1. PostgreSQL Database:${NC}"
echo "   Checking tables and data..."

# Check if PostgreSQL is accessible
if docker-compose exec -T postgres psql -U postgres -d blockchain_db -c "SELECT 1;" > /dev/null 2>&1; then
    echo -e "   ${GREEN}✓ PostgreSQL is accessible${NC}"
    
    # List all tables
    echo "   Tables:"
    docker-compose exec -T postgres psql -U postgres -d blockchain_db -c "\dt" 2>/dev/null | grep -v "List of relations" | grep -v "Schema" | grep -v "Name" | grep -v "Type" | grep -v "Owner" | grep -v "^-" | grep -v "^$" || echo "   No tables found"
    
    # Check blocks table
    if docker-compose exec -T postgres psql -U postgres -d blockchain_db -c "SELECT COUNT(*) FROM blocks;" > /dev/null 2>&1; then
        BLOCK_COUNT=$(docker-compose exec -T postgres psql -U postgres -d blockchain_db -t -c "SELECT COUNT(*) FROM blocks;" 2>/dev/null | xargs)
        echo -e "   ${GREEN}✓ blocks table: ${BLOCK_COUNT} records${NC}"
    else
        echo -e "   ${YELLOW}⚠ blocks table does not exist${NC}"
    fi
    
    # Check transactions table
    if docker-compose exec -T postgres psql -U postgres -d blockchain_db -c "SELECT COUNT(*) FROM transactions;" > /dev/null 2>&1; then
        TX_COUNT=$(docker-compose exec -T postgres psql -U postgres -d blockchain_db -t -c "SELECT COUNT(*) FROM transactions;" 2>/dev/null | xargs)
        echo -e "   ${GREEN}✓ transactions table: ${TX_COUNT} records${NC}"
    else
        echo -e "   ${YELLOW}⚠ transactions table does not exist${NC}"
    fi
else
    echo -e "   ${RED}✗ PostgreSQL is not accessible${NC}"
fi

echo ""

# 2. CometBFT (Blockchain State)
echo -e "${YELLOW}2. CometBFT Blockchain State:${NC}"
COMETBFT_HEIGHT=$(curl -s http://localhost:26657/status 2>/dev/null | jq -r '.result.sync_info.latest_block_height // "N/A"')
if [ "$COMETBFT_HEIGHT" != "N/A" ]; then
    echo -e "   ${GREEN}✓ Latest block height: ${COMETBFT_HEIGHT}${NC}"
    
    # Get block 1 info
    BLOCK_1_HASH=$(curl -s "http://localhost:26657/block?height=1" 2>/dev/null | jq -r '.result.block_id.hash // "N/A"')
    if [ "$BLOCK_1_HASH" != "N/A" ] && [ "$BLOCK_1_HASH" != "null" ]; then
        echo -e "   ${GREEN}✓ Block 1 hash: ${BLOCK_1_HASH:0:20}...${NC}"
    fi
else
    echo -e "   ${RED}✗ CometBFT is not accessible${NC}"
fi

echo ""

# 3. RocksDB (Cosmos SDK State)
echo -e "${YELLOW}3. RocksDB (Cosmos SDK State):${NC}"
# Check multiple possible paths
ROCKSDB_PATH=""
for path in "/app/block-chain-cosmos/data" "/data/blockchain" "/data/rocksdb" "/app/data"; do
    if docker-compose exec -T service-04-usc-blockchain-core test -d "$path" 2>/dev/null; then
        ROCKSDB_PATH="$path"
        break
    fi
done

if [ -n "$ROCKSDB_PATH" ]; then
    ROCKSDB_SIZE=$(docker-compose exec -T service-04-usc-blockchain-core du -sh "$ROCKSDB_PATH" 2>/dev/null | awk '{print $1}')
    echo -e "   ${GREEN}✓ RocksDB directory exists: ${ROCKSDB_PATH} (${ROCKSDB_SIZE})${NC}"
    
    # Check if there are files
    FILE_COUNT=$(docker-compose exec -T service-04-usc-blockchain-core find "$ROCKSDB_PATH" -type f 2>/dev/null | wc -l)
    echo -e "   ${GREEN}✓ Files in RocksDB: ${FILE_COUNT}${NC}"
    
    # Check application.db subdirectory
    if docker-compose exec -T service-04-usc-blockchain-core test -d "$ROCKSDB_PATH/application.db" 2>/dev/null; then
        APP_DB_SIZE=$(docker-compose exec -T service-04-usc-blockchain-core du -sh "$ROCKSDB_PATH/application.db" 2>/dev/null | awk '{print $1}')
        echo -e "   ${GREEN}✓ application.db size: ${APP_DB_SIZE}${NC}"
    fi
else
    echo -e "   ${YELLOW}⚠ RocksDB directory not found in common paths${NC}"
fi

echo ""

# 4. Redis (Cache)
echo -e "${YELLOW}4. Redis Cache:${NC}"
if docker-compose exec -T redis redis-cli ping > /dev/null 2>&1; then
    echo -e "   ${GREEN}✓ Redis is accessible${NC}"
    
    # Count keys
    KEY_COUNT=$(docker-compose exec -T redis redis-cli DBSIZE 2>/dev/null | xargs)
    echo -e "   ${GREEN}✓ Total keys: ${KEY_COUNT}${NC}"
    
    # Check block-related keys
    BLOCK_KEYS=$(docker-compose exec -T redis redis-cli KEYS "*block*" 2>/dev/null | wc -l)
    echo -e "   ${GREEN}✓ Block-related keys: ${BLOCK_KEYS}${NC}"
else
    echo -e "   ${RED}✗ Redis is not accessible${NC}"
fi

echo ""

# 5. Service-04 gRPC
echo -e "${YELLOW}5. Service-04 gRPC API:${NC}"
# Check if grpcurl is available
if ! command -v grpcurl &> /dev/null; then
    echo -e "   ${YELLOW}⚠ grpcurl not installed, checking via docker exec${NC}"
    if docker-compose exec -T service-04-usc-blockchain-core grpc_health_probe -addr=localhost:8004 > /dev/null 2>&1; then
        echo -e "   ${GREEN}✓ Service-04 is healthy (via health probe)${NC}"
    else
        echo -e "   ${YELLOW}⚠ Cannot verify gRPC health (grpcurl not available)${NC}"
    fi
elif grpcurl -plaintext localhost:8004 grpc.health.v1.Health/Check > /dev/null 2>&1; then
    echo -e "   ${GREEN}✓ Service-04 is healthy${NC}"
    
    # Try to get latest block
    LATEST_BLOCK=$(grpcurl -plaintext localhost:8004 service04.BlockOperationsService/GetLatestBlock 2>/dev/null | jq -r '.block_number // "N/A"')
    if [ "$LATEST_BLOCK" != "N/A" ] && [ "$LATEST_BLOCK" != "null" ]; then
        echo -e "   ${GREEN}✓ Latest block via gRPC: ${LATEST_BLOCK}${NC}"
    else
        echo -e "   ${YELLOW}⚠ Latest block not available via gRPC${NC}"
    fi
else
    echo -e "   ${YELLOW}⚠ Service-04 gRPC check failed (may need grpcurl or service restart)${NC}"
fi

echo ""
echo -e "${BLUE}=== 📊 SUMMARY ===${NC}"
echo "Check completed. Review results above."

