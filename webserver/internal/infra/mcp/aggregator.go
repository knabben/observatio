package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	"github.com/anthropics/anthropic-sdk-go"
)

// Conflict records a capability name claimed by more than one registered source. The
// first-registered source keeps the name; every later collision is recorded here instead of
// being silently merged or overwritten (research.md R6, FR-007).
type Conflict struct {
	CapabilityName string `json:"capabilityName"`
	WinningSource  string `json:"winningSource"`
	RejectedSource string `json:"rejectedSource"`
}

// SourceStatus is the read-only view of one registered source, as returned by
// GET /api/mcp/sources (contracts/mcp-sources-api.md).
type SourceStatus struct {
	Name         string       `json:"name"`
	Kind         SourceKind   `json:"kind"`
	Health       HealthStatus `json:"health"`
	Capabilities []string     `json:"capabilities"`
}

// Aggregator merges every registered ToolSource's capabilities into one unified set. Capability
// name ownership is resolved once, at construction, in registration order (research.md R6); which
// of those owned capabilities are actually callable is re-evaluated on every RenderTools/Dispatch
// call against each owning source's current health (research.md R4, US3), so a source recovering
// from an outage becomes usable again without rebuilding the Aggregator.
type Aggregator struct {
	mu        sync.RWMutex
	sources   []ToolSource          // registration order: local source first, then external sources in config order
	byName    map[string]ToolSource // source name -> source
	owner     map[string]Capability // capability name -> the Capability that won ownership
	conflicts []Conflict
}

// NewAggregator resolves capability ownership across sources (in the order given — the caller is
// responsible for putting the local source first) and returns the ready-to-use Aggregator.
func NewAggregator(sources ...ToolSource) *Aggregator {
	a := &Aggregator{
		sources: sources,
		byName:  make(map[string]ToolSource, len(sources)),
	}
	for _, src := range sources {
		a.byName[src.Name()] = src
	}
	a.owner, a.conflicts = resolveOwnership(sources)
	return a
}

// resolveOwnership walks sources in order, claiming each capability name for the first source
// that declares it; every later collision becomes a Conflict instead of overwriting the winner.
func resolveOwnership(sources []ToolSource) (map[string]Capability, []Conflict) {
	owner := make(map[string]Capability)
	var conflicts []Conflict
	for _, src := range sources {
		for _, cap := range src.Capabilities() {
			if !cap.ReadOnly { // defense in depth — sources are expected to have already filtered (research.md R5)
				continue
			}
			existing, taken := owner[cap.Name]
			if !taken {
				owner[cap.Name] = cap
				continue
			}
			conflicts = append(conflicts, Conflict{
				CapabilityName: cap.Name,
				WinningSource:  existing.SourceName,
				RejectedSource: cap.SourceName,
			})
		}
	}
	return owner, conflicts
}

// RenderTools returns the Anthropic tool schema for every capability whose owning source is
// currently healthy, sorted by name for a deterministic tool list across calls.
func (a *Aggregator) RenderTools() []anthropic.ToolUnionParam {
	a.mu.RLock()
	defer a.mu.RUnlock()

	names := make([]string, 0, len(a.owner))
	for name := range a.owner {
		names = append(names, name)
	}
	sort.Strings(names)

	tools := make([]anthropic.ToolUnionParam, 0, len(names))
	for _, name := range names {
		cap := a.owner[name]
		src, ok := a.byName[cap.SourceName]
		if !ok || src.Health().State != HealthHealthy {
			continue
		}
		tools = append(tools, cap.toAnthropicTool())
	}
	return tools
}

// Dispatch invokes the named capability against its owning source, returning that source's name
// alongside the result so the caller can attribute the response (research.md R7, FR-008). A
// capability whose owning source is currently unhealthy is reported as a tool-level error rather
// than attempted, matching how a Go-level Call failure is already reported.
func (a *Aggregator) Dispatch(ctx context.Context, capability string, args json.RawMessage) (result string, isError bool, sourceName string, err error) {
	a.mu.RLock()
	cap, ok := a.owner[capability]
	var src ToolSource
	if ok {
		src = a.byName[cap.SourceName]
	}
	a.mu.RUnlock()

	if !ok {
		return "", true, "", fmt.Errorf("unknown capability %q", capability)
	}
	if src == nil || src.Health().State != HealthHealthy {
		return fmt.Sprintf("tool source %q is currently unavailable", cap.SourceName), true, cap.SourceName, nil
	}

	result, isError, err = src.Call(ctx, capability, args)
	return result, isError, cap.SourceName, err
}

// Status returns every registered source's current health and capability list, plus any
// capability-name conflicts detected at construction (FR-003, contracts/mcp-sources-api.md). A
// source is always included, healthy or not (US3 AC2) — never omitted.
func (a *Aggregator) Status() ([]SourceStatus, []Conflict) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	statuses := make([]SourceStatus, 0, len(a.sources))
	for _, src := range a.sources {
		caps := src.Capabilities()
		names := make([]string, 0, len(caps))
		for _, c := range caps {
			names = append(names, c.Name)
		}
		statuses = append(statuses, SourceStatus{
			Name:         src.Name(),
			Kind:         src.Kind(),
			Health:       src.Health(),
			Capabilities: names,
		})
	}
	return statuses, a.conflicts
}

func (c Capability) toAnthropicTool() anthropic.ToolUnionParam {
	t := anthropic.ToolParam{
		Name:        c.Name,
		Description: anthropic.String(c.Description),
		InputSchema: anthropic.ToolInputSchemaParam{Properties: c.InputSchema},
	}
	return anthropic.ToolUnionParam{OfTool: &t}
}
