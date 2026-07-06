# Contract: Raw Object API

## GET /api/raw (NEW)

Returns the complete, unmodified Kubernetes object for a given GVR + namespace + name, bypassing
every curated DTO. Backs FR-011/FR-012 (the YAML tree tab).

**Query parameters** (all required):

| Param | Example |
|---|---|
| `group` | `cluster.x-k8s.io` |
| `version` | `v1beta1` |
| `resource` | `clusters` |
| `namespace` | `default` |
| `name` | `capi-workload` |

**Response 200**: the object exactly as the Kubernetes API server returns it (full `metadata`,
`spec`, `status`) — e.g. for `resource=clusters`:

```json
{
  "apiVersion": "cluster.x-k8s.io/v1beta1",
  "kind": "Cluster",
  "metadata": { "name": "capi-workload", "namespace": "default", "...": "..." },
  "spec": { "...": "every Spec field, including ones the curated Cluster DTO drops" },
  "status": { "...": "same as today's curated DTO — Status was already embedded wholesale" }
}
```

**Errors**:
- `400` if any required query parameter is missing, or `group`/`version`/`resource` don't resolve to
  a known GVR.
- `404` if the object doesn't exist (already deleted, wrong namespace/name).
- `500` on any other backend/API-server error, with an actionable message (Constitution I).

**Security note**: this endpoint performs a read-only `Get` scoped to whatever GVR/namespace/name is
requested — it grants no more access than the dynamic client already has (the same credentials
already used by the Docker infra fetchers and WS watchers); it does not expose cluster-wide list/watch,
only a single named object per request.

## Frontend GVR lookup (no new backend surface)

Each detail screen already implicitly knows its own Kind (Clusters screen → `clusters.cluster.x-k8s.io/v1beta1`,
Machines → `machines.cluster.x-k8s.io/v1beta1`, MachineDeployments →
`machinedeployments.cluster.x-k8s.io/v1beta1`; the infra variants use the provider-specific GVR
already defined for their WS object type in `webserver/internal/web/watchers/*.go`, mirrored as a
small frontend constant table). No backend discovery endpoint is needed — the caller already knows
what it's asking for.
