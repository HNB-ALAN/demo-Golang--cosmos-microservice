# 📊 Monitoring Templates Documentation

## Overview

Base templates để generate monitoring configuration cho tất cả USC services. Templates này cung cấp standard dashboards, alerts, và documentation cho mỗi service.

## 📁 Template Files

### 1. `base-dashboard-template.json`
**Purpose**: Base Grafana dashboard template cho tất cả services

**Template Variables**:
- `{{SERVICE_ID}}` - Service identifier (e.g., `service-02-auth`)
- `{{SERVICE_NAME}}` - Service display name (e.g., `Auth`)
- `{{METRIC_PREFIX}}` - Metric name prefix (e.g., `auth`)
- `{{METRICS_PORT}}` - Metrics endpoint port (e.g., `9002`)

**Standard Panels**:
1. Request Rate (QPS) - Requests per second
2. Latency (p50, p95, p99) - Latency percentiles
3. Error Rate - Errors per second
4. Service Health - Service up/down status

**Customization**: Services có thể thêm service-specific panels sau khi generate.

### 2. `base-alerts-template.yml`
**Purpose**: Base Prometheus alert rules template cho tất cả services

**Template Variables**:
- `{{SERVICE_ID}}` - Service identifier
- `{{SERVICE_NAME}}` - Service display name
- `{{METRIC_PREFIX}}` - Metric name prefix
- `{{LATENCY_THRESHOLD}}` - Latency threshold in seconds (default: 0.1)

**Standard Alerts**:
1. `{{SERVICE_NAME}}Down` - Service down alert (30s)
2. `{{SERVICE_NAME}}HighLatency` - High latency alert (5m)
3. `{{SERVICE_NAME}}HighErrorRate` - High error rate alert (2m)
4. `{{SERVICE_NAME}}UptimeBelowTarget` - Uptime below 99.9% (5m)
5. `{{SERVICE_NAME}}LowRequestRate` - Low request rate (10m)

**Customization**: Services có thể thêm service-specific alerts sau khi generate.

### 3. `monitoring-setup-template.md`
**Purpose**: Documentation template cho monitoring setup

**Template Variables**:
- `{{SERVICE_ID}}` - Service identifier
- `{{SERVICE_NAME}}` - Service display name
- `{{METRIC_PREFIX}}` - Metric name prefix
- `{{METRICS_PORT}}` - Metrics endpoint port
- `{{LATENCY_THRESHOLD}}` - Latency threshold

**Contents**:
- Prometheus configuration
- Grafana dashboard setup
- Alert configuration
- Metrics reference
- Prometheus queries
- SLOs
- Troubleshooting

### 4. `generate-monitoring.sh`
**Purpose**: Script để generate monitoring configs từ templates

**Usage**:
```bash
./generate-monitoring.sh <service-id> <service-name> <metric-prefix> <metrics-port> [latency-threshold]
```

**Examples**:
```bash
# Generate for Auth service
./generate-monitoring.sh service-02-auth Auth auth 9002

# Generate for Wallet service with custom latency threshold
./generate-monitoring.sh service-05-usc-wallet "USC Wallet" wallet 9005 0.2
```

**Output**:
- `SERVICES/{service-id}/grafana/dashboards/{service-id}-overview.json`
- `SERVICES/{service-id}/prometheus/alerts/{service-id}-alerts.yml`
- `SERVICES/{service-id}/docs/MONITORING_SETUP.md`

## 🚀 Usage Guide

### Step 1: Generate Base Configs

```bash
cd shared/monitoring/templates/

# Generate for a single service
./generate-monitoring.sh service-02-auth Auth auth 9002

# Generate for multiple services
./generate-monitoring.sh service-03-user User user 9003
./generate-monitoring.sh service-04-usc-blockchain-core "USC Blockchain" blockchain 9004 0.15
./generate-monitoring.sh service-05-usc-wallet "USC Wallet" wallet 9005 0.2
```

### Step 2: Customize (Optional)

Sau khi generate, services có thể customize:

**Dashboard Customization**:
```json
// SERVICES/service-XX/grafana/dashboards/service-XX-overview.json
{
  "panels": [
    // Base panels (from template)
    // ... standard panels ...
    
    // Service-specific panels
    {
      "title": "Custom Metric",
      "expr": "custom_metric_total"
    }
  ]
}
```

**Alerts Customization**:
```yaml
# SERVICES/service-XX/prometheus/alerts/service-XX-alerts.yml
groups:
  - name: service_XX_critical
    rules:
      # Base alerts (from template)
      # ... standard alerts ...
      
      # Service-specific alerts
      - alert: ServiceXXCustomAlert
        expr: custom_metric > threshold
```

### Step 3: Import to Grafana

1. Log into Grafana
2. Navigate to **Dashboards** → **Import**
3. Upload generated dashboard JSON
4. Select Prometheus data source
5. Click **Import**

### Step 4: Add Alerts to Prometheus

1. Copy alert rules to Prometheus:
```bash
cp SERVICES/service-XX/prometheus/alerts/service-XX-alerts.yml /etc/prometheus/alerts/
```

2. Add to `prometheus.yml`:
```yaml
rule_files:
  - "/etc/prometheus/alerts/service-XX-alerts.yml"
```

3. Reload Prometheus:
```bash
curl -X POST http://prometheus:9090/-/reload
```

## 📋 Service Configuration Mapping

### Standard Service Ports

| Service | Service ID | Metric Prefix | Metrics Port | Latency Threshold |
|---------|-----------|---------------|--------------|-------------------|
| Gateway | service-01-gateway | gateway | 9001 | 0.05 (50ms) |
| Auth | service-02-auth | auth | 9002 | 0.1 (100ms) |
| User | service-03-user | user | 9003 | 0.1 (100ms) |
| Blockchain | service-04-usc-blockchain-core | blockchain | 9004 | 0.15 (150ms) |
| Wallet | service-05-usc-wallet | wallet | 9005 | 0.2 (200ms) |
| Security | service-06-security | security | 9006 | 0.1 (100ms) |
| Caching | service-07-caching | caching | 9007 | 0.05 (50ms) |
| Monitoring | service-08-monitoring | monitoring | 9008 | 0.1 (100ms) |
| Social | service-09-social | social | 9009 | 0.1 (100ms) |
| Rewards | service-10-usc-bilateral-rewards | rewards | 9010 | 0.2 (200ms) |
| Content | service-11-content-management | content | 9011 | 0.15 (150ms) |
| Video | service-12-video-service | video | 9012 | 0.2 (200ms) |
| AI | service-13-ai-service | ai | 9013 | 0.3 (300ms) |
| Commerce | service-14-commerce-service | commerce | 9014 | 0.2 (200ms) |
| Notification | service-15-notification-service | notification | 9015 | 0.1 (100ms) |
| Search | service-16-search-service | search | 9016 | 0.2 (200ms) |
| Analytics | service-17-analytics-service | analytics | 9017 | 0.3 (300ms) |
| Moderation | service-18-moderation-service | moderation | 9018 | 0.2 (200ms) |
| Recommendation | service-19-recommendation-service | recommendation | 9019 | 0.2 (200ms) |
| Advertising | service-20-advertising-service | advertising | 9020 | 0.2 (200ms) |
| Admin | service-21-admin-service | admin | 9021 | 0.1 (100ms) |
| Kafka | service-22-kafka-messaging-service | kafka | 9022 | 0.1 (100ms) |

## 🔧 Template Variables Reference

### Required Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `{{SERVICE_ID}}` | Service identifier | `service-02-auth` |
| `{{SERVICE_NAME}}` | Service display name | `Auth` |
| `{{METRIC_PREFIX}}` | Metric name prefix | `auth` |
| `{{METRICS_PORT}}` | Metrics endpoint port | `9002` |

### Optional Variables

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `{{LATENCY_THRESHOLD}}` | Latency threshold in seconds | `0.1` | `0.2` |

## 📊 Standard Metrics Expected

Tất cả services nên expose các metrics sau:

### Required Metrics

1. `{metric_prefix}_requests_total` - Total requests
   - Labels: `method`, `status`
   - Type: Counter

2. `{metric_prefix}_request_duration_seconds` - Request duration
   - Labels: `method`
   - Type: Histogram

3. `{metric_prefix}_errors_total` - Total errors
   - Labels: `error_type`
   - Type: Counter

### Optional Metrics

- `{metric_prefix}_health_status` - Health status (0=down, 1=up)
- `{metric_prefix}_active_connections` - Active connections
- Service-specific metrics

## 🎯 Best Practices

### 1. Template Usage

- ✅ Generate từ templates cho consistency
- ✅ Customize sau khi generate cho service-specific needs
- ✅ Document service-specific metrics/alerts

### 2. Dashboard Customization

- ✅ Giữ standard panels từ template
- ✅ Thêm service-specific panels nếu cần
- ✅ Test dashboard với real metrics

### 3. Alerts Customization

- ✅ Giữ standard alerts từ template
- ✅ Thêm service-specific alerts nếu cần
- ✅ Test alert conditions với Prometheus

### 4. Documentation

- ✅ Update `MONITORING_SETUP.md` với service-specific info
- ✅ Document service-specific metrics
- ✅ Add troubleshooting steps nếu cần

## 🔄 Updating Templates

Khi cần update templates:

1. Update base templates trong `shared/monitoring/templates/`
2. Regenerate cho services cần update:
   ```bash
   ./generate-monitoring.sh service-XX-name "Service Name" prefix port
   ```
3. Services đã customize sẽ cần manual update

## 📚 Related Documentation

- [Monitoring Guide](../MONITORING_GUIDE.MD) - General monitoring guide
- [Prometheus Configuration](../prometheus.yml) - Centralized Prometheus config
- Service-specific monitoring docs trong `SERVICES/service-XX/docs/MONITORING_SETUP.md`

## 🎯 Service-01 Gateway Reference

Service-01 Gateway có customization examples:
- Gateway-specific metrics (GraphQL, gRPC, circuit breakers)
- Extended dashboard panels
- Custom alerts

Xem `SERVICES/service-01-gateway/` để reference customization patterns.

---

**Status**: ✅ **Templates Available** - Base templates sẵn sàng để generate monitoring configs cho tất cả services.

