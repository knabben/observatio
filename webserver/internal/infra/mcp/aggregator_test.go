package mcp

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeToolSource is a minimal ToolSource test double used across this package's tests.
type fakeToolSource struct {
	name   string
	kind   SourceKind
	caps   []Capability
	health HealthStatus
	callFn func(ctx context.Context, capability string, args json.RawMessage) (string, bool, error)
	calls  []string
}

func (f *fakeToolSource) Name() string               { return f.name }
func (f *fakeToolSource) Kind() SourceKind           { return f.kind }
func (f *fakeToolSource) Capabilities() []Capability { return f.caps }
func (f *fakeToolSource) Health() HealthStatus       { return f.health }
func (f *fakeToolSource) Call(ctx context.Context, capability string, args json.RawMessage) (string, bool, error) {
	f.calls = append(f.calls, capability)
	if f.callFn != nil {
		return f.callFn(ctx, capability, args)
	}
	return "ok", false, nil
}

func healthySource(name string, caps ...Capability) *fakeToolSource {
	return &fakeToolSource{name: name, kind: SourceKindExternal, caps: caps, health: HealthStatus{State: HealthHealthy}}
}

func cap(name, sourceName string) Capability {
	return Capability{Name: name, Description: name, SourceName: sourceName, ReadOnly: true}
}

func toolNames(t *testing.T, agg *Aggregator) []string {
	t.Helper()
	tools := agg.RenderTools()
	names := make([]string, 0, len(tools))
	for _, tool := range tools {
		names = append(names, tool.OfTool.Name)
	}
	return names
}

func TestAggregator_RenderTools_MergesHealthySources(t *testing.T) {
	a := healthySource("a", cap("cap-a", "a"))
	b := healthySource("b", cap("cap-b", "b"))

	agg := NewAggregator(a, b)

	assert.Equal(t, []string{"cap-a", "cap-b"}, toolNames(t, agg))
}

func TestAggregator_RenderTools_ExcludesUnhealthySource(t *testing.T) {
	a := healthySource("a", cap("cap-a", "a"))
	b := healthySource("b", cap("cap-b", "b"))
	b.health = HealthStatus{State: HealthUnhealthy, LastError: "connection refused"}

	agg := NewAggregator(a, b)

	assert.Equal(t, []string{"cap-a"}, toolNames(t, agg))
}

func TestAggregator_RenderTools_ReactsToHealthChangesWithoutRebuild(t *testing.T) {
	a := healthySource("a", cap("cap-a", "a"))
	b := healthySource("b", cap("cap-b", "b"))
	agg := NewAggregator(a, b)
	require.Equal(t, []string{"cap-a", "cap-b"}, toolNames(t, agg))

	b.health = HealthStatus{State: HealthUnhealthy}
	assert.Equal(t, []string{"cap-a"}, toolNames(t, agg), "an unhealthy source's capability must disappear from the next RenderTools call without reconstructing the Aggregator")

	b.health = HealthStatus{State: HealthHealthy}
	assert.Equal(t, []string{"cap-a", "cap-b"}, toolNames(t, agg), "a recovered source's capability must reappear without re-registration")
}

func TestAggregator_ConflictDetection_FirstRegisteredWins(t *testing.T) {
	a := healthySource("a", cap("shared", "a"))
	b := healthySource("b", cap("shared", "b"))

	agg := NewAggregator(a, b)

	assert.Equal(t, []string{"shared"}, toolNames(t, agg), "only one of the two colliding capabilities may ever be callable")

	_, conflicts := agg.Status()
	require.Len(t, conflicts, 1)
	assert.Equal(t, Conflict{CapabilityName: "shared", WinningSource: "a", RejectedSource: "b"}, conflicts[0])
}

func TestAggregator_ReadOnlyDefenseInDepth(t *testing.T) {
	a := healthySource("a", Capability{Name: "mutate", SourceName: "a", ReadOnly: false})

	agg := NewAggregator(a)

	assert.Empty(t, toolNames(t, agg), "a non-read-only capability must never become callable, even if a ToolSource incorrectly reports one")
}

func TestAggregator_Status_AlwaysListsEverySource(t *testing.T) {
	a := healthySource("a", cap("cap-a", "a"))
	b := healthySource("b", cap("cap-b", "b"))
	b.health = HealthStatus{State: HealthUnhealthy, LastError: "boom"}

	agg := NewAggregator(a, b)
	statuses, _ := agg.Status()

	require.Len(t, statuses, 2, "an unhealthy source must still be listed, never omitted")
	assert.Equal(t, "a", statuses[0].Name)
	assert.Equal(t, HealthHealthy, statuses[0].Health.State)
	assert.Equal(t, "b", statuses[1].Name)
	assert.Equal(t, HealthUnhealthy, statuses[1].Health.State)
	assert.Equal(t, []string{"cap-b"}, statuses[1].Capabilities, "status still reports what an unhealthy source would contribute")
}

func TestAggregator_Dispatch_RoutesToOwningSource(t *testing.T) {
	a := healthySource("a", cap("cap-a", "a"))
	a.callFn = func(_ context.Context, capability string, _ json.RawMessage) (string, bool, error) {
		return "result-from-a", false, nil
	}
	b := healthySource("b", cap("cap-b", "b"))

	agg := NewAggregator(a, b)

	result, isError, sourceName, err := agg.Dispatch(context.Background(), "cap-a", nil)
	require.NoError(t, err)
	assert.False(t, isError)
	assert.Equal(t, "result-from-a", result)
	assert.Equal(t, "a", sourceName)
	assert.Equal(t, []string{"cap-a"}, a.calls)
	assert.Empty(t, b.calls)
}

func TestAggregator_Dispatch_UnknownCapability(t *testing.T) {
	agg := NewAggregator(healthySource("a", cap("cap-a", "a")))

	_, _, _, err := agg.Dispatch(context.Background(), "does-not-exist", nil)
	assert.Error(t, err)
}

func TestAggregator_Dispatch_UnhealthySourceReportsErrorWithoutCalling(t *testing.T) {
	a := healthySource("a", cap("cap-a", "a"))
	a.health = HealthStatus{State: HealthUnhealthy}

	agg := NewAggregator(a)

	result, isError, sourceName, err := agg.Dispatch(context.Background(), "cap-a", nil)
	require.NoError(t, err)
	assert.True(t, isError)
	assert.Equal(t, "a", sourceName)
	assert.NotEmpty(t, result)
	assert.Empty(t, a.calls, "an unhealthy source's Call must never be invoked")
}
