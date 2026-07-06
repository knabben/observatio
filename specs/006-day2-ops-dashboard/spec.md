# Feature Specification: Day-2 Operations Dashboard

**Feature Branch**: `006-day2-ops-dashboard`
**Created**: 2026-07-06
**Status**: Draft
**Input**: User description: "The tools dashboard need to be the centralized place for 2days operations, this can be achieved to present the cluster mgmt + workloads on a layered format, separated by categories and layer debugging, a NEW FRONT dashboard must exist... [layered CAPI debugging model, Layer 1-5, from `clusterctl describe` down to node-level logs]"

## Overview

Today, understanding the health of a Cluster API environment in Observātiō means visiting separate
list pages for Clusters, Machines, and MachineDeployments and mentally combining what's seen there
with knowledge of how CAPI debugging normally works in practice: start at `clusterctl describe`
for object-level conditions, check Machine phase, drop into the provider-specific resource, then
controller logs, and — if none of that explains it — the node itself. That sequence lives in
operators' heads and terminals today; nothing in the dashboard reflects it.

This feature introduces a new, centralized Day-2 Operations dashboard that presents management-plane
and workload-plane objects together, organized by category, and — for any unhealthy object —
automatically walks the same layered path an experienced operator would run by hand, surfacing the
responsible layer and its evidence directly on the dashboard. It also proactively flags the
well-known slow-building CAPI risks (certificate expiry, stalled rollouts, provider/CRD version
skew, infrastructure drift) and distinguishes failures the system heals on its own from ones that
need escalating human attention, up to the most severe case: loss of the management cluster's CA
material, which no substitute CA can repair.

## Clarifications

### Session 2026-07-06

- Q: When a cluster is being debugged, several distinct issues can be present in a system at once
  (object-level, Machine-phase, provider-resource, controller-level). Must the operator still cross
  reference CLI tools to figure out which layer is at fault? → A: No — the tool MUST be able to
  identify the debugging path using only what's shown on the dashboard; this is the central
  requirement of the layered debugging user story (see US2, FR-004–FR-007), and the path shown must
  account for multiple simultaneously-contributing layers rather than stopping at the first one found
  (see Edge Cases).
- Q: Beyond the five debugging layers, are there additional known CAPI failure classes the dashboard
  should watch for proactively rather than only reactively? → A: Yes — four specific classes were
  called out: certificate expiry (CAPI-issued certs default to a 1-year lifetime and silently expire
  if rotation is missed), stalled MachineDeployment rollouts (new MachineSet healthy, old one blocked
  from scaling down by a PDB or finalizer), provider/CRD version skew (CRDs and controllers upgraded
  out of step), and infrastructure drift (out-of-band changes to provider-managed VMs/network/storage).
  These are captured as User Story 3.
- Q: Should identifying the failing layer require a drill-in click, or must the layers and overall
  status be visible on the first (landing) screen itself, with deeper investigation (logs, node
  access) reached only via further clicks? → A: The layers, status, and overall data MUST be
  presented directly on the first screen for every unhealthy object — no drill-in is required just
  to see *which* layer is implicated (US2's Acceptance Scenarios and FR-004/FR-005 already required
  this; this clarifies it applies to the landing screen itself, not a secondary page). Deeper
  artifacts — full evidence detail, and now also raw logs — are reached via a "deep dive" click from
  that same layer summary. Since Observātiō has no log-viewing capability today, one MUST be added,
  as a new item in the existing lateral navigation, reachable from the deep-dive action. Captured as
  a new User Story 5.
- Q: The requested deep-dive included "SSH" access to nodes. Observātiō has no SSH/terminal
  capability today, and CAPD (Docker) nodes and VM-based-provider nodes (vSphere) have very
  different access models — building a real in-browser SSH terminal means storing/managing SSH
  credentials or bastion access, a meaningfully larger security surface. What should "SSH" mean here?
  → A: Provider-aware logs, no live terminal. For VM-based providers, it shows SSH connection
  instructions (command + node address) for the operator to run themselves — Observātiō MUST NOT
  store or manage SSH credentials of any kind.
- Q: Which logs, specifically — Machine/node-level container output, or the CAPI/provider
  controllers' own reconciliation logs? → A: Controller logs — the same ones an operator would get
  from `kubectl logs -n capi-system deploy/capi-controller-manager` (and the equivalent
  provider-controller Deployment, e.g. `capd-controller-manager` in `capd-system`). This is Layer 4
  from the original request, retrieved via the standard Kubernetes Pod-log API (the same mechanism
  `kubectl logs` uses) — not a Docker-daemon-specific per-Machine-container log, and not a new
  external dependency. User Story 5 and FR-019–FR-023 are revised accordingly.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Centralized Day-2 landing view (Priority: P1)

An operator responsible for one or more Cluster API management clusters opens Observātiō and lands on a single dashboard that summarizes the health of both the management-plane objects (Clusters, ClusterClasses) and workload-plane objects (Machines, MachineDeployments, MachineSets) in one place, organized by category, instead of having to visit separate list pages to build a mental picture of overall health.

**Why this priority**: This is the entry point for every other scenario in this feature. Without a consolidated landing view, the layered debugging flow (US2) and proactive risk detection (US3) have nowhere to live.

**Independent Test**: Can be fully tested by opening the new dashboard route with a mix of healthy and unhealthy objects across categories and confirming every category (Clusters, Machines, MachineDeployments, and their infrastructure counterparts) is represented with an accurate rolled-up health count, without visiting any other page.

**Acceptance Scenarios**:

1. **Given** a management cluster with 3 Clusters, 10 Machines, and 2 MachineDeployments in various states, **When** the operator opens the Day-2 dashboard, **Then** they see all three categories summarized with counts of healthy/degraded/failed objects, without navigating away.
2. **Given** every object in the environment is healthy, **When** the operator opens the dashboard, **Then** the dashboard shows a clear "all clear" state per category rather than an empty or ambiguous screen.
3. **Given** an operator is viewing the dashboard, **When** they select a category (e.g., "Machines"), **Then** the dashboard narrows to show only that category's objects while remaining on the same screen (no full page navigation required).

---

### User Story 2 - Layered root-cause debugging path (Priority: P1)

When an object is unhealthy (a Machine stuck provisioning, a cluster that won't scale, an upgrade that stalls), the operator needs the dashboard itself to identify which layer of the system is responsible — from high-level object conditions, down through Machine phase, provider-specific infrastructure resource, and controller reconciliation activity — so they can go straight to the right layer instead of manually running `clusterctl describe`, `kubectl get machines`, `kubectl describe <provider-machine>`, and controller logs in sequence.

**Why this priority**: This is the core "Day-2 operations" value proposition explicitly requested: the tool must be able to identify the debugging path by looking only at the dashboard, removing the need for manual, sequential CLI investigation.

**Independent Test**: Can be fully tested by seeding a known failure at a specific layer (e.g., an infrastructure-provider error on a `DockerMachine`) and confirming the dashboard highlights that exact layer as the likely cause, with the conditions/evidence that led to the determination, using only the dashboard UI.

**Acceptance Scenarios**:

1. **Given** a Machine is stuck in `Provisioning` phase, **When** the operator opens that Machine's entry on the dashboard, **Then** the dashboard indicates the issue is at the "infrastructure provisioning" layer and surfaces the relevant condition (e.g., `WaitingForInfrastructure`) without requiring the operator to cross-reference the provider object manually.
2. **Given** a Machine is `Provisioned` but never reaches `Running`, **When** the operator inspects it on the dashboard, **Then** the dashboard indicates the issue is at the "bootstrap" layer (kubeadm/bootstrap provider), distinguishing it from an infrastructure-layer failure.
3. **Given** a provider-specific resource (e.g., `DockerMachine`, `VSphereMachine`) reports an error condition or event, **When** the operator views the parent Machine on the dashboard, **Then** that provider-level detail is shown inline as part of the same debugging path, not as a separate lookup.
4. **Given** the object-level layers (conditions, phase, provider resource) show no explanation for a stalled object, **When** the operator continues down the debugging path on the dashboard, **Then** the dashboard surfaces relevant controller-level reconciliation signals (e.g., recent error-level log activity for the responsible controller) tied to that object, so the operator is directed toward controller logs only when the higher layers were inconclusive.
5. **Given** an operator is debugging an object, **When** they view the suggested path, **Then** each step in the path is labeled with the layer it corresponds to (object conditions → Machine phase → provider resource → controller activity), so the operator always knows how deep into the stack they are looking.

---

### User Story 3 - Proactive risk detection (Priority: P2)

An operator wants the dashboard to flag known classes of slow-building problems before they cause an outage: certificates approaching expiry, MachineDeployment rollouts that have stalled mid-rollout, provider/CRD version mismatches, and infrastructure that has drifted from what CAPI expects — so these are caught during routine dashboard checks rather than discovered when something breaks.

**Why this priority**: These are well-understood, recurring CAPI failure classes that are otherwise invisible until they cause an incident. They are valuable but secondary to the reactive debugging flow in US2, since US2 covers the "something is already broken" case operators hit most often.

**Independent Test**: Can be fully tested by seeding each of the four risk conditions independently (a soon-to-expire certificate, a MachineDeployment with an old MachineSet stuck scaling down, a mismatched provider/CRD version pairing, and a manually modified provider VM) and confirming each is surfaced as a distinct, correctly categorized warning on the dashboard.

**Acceptance Scenarios**:

1. **Given** a cluster's certificates are within a configurable warning window of expiry, **When** the operator views the dashboard, **Then** a certificate-expiry warning is shown against that cluster with the expiry date.
2. **Given** a MachineDeployment rollout has a new MachineSet that is healthy but an old MachineSet that will not scale down, **When** the operator views that MachineDeployment on the dashboard, **Then** a "stalled rollout" warning is shown, along with the likely blocking cause (e.g., pod disruption budget or finalizer) if determinable from cluster state.
3. **Given** an installed provider's version and its CRDs are out of sync with what the running controller expects, **When** the operator views the dashboard, **Then** a version-skew warning is shown identifying the affected provider.
4. **Given** a provider-managed infrastructure resource's observed state no longer matches its CAPI spec (drift), **When** the operator views the corresponding object on the dashboard, **Then** a drift warning is shown against that object.

---

### User Story 4 - Failure-severity awareness (Priority: P3)

An operator wants the dashboard to distinguish between failures that the system will resolve on its own (e.g., a single worker node failing over via `MachineHealthCheck`) and failures that require human intervention with escalating urgency (a provider controller crash-looping, the management cluster itself degraded, or the management cluster's cluster secrets/CA lost) — so they can prioritize attention correctly instead of treating every red indicator the same way.

**Why this priority**: This shapes how alarming the dashboard should look for a given condition and prevents alert fatigue, but it is a refinement on top of US2/US3's detection — it does not introduce new detection capability, only better-communicated urgency.

**Independent Test**: Can be fully tested by simulating one condition from each severity level (a self-healing node failure, a crash-looping provider controller, a degraded management-cluster API server, and a missing cluster CA secret) and confirming the dashboard labels each with a distinct, correctly escalating severity and guidance.

**Acceptance Scenarios**:

1. **Given** a worker node fails and `MachineHealthCheck` is actively remediating it (Machine deleted, replacement provisioning), **When** the operator views the dashboard, **Then** the event is shown as informational/self-healing, not as a critical alert requiring action.
2. **Given** `MachineHealthCheck`'s `maxUnhealthy` threshold is breached (remediation paused because too many nodes are unhealthy at once), **When** the operator views the dashboard, **Then** this is escalated to a distinct "needs investigation" state, since it may indicate a network partition rather than independent node failures.
3. **Given** a provider controller is crash-looping, **When** the operator views the dashboard, **Then** the affected provider is flagged as degraded at the controller level, distinct from individual object-level failures.
4. **Given** the management cluster's API server is degraded (e.g., unreachable or returning errors consistent with lost etcd quorum), **When** the operator views the dashboard, **Then** the dashboard shows a top-level, hard-to-miss banner indicating that all lifecycle operations (scaling, upgrades, certificate rotation) are currently blocked, separate from individual object statuses.
5. **Given** the dashboard can determine that a cluster's CA secret is missing or inaccessible, **When** the operator views that cluster, **Then** the dashboard shows the highest-severity warning available, explaining that no new certificates can be issued or rotated for that cluster's nodes until the original CA is restored, and that a substitute CA cannot be used.

---

### User Story 5 - Deep-dive into controller logs (Priority: P2)

Once the layered debugging path has pointed an operator at "controller reconciliation activity" (Layer 4) as the likely cause, the operator wants one more click to go deeper: see the actual controller log output behind that evidence, without leaving Observātiō and without running `kubectl logs` by hand. Since no log-viewing capability exists today, this introduces a new "Logs" destination reachable from the dashboard's existing lateral navigation, and from a "deep dive" action on the debugging path when the controller-activity layer is implicated. It shows the relevant controller's own Pod log output (CAPI core in `capi-system`, or the specific infrastructure/bootstrap provider's controller in its own namespace, e.g. `capd-system`) — the same data `kubectl logs` would show, retrieved the same way, not node/Machine-level output.

**Why this priority**: This is the natural continuation of US2's debugging path — knowing *which* layer is at fault (US2) is necessary but sometimes insufficient when that layer is controller activity; operators periodically need the raw reconciliation log output behind it. It's priority P2 alongside proactive risk detection because, unlike US2, most stuck-object cases are already resolved by the layer identification alone, without needing to read controller logs directly.

**Independent Test**: Can be fully tested by seeding a case where the controller-activity layer is implicated (per US2 Acceptance Scenario 4), opening the deep-dive action, and confirming the relevant controller's log output streams into the new Logs view; and separately by confirming the Logs view is also reachable directly from the lateral navigation without first drilling into an object.

**Acceptance Scenarios**:

1. **Given** an object's debugging path implicates the controller-activity layer, **When** the operator chooses the deep-dive action, **Then** the new Logs view opens scoped to the responsible controller (CAPI core or the specific provider) and streams its Pod log output.
2. **Given** the Logs view is reachable from a debugging path, **When** an operator instead wants to browse controller logs without first drilling into a specific object, **Then** they can also reach it directly as its own item in the lateral navigation, choosing which controller to view.
3. **Given** a controller's logs cannot be retrieved (e.g., its Pod has been evicted/restarted and log history is gone, or the backend lacks permission), **When** the operator opens the Logs view for it, **Then** an explicit "logs unavailable" state is shown rather than a blank pane.
4. **Given** an operator wants node-level access beyond controller logs (the original request's Layer 5), **When** they look for it from a Machine's debugging path, **Then** they are shown SSH connection instructions (command + node address) for VM-based providers rather than a live terminal, and Observātiō does not store or transmit SSH credentials.

---

### Edge Cases

- What happens when the dashboard cannot reach the management cluster's API server at all (Level 3/4 failure)? The dashboard MUST distinguish "no data because the source is unreachable" from "no data because everything is healthy," and must not silently show a healthy-looking empty state in that case.
- How does the dashboard behave when an object exhibits conditions consistent with more than one layer failing at once (e.g., a provider crash-loop coinciding with a stuck Machine)? The layered path shown MUST reflect all contributing layers rather than only the first one found.
- What happens when certificate-expiry data, provider-version data, or drift data cannot be determined for a given object (e.g., insufficient permissions, unsupported provider)? The dashboard MUST show that the check could not be performed rather than omitting the risk category silently.
- What happens when a MachineHealthCheck's `maxUnhealthy` threshold is breached? The dashboard MUST surface this distinctly, since it may indicate a network partition rather than independent node failures, and remediation is intentionally paused.
- What happens when infrastructure has drifted out-of-band but CAPI still reports the object as otherwise healthy? The drift warning MUST be shown as an independent signal, not suppressed by an otherwise-healthy rollup.
- What happens when a cluster's CA secret is missing but the workload cluster is still serving traffic normally? The dashboard MUST still surface this as a top-severity warning, since the failure is otherwise invisible until the next certificate rotation or node replacement is attempted.
- What happens when a controller's logs cannot be retrieved (Pod restarted/evicted with no retained history, or insufficient permissions)? The Logs view MUST show an explicit "logs unavailable" state rather than a blank pane.
- What happens when an operator opens the node-access deep-dive for a VM-based-provider Machine expecting live output? It MUST clearly show that this is connection guidance, not live streamed output, so the operator isn't left waiting on a pane that will never populate.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST provide a single, new top-level dashboard view that consolidates management-plane objects (Clusters, ClusterClasses) and workload-plane objects (Machines, MachineDeployments, MachineSets) into one screen, organized by category.
- **FR-002**: The dashboard MUST present a rolled-up health summary (healthy / degraded / failed counts) per category, visible without navigating to a category's dedicated list page.
- **FR-003**: The dashboard MUST allow narrowing to a single category in place, without a full page navigation.
- **FR-004**: For any unhealthy object, the system MUST determine and display, directly on the first (landing) dashboard screen — not only after a drill-in click — which debugging layer(s) — object conditions, Machine phase, provider-specific resource, or controller reconciliation activity — most likely explain the failure, using only data already visible to the dashboard (no manual CLI steps required by the operator).
- **FR-005**: The system MUST label each layer shown in a debugging path with its position in the sequence (conditions → phase → provider resource → controller activity) so the operator understands how deep into the stack the shown evidence goes.
- **FR-006**: The system MUST surface provider-specific resource conditions/events (e.g., `DockerMachine`, `VSphereMachine`) inline against the parent Machine, without requiring a separate lookup.
- **FR-007**: The system MUST surface controller-level reconciliation signals (e.g., recent error-level activity for the responsible controller) tied to a specific stalled object, and MUST only foreground this layer when higher layers (conditions, phase, provider resource) do not already explain the failure.
- **FR-008**: The system MUST detect and flag certificates approaching expiry within a configurable warning window, shown against the affected cluster.
- **FR-009**: The system MUST detect MachineDeployment rollouts that have stalled (new MachineSet healthy, old MachineSet not scaling down) and flag the likely blocking cause when determinable (e.g., pod disruption budget, finalizer).
- **FR-010**: The system MUST detect version skew between an installed infrastructure/bootstrap provider and its CRDs, and flag the affected provider.
- **FR-011**: The system MUST detect infrastructure drift (a provider-managed resource's observed state no longer matching its CAPI spec) and flag the affected object.
- **FR-012**: The system MUST distinguish self-healing failures (e.g., `MachineHealthCheck`-driven node replacement in progress) from failures requiring human action, and MUST NOT present self-healing activity with the same urgency as an actionable alert.
- **FR-013**: The system MUST detect when `MachineHealthCheck` remediation has paused due to a breached `maxUnhealthy` threshold and escalate this as a distinct "needs investigation" condition.
- **FR-014**: The system MUST detect a crash-looping or otherwise degraded provider controller and flag it as a provider-level (not merely object-level) failure.
- **FR-015**: The system MUST detect management-cluster-level degradation (API server unreachable or erroring in a way consistent with lost etcd quorum) and present this as a top-level banner distinct from individual object statuses, since it implies all lifecycle operations are blocked.
- **FR-016**: The system MUST detect when a cluster's CA secret is missing or inaccessible and present this as the highest-severity warning for that cluster, explaining that certificate issuance/rotation is blocked and that the original CA cannot be substituted.
- **FR-017**: When the dashboard's data source (the management cluster's API server) is unreachable, the system MUST show an explicit "data unavailable" state rather than an empty or falsely healthy-looking state.
- **FR-018**: When a risk check (certificate expiry, version skew, drift) cannot be evaluated for a given object (e.g., unsupported provider, insufficient data), the system MUST show that the check could not be performed rather than omitting it silently.
- **FR-019**: The system MUST provide a new "Logs" destination, reachable both as its own item in the existing lateral navigation and as a "deep dive" action from an object's debugging path when the controller-activity layer is implicated, since no log-viewing capability exists in Observātiō today.
- **FR-020**: The Logs view MUST stream the relevant controller's Pod log output (CAPI core in `capi-system`, or the specific infrastructure/bootstrap provider's controller in its own namespace) via the standard Kubernetes Pod-log API — the same data and mechanism `kubectl logs` uses.
- **FR-021**: For a Machine's node-access deep-dive on a VM-based-provider, the system MUST show SSH connection instructions (command and node address) rather than live log output or an in-dashboard terminal.
- **FR-022**: The system MUST NOT store, manage, or transmit SSH credentials of any kind; VM-based-provider node access remains something the operator performs on their own machine using the shown instructions.
- **FR-023**: When a controller's logs cannot be retrieved, the Logs view MUST show an explicit "logs unavailable" state rather than a blank pane.

Machine/node-level log streaming (e.g., a CAPD Machine's underlying container output, or a VM-based
node's kubelet/containerd/cloud-init logs) is explicitly a **TODO for a future iteration**, not part
of this feature — FR-019–FR-023 cover controller logs only; node-level access in this feature is
limited to the static SSH connection instructions in FR-021.

### Key Entities

- **Debugging Layer**: One of the ordered stages used to explain an object's failure — object conditions, Machine phase, provider-specific resource, controller reconciliation activity. Has a position/order and the evidence (conditions, events, log signals) associated with it for a given object.
- **Health Rollup**: A per-category summary (healthy / degraded / failed counts) shown on the centralized dashboard.
- **Risk Warning**: A proactively detected issue (certificate expiry, stalled rollout, provider version skew, infrastructure drift) attached to a specific object, with a category, severity, and — where determinable — a likely cause.
- **Failure Severity**: A classification of a detected issue as self-healing, needs-investigation, provider-degraded, or management-cluster-critical, used to set the urgency shown to the operator.
- **Debugging Path**: The ordered set of Debugging Layers and their evidence shown for a specific unhealthy object, representing the "path" the operator would otherwise have had to walk manually across CLI tools.
- **Log View**: The new destination (lateral-navigation item and debugging-path deep-dive target) streaming a specific controller's (CAPI core or infrastructure/bootstrap provider) Pod log output.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: An operator can determine the overall health of all cluster-management and workload objects within 10 seconds of opening the dashboard, without visiting any other page.
- **SC-002**: For a stuck or failing object, an operator can identify the responsible debugging layer using only the dashboard — without running any CLI command — in under 60 seconds, for at least 90% of common failure scenarios (stuck provisioning, stuck bootstrap, provider error, controller crash).
- **SC-003**: Certificate-expiry warnings appear on the dashboard at least 30 days before expiry by default, with zero false negatives against known expiry dates in test data.
- **SC-004**: Stalled-rollout, version-skew, and drift warnings each appear within one monitoring cycle of the underlying condition existing in the cluster, with the likely cause shown for at least 80% of stalled-rollout cases where a blocking PDB or finalizer is present.
- **SC-005**: Self-healing events (e.g., MachineHealthCheck-driven remediation) are never shown with the same visual urgency as an actionable failure, verified across all defined severity levels.
- **SC-006**: When the management cluster's API server is unreachable, 100% of dashboard views show an explicit "data unavailable" banner rather than a false "all healthy" or blank state.
- **SC-007**: When a debugging path implicates the controller-activity layer, an operator can reach that controller's actual log output in one click, without leaving Observātiō and without running `kubectl` themselves.

## Assumptions

- The dashboard operates against a single management cluster per session (matching the existing product's scope); multi-management-cluster aggregation is out of scope for this feature.
- "Controller reconciliation activity" (FR-007) is derived from log/event data the backend already has access to or can reasonably collect (e.g., recent controller log lines or Kubernetes Events associated with an object); it does not require deploying new observability infrastructure (e.g., a metrics/tracing stack) as part of this feature.
- Certificate expiry, version-skew, and drift detection rely on data obtainable from the Kubernetes API (Secrets, CRD versions, object spec vs. status) without requiring SSH/node-level access.
- User Story 5's Logs view covers **controller** logs only (Layer 4 — CAPI core and infrastructure/bootstrap provider controllers), retrieved via the standard Kubernetes Pod-log API, the same mechanism and data `kubectl logs` provides — no new external dependency. Node-level access (the original request's Layer 5 — Machine/container output, cloud-init logs, kubelet journal) is limited in this feature to the static SSH connection instructions in FR-021; actual Machine/node-level log streaming (e.g., a CAPD Machine's container output) is an explicit **TODO for a future iteration**, not built now.
- "Configurable warning window" for certificate expiry defaults to 30 days and is not required to be end-user-configurable in this feature's first iteration; a fixed sensible default satisfies the requirement unless later clarified otherwise.
- Existing object detail screens (Clusters, Machines, MachineDeployments) and their per-object AI troubleshooting/YAML tree panels (delivered in a prior feature) remain the destination when an operator drills in from this new centralized dashboard; this feature adds the layered, categorized landing view and root-cause path on top of what already exists, it does not replace per-object detail screens.
