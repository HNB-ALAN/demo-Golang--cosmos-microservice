#!/bin/bash

# Generate monitoring configs for a service from templates
# Usage: ./generate-monitoring.sh <service-id> <service-name> <metric-prefix> <metrics-port> [latency-threshold]
#
# Example:
#   ./generate-monitoring.sh service-02-auth Auth auth 9002 0.1
#   ./generate-monitoring.sh service-05-usc-wallet "USC Wallet" wallet 9005 0.2

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check arguments
if [ $# -lt 4 ]; then
    echo -e "${RED}Error: Missing required arguments${NC}"
    echo "Usage: $0 <service-id> <service-name> <metric-prefix> <metrics-port> [latency-threshold]"
    echo ""
    echo "Arguments:"
    echo "  service-id        : Service identifier (e.g., service-02-auth)"
    echo "  service-name      : Service display name (e.g., Auth)"
    echo "  metric-prefix     : Metric name prefix (e.g., auth)"
    echo "  metrics-port      : Metrics endpoint port (e.g., 9002)"
    echo "  latency-threshold : Optional latency threshold in seconds (default: 0.1)"
    echo "  --no-sync         : Optional flag to skip auto-sync to centralized location"
    echo ""
    echo "Examples:"
    echo "  $0 service-02-auth Auth auth 9002"
    echo "  $0 service-05-usc-wallet \"USC Wallet\" wallet 9005 0.2"
    echo "  $0 service-02-auth Auth auth 9002 0.1 --no-sync  # Skip auto-sync"
    exit 1
fi

# Check for --no-sync flag and extract it
NO_SYNC=false
ARGS=()
for arg in "$@"; do
    if [ "$arg" = "--no-sync" ]; then
        NO_SYNC=true
    else
        ARGS+=("$arg")
    fi
done

# Set arguments without --no-sync
set -- "${ARGS[@]}"

# Check arguments again after removing --no-sync
if [ $# -lt 4 ]; then
    echo -e "${RED}Error: Missing required arguments${NC}"
    echo "Usage: $0 <service-id> <service-name> <metric-prefix> <metrics-port> [latency-threshold] [--no-sync]"
    echo ""
    echo "Arguments:"
    echo "  service-id        : Service identifier (e.g., service-02-auth)"
    echo "  service-name      : Service display name (e.g., Auth)"
    echo "  metric-prefix     : Metric name prefix (e.g., auth)"
    echo "  metrics-port      : Metrics endpoint port (e.g., 9002)"
    echo "  latency-threshold : Optional latency threshold in seconds (default: 0.1)"
    echo "  --no-sync         : Optional flag to skip auto-sync to centralized location"
    echo ""
    echo "Examples:"
    echo "  $0 service-02-auth Auth auth 9002"
    echo "  $0 service-05-usc-wallet \"USC Wallet\" wallet 9005 0.2"
    echo "  $0 service-02-auth Auth auth 9002 0.1 --no-sync  # Skip auto-sync"
    exit 1
fi

SERVICE_ID=$1
SERVICE_NAME=$2
METRIC_PREFIX=$3
METRICS_PORT=$4
LATENCY_THRESHOLD=${5:-0.1}

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# Go up from templates/ to shared/ to project root
# templates/ -> shared/monitoring/ -> shared/ -> project root
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
SERVICE_DIR="$PROJECT_ROOT/SERVICES/${SERVICE_ID}"

# Validate service directory exists
if [ ! -d "$SERVICE_DIR" ]; then
    echo -e "${RED}Error: Service directory not found: $SERVICE_DIR${NC}"
    exit 1
fi

echo -e "${GREEN}Generating monitoring configs for ${SERVICE_NAME}...${NC}"
echo "  Service ID: $SERVICE_ID"
echo "  Service Name: $SERVICE_NAME"
echo "  Metric Prefix: $METRIC_PREFIX"
echo "  Metrics Port: $METRICS_PORT"
echo "  Latency Threshold: ${LATENCY_THRESHOLD}s"
echo ""

# Create directories
mkdir -p "${SERVICE_DIR}/grafana/dashboards"
mkdir -p "${SERVICE_DIR}/prometheus/alerts"
mkdir -p "${SERVICE_DIR}/docs"

# Generate dashboard from template
echo -e "${YELLOW}Generating dashboard...${NC}"
sed -e "s/{{SERVICE_ID}}/${SERVICE_ID}/g" \
    -e "s/{{SERVICE_NAME}}/${SERVICE_NAME}/g" \
    -e "s/{{METRIC_PREFIX}}/${METRIC_PREFIX}/g" \
    -e "s/{{METRICS_PORT}}/${METRICS_PORT}/g" \
    "${SCRIPT_DIR}/base-dashboard-template.json" > \
    "${SERVICE_DIR}/grafana/dashboards/${SERVICE_ID}-overview.json"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}  ✓ Dashboard generated: ${SERVICE_DIR}/grafana/dashboards/${SERVICE_ID}-overview.json${NC}"
else
    echo -e "${RED}  ✗ Failed to generate dashboard${NC}"
    exit 1
fi

# Generate alerts from template
echo -e "${YELLOW}Generating alerts...${NC}"
sed -e "s/{{SERVICE_ID}}/${SERVICE_ID}/g" \
    -e "s/{{SERVICE_NAME}}/${SERVICE_NAME}/g" \
    -e "s/{{METRIC_PREFIX}}/${METRIC_PREFIX}/g" \
    -e "s/{{LATENCY_THRESHOLD}}/${LATENCY_THRESHOLD}/g" \
    "${SCRIPT_DIR}/base-alerts-template.yml" > \
    "${SERVICE_DIR}/prometheus/alerts/${SERVICE_ID}-alerts.yml"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}  ✓ Alerts generated: ${SERVICE_DIR}/prometheus/alerts/${SERVICE_ID}-alerts.yml${NC}"
else
    echo -e "${RED}  ✗ Failed to generate alerts${NC}"
    exit 1
fi

# Generate docs from template
echo -e "${YELLOW}Generating documentation...${NC}"
sed -e "s/{{SERVICE_ID}}/${SERVICE_ID}/g" \
    -e "s/{{SERVICE_NAME}}/${SERVICE_NAME}/g" \
    -e "s/{{METRIC_PREFIX}}/${METRIC_PREFIX}/g" \
    -e "s/{{METRICS_PORT}}/${METRICS_PORT}/g" \
    -e "s/{{LATENCY_THRESHOLD}}/${LATENCY_THRESHOLD}/g" \
    "${SCRIPT_DIR}/monitoring-setup-template.md" > \
    "${SERVICE_DIR}/docs/MONITORING_SETUP.md"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}  ✓ Documentation generated: ${SERVICE_DIR}/docs/MONITORING_SETUP.md${NC}"
else
    echo -e "${RED}  ✗ Failed to generate documentation${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}✅ Successfully generated monitoring configs for ${SERVICE_NAME}!${NC}"
echo ""

# Auto-sync to centralized location (unless --no-sync flag is provided)
if [ "$NO_SYNC" = false ]; then
    echo -e "${YELLOW}Syncing to centralized monitoring...${NC}"
    # Call sync script (which will copy alerts and dashboards)
    if [ -f "${SCRIPT_DIR}/sync-monitoring.sh" ]; then
        "${SCRIPT_DIR}/sync-monitoring.sh" > /dev/null 2>&1 || echo -e "${YELLOW}  ⚠️  Sync skipped (can run manually later)${NC}"
    else
        # Fallback: call copy scripts directly
        "${SCRIPT_DIR}/copy-service-alerts.sh" > /dev/null 2>&1 || true
        "${SCRIPT_DIR}/copy-service-dashboards.sh" > /dev/null 2>&1 || true
    fi
    echo ""
fi

echo "Next steps:"
echo "  1. Review and customize generated files if needed"
echo "  2. Verify metrics endpoint: curl http://localhost:${METRICS_PORT}/metrics"
if [ "$NO_SYNC" = true ]; then
    echo "  3. Sync to centralized: ./sync-monitoring.sh"
    echo "  4. Restart Prometheus & Grafana to load new configs"
else
    echo "  3. Restart Prometheus & Grafana to load new configs:"
    echo "     - Prometheus: docker-compose restart prometheus"
    echo "     - Grafana: docker-compose restart grafana"
fi
echo ""

