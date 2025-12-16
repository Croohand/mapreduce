#!/usr/bin/env bash
# Build script for MapReduce user operations
# $1 = transaction path (e.g., /path/to/transactions/txid)

set -e

TX_PATH="$1"
SRC_PATH="$TX_PATH/src"
BIN_PATH="$TX_PATH/bin"
MAIN_PATH="$SRC_PATH/main"
USERLIB_PATH="$SRC_PATH/mruserlib"

# Clean and create bin directory
rm -rf "$BIN_PATH"
mkdir -p "$BIN_PATH"

# Create go.mod for mruserlib
cat > "$USERLIB_PATH/go.mod" << 'EOF'
module mruserlib

go 1.21
EOF

# Create go.mod for main
cat > "$MAIN_PATH/go.mod" << 'EOF'
module main

go 1.21

require mruserlib v0.0.0

replace mruserlib => ../mruserlib
EOF

# Build the binary
cd "$MAIN_PATH"
go build -o "$BIN_PATH/main" .
