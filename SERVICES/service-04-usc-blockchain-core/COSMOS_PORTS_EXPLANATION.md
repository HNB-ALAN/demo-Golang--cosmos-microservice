# 🔌 COSMOS SDK PORTS - GIẢI THÍCH

## 📊 Ports 9090 và 1317 - Nguồn Gốc

### **✅ Đây là COSMOS SDK STANDARD PORTS**

Hai ports này **KHÔNG PHẢI** custom ports mà là **standard ports** trong Cosmos SDK ecosystem:

---

## 🔹 **Port 9090 - gRPC Query**

### **Nguồn gốc**
- ✅ **Cosmos SDK Standard**: Port 9090 là default cho gRPC Query server
- ✅ **Used by**: All Cosmos SDK chains (Osmosis, Cosmos Hub, etc.)
- ✅ **Purpose**: gRPC queries cho Cosmos SDK modules

### **Config trong code**
1. **`config/app.toml`** (line 55):
   ```toml
   [app.network]
   grpc_port = 9090
   ```

2. **`config/config.toml`** (line 12):
   ```toml
   grpc_address = "0.0.0.0:9090"
   ```

3. **`daemon/uscd/init.go`** (line 507):
   ```go
   address = "0.0.0.0:9090"  // Trong [grpc] section
   ```

4. **`daemon/uscd/start.go`** (line 51):
   ```go
   startCmd.Flags().String("grpc-address", "0.0.0.0:9090", "gRPC listen address")
   ```

### **Vị trí trong Cosmos SDK**
- **Cosmos SDK BaseApp**: Exposes gRPC server trên port 9090
- **Module Queries**: Tất cả Cosmos SDK modules (bank, auth, staking, etc.) register gRPC query services
- **Query Router**: `app.BaseApp.GRPCQueryRouter()` handles queries trên port này

---

## 🔹 **Port 1317 - REST API (LCD)**

### **Nguồn gốc**
- ✅ **Cosmos SDK Standard**: Port 1317 là default cho REST API (Light Client Daemon - LCD)
- ✅ **Used by**: All Cosmos SDK chains
- ✅ **Purpose**: REST API queries cho web clients, wallets, explorers

### **Config trong code**
1. **`config/app.toml`** (line 56):
   ```toml
   [app.network]
   rest_port = 1317
   ```

2. **`config/config.toml`** (line 11):
   ```toml
   api_address = "0.0.0.0:1317"
   ```

3. **`daemon/uscd/init.go`** (line 500):
   ```go
   address = "0.0.0.0:1317"  // Trong [api] section
   ```

4. **`daemon/uscd/start.go`** (line 52):
   ```go
   startCmd.Flags().String("api-address", "0.0.0.0:1317", "REST API listen address")
   ```

### **Vị trí trong Cosmos SDK**
- **Cosmos SDK BaseApp**: Exposes REST server trên port 1317
- **gRPC Gateway**: REST API là gRPC Gateway proxy (gRPC → REST translation)
- **Client Tools**: Cosmos CLI, wallets, explorers sử dụng port này

---

## 📋 **Tổng Quan Ports**

### **Cosmos SDK Ecosystem Standard Ports**

| Port | Service | Protocol | Standard |
|------|---------|----------|----------|
| **9090** | gRPC Query | gRPC/HTTP2 | ✅ Cosmos SDK |
| **1317** | REST API (LCD) | HTTP/JSON | ✅ Cosmos SDK |
| **26657** | CometBFT RPC | JSON-RPC | ✅ CometBFT |
| **26656** | CometBFT P2P | P2P | ✅ CometBFT |
| **26658** | ABCI (standalone) | Socket | ✅ CometBFT |

---

## 🎯 **Tại Sao Dùng Standard Ports?**

### **1. Ecosystem Compatibility**
- ✅ **Wallets**: Ledger, Keplr, Cosmostation hỗ trợ standard ports
- ✅ **Explorers**: Mintscan, BigDipper kết nối tới standard ports
- ✅ **Tools**: `gaiad`, `osmosisd`, `junod` đều dùng ports này

### **2. Developer Experience**
- ✅ **Documentation**: Cosmos SDK docs assume standard ports
- ✅ **Tutorials**: Tất cả tutorials sử dụng 9090 và 1317
- ✅ **SDKs**: Client SDKs (cosmos-py, cosmos-js) expect standard ports

### **3. Network Standard**
- ✅ **Interoperability**: Cross-chain tools work với standard ports
- ✅ **Monitoring**: Cosmos ecosystem monitoring tools expect these ports
- ✅ **Testing**: All Cosmos SDK test suites use standard ports

---

## 📂 **Files Configuring These Ports**

### **1. Configuration Files**
- ✅ `block-chain-cosmos/config/app.toml` - App-level config
- ✅ `block-chain-cosmos/config/config.toml` - Node-level config
- ✅ `block-chain-cosmos/config/client.toml` - Client config
- ✅ `block-chain-cosmos/test-home/config/app.toml` - Test config

### **2. Code Files**
- ✅ `daemon/uscd/init.go` - Creates config với standard ports
- ✅ `daemon/uscd/start.go` - Command flags với default ports
- ✅ `app/app.go` - BaseApp exposes servers trên standard ports

---

## ✅ **Kết Luận**

**Ports 9090 và 1317 là COSMOS SDK STANDARD PORTS**, không phải custom ports.

**Sử dụng standard ports giúp**:
- ✅ Compatibility với Cosmos ecosystem tools
- ✅ Developer experience tốt hơn
- ✅ Dễ dàng tích hợp với wallets, explorers
- ✅ Tuân thủ Cosmos SDK best practices

**Không nên thay đổi** các ports này trừ khi có lý do đặc biệt (conflict, security, etc.).

