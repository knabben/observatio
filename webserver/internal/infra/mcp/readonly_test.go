package mcp

import (
	"testing"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
)

func TestIsReadOnly(t *testing.T) {
	tests := []struct {
		name string
		tool *mcpsdk.Tool
		want bool
	}{
		{
			name: "explicit true is read-only",
			tool: &mcpsdk.Tool{Name: "list", Annotations: &mcpsdk.ToolAnnotations{ReadOnlyHint: true}},
			want: true,
		},
		{
			name: "explicit false is not read-only",
			tool: &mcpsdk.Tool{Name: "restore", Annotations: &mcpsdk.ToolAnnotations{ReadOnlyHint: false}},
			want: false,
		},
		{
			name: "missing annotations block is not read-only — fail closed",
			tool: &mcpsdk.Tool{Name: "mystery"},
			want: false,
		},
		{
			name: "nil tool is not read-only",
			tool: nil,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsReadOnly(tt.tool))
		})
	}
}

func TestTranslateTools_DropsNonReadOnlyTools(t *testing.T) {
	tools := []*mcpsdk.Tool{
		{Name: "list_backups", Description: "list backups", Annotations: &mcpsdk.ToolAnnotations{ReadOnlyHint: true}},
		{Name: "trigger_restore", Description: "trigger a restore", Annotations: &mcpsdk.ToolAnnotations{ReadOnlyHint: false}},
		{Name: "unannotated", Description: "no annotations at all"},
	}

	caps := translateTools("velero-mcp", tools)

	assert.Len(t, caps, 1)
	assert.Equal(t, "list_backups", caps[0].Name)
	assert.Equal(t, "velero-mcp", caps[0].SourceName)
	assert.True(t, caps[0].ReadOnly)
}
