<!--
SYNC IMPACT REPORT
==================
Version change: N/A (initial creation) → 1.0.0
Modified principles: N/A — first version
Added sections:
  - Core Principles (I–V)
  - Technology Stack
  - Development Workflow
  - Governance
Removed sections: N/A
Templates reviewed:
  - .specify/templates/plan-template.md ✅ Constitution Check section compatible
  - .specify/templates/spec-template.md ✅ User story / requirements structure aligns
  - .specify/templates/tasks-template.md ✅ Phase structure aligns with principle V (test-driven)
Follow-up TODOs: None — all placeholders resolved
-->

# Observātiō Constitution

## Core Principles

### I. Observability-First

Every cluster component MUST expose structured health and status data. Metrics, conditions,
events, and logs MUST be collected across all infrastructure layers — from the management
cluster down to individual Machines. Data that cannot be observed cannot be managed.

- All new API handlers MUST include structured logging with request context.
- Cluster, MachineDeployment, and Machine resources MUST expose condition-level health status.
- Error states MUST propagate to the dashboard with actionable, human-readable descriptions.

### II. Real-Time Visibility

Data presented to operators MUST reflect the current cluster state without requiring a manual
page refresh. WebSocket streaming is the primary delivery mechanism for live data.

- UI components displaying cluster state MUST consume WebSocket events, not poll REST endpoints.
- Watchers for Cluster, Machine, and MachineDeployment MUST emit delta updates on every change.
- Dashboard render latency from server event to UI update MUST stay below 2 seconds under
  normal operating load.

### III. ClusterAPI Resource Model Compliance

All features and domain models MUST align with the ClusterAPI resource hierarchy:
Cluster → MachineDeployment → Machine, with optional ClusterClass support.
No proprietary infrastructure-specific abstractions are permitted at the core domain layer.

- Models in `webserver/internal/infra/models/` MUST map directly to CAPI CRDs.
- Infrastructure provider details (vSphere, AWS, etc.) MUST be treated as opaque annotations,
  not first-class types in the core domain.
- CAPI resource watchers MUST handle both creation and deletion lifecycle events.

### IV. AI-Augmented Troubleshooting

The platform MUST provide AI-assisted diagnosis for cluster failure scenarios. LLM integration
MUST be grounded in actual cluster state — it MUST NOT generate advice disconnected from the
specific resource conditions currently observed.

- AI analysis MUST receive structured condition data from the affected resource before generating
  output; generic prompts without cluster context are not permitted.
- LLM tool definitions MUST be scoped to specific cluster resources to reduce hallucination risk.
- AI-generated recommendations MUST be presented as operator suggestions, never as automated
  remediation actions.
- Conversation history MUST be scoped to a single troubleshooting session and MUST NOT persist
  across sessions without explicit user consent.

### V. Test-Driven Quality

All backend and frontend changes MUST be accompanied by tests. Tests MUST be written before or
alongside implementation — not as an afterthought appended after the feature is complete.

- Backend Go packages MUST maintain `_test.go` coverage for all exported functions and handlers.
- Frontend components containing business logic MUST include Jest tests.
- Integration tests for cluster watchers and processors MUST use CAPI fake clients that faithfully
  represent the real API contract — mocks that diverge from actual CRD behaviour are not permitted.
- The full test suite (`make run-tests-backend` and `make run-tests-frontend`) MUST pass before
  any PR is merged.

## Technology Stack

**Backend**: Go 1.23+ with CAPI controller-runtime client
**Frontend**: Next.js 15, React 19, TypeScript 5, Mantine UI 7, TailwindCSS 4, XYFlow 12
**Real-time transport**: WebSocket — backend connection pool + frontend `react-use-websocket`
**AI integration**: Anthropic Claude API via `webserver/internal/infra/llm`
**Build tooling**: GNU Make — canonical targets: `build`, `run-backend`, `run-frontend`,
`run-tests-backend`, `run-tests-frontend`
**Frontend package manager**: pnpm

No new runtime dependency may be introduced without updating this section and providing
justification in the feature plan's Complexity Tracking table.

## Development Workflow

1. Features MUST be developed on a dedicated branch following the naming convention
   `###-feature-name` (e.g., `001-cluster-health-dashboard`).
2. The **Constitution Check** gate in `plan.md` MUST be completed and pass before Phase 0
   research begins; re-check MUST occur after Phase 1 design.
3. Backend and frontend components are developed and tested independently but MUST integrate
   successfully before a feature is marked complete.
4. `make build` MUST succeed before a pull request is opened.
5. All automated tests MUST pass before merge; test failures block merging with no exceptions.

## Governance

This constitution supersedes all other project practices and guidelines. In case of conflict,
the constitution wins.

Amendments require:
1. A documented rationale explaining why the change is necessary.
2. A semantic version bump: MAJOR for breaking governance changes or principle removals,
   MINOR for new principles or materially expanded guidance, PATCH for clarifications.
3. An updated `Last Amended` date.
4. A migration plan when the amendment invalidates existing features or workflows.

All feature plans MUST include a **Constitution Check** section that explicitly validates
compliance with each of the five principles. Any violation MUST be justified in the
Complexity Tracking table before the plan is approved.

**Version**: 1.0.0 | **Ratified**: 2026-04-25 | **Last Amended**: 2026-04-25
