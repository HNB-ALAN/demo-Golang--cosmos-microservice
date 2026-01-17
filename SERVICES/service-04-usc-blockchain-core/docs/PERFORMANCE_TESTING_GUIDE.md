# ⚡ Performance Testing Guide - Service-04 USC Blockchain Core

**Service**: service-04-usc-blockchain-core  
**Version**: 1.0.0  
**Last Updated**: 2025-11-10

---

## 📋 Overview

Comprehensive performance testing guide for Service-04 including load testing, stress testing, and benchmarking.

---

## 🎯 Performance Targets

### API Performance

| Metric | Target | Measurement |
|--------|--------|-------------|
| **p50 Latency** | <50ms | gRPC request duration |
| **p95 Latency** | <100ms | gRPC request duration |
| **p99 Latency** | <200ms | gRPC request duration |
| **Throughput** | 10,000 req/s | Requests per second |
| **Error Rate** | <0.1% | Error responses / Total requests |

### Blockchain Performance

| Metric | Target | Measurement |
|--------|--------|-------------|
| **Block Production Time** | <5s | Time to produce block |
| **Transaction Throughput** | 10,000+ TPS | Transactions per second |
| **Block Finality** | <10s | Time to finality (2-3 blocks) |
| **Transaction Submission** | <100ms p95 | Submit transaction latency |

### Infrastructure Performance

| Metric | Target | Measurement |
|--------|--------|-------------|
| **Database Query Time** | <50ms p95 | PostgreSQL query duration |
| **Cache Hit Rate** | >95% | Redis cache statistics |
| **Connection Pool Usage** | <80% | Active connections / Max |

---

## 🛠️ Testing Tools

### Recommended Tools

1. **ghz** - gRPC load testing tool
   ```bash
   go install github.com/bojand/ghz/cmd/ghz@latest
   ```

2. **k6** - Modern load testing tool
   ```bash
   # Install k6
   curl https://github.com/grafana/k6/releases/download/v0.47.0/k6-v0.47.0-linux-amd64.tar.gz -L | tar xvz
   ```

3. **Apache Bench (ab)** - HTTP load testing
   ```bash
   # Usually pre-installed
   ab -n 10000 -c 100 http://localhost:9004/metrics
   ```

4. **wrk** - HTTP benchmarking tool
   ```bash
   # Install wrk
   sudo apt-get install wrk
   ```

---

## 📊 Load Testing

### Test 1: Baseline Load Test

**Objective**: Measure performance under normal load

**Configuration**:
- **Concurrent Users**: 100
- **Duration**: 5 minutes
- **Request Rate**: 100 req/s
- **Target**: All 58 gRPC methods

**Command**:
```bash
# Using ghz
ghz --insecure \
  --proto ./proto/blockchain.proto \
  --call blockchain.v1.BlockOperationsService.GetBlock \
  -d '{"block_number":1}' \
  -n 10000 \
  -c 100 \
  localhost:8004
```

**Expected Results**:
- p95 latency: <100ms
- Error rate: <0.1%
- Throughput: ~100 req/s

### Test 2: Sustained Load Test

**Objective**: Measure performance under sustained load

**Configuration**:
- **Concurrent Users**: 500
- **Duration**: 30 minutes
- **Request Rate**: 500 req/s
- **Target**: All 58 gRPC methods

**Command**:
```bash
# Using ghz with multiple methods
ghz --insecure \
  --proto ./proto/blockchain.proto \
  --call blockchain.v1.BlockOperationsService.GetBlock \
  -d '{"block_number":1}' \
  -n 900000 \
  -c 500 \
  -t 30m \
  localhost:8004
```

**Expected Results**:
- p95 latency: <100ms
- Error rate: <0.1%
- No memory leaks
- Stable throughput

---

## 💥 Stress Testing

### Test 3: Peak Load Test

**Objective**: Find maximum capacity

**Configuration**:
- **Concurrent Users**: 1,000
- **Duration**: 10 minutes
- **Request Rate**: Gradually increase to 10,000 req/s
- **Target**: Critical methods (GetBlock, GetTransaction, GetWalletBalance)

**Command**:
```bash
# Using ghz with ramp-up
ghz --insecure \
  --proto ./proto/blockchain.proto \
  --call blockchain.v1.BlockOperationsService.GetBlock \
  -d '{"block_number":1}' \
  -n 6000000 \
  -c 1000 \
  -t 10m \
  --rps 10000 \
  localhost:8004
```

**Expected Results**:
- Identify breaking point
- Measure degradation
- Error rate acceptable until breaking point

### Test 4: Spike Test

**Objective**: Test behavior under sudden load spikes

**Configuration**:
- **Base Load**: 100 req/s
- **Spike**: 5,000 req/s for 1 minute
- **Duration**: 15 minutes
- **Spikes**: 3 spikes at 5min, 10min, 15min

**Command**:
```bash
# Using k6 script
k6 run --vus 100 --duration 15m spike-test.js
```

**Expected Results**:
- Service recovers quickly
- No cascading failures
- Latency returns to normal

---

## 🔬 Benchmark Testing

### Test 5: Method Benchmark

**Objective**: Benchmark individual methods

**Methods to Test**:
1. `GetBlock` - Block retrieval
2. `GetTransaction` - Transaction lookup
3. `GetWalletBalance` - Balance query
4. `SubmitTransaction` - Transaction submission
5. `ProduceBlock` - Block production

**Command**:
```bash
# Benchmark GetBlock
ghz --insecure \
  --proto ./proto/blockchain.proto \
  --call blockchain.v1.BlockOperationsService.GetBlock \
  -d '{"block_number":1}' \
  -n 10000 \
  -c 10 \
  --timeout 5s \
  localhost:8004
```

**Expected Results**:
- Individual method latencies
- Identify slow methods
- Optimization opportunities

### Test 6: Blockchain Operations Benchmark

**Objective**: Benchmark blockchain-specific operations

**Operations**:
1. Block production time
2. Transaction processing time
3. Block validation time
4. Balance query time

**Script**: `tests/benchmark-blockchain.sh`

**Expected Results**:
- Block production: <5s
- Transaction processing: <100ms
- Block validation: <200ms
- Balance query: <50ms

---

## 📈 Performance Test Scripts

### Script 1: Load Test Script

**File**: `tests/load-test.sh`

```bash
#!/bin/bash

# Load test configuration
CONCURRENT_USERS=100
DURATION=5m
TARGET="localhost:8004"

# Run load test
ghz --insecure \
  --proto ./proto/blockchain.proto \
  --call blockchain.v1.BlockOperationsService.GetBlock \
  -d '{"block_number":1}' \
  -n 30000 \
  -c $CONCURRENT_USERS \
  -t $DURATION \
  $TARGET
```

### Script 2: Stress Test Script

**File**: `tests/stress-test.sh`

```bash
#!/bin/bash

# Stress test configuration
MAX_CONCURRENT=1000
RAMP_UP=5m
DURATION=10m
TARGET="localhost:8004"

# Run stress test
ghz --insecure \
  --proto ./proto/blockchain.proto \
  --call blockchain.v1.BlockOperationsService.GetBlock \
  -d '{"block_number":1}' \
  -n 6000000 \
  -c $MAX_CONCURRENT \
  -t $DURATION \
  $TARGET
```

### Script 3: Benchmark Script

**File**: `tests/benchmark.sh`

```bash
#!/bin/bash

# Benchmark all methods
METHODS=(
  "blockchain.v1.BlockOperationsService.GetBlock"
  "blockchain.v1.TransactionOperationsService.GetTransaction"
  "blockchain.v1.USCCoinOperationsService.GetWalletBalance"
)

for method in "${METHODS[@]}"; do
  echo "Benchmarking $method..."
  ghz --insecure \
    --proto ./proto/blockchain.proto \
    --call $method \
    -n 10000 \
    -c 10 \
    localhost:8004
done
```

---

## 📊 Performance Metrics Collection

### During Testing

**Metrics to Collect**:
1. **Request Metrics**
   - Total requests
   - Successful requests
   - Failed requests
   - Request rate (req/s)

2. **Latency Metrics**
   - p50, p95, p99 latencies
   - Min, max, average latencies
   - Latency distribution

3. **Error Metrics**
   - Error rate
   - Error types
   - Error distribution

4. **System Metrics**
   - CPU usage
   - Memory usage
   - Network I/O
   - Disk I/O

5. **Application Metrics**
   - Database connections
   - Cache hit rate
   - Block production time
   - Transaction throughput

### Metrics Collection Tools

**Prometheus**:
```bash
# Query metrics during test
curl http://localhost:9004/metrics | grep grpc_server_request_duration_seconds
```

**Grafana**:
- Real-time dashboard during testing
- Monitor all metrics simultaneously

---

## 📋 Test Scenarios

### Scenario 1: Normal Operations

**Load**: 100 concurrent users, 100 req/s  
**Duration**: 5 minutes  
**Methods**: All 58 methods (round-robin)

**Expected**:
- p95 latency: <100ms
- Error rate: <0.1%
- All methods respond successfully

### Scenario 2: High Read Load

**Load**: 500 concurrent users, 500 req/s  
**Duration**: 10 minutes  
**Methods**: Read-only methods (GetBlock, GetTransaction, GetWalletBalance)

**Expected**:
- p95 latency: <50ms (read operations faster)
- Error rate: <0.1%
- Cache hit rate: >95%

### Scenario 3: High Write Load

**Load**: 200 concurrent users, 200 req/s  
**Duration**: 10 minutes  
**Methods**: Write methods (SubmitTransaction, ProduceBlock)

**Expected**:
- p95 latency: <200ms (write operations slower)
- Error rate: <0.1%
- Block production: <5s

### Scenario 4: Mixed Load

**Load**: 300 concurrent users, 300 req/s  
**Duration**: 15 minutes  
**Methods**: 70% read, 30% write

**Expected**:
- p95 latency: <100ms
- Error rate: <0.1%
- Balanced performance

---

## 🎯 Success Criteria

### Load Test Success

- ✅ p95 latency <100ms
- ✅ Error rate <0.1%
- ✅ Throughput meets target
- ✅ No memory leaks
- ✅ Stable performance

### Stress Test Success

- ✅ Breaking point identified
- ✅ Graceful degradation
- ✅ Service recovers after load
- ✅ No data corruption
- ✅ No cascading failures

### Benchmark Success

- ✅ All methods meet targets
- ✅ Slow methods identified
- ✅ Optimization opportunities documented

---

## 📝 Test Reports

### Report Template

**Test Report Structure**:
1. **Test Configuration**
   - Load parameters
   - Test duration
   - Methods tested

2. **Results Summary**
   - Throughput
   - Latency (p50, p95, p99)
   - Error rate

3. **Performance Analysis**
   - Bottlenecks identified
   - Optimization recommendations
   - Comparison with targets

4. **Recommendations**
   - Scaling recommendations
   - Configuration changes
   - Code optimizations

---

## 🔧 Performance Optimization

### Based on Test Results

1. **If Latency High**:
   - Increase connection pool
   - Optimize database queries
   - Add caching
   - Scale horizontally

2. **If Throughput Low**:
   - Increase concurrent connections
   - Optimize code paths
   - Reduce database queries
   - Use connection pooling

3. **If Error Rate High**:
   - Check resource limits
   - Review error handling
   - Increase timeouts
   - Scale resources

---

## 📊 Continuous Performance Testing

### Integration with CI/CD

**Pre-Production**:
- Run load tests before deployment
- Compare with baseline
- Fail if performance degrades

**Post-Deployment**:
- Monitor performance metrics
- Alert on performance degradation
- Run periodic benchmarks

---

## ✅ Verification Checklist

- [ ] Baseline load test completed
- [ ] Sustained load test completed
- [ ] Stress test completed
- [ ] Spike test completed
- [ ] Method benchmarks completed
- [ ] Blockchain operations benchmarked
- [ ] Performance targets met
- [ ] Test reports generated
- [ ] Optimization recommendations documented

---

**Status**: ✅ **GUIDE COMPLETE**  
**Last Updated**: 2025-11-10


