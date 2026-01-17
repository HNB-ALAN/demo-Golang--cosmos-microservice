#!/bin/bash
# USC Platform - Service-04 USC Blockchain Core - BLOCKCHAIN LAYER Redis Migration Rollback
# Database: Redis (Blockchain Layer)
# Purpose: Rollback blockchain mempool, consensus cache, and real-time data

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
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

# Function to execute Redis command
execute_redis() {
    local command="$1"
    local description="$2"
    
    log_message "INFO" "Executing: $description"
    
    if [ -n "$REDIS_PASSWORD" ]; then
        redis-cli -h $REDIS_HOST -p $REDIS_PORT -a $REDIS_PASSWORD -n $REDIS_DB $command
    else
        redis-cli -h $REDIS_HOST -p $REDIS_PORT -n $REDIS_DB $command
    fi
    
    if [ $? -eq 0 ]; then
        log_message "SUCCESS" "Successfully executed: $description"
    else
        log_message "ERROR" "Failed to execute: $description"
        exit 1
    fi
}

# Function to remove Redis keys
remove_redis_keys() {
    log_message "INFO" "Removing Redis keys for blockchain layer..."
    
    # Remove mempool keys
    execute_redis "DEL mempool:max_size" "Remove mempool max size"
    execute_redis "DEL mempool:timeout" "Remove mempool timeout"
    execute_redis "DEL mempool:cleanup_interval" "Remove mempool cleanup interval"
    
    # Remove consensus cache keys
    execute_redis "DEL consensus:timeout_commit" "Remove consensus timeout commit"
    execute_redis "DEL consensus:timeout_propose" "Remove consensus timeout propose"
    execute_redis "DEL consensus:timeout_prevote" "Remove consensus timeout prevote"
    execute_redis "DEL consensus:timeout_precommit" "Remove consensus timeout precommit"
    
    # Remove validator cache keys
    execute_redis "DEL validators:max_validators" "Remove max validators"
    execute_redis "DEL validators:min_commission_rate" "Remove min commission rate"
    execute_redis "DEL validators:max_commission_rate" "Remove max commission rate"
    execute_redis "DEL validators:unbonding_time" "Remove unbonding time"
    
    # Remove staking cache keys
    execute_redis "DEL staking:max_delegations" "Remove max delegations"
    execute_redis "DEL staking:min_delegation" "Remove min delegation"
    execute_redis "DEL staking:max_delegation" "Remove max delegation"
    
    # Remove bank module cache keys
    execute_redis "DEL bank:max_supply" "Remove max supply"
    execute_redis "DEL bank:denom" "Remove denomination"
    execute_redis "DEL bank:decimals" "Remove decimals"
    
    # Remove reward module cache keys
    execute_redis "DEL reward:reward_rate" "Remove reward rate"
    execute_redis "DEL reward:reward_pool" "Remove reward pool"
    execute_redis "DEL reward:distribution_interval" "Remove distribution interval"
    
    # Remove block module cache keys
    execute_redis "DEL block:max_block_size" "Remove max block size"
    execute_redis "DEL block:max_tx_size" "Remove max transaction size"
    execute_redis "DEL block:max_gas" "Remove max gas"
    
    # Remove NFT module cache keys
    execute_redis "DEL nft:max_supply" "Remove max NFT supply"
    execute_redis "DEL nft:max_metadata_size" "Remove max metadata size"
    execute_redis "DEL nft:max_attributes" "Remove max attributes"
    
    # Remove contract module cache keys
    execute_redis "DEL contract:max_contract_size" "Remove max contract size"
    execute_redis "DEL contract:max_contracts" "Remove max contracts"
    execute_redis "DEL contract:max_instantiate_per_address" "Remove max instantiate per address"
    
    # Remove network module cache keys
    execute_redis "DEL network:max_channels" "Remove max channels"
    execute_redis "DEL network:max_connections" "Remove max connections"
    execute_redis "DEL network:timeout" "Remove network timeout"
    
    # Remove bridge module cache keys
    execute_redis "DEL bridge:max_bridges" "Remove max bridges"
    execute_redis "DEL bridge:max_bridge_tokens" "Remove max bridge tokens"
    execute_redis "DEL bridge:bridge_fee_rate" "Remove bridge fee rate"
    
    # Remove monitoring module cache keys
    execute_redis "DEL monitoring:metrics_enabled" "Remove metrics enabled"
    execute_redis "DEL monitoring:metrics_interval" "Remove metrics interval"
    execute_redis "DEL monitoring:max_metrics_history" "Remove max metrics history"
    
    # Remove performance module cache keys
    execute_redis "DEL performance:optimization_enabled" "Remove optimization enabled"
    execute_redis "DEL performance:cache_size" "Remove cache size"
    execute_redis "DEL performance:max_cache_entries" "Remove max cache entries"
    
    # Remove store module cache keys
    execute_redis "DEL store:max_stores" "Remove max stores"
    execute_redis "DEL store:max_store_size" "Remove max store size"
    execute_redis "DEL store:max_items_per_store" "Remove max items per store"
    
    # Remove streaming module cache keys
    execute_redis "DEL streaming:max_streams" "Remove max streams"
    execute_redis "DEL streaming:max_stream_size" "Remove max stream size"
    execute_redis "DEL streaming:stream_timeout" "Remove stream timeout"
    
    # Remove certificate module cache keys
    execute_redis "DEL certificate:max_certificates" "Remove max certificates"
    execute_redis "DEL certificate:max_certificate_size" "Remove max certificate size"
    execute_redis "DEL certificate:certificate_expiry" "Remove certificate expiry"
    
    # Remove token module cache keys
    execute_redis "DEL token:max_tokens" "Remove max tokens"
    execute_redis "DEL token:max_token_supply" "Remove max token supply"
    execute_redis "DEL token:max_metadata_size" "Remove max metadata size"
}

# Function to remove Redis data structures
remove_redis_structures() {
    log_message "INFO" "Removing Redis data structures for blockchain layer..."
    
    # Remove mempool data structures
    execute_redis "DEL mempool:config" "Remove mempool config hash"
    
    # Remove consensus data structures
    execute_redis "DEL consensus:config" "Remove consensus config hash"
    
    # Remove validator data structures
    execute_redis "DEL validators:config" "Remove validators config hash"
    
    # Remove staking data structures
    execute_redis "DEL staking:config" "Remove staking config hash"
    
    # Remove bank data structures
    execute_redis "DEL bank:config" "Remove bank config hash"
    
    # Remove reward data structures
    execute_redis "DEL reward:config" "Remove reward config hash"
    
    # Remove block data structures
    execute_redis "DEL block:config" "Remove block config hash"
    
    # Remove NFT data structures
    execute_redis "DEL nft:config" "Remove NFT config hash"
    
    # Remove contract data structures
    execute_redis "DEL contract:config" "Remove contract config hash"
    
    # Remove network data structures
    execute_redis "DEL network:config" "Remove network config hash"
    
    # Remove bridge data structures
    execute_redis "DEL bridge:config" "Remove bridge config hash"
    
    # Remove monitoring data structures
    execute_redis "DEL monitoring:config" "Remove monitoring config hash"
    
    # Remove performance data structures
    execute_redis "DEL performance:config" "Remove performance config hash"
    
    # Remove store data structures
    execute_redis "DEL store:config" "Remove store config hash"
    
    # Remove streaming data structures
    execute_redis "DEL streaming:config" "Remove streaming config hash"
    
    # Remove certificate data structures
    execute_redis "DEL certificate:config" "Remove certificate config hash"
    
    # Remove token data structures
    execute_redis "DEL token:config" "Remove token config hash"
}

# Function to remove Redis indexes
remove_redis_indexes() {
    log_message "INFO" "Removing Redis indexes for blockchain layer..."
    
    # Remove mempool indexes
    execute_redis "DEL mempool:by_priority" "Remove mempool priority index"
    execute_redis "DEL mempool:by_timestamp" "Remove mempool timestamp index"
    
    # Remove consensus indexes
    execute_redis "DEL consensus:by_height" "Remove consensus height index"
    execute_redis "DEL consensus:by_timestamp" "Remove consensus timestamp index"
    
    # Remove validator indexes
    execute_redis "DEL validators:by_voting_power" "Remove validator voting power index"
    execute_redis "DEL validators:by_commission" "Remove validator commission index"
    
    # Remove staking indexes
    execute_redis "DEL staking:by_delegation" "Remove staking delegation index"
    execute_redis "DEL staking:by_height" "Remove staking height index"
    
    # Remove bank indexes
    execute_redis "DEL bank:by_balance" "Remove bank balance index"
    execute_redis "DEL bank:by_address" "Remove bank address index"
    
    # Remove reward indexes
    execute_redis "DEL reward:by_distribution" "Remove reward distribution index"
    execute_redis "DEL reward:by_interval" "Remove reward interval index"
    
    # Remove block indexes
    execute_redis "DEL block:by_height" "Remove block height index"
    execute_redis "DEL block:by_timestamp" "Remove block timestamp index"
    
    # Remove NFT indexes
    execute_redis "DEL nft:by_supply" "Remove NFT supply index"
    execute_redis "DEL nft:by_metadata" "Remove NFT metadata index"
    
    # Remove contract indexes
    execute_redis "DEL contract:by_size" "Remove contract size index"
    execute_redis "DEL contract:by_address" "Remove contract address index"
    
    # Remove network indexes
    execute_redis "DEL network:by_channels" "Remove network channels index"
    execute_redis "DEL network:by_connections" "Remove network connections index"
    
    # Remove bridge indexes
    execute_redis "DEL bridge:by_bridges" "Remove bridge bridges index"
    execute_redis "DEL bridge:by_tokens" "Remove bridge tokens index"
    
    # Remove monitoring indexes
    execute_redis "DEL monitoring:by_metrics" "Remove monitoring metrics index"
    execute_redis "DEL monitoring:by_timestamp" "Remove monitoring timestamp index"
    
    # Remove performance indexes
    execute_redis "DEL performance:by_cache" "Remove performance cache index"
    execute_redis "DEL performance:by_optimization" "Remove performance optimization index"
    
    # Remove store indexes
    execute_redis "DEL store:by_stores" "Remove store stores index"
    execute_redis "DEL store:by_size" "Remove store size index"
    
    # Remove streaming indexes
    execute_redis "DEL streaming:by_streams" "Remove streaming streams index"
    execute_redis "DEL streaming:by_size" "Remove streaming size index"
    
    # Remove certificate indexes
    execute_redis "DEL certificate:by_certificates" "Remove certificate certificates index"
    execute_redis "DEL certificate:by_expiry" "Remove certificate expiry index"
    
    # Remove token indexes
    execute_redis "DEL token:by_tokens" "Remove token tokens index"
    execute_redis "DEL token:by_supply" "Remove token supply index"
}

# Function to flush all blockchain data
flush_blockchain_data() {
    log_message "INFO" "Flushing all blockchain data..."
    
    # Flush all keys with blockchain prefix
    execute_redis "FLUSHDB" "Flush current database"
    
    log_message "SUCCESS" "All blockchain data flushed successfully!"
}

# Main execution
main() {
    log_message "INFO" "Starting Redis blockchain layer migration rollback..."
    
    # Remove Redis keys
    remove_redis_keys
    
    # Remove Redis data structures
    remove_redis_structures
    
    # Remove Redis indexes
    remove_redis_indexes
    
    # Flush all blockchain data
    flush_blockchain_data
    
    log_message "SUCCESS" "Redis blockchain layer migration rollback completed successfully!"
}

# Run main function
main "$@"
