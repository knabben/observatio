package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// localSourceName is both the built-in ToolSource's name and its one capability's name, matching
// today's "kubectl" tool exactly (FR-001).
const localSourceName = "kubectl"

// LocalToolSource wraps today's kubectl-backed capability as a real, in-process MCP server (per
// spec.md Clarifications, 2026-07-17) — connected to its own client over an in-memory transport,
// not a bespoke Go-interface adapter. It is structurally the same kind of thing as an external
// ToolSource, differing only in which transport its ClientSession was connected over
// (research.md R2).
type LocalToolSource struct {
	session *mcpsdk.ClientSession
}

// NewLocalToolSource builds the kubectl MCP server and connects an in-process client to it via
// mcp.NewInMemoryTransports (no OS subprocess). ctx bounds only the connection handshake.
func NewLocalToolSource(ctx context.Context) (*LocalToolSource, error) {
	server := mcpsdk.NewServer(&mcpsdk.Implementation{Name: localSourceName, Version: "v1"}, nil)
	server.AddTool(&mcpsdk.Tool{
		Name:        localSourceName,
		Description: "kubectl is a command-line tool for controlling Kubernetes clusters.",
		InputSchema: kubectlInputSchema,
		Annotations: &mcpsdk.ToolAnnotations{ReadOnlyHint: true},
	}, handleKubectlCall)

	serverTransport, clientTransport := mcpsdk.NewInMemoryTransports()
	if _, err := server.Connect(ctx, serverTransport, nil); err != nil {
		return nil, fmt.Errorf("starting local kubectl MCP server: %w", err)
	}

	client := mcpsdk.NewClient(&mcpsdk.Implementation{Name: "observatio-local-client", Version: "v1"}, nil)
	session, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		return nil, fmt.Errorf("connecting to local kubectl MCP server: %w", err)
	}
	return &LocalToolSource{session: session}, nil
}

// kubectlProperties is shared between the MCP server's full JSON-schema InputSchema (which needs
// the "type": "object" wrapper) and Capability.InputSchema (which holds only the properties
// object — see source.go and Aggregator.toAnthropicTool) so the two can never drift apart.
var kubectlProperties = map[string]interface{}{
	"command": map[string]interface{}{"type": "string"},
}

var kubectlInputSchema = map[string]any{
	"type":       "object",
	"properties": kubectlProperties,
}

func (s *LocalToolSource) Name() string     { return localSourceName }
func (s *LocalToolSource) Kind() SourceKind { return SourceKindLocal }

// Capabilities is static — the local source's one capability never changes at runtime, unlike an
// external source's cached, health-check-refreshed list.
func (s *LocalToolSource) Capabilities() []Capability {
	return []Capability{{
		Name:        localSourceName,
		Description: "kubectl is a command-line tool for controlling Kubernetes clusters.",
		InputSchema: kubectlProperties,
		SourceName:  localSourceName,
		ReadOnly:    true,
	}}
}

func (s *LocalToolSource) Call(ctx context.Context, capability string, args json.RawMessage) (string, bool, error) {
	res, err := s.session.CallTool(ctx, &mcpsdk.CallToolParams{Name: capability, Arguments: json.RawMessage(args)})
	if err != nil {
		return "", true, err
	}
	return extractText(res.Content), res.IsError, nil
}

// Health is always healthy: the local source has no network dependency of its own — its
// in-memory transport can't become unreachable the way an external source's can (research.md R4).
// A kubectl failure surfaces per-call via Call's isError result, exactly as it does today.
func (s *LocalToolSource) Health() HealthStatus {
	return HealthStatus{State: HealthHealthy}
}

// handleKubectlCall is the local MCP server's tool handler for "kubectl" — relocated from
// today's llm.RunKubectl unchanged: it shells out to the local kubectl binary with whatever
// command string the model produced. The output is returned even when the command fails, so a
// failure can be reported back to the model as a tool result (letting it explain the failure or
// retry) instead of aborting the whole exchange.
func handleKubectlCall(ctx context.Context, req *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
	var input struct {
		Command string `json:"command"`
	}
	if err := json.Unmarshal(req.Params.Arguments, &input); err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, "kubectl", strings.Fields(input.Command)...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		text := strings.TrimSpace(string(output) + "\n" + err.Error())
		return &mcpsdk.CallToolResult{
			Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: text}},
			IsError: true,
		}, nil
	}
	return &mcpsdk.CallToolResult{
		Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: string(output)}},
	}, nil
}
