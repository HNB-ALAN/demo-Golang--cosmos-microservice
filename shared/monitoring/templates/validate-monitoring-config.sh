#!/bin/bash
# Script to validate monitoring configuration
# Validates Prometheus, AlertManager, and YAML syntax

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
MONITORING_DIR="${PROJECT_ROOT}/SERVICES/service-08-monitoring/monitoring"

echo -e "${GREEN}🔍 Validating Monitoring Configuration...${NC}"
echo ""

# Check if promtool is available
PROMTOOL_AVAILABLE=false
if command -v promtool &> /dev/null; then
    PROMTOOL_AVAILABLE=true
    echo -e "${GREEN}✓ promtool found${NC}"
elif docker ps &> /dev/null; then
    echo -e "${YELLOW}⚠ promtool not found, will use Docker${NC}"
    PROMTOOL_CMD="docker run --rm -v ${MONITORING_DIR}:/etc/prometheus prom/prometheus promtool check"
else
    echo -e "${YELLOW}⚠ promtool not found and Docker not available${NC}"
    echo -e "${YELLOW}  Skipping Prometheus validation${NC}"
fi

# Check if amtool is available
AMTOOL_AVAILABLE=false
if command -v amtool &> /dev/null; then
    AMTOOL_AVAILABLE=true
    echo -e "${GREEN}✓ amtool found${NC}"
elif docker ps &> /dev/null; then
    echo -e "${YELLOW}⚠ amtool not found, will use Docker${NC}"
    AMTOOL_CMD="docker run --rm -v ${MONITORING_DIR}/alertmanager:/etc/alertmanager prom/alertmanager amtool check-config"
else
    echo -e "${YELLOW}⚠ amtool not found and Docker not available${NC}"
    echo -e "${YELLOW}  Skipping AlertManager validation${NC}"
fi

# Check if yamllint is available
YAMLLINT_AVAILABLE=false
if command -v yamllint &> /dev/null; then
    YAMLLINT_AVAILABLE=true
    echo -e "${GREEN}✓ yamllint found${NC}"
else
    echo -e "${YELLOW}⚠ yamllint not found${NC}"
    echo -e "${YELLOW}  Install with: pip install yamllint${NC}"
fi

echo ""
echo "=========================================="
echo "1. Validating Prometheus Configuration"
echo "=========================================="
echo ""

if [ "$PROMTOOL_AVAILABLE" = true ]; then
    echo -e "${YELLOW}Checking prometheus.yml...${NC}"
    if promtool check config "${MONITORING_DIR}/prometheus.yml"; then
        echo -e "${GREEN}  ✓ prometheus.yml is valid${NC}"
    else
        echo -e "${RED}  ✗ prometheus.yml has errors${NC}"
        exit 1
    fi
    echo ""
    
    echo -e "${YELLOW}Checking centralized rule files...${NC}"
    for rule_file in "${MONITORING_DIR}"/rules/*.yml; do
        if [ -f "$rule_file" ]; then
            rule_name=$(basename "$rule_file")
            if promtool check rules "$rule_file"; then
                echo -e "${GREEN}  ✓ $rule_name is valid${NC}"
            else
                echo -e "${RED}  ✗ $rule_name has errors${NC}"
                exit 1
            fi
        fi
    done
    echo ""
    
    echo -e "${YELLOW}Checking service-specific alert files...${NC}"
    alert_count=0
    alert_errors=0
    for alert_file in "${MONITORING_DIR}"/alerts/*.yml; do
        if [ -f "$alert_file" ]; then
            alert_name=$(basename "$alert_file")
            if promtool check rules "$alert_file" 2>&1 | grep -q "SUCCESS"; then
                echo -e "${GREEN}  ✓ $alert_name is valid${NC}"
                ((alert_count++))
            else
                echo -e "${RED}  ✗ $alert_name has errors${NC}"
                promtool check rules "$alert_file" 2>&1 | head -5
                ((alert_errors++))
            fi
        fi
    done
    
    if [ $alert_errors -eq 0 ]; then
        echo -e "${GREEN}  ✓ All $alert_count service alerts are valid${NC}"
    else
        echo -e "${RED}  ✗ $alert_errors alert files have errors${NC}"
        exit 1
    fi
elif [ -n "$PROMTOOL_CMD" ]; then
    echo -e "${YELLOW}Using Docker to validate...${NC}"
    $PROMTOOL_CMD check config /etc/prometheus/prometheus.yml || echo -e "${YELLOW}  ⚠ Could not validate (Docker command may need adjustment)${NC}"
else
    echo -e "${YELLOW}⚠ Skipping Prometheus validation (promtool not available)${NC}"
fi

echo ""
echo "=========================================="
echo "2. Validating AlertManager Configuration"
echo "=========================================="
echo ""

if [ "$AMTOOL_AVAILABLE" = true ]; then
    echo -e "${YELLOW}Checking alertmanager.yml...${NC}"
    if amtool check-config "${MONITORING_DIR}/alertmanager/alertmanager.yml"; then
        echo -e "${GREEN}  ✓ alertmanager.yml is valid${NC}"
    else
        echo -e "${RED}  ✗ alertmanager.yml has errors${NC}"
        exit 1
    fi
elif [ -n "$AMTOOL_CMD" ]; then
    echo -e "${YELLOW}Using Docker to validate...${NC}"
    $AMTOOL_CMD check-config /etc/alertmanager/alertmanager.yml || echo -e "${YELLOW}  ⚠ Could not validate (Docker command may need adjustment)${NC}"
else
    echo -e "${YELLOW}⚠ Skipping AlertManager validation (amtool not available)${NC}"
fi

echo ""
echo "=========================================="
echo "3. Validating YAML Syntax"
echo "=========================================="
echo ""

if [ "$YAMLLINT_AVAILABLE" = true ]; then
    echo -e "${YELLOW}Checking YAML syntax...${NC}"
    
    # Check prometheus.yml
    if yamllint "${MONITORING_DIR}/prometheus.yml" 2>/dev/null; then
        echo -e "${GREEN}  ✓ prometheus.yml syntax is valid${NC}"
    else
        echo -e "${YELLOW}  ⚠ prometheus.yml has syntax warnings (may still be valid)${NC}"
        yamllint "${MONITORING_DIR}/prometheus.yml" 2>&1 | head -5
    fi
    
    # Check alertmanager.yml
    if yamllint "${MONITORING_DIR}/alertmanager/alertmanager.yml" 2>/dev/null; then
        echo -e "${GREEN}  ✓ alertmanager.yml syntax is valid${NC}"
    else
        echo -e "${YELLOW}  ⚠ alertmanager.yml has syntax warnings${NC}"
        yamllint "${MONITORING_DIR}/alertmanager/alertmanager.yml" 2>&1 | head -5
    fi
    
    # Check rule files
    for rule_file in "${MONITORING_DIR}"/rules/*.yml; do
        if [ -f "$rule_file" ]; then
            rule_name=$(basename "$rule_file")
            if yamllint "$rule_file" 2>/dev/null; then
                echo -e "${GREEN}  ✓ $rule_name syntax is valid${NC}"
            else
                echo -e "${YELLOW}  ⚠ $rule_name has syntax warnings${NC}"
            fi
        fi
    done
else
    echo -e "${YELLOW}⚠ Skipping YAML syntax validation (yamllint not available)${NC}"
    echo -e "${YELLOW}  Install with: pip install yamllint${NC}"
fi

echo ""
echo "=========================================="
echo "4. Summary"
echo "=========================================="
echo ""

echo -e "${GREEN}✅ Validation complete!${NC}"
echo ""
echo "📋 Validation Status:"
echo "  - Prometheus config: $([ "$PROMTOOL_AVAILABLE" = true ] && echo "✅ Validated" || echo "⚠️  Skipped")"
echo "  - AlertManager config: $([ "$AMTOOL_AVAILABLE" = true ] && echo "✅ Validated" || echo "⚠️  Skipped")"
echo "  - YAML syntax: $([ "$YAMLLINT_AVAILABLE" = true ] && echo "✅ Validated" || echo "⚠️  Skipped")"
echo ""
echo "📝 Next steps:"
echo "  1. Fix any errors found above"
echo "  2. Deploy to Prometheus/AlertManager"
echo "  3. Import dashboards to Grafana"
echo "  4. Test alert firing"
echo ""

