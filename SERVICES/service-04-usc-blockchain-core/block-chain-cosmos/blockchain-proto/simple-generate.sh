#!/bin/bash

# Simple proto generation - no copying
echo "Simple proto generation..."

# Clean previous files
find . -name "*.pb.go" -type f -delete

# Generate each proto file with updated paths and package names
files=(
    "cosmos/base/v1beta1/coin.proto"
    "cosmos/base/query/v1beta1/pagination.proto"
    "cosmos/tx/v1beta1/tx.proto"
    "usc/usc_coin/v1/tx.proto"
    "usc/block/v1/tx.proto"
    "usc/store_bridge/v1/tx.proto"
    "usc/product_certificate/v1/tx.proto"
    "usc/smart_contract/v1/tx.proto"
    "usc/monitoring/v1/tx.proto"
    "usc/network/v1/tx.proto"
    "usc/nft_token/v1/tx.proto"
    "usc/performance/v1/tx.proto"
    "usc/store_network/v1/tx.proto"
    "usc/streaming/v1/tx.proto"
    "usc/custom_token/v1/tx.proto"
    "usc/validator/v1/tx.proto"
    "usc/transaction/v1/tx.proto"
)

for file in "${files[@]}"; do
    protoFile="$file"
    outputDir=$(dirname "$file")
    
    echo "Processing: $protoFile"
    
    # Generate with protoc (both Go and gRPC)
    protoc --proto_path=. --proto_path=third_party --go_out="$outputDir" --go_opt=paths=source_relative --go-grpc_out="$outputDir" --go-grpc_opt=paths=source_relative "$protoFile"
    
    echo "  Generated in: $outputDir"
done

echo "All proto files generated!"

# Show final structure
echo "Final structure:"
find . -name "*.pb.go" -type f | while read -r file; do
    relativePath=$(echo "$file" | sed 's|^\./||')
    echo "  - $relativePath"
done
