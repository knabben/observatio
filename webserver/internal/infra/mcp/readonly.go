package mcp

import (
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// IsReadOnly reports whether an MCP tool declares itself non-mutating via its readOnlyHint
// annotation. A missing annotations block, or an explicit readOnlyHint: false, is treated as NOT
// read-only — fail closed, never fail open (research.md R5, FR-009, SC-004).
//
// This trusts the declaring server's own annotation; it is a protocol-level allowlist, not a
// sandboxed runtime guarantee (research.md R5's documented "Known limitation").
func IsReadOnly(t *mcpsdk.Tool) bool {
	return t != nil && t.Annotations != nil && t.Annotations.ReadOnlyHint
}

// translateTools converts an external source's MCP tool list into this package's Capability
// type, dropping every tool that fails IsReadOnly before it ever reaches the aggregator.
func translateTools(sourceName string, tools []*mcpsdk.Tool) []Capability {
	var caps []Capability
	for _, t := range tools {
		if !IsReadOnly(t) {
			continue
		}
		schema, _ := t.InputSchema.(map[string]any)
		caps = append(caps, Capability{
			Name:        t.Name,
			Description: t.Description,
			InputSchema: schema,
			SourceName:  sourceName,
			ReadOnly:    true,
		})
	}
	return caps
}
