#!/bin/bash

# ========================================
# Service-04 USC Blockchain Core - Benchmark Script
# ========================================

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SERVICE_HOST="localhost"
SERVICE_PORT="8004"
SERVICE_ADDR="${SERVICE_HOST}:${SERVICE_PORT}"

# Benchmark parameters
REQUESTS_PER_METHOD=10000
CONCURRENT=10

echo -e "${BLUE}=== 🔬 BENCHMARK TEST - SERVICE-04 ===${NC}"
echo ""
echo -e "${YELLOW}Configuration:${NC}"
echo -e "  Target: ${SERVICE_ADDR}"
echo -e "  Requests per method: ${REQUESTS_PER_METHOD}"
echo -e "  Concurrent: ${CONCURRENT}"
echo ""

# Check if ghz is installed
if ! command -v ghz &> /dev/null; then
    echo -e "${RED}❌ ghz is not installed.${NC}"
    echo "Install: go install github.com/bojand/ghz/cmd/ghz@latest"
    exit 1
fi

# Check if service is running
echo -e "${BLUE}Checking service health...${NC}"
if grpcurl -plaintext "${SERVICE_ADDR}" grpc.health.v1.Health/Check > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Service is healthy${NC}"
else
    echo -e "${RED}❌ Service is not responding${NC}"
    exit 1
fi

echo ""
echo -e "${BLUE}Starting benchmark tests...${NC}"
echo ""

# Methods to benchmark
declare -A METHODS=(
    ["GetBlock"]='{"block_number":1}'
    ["GetTransaction"]='{"transaction_hash":"0xTxHash0000000000000000000000000000000005"}'
    ["GetWalletBalance"]='{"wallet_address":"0xWallet0000000000000000000000000003"}'
    ["GetLatestBlock"]='{}'
    ["GetNetworkInfo"]='{}'
)

RESULTS_FILE="benchmark-results.json"
echo "[]" > $RESULTS_FILE

for method_name in "${!METHODS[@]}"; do
    data="${METHODS[$method_name]}"
    
    echo -e "${YELLOW}Benchmarking: ${method_name}${NC}"
    
    # Determine service and method
    case $method_name in
        "GetBlock"|"GetLatestBlock")
            service="BlockOperationsService"
            ;;
        "GetTransaction")
            service="TransactionOperationsService"
            ;;
        "GetWalletBalance")
            service="USCCoinOperationsService"
            ;;
        "GetNetworkInfo")
            service="NetworkOperationsService"
            ;;
        *)
            echo -e "${RED}Unknown method: ${method_name}${NC}"
            continue
            ;;
    esac
    
    # Run benchmark
    ghz --insecure \
      --proto ./proto/blockchain.proto \
      --call "blockchain.v1.${service}.${method_name}" \
      -d "$data" \
      -n ${REQUESTS_PER_METHOD} \
      -c ${CONCURRENT} \
      --timeout 5s \
      --format json \
      ${SERVICE_ADDR} > "benchmark-${method_name}.json"
    
    echo -e "${GREEN}✅ ${method_name} benchmark completed${NC}"
    echo ""
done

echo -e "${GREEN}✅ All benchmarks completed!${NC}"
echo ""
echo -e "${BLUE}Results:${NC}"
for method_name in "${!METHODS[@]}"; do
    if [ -f "benchmark-${method_name}.json" ]; then
        if command -v jq &> /dev/null; then
            LATENCY_P95=$(jq -r '.latency.p95' "benchmark-${method_name}.json" 2>/dev/null || echo "N/A")
            RPS=$(jq -r '.rps' "benchmark-${method_name}.json" 2>/dev/null || echo "N/A")
            echo -e "  ${method_name}: p95=${LATENCY_P95}, RPS=${RPS}"
        else
            echo -e "  ${method_name}: benchmark-${method_name}.json"
        fi
    fi
done

echo ""
echo -e "${GREEN}🎉 Benchmark complete!${NC}"


