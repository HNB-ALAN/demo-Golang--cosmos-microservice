# 🚀 Quick Start - Monitoring Templates

## Generate Monitoring Configs

### Single Service

```bash
cd shared/monitoring/templates/

# Generate for a service
./generate-monitoring.sh <service-id> <service-name> <metric-prefix> <metrics-port> [latency-threshold]

# Example
./generate-monitoring.sh service-02-auth Auth auth 9002 0.1
```

### All Services

```bash
cd shared/monitoring/templates/

# Generate for all services (except 01, 08)
./generate-all-services.sh
```

## What Gets Generated

For each service, the script creates:

1. **Dashboard**: `SERVICES/{service-id}/grafana/dashboards/{service-id}-overview.json`
2. **Alerts**: `SERVICES/{service-id}/prometheus/alerts/{service-id}-alerts.yml`
3. **Docs**: `SERVICES/{service-id}/docs/MONITORING_SETUP.md`

## Next Steps

1. **Review generated files** - Customize if needed
2. **Import dashboard** - Upload to Grafana
3. **Add alerts** - Copy to Prometheus and reload
4. **Verify metrics** - Check `/metrics` endpoint

## Service Configuration

See `SERVICES_CONFIGURATION.md` for complete service mapping.

---

**Status**: ✅ Ready to use

