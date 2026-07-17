// Package mcp aggregates the assistant's tool sources — the built-in kubectl capability and any
// operator-registered external MCP servers — behind one unified, read-only capability set (see
// specs/009-mcp-server-client-aggregator).
package mcp

import (
	"context"
	"encoding/json"
	"time"
)

// SourceKind distinguishes the always-present built-in source from operator-registered ones.
type SourceKind string

const (
	SourceKindLocal    SourceKind = "local"
	SourceKindExternal SourceKind = "external"
)

// HealthState is one of a ToolSource's three possible health states (research.md R4).
type HealthState string

const (
	HealthUnknown   HealthState = "unknown"
	HealthHealthy   HealthState = "healthy"
	HealthUnhealthy HealthState = "unhealthy"
)

// HealthStatus is a ToolSource's current reachability, as reported to operators via
// GET /api/mcp/sources (contracts/mcp-sources-api.md).
type HealthStatus struct {
	State       HealthState `json:"state"`
	LastChecked time.Time   `json:"lastChecked,omitzero"`
	LastError   string      `json:"lastError,omitempty"`
}

// Capability is one named, read-only action contributed by exactly one ToolSource.
type Capability struct {
	Name        string
	Description string
	InputSchema map[string]interface{}
	// SourceName identifies which ToolSource contributed this capability — carried for
	// FR-008/SC-003 traceability (research.md R7), never exposed in the model-facing tool schema.
	SourceName string
	// ReadOnly is always true for any Capability that survives aggregation (research.md R5); a
	// capability that fails read-only verification never becomes a Capability value at all.
	ReadOnly bool
}

// ToolSource is the common shape both the built-in local capability and every operator-registered
// external MCP server implement, so the Aggregator can treat them identically (research.md R2).
type ToolSource interface {
	Name() string
	Kind() SourceKind

	// Capabilities returns the source's current, cached capability list. It must not block on
	// network I/O — external sources refresh this snapshot out-of-band via their health check
	// (research.md R4), so a slow or unreachable source never adds latency to a chat turn.
	Capabilities() []Capability

	// Call invokes the named capability with the given raw JSON arguments, returning the tool
	// result text and whether it represents a tool-level error (as opposed to a Go error, which
	// signals a protocol/transport failure).
	Call(ctx context.Context, capability string, args json.RawMessage) (result string, isError bool, err error)

	Health() HealthStatus
}
