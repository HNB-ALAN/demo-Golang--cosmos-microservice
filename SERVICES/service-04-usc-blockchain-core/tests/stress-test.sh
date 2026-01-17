#!/bin/bash

# ========================================
# Service-04 USC Blockchain Core - Stress Test Script
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

# Stress test parameters
MAX_CONCURRENT=${MAX_CONCURRENT:-1000}
DURATION=${DURATION:-10m}
RAMP_UP=${RAMP_UP:-5m}

echo -e "${BLUE}=== 💥 STRESS TEST - SERVICE-04 ===${NC}"
echo ""
echo -e "${YELLOW}Configuration:${NC}"
echo -e "  Target: ${SERVICE_ADDR}"
echo -e "  Max Concurrent Users: ${MAX_CONCURRENT}"
echo -e "  Duration: ${DURATION}"
echo -e "  Ramp-up: ${RAMP_UP}"
echo ""
echo -e "${YELLOW}⚠️  WARNING: This test will push the service to its limits${NC}"
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
echo -e "${BLUE}Starting stress test...${NC}"
echo ""

# Run stress test
ghz --insecure \
  --proto ./proto/blockchain.proto \
  --call blockchain.v1.BlockOperationsService.GetBlock \
  -d '{"block_number":1}' \
  -n 6000000 \
  -c ${MAX_CONCURRENT} \
  -t ${DURATION} \
  --timeout 10s \
  --format json \
  ${SERVICE_ADDR} > stress-test-results.json

echo ""
echo -e "${GREEN}✅ Stress test completed!${NC}"
echo -e "${BLUE}Results saved to: stress-test-results.json${NC}"

# Parse and display summary
if command -v jq &> /dev/null; then
    echo ""
    echo -e "${YELLOW}=== TEST SUMMARY ===${NC}"
    echo ""
    
    TOTAL=$(jq -r '.total' stress-test-results.json)
    RPS=$(jq -r '.rps' stress-test-results.json)
    LATENCY_P95=$(jq -r '.latency.p95' stress-test-results.json)
    ERRORS=$(jq -r '.statusCodeDist."0" // 0' stress-test-results.json)
    
    echo -e "Total Requests: ${TOTAL}"
    echo -e "Peak RPS: ${RPS}"
    echo -e "p95 Latency: ${LATENCY_P95}"
    echo -e "Errors: ${ERRORS}"
    echo ""
    echo -e "${YELLOW}Note: Stress test identifies breaking point${NC}"
    echo -e "${YELLOW}      Some degradation is expected at peak load${NC}"
fi

echo ""
echo -e "${GREEN}🎉 Stress test complete!${NC}"


