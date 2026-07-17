---

description: "Task list for feature implementation"
---

# Tasks: MCP Server Aggregation & Local Tool Server

**Input**: Design documents from `/specs/009-mcp-server-client-aggregator/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Constitution Principle V (Test-Driven Quality) mandates test coverage for all backend
and frontend changes. Per plan.md, all new logic lives in a new `webserver/internal/infra/mcp`
package; each file there gets its own `_test.go`, written first per user story (matching the
`_test.go`-per-file convention already used by `llm/conversation_test.go`).

**Organization**: Tasks are grouped by user story (US1–US3, priorities from spec.md). Per spec.md,
US1 explicitly is "the foundation the aggregation model is built on" for the other two stories, but
it is also independently testable on its own (its Independent Test requires zero external
sources) — so it is its own phase, not folded into Foundational. Foundational here is limited to
the `ToolSource`/`Aggregator` contract and plumbing that has no independent user-facing behavior
until a concrete `ToolSource` (US1's `LocalToolSource`) exists.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)

## Path Conventions

Web app: `webserver/` (Go backend), `front/` (Next.js frontend) — per plan.md Project Structure.
This feature adds one new backend package (`webserver/internal/infra/mcp`), one new REST endpoint,
and one new frontend status card; no new frontend route.

---

## Phase 1: Setup

**Purpose**: Add the new MCP client dependency before any aggregator code is written.

- [X] T001 Run `go get github.com/modelcontextprotocol/go-sdk` in `webserver/`, promote
      `sigs.k8s.io/yaml` from an indirect to a direct dependency in `webserver/go.mod` (used by the
      config loader), and run `go build ./...` to confirm the baseline still compiles before any
      `mcp` package code is added (research.md R1, R3)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Establish the `ToolSource`/`Capability`/`Aggregator` contract every story builds on.
No concrete `ToolSource` exists yet at the end of this phase — the aggregator is fully unit-tested
against fakes, but chat has zero usable capabilities until US1 registers one.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [X] T002 [P] Define `ToolSource` interface, `Capability`, `HealthStatus`, `SourceKind` types in
      `webserver/internal/infra/mcp/source.go` per data-model.md
- [X] T003 [P] Implement fail-closed `readOnlyHint`-based capability filtering in
      `webserver/internal/infra/mcp/readonly.go` per research.md R5 (depends on T002)
- [X] T004 Implement `Aggregator` — deterministic merge and conflict detection (research.md R6),
      `RenderTools()`, `Dispatch()` with source attribution (research.md R7), `Status()` — in
      `webserver/internal/infra/mcp/aggregator.go` (depends on T002, T003)
- [X] T005 [P] Go test for `Aggregator` merge/conflict-detection/`Dispatch`/attribution logic,
      using fake `ToolSource` test doubles, in `webserver/internal/infra/mcp/aggregator_test.go`
      (depends on T004)
- [X] T006 [P] Go test for read-only filtering — a capability with `readOnlyHint: true` is kept,
      `false` or omitted is dropped — in `webserver/internal/infra/mcp/readonly_test.go` (depends
      on T003)
- [X] T007 Wire `ObservationService` to hold a `*mcp.Aggregator` instead of a static tool list, and
      update `runToolCalls` to delegate to `Aggregator.Dispatch` (capturing the returned source
      name for later attribution) in `webserver/internal/infra/llm/observation.go` (depends on T004)

**Checkpoint**: Aggregator plumbing compiles and is unit-tested against fakes. Nothing is
user-visible yet — US1 makes the first real `ToolSource` available.

---

## Phase 3: User Story 1 - Existing tools become a built-in tool source (Priority: P1) 🎯 MVP

**Goal**: The existing kubectl-backed capability becomes a self-contained, named `LocalToolSource`,
functionally equivalent to today's behavior.

**Independent Test**: With no additional tool sources configured, open "Ask AI about this" on any
object and confirm the assistant performs the same kubectl-backed inspection it does today, now
presented as one named tool source (`GET /api/mcp/sources` — added in Phase 5 — later confirms
this too, but this story's own verification doesn't require it).

### Tests for User Story 1

> **Write this test FIRST, ensure it FAILS before implementation**

- [X] T008 [P] [US1] Go test for `LocalToolSource` — capability listing matches today's
      `KubectlTool()` schema, `Call` passes through success/failure output unchanged, `Health()` is
      always `healthy`, and its one capability carries `readOnlyHint: true` — in
      `webserver/internal/infra/mcp/local_test.go` (per spec.md Clarifications 2026-07-17,
      `LocalToolSource` is backed by a real `*mcp.ClientSession` over `mcp.NewInMemoryTransports()`,
      exercising the actual MCP protocol, not a bespoke adapter)

### Implementation for User Story 1

- [X] T009 [US1] Implement `LocalToolSource` in `webserver/internal/infra/mcp/local.go`: build an
      `mcp.Server` exposing one `kubectl` tool (`readOnlyHint: true`), relocating today's
      `KubectlTool()`/`RunKubectl()` logic out of `webserver/internal/infra/llm/tools.go` into that
      server's tool handler; connect an `*mcp.ClientSession` to it via `mcp.NewInMemoryTransports()`
      at construction time (research.md R2, depends on T002, T008 failing first)
- [X] T010 [US1] Remove the now-dead `RenderTools()`, `KubectlTool()`, `RunKubectl()`, and the
      hardcoded `"kubectl"` switch case in `webserver/internal/infra/llm/tools.go` and
      `observation.go`'s `runToolCalls` (depends on T007, T009)
- [X] T011 [US1] Register the `LocalToolSource` as the aggregator's always-present first source in
      `llm.NewObservationService` (`webserver/internal/infra/llm/observation.go`), so behavior with
      zero external sources configured is unchanged from today (depends on T009, T010)
- [X] T012 [US1] Manual verification: run quickstart.md US1 steps 1-4 (chat still answers via
      kubectl; a missing `kubectl` binary reports unavailable rather than hanging) — no new code,
      confirm via `make run-backend` + the chat panel (depends on T011)

**Checkpoint**: With no external sources configured, User Story 1 is fully functional and matches
today's behavior exactly.

---

## Phase 4: User Story 2 - Add a new tool source without changing code (Priority: P1)

**Goal**: An operator registers an external MCP tool source via a YAML config file and restart; no
Observātiō code change or release is required.

**Independent Test**: Register one additional external tool source pointing at a running, reachable
tool server; confirm the assistant's next answer can draw on a capability that source provides.

### Tests for User Story 2

> **Write these tests FIRST, ensure they FAIL before implementation**

- [X] T013 [P] [US2] Go test for `SourceConfig` YAML parsing and validation — unique names, the
      reserved `kubectl` name rejected, transport-kind-specific required fields, malformed YAML,
      and a missing-but-explicitly-set file path, all per contracts/tool-sources-config.md — in
      `webserver/internal/infra/mcp/config_test.go`
- [X] T014 [P] [US2] Go test for `MCPToolSource` — `tools/list` results translate into
      `Capabilities()`, `Call` translates into `tools/call` and returns its result/error, using an
      in-process fake MCP server (stdio transport, via the go-sdk's own server-side test support) —
      in `webserver/internal/infra/mcp/external_test.go`

### Implementation for User Story 2

- [X] T015 [P] [US2] Implement the `SourceConfig` loader and validation in
      `webserver/internal/infra/mcp/config.go` per contracts/tool-sources-config.md (depends on
      T013 failing first)
- [X] T016 [P] [US2] Implement `MCPToolSource` (stdio and streamable-HTTP transports via
      `github.com/modelcontextprotocol/go-sdk/mcp`, research.md R9) in
      `webserver/internal/infra/mcp/external.go` (depends on T002, T003, T014 failing first)
- [X] T017 [P] [US2] Add the `--tool-sources-config` flag and `TOOL_SOURCES_CONFIG` env-var
      override in `webserver/cmd/server.go`, matching the existing `--address`/`--dev` convention
      (depends on T015)
- [X] T018 [US2] On startup, load the config (T017's resolved path), build an `MCPToolSource` per
      `enabled: true` entry, and register each with the `Aggregator` after the local source —
      treating a config-load validation failure as a fatal startup error per
      contracts/tool-sources-config.md — in `webserver/internal/infra/llm/observation.go` (depends
      on T011, T015, T016, T017)
- [X] T019 [P] [US2] Go test: two fake external sources registered with an overlapping capability
      name produce exactly one `Conflict` entry and only one callable capability, in
      `webserver/internal/infra/mcp/aggregator_test.go` (depends on T004, T016)

**Checkpoint**: User Stories 1 AND 2 both work independently — registering, using, disabling, and
conflicting external sources all behave per quickstart.md US2.

---

## Phase 5: User Story 3 - Aggregate resilience when a source misbehaves or disappears (Priority: P2)

**Goal**: One unreachable/erroring external source doesn't block the others; its health is visible
to an operator; it recovers automatically once reachable again.

**Independent Test**: With two tool sources registered, stop one and confirm the assistant still
answers using the remaining healthy source, while the unhealthy one is visibly flagged.

### Tests for User Story 3

> **Write these tests FIRST, ensure they FAIL before implementation**

- [X] T020 [P] [US3] Go test for health-check state transitions (`unknown` → `healthy` →
      `unhealthy` → `healthy`) driven by a fixed polling interval, in
      `webserver/internal/infra/mcp/health_test.go`
- [X] T021 [P] [US3] Go test confirming `Aggregator.RenderTools()` excludes an unhealthy source's
      capabilities from the callable list while `Aggregator.Status()` still reports that source
      (as unhealthy, not omitted), in `webserver/internal/infra/mcp/aggregator_test.go`
- [X] T022 [P] [US3] Jest test for `ToolSourcesCard` — healthy/unhealthy/unknown badges, a conflict
      indicator when `conflicts` is non-empty, and a source with an empty capability list — in
      `front/app/ui/dashboard/components/ops/tool-sources-card.test.tsx` (written against
      contracts/mcp-sources-api.md's response shape, ahead of the card existing)

### Implementation for User Story 3

- [X] T023 [US3] Implement the background health-check goroutine (30s interval, `tools/list` probe)
      driving `HealthStatus` transitions in `webserver/internal/infra/mcp/health.go`; wire it into
      `MCPToolSource.Health()` (depends on T016, T020 failing first)
- [X] T024 [US3] Update `Aggregator.RenderTools()`/`Dispatch` to recompute the healthy-source
      capability set on every call rather than only at startup, so a recovered source becomes
      usable without an Observātiō restart (depends on T004, T023, T021 failing first)
- [X] T025 [P] [US3] Implement `GET /api/mcp/sources` (`HandleMCPSources`) per
      contracts/mcp-sources-api.md in `webserver/internal/web/handlers/system/mcp_sources.go`,
      calling `Aggregator.Status()` (depends on T004, T023)
- [X] T026 [US3] Register the new route in `webserver/internal/web/handlers/handlers.go`
      (`router.HandleFunc("/api/mcp/sources", system.HandleMCPSources).Methods("GET")`) (depends on
      T025)
- [X] T027 [P] [US3] Add the `use-tool-sources.ts` fetch hook (`SourceStatus`/`HealthStatus` types
      per data-model.md) in `front/app/ui/dashboard/shared/use-tool-sources.ts` (depends on T025's
      response shape / contracts/mcp-sources-api.md)
- [X] T028 [US3] Implement `ToolSourcesCard` in
      `front/app/ui/dashboard/components/ops/tool-sources-card.tsx`, cloned from
      `health-rollup-card.tsx`'s pattern (depends on T022 failing first, T027)
- [X] T029 [US3] Render `ToolSourcesCard` on the Day-2 Ops landing page in
      `front/app/ui/dashboard/components/ops/ops-dashboard.tsx` (depends on T028)

**Checkpoint**: All three user stories are independently functional; quickstart.md US3 passes.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [X] T030 [P] Run quickstart.md's US1-US3 and cross-cutting validation scenarios end-to-end
      against a local test MCP server (stdio and HTTP)
- [X] T031 [P] Verify `make build`, `make run-tests-backend`, and `make run-tests-frontend` all pass
- [X] T032 Flag, for a follow-up `/speckit-constitution` run, that the constitution's Technology
      Stack section should record the new `github.com/modelcontextprotocol/go-sdk` runtime
      dependency (plan.md Constitution Check follow-up note) — not performed automatically by this
      task list
- [X] T033 Annotate this tasks.md with any deviations discovered mid-implementation (project
      convention from features 004-008)

---

### Discovered mid-implementation

- **Mid-implementation clarification changed the local source's design (research.md R2,
  spec.md Clarifications 2026-07-17)**: the original plan had `LocalToolSource` as a plain Go
  interface adapter with no real MCP protocol involved. A `/speckit-clarify` session raised during
  T009 pointed out spec.md's original Input phrasing ("allows the existent tools to become one of
  the local mcp servers") intended the built-in capability to literally speak MCP. `local.go` was
  built as a real `mcp.Server` (via `mcp.AddTool`) connected to its own in-process client over
  `mcp.NewInMemoryTransports()` — no OS subprocess, but a genuine MCP handshake — instead. T008's
  test was written against this corrected design; no rework of already-committed code was needed
  since this was caught before T009's implementation, only after T008's test-writing had begun.
- **T007/T011/T018 required threading a shared `*mcp.Aggregator` through the HTTP call chain, not
  building it inside `llm.NewObservationService`**: inspecting `ws_client.go` mid-T007 revealed
  `registerClient` constructs a brand-new `ObservationService` per WebSocket connection (one per
  browser tab/reconnect). Building the Aggregator there — as tasks.md originally implied — would
  reconnect to every external MCP source (real network/subprocess handshakes) on every chat panel
  open, and leak a health-check goroutine per connection with no shutdown path. Fixed by building
  the Aggregator once in `cmd/server.go`'s `RunE` and threading it as an explicit parameter through
  `handlers.DefaultHandlers` → `startWebSocketHandlers` → `system.HandleChatbot` →
  `registerClient` → `llm.NewObservationService(aggregator)` — mirroring the existing
  `web.WithKubernetes(client, config)` pattern's spirit (shared, request-scoped dependencies) but
  as explicit parameters rather than context values, since the call chain here is a fixed Go
  function chain, not per-request middleware.
- **T018's "config loading + external source registration" ended up in `cmd/server.go` +
  `mcp.BuildExternalSources` (config.go), not `observation.go`**: a direct consequence of the
  point above — `observation.go` no longer builds an Aggregator at all, it only consumes the one
  it's given.
- **T023's `probe`/`recordProbeResult` split**: `MCPToolSource.probe` only refreshes the cached
  capability list and returns an error; a separate `recordProbeResult` maps that error to a
  `HealthStatus`. This split (not originally specified) was needed so the exact same probe logic
  could be reused for both the synchronous first probe in `NewMCPToolSource` (research.md R6 —
  capabilities must be known before Aggregator startup conflict resolution) and the recurring
  background loop in `health.go`, without duplicating the connect/list-tools code.
- **T022/T027 used the existing `useFetchState` + `app/lib/data.tsx` convention instead of a
  bespoke `use-tool-sources.ts` hook file**: `cluster-tabs.tsx`'s existing `getInfraCapabilities`/
  `useFetchState` pattern already does exactly what plan.md's proposed hook would have — a fetch
  function in `app/lib/data.tsx` plus the shared `useFetchState` hook — so `getMCPSources` was
  added there instead of creating a parallel, redundant hook abstraction.
- **T012's live verification was partial**: this environment has no `ANTHROPIC_API_KEY`, so the
  full Claude-mediated kubectl tool-call round trip (US1's actual chat flow) could not be exercised
  live. Substituted: (a) `local_test.go`'s `TestLocalToolSource_Call_*` tests exercise the real
  in-process MCP protocol round trip against the real `kubectl` binary — the same path a live chat
  turn would use, minus Claude's tool-selection step; (b) a manual WebSocket smoke test against a
  real `kind-capi-mgmt` cluster confirmed the full server startup (including the real local MCP
  handshake) and `/ws/analysis` request path work end-to-end up to the Anthropic API boundary,
  where it fails gracefully with the existing "AI assistant is not available" message rather than
  crashing.
- **Frontend test environment note (not a code deviation)**: this sandbox's default `node` (system
  package, v18.19.1) cannot load `jest.config.js` (`import` syntax) — the project's CI targets Node
  22 (`.github/workflows/build.yml`). A newer Node was available at `/snap/node/11776/bin` (v24) and
  was used instead to run `pnpm run test` and `make build` successfully. This is a pre-existing
  environment gap unrelated to this feature (confirmed via `git stash` — the same failure occurs on
  the unmodified branch under the system Node).
- **`tool-sources-card.test.tsx` needed `findAllByText`/length assertions, not `findByText`, for
  the realistic "source name equals its one capability's name" case** (`kubectl`/`kubectl`) — both
  render as separate text nodes, so a single-match query is ambiguous. Fixed in the two affected
  test cases; no component code change needed.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **Foundational (Phase 2)**: Depends on Setup — BLOCKS all user stories (the `Aggregator`/
  `ToolSource` contract is a genuine shared dependency; every story's `ToolSource` implementation
  plugs into it)
- **User Stories (Phase 3-5)**: All depend on Foundational (Phase 2) completion
  - US1 depends only on Foundational — it is the first concrete `ToolSource`
  - US2 depends only on Foundational for its own config/external-source logic (T013-T017, T019);
    T018 additionally depends on US1's T011 (the local-source registration point in
    `observation.go`) since it extends the same startup wiring
  - US3 depends on US2's `MCPToolSource` (T016) to have something whose health it checks (T023),
    and on Foundational's `Aggregator` (T024); its status endpoint (T025-T029) has no story
    dependency beyond Foundational's `Status()` method (T004)
- **Polish (Phase 6)**: Depends on all desired user stories being complete

### User Story Dependencies

- **US1 (P1)**: Depends only on Foundational. Fully independent — this is the MVP.
- **US2 (P1)**: Depends on Foundational directly; T018 also depends on US1's T011 since both edit
  the same startup registration sequence in `observation.go`.
- **US3 (P2)**: T023 extends US2's `MCPToolSource` (T016) with health checking, so it is sequenced
  after US2 in practice, even though nothing in spec.md makes US3 strictly blocked on US2's every
  task; T025-T029 (status endpoint + card) depend only on Foundational's `Aggregator.Status()`
  (T004) and T023.

### Parallel Opportunities

- T002 (types) can start immediately after T001 (Setup).
- Within Foundational, T005 and T006 (both tests) can run in parallel once their respective
  implementation tasks (T004, T003) land.
- Within each user story, `[P]`-marked test tasks run in parallel with each other (e.g. T013/T014,
  T020/T021/T022).
- US1's and US2's test-writing (T008 and T013/T014) touch entirely different files and can happen
  in parallel once Foundational is done.
- T015 (config.go) and T016 (external.go) are independent files and can be implemented in parallel
  once their respective tests (T013, T014) exist.

---

## Parallel Example: User Story 2

```bash
# Tests (parallel):
Task: "Go test for SourceConfig YAML parsing in webserver/internal/infra/mcp/config_test.go"
Task: "Go test for MCPToolSource against a fake stdio MCP server in webserver/internal/infra/mcp/external_test.go"

# Implementation (parallel, once tests above exist):
Task: "Implement SourceConfig loader in webserver/internal/infra/mcp/config.go"
Task: "Implement MCPToolSource in webserver/internal/infra/mcp/external.go"
Task: "Add --tool-sources-config flag in webserver/cmd/server.go"
```

---

## Implementation Strategy

### MVP First (Setup + Foundational + User Story 1)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL — the `ToolSource`/`Aggregator` contract nothing else
   works without)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: confirm against quickstart.md's US1 scenarios — behavior matches today's
   kubectl-only assistant exactly, now presented as one named source
5. Deploy/demo if ready

### Incremental Delivery

1. Setup + Foundational → ready
2. US1 (built-in tool source) → MVP, zero behavior change from today
3. US2 (register external sources via config) → the feature's flagship value: "aggregating new MCP
   server" becomes real
4. US3 (health/resilience) → makes multi-source aggregation trustworthy, not just possible
5. Each story is additive; US2 and US3 never change US1's behavior when no external sources are
   configured

### Parallel Team Strategy

Foundational must land first — it's a real shared blocker (the `Aggregator` every `ToolSource`
plugs into). After that: one developer on US1 (`local.go` + `llm` cleanup), one on US2
(`config.go` + `external.go` + `cmd/server.go`) — different files, integrating only at T018's
registration call in `observation.go`. US3 should follow US2 since its health checker directly
extends `MCPToolSource`.
