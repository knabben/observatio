# Quickstart: MCP Server Aggregation & Local Tool Server

Manual verification scenarios, one per user story, plus cross-cutting checks. No external MCP
server is required for US1; US2/US3 assume a reachable local test MCP server (a minimal stdio
"echo" tool server is sufficient to exercise registration/health without standing up `velero-mcp`).

## US1 — Existing tools become a built-in tool source

1. Start Observātiō with no `--tool-sources-config` / `TOOL_SOURCES_CONFIG` set.
2. Open "Ask AI about this" on any object and ask a question that requires inspecting live cluster
   state; confirm the assistant answers using the kubectl-backed capability exactly as it does
   today.
3. Call `GET /api/mcp/sources`; confirm exactly one entry appears, `kind: "local"`,
   `health.state: "healthy"`, with the kubectl capability listed.
4. Temporarily rename/remove the `kubectl` binary from `PATH`, ask the assistant a question that
   needs it; confirm it reports the capability is unavailable (a clear tool-result error) rather
   than hanging or crashing the chat turn.

## US2 — Add a new tool source without changing code

1. Stand up a reachable local test MCP server (stdio or HTTP) exposing at least one
   `readOnlyHint: true` tool.
2. Add it to the tool sources YAML config, set `TOOL_SOURCES_CONFIG` to that file, and restart
   Observātiō — no code change, no rebuild of the frontend or backend beyond the restart itself.
3. Call `GET /api/mcp/sources`; confirm the new source appears as `kind: "external"`,
   `health.state: "healthy"`, with its capability listed.
4. Ask the assistant a question only the new source's capability can answer; confirm the response
   is backed by that tool's output.
5. Register a second, distinct test source at the same time; ask a question spanning both;
   confirm the assistant draws on both within the same conversation.
6. Disable the second source (`enabled: false`) and restart; confirm the assistant no longer
   offers or uses its capability, and it disappears from `GET /api/mcp/sources`.
7. Register a second source that deliberately exposes a capability with the same name as the
   first; restart; confirm `GET /api/mcp/sources`'s `conflicts` array reports it, and only one of
   the two capabilities is actually callable.

## US3 — Aggregate resilience when a source misbehaves or disappears

1. With two external sources registered and both reachable, stop one (kill the subprocess / take
   the HTTP endpoint down).
2. Ask the assistant a question answerable only by the remaining healthy source; confirm it
   answers normally, without waiting on or being blocked by the unhealthy one.
3. Call `GET /api/mcp/sources`; confirm the stopped source shows `health.state: "unhealthy"` with
   a `lastError`, while the other source and the local source show `healthy`.
4. Restart the stopped source's process/endpoint without touching Observātiō's config or
   restarting Observātiō itself; wait past one health-check interval; confirm
   `GET /api/mcp/sources` shows it `healthy` again and the assistant can use its capabilities
   again, with no re-registration step.
5. Stop every registered source (including making kubectl unavailable, per US1 step 4); ask the
   assistant a question; confirm it reports it currently has no working capabilities rather than
   appearing to hang.

## Cross-cutting checks

- Confirm every capability listed anywhere in `GET /api/mcp/sources` corresponds to a tool the
  declaring source annotated `readOnlyHint: true`; register a test source with a tool that omits
  or sets `readOnlyHint: false` and confirm it never appears as a callable capability (SC-004).
- Confirm every chat response that used a tool can be traced to the source(s) it came from (check
  the chat UI's rendered attribution, or the WS frames in devtools) — SC-003.
- Confirm registering, disabling, or removing a source is only possible by editing the config file
  and restarting — there is no button, API, or chat-panel action that does it (FR-010).
- Run `make run-tests-backend` and `make run-tests-frontend`; confirm both pass, including the new
  `webserver/internal/infra/mcp` package tests (conflict detection, health-state transitions,
  read-only filtering, dispatch/attribution) and the new tool-sources status card test.
- Run `make build`; confirm it succeeds with the new `github.com/modelcontextprotocol/go-sdk`
  dependency vendored/resolved correctly.
