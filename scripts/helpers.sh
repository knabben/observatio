#!/usr/bin/env bash
# Shared helper library — sourced by every numbered script.
# Do not execute directly.

# ---------------------------------------------------------------------------
# Environment defaults
# ---------------------------------------------------------------------------
: "${MGMT_CLUSTER_NAME:=capi-mgmt}"
: "${CLUSTER_NAME:=capi-workload}"
: "${KUBERNETES_VERSION:=v1.30.0}"
: "${CONTROL_PLANE_MACHINE_COUNT:=1}"
: "${WORKER_MACHINE_COUNT:=1}"
: "${VELERO_NAMESPACE:=velero}"
: "${BACKUP_NAME:=capi-secrets-backup}"
: "${BACKUP_DATA_DIR:=.backup-data}"
: "${GENERATED_DIR:=generated}"
: "${SCRIPT_TIMEOUT:=300}"

# ---------------------------------------------------------------------------
# Colours
# ---------------------------------------------------------------------------
_BLUE='\033[0;34m'
_GREEN='\033[0;32m'
_YELLOW='\033[1;33m'
_RED='\033[0;31m'
_BOLD='\033[1m'
_RESET='\033[0m'

# ---------------------------------------------------------------------------
# Logging
# ---------------------------------------------------------------------------
info()    { echo -e "${_BLUE}[INFO]${_RESET}  $*"; }
success() { echo -e "${_GREEN}[OK]${_RESET}    $*"; }
warn()    { echo -e "${_YELLOW}[WARN]${_RESET}  $*"; }
error()   { echo -e "${_RED}[ERROR]${_RESET} $*" >&2; exit 1; }

mark() {
  # No-op unless both TIMING_LOG and RECORD_START_MS are exported by 01-record.sh.
  [[ -z "${TIMING_LOG:-}" || -z "${RECORD_START_MS:-}" ]] && return 0
  local now_ms offset_ms
  now_ms=$(date +%s%3N)
  offset_ms=$((now_ms - RECORD_START_MS))
  echo "{\"label\":\"$*\",\"offset_ms\":${offset_ms}}" >> "$TIMING_LOG"
}

step() {
  local n="$1"; shift
  mark "Step ${n} — $*"
  echo ""
  echo -e "${_BOLD}[STEP $n]${_RESET} $*"
  echo "────────────────────────────────────────"
}

banner() {
  local msg="$*"
  mark "$msg"
  local len=${#msg}
  local border
  border=$(printf '═%.0s' $(seq 1 $((len + 4))))
  echo ""
  echo -e "${_BOLD}╔${border}╗${_RESET}"
  echo -e "${_BOLD}║  ${msg}  ║${_RESET}"
  echo -e "${_BOLD}╚${border}╝${_RESET}"
  echo ""
}

# ---------------------------------------------------------------------------
# Tool checks
# ---------------------------------------------------------------------------
check_tool() {
  local name="$1"
  local min_ver="${2:-}"
  if ! command -v "$name" &>/dev/null; then
    echo -e "  ${_RED}[ERROR]${_RESET} $name — not found (install it before running the demo)"
    return 1
  fi
  local ver=""
  case "$name" in
    docker)     ver=$(docker version --format '{{.Client.Version}}' 2>/dev/null || true) ;;
    kind)       ver=$(kind version 2>/dev/null | awk '{print $2}' | tr -d 'v' || true) ;;
    kubectl)    ver=$(kubectl version --client -o json 2>/dev/null | \
                      python3 -c "import sys,json; print(json.load(sys.stdin)['clientVersion']['gitVersion'].lstrip('v'))" 2>/dev/null || true) ;;
    clusterctl) ver=$(clusterctl version -o short 2>/dev/null | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+' | tr -d 'v' || true) ;;
    velero)     ver=$(velero version --client-only 2>/dev/null | awk '/Version:/{print $2}' | tr -d 'v' || true) ;;
    *)          ver=$("$name" --version 2>/dev/null | head -1 | grep -oE '[0-9]+\.[0-9]+(\.[0-9]+)?' || true) ;;
  esac
  if [[ -z "$ver" ]]; then
    echo -e "  ${_YELLOW}[WARN]${_RESET}  $name — found (version unknown)"
    return 0
  fi
  echo -e "  ${_GREEN}[OK]${_RESET}    $name v$ver"
}

# ---------------------------------------------------------------------------
# Cluster / resource existence
# ---------------------------------------------------------------------------
cluster_exists() {
  kind get clusters 2>/dev/null | grep -q "^${1}$"
}

secret_exists() {
  local ns="$1" name="$2"
  kubectl get secret "$name" -n "$ns" &>/dev/null
}

# ---------------------------------------------------------------------------
# Wait functions
# ---------------------------------------------------------------------------
wait_for_pods() {
  local namespace="$1"
  local label="$2"
  local timeout="${3:-${SCRIPT_TIMEOUT}}"
  local elapsed=0
  info "Waiting for pods ready in namespace '$namespace' (label: $label)..."
  while [[ $elapsed -lt $timeout ]]; do
    local total ready
    total=$(kubectl get pods -n "$namespace" -l "$label" --no-headers 2>/dev/null | wc -l | tr -d ' ')
    ready=$(kubectl get pods -n "$namespace" -l "$label" --no-headers 2>/dev/null | \
            awk '{print $2}' | awk -F'/' '$1==$2{c++}END{print c+0}')
    if [[ "$total" -gt 0 && "$ready" -eq "$total" ]]; then
      success "All $total pod(s) ready in '$namespace' (${elapsed}s)"
      return 0
    fi
    sleep 5; elapsed=$((elapsed + 5))
    printf "\r${_BLUE}[INFO]${_RESET}  Waiting... ${elapsed}s / ${timeout}s (${ready:-0}/${total:-0} ready)   "
  done
  echo ""
  error "Timeout after ${timeout}s waiting for pods in '$namespace'"
}

wait_for_deployment() {
  local namespace="$1"
  local name="$2"
  local timeout="${3:-${SCRIPT_TIMEOUT}}"
  info "Waiting for deployment '$name' in '$namespace'..."
  kubectl rollout status deployment/"$name" -n "$namespace" --timeout="${timeout}s" && \
    success "Deployment '$name' ready"
}

wait_for_cluster_phase() {
  local cluster="$1"
  local phase="$2"
  local timeout="${3:-${SCRIPT_TIMEOUT}}"
  local elapsed=0
  info "Waiting for cluster '$cluster' to reach phase '$phase'..."
  while [[ $elapsed -lt $timeout ]]; do
    local current
    current=$(kubectl get cluster "$cluster" -o jsonpath='{.status.phase}' 2>/dev/null || true)
    if [[ "$current" == "$phase" ]]; then
      success "Cluster '$cluster' is $phase (${elapsed}s)"
      return 0
    fi
    sleep 10; elapsed=$((elapsed + 10))
    printf "\r${_BLUE}[INFO]${_RESET}  Phase: ${current:-Pending} → waiting for ${phase} ... ${elapsed}s   "
  done
  echo ""
  error "Timeout after ${timeout}s: cluster '$cluster' phase is '$(kubectl get cluster "$cluster" -o jsonpath='{.status.phase}' 2>/dev/null)'"
}

wait_for_machines_running() {
  local cluster="$1"
  local timeout="${2:-${SCRIPT_TIMEOUT}}"
  local elapsed=0
  info "Waiting for all Machines in cluster '$cluster' to be Running..."
  while [[ $elapsed -lt $timeout ]]; do
    local total not_running
    total=$(kubectl get machines -A -l "cluster.x-k8s.io/cluster-name=${cluster}" \
            --no-headers 2>/dev/null | wc -l | tr -d ' ')
    not_running=$(kubectl get machines -A -l "cluster.x-k8s.io/cluster-name=${cluster}" \
                  --no-headers 2>/dev/null | grep -vc "Running" || true)
    if [[ "$total" -gt 0 && "$not_running" -eq 0 ]]; then
      success "All $total Machine(s) Running (${elapsed}s)"
      return 0
    fi
    sleep 10; elapsed=$((elapsed + 10))
    printf "\r${_BLUE}[INFO]${_RESET}  Machines: $((total - not_running))/${total} Running ... ${elapsed}s   "
  done
  echo ""
  error "Timeout after ${timeout}s: not all Machines Running"
}

wait_for_velero_bsl() {
  local timeout="${1:-${SCRIPT_TIMEOUT}}"
  local elapsed=0
  info "Waiting for Velero BackupStorageLocation to be Available..."
  while [[ $elapsed -lt $timeout ]]; do
    local phase
    phase=$(kubectl get backupstoragelocation -n "$VELERO_NAMESPACE" \
            -o jsonpath='{.items[0].status.phase}' 2>/dev/null || true)
    if [[ "$phase" == "Available" ]]; then
      success "BackupStorageLocation is Available (${elapsed}s)"
      return 0
    fi
    sleep 5; elapsed=$((elapsed + 5))
    printf "\r${_BLUE}[INFO]${_RESET}  BSL phase: ${phase:-Unavailable} ... ${elapsed}s   "
  done
  echo ""
  error "Timeout after ${timeout}s: BSL not Available"
}

# ---------------------------------------------------------------------------
# Backup data export/import
# ---------------------------------------------------------------------------
export_backup_data() {
  # minio/minio:latest has no tar/find — use mc find + mc cat | base64 per object
  local backup_name="$1"
  local dest_dir="$2"
  mkdir -p "$dest_dir"
  local minio_pod
  minio_pod=$(kubectl get pod -n "$VELERO_NAMESPACE" -l app=minio \
              -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)
  if [[ -z "$minio_pod" ]]; then
    error "MinIO pod not found in namespace '$VELERO_NAMESPACE'"
  fi
  info "Exporting backup data from MinIO pod '$minio_pod' → $dest_dir"
  local objects
  objects=$(kubectl exec -n "$VELERO_NAMESPACE" "$minio_pod" -- \
    sh -c 'mc find local/velero 2>/dev/null')
  if [[ -z "$objects" ]]; then
    error "No backup files found in MinIO bucket for backup '$backup_name'"
  fi
  while IFS= read -r obj_full_path; do
    local rel_path="${obj_full_path#local/velero/}"
    local local_path="${dest_dir}/${rel_path}"
    mkdir -p "$(dirname "$local_path")"
    kubectl exec -n "$VELERO_NAMESPACE" "$minio_pod" -- \
      sh -c "mc cat '${obj_full_path}' | base64" | base64 -d > "$local_path"
  done <<< "$objects"
  local file_count
  file_count=$(find "$dest_dir" -type f | wc -l | tr -d ' ')
  success "Backup exported to ${dest_dir}/ (${file_count} files)"
}

import_backup_data() {
  # minio/minio:latest has no tar — stream each file via base64 into mc pipe
  local src_dir="$1"
  local minio_pod
  minio_pod=$(kubectl get pod -n "$VELERO_NAMESPACE" -l app=minio \
              -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)
  if [[ -z "$minio_pod" ]]; then
    error "MinIO pod not found in namespace '$VELERO_NAMESPACE'"
  fi
  info "Importing backup data into MinIO pod '$minio_pod'..."
  local file_count=0
  while IFS= read -r filepath; do
    local rel_path="${filepath#${src_dir}/}"
    base64 "$filepath" | kubectl exec -n "$VELERO_NAMESPACE" "$minio_pod" -i -- \
      sh -c "base64 -d | mc pipe 'local/velero/${rel_path}'"
    (( file_count++ )) || true
  done < <(find "$src_dir" -type f)
  success "Backup data imported into MinIO (${file_count} files)"
}

# ---------------------------------------------------------------------------
# MinIO / Velero install — workload-cluster architecture
# ---------------------------------------------------------------------------
# MinIO lives in the workload cluster so backup data survives management
# cluster loss. Velero in the management cluster reaches it via NodePort 30900.

install_minio_workload() {
  local kubeconfig="$1"
  local ns="${VELERO_NAMESPACE}"

  if kubectl --kubeconfig "$kubeconfig" get deployment minio -n "$ns" &>/dev/null; then
    info "MinIO deployment already exists in workload cluster"
  else
    info "Creating namespace '$ns' in workload cluster..."
  kubectl --kubeconfig "$kubeconfig" create namespace "$ns" \
    --dry-run=client -o yaml | kubectl --kubeconfig "$kubeconfig" apply -f -

  kubectl --kubeconfig "$kubeconfig" apply -n "$ns" -f - <<'MINIO_WL_EOF'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: minio
  labels:
    app: minio
spec:
  replicas: 1
  selector:
    matchLabels:
      app: minio
  template:
    metadata:
      labels:
        app: minio
    spec:
      containers:
      - name: minio
        image: minio/minio:latest
        args: ["server", "/data", "--console-address", ":9090"]
        env:
        - name: MINIO_ROOT_USER
          value: "minio"
        - name: MINIO_ROOT_PASSWORD
          value: "minio123"
        ports:
        - containerPort: 9000
        - containerPort: 9090
        volumeMounts:
        - name: storage
          mountPath: /data
      volumes:
      - name: storage
        emptyDir: {}
---
apiVersion: v1
kind: Service
metadata:
  name: minio
spec:
  type: NodePort
  selector:
    app: minio
  ports:
  - name: api
    port: 9000
    targetPort: 9000
    nodePort: 30900
MINIO_WL_EOF

  local elapsed=0 timeout="${SCRIPT_TIMEOUT}"
  while [[ $elapsed -lt $timeout ]]; do
    local total ready
    total=$(kubectl --kubeconfig "$kubeconfig" get pods -n "$ns" -l app=minio \
            --no-headers 2>/dev/null | wc -l | tr -d ' ')
    ready=$(kubectl --kubeconfig "$kubeconfig" get pods -n "$ns" -l app=minio \
            --no-headers 2>/dev/null | awk '{print $2}' | \
            awk -F'/' '$1==$2{c++}END{print c+0}')
    if [[ "$total" -gt 0 && "$ready" -eq "$total" ]]; then
      success "MinIO ready in workload cluster (${elapsed}s)"
      break
    fi
    sleep 5; elapsed=$((elapsed + 5))
    printf "\r${_BLUE}[INFO]${_RESET}  Waiting for MinIO in workload cluster... ${elapsed}s   "
  done
  [[ $elapsed -lt $timeout ]] || error "Timeout waiting for MinIO in workload cluster"
  fi  # end: deployment not yet installed

  # Always ensure the bucket exists — the early-return path above skipped this
  local minio_pod
  minio_pod=$(kubectl --kubeconfig "$kubeconfig" get pod -n "$ns" -l app=minio \
              -o jsonpath='{.items[0].metadata.name}')
  [[ -n "$minio_pod" ]] || error "MinIO pod not found in workload cluster namespace '$ns'"
  kubectl --kubeconfig "$kubeconfig" exec -n "$ns" "$minio_pod" -- \
    sh -c 'mc alias set local http://localhost:9000 minio minio123 --quiet && \
           mc mb --ignore-existing local/velero --quiet'
  success "MinIO bucket 'velero' ready in workload cluster"
}

get_workload_minio_url() {
  local kubeconfig="$1"
  local node_name node_ip
  node_name=$(kubectl --kubeconfig "$kubeconfig" get nodes \
    -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)
  [[ -n "$node_name" ]] || error "Cannot determine workload cluster node name"
  node_ip=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' \
    "$node_name" 2>/dev/null | tr -d '\n' | cut -d' ' -f1)
  [[ -n "$node_ip" ]] || error "Cannot get Docker IP for workload node '$node_name'"
  echo "http://${node_ip}:30900"
}

install_velero_in_mgmt() {
  local minio_url="$1"

  if kubectl get deployment velero -n "$VELERO_NAMESPACE" &>/dev/null; then
    info "Velero already installed — skipping"
    return 0
  fi

  local creds_file
  creds_file=$(mktemp)
  trap "rm -f $creds_file" RETURN
  printf '[default]\naws_access_key_id=minio\naws_secret_access_key=minio123\n' > "$creds_file"
  velero install \
    --provider aws \
    --plugins velero/velero-plugin-for-aws:v1.9.0 \
    --bucket velero \
    --secret-file "$creds_file" \
    --use-volume-snapshots=false \
    --namespace "$VELERO_NAMESPACE" \
    --backup-location-config \
      "region=minio,s3ForcePathStyle=true,s3Url=${minio_url}" \
    --wait
  wait_for_velero_bsl
  success "Velero installed — no credentials on disk"
}

# ---------------------------------------------------------------------------
# CAPI helpers
# ---------------------------------------------------------------------------
display_capi_secrets() {
  local cluster="$1"
  local secrets=("${cluster}-ca" "${cluster}-etcd" "${cluster}-proxy" "${cluster}-sa" "${cluster}-kubeconfig")
  echo ""
  info "CAPI certificate secrets for cluster '$cluster':"
  echo "  ┌─────────────────────────────────────────┬───────────────────────────────────────────┐"
  echo "  │ Secret Name                             │ Role                                      │"
  echo "  ├─────────────────────────────────────────┼───────────────────────────────────────────┤"
  echo "  │ ${cluster}-ca                           │ Cluster CA — signs all node certs         │"
  echo "  │ ${cluster}-etcd                         │ etcd CA — signs etcd peer/client certs    │"
  echo "  │ ${cluster}-proxy                        │ Front-proxy CA — aggregation layer        │"
  echo "  │ ${cluster}-sa                           │ Service account signing key pair          │"
  echo "  │ ${cluster}-kubeconfig                   │ Admin kubeconfig for workload cluster     │"
  echo "  └─────────────────────────────────────────┴───────────────────────────────────────────┘"
  echo ""
  for s in "${secrets[@]}"; do
    if secret_exists default "$s"; then
      local created
      created=$(kubectl get secret "$s" -n default -o jsonpath='{.metadata.creationTimestamp}' 2>/dev/null)
      echo -e "  ${_GREEN}✓${_RESET} $s  (created: $created)"
    else
      echo -e "  ${_RED}✗${_RESET} $s  — NOT FOUND"
    fi
  done
  echo ""
}

