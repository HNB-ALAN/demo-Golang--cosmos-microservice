# 🚀 Monitoring Deployment Guide

**Service-08: Monitoring Service** - Complete deployment guide for USC Platform monitoring infrastructure.

## 📋 Prerequisites

- Docker & Docker Compose (or Kubernetes)
- Prometheus instance
- Grafana instance
- AlertManager instance
- All 22 services running and exposing metrics on ports 9001-9022

## 🚀 Quick Start

### **Step 1: Integrate Service Alerts**

```bash
cd shared/monitoring/templates/
./integrate-service-alerts.sh
```

This creates symlinks for all service-specific alerts in `monitoring/alerts/`.

### **Step 2: Validate Configuration**

```bash
./validate-monitoring-config.sh
```

Validates Prometheus, AlertManager, and YAML syntax.

### **Step 3: Deploy to Prometheus**

#### **Option A: Docker Compose**

```yaml
# docker-compose.yml
services:
  prometheus:
    image: prom/prometheus:latest
    volumes:
      - ./SERVICES/service-08-monitoring/monitoring/prometheus.yml:/etc/prometheus/prometheus.yml
      - ./SERVICES/service-08-monitoring/monitoring/rules:/etc/prometheus/rules
      - ./SERVICES/service-08-monitoring/monitoring/alerts:/etc/prometheus/alerts
    ports:
      - "9090:9090"
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--web.console.libraries=/usr/share/prometheus/console_libraries'
      - '--web.console.templates=/usr/share/prometheus/consoles'
```

#### **Option B: Reload Existing Prometheus**

```bash
./shared/monitoring/templates/reload-prometheus.sh
```

### **Step 4: Deploy to AlertManager**

```yaml
# docker-compose.yml
services:
  alertmanager:
    image: prom/alertmanager:latest
    volumes:
      - ./SERVICES/service-08-monitoring/monitoring/alertmanager/alertmanager.yml:/etc/alertmanager/alertmanager.yml
      - ./SERVICES/service-08-monitoring/monitoring/alertmanager/templates:/etc/alertmanager/templates
    ports:
      - "9093:9093"
    command:
      - '--config.file=/etc/alertmanager/alertmanager.yml'
      - '--storage.path=/alertmanager'
```

### **Step 5: Import Dashboards to Grafana**

```bash
# Set credentials
export GRAFANA_URL="http://localhost:3000"
export GRAFANA_API_KEY="your-api-key"

# Import all dashboards
./shared/monitoring/templates/import-dashboards.sh
```

Or manually via Grafana UI:
1. Open Grafana → Dashboards → Import
2. Upload each JSON file from `SERVICES/service-*/grafana/dashboards/*.json`

### **Step 6: Test Monitoring Setup**

```bash
./shared/monitoring/templates/test-monitoring.sh
```

## 📊 Configuration Details

### **Prometheus Configuration**

**File**: `SERVICES/service-08-monitoring/monitoring/prometheus.yml`

**Key Features**:
- Scrapes all 22 services (ports 9001-9022)
- Loads 4 centralized rule files
- Loads 21 service-specific alert files
- Monitors infrastructure (databases, caches, queues)
- External monitoring via Blackbox exporter

**Rule Files**:
- `rules/alerting_rules.yml` - General alerts
- `rules/recording_rules.yml` - Pre-computed metrics
- `rules/slo_rules.yml` - SLO monitoring
- `rules/trading_rules.yml` - Trading-specific alerts
- `alerts/*.yml` - Service-specific alerts (21 files)

### **AlertManager Configuration**

**File**: `SERVICES/service-08-monitoring/monitoring/alertmanager/alertmanager.yml`

**Routing**:
- Critical alerts → PagerDuty
- Trading alerts → Slack #trading-alerts + Email
- Infrastructure alerts → Slack #infrastructure-alerts
- Business alerts → Slack #business-alerts
- SLO alerts → Slack #slo-alerts

**Notification Channels**:
- Slack (multiple channels)
- Email (team-specific)
- PagerDuty (critical escalation)

### **Grafana Dashboards**

**Service Dashboards** (21 dashboards):
- `SERVICES/service-*/grafana/dashboards/*-overview.json`

**Centralized Dashboards** (5 dashboards):
- `service-health-dashboard.json`
- `infrastructure-dashboard.json`
- `security-dashboard.json`
- `blockchain-dashboard.json`
- `usc-dashboard.json`

## 🔧 Maintenance

### **Adding New Service Alerts**

1. Generate alert file:
   ```bash
   cd shared/monitoring/templates/
   ./generate-monitoring.sh service-XX-name "Service Name" prefix 90XX 0.1
   ```

2. Integrate alerts:
   ```bash
   ./integrate-service-alerts.sh
   ```

3. Reload Prometheus:
   ```bash
   ./reload-prometheus.sh
   ```

### **Updating Alert Thresholds**

1. Edit alert file in `SERVICES/service-XX/prometheus/alerts/*.yml`
2. Reload Prometheus
3. Test alert firing

### **Updating Dashboards**

1. Edit dashboard JSON
2. Import to Grafana (overwrite existing)
3. Verify panels display correctly

## 🧪 Testing

### **Test Metrics Collection**

```bash
# Test each service endpoint
for port in {9001..9022}; do
    curl http://localhost:$port/metrics | head -5
done
```

### **Test Alert Loading**

```bash
# Check loaded rules
curl http://localhost:9090/api/v1/rules

# Check specific rule group
curl http://localhost:9090/api/v1/rules?rule_group=gateway_critical
```

### **Test Alert Firing**

```bash
# Check active alerts
curl http://localhost:9090/api/v1/alerts

# Check AlertManager alerts
curl http://localhost:9093/api/v2/alerts
```

## 📝 Verification Checklist

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

## 🐛 Troubleshooting

See `TROUBLESHOOTING.md` for detailed troubleshooting guide.

---

**Status**: ✅ Ready for deployment

**Scripts**: `shared/monitoring/templates/`

