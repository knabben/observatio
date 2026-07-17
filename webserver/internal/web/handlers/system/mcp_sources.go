package system

import (
	"net/http"

	mcpaggregator "github.com/knabben/observatio/webserver/internal/infra/mcp"
)

// mcpSourcesResponse is the GET /api/mcp/sources response shape
// (contracts/mcp-sources-api.md). Sources and Conflicts are never nil - an aggregator with zero
// external sources still reports the always-present local one, and zero conflicts marshal as
// `[]`, never `null`.
type mcpSourcesResponse struct {
	Sources   []mcpaggregator.SourceStatus `json:"sources"`
	Conflicts []mcpaggregator.Conflict     `json:"conflicts"`
}

// HandleMCPSources returns every registered tool source's health and capabilities, plus any
// capability-name conflicts detected at startup (FR-003, SC-005). aggregator is the shared,
// process-wide Aggregator - this handler never mutates it.
func HandleMCPSources(aggregator *mcpaggregator.Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statuses, conflicts := aggregator.Status()
		if statuses == nil {
			statuses = []mcpaggregator.SourceStatus{}
		}
		if conflicts == nil {
			conflicts = []mcpaggregator.Conflict{}
		}

		err := WriteResponse(w, mcpSourcesResponse{Sources: statuses, Conflicts: conflicts})
		if HandleError(w, http.StatusInternalServerError, err) {
			return
		}
	}
}
