#!/usr/bin/env bash
# Smoke-tests the built single-binary (output/observatio): launches it on a throwaway
# port, asserts the UI root and SPA fallback serve the embedded frontend (200, SPA shell)
# and a live API route responds, then tears the process down. Non-zero exit on any
# failure — a missing/stale embed or broken SPA fallback must fail the build, not pass
# silently. See specs/003-screen-ui-refactor/contracts/build-verification.md.
set -uo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BINARY_PATH="${BINARY_PATH:-${REPO_ROOT}/output/observatio}"
PORT="${VERIFY_PORT:-18080}"
BASE_URL="http://127.0.0.1:${PORT}"
LOG_FILE="$(mktemp)"
SERVER_PID=""
FAIL=0

cleanup() {
  local status=$?
  if [[ -n "${SERVER_PID}" ]] && kill -0 "${SERVER_PID}" 2>/dev/null; then
    kill "${SERVER_PID}" 2>/dev/null || true
    wait "${SERVER_PID}" 2>/dev/null || true
  fi
  rm -f "${LOG_FILE}"
  exit "${status}"
}
trap cleanup EXIT

if [[ ! -x "${BINARY_PATH}" ]]; then
  echo "[verify-binary] FAIL: binary not found or not executable at ${BINARY_PATH} (run 'make build' first)"
  exit 1
fi

echo "[verify-binary] Launching ${BINARY_PATH} on :${PORT}..."
"${BINARY_PATH}" serve --address ":${PORT}" >"${LOG_FILE}" 2>&1 &
SERVER_PID=$!

echo "[verify-binary] Waiting for the server to become ready..."
ready=false
for _ in $(seq 1 30); do
  if ! kill -0 "${SERVER_PID}" 2>/dev/null; then
    echo "[verify-binary] FAIL: server process exited before becoming ready"
    cat "${LOG_FILE}"
    exit 1
  fi
  if curl -sf -o /dev/null "${BASE_URL}/api/health"; then
    ready=true
    break
  fi
  sleep 0.5
done

if [[ "${ready}" != "true" ]]; then
  echo "[verify-binary] FAIL: server did not become ready within timeout"
  cat "${LOG_FILE}"
  exit 1
fi

check_status() {
  local desc="$1" path="$2" expected="$3"
  local code
  code=$(curl -s -o /dev/null -w '%{http_code}' "${BASE_URL}${path}")
  if [[ "${code}" != "${expected}" ]]; then
    echo "[verify-binary] FAIL: ${desc} — expected ${expected}, got ${code}"
    FAIL=1
  else
    echo "[verify-binary] OK: ${desc} (${code})"
  fi
}

check_contains() {
  local desc="$1" path="$2" needle="$3"
  local body
  body=$(curl -s "${BASE_URL}${path}")
  if [[ "${body}" != *"${needle}"* ]]; then
    echo "[verify-binary] FAIL: ${desc} — response did not contain '${needle}'"
    FAIL=1
  else
    echo "[verify-binary] OK: ${desc}"
  fi
}

check_status "UI root serves the embedded SPA shell"        "/"                             200
check_contains "UI root body is the embedded frontend"      "/"                             "observātiō"
check_status "unknown client route falls back to the SPA shell" "/dashboard/does-not-exist" 200
check_status "live API route responds same-origin"          "/api/health"                   200

if [[ "${FAIL}" -ne 0 ]]; then
  echo "[verify-binary] FAILED"
  exit 1
fi

echo "[verify-binary] PASSED — UI + API served same-origin from ${BASE_URL}"
