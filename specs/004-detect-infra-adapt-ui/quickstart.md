# Quickstart: Verifying Infrastructure Provider Detection

Manual verification steps once the feature is implemented, covering each acceptance scenario in
`spec.md`.

**Live validation (2026-07-06)**: Scenario 1 (Docker-only) was run end-to-end against the real
`kind-capi-mgmt` cluster (Docker/CAPD-backed) using a standalone build on a scratch port:
- `GET /api/infra/capabilities` → `{"docker":{"installed":true,"version":"v1.10.10"},"vsphere":{"installed":false,"version":""}}`.
- `GET /api/clusters/list` / `GET /api/machines/list` → each item carries `"provider":"docker"`.
- `GET /api/clusters/infra/list` (no `?provider=`) auto-selected Docker and returned the real
  `DockerCluster`; `GET /api/clusters/infra/list?provider=vsphere` correctly returned `404 vsphere
  infrastructure provider is not installed`.
- WebSocket `{"type":"cluster-infra-docker"}` and `{"type":"machine-infra-docker"}` both streamed
  the real Docker objects immediately (`ADDED` events) on connect.
Scenarios 2/3/4/5 (vSphere-only, mixed, unknown provider, no provider) are covered by the fake-client
and component test suites — no vSphere or multi-provider live cluster was available in this
environment to exercise those end-to-end.

## 1. Docker-only environment

1. Point the backend at a management cluster with only the Docker (CAPD) infrastructure provider
   installed and one or more Docker-backed clusters.
2. `curl localhost:8080/api/infra/capabilities` → expect `{"docker":{"installed":true,"version":"..."},"vsphere":{"installed":false,"version":""}}`.
3. Open the Clusters screen → expect a Docker infrastructure tab (not "vSphere Clusters"), and each
   row's provider badge reads "Docker vX.Y.Z".

## 2. vSphere-only environment

1. Point the backend at a management cluster with only vSphere installed.
2. `curl .../api/infra/capabilities` → `vsphere.installed=true`, `docker.installed=false`.
3. Open the Clusters screen → the existing "vSphere Clusters" tab and table render exactly as
   before (server, thumbprint, modules unchanged) — regression check for FR-010.

## 3. Mixed environment

1. Point the backend at a management cluster with both providers installed and at least one cluster
   of each kind.
2. Open the Clusters screen → both provider tabs are present; each shows only the clusters
   belonging to that provider; the main list's provider badges are correct per row.

## 4. Unknown / unsupported provider

1. Create (or fixture) a Cluster whose `infrastructureRef.kind` is neither `DockerCluster` nor
   `VSphereCluster`.
2. Confirm it still appears in the main Clusters list with an "Unknown" provider badge, and no
   provider-specific tab is forced open for it.

## 5. No supported provider installed

1. Point the backend at a management cluster with core Cluster API only (no infra provider CRDs).
2. `curl .../api/infra/capabilities` → both `installed` fields `false`.
3. Open the Clusters screen → a clear "no supported infrastructure provider detected" message is
   shown instead of an empty tab.

## Automated coverage

- Backend: `make run-tests-backend` — fake-client tests for `providerkind.FromKind`, the
  capability-filtering logic, and the Docker fetcher/handler.
- Frontend: `make run-tests-frontend` — Jest tests for dynamic tab rendering across all 5 scenarios
  above, provider badge rendering (incl. unknown), and the empty-provider message.
