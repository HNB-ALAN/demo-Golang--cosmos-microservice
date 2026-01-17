# BigQuery Credentials Configuration

## 📋 Cấu hình BigQuery cho Service 17 Analytics

### 🔧 Cách cấu hình:

**1. Tạo Service Account trong Google Cloud Console:**
- Vào [Google Cloud Console](https://console.cloud.google.com/)
- Chọn project `usc-analytics`
- Vào IAM & Admin > Service Accounts
- Tạo service account mới với tên `usc-analytics-service`
- Gán role `BigQuery Admin` hoặc `BigQuery Data Editor`

**2. Tạo và download JSON key:**
- Click vào service account vừa tạo
- Vào tab "Keys" > "Add Key" > "Create new key"
- Chọn JSON format và download

**3. Thay thế file credentials:**
- Thay thế nội dung file `service-account-key.json` bằng JSON key thật
- Đảm bảo file có quyền đọc: `chmod 600 service-account-key.json`

### 🚀 Environment Variables:

```yaml
environment:
  - BIGQUERY_PROJECT_ID=usc-analytics
  - BIGQUERY_DATASET=usc_analytics_warehouse
  - GOOGLE_APPLICATION_CREDENTIALS=/app/credentials/service-account-key.json
```

### 📁 Cấu trúc thư mục:

```
SERVICES/
├── credentials/
│   ├── service-account-key.json  # Google Cloud service account key
│   └── README.md                 # Hướng dẫn này
└── docker-compose.yml            # Đã cấu hình volume mount
```

### ⚠️ Lưu ý bảo mật:

- **KHÔNG commit** file `service-account-key.json` vào Git
- Thêm vào `.gitignore`: `SERVICES/credentials/service-account-key.json`
- Sử dụng secrets management trong production
- Rotate keys định kỳ

### 🧪 Test BigQuery connection:

```bash
# Restart service để load credentials mới
docker-compose down service-17-analytics-service
docker-compose up -d service-17-analytics-service

# Kiểm tra logs
docker-compose logs service-17-analytics-service | grep -i bigquery
```

### 📊 BigQuery Datasets sẽ được tạo:

- `usc_analytics_warehouse` - Main analytics warehouse
- `usc_analytics_events` - Event tracking
- `usc_analytics_revenue` - Revenue analytics
- `usc_analytics_performance` - Performance metrics
- `usc_analytics_ml` - ML model analytics
