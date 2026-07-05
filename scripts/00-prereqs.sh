#!/usr/bin/env bash
# Pre-conditions:  Docker daemon running; script has execute permission.
# Post-conditions: All five tools present at or above minimum versions (exit 0).
# Recovery:        Install any missing tool per README prerequisites section.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=scripts/helpers.sh
source "$SCRIPT_DIR/helpers.sh"

banner "Step 0 — Prerequisite Validation"

info "Checking required tools..."
echo ""

PASS=true

check_tool docker     "24.0" || PASS=false
check_tool kind       "0.23" || PASS=false
check_tool kubectl    "1.28" || PASS=false
check_tool clusterctl "1.7"  || PASS=false

echo ""

if [[ "$PASS" == "true" ]]; then
  success "All prerequisites satisfied — ready to run the demo"
else
  echo -e "${_RED}[ERROR]${_RESET} One or more prerequisites failed." >&2
  echo "        Install missing tools and re-run this script." >&2
  exit 1
fi
