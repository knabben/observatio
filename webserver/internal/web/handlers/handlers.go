package handlers

import (
	"context"
	"embed"
	"encoding/json"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/knabben/observatio/webserver/internal/web/handlers/kubernetes"
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

	// Start the websocker handlers
	startWebSocketHandlers(router)

	// Static React bundle hosting handler
	if !developmentMode {
		spa := system.SPAHandler{StaticFS: bundle, StaticPath: "build", IndexPath: "dashboard.html"}
		router.PathPrefix("/").Handler(spa)
	}
}

// startWebSocketHandlers initializes WebSocket routes and manages their corresponding client pool behaviors.
func startWebSocketHandlers(router *mux.Router) {
	pool := &system.ClientPool{
		Broadcast:  make(chan []byte),
		Register:   make(chan *system.WSClient),
		Unregister: make(chan *system.WSClient),
		Clients:    make(map[string]*system.WSClient),
	}

	go pool.Run(context.Background())

	router.HandleFunc("/ws/watcher", system.HandleWatcher).Methods("GET", "OPTIONS")
	router.HandleFunc("/ws/analysis", func(w http.ResponseWriter, r *http.Request) {
		system.HandleChatBot(pool, w, r)
	}).Methods("GET", "OPTIONS")
}
