#!/bin/bash

# Generate Go code from Protocol Buffers for USC Shared Library
# This script generates Go code from all .proto files in the proto/ directory

set -e

echo "🚀 Generating Go code from Protocol Buffers for USC Shared Library..."

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
PROTO_DIR="$PROJECT_ROOT/proto"

echo "📁 Project root: $PROJECT_ROOT"
echo "📁 Proto directory: $PROTO_DIR"

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "❌ protoc is not installed. Please install Protocol Buffers compiler."
    echo "   Install instructions: https://grpc.io/docs/protoc-installation/"
    exit 1
fi

# Check if Go protobuf plugins are installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo "📦 Installing protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "📦 Installing protoc-gen-go-grpc..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# Check if proto directory exists
if [ ! -d "$PROTO_DIR" ]; then
    echo "❌ Proto directory not found: $PROTO_DIR"
    exit 1
fi

# Check if there are any .proto files
PROTO_FILES=$(find "$PROTO_DIR" -name "*.proto" -type f)
if [ -z "$PROTO_FILES" ]; then
    echo "❌ No .proto files found in $PROTO_DIR"
    exit 1
fi

echo "📋 Found proto files:"
echo "$PROTO_FILES" | while read -r file; do
    echo "   - $(basename "$file")"
done

# Create output directory if it doesn't exist
OUTPUT_DIR="$PROTO_DIR"
echo "📁 Output directory: $OUTPUT_DIR"

# Generate Go code
echo "🔧 Generating Go code..."
protoc \
    --go_out="$OUTPUT_DIR" \
    --go_opt=paths=source_relative \
    --go-grpc_out="$OUTPUT_DIR" \
    --go-grpc_opt=paths=source_relative \
    --proto_path="$PROTO_DIR" \
    "$PROTO_DIR"/*.proto

# Verify generated files
echo "✅ Verifying generated files..."
GENERATED_FILES=$(find "$OUTPUT_DIR" -name "*.pb.go" -type f)
if [ -z "$GENERATED_FILES" ]; then
    echo "❌ No .pb.go files were generated"
    exit 1
fi

echo "📋 Generated files:"
echo "$GENERATED_FILES" | while read -r file; do
    echo "   - $(basename "$file")"
done

# Verify generated code compiles
echo "🔍 Verifying generated code compiles..."
cd "$PROJECT_ROOT"
if go build ./proto/...; then
    echo "✅ Generated code compiles successfully"
else
    echo "❌ Generated code failed to compile"
    exit 1
fi

# Run go mod tidy to ensure dependencies are correct
echo "🧹 Running go mod tidy..."
go mod tidy

echo ""
echo "🎉 Successfully generated Go code from Protocol Buffers!"
echo "📁 Generated files are in: $OUTPUT_DIR"
echo ""
echo "📋 Summary:"
echo "   - Proto files processed: $(echo "$PROTO_FILES" | wc -l)"
echo "   - Go files generated: $(echo "$GENERATED_FILES" | wc -l)"
echo ""
echo "💡 Next steps:"
echo "   1. Import the generated packages in your Go code"
echo "   2. Use the gRPC services and message types"
echo "   3. Run 'go mod tidy' if you encounter import issues"
