package mcp

import (
	"strings"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// extractText joins every TextContent block in an MCP tool result into one string, the shape
// this package's ToolSource.Call contract returns to its caller. Non-text content blocks (images,
// embedded resources) are skipped — no registered capability in this feature produces them.
func extractText(blocks []mcpsdk.Content) string {
	var sb strings.Builder
	for _, block := range blocks {
		if text, ok := block.(*mcpsdk.TextContent); ok {
			if sb.Len() > 0 {
				sb.WriteString("\n")
			}
			sb.WriteString(text.Text)
		}
	}
	return sb.String()
}
