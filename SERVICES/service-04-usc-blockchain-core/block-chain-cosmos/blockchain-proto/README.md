# 🚀 **BLOCKCHAIN PROTO - COSMOS SDK PROTOCOL BUFFERS**

## 📋 **TỔNG QUAN**

**`blockchain-proto/`** là thư mục chứa Protocol Buffer definitions cho Cosmos SDK blockchain layer. Cung cấp message definitions, query definitions, và generated Go code cho tất cả 15 USC modules.

---

## 🏗️ **CẤU TRÚC THƯ MỤC**

```
blockchain-proto/
├── cosmos/                    # Cosmos SDK base protocols
│   ├── base/v1beta1/         # Base types (Coin, DecCoin)
│   ├── base/query/v1beta1/   # Query pagination
│   └── tx/v1beta1/           # Transaction types
├── usc/                      # USC module protocols (15 modules)
│   ├── usc_coin/v1/          # USC Token operations
│   ├── block/v1/             # Block operations
│   ├── store_bridge/v1/      # Cross-chain bridge
│   ├── product_certificate/v1/ # Product certificates
│   ├── smart_contract/v1/    # Smart contract operations
│   ├── monitoring/v1/        # System monitoring
│   ├── network/v1/           # Network operations
│   ├── nft_token/v1/         # NFT operations
│   ├── performance/v1/       # Performance metrics
│   ├── custom_token/v1/      # Custom token operations
│   ├── store_network/v1/     # Network storage
│   ├── streaming/v1/         # Real-time streaming
│   ├── validator/v1/         # Validator operations
│   └── transaction/v1/       # Transaction management
├── third_party/              # Third-party protocols
│   └── gogoproto/             # gogo protobuf extensions
├── simple-generate.sh        # Bash script to generate .pb.go files
└── README.md                 # This documentation
```

---

## 🎯 **MỤC ĐÍCH SỬ DỤNG**

### **1. 🔄 COSMOS SDK INTEGRATION**

**`blockchain-proto/` được sử dụng trong:**
- **`x/*/keeper/msg_server.go`** - Message server implementations
- **`x/*/keeper/query_server.go`** - Query server implementations
- **`x/*/types/`** - Type definitions and codec registration

### **2. 📊 MAPPING VỚI X/ MODULES**

| **Blockchain Proto** | **X Module** | **Usage** |
|---------------------|--------------|-----------|
| `usc/usc_coin/v1/` | `x/usc_coin/keeper/` | USC token transfers, mint, burn |
| `usc/block/v1/` | `x/block/keeper/` | Block production, validation |
| `usc/store_bridge/v1/` | `x/store_bridge/keeper/` | Cross-chain bridge operations |
| `usc/product_certificate/v1/` | `x/product_certificate/keeper/` | Product tokenization |
| `usc/smart_contract/v1/` | `x/smart_contract/keeper/` | Smart contract deployment, execution |
| `usc/monitoring/v1/` | `x/monitoring/keeper/` | System monitoring, metrics |
| `usc/network/v1/` | `x/network/keeper/` | Network topology, metrics |
| `usc/nft_token/v1/` | `x/nft_token/keeper/` | NFT minting, trading, ownership |
| `usc/performance/v1/` | `x/performance/keeper/` | Performance metrics, analytics |
| `usc/custom_token/v1/` | `x/custom_token/keeper/` | Custom token operations |
| `usc/store_network/v1/` | `x/store_network/keeper/` | Network storage operations |
| `usc/streaming/v1/` | `x/streaming/keeper/` | Real-time data streaming |
| `usc/validator/v1/` | `x/validator/keeper/` | Validator registration, staking |
| `usc/transaction/v1/` | `x/transaction/keeper/` | Transaction management, validation |

---

## 🔧 **CÁCH SỬ DỤNG**

### **1. 📝 GENERATE PROTOCOL BUFFERS**

#### **A. Sử dụng Bash Script (Recommended)**
```bash
# Generate all .pb.go files
./simple-generate.sh
```

#### **B. Manual Generation**
```bash
# Generate specific module
protoc --proto_path=. --proto_path=third_party \
  --go_out=usc/usc_coin/v1 \
  --go_opt=paths=source_relative \
  --go-grpc_out=usc/usc_coin/v1 \
  --go-grpc_opt=paths=source_relative \
  usc/usc_coin/v1/tx.proto
```

### **2. 🔄 IMPORT VÀO X/ MODULES**

#### **A. Message Server Import**
```go
// x/usc_coin/keeper/msg_server.go
import (
    "context"
    sdk "github.com/cosmos/cosmos-sdk/types"
    
    // Local types
    "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/usc_coin/types"
    
    // Blockchain proto types
    blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/usc_coin/v1"
)

// MsgServer interface with blockchain proto types
type MsgServer interface {
    TransferUSC(context.Context, *blockchainproto.MsgTransferUSC) (*blockchainproto.MsgTransferUSCResponse, error)
    MintUSC(context.Context, *blockchainproto.MsgMintUSC) (*blockchainproto.MsgMintUSCResponse, error)
    BurnUSC(context.Context, *blockchainproto.MsgBurnUSC) (*blockchainproto.MsgBurnUSCResponse, error)
}
```

#### **B. Query Server Import**
```go
// x/usc_coin/keeper/query_server.go
import (
    "context"
    
    // Blockchain proto types
    blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/usc_coin/v1"
)

// Query server with blockchain proto types
func (k Keeper) GetUSCBalance(ctx context.Context, req *blockchainproto.QueryUSCBalanceRequest) (*blockchainproto.QueryUSCBalanceResponse, error) {
    // Implementation
}
```

### **3. 🎯 VÍ DỤ CỤ THỂ - TRANSACTION MODULE**

#### **A. Protocol Buffer Definition**
```protobuf
// usc/transaction/v1/tx.proto
syntax = "proto3";
package transaction.v1;

message MsgCreateTransaction {
  string creator = 1;
  string from_address = 2;
  string to_address = 3;
  cosmos.base.v1beta1.Coin amount = 4 [(gogoproto.nullable) = false];
  string transaction_type = 5;
  string data = 6;
  string memo = 7;
}

message MsgCreateTransactionResponse {
  bool success = 1;
  string transaction_hash = 2;
  string transaction_id = 3;
}
```

#### **B. Generated Go Code**
```go
// usc/transaction/v1/tx.pb.go
type MsgCreateTransaction struct {
    Creator         string      `protobuf:"bytes,1,opt,name=creator,proto3"`
    FromAddress     string      `protobuf:"bytes,2,opt,name=from_address,json=fromAddress,proto3"`
    ToAddress       string      `protobuf:"bytes,3,opt,name=to_address,json=toAddress,proto3"`
    Amount          *types.Coin `protobuf:"bytes,4,opt,name=amount,proto3"`
    TransactionType string      `protobuf:"bytes,5,opt,name=transaction_type,json=transactionType,proto3"`
    Data            string      `protobuf:"bytes,6,opt,name=data,proto3"`
    Memo            string      `protobuf:"bytes,7,opt,name=memo,proto3"`
}
```

#### **C. Usage in Cosmos SDK Module**
```go
// x/transaction/keeper/msg_server.go
func (k msgServer) CreateTransaction(ctx context.Context, msg *blockchainproto.MsgCreateTransaction) (*blockchainproto.MsgCreateTransactionResponse, error) {
    // 1. Validate message
    if err := msg.ValidateBasic(); err != nil {
        return nil, err
    }
    
    // 2. Process transaction creation
    transactionID, transactionHash, err := k.CreateTransaction(ctx, msg)
    if err != nil {
        return nil, err
    }
    
    // 3. Return response
    return &blockchainproto.MsgCreateTransactionResponse{
        Success:        true,
        TransactionId:  transactionID,
        TransactionHash: transactionHash,
    }, nil
}
```

---

## 📋 **DANH SÁCH MODULES**

### **🔥 CORE MODULES**
- **`usc/usc_coin/v1/`** - USC token operations (transfer, mint, burn)
- **`usc/block/v1/`** - Block production and validation
- **`usc/transaction/v1/`** - Transaction management and validation

### **🔥 BUSINESS MODULES**
- **`usc/nft_token/v1/`** - NFT operations (mint, trade, ownership)
- **`usc/smart_contract/v1/`** - Smart contract deployment and execution
- **`usc/validator/v1/`** - Validator management and staking
- **`usc/network/v1/`** - Network topology and metrics
- **`usc/product_certificate/v1/`** - Product tokenization

### **🔥 ADVANCED MODULES**
- **`usc/store_bridge/v1/`** - Cross-chain bridge operations
- **`usc/streaming/v1/`** - Real-time data streaming
- **`usc/custom_token/v1/`** - Custom token operations
- **`usc/store_network/v1/`** - Network storage operations
- **`usc/monitoring/v1/`** - System monitoring and metrics
- **`usc/performance/v1/`** - Performance metrics and analytics

---

## 🔄 **WORKFLOW SỬ DỤNG**

### **1. 📝 Define Protocol Buffer**
```protobuf
// usc/transaction/v1/tx.proto
syntax = "proto3";
package transaction.v1;

message MsgCreateTransaction {
  string creator = 1;
  string from_address = 2;
  string to_address = 3;
  cosmos.base.v1beta1.Coin amount = 4 [(gogoproto.nullable) = false];
  string transaction_type = 5;
  string data = 6;
  string memo = 7;
}
```

### **2. 🔧 Generate Go Code**
```bash
# Using script (recommended)
./simple-generate.sh

# Or manual generation
protoc --proto_path=. --proto_path=third_party \
  --go_out=usc/transaction/v1 \
  --go_opt=paths=source_relative \
  --go-grpc_out=usc/transaction/v1 \
  --go-grpc_opt=paths=source_relative \
  usc/transaction/v1/tx.proto
```

### **3. 📦 Import vào X Module**
```go
import blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/transaction/v1"
```

### **4. 🎯 Implement Message Server**
```go
func (k msgServer) CreateTransaction(ctx context.Context, msg *blockchainproto.MsgCreateTransaction) (*blockchainproto.MsgCreateTransactionResponse, error) {
    // Implementation
}
```

---

## ⚠️ **LƯU Ý QUAN TRỌNG**

### **1. 🔄 KHÁC VỚI PROTO/**

| **Aspect** | **blockchain-proto/** | **proto/** |
|------------|----------------------|------------|
| **Purpose** | Cosmos SDK blockchain operations | gRPC service operations |
| **Usage** | `x/*/keeper/msg_server.go` | `business/`, `handlers/`, `repository/` |
| **Messages** | `MsgTransferUSC`, `MsgMintUSC` | `GetWalletBalanceRequest`, `TransferUSCRequest` |
| **Layer** | Blockchain layer | Application layer |

### **2. 🎯 MAPPING ĐÚNG**

```
blockchain-proto/usc/usc_coin/v1/ → x/usc_coin/keeper/msg_server.go
blockchain-proto/usc/transaction/v1/ → x/transaction/keeper/msg_server.go
blockchain-proto/usc/validator/v1/ → x/validator/keeper/msg_server.go
blockchain-proto/usc/product_certificate/v1/ → x/product_certificate/keeper/msg_server.go
blockchain-proto/usc/custom_token/v1/ → x/custom_token/keeper/msg_server.go
```

### **3. 🔧 GENERATION REQUIREMENTS**

- **protoc** - Protocol Buffer compiler
- **protoc-gen-go** - Go code generator
- **protoc-gen-go-grpc** - gRPC Go code generator
- **gogoproto** - Go protobuf extensions
- **Bash** - For script execution (cross-platform)

---

## 🎯 **TÓM TẮT**

**`blockchain-proto/` cung cấp:**
- ✅ **Protocol Buffer definitions** - Cho 15 USC modules
- ✅ **Generated Go code** - Từ `.proto` thành `.pb.go`
- ✅ **Cosmos SDK integration** - Cho blockchain operations
- ✅ **Message definitions** - `MsgCreateTransaction`, `MsgTransferUSC`, etc.
- ✅ **Query definitions** - `QueryTransactionRequest`, `QueryUSCBalanceRequest`, etc.
- ✅ **Cross-platform scripts** - `simple-generate.sh` cho Linux/macOS/Windows

**Sử dụng trong:**
- ✅ **`x/*/keeper/msg_server.go`** - Message server implementations
- ✅ **`x/*/keeper/query_server.go`** - Query server implementations
- ✅ **Cosmos SDK modules** - Blockchain layer operations
- ✅ **Service-04 USC Blockchain Core** - Core blockchain functionality

**🚀 `blockchain-proto/` là foundation cho Cosmos SDK blockchain layer với 15 modules hoàn chỉnh!**