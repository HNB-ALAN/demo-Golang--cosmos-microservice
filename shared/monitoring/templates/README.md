# 📋 Monitoring Scripts - Quick Reference

## 🚀 Quick Start (Recommended)

### **Master Script (Simplified)**
```bash
# Use master script for all operations
./monitoring.sh <command> [options]
```

**Examples:**
```bash
# Generate service
./monitoring.sh generate service-02-auth Auth auth 9002 0.1

# Generate all services
./monitoring.sh generate-all

# Sync alerts & dashboards
./monitoring.sh sync

# Validate configs
./monitoring.sh validate

# Test monitoring stack
./monitoring.sh test

# Reload Prometheus
./monitoring.sh reload
```

---

## 🎯 Main Scripts (Sử dụng thường xuyên)

### **1. Generate Monitoring Configs**

#### **Generate một service**
```bash
./generate-monitoring.sh service-XX "Service Name" prefix 90XX 0.1
```
- Tạo dashboard, alerts, docs cho một service
- **Auto-sync** alerts & dashboards vào centralized location
- Dùng flag `--no-sync` để skip auto-sync

#### **Generate tất cả services**
```bash
./generate-all-services.sh
```
- Generate monitoring configs cho services 2-22
- **Auto-sync** sau khi generate xong

### **2. Sync to Centralized Location**

#### **Sync tất cả (Alerts + Dashboards)**
```bash
./sync-monitoring.sh
```
- **Unified script** - sync cả alerts và dashboards
- Gọi khi cần sync lại sau khi edit files

---

## 🔧 Utility Scripts (Dùng khi cần)

### **Copy Scripts (có thể dùng riêng)**

#### **Copy Alerts**
```bash
./copy-service-alerts.sh
```
- Copy alerts từ services → `service-08-monitoring/monitoring/alerts/`
- **Note**: Thường không cần chạy thủ công (auto-sync đã handle)

#### **Copy Dashboards**
```bash
./copy-service-dashboards.sh
```
- Copy dashboards từ services → `service-08-monitoring/monitoring/grafana/provisioning/dashboards/`
- **Note**: Thường không cần chạy thủ công (auto-sync đã handle)

### **Validation & Testing**

#### **Validate Monitoring Config**
```bash
./validate-monitoring-config.sh
```
- Validate Prometheus và AlertManager configs
- Check YAML syntax, PromQL syntax

#### **Test Monitoring Stack**
```bash
./test-monitoring.sh
```
- Test Prometheus, Grafana, AlertManager health
- Test service metrics endpoints

#### **Reload Prometheus**
```bash
./reload-prometheus.sh
```
- Reload Prometheus config (không cần restart)
- Verify loaded rules

#### **Import Dashboards**
```bash
./import-dashboards.sh
```
- Import dashboards vào Grafana via API
- **Note**: Thường không cần (Grafana auto-provisioning)

---

## 📊 Script Categories

### **🟢 Production Scripts (Thường dùng)**
1. `generate-monitoring.sh` - Generate monitoring configs
2. `generate-all-services.sh` - Generate tất cả services
3. `sync-monitoring.sh` - Sync alerts & dashboards

### **🟡 Utility Scripts (Dùng khi cần)**
4. `validate-monitoring-config.sh` - Validate configs
5. `test-monitoring.sh` - Test monitoring stack
6. `reload-prometheus.sh` - Reload Prometheus
7. `import-dashboards.sh` - Import dashboards (rarely used)

### **🔵 Copy Scripts (Internal - Auto-sync sử dụng)**
8. `copy-service-alerts.sh` - Copy alerts (internal)
9. `copy-service-dashboards.sh` - Copy dashboards (internal)

### **⚪ Legacy Scripts (Có thể loại bỏ)**
10. `integrate-service-alerts.sh` - Legacy (replaced by copy-service-alerts.sh)

---

## 🚀 Common Workflows

### **Workflow 1: Generate Service Mới**
```bash
# 1. Generate (auto-sync included)
./generate-monitoring.sh service-XX "Name" prefix 90XX 0.1

# 2. Restart services
docker-compose restart prometheus grafana
```

### **Workflow 2: Generate Tất Cả Services**
```bash
# 1. Generate tất cả (auto-sync included)
./generate-all-services.sh

# 2. Restart services
docker-compose restart prometheus grafana
```

### **Workflow 3: Edit Files và Sync Lại**
```bash
# 1. Edit files trong service directories
vim SERVICES/service-XX/grafana/dashboards/...

# 2. Sync lại
./sync-monitoring.sh

# 3. Restart services
docker-compose restart prometheus grafana
```

### **Workflow 4: Validate & Test**
```bash
# 1. Validate configs
./validate-monitoring-config.sh

# 2. Test monitoring stack
./test-monitoring.sh

# 3. Reload Prometheus (nếu cần)
./reload-prometheus.sh
```

---

## 📝 Script Summary

| Script | Purpose | Frequency | Auto-sync |
|--------|---------|-----------|-----------|
| `generate-monitoring.sh` | Generate một service | High | ✅ Yes |
| `generate-all-services.sh` | Generate tất cả | Medium | ✅ Yes |
| `sync-monitoring.sh` | Sync alerts & dashboards | Medium | N/A |
| `validate-monitoring-config.sh` | Validate configs | Low | ❌ No |
| `test-monitoring.sh` | Test stack | Low | ❌ No |
| `reload-prometheus.sh` | Reload Prometheus | Low | ❌ No |
| `import-dashboards.sh` | Import dashboards | Rare | ❌ No |
| `copy-service-alerts.sh` | Copy alerts (internal) | Auto | N/A |
| `copy-service-dashboards.sh` | Copy dashboards (internal) | Auto | N/A |

---

## 🎯 Quick Decision Tree

**"Tôi cần làm gì?"**

### **Option 1: Master Script (Recommended - Dễ nhất)**
```bash
./monitoring.sh <command>
```

- **Generate monitoring cho service mới?**
  → `./monitoring.sh generate service-XX "Name" prefix 90XX 0.1`

- **Generate cho tất cả services?**
  → `./monitoring.sh generate-all`

- **Edit files và cần sync lại?**
  → `./monitoring.sh sync`

- **Kiểm tra config có đúng không?**
  → `./monitoring.sh validate`

- **Test monitoring stack?**
  → `./monitoring.sh test`

- **Reload Prometheus?**
  → `./monitoring.sh reload`

### **Option 2: Direct Scripts (Nếu muốn dùng trực tiếp)**
- Generate service: `./generate-monitoring.sh`
- Generate all: `./generate-all-services.sh`
- Sync: `./sync-monitoring.sh`
- Validate: `./validate-monitoring-config.sh`
- Test: `./test-monitoring.sh`
- Reload: `./reload-prometheus.sh`

---

## 💡 Tips

1. **Thường chỉ cần 3 scripts chính**:
   - `generate-monitoring.sh` (generate một service)
   - `generate-all-services.sh` (generate tất cả)
   - `sync-monitoring.sh` (sync khi edit files)

2. **Auto-sync enabled**: Generate scripts tự động sync, không cần chạy copy scripts thủ công

3. **Internal scripts**: `copy-service-alerts.sh` và `copy-service-dashboards.sh` được gọi tự động, ít khi cần chạy thủ công

4. **Validation**: Chạy `validate-monitoring-config.sh` trước khi deploy

---

**Status**: ✅ **Simplified - Only 3 main scripts needed!**

