# Business Layer Refactor TODO - Service-04 (Theo Pattern Service-22)

**Date**: 2025-11-12  
**Purpose**: Bổ sung validator & metrics vào business layer theo pattern service-22  
**Total Groups**: 12 business service groups

---

## 📋 **TODO LIST**

### **🔧 INFRASTRUCTURE (Làm trước)**

#### **TODO-INFRA-1: Add Validator Methods** ⚠️
**File**: `internal/infrastructure/validation/validator_backend.go`

**Cần thêm các methods:**
- [ ] `ValidateWalletAddress(address string) error`
- [ ] `ValidateAmount(amount string) error`
- [ ] `ValidateBlockNumber(blockNumber int64) error`
- [ ] `ValidateTransactionHash(hash string) error`
- [ ] `ValidateContractAddress(address string) error`
- [ ] `ValidateTokenId(tokenId string) error`
- [ ] `ValidateValidatorAddress(address string) error`
- [ ] `ValidateGasPrice(gasPrice string) error`
- [ ] `ValidateGasLimit(gasLimit int64) error`

**Pattern:**
```go
func (v *Validator) ValidateWalletAddress(address string) error {
    if address == "" {
        return &validation.ValidationError{
            Field: "wallet_address",
            Message: "wallet address cannot be empty",
            Type: validation.ErrorTypeRequired,
        }
    }
    // Add format validation (0x prefix, length, etc.)
    return nil
}
```

---

#### **TODO-INFRA-2: Add Metrics Methods** ⚠️
**File**: `internal/infrastructure/metrics/metrics.go`

**Cần thêm các methods:**
- [ ] `RecordDuration(operation string, duration time.Duration)`
- [ ] `RecordSuccess(operation string, labels map[string]string)`
- [ ] `RecordFailure(operation, errorType string, labels map[string]string)`
- [ ] `RecordBlockCreated(blockNumber int64, blockHash string)`
- [ ] `RecordTransactionSubmitted(txHash string, from, to string)`
- [ ] `RecordUSCTransfer(amount string, from, to string)`
- [ ] `RecordContractDeployed(contractAddress string)`
- [ ] `RecordNFTMinted(tokenId, contractAddress string)`

**Pattern:**
```go
func (m *MetricsService) RecordDuration(operation string, duration time.Duration) {
    m.performanceMetrics.RecordRequest(duration, false)
}

func (m *MetricsService) RecordSuccess(operation string, labels map[string]string) {
    // Record success metric
}

func (m *MetricsService) RecordFailure(operation, errorType string, labels map[string]string) {
    // Record failure metric
}
```

---

### **📦 CONTAINER UPDATE**

#### **TODO-CONTAINER-1: Update Container Injection** ❌
**File**: `internal/application/container-grpc.go`

**Cần update `initializeBusiness()` method (line 260-300):**

**Before:**
```go
c.BlockService = blockbiz.NewService(
    c.BlockRepository,
    c.cosmosApp,
    c.blockchainStorage,
    c.logger
)
```

**After:**
```go
c.BlockService = blockbiz.NewService(
    c.BlockRepository,
    c.cosmosApp,
    c.blockchainStorage,
    c.logger,
    c.validator,  // ✅ Add
    c.metrics,    // ✅ Add
)
```

**Cần update cho 12 services:**
- [ ] Block Operations
- [ ] Transaction Operations
- [ ] USC Coin Operations
- [ ] Smart Contract Operations
- [ ] NFT Token Operations
- [ ] Custom Token Operations
- [ ] Product Certificate Operations
- [ ] Validator Operations
- [ ] Network Operations
- [ ] Streaming Operations
- [ ] Store Bridge Operations
- [ ] Store Network Operations

---

### **🎯 BUSINESS SERVICES (12 Groups)**

#### **TODO-GROUP-1: Block Operations** ❌
**File**: `internal/application/business/block_operations/block_operations_service.go`

**Tasks:**
1. [ ] Update `Service` struct: Add `validator *validation.Validator` và `metrics *metrics.MetricsService`
2. [ ] Update `NewService()`: Add validator và metrics parameters
3. [ ] Update `ProduceBlock()`:
   - [ ] Add `defer func() { s.metrics.RecordDuration("produce_block", time.Since(start)) }()`
   - [ ] Replace manual validation với `s.validator.ValidateBlockNumber()`
   - [ ] Add `s.metrics.RecordFailure()` cho validation errors
   - [ ] Add `s.metrics.RecordFailure()` cho repository errors
   - [ ] Add `s.metrics.RecordSuccess()` cho success
   - [ ] Replace response errors với `status.Errorf(codes.InvalidArgument, ...)`
4. [ ] Update `ValidateBlock()`: Same pattern
5. [ ] Update `GetBlock()`: Same pattern
6. [ ] Update `GetBlockByHash()`: Same pattern
7. [ ] Update `GetLatestBlock()`: Same pattern
8. [ ] Update `GetBlockRange()`: Same pattern
9. [ ] Add imports: `google.golang.org/grpc/codes`, `google.golang.org/grpc/status`

---

#### **TODO-GROUP-2: Transaction Operations** ❌
**File**: `internal/application/business/transaction_operations/transaction_operations_service.go`

**Tasks:**
1. [ ] Update `Service` struct: Add `validator *validation.Validator` và `metrics *metrics.MetricsService`
2. [ ] Update `NewService()`: Add validator và metrics parameters
3. [ ] Update `SubmitTransaction()`:
   - [ ] Add `defer func() { s.metrics.RecordDuration("submit_transaction", time.Since(start)) }()`
   - [ ] Replace manual validation với `s.validator.ValidateWalletAddress()`, `s.validator.ValidateAmount()`
   - [ ] Add `s.metrics.RecordFailure()` cho validation errors
   - [ ] Add `s.metrics.RecordFailure()` cho repository errors
   - [ ] Add `s.metrics.RecordSuccess()` cho success
   - [ ] Replace response errors với `status.Errorf(codes.InvalidArgument, ...)`
4. [ ] Update `GetTransaction()`: Same pattern
5. [ ] Update `GetTransactionStatus()`: Same pattern
6. [ ] Update `GetPendingTransactions()`: Same pattern
7. [ ] Update `EstimateTransactionFee()`: Same pattern
8. [ ] Add imports: `google.golang.org/grpc/codes`, `google.golang.org/grpc/status`

---

#### **TODO-GROUP-3: USC Coin Operations** ❌
**File**: `internal/application/business/usc_coin_operations/usc_coin_operations_service.go`

**Tasks:**
1. [ ] Update `Service` struct: Add `validator *validation.Validator` và `metrics *metrics.MetricsService`
2. [ ] Update `NewService()`: Add validator và metrics parameters
3. [ ] Update `GetUSCBalance()`:
   - [ ] Add `defer func() { s.metrics.RecordDuration("get_usc_balance", time.Since(start)) }()`
   - [ ] Replace manual validation với `s.validator.ValidateWalletAddress()`
   - [ ] Add `s.metrics.RecordFailure()` cho validation errors
   - [ ] Add `s.metrics.RecordFailure()` cho repository errors
   - [ ] Add `s.metrics.RecordSuccess()` cho success
   - [ ] Replace response errors với `status.Errorf(codes.InvalidArgument, ...)`
4. [ ] Update `TransferUSC()`:
   - [ ] Add `defer func() { s.metrics.RecordDuration("transfer_usc", time.Since(start)) }()`
   - [ ] Replace manual validation với `s.validator.ValidateWalletAddress()`, `s.validator.ValidateAmount()`
   - [ ] Add `s.metrics.RecordUSCTransfer()` cho success
   - [ ] Add `s.metrics.RecordFailure()` cho errors
   - [ ] Replace response errors với `status.Errorf(codes.InvalidArgument, ...)`
5. [ ] Update `GetUSCSupply()`: Same pattern
6. [ ] Update `GetTransactionHistory()`: Same pattern
7. [ ] Update `GetUSCTransactions()`: Same pattern
8. [ ] Add imports: `google.golang.org/grpc/codes`, `google.golang.org/grpc/status`

---

#### **TODO-GROUP-4: Smart Contract Operations** ❌
**File**: `internal/application/business/smart_contract_operations/smart_contract_operations_service.go`

**Tasks:**
1. [ ] Update `Service` struct: Add `validator *validation.Validator` và `metrics *metrics.MetricsService`
2. [ ] Update `NewService()`: Add validator và metrics parameters
3. [ ] Update `DeployContract()`:
   - [ ] Add `defer func() { s.metrics.RecordDuration("deploy_contract", time.Since(start)) }()`
   - [ ] Replace manual validation với `s.validator.ValidateWalletAddress()`, `s.validator.ValidateContractAddress()`
   - [ ] Add `s.metrics.RecordFailure()` cho validation errors
   - [ ] Add `s.metrics.RecordContractDeployed()` cho success
   - [ ] Replace response errors với `status.Errorf(codes.InvalidArgument, ...)`
4. [ ] Update `ExecuteContract()`: Same pattern
5. [ ] Update `QueryContract()`: Same pattern
6. [ ] Update `GetContractCode()`: Same pattern
7. [ ] Update `GetContractStorage()`: Same pattern
8. [ ] Add imports: `google.golang.org/grpc/codes`, `google.golang.org/grpc/status`

---

#### **TODO-GROUP-5: NFT Token Operations** ❌
**File**: `internal/application/business/nft_token_operations/nft_token_operations_service.go`

**Tasks:**
1. [ ] Update `Service` struct: Add `validator *validation.Validator` và `metrics *metrics.MetricsService`
2. [ ] Update `NewService()`: Add validator và metrics parameters
3. [ ] Update `MintNFT()`:
   - [ ] Add `defer func() { s.metrics.RecordDuration("mint_nft", time.Since(start)) }()`
   - [ ] Replace manual validation với `s.validator.ValidateContractAddress()`, `s.validator.ValidateWalletAddress()`
   - [ ] Add `s.metrics.RecordNFTMinted()` cho success
   - [ ] Add `s.metrics.RecordFailure()` cho errors
   - [ ] Replace response errors với `status.Errorf(codes.InvalidArgument, ...)`
4. [ ] Update `TransferNFT()`: Same pattern
5. [ ] Update `BurnNFT()`: Same pattern
6. [ ] Update `GetNFTInfo()`: Same pattern
7. [ ] Update `GetNFTsByOwner()`: Same pattern
8. [ ] Update `DeployNFTContract()`: Same pattern
9. [ ] Update `CreateNFTCollection()`: Same pattern
10. [ ] Add imports: `google.golang.org/grpc/codes`, `google.golang.org/grpc/status`

---

#### **TODO-GROUP-6: Custom Token Operations** ❌
**File**: `internal/application/business/custom_token_operations/custom_token_operations_service.go`

**Tasks:**
1. [ ] Update `Service` struct: Add `validator *validation.Validator` và `metrics *metrics.MetricsService`
2. [ ] Update `NewService()`: Add validator và metrics parameters
3. [ ] Update `CreateBlockchainToken()`:
   - [ ] Add `defer func() { s.metrics.RecordDuration("create_token", time.Since(start)) }()`
   - [ ] Replace manual validation với `s.validator.ValidateWalletAddress()`, `s.validator.ValidateAmount()`
   - [ ] Add `s.metrics.RecordFailure()` cho validation errors
   - [ ] Add `s.metrics.RecordSuccess()` cho success
   - [ ] Replace response errors với `status.Errorf(codes.InvalidArgument, ...)`
4. [ ] Update `MintTokens()`: Same pattern
5. [ ] Update `BurnTokens()`: Same pattern
6. [ ] Update `GetTokenInfo()`: Same pattern
7. [ ] Update `GetTokenBalance()`: Same pattern
8. [ ] Add imports: `google.golang.org/grpc/codes`, `google.golang.org/grpc/status`

---

#### **TODO-GROUP-7: Product Certificate Operations** ❌
**File**: `internal/application/business/product_certificate_operations/product_certificate_operations_service.go`

**Tasks:**
1. [ ] Update `Service` struct: Add `validator *validation.Validator` và `metrics *metrics.MetricsService`
2. [ ] Update `NewService()`: Add validator và metrics parameters
3. [ ] Update `CreateProductCertificate()`:
   - [ ] Add `defer func() { s.metrics.RecordDuration("create_certificate", time.Since(start)) }()`
   - [ ] Replace manual validation với `s.validator.ValidateWalletAddress()`
   - [ ] Add `s.metrics.RecordFailure()` cho validation errors
   - [ ] Add `s.metrics.RecordSuccess()` cho success
   - [ ] Replace response errors với `status.Errorf(codes.InvalidArgument, ...)`
4. [ ] Update `VerifyBlockchainProductCertificate()`: Same pattern
5. [ ] Update `TransferProductOwnership()`: Same pattern
6. [ ] Add imports: `google.golang.org/grpc/codes`, `google.golang.org/grpc/status`

---

#### **TODO-GROUP-8: Validator Operations** ❌
**File**: `internal/application/business/validator_operations/validator_operations_service.go`

**Tasks:**
1. [ ] Update `Service` struct: Add `validator *validation.Validator` và `metrics *metrics.MetricsService`
2. [ ] Update `NewService()`: Add validator và metrics parameters
3. [ ] Update `RegisterValidator()`:
   - [ ] Add `defer func() { s.metrics.RecordDuration("register_validator", time.Since(start)) }()`
   - [ ] Replace manual validation với `s.validator.ValidateValidatorAddress()`
   - [ ] Add `s.metrics.RecordFailure()` cho validation errors
   - [ ] Add `s.metrics.RecordSuccess()` cho success
   - [ ] Replace response errors với `status.Errorf(codes.InvalidArgument, ...)`
4. [ ] Update `GetValidators()`: Same pattern
5. [ ] Update `GetValidatorStatus()`: Same pattern
6. [ ] Update `StakeUSC()`: Same pattern
7. [ ] Update `UnstakeUSC()`: Same pattern
8. [ ] Add imports: `google.golang.org/grpc/codes`, `google.golang.org/grpc/status`

---

#### **TODO-GROUP-9: Network Operations** ❌
**File**: `internal/application/business/network_operations/network_operations_service.go`

**Tasks:**
1. [ ] Update `Service` struct: Add `validator *validation.Validator` và `metrics *metrics.MetricsService`
2. [ ] Update `NewService()`: Add validator và metrics parameters
3. [ ] Update `GetNetworkInfo()`:
   - [ ] Add `defer func() { s.metrics.RecordDuration("get_network_info", time.Since(start)) }()`
   - [ ] Add `s.metrics.RecordFailure()` cho errors
   - [ ] Add `s.metrics.RecordSuccess()` cho success
   - [ ] Replace response errors với `status.Errorf(codes.Internal, ...)`
4. [ ] Update `GetChainInfo()`: Same pattern
5. [ ] Update `GetPeers()`: Same pattern
6. [ ] Update `GetNetworkStats()`: Same pattern
7. [ ] Add imports: `google.golang.org/grpc/codes`, `google.golang.org/grpc/status`

---

#### **TODO-GROUP-10: Streaming Operations** ❌
**File**: `internal/application/business/streaming_operations/streaming_operations_service.go`

**Tasks:**
1. [ ] Update `Service` struct: Add `validator *validation.Validator` và `metrics *metrics.MetricsService`
2. [ ] Update `NewService()`: Add validator và metrics parameters
3. [ ] Update `StreamBlocks()`:
   - [ ] Add metrics recording cho stream operations
   - [ ] Add `s.metrics.RecordFailure()` cho errors
   - [ ] Replace response errors với `status.Errorf(codes.Internal, ...)`
4. [ ] Update `StreamTransactions()`: Same pattern
5. [ ] Update `StreamValidatorEvents()`: Same pattern
6. [ ] Update `StreamNetworkEvents()`: Same pattern
7. [ ] Add imports: `google.golang.org/grpc/codes`, `google.golang.org/grpc/status`

---

#### **TODO-GROUP-11: Store Bridge Operations** ❌
**File**: `internal/application/business/store_bridge_operations/store_bridge_operations_service.go`

**Tasks:**
1. [ ] Update `Service` struct: Add `validator *validation.Validator` và `metrics *metrics.MetricsService`
2. [ ] Update `NewService()`: Add validator và metrics parameters
3. [ ] Update `DeployStoreBridge()`:
   - [ ] Add `defer func() { s.metrics.RecordDuration("deploy_bridge", time.Since(start)) }()`
   - [ ] Replace manual validation với `s.validator.ValidateWalletAddress()`
   - [ ] Add `s.metrics.RecordFailure()` cho validation errors
   - [ ] Add `s.metrics.RecordSuccess()` cho success
   - [ ] Replace response errors với `status.Errorf(codes.InvalidArgument, ...)`
4. [ ] Update `RegisterStoreNetwork()`: Same pattern
5. [ ] Update `BridgeStoreTokenToUSC()`: Same pattern
6. [ ] Update `BridgeUSCToStoreToken()`: Same pattern
7. [ ] Update `GetStoreBridgeMetrics()`: Same pattern
8. [ ] Update `ValidateStoreBridge()`: Same pattern
9. [ ] Add imports: `google.golang.org/grpc/codes`, `google.golang.org/grpc/status`

---

#### **TODO-GROUP-12: Store Network Operations** ❌
**File**: `internal/application/business/store_network_operations/store_network_operations_service.go`

**Tasks:**
1. [ ] Update `Service` struct: Add `validator *validation.Validator` và `metrics *metrics.MetricsService`
2. [ ] Update `NewService()`: Add validator và metrics parameters
3. [ ] Update `SyncStoreNetworkState()`:
   - [ ] Add `defer func() { s.metrics.RecordDuration("sync_network_state", time.Since(start)) }()`
   - [ ] Replace manual validation với `s.validator.ValidateXXX()`
   - [ ] Add `s.metrics.RecordFailure()` cho validation errors
   - [ ] Add `s.metrics.RecordSuccess()` cho success
   - [ ] Replace response errors với `status.Errorf(codes.InvalidArgument, ...)`
4. [ ] Update `GetStoreNetworkInfo()`: Same pattern
5. [ ] Update `UpdateStoreBridgeConfig()`: Same pattern
6. [ ] Add imports: `google.golang.org/grpc/codes`, `google.golang.org/grpc/status`

---

## 📝 **PATTERN TEMPLATE (Service-22 Style)**

### **Service Struct:**
```go
type Service struct {
    repo              *repo.Repository
    cosmosApp         *app.USCApp
    blockchainStorage *storage.StateManager
    logger            *logging.Logger
    validator         *validation.Validator      // ✅ Add
    metrics           *metrics.MetricsService    // ✅ Add
}
```

### **NewService Constructor:**
```go
func NewService(
    repo *repo.Repository,
    cosmosApp *app.USCApp,
    blockchainStorage *storage.StateManager,
    logger *logging.Logger,
    validator *validation.Validator,      // ✅ Add
    metrics *metrics.MetricsService,      // ✅ Add
) *Service {
    return &Service{
        repo:              repo,
        cosmosApp:         cosmosApp,
        blockchainStorage: blockchainStorage,
        logger:            logger,
        validator:         validator,      // ✅ Add
        metrics:           metrics,        // ✅ Add
    }
}
```

### **Method Pattern:**
```go
func (s *Service) OperationName(ctx context.Context, req *proto.Request) (*proto.Response, error) {
    start := time.Now()
    defer func() {
        s.metrics.RecordDuration("operation_name", time.Since(start))
    }()

    correlationID := utils.GetCorrelationID(ctx)
    s.logger.Info("Operation description",
        logging.String("correlation_id", correlationID),
        logging.String("field", req.Field))

    // 1. Validation
    if err := s.validator.ValidateField(req.Field); err != nil {
        s.logger.Error("Validation failed",
            logging.String("correlation_id", correlationID),
            logging.Error(err))
        s.metrics.RecordFailure("operation_name", "validation_error", map[string]string{
            "field": req.Field,
        })
        return nil, status.Errorf(codes.InvalidArgument, "invalid field: %v", err)
    }

    // 2. Repository call
    resp, err := s.repo.OperationName(ctx, req)
    if err != nil {
        s.logger.Error("Repository operation failed",
            logging.String("correlation_id", correlationID),
            logging.Error(err))
        s.metrics.RecordFailure("operation_name", "repository_error", map[string]string{
            "field": req.Field,
        })
        return nil, status.Errorf(codes.Internal, "failed to operation: %v", err)
    }

    // 3. Success
    s.logger.Info("Operation completed successfully",
        logging.String("correlation_id", correlationID))
    s.metrics.RecordSuccess("operation_name", map[string]string{
        "field": req.Field,
    })

    return resp, nil
}
```

### **Imports:**
```go
import (
    "context"
    "time"

    "google.golang.org/grpc/codes"      // ✅ Add
    "google.golang.org/grpc/status"     // ✅ Add
    "service-04/internal/infrastructure/metrics"     // ✅ Add
    "service-04/internal/infrastructure/validation" // ✅ Add
    // ... other imports
)
```

---

## 🎯 **IMPLEMENTATION ORDER**

### **Phase 1: Infrastructure (Làm trước)**
1. ✅ TODO-INFRA-1: Add Validator Methods
2. ✅ TODO-INFRA-2: Add Metrics Methods

### **Phase 2: Container Update**
3. ✅ TODO-CONTAINER-1: Update Container Injection

### **Phase 3: Business Services (12 Groups)**
4. ✅ TODO-GROUP-1: Block Operations
5. ✅ TODO-GROUP-2: Transaction Operations
6. ✅ TODO-GROUP-3: USC Coin Operations
7. ✅ TODO-GROUP-4: Smart Contract Operations
8. ✅ TODO-GROUP-5: NFT Token Operations
9. ✅ TODO-GROUP-6: Custom Token Operations
10. ✅ TODO-GROUP-7: Product Certificate Operations
11. ✅ TODO-GROUP-8: Validator Operations
12. ✅ TODO-GROUP-9: Network Operations
13. ✅ TODO-GROUP-10: Streaming Operations
14. ✅ TODO-GROUP-11: Store Bridge Operations
15. ✅ TODO-GROUP-12: Store Network Operations

---

## ✅ **CHECKLIST PER GROUP**

Mỗi group cần hoàn thành:
- [ ] Update Service struct (add validator, metrics)
- [ ] Update NewService() constructor (add parameters)
- [ ] Update tất cả methods:
  - [ ] Add defer pattern cho metrics
  - [ ] Replace manual validation với validator service
  - [ ] Add RecordFailure() cho errors
  - [ ] Add RecordSuccess() cho success
  - [ ] Replace response errors với gRPC status codes
- [ ] Add imports (codes, status, validation, metrics)
- [ ] Test methods hoạt động đúng

---

**Last Updated**: 2025-11-12  
**Status**: 📋 **TODO CREATED** - Sẵn sàng bắt đầu implementation

