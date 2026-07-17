# Contract: Tool source registration file (`--tool-sources-config` / `TOOL_SOURCES_CONFIG`)

The operator-facing interface for FR-002/FR-004/FR-010 (research.md R3). A YAML file, read once at
process startup. Absence of the flag/env var is valid — the assistant runs with only the built-in
local source, matching today's behavior exactly (spec.md US1 Independent Test).

## Schema

```yaml
sources:
  - name: velero-mcp          # required, unique, non-empty
    enabled: true               # required
    transport:
      kind: stdio                # required: "stdio" | "http"
      command: velero-mcp        # required when kind: stdio
      args: ["--read-only"]      # optional, stdio only
      url: ""                     # required when kind: http; ignored for stdio
```

## Validation (fail-fast at startup, before the server starts serving traffic)

- `sources[].name` must be non-empty and unique across the file. A duplicate name is a startup
  error (distinct from the runtime capability-name conflict in `contracts/mcp-sources-api.md`,
  which degrades one source rather than failing startup — a duplicate *source* name is an operator
  typo, not a legitimate two-sources-overlap scenario).
- `sources[].name` must not equal the reserved built-in source name (`kubectl`) — that name is
  always the local source.
- `transport.kind` must be `stdio` or `http`; any other value is a startup error.
- `stdio` entries require `command`; `http` entries require `url`. A missing required field for
  the declared kind is a startup error.
- Malformed YAML, or a file path that doesn't exist when the flag/env var is explicitly set, is a
  startup error — Observātiō does not silently start with zero external sources when the operator
  clearly intended to configure some (distinct from "flag/env var omitted entirely," which is the
  valid zero-external-sources case above).

## Operator workflow (FR-002, FR-004, SC-001)

1. Add or edit an entry in the YAML file.
2. Restart the Observātiō process (`--tool-sources-config` is read once at startup, research.md
   R3) — no code change, no new build or release of Observātiō itself.
3. Confirm via `GET /api/mcp/sources` (`contracts/mcp-sources-api.md`) that the source appears
   with the expected health and capability list.

Disabling a source: set `enabled: false` (or delete its entry) and restart — after restart the
assistant no longer offers or uses that source's capabilities (FR-004).
