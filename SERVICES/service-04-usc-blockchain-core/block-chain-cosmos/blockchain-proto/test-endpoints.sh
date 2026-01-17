#!/bin/bash
# Test script for gRPC 9090 & REST API 1317

set -e

GRPC_ADDR="localhost:9090"
REST_ADDR="http://localhost:1317"

echo "🧪 Testing Cosmos SDK Endpoints"
echo "=================================="
echo ""

# Check if grpcurl is installed
if ! command -v grpcurl &> /dev/null; then
    echo "⚠️  grpcurl not found. Install with:"
    echo "   go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest"
    exit 1
fi

# Test 1: gRPC - List services
echo "1️⃣  Testing gRPC (9090) - List services..."
if grpcurl -plaintext "${GRPC_ADDR}" list 2>/dev/null | head -10; then
    echo "✅ gRPC endpoint accessible"
else
    echo "❌ gRPC endpoint not accessible"
fi
echo ""

# Test 2: gRPC - Query chain info
echo "2️⃣  Testing gRPC (9090) - Query chain info..."
if grpcurl -plaintext "${GRPC_ADDR}" list cosmos.base.tendermint.v1beta1.Service 2>/dev/null; then
    echo "✅ Chain query service available"
else
    echo "⚠️  Chain query service not found"
fi
echo ""

# Test 3: REST - Node info
echo "3️⃣  Testing REST API (1317) - Node info..."
if curl -s "${REST_ADDR}/cosmos/base/tendermint/v1beta1/node_info" | head -5; then
    echo "✅ REST endpoint accessible"
else
    echo "❌ REST endpoint not accessible"
fi
echo ""

# Test 4: REST - Latest block
echo "4️⃣  Testing REST API (1317) - Latest block..."
if curl -s "${REST_ADDR}/cosmos/base/tendermint/v1beta1/blocks/latest" | head -5; then
    echo "✅ Block query working"
else
    echo "⚠️  Block query failed"
fi
echo ""

# Test 5: USC-specific queries (if available)
echo "5️⃣  Testing USC Coin queries..."
echo "   Checking for usc.usc_coin.v1.Query service..."
if grpcurl -plaintext "${GRPC_ADDR}" list | grep -i usc; then
    echo "✅ USC services found"
    grpcurl -plaintext "${GRPC_ADDR}" list | grep -i usc
else
    echo "⚠️  USC services not yet registered"
fi
echo ""

echo "✅ Testing complete!"
echo ""
echo "💡 Usage examples:"
echo "   gRPC: grpcurl -plaintext ${GRPC_ADDR} list"
echo "   REST: curl ${REST_ADDR}/cosmos/base/tendermint/v1beta1/node_info"


