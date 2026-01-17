# ⛓️ **SERVICE-04: USC BLOCKCHAIN CORE SERVICE**
## **Priority: CRITICAL | Development Week: 4-5**

---

## 📋 **SERVICE OVERVIEW**

**USC Blockchain Core Service** - High-performance custom blockchain infrastructure cho **Universal Social Coin (USC)**. Enterprise-grade blockchain với Proof of Stake consensus, smart contracts, và Protocol Buffer optimization cho 50K-500M users. **Future-ready architecture** cho expansion thành 4 separate networks.

### **Core Responsibilities**
- USC blockchain network management via high-performance gRPC API
- Block production và validation cho single USC coin
- USC transaction processing và finality via gRPC streaming
- Smart contract execution optimized cho USC operations
- PoS consensus management cho USC network
- USC blockchain data integrity với gRPC-secured communication
- Network health monitoring và validator coordination
- Cross-chain bridge preparation cho future expansion
- **RETRY MECHANISMS** - Intelligent retry logic with exponential backoff for failed transactions
- **FALLBACK OPTIONS** - Graceful degradation when blockchain consensus is slow
- **HIGH AVAILABILITY** - Multi-node deployment with automatic failover

---

## 🏗️ **ARCHITECTURE DESIGN**

### **Technology Stack**
- **Language**: Go 1.24.4+ (files <222 lines)
- **Protocol**: gRPC + Protocol Buffers (high-performance binary serialization)
- **Consensus**: Proof of Stake (PoS) với Byzantine Fault Tolerance
- **Database**: RocksDB (blockchain data), PostgreSQL (metadata), Redis (cache)
- **Networking**: P2P networking với gossip protocol + gRPC service mesh
- **Cryptography**: Ed25519 signatures, SHA-256 hashing, gRPC security
- **Smart Contracts**: WebAssembly (WASM) runtime với gRPC interface
- **Monitoring**: OpenTelemetry, Prometheus metrics
- **Cosmos SDK**: v0.53.4 với 12 custom USC modules
- **Blockchain Layer**: Complete Cosmos SDK integration với custom modules

### **Service Dependencies**
```yaml
dependencies:
  upstream: 
    - Gateway Service (GraphQL queries)
    - User Service (user validation)
    - Auth Service (transaction signing)
  downstream:
    - RocksDB (blockchain state storage)
    - PostgreSQL (metadata, validator info)
    - Redis (mempool, consensus cache)
    - USC Wallet Service (balance updates)
    - USC Reward Service (reward distribution)
  blockchain:
    - Cosmos SDK v0.53.4 (blockchain infrastructure)
    - 12 USC Custom Modules (usc, nft, contract, validator, network, bridge, streaming, certificate, token, store, block, reward)
    - CometBFT (consensus engine)
    - Protocol Buffers (blockchain messages)
```

### **Performance Requirements**
- **TPS**: 10,000+ transactions per second
- **Block Time**: 3 seconds average
- **Finality**: <10 seconds (2-3 block confirmations)
- **Latency**: <100ms for transaction submission
- **Availability**: 99.99% network uptime

### **🔧 OPTIMIZATION STRATEGIES**
- **Retry Mechanisms**: Exponential backoff with maximum 3 retries for failed transactions
- **Fallback Options**: Graceful degradation when consensus is slow or unavailable
- **High Availability**: Multi-node deployment with automatic failover and load balancing
- **Connection Pooling**: Optimized database connections for high throughput
- **Caching Strategy**: Redis-based caching for frequently accessed blockchain data

---

## 🔧 **API DEFINITIONS**

### **gRPC Service Definition**

```protobuf
// USC Blockchain Core Service
syntax = "proto3";
package usc.blockchain;

import "google/protobuf/timestamp.proto";
import "google/protobuf/any.proto";

service USCBlockchainService {
  // Block Operations
  rpc ProduceBlock(ProduceBlockRequest) returns (ProduceBlockResponse);
  rpc ValidateBlock(ValidateBlockRequest) returns (ValidateBlockResponse);
  rpc GetBlock(GetBlockRequest) returns (GetBlockResponse);
  rpc GetBlockByHash(GetBlockByHashRequest) returns (GetBlockResponse);
  rpc GetLatestBlock(GetLatestBlockRequest) returns (GetBlockResponse);
  rpc GetBlockRange(GetBlockRangeRequest) returns (GetBlockRangeResponse);
  
  // Transaction Operations (Blockchain-level - Raw blockchain data)
  rpc SubmitTransaction(SubmitTransactionRequest) returns (SubmitTransactionResponse);
  rpc GetTransaction(GetTransactionRequest) returns (GetTransactionResponse); // Raw blockchain transaction data
  rpc GetTransactionStatus(GetTransactionStatusRequest) returns (TransactionStatusResponse);
  rpc GetPendingTransactions(GetPendingTransactionsRequest) returns (GetPendingTransactionsResponse);
  rpc EstimateTransactionFee(EstimateTransactionFeeRequest) returns (EstimateTransactionFeeResponse);
  
  // USC Coin Operations (Cosmos SDK Integration)
  rpc GetUSCBalance(GetUSCBalanceRequest) returns (GetUSCBalanceResponse);
  rpc TransferUSC(TransferUSCRequest) returns (TransferUSCResponse);
  rpc GetUSCSupply(GetUSCSupplyRequest) returns (GetUSCSupplyResponse);
  rpc GetUSCTransactionHistory(GetUSCTransactionHistoryRequest) returns (GetUSCTransactionHistoryResponse);
  rpc MintUSC(MintUSCRequest) returns (MintUSCResponse);
  rpc BurnUSC(BurnUSCRequest) returns (BurnUSCResponse);
  
  // Smart Contract Operations
  rpc DeployContract(DeployContractRequest) returns (DeployContractResponse);
  rpc ExecuteContract(ExecuteContractRequest) returns (ExecuteContractResponse);
  rpc QueryContract(QueryContractRequest) returns (QueryContractResponse);
  rpc GetContractCode(GetContractCodeRequest) returns (GetContractCodeResponse);
  rpc GetContractStorage(GetContractStorageRequest) returns (GetContractStorageResponse);
  
  // ===== NFT & TOKEN CREATION SUPPORT =====
  
  // NFT Contract Management (Core blockchain operations)
  rpc DeployNFTContract(DeployNFTContractRequest) returns (DeployNFTContractResponse);
  rpc CreateNFTCollection(CreateNFTCollectionRequest) returns (CreateNFTCollectionResponse);
  rpc MintNFT(MintNFTRequest) returns (MintNFTResponse); // Core blockchain NFT minting
  rpc TransferNFT(TransferNFTRequest) returns (TransferNFTResponse); // Core blockchain NFT transfer
  rpc BurnNFT(BurnNFTRequest) returns (BurnNFTResponse);
  rpc GetNFTInfo(GetNFTInfoRequest) returns (GetNFTInfoResponse);
  rpc GetNFTsByOwner(GetNFTsByOwnerRequest) returns (GetNFTsByOwnerResponse);
  
  // Token Creation & Management
  rpc CreateCustomToken(CreateCustomTokenRequest) returns (CreateCustomTokenResponse);
  rpc MintTokens(MintTokensRequest) returns (MintTokensResponse);
  rpc BurnTokens(BurnTokensRequest) returns (BurnTokensResponse);
  rpc GetTokenInfo(GetTokenInfoRequest) returns (GetTokenInfoResponse);
  rpc GetTokenBalance(GetTokenBalanceRequest) returns (GetTokenBalanceResponse);
  
  // Product Tokenization Support
  rpc CreateProductCertificate(CreateProductCertificateRequest) returns (CreateProductCertificateResponse);
  rpc VerifyProductCertificate(VerifyProductCertificateRequest) returns (VerifyProductCertificateResponse);
  rpc TransferProductOwnership(TransferProductOwnershipRequest) returns (TransferProductOwnershipResponse);
  
  // Validator Operations
  rpc RegisterValidator(RegisterValidatorRequest) returns (RegisterValidatorResponse);
  rpc GetValidators(GetValidatorsRequest) returns (GetValidatorsResponse);
  rpc GetValidatorStatus(GetValidatorStatusRequest) returns (GetValidatorStatusResponse);
  rpc StakeUSC(StakeUSCRequest) returns (StakeUSCResponse);
  rpc UnstakeUSC(UnstakeUSCRequest) returns (UnstakeUSCResponse);
  
  // Network Information
  rpc GetNetworkInfo(GetNetworkInfoRequest) returns (GetNetworkInfoResponse);
  rpc GetChainInfo(GetChainInfoRequest) returns (GetChainInfoResponse);
  rpc GetPeers(GetPeersRequest) returns (GetPeersResponse);
  rpc GetNetworkStats(GetNetworkStatsRequest) returns (GetNetworkStatsResponse);
  
  // Real-time Streaming
  rpc StreamBlocks(StreamBlocksRequest) returns (stream BlockStream);
  rpc StreamTransactions(StreamTransactionsRequest) returns (stream TransactionStream);
  rpc StreamValidatorEvents(StreamValidatorEventsRequest) returns (stream ValidatorEvent);
  rpc StreamNetworkEvents(StreamNetworkEventsRequest) returns (stream NetworkEvent);
  
  // ===== STORE NETWORK BRIDGE OPERATIONS =====
  
  // Store Network Bridge Management
  rpc DeployStoreBridge(DeployStoreBridgeRequest) returns (DeployStoreBridgeResponse);
  rpc RegisterStoreNetwork(RegisterStoreNetworkRequest) returns (RegisterStoreNetworkResponse);
  rpc BridgeStoreTokenToUSC(BridgeStoreTokenToUSCRequest) returns (BridgeStoreTokenToUSCResponse);
  rpc BridgeUSCToStoreToken(BridgeUSCToStoreTokenRequest) returns (BridgeUSCToStoreTokenResponse);
  rpc GetStoreBridgeMetrics(GetStoreBridgeMetricsRequest) returns (GetStoreBridgeMetricsResponse);
  rpc ValidateStoreBridge(ValidateStoreBridgeRequest) returns (ValidateStoreBridgeResponse);
  
  // Store Network Integration
  rpc SyncStoreNetworkState(SyncStoreNetworkStateRequest) returns (SyncStoreNetworkStateResponse);
  rpc GetStoreNetworkInfo(GetStoreNetworkInfoRequest) returns (GetStoreNetworkInfoResponse);
  rpc UpdateStoreBridgeConfig(UpdateStoreBridgeConfigRequest) returns (UpdateStoreBridgeConfigResponse);
}

// Core USC Types
message USCBlock {
  string block_hash = 1;
  string previous_hash = 2;
  int64 block_number = 3;
  int64 timestamp = 4;
  string merkle_root = 5;
  string state_root = 6;
  string validator_address = 7;
  string validator_signature = 8;
  repeated USCTransaction transactions = 9;
  BlockMetadata metadata = 10;
  ConsensusInfo consensus_info = 11;
}

message BlockMetadata {
  int32 transaction_count = 1;
  string total_usc_transferred = 2;
  int64 gas_used = 3;
  int64 gas_limit = 4;
  int32 block_size_bytes = 5;
  string difficulty = 6;
}

message ConsensusInfo {
  repeated ValidatorVote votes = 1;
  string consensus_round = 2;
  bool is_finalized = 3;
  int32 confirmation_count = 4;
}

message ValidatorVote {
  string validator_address = 1;
  string vote_signature = 2;
  VoteType vote_type = 3;
  int64 voting_power = 4;
}

enum VoteType {
  PROPOSE = 0;
  PREVOTE = 1;
  PRECOMMIT = 2;
}

// USC Transaction Types
message USCTransaction {
  string tx_hash = 1;
  string from_address = 2;
  string to_address = 3;
  string amount = 4;                   // USC amount (string for precision)
  string fee = 5;                      // Transaction fee in USC
  int64 nonce = 6;
  int64 gas_limit = 7;
  int64 gas_price = 8;
  string memo = 9;
  USCTransactionType tx_type = 10;
  TransactionStatus status = 11;
  string signature = 12;
  google.protobuf.timestamp timestamp = 13;
  int64 block_number = 14;
  int32 transaction_index = 15;
  google.protobuf.Any transaction_data = 16; // Type-specific data
}

enum USCTransactionType {
  TRANSFER = 0;
  MINT = 1;                    // USC minting operations
  BURN = 2;                    // USC burning operations
  SOCIAL_REWARD = 3;           // Social interaction rewards
  VIDEO_REWARD = 4;            // Video completion rewards
  NFT_CREATE = 5;              // NFT creation
  NFT_TRANSFER = 6;            // NFT trading
  ONLINE_REWARD = 7;           // Online activity rewards
  STAKING = 8;                 // USC staking operations
  CONTRACT_DEPLOY = 9;         // Smart contract deployment
  CONTRACT_EXECUTE = 10;       // Smart contract execution
  VALIDATOR_OPERATION = 9;     // Validator operations
  GOVERNANCE = 10;             // Governance voting
  
  // ===== EXPANDED NFT & TOKEN OPERATIONS =====
  NFT_MINT = 11;               // Mint new NFT
  NFT_BURN = 12;               // Burn NFT
  NFT_COLLECTION_CREATE = 13;  // Create NFT collection
  TOKEN_CREATE = 14;           // Create custom token
  TOKEN_MINT = 15;             // Mint tokens
  TOKEN_BURN = 16;             // Burn tokens
  PRODUCT_TOKENIZE = 17;       // Tokenize physical product
  PRODUCT_CERTIFICATE = 18;    // Create product certificate
  PRODUCT_REDEEM = 19;         // Redeem tokenized product
  ROYALTY_PAYMENT = 20;        // NFT royalty distribution
}

enum TransactionStatus {
  BLOCKCHAIN_TRANSACTION_PENDING = 0;
  BLOCKCHAIN_TRANSACTION_CONFIRMED = 1;
  BLOCKCHAIN_TRANSACTION_FAILED = 2;
  BLOCKCHAIN_TRANSACTION_REVERTED = 3;
}

// USC-specific Transaction Data
message SocialRewardData {
  string interaction_type = 1;        // "like", "comment", "share", "friend"
  string content_id = 2;
  string content_creator = 3;
  bool bilateral_reward = 4;
  string quality_multiplier = 5;
  string reward_pool_id = 6;
}

message VideoRewardData {
  string video_id = 1;
  string creator_id = 2;
  int32 completion_percentage = 3;
  int32 watch_duration_seconds = 4;
  string video_quality = 5;           // HD, 4K for bonus calculation
  bool is_live_stream = 6;
}

message NFTData {
  string nft_id = 1;
  string collection_id = 2;
  string metadata_uri = 3;
  string royalty_percentage = 4;
  bool is_creation = 5;               // true for creation, false for transfer
  string trade_price = 6;             // if trading
}

// ===== COMPREHENSIVE NFT & TOKEN TYPES =====

// NFT Contract & Collection
message NFTContract {
  string contract_address = 1;
  string contract_name = 2;
  string contract_symbol = 3;
  string creator_address = 4;
  NFTStandard standard = 5;           // ERC721, ERC1155, USC_NFT
  string base_uri = 6;
  int64 max_supply = 7;
  int64 current_supply = 8;
  string mint_price = 9;              // Price in USC
  bool is_paused = 10;
  NFTContractMetadata metadata = 11;
  google.protobuf.Timestamp created_at = 12;
}

enum NFTStandard {
  USC_NFT = 0;                        // Native USC NFT standard
  ERC721_COMPATIBLE = 1;              // ERC721 compatible
  ERC1155_COMPATIBLE = 2;             // ERC1155 compatible
}

message NFTContractMetadata {
  string description = 1;
  string image_url = 2;
  string external_url = 3;
  string banner_image_url = 4;
  string discord_url = 5;
  string twitter_url = 6;
  map<string, string> custom_fields = 7;
}

// Individual NFT
message NFTToken {
  string token_id = 1;
  string contract_address = 2;
  string owner_address = 3;
  string creator_address = 4;
  
  // Metadata
  string name = 5;
  string description = 6;
  string image_url = 7;
  string animation_url = 8;
  string external_url = 9;
  repeated NFTAttribute attributes = 10;
  
  // Royalty & Economics
  string royalty_percentage = 11;     // Creator royalty
  string royalty_recipient = 12;      // Royalty recipient address
  
  // Status
  NFTTokenStatus status = 13;
  bool is_locked = 14;               // Locked for trading
  
  // Timestamps
  google.protobuf.Timestamp minted_at = 15;
  google.protobuf.Timestamp last_transferred = 16;
  
  // Associated Product (for commerce)
  string associated_product_id = 17;
  ProductTokenizationInfo tokenization_info = 18;
}

message NFTAttribute {
  string trait_type = 1;
  string value = 2;
  string display_type = 3;           // "string", "number", "boost_percentage"
  float max_value = 4;               // For numeric traits
}

enum NFTTokenStatus {
  ACTIVE_TOKEN = 0;
  BURNED = 1;
  FROZEN = 2;
  PENDING_TRANSFER = 3;
}

// Product Tokenization
message ProductTokenizationInfo {
  string product_id = 1;
  TokenizationType tokenization_type = 2;
  string certificate_hash = 3;       // On-chain certificate hash
  bool is_redeemable = 4;            // Can redeem for physical product
  google.protobuf.Timestamp redemption_deadline = 5;
  string redemption_instructions = 6;
  ProductCertificateData certificate_data = 7;
}

enum TokenizationType {
  AUTHENTICITY_CERTIFICATE = 0;      // Product authenticity
  OWNERSHIP_CERTIFICATE = 1;         // Ownership proof
  WARRANTY_TOKEN = 2;                // Warranty certificate
  PROVENANCE_RECORD = 3;             // Supply chain tracking
  LICENSE_TOKEN = 4;                 // Usage license
}

message ProductCertificateData {
  string certificate_id = 1;
  string issuer = 2;
  string verification_method = 3;
  map<string, string> product_details = 4;
  string authentication_code = 5;
  google.protobuf.Timestamp issued_at = 6;
  google.protobuf.Timestamp expires_at = 7;
}

// Custom Token Creation
message CustomToken {
  string token_address = 1;
  string token_name = 2;
  string token_symbol = 3;
  string creator_address = 4;
  int32 decimals = 5;
  string total_supply = 6;
  string current_supply = 7;
  bool is_mintable = 8;
  bool is_burnable = 9;
  bool is_pausable = 10;
  TokenType token_type = 11;
  map<string, string> metadata = 12;
  google.protobuf.Timestamp created_at = 13;
}

enum TokenType {
  UTILITY_TOKEN = 0;                 // Utility token
  GOVERNANCE_TOKEN = 1;              // Governance voting
  REWARD_TOKEN = 2;                  // Reward distribution
  PRODUCT_TOKEN = 3;                 // Product-backed token
  LOYALTY_TOKEN = 4;                 // Loyalty program
}

message OnlineActivityData {
  string activity_type = 1;           // "login", "active_time", "engagement"
  int32 duration_minutes = 2;
  string activity_score = 3;
  bool streak_bonus = 4;
  int32 consecutive_days = 5;
}

// Validator Information
message Validator {
  string validator_address = 1;
  string public_key = 2;
  string voting_power = 3;            // Staked USC amount
  ValidatorStatus status = 4;
  ValidatorMetadata metadata = 5;
  ValidatorPerformance performance = 6;
  google.protobuf.timestamp joined_at = 7;
  google.protobuf.timestamp last_active = 8;
}

enum ValidatorStatus {
  ACTIVE = 0;
  INACTIVE = 1;
  JAILED = 2;
  TOMBSTONED = 3;
}

message ValidatorMetadata {
  string moniker = 1;                 // Validator name
  string website = 2;
  string description = 3;
  string commission_rate = 4;         // Commission percentage
  string min_self_delegation = 5;
}

message ValidatorPerformance {
  int64 blocks_proposed = 1;
  int64 blocks_signed = 2;
  int64 blocks_missed = 3;
  string uptime_percentage = 4;
  string total_commissions_earned = 5;
  int32 jail_count = 6;
}

// Smart Contract Operations
message SmartContract {
  string contract_address = 1;
  string creator_address = 2;
  string code_hash = 3;
  bytes bytecode = 4;
  string abi = 5;                     // JSON ABI
  ContractType contract_type = 6;
  ContractMetadata metadata = 7;
  google.protobuf.timestamp deployed_at = 8;
  int64 deployment_block = 9;
}

enum ContractType {
  USC_TOKEN = 0;                      // USC token contract
  NFT_COLLECTION = 1;                 // NFT collection contract
  REWARD_DISTRIBUTION = 2;            // Reward distribution contract
  GOVERNANCE = 3;                     // Governance contract
  STAKING = 4;                        // Staking contract
  CUSTOM = 5;                         // Custom user contract
  
  // ===== EXPANDED CONTRACT TYPES =====
  CUSTOM_TOKEN = 6;                   // Custom ERC20-style token
  PRODUCT_CERTIFICATE = 7;           // Product authenticity certificate
  MARKETPLACE = 8;                    // NFT marketplace contract
  ROYALTY_DISTRIBUTOR = 9;           // Royalty distribution contract
  MULTI_SIG_WALLET = 10;             // Multi-signature wallet
  AUCTION_HOUSE = 11;                // Auction contract
  ESCROW = 12;                       // Escrow contract
}

message ContractMetadata {
  string name = 1;
  string description = 2;
  string version = 3;
  repeated string tags = 4;
  bool is_verified = 5;
  string verification_status = 6;
}

// Network Information
message NetworkInfo {
  string network_id = 1;
  string chain_id = 2;
  int64 latest_block_number = 3;
  string latest_block_hash = 4;
  int64 total_validators = 5;
  int64 active_validators = 6;
  string total_staked_usc = 7;
  string circulating_supply = 8;
  string total_supply = 9;
  NetworkHealth health = 10;
}

message NetworkHealth {
  string status = 1;                  // "healthy", "degraded", "critical"
  float block_time_average = 2;       // Average block time in seconds
  float transaction_throughput = 3;   // TPS
  int32 pending_transactions = 4;
  float network_hashrate = 5;
  int32 peer_count = 6;
}

// Streaming Types
message BlockStream {
  USCBlock block = 1;
  StreamEventType event_type = 2;
  google.protobuf.timestamp timestamp = 3;
}

message TransactionStream {
  USCTransaction transaction = 1;
  StreamEventType event_type = 2;
  google.protobuf.timestamp timestamp = 3;
}

message ValidatorEvent {
  string validator_address = 1;
  ValidatorEventType event_type = 2;
  google.protobuf.Any event_data = 3;
  google.protobuf.timestamp timestamp = 4;
}

enum ValidatorEventType {
  VALIDATOR_JOINED = 0;
  VALIDATOR_LEFT = 1;
  VALIDATOR_JAILED = 2;
  VALIDATOR_UNJAILED = 3;
  BLOCK_PROPOSED = 4;
  BLOCK_SIGNED = 5;
  BLOCK_MISSED = 6;
}

enum StreamEventType {
  CREATED = 0;
  CONFIRMED = 1;
  FINALIZED = 2;
  FAILED = 3;
}

// Request/Response Messages
message SubmitTransactionRequest {
  USCTransaction transaction = 1;
  bool broadcast = 2;                 // Whether to broadcast immediately
  bool wait_for_confirmation = 3;     // Whether to wait for block inclusion
}

message SubmitTransactionResponse {
  bool success = 1;
  string error_message = 2;
  string transaction_hash = 3;
  int64 block_number = 4;
  string estimated_confirmation_time = 5;
  TransactionReceipt receipt = 6;
}

message TransactionReceipt {
  string transaction_hash = 1;
  int64 block_number = 2;
  int32 transaction_index = 3;
  int64 gas_used = 4;
  string gas_price = 5;
  string fee_paid = 6;
  TransactionStatus status = 7;
  repeated EventLog logs = 8;
}

message EventLog {
  string contract_address = 1;
  repeated string topics = 2;
  bytes data = 3;
  int32 log_index = 4;
}

message GetUSCBalanceRequest {
  string address = 1;
  int64 block_number = 2;             // Optional: balance at specific block
}

message GetUSCBalanceResponse {
  string balance = 1;                 // Available balance
  string pending_balance = 2;         // Pending transactions
  string staked_balance = 3;          // Staked USC
  string total_balance = 4;           // Total balance
  int64 nonce = 5;                    // Account nonce
}

message TransferUSCRequest {
  string from_address = 1;
  string to_address = 2;
  string amount = 3;
  string memo = 4;
  int64 gas_limit = 5;
  int64 gas_price = 6;
  string private_key = 7;             // For signing (should be encrypted)
}

message TransferUSCResponse {
  bool success = 1;
  string error_message = 2;
  string transaction_hash = 3;
  TransactionReceipt receipt = 4;
}
```

---

## 🗄️ **DATABASE SCHEMA**

### **RocksDB Schema (Blockchain State)**

```yaml
# RocksDB Key-Value Store for USC Blockchain
# Optimized for high-performance blockchain operations

rocksdb_schemas:
  # Block Storage
  blocks:
    key_format: "block:{block_number}"
    value: "protobuf serialized USCBlock"
    example: "block:123456 -> USCBlock{...}"
    
  block_by_hash:
    key_format: "block_hash:{block_hash}"
    value: "block_number"
    example: "block_hash:0xabc123... -> 123456"
    
  # Transaction Storage
  transactions:
    key_format: "tx:{tx_hash}"
    value: "protobuf serialized USCTransaction"
    example: "tx:0xdef456... -> USCTransaction{...}"
    
  tx_by_block:
    key_format: "block_txs:{block_number}"
    value: "list of transaction hashes"
    example: "block_txs:123456 -> [0xdef456..., 0xghi789...]"
    
  # Account State
  account_balance:
    key_format: "balance:{address}"
    value: "USC balance (string)"
    example: "balance:usc1abc... -> 1000.500000000000000000"
    
  account_nonce:
    key_format: "nonce:{address}"
    value: "account nonce (uint64)"
    example: "nonce:usc1abc... -> 42"
    
  account_staked:
    key_format: "staked:{address}"
    value: "staked USC amount (string)"
    example: "staked:usc1abc... -> 500.000000000000000000"
    
  # Smart Contract State
  contract_code:
    key_format: "code:{contract_address}"
    value: "contract bytecode"
    example: "code:usc1contract... -> [bytecode]"
    
  contract_storage:
    key_format: "storage:{contract_address}:{storage_key}"
    value: "storage value"
    example: "storage:usc1contract...:0x01 -> storage_value"
    
  # Validator State
  validator_info:
    key_format: "validator:{validator_address}"
    value: "protobuf serialized Validator"
    example: "validator:usc1val... -> Validator{...}"
    
  validator_power:
    key_format: "power:{validator_address}"
    value: "voting power (uint64)"
    example: "power:usc1val... -> 1000000000000000000"
    
  # Consensus State
  latest_block:
    key_format: "latest_block"
    value: "latest block number"
    example: "latest_block -> 123456"
    
  chain_state:
    key_format: "chain_state"
    value: "protobuf serialized ChainState"
    example: "chain_state -> ChainState{...}"
    
  # Indexing for Fast Queries
  tx_by_address:
    key_format: "addr_tx:{address}:{block_number}:{tx_index}"
    value: "transaction hash"
    example: "addr_tx:usc1abc...:123456:0 -> 0xdef456..."
    
  block_rewards:
    key_format: "rewards:{block_number}"
    value: "total rewards distributed in block"
    example: "rewards:123456 -> 100.000000000000000000"
```

### **PostgreSQL Schema (Metadata & Analytics)**

```sql
-- USC Blockchain Metadata Tables

-- Chain Information
CREATE TABLE usc_chain_info (
    id SERIAL PRIMARY KEY,
    chain_id VARCHAR(50) NOT NULL,
    network_id VARCHAR(50) NOT NULL,
    genesis_hash VARCHAR(66) NOT NULL,
    genesis_time TIMESTAMP NOT NULL,
    latest_block_number BIGINT NOT NULL DEFAULT 0,
    latest_block_hash VARCHAR(66),
    total_validators INTEGER DEFAULT 0,
    active_validators INTEGER DEFAULT 0,
    total_supply DECIMAL(28,18) DEFAULT 1000000000,
    circulating_supply DECIMAL(28,18) DEFAULT 0,
    total_staked DECIMAL(28,18) DEFAULT 0,
    updated_at TIMESTAMP DEFAULT NOW(),
    
    INDEX idx_chain_latest_block (latest_block_number),
    INDEX idx_chain_updated (updated_at)
);

-- Validator Registry
CREATE TABLE usc_validators (
    validator_address VARCHAR(42) PRIMARY KEY,
    public_key VARCHAR(88) NOT NULL,
    moniker VARCHAR(100),
    website VARCHAR(255),
    description TEXT,
    commission_rate DECIMAL(5,4) NOT NULL DEFAULT 0.1000,
    min_self_delegation DECIMAL(28,18) NOT NULL,
    voting_power DECIMAL(28,18) NOT NULL DEFAULT 0,
    status INTEGER NOT NULL DEFAULT 0,
    jailed BOOLEAN DEFAULT FALSE,
    jail_until TIMESTAMP,
    tombstoned BOOLEAN DEFAULT FALSE,
    blocks_proposed BIGINT DEFAULT 0,
    blocks_signed BIGINT DEFAULT 0,
    blocks_missed BIGINT DEFAULT 0,
    total_commissions DECIMAL(28,18) DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    INDEX idx_validators_status (status),
    INDEX idx_validators_power (voting_power),
    INDEX idx_validators_performance (blocks_signed, blocks_missed),
    INDEX idx_validators_jailed (jailed, jail_until)
);

-- Block Analytics
CREATE TABLE usc_block_analytics (
    block_number BIGINT PRIMARY KEY,
    block_hash VARCHAR(66) UNIQUE NOT NULL,
    validator_address VARCHAR(42) NOT NULL REFERENCES usc_validators(validator_address),
    timestamp TIMESTAMP NOT NULL,
    transaction_count INTEGER NOT NULL DEFAULT 0,
    total_usc_transferred DECIMAL(28,18) DEFAULT 0,
    gas_used BIGINT DEFAULT 0,
    gas_limit BIGINT DEFAULT 0,
    block_size_bytes INTEGER DEFAULT 0,
    processing_time_ms INTEGER,
    is_finalized BOOLEAN DEFAULT FALSE,
    finalized_at TIMESTAMP,
    
    INDEX idx_block_validator (validator_address),
    INDEX idx_block_timestamp (timestamp),
    INDEX idx_block_finalized (is_finalized, finalized_at),
    INDEX idx_block_metrics (transaction_count, gas_used)
);

-- Transaction Analytics
CREATE TABLE usc_transaction_analytics (
    tx_hash VARCHAR(66) PRIMARY KEY,
    block_number BIGINT REFERENCES usc_block_analytics(block_number),
    tx_index INTEGER NOT NULL,
    from_address VARCHAR(42) NOT NULL,
    to_address VARCHAR(42) NOT NULL,
    amount DECIMAL(28,18) NOT NULL,
    fee DECIMAL(28,18) NOT NULL,
    gas_used BIGINT,
    gas_price BIGINT,
    tx_type INTEGER NOT NULL,
    status INTEGER NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    processing_time_ms INTEGER,
    error_message TEXT,
    
    INDEX idx_tx_block (block_number, tx_index),
    INDEX idx_tx_from (from_address),
    INDEX idx_tx_to (to_address),
    INDEX idx_tx_type (tx_type),
    INDEX idx_tx_status (status),
    INDEX idx_tx_timestamp (timestamp),
    INDEX idx_tx_amount (amount),
    
    -- Partition by month for performance
    PARTITION BY RANGE (timestamp)
);

-- Create monthly partitions for transaction analytics
CREATE TABLE usc_transaction_analytics_2024_01 PARTITION OF usc_transaction_analytics
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');
CREATE TABLE usc_transaction_analytics_2024_02 PARTITION OF usc_transaction_analytics
    FOR VALUES FROM ('2024-02-01') TO ('2024-03-01');
-- Continue for each month...

-- Cosmos SDK Integration Tables
CREATE TABLE cosmos_sdk_modules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    module_name VARCHAR(64) UNIQUE NOT NULL,
    module_version VARCHAR(32) NOT NULL,
    module_address VARCHAR(64) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE cosmos_sdk_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    module_name VARCHAR(64) NOT NULL,
    event_type VARCHAR(64) NOT NULL,
    event_data JSONB,
    block_height BIGINT NOT NULL,
    tx_hash VARCHAR(64),
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Smart Contract Registry
CREATE TABLE usc_smart_contracts (
    contract_address VARCHAR(42) PRIMARY KEY,
    creator_address VARCHAR(42) NOT NULL,
    code_hash VARCHAR(66) NOT NULL,
    contract_type INTEGER NOT NULL,
    name VARCHAR(100),
    description TEXT,
    version VARCHAR(20),
    abi TEXT,
    bytecode BYTEA,
    is_verified BOOLEAN DEFAULT FALSE,
    verification_date TIMESTAMP,
    deployed_at TIMESTAMP NOT NULL,
    deployment_block BIGINT NOT NULL,
    deployment_tx_hash VARCHAR(66) NOT NULL,
    total_calls BIGINT DEFAULT 0,
    total_gas_used BIGINT DEFAULT 0,
    
    INDEX idx_contracts_creator (creator_address),
    INDEX idx_contracts_type (contract_type),
    INDEX idx_contracts_deployed (deployed_at),
    INDEX idx_contracts_verified (is_verified),
    INDEX idx_contracts_usage (total_calls, total_gas_used)
);

-- ===== NFT & TOKEN INFRASTRUCTURE TABLES =====

-- NFT Contracts (Collections)
CREATE TABLE usc_nft_contracts (
    contract_address VARCHAR(42) PRIMARY KEY,
    contract_name VARCHAR(200) NOT NULL,
    contract_symbol VARCHAR(20) NOT NULL,
    creator_address VARCHAR(42) NOT NULL,
    standard INTEGER NOT NULL,           -- NFTStandard enum
    base_uri TEXT,
    max_supply BIGINT,
    current_supply BIGINT DEFAULT 0,
    mint_price DECIMAL(28,18),
    is_paused BOOLEAN DEFAULT FALSE,
    
    -- Metadata
    description TEXT,
    image_url TEXT,
    external_url TEXT,
    banner_image_url TEXT,
    
    -- Blockchain Info
    deployed_at TIMESTAMP NOT NULL,
    deployment_block BIGINT NOT NULL,
    deployment_tx_hash VARCHAR(66) NOT NULL,
    
    INDEX idx_nft_contracts_creator (creator_address),
    INDEX idx_nft_contracts_standard (standard),
    INDEX idx_nft_contracts_supply (current_supply, max_supply),
    INDEX idx_nft_contracts_deployed (deployed_at),
    
    FOREIGN KEY (contract_address) REFERENCES usc_smart_contracts(contract_address)
);

-- Individual NFT Tokens
CREATE TABLE usc_nft_tokens (
    token_id VARCHAR(100) NOT NULL,
    contract_address VARCHAR(42) NOT NULL,
    owner_address VARCHAR(42) NOT NULL,
    creator_address VARCHAR(42) NOT NULL,
    
    -- Metadata
    name VARCHAR(500),
    description TEXT,
    image_url TEXT,
    animation_url TEXT,
    external_url TEXT,
    metadata_uri TEXT,
    
    -- Economics
    royalty_percentage DECIMAL(5,4) DEFAULT 0,
    royalty_recipient VARCHAR(42),
    
    -- Status
    status INTEGER DEFAULT 0,           -- NFTTokenStatus enum
    is_locked BOOLEAN DEFAULT FALSE,
    
    -- Product Association
    associated_product_id UUID,
    tokenization_type INTEGER,          -- TokenizationType enum
    certificate_hash VARCHAR(66),
    is_redeemable BOOLEAN DEFAULT FALSE,
    redemption_deadline TIMESTAMP,
    
    -- Timestamps
    minted_at TIMESTAMP NOT NULL,
    minted_block BIGINT NOT NULL,
    minted_tx_hash VARCHAR(66) NOT NULL,
    last_transferred TIMESTAMP,
    
    PRIMARY KEY (token_id, contract_address),
    INDEX idx_nft_tokens_owner (owner_address),
    INDEX idx_nft_tokens_creator (creator_address),
    INDEX idx_nft_tokens_product (associated_product_id),
    INDEX idx_nft_tokens_status (status),
    INDEX idx_nft_tokens_minted (minted_at),
    
    FOREIGN KEY (contract_address) REFERENCES usc_nft_contracts(contract_address)
);

-- NFT Attributes
CREATE TABLE usc_nft_attributes (
    token_id VARCHAR(100) NOT NULL,
    contract_address VARCHAR(42) NOT NULL,
    trait_type VARCHAR(100) NOT NULL,
    trait_value VARCHAR(500) NOT NULL,
    display_type VARCHAR(50),           -- "string", "number", "boost_percentage"
    max_value DECIMAL(15,6),           -- For numeric traits
    
    INDEX idx_nft_attributes_token (token_id, contract_address),
    INDEX idx_nft_attributes_trait (trait_type),
    INDEX idx_nft_attributes_value (trait_value),
    
    FOREIGN KEY (token_id, contract_address) REFERENCES usc_nft_tokens(token_id, contract_address) ON DELETE CASCADE
);

-- Custom Tokens
CREATE TABLE usc_custom_tokens (
    token_address VARCHAR(42) PRIMARY KEY,
    token_name VARCHAR(200) NOT NULL,
    token_symbol VARCHAR(20) NOT NULL,
    creator_address VARCHAR(42) NOT NULL,
    decimals INTEGER NOT NULL DEFAULT 18,
    total_supply DECIMAL(38,18) NOT NULL,
    current_supply DECIMAL(38,18) DEFAULT 0,
    is_mintable BOOLEAN DEFAULT TRUE,
    is_burnable BOOLEAN DEFAULT TRUE,
    is_pausable BOOLEAN DEFAULT TRUE,
    token_type INTEGER NOT NULL,        -- TokenType enum
    
    -- Status
    is_paused BOOLEAN DEFAULT FALSE,
    
    -- Blockchain Info
    deployed_at TIMESTAMP NOT NULL,
    deployment_block BIGINT NOT NULL,
    deployment_tx_hash VARCHAR(66) NOT NULL,
    
    INDEX idx_custom_tokens_creator (creator_address),
    INDEX idx_custom_tokens_type (token_type),
    INDEX idx_custom_tokens_symbol (token_symbol),
    INDEX idx_custom_tokens_deployed (deployed_at),
    
    FOREIGN KEY (token_address) REFERENCES usc_smart_contracts(contract_address)
);

-- Token Balances
CREATE TABLE usc_token_balances (
    token_address VARCHAR(42) NOT NULL,
    owner_address VARCHAR(42) NOT NULL,
    balance DECIMAL(38,18) NOT NULL DEFAULT 0,
    last_updated TIMESTAMP DEFAULT NOW(),
    last_updated_block BIGINT,
    last_updated_tx_hash VARCHAR(66),
    
    PRIMARY KEY (token_address, owner_address),
    INDEX idx_token_balances_owner (owner_address),
    INDEX idx_token_balances_balance (balance),
    INDEX idx_token_balances_updated (last_updated),
    
    FOREIGN KEY (token_address) REFERENCES usc_custom_tokens(token_address)
);

-- Product Certificates (Tokenized Products)
CREATE TABLE usc_product_certificates (
    certificate_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL,
    token_id VARCHAR(100) NOT NULL,
    contract_address VARCHAR(42) NOT NULL,
    
    -- Certificate Details
    certificate_type INTEGER NOT NULL,  -- TokenizationType enum
    issuer_address VARCHAR(42) NOT NULL,
    verification_method TEXT,
    authentication_code VARCHAR(100),
    
    -- Product Details (stored as JSONB for flexibility)
    product_details JSONB,
    
    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    is_redeemed BOOLEAN DEFAULT FALSE,
    
    -- Timestamps
    issued_at TIMESTAMP NOT NULL,
    issued_block BIGINT NOT NULL,
    issued_tx_hash VARCHAR(66) NOT NULL,
    expires_at TIMESTAMP,
    redeemed_at TIMESTAMP,
    redeemed_tx_hash VARCHAR(66),
    
    INDEX idx_product_certs_product (product_id),
    INDEX idx_product_certs_token (token_id, contract_address),
    INDEX idx_product_certs_type (certificate_type),
    INDEX idx_product_certs_issuer (issuer_address),
    INDEX idx_product_certs_status (is_active, is_redeemed),
    INDEX idx_product_certs_issued (issued_at),
    
    FOREIGN KEY (token_id, contract_address) REFERENCES usc_nft_tokens(token_id, contract_address)
);

-- Network Performance Metrics
CREATE TABLE usc_network_metrics (
    id BIGSERIAL PRIMARY KEY,
    block_number BIGINT NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    
    -- Performance Metrics
    block_time_seconds DECIMAL(8,3) NOT NULL,
    transactions_per_second DECIMAL(10,3) NOT NULL,
    gas_price_gwei DECIMAL(12,3) NOT NULL,
    pending_tx_count INTEGER NOT NULL,
    
    -- Network Health
    validator_count INTEGER NOT NULL,
    online_validator_count INTEGER NOT NULL,
    network_hashrate DECIMAL(20,3),
    peer_count INTEGER NOT NULL,
    
    -- Economic Metrics
    total_supply DECIMAL(28,18) NOT NULL,
    circulating_supply DECIMAL(28,18) NOT NULL,
    total_staked DECIMAL(28,18) NOT NULL,
    staking_ratio DECIMAL(5,4) NOT NULL,
    
    INDEX idx_network_metrics_block (block_number),
    INDEX idx_network_metrics_timestamp (timestamp),
    INDEX idx_network_metrics_performance (block_time_seconds, transactions_per_second)
);

-- USC Address Analytics
CREATE TABLE usc_address_analytics (
    address VARCHAR(42) PRIMARY KEY,
    address_type INTEGER NOT NULL,     -- 0=user, 1=contract, 2=validator
    balance DECIMAL(28,18) DEFAULT 0,
    staked_balance DECIMAL(28,18) DEFAULT 0,
    nonce BIGINT DEFAULT 0,
    
    -- Transaction Statistics
    total_sent_count BIGINT DEFAULT 0,
    total_received_count BIGINT DEFAULT 0,
    total_sent_amount DECIMAL(28,18) DEFAULT 0,
    total_received_amount DECIMAL(28,18) DEFAULT 0,
    
    -- Activity Statistics
    first_tx_block BIGINT,
    last_tx_block BIGINT,
    most_active_day DATE,
    max_daily_tx_count INTEGER DEFAULT 0,
    
    updated_at TIMESTAMP DEFAULT NOW(),
    
    INDEX idx_address_type (address_type),
    INDEX idx_address_balance (balance),
    INDEX idx_address_activity (last_tx_block, total_sent_count)
);
```

### **Redis Schema (Mempool & Cache)**

```yaml
# Redis Keys for USC Blockchain
redis_keys:
  # Mempool (Pending Transactions)
  mempool_tx:{tx_hash}:
    ttl: 3600  # 1 hour
    data: "serialized pending transaction"
    
  mempool_by_sender:{address}:
    ttl: 3600  # 1 hour
    data: "list of pending tx hashes from address"
    
  mempool_stats:
    ttl: 60  # 1 minute
    data: "mempool statistics (count, size, fee stats)"
    
  # Consensus Cache
  consensus_round:{block_number}:
    ttl: 300  # 5 minutes
    data: "consensus round information"
    
  validator_votes:{block_number}:
    ttl: 300  # 5 minutes
    data: "validator votes for block"
    
  # Block Cache
  latest_blocks:
    ttl: 30  # 30 seconds
    data: "list of latest 10 block hashes"
    
  block_cache:{block_hash}:
    ttl: 3600  # 1 hour
    data: "serialized block data"
    
  # Transaction Cache
  tx_cache:{tx_hash}:
    ttl: 1800  # 30 minutes
    data: "serialized transaction data"
    
  address_balance:{address}:
    ttl: 60  # 1 minute
    data: "cached balance information"
    
  # Network Status Cache
  network_status:
    ttl: 30  # 30 seconds
    data: "current network status and health"
    
  validator_status:
    ttl: 60  # 1 minute
    data: "all validators status and voting power"
    
  # Smart Contract Cache
  contract_code:{contract_address}:
    ttl: 3600  # 1 hour
    data: "contract bytecode and metadata"
    
  contract_state:{contract_address}:{state_key}:
    ttl: 300  # 5 minutes
    data: "contract state value"
    
  # Performance Metrics
  blockchain_metrics:
    ttl: 60  # 1 minute
    data: "current blockchain performance metrics"
    
  tps_counter:
    ttl: 60  # 1 minute
    data: "transactions per second counter"
```

---

## 💻 **CORE FEATURES (NO CODE)**

### **1. USC Block Production & Validation**

```go
// USC Block Production Implementation
package blockchain

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "time"
    pb "usc-blockchain/proto"
)

type USCBlockchainService struct {
    pb.UnimplementedUSCBlockchainServiceServer
    rocksDB         *RocksDBClient
    postgres        *PostgreSQLClient
    redis           *RedisClient
    consensusEngine *PoSConsensusEngine
    validatorMgr    *ValidatorManager
    contractEngine  *WASMContractEngine
    p2pNetwork      *P2PNetwork
    cosmosApp       *app.USCApp
    blockchainStorage *storage.StateManager
    rocksDBManager  *storage.RocksDBManager
}

// Produce USC Block with PoS Consensus
func (bs *USCBlockchainService) ProduceBlock(ctx context.Context, req *pb.ProduceBlockRequest) (*pb.ProduceBlockResponse, error) {
    // Verify validator is authorized to produce block
    validator, err := bs.validatorMgr.GetValidator(req.ValidatorAddress)
    if err != nil {
        return &pb.ProduceBlockResponse{
            Success: false,
            Error:   fmt.Sprintf("Invalid validator: %v", err),
        }, nil
    }
    
    if !bs.consensusEngine.IsValidatorTurn(validator.ValidatorAddress) {
        return &pb.ProduceBlockResponse{
            Success: false,
            Error:   "Not validator's turn to produce block",
        }, nil
    }
    
    // Get previous block
    previousBlock, err := bs.GetLatestBlock(ctx, &pb.GetLatestBlockRequest{})
    if err != nil {
        return nil, fmt.Errorf("failed to get previous block: %v", err)
    }
    
    // Create new block
    newBlock := &pb.USCBlock{
        BlockNumber:      previousBlock.Block.BlockNumber + 1,
        PreviousHash:     previousBlock.Block.BlockHash,
        Timestamp:        time.Now().Unix(),
        ValidatorAddress: req.ValidatorAddress,
        Transactions:     make([]*pb.USCTransaction, 0),
        Metadata:         &pb.BlockMetadata{},
        ConsensusInfo:    &pb.ConsensusInfo{},
    }
    
    // Get pending transactions from mempool
    pendingTxs, err := bs.getPendingTransactions(req.MaxTransactions)
    if err != nil {
        return nil, fmt.Errorf("failed to get pending transactions: %v", err)
    }
    
    // Execute transactions and update state
    validTxs := make([]*pb.USCTransaction, 0)
    totalGasUsed := int64(0)
    totalFees := "0"
    totalRewards := "0"
    
    for _, tx := range pendingTxs {
        // Validate transaction
        if !bs.validateTransaction(tx) {
            continue
        }
        
        // Execute transaction
        receipt, err := bs.executeTransaction(tx, newBlock.BlockNumber)
        if err != nil {
            // Mark transaction as failed
            tx.Status = pb.TransactionStatus_FAILED
            bs.updateTransactionStatus(tx.TxHash, pb.TransactionStatus_FAILED, err.Error())
            continue
        }
        
        // Add to block
        tx.Status = pb.TransactionStatus_CONFIRMED
        tx.BlockNumber = newBlock.BlockNumber
        tx.TransactionIndex = int32(len(validTxs))
        validTxs = append(validTxs, tx)
        
        totalGasUsed += receipt.GasUsed
        totalFees = addUSCAmounts(totalFees, tx.Fee)
        
        // Track rewards for USC reward transactions
        if tx.TxType == pb.USCTransactionType_SOCIAL_REWARD ||
           tx.TxType == pb.USCTransactionType_VIDEO_REWARD ||
           tx.TxType == pb.USCTransactionType_ONLINE_REWARD {
            totalRewards = addUSCAmounts(totalRewards, tx.Amount)
        }
    }
    
    newBlock.Transactions = validTxs
    
    // Calculate Merkle root
    newBlock.MerkleRoot = bs.calculateMerkleRoot(validTxs)
    
    // Calculate state root after all transactions
    newBlock.StateRoot = bs.calculateStateRoot()
    
    // Update block metadata
    newBlock.Metadata = &pb.BlockMetadata{
        TransactionCount:        int32(len(validTxs)),
        TotalUscTransferred:     bs.calculateTotalTransferred(validTxs),
        TotalFeesCollected:      totalFees,
        TotalRewardsDistributed: totalRewards,
        GasUsed:                 totalGasUsed,
        GasLimit:                req.GasLimit,
        BlockSizeBytes:          bs.calculateBlockSize(newBlock),
    }
    
    // Calculate block hash
    newBlock.BlockHash = bs.calculateBlockHash(newBlock)
    
    // Sign block
    signature, err := bs.signBlock(newBlock, validator.PrivateKey)
    if err != nil {
        return nil, fmt.Errorf("failed to sign block: %v", err)
    }
    newBlock.ValidatorSignature = signature
    
    // Validate block before committing
    isValid, err := bs.ValidateBlock(ctx, &pb.ValidateBlockRequest{Block: newBlock})
    if err != nil || !isValid.IsValid {
        return &pb.ProduceBlockResponse{
            Success: false,
            Error:   fmt.Sprintf("Block validation failed: %v", err),
        }, nil
    }
    
    // Commit block to storage
    err = bs.commitBlock(newBlock)
    if err != nil {
        return nil, fmt.Errorf("failed to commit block: %v", err)
    }
    
    // Broadcast block to network
    go bs.p2pNetwork.BroadcastBlock(newBlock)
    
    // Update validator performance
    bs.updateValidatorPerformance(req.ValidatorAddress, "block_produced")
    
    // Update chain info
    bs.updateChainInfo(newBlock)
    
    // Distribute block rewards to validator
    go bs.distributeBlockReward(req.ValidatorAddress, totalFees)
    
    return &pb.ProduceBlockResponse{
        Success: true,
        Block:   newBlock,
        Message: fmt.Sprintf("Block %d produced successfully", newBlock.BlockNumber),
    }, nil
}

// Validate USC Block
func (bs *USCBlockchainService) ValidateBlock(ctx context.Context, req *pb.ValidateBlockRequest) (*pb.ValidateBlockResponse, error) {
    block := req.Block
    
    // Basic block validation
    if block.BlockNumber <= 0 {
        return &pb.ValidateBlockResponse{
            IsValid: false,
            Error:   "Invalid block number",
        }, nil
    }
    
    // Validate previous block hash
    if block.BlockNumber > 1 {
        previousBlock, err := bs.GetBlock(ctx, &pb.GetBlockRequest{
            BlockNumber: block.BlockNumber - 1,
        })
        if err != nil {
            return &pb.ValidateBlockResponse{
                IsValid: false,
                Error:   "Previous block not found",
            }, nil
        }
        
        if block.PreviousHash != previousBlock.Block.BlockHash {
            return &pb.ValidateBlockResponse{
                IsValid: false,
                Error:   "Invalid previous block hash",
            }, nil
        }
    }
    
    // Validate validator signature
    isValidSignature := bs.verifyValidatorSignature(block)
    if !isValidSignature {
        return &pb.ValidateBlockResponse{
            IsValid: false,
            Error:   "Invalid validator signature",
        }, nil
    }
    
    // Validate validator authority
    isAuthorized, err := bs.consensusEngine.ValidateValidatorAuthority(
        block.ValidatorAddress, 
        block.BlockNumber,
    )
    if err != nil || !isAuthorized {
        return &pb.ValidateBlockResponse{
            IsValid: false,
            Error:   "Validator not authorized for this block",
        }, nil
    }
    
    // Validate transactions
    for i, tx := range block.Transactions {
        isValidTx, err := bs.validateTransactionInBlock(tx, block, i)
        if err != nil || !isValidTx {
            return &pb.ValidateBlockResponse{
                IsValid: false,
                Error:   fmt.Sprintf("Invalid transaction at index %d: %v", i, err),
            }, nil
        }
    }
    
    // Validate Merkle root
    calculatedMerkleRoot := bs.calculateMerkleRoot(block.Transactions)
    if calculatedMerkleRoot != block.MerkleRoot {
        return &pb.ValidateBlockResponse{
            IsValid: false,
            Error:   "Invalid Merkle root",
        }, nil
    }
    
    // Validate state root
    calculatedStateRoot := bs.calculateStateRootForBlock(block)
    if calculatedStateRoot != block.StateRoot {
        return &pb.ValidateBlockResponse{
            IsValid: false,
            Error:   "Invalid state root",
        }, nil
    }
    
    // Validate block hash
    calculatedHash := bs.calculateBlockHash(block)
    if calculatedHash != block.BlockHash {
        return &pb.ValidateBlockResponse{
            IsValid: false,
            Error:   "Invalid block hash",
        }, nil
    }
    
    // Validate consensus rules
    consensusValid, err := bs.consensusEngine.ValidateBlock(block)
    if err != nil || !consensusValid {
        return &pb.ValidateBlockResponse{
            IsValid: false,
            Error:   fmt.Sprintf("Consensus validation failed: %v", err),
        }, nil
    }
    
    return &pb.ValidateBlockResponse{
        IsValid: true,
        Message: "Block validation successful",
    }, nil
}

// Execute USC Transaction
func (bs *USCBlockchainService) executeTransaction(tx *pb.USCTransaction, blockNumber int64) (*pb.TransactionReceipt, error) {
    receipt := &pb.TransactionReceipt{
        TransactionHash: tx.TxHash,
        BlockNumber:     blockNumber,
        Status:          pb.TransactionStatus_CONFIRMED,
        Logs:            make([]*pb.EventLog, 0),
    }
    
    // Validate transaction signature
    if !bs.verifyTransactionSignature(tx) {
        return nil, fmt.Errorf("invalid transaction signature")
    }
    
    // Check sender balance
    senderBalance, err := bs.getUSCBalance(tx.FromAddress)
    if err != nil {
        return nil, fmt.Errorf("failed to get sender balance: %v", err)
    }
    
    totalAmount := addUSCAmounts(tx.Amount, tx.Fee)
    if compareUSCAmounts(senderBalance, totalAmount) < 0 {
        return nil, fmt.Errorf("insufficient balance")
    }
    
    // Check and update nonce
    expectedNonce, err := bs.getAccountNonce(tx.FromAddress)
    if err != nil {
        return nil, fmt.Errorf("failed to get account nonce: %v", err)
    }
    
    if tx.Nonce != expectedNonce {
        return nil, fmt.Errorf("invalid nonce: expected %d, got %d", expectedNonce, tx.Nonce)
    }
    
    // Execute based on transaction type
    switch tx.TxType {
    case pb.USCTransactionType_TRANSFER:
        err = bs.executeUSCTransfer(tx, receipt)
        
    case pb.USCTransactionType_SOCIAL_REWARD:
        err = bs.executeSocialReward(tx, receipt)
        
    case pb.USCTransactionType_VIDEO_REWARD:
        err = bs.executeVideoReward(tx, receipt)
        
    case pb.USCTransactionType_NFT_CREATE:
        err = bs.executeNFTCreation(tx, receipt)
        
    case pb.USCTransactionType_CONTRACT_EXECUTE:
        err = bs.executeSmartContract(tx, receipt)
        
    case pb.USCTransactionType_STAKING:
        err = bs.executeStaking(tx, receipt)
        
    default:
        return nil, fmt.Errorf("unsupported transaction type: %v", tx.TxType)
    }
    
    if err != nil {
        receipt.Status = pb.TransactionStatus_FAILED
        return receipt, err
    }
    
    // Update account nonce
    bs.setAccountNonce(tx.FromAddress, expectedNonce+1)
    
    // Deduct fee from sender
    bs.deductUSCBalance(tx.FromAddress, tx.Fee)
    
    // Calculate gas used (simplified)
    receipt.GasUsed = bs.calculateGasUsed(tx)
    receipt.GasPrice = fmt.Sprintf("%d", tx.GasPrice)
    receipt.FeePaid = tx.Fee
    
    return receipt, nil
}
```

### **2. USC Transaction Processing**

```go
// USC Transaction Processing Implementation
package transactions

import (
    "context"
    "fmt"
    "time"
    pb "usc-blockchain/proto"
)

// Submit USC Transaction
func (bs *USCBlockchainService) SubmitTransaction(ctx context.Context, req *pb.SubmitTransactionRequest) (*pb.SubmitTransactionResponse, error) {
    tx := req.Transaction
    
    // Generate transaction hash if not provided
    if tx.TxHash == "" {
        tx.TxHash = bs.generateTransactionHash(tx)
    }
    
    // Set timestamp if not provided
    if tx.Timestamp == nil {
        tx.Timestamp = timestamppb.New(time.Now())
    }
    
    // Basic transaction validation
    if err := bs.validateTransactionBasic(tx); err != nil {
        return &pb.SubmitTransactionResponse{
            Success:      false,
            ErrorMessage: fmt.Sprintf("Transaction validation failed: %v", err),
        }, nil
    }
    
    // Check if transaction already exists
    existingTx, _ := bs.GetTransaction(ctx, &pb.GetTransactionRequest{
        TxHash: tx.TxHash,
    })
    if existingTx != nil && existingTx.Success {
        return &pb.SubmitTransactionResponse{
            Success:         true,
            TransactionHash: tx.TxHash,
            Message:         "Transaction already exists",
        }, nil
    }
    
    // Validate transaction signature
    if !bs.verifyTransactionSignature(tx) {
        return &pb.SubmitTransactionResponse{
            Success:      false,
            ErrorMessage: "Invalid transaction signature",
        }, nil
    }
    
    // Check sender balance and nonce
    if err := bs.validateTransactionState(tx); err != nil {
        return &pb.SubmitTransactionResponse{
            Success:      false,
            ErrorMessage: fmt.Sprintf("Transaction state validation failed: %v", err),
        }, nil
    }
    
    // Add to mempool with Cosmos SDK integration
    tx.Status = pb.TransactionStatus_PENDING
    err := bs.addToMempool(tx)
    if err != nil {
        return nil, fmt.Errorf("failed to add transaction to mempool: %v", err)
    }
    
    // Store in Cosmos SDK state
    if bs.blockchainStorage != nil {
        stateKey := fmt.Sprintf("tx:%s", tx.TxHash)
        txData, _ := json.Marshal(tx)
        bs.blockchainStorage.SetState(storage.StateKey(stateKey), txData)
    }
    
    response := &pb.SubmitTransactionResponse{
        Success:         true,
        TransactionHash: tx.TxHash,
        Message:         "Transaction submitted successfully",
    }
    
    // Broadcast to network if requested
    if req.Broadcast {
        go bs.p2pNetwork.BroadcastTransaction(tx)
    }
    
    // Wait for confirmation if requested
    if req.WaitForConfirmation {
        confirmationCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
        defer cancel()
        
        confirmed, blockNumber := bs.waitForTransactionConfirmation(confirmationCtx, tx.TxHash)
        if confirmed {
            response.BlockNumber = blockNumber
            response.EstimatedConfirmationTime = "confirmed"
            
            // Get transaction receipt
            receipt, err := bs.getTransactionReceipt(tx.TxHash)
            if err == nil {
                response.Receipt = receipt
            }
        } else {
            response.EstimatedConfirmationTime = bs.estimateConfirmationTime()
        }
    } else {
        response.EstimatedConfirmationTime = bs.estimateConfirmationTime()
    }
    
    return response, nil
}

// Execute USC Transfer
func (bs *USCBlockchainService) executeUSCTransfer(tx *pb.USCTransaction, receipt *pb.TransactionReceipt) error {
    // Deduct USC from sender
    err := bs.deductUSCBalance(tx.FromAddress, tx.Amount)
    if err != nil {
        return fmt.Errorf("failed to deduct USC from sender: %v", err)
    }
    
    // Add USC to recipient
    err = bs.addUSCBalance(tx.ToAddress, tx.Amount)
    if err != nil {
        // Rollback sender deduction
        bs.addUSCBalance(tx.FromAddress, tx.Amount)
        return fmt.Errorf("failed to add USC to recipient: %v", err)
    }
    
    // Create transfer event log
    transferLog := &pb.EventLog{
        ContractAddress: "usc_token",
        Topics:          []string{"Transfer", tx.FromAddress, tx.ToAddress},
        Data:            []byte(tx.Amount),
        LogIndex:        0,
    }
    receipt.Logs = append(receipt.Logs, transferLog)
    
    // Update address analytics
    go bs.updateAddressAnalytics(tx.FromAddress, "sent", tx.Amount)
    go bs.updateAddressAnalytics(tx.ToAddress, "received", tx.Amount)
    
    return nil
}

// Execute Social Reward Distribution
func (bs *USCBlockchainService) executeSocialReward(tx *pb.USCTransaction, receipt *pb.TransactionReceipt) error {
    // Parse social reward data
    var socialData *pb.SocialRewardData
    if err := tx.TransactionData.UnmarshalTo(socialData); err != nil {
        return fmt.Errorf("failed to parse social reward data: %v", err)
    }
    
    // Validate reward pool has sufficient balance
    rewardPoolBalance, err := bs.getUSCBalance("system_rewards_pool")
    if err != nil {
        return fmt.Errorf("failed to get reward pool balance: %v", err)
    }
    
    if compareUSCAmounts(rewardPoolBalance, tx.Amount) < 0 {
        return fmt.Errorf("insufficient reward pool balance")
    }
    
    // Deduct from reward pool
    err = bs.deductUSCBalance("system_rewards_pool", tx.Amount)
    if err != nil {
        return fmt.Errorf("failed to deduct from reward pool: %v", err)
    }
    
    // Add to recipient
    err = bs.addUSCBalance(tx.ToAddress, tx.Amount)
    if err != nil {
        // Rollback reward pool deduction
        bs.addUSCBalance("system_rewards_pool", tx.Amount)
        return fmt.Errorf("failed to add USC to recipient: %v", err)
    }
    
    // Create reward event log
    rewardLog := &pb.EventLog{
        ContractAddress: "social_rewards",
        Topics:          []string{"SocialReward", socialData.InteractionType, tx.ToAddress},
        Data:            []byte(fmt.Sprintf("%s:%s", tx.Amount, socialData.ContentId)),
        LogIndex:        0,
    }
    receipt.Logs = append(receipt.Logs, rewardLog)
    
    // Update reward analytics
    go bs.updateRewardAnalytics(tx.ToAddress, "social", tx.Amount, socialData.InteractionType)
    
    return nil
}

// Execute Video Reward Distribution
func (bs *USCBlockchainService) executeVideoReward(tx *pb.USCTransaction, receipt *pb.TransactionReceipt) error {
    // Parse video reward data
    var videoData *pb.VideoRewardData
    if err := tx.TransactionData.UnmarshalTo(videoData); err != nil {
        return fmt.Errorf("failed to parse video reward data: %v", err)
    }
    
    // Validate minimum completion percentage for reward
    if videoData.CompletionPercentage < 25 {
        return fmt.Errorf("insufficient completion percentage: %d%% (minimum 25%%)", videoData.CompletionPercentage)
    }
    
    // Calculate reward based on completion percentage
    baseReward := tx.Amount
    if videoData.CompletionPercentage >= 95 {
        // 50% bonus for full completion
        bonusAmount := multiplyUSCAmount(baseReward, "1.5")
        baseReward = bonusAmount
    }
    
    // Quality bonus for HD/4K videos
    if videoData.VideoQuality == "4K" {
        bonusAmount := multiplyUSCAmount(baseReward, "1.2")
        baseReward = bonusAmount
    } else if videoData.VideoQuality == "HD" {
        bonusAmount := multiplyUSCAmount(baseReward, "1.1")
        baseReward = bonusAmount
    }
    
    // Deduct from reward pool
    err := bs.deductUSCBalance("system_rewards_pool", baseReward)
    if err != nil {
        return fmt.Errorf("failed to deduct from reward pool: %v", err)
    }
    
    // Add to recipient
    err = bs.addUSCBalance(tx.ToAddress, baseReward)
    if err != nil {
        // Rollback reward pool deduction
        bs.addUSCBalance("system_rewards_pool", baseReward)
        return fmt.Errorf("failed to add USC to recipient: %v", err)
    }
    
    // Create video reward event log
    videoLog := &pb.EventLog{
        ContractAddress: "video_rewards",
        Topics:          []string{"VideoReward", videoData.VideoId, tx.ToAddress},
        Data:            []byte(fmt.Sprintf("%s:%d%%", baseReward, videoData.CompletionPercentage)),
        LogIndex:        0,
    }
    receipt.Logs = append(receipt.Logs, videoLog)
    
    // Update reward analytics
    go bs.updateRewardAnalytics(tx.ToAddress, "video", baseReward, videoData.VideoId)
    
    return nil
}

// Transfer USC between addresses with Cosmos SDK integration
func (bs *USCBlockchainService) TransferUSC(ctx context.Context, req *pb.TransferUSCRequest) (*pb.TransferUSCResponse, error) {
    // Create USC transaction
    tx := &pb.USCTransaction{
        TxHash:      "",  // Will be generated
        FromAddress: req.FromAddress,
        ToAddress:   req.ToAddress,
        Amount:      req.Amount,
        Fee:         "0",  // Will be calculated
        Nonce:       0,    // Will be set
        GasLimit:    req.GasLimit,
        GasPrice:    req.GasPrice,
        Memo:        req.Memo,
        TxType:      pb.USCTransactionType_TRANSFER,
        Status:      pb.TransactionStatus_PENDING,
        Timestamp:   timestamppb.New(time.Now()),
    }
    
    // Use Cosmos SDK USC module for transfer
    if bs.cosmosApp != nil {
        // Create Cosmos SDK transfer message
        transferMsg := &blockchainproto.MsgTransferUSC{
            FromAddress: req.FromAddress,
            ToAddress:   req.ToAddress,
            Amount:      &types.Coin{Denom: "usc", Amount: req.Amount},
            Memo:        req.Memo,
        }
        
        // Execute transfer via Cosmos SDK
        _, err := bs.cosmosApp.USCKeeper.TransferUSC(ctx, transferMsg)
        if err != nil {
            return nil, fmt.Errorf("Cosmos SDK transfer failed: %v", err)
        }
    }
    
    // Get account nonce
    nonce, err := bs.getAccountNonce(req.FromAddress)
    if err != nil {
        return nil, fmt.Errorf("failed to get account nonce: %v", err)
    }
    tx.Nonce = nonce
    
    // Calculate transaction fee
    fee, err := bs.EstimateTransactionFee(ctx, &pb.EstimateTransactionFeeRequest{
        Transaction: tx,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to estimate transaction fee: %v", err)
    }
    tx.Fee = fee.EstimatedFee
    
    // Generate transaction hash
    tx.TxHash = bs.generateTransactionHash(tx)
    
    // Sign transaction
    signature, err := bs.signTransaction(tx, req.PrivateKey)
    if err != nil {
        return nil, fmt.Errorf("failed to sign transaction: %v", err)
    }
    tx.Signature = signature
    
    // Submit transaction
    submitReq := &pb.SubmitTransactionRequest{
        Transaction:           tx,
        Broadcast:            true,
        WaitForConfirmation:  true,
    }
    
    submitResp, err := bs.SubmitTransaction(ctx, submitReq)
    if err != nil {
        return nil, err
    }
    
    if !submitResp.Success {
        return &pb.TransferUSCResponse{
            Success:      false,
            ErrorMessage: submitResp.ErrorMessage,
        }, nil
    }
    
    return &pb.TransferUSCResponse{
        Success:         true,
        TransactionHash: submitResp.TransactionHash,
        Receipt:         submitResp.Receipt,
    }, nil
}
```

### **3. Proof-of-Stake Consensus Engine**

```go
// Proof-of-Stake Consensus Implementation
package consensus

import (
    "context"
    "fmt"
    "math/big"
    "sort"
    "time"
    pb "usc-blockchain/proto"
)

type PoSConsensusEngine struct {
    blockchain      *USCBlockchainService
    validatorMgr    *ValidatorManager
    stakingMgr      *StakingManager
    epochManager    *EpochManager
    slashingMgr     *SlashingManager
}

// Validator Selection for Block Production
func (pos *PoSConsensusEngine) SelectBlockProducer(blockNumber int64) (*pb.Validator, error) {
    // Get current epoch
    epoch := pos.epochManager.GetEpochForBlock(blockNumber)
    
    // Get active validators for this epoch
    validators, err := pos.validatorMgr.GetActiveValidators(epoch)
    if err != nil {
        return nil, fmt.Errorf("failed to get active validators: %v", err)
    }
    
    if len(validators) == 0 {
        return nil, fmt.Errorf("no active validators found")
    }
    
    // Calculate validator proposer using deterministic algorithm
    // Based on block number and validator voting power
    totalVotingPower := pos.calculateTotalVotingPower(validators)
    seed := pos.generateSeed(blockNumber)
    
    // Weighted random selection based on stake
    selectedValidator := pos.weightedRandomSelection(validators, totalVotingPower, seed)
    
    return selectedValidator, nil
}

// Validate Validator Authority
func (pos *PoSConsensusEngine) ValidateValidatorAuthority(validatorAddress string, blockNumber int64) (bool, error) {
    // Get validator info
    validator, err := pos.validatorMgr.GetValidator(validatorAddress)
    if err != nil {
        return false, fmt.Errorf("validator not found: %v", err)
    }
    
    // Check if validator is active
    if validator.Status != pb.ValidatorStatus_ACTIVE {
        return false, fmt.Errorf("validator is not active: %v", validator.Status)
    }
    
    // Check if validator is jailed
    if validator.Jailed {
        return false, fmt.Errorf("validator is jailed")
    }
    
    // Check minimum stake requirement
    minStake := pos.getMinimumStakeRequirement()
    if compareUSCAmounts(validator.VotingPower, minStake) < 0 {
        return false, fmt.Errorf("validator stake below minimum: %s < %s", validator.VotingPower, minStake)
    }
    
    // Check if it's validator's turn to produce block
    expectedValidator, err := pos.SelectBlockProducer(blockNumber)
    if err != nil {
        return false, err
    }
    
    if expectedValidator.ValidatorAddress != validatorAddress {
        return false, fmt.Errorf("not validator's turn: expected %s, got %s", 
            expectedValidator.ValidatorAddress, validatorAddress)
    }
    
    return true, nil
}

// Process Validator Votes for Block Finalization
func (pos *PoSConsensusEngine) ProcessValidatorVotes(blockNumber int64, votes []*pb.ValidatorVote) (*pb.ConsensusResult, error) {
    // Get active validators for this block
    epoch := pos.epochManager.GetEpochForBlock(blockNumber)
    activeValidators, err := pos.validatorMgr.GetActiveValidators(epoch)
    if err != nil {
        return nil, err
    }
    
    totalVotingPower := pos.calculateTotalVotingPower(activeValidators)
    
    // Count votes by type
    votesCounted := make(map[pb.VoteType]string) // vote type -> total voting power
    votesCounted[pb.VoteType_PROPOSE] = "0"
    votesCounted[pb.VoteType_PREVOTE] = "0"
    votesCounted[pb.VoteType_PRECOMMIT] = "0"
    
    // Validate and count each vote
    for _, vote := range votes {
        // Verify vote signature
        if !pos.verifyVoteSignature(vote, blockNumber) {
            continue
        }
        
        // Get validator
        validator, err := pos.validatorMgr.GetValidator(vote.ValidatorAddress)
        if err != nil {
            continue
        }
        
        // Check validator is active
        if validator.Status != pb.ValidatorStatus_ACTIVE {
            continue
        }
        
        // Add voting power to vote count
        currentPower := votesCounted[vote.VoteType]
        newPower := addUSCAmounts(currentPower, fmt.Sprintf("%d", vote.VotingPower))
        votesCounted[vote.VoteType] = newPower
    }
    
    // Calculate thresholds (2/3 majority required)
    twoThirdsThreshold := multiplyUSCAmount(totalVotingPower, "0.6667") // 2/3
    
    // Check if block reaches consensus
    precommitPower := votesCounted[pb.VoteType_PRECOMMIT]
    blockFinalized := compareUSCAmounts(precommitPower, twoThirdsThreshold) >= 0
    
    result := &pb.ConsensusResult{
        BlockNumber:     blockNumber,
        IsFinalized:     blockFinalized,
        TotalVotingPower: totalVotingPower,
        Votes:           votes,
        VotePowerByType: map[string]string{
            "propose":   votesCounted[pb.VoteType_PROPOSE],
            "prevote":   votesCounted[pb.VoteType_PREVOTE],
            "precommit": votesCounted[pb.VoteType_PRECOMMIT],
        },
        FinalizedAt: timestamppb.New(time.Now()),
    }
    
    if blockFinalized {
        // Update validator performance for voters
        pos.updateValidatorPerformances(votes, true)
        
        // Slash non-voting validators
        pos.handleNonVotingValidators(activeValidators, votes, blockNumber)
    }
    
    return result, nil
}

// Stake USC for Validator
func (pos *PoSConsensusEngine) StakeUSC(ctx context.Context, req *pb.StakeUSCRequest) (*pb.StakeUSCResponse, error) {
    // Validate staking amount
    minStake := pos.getMinimumStakeRequirement()
    if compareUSCAmounts(req.Amount, minStake) < 0 {
        return &pb.StakeUSCResponse{
            Success: false,
            Error:   fmt.Sprintf("Stake amount below minimum: %s < %s", req.Amount, minStake),
        }, nil
    }
    
    // Check user USC balance
    balance, err := pos.blockchain.GetUSCBalance(ctx, &pb.GetUSCBalanceRequest{
        Address: req.DelegatorAddress,
    })
    if err != nil {
        return nil, err
    }
    
    if compareUSCAmounts(balance.Balance, req.Amount) < 0 {
        return &pb.StakeUSCResponse{
            Success: false,
            Error:   "Insufficient USC balance for staking",
        }, nil
    }
    
    // Create staking transaction
    stakingTx := &pb.USCTransaction{
        TxHash:      "",
        FromAddress: req.DelegatorAddress,
        ToAddress:   "staking_contract",
        Amount:      req.Amount,
        TxType:      pb.USCTransactionType_STAKING,
        Status:      pb.TransactionStatus_PENDING,
        Timestamp:   timestamppb.New(time.Now()),
    }
    
    // Add staking data
    stakingData := &pb.StakingData{
        ValidatorAddress:  req.ValidatorAddress,
        DelegatorAddress:  req.DelegatorAddress,
        StakeAmount:      req.Amount,
        StakingType:      "delegation",
        LockPeriod:       req.LockPeriod,
    }
    
    stakingTx.TransactionData, _ = anypb.New(stakingData)
    
    // Submit staking transaction
    submitReq := &pb.SubmitTransactionRequest{
        Transaction:          stakingTx,
        Broadcast:           true,
        WaitForConfirmation: true,
    }
    
    submitResp, err := pos.blockchain.SubmitTransaction(ctx, submitReq)
    if err != nil {
        return nil, err
    }
    
    if !submitResp.Success {
        return &pb.StakeUSCResponse{
            Success: false,
            Error:   submitResp.ErrorMessage,
        }, nil
    }
    
    // Update validator voting power
    err = pos.updateValidatorVotingPower(req.ValidatorAddress, req.Amount, "add")
    if err != nil {
        return nil, fmt.Errorf("failed to update validator voting power: %v", err)
    }
    
    // Update delegator stake record
    err = pos.stakingMgr.AddDelegation(req.DelegatorAddress, req.ValidatorAddress, req.Amount)
    if err != nil {
        return nil, fmt.Errorf("failed to record delegation: %v", err)
    }
    
    // Calculate staking rewards
    expectedRewards := pos.calculateStakingRewards(req.Amount, req.LockPeriod)
    
    return &pb.StakeUSCResponse{
        Success:          true,
        TransactionHash:  submitResp.TransactionHash,
        ExpectedRewards:  expectedRewards,
        UnlockTime:       time.Now().Add(time.Duration(req.LockPeriod) * 24 * time.Hour),
        Message:          fmt.Sprintf("Successfully staked %s USC", req.Amount),
    }, nil
}

// Handle Validator Slashing for Misbehavior
func (pos *PoSConsensusEngine) SlashValidator(validatorAddress string, slashType pb.SlashType, evidence *pb.SlashingEvidence) error {
    validator, err := pos.validatorMgr.GetValidator(validatorAddress)
    if err != nil {
        return fmt.Errorf("validator not found: %v", err)
    }
    
    // Calculate slash percentage based on violation type
    var slashPercentage string
    var jailDuration time.Duration
    
    switch slashType {
    case pb.SlashType_DOUBLE_SIGN:
        slashPercentage = "0.05"  // 5% slash for double signing
        jailDuration = 7 * 24 * time.Hour  // 7 days jail
        
    case pb.SlashType_DOWNTIME:
        slashPercentage = "0.01"  // 1% slash for downtime
        jailDuration = 1 * 24 * time.Hour  // 1 day jail
        
    case pb.SlashType_INVALID_BLOCK:
        slashPercentage = "0.02"  // 2% slash for invalid block
        jailDuration = 3 * 24 * time.Hour  // 3 days jail
        
    default:
        return fmt.Errorf("unknown slash type: %v", slashType)
    }
    
    // Calculate slash amount
    slashAmount := multiplyUSCAmount(validator.VotingPower, slashPercentage)
    
    // Slash validator stake
    newVotingPower := subtractUSCAmounts(validator.VotingPower, slashAmount)
    validator.VotingPower = newVotingPower
    
    // Jail validator
    validator.Jailed = true
    validator.JailUntil = time.Now().Add(jailDuration)
    validator.Status = pb.ValidatorStatus_JAILED
    
    // Update validator record
    err = pos.validatorMgr.UpdateValidator(validator)
    if err != nil {
        return fmt.Errorf("failed to update validator: %v", err)
    }
    
    // Distribute slashed amount to reward pool
    pos.blockchain.addUSCBalance("system_rewards_pool", slashAmount)
    
    // Record slashing event
    slashingEvent := &pb.SlashingEvent{
        ValidatorAddress: validatorAddress,
        SlashType:       slashType,
        SlashAmount:     slashAmount,
        Evidence:        evidence,
        JailUntil:       timestamppb.New(validator.JailUntil),
        Timestamp:       timestamppb.New(time.Now()),
    }
    
    err = pos.slashingMgr.RecordSlashingEvent(slashingEvent)
    if err != nil {
        return fmt.Errorf("failed to record slashing event: %v", err)
    }
    
    // Emit slashing event
    pos.blockchain.emitValidatorEvent(validatorAddress, pb.ValidatorEventType_VALIDATOR_JAILED, slashingEvent)
    
    return nil
}

// Weighted Random Validator Selection
func (pos *PoSConsensusEngine) weightedRandomSelection(validators []*pb.Validator, totalPower string, seed int64) *pb.Validator {
    // Convert total power to big.Int for precision
    totalPowerBig, _ := new(big.Int).SetString(totalPower, 10)
    
    // Generate deterministic random number based on seed
    randomValue := new(big.Int).Mod(big.NewInt(seed), totalPowerBig)
    
    // Select validator based on cumulative stake
    cumulativeStake := big.NewInt(0)
    
    for _, validator := range validators {
        stakeBig, _ := new(big.Int).SetString(validator.VotingPower, 10)
        cumulativeStake.Add(cumulativeStake, stakeBig)
        
        if randomValue.Cmp(cumulativeStake) < 0 {
            return validator
        }
    }
    
    // Fallback to first validator (should never happen)
    return validators[0]
}
```

---

## 📊 **PERFORMANCE METRICS**

### **Service Level Objectives (SLOs)**

```yaml
performance_targets:
  blockchain_operations:
    block_production_time: "<3 seconds average"
    transaction_throughput: "10,000+ TPS sustained"
    block_finality: "<10 seconds (2-3 confirmations)"
    transaction_submission_p95: "<100ms"
    
  consensus:
    validator_selection_time: "<1 second"
    consensus_rounds: "<5 rounds per block"
    voting_participation: ">90% validators"
    
  data_operations:
    balance_query_p95: "<50ms"
    transaction_lookup_p95: "<100ms"
    block_retrieval_p95: "<200ms"
    
  network:
    p2p_message_propagation: "<2 seconds"
    network_sync_time: "<30 minutes for new nodes"
    peer_discovery: "<10 seconds"
    
  availability:
    network_uptime: "99.99%"
    validator_uptime: ">95% average"
    api_availability: "99.99%"
```

---

## 🚀 **DEPLOYMENT & SCALING**

### **Kubernetes Configuration**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: usc-blockchain-core
spec:
  replicas: 21  # Minimum validators for BFT consensus
  selector:
    matchLabels:
      app: usc-blockchain-core
  template:
    spec:
      containers:
      - name: usc-blockchain-core
        image: usc-blockchain-core:latest
        resources:
          requests:
            memory: "8Gi"
            cpu: "4000m"
          limits:
            memory: "16Gi"
            cpu: "8000m"
        env:
        - name: ROCKSDB_PATH
          value: "/data/rocksdb"
        - name: POSTGRES_URL
          valueFrom:
            secretKeyRef:
              name: postgres-secret
              key: blockchain-db-url
        ports:
        - containerPort: 50053
          name: grpc
        - containerPort: 9093
          name: metrics
        - containerPort: 26656
          name: p2p
        - containerPort: 26657
          name: rpc
        volumeMounts:
        - name: blockchain-data
          mountPath: /data
      volumes:
      - name: blockchain-data
        persistentVolumeClaim:
          claimName: usc-blockchain-pvc
```

---

## ✅ **SUCCESS CRITERIA**

### **Functional Requirements**
- ✅ USC blockchain producing blocks consistently
- ✅ PoS consensus với Byzantine fault tolerance  
- ✅ USC transaction processing với all reward types
- ✅ Smart contract execution với WASM runtime
- ✅ Validator management và staking operations
- ✅ Cosmos SDK v0.53.4 integration với 12 custom modules
- ✅ Blockchain state management với RocksDB + PostgreSQL
- ✅ Protocol Buffer integration cho blockchain messages
- ✅ Multi-tier storage architecture (RocksDB + Redis + PostgreSQL)

### **Performance Requirements**
- ✅ 10,000+ TPS sustained throughput
- ✅ <3 seconds average block time
- ✅ <10 seconds block finality
- ✅ <100ms transaction submission latency

### **Security Requirements**
- ✅ Cryptographic security với Ed25519 signatures
- ✅ BFT consensus protecting against <1/3 malicious validators
- ✅ Slashing protection against validator misbehavior
- ✅ Transaction replay protection với nonces

**USC Blockchain Core Service provides robust, high-performance foundation cho entire USC ecosystem với enterprise-grade security và scalability!** ⛓️⚡
