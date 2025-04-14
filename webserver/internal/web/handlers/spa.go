package handlers

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

type SPAHandler struct {
	staticFS   embed.FS
	staticPath string
	indexPath  string
}

func (h SPAHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path = filepath.Join(h.staticPath, path)
	if _, err = h.staticFS.Open(path); os.IsNotExist(err) {
		index, err := h.staticFS.ReadFile(filepath.Join(h.staticPath, h.indexPath))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusAccepted)
		w.Write(index)
		return

	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var files fs.FS
	files, err = fs.Sub(h.staticFS, h.staticPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.FS(files)).ServeHTTP(w, r)
}
