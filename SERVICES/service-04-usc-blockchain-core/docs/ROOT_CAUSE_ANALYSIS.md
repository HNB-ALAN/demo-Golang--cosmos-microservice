# 🔍 **ROOT CAUSE ANALYSIS - Product Certificate Operations**

**Ngày phân tích**: 2025-11-12  
**Status**: ✅ **VẤN ĐỀ GỐC RỄ ĐÃ XÁC ĐỊNH**

---

## 🎯 **3 VẤN ĐỀ GỐC RỄ**

### **1. GetSDKContext sử dụng READ-ONLY CONTEXT** ⚠️ **CRITICAL**

**Vấn đề**:
- `GetSDKContext` fallback sử dụng `NewContext(true)` - **READ-ONLY CONTEXT**
- Khi gọi `SetCertificate` trên read-only context, nó **KHÔNG COMMIT** vào RocksDB
- Khi gọi `GetCertificate` sau đó, nó đọc từ **committed state** (chưa có certificate mới)
- → Certificate được "tạo" nhưng không tồn tại trong keeper

**Code hiện tại**:
```go
// cosmos_context.go:46
sdkCtx = cosmosApp.BaseApp.NewContext(true) // READ-ONLY!
```

**Giải pháp**:
- Tạo helper function `GetSDKContextForWrite` sử dụng `NewContext(false)` cho write operations
- Hoặc sử dụng transaction context nếu có

---

### **2. ProductCertificate không có Protobuf Tags** ⚠️ **CRITICAL**

**Vấn đề**:
- `ProductCertificate` struct chỉ có `json` tags, **KHÔNG CÓ protobuf tags**
- Khi unmarshal từ RocksDB, protobuf codec cần protobuf tags
- → Panic: `"protobuf tag not enough fields in ProductCertificate.ID: "`

**Code hiện tại**:
```go
// types/types.go:37-47
type ProductCertificate struct {
    ID         string `json:"id"`           // ❌ Không có protobuf tag
    ProductID  string `json:"product_id"`   // ❌ Không có protobuf tag
    Owner      string `json:"owner"`        // ❌ Không có protobuf tag
    // ...
}
```

**Giải pháp**:
- Thêm protobuf tags vào ProductCertificate struct
- Hoặc sử dụng protobuf definition từ `blockchain-proto/usc/product_certificate/v1/tx.proto`

---

### **3. Duplicate Certificate Creation** ⚠️ **LOGIC ERROR**

**Vấn đề**:
- Business service (`createCertificateOnBlockchain`) tạo certificate với ID mới (timestamp)
- Repository (`createCertificateOnKeeper`) cũng tạo certificate với ID mới (timestamp khác)
- → **2 certificate ID khác nhau!**
- Business service tạo certificate nhưng không dùng, chỉ gọi repository để tạo lại

**Code hiện tại**:
```go
// business service:83
certificateId, err := s.createCertificateOnBlockchain(ctx, req) // ID: cert_PRD-001_1762958699

// business service:94
repoResult, repoErr := s.repo.CreateProductCertificate(ctx, req) // ID: cert_PRD-001_1762958700 (khác!)
```

**Giải pháp**:
- **Loại bỏ duplicate creation**: Business service không nên tạo certificate, chỉ gọi repository
- Repository là single source of truth cho certificate creation

---

## 🔧 **GIẢI PHÁP TẬN GỐC**

### **Fix 1: Tạo Writable Context Helper**

```go
// cosmos_context.go
// GetSDKContextForWrite creates a writable sdk.Context for write operations
func GetSDKContextForWrite(ctx context.Context, cosmosApp *app.USCApp, logger *logging.Logger) (sdk.Context, error) {
    if cosmosApp == nil || cosmosApp.BaseApp == nil {
        return sdk.Context{}, errors.New("cosmosApp not initialized")
    }

    // Try to unwrap SDK context from gRPC context first
    var sdkCtx sdk.Context
    func() {
        defer func() {
            if panicVal := recover(); panicVal != nil {
                // Context is not wrapped, will use fallback
            }
        }()
        sdkCtx = sdk.UnwrapSDKContext(ctx)
    }()

    // If context is not wrapped, use NewContext(false) for write operations
    if sdkCtx.IsZero() {
        if logger != nil {
            logger.Debug("Context not wrapped, using BaseApp.NewContext(false) for write operation")
        }
        // Use NewContext(false) to allow writes (will be committed on next block)
        sdkCtx = cosmosApp.BaseApp.NewContext(false)
    }

    return sdkCtx, nil
}
```

### **Fix 2: Thêm Protobuf Tags vào ProductCertificate**

```go
// types/types.go
type ProductCertificate struct {
    ID         string `json:"id" protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
    ProductID  string `json:"product_id" protobuf:"bytes,2,opt,name=product_id,json=productId,proto3" json:"product_id"`
    Owner      string `json:"owner" protobuf:"bytes,3,opt,name=owner,proto3" json:"owner"`
    Status     string `json:"status" protobuf:"bytes,4,opt,name=status,proto3" json:"status"`
    Metadata   string `json:"metadata" protobuf:"bytes,5,opt,name=metadata,proto3" json:"metadata"`
    CreatedAt  int64  `json:"created_at" protobuf:"varint,6,opt,name=created_at,json=createdAt,proto3" json:"created_at"`
    UpdatedAt  int64  `json:"updated_at" protobuf:"varint,7,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at"`
    VerifiedAt int64  `json:"verified_at" protobuf:"varint,8,opt,name=verified_at,json=verifiedAt,proto3" json:"verified_at"`
    ExpiresAt  int64  `json:"expires_at" protobuf:"varint,9,opt,name=expires_at,json=expiresAt,proto3" json:"expires_at"`
}
```

### **Fix 3: Loại bỏ Duplicate Creation**

```go
// business service: CreateProductCertificate
// ❌ REMOVE: createCertificateOnBlockchain call
// ✅ ONLY: Call repository directly
return s.repo.CreateProductCertificate(ctx, req)
```

---

## 📊 **IMPACT ANALYSIS**

### **Current Behavior**:
1. Business service tạo certificate trên keeper (read-only context) → **KHÔNG COMMIT**
2. Repository tạo certificate mới với ID khác → **COMMIT** (nhưng ID khác)
3. Verify/Transfer tìm certificate với ID từ business service → **KHÔNG TÌM THẤY**
4. Panic khi unmarshal → **FALLBACK** về database

### **After Fix**:
1. Repository tạo certificate trên keeper (writable context) → **COMMIT**
2. Certificate có trong keeper với đúng ID
3. Verify/Transfer tìm thấy certificate → **SUCCESS**
4. Không panic → **NO FALLBACK NEEDED**

---

## 🎯 **PRIORITY**

1. **Fix 1 (Writable Context)**: ⚠️ **CRITICAL** - Certificate không được commit
2. **Fix 2 (Protobuf Tags)**: ⚠️ **CRITICAL** - Panic khi unmarshal
3. **Fix 3 (Duplicate Creation)**: ⚠️ **HIGH** - Logic error, waste resources

---

**Last Updated**: 2025-11-12  
**Status**: ✅ **ROOT CAUSES IDENTIFIED**

