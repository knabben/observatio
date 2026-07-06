# Contract: Logs API

Backs the new "Logs" destination (User Story 5): controller Pod-log streaming plus static
node-access instructions. No Docker daemon access, no SSH credentials — see research.md R10.

## `GET /api/logs/controller`

### Request

```
GET /api/logs/controller?namespace=capi-system&deployment=capi-controller-manager&follow=true
```

`namespace`/`deployment` identify the controller (`ControllerRef` minus the resolved Pod name,
which the handler resolves server-side from the Deployment's current Pod). `follow=true` keeps the
response streaming (chunked transfer), matching `kubectl logs -f` semantics; omit for a bounded
snapshot.

### Response `200 OK`

Streamed (or snapshot) plain-text log lines, identical in content to what
`kubectl logs -n <namespace> deploy/<deployment>` would print — this handler is a thin wrapper
around the same Kubernetes Pod-log subresource, not a reformatted/filtered view.

### Errors

- `400 Bad Request` — missing `namespace`/`deployment`.
- `404 Not Found` — no Pod currently backs the given Deployment (e.g., mid-rollout, crash-looping
  with no ready replica) — this itself doubles as evidence for FR-014's controller-degraded signal.
- `503 Service Unavailable` — logs could not be retrieved (e.g., Pod has no retained log history);
  maps to the frontend's FR-023 "logs unavailable" state.

## `GET /api/logs/node-access`

### Request

```
GET /api/logs/node-access?group=cluster.x-k8s.io&version=v1beta1&resource=machines&namespace=default&name=worker-0
```

Same GVR + namespace/name shape as `/api/raw` and `/api/day2ops/detail`, scoped to a Machine.

### Response `200 OK`

```json
{
  "objectRef": {"group": "cluster.x-k8s.io", "version": "v1beta1", "resource": "machines", "namespace": "default", "name": "worker-0"},
  "command": "ssh capi@10.0.1.23",
  "note": "Observātiō does not store or manage SSH credentials. Run this command from your own machine."
}
```

Sourced entirely from `Machine.status.addresses` already available to the backend; never contacts
the node, never handles a key or password.

### Errors

- `400 Bad Request` — missing/invalid GVR or name.
- `404 Not Found` — Machine does not exist, or has no recorded address yet.
