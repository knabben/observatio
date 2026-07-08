package kubernetes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HandleControllerLogs_MissingQueryParams(t *testing.T) {
	cases := []struct {
		name  string
		query string
	}{
		{name: "missing both", query: ""},
		{name: "missing deployment", query: "?namespace=capi-system"},
		{name: "missing namespace", query: "?deployment=capi-controller-manager"},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/logs/controller"+tt.query, nil)
			rec := httptest.NewRecorder()

			HandleControllerLogs(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	}
}

func Test_HandleNodeAccess_MissingQueryParams(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/logs/node-access", nil)
	rec := httptest.NewRecorder()

	HandleNodeAccess(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
