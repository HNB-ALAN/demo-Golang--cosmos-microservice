#!/bin/bash
# USC Platform - Service-04 USC Blockchain Core - BLOCKCHAIN LAYER Redis Migration
# Database: Redis (Blockchain Layer)
# Purpose: Blockchain mempool, consensus cache, and real-time data
# Architecture: High-speed caching for blockchain operations

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

# Function to create Redis keys
create_redis_keys() {
    log_message "INFO" "Creating Redis keys for blockchain layer..."
    
    # Mempool keys
    execute_redis "SET mempool:max_size 10000" "Set mempool max size"
    execute_redis "SET mempool:timeout 300" "Set mempool timeout (5 minutes)"
    execute_redis "SET mempool:cleanup_interval 60" "Set mempool cleanup interval (1 minute)"
    
    # Consensus cache keys
    execute_redis "SET consensus:timeout_commit 1000" "Set consensus timeout commit (1s)"
    execute_redis "SET consensus:timeout_propose 3000" "Set consensus timeout propose (3s)"
    execute_redis "SET consensus:timeout_prevote 1000" "Set consensus timeout prevote (1s)"
    execute_redis "SET consensus:timeout_precommit 1000" "Set consensus timeout precommit (1s)"
    
    # Validator cache keys
    execute_redis "SET validators:max_validators 100" "Set max validators"
    execute_redis "SET validators:min_commission_rate 0.05" "Set min commission rate"
    execute_redis "SET validators:max_commission_rate 0.20" "Set max commission rate"
    execute_redis "SET validators:unbonding_time 1814400" "Set unbonding time (21 days in seconds)"
    
    # Staking cache keys
    execute_redis "SET staking:max_delegations 1000" "Set max delegations"
    execute_redis "SET staking:min_delegation 1000000" "Set min delegation (1 USC with 18 decimals)"
    execute_redis "SET staking:max_delegation 1000000000000000000000000000" "Set max delegation (1B USC)"
    
    # Bank module cache keys
    execute_redis "SET bank:max_supply 1000000000000000000000000000" "Set max supply (1B USC)"
    execute_redis "SET bank:denom usc" "Set denomination"
    execute_redis "SET bank:decimals 18" "Set decimals"
    
    # Reward module cache keys
    execute_redis "SET reward:reward_rate 0.1" "Set reward rate"
    execute_redis "SET reward:reward_pool 100000000000000000000000000" "Set reward pool (100M USC)"
    execute_redis "SET reward:distribution_interval 86400" "Set distribution interval (24 hours in seconds)"
    
    # Block module cache keys
    execute_redis "SET block:max_block_size 1048576" "Set max block size (1MB)"
    execute_redis "SET block:max_tx_size 1048576" "Set max transaction size (1MB)"
    execute_redis "SET block:max_gas 10000000" "Set max gas"
    
    # NFT module cache keys
    execute_redis "SET nft:max_supply 1000000" "Set max NFT supply"
    execute_redis "SET nft:max_metadata_size 1024" "Set max metadata size"
    execute_redis "SET nft:max_attributes 100" "Set max attributes"
    
    # Contract module cache keys
    execute_redis "SET contract:max_contract_size 500000" "Set max contract size (500KB)"
    execute_redis "SET contract:max_contracts 10000" "Set max contracts"
    execute_redis "SET contract:max_instantiate_per_address 100" "Set max instantiate per address"
    
    # Network module cache keys
    execute_redis "SET network:max_channels 100" "Set max channels"
    execute_redis "SET network:max_connections 1000" "Set max connections"
    execute_redis "SET network:timeout 30" "Set network timeout (30 seconds)"
    
    # Bridge module cache keys
    execute_redis "SET bridge:max_bridges 10" "Set max bridges"
    execute_redis "SET bridge:max_bridge_tokens 1000" "Set max bridge tokens"
    execute_redis "SET bridge:bridge_fee_rate 0.001" "Set bridge fee rate"
    
    # Monitoring module cache keys
    execute_redis "SET monitoring:metrics_enabled 1" "Enable metrics"
    execute_redis "SET monitoring:metrics_interval 60" "Set metrics interval (1 minute)"
    execute_redis "SET monitoring:max_metrics_history 1000" "Set max metrics history"
    
    # Performance module cache keys
    execute_redis "SET performance:optimization_enabled 1" "Enable optimization"
    execute_redis "SET performance:cache_size 134217728" "Set cache size (128MB)"
    execute_redis "SET performance:max_cache_entries 1000000" "Set max cache entries"
    
    # Store module cache keys
    execute_redis "SET store:max_stores 1000" "Set max stores"
    execute_redis "SET store:max_store_size 10485760" "Set max store size (10MB)"
    execute_redis "SET store:max_items_per_store 100000" "Set max items per store"
    
    # Streaming module cache keys
    execute_redis "SET streaming:max_streams 100" "Set max streams"
    execute_redis "SET streaming:max_stream_size 104857600" "Set max stream size (100MB)"
    execute_redis "SET streaming:stream_timeout 3600" "Set stream timeout (1 hour)"
    
    # Certificate module cache keys
    execute_redis "SET certificate:max_certificates 10000" "Set max certificates"
    execute_redis "SET certificate:max_certificate_size 1024" "Set max certificate size"
    execute_redis "SET certificate:certificate_expiry 31536000" "Set certificate expiry (365 days in seconds)"
    
    # Token module cache keys
    execute_redis "SET token:max_tokens 1000" "Set max tokens"
    execute_redis "SET token:max_token_supply 1000000000000000000000000000" "Set max token supply (1B tokens)"
    execute_redis "SET token:max_metadata_size 1024" "Set max metadata size"
}

# Function to create Redis data structures
create_redis_structures() {
    log_message "INFO" "Creating Redis data structures for blockchain layer..."
    
    # Create mempool data structures
    execute_redis "HSET mempool:config max_size 10000 timeout 300 cleanup_interval 60" "Create mempool config hash"
    
    # Create consensus data structures
    execute_redis "HSET consensus:config timeout_commit 1000 timeout_propose 3000 timeout_prevote 1000 timeout_precommit 1000" "Create consensus config hash"
    
    # Create validator data structures
    execute_redis "HSET validators:config max_validators 100 min_commission_rate 0.05 max_commission_rate 0.20 unbonding_time 1814400" "Create validators config hash"
    
    # Create staking data structures
    execute_redis "HSET staking:config max_delegations 1000 min_delegation 1000000 max_delegation 1000000000000000000000000000" "Create staking config hash"
    
    # Create bank data structures
    execute_redis "HSET bank:config max_supply 1000000000000000000000000000 denom usc decimals 18" "Create bank config hash"
    
    # Create reward data structures
    execute_redis "HSET reward:config reward_rate 0.1 reward_pool 100000000000000000000000000 distribution_interval 86400" "Create reward config hash"
    
    # Create block data structures
    execute_redis "HSET block:config max_block_size 1048576 max_tx_size 1048576 max_gas 10000000" "Create block config hash"
    
    # Create NFT data structures
    execute_redis "HSET nft:config max_supply 1000000 max_metadata_size 1024 max_attributes 100" "Create NFT config hash"
    
    # Create contract data structures
    execute_redis "HSET contract:config max_contract_size 500000 max_contracts 10000 max_instantiate_per_address 100" "Create contract config hash"
    
    # Create network data structures
    execute_redis "HSET network:config max_channels 100 max_connections 1000 timeout 30" "Create network config hash"
    
    # Create bridge data structures
    execute_redis "HSET bridge:config max_bridges 10 max_bridge_tokens 1000 bridge_fee_rate 0.001" "Create bridge config hash"
    
    # Create monitoring data structures
    execute_redis "HSET monitoring:config metrics_enabled 1 metrics_interval 60 max_metrics_history 1000" "Create monitoring config hash"
    
    # Create performance data structures
    execute_redis "HSET performance:config optimization_enabled 1 cache_size 134217728 max_cache_entries 1000000" "Create performance config hash"
    
    # Create store data structures
    execute_redis "HSET store:config max_stores 1000 max_store_size 10485760 max_items_per_store 100000" "Create store config hash"
    
    # Create streaming data structures
    execute_redis "HSET streaming:config max_streams 100 max_stream_size 104857600 stream_timeout 3600" "Create streaming config hash"
    
    # Create certificate data structures
    execute_redis "HSET certificate:config max_certificates 10000 max_certificate_size 1024 certificate_expiry 31536000" "Create certificate config hash"
    
    # Create token data structures
    execute_redis "HSET token:config max_tokens 1000 max_token_supply 1000000000000000000000000000 max_metadata_size 1024" "Create token config hash"
}

# Function to create Redis indexes
create_redis_indexes() {
    log_message "INFO" "Creating Redis indexes for blockchain layer..."
    
    # Create mempool indexes
    execute_redis "ZADD mempool:by_priority 0 default" "Create mempool priority index"
    execute_redis "ZADD mempool:by_timestamp 0 $(date +%s)" "Create mempool timestamp index"
    
    # Create consensus indexes
    execute_redis "ZADD consensus:by_height 0 0" "Create consensus height index"
    execute_redis "ZADD consensus:by_timestamp 0 $(date +%s)" "Create consensus timestamp index"
    
    # Create validator indexes
    execute_redis "ZADD validators:by_voting_power 0 default" "Create validator voting power index"
    execute_redis "ZADD validators:by_commission 0 0.05" "Create validator commission index"
    
    # Create staking indexes
    execute_redis "ZADD staking:by_delegation 0 default" "Create staking delegation index"
    execute_redis "ZADD staking:by_height 0 0" "Create staking height index"
    
    # Create bank indexes
    execute_redis "ZADD bank:by_balance 0 default" "Create bank balance index"
    execute_redis "ZADD bank:by_address 0 default" "Create bank address index"
    
    # Create reward indexes
    execute_redis "ZADD reward:by_distribution 0 default" "Create reward distribution index"
    execute_redis "ZADD reward:by_interval 0 86400" "Create reward interval index"
    
    # Create block indexes
    execute_redis "ZADD block:by_height 0 0" "Create block height index"
    execute_redis "ZADD block:by_timestamp 0 $(date +%s)" "Create block timestamp index"
    
    # Create NFT indexes
    execute_redis "ZADD nft:by_supply 0 0" "Create NFT supply index"
    execute_redis "ZADD nft:by_metadata 0 default" "Create NFT metadata index"
    
    # Create contract indexes
    execute_redis "ZADD contract:by_size 0 0" "Create contract size index"
    execute_redis "ZADD contract:by_address 0 default" "Create contract address index"
    
    # Create network indexes
    execute_redis "ZADD network:by_channels 0 0" "Create network channels index"
    execute_redis "ZADD network:by_connections 0 0" "Create network connections index"
    
    # Create bridge indexes
    execute_redis "ZADD bridge:by_bridges 0 0" "Create bridge bridges index"
    execute_redis "ZADD bridge:by_tokens 0 0" "Create bridge tokens index"
    
    # Create monitoring indexes
    execute_redis "ZADD monitoring:by_metrics 0 0" "Create monitoring metrics index"
    execute_redis "ZADD monitoring:by_timestamp 0 $(date +%s)" "Create monitoring timestamp index"
    
    # Create performance indexes
    execute_redis "ZADD performance:by_cache 0 0" "Create performance cache index"
    execute_redis "ZADD performance:by_optimization 0 0" "Create performance optimization index"
    
    # Create store indexes
    execute_redis "ZADD store:by_stores 0 0" "Create store stores index"
    execute_redis "ZADD store:by_size 0 0" "Create store size index"
    
    # Create streaming indexes
    execute_redis "ZADD streaming:by_streams 0 0" "Create streaming streams index"
    execute_redis "ZADD streaming:by_size 0 0" "Create streaming size index"
    
    # Create certificate indexes
    execute_redis "ZADD certificate:by_certificates 0 0" "Create certificate certificates index"
    execute_redis "ZADD certificate:by_expiry 0 31536000" "Create certificate expiry index"
    
    # Create token indexes
    execute_redis "ZADD token:by_tokens 0 0" "Create token tokens index"
    execute_redis "ZADD token:by_supply 0 0" "Create token supply index"
}

# Main execution
main() {
    log_message "INFO" "Starting Redis blockchain layer migration..."
    
    # Create Redis keys
    create_redis_keys
    
    # Create Redis data structures
    create_redis_structures
    
    # Create Redis indexes
    create_redis_indexes
    
    log_message "SUCCESS" "Redis blockchain layer migration completed successfully!"
}

# Run main function
main "$@"
