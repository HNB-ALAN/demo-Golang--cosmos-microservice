# Service-04 Architecture Patterns and Conventions

## 📋 Table of Contents

1. [Overview](#overview)
2. [Architecture Layers](#architecture-layers)
3. [Dependency Injection](#dependency-injection)
4. [Correlation IDs](#correlation-ids)
5. [Error Handling](#error-handling)
6. [Logging Conventions](#logging-conventions)
7. [Code Organization](#code-organization)
8. [Best Practices](#best-practices)

---

## Overview

Service-04 (USC Blockchain Core) là **blockchain service core** trong USC ecosystem, cung cấp 12 gRPC services cho blockchain operations. Service này follow **layered architecture** pattern với clear separation of concerns:

### Architecture Components

1. **Service Layer (gRPC API)**: 
   - Gateway Service → USC Blockchain Core → 12 gRPC Services (Handlers)
   
2. **Application Layer**:
   - Business Services (business logic, validation, logging)
   - Repository Layer (data access, Cosmos SDK operations)
   
3. **Blockchain Layer (Cosmos SDK)**:
   - USC App với 14 Custom Modules
   - RocksDB Storage (blockchain state)
   - PostgreSQL Analytics (analytics data)
   
4. **Infrastructure**:
   - PostgreSQL (database)
   - Redis (cache)
   - Kafka (messaging)
   - Monitoring (observability)

```
┌─────────────────────────────────────────────────────────────┐
│              Service Layer (gRPC API)                        │
│  ┌──────────────┐      ┌──────────────────┐                │
│  │Gateway Service│ ───▶ │USC Blockchain   │                │
│  │              │      │Core (Service-04) │                │
│  └──────────────┘      └────────┬─────────┘                │
│                                  │                           │
│                                  ▼                           │
│                          ┌──────────────────┐               │
│                          │ 12 gRPC Services │               │
│                          │ (Handlers Layer) │               │
│                          └──────────────────┘               │
└────────────────────────────────┬────────────────────────────┘
                                 │
                    ┌─────────────┴─────────────┐
                    │                           │
        ┌───────────▼──────────┐   ┌───────────▼──────────┐
        │  Application Layer   │   │ Blockchain Layer     │
        │  (Business Services)  │   │ (Cosmos SDK)         │
        │                       │   │                      │
        │  - Business logic    │   │  - USC App          │
        │  - Validation        │   │  - 14 Custom Modules│
        │  - Logging/Metrics   │   │  - RocksDB Storage  │
        │  - Correlation IDs   │   │  - PostgreSQL       │
        └───────────┬───────────┘   │    Analytics        │
                    │               └───────────┬──────────┘
                    │                           │
        ┌───────────▼──────────┐               │
        │  Repository Layer    │◀──────────────┘
        │                      │
        │  - Data access       │
        │  - Cosmos SDK ops    │
        │  - Database ops      │
        │  - Error handling    │
        └───────────┬───────────┘
                    │
        ┌───────────▼──────────┐
        │  Infrastructure      │
        │                      │
        │  - PostgreSQL DB    │
        │  - Redis Cache       │
        │  - RocksDB (State)   │
        │  - Kafka Messaging   │
        │  - Monitoring        │
        └──────────────────────┘
```

---

## Architecture Layers

### 1. Handlers Layer (`internal/application/handlers/`)

**Purpose**: Pure delegation layer - chỉ chuyển tiếp gRPC requests xuống Business Service.

**Responsibilities**:
- ✅ Receive gRPC requests
- ✅ Delegate to business service (pure delegation)
- ❌ **NO business logic**
- ❌ **NO logging** (delegated to business layer)
- ❌ **NO metrics** (delegated to business layer)
- ❌ **NO transformation** (pure pass-through)

**Pattern**:
```go
func (h *Handlers) MethodName(ctx context.Context, req *proto.Request) (*proto.Response, error) {
    // Pure delegation - no business logic
    return h.service.MethodName(ctx, req)
}
```

**Example**:
```go
// internal/application/handlers/block_operations/block_operations_handlers.go
func (h *Handlers) ProduceBlock(ctx context.Context, req *proto.ProduceBlockRequest) (*proto.ProduceBlockResponse, error) {
    return h.service.ProduceBlock(ctx, req)
}
```

---

### 2. Business Layer (`internal/application/business/`)

**Purpose**: Contains all business logic, validation, metrics, and logging.

**Responsibilities**:
- ✅ Business logic and validation
- ✅ Logging with correlation IDs
- ✅ Performance metrics
- ✅ Error handling and transformation
- ✅ Orchestration between repositories

**Pattern**:
```go
func (s *Service) MethodName(ctx context.Context, req *proto.Request) (*proto.Response, error) {
    start := time.Now()
    correlationID := utils.GetCorrelationID(ctx)
    
    s.logger.Info("MethodName in business service",
        logging.String("correlation_id", correlationID),
        logging.String("field", req.Field))
    
    // Business validation
    if req.Field == "" {
        s.logger.Warn("Empty field in MethodName request",
            logging.String("correlation_id", correlationID),
            logging.String("service", "service_name"))
        return &proto.Response{
            Success:      false,
            ErrorMessage: "field is required",
        }, nil
    }
    
    // Delegate to repository
    return s.repo.MethodName(ctx, req)
}
```

**Example**:
```go
// internal/application/business/block_operations/block_operations_service.go
func (s *Service) ProduceBlock(ctx context.Context, req *proto.ProduceBlockRequest) (*proto.ProduceBlockResponse, error) {
    start := time.Now()
    correlationID := utils.GetCorrelationID(ctx)
    s.logger.Info("Producing block in business service",
        logging.String("correlation_id", correlationID),
        logging.String("validator", req.ValidatorId))
    
    // Business validation
    if req.ValidatorId == "" {
        s.logger.Warn("Empty validator ID in ProduceBlock request",
            logging.String("correlation_id", correlationID),
            logging.String("service", "block_operations"))
        return &proto.ProduceBlockResponse{
            Success:      false,
            ErrorMessage: "validator_id is required",
        }, nil
    }
    
    // Delegate to repository
    return s.repo.ProduceBlock(ctx, req)
}
```

---

### 3. Repository Layer (`internal/application/repository/`)

**Purpose**: Data access layer that handles database and blockchain operations.

**Responsibilities**:
- ✅ Database operations (PostgreSQL)
- ✅ Cosmos SDK blockchain operations
- ✅ Caching (Redis)
- ✅ Error standardization
- ✅ Logging with correlation IDs

**Pattern**:
```go
func (r *Repository) MethodName(ctx context.Context, req *proto.Request) (*proto.Response, error) {
    startTime := time.Now()
    correlationID := utils.GetCorrelationID(ctx)
    
    r.logger.Info("MethodName in repository",
        logging.String("correlation_id", correlationID),
        logging.String("field", req.Field))
    
    // Try Cosmos SDK blockchain first
    if r.cosmosApp != nil {
        result, err := r.methodNameOnKeeper(ctx, req)
        if err == nil {
            return result, nil
        }
        r.logger.Warn("Failed on blockchain, falling back to database",
            logging.String("correlation_id", correlationID),
            logging.Error(err))
    }
    
    // Fallback to database
    return r.methodNameOnDatabase(ctx, req)
}
```

**Error Handling Pattern**:
```go
import repoerrors "service-04/internal/application/repository"

// Use standardized errors
if err != nil {
    return nil, repoerrors.NewError(repoerrors.ErrBlockNotFound, "block not found", err)
}
```

**Example**:
```go
// internal/application/repository/block_operations/block_operations_repository.go
func (r *Repository) ProduceBlock(ctx context.Context, req *proto.ProduceBlockRequest) (*proto.ProduceBlockResponse, error) {
    startTime := time.Now()
    correlationID := utils.GetCorrelationID(ctx)
    r.logger.Info("Producing block in repository",
        logging.String("correlation_id", correlationID),
        logging.String("validator", req.ValidatorId))
    
    // Try Cosmos SDK blockchain first
    if r.cosmosApp != nil {
        result, err := r.produceBlockOnKeeper(ctx, req)
        if err == nil {
            return result, nil
        }
        r.logger.Warn("Failed to produce block on blockchain, falling back to database",
            logging.String("correlation_id", correlationID),
            logging.Error(err))
    }
    
    // Fallback to database
    return r.produceBlockOnDatabase(ctx, req)
}
```

---

## Dependency Injection

### Container Pattern

All dependencies are managed through a centralized `Container`:

**Location**: `internal/application/container-grpc.go`

**Structure**:
```go
type Container struct {
    // Infrastructure
    config            *config.Config
    logger            *logging.Logger
    db                *database.PostgreSQLManager
    cosmosApp         *app.USCApp
    blockchainStorage *storage.StateManager
    
    // Domain-specific components (12 domains)
    BlockRepository *blockrepo.Repository
    BlockService    *blockbiz.Service
    BlockHandlers   *blockhandlers.Handlers
    
    // ... (11 more domains)
}
```

**Initialization Order**:
1. Repositories (data access)
2. Business Services (business logic)
3. Handlers (API layer)

**Example**:
```go
// Initialize repositories first
c.BlockRepository = blockrepo.NewRepository(c.db, c.cosmosApp, c.blockchainStorage, c.redisManager, c.logger)

// Then business services
c.BlockService = blockbiz.NewService(c.BlockRepository, c.cosmosApp, c.blockchainStorage, c.logger)

// Finally handlers
c.BlockHandlers = blockhandlers.NewHandlers(c.BlockService)
```

---

## Correlation IDs

### Purpose

Correlation IDs enable **end-to-end request tracing** across all layers (API → Business → Repository → Database).

### Implementation

**Utility**: `internal/application/utils/context_helpers.go`

**Usage Pattern**:
```go
correlationID := utils.GetCorrelationID(ctx)
s.logger.Info("Operation in service",
    logging.String("correlation_id", correlationID),
    logging.String("field", req.Field))
```

### Support

- ✅ **gRPC metadata**: `x-correlation-id` header
- ✅ **Context values**: Standard context.Value()
- ✅ **Both incoming and outgoing**: Supports both directions

### Coverage

- ✅ **Handlers**: Not required (pure delegation)
- ✅ **Business Layer**: 100% (12/12 files)
- ✅ **Repository Layer**: 100% (12/12 files, 58 methods)

---

## Error Handling

### Standardized Errors

**Location**: `internal/application/repository/errors.go`

**Error Code Ranges**:
- `2000-2099`: Block operations
- `2100-2199`: Transaction operations
- `2200-2299`: USC Coin operations
- `2300-2399`: Smart Contract operations
- `2400-2499`: NFT Token operations
- `2500-2599`: Custom Token operations
- `2600-2699`: Product Certificate operations
- `2700-2799`: Validator operations
- `2800-2899`: Network operations
- `2900-2999`: Streaming operations
- `3000-3099`: Store Bridge operations
- `3100-3199`: Store Network operations
- `3200-3299`: Common repository errors

**Usage Pattern**:
```go
import repoerrors "service-04/internal/application/repository"

// Create standardized error
if err != nil {
    return nil, repoerrors.NewError(repoerrors.ErrBlockNotFound, "block not found", err)
}

// Check error type
if repoerrors.IsError(err, repoerrors.ErrBlockNotFound) {
    // Handle specific error
}
```

**Coverage**: 100% (12/12 repository files, ~80+ errors replaced)

---

## Logging Conventions

### Structured Logging

All logging uses the shared library's structured logging:

```go
import "github.com/usc-platform/shared/logging"

s.logger.Info("Operation description",
    logging.String("correlation_id", correlationID),
    logging.String("field", value),
    logging.Int("count", count),
    logging.Duration("duration", duration),
    logging.Error(err))
```

### Log Levels

- **Info**: Normal operation flow
- **Warn**: Recoverable issues, fallbacks
- **Error**: Errors that need attention
- **Debug**: Detailed debugging information

### Required Fields

- ✅ **correlation_id**: Always include for traceability
- ✅ **service**: Service name for filtering
- ✅ **operation**: Operation name for context

### Example

```go
correlationID := utils.GetCorrelationID(ctx)
s.logger.Info("Producing block in business service",
    logging.String("correlation_id", correlationID),
    logging.String("validator", req.ValidatorId),
    logging.Int("transaction_count", len(req.TransactionHashes)))

s.logger.Warn("Failed to produce block on blockchain, falling back to database",
    logging.String("correlation_id", correlationID),
    logging.Error(err))
```

---

## Code Organization

### Directory Structure

```
internal/application/
├── handlers/              # gRPC handlers (12 domains)
│   ├── block_operations/
│   ├── transaction_operations/
│   └── ...
├── business/              # Business logic (12 domains)
│   ├── block_operations/
│   ├── transaction_operations/
│   └── ...
├── repository/            # Data access (12 domains)
│   ├── block_operations/
│   ├── transaction_operations/
│   ├── errors.go          # Standardized errors
│   └── ...
├── utils/                 # Shared utilities
│   └── context_helpers.go # Correlation ID utilities
└── container-grpc.go      # Dependency injection container
```

### Domain Organization

Each domain follows the same structure:

```
{domain}_operations/
├── {domain}_operations_handlers.go    # Handlers
├── {domain}_operations_service.go     # Business logic
└── {domain}_operations_repository.go  # Data access
```

### 12 Domains

1. **Block Operations**: Block production, validation, querying
2. **Transaction Operations**: Transaction submission, querying
3. **USC Coin Operations**: USC token transfers, balance queries
4. **Smart Contract Operations**: Contract deployment, execution
5. **NFT Token Operations**: NFT minting, transfer, burning
6. **Custom Token Operations**: Custom token creation, minting
7. **Product Certificate Operations**: Certificate creation, verification
8. **Validator Operations**: Validator registration, staking
9. **Network Operations**: Network info, peer management
10. **Streaming Operations**: Real-time data streaming
11. **Store Bridge Operations**: Cross-chain bridge operations
12. **Store Network Operations**: Network state synchronization

---

## Best Practices

### 1. Layer Separation

- ✅ **Handlers**: Pure delegation only
- ✅ **Business**: All business logic, validation, metrics
- ✅ **Repository**: Data access only

### 2. Error Handling

- ✅ Use standardized errors from `repository/errors.go`
- ✅ Always include correlation IDs in error logs
- ✅ Provide meaningful error messages

### 3. Logging

- ✅ Always include correlation IDs
- ✅ Use structured logging (key-value pairs)
- ✅ Log at appropriate levels (Info/Warn/Error/Debug)

### 4. Performance

- ✅ Try Cosmos SDK blockchain first, fallback to database
- ✅ Use Redis for caching when appropriate
- ✅ Record performance metrics in business layer

### 5. Code Quality

- ✅ Keep files under 333 lines (project rule)
- ✅ Use panic recovery in critical paths
- ✅ Add logging before panic statements
- ✅ Document complex logic

### 6. Testing

- ✅ Test each layer independently
- ✅ Mock dependencies in tests
- ✅ Test error paths and edge cases

---

## Summary

Service-04 follows a **clean, layered architecture** with:

- ✅ **Clear separation of concerns** (Handlers → Business → Repository)
- ✅ **End-to-end traceability** (Correlation IDs)
- ✅ **Standardized error handling** (Error codes and messages)
- ✅ **Consistent logging** (Structured logging with correlation IDs)
- ✅ **Dependency injection** (Centralized Container)
- ✅ **12 domain services** (Each with handlers, business, repository)

This architecture ensures **maintainability**, **testability**, and **observability** across the entire service.

---

**Last Updated**: 2024
**Version**: 1.0

