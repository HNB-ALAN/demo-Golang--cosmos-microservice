# Business Layer Comparison: Service-04 vs Service-22

**Date**: 2025-11-12  
**Purpose**: Đối chiếu business layer của service-04 với service-22 để xác định phần còn thiếu

---

## 📊 **SO SÁNH TỔNG QUAN**

### **Service-22 (Reference Pattern)** ✅

**Components trong Business Service:**
```go
type Service struct {
    cfg       *config.Config
    log       *logging.Logger
    repo      *repo.Repository
    validator *validation.Validator      // ✅ Có validator service
    metrics   *metrics.MetricsService    // ✅ Có metrics service
}
```

**Pattern trong mỗi method:**
1. ✅ `start := time.Now()` + `defer func() { s.metrics.RecordDuration(...) }()`
2. ✅ Input validation sử dụng `s.validator.ValidateXXX()`
3. ✅ Logging với correlation ID
4. ✅ Metrics recording: `RecordDuration()`, `RecordFailure()`, `RecordSuccess()`
5. ✅ Error handling với gRPC status codes: `status.Errorf(codes.InvalidArgument, ...)`
6. ✅ Delegate to repository (không có direct infrastructure access)

---

### **Service-04 (Current State)** ⚠️

**Components trong Business Service:**
```go
type Service struct {
    repo              *repo.Repository
    cosmosApp         *app.USCApp
    blockchainStorage *storage.StateManager
    logger            *logging.Logger
    // ❌ Không có validator service
    // ❌ Không có metrics service
}
```

**Pattern trong mỗi method:**
1. ✅ `start := time.Now()` (có)
2. ⚠️ `defer func() { s.metrics.RecordDuration(...) }()` (KHÔNG có - chỉ có manual recording)
3. ⚠️ Input validation: Manual `if req.Field == ""` (KHÔNG có validator service)
4. ✅ Logging với correlation ID (có)
5. ⚠️ Metrics recording: `utils.RecordPerformanceMetric()` (helper function, không consistent)
6. ⚠️ Error handling: Return response với `ErrorMessage` (KHÔNG dùng gRPC status codes)
7. ✅ Delegate to repository (đã refactor xong)

---

## 🔍 **CHI TIẾT SO SÁNH**

### **1. Validator Service** ❌

**Service-22:**
```go
// Infrastructure: validation/validator.go
type Validator struct {
    // Validation logic
}

func (v *Validator) ValidateTopicName(name string) error {
    // Structured validation
}

// Business layer usage:
if err := s.validator.ValidateTopicName(req.Topic); err != nil {
    s.metrics.RecordFailure("create_topic", "validation_error", ...)
    return nil, status.Errorf(codes.InvalidArgument, "invalid topic name: %v", err)
}
```

**Service-04:**
```go
// Manual validation:
if req.WalletAddress == "" {
    s.logger.Warn("Empty wallet address", ...)
    return &proto.GetWalletBalanceResponse{
        Success: false,
        ErrorMessage: "wallet address is required",
    }, nil
}
```

**Thiếu:**
- ❌ Không có validator service/infrastructure
- ❌ Manual validation thay vì structured validation
- ❌ Không có validation error metrics

---

### **2. Metrics Service** ⚠️

**Service-22:**
```go
// Infrastructure: metrics/metrics_service.go
type MetricsService struct {
    // Metrics recording
}

func (m *MetricsService) RecordDuration(operation string, duration time.Duration) {
    // Record duration metric
}

func (m *MetricsService) RecordFailure(operation, errorType string, tags map[string]string) {
    // Record failure metric
}

func (m *MetricsService) RecordSuccess(operation string, tags map[string]string) {
    // Record success metric
}

// Business layer usage:
start := time.Now()
defer func() {
    s.metrics.RecordDuration("create_topic", time.Since(start))
}()

if err != nil {
    s.metrics.RecordFailure("create_topic", "repository_error", ...)
}
```

**Service-04:**
```go
// Helper function: utils.RecordPerformanceMetric()
start := time.Now()
// ... operation ...
if utils.IsCosmosAppAvailable(s.cosmosApp) && err == nil {
    _ = utils.RecordPerformanceMetric(ctx, s.cosmosApp, s.logger, start, ...)
}
```

**Thiếu:**
- ⚠️ Không có dedicated metrics service
- ⚠️ Không có defer pattern (metrics chỉ record khi success)
- ⚠️ Không có `RecordFailure()` và `RecordSuccess()` methods
- ⚠️ Metrics recording không consistent (chỉ một số methods có)

---

### **3. Error Handling** ⚠️

**Service-22:**
```go
// Use gRPC status codes
import "google.golang.org/grpc/codes"
import "google.golang.org/grpc/status"

if err := s.validator.ValidateTopicName(req.Topic); err != nil {
    return nil, status.Errorf(codes.InvalidArgument, "invalid topic name: %v", err)
}

if err != nil {
    return nil, status.Errorf(codes.Internal, "failed to create topic: %v", err)
}
```

**Service-04:**
```go
// Return response với ErrorMessage
if req.WalletAddress == "" {
    return &proto.GetWalletBalanceResponse{
        Success: false,
        ErrorMessage: "wallet address is required",
    }, nil
}

if err != nil {
    return &proto.GetWalletBalanceResponse{
        Success: false,
        ErrorMessage: err.Error(),
    }, nil
}
```

**Thiếu:**
- ⚠️ Không sử dụng gRPC status codes
- ⚠️ Error handling không consistent với gRPC best practices

---

### **4. Metrics Recording Pattern** ⚠️

**Service-22:**
```go
func (s *Service) CreateTopic(ctx context.Context, req *proto.CreateTopicRequest) (*proto.CreateTopicResponse, error) {
    start := time.Now()
    defer func() {
        s.metrics.RecordDuration("create_topic", time.Since(start))
    }()

    // ... validation ...

    // ... repository call ...

    if err != nil {
        s.metrics.RecordFailure("create_topic", "repository_error", ...)
        return nil, err
    }

    s.metrics.RecordSuccess("create_topic", ...)
    return resp, nil
}
```

**Service-04:**
```go
func (s *Service) GetUSCBalance(ctx context.Context, req *proto.GetWalletBalanceRequest) (*proto.GetWalletBalanceResponse, error) {
    start := time.Now()
    // ... validation ...
    // ... repository call ...

    // Metrics chỉ record khi success
    if utils.IsCosmosAppAvailable(s.cosmosApp) && err == nil && response != nil {
        _ = utils.RecordPerformanceMetric(ctx, s.cosmosApp, s.logger, start, ...)
    }

    return response, err
}
```

**Thiếu:**
- ⚠️ Không có defer pattern (metrics không được record khi có error)
- ⚠️ Không có failure metrics
- ⚠️ Metrics recording không consistent

---

## 📋 **CHECKLIST: PHẦN CÒN THIẾU**

### **1. Validator Infrastructure** ❌

- [ ] Tạo `internal/infrastructure/validation/validator.go`
- [ ] Implement validation methods cho từng domain:
  - [ ] `ValidateWalletAddress(address string) error`
  - [ ] `ValidateAmount(amount string) error`
  - [ ] `ValidateBlockNumber(blockNumber int64) error`
  - [ ] `ValidateTransactionHash(hash string) error`
  - [ ] `ValidateContractAddress(address string) error`
  - [ ] `ValidateTokenId(tokenId string) error`
- [ ] Inject validator vào business services
- [ ] Replace manual validation với validator service

---

### **2. Metrics Infrastructure** ⚠️

- [ ] Tạo `internal/infrastructure/metrics/metrics_service.go`
- [ ] Implement metrics methods:
  - [ ] `RecordDuration(operation string, duration time.Duration)`
  - [ ] `RecordFailure(operation, errorType string, tags map[string]string)`
  - [ ] `RecordSuccess(operation string, tags map[string]string)`
  - [ ] `RecordBlockchainOperation(operation, result string, tags map[string]string)`
- [ ] Inject metrics service vào business services
- [ ] Replace `utils.RecordPerformanceMetric()` với metrics service
- [ ] Add defer pattern cho tất cả methods

---

### **3. Error Handling** ⚠️

- [ ] Import gRPC status codes: `google.golang.org/grpc/codes`, `google.golang.org/grpc/status`
- [ ] Replace response error messages với gRPC status codes:
  - [ ] `codes.InvalidArgument` cho validation errors
  - [ ] `codes.NotFound` cho not found errors
  - [ ] `codes.Internal` cho repository errors
  - [ ] `codes.FailedPrecondition` cho business rule violations
- [ ] Update tất cả business methods để return gRPC status errors

---

### **4. Metrics Recording Pattern** ⚠️

- [ ] Add defer pattern cho tất cả methods:
  ```go
  start := time.Now()
  defer func() {
      s.metrics.RecordDuration("operation_name", time.Since(start))
  }()
  ```
- [ ] Add failure metrics recording:
  ```go
  if err != nil {
      s.metrics.RecordFailure("operation_name", "error_type", tags)
      return nil, status.Errorf(codes.Internal, "...")
  }
  ```
- [ ] Add success metrics recording:
  ```go
  s.metrics.RecordSuccess("operation_name", tags)
  return resp, nil
  ```

---

### **5. Consistency Check** ⚠️

- [ ] Tất cả methods có correlation ID logging
- [ ] Tất cả methods có metrics recording (duration + success/failure)
- [ ] Tất cả methods có structured validation
- [ ] Tất cả methods có consistent error handling
- [ ] Tất cả methods delegate to repository (không có direct infrastructure access)

---

## 🎯 **PRIORITY**

### **High Priority** (Core Infrastructure)
1. ✅ **Validator Infrastructure** - Cần cho structured validation
2. ✅ **Metrics Infrastructure** - Cần cho consistent metrics recording
3. ✅ **Error Handling** - Cần cho gRPC best practices

### **Medium Priority** (Pattern Consistency)
4. ⚠️ **Metrics Recording Pattern** - Cần defer pattern cho tất cả methods
5. ⚠️ **Validation Pattern** - Replace manual validation với validator service

### **Low Priority** (Nice to Have)
6. ⚠️ **Additional Metrics** - Custom metrics cho blockchain operations
7. ⚠️ **Validation Rules** - More comprehensive validation rules

---

## 📝 **IMPLEMENTATION NOTES**

### **Validator Service Structure**
```go
// internal/infrastructure/validation/validator.go
package validation

type Validator struct {
    // Validation rules
}

func NewValidator() *Validator {
    return &Validator{}
}

func (v *Validator) ValidateWalletAddress(address string) error {
    // Validation logic
}

func (v *Validator) ValidateAmount(amount string) error {
    // Validation logic
}
```

### **Metrics Service Structure**
```go
// internal/infrastructure/metrics/metrics_service.go
package metrics

type MetricsService struct {
    // Metrics recording
}

func NewMetricsService() *MetricsService {
    return &MetricsService{}
}

func (m *MetricsService) RecordDuration(operation string, duration time.Duration) {
    // Record duration
}

func (m *MetricsService) RecordFailure(operation, errorType string, tags map[string]string) {
    // Record failure
}

func (m *MetricsService) RecordSuccess(operation string, tags map[string]string) {
    // Record success
}
```

---

## ✅ **KẾT LUẬN**

**Service-04 cần bổ sung:**
1. ❌ **Validator Infrastructure** - Structured validation service
2. ⚠️ **Metrics Infrastructure** - Dedicated metrics service với defer pattern
3. ⚠️ **Error Handling** - gRPC status codes thay vì response error messages
4. ⚠️ **Metrics Pattern** - Consistent defer pattern cho tất cả methods

**Service-04 đã có:**
- ✅ Correlation ID logging
- ✅ Performance metrics (nhưng không consistent)
- ✅ Repository delegation (đã refactor xong)
- ✅ Input validation (nhưng manual)

---

**Last Updated**: 2025-11-12  
**Status**: ⚠️ **CẦN BỔ SUNG** - Validator + Metrics Infrastructure + Error Handling Pattern

