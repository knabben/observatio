#!/usr/bin/env bash
# Pre-conditions:  None — safe to run from any demo stage.
# Post-conditions: No kind clusters, no orphaned CAPD containers, no generated files.
# Recovery:        Re-run (handles partial state). Last resort: kind delete clusters --all
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=scripts/helpers.sh
source "$SCRIPT_DIR/helpers.sh"

banner "Step 8 — Full Teardown"

step 1 "Delete workload cluster via CAPI (if management cluster is up)"
if cluster_exists "$MGMT_CLUSTER_NAME"; then
  if kubectl get cluster "$CLUSTER_NAME" -n default &>/dev/null; then
    info "Deleting workload cluster '$CLUSTER_NAME' via CAPI..."
    kubectl delete cluster "$CLUSTER_NAME" -n default --timeout=120s 2>/dev/null || \
      warn "CAPI deletion timed out or failed — will clean up Docker containers directly"
  else
    info "No workload cluster object found — skipping CAPI deletion"
  fi
else
  info "Management cluster not running — skipping CAPI deletion"
fi

step 2 "Delete management kind cluster"
if cluster_exists "$MGMT_CLUSTER_NAME"; then
  kind delete cluster --name "$MGMT_CLUSTER_NAME"
  success "Management cluster '$MGMT_CLUSTER_NAME' deleted"
else
  info "Management cluster '$MGMT_CLUSTER_NAME' not found — skipping"
fi

step 3 "Remove orphaned CAPD workload cluster containers"
ORPHANS=$(docker ps -aq --filter "name=${CLUSTER_NAME}" 2>/dev/null || true)
if [[ -n "$ORPHANS" ]]; then
  info "Removing orphaned containers for cluster '$CLUSTER_NAME'..."
  # shellcheck disable=SC2086
  docker rm -f $ORPHANS
  success "Orphaned containers removed"
else
  info "No orphaned containers found for '$CLUSTER_NAME'"
fi

step 4 "Remove generated files"
for path in "$BACKUP_DATA_DIR" "$GENERATED_DIR"; do
  if [[ -d "$path" ]]; then
    rm -rf "$path"
    success "Removed $path/"
  else
    info "$path/ not found — skipping"
  fi
done

banner "Environment clean ✓"
success "All demo resources removed. Workstation is back to a clean state."
