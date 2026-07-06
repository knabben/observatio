# Contract: `GET /api/day2ops/detail`

Scoped, on-demand REST exception (research.md R9), mirroring the pattern already established by
`GET /api/raw` in feature 005: WS delivers the live rollup/summary; this endpoint hydrates the full
evidence list for one specific object only when an operator expands it on the dashboard, avoiding
recomputing/re-pushing large per-object payloads to every connected client on every change.

## Request

```
GET /api/day2ops/detail?group=&version=&resource=&namespace=&name=
```

Same GVR + namespace/name query shape as `/api/raw` for consistency.

## Response `200 OK`

```json
{
  "objectRef": {"group": "cluster.x-k8s.io", "version": "v1beta1", "resource": "machines", "namespace": "default", "name": "worker-0"},
  "path": {
    "layers": [
      {"layer": "conditions", "status": "implicated", "evidence": ["Ready=False: WaitingForInfrastructure"], "source": "Machine/worker-0"},
      {"layer": "phase", "status": "implicated", "evidence": ["Phase=Provisioning"], "source": "Machine/worker-0"},
      {"layer": "provider_resource", "status": "implicated", "evidence": ["Ready=False: VM creation failed - insufficient resources"], "source": "DockerMachine/worker-0"},
      {"layer": "controller_activity", "status": "inconclusive", "evidence": [], "source": ""}
    ],
    "summary": "Waiting on infrastructure provisioning (DockerMachine: VM creation failed)"
  }
}
```

`controller_activity` is only populated (non-empty `evidence`) when the three earlier layers were
all `inconclusive` (FR-007); here it stays empty because `provider_resource` already explains the
failure.

## Errors

- `400 Bad Request` — missing/invalid GVR or name, same validation shape as `/api/raw`.
- `404 Not Found` — object does not exist.
- `500 Internal Server Error` — underlying API-server call failed for a reason other than
  not-found (distinct from the WS-level `SourceUnavailable` flag, which reflects the aggregator's
  own persistent-connection health).
