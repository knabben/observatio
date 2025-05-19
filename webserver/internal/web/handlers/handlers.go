package handlers

import (
	"embed"
	"encoding/json"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/knabben/observatio/webserver/internal/web/handlers/kubernetes"
	"github.com/knabben/observatio/webserver/internal/web/handlers/llm"
	"github.com/knabben/observatio/webserver/internal/web/handlers/system"
)

//go:embed build/*
var bundle embed.FS

func DefaultHandlers(router *mux.Router, developmentMode bool) {
	// Generic handlers, healthcheck, version, etc.
	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(map[string]bool{"ok": true}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	// Cluster API handlers
	router.HandleFunc("/api/clusters/list", kubernetes.HandleClusterList).Methods("GET")
	router.HandleFunc("/api/clusters/infra/list", kubernetes.HandleClusterInfraList).Methods("GET")

	// Cluster API dashboard Handlers
	router.HandleFunc("/api/clusters/info", kubernetes.HandleClusterInfo).Methods("GET")
	router.HandleFunc("/api/clusters/components", kubernetes.HandleComponentsVersion).Methods("GET")
	router.HandleFunc("/api/clusters/summary", kubernetes.HandleSummaryCluster).Methods("GET")
	router.HandleFunc("/api/clusters/classes", kubernetes.HandleClusterClasses).Methods("GET")
	router.HandleFunc("/api/clusters/topology", kubernetes.HandleClusterTopology).Methods("GET")

	// Cluster API Machine Deployments Handlers
	router.HandleFunc("/api/machinesdeployment/list", kubernetes.HandleMachineDeployments).Methods("GET")

	// Cluster API Machine Handlers
	router.HandleFunc("/api/machines/list", kubernetes.HandleMachines).Methods("GET")
	router.HandleFunc("/api/machines/infra/list", kubernetes.HandleMachineInfra).Methods("GET")

	// Anthropic LLM handlers
	router.HandleFunc("/api/analysis", llm.HandleClaude).Methods("POST", "OPTIONS")

	// Websocket Handler for object watchers.
	router.HandleFunc("/ws", system.HandleWebsocket)

	// Static React bundle hosting handler
	if !developmentMode {
		spa := system.SPAHandler{StaticFS: bundle, StaticPath: "build", IndexPath: "dashboard.html"}
		router.PathPrefix("/").Handler(spa)
	}
}
