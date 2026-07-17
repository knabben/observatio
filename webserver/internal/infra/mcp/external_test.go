package mcp

import (
	"context"
	"testing"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newFakeServerTransportFactory starts an in-memory MCP server exposing the given tools (each
// echoing back its arguments as text) and returns a transportFactory that connects to it — the
// same shape MCPToolSource uses for a real stdio/HTTP source, standing in for one in tests.
func newFakeServerTransportFactory(t *testing.T, tools ...*mcpsdk.Tool) transportFactory {
	t.Helper()
	server := mcpsdk.NewServer(&mcpsdk.Implementation{Name: "fake", Version: "v1"}, nil)
	for _, tool := range tools {
		server.AddTool(tool, func(_ context.Context, req *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
			return &mcpsdk.CallToolResult{
				Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: "echo:" + string(req.Params.Arguments)}},
			}, nil
		})
	}

	// A fresh in-memory transport pair per factory call, matching how a real stdio source gets a
	// new process (and a real HTTP source a new connection) on every (re)connect attempt.
	return func() mcpsdk.Transport {
		serverTransport, clientTransport := mcpsdk.NewInMemoryTransports()
		go func() {
			_, _ = server.Connect(context.Background(), serverTransport, nil)
		}()
		return clientTransport
	}
}

func readOnlyTool(name string) *mcpsdk.Tool {
	return &mcpsdk.Tool{
		Name:        name,
		Description: name,
		InputSchema: map[string]any{"type": "object"},
		Annotations: &mcpsdk.ToolAnnotations{ReadOnlyHint: true},
	}
}

func TestMCPToolSource_ProbePopulatesCapabilities(t *testing.T) {
	factory := newFakeServerTransportFactory(t, readOnlyTool("list_backups"), readOnlyTool("describe_backup"))
	src := newMCPToolSource("velero-mcp", factory)

	require.NoError(t, src.probe(context.Background()))

	caps := src.Capabilities()
	names := []string{caps[0].Name, caps[1].Name}
	assert.ElementsMatch(t, []string{"list_backups", "describe_backup"}, names)
	for _, c := range caps {
		assert.Equal(t, "velero-mcp", c.SourceName)
		assert.True(t, c.ReadOnly)
	}
}

func TestMCPToolSource_ProbeDropsNonReadOnlyTools(t *testing.T) {
	mutating := &mcpsdk.Tool{
		Name:        "trigger_restore",
		InputSchema: map[string]any{"type": "object"},
		Annotations: &mcpsdk.ToolAnnotations{ReadOnlyHint: false},
	}
	factory := newFakeServerTransportFactory(t, readOnlyTool("list_backups"), mutating)
	src := newMCPToolSource("velero-mcp", factory)

	require.NoError(t, src.probe(context.Background()))

	caps := src.Capabilities()
	require.Len(t, caps, 1)
	assert.Equal(t, "list_backups", caps[0].Name)
}

func TestMCPToolSource_ProbeFailureReportsUnhealthy(t *testing.T) {
	// A stdio transport pointed at a nonexistent binary fails fast at Command.Start(), unlike an
	// unpaired in-memory transport (which would block forever waiting for a response nobody sends).
	factory := newStdioTransportFactory("this-command-does-not-exist-xyz", nil)
	src := newMCPToolSource("broken", factory)

	err := src.probe(context.Background())
	assert.Error(t, err)
}

func TestMCPToolSource_Call_ReturnsToolResultText(t *testing.T) {
	factory := newFakeServerTransportFactory(t, readOnlyTool("list_backups"))
	src := newMCPToolSource("velero-mcp", factory)
	require.NoError(t, src.probe(context.Background()))

	result, isError, err := src.Call(context.Background(), "list_backups", []byte(`{"namespace":"default"}`))
	require.NoError(t, err)
	assert.False(t, isError)
	assert.Contains(t, result, "namespace")
}

func TestMCPToolSource_RecordProbeResult_SuccessIsHealthy(t *testing.T) {
	factory := newFakeServerTransportFactory(t, readOnlyTool("list_backups"))
	src := newMCPToolSource("velero-mcp", factory)

	src.recordProbeResult(src.probe(context.Background()))

	assert.Equal(t, HealthHealthy, src.Health().State)
}

func TestMCPToolSource_RecordProbeResult_FailureIsUnhealthy(t *testing.T) {
	factory := newStdioTransportFactory("this-command-does-not-exist-xyz", nil)
	src := newMCPToolSource("broken", factory)

	src.recordProbeResult(src.probe(context.Background()))

	health := src.Health()
	assert.Equal(t, HealthUnhealthy, health.State)
	assert.NotEmpty(t, health.LastError)
}

func TestMCPToolSource_Name_Kind(t *testing.T) {
	src := newMCPToolSource("velero-mcp", newFakeServerTransportFactory(t, readOnlyTool("list_backups")))
	assert.Equal(t, "velero-mcp", src.Name())
	assert.Equal(t, SourceKindExternal, src.Kind())
}

func TestAggregator_ConflictDetection_TwoRealExternalSources(t *testing.T) {
	a := newMCPToolSource("velero-mcp", newFakeServerTransportFactory(t, readOnlyTool("list_backups")))
	require.NoError(t, a.probe(context.Background()))
	a.recordProbeResult(nil)

	b := newMCPToolSource("velero-mcp-mirror", newFakeServerTransportFactory(t, readOnlyTool("list_backups")))
	require.NoError(t, b.probe(context.Background()))
	b.recordProbeResult(nil)

	agg := NewAggregator(a, b)

	assert.Equal(t, []string{"list_backups"}, toolNames(t, agg), "only one of the two identically-named capabilities may ever be callable")

	statuses, conflicts := agg.Status()
	require.Len(t, statuses, 2, "both sources are still listed, never omitted, even though one lost the naming conflict")
	require.Len(t, conflicts, 1)
	assert.Equal(t, Conflict{CapabilityName: "list_backups", WinningSource: "velero-mcp", RejectedSource: "velero-mcp-mirror"}, conflicts[0])
}
