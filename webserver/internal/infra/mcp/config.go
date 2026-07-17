package mcp

import (
	"context"
	"fmt"
	"os"

	"sigs.k8s.io/yaml"
)

// TransportKind is one of the two transports an external tool source may be reached over
// (research.md R9).
type TransportKind string

const (
	TransportStdio TransportKind = "stdio"
	TransportHTTP  TransportKind = "http"
)

// TransportConfig describes how to reach one external tool source, per
// contracts/tool-sources-config.md.
type TransportConfig struct {
	Kind    TransportKind `json:"kind"`
	Command string        `json:"command,omitempty"`
	Args    []string      `json:"args,omitempty"`
	URL     string        `json:"url,omitempty"`
}

// SourceEntry is one operator-declared external tool source.
type SourceEntry struct {
	Name      string          `json:"name"`
	Enabled   bool            `json:"enabled"`
	Transport TransportConfig `json:"transport"`
}

// SourceConfig is the top-level shape of the file at --tool-sources-config/TOOL_SOURCES_CONFIG.
type SourceConfig struct {
	Sources []SourceEntry `json:"sources"`
}

// LoadSourceConfig reads and validates the tool sources config file at path (research.md R3,
// contracts/tool-sources-config.md). It fails fast: a malformed file, or one that fails
// validation, is returned as an error rather than silently registering zero sources — this is
// distinct from the flag/env var being unset entirely, which is the valid zero-external-sources
// case and is handled by the caller, not this function.
func LoadSourceConfig(path string) (*SourceConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading tool sources config %q: %w", path, err)
	}
	var cfg SourceConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing tool sources config %q: %w", path, err)
	}
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("validating tool sources config %q: %w", path, err)
	}
	return &cfg, nil
}

// BuildExternalSources loads the tool sources config at path (if non-empty) and constructs one
// MCPToolSource per enabled entry, performing each source's first health probe synchronously so
// its capabilities participate in the Aggregator's startup conflict resolution (research.md R6).
// An empty path is valid - it's the zero-external-sources case (US1's Independent Test) - and
// returns (nil, nil), not an error.
func BuildExternalSources(ctx context.Context, path string) ([]*MCPToolSource, error) {
	if path == "" {
		return nil, nil
	}
	cfg, err := LoadSourceConfig(path)
	if err != nil {
		return nil, err
	}
	var sources []*MCPToolSource
	for _, entry := range cfg.Sources {
		if !entry.Enabled {
			continue
		}
		sources = append(sources, NewMCPToolSource(ctx, entry))
	}
	return sources, nil
}

func (c *SourceConfig) validate() error {
	seen := make(map[string]bool, len(c.Sources))
	for i, s := range c.Sources {
		if s.Name == "" {
			return fmt.Errorf("sources[%d]: name is required", i)
		}
		if s.Name == localSourceName {
			return fmt.Errorf("sources[%d]: name %q is reserved for the built-in local source", i, s.Name)
		}
		if seen[s.Name] {
			return fmt.Errorf("sources[%d]: duplicate source name %q", i, s.Name)
		}
		seen[s.Name] = true

		switch s.Transport.Kind {
		case TransportStdio:
			if s.Transport.Command == "" {
				return fmt.Errorf("sources[%d] (%s): transport.command is required for stdio transport", i, s.Name)
			}
		case TransportHTTP:
			if s.Transport.URL == "" {
				return fmt.Errorf("sources[%d] (%s): transport.url is required for http transport", i, s.Name)
			}
		default:
			return fmt.Errorf("sources[%d] (%s): transport.kind must be %q or %q, got %q", i, s.Name, TransportStdio, TransportHTTP, s.Transport.Kind)
		}
	}
	return nil
}
