# 🔄 Auto-Sync Monitoring Files

## Overview

Kể từ bây giờ, các script generate sẽ **tự động sync** alerts và dashboards vào centralized location sau khi generate. Bạn không cần chạy thủ công `copy-service-alerts.sh` và `copy-service-dashboards.sh` nữa!

---

## ✅ Auto-Sync Features

### **1. generate-monitoring.sh**

**Auto-sync mặc định** sau khi generate:
```bash
./generate-monitoring.sh service-02-auth Auth auth 9002 0.1
# ✅ Tự động copy alerts và dashboards vào centralized location
```

**Skip auto-sync** (nếu cần):
```bash
./generate-monitoring.sh service-02-auth Auth auth 9002 0.1 --no-sync
# ⏭️  Chỉ generate, không sync
```

### **2. generate-all-services.sh**

**Auto-sync tự động** sau khi generate tất cả services:
```bash
./generate-all-services.sh
# ✅ Generate tất cả services + auto-sync alerts và dashboards
```

### **3. sync-monitoring.sh** (Manual)

**Chạy thủ công** khi cần sync lại:
```bash
./sync-monitoring.sh
# ✅ Copy tất cả alerts và dashboards vào centralized location
```

---

## 📋 Workflow

### **Generate Service Mới**

```bash
# 1. Generate monitoring configs (auto-sync included)
./generate-monitoring.sh service-XX-new-service "New Service" newservice 90XX 0.1

# 2. Restart Prometheus & Grafana để load configs
docker-compose restart prometheus grafana
```

### **Generate Tất Cả Services**

```bash
# 1. Generate tất cả (auto-sync included)
./generate-all-services.sh

# 2. Restart Prometheus & Grafana
docker-compose restart prometheus grafana
```

### **Sync Manual (Khi Cần)**

```bash
# Sync alerts và dashboards
./sync-monitoring.sh

# Hoặc sync riêng
./copy-service-alerts.sh
./copy-service-dashboards.sh
```

---

## 🎯 Benefits

### **Trước (Manual)**
```bash
# 1. Generate
./generate-monitoring.sh service-02-auth Auth auth 9002

# 2. Sync thủ công
./copy-service-alerts.sh      # ❌ Phải nhớ chạy
./copy-service-dashboards.sh  # ❌ Phải nhớ chạy

# 3. Restart
docker-compose restart prometheus grafana
```

### **Bây Giờ (Auto)**
```bash
# 1. Generate (auto-sync included)
./generate-monitoring.sh service-02-auth Auth auth 9002
# ✅ Auto-sync alerts và dashboards

# 2. Restart
docker-compose restart prometheus grafana
```

---

## 🔧 Technical Details

### **Auto-Sync Logic**

1. **generate-monitoring.sh**:
   - Sau khi generate dashboard, alerts, docs
   - Tự động gọi `sync-monitoring.sh`
   - Có thể skip với flag `--no-sync`

2. **generate-all-services.sh**:
   - Sau khi generate tất cả services
   - Tự động gọi `sync-monitoring.sh` một lần

3. **sync-monitoring.sh**:
   - Gọi `copy-service-alerts.sh`
   - Gọi `copy-service-dashboards.sh`
   - Hiển thị summary

---

## 📝 Notes

- **Auto-sync mặc định**: Bật cho tất cả generate operations
- **Skip sync**: Dùng `--no-sync` flag nếu cần
- **Manual sync**: Dùng `sync-monitoring.sh` khi cần sync lại
- **Restart required**: Vẫn cần restart Prometheus/Grafana sau khi sync

---

**Status**: ✅ **Auto-sync Enabled - No Manual Steps Required!**

