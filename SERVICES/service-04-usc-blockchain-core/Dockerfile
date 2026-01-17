# Multi-stage build for USC USC Blockchain Core Service
# service-04: USC Blockchain Core Service

# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies including RocksDB
# Note: rocksdb-dev package provides both headers and libraries for RocksDB
# Split into multiple steps for better caching and error handling
RUN apk update && \
    apk add --no-cache --virtual .build-deps \
    git ca-certificates tzdata protobuf-dev protoc build-base \
    gcc g++ make cmake zlib-dev bzip2-dev snappy-dev lz4-dev zstd-dev perl \
    pkgconfig && \
    apk add --no-cache rocksdb-dev && \
    # Install protoc plugins
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Set CGO environment variables globally for RocksDB
# COSMOS SDK 0.53.4: RocksDB requires CGO with proper linking flags
ENV CGO_ENABLED=1
ENV CGO_CFLAGS="-I/usr/include"
# Link against RocksDB and all required compression libraries
ENV CGO_LDFLAGS="-L/usr/lib -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4 -lzstd -lpthread"
# Build optimization: Use all CPU cores for parallel compilation
ENV GOFLAGS="-buildmode=default"
# Enable Go build cache
ENV GOCACHE=/root/.cache/go-build

# Set working directory
WORKDIR /app

# Copy go mod files first
COPY SERVICES/service-04-usc-blockchain-core/go.mod SERVICES/service-04-usc-blockchain-core/go.sum ./

# Copy shared library
COPY shared ./shared

# Copy proto directory (needed for replace directive)
COPY SERVICES/service-04-usc-blockchain-core/proto ./proto

# Copy block-chain-cosmos before go mod download
COPY SERVICES/service-04-usc-blockchain-core/block-chain-cosmos ./block-chain-cosmos

# Fix go.mod replace path for Docker context
RUN sed -i 's|../../shared|./shared|g' go.mod

# Download dependencies (optimized with parallel downloads)
RUN go mod download && go mod tidy

# Copy source code (excluding go.mod and go.sum to keep our fixed versions)
# Note: proto already copied above
COPY SERVICES/service-04-usc-blockchain-core/cmd ./cmd
COPY SERVICES/service-04-usc-blockchain-core/api ./api
COPY SERVICES/service-04-usc-blockchain-core/internal ./internal
COPY SERVICES/service-04-usc-blockchain-core/configs ./configs
COPY SERVICES/service-04-usc-blockchain-core/schemas ./schemas
COPY SERVICES/service-04-usc-blockchain-core/tests ./tests
COPY SERVICES/service-04-usc-blockchain-core/tools ./tools
COPY SERVICES/service-04-usc-blockchain-core/docs ./docs
COPY SERVICES/service-04-usc-blockchain-core/migrations ./migrations

COPY SERVICES/service-04-usc-blockchain-core/README.md ./

# Enforce vendoring after all sources are present
# Note: Vendor is cached, only runs if go.mod/go.sum change
RUN go mod tidy && go mod vendor

# Note: All proto files are already generated and included in the build

# Build the application with CGO enabled for RocksDB production
# Note: rocksdb_legacy tag is required for grocksdb v1.10.1 compatibility
# Note: rocksdb build tag is required for cometbft-db RocksDB support
# COSMOS SDK 0.53.4: Ensure CGO flags are set for RocksDB linking
RUN GOOS=linux GOARCH=amd64 \
    CGO_ENABLED=1 \
    CGO_CFLAGS="-I/usr/include" \
    CGO_LDFLAGS="-L/usr/lib -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4 -lzstd -lpthread" \
    go build \
    -p $(nproc) \
    -tags "rocksdb_legacy rocksdb" \
    -ldflags="-w -s -X main.Version=${VERSION:-dev} -X main.BuildTime=$(date -u '+%Y-%m-%d_%H:%M:%S') -X main.CommitHash=${COMMIT_HASH:-unknown}" \
    -o service-04-service ./cmd/main.go

# Vendor dependencies for block-chain-cosmos to keep versions pinned
# Optimization: Only vendor if go.mod changed (Docker layer cache)
RUN cd block-chain-cosmos && \
    go mod download && \
    go mod tidy && \
    go mod vendor

# Build Cosmos SDK daemon with CGO enabled for RocksDB production
# Note: rocksdb_legacy tag is required for grocksdb v1.10.1 compatibility
# Note: rocksdb build tag is required for cometbft-db RocksDB support
# COSMOS SDK 0.53.4: Ensure CGO flags are set for RocksDB linking
# Optimization: Use parallel builds for faster compilation
RUN cd block-chain-cosmos/daemon/uscd && \
    GOOS=linux GOARCH=amd64 \
    CGO_ENABLED=1 \
    CGO_CFLAGS="-I/usr/include" \
    CGO_LDFLAGS="-L/usr/lib -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4 -lzstd -lpthread" \
    go build \
    -p $(nproc) \
    -tags "rocksdb_legacy rocksdb" \
    -ldflags="-w -s" \
    -o /app/block-chain-cosmos/daemon/uscd/uscd ./cmd/main.go

# Final stage
FROM alpine:latest

# Install runtime dependencies including RocksDB
# Note: rocksdb-dev includes both development headers and runtime libraries
RUN apk update && \
    apk --no-cache add ca-certificates tzdata wget postgresql-client redis bash \
    libstdc++ rocksdb-dev && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.24/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe && \
    rm -rf /tmp/* /var/cache/apk/*

# Create non-root user (standardized)
RUN addgroup -g 1001 -S uscgroup && \
    adduser -u 1001 -S uscuser -G uscgroup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/service-04-service .

# Copy configuration files
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/configs ./config
COPY --from=builder /app/schemas ./schemas
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/block-chain-cosmos ./block-chain-cosmos

# Make init-database script executable
RUN chmod +x /app/migrations/init-database.sh

# Copy uscd binary from builder stage
COPY --from=builder /app/block-chain-cosmos/daemon/uscd/uscd /app/block-chain-cosmos/daemon/uscd/uscd

# Create necessary directories
RUN mkdir -p /app/logs /app/data /data/rocksdb /data/rocksdb/wal /backups/rocksdb /app/block-chain-cosmos/data && \
    chown -R uscuser:uscgroup /app /data /backups && \
    chmod +x /app/block-chain-cosmos/daemon/uscd/uscd

# Switch to non-root user
USER uscuser

# Expose ports (USC Platform standard)
EXPOSE 8004   
# Main service gRPC port
EXPOSE 9004   
# Metrics endpoint
# Note: Ports 9090 (gRPC) and 1317 (REST) removed - not exposed by ABCI server
#       Query interfaces available via CometBFT RPC (26657) on service-04-cometbft

# Health check using gRPC health probe (FIXED PATH)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD /bin/grpc_health_probe -addr=0.0.0.0:8004 || exit 1

# Create startup script for Standalone ABCI Server
RUN echo '#!/bin/bash' > /app/start.sh && \
    echo 'set -e' >> /app/start.sh && \
    echo './migrations/init-database.sh' >> /app/start.sh && \
    echo '# Initialize blockchain if config directory does not exist' >> /app/start.sh && \
    echo 'if [ ! -d "/app/block-chain-cosmos/data/config" ] || [ ! -f "/app/block-chain-cosmos/data/config/genesis.json" ]; then' >> /app/start.sh && \
    echo '  echo "🔧 Initializing blockchain with fresh genesis.json..."' >> /app/start.sh && \
    echo '  cd /app/block-chain-cosmos/daemon/uscd && ./uscd init validator --chain-id usc-1 --home /app/block-chain-cosmos/data' >> /app/start.sh && \
    echo '  echo "✅ Blockchain initialized successfully"' >> /app/start.sh && \
    echo 'fi' >> /app/start.sh && \
    echo 'cd /app/block-chain-cosmos/daemon/uscd && ./uscd abci-server --abci-address tcp://0.0.0.0:26658 --db-backend rocksdb --db-dir /app/block-chain-cosmos/data --home /app/block-chain-cosmos/data > /app/logs/abci-server.log 2>&1 &' >> /app/start.sh && \
    echo 'ABCI_PID=$!' >> /app/start.sh && \
    echo 'sleep 2' >> /app/start.sh && \
    echo 'trap "kill $ABCI_PID 2>/dev/null || true" EXIT' >> /app/start.sh && \
    echo 'exec ./service-04-service' >> /app/start.sh && \
    chmod +x /app/start.sh

# Run the startup script
CMD ["/app/start.sh"]
