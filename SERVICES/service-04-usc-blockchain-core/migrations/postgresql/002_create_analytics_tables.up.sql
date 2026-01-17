-- USC Platform - Service-04 USC Blockchain Core - Analytics Tables Migration
-- Database: PostgreSQL (Primary)
-- Purpose: Create analytics tables for transaction and block analytics
-- Architecture: Analytics tables for reporting and monitoring

-- ============================================================================
-- ANALYTICS TABLE DEFINITIONS
-- ============================================================================

-- Transaction Analytics Table
-- Purpose: Store transaction analytics data for reporting and monitoring
CREATE TABLE IF NOT EXISTS usc_transaction_analytics (
    transaction_hash VARCHAR(66) PRIMARY KEY,
    from_address VARCHAR(42) NOT NULL,
    to_address VARCHAR(42) NOT NULL,
    amount DECIMAL(28,18) NOT NULL,
    gas_price DECIMAL(28,18) NOT NULL,
    gas_limit BIGINT NOT NULL,
    gas_used BIGINT,
    timestamp TIMESTAMP NOT NULL,
    status VARCHAR(50) NOT NULL,
    block_number BIGINT,
    transaction_index INTEGER,
    block_hash VARCHAR(66),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for Transaction Analytics
CREATE INDEX IF NOT EXISTS idx_tx_analytics_block ON usc_transaction_analytics(block_number, transaction_index);
CREATE INDEX IF NOT EXISTS idx_tx_analytics_from ON usc_transaction_analytics(from_address);
CREATE INDEX IF NOT EXISTS idx_tx_analytics_to ON usc_transaction_analytics(to_address);
CREATE INDEX IF NOT EXISTS idx_tx_analytics_status ON usc_transaction_analytics(status);
CREATE INDEX IF NOT EXISTS idx_tx_analytics_timestamp ON usc_transaction_analytics(timestamp);
CREATE INDEX IF NOT EXISTS idx_tx_analytics_block_hash ON usc_transaction_analytics(block_hash);

-- Block Analytics Table
-- Purpose: Store block analytics data for reporting and monitoring
CREATE TABLE IF NOT EXISTS usc_block_analytics (
    block_number BIGINT PRIMARY KEY,
    block_hash VARCHAR(66) UNIQUE NOT NULL,
    validator_address VARCHAR(42) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    transaction_count INTEGER NOT NULL DEFAULT 0,
    total_usc_transferred DECIMAL(28,18) DEFAULT 0,
    gas_used BIGINT DEFAULT 0,
    gas_limit BIGINT DEFAULT 0,
    block_size_bytes INTEGER DEFAULT 0,
    processing_time_ms INTEGER,
    is_finalized BOOLEAN DEFAULT FALSE,
    finalized_at TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for Block Analytics
CREATE INDEX IF NOT EXISTS idx_block_analytics_validator ON usc_block_analytics(validator_address);
CREATE INDEX IF NOT EXISTS idx_block_analytics_timestamp ON usc_block_analytics(timestamp);
CREATE INDEX IF NOT EXISTS idx_block_analytics_finalized ON usc_block_analytics(is_finalized, finalized_at);
CREATE INDEX IF NOT EXISTS idx_block_analytics_metrics ON usc_block_analytics(transaction_count, gas_used);

-- USC Coin Analytics Table (if not exists)
-- Purpose: Store USC coin transaction analytics
CREATE TABLE IF NOT EXISTS usc_coin_analytics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    wallet_address VARCHAR(42) NOT NULL,
    amount DECIMAL(28,18) NOT NULL,
    transaction_type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    from_address VARCHAR(42),
    to_address VARCHAR(42),
    transaction_hash VARCHAR(66),
    block_number BIGINT,
    block_hash VARCHAR(66),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for USC Coin Analytics
CREATE INDEX IF NOT EXISTS idx_usc_coin_analytics_wallet ON usc_coin_analytics(wallet_address);
CREATE INDEX IF NOT EXISTS idx_usc_coin_analytics_type ON usc_coin_analytics(transaction_type);
CREATE INDEX IF NOT EXISTS idx_usc_coin_analytics_status ON usc_coin_analytics(status);
CREATE INDEX IF NOT EXISTS idx_usc_coin_analytics_timestamp ON usc_coin_analytics(timestamp);
CREATE INDEX IF NOT EXISTS idx_usc_coin_analytics_tx_hash ON usc_coin_analytics(transaction_hash);

-- Smart Contract Analytics Table
-- Purpose: Store smart contract deployment analytics
CREATE TABLE IF NOT EXISTS usc_smart_contract_analytics (
    contract_address VARCHAR(42) PRIMARY KEY,
    contract_name VARCHAR(255) NOT NULL,
    bytecode TEXT,
    abi TEXT,
    from_address VARCHAR(42) NOT NULL,
    gas_price VARCHAR(50),
    gas_limit BIGINT,
    transaction_hash VARCHAR(66),
    status VARCHAR(50) NOT NULL,
    deployed_at TIMESTAMP NOT NULL,
    last_updated TIMESTAMP NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for Smart Contract Analytics
CREATE INDEX IF NOT EXISTS idx_smart_contract_analytics_from ON usc_smart_contract_analytics(from_address);
CREATE INDEX IF NOT EXISTS idx_smart_contract_analytics_status ON usc_smart_contract_analytics(status);
CREATE INDEX IF NOT EXISTS idx_smart_contract_analytics_tx_hash ON usc_smart_contract_analytics(transaction_hash);
CREATE INDEX IF NOT EXISTS idx_smart_contract_analytics_deployed_at ON usc_smart_contract_analytics(deployed_at);

-- Contract Execution Analytics Table
-- Purpose: Store smart contract execution analytics
CREATE TABLE IF NOT EXISTS usc_contract_execution_analytics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    contract_address VARCHAR(42) NOT NULL,
    function_name VARCHAR(255) NOT NULL,
    from_address VARCHAR(42) NOT NULL,
    gas_price VARCHAR(50),
    gas_limit BIGINT,
    transaction_hash VARCHAR(66),
    status VARCHAR(50) NOT NULL,
    executed_at TIMESTAMP NOT NULL,
    return_value TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for Contract Execution Analytics
CREATE INDEX IF NOT EXISTS idx_contract_execution_analytics_contract ON usc_contract_execution_analytics(contract_address);
CREATE INDEX IF NOT EXISTS idx_contract_execution_analytics_from ON usc_contract_execution_analytics(from_address);
CREATE INDEX IF NOT EXISTS idx_contract_execution_analytics_status ON usc_contract_execution_analytics(status);
CREATE INDEX IF NOT EXISTS idx_contract_execution_analytics_tx_hash ON usc_contract_execution_analytics(transaction_hash);
CREATE INDEX IF NOT EXISTS idx_contract_execution_analytics_executed_at ON usc_contract_execution_analytics(executed_at);

-- Validator Analytics Table
-- Purpose: Store validator registration and status analytics
CREATE TABLE IF NOT EXISTS usc_validator_analytics (
    validator_address VARCHAR(42) PRIMARY KEY,
    validator_name VARCHAR(255) NOT NULL,
    validator_public_key TEXT,
    commission_rate VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    registered_at TIMESTAMP NOT NULL,
    last_updated TIMESTAMP NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for Validator Analytics
CREATE INDEX IF NOT EXISTS idx_validator_analytics_status ON usc_validator_analytics(status);
CREATE INDEX IF NOT EXISTS idx_validator_analytics_registered_at ON usc_validator_analytics(registered_at);

-- Staking Analytics Table
-- Purpose: Store staking transaction analytics
CREATE TABLE IF NOT EXISTS usc_staking_analytics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    validator_address VARCHAR(42) NOT NULL,
    staker_address VARCHAR(42) NOT NULL,
    amount DECIMAL(28,18) NOT NULL,
    transaction_type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    transaction_hash VARCHAR(66),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for Staking Analytics
CREATE INDEX IF NOT EXISTS idx_staking_analytics_validator ON usc_staking_analytics(validator_address);
CREATE INDEX IF NOT EXISTS idx_staking_analytics_staker ON usc_staking_analytics(staker_address);
CREATE INDEX IF NOT EXISTS idx_staking_analytics_type ON usc_staking_analytics(transaction_type);
CREATE INDEX IF NOT EXISTS idx_staking_analytics_status ON usc_staking_analytics(status);
CREATE INDEX IF NOT EXISTS idx_staking_analytics_timestamp ON usc_staking_analytics(timestamp);
CREATE INDEX IF NOT EXISTS idx_staking_analytics_tx_hash ON usc_staking_analytics(transaction_hash);


