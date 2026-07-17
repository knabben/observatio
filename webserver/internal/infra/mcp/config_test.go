package mcp

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "tool-sources.yaml")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func TestLoadSourceConfig_ValidStdioAndHTTP(t *testing.T) {
	path := writeConfig(t, `
sources:
  - name: velero-mcp
    enabled: true
    transport:
      kind: http
      url: http://velero-mcp.velero.svc.cluster.local:8080
  - name: local-test
    enabled: false
    transport:
      kind: stdio
      command: some-mcp-server
      args: ["--read-only"]
`)

	cfg, err := LoadSourceConfig(path)
	require.NoError(t, err)
	require.Len(t, cfg.Sources, 2)

	assert.Equal(t, "velero-mcp", cfg.Sources[0].Name)
	assert.True(t, cfg.Sources[0].Enabled)
	assert.Equal(t, TransportHTTP, cfg.Sources[0].Transport.Kind)
	assert.Equal(t, "http://velero-mcp.velero.svc.cluster.local:8080", cfg.Sources[0].Transport.URL)

	assert.Equal(t, "local-test", cfg.Sources[1].Name)
	assert.False(t, cfg.Sources[1].Enabled)
	assert.Equal(t, TransportStdio, cfg.Sources[1].Transport.Kind)
	assert.Equal(t, "some-mcp-server", cfg.Sources[1].Transport.Command)
	assert.Equal(t, []string{"--read-only"}, cfg.Sources[1].Transport.Args)
}

func TestLoadSourceConfig_RejectsDuplicateNames(t *testing.T) {
	path := writeConfig(t, `
sources:
  - name: dup
    enabled: true
    transport: {kind: http, url: "http://a"}
  - name: dup
    enabled: true
    transport: {kind: http, url: "http://b"}
`)

	_, err := LoadSourceConfig(path)
	assert.ErrorContains(t, err, "duplicate")
}

func TestLoadSourceConfig_RejectsReservedLocalName(t *testing.T) {
	path := writeConfig(t, `
sources:
  - name: kubectl
    enabled: true
    transport: {kind: http, url: "http://a"}
`)

	_, err := LoadSourceConfig(path)
	assert.ErrorContains(t, err, "reserved")
}

func TestLoadSourceConfig_RejectsMissingName(t *testing.T) {
	path := writeConfig(t, `
sources:
  - enabled: true
    transport: {kind: http, url: "http://a"}
`)

	_, err := LoadSourceConfig(path)
	assert.ErrorContains(t, err, "name")
}

func TestLoadSourceConfig_RejectsUnknownTransportKind(t *testing.T) {
	path := writeConfig(t, `
sources:
  - name: bad
    enabled: true
    transport: {kind: carrier-pigeon}
`)

	_, err := LoadSourceConfig(path)
	assert.ErrorContains(t, err, "transport.kind")
}

func TestLoadSourceConfig_RejectsStdioWithoutCommand(t *testing.T) {
	path := writeConfig(t, `
sources:
  - name: bad
    enabled: true
    transport: {kind: stdio}
`)

	_, err := LoadSourceConfig(path)
	assert.ErrorContains(t, err, "command")
}

func TestLoadSourceConfig_RejectsHTTPWithoutURL(t *testing.T) {
	path := writeConfig(t, `
sources:
  - name: bad
    enabled: true
    transport: {kind: http}
`)

	_, err := LoadSourceConfig(path)
	assert.ErrorContains(t, err, "url")
}

func TestLoadSourceConfig_RejectsMalformedYAML(t *testing.T) {
	path := writeConfig(t, "sources: [this is not valid yaml")

	_, err := LoadSourceConfig(path)
	assert.Error(t, err)
}

func TestLoadSourceConfig_RejectsMissingFile(t *testing.T) {
	_, err := LoadSourceConfig(filepath.Join(t.TempDir(), "does-not-exist.yaml"))
	assert.Error(t, err)
}
