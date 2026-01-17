# 📊 USC Social Platform - Monitoring System Guide

## 🎯 Overview

The USC Social Platform monitoring system provides comprehensive observability across all 22 microservices using Prometheus, Grafana, and custom Go monitoring packages.

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Services      │    │   Prometheus    │    │    Grafana      │
│   (22 services) │───▶│   (Metrics)     │───▶│  (Dashboards)   │
│                 │    │                 │    │                 │
│ Port 9001-9022  │    │ Port 9091       │    │ Port 3000       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 📁 File Structure

```
shared/monitoring/
├── monitoring.go           # Go package for monitoring functionality
├── prometheus.yml          # Prometheus configuration (reference)
├── grafana/
│   ├── datasources/
│   │   └── prometheus.yml  # Grafana datasource config
│   └── dashboards/         # Dashboard definitions
│       └── usc-dashboard.json  # Platform-wide dashboard
├── templates/              # 🆕 NEW - Monitoring templates
│   ├── base-dashboard-template.json
│   ├── base-alerts-template.yml
│   ├── monitoring-setup-template.md
│   ├── generate-monitoring.sh
│   └── README_TEMPLATES.md
└── MONITORING_GUIDE.md     # This guide
```

## 🚀 Quick Start

### 1. Start Monitoring Infrastructure

```bash
# Start Prometheus and Grafana
docker-compose up -d prometheus grafana

# Verify services are running
curl http://localhost:9091  # Prometheus
curl http://localhost:3000  # Grafana (admin/admin)
```

### 2. Generate Service Monitoring Configs

```bash
# Generate monitoring configs from templates
cd shared/monitoring/templates/

# Generate for a service
./generate-monitoring.sh service-02-auth Auth auth 9002

# See README_TEMPLATES.md for full documentation
```

### 3. Start Services with Monitoring

```bash
# Start services that have monitoring enabled
docker-compose up -d service-01-gateway service-06-security service-08-monitoring service-21-admin-service
```

## 📊 Services with Monitoring Integration

### Tier 4 Monitoring Services (4/22):

| Service | Role | Features |
|---------|------|----------|
| **service-08-monitoring** | Primary Monitoring Hub | System metrics collection, alerting |
| **service-21-admin-service** | Admin Dashboards | Dashboard management, visualization |
| **service-01-gateway** | Request Tracing | Distributed tracing, performance monitoring |
| **service-06-security** | Security Monitoring | Threat detection, security alerts |

## 🔧 Configuration

### Prometheus Configuration (`prometheus.yml`)

```yaml
# Scrape all 22 microservices
- job_name: 'usc-gateway'
  static_configs:
    - targets: ['service-01-gateway:9001']
  scrape_interval: 5s

# Infrastructure monitoring
- job_name: 'postgres'
  static_configs:
    - targets: ['postgres:5432']
```

### Grafana Datasource (`grafana/datasources/prometheus.yml`)

```yaml
datasources:
  - name: Prometheus
    type: prometheus
    url: http://prometheus:9091
    isDefault: true
```

## 💻 Go Code Usage

### 1. Basic Metrics Collection

```go
import "github.com/usc-platform/shared/monitoring"

// Initialize monitoring service
monitoringService, err := monitoring.NewMonitoringService(cfg, logger)
if err != nil {
    log.Fatal("Failed to create monitoring service:", err)
}

// Record basic metrics
monitoringService.RecordRequest(latency, isError)
monitoringService.RecordCacheOperation(isHit)
monitoringService.RecordDatabaseQuery(duration)
```

### 2. Advanced Monitoring Operations

```go
// Send metrics to Prometheus
err := monitoringService.SendToPrometheus(ctx, metrics)

// Create alerts
err := monitoringService.CreateAlert(ctx, 
    "High Error Rate", 
    "Error rate exceeds threshold", 
    "critical", 
    0.05, 
    "error_rate > 0.05")

// Distributed tracing
trace, err := monitoringService.StartTrace(ctx, 
    "request-123", 
    "service-name", 
    "operation")
defer monitoringService.FinishTrace(ctx, trace.TraceID, tags)
```

### 3. Service-Specific Monitoring

#### Primary Monitoring Hub (service-08):

```go
// Collect system-wide metrics
func (pmh *PrimaryMonitoringHandler) CollectSystemMetrics(ctx context.Context) error {
    metrics := pmh.service.GetAllMetrics()
    return pmh.service.SendToPrometheus(ctx, metrics)
}

// Generate comprehensive reports
report, err := pmh.service.GenerateReport(ctx, startTime, endTime)
```

#### Admin Service (service-21):

```go
// Create monitoring dashboards
func (mdh *MonitoringDashboardHandler) CreateDashboard(ctx context.Context, name, description string) error {
    // Creates Grafana dashboard
    return mdh.service.CreateDashboard(ctx, name, description)
}
```

## 📈 Monitoring Features

### 1. Metrics Collection

- **Request Latency**: Response time tracking
- **Cache Performance**: Hit/miss ratios
- **Database Queries**: Query execution times
- **Error Rates**: Failure tracking
- **Throughput**: Requests per second

### 2. Distributed Tracing

- **Request Flow**: Track requests across services
- **Service Dependencies**: Map service interactions
- **Performance Analysis**: Identify bottlenecks
- **Error Propagation**: Trace error sources

### 3. Alerting System

- **Threshold Alerts**: Configurable thresholds
- **Severity Levels**: Critical, warning, info
- **Condition-based**: Complex alert conditions
- **Real-time Notifications**: Immediate alerts

### 4. Dashboard Management

- **Performance Dashboards**: Service performance metrics
- **Infrastructure Dashboards**: Database, cache, queue status
- **Security Dashboards**: Threat detection, access logs
- **Custom Dashboards**: Service-specific visualizations

## 🔍 Monitoring Endpoints

### Service Metrics Endpoints

| Service | Metrics Endpoint | Port |
|---------|------------------|------|
| service-01-gateway | `/metrics` | 9001 |
| service-06-security | `/metrics` | 9006 |
| service-08-monitoring | `/metrics` | 9008 |
| service-21-admin-service | `/metrics` | 9021 |

### Infrastructure Endpoints

| Service | Endpoint | Port |
|---------|----------|------|
| Prometheus | `/` | 9091 |
| Grafana | `/` | 3000 |
| PostgreSQL | `/metrics` | 5432 |
| Redis | `/metrics` | 6379 |

## 🚨 Alerting Examples

### 1. High Error Rate Alert

```go
// Create alert for high error rate
err := monitoringService.CreateAlert(ctx,
    "High Error Rate",
    "Service error rate exceeds 5%",
    "critical",
    0.05,
    "error_rate > 0.05")
```

### 2. High Latency Alert

```go
// Create alert for high latency
err := monitoringService.CreateAlert(ctx,
    "High Latency",
    "Service response time exceeds 1 second",
    "warning",
    1000.0,
    "avg_latency > 1000")
```

### 3. Security Threat Alert

```go
// Create security alert
err := monitoringService.CreateAlert(ctx,
    "Security Threat",
    "Multiple failed authentication attempts",
    "critical",
    5.0,
    "failed_auth_attempts > 5")
```

## 📊 Dashboard Examples

### 1. Performance Dashboard

- **Request Rate**: Requests per second
- **Response Time**: Average, P50, P95, P99 latencies
- **Error Rate**: Percentage of failed requests
- **Throughput**: Data processed per second

### 2. Infrastructure Dashboard

- **Database Connections**: Active connections
- **Cache Hit Rate**: Cache performance
- **Queue Depth**: Message queue status
- **Resource Usage**: CPU, memory, disk

### 3. Security Dashboard

- **Authentication Failures**: Failed login attempts
- **Suspicious Activity**: Unusual access patterns
- **Threat Detection**: Security alerts
- **Access Logs**: User activity tracking

## 🔧 Troubleshooting

### Common Issues

#### 1. Prometheus Not Scraping Metrics

```bash
# Check Prometheus targets
curl http://localhost:9091/api/v1/targets

# Verify service metrics endpoint
curl http://service-01-gateway:9001/metrics
```

#### 2. Grafana Not Connecting to Prometheus

```bash
# Check Grafana datasource
curl http://localhost:3000/api/datasources

# Verify Prometheus is accessible
curl http://prometheus:9091/api/v1/query?query=up
```

#### 3. Services Not Sending Metrics

```go
// Check monitoring service initialization
monitoringService, err := monitoring.NewMonitoringService(cfg, logger)
if err != nil {
    log.Printf("Monitoring service error: %v", err)
}

// Verify metrics are being recorded
metrics := monitoringService.GetAllMetrics()
log.Printf("Current metrics: %+v", metrics)
```

### Debug Commands

```bash
# Check Prometheus configuration
docker exec prometheus cat /etc/prometheus/prometheus.yml

# View Prometheus logs
docker logs prometheus

# Check Grafana configuration
docker exec grafana cat /etc/grafana/grafana.ini

# View Grafana logs
docker logs grafana
```

## 📚 Best Practices

### 1. Metrics Naming

- Use consistent naming conventions
- Include service name in metric names
- Use descriptive labels

```go
// Good
monitoringService.RecordRequest("gateway_request_duration", latency, map[string]string{
    "service": "gateway",
    "endpoint": "/api/v1/users",
    "method": "GET",
})

// Bad
monitoringService.RecordRequest("req_time", latency, nil)
```

### 2. Alert Thresholds

- Set realistic thresholds
- Use different severity levels
- Test alert conditions

```go
// Good - realistic thresholds
monitoringService.CreateAlert(ctx, "High Error Rate", "Error rate > 1%", "warning", 0.01, "error_rate > 0.01")
monitoringService.CreateAlert(ctx, "Critical Error Rate", "Error rate > 5%", "critical", 0.05, "error_rate > 0.05")

// Bad - unrealistic thresholds
monitoringService.CreateAlert(ctx, "Any Error", "Error rate > 0%", "critical", 0.0, "error_rate > 0")
```

### 3. Dashboard Design

- Keep dashboards focused
- Use appropriate visualizations
- Include time ranges
- Add descriptions

### 4. Performance Considerations

- Don't over-instrument
- Use sampling for high-volume metrics
- Batch metric updates
- Monitor monitoring overhead

## 🔄 Maintenance

### Regular Tasks

1. **Review Alert Thresholds**: Monthly review of alert effectiveness
2. **Dashboard Updates**: Quarterly dashboard improvements
3. **Metric Cleanup**: Remove unused metrics
4. **Performance Tuning**: Optimize monitoring overhead

### Backup and Recovery

```bash
# Backup Prometheus data
docker exec prometheus tar -czf /backup/prometheus-$(date +%Y%m%d).tar.gz /prometheus

# Backup Grafana configuration
docker exec grafana tar -czf /backup/grafana-$(date +%Y%m%d).tar.gz /var/lib/grafana
```

## 📞 Support

For monitoring system issues:

1. Check this guide first
2. Review service logs
3. Verify configuration files
4. Test individual components
5. Contact development team

## 📊 Monitoring Templates

### Generate Service Monitoring Configs

Use templates để generate monitoring configs cho services:

```bash
cd shared/monitoring/templates/

# Generate for single service
./generate-monitoring.sh service-02-auth Auth auth 9002

# Generate for all services
./generate-all-services.sh
```

### Template Documentation

- **Templates Guide**: [templates/README_TEMPLATES.md](templates/README_TEMPLATES.md)
- **Implementation**: [templates/MONITORING_TEMPLATES_IMPLEMENTATION.md](templates/MONITORING_TEMPLATES_IMPLEMENTATION.md)
- **Service Config**: [templates/SERVICES_CONFIGURATION.md](templates/SERVICES_CONFIGURATION.md)

### Template Files

- `base-dashboard-template.json` - Base Grafana dashboard
- `base-alerts-template.yml` - Base Prometheus alerts
- `monitoring-setup-template.md` - Documentation template
- `generate-monitoring.sh` - Generation script

## 🔗 Related Documentation

- [Shared Library Documentation](../README.md)
- [Service Architecture Guide](../../SERVICES/SHARED-INTEGRATION-WORKFLOW.md)
- [Deployment Guide](../../SERVICES/DEPLOYMENT-WORKFLOW.md)
- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
- [Monitoring Templates](../../shared/monitoring/templates/README_TEMPLATES.md)

---

**Last Updated**: $(date)
**Version**: 1.1 (Added Templates)
**Maintainer**: USC Development Team
