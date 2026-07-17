# Data Model: MCP Server Aggregation & Local Tool Server

New types live in a new package, `webserver/internal/infra/mcp` (research.md R2) — the aggregator's
own domain model, independent of both the Anthropic SDK's tool-schema types and the upstream MCP
SDK's wire types. `webserver/internal/infra/llm` depends on this package instead of hardcoding a
single tool.

## ToolSource (interface)

The common shape both the built-in and external sources implement.

```go
type ToolSource interface {
    Name() string
    Kind() SourceKind // "local" | "external"
    Capabilities() []Capability
    Call(ctx context.Context, capability string, args json.RawMessage) (result string, isError bool, err error)
    Health() HealthStatus
}
```

## Capability

One named, read-only action contributed by exactly one `ToolSource` (spec.md Key Entities).

| Field | Type | Notes |
|---|---|---|
| `Name` | `string` | Model-facing tool name; unique across all healthy sources post-conflict-resolution (research.md R6) |
| `Description` | `string` | Passed through to the Claude tool schema unchanged |
| `InputSchema` | `map[string]interface{}` | JSON-schema fragment, same shape `KubectlTool()` already builds today |
| `SourceName` | `string` | Which `ToolSource` contributed it — carried for FR-008/SC-003 traceability (research.md R7), not shown to the model |
| `ReadOnly` | `bool` | Always `true` for any capability that survives aggregation (research.md R5) — capabilities that fail this check never reach a `Capability` value at all |

## HealthStatus

| Field | Type | Notes |
|---|---|---|
| `State` | `string` | One of `healthy`, `unhealthy`, `unknown` (research.md R4) |
| `LastChecked` | `time.Time` | Zero value when `State == "unknown"` |
| `LastError` | `string` | Empty unless `State == "unhealthy"` |

## SourceConfig (YAML, operator-authored)

Decoded from the file at `--tool-sources-config` / `TOOL_SOURCES_CONFIG` (research.md R3). Only
external sources are declared here — the built-in local source is always present and is not
config-driven.

```yaml
sources:
  - name: velero-mcp
    enabled: true
    transport:
      kind: stdio        # or "http"
      command: velero-mcp # stdio only
      args: []            # stdio only
      url: ""              # http only
```

| Field | Type | Notes |
|---|---|---|
| `Name` | `string` | Must be non-empty and unique among configured sources; duplicate names fail config load at startup (distinct from the *capability*-name conflict in research.md R6, which is resolved at runtime, not startup) |
| `Enabled` | `bool` | `false` entries are parsed but never instantiated as a `ToolSource` — lets an operator keep a source's config around without removing it |
| `Transport.Kind` | `string` | `stdio` or `http` (research.md R9) |
| `Transport.Command` / `Args` | `string` / `[]string` | `stdio` only |
| `Transport.URL` | `string` | `http` only |

## Aggregator

Not persisted — an in-memory object owned by `ObservationService`, built once at startup from the
local source plus every `enabled: true` entry in `SourceConfig`.

| Field | Type | Notes |
|---|---|---|
| `sources` | `[]ToolSource` | Local source first, then external sources in config-file order (research.md R6) |
| `capabilityIndex` | `map[string]Capability` | The winning capability per name after conflict resolution |
| `conflicts` | `[]Conflict` | Recorded, not fatal (research.md R6) |

### Conflict

| Field | Type | Notes |
|---|---|---|
| `CapabilityName` | `string` | The colliding name |
| `WinningSource` | `string` | The source whose capability was kept (first-registered) |
| `RejectedSource` | `string` | The source whose same-named capability was excluded |

## Aggregator methods used by `webserver/internal/infra/llm`

```go
func (a *Aggregator) RenderTools() []anthropic.ToolUnionParam   // replaces llm.RenderTools()
func (a *Aggregator) Dispatch(ctx context.Context, name string, args json.RawMessage) (result string, isError bool, sourceName string, err error)
func (a *Aggregator) Status() []SourceStatus                     // for GET /api/mcp/sources (research.md R8)
```

`Dispatch`'s returned `sourceName` is threaded into the existing `runToolCalls`/`StreamChatWithAgent`
flow in `observation.go` so a chat response can be attributed to its source (research.md R7) without
changing the model-facing tool schema.

### SourceStatus (REST response shape — see `contracts/mcp-sources-api.md`)

| Field | Type | Notes |
|---|---|---|
| `Name` | `string` | |
| `Kind` | `string` | `local` \| `external` |
| `Health` | `HealthStatus` | |
| `Capabilities` | `[]string` | Names only — full schema stays internal to the tool-calling loop |

## Frontend types (`front/app/ui/dashboard/shared/use-tool-sources.ts`, new)

Mirrors the backend JSON shape exactly, following the existing `use-day2-ops.ts` convention:

```ts
export type SourceKind = 'local' | 'external';
export type HealthState = 'healthy' | 'unhealthy' | 'unknown';

export type HealthStatus = {
  state: HealthState;
  lastChecked?: string; // ISO timestamp, omitted when state is "unknown"
  lastError?: string;
};

export type SourceStatus = {
  name: string;
  kind: SourceKind;
  health: HealthStatus;
  capabilities: string[];
};
```
