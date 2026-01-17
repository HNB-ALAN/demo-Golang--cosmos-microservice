-- USC Platform - Service-04 USC Blockchain Core - Analytics Tables Migration (Rollback)
-- Purpose: Drop analytics tables

-- Drop indexes first
DROP INDEX IF EXISTS idx_tx_analytics_block;
DROP INDEX IF EXISTS idx_tx_analytics_from;
DROP INDEX IF EXISTS idx_tx_analytics_to;
DROP INDEX IF EXISTS idx_tx_analytics_type;
DROP INDEX IF EXISTS idx_tx_analytics_status;
DROP INDEX IF EXISTS idx_tx_analytics_timestamp;

DROP INDEX IF EXISTS idx_block_analytics_validator;
DROP INDEX IF EXISTS idx_block_analytics_timestamp;
DROP INDEX IF EXISTS idx_block_analytics_finalized;
DROP INDEX IF EXISTS idx_block_analytics_metrics;

DROP INDEX IF EXISTS idx_usc_coin_analytics_wallet;
DROP INDEX IF EXISTS idx_usc_coin_analytics_type;
DROP INDEX IF EXISTS idx_usc_coin_analytics_status;
DROP INDEX IF EXISTS idx_usc_coin_analytics_timestamp;
DROP INDEX IF EXISTS idx_usc_coin_analytics_tx_hash;

DROP INDEX IF EXISTS idx_smart_contract_analytics_from;
DROP INDEX IF EXISTS idx_smart_contract_analytics_status;
DROP INDEX IF EXISTS idx_smart_contract_analytics_tx_hash;
DROP INDEX IF EXISTS idx_smart_contract_analytics_deployed_at;

DROP INDEX IF EXISTS idx_contract_execution_analytics_contract;
DROP INDEX IF EXISTS idx_contract_execution_analytics_from;
DROP INDEX IF EXISTS idx_contract_execution_analytics_status;
DROP INDEX IF EXISTS idx_contract_execution_analytics_tx_hash;
DROP INDEX IF EXISTS idx_contract_execution_analytics_executed_at;

DROP INDEX IF EXISTS idx_validator_analytics_status;
DROP INDEX IF EXISTS idx_validator_analytics_registered_at;

DROP INDEX IF EXISTS idx_staking_analytics_validator;
DROP INDEX IF EXISTS idx_staking_analytics_staker;
DROP INDEX IF EXISTS idx_staking_analytics_type;
DROP INDEX IF EXISTS idx_staking_analytics_status;
DROP INDEX IF EXISTS idx_staking_analytics_timestamp;
DROP INDEX IF EXISTS idx_staking_analytics_tx_hash;

-- Drop tables
DROP TABLE IF EXISTS usc_transaction_analytics;
DROP TABLE IF EXISTS usc_block_analytics;
DROP TABLE IF EXISTS usc_coin_analytics;
DROP TABLE IF EXISTS usc_smart_contract_analytics;
DROP TABLE IF EXISTS usc_contract_execution_analytics;
DROP TABLE IF EXISTS usc_validator_analytics;
DROP TABLE IF EXISTS usc_staking_analytics;


