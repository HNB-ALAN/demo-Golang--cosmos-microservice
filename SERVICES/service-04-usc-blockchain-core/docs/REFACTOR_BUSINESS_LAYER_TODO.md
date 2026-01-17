# 🛠️ **TODO: REFACTOR BUSINESS LAYER - 12 NHÓM**

**Ngày tạo**: 2025-11-12  
**Mục đích**: Refactor Business layer để follow pattern của Service-22 (loại bỏ duplicate code với Repository)  
**Pattern**: Business chỉ orchestrate (validation, logging, metrics, delegate), Repository là single source of truth

---

## 📋 **TỔNG QUAN**

**12 nhóm cần refactor**:
1. Block Operations
2. Transaction Operations
3. USC Coin Operations
4. Smart Contract Operations
5. NFT Token Operations
6. Custom Token Operations
7. Product Certificate Operations
8. Validator Operations
9. Network Operations
10. Streaming Operations
11. Store Bridge Operations
12. Store Network Operations

**Nguyên tắc chung**:
- ❌ **XÓA** tất cả logic tương tác trực tiếp với Keeper từ Business layer
- ✅ **GIỮ** validation, logging, metrics, delegate to repository
- ✅ Repository là **single source of truth** cho data access

---

## 📦 **GROUP 1: BLOCK OPERATIONS**

**File**: `internal/application/business/block_operations/block_operations_service.go`

### **TODO**:
- [ ] Xóa `getBlockFromKeeper` (line ~336)
- [ ] Xóa `getBlockByHashFromKeeper` (line ~352)
- [ ] Xóa `getLatestBlockFromKeeper` (line ~367)
- [ ] Xóa helper `getSDKContext` nếu chỉ dùng cho Keeper operations
- [ ] Update `GetBlock` để chỉ delegate đến `s.repo.GetBlock`
- [ ] Update `GetBlockByHash` để chỉ delegate đến `s.repo.GetBlockByHash`
- [ ] Update `GetLatestBlock` để chỉ delegate đến `s.repo.GetLatestBlock`
- [ ] Giữ lại: validation, logging, metrics recording

### **Methods cần refactor**:
- `GetBlock` → Chỉ delegate, không gọi `getBlockFromKeeper`
- `GetBlockByHash` → Chỉ delegate, không gọi `getBlockByHashFromKeeper`
- `GetLatestBlock` → Chỉ delegate, không gọi `getLatestBlockFromKeeper`

---

## 📦 **GROUP 2: TRANSACTION OPERATIONS**

**File**: `internal/application/business/transaction_operations/transaction_operations_service.go`

### **TODO**:
- [ ] Xóa `submitTransactionOnBlockchain` (line ~248)
- [ ] Xóa `getTransactionFromKeeper` (line ~308)
- [ ] Xóa `getTransactionStatusFromKeeper` (line ~325)
- [ ] Xóa `getPendingTransactionsFromKeeper` (line ~364)
- [ ] Xóa helper `getSDKContext` nếu chỉ dùng cho Keeper operations
- [ ] Update `SubmitTransaction` để chỉ delegate đến `s.repo.SubmitTransaction`
- [ ] Update `GetTransaction` để chỉ delegate đến `s.repo.GetTransaction`
- [ ] Update `GetTransactionStatus` để chỉ delegate đến `s.repo.GetTransactionStatus`
- [ ] Update `GetPendingTransactions` để chỉ delegate đến `s.repo.GetPendingTransactions`
- [ ] Giữ lại: validation, logging, metrics recording

### **Methods cần refactor**:
- `SubmitTransaction` → Chỉ delegate, không gọi `submitTransactionOnBlockchain`
- `GetTransaction` → Chỉ delegate, không gọi `getTransactionFromKeeper`
- `GetTransactionStatus` → Chỉ delegate, không gọi `getTransactionStatusFromKeeper`
- `GetPendingTransactions` → Chỉ delegate, không gọi `getPendingTransactionsFromKeeper`

---

## 📦 **GROUP 3: USC COIN OPERATIONS**

**File**: `internal/application/business/usc_coin_operations/usc_coin_operations_service.go`

### **TODO**:
- [ ] Kiểm tra xem có methods nào tương tác trực tiếp với Keeper không
- [ ] Nếu có, xóa tất cả `*OnBlockchain` và `*FromKeeper` methods
- [ ] Update tất cả methods để chỉ delegate đến repository
- [ ] Giữ lại: validation, logging, metrics recording

### **Methods cần kiểm tra**:
- `GetUSCBalance`
- `TransferUSC`
- `GetTransactionHistory`
- Tất cả methods khác

---

## 📦 **GROUP 4: SMART CONTRACT OPERATIONS**

**File**: `internal/application/business/smart_contract_operations/smart_contract_operations_service.go`

### **TODO**:
- [ ] Xóa `deployContractOnBlockchain` (line ~deployContractOnBlockchain)
- [ ] Xóa `executeContractOnBlockchain` (line ~executeContractOnBlockchain)
- [ ] Xóa `queryContractFromKeeper` (line ~queryContractFromKeeper)
- [ ] Xóa `getContractCodeFromKeeper` (line ~getContractCodeFromKeeper)
- [ ] Xóa `getContractStorageFromKeeper` (line ~getContractStorageFromKeeper)
- [ ] Xóa helper `getSDKContext` nếu chỉ dùng cho Keeper operations
- [ ] Update `DeployContract` để chỉ delegate đến `s.repo.DeployContract`
- [ ] Update `ExecuteContract` để chỉ delegate đến `s.repo.ExecuteContract`
- [ ] Update `QueryContract` để chỉ delegate đến `s.repo.QueryContract`
- [ ] Update `GetContractCode` để chỉ delegate đến `s.repo.GetContractCode`
- [ ] Update `GetContractStorage` để chỉ delegate đến `s.repo.GetContractStorage`
- [ ] Giữ lại: validation, logging, metrics recording

### **Methods cần refactor**:
- `DeployContract` → Chỉ delegate, không gọi `deployContractOnBlockchain`
- `ExecuteContract` → Chỉ delegate, không gọi `executeContractOnBlockchain`
- `QueryContract` → Chỉ delegate, không gọi `queryContractFromKeeper`
- `GetContractCode` → Chỉ delegate, không gọi `getContractCodeFromKeeper`
- `GetContractStorage` → Chỉ delegate, không gọi `getContractStorageFromKeeper`

---

## 📦 **GROUP 5: NFT TOKEN OPERATIONS**

**File**: `internal/application/business/nft_token_operations/nft_token_operations_service.go`

### **TODO**:
- [ ] Xóa `mintNFTOnBlockchain` (line ~186)
- [ ] Xóa `getNFTInfoFromKeeper` (line ~258)
- [ ] Xóa `getNFTsByOwnerFromKeeper` (line ~getNFTsByOwnerFromKeeper)
- [ ] Xóa `deployNFTContractOnBlockchain` (line ~deployNFTContractOnBlockchain)
- [ ] Xóa `createNFTCollectionOnBlockchain` (line ~createNFTCollectionOnBlockchain)
- [ ] Xóa `burnNFTOnBlockchain` (line ~burnNFTOnBlockchain)
- [ ] Xóa helper `getSDKContext` nếu chỉ dùng cho Keeper operations
- [ ] Update `MintNFT` để chỉ delegate đến `s.repo.MintNFT`
- [ ] Update `GetNFTInfo` để chỉ delegate đến `s.repo.GetNFTInfo`
- [ ] Update `GetNFTsByOwner` để chỉ delegate đến `s.repo.GetNFTsByOwner`
- [ ] Update `DeployNFTContract` để chỉ delegate đến `s.repo.DeployNFTContract`
- [ ] Update `CreateNFTCollection` để chỉ delegate đến `s.repo.CreateNFTCollection`
- [ ] Update `BurnNFT` để chỉ delegate đến `s.repo.BurnNFT`
- [ ] Giữ lại: validation, logging, metrics recording

### **Methods cần refactor**:
- `MintNFT` → Chỉ delegate, không gọi `mintNFTOnBlockchain`
- `GetNFTInfo` → Chỉ delegate, không gọi `getNFTInfoFromKeeper`
- `GetNFTsByOwner` → Chỉ delegate, không gọi `getNFTsByOwnerFromKeeper`
- `DeployNFTContract` → Chỉ delegate, không gọi `deployNFTContractOnBlockchain`
- `CreateNFTCollection` → Chỉ delegate, không gọi `createNFTCollectionOnBlockchain`
- `BurnNFT` → Chỉ delegate, không gọi `burnNFTOnBlockchain`

---

## 📦 **GROUP 6: CUSTOM TOKEN OPERATIONS**

**File**: `internal/application/business/custom_token_operations/custom_token_operations_service.go`

### **TODO**:
- [ ] Xóa `createTokenOnBlockchain` (line ~244)
- [ ] Xóa `mintTokensOnBlockchain` (line ~288)
- [ ] Xóa `getTokenBalanceFromKeeper` (line ~344)
- [ ] Xóa `getTokenInfoFromKeeper` (line ~398)
- [ ] Xóa `burnTokensOnBlockchain` (line ~491)
- [ ] Xóa helper `getSDKContext` nếu chỉ dùng cho Keeper operations
- [ ] Update `CreateBlockchainToken` để chỉ delegate đến `s.repo.CreateBlockchainToken`
- [ ] Update `MintTokens` để chỉ delegate đến `s.repo.MintTokens`
- [ ] Update `GetTokenBalance` để chỉ delegate đến `s.repo.GetTokenBalance`
- [ ] Update `GetTokenInfo` để chỉ delegate đến `s.repo.GetTokenInfo`
- [ ] Update `BurnTokens` để chỉ delegate đến `s.repo.BurnTokens`
- [ ] Giữ lại: validation, logging, metrics recording

### **Methods cần refactor**:
- `CreateBlockchainToken` → Chỉ delegate, không gọi `createTokenOnBlockchain`
- `MintTokens` → Chỉ delegate, không gọi `mintTokensOnBlockchain`
- `GetTokenBalance` → Chỉ delegate, không gọi `getTokenBalanceFromKeeper`
- `GetTokenInfo` → Chỉ delegate, không gọi `getTokenInfoFromKeeper`
- `BurnTokens` → Chỉ delegate, không gọi `burnTokensOnBlockchain`

---

## 📦 **GROUP 7: PRODUCT CERTIFICATE OPERATIONS**

**File**: `internal/application/business/product_certificate_operations/product_certificate_operations_service.go`

### **TODO**:
- [ ] Xóa `createCertificateOnBlockchain` (line ~225) - **DEAD CODE** (đã fix trước đó)
- [ ] Xóa `verifyCertificateFromKeeper` (line ~260)
- [ ] Xóa `transferOwnershipOnBlockchain` (line ~306)
- [ ] Xóa helper `getSDKContext` nếu chỉ dùng cho Keeper operations
- [ ] Update `CreateProductCertificate` - ✅ **ĐÃ FIX** (chỉ delegate đến repository)
- [ ] Update `VerifyBlockchainProductCertificate` để chỉ delegate đến `s.repo.VerifyBlockchainProductCertificate`
- [ ] Update `TransferProductOwnership` để chỉ delegate đến `s.repo.TransferProductOwnership`
- [ ] Giữ lại: validation, logging, metrics recording

### **Methods cần refactor**:
- `CreateProductCertificate` → ✅ **ĐÃ FIX** (chỉ delegate)
- `VerifyBlockchainProductCertificate` → Chỉ delegate, không gọi `verifyCertificateFromKeeper`
- `TransferProductOwnership` → Chỉ delegate, không gọi `transferOwnershipOnBlockchain`

---

## 📦 **GROUP 8: VALIDATOR OPERATIONS**

**File**: `internal/application/business/validator_operations/validator_operations_service.go`

### **TODO**:
- [ ] Xóa `registerValidatorOnBlockchain` (line ~152)
- [ ] Xóa `getValidatorsFromKeeper` (line ~209)
- [ ] Xóa `getValidatorStatusFromKeeper` (line ~240)
- [ ] Xóa helper `convertValidatorToProto` nếu chỉ dùng cho Keeper operations
- [ ] Xóa helper `getSDKContext` nếu chỉ dùng cho Keeper operations
- [ ] Update `RegisterValidator` để chỉ delegate đến `s.repo.RegisterValidator`
- [ ] Update `GetValidators` để chỉ delegate đến `s.repo.GetValidators`
- [ ] Update `GetValidatorStatus` để chỉ delegate đến `s.repo.GetValidatorStatus`
- [ ] `StakeUSC` và `UnstakeUSC` - ✅ **ĐÃ ĐÚNG** (chỉ delegate)
- [ ] Giữ lại: validation, logging, metrics recording

### **Methods cần refactor**:
- `RegisterValidator` → Chỉ delegate, không gọi `registerValidatorOnBlockchain`
- `GetValidators` → Chỉ delegate, không gọi `getValidatorsFromKeeper`
- `GetValidatorStatus` → Chỉ delegate, không gọi `getValidatorStatusFromKeeper`
- `StakeUSC` → ✅ **ĐÃ ĐÚNG** (chỉ delegate)
- `UnstakeUSC` → ✅ **ĐÃ ĐÚNG** (chỉ delegate)

---

## 📦 **GROUP 9: NETWORK OPERATIONS**

**File**: `internal/application/business/network_operations/network_operations_service.go`

### **TODO**:
- [ ] Xóa `getNetworkInfoFromKeeper` (line ~124)
- [ ] Xóa `getPeersFromKeeper` (line ~158)
- [ ] Xóa helper `getSDKContext` nếu chỉ dùng cho Keeper operations
- [ ] Update `GetNetworkInfo` để chỉ delegate đến `s.repo.GetNetworkInfo`
- [ ] Update `GetPeers` để chỉ delegate đến `s.repo.GetPeers`
- [ ] Giữ lại: validation, logging, metrics recording

### **Methods cần refactor**:
- `GetNetworkInfo` → Chỉ delegate, không gọi `getNetworkInfoFromKeeper`
- `GetPeers` → Chỉ delegate, không gọi `getPeersFromKeeper`

---

## 📦 **GROUP 10: STREAMING OPERATIONS**

**File**: `internal/application/business/streaming_operations/streaming_operations_service.go`

### **TODO**:
- [ ] Kiểm tra xem có methods nào tương tác trực tiếp với Keeper không
- [ ] Nếu có, xóa tất cả `*OnBlockchain` và `*FromKeeper` methods
- [ ] Update tất cả methods để chỉ delegate đến repository
- [ ] Giữ lại: validation, logging, metrics recording

### **Methods cần kiểm tra**:
- Tất cả methods trong file

---

## 📦 **GROUP 11: STORE BRIDGE OPERATIONS**

**File**: `internal/application/business/store_bridge_operations/store_bridge_operations_service.go`

### **TODO**:
- [ ] Kiểm tra xem có methods nào tương tác trực tiếp với Keeper không
- [ ] Nếu có, xóa tất cả `*OnBlockchain` và `*FromKeeper` methods
- [ ] Update tất cả methods để chỉ delegate đến repository
- [ ] Giữ lại: validation, logging, metrics recording

### **Methods cần kiểm tra**:
- `DeployStoreBridge`
- `BridgeStoreTokenToUSC`
- `BridgeUSCToStoreToken`
- Tất cả methods khác

---

## 📦 **GROUP 12: STORE NETWORK OPERATIONS**

**File**: `internal/application/business/store_network_operations/store_network_operations_service.go`

### **TODO**:
- [ ] Kiểm tra xem có methods nào tương tác trực tiếp với Keeper không
- [ ] Nếu có, xóa tất cả `*OnBlockchain` và `*FromKeeper` methods
- [ ] Update tất cả methods để chỉ delegate đến repository
- [ ] Giữ lại: validation, logging, metrics recording

### **Methods cần kiểm tra**:
- `RegisterStoreNetwork`
- `SyncStoreNetworkState`
- Tất cả methods khác

---

## ✅ **CHECKLIST SAU KHI REFACTOR**

### **Kiểm tra từng file Business Service**:
- [ ] Không còn `*OnBlockchain` methods
- [ ] Không còn `*FromKeeper` methods
- [ ] Không còn `*OnKeeper` methods
- [ ] Không còn `getSDKContext` helper nếu chỉ dùng cho Keeper operations
- [ ] Tất cả methods chỉ có: validation, logging, metrics, delegate to repository
- [ ] Repository là single source of truth cho data access

### **Kiểm tra tổng thể**:
- [ ] Chạy linter: `golangci-lint run`
- [ ] Chạy tests: `go test ./...`
- [ ] Verify không có duplicate code giữa Business và Repository
- [ ] Verify Business layer không tương tác trực tiếp với `cosmosApp.*Keeper`

---

## 📊 **THỐNG KÊ**

**Tổng số methods cần refactor**: ~30+ methods  
**Tổng số helper functions cần xóa**: ~15+ functions  
**Tổng số lines code cần xóa**: ~500+ lines (ước tính)

---

## 🎯 **MỤC TIÊU**

Sau khi refactor xong:
- ✅ Business layer chỉ orchestrate (validation, logging, metrics, delegate)
- ✅ Repository là single source of truth cho data access
- ✅ Không còn duplicate code
- ✅ Clear separation of concerns
- ✅ Follow pattern của Service-22

---

**Last Updated**: 2025-11-12  
**Status**: 📋 **TODO CREATED** - Sẵn sàng bắt đầu refactor

