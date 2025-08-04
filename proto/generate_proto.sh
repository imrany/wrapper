#!/bin/bash

# Setup script for protoc
set -e

# Define source and output directories relative to proto/
PROTO_SRC=api/v1
PROTO_OUT=gen
GOOGLEAPIS_DIR=../googleapis  # Adjust this path if needed

echo "🔄 Starting proto generation..."
echo "📁 Source directory: $PROTO_SRC"
echo "📁 Output directory: $PROTO_OUT"
echo "📁 Google APIs directory: $GOOGLEAPIS_DIR"

# Create output directory if it doesn't exist
mkdir -p "$PROTO_OUT"

# Loop through all .proto files in the source directory
for file in $PROTO_SRC/*.proto; do
  filename=$(basename "$file")
  echo "🛠️  Generating code for: $filename"

  # Generate gRPC and Go types
  protoc \
    --proto_path=. \
    --proto_path="$GOOGLEAPIS_DIR" \
    --go_out="$PROTO_OUT" \
    --go-grpc_out="$PROTO_OUT" \
    --go_opt=paths=source_relative \
    --go-grpc_opt=paths=source_relative \
    "$file"

  # Generate grpc-gateway and OpenAPI docs
  protoc \
    --proto_path=. \
    --proto_path="$GOOGLEAPIS_DIR" \
    --grpc-gateway_out="$PROTO_OUT" \
    --grpc-gateway_opt=paths=source_relative \
    --openapiv2_out="$PROTO_OUT" \
    "$file"

  echo "✅ Done: $filename"
done

echo "🎉 All proto files generated successfully!"
