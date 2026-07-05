#!/usr/bin/env bash
# Pre-conditions:  00-prereqs.sh exited 0; Docker daemon running.
# Post-conditions: kind cluster '$MGMT_CLUSTER_NAME' running; all four CAPI
#                  controller namespaces have Ready pods.
# Recovery:        Re-run this script (idempotent). If stuck, run 08-cleanup.sh first.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=scripts/helpers.sh
source "$SCRIPT_DIR/helpers.sh"

banner "Step 1 — Management Cluster Bootstrap"

step 1 "Create kind management cluster: $MGMT_CLUSTER_NAME"
if cluster_exists "$MGMT_CLUSTER_NAME"; then
  info "Cluster '$MGMT_CLUSTER_NAME' already exists — skipping creation"
else
  kind create cluster \
    --name "$MGMT_CLUSTER_NAME" \
    --config "$SCRIPT_DIR/config/kind-mgmt.yaml"
  success "kind cluster '$MGMT_CLUSTER_NAME' created"
fi

step 2 "Initialize Cluster API with Docker provider"
if kubectl get deployment capi-controller-manager -n capi-system &>/dev/null; then
  info "CAPI already initialized — skipping"
else
  CLUSTER_TOPOLOGY=true clusterctl init --infrastructure docker
  success "CAPI initialized"
fi

step 3 "Wait for CAPI controllers to be ready"
CAPI_NAMESPACES=(
  "capi-system:control-plane=controller-manager"
  "capd-system:control-plane=controller-manager"
  "capi-kubeadm-bootstrap-system:control-plane=controller-manager"
  "capi-kubeadm-control-plane-system:control-plane=controller-manager"
)
for entry in "${CAPI_NAMESPACES[@]}"; do
  ns="${entry%%:*}"
  label="${entry##*:}"
  wait_for_pods "$ns" "$label"
done

banner "Management cluster ready ✓"
info "Context: $(kubectl config current-context)"
info "CAPI namespaces:"
kubectl get pods -A --field-selector metadata.namespace!=kube-system \
  --no-headers 2>/dev/null | grep "controller-manager" | \
  awk '{printf "  %-45s %s\n", $1"/"$2, $4}'
echo ""
success "Run ./02-workload-cluster.sh to provision a workload cluster"
