# Implementation Plan: MCP Server Aggregation & Local Tool Server

**Branch**: `009-mcp-server-client-aggregator` | **Date**: 2026-07-17 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/009-mcp-server-client-aggregator/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

Today the "Ask AI about this" panel's assistant (`webserver/internal/infra/llm`) has exactly one
hard-wired tool — `kubectl`, defined in `tools.go` and dispatched via a `switch` in
`observation.go`'s `runToolCalls`. This feature replaces that single hard-wired tool with a small
aggregator: a new `webserver/internal/infra/mcp` package defines a `ToolSource` interface, wraps
the existing kubectl capability as an in-process `LocalToolSource` (US1), and lets operators
register additional external MCP servers (stdio or streamable-HTTP) via a YAML config file loaded
once at startup (US2, no runtime API — FR-010). The aggregator merges every healthy source's
capabilities into one unified tool list handed to Claude, filters out anything not declared
`readOnlyHint: true` by its source (FR-009/SC-004), detects and reports capability-name conflicts
at startup rather than silently merging them (FR-007), and runs a background health check per
external source so one misbehaving source degrades gracefully instead of blocking the others
(US3). A new `GET /api/mcp/sources` REST endpoint and a matching frontend status card give
operators visibility into registered sources, their health, and their capabilities (FR-003,
SC-005) without adding a mutation surface.

## Technical Context

**Language/Version**: Go 1.25 (backend, matches `webserver/go.mod`), TypeScript 5 / React 19 /
Next.js 15 (frontend — status view only, no new route)
**Primary Dependencies**: `github.com/anthropics/anthropic-sdk-go` (existing); NEW:
`github.com/modelcontextprotocol/go-sdk` (MCP client, `stdio` + streamable-HTTP transports —
research.md R1); `sigs.k8s.io/yaml` (already an indirect dependency, promoted to direct for
config-file parsing — research.md R3)
**Storage**: N/A — the tool-sources config is a YAML file read once at process startup, not a
database; no persistence is added
**Testing**: `go test` + `testify` (existing pattern, e.g. `conversation_test.go`); Jest for the
new frontend status card, matching the existing `ops/` card test convention
**Target Platform**: Existing Observātiō web app (Next.js static export embedded in the Go binary)
**Project Type**: Web application (existing `webserver/` + `front/` structure)
**Performance Goals**: Health checks run out-of-band on a background interval (research.md R4),
never inline with a chat turn — a slow/unreachable external source must not add latency to "Ask
AI" turns that don't need it
**Constraints**: Registration/modification/removal of a tool source is config-file + restart only,
never a runtime API or chat-panel action (FR-010); no capability reaches the model unless verified
`readOnlyHint: true` at its source (FR-009/SC-004, fail-closed per research.md R5); with zero
external sources configured, behavior must be unchanged from today (US1 Independent Test)
**Scale/Scope**: A reasonable number of simultaneously registered tool sources (low tens per
spec.md Assumptions) — not designed to scale to hundreds

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **I. Observability & Data Consolidation** — PASS. Tool source health/capabilities are structured
  data, consolidated into one status view (`GET /api/mcp/sources`) rather than scattered across
  logs; every chat response is traceable to the source(s) that produced it (FR-008/SC-003,
  research.md R7).
- **II. Real-Time Visibility** — PASS with a scoped, documented exception. This principle's own
  bullet list mandates WebSocket delivery for *cluster state* specifically ("Watchers for Cluster,
  Machine, and MachineDeployment"). Tool-source health is assistant-tooling/configuration state,
  not CAPI cluster state, so the new status endpoint is deliberately a polled REST read
  (research.md R8), matching the existing `/api/infra/capabilities` precedent — no cluster-state
  UI in this feature bypasses WebSocket delivery.
- **III. ClusterAPI Resource Model Compliance** — N/A. Tool sources and capabilities are not CAPI
  resources and introduce no infrastructure-provider abstraction into the core domain layer; the
  new `mcp` package is fully independent of `internal/infra/models`.
- **IV. AI-Augmented Troubleshooting** — PASS. Every registered capability must be verified
  read-only before it's callable (FR-009), directly reinforcing the existing "AI-generated
  recommendations MUST be presented as operator suggestions, never as automated remediation
  actions" constraint by making it structurally impossible for the assistant to mutate state
  through any registered tool, built-in or external. Conversation scoping (`ConversationManager`)
  is unchanged.
- **V. Test-Driven Quality** — PASS. New `webserver/internal/infra/mcp` package gets `_test.go`
  coverage for conflict detection, health-state transitions, read-only filtering, and dispatch/
  attribution (data-model.md); the new frontend status card gets a Jest test matching the existing
  `ops/` card convention (e.g. `backup-health-card.test.tsx`). `make run-tests-backend` and
  `make run-tests-frontend` must both pass before merge.

**Result**: No unjustified violations. One new runtime dependency
(`github.com/modelcontextprotocol/go-sdk`) is introduced and justified in Complexity Tracking below,
per the constitution's Technology Stack rule ("No new runtime dependency may be introduced without
updating this section and providing justification..."). This plan documents the justification here;
recording the dependency in the constitution's own Technology Stack section is a follow-up action
for `/speckit-constitution`, not performed as a side effect of this plan.

## Project Structure

### Documentation (this feature)

```text
specs/009-mcp-server-client-aggregator/
├── plan.md              # This file (/speckit-plan command output)
├── research.md          # Phase 0 output (/speckit-plan command)
├── data-model.md         # Phase 1 output (/speckit-plan command)
├── quickstart.md         # Phase 1 output (/speckit-plan command)
├── contracts/            # Phase 1 output (/speckit-plan command)
│   ├── mcp-sources-api.md
│   └── tool-sources-config.md
└── tasks.md              # Phase 2 output (/speckit-tasks command - NOT created by /speckit-plan)
```

### Source Code (repository root)

```text
webserver/
├── internal/infra/mcp/                    # NEW package (research.md R2)
│   ├── source.go                          # ToolSource interface, Capability, HealthStatus types
│   ├── local.go                           # LocalToolSource — kubectl exposed via a real mcp.Server,
│   │                                      #   connected in-process over mcp.NewInMemoryTransports()
│   │                                      #   (spec.md Clarifications 2026-07-17, research.md R2)
│   ├── external.go                        # MCPToolSource — go-sdk-backed client per registered source
│   ├── config.go                          # SourceConfig YAML loader + validation (contracts/tool-sources-config.md)
│   ├── health.go                          # background health-check goroutine, state transitions
│   ├── readonly.go                        # readOnlyHint-based fail-closed filtering (research.md R5)
│   ├── aggregator.go                      # merge, conflict detection (R6), RenderTools/Dispatch/Status
│   └── *_test.go                          # NEW — one per file above, following existing testify convention
│
├── internal/infra/llm/                    # MODIFY existing package
│   ├── observation.go                     # ObservationService holds *mcp.Aggregator instead of
│   │                                      #   a static []anthropic.ToolUnionParam; runToolCalls
│   │                                      #   delegates to Aggregator.Dispatch instead of a switch
│   └── tools.go                           # MODIFY: KubectlTool()/RunKubectl() logic relocates into
│                                          #   infra/mcp/local.go; this file shrinks or is removed
│
├── internal/web/handlers/
│   ├── handlers.go                        # MODIFY: register `GET /api/mcp/sources`
│   └── system/mcp_sources.go              # NEW: HandleMCPSources (contracts/mcp-sources-api.md)
│                                          #   pattern-matched to kubernetes/cluster.go's
│                                          #   HandleInfraCapabilities
│
└── cmd/server.go                          # MODIFY: --tool-sources-config flag / TOOL_SOURCES_CONFIG
                                           #   env var (matching existing --address/--dev convention),
                                           #   passed through to llm.NewObservationService

front/
├── app/ui/dashboard/shared/
│   └── use-tool-sources.ts                # NEW: fetch hook for GET /api/mcp/sources (data-model.md)
├── app/ui/dashboard/components/ops/
│   ├── tool-sources-card.tsx              # NEW: status card, cloned from health-rollup-card.tsx
│   │                                      #   pattern (health badge per source, conflict indicator)
│   ├── tool-sources-card.test.tsx         # NEW
│   └── ops-dashboard.tsx                  # MODIFY: render ToolSourcesCard
```

**Structure Decision**: New capability aggregation logic gets its own package
(`webserver/internal/infra/mcp`) rather than living inside `internal/infra/llm`, because it is a
general "aggregate tool sources" concern independent of the Anthropic-specific chat/streaming code
that already fills `llm` — `llm` becomes a *consumer* of `mcp`'s aggregator, the same relationship
`llm` already has with `clusterapi` for cluster data. The status view is added to the existing
Day-2 Ops `ops/` card family (frontend) rather than a new page/route, since FR-003's "view
registered tool sources" is a small, read-only status surface that fits the same visual slot as
the existing infra-capabilities-style cards — no new frontend route is introduced.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| New runtime dependency: `github.com/modelcontextprotocol/go-sdk` | FR-002 requires registering real external MCP tool sources; the MCP protocol's JSON-RPC framing, capability negotiation, and stdio/HTTP transports are substantial, spec-governed surface area (research.md R1) | Hand-rolling a JSON-RPC/MCP client — rejected: reimplements protocol plumbing with no product-differentiating value and an ongoing maintenance burden every time the spec changes |
| New config-loading pattern: YAML file read at startup (first in this repo) | FR-002/FR-010 require registering sources "without a code change or new release" while keeping registration an administrative, not runtime, action (research.md R3) | Per-source environment variables — rejected: doesn't scale past 1-2 sources and can't cleanly express nested per-transport fields (research.md R3) |
