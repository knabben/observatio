# Research: MCP Server Aggregation & Local Tool Server

## R1: Adopt the official `modelcontextprotocol/go-sdk` as the MCP client library

**Decision**: Add `github.com/modelcontextprotocol/go-sdk` (package path
`github.com/modelcontextprotocol/go-sdk/mcp`) as the MCP client implementation used to talk to
externally-registered tool sources.

**Rationale**: No MCP client exists anywhere in this repo today (confirmed: zero hits for
`modelcontextprotocol|mcp-go|mark3labs` across the whole tree, and `webserver/go.mod` has no MCP
dependency). The protocol itself — JSON-RPC 2.0 framing, capability negotiation
(`initialize`/`initialized`), `tools/list`, `tools/call`, and both `stdio` and streamable-HTTP
transports — is meaningful surface area with edge cases (cancellation, progress notifications,
malformed-server handling) that this feature has no reason to hand-roll. The official SDK is
maintained by the same stewards who define the spec, so it tracks protocol changes directly rather
than through a third party's interpretation of them.

**Alternatives considered**:
- `github.com/mark3labs/mcp-go` — a mature, widely-used community client. Rejected as the primary
  choice only because it is a third-party implementation of a spec still evolving under a
  different steward; the official SDK is the safer long-term bet for a dependency this feature
  will keep for its lifetime. Worth revisiting if the official SDK proves to lag the spec in
  practice.
- Hand-rolled JSON-RPC-over-stdio/HTTP client — rejected: reimplements transport framing and
  capability negotiation the SDK already provides, for zero product-differentiating benefit, and
  becomes an ongoing maintenance burden every time the spec adds a feature.

## R2: New `webserver/internal/infra/mcp` package; the built-in kubectl tool is a real, in-process MCP server

**Decision** (revised 2026-07-17 — see spec.md Clarifications): Introduce a small `ToolSource`
interface (list capabilities, call a capability by name, report health) in a new package,
`webserver/internal/infra/mcp`. The existing `kubectl`-backed capability
(`webserver/internal/infra/llm/tools.go`) is wrapped by a real `mcp.Server` (from R1's SDK,
`mcp.NewServer` + `mcp.AddTool`) exposing one `kubectl` tool; `LocalToolSource` connects to that
server via `mcp.NewInMemoryTransports()` — an in-process pipe pair, no OS subprocess — and consumes
it through the exact same `*mcp.ClientSession` machinery an external source uses. External sources
are backed by an `MCPToolSource` that holds a live `*mcp.ClientSession` over a real `stdio` or HTTP
transport and translates `ListTools`/`CallTool` into the same `ToolSource` interface; the local
source's `mcp.ClientSession` differs only in which transport it was connected over.

**Rationale**: spec.md's original Input phrasing — "allows the existent tools to become one of the
local mcp servers" — and the Clarifications session confirm the built-in capability should
literally speak the MCP protocol (`initialize`, `tools/list`, `tools/call`), not merely be adapted
to look like it does from the aggregator's point of view. `mcp.NewInMemoryTransports()` delivers
that without the downside an actual subprocess would add (process lifecycle management, a new
"unhealthy" failure mode for a capability that has no external dependency today) — connection setup
is in-memory channels, not a spawned OS process, so it's exactly as reliable as calling a Go
function directly. This also means `local.go`'s `mcp.Server` is, structurally, reusable by other
MCP clients later (spec.md Assumptions) even though exposing it externally remains out of scope for
this feature — and it eliminates a second, parallel `ToolSource` implementation: `local.go` and
`external.go` both ultimately wrap an `*mcp.ClientSession`, differing only in transport
construction.

**Alternatives considered**:
- A plain in-process Go interface with no MCP protocol involvement (the original R2 decision) —
  reconsidered per the Clarifications session: it satisfied FR-001's literal wording but not the
  spec's original intent, and forced two structurally different `ToolSource` implementations
  (local vs. external) instead of one.
- Run the built-in kubectl tool as an actual local MCP **subprocess** server (stdio, a second OS
  process) — rejected: adds process-management complexity and a new failure mode for zero
  additional protocol fidelity over the in-memory-transport approach; `mcp.NewInMemoryTransports()`
  gets the same "real MCP server" property without spawning anything.

## R3: Tool sources are registered via a YAML config file loaded once at startup

**Decision**: External tool sources are declared in a YAML file whose path comes from a new
`--tool-sources-config` flag / `TOOL_SOURCES_CONFIG` env var (mirroring the existing
flag+env-override convention already used for `--address`/`--dev` and `ANTHROPIC_MODEL` in
`webserver/cmd/server.go` and `observation.go`). Each entry declares `name`, `transport` (`stdio`
with a command + args, or `http` with a URL), and `enabled`. The file is read once at process
startup; there is no REST or WebSocket endpoint to add, edit, or remove a source at runtime.

**Rationale**: FR-010 explicitly restricts registering/modifying/removing a tool source to
"administrative configuration of Observātiō," not an end-user or chat-panel action, and spec.md's
Assumptions section confirms registration is "an administrative/deployment-time concern
(config-driven), not a per-conversation or per-user runtime action." A config file edited before a
restart satisfies FR-002's "without requiring a code change or new release" — a restart is neither.
This is the first file-based config loader in the repo (confirmed: no `yaml.Unmarshal`,
`yaml.NewDecoder`, or `viper` usage anywhere in `webserver/`); it is flagged in plan.md's
Complexity Tracking per the constitution's "no new runtime dependency without justification" gate
(the YAML parser itself — `sigs.k8s.io/yaml`, already an indirect dependency — becomes direct).

**Alternatives considered**:
- Per-source environment variables (`TOOL_SOURCE_1_NAME=...`, `TOOL_SOURCE_1_URL=...`) — rejected:
  doesn't scale past one or two sources and can't cleanly express nested per-transport fields.
- A database-backed runtime registration API — rejected outright by FR-010.

## R4: Background health polling per external source; local source is always healthy

**Decision**: Each external `MCPToolSource` runs a background goroutine that calls `tools/list`
against its session on a fixed interval (30s) after the initial connection, tracking one of three
states: `healthy`, `unhealthy`, or `unknown` (before the first check completes). A failing check
excludes that source from the next capability set the aggregator hands to Claude, but the source
stays registered — no operator action is needed for it to reappear once a check succeeds again.
The built-in `LocalToolSource` is always reported `healthy` (no network dependency; kubectl
failures surface per-call the same way they do today).

**Rationale**: Directly implements US3: "the assistant [keeps] working with whatever sources are
still healthy" (AC1), health is visible to an operator (AC2), and a recovered source becomes
available again "without requiring the operator to re-register it" (AC3). Polling out-of-band
rather than inline keeps a slow/unreachable source from adding latency to every "Ask AI" turn.

**Alternatives considered**:
- Health-check on every chat turn before building the tool list — rejected: one slow source would
  add latency to every turn, including ones that don't need that source.
- No health tracking, just let `tools/call` fail when it fails — rejected: doesn't satisfy FR-003
  ("view... each source's health/reachability") or US3 AC2's explicit "shown as unhealthy."

## R5: Read-only enforcement is fail-closed on MCP tool annotations

**Decision**: A capability is only added to the aggregated tool list if its declaring source marks
it `readOnlyHint: true` in the MCP tool's `annotations` (the standard MCP mechanism for a server to
describe its own tool's side effects). The built-in `kubectl` capability is hardcoded with the
equivalent guarantee, matching its existing (implicit) behavior. Any capability without an explicit
`readOnlyHint: true` — including one that omits annotations entirely — is dropped before it ever
reaches Claude. This is fail-closed: ambiguity excludes a capability, it never includes one.

**Rationale**: Directly satisfies FR-009 ("MUST NOT execute a capability that would change cluster
or backup state") and SC-004 ("100% of capabilities... verified non-mutating before being made
callable"), using the standard field the protocol already gives servers for declaring this rather
than inventing a parallel mechanism.

**Known limitation** (worth surfacing to operators, not solving here): this trusts the annotation
as declared by the source's own server — it is a protocol-level allowlist, not a sandboxed runtime
guarantee. A misconfigured or malicious external server could mis-declare a mutating tool as
read-only. Registering an external source is already an administrative action gated by FR-010; this
plan treats "an operator chose to trust this source's declared annotations" as an acceptable
consequence of that administrative act, not a gap to close in this iteration.

**Alternatives considered**:
- Trust every capability from any registered source — rejected outright by FR-009.
- A second, hand-maintained per-tool allowlist in the YAML config, independent of server
  annotations — rejected as redundant operator burden that creates a second list to keep in sync
  with the first; the protocol's own annotation field already exists for exactly this purpose.

## R6: Deterministic, config-order conflict detection — reject, don't merge

**Decision**: The aggregator builds its capability-name → source map once at startup, walking
sources in a fixed order (built-in local source first, then external sources in config-file
order). If a later source declares a capability name already claimed by an earlier source, the
later one is excluded from the tool list, and the conflict is recorded as a structured entry
(surfaced via logs and the status endpoint from R8) — never silently merged or overwritten.

**Rationale**: FR-007 requires conflicts to be "surface[d]... to the operator at registration time
rather than silently picking one or merging them," and FR-005 requires the tool list to "stay
unambiguous." A conflict degrades one source's one capability, not the whole aggregator.

**Alternatives considered**:
- Last-registered-wins silent overwrite — explicitly ruled out by FR-007.
- Refuse to start the process on any conflict — rejected as disproportionate; US3's whole premise
  is that one misbehaving/overlapping source shouldn't take the rest down.

## R7: Source attribution travels out-of-band, not as a name prefix

**Decision**: Capability names presented to Claude stay exactly as declared (no
`sourcename__toolname` prefixing). The aggregator's internal dispatch table tags every
`tools/call` result with the source it came from, and that attribution is threaded back through
`StreamChatWithAgent` so each chat response can be traced to its source(s) — satisfying FR-008/
SC-003 — without leaking wiring details into the model's own tool-selection reasoning.

**Rationale**: FR-005 explicitly frames the union of capabilities as one set an operator/model
"should not need to know which source backs" to use; prefixing names would contradict that by
making source identity part of the tool schema itself. Attribution only needs to reach the human
reading the final answer, not the model choosing which tool to call.

**Alternatives considered**:
- Prefix every tool name with its source (`velero_mcp__list_backups`) — rejected as the primary
  mechanism for the reason above; still available later as a targeted fallback if two sources ever
  need simultaneous disambiguation mid-conversation, which no current acceptance scenario requires.

## R8: Status visibility is a polled REST endpoint, not a new WebSocket push

**Decision**: Add `GET /api/mcp/sources`, pattern-matched to the existing
`/api/infra/capabilities` handler (`webserver/internal/web/handlers/kubernetes/cluster.go`):
simple `router.HandleFunc` registration, JSON response, no new middleware. It returns each
registered source's name, kind (`local`/`external`), health, and current capability list. The new
frontend status view (cloned from the `ops/` dashboard card pattern) fetches this on mount rather
than subscribing over the existing `/ws/analysis` or `/ws/watcher` sockets.

**Rationale**: Constitution Principle II mandates WebSocket delivery for *cluster state*
specifically — its own bullet list scopes this to "Watchers for Cluster, Machine, and
MachineDeployment." Tool-source health is assistant-tooling/configuration state, not CAPI cluster
state, so it falls outside that principle's literal scope (documented explicitly in this plan's
Constitution Check rather than left as a silent gap). `/api/infra/capabilities` is direct
precedent for exactly this kind of "what's installed/available" REST read in this codebase.

**Alternatives considered**:
- Push tool-source health over a new or existing WebSocket type — rejected: adds a new live
  channel and payload type for data that changes on the order of minutes (a health check interval),
  not seconds; a poll-on-view is simpler and has clear precedent.

## R9: Support both `stdio` and streamable-HTTP transports for external sources

**Decision**: v1 supports registering an external source over `stdio` (Observātiō spawns a local
subprocess, e.g. a locally-installed `velero-mcp` binary) or streamable HTTP (a remotely-hosted MCP
server). Legacy HTTP+SSE transport is not implemented, since the MCP spec has deprecated it in
favor of streamable HTTP.

**Rationale**: Both transports are first-class in the MCP spec and supported at comparable
implementation cost by the R1 SDK. Restricting to `stdio` only would preclude any remotely-hosted
tool source, which contradicts "aggregating new MCP server" (spec.md Input) as a general-purpose
capability rather than a one-off integration for a single local binary.

**Alternatives considered**:
- `stdio`-only for v1 — rejected: `velero-mcp` (spec.md Assumptions) runs as a workload inside the
  management cluster and is reached over HTTP (a `ClusterIP` Service URL), not spawned as a local
  subprocess by Observātiō — confirming HTTP, not `stdio`, is the transport that matters for the
  concrete real-world example this feature targets. Limiting to `stdio` would make "register any
  external tool source" (FR-002) false for that case, and for any `kmcp`-hosted server reached the
  same way.
