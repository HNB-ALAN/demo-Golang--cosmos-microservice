# 📊 Services Configuration for Monitoring Templates

## Service Mapping

Sử dụng file này để generate monitoring configs cho tất cả services.

## Quick Generate Script

```bash
#!/bin/bash
# Generate monitoring configs for all services

cd shared/monitoring/templates/

# Service 02: Auth
./generate-monitoring.sh service-02-auth Auth auth 9002 0.1

# Service 03: User
./generate-monitoring.sh service-03-user User user 9003 0.1

# Service 04: USC Blockchain Core
./generate-monitoring.sh service-04-usc-blockchain-core "USC Blockchain" blockchain 9004 0.15

# Service 05: USC Wallet
./generate-monitoring.sh service-05-usc-wallet "USC Wallet" wallet 9005 0.2

# Service 06: Security
./generate-monitoring.sh service-06-security Security security 9006 0.1

# Service 07: Caching
./generate-monitoring.sh service-07-caching Caching caching 9007 0.05

# Service 09: Social
./generate-monitoring.sh service-09-social Social social 9009 0.1

# Service 10: USC Bilateral Rewards
./generate-monitoring.sh service-10-usc-bilateral-rewards "USC Rewards" rewards 9010 0.2

# Service 11: Content Management
./generate-monitoring.sh service-11-content-management "Content Management" content 9011 0.15

# Service 12: Video Service
./generate-monitoring.sh service-12-video-service Video video 9012 0.2

# Service 13: AI Service
./generate-monitoring.sh service-13-ai-service AI ai 9013 0.3

# Service 14: Commerce Service
./generate-monitoring.sh service-14-commerce-service Commerce commerce 9014 0.2

# Service 15: Notification Service
./generate-monitoring.sh service-15-notification-service Notification notification 9015 0.1

# Service 16: Search Service
./generate-monitoring.sh service-16-search-service Search search 9016 0.2

# Service 17: Analytics Service
./generate-monitoring.sh service-17-analytics-service Analytics analytics 9017 0.3

# Service 18: Moderation Service
./generate-monitoring.sh service-18-moderation-service Moderation moderation 9018 0.2

# Service 19: Recommendation Service
./generate-monitoring.sh service-19-recommendation-service Recommendation recommendation 9019 0.2

# Service 20: Advertising Service
./generate-monitoring.sh service-20-advertising-service Advertising advertising 9020 0.2

# Service 21: Admin Service
./generate-monitoring.sh service-21-admin-service Admin admin 9021 0.1

# Service 22: Kafka Messaging Service
./generate-monitoring.sh service-22-kafka-messaging-service "Kafka Messaging" kafka 9022 0.1
```

## Service Configuration Table

| # | Service ID | Service Name | Metric Prefix | Port | Latency Threshold | Notes |
|---|-----------|--------------|---------------|------|-------------------|-------|
| 01 | service-01-gateway | Gateway | gateway | 9001 | 0.05 | ✅ Already customized |
| 02 | service-02-auth | Auth | auth | 9002 | 0.1 | ✅ Generated |
| 03 | service-03-user | User | user | 9003 | 0.1 | ⏳ Pending |
| 04 | service-04-usc-blockchain-core | USC Blockchain | blockchain | 9004 | 0.15 | ⏳ Pending |
| 05 | service-05-usc-wallet | USC Wallet | wallet | 9005 | 0.2 | ⏳ Pending |
| 06 | service-06-security | Security | security | 9006 | 0.1 | ⏳ Pending |
| 07 | service-07-caching | Caching | caching | 9007 | 0.05 | ⏳ Pending |
| 08 | service-08-monitoring | Monitoring | monitoring | 9008 | 0.1 | ✅ Centralized (has own config) |
| 09 | service-09-social | Social | social | 9009 | 0.1 | ⏳ Pending |
| 10 | service-10-usc-bilateral-rewards | USC Rewards | rewards | 9010 | 0.2 | ⏳ Pending |
| 11 | service-11-content-management | Content Management | content | 9011 | 0.15 | ⏳ Pending |
| 12 | service-12-video-service | Video | video | 9012 | 0.2 | ⏳ Pending |
| 13 | service-13-ai-service | AI | ai | 9013 | 0.3 | ⏳ Pending |
| 14 | service-14-commerce-service | Commerce | commerce | 9014 | 0.2 | ⏳ Pending |
| 15 | service-15-notification-service | Notification | notification | 9015 | 0.1 | ⏳ Pending |
| 16 | service-16-search-service | Search | search | 9016 | 0.2 | ⏳ Pending |
| 17 | service-17-analytics-service | Analytics | analytics | 9017 | 0.3 | ⏳ Pending |
| 18 | service-18-moderation-service | Moderation | moderation | 9018 | 0.2 | ⏳ Pending |
| 19 | service-19-recommendation-service | Recommendation | recommendation | 9019 | 0.2 | ⏳ Pending |
| 20 | service-20-advertising-service | Advertising | advertising | 9020 | 0.2 | ⏳ Pending |
| 21 | service-21-admin-service | Admin | admin | 9021 | 0.1 | ⏳ Pending |
| 22 | service-22-kafka-messaging-service | Kafka Messaging | kafka | 9022 | 0.1 | ⏳ Pending |

## Notes

- **Service-01**: Already has customized dashboards/alerts (keep as-is, use as reference)
- **Service-08**: Has centralized monitoring config (keep as-is)
- **Services 2-22**: Generate from templates (except 01, 08)

## Generate All Script

Save as `generate-all-services.sh`:

```bash
#!/bin/bash
# Generate monitoring configs for all services (except 01, 08)

cd shared/monitoring/templates/

echo "Generating monitoring configs for all services..."
echo ""

# Services 2-7
./generate-monitoring.sh service-02-auth Auth auth 9002 0.1
./generate-monitoring.sh service-03-user User user 9003 0.1
./generate-monitoring.sh service-04-usc-blockchain-core "USC Blockchain" blockchain 9004 0.15
./generate-monitoring.sh service-05-usc-wallet "USC Wallet" wallet 9005 0.2
./generate-monitoring.sh service-06-security Security security 9006 0.1
./generate-monitoring.sh service-07-caching Caching caching 9007 0.05

# Services 9-22
./generate-monitoring.sh service-09-social Social social 9009 0.1
./generate-monitoring.sh service-10-usc-bilateral-rewards "USC Rewards" rewards 9010 0.2
./generate-monitoring.sh service-11-content-management "Content Management" content 9011 0.15
./generate-monitoring.sh service-12-video-service Video video 9012 0.2
./generate-monitoring.sh service-13-ai-service AI ai 9013 0.3
./generate-monitoring.sh service-14-commerce-service Commerce commerce 9014 0.2
./generate-monitoring.sh service-15-notification-service Notification notification 9015 0.1
./generate-monitoring.sh service-16-search-service Search search 9016 0.2
./generate-monitoring.sh service-17-analytics-service Analytics analytics 9017 0.3
./generate-monitoring.sh service-18-moderation-service Moderation moderation 9018 0.2
./generate-monitoring.sh service-19-recommendation-service Recommendation recommendation 9019 0.2
./generate-monitoring.sh service-20-advertising-service Advertising advertising 9020 0.2
./generate-monitoring.sh service-21-admin-service Admin admin 9021 0.1
./generate-monitoring.sh service-22-kafka-messaging-service "Kafka Messaging" kafka 9022 0.1

echo ""
echo "✅ All services generated successfully!"
```

