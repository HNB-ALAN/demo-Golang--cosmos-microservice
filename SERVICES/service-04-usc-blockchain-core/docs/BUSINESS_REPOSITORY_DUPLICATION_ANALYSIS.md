# 🔍 **PHÂN TÍCH TRÙNG LẶP: BUSINESS vs REPOSITORY**

**Ngày phân tích**: 2025-11-12  
**Status**: ⚠️ **CÓ TRÙNG LẶP CHỨC NĂNG - CẦN REFACTOR**

---

## 🎯 **VẤN ĐỀ PHÁT HIỆN**

Business layer đang có **logic trùng lặp** với Repository layer trong việc tương tác trực tiếp với Cosmos SDK Keeper.

---

## 📊 **CHI TIẾT TRÙNG LẶP**

### **1. Product Certificate Operations** ⚠️ **CRITICAL**

**Business Service** (`product_certificate_operations_service.go`):
- ✅ `CreateProductCertificate`: Đã fix - chỉ delegate đến repository (line 83)
- ❌ `verifyCertificateFromKeeper` (line 259-303): **DUPLICATE** với repository
- ❌ `transferOwnershipOnBlockchain` (line 305-356): **DUPLICATE** với repository
- ❌ `createCertificateOnBlockchain` (line 224-257): **DEAD CODE** (không được gọi)

**Repository** (`product_certificate_operations_repository.go`):
- ✅ `createCertificateOnKeeper`: Single source of truth
- ✅ `verifyCertificateFromKeeper`: Single source of truth
- ✅ `transferOwnershipOnKeeper`: Single source of truth

**Vấn đề**: Business service vẫn có logic tương tác trực tiếp với keeper, tạo duplicate code.

---

### **2. Validator Operations** ⚠️ **CRITICAL**

**Business Service** (`validator_operations_service.go`):
- ❌ `registerValidatorOnBlockchain` (line 152-199): **DUPLICATE** với repository
- ❌ `getValidatorsFromKeeper` (line 209-238): **DUPLICATE** với repository
- ❌ `getValidatorStatusFromKeeper` (line 240-259): **DUPLICATE** với repository

**Repository** (`validator_operations_repository.go`):
- ✅ `registerValidatorOnKeeper`: Single source of truth
- ✅ `getValidatorsFromKeeper`: Single source of truth
- ✅ `getValidatorStatusFromKeeper`: Single source of truth

**Vấn đề**: Business service có logic duplicate với repository.

---

### **3. Custom Token Operations** ⚠️ **CRITICAL**

**Business Service** (`custom_token_operations_service.go`):
- ❌ `createTokenOnBlockchain` (line 243-285): **DUPLICATE** với repository
- ❌ `mintTokensOnBlockchain` (line 287-341): **DUPLICATE** với repository
- ❌ `getTokenBalanceFromKeeper` (line 343-395): **DUPLICATE** với repository
- ❌ `getTokenInfoFromKeeper` (line 397-430): **DUPLICATE** với repository
- ❌ `burnTokensOnBlockchain` (line 490-524): **DUPLICATE** với repository

**Repository** (`custom_token_operations_repository.go`):
- ✅ `createTokenOnKeeper`: Single source of truth
- ✅ `mintTokensOnKeeper`: Single source of truth
- ✅ `getTokenBalanceFromKeeper`: Single source of truth
- ✅ `getTokenInfoFromKeeper`: Single source of truth
- ✅ `burnTokensOnKeeper`: Single source of truth

**Vấn đề**: Business service có logic duplicate hoàn toàn với repository.

---

### **4. NFT Token Operations** ⚠️ **CRITICAL**

**Business Service** (`nft_token_operations_service.go`):
- ❌ `mintNFTOnBlockchain`: **DUPLICATE** với repository

**Repository** (`nft_token_operations_repository.go`):
- ✅ `mintNFTOnKeeper`: Single source of truth

**Vấn đề**: Business service có logic duplicate với repository.

---

### **5. Smart Contract Operations** ⚠️ **CRITICAL**

**Business Service** (`smart_contract_operations_service.go`):
- ❌ `deployContractOnBlockchain`: **DUPLICATE** với repository
- ❌ `executeContractOnBlockchain`: **DUPLICATE** với repository

**Repository** (`smart_contract_operations_repository.go`):
- ✅ `deployContractOnKeeper`: Single source of truth
- ✅ `executeContractOnKeeper`: Single source of truth

**Vấn đề**: Business service có logic duplicate với repository.

---

## 🎯 **NGUYÊN TẮC ĐÚNG**

### **Business Layer** nên:
- ✅ Validate input
- ✅ Orchestrate calls to repository
- ✅ Handle business rules (không phải data access)
- ✅ Logging và error handling ở business level
- ❌ **KHÔNG** tương tác trực tiếp với Keeper/Blockchain

### **Repository Layer** nên:
- ✅ Single source of truth cho data access
- ✅ Tương tác với Keeper (RocksDB)
- ✅ Tương tác với PostgreSQL
- ✅ Priority-based fallback (Keeper → PostgreSQL)
- ✅ Dual-write pattern (Keeper + PostgreSQL)

---

## 🛠️ **GIẢI PHÁP**

### **Option 1: Refactor Business Layer (Khuyến nghị)**

**Loại bỏ tất cả logic tương tác trực tiếp với Keeper từ Business layer**:

1. **Product Certificate Operations**:
   - ❌ Xóa `verifyCertificateFromKeeper` từ business service
   - ❌ Xóa `transferOwnershipOnBlockchain` từ business service
   - ❌ Xóa `createCertificateOnBlockchain` từ business service (dead code)
   - ✅ Business service chỉ validate và delegate đến repository

2. **Validator Operations**:
   - ❌ Xóa `registerValidatorOnBlockchain` từ business service
   - ❌ Xóa `getValidatorsFromKeeper` từ business service
   - ❌ Xóa `getValidatorStatusFromKeeper` từ business service
   - ✅ Business service chỉ validate và delegate đến repository

3. **Custom Token Operations**:
   - ❌ Xóa tất cả `*OnBlockchain` và `*FromKeeper` methods từ business service
   - ✅ Business service chỉ validate và delegate đến repository

4. **NFT Token Operations**:
   - ❌ Xóa `mintNFTOnBlockchain` từ business service
   - ✅ Business service chỉ validate và delegate đến repository

5. **Smart Contract Operations**:
   - ❌ Xóa `deployContractOnBlockchain` và `executeContractOnBlockchain` từ business service
   - ✅ Business service chỉ validate và delegate đến repository

### **Option 2: Giữ nguyên (Không khuyến nghị)**

- ⚠️ Duplicate code
- ⚠️ Confusion về responsibility
- ⚠️ Khó maintain
- ⚠️ Risk of inconsistency

---

## 📋 **CHECKLIST REFACTOR**

### **Product Certificate Operations**
- [ ] Xóa `verifyCertificateFromKeeper` từ business service
- [ ] Xóa `transferOwnershipOnBlockchain` từ business service
- [ ] Xóa `createCertificateOnBlockchain` từ business service
- [ ] Update `VerifyBlockchainProductCertificate` để chỉ delegate đến repository
- [ ] Update `TransferProductOwnership` để chỉ delegate đến repository

### **Validator Operations**
- [ ] Xóa `registerValidatorOnBlockchain` từ business service
- [ ] Xóa `getValidatorsFromKeeper` từ business service
- [ ] Xóa `getValidatorStatusFromKeeper` từ business service
- [ ] Update `RegisterValidator` để chỉ delegate đến repository
- [ ] Update `GetValidators` để chỉ delegate đến repository
- [ ] Update `GetValidatorStatus` để chỉ delegate đến repository

### **Custom Token Operations**
- [ ] Xóa `createTokenOnBlockchain` từ business service
- [ ] Xóa `mintTokensOnBlockchain` từ business service
- [ ] Xóa `getTokenBalanceFromKeeper` từ business service
- [ ] Xóa `getTokenInfoFromKeeper` từ business service
- [ ] Xóa `burnTokensOnBlockchain` từ business service
- [ ] Update tất cả methods để chỉ delegate đến repository

### **NFT Token Operations**
- [ ] Xóa `mintNFTOnBlockchain` từ business service
- [ ] Update `MintNFT` để chỉ delegate đến repository

### **Smart Contract Operations**
- [ ] Xóa `deployContractOnBlockchain` từ business service
- [ ] Xóa `executeContractOnBlockchain` từ business service
- [ ] Update `DeployContract` để chỉ delegate đến repository
- [ ] Update `ExecuteContract` để chỉ delegate đến repository

---

## ✅ **KẾT LUẬN**

**CÓ TRÙNG LẶP CHỨC NĂNG** giữa Business và Repository layer.

**Khuyến nghị**: Refactor Business layer để loại bỏ tất cả logic tương tác trực tiếp với Keeper, chỉ giữ lại:
- Input validation
- Business rules
- Orchestration (delegate to repository)
- Logging và error handling

**Repository layer** sẽ là **single source of truth** cho tất cả data access operations.

---

**Last Updated**: 2025-11-12  
**Status**: ⚠️ **ACTION REQUIRED** - Cần refactor Business layer

