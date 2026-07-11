#!/usr/bin/env bash
# From monorepo: prefer ../../scripts/dev.sh
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
exec bash "$ROOT/scripts/dev.sh"
