package mcp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalToolSource_Capabilities(t *testing.T) {
	src, err := NewLocalToolSource(context.Background())
	require.NoError(t, err)

	assert.Equal(t, "kubectl", src.Name())
	assert.Equal(t, SourceKindLocal, src.Kind())

	caps := src.Capabilities()
	require.Len(t, caps, 1)
	assert.Equal(t, "kubectl", caps[0].Name)
	assert.Equal(t, "kubectl", caps[0].SourceName)
	assert.True(t, caps[0].ReadOnly, "the built-in capability must carry readOnlyHint: true")
}

func TestLocalToolSource_Health_AlwaysHealthy(t *testing.T) {
	src, err := NewLocalToolSource(context.Background())
	require.NoError(t, err)

	assert.Equal(t, HealthHealthy, src.Health().State)
}

func TestLocalToolSource_Call_SuccessPassesThroughOutput(t *testing.T) {
	src, err := NewLocalToolSource(context.Background())
	require.NoError(t, err)

	// "kubectl version --client" doesn't need a live cluster and is read-only, so it's a stable
	// smoke test for the real MCP round trip (server tool handler -> in-memory transport -> client).
	result, isError, err := src.Call(context.Background(), "kubectl", []byte(`{"command":"version --client"}`))
	require.NoError(t, err)
	assert.False(t, isError)
	assert.NotEmpty(t, result)
}

func TestLocalToolSource_Call_FailureReportsAsToolErrorNotGoError(t *testing.T) {
	src, err := NewLocalToolSource(context.Background())
	require.NoError(t, err)

	result, isError, err := src.Call(context.Background(), "kubectl", []byte(`{"command":"this-is-not-a-real-subcommand"}`))
	require.NoError(t, err, "a kubectl failure must be reported back as a tool result, not a Go error, per today's existing behavior")
	assert.True(t, isError)
	assert.NotEmpty(t, result)
}
