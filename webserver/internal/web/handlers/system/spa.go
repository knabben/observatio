package system

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

type SPAHandler struct {
	StaticFS   embed.FS
	StaticPath string
	IndexPath  string
}

func (h SPAHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path = filepath.Join(h.StaticPath, path)
	if _, err = h.StaticFS.Open(path); os.IsNotExist(err) {
		index, err := h.StaticFS.ReadFile(filepath.Join(h.StaticPath, h.IndexPath))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusAccepted)
		if _, err = w.Write(index); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return

	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var files fs.FS
	files, err = fs.Sub(h.StaticFS, h.StaticPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.FS(files)).ServeHTTP(w, r)
}
