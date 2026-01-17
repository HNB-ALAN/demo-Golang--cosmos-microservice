# 🚀 Production Deployment Guide - Service-04 USC Blockchain Core

**Service**: service-04-usc-blockchain-core  
**Version**: 1.0.0  
**Last Updated**: 2025-11-10

---

## 📋 Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Prerequisites](#prerequisites)
4. [Environment Configuration](#environment-configuration)
5. [Deployment Steps](#deployment-steps)
6. [Health Checks](#health-checks)
7. [Monitoring & Observability](#monitoring--observability)
8. [SLOs & Performance Targets](#slos--performance-targets)
9. [Runbooks](#runbooks)
10. [Disaster Recovery](#disaster-recovery)
11. [Security Checklist](#security-checklist)

---

## 🎯 Overview

Service-04 is the **USC Blockchain Core** service providing blockchain infrastructure for the USC Social Media Platform. It integrates Cosmos SDK v0.53.4 with CometBFT v0.38.19 for consensus.

### Key Features
- **12 gRPC Services**: Block, Transaction, USC Coin, NFT, Smart Contract, Validator, Network, Streaming, Custom Token, Product Certificate, Store Bridge, Store Network
- **14 Custom Modules**: USC, Reward, NFT, Contract, Validator, Network, Bridge, Streaming, Certificate, Store, Token, Block, Performance, Monitoring
- **58 gRPC Methods**: Complete blockchain operations API
- **Hybrid Storage**: Cosmos SDK (RocksDB) + PostgreSQL (analytics)

---

## 🏗️ Architecture

### Deployment Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Load Balancer / Gateway                  │
│                  (Service-01 Gateway)                       │
└──────────────────────┬──────────────────────────────────────┘
                       │
        ┌──────────────┴──────────────┐
        │                            │
┌───────▼────────┐          ┌─────────▼────────┐
│ Service-04-1   │          │ Service-04-2     │
│ (Primary)      │          │ (Secondary)       │
│ Port: 8004     │          │ Port: 8004       │
└───────┬────────┘          └─────────┬────────┘
        │                            │
        └──────────────┬──────────────┘
                       │
        ┌──────────────┴──────────────┐
        │                            │
┌───────▼────────┐          ┌─────────▼────────┐
│ CometBFT-1     │          │ CometBFT-2      │
│ Port: 26657    │          │ Port: 26657     │
└───────┬────────┘          └─────────┬────────┘
        │                            │
        └──────────────┬──────────────┘
                       │
        ┌──────────────┴──────────────┐
        │                            │
┌───────▼────────┐          ┌─────────▼────────┐
│ PostgreSQL     │          │ Redis           │
│ Port: 5432     │          │ Port: 6379      │
└────────────────┘          └────────────────┘
```

### Service Components

1. **Service-04 Application**
   - gRPC Server (Port 8004)
   - Cosmos SDK Integration
   - Business Logic Layer
   - Repository Layer

2. **CometBFT Consensus Engine**
   - RPC Server (Port 26657)
   - P2P Network (Port 26656)
   - ABCI Server (Port 26658)

3. **Storage**
   - RocksDB (Blockchain state)
   - PostgreSQL (Analytics, fallback)
   - Redis (Caching, sessions)

---

## ✅ Prerequisites

### Infrastructure Requirements

- **Kubernetes Cluster** (recommended) or Docker Swarm
- **PostgreSQL 14+** (for analytics and fallback)
- **Redis 6+** (for caching)
- **Persistent Storage** (for RocksDB blockchain data)
- **Network**: Internal service mesh connectivity

### Resource Requirements

#### Minimum (Development)
- **CPU**: 2 cores
- **Memory**: 4 GB
- **Storage**: 50 GB (blockchain data)
- **Network**: 100 Mbps

#### Recommended (Production)
- **CPU**: 4-8 cores
- **Memory**: 16-32 GB
- **Storage**: 500 GB+ (blockchain data, growing)
- **Network**: 1 Gbps

#### High Availability (Production)
- **CPU**: 8-16 cores per instance
- **Memory**: 32-64 GB per instance
- **Storage**: 1 TB+ per instance
- **Instances**: 3+ (for consensus)

---

## 🔧 Environment Configuration

### Required Environment Variables

```bash
# Service Configuration
SERVICE_NAME=service-04
SERVICE_VERSION=1.0.0
ENVIRONMENT=production

# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8004
GRPC_PORT=8004

# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=${DB_PASSWORD}  # From secrets
DB_NAME=blockchain_db
DB_SSLMODE=require  # Production: require SSL

# Redis Configuration
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=${REDIS_PASSWORD}  # From secrets
REDIS_DB=0

# JWT Configuration
JWT_SECRET=${JWT_SECRET}  # From secrets (≥32 chars)
JWT_EXPIRY=24h
REFRESH_EXPIRY=168h
JWT_ISSUER=usc-platform

# Cosmos SDK Configuration
COSMOS_SDK_ENABLED=true
CHAIN_ID=usc-1

# CometBFT Configuration
COMETBFT_RPC_PORT=26657
COMETBFT_P2P_PORT=26656
COMETBFT_ABCI_PORT=26658

# Monitoring
METRICS_ENABLED=true
METRICS_PORT=9004
METRICS_PATH=/metrics

# Logging
LOG_LEVEL=info  # Production: info, warn, error
LOG_FORMAT=json
LOG_OUTPUT=stdout

# Kafka (if enabled)
KAFKA_BROKERS=kafka:9092
KAFKA_GROUP_ID=service-04
```

### Secrets Management

**CRITICAL**: Never hardcode secrets in config files. Use:
- Kubernetes Secrets
- HashiCorp Vault
- AWS Secrets Manager
- Environment variables from secure storage

**Required Secrets**:
- `DB_PASSWORD` - PostgreSQL password
- `REDIS_PASSWORD` - Redis password (if enabled)
- `JWT_SECRET` - JWT signing secret (≥32 characters)

---

## 🚀 Deployment Steps

### Step 1: Prepare Infrastructure

```bash
# 1. Create PostgreSQL database
createdb blockchain_db

# 2. Create Redis instance
# (Managed service or Docker container)

# 3. Prepare persistent storage for RocksDB
# (Kubernetes PVC or Docker volume)
```

### Step 2: Configure Environment

```bash
# 1. Set environment variables
export DB_PASSWORD="your-secure-password"
export JWT_SECRET="your-super-secret-jwt-key-min-32-chars"
export REDIS_PASSWORD="your-redis-password"

# 2. Update config.yaml (or use environment variables)
# See Environment Configuration section
```

### Step 3: Build Docker Image

```bash
cd SERVICES/service-04-usc-blockchain-core

# Build image
docker build -t service-04:1.0.0 .

# Tag for registry
docker tag service-04:1.0.0 registry.example.com/service-04:1.0.0

# Push to registry
docker push registry.example.com/service-04:1.0.0
```

### Step 4: Deploy with Docker Compose

```bash
# Start services
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f service-04
```

### Step 5: Deploy with Kubernetes

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: service-04
  namespace: usc-platform
spec:
  replicas: 3
  selector:
    matchLabels:
      app: service-04
  template:
    metadata:
      labels:
        app: service-04
    spec:
      containers:
      - name: service-04
        image: registry.example.com/service-04:1.0.0
        ports:
        - containerPort: 8004
          name: grpc
        - containerPort: 9004
          name: metrics
        env:
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: service-04-secrets
              key: db-password
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: service-04-secrets
              key: jwt-secret
        volumeMounts:
        - name: blockchain-data
          mountPath: /data/blockchain
        resources:
          requests:
            cpu: "2"
            memory: "4Gi"
          limits:
            cpu: "8"
            memory: "16Gi"
      volumes:
      - name: blockchain-data
        persistentVolumeClaim:
          claimName: service-04-pvc
```

### Step 6: Verify Deployment

```bash
# Health check
grpcurl -plaintext localhost:8004 grpc.health.v1.Health/Check

# Expected response:
# {
#   "status": "SERVING"
# }
```

---

## 🏥 Health Checks

### gRPC Health Check

**Endpoint**: `grpc.health.v1.Health/Check`

```bash
# Check health
grpcurl -plaintext localhost:8004 grpc.health.v1.Health/Check

# Expected: {"status": "SERVING"}
```

### HTTP Health Check (if enabled)

**Endpoint**: `http://localhost:9004/health`

```bash
curl http://localhost:9004/health
```

### Health Check Criteria

Service is healthy if:
- ✅ gRPC server is listening on port 8004
- ✅ PostgreSQL connection is active
- ✅ Redis connection is active (if enabled)
- ✅ Cosmos SDK app is initialized
- ✅ CometBFT is synced

### Kubernetes Liveness Probe

```yaml
livenessProbe:
  exec:
    command:
    - grpcurl
    - -plaintext
    - localhost:8004
    - grpc.health.v1.Health/Check
  initialDelaySeconds: 30
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 3
```

### Kubernetes Readiness Probe

```yaml
readinessProbe:
  exec:
    command:
    - grpcurl
    - -plaintext
    - localhost:8004
    - grpc.health.v1.Health/Check
  initialDelaySeconds: 10
  periodSeconds: 5
  timeoutSeconds: 3
  failureThreshold: 3
```

---

## 📊 Monitoring & Observability

### Metrics Endpoint

**Endpoint**: `http://localhost:9004/metrics`

**Key Metrics**:
- `grpc_server_requests_total` - Total gRPC requests
- `grpc_server_request_duration_seconds` - Request latency
- `blockchain_blocks_produced_total` - Blocks produced
- `blockchain_transactions_total` - Transactions processed
- `database_connections_active` - Active DB connections
- `redis_connections_active` - Active Redis connections

### Prometheus Configuration

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'service-04'
    static_configs:
      - targets: ['service-04:9004']
    metrics_path: '/metrics'
    scrape_interval: 15s
```

### Grafana Dashboards

**Recommended Dashboards**:
1. **Service Overview**
   - Request rate
   - Error rate
   - Latency (p50, p95, p99)
   - Active connections

2. **Blockchain Metrics**
   - Blocks produced per hour
   - Transactions per second
   - Block production time
   - Chain height

3. **Infrastructure**
   - CPU usage
   - Memory usage
   - Disk I/O
   - Network traffic

### Logging

**Log Format**: JSON (structured logging)

**Log Levels**:
- `DEBUG` - Development only
- `INFO` - General information
- `WARN` - Warnings
- `ERROR` - Errors

**Log Aggregation**: ELK Stack, Loki, or CloudWatch

**Correlation IDs**: All logs include `correlation_id` for request tracing

---

## 🎯 SLOs & Performance Targets

### Service Level Objectives (SLOs)

#### Availability
- **Target**: 99.9% uptime
- **Measurement**: Health check success rate
- **Window**: 30 days rolling

#### Latency
- **p50**: <50ms
- **p95**: <100ms
- **p99**: <200ms
- **Measurement**: gRPC request duration

#### Throughput
- **Target**: 10,000 requests/second
- **Measurement**: Requests per second

#### Error Rate
- **Target**: <0.1% error rate
- **Measurement**: gRPC error responses

### Blockchain-Specific SLOs

#### Block Production
- **Block Time**: <5 seconds (configured)
- **Block Production Success Rate**: >99%
- **Transaction Finality**: <10 seconds (2-3 blocks)

#### Transaction Processing
- **Transaction Throughput**: >10,000 TPS
- **Transaction Success Rate**: >99.9%
- **Mempool Size**: <100,000 transactions

### Performance Targets

| Metric | Target | Measurement |
|--------|--------|-------------|
| API Response Time (p95) | <100ms | gRPC request duration |
| Block Production Time | <5s | Block creation to finality |
| Database Query Time | <50ms | PostgreSQL query duration |
| Cache Hit Rate | >95% | Redis cache statistics |
| Error Rate | <0.1% | gRPC error responses |

---

## 📖 Runbooks

### Runbook 1: Service Unavailable

**Symptoms**:
- Health check failing
- gRPC requests timing out
- Service not responding

**Steps**:
1. Check service logs: `docker-compose logs service-04`
2. Check health endpoint: `grpcurl -plaintext localhost:8004 grpc.health.v1.Health/Check`
3. Check database connectivity: `psql -h postgres -U postgres -d blockchain_db`
4. Check Redis connectivity: `redis-cli -h redis ping`
5. Check CometBFT status: `curl http://localhost:26657/status`
6. Restart service if needed: `docker-compose restart service-04`

### Runbook 2: High Latency

**Symptoms**:
- p95 latency >100ms
- Slow response times
- Timeout errors

**Steps**:
1. Check metrics: `curl http://localhost:9004/metrics | grep latency`
2. Check database performance: `SELECT * FROM pg_stat_activity;`
3. Check Redis performance: `redis-cli --latency`
4. Check CometBFT sync status: `curl http://localhost:26657/status`
5. Scale up if needed: Increase CPU/memory resources
6. Optimize queries: Review slow queries in PostgreSQL

### Runbook 3: Database Connection Issues

**Symptoms**:
- Database connection errors
- "too many connections" errors
- Slow database queries

**Steps**:
1. Check connection pool: `SELECT count(*) FROM pg_stat_activity;`
2. Check max connections: `SHOW max_connections;`
3. Increase connection pool in config: `max_conns: 20`
4. Restart service: `docker-compose restart service-04`
5. Monitor connection usage

### Runbook 4: Blockchain Sync Issues

**Symptoms**:
- Blocks not producing
- Chain height not increasing
- Validator not participating

**Steps**:
1. Check CometBFT status: `curl http://localhost:26657/status`
2. Check validator status: `curl http://localhost:26657/validators`
3. Check logs: `docker-compose logs cometbft`
4. Check network connectivity: `ping other-validators`
5. Restart CometBFT if needed: `docker-compose restart cometbft`

### Runbook 5: Memory Leak

**Symptoms**:
- Memory usage continuously increasing
- OOM (Out of Memory) kills
- Service crashes

**Steps**:
1. Check memory usage: `docker stats service-04`
2. Check for goroutine leaks: `curl http://localhost:9004/debug/pprof/goroutine`
3. Review recent code changes
4. Increase memory limits temporarily
5. Restart service: `docker-compose restart service-04`
6. Investigate root cause

---

## 🔄 Disaster Recovery

### Backup Strategy

#### Database Backup
```bash
# Daily PostgreSQL backup
pg_dump -h postgres -U postgres blockchain_db > backup_$(date +%Y%m%d).sql

# Restore from backup
psql -h postgres -U postgres blockchain_db < backup_20251110.sql
```

#### Blockchain State Backup
```bash
# Backup RocksDB data
tar -czf blockchain_data_$(date +%Y%m%d).tar.gz /data/blockchain

# Restore blockchain data
tar -xzf blockchain_data_20251110.tar.gz -C /data/blockchain
```

### Recovery Procedures

#### Full Service Recovery
1. Restore database from backup
2. Restore blockchain state from backup
3. Start CometBFT and sync with network
4. Start Service-04
5. Verify health checks
6. Monitor for 24 hours

#### Partial Recovery (Database Only)
1. Restore database from backup
2. Restart Service-04
3. Service will use database fallback
4. Blockchain will continue from current state

### RTO/RPO Targets

- **RTO (Recovery Time Objective)**: <1 hour
- **RPO (Recovery Point Objective)**: <15 minutes (last backup)

---

## 🔐 Security Checklist

### Pre-Production Security Checklist

- [ ] All secrets moved to environment variables
- [ ] JWT secret ≥32 characters
- [ ] Database SSL enabled (`sslmode: require`)
- [ ] Redis password set (if enabled)
- [ ] TLS/mTLS configured for gRPC
- [ ] Rate limiting enabled
- [ ] Auth interceptor implemented
- [ ] Input validation enabled
- [ ] Error messages sanitized
- [ ] Logging configured (no sensitive data)
- [ ] Security audit completed
- [ ] Dependencies updated (no known vulnerabilities)
- [ ] Network policies configured
- [ ] Firewall rules configured

### Ongoing Security

- [ ] Regular security audits (quarterly)
- [ ] Dependency updates (monthly)
- [ ] Security patches applied promptly
- [ ] Access logs reviewed
- [ ] Anomaly detection enabled

---

## 📝 Additional Resources

### Documentation
- [Architecture Patterns](./ARCHITECTURE_PATTERNS_AND_CONVENTIONS.md)
- [Security Audit Report](./SECURITY_AUDIT_REPORT.md)
- [API Documentation](./api.md)

### Support
- **Service Owner**: USC Platform Team
- **On-Call**: [Contact Information]
- **Documentation**: [Link to Wiki]

---

**Last Updated**: 2025-11-10  
**Version**: 1.0.0  
**Status**: Production Ready ✅


