#!/bin/bash
# Script to test monitoring setup
# Tests metrics collection, alert loading, and dashboard availability

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

# Configuration
PROMETHEUS_URL="${PROMETHEUS_URL:-http://localhost:9090}"
GRAFANA_URL="${GRAFANA_URL:-http://localhost:3000}"
ALERTMANAGER_URL="${ALERTMANAGER_URL:-http://localhost:9093}"

echo -e "${GREEN}🧪 Testing Monitoring Setup...${NC}"
echo ""

# Test 1: Prometheus Health
echo "=========================================="
echo "1. Testing Prometheus"
echo "=========================================="
echo ""

if curl -s -f "${PROMETHEUS_URL}/-/healthy" > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Prometheus is healthy${NC}"
else
    echo -e "${RED}✗ Prometheus is not accessible${NC}"
    exit 1
fi

# Test 2: Prometheus Rules
echo -e "${YELLOW}Checking loaded rules...${NC}"
RULES_RESPONSE=$(curl -s "${PROMETHEUS_URL}/api/v1/rules" 2>/dev/null || echo "")

if [ -n "$RULES_RESPONSE" ] && echo "$RULES_RESPONSE" | grep -q "groups"; then
    GROUP_COUNT=$(echo "$RULES_RESPONSE" | grep -o '"name":"[^"]*"' | wc -l)
    echo -e "${GREEN}✓ Rules loaded: ${GROUP_COUNT} rule groups${NC}"
else
    echo -e "${YELLOW}⚠ No rules found or cannot access${NC}"
fi

# Test 3: Prometheus Targets
echo -e "${YELLOW}Checking targets...${NC}"
TARGETS_RESPONSE=$(curl -s "${PROMETHEUS_URL}/api/v1/targets" 2>/dev/null || echo "")

if [ -n "$TARGETS_RESPONSE" ] && echo "$TARGETS_RESPONSE" | grep -q "activeTargets"; then
    ACTIVE_COUNT=$(echo "$TARGETS_RESPONSE" | grep -o '"health":"up"' | wc -l)
    echo -e "${GREEN}✓ Active targets: ${ACTIVE_COUNT}${NC}"
else
    echo -e "${YELLOW}⚠ Cannot check targets${NC}"
fi

echo ""

# Test 4: Grafana Health
echo "=========================================="
echo "2. Testing Grafana"
echo "=========================================="
echo ""

if curl -s -f "${GRAFANA_URL}/api/health" > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Grafana is healthy${NC}"
    
    # Check dashboards
    echo -e "${YELLOW}Checking dashboards...${NC}"
    DASHBOARDS=$(curl -s -u admin:admin "${GRAFANA_URL}/api/search?type=dash-db" 2>/dev/null || echo "")
    
    if [ -n "$DASHBOARDS" ] && echo "$DASHBOARDS" | grep -q "title"; then
        DASHBOARD_COUNT=$(echo "$DASHBOARDS" | grep -o '"title"' | wc -l)
        echo -e "${GREEN}✓ Dashboards found: ${DASHBOARD_COUNT}${NC}"
    else
        echo -e "${YELLOW}⚠ Cannot check dashboards (may need authentication)${NC}"
    fi
else
    echo -e "${RED}✗ Grafana is not accessible${NC}"
fi

echo ""

# Test 5: AlertManager Health
echo "=========================================="
echo "3. Testing AlertManager"
echo "=========================================="
echo ""

if curl -s -f "${ALERTMANAGER_URL}/-/healthy" > /dev/null 2>&1; then
    echo -e "${GREEN}✓ AlertManager is healthy${NC}"
    
    # Check alerts
    echo -e "${YELLOW}Checking alerts...${NC}"
    ALERTS_RESPONSE=$(curl -s "${ALERTMANAGER_URL}/api/v2/alerts" 2>/dev/null || echo "")
    
    if [ -n "$ALERTS_RESPONSE" ]; then
        ALERT_COUNT=$(echo "$ALERTS_RESPONSE" | grep -o '"status"' | wc -l)
        echo -e "${GREEN}✓ Alerts in AlertManager: ${ALERT_COUNT}${NC}"
    else
        echo -e "${YELLOW}⚠ Cannot check alerts${NC}"
    fi
else
    echo -e "${RED}✗ AlertManager is not accessible${NC}"
fi

echo ""

# Test 6: Service Metrics Endpoints
echo "=========================================="
echo "4. Testing Service Metrics Endpoints"
echo "=========================================="
echo ""

echo -e "${YELLOW}Checking service metrics endpoints (ports 9001-9022)...${NC}"

success_count=0
failed_count=0

for port in {9001..9022}; do
    if curl -s -f --max-time 2 "http://localhost:${port}/metrics" > /dev/null 2>&1; then
        echo -e "${GREEN}  ✓ Port ${port}: Accessible${NC}"
        ((success_count++))
    else
        echo -e "${YELLOW}  ⚠ Port ${port}: Not accessible (may be normal if service not running)${NC}"
        ((failed_count++))
    fi
done

echo ""
echo -e "Accessible: ${success_count} services"
echo -e "Not accessible: ${failed_count} services (may be normal)"
echo ""

# Summary
echo "=========================================="
echo "Test Summary"
echo "=========================================="
echo ""

echo -e "${GREEN}✅ Monitoring components tested${NC}"
echo ""
echo "📋 Next steps:"
echo "  1. Verify all services are running"
echo "  2. Check Prometheus targets are up"
echo "  3. Verify dashboards in Grafana"
echo "  4. Test alert firing"
echo ""

