# Gap analysis: Observātiō vs. the CAPI companion guide ("When Clusters Break")

**Source guide**: `knabben/velero-mcp` — `docs/companion-guide.md`
**Compared against**: this repo as of `006-day2-ops-dashboard` (spec.md/plan.md/tasks.md) plus the
existing `005-object-viz-ai-panel` AI troubleshooting feature.
**Date**: 2026-07-06

## Summary

Feature 006 already closes most of the companion guide's "Day-2 Operations and Debugging" and "Failure
Modes" sections: the layered debugging path (Layers 1-4: conditions → Machine phase → provider resource
→ controller activity), certificate-expiry, stalled-rollout, version-skew, and drift detection, MHC
self-healing vs. escalation, provider-controller crash detection, a management-cluster
API-server-unreachable banner, and a CA-secret-missing top-severity warning. What's *not* covered is the
guide's entire **Recovery Strategy with Velero** section (Section 5) — Observātiō has zero backup/DR
capability today — plus a handful of adjacent gaps in metrics, AI tooling, and object coverage. Those are
catalogued below; see `docs/proposals/` for the corresponding next-step specification statements.

## Objects/CRDs missing from Observātiō's data model

- **Velero CRDs, entirely**: `Backup`, `Schedule`, `Restore`, `BackupStorageLocation`,
  `VolumeSnapshotLocation`, `PodVolumeBackup`/`PodVolumeRestore`, `DeleteBackupRequest`.
- **The five CAPI secrets aren't individually modeled.** Spec 006 only names "the cluster CA secret"
  (FR-016); the etcd CA, front-proxy CA, service-account signing key, and admin kubeconfig aren't
  tracked as distinct entities with their own existence/health checks, even though the guide treats all
  five as independently critical.
- **`Cluster.spec.paused` isn't surfaced as first-class state anywhere.** It's needed for the guide's
  recovery steps 2.5 and 4 (pause CAPI reconciliation before restoring, unpause after).
- **No `KubeadmControlPlane` (KCP) watcher.** Plan.md's watched kinds are Cluster, ClusterClass,
  Machine, MachineDeployment, MachineSet, MachineHealthCheck — not KCP — so nothing attributes Level-3
  management-cluster degradation to etcd quorum specifically; today it's only generic
  API-server-unreachable (FR-015).

## Functionalities missing

- **No Prometheus metrics ingestion.** The guide dedicates a full subsection to controller `/metrics`,
  the `/debug/flags/v` dynamic log-level endpoint, and pprof profiling. Observātiō has none of this —
  controller-level degradation (Level 2) is only inferred from Pod/Deployment status, not reconciliation
  latency, queue depth, or error-rate trends.
- **No backup/DR workflow of any kind.** No "is there a recent backup," no RPO/RTO tracking, no guided
  Detect → Bootstrap → Pause → Restore → Unpause → Adopt runbook, no pause/unpause action.
- **No outbound alerting** (Slack/email/webhook). The guide recommends Grafana alerts on reconciliation
  failures and queue saturation; Observātiō's severity banners (006/US4) are push-to-UI only.
- **No historical/trend state.** The product is intentionally stateless/live-derived (plan.md:
  "Storage: N/A"), so there's no incident timeline or trend view of past severity events.
- **No DR-readiness testing.** The guide's "an untested backup is not a backup" principle has no
  counterpart — nothing schedules or verifies periodic recovery-test runs.

## Integrations missing

- **Velero** — the guide's entire recovery strategy; zero references anywhere in this codebase
  (`webserver/`, `front/`).
- **MCP, in either direction.** The AI agent's only tool is a hand-rolled `kubectl` exec
  (`webserver/internal/infra/llm/tools.go`). Observātiō is neither an MCP client (consuming the
  referenced `velero-mcp` server) nor an MCP server (exposing its own aggregated Day-2/backup state to
  other agents).
- **No S3-compatible object-storage awareness** (MinIO/S3/GCS) — irrelevant today since there's no
  backup feature yet.

## AI gaps

- The existing "Ask AI about this" panel (feature 005) is **object-scoped only** — there's no
  cluster-wide/global assistant a user can ask "are we protected? when did we last back up?" without an
  object already in view.
- Exactly one tool exists (`RunKubectl`, `tools.go:36`), and it execs whatever command string the model
  produces with no read/write guard — worth flagging as an existing risk independent of any Velero work,
  since the rest of the product's constitution is read-only by design.
- **No Velero-aware tool** — can't list/describe backups or restores, can't check
  `BackupStorageLocation` health, can't correlate "CA secret missing" (006/severity.go) with "last
  backup age" to tell the operator whether recovery is actually possible right now.
- **No multi-agent/per-domain specialization** the way kagent has (separate k8s/Istio/observability/Argo
  agents, each with their own MCP tool servers, CRD-based agent registry, chat UI). Observātiō has a
  single fixed agent and one tool; `agent.go`'s `initializeAgents()` anticipated more agents but only
  ever registers one (`cluster-agent`).

## Services missing

- **No backup-orchestration service** (trigger/schedule backups, poll Velero objects, surface
  `BackupStorageLocation` reachability).
- **No recovery-runbook service** encoding the guide's 5-step sequence as a guided, semi-automated flow
  — doing this for real requires a deliberate, explicit expansion of the product's current read-only
  constitution, since pausing/unpausing a Cluster and triggering a Velero restore are both writes.
- **No metrics-ingestion/query service** to back a real controller-health panel beyond Pod-Ready
  heuristics.

## Proposals

See the individual quick-specification statements in `docs/proposals/`, one per proposal, written so
each can be handed to `/speckit-specify` directly if/when it's formalized into a numbered feature:

1. `01-velero-integration.md` — Velero backup/restore state surfaced on the Day-2 Ops dashboard.
2. `02-mcp-server-integration.md` — Observātiō as an MCP client and/or server.
3. `03-velero-management-dashboard.md` — a dedicated Backups/Schedules/Restores management view.
4. `04-controller-metrics.md` — Prometheus-backed controller health.
5. `05-etcd-control-plane-health.md` — KubeadmControlPlane/etcd-quorum-aware Level-3 detection.
6. `06-first-class-capi-objects.md` — first-class list/detail pages for MHC, KCP, MachineSet, and
   ClusterClass, which today only exist as internal signals or aren't watched at all.
