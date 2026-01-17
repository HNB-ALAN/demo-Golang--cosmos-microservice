#!/bin/bash
# Script to COPY service alert files instead of symlinking
# Docker volume mounts don't follow symlinks, so we copy files instead

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
ALERTS_DIR="${PROJECT_ROOT}/SERVICES/service-08-monitoring/monitoring/alerts"

echo -e "${GREEN}📋 Copying Service Alert Files...${NC}"
echo ""
echo "Project root: ${PROJECT_ROOT}"
echo "Alerts directory: ${ALERTS_DIR}"
echo ""

# Create alerts directory if it doesn't exist
mkdir -p "${ALERTS_DIR}"

# Remove all existing files/symlinks in alerts directory
echo -e "${YELLOW}Cleaning existing alerts directory...${NC}"
rm -f "${ALERTS_DIR}"/*.yml
echo -e "${GREEN}✓ Cleaned${NC}"
echo ""

# Counter
COPIED=0

# Find and copy all service alert files using find -exec (most reliable)
echo -e "${YELLOW}Copying alert files from services...${NC}"
cd "${PROJECT_ROOT}"

# Use find with -exec to copy files directly (most reliable method)
find SERVICES/service-*/prometheus/alerts -name "*.yml" -type f -exec sh -c '
    for file; do
        alert_file=$(basename "$file")
        dest_path="'"${ALERTS_DIR}"'/${alert_file}"
        cp "$file" "$dest_path"
        echo -e "\033[0;32m✓ Copied: ${alert_file}\033[0m"
    done
' _ {} +

# Count copied files
COPIED=$(ls -1 "${ALERTS_DIR}"/*.yml 2>/dev/null | wc -l)

echo ""
if [ "$COPIED" -eq 0 ]; then
    echo -e "${RED}⚠️  No alert files found!${NC}"
    echo "Expected location: SERVICES/service-*/prometheus/alerts/*.yml"
    exit 1
fi

echo -e "${GREEN}✅ Copy Complete!${NC}"
echo "   - Files copied: ${COPIED}"
echo ""
echo "Alerts directory now contains:"
ls -1 "${ALERTS_DIR}"/*.yml 2>/dev/null | wc -l | xargs echo "   files"
echo ""
echo "Sample files:"
ls -1 "${ALERTS_DIR}"/*.yml 2>/dev/null | head -5 | sed 's/^/   /'
echo ""
echo -e "${BLUE}Next steps:${NC}"
echo "   1. Restart Prometheus: docker-compose restart prometheus"
echo "   2. Or reload Prometheus: curl -X POST http://localhost:9091/-/reload"
echo ""
