package system

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

type SPAHandler struct {
	StaticFS   embed.FS
	StaticPath string
	IndexPath  string
}

// ServeHTTP serves an embedded static asset when the request path resolves to one, and
// otherwise falls back to the SPA shell (IndexPath) with a 200 — covering the root path
// and any unknown client-side route — so the client-side router always gets HTML to
// bootstrap from instead of a 404 or a directory listing.
func (h SPAHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	files, err := fs.Sub(h.StaticFS, h.StaticPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	requestPath := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
	if requestPath != "" && requestPath != "." {
		if f, openErr := files.Open(requestPath); openErr == nil {
			_ = f.Close()
			http.FileServer(http.FS(files)).ServeHTTP(w, r)
			return
		}
	}

	h.serveIndex(w, files)
}

func (h SPAHandler) serveIndex(w http.ResponseWriter, files fs.FS) {
	index, err := fs.ReadFile(files, h.IndexPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(index); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
