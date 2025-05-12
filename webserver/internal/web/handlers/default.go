package handlers

import (
	"embed"
	"encoding/json"

	"net/http"

	"k8s.io/apimachinery/pkg/runtime"

	capv "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterctlv1 "sigs.k8s.io/cluster-api/cmd/clusterctl/api/v1alpha3"

	"github.com/gorilla/mux"
)

//go:embed build/*
var bundle embed.FS

var (
	scheme = runtime.NewScheme()
	_      = clusterctlv1.AddToScheme(scheme) // Register Cluster API types
	_      = clusterv1.AddToScheme(scheme)    // Register Cluster API types
	_      = capv.AddToScheme(scheme)         // Register Cluster API types
)

func DefaultHandlers(router *mux.Router, developmentMode bool) {
	// Generic handlers, healthcheck, version, etc.
	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(map[string]bool{"ok": true}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	// Cluster API handlers
	router.HandleFunc("/api/clusters/list", handleClusterList).Methods("GET")
	router.HandleFunc("/api/clusters/infra/list", handleClusterInfraList).Methods("GET")

	// Cluster API dashboard Handlers
	router.HandleFunc("/api/clusters/info", handleClusterInfo).Methods("GET")
	router.HandleFunc("/api/clusters/components", handleComponentsVersion).Methods("GET")
	router.HandleFunc("/api/clusters/summary", handleSummaryCluster).Methods("GET")
	router.HandleFunc("/api/clusters/classes", handleClusterClasses).Methods("GET")

	// Cluster API Machine Deployments Handlers
	router.HandleFunc("/api/machinesdeployment/list", handleMachineDeployments).Methods("GET")

	// Cluster API Machine Handlers
	router.HandleFunc("/api/machines/list", handleMachines).Methods("GET")
	router.HandleFunc("/api/machines/infra/list", handleMachineInfra).Methods("GET")

	// Websocket Handler for object watchers.
	router.HandleFunc("/ws", handleWebsocket)

	// Static React bundle hosting handler
	if !developmentMode {
		spa := SPAHandler{staticFS: bundle, staticPath: "build", indexPath: "dashboard.html"}
		router.PathPrefix("/").Handler(spa)
	}
}
