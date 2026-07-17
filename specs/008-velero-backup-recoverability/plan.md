# Implementation Plan: Velero Backup Recoverability Awareness

**Branch**: `008-velero-backup-recoverability` | **Date**: 2026-07-09 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/008-velero-backup-recoverability/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

Give the Day-2 Ops dashboard (006) visibility into Velero: watch `Backup`, `Schedule`, `Restore`,
and `BackupStorageLocation` objects the same way every other CAPI kind is watched (dynamic-client
GVR watch, decoded without a new typed dependency), extend the existing Day-2 Ops aggregator
(`WatchDay2Ops`/`assembleData`) to compute a Backup Health summary and per-cluster backup coverage,
and cross-reference that coverage with 006's existing CA-secret-missing severity so the dashboard
distinguishes "CA lost but a recent backup covers this cluster" from "CA lost, unrecoverable."
Reconciliation-pause visibility (`Cluster.spec.paused`) turns out to already be fully implemented
on the existing Cluster detail page (research.md R6) — no new work needed there. Whether to add an
active pause/unpause control is explicitly deferred (spec.md Assumptions); this plan is read-only.

## Technical Context

**Language/Version**: Go 1.23+ (backend, matches constitution), TypeScript 5 / React 19 (frontend)
**Primary Dependencies**: `k8s.io/client-go` dynamic client (existing) — no new Velero Go module
dependency (research.md R1); Mantine UI 7 (existing, frontend)
**Storage**: N/A — read-only live view over the Kubernetes API, no persistence added
**Testing**: `go test` with CAPI/dynamic fake clients (backend), Jest + Testing Library (frontend)
**Target Platform**: Existing Observātiō web app (Next.js static export embedded in the Go binary)
**Project Type**: Web application (existing `webserver/` + `front/` structure)
**Performance Goals**: Consistent with constitution Principle II — dashboard update latency from
server event to UI stays below 2 seconds under normal load; no new polling introduced
**Constraints**: Read-only — no mutating calls against Velero, CAPI, or cluster resources (spec.md
FR-013); must degrade gracefully (explicit "not available" state, not an error) when Velero is not
installed, matching the corrected pattern from 006's live vSphere-CRD-missing bug fix
**Scale/Scope**: One new aggregated data section (Backup Health) on the existing Day-2 Ops landing
page, plus an enrichment of one existing severity computation (CA-secret-missing); no new frontend
routes

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **I. Observability & Data Consolidation** — PASS. Backup/Restore/Schedule/BackupStorageLocation
  data is correlated to the owning Cluster (coverage matching, research.md R5) rather than shown as
  a disconnected Velero object list, consolidating it into the same view as Cluster/Machine health.
- **II. Real-Time Visibility** — PASS. Extends the existing WebSocket-driven `WatchDay2Ops`
  aggregator; no new polling or REST endpoints for live data.
- **III. ClusterAPI Resource Model Compliance** — PASS. Velero objects are not CAPI resources and
  are not modeled as such; they are treated as an external signal correlated *to* the existing
  Cluster model (via namespace/label matching), consistent with treating non-CAPI infrastructure
  details as data feeding the core domain, not becoming first-class domain types themselves.
- **IV. AI-Augmented Troubleshooting** — N/A for this feature. No new AI panel surface is added;
  existing AI context building for Cluster objects is unaffected. (If a future iteration wants the
  AI panel to reason about recoverability, the enriched severity `Reason` text is already
  structured plain language it could consume unchanged.)
- **V. Test-Driven Quality** — PASS. New pure classification functions (coverage matching, backup
  health rollup) follow the existing `day2ops` package's tested-pure-function convention
  (`severity.go`, `risk_certexpiry.go`); watcher/store changes get fake-dynamic-client tests
  matching `day2ops_test.go`'s existing pattern.

**Result**: No violations. Complexity Tracking table not needed.

## Project Structure

### Documentation (this feature)

```text
specs/008-velero-backup-recoverability/
├── plan.md              # This file (/speckit-plan command output)
├── research.md          # Phase 0 output (/speckit-plan command)
├── data-model.md         # Phase 1 output (/speckit-plan command)
├── quickstart.md         # Phase 1 output (/speckit-plan command)
├── contracts/            # Phase 1 output (/speckit-plan command)
└── tasks.md              # Phase 2 output (/speckit-tasks command - NOT created by /speckit-plan)
```

### Source Code (repository root)

```text
webserver/
├── internal/infra/clusterapi/day2ops/
│   ├── types.go                    # MODIFY: BackupHealth, ClusterBackupCoverage,
│   │                                #   RecoveryInfo, CategoryBackup, RiskBackupCoverage types
│   ├── severity.go                 # MODIFY: ComputeCASecretMissingSeverity gains recovery info
│   ├── backuphealth.go             # NEW: ComputeBackupHealth, ComputeClusterBackupCoverage
│   ├── backuphealth_test.go        # NEW
│   └── severity_test.go            # MODIFY: recovery-info cases
├── internal/web/watchers/
│   └── day2ops.go                  # MODIFY: add velero GVRs (gated on CRD presence), store
│                                   #   maps for Backup/Restore/Schedule/BackupStorageLocation,
│                                   #   wire ComputeBackupHealth + coverage into assembleData
└── internal/web/watchers/day2ops_test.go  # MODIFY: velero fan-in cases

front/
├── app/ui/dashboard/shared/use-day2-ops.ts       # MODIFY: BackupHealth types, `backup` category
├── app/ui/dashboard/components/ops/
│   ├── backup-health-card.tsx                    # NEW: dedicated card (data shape doesn't
│   │                                              #   fit the existing 3-bucket rollup card)
│   ├── backup-health-card.test.tsx               # NEW
│   ├── ops-dashboard.tsx                         # MODIFY: render BackupHealthCard
│   └── severity-banner.tsx                       # MODIFY (if needed): render recoverable/
│                                                  #   unrecoverable distinctly, not just prose
```

No changes needed to `webserver/internal/web/handlers/system/websocket.go`, `resource-gvr.ts`, or
`nav-links.tsx` — this feature extends the existing `day2ops` WS type, it does not add a new one
(research.md R2). No new frontend routes are added.

**Structure Decision**: Extend the existing Day-2 Ops vertical slice (backend `day2ops` domain
package + `watchers/day2ops.go` aggregator + frontend `ops/` components) rather than introducing a
parallel Velero-specific pipeline, since every new signal here exists specifically to enrich that
dashboard's rollups and severities (research.md R2) — mirroring how MachineSet and
MachineHealthCheck were added internally to `day2ops.go` in 006 rather than as standalone watchers.

## Complexity Tracking

*No violations — table not needed.*
