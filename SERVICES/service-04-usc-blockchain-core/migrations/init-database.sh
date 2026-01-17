#!/bin/bash
# Database Initialization Script
# Description: Initialize database for SERVICE-04-USC-BLOCKCHAIN-CORE
# Created: 2024-01-01

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
DB_HOST=${DB_HOST:-postgres}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-blockchain_db}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-password}

# Redis Configuration
REDIS_HOST=${REDIS_HOST:-redis}
REDIS_PORT=${REDIS_PORT:-6379}
REDIS_PASSWORD=${REDIS_PASSWORD:-}
REDIS_DB=${REDIS_DB:-0}


# Function to log messages
log_message() {
    local level=$1
    local message=$2
    
    case $level in
        "INFO")
            echo -e "${BLUE}[INFO]${NC} $message"
            ;;
        "WARN")
            echo -e "${YELLOW}[WARN]${NC} $message"
            ;;
        "ERROR")
            echo -e "${RED}[ERROR]${NC} $message"
            ;;
        "SUCCESS")
            echo -e "${GREEN}[SUCCESS]${NC} $message"
            ;;
    esac
}

# Function to create database (legacy - now handled by check_database_exists)
create_database() {
    log_message "INFO" "Creating database $DB_NAME"
    
    # Use the new check_database_exists function
    check_database_exists
}

# Function to run application layer migrations
run_application_migrations() {
    log_message "INFO" "Running application layer migrations..."
    
    # Run PostgreSQL migrations for application layer
    if [ -f "migrations/postgresql/001_create_blockchain_tables.up.sql" ]; then
        log_message "INFO" "Running application layer PostgreSQL migration..."
        PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f migrations/postgresql/001_create_blockchain_tables.up.sql
        
        if [ $? -eq 0 ]; then
            log_message "SUCCESS" "Application layer PostgreSQL migration completed successfully"
        else
            log_message "ERROR" "Application layer PostgreSQL migration failed"
            exit 1
        fi
    else
        log_message "WARN" "Application layer PostgreSQL migration file not found"
    fi
    
    
    # Run Analytics tables migration for application layer
    if [ -f "migrations/postgresql/002_create_analytics_tables.up.sql" ]; then
        log_message "INFO" "Running Analytics tables migration..."
        PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f migrations/postgresql/002_create_analytics_tables.up.sql
        
        if [ $? -eq 0 ]; then
            log_message "SUCCESS" "Analytics tables migration completed successfully"
        else
            log_message "ERROR" "Analytics tables migration failed"
            exit 1
        fi
    else
        log_message "WARN" "Analytics tables migration file not found"
    fi
}

# Function to run blockchain layer migrations
# NOTE: Blockchain layer PostgreSQL migrations have been removed as they are not used
# RocksDB state is managed by Cosmos SDK, not via migration scripts
run_blockchain_migrations() {
    log_message "INFO" "Blockchain layer PostgreSQL migrations removed (not used)"
    log_message "INFO" "RocksDB state is managed by Cosmos SDK, skipping RocksDB migrations"
}

# Function to check if database exists
check_database_exists() {
    log_message "INFO" "Checking if database exists"
    
    # Check if database exists
    local db_exists=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -t -c "SELECT EXISTS (SELECT FROM pg_database WHERE datname = '$DB_NAME');" 2>/dev/null | tr -d ' ' || echo "f")
    
    if [ "$db_exists" = "t" ]; then
        log_message "SUCCESS" "Database $DB_NAME exists"
        return 0
    else
        log_message "INFO" "Database $DB_NAME does not exist, creating it"
        PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -c "CREATE DATABASE $DB_NAME;" 2>/dev/null || {
            log_message "ERROR" "Failed to create database $DB_NAME"
            return 1
        }
        log_message "SUCCESS" "Database $DB_NAME created successfully"
        return 0
    fi
}

# Function to run PostgreSQL migrations
run_postgresql_migrations() {
    log_message "INFO" "Checking if PostgreSQL migrations are needed"
    
    # First check if database exists
    if ! check_database_exists; then
        log_message "ERROR" "Database check failed"
        return 1
    fi
    
    # Check if blockchain_transactions table exists
    local transactions_exists=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'blockchain_transactions');" 2>/dev/null | tr -d ' ' || echo "f")
    
    if [ "$transactions_exists" = "t" ]; then
        log_message "INFO" "PostgreSQL tables already exist, skipping migrations"
        return 0
    fi
    
    log_message "INFO" "Running PostgreSQL migrations"
    
    # Run PostgreSQL migrations
    if [ -d "migrations/postgresql" ]; then
        for migration_file in migrations/postgresql/*.up.sql; do
            if [ -f "$migration_file" ]; then
                log_message "INFO" "Running $(basename "$migration_file")"
                PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$migration_file"
            fi
        done
    fi    
    
    log_message "SUCCESS" "PostgreSQL migrations completed"
}

# Function to run Redis migrations for application layer
run_redis_migrations() {
    log_message "INFO" "Running application layer Redis migrations..."
    
    # Run Redis migrations for application layer
    if [ -f "migrations/redis/001_initial_redis_schema.up.sh" ]; then
        log_message "INFO" "Running application layer Redis migration..."
        chmod +x migrations/redis/001_initial_redis_schema.up.sh
        ./migrations/redis/001_initial_redis_schema.up.sh
        
        if [ $? -eq 0 ]; then
            log_message "SUCCESS" "Application layer Redis migration completed successfully"
        else
            log_message "ERROR" "Application layer Redis migration failed"
            exit 1
        fi
    else
        log_message "WARN" "Application layer Redis migration file not found"
    fi
    
    # Run Cosmos SDK modules Redis migration for application layer
    if [ -f "migrations/redis/002_cosmos_modules_cache.up.sh" ]; then
        log_message "INFO" "Running Cosmos SDK modules Redis migration for application layer..."
        chmod +x migrations/redis/002_cosmos_modules_cache.up.sh
        ./migrations/redis/002_cosmos_modules_cache.up.sh
        
        if [ $? -eq 0 ]; then
            log_message "SUCCESS" "Cosmos SDK modules Redis migration for application layer completed successfully"
        else
            log_message "ERROR" "Cosmos SDK modules Redis migration for application layer failed"
            exit 1
        fi
    else
        log_message "WARN" "Cosmos SDK modules Redis migration file not found"
    fi
    
    log_message "SUCCESS" "Application layer Redis migrations completed"
}


# Function to validate migrations
validate_migrations() {
    log_message "INFO" "Validating migrations"
    
    # Check if all required tables exist (must match actual migrations)
    local required_tables=(
        "blocks"
        "transactions"
        "smart_contracts"
        "nft_collections"
        "nfts"
        "custom_tokens"
        "product_certificates"
        "validators"
        "staking"
        "store_bridges"
        "store_networks"
        "bridge_transactions"
        "network_sync_logs"
    )
    
    for table in "${required_tables[@]}"; do
        local exists=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = '$table');" | tr -d ' ')
        
        if [ "$exists" = "t" ]; then
            log_message "SUCCESS" "Table $table exists"
        else
            log_message "ERROR" "Table $table does not exist"
            exit 1
        fi
    done
    
    log_message "SUCCESS" "All required tables exist"
}

# Main function
main() {
    log_message "INFO" "Starting database initialization for SERVICE-04-USC-BLOCKCHAIN-CORE"
    
    # APPLICATION LAYER - Business Logic Database
    log_message "INFO" "=== APPLICATION LAYER INITIALIZATION ==="
    create_database
    run_application_migrations
    run_redis_migrations
    validate_migrations
    
    # BLOCKCHAIN LAYER - Consensus & Blockchain Database
    # NOTE: Blockchain layer PostgreSQL migrations removed (not used)
    # Only RocksDB is used for blockchain state (managed by Cosmos SDK)
    log_message "INFO" "=== BLOCKCHAIN LAYER INITIALIZATION ==="
    log_message "INFO" "Blockchain layer uses RocksDB only (managed by Cosmos SDK)"
    # create_blockchain_database  # Commented out - not used
    # run_blockchain_migrations   # Commented out - migrations removed
    
    log_message "SUCCESS" "Database initialization completed successfully"
    log_message "INFO" "Service-04-USC-BLOCKCHAIN-CORE database is ready for use"
    log_message "INFO" "Application Layer: $DB_NAME (Business Logic) + Redis (Cache & Real-time)"
    log_message "INFO" "Blockchain Layer: RocksDB (State - managed by Cosmos SDK)"
}

# Run main function
main "$@"