# 📊 MONITORING & PERFORMANCE MODULES - USAGE GUIDE

## 📋 TỔNG QUAN

**Monitoring** và **Performance** là 2 observability modules tự động thu thập metrics từ blockchain. Chúng có 3 cách sử dụng:

1. **Tự động** - Chạy trong BeginBlock/EndBlock (auto-collection)
2. **Query** - Query metrics qua Cosmos SDK Query interface
3. **Manual Record** - Record metrics thủ công qua Message handlers

---

## 🔄 **CÁCH 1: TỰ ĐỘNG (AUTO-COLLECTION)**

### **Cách hoạt động**

Monitoring và Performance tự động chạy mỗi block trong `BeginBlock()` và `EndBlock()`:

```go
// x/monitoring/abci.go
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
    // 1. Perform health checks
    performHealthChecks(ctx, k)
    
    // 2. Collect performance metrics
    collectPerformanceMetrics(ctx, k)
    
    // 3. Process alerts
    processAlerts(ctx, k)
    
    // 4. Update system health
    updateSystemHealth(ctx, k)
}

func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
    // 1. Finalize metrics collection
    finalizeMetricsCollection(ctx, k)
    
    // 2. Evaluate alert conditions
    evaluateAlertConditions(ctx, k)
    
    // 3. Generate system health report
    generateSystemHealthReport(ctx, k)
    
    // 4. Cleanup old data
    cleanupOldData(ctx, k)
    
    return []abci.ValidatorUpdate{}
}
```

### **Metrics được tự động thu thập**

**Monitoring Module**:
- Health checks cho tất cả services
- System health status
- Alert conditions evaluation
- Performance data collection

**Performance Module**:
- CPU usage
- Memory usage
- Network I/O
- Disk I/O
- Response times
- Error rates
- Throughput

### **Không cần làm gì**

Modules tự động chạy khi blockchain process blocks. Metrics được lưu vào blockchain state.

---

## 🔍 **CÁCH 2: QUERY METRICS (Cosmos SDK Query)**

### **Query Monitoring Metrics**

#### **2.1. Query Monitoring Session**

```go
// Trong business service hoặc handler
import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    monitoringtypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/monitoring/types"
)

func (s *Service) GetMonitoringSession(ctx context.Context, monitoringID string) (*monitoringtypes.MonitoringConfig, error) {
    sdkCtx, err := s.getSDKContext(ctx)
    if err != nil {
        return nil, err
    }
    
    // Query monitoring config từ keeper
    config, err := s.cosmosApp.MonitoringKeeper.GetMonitoringConfig(sdkCtx, monitoringID)
    if err != nil {
        return nil, err
    }
    
    return &config, nil
}
```

#### **2.2. Query All Monitoring Sessions**

```go
func (s *Service) GetAllMonitoringSessions(ctx context.Context) ([]monitoringtypes.MonitoringConfig, error) {
    sdkCtx, err := s.getSDKContext(ctx)
    if err != nil {
        return nil, err
    }
    
    // Get all monitoring configs
    configs := s.cosmosApp.MonitoringKeeper.GetAllMonitoringConfigs(sdkCtx)
    
    return configs, nil
}
```

#### **2.3. Query Metrics**

```go
func (s *Service) GetMetrics(ctx context.Context, monitoringID string) ([]monitoringtypes.Metric, error) {
    sdkCtx, err := s.getSDKContext(ctx)
    if err != nil {
        return nil, err
    }
    
    // Get all metrics
    allMetrics := s.cosmosApp.MonitoringKeeper.GetAllMetrics(sdkCtx)
    
    // Filter by monitoring ID if needed
    var filtered []monitoringtypes.Metric
    for _, metric := range allMetrics {
        // Filter logic here
        filtered = append(filtered, metric)
    }
    
    return filtered, nil
}
```

#### **2.4. Query Performance Metrics**

```go
import (
    performancetypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/performance/types"
)

func (s *Service) GetPerformanceMetrics(ctx context.Context, serviceName string) ([]performancetypes.PerformanceMetric, error) {
    sdkCtx, err := s.getSDKContext(ctx)
    if err != nil {
        return nil, err
    }
    
    // Get all performance metrics
    allMetrics := s.cosmosApp.PerformanceKeeper.GetAllPerformanceMetrics(sdkCtx)
    
    // Filter by service name
    var filtered []performancetypes.PerformanceMetric
    for _, metric := range allMetrics {
        if metric.Tags["service"] == serviceName {
            filtered = append(filtered, metric)
        }
    }
    
    return filtered, nil
}
```

#### **2.5. Query System Health**

```go
func (s *Service) GetSystemHealth(ctx context.Context, serviceName string) (*monitoringtypes.SystemHealth, error) {
    sdkCtx, err := s.getSDKContext(ctx)
    if err != nil {
        return nil, err
    }
    
    // Get system health
    health, err := s.cosmosApp.MonitoringKeeper.GetSystemHealth(sdkCtx, serviceName)
    if err != nil {
        return nil, err
    }
    
    return &health, nil
}
```

### **Query qua gRPC (Cosmos SDK Query Server)**

Monitoring và Performance đã có `query_server.go` với các methods:

**Monitoring Query Methods**:
- `QueryMonitoring` - Query single monitoring session
- `QueryMonitoringSessions` - Query all monitoring sessions
- `QueryMetrics` - Query metrics for a monitoring session
- `QueryMonitoringStats` - Query monitoring statistics

**Performance Query Methods**:
- `QueryMetrics` - Query performance metrics
- `QueryMetricsList` - Query multiple performance metrics
- `QueryPerformanceStats` - Query performance statistics
- `QueryOptimization` - Query optimization results

**Cách sử dụng qua gRPC**:

```go
// Trong handler hoặc client
import (
    "google.golang.org/grpc"
    monitoringproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/monitoring/v1/usc/monitoring/v1"
)

// Connect to Cosmos SDK gRPC server
conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
if err != nil {
    return err
}
defer conn.Close()

client := monitoringproto.NewQueryClient(conn)

// Query monitoring session
req := &monitoringproto.QueryMonitoringRequest{
    MonitoringId: "service-04-usc-blockchain-core",
}

resp, err := client.QueryMonitoring(ctx, req)
if err != nil {
    return err
}

// Use resp.Session
```

---

## ✍️ **CÁCH 3: RECORD METRICS THỦ CÔNG (Manual Recording)**

### **Record Monitoring Metrics**

#### **3.1. Create Monitoring Config**

```go
import (
    monitoringtypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/monitoring/types"
    "time"
)

func (s *Service) CreateMonitoringConfig(ctx context.Context, serviceName string) error {
    sdkCtx, err := s.getSDKContext(ctx)
    if err != nil {
        return err
    }
    
    // Create monitoring config
    config := monitoringtypes.MonitoringConfig{
        ID:            fmt.Sprintf("monitoring_%s_%d", serviceName, time.Now().Unix()),
        ServiceName:   serviceName,
        Enabled:       true,
        CheckInterval: 5 * time.Second,
        RetentionPeriod: 30 * 24 * time.Hour, // 30 days
        CreatedAt:    time.Now(),
        UpdatedAt:     time.Now(),
    }
    
    // Set config in keeper
    if err := s.cosmosApp.MonitoringKeeper.SetMonitoringConfig(sdkCtx, config); err != nil {
        return err
    }
    
    return nil
}
```

#### **3.2. Record Metric**

```go
func (s *Service) RecordMetric(ctx context.Context, monitoringID, metricName string, value int64) error {
    sdkCtx, err := s.getSDKContext(ctx)
    if err != nil {
        return err
    }
    
    // Create metric
    metric := monitoringtypes.Metric{
        ID:          fmt.Sprintf("metric_%s_%d", metricName, time.Now().Unix()),
        Name:        metricName,
        Value:       value,
        Unit:        "count",
        Timestamp:   time.Now(),
        Tags:        map[string]string{"monitoring_id": monitoringID},
        Description: fmt.Sprintf("Metric %s for %s", metricName, monitoringID),
    }
    
    // Set metric in keeper
    if err := s.cosmosApp.MonitoringKeeper.SetMetric(sdkCtx, metric); err != nil {
        return err
    }
    
    return nil
}
```

#### **3.3. Record Performance Metric**

```go
import (
    performancetypes "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/performance/types"
)

func (s *Service) RecordPerformanceMetric(ctx context.Context, serviceName, metricName string, value int64, category string) error {
    sdkCtx, err := s.getSDKContext(ctx)
    if err != nil {
        return err
    }
    
    // Create performance metric
    metric := performancetypes.PerformanceMetric{
        ID:          fmt.Sprintf("perf_%s_%s_%d", serviceName, metricName, time.Now().Unix()),
        Name:        metricName,
        Value:       value,
        Unit:        "percent",
        Timestamp:   time.Now(),
        Category:    category, // "cpu", "memory", "network", "disk"
        Tags:        map[string]string{"service": serviceName},
        Description: fmt.Sprintf("Performance metric %s for %s", metricName, serviceName),
    }
    
    // Set metric in keeper
    if err := s.cosmosApp.PerformanceKeeper.SetPerformanceMetric(sdkCtx, metric); err != nil {
        return err
    }
    
    return nil
}
```

#### **3.4. Record Performance Data (Monitoring Module)**

```go
func (s *Service) RecordPerformanceData(ctx context.Context, serviceName, metricName string, value int64, unit string) error {
    sdkCtx, err := s.getSDKContext(ctx)
    if err != nil {
        return err
    }
    
    // Create performance data
    perfData := monitoringtypes.PerformanceData{
        ID:          fmt.Sprintf("perf_data_%s_%d", serviceName, time.Now().Unix()),
        ServiceName: serviceName,
        MetricName:  metricName,
        Value:       value,
        Unit:        unit,
        Timestamp:   time.Now(),
        Metadata:    make(map[string]string),
    }
    
    // Set performance data in keeper
    if err := s.cosmosApp.MonitoringKeeper.SetPerformanceData(sdkCtx, perfData); err != nil {
        return err
    }
    
    return nil
}
```

---

## 🎯 **VÍ DỤ TÍCH HỢP VÀO BUSINESS SERVICES**

### **Ví dụ 1: Record Block Performance Metrics**

```go
// Trong block_operations_service.go
func (s *Service) ProduceBlock(ctx context.Context, req *proto.ProduceBlockRequest) (*proto.ProduceBlockResponse, error) {
    start := time.Now()
    
    // ... block production logic ...
    
    // Record performance metric
    duration := time.Since(start).Milliseconds()
    if s.cosmosApp != nil {
        sdkCtx, _ := s.getSDKContext(ctx)
        metric := performancetypes.PerformanceMetric{
            ID:          fmt.Sprintf("block_produce_%d", time.Now().Unix()),
            Name:        "block_production_time",
            Value:       duration,
            Unit:        "milliseconds",
            Timestamp:   time.Now(),
            Category:    "performance",
            Tags:        map[string]string{"service": "block_operations", "validator": req.ValidatorId},
            Description: "Block production time",
        }
        _ = s.cosmosApp.PerformanceKeeper.SetPerformanceMetric(sdkCtx, metric)
    }
    
    return response, nil
}
```

### **Ví dụ 2: Record Network Health Metrics**

```go
// Trong network_operations_service.go
func (s *Service) GetNetworkInfo(ctx context.Context) (*proto.GetNetworkInfoResponse, error) {
    // ... get network info logic ...
    
    // Record network health metric
    if s.cosmosApp != nil {
        sdkCtx, _ := s.getSDKContext(ctx)
        health := monitoringtypes.SystemHealth{
            ID:        fmt.Sprintf("network_health_%d", time.Now().Unix()),
            Status:    "healthy",
            Score:     95.0,
            Timestamp: time.Now(),
            Components: []monitoringtypes.ComponentHealth{
                {
                    Name:      "network",
                    Status:    "healthy",
                    Score:     95.0,
                    LastCheck: time.Now(),
                    Message:   "Network is running normally",
                },
            },
            Summary: "Network is healthy",
        }
        _ = s.cosmosApp.MonitoringKeeper.SetSystemHealth(sdkCtx, health)
    }
    
    return response, nil
}
```

### **Ví dụ 3: Query Metrics trong Business Service**

```go
// Trong network_operations_service.go
func (s *Service) GetNetworkStats(ctx context.Context, req *proto.GetNetworkStatsRequest) (*proto.GetNetworkStatsResponse, error) {
    // Query monitoring metrics
    var metrics []monitoringtypes.Metric
    if s.cosmosApp != nil {
        sdkCtx, _ := s.getSDKContext(ctx)
        allMetrics := s.cosmosApp.MonitoringKeeper.GetAllMetrics(sdkCtx)
        
        // Filter by time range and metric type
        for _, metric := range allMetrics {
            if metric.Tags["service"] == "network" {
                metrics = append(metrics, metric)
            }
        }
    }
    
    // Convert to proto response
    // ...
    
    return response, nil
}
```

---

## 📊 **KEEPER METHODS AVAILABLE**

### **MonitoringKeeper Methods**

```go
// Config Management
GetMonitoringConfig(ctx sdk.Context, id string) (MonitoringConfig, error)
SetMonitoringConfig(ctx sdk.Context, config MonitoringConfig) error
GetAllMonitoringConfigs(ctx sdk.Context) []MonitoringConfig

// Metrics
GetMetric(ctx sdk.Context, id string) (Metric, error)
SetMetric(ctx sdk.Context, metric Metric) error
GetAllMetrics(ctx sdk.Context) []Metric

// Alerts
GetAlert(ctx sdk.Context, id string) (Alert, error)
SetAlert(ctx sdk.Context, alert Alert) error
GetAllAlerts(ctx sdk.Context) []Alert

// Performance Data
GetPerformanceData(ctx sdk.Context, id string) (PerformanceData, error)
SetPerformanceData(ctx sdk.Context, data PerformanceData) error
GetAllPerformanceData(ctx sdk.Context) []PerformanceData

// System Health
GetSystemHealth(ctx sdk.Context, serviceName string) (SystemHealth, error)
SetSystemHealth(ctx sdk.Context, health SystemHealth) error
GetAllSystemHealth(ctx sdk.Context) []SystemHealth
```

### **PerformanceKeeper Methods**

```go
// Performance Metrics
GetPerformanceMetric(ctx sdk.Context, id string) (PerformanceMetric, error)
SetPerformanceMetric(ctx sdk.Context, metric PerformanceMetric) error
GetAllPerformanceMetrics(ctx sdk.Context) []PerformanceMetric

// Benchmarks
GetBenchmark(ctx sdk.Context, id string) (Benchmark, error)
SetBenchmark(ctx sdk.Context, benchmark Benchmark) error
GetAllBenchmarks(ctx sdk.Context) []Benchmark

// Optimizations
GetOptimization(ctx sdk.Context, id string) (Optimization, error)
SetOptimization(ctx sdk.Context, optimization Optimization) error
GetAllOptimizations(ctx sdk.Context) []Optimization

// Performance Profiles
GetPerformanceProfile(ctx sdk.Context, id string) (PerformanceProfile, error)
SetPerformanceProfile(ctx sdk.Context, profile PerformanceProfile) error
GetAllPerformanceProfiles(ctx sdk.Context) []PerformanceProfile

// Performance Reports
GetPerformanceReport(ctx sdk.Context, id string) (PerformanceReport, error)
SetPerformanceReport(ctx sdk.Context, report PerformanceReport) error
GetAllPerformanceReports(ctx sdk.Context) []PerformanceReport
```

---

## 🔗 **TÍCH HỢP VỚI CÁC MODULES KHÁC**

### **Block Operations + Performance**

```go
// Trong block_operations_service.go
func (s *Service) ProduceBlock(ctx context.Context, req *proto.ProduceBlockRequest) (*proto.ProduceBlockResponse, error) {
    start := time.Now()
    
    // Produce block
    response, err := s.repo.ProduceBlock(ctx, req)
    
    // Record performance metric
    if s.cosmosApp != nil && err == nil {
        sdkCtx, _ := s.getSDKContext(ctx)
        duration := time.Since(start).Milliseconds()
        
        metric := performancetypes.PerformanceMetric{
            ID:          fmt.Sprintf("block_produce_%d", time.Now().Unix()),
            Name:        "block_production_time",
            Value:       duration,
            Unit:        "milliseconds",
            Timestamp:   time.Now(),
            Category:    "performance",
            Tags:        map[string]string{"service": "block_operations"},
        }
        _ = s.cosmosApp.PerformanceKeeper.SetPerformanceMetric(sdkCtx, metric)
    }
    
    return response, err
}
```

### **Network Operations + Monitoring**

```go
// Trong network_operations_service.go
func (s *Service) GetNetworkInfo(ctx context.Context) (*proto.GetNetworkInfoResponse, error) {
    // Get network info
    response, err := s.repo.GetNetworkInfo(ctx)
    
    // Record monitoring metric
    if s.cosmosApp != nil && err == nil {
        sdkCtx, _ := s.getSDKContext(ctx)
        
        metric := monitoringtypes.Metric{
            ID:          fmt.Sprintf("network_info_%d", time.Now().Unix()),
            Name:        "network_info_query",
            Value:       1,
            Unit:        "count",
            Timestamp:   time.Now(),
            Tags:        map[string]string{"service": "network_operations"},
        }
        _ = s.cosmosApp.MonitoringKeeper.SetMetric(sdkCtx, metric)
    }
    
    return response, err
}
```

### **Streaming Operations + Performance**

```go
// Trong streaming_operations_service.go
func (s *Service) StreamBlocks(ctx context.Context, req *proto.StreamBlocksRequest) (*proto.StreamBlocksResponse, error) {
    start := time.Now()
    
    // Stream blocks
    response, err := s.repo.StreamBlocks(ctx, req)
    
    // Record performance metric
    if s.cosmosApp != nil && err == nil {
        sdkCtx, _ := s.getSDKContext(ctx)
        duration := time.Since(start).Milliseconds()
        
        metric := performancetypes.PerformanceMetric{
            ID:          fmt.Sprintf("stream_blocks_%d", time.Now().Unix()),
            Name:        "stream_blocks_latency",
            Value:       duration,
            Unit:        "milliseconds",
            Timestamp:   time.Now(),
            Category:    "network",
            Tags:        map[string]string{"service": "streaming_operations", "client": req.ClientId},
        }
        _ = s.cosmosApp.PerformanceKeeper.SetPerformanceMetric(sdkCtx, metric)
    }
    
    return response, err
}
```

---

## 🎯 **BEST PRACTICES**

### **1. Sử dụng Auto-Collection cho System Metrics**
- Để monitoring và performance tự động thu thập system metrics
- Không cần record thủ công cho basic metrics

### **2. Record Manual Metrics cho Business Logic**
- Record custom metrics cho business operations
- Ví dụ: block production time, transaction processing time

### **3. Query Metrics khi cần Analytics**
- Query metrics từ keeper khi cần analytics
- Sử dụng filters để query specific metrics

### **4. Sử dụng Tags để Filter**
- Thêm tags vào metrics để dễ filter
- Ví dụ: `tags["service"] = "block_operations"`

### **5. Cleanup Old Data**
- Monitoring và Performance tự động cleanup old data trong EndBlock
- Có thể configure retention period trong MonitoringConfig

---

## 📝 **TÓM TẮT**

### **Monitoring Module**
- **Auto-collection**: ✅ Tự động trong BeginBlock/EndBlock
- **Query**: ✅ Query qua MonitoringKeeper methods
- **Manual Record**: ✅ Record qua MonitoringKeeper.SetMetric()

### **Performance Module**
- **Auto-collection**: ✅ Tự động trong BeginBlock/EndBlock
- **Query**: ✅ Query qua PerformanceKeeper methods
- **Manual Record**: ✅ Record qua PerformanceKeeper.SetPerformanceMetric()

### **Cách sử dụng phổ biến**
1. **Auto-collection**: Để modules tự động thu thập (không cần làm gì)
2. **Query trong Business Services**: Query metrics khi cần analytics
3. **Record Custom Metrics**: Record metrics cho business operations cụ thể

---

**🎉 Monitoring & Performance modules đã sẵn sàng sử dụng!**


