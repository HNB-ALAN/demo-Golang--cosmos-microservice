# 🔧 Monitoring Troubleshooting Guide

**Service-08: Monitoring Service** - Common issues and solutions.

## 🐛 Common Issues

### **1. Prometheus Not Scraping Services**

**Symptoms**:
- No metrics in Prometheus
- Targets show as "down"
- Empty query results

**Diagnosis**:
```bash
# Check Prometheus targets
curl http://localhost:9090/api/v1/targets

# Check service metrics endpoint
curl http://localhost:9001/metrics
```

**Solutions**:
1. **Service not running**: Start the service
2. **Wrong port**: Check service exposes metrics on correct port
3. **Network issue**: Check connectivity between Prometheus and service
4. **Firewall**: Check firewall rules
5. **Service discovery**: Check Prometheus scrape config

---

### **2. Rules Not Loading**

**Symptoms**:
- Rules not visible in Prometheus UI
- Alerts not firing
- Error in Prometheus logs

**Diagnosis**:
```bash
# Check loaded rules
curl http://localhost:9090/api/v1/rules

# Check Prometheus logs
docker logs prometheus | grep -i error
```

**Solutions**:
1. **Syntax error**: Validate with `promtool check rules`
2. **File path wrong**: Check `prometheus.yml` rule_files paths
3. **File permissions**: Ensure Prometheus can read files
4. **YAML format**: Check YAML syntax
5. **Reload Prometheus**: `curl -X POST http://localhost:9090/-/reload`

---

### **3. Alerts Not Firing**

**Symptoms**:
- Alerts defined but not firing
- Conditions met but no alert

**Diagnosis**:
```bash
# Check active alerts
curl http://localhost:9090/api/v1/alerts

# Check rule evaluation
curl http://localhost:9090/api/v1/rules
```

**Solutions**:
1. **Rule expression wrong**: Check PromQL expression
2. **Metrics not available**: Check if metrics exist
3. **Threshold too high**: Adjust threshold
4. **For duration not met**: Check `for` duration
5. **AlertManager not connected**: Check AlertManager config

---

### **4. AlertManager Not Receiving Alerts**

**Symptoms**:
- Alerts in Prometheus but not in AlertManager
- No notifications sent

**Diagnosis**:
```bash
# Check AlertManager alerts
curl http://localhost:9093/api/v2/alerts

# Check AlertManager config
curl http://localhost:9093/api/v1/status
```

**Solutions**:
1. **AlertManager not connected**: Check Prometheus alerting config
2. **Config error**: Validate with `amtool check-config`
3. **Network issue**: Check connectivity
4. **AlertManager down**: Check service status

---

### **5. Dashboards Not Showing Data**

**Symptoms**:
- Dashboard empty
- "No data" message
- Panels not loading

**Diagnosis**:
```bash
# Check Grafana data source
curl http://localhost:3000/api/datasources

# Check Prometheus query
curl "http://localhost:9090/api/v1/query?query=up"
```

**Solutions**:
1. **Data source not configured**: Configure Prometheus data source
2. **Query wrong**: Check PromQL query in dashboard
3. **Time range**: Check time range selection
4. **Metrics not available**: Check if metrics exist
5. **Permissions**: Check Grafana permissions

---

### **6. Service Metrics Not Available**

**Symptoms**:
- `/metrics` endpoint returns 404
- No metrics in Prometheus
- Service running but no metrics

**Diagnosis**:
```bash
# Test metrics endpoint
curl http://localhost:9001/metrics

# Check service logs
docker logs service-01-gateway | grep -i metric
```

**Solutions**:
1. **Metrics not exposed**: Add metrics endpoint to service
2. **Wrong port**: Check service exposes metrics on correct port
3. **Path wrong**: Check metrics path (usually `/metrics`)
4. **Service not instrumented**: Add Prometheus client library
5. **Service not started**: Start the service

---

### **7. High Memory Usage in Prometheus**

**Symptoms**:
- Prometheus using too much memory
- OOM kills
- Slow queries

**Solutions**:
1. **Reduce retention**: Lower retention period
2. **Increase memory**: Allocate more memory
3. **Optimize queries**: Use recording rules
4. **Reduce scrape interval**: Increase scrape interval
5. **Use remote storage**: Offload to remote storage

---

### **8. Alert Fatigue**

**Symptoms**:
- Too many alerts
- Same alert firing repeatedly
- Alerts not actionable

**Solutions**:
1. **Tune thresholds**: Adjust alert thresholds
2. **Group alerts**: Use AlertManager grouping
3. **Inhibit rules**: Suppress redundant alerts
4. **Silence alerts**: Temporarily silence noisy alerts
5. **Review alerts**: Remove unnecessary alerts

---

## 🔍 Diagnostic Commands

### **Prometheus**

```bash
# Check health
curl http://localhost:9090/-/healthy

# Check config
curl http://localhost:9090/api/v1/status/config

# Check targets
curl http://localhost:9090/api/v1/targets

# Check rules
curl http://localhost:9090/api/v1/rules

# Check alerts
curl http://localhost:9090/api/v1/alerts

# Query metrics
curl "http://localhost:9090/api/v1/query?query=up"
```

### **AlertManager**

```bash
# Check health
curl http://localhost:9093/-/healthy

# Check config
curl http://localhost:9093/api/v1/status

# Check alerts
curl http://localhost:9093/api/v2/alerts

# Check silences
curl http://localhost:9093/api/v2/silences
```

### **Grafana**

```bash
# Check health
curl http://localhost:3000/api/health

# Check data sources
curl http://localhost:3000/api/datasources

# Check dashboards
curl http://localhost:3000/api/search?type=dash-db
```

---

## 📊 Performance Tuning

### **Prometheus**

1. **Scrape Interval**: Adjust based on service needs
   - Critical services: 10s
   - Standard services: 15s
   - Infrastructure: 30s

2. **Retention**: Balance storage vs retention
   - High-resolution: 1 month
   - Aggregated: 1 year

3. **Recording Rules**: Pre-compute expensive queries

### **AlertManager**

1. **Grouping**: Group related alerts
2. **Throttling**: Prevent alert spam
3. **Inhibition**: Suppress redundant alerts

### **Grafana**

1. **Refresh Interval**: Adjust based on needs
2. **Query Optimization**: Use recording rules
3. **Dashboard Organization**: Group related dashboards

---

## 📝 Support

For additional support:
- Check Prometheus logs: `docker logs prometheus`
- Check AlertManager logs: `docker logs alertmanager`
- Check Grafana logs: `docker logs grafana`
- Review documentation: `README.md`
- Contact monitoring team

---

**Status**: ✅ Complete

**Location**: `shared/monitoring/TROUBLESHOOTING.md`

