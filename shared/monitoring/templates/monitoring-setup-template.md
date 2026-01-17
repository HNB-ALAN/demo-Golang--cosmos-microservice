# 📊 Monitoring Setup Guide - {{SERVICE_NAME}}

## Overview

Hướng dẫn setup monitoring cho {{SERVICE_NAME}} service với Prometheus và Grafana dashboards.

## Prerequisites

1. **Prometheus** - Metrics collection và alerting
2. **Grafana** - Visualization và dashboards
3. **AlertManager** - Alert routing và notification

## Prometheus Configuration

### 1. Scrape Configuration

Thêm vào `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: '{{SERVICE_ID}}'
    scrape_interval: 15s
    scrape_timeout: 10s
    metrics_path: '/metrics'
    static_configs:
      - targets: ['{{SERVICE_ID}}:{{METRICS_PORT}}']
        labels:
          service: '{{METRIC_PREFIX}}'
          environment: 'production'
```

### 2. Alert Rules

Copy alert rules vào Prometheus:

```bash
# Copy alert rules
cp SERVICES/{{SERVICE_ID}}/prometheus/alerts/{{SERVICE_ID}}-alerts.yml /etc/prometheus/alerts/

# Add to prometheus.yml
alerting:
  alertmanagers:
    - static_configs:
        - targets:
          - alertmanager:9093
rule_files:
  - "/etc/prometheus/alerts/{{SERVICE_ID}}-alerts.yml"
```

### 3. Reload Prometheus

```bash
# Reload configuration
curl -X POST http://prometheus:9090/-/reload
```

## Grafana Dashboard Setup

### 1. Import Dashboard

1. Log into Grafana
2. Navigate to **Dashboards** → **Import**
3. Upload `grafana/dashboards/{{SERVICE_ID}}-overview.json`
4. Select Prometheus data source
5. Click **Import**

### 2. Dashboard URL

After import, dashboard will be available at:
```
http://grafana:3000/d/{{SERVICE_ID}}-overview
```

### 3. Dashboard Panels

**{{SERVICE_NAME}} Overview Dashboard** includes:

1. **Request Rate (QPS)**
   - Requests per second by method

2. **Latency (p50, p95, p99)**
   - Latency percentiles for requests
   - Target: p95 < {{LATENCY_THRESHOLD}}s

3. **Error Rate**
   - Errors by error type

4. **Service Health**
   - Service health status (1=up, 0=down)

## Alert Configuration

### Critical Alerts

#### 1. {{SERVICE_NAME}}Down
- **Trigger**: `up{job="{{SERVICE_ID}}"} == 0` for 30s
- **Severity**: Critical
- **Action**: Immediate on-call response

#### 2. {{SERVICE_NAME}}HighLatency
- **Trigger**: p95 latency > {{LATENCY_THRESHOLD}}s for 5m
- **Severity**: Warning
- **Action**: Check service health, dependencies

#### 3. {{SERVICE_NAME}}HighErrorRate
- **Trigger**: Error rate > 1% for 2m
- **Severity**: Warning
- **Action**: Check logs for error patterns

### Availability Alerts

#### 4. {{SERVICE_NAME}}UptimeBelowTarget
- **Trigger**: Uptime < 99.9% for 5m
- **Severity**: Critical
- **Action**: Immediate investigation

### Performance Alerts

#### 5. {{SERVICE_NAME}}LowRequestRate
- **Trigger**: Request rate < 1 req/sec for 10m
- **Severity**: Info
- **Action**: Monitor for service degradation

## Metrics Reference

### Standard Metrics

- `{{METRIC_PREFIX}}_requests_total` - Total requests
  - Labels: `method`, `status`
- `{{METRIC_PREFIX}}_request_duration_seconds` - Request duration
  - Labels: `method`
- `{{METRIC_PREFIX}}_errors_total` - Total errors
  - Labels: `error_type`

## Prometheus Queries

### Request Rate
```promql
# Requests per second
sum(rate({{METRIC_PREFIX}}_requests_total[5m])) by (method)
```

### Latency
```promql
# p95 latency
histogram_quantile(0.95, 
  sum(rate({{METRIC_PREFIX}}_request_duration_seconds_bucket[5m])) by (le, method)
)
```

### Error Rate
```promql
# Error rate
sum(rate({{METRIC_PREFIX}}_errors_total[5m])) by (error_type)
/
sum(rate({{METRIC_PREFIX}}_requests_total[5m]))
```

### Service Health
```promql
# Service health status
up{job="{{SERVICE_ID}}"}
```

## SLOs (Service Level Objectives)

### Availability
- **Target**: 99.99% uptime
- **Measurement**: `up{job="{{SERVICE_ID}}"}`
- **Alert**: Uptime < 99.9% for 5 minutes

### Latency
- **p95**: < {{LATENCY_THRESHOLD}}s
- **Alert**: p95 > threshold for 5 minutes

### Error Rate
- **Target**: < 0.1%
- **Measurement**: Errors / Total requests
- **Alert**: Error rate > 1% for 2 minutes

## Troubleshooting

### Dashboard Not Loading

1. Check Prometheus data source in Grafana
2. Verify metrics endpoint: `curl http://localhost:{{METRICS_PORT}}/metrics`
3. Check Prometheus scrape configuration

### Alerts Not Firing

1. Verify alert rules are loaded: `curl http://prometheus:9090/api/v1/rules`
2. Check AlertManager configuration
3. Verify notification channels are configured

### Metrics Missing

1. Check service metrics endpoint: `curl http://localhost:{{METRICS_PORT}}/metrics`
2. Verify Prometheus is scraping service
3. Check service logs for metric collection errors

---

**Status**: ✅ **Generated from template** - This monitoring setup was generated from base templates. Customize as needed for service-specific requirements.

