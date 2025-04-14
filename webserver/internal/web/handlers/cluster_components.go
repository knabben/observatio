package handlers

import (
	"encoding/json"
	"net/http"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"

	clusterctlv1 "sigs.k8s.io/cluster-api/cmd/clusterctl/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	scheme = runtime.NewScheme()
	_      = clusterctlv1.AddToScheme(scheme) // Register Cluster API types
)

type components struct {
	Name    string `json:"name"`
	Kind    string `json:"kind"`
	Version string `json:"version"`
}

func handleComponentsVersion(w http.ResponseWriter, r *http.Request) {
	cfg := r.Context().Value("config").(*rest.Config)
	cli, err := client.New(cfg, client.Options{Scheme: scheme})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var providers clusterctlv1.ProviderList
	if err = cli.List(r.Context(), &providers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response []components
	for _, r := range providers.Items {
		response = append(response, components{Name: r.Name, Kind: r.Type, Version: r.Version})
	}

	if response, err := json.Marshal(&response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}
