# 📊 Prometheus Configuration Explanation

## Overview

Có **2 Prometheus config files** trong hệ thống, mỗi file có mục đích khác nhau:

## 1. `shared/monitoring/prometheus.yml` (Reference Config)

**Purpose**: Reference configuration cho templates và documentation

**Location**: `shared/monitoring/prometheus.yml`

**Characteristics**:
- ✅ **Standard ports**: 9001-9022 cho tất cả 22 services
- ✅ **Correct service names**: Đầy đủ và chính xác
- ✅ **All 22 services**: Đầy đủ tất cả services
- ✅ **Infrastructure services**: Postgres, Redis, ClickHouse, etc.
- ✅ **Simple structure**: Dễ đọc và maintain

**Usage**:
- Reference cho templates
- Documentation
- Development
- Testing

**DO NOT USE** trong production Prometheus instance (chỉ dùng để reference)

## 2. `SERVICES/service-08-monitoring/monitoring/prometheus.yml` (Operational Config)

**Purpose**: Operational configuration cho actual Prometheus instance

**Location**: `SERVICES/service-08-monitoring/monitoring/prometheus.yml`

**Characteristics**:
- ✅ **Standard ports**: 9001-9022 cho tất cả 22 services (đã fix)
- ✅ **Correct service names**: Đầy đủ và chính xác (đã fix)
- ✅ **All 22 services**: Đầy đủ tất cả services (đã fix)
- ✅ **Infrastructure monitoring**: Node exporter, cAdvisor, kube-state-metrics
- ✅ **Advanced features**:
  - AlertManager configuration
  - Alert rules (alerting_rules.yml, recording_rules.yml, slo_rules.yml, trading_rules.yml)
  - External labels
  - Blackbox exporter
  - Trading exchange monitoring
  - Storage monitoring (RocksDB, Hazelcast)

**Usage**:
- **Production Prometheus instance**
- Centralized monitoring
- Platform-wide observability

**USE THIS** cho actual Prometheus instance trong production

## Key Differences

| Feature | shared/monitoring | service-08-monitoring |
|---------|-------------------|----------------------|
| **Purpose** | Reference | Operational |
| **Ports** | 9001-9022 | 9001-9022 (fixed) |
| **Service Names** | Correct | Correct (fixed) |
| **Alert Rules** | None | Yes (4 rule files) |
| **AlertManager** | None | Configured |
| **Infrastructure** | Basic | Advanced |
| **External Labels** | None | Yes |
| **Blackbox Exporter** | None | Yes |
| **Trading Monitoring** | None | Yes |

## Port Mapping (Standard)

Tất cả services sử dụng ports **9001-9022** cho metrics:

| Service | Port | Service Name |
|---------|------|--------------|
| service-01-gateway | 9001 | service-01-gateway:9001 |
| service-02-auth | 9002 | service-02-auth:9002 |
| service-03-user | 9003 | service-03-user:9003 |
| service-04-usc-blockchain-core | 9004 | service-04-usc-blockchain-core:9004 |
| service-05-usc-wallet | 9005 | service-05-usc-wallet:9005 |
| service-06-security | 9006 | service-06-security:9006 |
| service-07-caching | 9007 | service-07-caching:9007 |
| service-08-monitoring | 9008 | service-08-monitoring:9008 |
| service-09-social | 9009 | service-09-social:9009 |
| service-10-usc-bilateral-rewards | 9010 | service-10-usc-bilateral-rewards:9010 |
| service-11-content-management | 9011 | service-11-content-management:9011 |
| service-12-video-service | 9012 | service-12-video-service:9012 |
| service-13-ai-service | 9013 | service-13-ai-service:9013 |
| service-14-commerce-service | 9014 | service-14-commerce-service:9014 |
| service-15-notification-service | 9015 | service-15-notification-service:9015 |
| service-16-search-service | 9016 | service-16-search-service:9016 |
| service-17-analytics-service | 9017 | service-17-analytics-service:9017 |
| service-18-moderation-service | 9018 | service-18-moderation-service:9018 |
| service-19-recommendation-service | 9019 | service-19-recommendation-service:9019 |
| service-20-advertising-service | 9020 | service-20-advertising-service:9020 |
| service-21-admin-service | 9021 | service-21-admin-service:9021 |
| service-22-kafka-messaging-service | 9022 | service-22-kafka-messaging-service:9022 |

## Fixes Applied

### ✅ Fixed Issues (2025-11-04)

1. **Ports**: Updated all service ports from `8081, 9091-9110` to `9001-9022`
2. **Service Names**: Fixed all service names to match actual service names:
   - `service-04-blockchain` → `service-04-usc-blockchain-core`
   - `service-05-wallet` → `service-05-usc-wallet`
   - `service-10-rewards` → `service-10-usc-bilateral-rewards`
   - `service-11-content` → `service-11-content-management`
   - `service-12-video` → `service-12-video-service`
   - `service-13-ai` → `service-13-ai-service`
   - `service-14-commerce` → `service-14-commerce-service`
   - `service-15-notification` → `service-15-notification-service`
   - `service-16-search` → `service-16-search-service`
   - `service-17-analytics` → `service-17-analytics-service`
   - `service-18-moderation` → `service-18-moderation-service`
   - `service-19-recommendation` → `service-19-recommendation-service`
   - `service-20-advertising` → `service-20-advertising-service`
   - `service-21-admin` → `service-21-admin-service`
3. **Missing Services**: Added:
   - `service-08-monitoring` (self-monitoring)
   - `service-22-kafka-messaging-service`

## Best Practices

### When to Update

1. **shared/monitoring/prometheus.yml**:
   - Update khi có service mới
   - Update khi thay đổi port standard
   - Update khi có service name changes

2. **service-08-monitoring/monitoring/prometheus.yml**:
   - Update khi có service mới
   - Update infrastructure monitoring
   - Update alert rules
   - Update AlertManager config

### Synchronization

- **Always sync** service ports và names giữa 2 configs
- **Keep** shared/monitoring/prometheus.yml simple (reference)
- **Enhance** service-08-monitoring/monitoring/prometheus.yml với advanced features

## Testing

After updating configs, verify:

```bash
# Check Prometheus config syntax
promtool check config prometheus.yml

# Test scrape targets
curl http://service-01-gateway:9001/metrics
curl http://service-02-auth:9002/metrics
# ... test all services

# Reload Prometheus
curl -X POST http://prometheus:9090/-/reload
```

---

**Last Updated**: 2025-11-04
**Status**: ✅ **All Issues Fixed**

