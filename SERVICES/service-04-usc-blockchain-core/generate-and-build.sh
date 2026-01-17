#!/bin/bash

# 🔐 Service - Generate Proto Code and Build
# This script generates Go code from proto files and builds the service

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SERVICE_DIR="$SCRIPT_DIR"
PROTO_DIR="$SERVICE_DIR/proto"

print_status "Starting Service build process..."
print_status "Service directory: $SERVICE_DIR"
print_status "Proto directory: $PROTO_DIR"

# Check if we're in the correct directory
if [ ! -d "$PROTO_DIR" ]; then
    print_error "Proto directory not found: $PROTO_DIR"
    exit 1
fi

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    print_error "protoc is not installed. Please install Protocol Buffers compiler."
    print_status "Install instructions:"
    print_status "  - Windows: Download from https://github.com/protocolbuffers/protobuf/releases"
    print_status "  - macOS: brew install protobuf"
    print_status "  - Ubuntu: apt-get install protobuf-compiler"
    exit 1
fi

# Check if protoc-gen-go is installed
if ! command -v protoc-gen-go &> /dev/null; then
    print_error "protoc-gen-go is not installed. Please install it:"
    print_status "go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
    exit 1
fi

# Check if protoc-gen-go-grpc is installed
if ! command -v protoc-gen-go-grpc &> /dev/null; then
    print_error "protoc-gen-go-grpc is not installed. Please install it:"
    print_status "go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
    exit 1
fi

print_success "All required tools are installed"

# Change to proto directory
cd "$PROTO_DIR"

print_status "Generating Go code from proto files..."

# List all proto files
PROTO_FILES=$(find . -name "*.proto" -type f | sort)
print_status "Found proto files:"
for file in $PROTO_FILES; do
    echo "  - $file"
done

# Generate Go code from all proto files
print_status "Running protoc to generate Go code..."

# Generate protobuf Go code and gRPC code
protoc \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    *.proto

if [ $? -eq 0 ]; then
    print_success "Proto code generation completed successfully"
else
    print_error "Proto code generation failed"
    exit 1
fi

# List generated files
print_status "Generated Go files:"
for file in *.pb.go; do
    if [ -f "$file" ]; then
        echo "  - $file"
    fi
done

# Change back to service directory
cd "$SERVICE_DIR"

print_status "Running go mod tidy..."
go mod tidy

if [ $? -eq 0 ]; then
    print_success "Go modules updated successfully"
else
    print_error "Go mod tidy failed"
    exit 1
fi

print_status "Building the service..."

# Build the service
go build -buildvcs=false ./...

if [ $? -eq 0 ]; then
    print_success "service built successfully!"
    print_status "Build artifacts:"
    
    # Check for built binaries
    if [ -f "cmd.exe" ]; then
        echo "  - cmd.exe (Windows binary)"
    fi
    
    if [ -f "cmd" ]; then
        echo "  - cmd (Unix binary)"
    fi
    
    # List any other build artifacts
    find . -name "*.exe" -o -name "*.bin" -o -name "service-04-usc-blockchain-core*" | while read -r artifact; do
        if [ -f "$artifact" ]; then
            echo "  - $(basename "$artifact")"
        fi
    done
    
else
    print_error "Build failed"
    exit 1
fi

print_success "🎉 Service build process completed successfully!"
print_status "You can now run the service with:"
print_status "  ./cmd (on Unix/Linux/macOS)"
print_status "  ./cmd.exe (on Windows)"

# Optional: Run tests if they exist
if [ -d "tests" ] || [ -d "test" ]; then
    print_status "Running tests..."
    go test ./...
    if [ $? -eq 0 ]; then
        print_success "All tests passed!"
    else
        print_warning "Some tests failed, but build was successful"
    fi
fi

print_status "Build process completed at $(date)"
