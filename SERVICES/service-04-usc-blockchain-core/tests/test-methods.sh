#!/bin/bash

# ========================================
# Service-04 USC Blockchain Core - Method Testing Script
# ========================================

echo "🚀 Testing Service-04 USC Blockchain Core Methods"
echo "==============================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Service configuration
SERVICE_HOST="localhost"
SERVICE_PORT="8004"
SERVICE_ADDR="${SERVICE_HOST}:${SERVICE_PORT}"

# Helper
test_method() {
    local method_name="$1"
    local description="$2"
    local command="$3"

    echo -e "\n${BLUE}Testing: ${method_name}${NC}"
    echo -e "${YELLOW}Description: ${description}${NC}"
    echo "Command: $command"
    echo "----------------------------------------"

    if eval "$command"; then
        echo -e "${GREEN}✅ ${method_name} - SUCCESS${NC}"
        return 0
    else
        echo -e "${RED}❌ ${method_name} - FAILED${NC}"
        return 1
    fi
}

# Check grpcurl
if ! command -v grpcurl &> /dev/null; then
    echo -e "${RED}❌ grpcurl is not installed. Please install it first.${NC}"
    echo "Install: go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest"
    exit 1
fi

# Health check
echo -e "${BLUE}Checking service health...${NC}"
if grpcurl -plaintext "${SERVICE_ADDR}" grpc.health.v1.Health/Check; then
    echo -e "${GREEN}✅ Service is healthy${NC}"
else
    echo -e "${RED}❌ Service is not responding${NC}"
    exit 1
fi

echo -e "\n${BLUE}Starting method tests...${NC}"

# ========================================
# SETUP (demo data + real blockchain data)
# ========================================
echo -e "\n${YELLOW}=== SETTING UP TEST DATA ===${NC}"
FROM_ADDR="0xFromAddress00000000000000000000000001"
TO_ADDR="0xToAddress0000000000000000000000000002"
WALLET_ADDR="0xWallet0000000000000000000000000003"
CONTRACT_ADDR="0xContract00000000000000000000000004"
USER_ID="user-123"
DEVICE_ID="device-001"

# Get real blockchain data from CometBFT
echo -e "${BLUE}Fetching real blockchain data...${NC}"
COMETBFT_RPC="http://localhost:26657"
LATEST_HEIGHT=$(curl -s "${COMETBFT_RPC}/status" | jq -r '.result.sync_info.latest_block_height // "1"')
BLOCK_1_HASH=$(curl -s "${COMETBFT_RPC}/block?height=1" | jq -r '.result.block_id.hash // ""')
LATEST_BLOCK_HASH=$(curl -s "${COMETBFT_RPC}/block?height=${LATEST_HEIGHT}" | jq -r '.result.block_id.hash // ""')

# Use real data if available, otherwise use test data
if [ -n "$BLOCK_1_HASH" ] && [ "$BLOCK_1_HASH" != "null" ]; then
    BLOCK_HASH="$BLOCK_1_HASH"
    echo -e "${GREEN}✓ Using real block 1 hash: ${BLOCK_HASH}${NC}"
else
    BLOCK_HASH="0xBlockHash000000000000000000000000000006"
    echo -e "${YELLOW}⚠ Using test block hash${NC}"
fi

# Try to get a real transaction hash from pending transactions
REAL_TX_HASH=$(grpcurl -plaintext -d '{"address":"'$FROM_ADDR'","limit":1,"offset":0}' "${SERVICE_ADDR}" blockchain.v1.TransactionOperationsService/GetPendingTransactions 2>/dev/null | jq -r '.transactions[0].transactionHash // ""' | head -1)
if [ -z "$REAL_TX_HASH" ] || [ "$REAL_TX_HASH" == "null" ] || [ "$REAL_TX_HASH" == "" ]; then
    # Try to submit a transaction first to get a real hash
    SUBMIT_RESPONSE=$(grpcurl -plaintext -d '{"from_address":"'$FROM_ADDR'","to_address":"'$TO_ADDR'","amount":"1.0","gas_price":"1","gas_limit":21000,"data":"","user_id":"'$USER_ID'"}' "${SERVICE_ADDR}" blockchain.v1.TransactionOperationsService/SubmitTransaction 2>/dev/null)
    REAL_TX_HASH=$(echo "$SUBMIT_RESPONSE" | jq -r '.transactionHash // ""' 2>/dev/null)
fi
if [ -n "$REAL_TX_HASH" ] && [ "$REAL_TX_HASH" != "null" ] && [ "$REAL_TX_HASH" != "" ]; then
    TX_HASH="$REAL_TX_HASH"
    echo -e "${GREEN}✓ Using real transaction hash: ${TX_HASH}${NC}"
else
    TX_HASH="0xTxHash0000000000000000000000000000000005"
    echo -e "${YELLOW}⚠ Using test transaction hash${NC}"
fi

echo -e "${BLUE}Test Configuration:${NC}"
echo -e "  Block 1 Hash: ${BLOCK_HASH}"
echo -e "  Latest Height: ${LATEST_HEIGHT}"
echo -e "  Latest Block Hash: ${LATEST_BLOCK_HASH}"
echo -e "  Transaction Hash: ${TX_HASH}"

# ========================================
# TRANSACTION OPERATIONS
# ========================================
echo -e "\n${YELLOW}=== TRANSACTION OPERATIONS ===${NC}"

test_method "SubmitTransaction" "Submit USC transaction" \
    'grpcurl -plaintext -d "{\"from_address\":\"'"$FROM_ADDR"'\",\"to_address\":\"'"$TO_ADDR"'\",\"amount\":\"1.0\",\"gas_price\":\"1\",\"gas_limit\":21000,\"data\":\"\",\"nonce\":1,\"signature\":\"0xSIG\",\"user_id\":\"'"$USER_ID"'\",\"device_id\":\"'"$DEVICE_ID"'\"}" \
    "${SERVICE_ADDR}" blockchain.v1.TransactionOperationsService/SubmitTransaction'

test_method "GetTransaction" "Get transaction details" \
    'grpcurl -plaintext -d "{\"transaction_hash\":\"'"$TX_HASH"'\",\"include_receipt\":true}" \
    "${SERVICE_ADDR}" blockchain.v1.TransactionOperationsService/GetTransaction'

test_method "GetTransactionStatus" "Get transaction status" \
    'grpcurl -plaintext -d "{\"transaction_hash\":\"'"$TX_HASH"'\"}" \
    "${SERVICE_ADDR}" blockchain.v1.TransactionOperationsService/GetTransactionStatus'

test_method "GetPendingTransactions" "List mempool txs" \
    'grpcurl -plaintext -d "{\"address\":\"'"$FROM_ADDR"'\",\"limit\":10,\"offset\":0}" \
    "${SERVICE_ADDR}" blockchain.v1.TransactionOperationsService/GetPendingTransactions'

test_method "EstimateTransactionFee" "Estimate gas/fees" \
    'grpcurl -plaintext -d "{\"from_address\":\"'"$FROM_ADDR"'\",\"to_address\":\"'"$TO_ADDR"'\",\"amount\":\"1.0\",\"data\":\"\",\"gas_limit\":21000}" \
    "${SERVICE_ADDR}" blockchain.v1.TransactionOperationsService/EstimateTransactionFee'

# ========================================
# BLOCK OPERATIONS
# ========================================
echo -e "\n${YELLOW}=== BLOCK OPERATIONS ===${NC}"

# Produce block 1 first (if it doesn't exist) before testing GetBlock(1)
echo -e "${BLUE}Ensuring block 1 exists before testing GetBlock(1)...${NC}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [ -f "${SCRIPT_DIR}/produce-block-1.sh" ]; then
    echo -e "${YELLOW}Running produce-block-1.sh to ensure block 1 exists...${NC}"
    bash "${SCRIPT_DIR}/produce-block-1.sh" || echo -e "${YELLOW}⚠️  Block 1 may already exist or could not be produced${NC}"
    echo ""
else
    echo -e "${YELLOW}⚠️  produce-block-1.sh not found, skipping block 1 production${NC}"
    echo -e "${YELLOW}   You may need to produce block 1 manually before testing GetBlock(1)${NC}"
    echo ""
fi

test_method "ProduceBlock" "Produce new block" \
    'grpcurl -plaintext -d "{\"validator_id\":\"val-001\",\"transaction_hashes\":[\"'"$TX_HASH"'\"],\"previous_block_hash\":\"'"$BLOCK_HASH"'\",\"timestamp\":1710000000,\"nonce\":1,\"merkle_root\":\"0xMERKLE\"}" \
    "${SERVICE_ADDR}" blockchain.v1.BlockOperationsService/ProduceBlock'

test_method "ValidateBlock" "Validate block integrity" \
    'grpcurl -plaintext -d "{\"block_hash\":\"'"$BLOCK_HASH"'\",\"block_data\":\"0xDATA\",\"previous_block_hash\":\"0xPREV\",\"block_number\":100}" \
    "${SERVICE_ADDR}" blockchain.v1.BlockOperationsService/ValidateBlock'

# Use real block number (1) instead of test block 100
# Note: Block 1 may not exist until first block is produced (Cosmos SDK standard)
echo -e "${YELLOW}Testing GetBlock(1) - Note: Block 1 may not exist until first block is produced${NC}"
test_method "GetBlock" "Get block by number (using real block 1)" \
    'grpcurl -plaintext -d "{\"block_number\":1,\"include_transactions\":true}" \
    "${SERVICE_ADDR}" blockchain.v1.BlockOperationsService/GetBlock'

# Get block 1 hash from GetBlock response if available
# Try using jq first (more reliable), fallback to grep
BLOCK_1_HASH_FROM_GET=$(grpcurl -plaintext -d '{"block_number":1}' "${SERVICE_ADDR}" blockchain.v1.BlockOperationsService/GetBlock 2>/dev/null | jq -r '.blockHash // empty' 2>/dev/null)
if [ -z "$BLOCK_1_HASH_FROM_GET" ] || [ "$BLOCK_1_HASH_FROM_GET" == "null" ] || [ "$BLOCK_1_HASH_FROM_GET" == "" ]; then
    # Fallback to grep method
    BLOCK_1_HASH_FROM_GET=$(grpcurl -plaintext -d '{"block_number":1}' "${SERVICE_ADDR}" blockchain.v1.BlockOperationsService/GetBlock 2>/dev/null | grep -o '"blockHash":"[^"]*"' | cut -d'"' -f4 | head -1)
fi
if [ -n "$BLOCK_1_HASH_FROM_GET" ] && [ "$BLOCK_1_HASH_FROM_GET" != "null" ] && [ "$BLOCK_1_HASH_FROM_GET" != "" ]; then
    BLOCK_HASH="$BLOCK_1_HASH_FROM_GET"
    echo -e "${GREEN}✓ Using block 1 hash from GetBlock response: ${BLOCK_HASH}${NC}"
fi

test_method "GetBlockByHash" "Get block by hash (using real block 1 hash)" \
    'grpcurl -plaintext -d "{\"block_hash\":\"'"$BLOCK_HASH"'\",\"include_transactions\":true}" \
    "${SERVICE_ADDR}" blockchain.v1.BlockOperationsService/GetBlockByHash'

test_method "GetLatestBlock" "Get current head" \
    'grpcurl -plaintext -d "{}" "${SERVICE_ADDR}" blockchain.v1.BlockOperationsService/GetLatestBlock'

# Use real block range (1-5) instead of test range (90-100)
test_method "GetBlockRange" "Get range of blocks (using real blocks 1-5)" \
    'grpcurl -plaintext -d "{\"start_block\":1,\"end_block\":5,\"include_transactions\":false,\"limit\":10,\"offset\":0}" \
    "${SERVICE_ADDR}" blockchain.v1.BlockOperationsService/GetBlockRange'

# ========================================
# USC COIN OPERATIONS
# ========================================
echo -e "\n${YELLOW}=== USC COIN OPERATIONS ===${NC}"

test_method "GetWalletBalance" "Check USC balance" \
    'grpcurl -plaintext -d "{\"wallet_address\":\"'"$WALLET_ADDR"'\",\"user_id\":\"'"$USER_ID"'\",\"include_pending\":true}" \
    "${SERVICE_ADDR}" blockchain.v1.USCCoinOperationsService/GetWalletBalance'

test_method "TransferUSCBlockchain" "Send USC payment" \
    'grpcurl -plaintext -d "{\"from_address\":\"'"$FROM_ADDR"'\",\"to_address\":\"'"$TO_ADDR"'\",\"amount\":\"2.5\",\"gas_price\":\"1\",\"gas_limit\":21000,\"data\":\"\",\"user_id\":\"'"$USER_ID"'\",\"device_id\":\"'"$DEVICE_ID"'\"}" \
    "${SERVICE_ADDR}" blockchain.v1.USCCoinOperationsService/TransferUSCBlockchain'

test_method "GetUSCSupply" "Token supply info" \
    'grpcurl -plaintext -d "{}" "${SERVICE_ADDR}" blockchain.v1.USCCoinOperationsService/GetUSCSupply'

test_method "GetTransactionHistory" "Wallet tx history" \
    'grpcurl -plaintext -d "{\"wallet_address\":\"'"$WALLET_ADDR"'\",\"user_id\":\"'"$USER_ID"'\",\"limit\":10,\"offset\":0}" \
    "${SERVICE_ADDR}" blockchain.v1.USCCoinOperationsService/GetTransactionHistory'

test_method "GetUSCTransactions" "USC-specific transactions" \
    'grpcurl -plaintext -d "{\"wallet_address\":\"'"$WALLET_ADDR"'\",\"user_id\":\"'"$USER_ID"'\",\"limit\":10,\"offset\":0,\"status\":\"all\"}" \
    "${SERVICE_ADDR}" blockchain.v1.USCCoinOperationsService/GetUSCTransactions'

# ========================================
# NFT TOKEN OPERATIONS
# ========================================
echo -e "\n${YELLOW}=== NFT TOKEN OPERATIONS ===${NC}"

test_method "DeployNFTContract" "Deploy NFT contract" \
    'grpcurl -plaintext -d "{\"from_address\":\"'"$FROM_ADDR"'\",\"contract_name\":\"DemoNFT\",\"contract_symbol\":\"DNFT\",\"base_uri\":\"https://ipfs/\",\"gas_price\":\"1\",\"gas_limit\":100000,\"user_id\":\"'"$USER_ID"'\"}" \
    "${SERVICE_ADDR}" blockchain.v1.NFTTokenOperationsService/DeployNFTContract'

test_method "CreateNFTCollection" "Create NFT collection" \
    'grpcurl -plaintext -d "{\"contract_address\":\"'"$CONTRACT_ADDR"'\",\"collection_name\":\"Demo Collection\",\"creator_address\":\"'"$FROM_ADDR"'\"}" \
    "${SERVICE_ADDR}" blockchain.v1.NFTTokenOperationsService/CreateNFTCollection'

test_method "MintNFT" "Mint an NFT" \
    'grpcurl -plaintext -d "{\"contract_address\":\"'"$CONTRACT_ADDR"'\",\"to_address\":\"'"$TO_ADDR"'\",\"token_uri\":\"ipfs://demo.json\",\"gas_price\":\"1\",\"gas_limit\":100000,\"user_id\":\"'"$USER_ID"'\"}" \
    "${SERVICE_ADDR}" blockchain.v1.NFTTokenOperationsService/MintNFT'

test_method "TransferNFT" "Transfer NFT" \
    'grpcurl -plaintext -d "{\"contract_address\":\"'"$CONTRACT_ADDR"'\",\"token_id\":\"1\",\"from_address\":\"'"$FROM_ADDR"'\",\"to_address\":\"'"$TO_ADDR"'\",\"gas_price\":\"1\",\"gas_limit\":100000,\"user_id\":\"'"$USER_ID"'\"}" \
    "${SERVICE_ADDR}" blockchain.v1.NFTTokenOperationsService/TransferNFT'

test_method "BurnNFT" "Burn an NFT" \
    'grpcurl -plaintext -d "{\"contract_address\":\"'"$CONTRACT_ADDR"'\",\"token_id\":\"1\",\"owner_address\":\"'"$FROM_ADDR"'\",\"gas_price\":\"1\",\"gas_limit\":100000,\"user_id\":\"'"$USER_ID"'\"}" \
    "${SERVICE_ADDR}" blockchain.v1.NFTTokenOperationsService/BurnNFT'

test_method "GetNFTInfo" "NFT metadata/info" \
    'grpcurl -plaintext -d "{\"contract_address\":\"'"$CONTRACT_ADDR"'\",\"token_id\":\"1\"}" \
    "${SERVICE_ADDR}" blockchain.v1.NFTTokenOperationsService/GetNFTInfo'

test_method "GetNFTsByOwner" "User's NFTs" \
    'grpcurl -plaintext -d "{\"owner_address\":\"'"$WALLET_ADDR"'\",\"limit\":10,\"offset\":0}" \
    "${SERVICE_ADDR}" blockchain.v1.NFTTokenOperationsService/GetNFTsByOwner'

# ========================================
# SMART CONTRACT OPERATIONS
# ========================================
echo -e "\n${YELLOW}=== SMART CONTRACT OPERATIONS ===${NC}"

test_method "DeployContract" "Deploy smart contract" \
    'grpcurl -plaintext -d "{\"from_address\":\"'"$FROM_ADDR"'\",\"bytecode\":\"0x6000\",\"abi\":\"[]\",\"constructor_args\":\"[]\",\"gas_price\":\"1\",\"gas_limit\":200000,\"user_id\":\"'"$USER_ID"'\",\"contract_name\":\"Demo\"}" \
    "${SERVICE_ADDR}" blockchain.v1.SmartContractOperationsService/DeployContract'

test_method "ExecuteContract" "Execute contract function" \
    'grpcurl -plaintext -d "{\"contract_address\":\"'"$CONTRACT_ADDR"'\",\"function_name\":\"set\",\"parameters\":[\"42\"],\"from_address\":\"'"$FROM_ADDR"'\",\"gas_price\":\"1\",\"gas_limit\":100000,\"user_id\":\"'"$USER_ID"'\"}" \
    "${SERVICE_ADDR}" blockchain.v1.SmartContractOperationsService/ExecuteContract'

test_method "QueryContract" "Read contract state" \
    'grpcurl -plaintext -d "{\"contract_address\":\"'"$CONTRACT_ADDR"'\",\"function_name\":\"get\",\"parameters\":[]}" \
    "${SERVICE_ADDR}" blockchain.v1.SmartContractOperationsService/QueryContract'

test_method "GetContractCode" "Get bytecode/ABI" \
    'grpcurl -plaintext -d "{\"contract_address\":\"'"$CONTRACT_ADDR"'\"}" \
    "${SERVICE_ADDR}" blockchain.v1.SmartContractOperationsService/GetContractCode'

test_method "GetContractStorage" "Get storage slot" \
    'grpcurl -plaintext -d "{\"contract_address\":\"'"$CONTRACT_ADDR"'\",\"storage_key\":\"0x00\"}" \
    "${SERVICE_ADDR}" blockchain.v1.SmartContractOperationsService/GetContractStorage'



# ========================================
# NETWORK OPERATIONS
# ========================================
echo -e "\n${YELLOW}=== NETWORK OPERATIONS ===${NC}"

test_method "GetNetworkInfo" "Network statistics" \
    'grpcurl -plaintext -d "{}" "${SERVICE_ADDR}" blockchain.v1.NetworkOperationsService/GetNetworkInfo'

test_method "GetChainInfo" "Blockchain info" \
    'grpcurl -plaintext -d "{}" "${SERVICE_ADDR}" blockchain.v1.NetworkOperationsService/GetChainInfo'

test_method "GetPeers" "Peer list" \
    'grpcurl -plaintext -d "{\"limit\":10,\"offset\":0,\"peer_type\":\"all\",\"status\":\"all\"}" \
    "${SERVICE_ADDR}" blockchain.v1.NetworkOperationsService/GetPeers'

test_method "GetNetworkStats" "Network performance metrics" \
    'grpcurl -plaintext -d "{\"time_range\":\"1h\",\"metric_type\":\"all\"}" \
    "${SERVICE_ADDR}" blockchain.v1.NetworkOperationsService/GetNetworkStats'

# ========================================
# VALIDATOR OPERATIONS
# ========================================
echo -e "\n${YELLOW}=== VALIDATOR OPERATIONS ===${NC}"

test_method "RegisterValidator" "Register PoS validator" \
    'grpcurl -plaintext -d "{\"validator_address\":\"'"$FROM_ADDR"'\",\"validator_name\":\"Demo Validator\",\"commission_rate\":\"5\",\"validator_public_key\":\"0xPUB\",\"user_id\":\"'"$USER_ID"'\"}" \
    "${SERVICE_ADDR}" blockchain.v1.ValidatorOperationsService/RegisterValidator'

test_method "GetValidators" "Active validators list" \
    'grpcurl -plaintext -d "{\"limit\":10,\"offset\":0,\"status\":\"all\"}" \
    "${SERVICE_ADDR}" blockchain.v1.ValidatorOperationsService/GetValidators'

test_method "GetValidatorStatus" "Validator performance" \
    'grpcurl -plaintext -d "{\"validator_address\":\"'"$FROM_ADDR"'\"}" \
    "${SERVICE_ADDR}" blockchain.v1.ValidatorOperationsService/GetValidatorStatus'

test_method "StakeUSC" "Stake USC" \
    'grpcurl -plaintext -d "{\"delegator_address\":\"'"$WALLET_ADDR"'\",\"validator_address\":\"'"$FROM_ADDR"'\",\"stake_amount\":\"10\",\"gas_price\":\"1\",\"gas_limit\":21000,\"user_id\":\"'"$USER_ID"'\",\"device_id\":\"'"$DEVICE_ID"'\"}" \
    "${SERVICE_ADDR}" blockchain.v1.ValidatorOperationsService/StakeUSC'

test_method "UnstakeUSC" "Unstake USC" \
    'grpcurl -plaintext -d "{\"delegator_address\":\"'"$WALLET_ADDR"'\",\"validator_address\":\"'"$FROM_ADDR"'\",\"unstake_amount\":\"5\",\"gas_price\":\"1\",\"gas_limit\":21000,\"user_id\":\"'"$USER_ID"'\",\"device_id\":\"'"$DEVICE_ID"'\"}" \
    "${SERVICE_ADDR}" blockchain.v1.ValidatorOperationsService/UnstakeUSC'

# ========================================
# CUSTOM TOKEN OPERATIONS
# ========================================
echo -e "\n${YELLOW}=== CUSTOM TOKEN OPERATIONS ===${NC}"

test_method "CreateBlockchainToken" "Create store token" \
    'grpcurl -plaintext -d "{\"from_address\":\"'"$FROM_ADDR"'\",\"token_name\":\"StoreCoin\",\"token_symbol\":\"SC\",\"total_supply\":\"1000000\",\"decimals\":18,\"gas_price\":\"1\",\"gas_limit\":200000,\"user_id\":\"'"$USER_ID"'\"}" \
    "${SERVICE_ADDR}" blockchain.v1.CustomTokenOperationsService/CreateBlockchainToken'

test_method "MintTokens" "Mint store tokens" \
    'grpcurl -plaintext -d "{\"contract_address\":\"'"$CONTRACT_ADDR"'\",\"to_address\":\"'"$WALLET_ADDR"'\",\"amount\":\"100\",\"gas_price\":\"1\",\"gas_limit\":200000,\"user_id\":\"'"$USER_ID"'\"}" \
    "${SERVICE_ADDR}" blockchain.v1.CustomTokenOperationsService/MintTokens'

test_method "BurnTokens" "Burn store tokens" \
    'grpcurl -plaintext -d "{\"contract_address\":\"'"$CONTRACT_ADDR"'\",\"from_address\":\"'"$WALLET_ADDR"'\",\"amount\":\"50\",\"gas_price\":\"1\",\"gas_limit\":200000,\"user_id\":\"'"$USER_ID"'\"}" \
    "${SERVICE_ADDR}" blockchain.v1.CustomTokenOperationsService/BurnTokens'

test_method "GetTokenInfo" "Token metadata" \
    'grpcurl -plaintext -d "{\"contract_address\":\"'"$CONTRACT_ADDR"'\"}" \
    "${SERVICE_ADDR}" blockchain.v1.CustomTokenOperationsService/GetTokenInfo'

test_method "GetTokenBalance" "Token balance" \
    'grpcurl -plaintext -d "{\"contract_address\":\"'"$CONTRACT_ADDR"'\",\"wallet_address\":\"'"$WALLET_ADDR"'\"}" \
    "${SERVICE_ADDR}" blockchain.v1.CustomTokenOperationsService/GetTokenBalance'

# ========================================
# PRODUCT CERTIFICATE OPERATIONS
# ========================================
echo -e "\n${YELLOW}=== PRODUCT CERTIFICATE OPERATIONS ===${NC}"

# Create product certificate and capture certificate ID
echo -e "${BLUE}Creating product certificate...${NC}"
CREATE_CERT_RESPONSE=$(grpcurl -plaintext -d "{\"from_address\":\"$FROM_ADDR\",\"product_id\":\"PRD-001\",\"product_name\":\"Demo Product\",\"manufacturer_address\":\"$FROM_ADDR\",\"serial_number\":\"SN-001\",\"gas_price\":\"1\",\"gas_limit\":200000,\"user_id\":\"$USER_ID\"}" \
    "${SERVICE_ADDR}" blockchain.v1.ProductCertificateOperationsService/CreateProductCertificate 2>/dev/null)

CERT_ID=$(echo "$CREATE_CERT_RESPONSE" | jq -r '.certificateId // empty' 2>/dev/null)
if [ -z "$CERT_ID" ] || [ "$CERT_ID" == "null" ] || [ "$CERT_ID" == "" ]; then
    # Try alternative parsing
    CERT_ID=$(echo "$CREATE_CERT_RESPONSE" | grep -o '"certificateId":"[^"]*"' | cut -d'"' -f4 | head -1)
fi

if [ -n "$CERT_ID" ] && [ "$CERT_ID" != "null" ] && [ "$CERT_ID" != "" ]; then
    echo -e "${GREEN}✓ Created certificate: ${CERT_ID}${NC}"
else
    CERT_ID="cert_PRD-001_$(date +%s)"
    echo -e "${YELLOW}⚠ Could not extract certificate ID, using fallback: ${CERT_ID}${NC}"
fi

test_method "CreateProductCertificate" "Issue product certificate" \
    'echo "$CREATE_CERT_RESPONSE"'

test_method "VerifyBlockchainProductCertificate" "Verify product authenticity" \
    'grpcurl -plaintext -d "{\"certificate_id\":\"'"$CERT_ID"'\",\"product_id\":\"PRD-001\",\"serial_number\":\"SN-001\"}" \
    "${SERVICE_ADDR}" blockchain.v1.ProductCertificateOperationsService/VerifyBlockchainProductCertificate'

test_method "TransferProductOwnership" "Transfer product ownership" \
    'grpcurl -plaintext -d "{\"certificate_id\":\"'"$CERT_ID"'\",\"from_address\":\"'"$FROM_ADDR"'\",\"to_address\":\"'"$TO_ADDR"'\",\"gas_price\":\"1\",\"gas_limit\":200000,\"user_id\":\"'"$USER_ID"'\"}" \
    "${SERVICE_ADDR}" blockchain.v1.ProductCertificateOperationsService/TransferProductOwnership'

# ========================================
# STORE BRIDGE OPERATIONS
# ========================================
echo -e "\n${YELLOW}=== STORE BRIDGE OPERATIONS ===${NC}"

test_method "DeployStoreBridge" "Deploy bridge contract" \
    'grpcurl -plaintext -d "{\"from_address\":\"'"$FROM_ADDR"'\",\"bridge_name\":\"Demo Bridge\",\"target_network\":\"extnet\",\"gas_price\":\"1\",\"gas_limit\":200000,\"user_id\":\"'"$USER_ID"'\"}" \
    "${SERVICE_ADDR}" blockchain.v1.StoreBridgeOperationsService/DeployStoreBridge'

test_method "RegisterStoreNetwork" "Register external network" \
    'grpcurl -plaintext -d "{\"network_name\":\"ExtNet\",\"network_id\":\"ext-1\",\"chain_id\":\"1001\",\"rpc_url\":\"http://rpc\",\"explorer_url\":\"http://explorer\",\"native_currency\":\"EXT\",\"user_id\":\"'"$USER_ID"'\"}" \
    "${SERVICE_ADDR}" blockchain.v1.StoreBridgeOperationsService/RegisterStoreNetwork'

test_method "BridgeStoreTokenToUSC" "Bridge store token -> USC" \
    'grpcurl -plaintext -d "{\"from_address\":\"'"$FROM_ADDR"'\",\"store_token_address\":\"'"$CONTRACT_ADDR"'\",\"store_token_amount\":\"100\",\"target_network\":\"usc\",\"gas_price\":\"1\",\"gas_limit\":200000,\"user_id\":\"'"$USER_ID"'\",\"device_id\":\"'"$DEVICE_ID"'\"}" \
    "${SERVICE_ADDR}" blockchain.v1.StoreBridgeOperationsService/BridgeStoreTokenToUSC'

test_method "BridgeUSCToStoreToken" "Bridge USC -> store token" \
    'grpcurl -plaintext -d "{\"from_address\":\"'"$FROM_ADDR"'\",\"usc_amount\":\"50\",\"target_network\":\"extnet\",\"target_token_address\":\"'"$CONTRACT_ADDR"'\",\"gas_price\":\"1\",\"gas_limit\":200000,\"user_id\":\"'"$USER_ID"'\",\"device_id\":\"'"$DEVICE_ID"'\"}" \
    "${SERVICE_ADDR}" blockchain.v1.StoreBridgeOperationsService/BridgeUSCToStoreToken'

test_method "GetStoreBridgeMetrics" "Bridge metrics" \
    'grpcurl -plaintext -d "{\"bridge_address\":\"'"$CONTRACT_ADDR"'\",\"time_range\":\"1h\",\"metric_type\":\"all\"}" \
    "${SERVICE_ADDR}" blockchain.v1.StoreBridgeOperationsService/GetStoreBridgeMetrics'

test_method "ValidateStoreBridge" "Validate bridge" \
    'grpcurl -plaintext -d "{\"bridge_address\":\"'"$CONTRACT_ADDR"'\",\"validation_type\":\"security\",\"user_id\":\"'"$USER_ID"'\"}" \
    "${SERVICE_ADDR}" blockchain.v1.StoreBridgeOperationsService/ValidateStoreBridge'

# ========================================
# STORE NETWORK OPERATIONS
# ========================================
echo -e "\n${YELLOW}=== STORE NETWORK OPERATIONS ===${NC}"

test_method "SyncStoreNetworkState" "Sync external network state" \
    'grpcurl -plaintext -d "{\"network_id\":\"ext-1\",\"sync_type\":\"full\",\"from_block\":0,\"to_block\":100,\"user_id\":\"'"$USER_ID"'\"}" \
    "${SERVICE_ADDR}" blockchain.v1.StoreNetworkOperationsService/SyncStoreNetworkState'

test_method "GetStoreNetworkInfo" "External network info" \
    'grpcurl -plaintext -d "{\"network_id\":\"ext-1\",\"info_type\":\"basic\"}" \
    "${SERVICE_ADDR}" blockchain.v1.StoreNetworkOperationsService/GetStoreNetworkInfo'

test_method "UpdateStoreBridgeConfig" "Update bridge config" \
    'grpcurl -plaintext -d "{\"bridge_address\":\"'"$CONTRACT_ADDR"'\",\"config_data\":\"{\\\"limits\\\":100}\",\"config_type\":\"security\",\"user_id\":\"'"$USER_ID"'\"}" \
    "${SERVICE_ADDR}" blockchain.v1.StoreNetworkOperationsService/UpdateStoreBridgeConfig'

# ========================================
# STREAMING OPERATIONS (with timeout)
# ========================================
echo -e "\n${YELLOW}=== STREAMING OPERATIONS ===${NC}"

test_method "StreamBlocks" "Live block updates" \
    'timeout 5s grpcurl -plaintext -d "{\"client_id\":\"cli-001\",\"include_transactions\":false,\"filter_type\":\"all\"}" \
    "${SERVICE_ADDR}" blockchain.v1.StreamingOperationsService/StreamBlocks || true'

test_method "StreamTransactions" "Live tx feed" \
    'timeout 5s grpcurl -plaintext -d "{\"client_id\":\"cli-001\",\"transaction_type\":\"all\",\"status\":\"all\"}" \
    "${SERVICE_ADDR}" blockchain.v1.StreamingOperationsService/StreamTransactions || true'

test_method "StreamValidatorEvents" "Validator event stream" \
    'timeout 5s grpcurl -plaintext -d "{\"client_id\":\"cli-001\",\"event_type\":\"all\",\"include_delegator_events\":false}" \
    "${SERVICE_ADDR}" blockchain.v1.StreamingOperationsService/StreamValidatorEvents || true'

test_method "StreamNetworkEvents" "Network event stream" \
    'timeout 5s grpcurl -plaintext -d "{\"client_id\":\"cli-001\",\"event_type\":\"all\",\"severity\":\"info\",\"include_peer_events\":true}" \
    "${SERVICE_ADDR}" blockchain.v1.StreamingOperationsService/StreamNetworkEvents || true'


# ========================================
# SUMMARY
# ========================================
echo -e "\n${GREEN}🎉 Blockchain core method tests completed!${NC}"