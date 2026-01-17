#!/bin/bash

# Generate monitoring configs for all services (except 01, 08)
# Service-01: Already customized (keep as-is)
# Service-08: Centralized monitoring (has own config)

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "🚀 Generating monitoring configs for all services..."
echo ""

# Services 2-7
echo "📦 Generating for services 2-7..."
./generate-monitoring.sh service-02-auth Auth auth 9002 0.1
./generate-monitoring.sh service-03-user User user 9003 0.1
./generate-monitoring.sh service-04-usc-blockchain-core "USC Blockchain" blockchain 9004 0.15
./generate-monitoring.sh service-05-usc-wallet "USC Wallet" wallet 9005 0.2
./generate-monitoring.sh service-06-security Security security 9006 0.1
./generate-monitoring.sh service-07-caching Caching caching 9007 0.05

echo ""
echo "📦 Generating for services 9-22..."
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
echo ""

# Auto-sync to centralized location
echo "🔄 Syncing to centralized monitoring..."
if [ -f "${SCRIPT_DIR}/sync-monitoring.sh" ]; then
    "${SCRIPT_DIR}/sync-monitoring.sh"
else
    echo "  ⚠️  Sync script not found, running copy scripts manually..."
    "${SCRIPT_DIR}/copy-service-alerts.sh"
    "${SCRIPT_DIR}/copy-service-dashboards.sh"
fi

echo ""
echo "Note:"
echo "  - Service-01 (Gateway): Already customized (keep as-is)"
echo "  - Service-08 (Monitoring): Centralized monitoring (has own config)"
echo "  - Services 2-22: Generated from templates"
echo ""
echo "Next steps:"
echo "  1. Review generated configs"
echo "  2. Customize service-specific metrics/alerts if needed"
echo "  3. Restart Prometheus & Grafana to load new configs:"
echo "     - docker-compose restart prometheus"
echo "     - docker-compose restart grafana"
echo ""

