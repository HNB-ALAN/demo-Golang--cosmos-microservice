# 🔧 **ERROR STANDARDIZATION GUIDE**

## 📊 **OVERVIEW**

Hướng dẫn chuẩn hóa error messages trong repository layer để đảm bảo consistency và dễ bảo trì.

**Date**: 2025-01-05  
**Status**: ✅ **READY FOR IMPLEMENTATION**

---

## ✅ **ERROR INFRASTRUCTURE CREATED**

### **File**: `internal/application/repository/errors.go`

**Features**:
- ✅ Error code constants (2000-3299 range)
- ✅ Error messages map
- ✅ Helper functions for common patterns
- ✅ Retryable error detection
- ✅ Integration with shared errors package

**Error Code Ranges**:
- **2000-2099**: Block operations errors
- **2100-2199**: Transaction operations errors
- **2200-2299**: USC Coin operations errors
- **2300-2399**: Smart Contract operations errors
- **2400-2499**: NFT Token operations errors
- **2500-2599**: Custom Token operations errors
- **2600-2699**: Product Certificate operations errors
- **2700-2799**: Validator operations errors
- **2800-2899**: Network operations errors
- **2900-2999**: Streaming operations errors
- **3000-3099**: Store Bridge operations errors
- **3100-3199**: Store Network operations errors
- **3200-3299**: Common repository errors

---

## 📋 **IMPLEMENTATION PATTERN**

### **Pattern 1: Use Error Constants** (Recommended)

**Before**:
```go
if blockNumber <= 0 {
    return &proto.GetBlockResponse{}, fmt.Errorf("block_number must be greater than 0")
}
```

**After**:
```go
import repoerrors "service-04/internal/application/repository"

if blockNumber <= 0 {
    return nil, repoerrors.NewValidationError("block_number", "must be greater than 0")
}
```

---

### **Pattern 2: Wrap Existing Errors**

**Before**:
```go
if err != nil {
    return nil, fmt.Errorf("failed to get SDK context: %w", err)
}
```

**After**:
```go
import repoerrors "service-04/internal/application/repository"

if err != nil {
    return nil, repoerrors.NewBlockchainError("get_sdk_context", err)
}
```

---

### **Pattern 3: Not Found Errors**

**Before**:
```go
if block == nil {
    return nil, fmt.Errorf("block not found: block_number=%d", blockNumber)
}
```

**After**:
```go
import repoerrors "service-04/internal/application/repository"

if block == nil {
    return nil, repoerrors.NewNotFoundError("block", fmt.Sprintf("block_number=%d", blockNumber))
}
```

---

## 🎯 **IMPLEMENTATION PRIORITY**

### **Priority 1: Common Patterns** 🔴

**Focus**: Replace common error patterns first

**Files to Update**:
1. `block_operations_repository.go` - ~16 error patterns
2. `transaction_operations_repository.go` - ~7 error patterns
3. `usc_coin_operations_repository.go` - Error patterns
4. `smart_contract_operations_repository.go` - Error patterns
5. `nft_token_operations_repository.go` - Error patterns
6. `custom_token_operations_repository.go` - Error patterns

**Estimated Errors**: ~50-70 error statements

---

### **Priority 2: Validation Errors** 🟡

**Focus**: Standardize validation error messages

**Pattern**:
```go
// Before
if req.RequiredField == "" {
    return nil, fmt.Errorf("required_field is required")
}

// After
if req.RequiredField == "" {
    return nil, repoerrors.NewValidationError("required_field", "is required")
}
```

**Estimated Errors**: ~30-40 validation errors

---

### **Priority 3: Database/Blockchain Errors** 🟢

**Focus**: Wrap database and blockchain errors

**Pattern**:
```go
// Before
if err != nil {
    return nil, fmt.Errorf("failed to query database: %w", err)
}

// After
if err != nil {
    return nil, repoerrors.NewDatabaseError("query", err)
}
```

**Estimated Errors**: ~20-30 database/blockchain errors

---

## 📝 **EXAMPLE IMPLEMENTATION**

### **Example: block_operations_repository.go**

**Before**:
```go
func (r *Repository) GetBlock(ctx context.Context, req *proto.GetBlockRequest) (*proto.GetBlockResponse, error) {
    if req.BlockNumber <= 0 {
        return &proto.GetBlockResponse{}, fmt.Errorf("block_number must be greater than 0")
    }
    
    block, err := r.getBlockFromKeeper(ctx, req.BlockNumber)
    if err != nil {
        return nil, fmt.Errorf("block not found in keeper: block_number=%d, error=%w", req.BlockNumber, err)
    }
    
    if block == nil {
        return nil, fmt.Errorf("block is zero value from keeper: block_number=%d", req.BlockNumber)
    }
    
    return &proto.GetBlockResponse{
        Block: convertBlockToProto(block),
    }, nil
}
```

**After**:
```go
import repoerrors "service-04/internal/application/repository"

func (r *Repository) GetBlock(ctx context.Context, req *proto.GetBlockRequest) (*proto.GetBlockResponse, error) {
    correlationID := utils.GetCorrelationID(ctx)
    
    if req.BlockNumber <= 0 {
        return nil, repoerrors.NewValidationError("block_number", "must be greater than 0")
    }
    
    block, err := r.getBlockFromKeeper(ctx, req.BlockNumber)
    if err != nil {
        return nil, repoerrors.WrapRepositoryError(ErrBlockNotFound, err, 
            fmt.Sprintf("block_number=%d", req.BlockNumber))
    }
    
    if block == nil {
        return nil, repoerrors.NewNotFoundError("block", fmt.Sprintf("block_number=%d", req.BlockNumber))
    }
    
    return &proto.GetBlockResponse{
        Block: convertBlockToProto(block),
    }, nil
}
```

---

## 🔧 **AUTOMATION SCRIPT**

### **Script to Replace Common Patterns** (Optional)

**Note**: Manual review recommended for accuracy

```bash
#!/bin/bash
# Replace common error patterns with standardized errors
# Usage: ./standardize_errors.sh

# Pattern 1: Validation errors
sed -i 's/fmt\.Errorf("\(.*\)_is required")/repoerrors.NewValidationError("\1", "is required")/g' \
    internal/application/repository/*/*_repository.go

# Pattern 2: Not found errors
sed -i 's/fmt\.Errorf("\(.*\) not found")/repoerrors.NewNotFoundError("\1", identifier)/g' \
    internal/application/repository/*/*_repository.go

# Pattern 3: Database errors
sed -i 's/fmt\.Errorf("failed to \(.*\): %w", err)/repoerrors.NewDatabaseError("\1", err)/g' \
    internal/application/repository/*/*_repository.go
```

---

## 📊 **BENEFITS**

### **1. Consistency** ✅
- All errors use same format
- Error codes are standardized
- Easy to identify error types

### **2. Maintainability** ✅
- Centralized error definitions
- Easy to update error messages
- Clear error categorization

### **3. Observability** ✅
- Error codes for monitoring
- Retryable error detection
- Better error tracking

---

## 🎯 **IMPLEMENTATION CHECKLIST**

### **Phase 1: Setup** ✅
- [x] Create error constants file
- [x] Define error code ranges
- [x] Create helper functions
- [x] Document implementation pattern

### **Phase 2: Implementation** ⏳
- [ ] Replace errors in block_operations_repository.go
- [ ] Replace errors in transaction_operations_repository.go
- [ ] Replace errors in usc_coin_operations_repository.go
- [ ] Replace errors in smart_contract_operations_repository.go
- [ ] Replace errors in nft_token_operations_repository.go
- [ ] Replace errors in custom_token_operations_repository.go
- [ ] Replace errors in product_certificate_operations_repository.go
- [ ] Replace errors in validator_operations_repository.go
- [ ] Replace errors in network_operations_repository.go
- [ ] Replace errors in streaming_operations_repository.go
- [ ] Replace errors in store_bridge_operations_repository.go
- [ ] Replace errors in store_network_operations_repository.go

### **Phase 3: Verification** ⏳
- [ ] Test error handling
- [ ] Verify error codes
- [ ] Check error messages
- [ ] Test retryable error detection

---

## 📈 **METRICS**

| Metric | Before | After | Progress |
|--------|--------|-------|----------|
| **Error Constants** | 0 | 50+ | ✅ Complete |
| **Helper Functions** | 0 | 5 | ✅ Complete |
| **Standardized Errors** | 0 | 0 | ⏳ Ready |
| **Documentation** | 0 | 1 | ✅ Complete |

---

## 🚀 **NEXT STEPS**

1. **Review Error Infrastructure** - Verify error constants cover all cases
2. **Start with Priority 1** - Replace common error patterns
3. **Test Error Handling** - Verify errors work correctly
4. **Gradually Expand** - Replace all errors over time

---

**Last Updated**: 2025-01-05  
**Status**: ✅ **READY FOR IMPLEMENTATION**

