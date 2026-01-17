# 📋 TODO Tiếp Theo - Monitoring & Service Integration

**Date**: 2025-11-04  
**Status**: Monitoring fixes completed, ready for next steps

## ✅ Đã Hoàn Thành

### **Monitoring Configuration**
- ✅ Service alerts integration (21 symlinks created)
- ✅ prometheus.yml updated with `alerts/*.yml`
- ✅ Directory structure fixed
- ✅ AlertManager verified
- ✅ All 22 services have monitoring configs

## 🔥 TODO Tiếp Theo (Priority Order)

### **PRIORITY 1: Validate Monitoring Config** 🔴 HIGH

**Timeline**: 1-2 hours  
**Status**: ⏳ Pending

#### **1.1 Validate Prometheus Configuration**

```bash
# Install promtool if not available
# Docker: docker run --rm -v $(pwd):/etc/prometheus prom/prometheus promtool check config /etc/prometheus/prometheus.yml

# Validate config
promtool check config SERVICES/service-08-monitoring/monitoring/prometheus.yml

# Check rules
promtool check rules SERVICES/service-08-monitoring/monitoring/rules/*.yml
promtool check rules SERVICES/service-08-monitoring/monitoring/alerts/*.yml
```

**Expected Output**:
- ✅ Config is valid
- ✅ All rules are valid
- ✅ No syntax errors

#### **1.2 Validate AlertManager Configuration**

```bash
# Install amtool if not available
# Docker: docker run --rm -v $(pwd):/etc/alertmanager prom/alertmanager amtool check-config /etc/alertmanager/alertmanager.yml

# Validate config
amtool check-config SERVICES/service-08-monitoring/monitoring/alertmanager/alertmanager.yml
```

**Expected Output**:
- ✅ Config is valid
- ✅ Templates are valid
- ✅ Receivers are valid

#### **1.3 Validate YAML Syntax**

```bash
# Check all YAML files
yamllint SERVICES/service-08-monitoring/monitoring/prometheus.yml
yamllint SERVICES/service-08-monitoring/monitoring/alertmanager/alertmanager.yml
yamllint SERVICES/service-08-monitoring/monitoring/rules/*.yml
yamllint SERVICES/service-08-monitoring/monitoring/alerts/*.yml
```

**Deliverables**:
- [ ] Validation script created
- [ ] All configs validated
- [ ] Issues fixed (if any)

---

### **PRIORITY 2: Deploy & Test Monitoring** 🟡 MEDIUM

**Timeline**: 2-3 hours  
**Status**: ⏳ Pending

#### **2.1 Import Dashboards to Grafana**

**Option A: Via Grafana UI**
1. Open Grafana → Dashboards → Import
2. Upload each dashboard JSON file
3. Configure data source (Prometheus)
4. Save dashboard

**Option B: Via Grafana API**

```bash
# Create import script
cat > import-dashboards.sh << 'EOF'
#!/bin/bash
GRAFANA_URL="http://grafana:3000"
GRAFANA_API_KEY="your-api-key"

for dashboard in SERVICES/service-*/grafana/dashboards/*-overview.json; do
    if [ -f "$dashboard" ]; then
        echo "Importing: $dashboard"
        curl -X POST \
            -H "Authorization: Bearer $GRAFANA_API_KEY" \
            -H "Content-Type: application/json" \
            -d @"$dashboard" \
            "$GRAFANA_URL/api/dashboards/db"
    fi
done
EOF
chmod +x import-dashboards.sh
```

**Deliverables**:
- [ ] All 21 dashboards imported
- [ ] Dashboards accessible in Grafana
- [ ] Data sources configured

#### **2.2 Add Alert Rules to Prometheus**

**Option A: Copy to Prometheus Config Directory**

```bash
# Copy alert rules
cp SERVICES/service-08-monitoring/monitoring/rules/*.yml /etc/prometheus/rules/
cp SERVICES/service-08-monitoring/monitoring/alerts/*.yml /etc/prometheus/alerts/

# Reload Prometheus
curl -X POST http://prometheus:9090/-/reload
```

**Option B: Use Volume Mounts (Docker/Kubernetes)**

```yaml
# docker-compose.yml or Kubernetes
volumes:
  - ./SERVICES/service-08-monitoring/monitoring/rules:/etc/prometheus/rules
  - ./SERVICES/service-08-monitoring/monitoring/alerts:/etc/prometheus/alerts
```

**Deliverables**:
- [ ] Alert rules loaded in Prometheus
- [ ] Prometheus reloaded successfully
- [ ] Rules visible in Prometheus UI

#### **2.3 Test Alert Loading**

```bash
# Check loaded rules
curl http://prometheus:9090/api/v1/rules

# Check specific service alerts
curl http://prometheus:9090/api/v1/rules?rule_group=gateway_critical

# Verify all 21 service alerts are loaded
# Should see 25 rule files total (4 centralized + 21 service-specific)
```

**Deliverables**:
- [ ] All rules loaded
- [ ] Service alerts visible
- [ ] No errors in Prometheus logs

#### **2.4 Test Alert Firing**

```bash
# Trigger test alert (if possible)
# Or wait for actual service issues

# Check AlertManager
curl http://alertmanager:9093/api/v2/alerts

# Check notification channels
# Verify Slack/PagerDuty/Email receive alerts
```

**Deliverables**:
- [ ] Test alerts fire correctly
- [ ] AlertManager routing works
- [ ] Notifications delivered

#### **2.5 Verify Metrics Collection**

```bash
# Test metrics endpoints for each service
for port in {9001..9022}; do
    echo "Testing service on port $port..."
    curl -s http://localhost:$port/metrics | head -5
done
```

**Deliverables**:
- [ ] All services expose metrics
- [ ] Prometheus scraping successfully
- [ ] Metrics visible in Grafana

---

### **PRIORITY 3: Documentation & Cleanup** 🔵 LOW

**Timeline**: 1-2 hours  
**Status**: ⏳ Pending

#### **3.1 Update Monitoring Documentation**

- [ ] Update `README.md` with deployment steps
- [ ] Document alert runbooks
- [ ] Create troubleshooting guide
- [ ] Document dashboard usage

#### **3.2 Create Deployment Guide**

**File**: `shared/monitoring/DEPLOYMENT_GUIDE.md`

**Contents**:
- Prerequisites
- Installation steps
- Configuration
- Validation
- Troubleshooting

#### **3.3 Document Alert Runbooks**

**File**: `shared/monitoring/ALERT_RUNBOOKS.md`

**Contents**:
- Gateway alerts
- Service alerts
- Infrastructure alerts
- Resolution steps

#### **3.4 Cleanup**

- [ ] Remove temporary files
- [ ] Organize documentation
- [ ] Update TODO status
- [ ] Create summary report

---

## 📊 Implementation Checklist

### **Phase 1: Validation** (1-2h)
- [ ] Install promtool/amtool
- [ ] Validate Prometheus config
- [ ] Validate AlertManager config
- [ ] Validate YAML syntax
- [ ] Fix any issues found

### **Phase 2: Deployment** (2-3h)
- [ ] Import dashboards to Grafana
- [ ] Add alert rules to Prometheus
- [ ] Test alert loading
- [ ] Test alert firing
- [ ] Verify metrics collection

### **Phase 3: Documentation** (1-2h)
- [ ] Update monitoring docs
- [ ] Create deployment guide
- [ ] Document alert runbooks
- [ ] Cleanup and organize

---

## 🎯 Quick Start

### **Immediate Next Steps (1 hour)**

1. **Validate Configs**:
   ```bash
   # Check if promtool/amtool available
   which promtool amtool
   
   # If not, install or use Docker
   # Validate configs
   ```

2. **Test Integration**:
   ```bash
   # Run integration script again to verify
   ./shared/monitoring/templates/integrate-service-alerts.sh
   
   # Check symlinks
   ls -la SERVICES/service-08-monitoring/monitoring/alerts/
   ```

3. **Prepare for Deployment**:
   ```bash
   # Create deployment checklist
   # Prepare Grafana import script
   # Prepare Prometheus reload script
   ```

---

## 📝 Notes

### **Current Status**
- ✅ Monitoring configs complete
- ✅ All services integrated
- ⏳ Ready for validation and deployment

### **Dependencies**
- Prometheus instance running
- Grafana instance running
- AlertManager instance running
- Services exposing metrics on ports 9001-9022

### **Timeline**
- **Validation**: 1-2 hours
- **Deployment**: 2-3 hours
- **Documentation**: 1-2 hours
- **Total**: 4-7 hours

---

**Status**: ⏳ **Ready for Phase 1: Validation**

**Next Action**: Validate Prometheus and AlertManager configurations

