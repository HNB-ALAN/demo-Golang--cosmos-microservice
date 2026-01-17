# 🚀 **X/ MODULES - COSMOS SDK USC MODULES**

## 📋 **TỔNG QUAN**

**`x/`** là thư mục chứa 12 Cosmos SDK modules cho USC blockchain platform. Mỗi module implement Cosmos SDK pattern với keeper, message server, query server, và type definitions.

---

## 🏗️ **CẤU TRÚC THƯ MỤC**

```
x/
├── usc/                      # USC Token Module (Core)
│   ├── keeper/
│   │   ├── keeper.go        # State management
│   │   ├── msg_server.go    # Message server
│   │   └── query_server.go  # Query server
│   ├── types/
│   │   ├── types.go         # Data types
│   │   ├── keys.go          # Store keys
│   │   ├── codec.go         # Codec registration
│   │   └── genesis.go       # Genesis state
│   ├── module.go            # Module definition
│   └── abci.go              # ABCI handlers
├── nft/                     # NFT Module
├── contract/                # Smart Contract Module
├── validator/               # Validator Module
├── network/                 # Network Module
├── bridge/                  # Cross-chain Bridge Module
├── streaming/               # Streaming Module
├── certificate/             # Certificate Module
├── token/                   # Custom Token Module
├── store/                   # Store Module
├── block/                   # Block Module
├── reward/                  # Reward Module
└── modules_test.go          # Module tests
```

---

## 🎯 **MỤC ĐÍCH SỬ DỤNG**

### **1. 🔄 COSMOS SDK INTEGRATION**

**`x/` modules được sử dụng trong:**
- **`block-chain-cosmos/app/app.go`** - Module registration
- **`block-chain-cosmos/blockchain-proto/`** - Protocol Buffer integration
- **Cosmos SDK blockchain operations** - Message processing, state management

### **2. 📊 MAPPING VỚI BLOCKCHAIN-PROTO**

| **X Module** | **Blockchain Proto** | **Purpose** |
|--------------|---------------------|-------------|
| `x/usc/` | `blockchain-proto/usc/coin/v1/` | USC token transfers, mint, burn |
| `x/nft/` | `blockchain-proto/usc/nft/v1/` | NFT minting, trading, ownership |
| `x/contract/` | `blockchain-proto/usc/contract/v1/` | Smart contract deployment, execution |
| `x/validator/` | `blockchain-proto/usc/validator/v1/` | Validator registration, staking |
| `x/network/` | `blockchain-proto/usc/network/v1/` | Network topology, metrics |
| `x/bridge/` | `blockchain-proto/usc/bridge/v1/` | Cross-chain bridge operations |
| `x/streaming/` | `blockchain-proto/usc/streaming/v1/` | Real-time data streaming |
| `x/certificate/` | `blockchain-proto/usc/certificate/v1/` | Product tokenization |
| `x/token/` | `blockchain-proto/usc/token/v1/` | Custom token creation |
| `x/store/` | `blockchain-proto/usc/store/v1/` | Data storage operations |
| `x/block/` | `blockchain-proto/usc/block/v1/` | Block production, validation |
| `x/reward/` | `blockchain-proto/usc/reward/v1/` | USC reward distribution |

---

## 🔧 **CÁCH SỬ DỤNG**

### **1. 📝 IMPORT BLOCKCHAIN-PROTO**

#### **A. Message Server Import**
```go
// x/usc/keeper/msg_server.go
import (
    "context"
    sdk "github.com/cosmos/cosmos-sdk/types"
    
    // Local types
    "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/x/usc/types"
    
    // Blockchain proto types
    blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/coin/v1"
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
// x/usc/keeper/query_server.go
import (
    "context"
    
    // Blockchain proto types
    blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/coin/v1"
)

// Query server with blockchain proto types
func (k Keeper) GetUSCBalance(ctx context.Context, req *blockchainproto.QueryUSCBalanceRequest) (*blockchainproto.QueryUSCBalanceResponse, error) {
    // Implementation
}
```

### **2. 🎯 VÍ DỤ CỤ THỂ - USC MODULE**

#### **A. Message Server Implementation**
```go
// x/usc/keeper/msg_server.go
func (k msgServer) TransferUSC(ctx context.Context, msg *blockchainproto.MsgTransferUSC) (*blockchainproto.MsgTransferUSCResponse, error) {
    // 1. Validate message
    if err := msg.ValidateBasic(); err != nil {
        return nil, err
    }
    
    // 2. Process blockchain transfer
    err := k.TransferUSC(ctx, msg.FromAddress, msg.ToAddress, msg.Amount)
    if err != nil {
        return nil, err
    }
    
    // 3. Return response
    return &blockchainproto.MsgTransferUSCResponse{
        Success: true,
        TransactionHash: "0x123...",
    }, nil
}
```

#### **B. Query Server Implementation**
```go
// x/usc/keeper/query_server.go
func (k Keeper) GetUSCBalance(ctx context.Context, req *blockchainproto.QueryUSCBalanceRequest) (*blockchainproto.QueryUSCBalanceResponse, error) {
    // Query balance from blockchain state
    balance := k.GetBalance(ctx, req.Address)
    
    return &blockchainproto.QueryUSCBalanceResponse{
        Balance: balance,
    }, nil
}
```

#### **C. Keeper Implementation**
```go
// x/usc/keeper/keeper.go
type Keeper struct {
    cdc        codec.BinaryCodec
    storeKey   storetypes.StoreKey
    paramSpace paramtypes.Subspace
    bk         keeper.Keeper
}

func (k Keeper) TransferUSC(ctx context.Context, from, to string, amount *types.Coin) error {
    // Blockchain transfer logic
    return k.bk.SendCoins(ctx, from, to, sdk.NewCoins(*amount))
}
```

### **3. 🔄 MODULE REGISTRATION**

#### **A. App Registration**
```go
// block-chain-cosmos/app/app.go
func NewUSCApp() *USCApp {
    app := &USCApp{
        BaseApp: baseapp.NewBaseApp(...),
    }
    
    // Register USC modules
    app.ModuleManager = module.NewManager(
        usc.NewAppModule(appCodec, app.USCKeeper),
        nft.NewAppModule(appCodec, app.NFTKeeper),
        contract.NewAppModule(appCodec, app.ContractKeeper),
        // ... other modules
    )
    
    return app
}
```

#### **B. Module Definition**
```go
// x/usc/module.go
type AppModule struct {
    AppModuleBasic
    keeper keeper.Keeper
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
    types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
    types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServerImpl(am.keeper))
}
```

---

## 📋 **DANH SÁCH MODULES**

### **🔥 CORE MODULES**
- **`x/usc/`** - USC token operations (transfer, mint, burn)
- **`x/reward/`** - USC reward distribution
- **`x/block/`** - Block production and validation

### **🔥 BUSINESS MODULES**
- **`x/nft/`** - NFT operations (mint, trade, ownership)
- **`x/contract/`** - Smart contract deployment and execution
- **`x/validator/`** - Validator management and staking
- **`x/network/`** - Network topology and metrics

### **🔥 ADVANCED MODULES**
- **`x/bridge/`** - Cross-chain bridge operations
- **`x/streaming/`** - Real-time data streaming
- **`x/certificate/`** - Product tokenization
- **`x/token/`** - Custom token creation
- **`x/store/`** - Data storage operations

---

## 🔄 **WORKFLOW SỬ DỤNG**

### **1. 📝 Define Protocol Buffer**
```protobuf
// blockchain-proto/usc/coin/v1/tx.proto
message MsgTransferUSC {
  string from_address = 1;
  string to_address = 2;
  cosmos.base.v1beta1.Coin amount = 3;
  string memo = 4;
}
```

### **2. 🔧 Generate Go Code**
```bash
protoc --proto_path=blockchain-proto --go_out=blockchain-proto/usc/coin/v1 usc/coin/v1/tx.proto
```

### **3. 📦 Import vào X Module**
```go
import blockchainproto "github.com/usc-platform/usc-social-xxx-app/SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/blockchain-proto/usc/coin/v1"
```

### **4. 🎯 Implement Message Server**
```go
func (k msgServer) TransferUSC(ctx context.Context, msg *blockchainproto.MsgTransferUSC) (*blockchainproto.MsgTransferUSCResponse, error) {
    // Implementation
}
```

### **5. 🔄 Register Module**
```go
// app/app.go
app.ModuleManager = module.NewManager(
    usc.NewAppModule(appCodec, app.USCKeeper),
    // ... other modules
)
```

---

## ⚠️ **LƯU Ý QUAN TRỌNG**

### **1. 🔄 KHÁC VỚI BUSINESS LAYER**

| **Aspect** | **x/ modules** | **business/ services** |
|------------|----------------|------------------------|
| **Purpose** | Cosmos SDK blockchain operations | gRPC service operations |
| **Usage** | `blockchain-proto/` integration | `proto/` integration |
| **Messages** | `MsgTransferUSC`, `MsgMintUSC` | `GetWalletBalanceRequest`, `TransferUSCRequest` |
| **Layer** | Blockchain layer | Application layer |

### **2. 🎯 MAPPING ĐÚNG**

```
blockchain-proto/usc/coin/v1/ → x/usc/keeper/msg_server.go
blockchain-proto/usc/nft/v1/ → x/nft/keeper/msg_server.go
blockchain-proto/usc/contract/v1/ → x/contract/keeper/msg_server.go
```

### **3. 🔧 COSMOS SDK PATTERN**

Mỗi module phải có:
- ✅ **`keeper/`** - State management
- ✅ **`types/`** - Type definitions
- ✅ **`module.go`** - Module definition
- ✅ **`abci.go`** - ABCI handlers
- ✅ **Protocol Buffer integration** - Import từ `blockchain-proto/`

---

## 🎯 **TÓM TẮT**

**`x/` modules cung cấp:**
- ✅ **Cosmos SDK modules** - 12 USC modules
- ✅ **Blockchain operations** - Message processing, state management
- ✅ **Protocol Buffer integration** - Import từ `blockchain-proto/`
- ✅ **Module registration** - Trong `app/app.go`
- ✅ **Keeper pattern** - State management, message server, query server

**Sử dụng trong:**
- ✅ **`block-chain-cosmos/app/app.go`** - Module registration
- ✅ **`blockchain-proto/`** - Protocol Buffer integration
- ✅ **Cosmos SDK blockchain** - Blockchain layer operations

**🚀 `x/` modules là core của Cosmos SDK blockchain layer!**
