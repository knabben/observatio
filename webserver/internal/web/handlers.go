package web

import (
	"embed"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

//go:embed build/*
var bundle embed.FS

func DefaultHandlers(router *mux.Router, developmentMode bool) {
	// Generic handlers, healthcheck, version, etc.
	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	// Cluster API handlers
	router.HandleFunc("/api/clusters/info", handleClusterInfo).Methods("GET")

	// Static React bundle hosting handler
	if !developmentMode {
		spa := SPAHandler{staticFS: bundle, staticPath: "build", indexPath: "index.html"}
		router.PathPrefix("/front").Handler(spa)
	}
}
