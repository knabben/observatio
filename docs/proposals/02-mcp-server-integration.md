# Quick specification statement: MCP server integration

Observātiō already has one AI capability: a per-object "Ask AI about this" panel (feature 005,
`webserver/internal/infra/llm/agent.go` + `tools.go`) built on a single Claude agent with exactly one
tool — `RunKubectl`, which shells out to the local `kubectl` binary with whatever command string the
model produces. There's no Model Context Protocol integration in either direction today.

Two distinct additions are being considered under this heading, and they should be scoped and decided
separately rather than bundled as one feature:

1. **Observātiō as an MCP client.** Add the community `velero-mcp` server (referenced by the companion
   guide this analysis is based on) as a tool source for the existing AI agent, alongside or in place of
   the ad hoc `kubectl` tool. This gives the current per-object AI panel real backup/restore verbs — list
   backups, describe a backup, check `BackupStorageLocation` health, kick off a restore — without
   hand-rolling each one as a bespoke Go tool the way `KubectlTool()` is today.
2. **Observātiō as an MCP server.** Expose Observātiō's own aggregated Day-2 Ops state (rollups, debug
   paths, risk warnings, severity, and — once proposal 01 lands — backup health) as MCP resources/tools,
   so external agents can query CAPI+Velero-specific context directly from Observātiō instead of
   re-deriving it themselves. This is the concrete shape of "a replacement for kagent" for this domain:
   rather than running kagent's generic Kubernetes agent against a management cluster, other tooling
   (including kagent itself, or Claude Desktop) points at Observātiō's own domain-specific MCP endpoint
   and gets CAPI-aware, already-correlated answers.

Planning should determine which of these two (or both, and in what order) is worth building first, and
should account for the fact that `RunKubectl`'s current unguarded command execution is a pre-existing gap
worth addressing regardless of MCP — an MCP-based tool surface will need its own read/write boundary
given this product's read-only design.
