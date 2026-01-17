-- USC Platform - Service-04 USC Blockchain Core - PostgreSQL Migration
-- Database: PostgreSQL (Primary)
-- Purpose: Blockchain infrastructure, USC transactions, smart contracts, NFTs, validators
-- Architecture: ACID compliance for financial data, complex queries for blockchain analytics

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================================================
-- TABLE DEFINITIONS
-- ============================================================================

-- Blocks table - Blockchain block data
CREATE TABLE IF NOT EXISTS blocks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    block_hash VARCHAR(66) UNIQUE NOT NULL,
    block_number BIGINT UNIQUE NOT NULL,
    previous_block_hash VARCHAR(66),
    merkle_root VARCHAR(66),
    timestamp BIGINT NOT NULL,
    nonce BIGINT NOT NULL,
    validator_address VARCHAR(42),
    gas_used BIGINT DEFAULT 0,
    gas_limit BIGINT DEFAULT 0,
    extra_data TEXT,
    block_data TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Transactions table - USC transaction records
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_hash VARCHAR(66) UNIQUE NOT NULL,
    block_number BIGINT,
    block_hash VARCHAR(66),
    transaction_index INTEGER,
    from_address VARCHAR(42) NOT NULL,
    to_address VARCHAR(42),
    amount DECIMAL(28,18) NOT NULL,
    gas_price DECIMAL(28,18) NOT NULL,
    gas_limit BIGINT NOT NULL,
    gas_used BIGINT DEFAULT 0,
    data TEXT,
    nonce BIGINT NOT NULL,
    signature TEXT NOT NULL,
    status INTEGER DEFAULT 0, -- 0=Pending, 1=Confirmed, 2=Failed
    user_id UUID,
    device_id VARCHAR(255),
    memo TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    confirmed_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Smart contracts table - Deployed smart contracts
CREATE TABLE IF NOT EXISTS smart_contracts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    contract_address VARCHAR(42) UNIQUE NOT NULL,
    contract_name VARCHAR(255) NOT NULL,
    contract_version VARCHAR(50),
    bytecode TEXT NOT NULL,
    abi TEXT NOT NULL,
    source_code TEXT,
    compiler_version VARCHAR(50),
    constructor_args TEXT,
    deployer_address VARCHAR(42) NOT NULL,
    deployment_transaction_hash VARCHAR(66) NOT NULL,
    gas_used BIGINT,
    is_verified BOOLEAN DEFAULT FALSE,
    description TEXT,
    user_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- NFT collections table - NFT collection metadata
CREATE TABLE IF NOT EXISTS nft_collections (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    collection_id VARCHAR(255) UNIQUE NOT NULL,
    contract_address VARCHAR(42) NOT NULL,
    collection_name VARCHAR(255) NOT NULL,
    collection_description TEXT,
    collection_image_url TEXT,
    collection_banner_url TEXT,
    creator_address VARCHAR(42) NOT NULL,
    royalty_percentage DECIMAL(5,2) DEFAULT 0,
    category VARCHAR(100),
    tags TEXT[],
    user_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- NFTs table - Individual NFT tokens
CREATE TABLE IF NOT EXISTS nfts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    token_id VARCHAR(255) NOT NULL,
    contract_address VARCHAR(42) NOT NULL,
    collection_id VARCHAR(255),
    owner_address VARCHAR(42) NOT NULL,
    token_uri TEXT,
    metadata TEXT, -- JSON metadata
    name VARCHAR(255),
    description TEXT,
    image_url TEXT,
    animation_url TEXT,
    attributes TEXT, -- JSON attributes
    creator_address VARCHAR(42),
    transfer_count BIGINT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_transferred_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(contract_address, token_id)
);

-- Custom tokens table - Store coins and custom tokens
CREATE TABLE IF NOT EXISTS custom_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    contract_address VARCHAR(42) UNIQUE NOT NULL,
    token_name VARCHAR(255) NOT NULL,
    token_symbol VARCHAR(50) NOT NULL,
    total_supply DECIMAL(28,18) NOT NULL,
    circulating_supply DECIMAL(28,18) DEFAULT 0,
    decimals INTEGER DEFAULT 18,
    owner_address VARCHAR(42) NOT NULL,
    is_mintable BOOLEAN DEFAULT TRUE,
    is_burnable BOOLEAN DEFAULT TRUE,
    token_description TEXT,
    token_image_url TEXT,
    store_id VARCHAR(255),
    deployment_transaction_hash VARCHAR(66) NOT NULL,
    gas_used BIGINT,
    user_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Product certificates table - Product authenticity certificates
CREATE TABLE IF NOT EXISTS product_certificates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    certificate_id VARCHAR(255) UNIQUE NOT NULL,
    product_id VARCHAR(255) NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    product_description TEXT,
    manufacturer_address VARCHAR(42) NOT NULL,
    serial_number VARCHAR(255),
    batch_number VARCHAR(255),
    manufacturing_date DATE,
    expiration_date DATE,
    expires_at BIGINT, -- Unix timestamp for code compatibility
    product_metadata TEXT, -- JSON metadata (legacy)
    metadata TEXT, -- JSON metadata (used by code)
    certificate_type VARCHAR(100),
    certificate_hash VARCHAR(66),
    current_owner_address VARCHAR(42) NOT NULL,
    status VARCHAR(50) DEFAULT 'active',
    deployment_transaction_hash VARCHAR(66) NOT NULL,
    gas_used BIGINT,
    user_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Product certificate ownership history table - Track ownership transfers
CREATE TABLE IF NOT EXISTS product_certificate_ownership_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    certificate_id VARCHAR(255) NOT NULL,
    from_address VARCHAR(42) NOT NULL,
    to_address VARCHAR(42) NOT NULL,
    transaction_hash VARCHAR(66) NOT NULL,
    transferred_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT product_certificate_ownership_history_unique UNIQUE (certificate_id, transaction_hash)
);

-- Validators table - PoS validators
CREATE TABLE IF NOT EXISTS validators (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    validator_id VARCHAR(255) UNIQUE NOT NULL,
    validator_address VARCHAR(42) UNIQUE NOT NULL,
    validator_name VARCHAR(255) NOT NULL,
    description TEXT,
    website_url TEXT,
    logo_url TEXT,
    commission_rate DECIMAL(5,2) NOT NULL,
    min_delegation DECIMAL(28,18) DEFAULT 0,
    validator_public_key TEXT,
    status VARCHAR(20) DEFAULT 'active', -- active, inactive, jailed
    stake_amount DECIMAL(28,18) DEFAULT 0,
    delegated_amount DECIMAL(28,18) DEFAULT 0,
    uptime_percentage DECIMAL(5,2) DEFAULT 0,
    blocks_proposed BIGINT DEFAULT 0,
    blocks_missed BIGINT DEFAULT 0,
    rewards_earned DECIMAL(28,18) DEFAULT 0,
    jail_reason TEXT,
    jailed_until TIMESTAMP WITH TIME ZONE,
    user_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_active_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Staking table - USC staking records
CREATE TABLE IF NOT EXISTS staking (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    delegator_address VARCHAR(42) NOT NULL,
    validator_address VARCHAR(42) NOT NULL,
    stake_amount DECIMAL(28,18) NOT NULL,
    stake_type VARCHAR(20) NOT NULL, -- delegate, redelegate, undelegate
    transaction_hash VARCHAR(66) NOT NULL,
    gas_used BIGINT,
    user_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    unlock_time TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT staking_unique_delegation UNIQUE (delegator_address, validator_address, transaction_hash)
);

-- Store bridges table - Cross-chain bridge contracts
CREATE TABLE IF NOT EXISTS store_bridges (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bridge_address VARCHAR(42) UNIQUE NOT NULL,
    bridge_name VARCHAR(255) NOT NULL,
    bridge_description TEXT,
    target_network VARCHAR(100) NOT NULL,
    target_chain_id VARCHAR(50) NOT NULL,
    bridge_config TEXT, -- JSON configuration
    bridge_status VARCHAR(20) DEFAULT 'active', -- active, inactive, maintenance
    deployment_transaction_hash VARCHAR(66) NOT NULL,
    gas_used BIGINT,
    user_id UUID,
    store_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Store networks table - External network integration
CREATE TABLE IF NOT EXISTS store_networks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    network_id VARCHAR(255) UNIQUE NOT NULL,
    network_name VARCHAR(255) NOT NULL,
    network_type VARCHAR(50) NOT NULL, -- ethereum, binance, polygon, etc.
    chain_id VARCHAR(50) NOT NULL,
    rpc_url TEXT NOT NULL,
    explorer_url TEXT,
    native_currency VARCHAR(20) DEFAULT 'ETH',
    network_config TEXT, -- JSON configuration
    network_status VARCHAR(20) DEFAULT 'active', -- active, inactive, maintenance
    sync_status VARCHAR(20) DEFAULT 'synced', -- synced, syncing, behind
    current_block BIGINT DEFAULT 0,
    latest_block BIGINT DEFAULT 0,
    blocks_behind BIGINT DEFAULT 0,
    network_health VARCHAR(20) DEFAULT 'healthy', -- healthy, degraded, down
    last_sync TIMESTAMP WITH TIME ZONE,
    user_id UUID,
    store_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Bridge transactions table - Cross-chain transactions
CREATE TABLE IF NOT EXISTS bridge_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_hash VARCHAR(66) UNIQUE NOT NULL,
    bridge_address VARCHAR(42) NOT NULL,
    from_address VARCHAR(42) NOT NULL,
    to_address VARCHAR(42),
    source_network VARCHAR(100) NOT NULL,
    target_network VARCHAR(100) NOT NULL,
    source_token_address VARCHAR(42),
    target_token_address VARCHAR(42),
    source_amount DECIMAL(28,18) NOT NULL,
    target_amount DECIMAL(28,18),
    bridge_fee DECIMAL(28,18) DEFAULT 0,
    transaction_type VARCHAR(50) NOT NULL, -- token_to_usc, usc_to_token
    status VARCHAR(20) DEFAULT 'pending', -- pending, completed, failed, cancelled
    gas_used BIGINT,
    gas_price DECIMAL(28,18),
    user_id UUID,
    device_id VARCHAR(255),
    memo TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Network sync logs table - Network synchronization tracking
CREATE TABLE IF NOT EXISTS network_sync_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    network_id VARCHAR(255) NOT NULL,
    sync_id VARCHAR(255) UNIQUE NOT NULL,
    sync_type VARCHAR(50) NOT NULL, -- full, incremental, selective
    from_block BIGINT,
    to_block BIGINT,
    blocks_synced BIGINT DEFAULT 0,
    transactions_synced BIGINT DEFAULT 0,
    sync_status VARCHAR(20) DEFAULT 'in_progress', -- in_progress, completed, failed
    sync_started_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    sync_completed_at TIMESTAMP WITH TIME ZONE,
    error_message TEXT,
    user_id UUID,
    store_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ============================================================================
-- INDEXES
-- ============================================================================

-- Blocks indexes
CREATE INDEX IF NOT EXISTS idx_blocks_hash ON blocks(block_hash);
CREATE INDEX IF NOT EXISTS idx_blocks_number ON blocks(block_number);
CREATE INDEX IF NOT EXISTS idx_blocks_timestamp ON blocks(timestamp);
CREATE INDEX IF NOT EXISTS idx_blocks_validator ON blocks(validator_address);

-- Transactions indexes
CREATE INDEX IF NOT EXISTS idx_transactions_hash ON transactions(transaction_hash);
CREATE INDEX IF NOT EXISTS idx_transactions_from ON transactions(from_address);
CREATE INDEX IF NOT EXISTS idx_transactions_to ON transactions(to_address);
CREATE INDEX IF NOT EXISTS idx_transactions_block ON transactions(block_number);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
CREATE INDEX IF NOT EXISTS idx_transactions_user ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_created ON transactions(created_at);

-- Smart contracts indexes
CREATE INDEX IF NOT EXISTS idx_contracts_address ON smart_contracts(contract_address);
CREATE INDEX IF NOT EXISTS idx_contracts_deployer ON smart_contracts(deployer_address);
CREATE INDEX IF NOT EXISTS idx_contracts_verified ON smart_contracts(is_verified);
CREATE INDEX IF NOT EXISTS idx_contracts_user ON smart_contracts(user_id);

-- NFT collections indexes
CREATE INDEX IF NOT EXISTS idx_nft_collections_id ON nft_collections(collection_id);
CREATE INDEX IF NOT EXISTS idx_nft_collections_contract ON nft_collections(contract_address);
CREATE INDEX IF NOT EXISTS idx_nft_collections_creator ON nft_collections(creator_address);
CREATE INDEX IF NOT EXISTS idx_nft_collections_user ON nft_collections(user_id);

-- NFTs indexes
CREATE INDEX IF NOT EXISTS idx_nfts_token ON nfts(token_id);
CREATE INDEX IF NOT EXISTS idx_nfts_contract ON nfts(contract_address);
CREATE INDEX IF NOT EXISTS idx_nfts_owner ON nfts(owner_address);
CREATE INDEX IF NOT EXISTS idx_nfts_collection ON nfts(collection_id);
CREATE INDEX IF NOT EXISTS idx_nfts_creator ON nfts(creator_address);

-- Custom tokens indexes
CREATE INDEX IF NOT EXISTS idx_tokens_address ON custom_tokens(contract_address);
CREATE INDEX IF NOT EXISTS idx_tokens_owner ON custom_tokens(owner_address);
CREATE INDEX IF NOT EXISTS idx_tokens_store ON custom_tokens(store_id);
CREATE INDEX IF NOT EXISTS idx_tokens_user ON custom_tokens(user_id);

-- Product certificates indexes
CREATE INDEX IF NOT EXISTS idx_certificates_id ON product_certificates(certificate_id);
CREATE INDEX IF NOT EXISTS idx_certificates_product ON product_certificates(product_id);
CREATE INDEX IF NOT EXISTS idx_certificates_manufacturer ON product_certificates(manufacturer_address);
CREATE INDEX IF NOT EXISTS idx_certificates_current_owner ON product_certificates(current_owner_address);
CREATE INDEX IF NOT EXISTS idx_certificates_status ON product_certificates(status);
CREATE INDEX IF NOT EXISTS idx_certificates_user ON product_certificates(user_id);

-- Product certificate ownership history indexes
CREATE INDEX IF NOT EXISTS idx_certificate_ownership_history_certificate ON product_certificate_ownership_history(certificate_id);
CREATE INDEX IF NOT EXISTS idx_certificate_ownership_history_from ON product_certificate_ownership_history(from_address);
CREATE INDEX IF NOT EXISTS idx_certificate_ownership_history_to ON product_certificate_ownership_history(to_address);
CREATE INDEX IF NOT EXISTS idx_certificate_ownership_history_tx_hash ON product_certificate_ownership_history(transaction_hash);
CREATE INDEX IF NOT EXISTS idx_certificate_ownership_history_transferred_at ON product_certificate_ownership_history(transferred_at);

-- Validators indexes
CREATE INDEX IF NOT EXISTS idx_validators_id ON validators(validator_id);
CREATE INDEX IF NOT EXISTS idx_validators_address ON validators(validator_address);
CREATE INDEX IF NOT EXISTS idx_validators_status ON validators(status);
CREATE INDEX IF NOT EXISTS idx_validators_user ON validators(user_id);

-- Staking indexes
CREATE INDEX IF NOT EXISTS idx_staking_delegator ON staking(delegator_address);
CREATE INDEX IF NOT EXISTS idx_staking_validator ON staking(validator_address);
CREATE INDEX IF NOT EXISTS idx_staking_type ON staking(stake_type);
CREATE INDEX IF NOT EXISTS idx_staking_user ON staking(user_id);
CREATE INDEX IF NOT EXISTS idx_staking_unique_delegation ON staking(delegator_address, validator_address, transaction_hash);

-- Store bridges indexes
CREATE INDEX IF NOT EXISTS idx_store_bridges_address ON store_bridges(bridge_address);
CREATE INDEX IF NOT EXISTS idx_store_bridges_target_network ON store_bridges(target_network);
CREATE INDEX IF NOT EXISTS idx_store_bridges_status ON store_bridges(bridge_status);
CREATE INDEX IF NOT EXISTS idx_store_bridges_user_id ON store_bridges(user_id);
CREATE INDEX IF NOT EXISTS idx_store_bridges_store_id ON store_bridges(store_id);

-- Store networks indexes
CREATE INDEX IF NOT EXISTS idx_store_networks_network_id ON store_networks(network_id);
CREATE INDEX IF NOT EXISTS idx_store_networks_network_type ON store_networks(network_type);
CREATE INDEX IF NOT EXISTS idx_store_networks_status ON store_networks(network_status);
CREATE INDEX IF NOT EXISTS idx_store_networks_sync_status ON store_networks(sync_status);
CREATE INDEX IF NOT EXISTS idx_store_networks_user_id ON store_networks(user_id);
CREATE INDEX IF NOT EXISTS idx_store_networks_store_id ON store_networks(store_id);

-- Bridge transactions indexes
CREATE INDEX IF NOT EXISTS idx_bridge_transactions_hash ON bridge_transactions(transaction_hash);
CREATE INDEX IF NOT EXISTS idx_bridge_transactions_bridge ON bridge_transactions(bridge_address);
CREATE INDEX IF NOT EXISTS idx_bridge_transactions_from ON bridge_transactions(from_address);
CREATE INDEX IF NOT EXISTS idx_bridge_transactions_networks ON bridge_transactions(source_network, target_network);
CREATE INDEX IF NOT EXISTS idx_bridge_transactions_status ON bridge_transactions(status);
CREATE INDEX IF NOT EXISTS idx_bridge_transactions_type ON bridge_transactions(transaction_type);
CREATE INDEX IF NOT EXISTS idx_bridge_transactions_user_id ON bridge_transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_bridge_transactions_created ON bridge_transactions(created_at);

-- Network sync logs indexes
CREATE INDEX IF NOT EXISTS idx_network_sync_logs_network ON network_sync_logs(network_id);
CREATE INDEX IF NOT EXISTS idx_network_sync_logs_sync_id ON network_sync_logs(sync_id);
CREATE INDEX IF NOT EXISTS idx_network_sync_logs_status ON network_sync_logs(sync_status);
CREATE INDEX IF NOT EXISTS idx_network_sync_logs_user_id ON network_sync_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_network_sync_logs_store_id ON network_sync_logs(store_id);
CREATE INDEX IF NOT EXISTS idx_network_sync_logs_created ON network_sync_logs(created_at);

-- ============================================================================
-- TRIGGERS
-- ============================================================================

-- Create function for updated_at timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Blocks triggers
DROP TRIGGER IF EXISTS update_blocks_updated_at ON blocks;
DROP TRIGGER IF EXISTS update_blocks_updated_at ON blocks;
DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;
DROP TRIGGER IF EXISTS update_contracts_updated_at ON smart_contracts;
DROP TRIGGER IF EXISTS update_collections_updated_at ON nft_collections;
DROP TRIGGER IF EXISTS update_nfts_updated_at ON nfts;
DROP TRIGGER IF EXISTS update_tokens_updated_at ON custom_tokens;
DROP TRIGGER IF EXISTS update_certificates_updated_at ON product_certificates;
DROP TRIGGER IF EXISTS update_validators_updated_at ON validators;
DROP TRIGGER IF EXISTS update_staking_updated_at ON staking;
DROP TRIGGER IF EXISTS update_store_bridges_updated_at ON store_bridges;
DROP TRIGGER IF EXISTS update_store_networks_updated_at ON store_networks;
DROP TRIGGER IF EXISTS update_bridge_transactions_updated_at ON bridge_transactions;
CREATE TRIGGER update_blocks_updated_at BEFORE UPDATE ON blocks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Transactions triggers
DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;
DROP TRIGGER IF EXISTS update_blocks_updated_at ON blocks;
DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;
DROP TRIGGER IF EXISTS update_contracts_updated_at ON smart_contracts;
DROP TRIGGER IF EXISTS update_collections_updated_at ON nft_collections;
DROP TRIGGER IF EXISTS update_nfts_updated_at ON nfts;
DROP TRIGGER IF EXISTS update_tokens_updated_at ON custom_tokens;
DROP TRIGGER IF EXISTS update_certificates_updated_at ON product_certificates;
DROP TRIGGER IF EXISTS update_validators_updated_at ON validators;
DROP TRIGGER IF EXISTS update_staking_updated_at ON staking;
DROP TRIGGER IF EXISTS update_store_bridges_updated_at ON store_bridges;
DROP TRIGGER IF EXISTS update_store_networks_updated_at ON store_networks;
DROP TRIGGER IF EXISTS update_bridge_transactions_updated_at ON bridge_transactions;
CREATE TRIGGER update_transactions_updated_at BEFORE UPDATE ON transactions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Smart contracts triggers
DROP TRIGGER IF EXISTS update_contracts_updated_at ON contracts;
DROP TRIGGER IF EXISTS update_blocks_updated_at ON blocks;
DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;
DROP TRIGGER IF EXISTS update_contracts_updated_at ON smart_contracts;
DROP TRIGGER IF EXISTS update_collections_updated_at ON nft_collections;
DROP TRIGGER IF EXISTS update_nfts_updated_at ON nfts;
DROP TRIGGER IF EXISTS update_tokens_updated_at ON custom_tokens;
DROP TRIGGER IF EXISTS update_certificates_updated_at ON product_certificates;
DROP TRIGGER IF EXISTS update_validators_updated_at ON validators;
DROP TRIGGER IF EXISTS update_staking_updated_at ON staking;
DROP TRIGGER IF EXISTS update_store_bridges_updated_at ON store_bridges;
DROP TRIGGER IF EXISTS update_store_networks_updated_at ON store_networks;
DROP TRIGGER IF EXISTS update_bridge_transactions_updated_at ON bridge_transactions;
CREATE TRIGGER update_contracts_updated_at BEFORE UPDATE ON smart_contracts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- NFT collections triggers
DROP TRIGGER IF EXISTS update_collections_updated_at ON collections;
DROP TRIGGER IF EXISTS update_blocks_updated_at ON blocks;
DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;
DROP TRIGGER IF EXISTS update_contracts_updated_at ON smart_contracts;
DROP TRIGGER IF EXISTS update_collections_updated_at ON nft_collections;
DROP TRIGGER IF EXISTS update_nfts_updated_at ON nfts;
DROP TRIGGER IF EXISTS update_tokens_updated_at ON custom_tokens;
DROP TRIGGER IF EXISTS update_certificates_updated_at ON product_certificates;
DROP TRIGGER IF EXISTS update_validators_updated_at ON validators;
DROP TRIGGER IF EXISTS update_staking_updated_at ON staking;
DROP TRIGGER IF EXISTS update_store_bridges_updated_at ON store_bridges;
DROP TRIGGER IF EXISTS update_store_networks_updated_at ON store_networks;
DROP TRIGGER IF EXISTS update_bridge_transactions_updated_at ON bridge_transactions;
CREATE TRIGGER update_collections_updated_at BEFORE UPDATE ON nft_collections
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- NFTs triggers
DROP TRIGGER IF EXISTS update_nfts_updated_at ON nfts;
DROP TRIGGER IF EXISTS update_blocks_updated_at ON blocks;
DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;
DROP TRIGGER IF EXISTS update_contracts_updated_at ON smart_contracts;
DROP TRIGGER IF EXISTS update_collections_updated_at ON nft_collections;
DROP TRIGGER IF EXISTS update_nfts_updated_at ON nfts;
DROP TRIGGER IF EXISTS update_tokens_updated_at ON custom_tokens;
DROP TRIGGER IF EXISTS update_certificates_updated_at ON product_certificates;
DROP TRIGGER IF EXISTS update_validators_updated_at ON validators;
DROP TRIGGER IF EXISTS update_staking_updated_at ON staking;
DROP TRIGGER IF EXISTS update_store_bridges_updated_at ON store_bridges;
DROP TRIGGER IF EXISTS update_store_networks_updated_at ON store_networks;
DROP TRIGGER IF EXISTS update_bridge_transactions_updated_at ON bridge_transactions;
CREATE TRIGGER update_nfts_updated_at BEFORE UPDATE ON nfts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Custom tokens triggers
DROP TRIGGER IF EXISTS update_tokens_updated_at ON tokens;
DROP TRIGGER IF EXISTS update_blocks_updated_at ON blocks;
DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;
DROP TRIGGER IF EXISTS update_contracts_updated_at ON smart_contracts;
DROP TRIGGER IF EXISTS update_collections_updated_at ON nft_collections;
DROP TRIGGER IF EXISTS update_nfts_updated_at ON nfts;
DROP TRIGGER IF EXISTS update_tokens_updated_at ON custom_tokens;
DROP TRIGGER IF EXISTS update_certificates_updated_at ON product_certificates;
DROP TRIGGER IF EXISTS update_validators_updated_at ON validators;
DROP TRIGGER IF EXISTS update_staking_updated_at ON staking;
DROP TRIGGER IF EXISTS update_store_bridges_updated_at ON store_bridges;
DROP TRIGGER IF EXISTS update_store_networks_updated_at ON store_networks;
DROP TRIGGER IF EXISTS update_bridge_transactions_updated_at ON bridge_transactions;
CREATE TRIGGER update_tokens_updated_at BEFORE UPDATE ON custom_tokens
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Product certificates triggers
DROP TRIGGER IF EXISTS update_certificates_updated_at ON certificates;
DROP TRIGGER IF EXISTS update_blocks_updated_at ON blocks;
DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;
DROP TRIGGER IF EXISTS update_contracts_updated_at ON smart_contracts;
DROP TRIGGER IF EXISTS update_collections_updated_at ON nft_collections;
DROP TRIGGER IF EXISTS update_nfts_updated_at ON nfts;
DROP TRIGGER IF EXISTS update_tokens_updated_at ON custom_tokens;
DROP TRIGGER IF EXISTS update_certificates_updated_at ON product_certificates;
DROP TRIGGER IF EXISTS update_validators_updated_at ON validators;
DROP TRIGGER IF EXISTS update_staking_updated_at ON staking;
DROP TRIGGER IF EXISTS update_store_bridges_updated_at ON store_bridges;
DROP TRIGGER IF EXISTS update_store_networks_updated_at ON store_networks;
DROP TRIGGER IF EXISTS update_bridge_transactions_updated_at ON bridge_transactions;
CREATE TRIGGER update_certificates_updated_at BEFORE UPDATE ON product_certificates
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Validators triggers
DROP TRIGGER IF EXISTS update_validators_updated_at ON validators;
DROP TRIGGER IF EXISTS update_blocks_updated_at ON blocks;
DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;
DROP TRIGGER IF EXISTS update_contracts_updated_at ON smart_contracts;
DROP TRIGGER IF EXISTS update_collections_updated_at ON nft_collections;
DROP TRIGGER IF EXISTS update_nfts_updated_at ON nfts;
DROP TRIGGER IF EXISTS update_tokens_updated_at ON custom_tokens;
DROP TRIGGER IF EXISTS update_certificates_updated_at ON product_certificates;
DROP TRIGGER IF EXISTS update_validators_updated_at ON validators;
DROP TRIGGER IF EXISTS update_staking_updated_at ON staking;
DROP TRIGGER IF EXISTS update_store_bridges_updated_at ON store_bridges;
DROP TRIGGER IF EXISTS update_store_networks_updated_at ON store_networks;
DROP TRIGGER IF EXISTS update_bridge_transactions_updated_at ON bridge_transactions;
CREATE TRIGGER update_validators_updated_at BEFORE UPDATE ON validators
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Staking triggers
DROP TRIGGER IF EXISTS update_staking_updated_at ON staking;
DROP TRIGGER IF EXISTS update_blocks_updated_at ON blocks;
DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;
DROP TRIGGER IF EXISTS update_contracts_updated_at ON smart_contracts;
DROP TRIGGER IF EXISTS update_collections_updated_at ON nft_collections;
DROP TRIGGER IF EXISTS update_nfts_updated_at ON nfts;
DROP TRIGGER IF EXISTS update_tokens_updated_at ON custom_tokens;
DROP TRIGGER IF EXISTS update_certificates_updated_at ON product_certificates;
DROP TRIGGER IF EXISTS update_validators_updated_at ON validators;
DROP TRIGGER IF EXISTS update_staking_updated_at ON staking;
DROP TRIGGER IF EXISTS update_store_bridges_updated_at ON store_bridges;
DROP TRIGGER IF EXISTS update_store_networks_updated_at ON store_networks;
DROP TRIGGER IF EXISTS update_bridge_transactions_updated_at ON bridge_transactions;
CREATE TRIGGER update_staking_updated_at BEFORE UPDATE ON staking
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Store bridges triggers
DROP TRIGGER IF EXISTS update_store_bridges_updated_at ON store_bridges;
DROP TRIGGER IF EXISTS update_blocks_updated_at ON blocks;
DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;
DROP TRIGGER IF EXISTS update_contracts_updated_at ON smart_contracts;
DROP TRIGGER IF EXISTS update_collections_updated_at ON nft_collections;
DROP TRIGGER IF EXISTS update_nfts_updated_at ON nfts;
DROP TRIGGER IF EXISTS update_tokens_updated_at ON custom_tokens;
DROP TRIGGER IF EXISTS update_certificates_updated_at ON product_certificates;
DROP TRIGGER IF EXISTS update_validators_updated_at ON validators;
DROP TRIGGER IF EXISTS update_staking_updated_at ON staking;
DROP TRIGGER IF EXISTS update_store_bridges_updated_at ON store_bridges;
DROP TRIGGER IF EXISTS update_store_networks_updated_at ON store_networks;
DROP TRIGGER IF EXISTS update_bridge_transactions_updated_at ON bridge_transactions;
CREATE TRIGGER update_store_bridges_updated_at BEFORE UPDATE ON store_bridges
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Store networks triggers
DROP TRIGGER IF EXISTS update_store_networks_updated_at ON store_networks;
DROP TRIGGER IF EXISTS update_blocks_updated_at ON blocks;
DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;
DROP TRIGGER IF EXISTS update_contracts_updated_at ON smart_contracts;
DROP TRIGGER IF EXISTS update_collections_updated_at ON nft_collections;
DROP TRIGGER IF EXISTS update_nfts_updated_at ON nfts;
DROP TRIGGER IF EXISTS update_tokens_updated_at ON custom_tokens;
DROP TRIGGER IF EXISTS update_certificates_updated_at ON product_certificates;
DROP TRIGGER IF EXISTS update_validators_updated_at ON validators;
DROP TRIGGER IF EXISTS update_staking_updated_at ON staking;
DROP TRIGGER IF EXISTS update_store_bridges_updated_at ON store_bridges;
DROP TRIGGER IF EXISTS update_store_networks_updated_at ON store_networks;
DROP TRIGGER IF EXISTS update_bridge_transactions_updated_at ON bridge_transactions;
CREATE TRIGGER update_store_networks_updated_at BEFORE UPDATE ON store_networks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Bridge transactions triggers
DROP TRIGGER IF EXISTS update_bridge_transactions_updated_at ON bridge_transactions;
DROP TRIGGER IF EXISTS update_blocks_updated_at ON blocks;
DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;
DROP TRIGGER IF EXISTS update_contracts_updated_at ON smart_contracts;
DROP TRIGGER IF EXISTS update_collections_updated_at ON nft_collections;
DROP TRIGGER IF EXISTS update_nfts_updated_at ON nfts;
DROP TRIGGER IF EXISTS update_tokens_updated_at ON custom_tokens;
DROP TRIGGER IF EXISTS update_certificates_updated_at ON product_certificates;
DROP TRIGGER IF EXISTS update_validators_updated_at ON validators;
DROP TRIGGER IF EXISTS update_staking_updated_at ON staking;
DROP TRIGGER IF EXISTS update_store_bridges_updated_at ON store_bridges;
DROP TRIGGER IF EXISTS update_store_networks_updated_at ON store_networks;
DROP TRIGGER IF EXISTS update_bridge_transactions_updated_at ON bridge_transactions;
CREATE TRIGGER update_bridge_transactions_updated_at BEFORE UPDATE ON bridge_transactions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- FUNCTIONS
-- ============================================================================

-- Create function for blockchain analytics
CREATE OR REPLACE FUNCTION get_blockchain_analytics(
    p_start_date TIMESTAMP WITH TIME ZONE,
    p_end_date TIMESTAMP WITH TIME ZONE
)
RETURNS TABLE (
    total_blocks BIGINT,
    total_transactions BIGINT,
    total_volume DECIMAL(28,18),
    average_block_time DECIMAL(10,2),
    active_validators BIGINT,
    total_stake DECIMAL(28,18)
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        COUNT(DISTINCT b.id)::BIGINT as total_blocks,
        COUNT(DISTINCT t.id)::BIGINT as total_transactions,
        COALESCE(SUM(t.amount), 0) as total_volume,
        COALESCE(AVG(b.timestamp - LAG(b.timestamp) OVER (ORDER BY b.block_number)), 0) as average_block_time,
        COUNT(DISTINCT v.id)::BIGINT as active_validators,
        COALESCE(SUM(v.stake_amount), 0) as total_stake
    FROM blocks b
    LEFT JOIN transactions t ON b.block_number = t.block_number
    LEFT JOIN validators v ON v.status = 'active'
    WHERE b.created_at BETWEEN p_start_date AND p_end_date;
END;
$$ language 'plpgsql';

-- Create function for transaction validation
CREATE OR REPLACE FUNCTION validate_transaction(
    p_transaction_hash VARCHAR(66),
    p_from_address VARCHAR(42),
    p_to_address VARCHAR(42),
    p_amount DECIMAL(28,18)
)
RETURNS BOOLEAN AS $$
DECLARE
    v_exists BOOLEAN;
    v_balance DECIMAL(28,18);
BEGIN
    -- Check if transaction already exists
    SELECT EXISTS(SELECT 1 FROM transactions WHERE transaction_hash = p_transaction_hash) INTO v_exists;
    IF v_exists THEN
        RETURN FALSE;
    END IF;
    
    -- Check balance (simplified - in real implementation, this would be more complex)
    -- This is a placeholder for actual balance checking logic
    RETURN TRUE;
END;
$$ language 'plpgsql';

-- Create function for NFT ownership tracking
CREATE OR REPLACE FUNCTION track_nft_ownership(
    p_token_id VARCHAR(255),
    p_contract_address VARCHAR(42),
    p_from_address VARCHAR(42),
    p_to_address VARCHAR(42),
    p_transaction_hash VARCHAR(66)
)
RETURNS VOID AS $$
BEGIN
    -- Update NFT ownership
    UPDATE nfts 
    SET 
        owner_address = p_to_address,
        last_transferred_at = NOW(),
        transfer_count = transfer_count + 1
    WHERE token_id = p_token_id AND contract_address = p_contract_address;
    
    -- Log the transfer (could be in a separate transfers table)
    -- This is a simplified version
END;
$$ language 'plpgsql';