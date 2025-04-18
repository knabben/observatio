package handlers

import (
	"embed"
	"encoding/json"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"

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
)

func DefaultHandlers(router *mux.Router, developmentMode bool) {
	// Generic handlers, healthcheck, version, etc.
	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	// Cluster API handlers
	router.HandleFunc("/api/clusters/list", handleClusterList).Methods("GET")

	// Cluster API dashboard handlers
	router.HandleFunc("/api/clusters/info", handleClusterInfo).Methods("GET")
	router.HandleFunc("/api/clusters/components", handleComponentsVersion).Methods("GET")
	router.HandleFunc("/api/clusters/summary", handleSummaryCluster).Methods("GET")
	router.HandleFunc("/api/clusters/classes", handleClusterClasses).Methods("GET")

	// Static React bundle hosting handler
	if !developmentMode {
		spa := SPAHandler{staticFS: bundle, staticPath: "build", indexPath: "index.html"}
		router.PathPrefix("/front").Handler(spa)
	}
}
