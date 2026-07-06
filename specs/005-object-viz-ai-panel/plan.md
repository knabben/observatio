# Implementation Plan: Object YAML View & Global AI Troubleshooting Panel

**Branch**: `005-object-viz-ai-panel` | **Date**: 2026-07-06 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/005-object-viz-ai-panel/spec.md`

## Summary

Two changes to how operators inspect objects and get AI help. First, every object detail screen
(Cluster, Cluster Infrastructure — Docker and vSphere, Machine, Machine Infrastructure — Docker and
vSphere, Machine Deployment) gains a "YAML" tab rendering the *complete* underlying object as an
expandable/collapsible tree — not the curated subset each screen's "Specification" tab shows today.
Second, AI troubleshooting stops being a tab embedded inside one object type's detail view and
becomes a single, app-wide collapsible side panel (reusing the `Drawer` pattern already used for
mobile navigation) reachable from anywhere, auto-populated with a richer description of whatever
object is currently in view, restyled onto the dashboard's existing theme tokens, and given a
one-click "Ask AI about this" quick-action on every detail screen.

## Technical Context

**Language/Version**: Go 1.25 (backend, `webserver/`); TypeScript 5.9, React 19, Next.js 15.3 (frontend, `front/`)
**Primary Dependencies**: No new dependencies. Backend: `k8s.io/client-go` dynamic client (already used by the Docker infra fetchers and WebSocket watchers) for a new generic raw-object `Get`. Frontend: Mantine core's built-in `Tree`/`useTree`/`TreeNodeData` (already present in the pinned `@mantine/core` 7.17.8 — confirmed via package inspection, no version bump needed) for the tree tab; the existing `Drawer` component (already used by `sidenav.tsx` for mobile nav) for the collapsible AI panel.
**Storage**: N/A — stateless; raw object fetched on-demand, AI conversation kept in-memory (React state) for the browser session only, never persisted (Constitution IV: conversation history must not persist across sessions without explicit consent)
**Testing**: Go `testing` + `testify` (existing pattern); Jest 29 + `@testing-library/react` via the shared `test-render.tsx` helper (existing pattern)
**Target Platform**: Same single-binary deployment — Next.js static export embedded in and served by the Go binary
**Project Type**: Web application — both backend (`webserver/`) and frontend (`front/`) are touched
**Performance Goals**: Raw-object tab fetch is on-demand (only when that tab is first opened) plus one re-fetch per live `resourceVersion` change observed on the object's existing WebSocket stream — no independent polling loop
**Constraints**: New raw-object endpoint is read-only (`Get`, no mutation); AI panel conversation state is in-memory only, not persisted to any storage; no new runtime dependency (Technology Stack clause) — both the tree UI and the collapsible panel reuse components already in the dependency tree
**Scale/Scope**: 6 detail-screen variants gain a YAML tab (Cluster, ClusterInfra×2 providers, Machine, MachineInfra×2 providers, MachineDeployment); 1 new backend endpoint; the AI panel moves from N per-object embeds to 1 app-wide instance; 16 FRs, 7 success criteria

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Constitution v1.1.0 — five principles:

| Principle | Impact | Verdict |
|-----------|--------|---------|
| **I. Observability & Data Consolidation** | Strengthens observability directly: operators can now see every field the backend returns, not just what a screen author chose to surface, and status stays consolidated (not scattered) in the Specification tab. | ✅ PASS |
| **II. Real-Time Visibility** | The raw-object tab uses REST (`Get`), not WebSocket, for the on-demand tree — a deliberate, scoped exception. It is *driven by* the object's existing WS stream (re-fetches when that stream's `resourceVersion` changes) rather than independent polling, and only activates when the operator opens the tab; the curated Specification tab and every list screen remain 100% WS-delivered, unchanged. | ✅ PASS (scoped, WS-triggered REST hydration — see research.md R2) |
| **III. ClusterAPI Resource Model Compliance** | The new raw-object endpoint is a thin, generic `(group, version, resource, namespace, name) → Get` passthrough — no proprietary domain type, no infra-provider-specific branching. Existing curated models (Cluster, Machine, etc.) are untouched. | ✅ PASS |
| **IV. AI-Augmented Troubleshooting** | Directly strengthens this principle: today's AI context is arguably already thin (a bare condition-reason string); auto-populating identity + status/conditions + key spec fields moves it toward genuine "structured condition data... before generating output." Conversation stays in-memory per session, never persisted — no new consent/retention concern. | ✅ PASS |
| **V. Test-Driven Quality** | New Go tests for the raw-object handler's GVR-parsing/validation logic; new Jest tests for the AI panel context (auto-context refresh, edit-locks-prefill), the object-tree conversion utility, and the quick-action wiring. | ✅ PASS |

**Result**: No violations. No entries required in Complexity Tracking.

## Project Structure

### Documentation (this feature)

```text
specs/005-object-viz-ai-panel/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
│   └── raw-object-api.md
└── checklists/
    └── requirements.md  # Spec quality checklist (from /speckit-specify)
```

### Source Code (repository root)

```text
webserver/
└── internal/web/handlers/
    └── kubernetes/
        └── raw.go                          # NEW: HandleRawObject — generic dynamic-client Get by
                                              # group/version/resource/namespace/name; registered as
                                              # GET /api/raw in handlers.go

front/
├── app/ui/dashboard/
│   ├── ai-panel/                            # NEW: app-wide AI panel, replacing the per-object embed
│   │   ├── ai-panel-context.tsx             # Global open/closed + conversation + auto-context state
│   │   ├── ai-panel.tsx                     # The Drawer-based panel UI (restyled onto theme tokens),
│   │   │                                    # built from the existing ChatBot markup in
│   │   │                                    # base/ai-troubleshooting.tsx
│   │   ├── ai-panel-trigger.tsx             # Global open control, mounted in dashboard/layout.tsx
│   │   ├── use-current-object-context.ts    # Hook each detail screen calls to register/unregister
│   │   │                                    # its current object with the panel's auto-context
│   │   └── ask-ai-button.tsx                # Per-object-screen quick-action ("Ask AI about this")
│   ├── shared/
│   │   ├── object-tree.tsx                  # NEW: renders TreeNodeData[] via Mantine's Tree,
│   │   │                                    # fetches/refreshes the raw object on demand
│   │   └── to-tree-data.ts                  # NEW: arbitrary JSON -> Mantine TreeNodeData[] utility
│   └── base/
│       ├── details.tsx                      # Unchanged (ObjectDetails still just renders tabs)
│       └── ai-troubleshooting.tsx            # REMOVED — replaced by ai-panel/*
├── app/dashboard/layout.tsx                  # Mount <AIPanelProvider> + <AIPanelTrigger>
└── app/ui/dashboard/components/
    ├── clusters/{details.tsx, infra/{infra-details.tsx, docker-details.tsx}}         # add YAML tab, AskAIButton; drop AI Troubleshooting tab
    ├── machines/{details.tsx, infra/{infra-details.tsx, docker-details.tsx}}         # same
    └── mds/details.tsx                                                              # same (if not already present, extend existing detail component)
```

**Structure Decision**: Existing single Go binary (`webserver/`) + existing single Next.js frontend
(`front/`) — no new top-level project, no new dependency. Backend work is a single small generic
handler. Frontend work replaces the per-object-embedded AI section with a new small `ai-panel/`
module (context + Drawer UI + trigger + quick-action) reusing the existing `Drawer` pattern and
the `InfraCapabilityContext`-style context-provider pattern already established in feature `004`,
plus a new `shared/object-tree.tsx` built on Mantine's already-present `Tree` component. Every
existing detail component (`ClusterDetails`, `ClusterInfraDetails`, `ClusterInfraDockerDetails`,
`MachineDetails`, `MachineInfraDetails`, `MachineInfraDockerDetails`, and the Machine Deployment
details component) gets the same two edits: drop the "AI Troubleshooting" tab (folding its
conditions table into "Specification"), add the new "YAML" tab.

## Complexity Tracking

> No Constitution Check violations — this section is intentionally empty.
