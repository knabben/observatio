package system

import (
	"embed"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//go:embed testdata/site
var testBuildFS embed.FS

func newTestSPAHandler() SPAHandler {
	return SPAHandler{StaticFS: testBuildFS, StaticPath: "testdata/site", IndexPath: "dashboard.html"}
}

func TestSPAHandler_ServesShellAtRoot(t *testing.T) {
	h := newTestSPAHandler()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 at root, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "dashboard-shell") {
		t.Fatalf("expected embedded SPA shell content, got: %s", rec.Body.String())
	}
}

func TestSPAHandler_FallsBackToShellForUnknownClientRoute(t *testing.T) {
	h := newTestSPAHandler()
	req := httptest.NewRequest(http.MethodGet, "/dashboard/clusters/some-cluster", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 SPA fallback for an unknown client route, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "dashboard-shell") {
		t.Fatalf("expected embedded SPA shell content, got: %s", rec.Body.String())
	}
}

func TestSPAHandler_ServesKnownStaticAssetDirectly(t *testing.T) {
	h := newTestSPAHandler()
	req := httptest.NewRequest(http.MethodGet, "/static/app.js", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for a known static asset, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "test asset") {
		t.Fatalf("expected the actual asset content, not the SPA shell, got: %s", rec.Body.String())
	}
}
