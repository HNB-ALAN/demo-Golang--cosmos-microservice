# 🚀 Priority 2: Deploy & Test Monitoring - Deployment Guide

**Date**: 2025-11-04  
**Status**: Ready for deployment

## 📋 Overview

This guide covers deploying monitoring configuration to Prometheus, Grafana, and AlertManager, and testing the setup.

## 🎯 Deployment Steps

### **Step 1: Import Dashboards to Grafana**

#### **Option A: Automated (Script)**

```bash
# Set Grafana credentials
export GRAFANA_URL="http://localhost:3000"
export GRAFANA_API_KEY="your-api-key"  # Optional, can use username/password

# Run import script
./shared/monitoring/templates/import-dashboards.sh
```

#### **Option B: Manual (UI)**

1. Open Grafana: `http://localhost:3000`
2. Navigate to: **Dashboards → Import**
3. For each service, upload:
   - `SERVICES/service-01-gateway/grafana/dashboards/gateway-overview.json`
   - `SERVICES/service-02-auth/grafana/dashboards/service-02-auth-overview.json`
   - ... (all 21 dashboards)

#### **Option C: Grafana API**

```bash
# Get API key first
curl -X POST http://localhost:3000/api/auth/keys \
  -H "Content-Type: application/json" \
  -d '{"name":"monitoring","role":"Admin","secondsToLive":86400}'

# Import dashboard
curl -X POST http://localhost:3000/api/dashboards/db \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d @SERVICES/service-01-gateway/grafana/dashboards/gateway-overview.json
```

**Expected Result**: 21 dashboards imported successfully

---

### **Step 2: Add Alert Rules to Prometheus**

#### **Option A: Docker Volume Mount** (Recommended)

```yaml
# docker-compose.yml
services:
  prometheus:
    image: prom/prometheus:latest
    volumes:
      - ./SERVICES/service-08-monitoring/monitoring/prometheus.yml:/etc/prometheus/prometheus.yml
      - ./SERVICES/service-08-monitoring/monitoring/rules:/etc/prometheus/rules
      - ./SERVICES/service-08-monitoring/monitoring/alerts:/etc/prometheus/alerts
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
```

#### **Option B: Copy Files**

```bash
# Copy centralized rules
cp SERVICES/service-08-monitoring/monitoring/rules/*.yml /etc/prometheus/rules/

# Copy service-specific alerts
cp SERVICES/service-08-monitoring/monitoring/alerts/*.yml /etc/prometheus/alerts/

# Reload Prometheus
curl -X POST http://localhost:9090/-/reload
```

#### **Option C: Script**

```bash
# Use reload script
./shared/monitoring/templates/reload-prometheus.sh

# Or verify rules
METHOD=verify ./shared/monitoring/templates/reload-prometheus.sh
```

**Expected Result**: 25 rule files loaded (4 centralized + 21 service-specific)

---

### **Step 3: Configure AlertManager**

#### **Docker Volume Mount**

```yaml
# docker-compose.yml
services:
  alertmanager:
    image: prom/alertmanager:latest
    volumes:
      - ./SERVICES/service-08-monitoring/monitoring/alertmanager/alertmanager.yml:/etc/alertmanager/alertmanager.yml
      - ./SERVICES/service-08-monitoring/monitoring/alertmanager/templates:/etc/alertmanager/templates
    command:
      - '--config.file=/etc/alertmanager/alertmanager.yml'
      - '--storage.path=/alertmanager'
```

#### **Update Notification Channels**

Edit `alertmanager.yml`:
- Update Slack webhook URL
- Update email SMTP settings
- Update PagerDuty service key

**Expected Result**: AlertManager routes alerts correctly

---

### **Step 4: Test Alert Loading**

```bash
# Check loaded rules
curl http://localhost:9090/api/v1/rules

# Check specific rule group
curl http://localhost:9090/api/v1/rules?rule_group=gateway_critical

# Check Prometheus targets
curl http://localhost:9090/api/v1/targets
```

**Expected Result**: 
- All rule groups visible
- Service targets are up
- Rules are evaluated

---

### **Step 5: Test Alert Firing**

#### **Manual Test**

```bash
# Check active alerts
curl http://localhost:9090/api/v1/alerts

# Check AlertManager alerts
curl http://localhost:9093/api/v2/alerts

# Force test alert (if possible)
# Or wait for actual service issues
```

#### **Verify Notifications**

- Check Slack channels
- Check email inbox
- Check PagerDuty (if configured)

**Expected Result**: Alerts fire and notifications are delivered

---

### **Step 6: Verify Metrics Collection**

```bash
# Test metrics endpoints
for port in {9001..9022}; do
    echo "Testing port $port..."
    curl -s http://localhost:$port/metrics | head -5
done

# Or use test script
./shared/monitoring/templates/test-monitoring.sh
```

**Expected Result**: All services expose metrics

---

## 🧪 Testing Script

A comprehensive test script is available:

```bash
./shared/monitoring/templates/test-monitoring.sh
```

**Tests**:
1. Prometheus health
2. Prometheus rules
3. Prometheus targets
4. Grafana health
5. Grafana dashboards
6. AlertManager health
7. AlertManager alerts
8. Service metrics endpoints

---

## 📊 Verification Checklist

After deployment:

- [ ] All 21 dashboards imported to Grafana
- [ ] Prometheus config reloaded
- [ ] All 25 rule files loaded
- [ ] AlertManager configured
- [ ] Notification channels configured
- [ ] Service metrics accessible
- [ ] Prometheus scraping successfully
- [ ] Rules evaluated correctly
- [ ] Alerts fire when conditions met
- [ ] Notifications delivered

---

## 🐛 Troubleshooting

### **Dashboards Not Showing**

1. Check Grafana data source is configured
2. Verify Prometheus is accessible from Grafana
3. Check dashboard JSON syntax
4. Verify permissions

### **Rules Not Loading**

1. Check Prometheus logs
2. Verify file paths in `prometheus.yml`
3. Check rule file syntax with `promtool`
4. Verify file permissions

### **Alerts Not Firing**

1. Check rule expressions are correct
2. Verify metrics are being collected
3. Check alert thresholds
4. Verify AlertManager is connected

### **Metrics Not Collected**

1. Check service is running
2. Verify metrics endpoint is accessible
3. Check Prometheus scrape config
4. Verify network connectivity

---

## 📝 Next Steps

After successful deployment:

1. **Monitor in staging** for 24-48 hours
2. **Tune alert thresholds** based on actual metrics
3. **Update runbooks** with actual resolution steps
4. **Document alert response procedures**
5. **Set up on-call rotation**

---

**Status**: ⏳ **Ready for deployment**

**Scripts Available**:
- `import-dashboards.sh` - Import dashboards
- `reload-prometheus.sh` - Reload Prometheus
- `test-monitoring.sh` - Test monitoring setup

