#!/usr/bin/env bash
# UI only with MSW mocks (no kernel required)
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT/ui"
[ -d node_modules ] || npm install
echo "Open http://127.0.0.1:5173 (mocks)"
VITE_USE_MOCKS=1 npm run dev -- --host 127.0.0.1 --port 5173
