# Contract: Tool Source status — `GET /api/mcp/sources`

New read-only REST endpoint (research.md R8), pattern-matched to the existing
`GET /api/infra/capabilities` handler. No request body, no query parameters. Satisfies FR-003 and
SC-005.

## Response — `200 OK`

```jsonc
{
  "sources": [
    {
      "name": "kubectl",
      "kind": "local",
      "health": {
        "state": "healthy"
      },
      "capabilities": ["kubectl"]
    },
    {
      "name": "velero-mcp",
      "kind": "external",
      "health": {
        "state": "unhealthy",
        "lastChecked": "2026-07-17T14:02:11Z",
        "lastError": "dial tcp: connection refused"
      },
      "capabilities": ["list_backups", "describe_backup"]
    }
  ],
  "conflicts": [
    {
      "capabilityName": "list_backups",
      "winningSource": "velero-mcp",
      "rejectedSource": "velero-mcp-mirror"
    }
  ]
}
```

- `sources` always includes exactly one `kind: "local"` entry (the built-in kubectl capability,
  FR-001) plus one entry per `enabled: true` source in the operator's config, regardless of current
  health — an unhealthy source is still listed (`health.state: "unhealthy"`), never omitted
  (US3 AC2).
- `sources[].capabilities` lists only capability names that survived read-only verification
  (research.md R5) and conflict resolution (research.md R6) for that source; a source that lost
  every capability to conflicts can still appear with an empty `capabilities` array.
- `conflicts` is empty (`[]`, never omitted) when no capability-name collisions were detected at
  startup.
- `health.lastChecked` is present whenever the source has been probed at least once (healthy or
  unhealthy) and omitted only while it's still `"unknown"` (never probed — e.g. right at startup,
  before an external source's first health check completes). `health.lastError` is omitted unless
  the most recent probe failed — mirroring the existing "omit rather than null-out" convention
  used by 008's `recoveryInfo`.

## Consumer contract

- Frontend `use-tool-sources.ts` fetches this once on mount (no polling loop, no WebSocket
  subscription — research.md R8); a manual refresh (e.g. re-opening the status view) re-fetches.
- A new status card (cloned from the `ops/` dashboard card pattern, e.g.
  `front/app/ui/dashboard/components/ops/tool-sources-card.tsx`) renders one row per `sources[]`
  entry with a health badge, and a distinct "conflict" indicator when `conflicts` is non-empty for
  a capability that entry would otherwise have exposed.
- This endpoint is read-only and carries no way to register, disable, or remove a source — that
  remains config-file-only (FR-010, research.md R3). It exists purely so an operator can *observe*
  what the config produced without reading server logs.
