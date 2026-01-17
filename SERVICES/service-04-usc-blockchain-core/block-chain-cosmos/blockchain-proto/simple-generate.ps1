# Simple proto generation - no copying
Write-Host "Simple proto generation..." -ForegroundColor Green

# Clean previous files
Get-ChildItem -Path "." -Recurse -Filter "*.pb.go" | Remove-Item -Force -ErrorAction SilentlyContinue

# Generate each proto file with updated paths and package names
$files = @(
    "cosmos/base/v1beta1/coin.proto",
    "cosmos/base/query/v1beta1/pagination.proto",
    "cosmos/tx/v1beta1/tx.proto",
    "usc/usc_coin/v1/tx.proto",
    "usc/block/v1/tx.proto",
    "usc/store_bridge/v1/tx.proto",
    "usc/product_certificate/v1/tx.proto",
    "usc/smart_contract/v1/tx.proto",
    "usc/monitoring/v1/tx.proto",
    "usc/network/v1/tx.proto",
    "usc/nft_token/v1/tx.proto",
    "usc/performance/v1/tx.proto",
    "usc/store_network/v1/tx.proto",
    "usc/streaming/v1/tx.proto",
    "usc/custom_token/v1/tx.proto",
    "usc/validator/v1/tx.proto",
    "usc/transaction/v1/tx.proto"
)

foreach ($file in $files) {
    $protoFile = $file
    $outputDir = Split-Path $file -Parent
    
    Write-Host "Processing: $protoFile" -ForegroundColor Yellow
    
    # Check if proto file exists
    if (Test-Path $protoFile) {
        # Generate with protoc (both Go and gRPC)
        $protocCmd = "protoc --proto_path=. --proto_path=third_party --go_out=`"$outputDir`" --go_opt=paths=source_relative --go-grpc_out=`"$outputDir`" --go-grpc_opt=paths=source_relative `"$protoFile`""
        
        try {
            Invoke-Expression $protocCmd
            Write-Host "  Generated in: $outputDir" -ForegroundColor Green
        }
        catch {
            Write-Host "  Error generating $protoFile : $($_.Exception.Message)" -ForegroundColor Red
        }
    }
    else {
        Write-Host "  File not found: $protoFile" -ForegroundColor Red
    }
}

Write-Host "All proto files generated!" -ForegroundColor Green

# Show final structure
Write-Host "Final structure:" -ForegroundColor Yellow
Get-ChildItem -Path "." -Recurse -Filter "*.pb.go" | ForEach-Object { 
    $relativePath = $_.FullName.Replace((Get-Location).Path + "\", "")
    Write-Host "  - $relativePath" -ForegroundColor Cyan 
}
