# Business Layer Missing Components - Service-04 vs Service-22

**Date**: 2025-11-12  
**Status**: ⚠️ **CẦN BỔ SUNG**

---

## 📊 **TÓM TẮT**

### **✅ Service-04 ĐÃ CÓ:**
1. ✅ **Validator Infrastructure** - `internal/infrastructure/validation/validator_backend.go`
2. ✅ **Metrics Infrastructure** - `internal/infrastructure/metrics/metrics.go`
3. ✅ **Container có validator và metrics** - `container-grpc.go` (line 76-77)

### **❌ Service-04 CÒN THIẾU:**
1. ❌ **Business services chưa được inject validator và metrics**
2. ❌ **Business methods chưa sử dụng validator service**
3. ❌ **Business methods chưa sử dụng metrics service với defer pattern**
4. ❌ **Error handling chưa dùng gRPC status codes**

---

## 🔍 **CHI TIẾT VẤN ĐỀ**

### **1. Container có nhưng Business Services chưa nhận** ❌

**Container (container-grpc.go):**
```go
type Container struct {
    validator *validation.Validator  // ✅ Có
    metrics   *metrics.MetricsService // ✅ Có
    // ...
}

// initializeBusiness() - Line 260-300
c.BlockService = blockbiz.NewService(
    c.BlockRepository, 
    c.cosmosApp, 
    c.blockchainStorage, 
    c.logger
    // ❌ THIẾU: c.validator, c.metrics
)
```

**Service-22 (Reference):**
```go
func NewService(
    cfg *config.Config,
    log *logging.Logger,
    repository *repo.Repository,
    validator *validation.Validator,      // ✅ Có
    metricsService *metrics.MetricsService, // ✅ Có
) *Service {
    return &Service{
        cfg: cfg,
        log: log,
        repo: repository,
        validator: validator,      // ✅ Inject
        metrics: metricsService,   // ✅ Inject
    }
}
```

**Service-04 (Current):**
```go
func NewService(
    repo *repo.Repository,
    cosmosApp *app.USCApp,
    blockchainStorage *storage.StateManager,
    logger *logging.Logger
    // ❌ THIẾU: validator, metrics
) *Service {
    return &Service{
        repo: repo,
        cosmosApp: cosmosApp,
        blockchainStorage: blockchainStorage,
        logger: logger,
        // ❌ THIẾU: validator, metrics fields
    }
}
```

---

### **2. Validator Methods chưa đầy đủ** ⚠️

**Service-22 có:**
- `ValidateTopicName()`
- `ValidateConsumerGroup()`
- `ValidateMessageStructure()`
- `ValidatePartitionCount()`
- `ValidateHeaders()`

**Service-04 có:**
- `ValidateTransaction()` - Generic
- `ValidateModelInput()` - AI-related
- `ValidateQuery()` - Search-related

**Service-04 THIẾU:**
- ❌ `ValidateWalletAddress(address string) error`
- ❌ `ValidateAmount(amount string) error`
- ❌ `ValidateBlockNumber(blockNumber int64) error`
- ❌ `ValidateTransactionHash(hash string) error`
- ❌ `ValidateContractAddress(address string) error`
- ❌ `ValidateTokenId(tokenId string) error`
- ❌ `ValidateValidatorAddress(address string) error`

---

### **3. Metrics Service Methods chưa đầy đủ** ⚠️

**Service-22 có:**
```go
func (m *MetricsService) RecordDuration(operation string, duration time.Duration)
func (m *MetricsService) RecordSuccess(operation string, labels map[string]string)
func (m *MetricsService) RecordFailure(operation, errorType string, labels map[string]string)
func (m *MetricsService) RecordTopicCreated(...)
func (m *MetricsService) RecordMessagePublished(...)
```

**Service-04 có:**
```go
func (m *MetricsService) RecordRequest(latency time.Duration, isError bool)
func (m *MetricsService) RecordCacheOperation(isHit bool)
func (m *MetricsService) RecordDatabaseQuery(latency time.Duration)
func (m *MetricsService) RecordBlockchainMetrics(...)
```

**Service-04 THIẾU:**
- ❌ `RecordDuration(operation string, duration time.Duration)` - Defer pattern
- ❌ `RecordSuccess(operation string, labels map[string]string)`
- ❌ `RecordFailure(operation, errorType string, labels map[string]string)`
- ❌ Blockchain-specific metrics methods:
  - ❌ `RecordBlockCreated(blockNumber int64, blockHash string)`
  - ❌ `RecordTransactionSubmitted(txHash string, from, to string)`
  - ❌ `RecordUSCTransfer(amount string, from, to string)`
  - ❌ `RecordContractDeployed(contractAddress string)`
  - ❌ `RecordNFTMinted(tokenId, contractAddress string)`

---

### **4. Business Methods Pattern chưa đúng** ⚠️

**Service-22 Pattern (ĐÚNG):**
```go
func (s *Service) CreateTopic(ctx context.Context, req *proto.CreateTopicRequest) (*proto.CreateTopicResponse, error) {
    start := time.Now()
    defer func() {
        s.metrics.RecordDuration("create_topic", time.Since(start))
    }()

    // 1. Validation
    if err := s.validator.ValidateTopicName(req.TopicName); err != nil {
        s.metrics.RecordFailure("create_topic", "validation_error", ...)
        return nil, status.Errorf(codes.InvalidArgument, "invalid topic name: %v", err)
    }

    // 2. Log
    s.log.Info("Creating Kafka topic", ...)

    // 3. Repository
    resp, err := s.repo.CreateTopic(ctx, req)
    if err != nil {
        s.metrics.RecordFailure("create_topic", "repository_error", ...)
        return nil, status.Errorf(codes.Internal, "failed to create topic: %v", err)
    }

    // 4. Success metrics
    s.metrics.RecordSuccess("create_topic", ...)
    return resp, nil
}
```

**Service-04 Pattern (CHƯA ĐÚNG):**
```go
func (s *Service) GetUSCBalance(ctx context.Context, req *proto.GetWalletBalanceRequest) (*proto.GetWalletBalanceResponse, error) {
    start := time.Now()
    correlationID := utils.GetCorrelationID(ctx)
    s.logger.Info("Getting USC balance", ...)

    // 1. Manual validation (KHÔNG dùng validator service)
    if req.WalletAddress == "" {
        s.logger.Warn("Empty wallet address", ...)
        return &proto.GetWalletBalanceResponse{
            Success: false,
            ErrorMessage: "wallet address is required", // ❌ Không dùng gRPC status
        }, nil
    }

    // 2. Repository
    response, err := s.repo.GetUSCBalance(ctx, req)

    // 3. Metrics chỉ record khi success (KHÔNG có defer)
    if utils.IsCosmosAppAvailable(s.cosmosApp) && err == nil && response != nil {
        _ = utils.RecordPerformanceMetric(...) // ❌ Helper function, không consistent
    }

    return response, err
}
```

---

## 📋 **CHECKLIST: CẦN BỔ SUNG**

### **1. Update Business Service Constructors** ❌

**Cần update 12 business services:**
- [ ] `block_operations/block_operations_service.go`
- [ ] `transaction_operations/transaction_operations_service.go`
- [ ] `usc_coin_operations/usc_coin_operations_service.go`
- [ ] `smart_contract_operations/smart_contract_operations_service.go`
- [ ] `nft_token_operations/nft_token_operations_service.go`
- [ ] `custom_token_operations/custom_token_operations_service.go`
- [ ] `product_certificate_operations/product_certificate_operations_service.go`
- [ ] `validator_operations/validator_operations_service.go`
- [ ] `network_operations/network_operations_service.go`
- [ ] `streaming_operations/streaming_operations_service.go`
- [ ] `store_bridge_operations/store_bridge_operations_service.go`
- [ ] `store_network_operations/store_network_operations_service.go`

**Thay đổi:**
```go
// Before
func NewService(repo *repo.Repository, cosmosApp *app.USCApp, blockchainStorage *storage.StateManager, logger *logging.Logger) *Service

// After
func NewService(
    repo *repo.Repository,
    cosmosApp *app.USCApp,
    blockchainStorage *storage.StateManager,
    logger *logging.Logger,
    validator *validation.Validator,      // ✅ Add
    metrics *metrics.MetricsService,      // ✅ Add
) *Service
```

---

### **2. Update Container Injection** ❌

**File:** `internal/application/container-grpc.go`

**Update `initializeBusiness()` method (line 260-300):**
```go
// Before
c.BlockService = blockbiz.NewService(c.BlockRepository, c.cosmosApp, c.blockchainStorage, c.logger)

// After
c.BlockService = blockbiz.NewService(
    c.BlockRepository,
    c.cosmosApp,
    c.blockchainStorage,
    c.logger,
    c.validator,  // ✅ Add
    c.metrics,    // ✅ Add
)
```

**Cần update cho 12 services.**

---

### **3. Add Validator Methods** ⚠️

**File:** `internal/infrastructure/validation/validator_backend.go`

**Cần thêm:**
- [ ] `ValidateWalletAddress(address string) error`
- [ ] `ValidateAmount(amount string) error`
- [ ] `ValidateBlockNumber(blockNumber int64) error`
- [ ] `ValidateTransactionHash(hash string) error`
- [ ] `ValidateContractAddress(address string) error`
- [ ] `ValidateTokenId(tokenId string) error`
- [ ] `ValidateValidatorAddress(address string) error`
- [ ] `ValidateGasPrice(gasPrice string) error`
- [ ] `ValidateGasLimit(gasLimit int64) error`

---

### **4. Add Metrics Methods** ⚠️

**File:** `internal/infrastructure/metrics/metrics.go`

**Cần thêm:**
- [ ] `RecordDuration(operation string, duration time.Duration)`
- [ ] `RecordSuccess(operation string, labels map[string]string)`
- [ ] `RecordFailure(operation, errorType string, labels map[string]string)`
- [ ] `RecordBlockCreated(blockNumber int64, blockHash string)`
- [ ] `RecordTransactionSubmitted(txHash string, from, to string)`
- [ ] `RecordUSCTransfer(amount string, from, to string)`
- [ ] `RecordContractDeployed(contractAddress string)`
- [ ] `RecordNFTMinted(tokenId, contractAddress string)`

---

### **5. Update Business Methods Pattern** ⚠️

**Cần update tất cả business methods để:**

1. **Add defer pattern:**
```go
start := time.Now()
defer func() {
    s.metrics.RecordDuration("operation_name", time.Since(start))
}()
```

2. **Replace manual validation với validator:**
```go
// Before
if req.WalletAddress == "" {
    return &proto.Response{ErrorMessage: "wallet address is required"}, nil
}

// After
if err := s.validator.ValidateWalletAddress(req.WalletAddress); err != nil {
    s.metrics.RecordFailure("operation_name", "validation_error", ...)
    return nil, status.Errorf(codes.InvalidArgument, "invalid wallet address: %v", err)
}
```

3. **Add failure metrics:**
```go
if err != nil {
    s.metrics.RecordFailure("operation_name", "repository_error", ...)
    return nil, status.Errorf(codes.Internal, "failed: %v", err)
}
```

4. **Add success metrics:**
```go
s.metrics.RecordSuccess("operation_name", map[string]string{
    "field": "value",
})
```

5. **Replace response errors với gRPC status:**
```go
// Before
return &proto.Response{Success: false, ErrorMessage: "error"}, nil

// After
return nil, status.Errorf(codes.InvalidArgument, "error")
```

---

### **6. Import gRPC Status Codes** ⚠️

**Cần thêm vào tất cả business service files:**
```go
import (
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)
```

---

## 🎯 **PRIORITY**

### **🔴 HIGH PRIORITY**
1. ✅ **Update Business Service Constructors** - Inject validator và metrics
2. ✅ **Update Container Injection** - Pass validator và metrics vào business services
3. ✅ **Add Validator Methods** - Blockchain-specific validation methods
4. ✅ **Add Metrics Methods** - RecordDuration, RecordSuccess, RecordFailure

### **🟡 MEDIUM PRIORITY**
5. ⚠️ **Update Business Methods Pattern** - Defer pattern, validator usage, gRPC status codes
6. ⚠️ **Import gRPC Status Codes** - Add imports

---

## 📝 **IMPLEMENTATION ORDER**

1. **Step 1**: Add validator methods cho blockchain operations
2. **Step 2**: Add metrics methods (RecordDuration, RecordSuccess, RecordFailure)
3. **Step 3**: Update business service constructors (add validator, metrics fields)
4. **Step 4**: Update container injection (pass validator, metrics)
5. **Step 5**: Update business methods pattern (defer, validator, metrics, gRPC status)

---

## ✅ **KẾT LUẬN**

**Service-04 có infrastructure nhưng chưa sử dụng:**
- ✅ Validator infrastructure: CÓ
- ✅ Metrics infrastructure: CÓ
- ❌ Business services chưa inject: THIẾU
- ❌ Business methods chưa sử dụng: THIẾU
- ❌ Validator methods chưa đầy đủ: THIẾU
- ❌ Metrics methods chưa đầy đủ: THIẾU
- ❌ Pattern chưa đúng: THIẾU

**Cần bổ sung để align với Service-22 pattern.**

---

**Last Updated**: 2025-11-12  
**Status**: ⚠️ **CẦN BỔ SUNG** - Infrastructure có nhưng chưa được sử dụng

