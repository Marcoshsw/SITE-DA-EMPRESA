#!/usr/bin/env sh
set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
export PORT=8080
export GIN_MODE=release
export FRONTEND_DIR="$SCRIPT_DIR/frontend"
export CORS_ORIGINS='*'
cd "$SCRIPT_DIR/backend"
go run .
