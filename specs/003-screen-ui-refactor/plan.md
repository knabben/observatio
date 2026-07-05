# Implementation Plan: Screen Refactoring & UI Tech-Debt Remediation

**Branch**: `003-screen-ui-refactor` | **Date**: 2026-07-05 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/003-screen-ui-refactor/spec.md`

## Summary

Remediate user-facing defects and structural tech debt across the five dashboard screens (Dashboard
overview, Clusters, Machine Deployments, Machines, and the shared navigation/AI shell). The work
makes every screen tolerant of partial/empty/error data, corrects broken and non-responsive layouts,
fixes navigation/accessibility/status legibility, consolidates four copy-pasted resource-area
component sets into shared configurable building blocks, and centralizes theming. No new product
capability is introduced. The frontend addresses the backend **same-origin** (it is embedded in and served
by the single Go binary), and a new `make verify-binary` target proves the embed+serve pipeline end-to-end
(US6). All existing dependencies are refreshed to the latest **within their current major** (frontend +
backend; breaking majors deferred). Every changed component ships with Jest tests covering the partial-data,
empty-collection, and error-state cases.

## Technical Context

**Language/Version**: TypeScript 5, React 19, Next.js 15.3 (App Router, static export — `output: "export"`, `distDir: ./output`)
**Primary Dependencies**: Mantine UI 7.17 (`@mantine/core`, `@mantine/charts`), TailwindCSS 4, XYFlow 12 (`@xyflow/react`), Recharts 2, `react-use-websocket` 4.13, `@tabler/icons-react`, `@heroicons/react`, `uuid` 11 (already present)
**Storage**: N/A — stateless SPA; all state derived from live WebSocket stream + REST reads
**Testing**: Jest 29 + `ts-jest` + `jest-environment-jsdom` + `@testing-library/react`; shared `render()` helper wraps `MantineProvider` (`app/ui/dashboard/utils/test-render.tsx`); run via `make run-tests-frontend` (`pnpm test`)
**Target Platform**: Static-exported SPA (served as static assets); modern evergreen browsers, desktop/laptop primary, tablet usable
**Project Type**: Web application — frontend only for this feature (`front/`)
**Performance Goals**: Event-to-UI render latency < 2s (Constitution II); loading state resolves within a bounded threshold (FR-003, target 10s)
**Constraints**: WebSocket MUST remain the primary live transport (no switch to polling — Constitution II); the SPA is embedded in and served by the single Go binary, so the frontend addresses the API/WS **same-origin** from `window.location` (FR-036) — `NEXT_PUBLIC_*` is a dev-only override; no backend API/WS message-shape changes; WCAG AA contrast + full keyboard operability (FR-018–FR-021)
**Dependency currency**: refresh all deps to latest **within current major** — frontend via `pnpm update` (+ within-major `next`), backend via targeted `go get` (`gin`, `gorilla/websocket`, `cobra`, `testify`) + `go mod tidy`; defer breaking majors (Mantine 9, Jest 30, anthropic v1, k8s 0.36 / CAPI 1.13 / controller-runtime 0.24). Gate on tests + `make build` green (R11)
**Scale/Scope**: 5 screens, ~50 component files, 4 duplicated resource areas (clusters, clusters/infra, machines+machines/infra, mds) to consolidate; 39 FRs, 11 success criteria; +1 single-binary build/verify concern (US6)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Constitution v1.1.0 — five principles:

| Principle | Impact | Verdict |
|-----------|--------|---------|
| **I. Observability & Data Consolidation** | Adds explicit empty/error/unknown states and surfaces failure info (`failureReason`/`failureMessage`) previously dropped; consolidates status into a single legible tri-state indicator. Strengthens observability. | ✅ PASS |
| **II. Real-Time Visibility** | WebSocket stays the primary transport; work fixes bounded reconnect + error surfacing and MUST NOT introduce REST polling for live views. Render-latency budget (<2s) preserved. | ✅ PASS |
| **III. ClusterAPI Resource Model Compliance** | Shared components are parameterized over the existing CAPI-mapped types (Cluster → MachineDeployment → Machine); infra provider fields remain opaque, not promoted to core domain. No model changes. | ✅ PASS |
| **IV. AI-Augmented Troubleshooting** | AI panel keeps cluster-condition grounding and single-session scope; work removes unsafe content rendering (XSS) and layout misuse, and adds open-socket/error handling. AI remains suggestion-only. | ✅ PASS |
| **V. Test-Driven Quality** | Every changed component gets Jest tests for partial-data, empty-collection, and error states; `make run-tests-frontend` must pass before merge. Tests authored alongside each refactor. | ✅ PASS |

**Result**: No violations. No entries required in Complexity Tracking.

**Post-clarification notes**: (a) The within-major dependency refresh (R11) introduces **no new** runtime
dependency and keeps Mantine on v7 — consistent with this refactor and the constitution's Technology Stack
clause (no stack change to justify). (b) The single-binary embed + `make verify-binary` work (US6) touches
the Makefile and SPA-serving path but makes no CAPI/API change — it strengthens Principle V (an automated,
end-to-end release gate) and leaves Principles I–IV intact.

## Project Structure

### Documentation (this feature)

```text
specs/003-screen-ui-refactor/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output (frontend view-models & configs)
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (UI/data contracts)
│   ├── data-channels.md
│   ├── shared-components.md
│   ├── environment-config.md   # revised: same-origin addressing
│   └── build-verification.md   # single-binary embed + make verify-binary
└── checklists/
    └── requirements.md  # Spec quality checklist (from /speckit-specify)
```

### Source Code (repository root)

```text
front/
├── app/
│   ├── layout.tsx                     # Root: color-scheme + theme + fonts (US4)
│   ├── globals.css                    # Remove dead font vars, reconcile scheme (US4)
│   ├── error.tsx                      # NEW: root error boundary (US1)
│   ├── styles/
│   │   ├── fonts.ts                   # Correct font identifiers/weights (US4)
│   │   └── theme.ts                   # NEW: Mantine theme — colors, status tokens (US4)
│   ├── lib/
│   │   ├── config.ts                  # NEW: env-driven endpoint config (FR-036)
│   │   ├── data.tsx                   # res.ok checks, error propagation (US1)
│   │   └── websocket.tsx              # bounded reconnect, no empty-frame wipe (US1)
│   ├── dashboard/
│   │   ├── layout.tsx                 # responsive scroll containment (US2)
│   │   ├── page.tsx                   # responsive grid; topology (US2)
│   │   ├── {clusters,machines,machinedeployments}/  # thin route shells
│   │   └── utils.tsx                  # null-safe filter (US1)
│   └── ui/dashboard/
│       ├── sidenav.tsx, nav-links.tsx, search.tsx   # nav + a11y + search (US2/US3/US5)
│       ├── base/
│       │   ├── lister.tsx             # loading/empty/error state machine (US1)
│       │   ├── details.tsx            # shared detail header (US4)
│       │   ├── types.tsx              # shared config types (US4)
│       │   └── ai-troubleshooting.tsx # safe render, containment, collapse (US5)
│       ├── shared/                    # NEW: consolidated building blocks (US4)
│       │   ├── object-table.tsx       # config-driven table
│       │   ├── status-indicator.tsx   # tri-state indicator
│       │   ├── empty-state.tsx        # labeled empty state
│       │   ├── error-state.tsx        # error + retry
│       │   └── resource-hooks.ts      # shared fetch hook (dedupe US4)
│       ├── components/{clusters,machines,mds}/  # collapse to config over shared/
│       └── utils/{header,loader,panel}.tsx      # center loader, tokens (US2/US4)
├── package.json                        # deps refreshed within-major (R11)
└── (jest.config.js, next.config.ts unchanged)

webserver/
├── internal/web/handlers/
│   ├── handlers.go                     # //go:embed build/* (unchanged mechanism)
│   ├── build/                          # embed target (populated by make build)
│   └── system/spa.go                   # SPA fallback → index.html for client routes (US6)
├── go.mod / go.sum                     # deps refreshed within-major (R11)

Makefile                                # NEW: verify-binary target (US6, FR-037–FR-039)
```

**Structure Decision**: Single existing Next.js App-Router frontend under `front/`. This feature adds a
`front/app/ui/dashboard/shared/` directory for the consolidated building blocks, a `front/app/styles/theme.ts`
and `front/app/lib/config.ts` for centralized theming and same-origin endpoint configuration, and a root
`front/app/error.tsx`. The per-resource `components/{clusters,machines,mds}` files shrink to configuration
that composes the shared blocks. Beyond the frontend, three deliberate cross-layer touches are in scope:
(1) a `make verify-binary` target and possible SPA-fallback confirmation in `webserver/.../system/spa.go`
for the single-binary embed check (US6); (2) within-major dependency refreshes to `front/package.json` and
`webserver/go.mod`; no backend **API/message-shape** changes are made.

## Complexity Tracking

> No Constitution Check violations — this section is intentionally empty.
