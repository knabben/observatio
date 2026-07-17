package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"sync"
	"time"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// transportFactory builds a fresh mcp.Transport on every call — a stdio source needs a new
// *exec.Cmd per (re)connect attempt (an exec.Cmd can't be reused after it exits), and modeling an
// HTTP source's reconnect the same way keeps both transports behind one interface.
type transportFactory func() mcpsdk.Transport

func newStdioTransportFactory(command string, args []string) transportFactory {
	return func() mcpsdk.Transport {
		return &mcpsdk.CommandTransport{Command: exec.Command(command, args...)}
	}
}

func newHTTPTransportFactory(url string) transportFactory {
	return func() mcpsdk.Transport {
		return &mcpsdk.StreamableClientTransport{Endpoint: url}
	}
}

// MCPToolSource is an operator-registered external tool source, reached over a real MCP
// transport (stdio or streamable HTTP — research.md R9). Its capability list and health are
// refreshed out-of-band by a background health check (health.go, research.md R4), never inline
// with a chat turn.
type MCPToolSource struct {
	name    string
	factory transportFactory
	client  *mcpsdk.Client

	mu      sync.RWMutex
	session *mcpsdk.ClientSession
	caps    []Capability
	health  HealthStatus
}

func newMCPToolSource(name string, factory transportFactory) *MCPToolSource {
	return &MCPToolSource{
		name:    name,
		factory: factory,
		client:  mcpsdk.NewClient(&mcpsdk.Implementation{Name: "observatio-aggregator", Version: "v1"}, nil),
		health:  HealthStatus{State: HealthUnknown},
	}
}

// NewMCPToolSource constructs an external tool source over stdio or streamable HTTP per its
// TransportConfig (contracts/tool-sources-config.md) and performs the first health probe
// synchronously, so its capabilities participate in the Aggregator's startup conflict resolution
// (research.md R6) — a source that fails this first probe is still returned (unhealthy, no
// capabilities yet), never an error, since a single unreachable source at startup must not
// prevent the rest of the aggregator from coming up (US3).
func NewMCPToolSource(ctx context.Context, entry SourceEntry) *MCPToolSource {
	var factory transportFactory
	switch entry.Transport.Kind {
	case TransportStdio:
		factory = newStdioTransportFactory(entry.Transport.Command, entry.Transport.Args)
	case TransportHTTP:
		factory = newHTTPTransportFactory(entry.Transport.URL)
	}
	src := newMCPToolSource(entry.Name, factory)
	src.recordProbeResult(src.probe(ctx))
	return src
}

func (s *MCPToolSource) Name() string     { return s.name }
func (s *MCPToolSource) Kind() SourceKind { return SourceKindExternal }

func (s *MCPToolSource) Capabilities() []Capability {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.caps
}

func (s *MCPToolSource) Health() HealthStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.health
}

// Call invokes the named capability against the source's current session, reconnecting first if
// there isn't one (e.g. the previous session dropped between health checks).
func (s *MCPToolSource) Call(ctx context.Context, capability string, args json.RawMessage) (string, bool, error) {
	session, err := s.ensureSession(ctx)
	if err != nil {
		return "", true, err
	}

	var arguments any
	if len(args) > 0 {
		if err := json.Unmarshal(args, &arguments); err != nil {
			return "", true, err
		}
	}

	res, err := session.CallTool(ctx, &mcpsdk.CallToolParams{Name: capability, Arguments: arguments})
	if err != nil {
		s.dropSession()
		return "", true, err
	}
	return extractText(res.Content), res.IsError, nil
}

// probe connects (if needed) and calls tools/list, refreshing the cached capability list on
// success. It never mutates health itself — callers decide how a probe result maps to
// HealthStatus (recordProbeResult), so this method can be reused identically by the synchronous
// startup probe and the background health-check loop (health.go).
func (s *MCPToolSource) probe(ctx context.Context) error {
	session, err := s.ensureSession(ctx)
	if err != nil {
		return err
	}
	result, err := session.ListTools(ctx, nil)
	if err != nil {
		s.dropSession()
		return err
	}
	caps := translateTools(s.name, result.Tools)
	s.mu.Lock()
	s.caps = caps
	s.mu.Unlock()
	return nil
}

func (s *MCPToolSource) recordProbeResult(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err != nil {
		s.health = HealthStatus{State: HealthUnhealthy, LastChecked: time.Now(), LastError: err.Error()}
		return
	}
	s.health = HealthStatus{State: HealthHealthy, LastChecked: time.Now()}
}

func (s *MCPToolSource) ensureSession(ctx context.Context) (*mcpsdk.ClientSession, error) {
	s.mu.RLock()
	session := s.session
	s.mu.RUnlock()
	if session != nil {
		return session, nil
	}

	newSession, err := s.client.Connect(ctx, s.factory(), nil)
	if err != nil {
		return nil, fmt.Errorf("connecting to tool source %q: %w", s.name, err)
	}
	s.mu.Lock()
	s.session = newSession
	s.mu.Unlock()
	return newSession, nil
}

func (s *MCPToolSource) dropSession() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.session = nil
}
