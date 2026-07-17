package system

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mcpaggregator "github.com/knabben/observatio/webserver/internal/infra/mcp"
)

func TestHandleMCPSources_ReportsLocalSourceAndEmptyConflicts(t *testing.T) {
	local, err := mcpaggregator.NewLocalToolSource(context.Background())
	if err != nil {
		t.Fatalf("NewLocalToolSource: %v", err)
	}
	aggregator := mcpaggregator.NewAggregator(local)

	req := httptest.NewRequest(http.MethodGet, "/api/mcp/sources", nil)
	rec := httptest.NewRecorder()
	HandleMCPSources(aggregator)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var body struct {
		Sources []struct {
			Name   string `json:"name"`
			Kind   string `json:"kind"`
			Health struct {
				State string `json:"state"`
			} `json:"health"`
			Capabilities []string `json:"capabilities"`
		} `json:"sources"`
		Conflicts []any `json:"conflicts"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v, body=%s", err, rec.Body.String())
	}

	if len(body.Sources) != 1 {
		t.Fatalf("expected exactly one source (the built-in local one), got %d", len(body.Sources))
	}
	got := body.Sources[0]
	if got.Name != "kubectl" || got.Kind != "local" || got.Health.State != "healthy" {
		t.Fatalf("unexpected local source status: %+v", got)
	}
	if len(got.Capabilities) != 1 || got.Capabilities[0] != "kubectl" {
		t.Fatalf("expected the single kubectl capability, got %v", got.Capabilities)
	}
	if body.Conflicts == nil || len(body.Conflicts) != 0 {
		t.Fatalf("expected conflicts to marshal as an empty array, got %#v", body.Conflicts)
	}
}
