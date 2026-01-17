-- USC Platform - Service-04 USC Blockchain Core - PostgreSQL Migration Rollback
-- Database: PostgreSQL (Primary)
-- Purpose: Rollback blockchain tables and functions

-- Drop functions first
DROP FUNCTION IF EXISTS track_nft_ownership(VARCHAR(255), VARCHAR(42), VARCHAR(42), VARCHAR(42), VARCHAR(66));
DROP FUNCTION IF EXISTS validate_transaction(VARCHAR(66), VARCHAR(42), VARCHAR(42), DECIMAL(28,18));
DROP FUNCTION IF EXISTS get_blockchain_analytics(TIMESTAMP WITH TIME ZONE, TIMESTAMP WITH TIME ZONE);
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop triggers
DROP TRIGGER IF EXISTS update_bridge_tx_updated_at ON bridge_transactions;
DROP TRIGGER IF EXISTS update_networks_updated_at ON store_networks;
DROP TRIGGER IF EXISTS update_bridges_updated_at ON store_bridges;
DROP TRIGGER IF EXISTS update_staking_updated_at ON staking;
DROP TRIGGER IF EXISTS update_validators_updated_at ON validators;
DROP TRIGGER IF EXISTS update_certificates_updated_at ON product_certificates;
DROP TRIGGER IF EXISTS update_tokens_updated_at ON custom_tokens;
DROP TRIGGER IF EXISTS update_nfts_updated_at ON nfts;
DROP TRIGGER IF EXISTS update_collections_updated_at ON nft_collections;
DROP TRIGGER IF EXISTS update_contracts_updated_at ON smart_contracts;
DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;
DROP TRIGGER IF EXISTS update_blocks_updated_at ON blocks;

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS bridge_transactions;
DROP TABLE IF EXISTS store_networks;
DROP TABLE IF EXISTS store_bridges;
DROP TABLE IF EXISTS staking;
DROP TABLE IF EXISTS validators;
DROP TABLE IF EXISTS product_certificate_ownership_history;
DROP TABLE IF EXISTS product_certificates;
DROP TABLE IF EXISTS custom_tokens;
DROP TABLE IF EXISTS nfts;
DROP TABLE IF EXISTS nft_collections;
DROP TABLE IF EXISTS smart_contracts;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS blocks;

-- Drop extensions (optional - may be used by other services)
-- DROP EXTENSION IF EXISTS "pgcrypto";
-- DROP EXTENSION IF EXISTS "uuid-ossp";
