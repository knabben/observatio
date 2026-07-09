# Feature Specification: First-Class Pages for MachineHealthCheck, KubeadmControlPlane, MachineSet, and ClusterClass

**Feature Branch**: `007-capi-object-pages`
**Created**: 2026-07-08
**Status**: Draft
**Input**: User description: "Give `MachineHealthCheck`, `KubeadmControlPlane`, `MachineSet`, and `ClusterClass` their own first-class list + detail pages and per-object AI panel, mirroring the existing Clusters/Machines/MachineDeployments pattern... instead of leaving them as backend-only signals that only surface indirectly through rollups, severity, or debug-path evidence."

## Overview

Observātiō's dashboard today lets operators independently browse only three CAPI kinds — Clusters,
Machines, and MachineDeployments — each with its own live list page, a detail screen (status/spec,
conditions, full YAML), and a one-click "Ask AI about this" action. Four other operationally
important kinds have no such page:

- **MachineHealthCheck (MHC)** defines the remediation policy — timeouts, the `maxUnhealthy`
  threshold, which Machines it targets — behind every self-healing-vs-needs-investigation call the
  Day-2 Ops dashboard makes (006/US4). Today it's watched only as an internal signal feeding that
  classification; an operator can't see the policy itself anywhere.
- **KubeadmControlPlane (KCP)** carries control-plane replica health and etcd conditions — the most
  precise available signal for "is the management cluster's control plane actually healthy," and
  isn't watched or surfaced anywhere today.
- **MachineSet** sits between MachineDeployment and Machine and already drives the Day-2 Ops
  stalled-rollout check (006/US3), but exists only as an internal detail of that computation.
- **ClusterClass** has a working backend fetcher and a compact read-only widget embedded on the
  main dashboard, but no dedicated page, live updates, or detail view of its own.

This feature gives all four the same first-class treatment as Clusters/Machines/MachineDeployments:
a live list page, a detail screen showing their status/conditions and complete YAML, and a
per-object "Ask AI about this" action — so operators can inspect the underlying policy or state
directly, independent of (and to understand *why*) any conclusion the Day-2 Ops dashboard reaches.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Inspect MachineHealthCheck remediation policy (Priority: P1)

An operator wants to see exactly which Machines a MachineHealthCheck targets, its unhealthy-node
timeout, and its `maxUnhealthy` threshold — the policy behind every self-healing/needs-investigation
call the Day-2 Ops dashboard already makes — without having to reason about it only from the
dashboard's summarized conclusion.

**Why this priority**: MHC policy is the single most directly-requested gap — it's the concrete
mechanism behind an existing, already-shipped dashboard feature (006/US4), so understanding it
first unblocks operators from having to trust the dashboard's classification blindly.

**Independent Test**: Can be fully tested by navigating to the new MachineHealthCheck page,
confirming every MachineHealthCheck in the cluster is listed live, and opening one to see its
target selector, timeouts, `maxUnhealthy`, current remediation status, and full YAML.

**Acceptance Scenarios**:

1. **Given** one or more MachineHealthChecks exist in the cluster, **When** the operator opens the
   new MachineHealthCheck page, **Then** they see a live list of all MachineHealthChecks, matching
   the same list behavior (search, live updates, empty/error states) as the existing Machines page.
2. **Given** an operator selects a MachineHealthCheck from the list, **When** the detail screen
   opens, **Then** they see its target selector, `maxUnhealthy` threshold, unhealthy-condition
   timeouts, and current remediation status (e.g. current healthy count vs. expected) on the
   Specification tab, and its complete raw object on a YAML tab.
3. **Given** an operator is viewing a MachineHealthCheck's detail screen, **When** they use the
   "Ask AI about this" action, **Then** the AI panel opens pre-loaded with that MachineHealthCheck's
   identity and key fields, the same as it does today for a Machine or Cluster.

---

### User Story 2 - Inspect KubeadmControlPlane / etcd health (Priority: P2)

An operator wants to see control-plane replica counts, readiness, and etcd cluster health directly,
rather than only inferring management-cluster health from the Day-2 Ops dashboard's general
"API server unreachable"-style banner (006/FR-015).

**Why this priority**: KCP is the most severe blind spot today — it isn't watched at all, so there
is currently no way to see etcd/control-plane state independent of the coarse "management cluster
degraded" signal.

**Independent Test**: Can be fully tested by navigating to the new KubeadmControlPlane page and
confirming a KCP object's replica counts, readiness, and etcd-related conditions are visible.

**Acceptance Scenarios**:

1. **Given** a cluster using KubeadmControlPlane exists, **When** the operator opens the new
   KubeadmControlPlane page, **Then** they see a live list of KubeadmControlPlane objects.
2. **Given** an operator opens a KubeadmControlPlane's detail screen, **Then** they see its desired
   vs. ready replica counts and its status conditions (including etcd-related conditions when
   present) on the Specification tab, and its complete raw object on a YAML tab.
3. **Given** a cluster does not use KubeadmControlPlane (e.g. a different control-plane provider),
   **When** the operator opens the KubeadmControlPlane page, **Then** they see a clear empty state,
   not an error.

---

### User Story 3 - Inspect MachineSet rollout state (Priority: P3)

An operator wants to see MachineSet replica counts and age directly — the same data already used
internally to detect stalled rollouts (006/US3) — to understand a MachineDeployment rollout's
progress in more detail than the parent MachineDeployment view shows.

**Why this priority**: Valuable but already partially visible indirectly through the Day-2 Ops
stalled-rollout warning; a dedicated page adds direct visibility without being the primary gap.

**Independent Test**: Can be fully tested by navigating to the new MachineSet page and confirming
MachineSets are listed live with accurate replica counts, and that a stalled MachineSet's Machine
count matches what the Day-2 Ops dashboard already reports for it.

**Acceptance Scenarios**:

1. **Given** one or more MachineSets exist, **When** the operator opens the new MachineSet page,
   **Then** they see a live list showing each MachineSet's replica counts and owning
   MachineDeployment.
2. **Given** an operator opens a MachineSet's detail screen, **Then** they see its full status
   (replicas, ready replicas, available replicas, conditions) and complete raw object.

---

### User Story 4 - Browse ClusterClass as a first-class page (Priority: P4)

An operator wants to open a dedicated ClusterClass page with live updates and a detail view (status,
conditions, full YAML, Ask AI), the same as any other CAPI kind, instead of only seeing a compact
read-only table embedded on the main dashboard.

**Why this priority**: ClusterClass already has a working backend and a visible (if minimal) surface
today, so it's the smallest gap of the four — lowest priority, but still a real, requested gap for
consistency across all first-class kinds.

**Independent Test**: Can be fully tested by navigating to the new ClusterClass page and confirming
live list and detail behavior consistent with the other three new pages.

**Acceptance Scenarios**:

1. **Given** one or more ClusterClasses exist, **When** the operator opens the new ClusterClass
   page, **Then** they see a live list of ClusterClasses.
2. **Given** an operator opens a ClusterClass's detail screen, **Then** they see its status/reference
   fields and complete raw object, and can use "Ask AI about this."
3. **Given** the existing compact ClusterClass widget on the main dashboard, **When** this feature
   ships, **Then** the widget remains on the main dashboard as a rollup (consistent with Clusters/
   Machines/MachineDeployments, which also keep a summary on the main dashboard alongside their own
   dedicated pages) and is not removed.

---

### Edge Cases

- What happens when a MachineHealthCheck, KubeadmControlPlane, MachineSet, or ClusterClass list is
  empty (none exist, or the CRD/feature isn't in use in this environment)? Each page MUST show a
  clear "none found" empty state, not an error or a blank screen.
- What happens when the KubeadmControlPlane CRD isn't installed at all (a non-kubeadm control-plane
  provider is in use)? The page MUST degrade to an empty/unavailable state rather than failing the
  whole dashboard (matching 006's fix for provider-specific CRDs that may not exist).
- What happens when a MachineHealthCheck's target selector matches zero Machines? The detail screen
  MUST still show the policy itself (timeouts, threshold) rather than only an empty target list.
- What happens when a KubeadmControlPlane exists but has not yet reported etcd-related conditions
  (e.g. very early in provisioning)? The Specification tab MUST show those fields as
  not-yet-available rather than blank or misleading.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST provide a live list page for MachineHealthCheck objects, consistent
  in behavior (live updates, search, empty/error states) with the existing Machines/Clusters/
  MachineDeployments list pages.
- **FR-002**: The system MUST provide a live list page for KubeadmControlPlane objects, with the
  same list behavior as FR-001.
- **FR-003**: The system MUST provide a live list page for MachineSet objects, with the same list
  behavior as FR-001.
- **FR-004**: The system MUST provide a live list page for ClusterClass objects, with the same list
  behavior as FR-001, in addition to (not replacing) the existing compact ClusterClass widget on the
  main dashboard.
- **FR-005**: Each of the four new pages MUST provide a detail screen for a selected object showing
  its status/spec fields on a "Specification" tab and its complete raw object on a "YAML" tab,
  consistent with the existing detail-screen pattern (005/006).
- **FR-006**: The MachineHealthCheck detail screen MUST show its target selector, unhealthy-node
  timeout(s), `maxUnhealthy` threshold, and current remediation status (expected vs. currently
  healthy Machine counts).
- **FR-007**: The KubeadmControlPlane detail screen MUST show desired vs. ready replica counts and
  its status conditions, including etcd-related conditions when the object reports them.
- **FR-008**: The MachineSet detail screen MUST show replica, ready-replica, and available-replica
  counts, its owning MachineDeployment, and its status conditions.
- **FR-009**: The ClusterClass detail screen MUST show its status/reference fields already available
  from the existing ClusterClass backend data.
- **FR-010**: Each of the four new detail screens MUST provide a one-click "Ask AI about this"
  action that opens the existing global AI panel pre-loaded with that object's identity and key
  fields, the same as the existing Cluster/Machine/MachineDeployment detail screens.
- **FR-011**: Each of the four new list/detail pages MUST be reachable from the dashboard's lateral
  navigation, alongside the existing Clusters/Machine Deployments/Machines/Logs entries.
- **FR-012**: When a targeted kind's CRD is not installed or no objects of that kind exist, its page
  MUST show an explicit empty/unavailable state rather than an error or a blank screen (FR-017-style
  precedent from 006).

### Key Entities

- **MachineHealthCheck**: A first-class CAPI remediation policy object — target selector,
  unhealthy-node timeouts, `maxUnhealthy` threshold, and live remediation status.
- **KubeadmControlPlane**: A first-class CAPI control-plane object — desired/ready replica counts
  and status conditions, including etcd health when reported.
- **MachineSet**: A first-class CAPI object between MachineDeployment and Machine — replica counts
  and status conditions for one specific rollout generation.
- **ClusterClass**: An existing, already-modeled CAPI object — promoted from a dashboard-only widget
  to a first-class page with live updates and a detail view.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: An operator can find and open any MachineHealthCheck's remediation policy (target
  selector, timeouts, `maxUnhealthy`) in under 15 seconds from the main dashboard, without reading
  raw YAML or asking the AI panel.
- **SC-002**: An operator can determine a KubeadmControlPlane's replica health and etcd condition
  status directly, without inferring it from the Day-2 Ops dashboard's general management-cluster
  banner.
- **SC-003**: All four new list pages show live updates within the same latency budget as existing
  list pages (Constitution Principle II: under 2 seconds from a server-side change to UI update).
- **SC-004**: 100% of the four new pages show an explicit empty/unavailable state (never a blank
  screen or unhandled error) when their underlying kind has zero objects or an uninstalled CRD.
- **SC-005**: Every one of the four new detail screens supports the same one-click "Ask AI about
  this" action already available on Cluster/Machine/MachineDeployment detail screens, with no
  reduced functionality.

## Assumptions

- The four new pages follow the exact same list+detail+YAML+AI-panel pattern already established by
  Clusters/Machines/MachineDeployments (005/006) — this feature does not introduce new UI patterns,
  only extends the existing one to four more kinds, per the feature description's own framing.
- MachineHealthCheck and MachineSet are already watched internally (feeding the Day-2 Ops
  dashboard's 006 severity/risk classification); this feature adds first-class, independently
  browsable pages for them but does not change how 006 already consumes them internally.
- KubeadmControlPlane is watched for the first time by this feature. Read-only; no lifecycle actions
  (scaling, upgrades) are introduced — consistent with the dashboard's observability-only scope.
- The existing compact ClusterClass widget on the main dashboard is retained as-is; this feature
  only adds a new dedicated page alongside it, matching how Clusters/Machines/MachineDeployments
  already have both a main-dashboard rollup and their own dedicated pages.
- No new AI-panel capability is introduced; the four new detail screens reuse the existing global AI
  panel and per-object context mechanism from feature 005 unchanged.
