#!/bin/bash
# USC Platform - Service-04 USC Blockchain Core - BLOCKCHAIN LAYER Redis Cosmos Modules Cache Migration Rollback
# Database: Redis (Blockchain Layer)
# Purpose: Remove Cosmos SDK modules cache and analytics
# Architecture: Cleanup high-speed caching for Cosmos SDK module operations

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REDIS_HOST=${REDIS_HOST:-localhost}
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
        *)
            echo -e "${NC}[LOG]${NC} $message"
            ;;
    esac
}

# Function to execute Redis command
execute_redis() {
    local command=$1
    local description=$2
    
    log_message "INFO" "Executing: $description"
    
    if [ -n "$REDIS_PASSWORD" ]; then
        redis-cli -h $REDIS_HOST -p $REDIS_PORT -a $REDIS_PASSWORD -n $REDIS_DB $command
    else
        redis-cli -h $REDIS_HOST -p $REDIS_PORT -n $REDIS_DB $command
    fi
    
    if [ $? -eq 0 ]; then
        log_message "SUCCESS" "$description completed successfully"
    else
        log_message "ERROR" "$description failed"
        exit 1
    fi
}

# Function to remove Cosmos SDK modules cache keys
remove_cosmos_modules_cache() {
    log_message "INFO" "Removing Cosmos SDK modules cache keys..."
    
    # USC Module cache
    execute_redis "DEL usc:module:enabled" "Remove USC module enabled"
    execute_redis "DEL usc:module:version" "Remove USC module version"
    execute_redis "DEL usc:module:config" "Remove USC module config"
    execute_redis "DEL usc:module:keeper_address" "Remove USC module keeper address"
    
    # Reward Module cache
    execute_redis "DEL reward:module:enabled" "Remove Reward module enabled"
    execute_redis "DEL reward:module:version" "Remove Reward module version"
    execute_redis "DEL reward:module:config" "Remove Reward module config"
    execute_redis "DEL reward:module:keeper_address" "Remove Reward module keeper address"
    
    # NFT Module cache
    execute_redis "DEL nft:module:enabled" "Remove NFT module enabled"
    execute_redis "DEL nft:module:version" "Remove NFT module version"
    execute_redis "DEL nft:module:config" "Remove NFT module config"
    execute_redis "DEL nft:module:keeper_address" "Remove NFT module keeper address"
    
    # Contract Module cache
    execute_redis "DEL contract:module:enabled" "Remove Contract module enabled"
    execute_redis "DEL contract:module:version" "Remove Contract module version"
    execute_redis "DEL contract:module:config" "Remove Contract module config"
    execute_redis "DEL contract:module:keeper_address" "Remove Contract module keeper address"
    
    # Validator Module cache
    execute_redis "DEL validator:module:enabled" "Remove Validator module enabled"
    execute_redis "DEL validator:module:version" "Remove Validator module version"
    execute_redis "DEL validator:module:config" "Remove Validator module config"
    execute_redis "DEL validator:module:keeper_address" "Remove Validator module keeper address"
    
    # Network Module cache
    execute_redis "DEL network:module:enabled" "Remove Network module enabled"
    execute_redis "DEL network:module:version" "Remove Network module version"
    execute_redis "DEL network:module:config" "Remove Network module config"
    execute_redis "DEL network:module:keeper_address" "Remove Network module keeper address"
    
    # Bridge Module cache
    execute_redis "DEL bridge:module:enabled" "Remove Bridge module enabled"
    execute_redis "DEL bridge:module:version" "Remove Bridge module version"
    execute_redis "DEL bridge:module:config" "Remove Bridge module config"
    execute_redis "DEL bridge:module:keeper_address" "Remove Bridge module keeper address"
    
    # Streaming Module cache
    execute_redis "DEL streaming:module:enabled" "Remove Streaming module enabled"
    execute_redis "DEL streaming:module:version" "Remove Streaming module version"
    execute_redis "DEL streaming:module:config" "Remove Streaming module config"
    execute_redis "DEL streaming:module:keeper_address" "Remove Streaming module keeper address"
    
    # Certificate Module cache
    execute_redis "DEL certificate:module:enabled" "Remove Certificate module enabled"
    execute_redis "DEL certificate:module:version" "Remove Certificate module version"
    execute_redis "DEL certificate:module:config" "Remove Certificate module config"
    execute_redis "DEL certificate:module:keeper_address" "Remove Certificate module keeper address"
    
    # Store Module cache
    execute_redis "DEL store:module:enabled" "Remove Store module enabled"
    execute_redis "DEL store:module:version" "Remove Store module version"
    execute_redis "DEL store:module:config" "Remove Store module config"
    execute_redis "DEL store:module:keeper_address" "Remove Store module keeper address"
    
    # Token Module cache
    execute_redis "DEL token:module:enabled" "Remove Token module enabled"
    execute_redis "DEL token:module:version" "Remove Token module version"
    execute_redis "DEL token:module:config" "Remove Token module config"
    execute_redis "DEL token:module:keeper_address" "Remove Token module keeper address"
    
    # Block Module cache
    execute_redis "DEL block:module:enabled" "Remove Block module enabled"
    execute_redis "DEL block:module:version" "Remove Block module version"
    execute_redis "DEL block:module:config" "Remove Block module config"
    execute_redis "DEL block:module:keeper_address" "Remove Block module keeper address"
}

# Function to remove module analytics cache
remove_module_analytics_cache() {
    log_message "INFO" "Removing module analytics cache keys..."
    
    # Module performance metrics
    execute_redis "DEL analytics:module:performance:enabled" "Remove module performance analytics"
    execute_redis "DEL analytics:module:performance:interval" "Remove performance analytics interval"
    execute_redis "DEL analytics:module:performance:retention" "Remove performance analytics retention"
    
    # Module error tracking
    execute_redis "DEL analytics:module:errors:enabled" "Remove module error tracking"
    execute_redis "DEL analytics:module:errors:max_errors" "Remove max error tracking count"
    execute_redis "DEL analytics:module:errors:retention" "Remove error tracking retention"
    
    # Module state tracking
    execute_redis "DEL analytics:module:state:enabled" "Remove module state tracking"
    execute_redis "DEL analytics:module:state:interval" "Remove state tracking interval"
    execute_redis "DEL analytics:module:state:retention" "Remove state tracking retention"
}

# Function to remove module events cache
remove_module_events_cache() {
    log_message "INFO" "Removing module events cache keys..."
    
    # Module event streams
    execute_redis "DEL events:module:streams:enabled" "Remove module event streams"
    execute_redis "DEL events:module:streams:max_streams" "Remove max event streams"
    execute_redis "DEL events:module:streams:timeout" "Remove event stream timeout"
    
    # Module event processing
    execute_redis "DEL events:module:processing:enabled" "Remove module event processing"
    execute_redis "DEL events:module:processing:batch_size" "Remove event processing batch size"
    execute_redis "DEL events:module:processing:workers" "Remove event processing workers"
}

# Function to remove module keeper cache
remove_module_keeper_cache() {
    log_message "INFO" "Removing module keeper cache keys..."
    
    # Keeper state management
    execute_redis "DEL keeper:state:enabled" "Remove keeper state management"
    execute_redis "DEL keeper:state:sync_interval" "Remove keeper state sync interval"
    execute_redis "DEL keeper:state:retention" "Remove keeper state retention"
    
    # Keeper health monitoring
    execute_redis "DEL keeper:health:enabled" "Remove keeper health monitoring"
    execute_redis "DEL keeper:health:check_interval" "Remove keeper health check interval"
    execute_redis "DEL keeper:health:timeout" "Remove keeper health check timeout"
    
    # Keeper performance tracking
    execute_redis "DEL keeper:performance:enabled" "Remove keeper performance tracking"
    execute_redis "DEL keeper:performance:metrics_interval" "Remove keeper performance metrics interval"
    execute_redis "DEL keeper:performance:retention" "Remove keeper performance retention"
}

# Main execution
log_message "INFO" "Starting Cosmos SDK modules cache migration rollback..."

# Remove module cache keys
remove_cosmos_modules_cache

# Remove analytics cache keys
remove_module_analytics_cache

# Remove events cache keys
remove_module_events_cache

# Remove keeper cache keys
remove_module_keeper_cache

log_message "SUCCESS" "Cosmos SDK modules cache migration rollback completed successfully"
