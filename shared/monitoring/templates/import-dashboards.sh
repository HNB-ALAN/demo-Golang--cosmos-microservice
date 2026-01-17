#!/bin/bash
# Script to import Grafana dashboards
# Supports both UI and API methods

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
GRAFANA_URL="${GRAFANA_URL:-http://localhost:3000}"
GRAFANA_API_KEY="${GRAFANA_API_KEY:-}"
GRAFANA_USER="${GRAFANA_USER:-admin}"
GRAFANA_PASSWORD="${GRAFANA_PASSWORD:-admin}"

echo -e "${GREEN}📊 Importing Grafana Dashboards...${NC}"
echo ""

# Check if Grafana is accessible
if ! curl -s -f "${GRAFANA_URL}/api/health" > /dev/null 2>&1; then
    echo -e "${RED}❌ Cannot connect to Grafana at ${GRAFANA_URL}${NC}"
    echo -e "${YELLOW}Please ensure Grafana is running and accessible${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Grafana is accessible at ${GRAFANA_URL}${NC}"
echo ""

# Get API key if not provided
if [ -z "$GRAFANA_API_KEY" ]; then
    echo -e "${YELLOW}API key not provided, attempting to authenticate...${NC}"
    
    # Try to get API key via login
    AUTH_RESPONSE=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d "{\"user\":\"${GRAFANA_USER}\",\"password\":\"${GRAFANA_PASSWORD}\"}" \
        "${GRAFANA_URL}/api/login")
    
    if echo "$AUTH_RESPONSE" | grep -q "Logged in"; then
        echo -e "${GREEN}✓ Authenticated successfully${NC}"
    else
        echo -e "${YELLOW}⚠ Could not authenticate, will use provided credentials${NC}"
        echo -e "${YELLOW}Please provide GRAFANA_API_KEY environment variable${NC}"
        echo ""
        echo -e "${BLUE}Alternative: Import dashboards manually via Grafana UI:${NC}"
        echo "  1. Open Grafana → Dashboards → Import"
        echo "  2. Upload each JSON file from:"
        echo "     SERVICES/service-*/grafana/dashboards/*-overview.json"
        exit 0
    fi
fi

# Import dashboards
echo -e "${YELLOW}Importing dashboards...${NC}"
echo ""

count=0
success_count=0
failed_count=0

for dashboard_file in "${PROJECT_ROOT}"/SERVICES/service-*/grafana/dashboards/*-overview.json; do
    if [ -f "$dashboard_file" ]; then
        dashboard_name=$(basename "$dashboard_file")
        service_name=$(basename "$(dirname "$(dirname "$dashboard_file")")")
        
        echo -e "${YELLOW}Importing: ${dashboard_name}${NC}"
        
        # Prepare dashboard JSON
        # Grafana API expects dashboard object with dashboard field
        DASHBOARD_JSON=$(cat "$dashboard_file")
        
        # Create API payload
        PAYLOAD=$(cat <<EOF
{
  "dashboard": ${DASHBOARD_JSON},
  "overwrite": true,
  "inputs": []
}
EOF
)
        
        # Import via API
        if [ -n "$GRAFANA_API_KEY" ]; then
            RESPONSE=$(curl -s -w "\n%{http_code}" -X POST \
                -H "Authorization: Bearer ${GRAFANA_API_KEY}" \
                -H "Content-Type: application/json" \
                -d "$PAYLOAD" \
                "${GRAFANA_URL}/api/dashboards/db")
        else
            RESPONSE=$(curl -s -w "\n%{http_code}" -X POST \
                -u "${GRAFANA_USER}:${GRAFANA_PASSWORD}" \
                -H "Content-Type: application/json" \
                -d "$PAYLOAD" \
                "${GRAFANA_URL}/api/dashboards/db")
        fi
        
        HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
        BODY=$(echo "$RESPONSE" | sed '$d')
        
        if [ "$HTTP_CODE" = "200" ] || [ "$HTTP_CODE" = "201" ]; then
            echo -e "${GREEN}  ✓ ${dashboard_name} imported successfully${NC}"
            ((success_count++))
        else
            echo -e "${RED}  ✗ ${dashboard_name} failed (HTTP ${HTTP_CODE})${NC}"
            echo "$BODY" | head -3
            ((failed_count++))
        fi
        
        ((count++))
    fi
done

echo ""
echo "=========================================="
echo "Import Summary"
echo "=========================================="
echo ""
echo -e "Total dashboards: ${count}"
echo -e "${GREEN}Successfully imported: ${success_count}${NC}"
if [ $failed_count -gt 0 ]; then
    echo -e "${RED}Failed: ${failed_count}${NC}"
fi

echo ""
if [ $success_count -eq $count ]; then
    echo -e "${GREEN}✅ All dashboards imported successfully!${NC}"
    echo ""
    echo "📋 Next steps:"
    echo "  1. Open Grafana: ${GRAFANA_URL}"
    echo "  2. Navigate to Dashboards"
    echo "  3. Verify all dashboards are visible"
    echo "  4. Configure data sources if needed"
else
    echo -e "${YELLOW}⚠ Some dashboards failed to import${NC}"
    echo ""
    echo "📋 Troubleshooting:"
    echo "  1. Check Grafana API key/credentials"
    echo "  2. Verify dashboard JSON syntax"
    echo "  3. Check Grafana logs"
    echo "  4. Try manual import via UI"
fi

echo ""

