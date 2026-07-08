# Quickstart: Day-2 Operations Dashboard

Manual verification scenarios, one per user story, against a `kind-capi-mgmt`-style management
cluster with the Docker infrastructure provider (CAPD) enabled.

## US1 — Centralized landing view

1. Ensure a mix of healthy and unhealthy objects exists (e.g. one healthy Cluster, one Machine
   stuck `Provisioning`).
2. Open the dashboard (`front/app/dashboard/page.tsx`).
3. Confirm Cluster, MachineDeployment, and Machine category rollups are all visible without
   navigating to `/dashboard/clusters`, `/dashboard/machines`, etc.
4. Stop all objects' issues (or point at an all-healthy environment) and confirm each category
   shows a clear "all clear" state, not an empty/blank card.
5. Click a category to narrow the view in place; confirm the URL/page does not fully navigate away.

## US2 — Layered debugging path

1. Seed a Machine stuck in `Provisioning` (e.g. patch its `DockerMachine` to a failing image, or
   let CAPD naturally stall by requesting more replicas than local Docker resources allow).
2. Open that Machine's rollup entry on the dashboard; confirm the path shows
   `conditions → phase → provider_resource` implicated, with the DockerMachine's condition/message
   inline (FR-004, FR-006).
3. Seed a Machine that is `Provisioned` but never reaches `Running` (e.g. break the bootstrap
   token or kubeadm join config); confirm the path is labeled at the bootstrap/phase layer, and is
   visually distinguishable from the infrastructure-layer case in step 2.
4. Seed a case where conditions/phase/provider-resource are all inconclusive (hard to construct
   naturally — approximate by injecting a Kubernetes Event with an error reason against a Machine
   that otherwise looks healthy) and confirm `controller_activity` becomes populated only in this
   case, not in steps 2–3.
5. Confirm every layer shown is labeled with its position in the sequence.

## US3 — Proactive risk detection

1. **Cert expiry**: patch a test cluster's `<cluster>-ca` Secret's certificate (or use a
   short-lived test CA) to expire within the 30-day default window; confirm a `cert_expiry`
   warning appears with the correct date.
2. **Stalled rollout**: trigger a MachineDeployment rollout, then apply a PodDisruptionBudget on
   the workload cluster that blocks eviction of the old MachineSet's node; confirm a
   `stalled_rollout` warning appears once the grace period elapses, naming the blocking PDB.
3. **Version skew**: (best-effort per research.md R6) install a provider version whose CRDs are
   intentionally older/newer than the controller expects; confirm a `version_skew` warning appears,
   and confirm the UI does not overstate certainty (per R6, this is heuristic).
4. **Drift**: manually edit a CAPD-managed Docker container's config out-of-band (or simulate via a
   stale `observedGeneration`); confirm a `drift` warning appears against the affected object.
5. For each of the four, verify a `not_evaluable` / "check could not be performed" state is shown
   when the underlying data can't be read (e.g., insufficient RBAC on the Secret) rather than the
   warning being silently omitted (FR-018).

## US4 — Failure-severity awareness

1. **Self-healing**: with a `MachineHealthCheck` configured, kill a worker node's kubelet; confirm
   the dashboard shows the resulting remediation as informational, not a red alert.
2. **maxUnhealthy breach**: fail enough nodes simultaneously to breach the configured
   `maxUnhealthy` threshold; confirm the dashboard escalates to "needs investigation" instead of
   showing continued self-healing.
3. **Provider controller crash-loop**: force-crash the CAPD controller pod repeatedly (e.g. inject
   a bad flag) and confirm a provider-degraded flag appears, distinct from any single object's
   status.
4. **Management cluster degraded**: temporarily block the backend's network path to the management
   cluster's API server; confirm a top-level, hard-to-miss banner appears stating lifecycle
   operations are blocked, and that it disappears once connectivity is restored.
5. **CA secret loss**: delete (in a disposable test environment only) a cluster's `<cluster>-ca`
   Secret; confirm the highest-severity warning appears against that cluster, and that the
   workload cluster's own continued healthy-looking status does not suppress it (Edge Cases).

## US5 — Deep-dive into controller logs

1. Seed the "controller_activity implicated" case from US2 step 4 (inject an error Event against an
   otherwise-inconclusive object). From that object's debugging path, choose the deep-dive action;
   confirm the Logs view opens scoped to the correct controller (CAPI core or provider) and streams
   its Pod's actual log output — the same content `kubectl logs -n <ns> deploy/<name>` would show.
2. Navigate to the Logs view directly from the lateral navigation (without drilling into an object
   first); confirm a controller can be chosen manually.
3. Scale the relevant controller Deployment to 0 replicas (or otherwise ensure no Pod backs it)
   and confirm the Logs view shows an explicit "logs unavailable" state, not a blank pane.
4. From a VM-based-provider Machine's debugging path, open the node-access deep-dive; confirm it
   shows only the SSH command + node address (no live output, no terminal), and confirm no
   credential material appears anywhere in the request/response.

## Cross-cutting checks

- With the management cluster's API server unreachable at startup, confirm `SourceUnavailable`
  drives an explicit "data unavailable" banner (FR-017) rather than a false "all healthy" or blank
  dashboard.
- Confirm rollup/severity updates arrive over the existing WebSocket connection (check browser
  devtools' WS frames) rather than via polling REST calls, and that only the on-demand detail
  drill-in (`GET /api/day2ops/detail`) appears as a REST call, matching the scoped exception in
  research.md R9.
