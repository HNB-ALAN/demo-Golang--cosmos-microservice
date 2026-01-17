#!/bin/bash
# Script to reload Prometheus configuration
# Adds alert rules and reloads Prometheus

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
MONITORING_DIR="${PROJECT_ROOT}/SERVICES/service-08-monitoring/monitoring"

# Configuration
PROMETHEUS_URL="${PROMETHEUS_URL:-http://localhost:9090}"
PROMETHEUS_RULES_DIR="${PROMETHEUS_RULES_DIR:-/etc/prometheus/rules}"
PROMETHEUS_ALERTS_DIR="${PROMETHEUS_ALERTS_DIR:-/etc/prometheus/alerts}"

echo -e "${GREEN}🔄 Reloading Prometheus Configuration...${NC}"
echo ""

# Check if Prometheus is accessible
if ! curl -s -f "${PROMETHEUS_URL}/-/healthy" > /dev/null 2>&1; then
    echo -e "${RED}❌ Cannot connect to Prometheus at ${PROMETHEUS_URL}${NC}"
    echo -e "${YELLOW}Please ensure Prometheus is running and accessible${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Prometheus is accessible at ${PROMETHEUS_URL}${NC}"
echo ""

# Method selection
METHOD="${METHOD:-reload}"

if [ "$METHOD" = "reload" ]; then
    echo -e "${YELLOW}Reloading Prometheus configuration...${NC}"
    
    # Reload Prometheus
    RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${PROMETHEUS_URL}/-/reload")
    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    
    if [ "$HTTP_CODE" = "200" ] || [ "$HTTP_CODE" = "204" ]; then
        echo -e "${GREEN}✓ Prometheus reloaded successfully${NC}"
    else
        echo -e "${RED}✗ Failed to reload Prometheus (HTTP ${HTTP_CODE})${NC}"
        exit 1
    fi
    
elif [ "$METHOD" = "copy" ]; then
    echo -e "${YELLOW}Copying rule files...${NC}"
    echo ""
    echo -e "${BLUE}Note: This requires access to Prometheus server${NC}"
    echo -e "${BLUE}Copy commands:${NC}"
    echo ""
    echo "# Copy centralized rules"
    echo "cp ${MONITORING_DIR}/rules/*.yml ${PROMETHEUS_RULES_DIR}/"
    echo ""
    echo "# Copy service-specific alerts"
    echo "cp ${MONITORING_DIR}/alerts/*.yml ${PROMETHEUS_ALERTS_DIR}/"
    echo ""
    echo "# Then reload Prometheus"
    echo "curl -X POST ${PROMETHEUS_URL}/-/reload"
    echo ""
    
elif [ "$METHOD" = "verify" ]; then
    echo -e "${YELLOW}Verifying loaded rules...${NC}"
    echo ""
    
    # Get all rules
    RULES_RESPONSE=$(curl -s "${PROMETHEUS_URL}/api/v1/rules")
    
    if echo "$RULES_RESPONSE" | grep -q "groups"; then
        echo -e "${GREEN}✓ Rules are loaded${NC}"
        echo ""
        echo "Rule groups found:"
        echo "$RULES_RESPONSE" | grep -o '"name":"[^"]*"' | sort -u | head -10
    else
        echo -e "${YELLOW}⚠ No rules found or Prometheus not responding correctly${NC}"
    fi
    
    # Count rule groups
    GROUP_COUNT=$(echo "$RULES_RESPONSE" | grep -o '"name":"[^"]*"' | wc -l)
    echo ""
    echo "Total rule groups: ${GROUP_COUNT}"
    
else
    echo -e "${RED}Invalid method: ${METHOD}${NC}"
    echo "Available methods: reload, copy, verify"
    exit 1
fi

echo ""
echo "=========================================="
echo "Next Steps"
echo "=========================================="
echo ""
echo "1. Verify rules are loaded:"
echo "   curl ${PROMETHEUS_URL}/api/v1/rules"
echo ""
echo "2. Check Prometheus targets:"
echo "   curl ${PROMETHEUS_URL}/api/v1/targets"
echo ""
echo "3. View rules in Prometheus UI:"
echo "   ${PROMETHEUS_URL}/rules"
echo ""
echo "4. Test alerts:"
echo "   ${PROMETHEUS_URL}/alerts"
echo ""

