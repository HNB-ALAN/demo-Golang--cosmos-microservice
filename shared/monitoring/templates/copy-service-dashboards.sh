#!/bin/bash
# Script to COPY service dashboard files instead of importing via API
# Docker volume mounts don't follow symlinks, so we copy files for provisioning

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
# Grafana provisioning expects dashboards in provisioning/dashboards/ directory
DASHBOARDS_DIR="${PROJECT_ROOT}/SERVICES/service-08-monitoring/monitoring/grafana/provisioning/dashboards"

echo -e "${GREEN}📊 Copying Service Dashboard Files...${NC}"
echo ""
echo "Project root: ${PROJECT_ROOT}"
echo "Dashboards directory: ${DASHBOARDS_DIR}"
echo ""

# Create dashboards directory if it doesn't exist
mkdir -p "${DASHBOARDS_DIR}"

# Remove existing service dashboard files (but keep dashboard.yml config)
echo -e "${YELLOW}Cleaning existing service dashboards...${NC}"
# Remove all JSON files except dashboard.yml (config file)
find "${DASHBOARDS_DIR}" -maxdepth 1 -name "*.json" 2>/dev/null | while read file; do
    if [ -f "$file" ]; then
        rm -f "$file"
    fi
done
echo -e "${GREEN}✓ Cleaned${NC}"
echo ""

# Counter
COPIED=0

# Find and copy all service dashboard files using find -exec (most reliable)
echo -e "${YELLOW}Copying dashboard files from services...${NC}"
cd "${PROJECT_ROOT}"

# Use find with -exec to copy files directly (most reliable method)
# Find all dashboard files (both *-overview.json and gateway-overview.json)
find SERVICES/service-*/grafana/dashboards -type f \( -name "*-overview.json" -o -name "gateway-overview.json" \) -exec sh -c '
    for file; do
        dashboard_name=$(basename "$file")
        dest_path="'"${DASHBOARDS_DIR}"'/${dashboard_name}"
        
        # Copy file
        cp "$file" "$dest_path"
        echo -e "\033[0;32m✓ Copied: ${dashboard_name}\033[0m"
    done
' _ {} +

# Count copied files
COPIED=$(find "${DASHBOARDS_DIR}" -maxdepth 1 -name "*.json" ! -name "dashboard.yml" 2>/dev/null | wc -l)

echo ""
if [ "$COPIED" -eq 0 ]; then
    echo -e "${RED}⚠️  No dashboard files found!${NC}"
    echo "Expected location: SERVICES/service-*/grafana/dashboards/*-overview.json"
    echo "or: SERVICES/service-01-gateway/grafana/dashboards/gateway-overview.json"
    exit 1
fi

echo -e "${GREEN}✅ Copy Complete!${NC}"
echo "   - Service dashboards copied: ${COPIED}"
echo ""
echo "Dashboards directory now contains:"
find "${DASHBOARDS_DIR}" -maxdepth 1 -name "*.json" ! -name "dashboard.yml" 2>/dev/null | wc -l | xargs echo "   dashboard files"
echo ""
echo "Sample files:"
find "${DASHBOARDS_DIR}" -maxdepth 1 -name "*.json" ! -name "dashboard.yml" 2>/dev/null | head -5 | sed 's|^.*/||' | sed 's/^/   /'
echo ""
echo -e "${BLUE}Next steps:${NC}"
echo "   1. Grafana will auto-reload dashboards (provisioning enabled)"
echo "   2. Or restart Grafana: docker-compose restart grafana"
echo "   3. Verify in Grafana UI: http://localhost:3000"
echo ""

