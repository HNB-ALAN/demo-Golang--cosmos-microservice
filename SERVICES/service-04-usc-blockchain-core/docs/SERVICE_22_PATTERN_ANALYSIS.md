# 📊 **PHÂN TÍCH PATTERN: SERVICE-22 (KAFKA MESSAGING)**

**Ngày phân tích**: 2025-11-12  
**Mục đích**: Hiểu cách service-22 tổ chức code để áp dụng pattern đúng cho service-04

---

## 🎯 **KIẾN TRÚC 3 LAYER CỦA SERVICE-22**

### **1. Handlers Layer** (Thin Wrapper)

**Trách nhiệm**:
- ✅ Nhận gRPC requests
- ✅ Gọi business service
- ✅ Return gRPC responses
- ❌ **KHÔNG** có business logic
- ❌ **KHÔNG** có data access logic

**Ví dụ** (`consumer_group_handlers.go`):
```go
func (h *Handlers) CreateConsumerGroup(ctx context.Context, req *proto.CreateConsumerGroupRequest) (*proto.CreateConsumerGroupResponse, error) {
    return h.service.CreateConsumerGroup(ctx, req)
}
```

**Pattern**: Handlers chỉ là thin wrapper, delegate ngay lập tức đến business service.

---

### **2. Business Layer** (Orchestration)

**Trách nhiệm**:
- ✅ **Input validation** (sử dụng validator)
- ✅ **Logging** (business operations)
- ✅ **Metrics recording** (performance, success/failure)
- ✅ **Orchestration** (delegate to repository)
- ❌ **KHÔNG** tương tác trực tiếp với infrastructure (Kafka, DB, Redis)
- ❌ **KHÔNG** có data access logic

**Ví dụ** (`consumer_group_service.go`):
```go
func (s *Service) CreateConsumerGroup(ctx context.Context, req *proto.CreateConsumerGroupRequest) (*proto.CreateConsumerGroupResponse, error) {
    start := time.Now()
    defer func() {
        s.metrics.RecordDuration("create_consumer_group", time.Since(start))
    }()

    // 1. Input validation
    if err := s.validator.ValidateConsumerGroup(req.GroupId); err != nil {
        s.log.Error("Consumer group validation failed", ...)
        s.metrics.RecordFailure("create_consumer_group", "validation_error", ...)
        return &proto.CreateConsumerGroupResponse{Success: false, ErrorMessage: err.Error()}, nil
    }

    // 2. Log business operation
    s.log.Info("Creating consumer group", ...)

    // 3. Call repository (DELEGATE - không có logic tương tác trực tiếp với Kafka/DB)
    resp, err := s.repo.CreateConsumerGroup(ctx, req)
    if err != nil {
        s.log.Error("Failed to create consumer group", ...)
        s.metrics.RecordFailure("create_consumer_group", "repository_error", ...)
        return &proto.CreateConsumerGroupResponse{Success: false, ErrorMessage: err.Error()}, nil
    }

    // 4. Log success
    s.log.Info("Consumer group created successfully", ...)

    // 5. Record success metrics
    s.metrics.RecordConsumerGroupCreated(req.GroupId, req.OwnerService)

    return resp, nil
}
```

**Pattern**: Business layer chỉ orchestrate, validate, log, và record metrics. **KHÔNG** có logic tương tác trực tiếp với infrastructure.

---

### **3. Repository Layer** (Data Access)

**Trách nhiệm**:
- ✅ **Single source of truth** cho data access
- ✅ Tương tác với **Kafka** (publish messages, manage topics)
- ✅ Tương tác với **PostgreSQL** (persist metadata)
- ✅ Tương tác với **Redis** (caching)
- ✅ Priority-based fallback (Kafka → PostgreSQL)
- ✅ Dual-write pattern (Kafka + PostgreSQL)

**Ví dụ** (`consumer_group_repository.go`):
```go
func (r *Repository) CreateConsumerGroup(ctx context.Context, req *proto.CreateConsumerGroupRequest) (*proto.CreateConsumerGroupResponse, error) {
    // 1. Check cache (Redis)
    if r.redisManager != nil {
        cacheKey := fmt.Sprintf("kafka:consumer_group:%s", req.GroupId)
        if cached, err := r.redisManager.Get(ctx, cacheKey); err == nil && cached != "" {
            // Return cached data
            ...
        }
    }

    // 2. Insert into PostgreSQL
    err = r.db.GetPostgres().QueryRowContext(ctx, `
        INSERT INTO kafka_consumer_groups (...)
        VALUES (...)
        ON CONFLICT (group_id) 
        DO UPDATE SET updated_at = NOW()
        RETURNING ...
    `, ...).Scan(...)

    // 3. Cache in Redis
    if r.redisManager != nil {
        cacheKey := fmt.Sprintf("kafka:consumer_group:%s", groupID)
        _ = r.redisManager.Set(ctx, cacheKey, string(groupInfoJSON), 24*time.Hour)
    }

    return &proto.CreateConsumerGroupResponse{...}, nil
}
```

**Pattern**: Repository là **single source of truth** cho tất cả data access operations. Business layer **KHÔNG** tương tác trực tiếp với infrastructure.

---

## 🔍 **SO SÁNH VỚI SERVICE-04**

### **SERVICE-22 (ĐÚNG)** ✅

```
Handlers → Business → Repository → Infrastructure (Kafka/DB/Redis)
   ↓          ↓            ↓
Thin      Validate     Data Access
Wrapper   + Log        (Single Source
         + Metrics     of Truth)
         + Delegate
```

**Đặc điểm**:
- ✅ Business layer **KHÔNG** có logic tương tác trực tiếp với infrastructure
- ✅ Repository là **single source of truth** cho data access
- ✅ Clear separation of concerns
- ✅ No code duplication

---

### **SERVICE-04 (SAI - CÓ TRÙNG LẶP)** ❌

```
Handlers → Business → Repository → Infrastructure (Keeper/DB)
   ↓          ↓            ↓
Thin      Validate     Data Access
Wrapper   + Log        (Keeper + DB)
         + Metrics
         + TRỰC TIẾP   ← ❌ DUPLICATE!
         tương tác
         với Keeper
```

**Đặc điểm**:
- ❌ Business layer **CÓ** logic tương tác trực tiếp với Keeper (duplicate với Repository)
- ❌ Repository cũng tương tác với Keeper
- ❌ Code duplication (~15+ hàm trùng lặp)
- ❌ Confusion về responsibility

---

## 🛠️ **GIẢI PHÁP CHO SERVICE-04**

### **Refactor Business Layer**

**Loại bỏ tất cả logic tương tác trực tiếp với Keeper từ Business layer**:

1. **Product Certificate Operations**:
   - ❌ Xóa `verifyCertificateFromKeeper`
   - ❌ Xóa `transferOwnershipOnBlockchain`
   - ❌ Xóa `createCertificateOnBlockchain` (dead code)
   - ✅ Chỉ giữ: validation, logging, metrics, delegate to repository

2. **Validator Operations**:
   - ❌ Xóa `registerValidatorOnBlockchain`
   - ❌ Xóa `getValidatorsFromKeeper`
   - ❌ Xóa `getValidatorStatusFromKeeper`
   - ✅ Chỉ giữ: validation, logging, metrics, delegate to repository

3. **Custom Token Operations**:
   - ❌ Xóa tất cả `*OnBlockchain` và `*FromKeeper` methods
   - ✅ Chỉ giữ: validation, logging, metrics, delegate to repository

4. **NFT Token Operations**:
   - ❌ Xóa `mintNFTOnBlockchain`, `getNFTInfoFromKeeper`, etc.
   - ✅ Chỉ giữ: validation, logging, metrics, delegate to repository

5. **Smart Contract Operations**:
   - ❌ Xóa `deployContractOnBlockchain`, `executeContractOnBlockchain`, etc.
   - ✅ Chỉ giữ: validation, logging, metrics, delegate to repository

6. **Transaction Operations**:
   - ❌ Xóa `submitTransactionOnBlockchain`, `getTransactionFromKeeper`, etc.
   - ✅ Chỉ giữ: validation, logging, metrics, delegate to repository

7. **Block Operations**:
   - ❌ Xóa `getBlockFromKeeper`, `getBlockByHashFromKeeper`, etc.
   - ✅ Chỉ giữ: validation, logging, metrics, delegate to repository

---

## 📋 **PATTERN ĐÚNG SAU KHI REFACTOR**

### **Business Layer Pattern** (Giống Service-22)

```go
func (s *Service) CreateProductCertificate(ctx context.Context, req *proto.CreateProductCertificateRequest) (*proto.CreateProductCertificateResponse, error) {
    start := time.Now()
    defer func() {
        s.metrics.RecordDuration("create_product_certificate", time.Since(start))
    }()

    // 1. Input validation
    if req.ProductId == "" {
        s.logger.Warn("Empty product ID", ...)
        s.metrics.RecordFailure("create_product_certificate", "validation_error", ...)
        return &proto.CreateProductCertificateResponse{Status: 2, ErrorMessage: "product_id is required"}, nil
    }

    // 2. Log business operation
    s.logger.Info("Creating product certificate", ...)

    // 3. Call repository (DELEGATE - không có logic tương tác trực tiếp với Keeper)
    resp, err := s.repo.CreateProductCertificate(ctx, req)
    if err != nil {
        s.logger.Error("Failed to create product certificate", ...)
        s.metrics.RecordFailure("create_product_certificate", "repository_error", ...)
        return &proto.CreateProductCertificateResponse{Status: 2, ErrorMessage: err.Error()}, nil
    }

    // 4. Log success
    s.logger.Info("Product certificate created successfully", ...)

    // 5. Record success metrics
    s.metrics.RecordProductCertificateCreated(resp.CertificateId, req.ProductId)

    return resp, nil
}
```

**Key Points**:
- ✅ **KHÔNG** có `createCertificateOnBlockchain` hoặc `createCertificateOnKeeper`
- ✅ **KHÔNG** có logic tương tác trực tiếp với `cosmosApp.ProductCertificateKeeper`
- ✅ Chỉ có: validation, logging, metrics, delegate to repository
- ✅ Repository là **single source of truth** cho data access

---

## ✅ **KẾT LUẬN**

**Service-22 sử dụng pattern đúng**:
- Handlers: Thin wrapper
- Business: Orchestration (validation, logging, metrics, delegate)
- Repository: Data access (single source of truth)

**Service-04 cần refactor** để follow pattern tương tự:
- Loại bỏ tất cả logic tương tác trực tiếp với Keeper từ Business layer
- Repository là **single source of truth** cho tất cả data access operations
- Business layer chỉ orchestrate, validate, log, và record metrics

---

**Last Updated**: 2025-11-12  
**Status**: ✅ **PATTERN ĐÚNG ĐÃ XÁC ĐỊNH** - Cần refactor Service-04

