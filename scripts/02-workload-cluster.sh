#!/usr/bin/env bash
# Pre-conditions:  01-bootstrap.sh exited 0.
# Post-conditions: Workload cluster '$CLUSTER_NAME' Provisioned; all Machines
#                  Running; kindnet CNI installed; nodes Ready; kubeconfig saved;
#                  five CA secrets visible.
# Recovery:        Re-run (idempotent). If machines stuck >10 min, run 08-cleanup.sh.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=scripts/helpers.sh
source "$SCRIPT_DIR/helpers.sh"

mkdir -p "$GENERATED_DIR"

banner "Step 2 — Workload Cluster Provisioning"

step 1 "Generate workload cluster manifest"
MANIFEST="$GENERATED_DIR/${CLUSTER_NAME}.yaml"

if kubectl get cluster "$CLUSTER_NAME" -n default &>/dev/null; then
  info "Cluster '$CLUSTER_NAME' already exists — skipping manifest generation and apply"
else
  info "Generating manifest for cluster '$CLUSTER_NAME' (k8s $KUBERNETES_VERSION)..."
  clusterctl generate cluster "$CLUSTER_NAME" \
    --infrastructure docker \
    --flavor development \
    --kubernetes-version "$KUBERNETES_VERSION" \
    --control-plane-machine-count "$CONTROL_PLANE_MACHINE_COUNT" \
    --worker-machine-count "$WORKER_MACHINE_COUNT" \
    > "$MANIFEST"
  success "Manifest written to $MANIFEST"

  step 2 "Apply workload cluster manifest"
  kubectl apply -f "$MANIFEST"
  success "Cluster objects applied"
fi

step 3 "Wait for cluster to provision"
wait_for_cluster_phase "$CLUSTER_NAME" "Provisioned" 600
wait_for_machines_running "$CLUSTER_NAME" 600

step 4 "Save workload cluster kubeconfig"
KUBECONFIG_FILE="$GENERATED_DIR/${CLUSTER_NAME}.kubeconfig"
if [[ ! -f "$KUBECONFIG_FILE" ]]; then
  clusterctl get kubeconfig "$CLUSTER_NAME" > "$KUBECONFIG_FILE"
  success "Kubeconfig saved to $KUBECONFIG_FILE"
else
  info "Kubeconfig already exists at $KUBECONFIG_FILE"
fi

step 5 "Install kindnet CNI"
KINDNET_URL="https://raw.githubusercontent.com/aojea/kindnet/main/install-kindnet.yaml"
if KUBECONFIG="$KUBECONFIG_FILE" kubectl get daemonset kindnet -n kube-system &>/dev/null; then
  info "kindnet already installed — skipping"
else
  info "Applying kindnet (compatible with CAPD Docker nodes)..."
  KUBECONFIG="$KUBECONFIG_FILE" kubectl apply -f "$KINDNET_URL"
  success "kindnet applied"
fi
info "Waiting for nodes to become Ready..."
KUBECONFIG="$KUBECONFIG_FILE" kubectl wait node --all --for=condition=Ready --timeout=180s
success "All nodes Ready"

step 7 "Display CAPI object tree"
echo ""
clusterctl describe cluster "$CLUSTER_NAME" 2>/dev/null || \
  kubectl get cluster,machines,dockercluster,dockermachines -A 2>/dev/null
echo ""

step 8 "Display CAPI certificate secrets (the Hidden Critical State)"
display_capi_secrets "$CLUSTER_NAME"

banner "Workload cluster ready ✓"
info "Cluster:    $CLUSTER_NAME"
info "Kubeconfig: $KUBECONFIG_FILE"
info ""
info "Try it:"
info "  kubectl --kubeconfig $KUBECONFIG_FILE get nodes"
echo ""
