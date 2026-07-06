# Contract: Infrastructure Detection API

## GET /api/infra/capabilities (NEW)

Returns which infrastructure providers are installed in the connected environment, and their version.
Backs FR-002, FR-004, FR-007, FR-009, FR-011, FR-012 — the frontend uses this response to decide
which listing tabs to render and never performs its own detection.

**Response 200**:

```json
{
  "docker": { "installed": true, "version": "v1.9.6" },
  "vsphere": { "installed": false, "version": "" }
}
```

- Both `docker` and `vsphere` keys are always present.
- `version` is `""` when `installed` is `false`.
- If neither provider is installed, both `installed` fields are `false` — the frontend renders the
  FR-009 "no supported infrastructure provider detected" message instead of any tab.

## GET /api/clusters/list (existing, extended)

Each cluster item gains a `provider` field (FR-001, FR-003):

```json
{
  "metadata": { "name": "workload-a" },
  "infrastructureRef": { "kind": "DockerCluster", "name": "workload-a" },
  "provider": "docker"
}
```

`provider` is one of `"docker"`, `"vsphere"`, `"unknown"` (FR-006) — derived server-side from
`infrastructureRef.kind`; never absent.

## GET /api/clusters/infra/list?provider={docker|vsphere} (existing, extended)

- `provider` query parameter is **optional**. When omitted, the backend uses the first provider
  reported as `installed` by `/api/infra/capabilities` (FR-004, FR-012).
- When `provider=vsphere`, response shape is byte-for-byte unchanged from today (FR-010) —
  vSphere-specific fields (`server`, `thumbprint`, `modules`).
- When `provider=docker`, response is the new Docker-equivalent shape (FR-005):

```json
{
  "total": 2,
  "failing": 0,
  "clusters": [
    { "metadata": { "name": "workload-a" }, "cluster": "workload-a", "loadBalancerIP": "172.18.0.5", "status": { "ready": true } }
  ]
}
```

- If the requested (or auto-selected) provider is not installed, respond `404` with a clear error
  body — never a silently empty `200`.

## GET /api/machines/list, GET /api/machines/infra/list?provider={docker|vsphere} (existing, extended)

Same `provider` field addition on `GET /api/machines/list` items, and the same
`?provider=` dispatch/fallback behavior on the infra-list endpoint, mirroring the Clusters contract
above (FR-008).

## Error handling

- Unknown `?provider=` value (anything other than `docker`/`vsphere`) → `400` with a message listing
  the supported values.
- Any downstream list failure (e.g., CRD not actually queryable despite being reported installed) →
  `500` with an actionable message, per Constitution I ("Error states MUST propagate to the dashboard
  with actionable, human-readable descriptions").
