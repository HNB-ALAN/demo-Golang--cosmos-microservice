#!/bin/bash
# Unified script to sync all monitoring files (alerts + dashboards)
# Automatically called after generating monitoring configs or can be run manually

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo -e "${GREEN}🔄 Syncing Monitoring Files...${NC}"
echo ""

# Run copy-service-alerts.sh
echo -e "${YELLOW}Step 1: Copying service alerts...${NC}"
"${SCRIPT_DIR}/copy-service-alerts.sh"
echo ""

# Run copy-service-dashboards.sh
echo -e "${YELLOW}Step 2: Copying service dashboards...${NC}"
"${SCRIPT_DIR}/copy-service-dashboards.sh"
echo ""

echo -e "${GREEN}✅ Monitoring Sync Complete!${NC}"
echo ""
echo -e "${BLUE}Next steps:${NC}"
echo "   1. Restart Prometheus: docker-compose restart prometheus"
echo "   2. Restart Grafana: docker-compose restart grafana"
echo "   3. Or reload without restart:"
echo "      - Prometheus: curl -X POST http://localhost:9091/-/reload"
echo "      - Grafana: Auto-reloads every 10 seconds"
echo ""

