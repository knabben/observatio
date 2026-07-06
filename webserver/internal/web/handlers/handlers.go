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

	// Cluster API handlers. Each Cluster/Machine item carries a derived "provider"
	// (docker/vsphere/unknown); the /infra/list routes accept an optional ?provider=
	// query param and otherwise auto-select the first provider reported installed by
	// /api/infra/capabilities, returning 404 if the resolved/requested provider isn't installed.
	router.HandleFunc("/api/clusters/list", kubernetes.HandleClusterList).Methods("GET")
	router.HandleFunc("/api/clusters/infra/list", kubernetes.HandleClusterInfraList).Methods("GET")

	// Infrastructure provider detection: which of Docker/vSphere are installed in the
	// connected environment, and their version (see specs/004-detect-infra-adapt-ui).
	router.HandleFunc("/api/infra/capabilities", kubernetes.HandleInfraCapabilities).Methods("GET")

	// Raw object passthrough for the object detail screens' YAML tree tab: returns the
	// complete Kubernetes object for ?group=&version=&resource=&namespace=&name=, bypassing
	// every curated DTO (see specs/005-object-viz-ai-panel/contracts/raw-object-api.md).
	router.HandleFunc("/api/raw", kubernetes.HandleRawObject).Methods("GET")

	// Cluster API dashboard Handlers
	router.HandleFunc("/api/clusters/info", kubernetes.HandleClusterInfo).Methods("GET")
	router.HandleFunc("/api/clusters/components", kubernetes.HandleComponentsVersion).Methods("GET")
	router.HandleFunc("/api/clusters/summary", kubernetes.HandleSummaryCluster).Methods("GET")
	router.HandleFunc("/api/clusters/classes", kubernetes.HandleClusterClasses).Methods("GET")
	router.HandleFunc("/api/clusters/topology", kubernetes.HandleClusterTopology).Methods("GET")

	// Cluster API Machine Deployments Handlers
	router.HandleFunc("/api/machinesdeployment/list", kubernetes.HandleMachineDeployments).Methods("GET")

	// Cluster API Machine Handlers (same provider dispatch as /api/clusters/infra/list)
	router.HandleFunc("/api/machines/list", kubernetes.HandleMachines).Methods("GET")
	router.HandleFunc("/api/machines/infra/list", kubernetes.HandleMachineInfra).Methods("GET")

	// Start the websocket handlers. Live resource lists (Clusters/Machines listing
	// screens) stream over /ws/watcher, not these REST endpoints — object types include
	// "cluster-infra-docker"/"machine-infra-docker" alongside the existing vSphere ones.
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
		system.HandleChatbot(pool, w, r)
	}).Methods("GET", "OPTIONS")
}
