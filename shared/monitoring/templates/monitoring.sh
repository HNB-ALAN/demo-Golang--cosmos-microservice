#!/bin/bash
# Master Monitoring Script - Simplified interface for all monitoring operations
# Usage: ./monitoring.sh <command> [options]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Show usage
usage() {
    echo "📊 Monitoring Scripts - Master Command"
    echo ""
    echo "Usage: $0 <command> [options]"
    echo ""
    echo "Commands:"
    echo "  generate <service-id> <name> <prefix> <port> [threshold] [--no-sync]"
    echo "    Generate monitoring configs for a service"
    echo "    Example: $0 generate service-02-auth Auth auth 9002 0.1"
    echo ""
    echo "  generate-all"
    echo "    Generate monitoring configs for all services"
    echo ""
    echo "  sync"
    echo "    Sync alerts and dashboards to centralized location"
    echo ""
    echo "  validate"
    echo "    Validate monitoring configurations"
    echo ""
    echo "  test"
    echo "    Test monitoring stack (Prometheus, Grafana, AlertManager)"
    echo ""
    echo "  reload"
    echo "    Reload Prometheus configuration"
    echo ""
    echo "  help"
    echo "    Show this help message"
    echo ""
    exit 1
}

# Main command handler
case "$1" in
    generate)
        shift
        if [ $# -lt 4 ]; then
            echo -e "${YELLOW}Usage: $0 generate <service-id> <name> <prefix> <port> [threshold] [--no-sync]${NC}"
            exit 1
        fi
        ./generate-monitoring.sh "$@"
        ;;
    
    generate-all)
        ./generate-all-services.sh
        ;;
    
    sync)
        ./sync-monitoring.sh
        ;;
    
    validate)
        ./validate-monitoring-config.sh
        ;;
    
    test)
        ./test-monitoring.sh
        ;;
    
    reload)
        ./reload-prometheus.sh
        ;;
    
    help|--help|-h)
        usage
        ;;
    
    *)
        echo -e "${YELLOW}Unknown command: $1${NC}"
        echo ""
        usage
        ;;
esac

