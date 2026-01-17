# 📘 Alert Runbooks

**Service-08: Monitoring Service** - Alert response procedures and resolution steps.

## 🔴 Critical Alerts

### **ServiceDown**

**Alert**: `ServiceDown`  
**Severity**: Critical  
**Description**: Service has been down for more than 1 minute

**Resolution Steps**:
1. Check service status: `docker ps | grep service-XX` or `kubectl get pods | grep service-XX`
2. Check service logs: `docker logs service-XX` or `kubectl logs service-XX`
3. Check resource usage: CPU, memory, disk
4. Check dependencies: Database, cache, message queue
5. Restart service if needed
6. Escalate if issue persists

**Service-Specific**:
- **Gateway**: Check all backend services, circuit breakers
- **Auth**: Check database, JWT service
- **Blockchain**: Check network, validators
- **Wallet**: Check blockchain connection, database

---

### **GatewayDown**

**Alert**: `GatewayDown`  
**Severity**: Critical  
**Description**: Gateway service is down - entire platform affected

**Immediate Actions**:
1. **Check service status**: `docker ps | grep service-01-gateway`
2. **Check logs**: `docker logs service-01-gateway --tail 100`
3. **Check dependencies**:
   - Redis connection
   - Kafka connection
   - Database connection
4. **Check resource limits**: CPU, memory
5. **Restart service**: `docker restart service-01-gateway`
6. **Escalate to on-call**: If not resolved in 5 minutes

**Post-Incident**:
- Root cause analysis
- Update monitoring thresholds
- Document incident

---

### **GatewayHighLatency**

**Alert**: `GatewayHighLatency`  
**Severity**: Warning  
**Description**: Gateway p95 latency exceeds threshold (GraphQL >50ms or gRPC >100ms)

**Resolution Steps**:
1. Check current latency: Grafana dashboard
2. Identify slow operations: Check logs for slow queries
3. Check backend services: Are downstream services slow?
4. Check circuit breakers: Are any services failing?
5. Check cache hit rate: Is caching working?
6. Check database queries: Are queries optimized?
7. Scale if needed: Add more instances

**Thresholds**:
- GraphQL: >50ms p95
- gRPC: >100ms p95

---

### **GatewayHighErrorRate**

**Alert**: `GatewayHighErrorRate`  
**Severity**: Warning  
**Description**: Gateway error rate exceeds 1% for 2 minutes

**Resolution Steps**:
1. Check error types: Grafana dashboard
2. Check backend services: Are services returning errors?
3. Check authentication: Are auth failures increasing?
4. Check rate limiting: Are requests being rate limited?
5. Check database: Are queries failing?
6. Check external services: Are third-party services down?

**Common Causes**:
- Backend service failures
- Database connection issues
- Authentication problems
- Invalid requests

---

### **GatewayCircuitBreakerOpen**

**Alert**: `GatewayCircuitBreakerOpen`  
**Severity**: Warning  
**Description**: Circuit breaker is open for a service

**Resolution Steps**:
1. Identify affected service: Check alert labels
2. Check service health: Is the service down?
3. Check service logs: What errors are occurring?
4. Check service metrics: Response time, error rate
5. Fix service issue: Resolve root cause
6. Wait for circuit breaker to close: Usually auto-closes after service recovers

**Circuit Breaker States**:
- **Closed**: Normal operation
- **Open**: Failing fast, not calling service
- **Half-Open**: Testing if service recovered

---

## 🟡 Warning Alerts

### **ServiceHighLatency**

**Alert**: `ServiceHighLatency`  
**Severity**: Warning  
**Description**: Service p95 latency exceeds threshold

**Resolution Steps**:
1. Check service performance: Grafana dashboard
2. Check database queries: Slow queries?
3. Check external dependencies: External API calls slow?
4. Check resource usage: CPU, memory
5. Optimize code: Identify bottlenecks
6. Scale horizontally: Add more instances

**Thresholds** (varies by service):
- Gateway: 0.05s
- Auth: 0.1s
- Blockchain: 0.15s
- Wallet: 0.2s
- Video: 0.2s
- AI: 0.3s

---

### **ServiceHighErrorRate**

**Alert**: `ServiceHighErrorRate`  
**Severity**: Warning  
**Description**: Service error rate exceeds 1%

**Resolution Steps**:
1. Check error logs: What errors are occurring?
2. Check error types: 4xx vs 5xx errors
3. Check input validation: Invalid requests?
4. Check dependencies: Downstream services failing?
5. Check database: Connection issues?
6. Fix root cause: Address the issue

---

### **ServiceUptimeBelowTarget**

**Alert**: `ServiceUptimeBelowTarget`  
**Severity**: Critical  
**Description**: Service uptime below 99.9% target

**Resolution Steps**:
1. Check service availability: How often is it down?
2. Check incident history: What caused downtime?
3. Improve reliability: Fix recurring issues
4. Add redundancy: Multiple instances
5. Improve monitoring: Better alerting

---

## 🔵 Info Alerts

### **ServiceLowRequestRate**

**Alert**: `ServiceLowRequestRate`  
**Severity**: Info  
**Description**: Service receiving less than 1 request/second for 10 minutes

**Resolution Steps**:
1. Check if this is expected: Off-peak hours?
2. Check service health: Is service still working?
3. Check upstream services: Are they calling this service?
4. Check routing: Is traffic being routed correctly?

**Note**: This may be normal during off-peak hours.

---

## 📋 Service-Specific Runbooks

### **Gateway Service**

- **Circuit Breaker**: Check backend services
- **Connection Count**: Check if approaching limits
- **GraphQL Errors**: Check query complexity, validation
- **gRPC Errors**: Check service connectivity

### **Auth Service**

- **Authentication Failures**: Check user credentials, JWT
- **Rate Limiting**: Check if being attacked
- **Database Issues**: Check connection pool

### **Blockchain Service**

- **Network Issues**: Check validator status
- **Transaction Failures**: Check network congestion
- **Block Production**: Check validator health

### **Wallet Service**

- **Balance Mismatches**: Check blockchain sync
- **Transaction Failures**: Check blockchain connection
- **High Latency**: Check blockchain network

---

## 🚨 Escalation Procedures

### **Level 1: On-Call Engineer**
- Responds within 15 minutes
- Attempts basic resolution
- Escalates if not resolved in 30 minutes

### **Level 2: Senior Engineer**
- Responds within 5 minutes
- Deep technical investigation
- Escalates if not resolved in 1 hour

### **Level 3: Engineering Manager**
- Responds immediately
- Coordinates resources
- Makes critical decisions

### **Critical Alerts**
- Gateway down
- All services down
- Data breach suspected
- Financial system failure

---

## 📝 Post-Incident Actions

After resolving an alert:

1. **Document incident**: What happened, root cause, resolution
2. **Update runbook**: Add new resolution steps
3. **Review monitoring**: Adjust thresholds if needed
4. **Improve resilience**: Prevent recurrence
5. **Share learnings**: Team knowledge sharing

---

**Status**: ✅ Complete

**Location**: `shared/monitoring/ALERT_RUNBOOKS.md`

