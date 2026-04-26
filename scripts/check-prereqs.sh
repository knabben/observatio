#!/usr/bin/env bash
# Validates required tools before any build step runs.
# Usage:
#   check-prereqs.sh           — check all tools
#   check-prereqs.sh --go      — check Go only
#   check-prereqs.sh --node    — check Node + pnpm only
#
# Exit 0 if all requested tools pass; non-zero otherwise.

set -euo pipefail

MIN_GO_MAJOR=1
MIN_GO_MINOR=24
MIN_NODE_MAJOR=22

CHECK_GO=false
CHECK_NODE=false

# If no flags given, check everything
if [[ $# -eq 0 ]]; then
    CHECK_GO=true
    CHECK_NODE=true
else
    for arg in "$@"; do
        case "$arg" in
            --go)   CHECK_GO=true ;;
            --node) CHECK_NODE=true ;;
            *) echo "[check] Unknown flag: $arg" >&2; exit 1 ;;
        esac
    done
fi

FAILED=false

# --- version helpers ---

version_gte() {
    # Returns 0 (true) if $1 >= $2 using sort -V
    printf '%s\n%s\n' "$2" "$1" | sort -V -C
}

# --- Go ---
check_go() {
    if ! command -v go >/dev/null 2>&1; then
        echo "[check] go      MISSING — install from https://go.dev/doc/install"
        return 1
    fi

    local raw
    raw=$(go version 2>&1)
    # "go version go1.24.1 linux/amd64" → "1.24.1"
    local ver
    ver=$(echo "$raw" | grep -oE 'go[0-9]+\.[0-9]+(\.[0-9]+)?' | head -1 | sed 's/go//')

    local required="${MIN_GO_MAJOR}.${MIN_GO_MINOR}"
    if version_gte "$ver" "$required"; then
        echo "[check] go      ${ver} ✓"
    else
        echo "[check] go      ${ver} ✗  (need ≥ ${required} — update at https://go.dev/doc/install)"
        return 1
    fi
}

# --- Node ---
check_node() {
    if ! command -v node >/dev/null 2>&1; then
        echo "[check] node    MISSING — install Node ${MIN_NODE_MAJOR} LTS from https://nodejs.org or via nvm"
        return 1
    fi

    local ver
    ver=$(node --version 2>&1 | sed 's/v//')
    local major
    major=$(echo "$ver" | cut -d. -f1)

    if [[ "$major" -ge "$MIN_NODE_MAJOR" ]]; then
        echo "[check] node    ${ver} ✓"
    else
        echo "[check] node    ${ver} ✗  (need ≥ ${MIN_NODE_MAJOR} — https://nodejs.org)"
        return 1
    fi
}

# --- pnpm ---
check_pnpm() {
    if ! command -v pnpm >/dev/null 2>&1; then
        echo "[check] pnpm    MISSING — run: npm install -g pnpm"
        return 1
    fi
    local ver
    ver=$(pnpm --version 2>&1)
    echo "[check] pnpm    ${ver} ✓"
}

# --- run requested checks ---
if $CHECK_GO; then
    check_go || FAILED=true
fi

if $CHECK_NODE; then
    check_node  || FAILED=true
    check_pnpm  || FAILED=true
fi

if $FAILED; then
    echo ""
    echo "[check] One or more prerequisites are missing or outdated. Resolve the issues above and retry."
    exit 1
fi

echo "[check] All prerequisites satisfied."
