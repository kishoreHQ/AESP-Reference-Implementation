#!/usr/bin/env bash
# Start Agent OS kernel + Mission Control UI (monorepo)
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
KERNEL_PORT="${KERNEL_PORT:-8080}"
UI_PORT="${UI_PORT:-5173}"
USE_MOCKS="${USE_MOCKS:-0}"

cd "$ROOT"

echo "╔══════════════════════════════════════════════════════╗"
echo "║  AESP Agent OS monorepo — local dev                  ║"
echo "╚══════════════════════════════════════════════════════╝"
echo "  Kernel  →  http://127.0.0.1:${KERNEL_PORT}"
echo "  UI      →  http://127.0.0.1:${UI_PORT}"
echo "  Mocks   →  USE_MOCKS=${USE_MOCKS} (0=live Host API)"
echo ""

# Build kernel if needed
if [ ! -x bin/aespd ]; then
  echo "→ Building aespd..."
  go build -o bin/aespd ./cmd/aespd
fi

# Install UI deps if needed
if [ ! -d ui/node_modules ]; then
  echo "→ npm install (ui/)..."
  npm --prefix ui install
fi

cleanup() {
  echo ""
  echo "Stopping..."
  [[ -n "${KERNEL_PID:-}" ]] && kill "$KERNEL_PID" 2>/dev/null || true
  exit 0
}
trap cleanup INT TERM

# Start kernel if not already up
if curl -sf "http://127.0.0.1:${KERNEL_PORT}/api/v1/health" >/dev/null 2>&1 \
   || curl -sf "http://127.0.0.1:${KERNEL_PORT}/health" >/dev/null 2>&1; then
  echo "→ Kernel already listening on :${KERNEL_PORT}"
else
  echo "→ Starting aespd on :${KERNEL_PORT}"
  AESP_WORKSPACE="${AESP_WORKSPACE:-$ROOT/.aesp-workspace}" \
    ./bin/aespd serve ":${KERNEL_PORT}" &
  KERNEL_PID=$!
  sleep 1
fi

echo "→ Starting Mission Control UI on :${UI_PORT}"
echo ""
echo "  Open:  http://127.0.0.1:${UI_PORT}"
echo "  Stop:  Ctrl+C"
echo ""

cd "$ROOT/ui"
export VITE_USE_MOCKS="$USE_MOCKS"
export VITE_API_BASE="${VITE_API_BASE:-/api}"
npm run dev -- --host 127.0.0.1 --port "$UI_PORT"
