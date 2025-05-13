package clusterapi

import (
	"context"
	"fmt"
	"strconv"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/fetchers"
	"github.com/knabben/observatio/webserver/internal/infra/models"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	utilnet "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/client-go/rest"
	clusterctlv1 "sigs.k8s.io/cluster-api/cmd/clusterctl/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ClusterSummary defines the summary of cluster states from a kubeconfig.
type ClusterSummary struct {
	// Provisioned stores the number of running clusters.
	ClusterProvisioned int `json:"clusterProvisioned"`

	// Failed stores the number of failing clusters.
	ClusterFailed int `json:"clusterFailed"`

	// MachineDeploymentProvisioned stores the number of running machine deployments within the cluster.
	MachineDeploymentProvisioned int `json:"machineDeploymentProvisioned"`

	// MachineDeploymentFailed stores the number of machine deployments that have failed.
	MachineDeploymentFailed int `json:"machineDeploymentFailed"`

	// MachineProvisioned represents the count of successfully provisioned machines.
	MachineProvisioned int `json:"machineProvisioned"`

	// MachineFailed stores the number of failing machines.
	MachineFailed int `json:"machineFailed"`
}

// GenerateClusterSummary returns the entire cluster summary from a kubeconfig.
func GenerateClusterSummary(ctx context.Context, c client.Client) (summary ClusterSummary, err error) {
	// fetch clusters objects
	var clusterResponse models.ClusterResponse
	if clusterResponse, err = fetchers.FetchClusters(ctx, c); err != nil {
		return summary, err
	}

	// fetch machine deployment objects
	var mdResponse models.MachineDeploymentResponse
	if mdResponse, err = fetchers.FetchMachineDeployment(ctx, c); err != nil {
		return summary, err
	}

	// fetch machine objects
	var machineResponse models.MachineResponse
	if machineResponse, err = fetchers.FetchMachines(ctx, c); err != nil {
		return summary, err
	}

	return ClusterSummary{
		ClusterProvisioned:           clusterResponse.Total - clusterResponse.Failing,
		ClusterFailed:                clusterResponse.Failing,
		MachineProvisioned:           machineResponse.Total - machineResponse.Failing,
		MachineFailed:                machineResponse.Failing,
		MachineDeploymentProvisioned: mdResponse.Total - mdResponse.Failing,
		MachineDeploymentFailed:      mdResponse.Failing,
	}, nil
}

// Services defines the service name and path.
type Services struct {
	// Name is the services name.
	Name string `json:"name"`

	// Path is the accessible core service path.
	Path string `json:"path"`
}

// FindServices returns a list of internal mgmt cluster services.
func FindServices(ctx context.Context, c client.Client, namespace string) ([]Services, error) {
	var (
		servicesList = corev1.ServiceList{}
		cfg          = ctx.Value("config").(*rest.Config)
	)

	services := []Services{
		{
			Name: "control-plane",
			Path: cfg.Host,
		},
	}
	labels := client.MatchingLabels{
		"kubernetes.io/cluster-service": "true",
	}
	if err := c.List(ctx, &servicesList, client.InNamespace(namespace), labels); err != nil {
		return nil, err
	}
	for _, svc := range servicesList.Items {
		var link string
		if len(svc.Status.LoadBalancer.Ingress) > 0 {
			ingress := svc.Status.LoadBalancer.Ingress[0]
			ip := ingress.IP
			if ip == "" {
				ip = ingress.Hostname
			}
			for _, port := range svc.Spec.Ports {
				link += "http://" + ip + ":" + strconv.Itoa(int(port.Port)) + " "
			}
		} else {
			name := svc.Name
			if len(svc.Spec.Ports) > 0 {
				port := svc.Spec.Ports[0]
				scheme := ""
				if port.Name == "https" || port.Port == 443 {
					scheme = "https"
				}
				name = utilnet.JoinSchemeNamePort(scheme, svc.Name, port.Name)
			}
			if len(svc.GroupVersionKind().Group) == 0 {
				link = cfg.Host + "/api" + svc.GroupVersionKind().Version + "/namespaces/" + svc.Namespace + "/services/" + name + "/proxy"
			} else {
				link = cfg.Host + "/api" + svc.GroupVersionKind().Group + "/" + svc.GroupVersionKind().Version + "/namespaces/" + svc.Namespace + "/services/" + name + "/proxy"
			}
		}
		name := svc.Labels["kubernetes.io/name"]
		if len(name) == 0 {
			name = svc.Name
		}
		services = append(services, Services{Name: name, Path: link})
	}
	return services, nil
}

// Components stores the internal component details.
type Components struct {
	// Name is the name of the component.
	Name string `json:"name"`

	// ProviderName is the provider name of the component.
	ProviderName string `json:"providerName"`

	// Kind is the CAPI type of component.
	Kind string `json:"kind"`

	// Version is the CAPI version of the component.
	Version string `json:"version"`
}

// GenerateComponentVersions return a list of ClusterAPI components and versions.
func GenerateComponentVersions(ctx context.Context, c client.Client) (components []Components, err error) {
	var providers clusterctlv1.ProviderList
	if err := c.List(ctx, &providers); err != nil {
		return components, err
	}

	for _, r := range providers.Items {
		components = append(components, Components{
			ProviderName: r.ProviderName, Name: r.Name, Kind: r.Type, Version: r.Version,
		})
	}
	return components, nil
}

type ClusterTopology struct {
}

func GenerateClusterTopology(ctx context.Context, c client.Client) (topology ClusterTopology, err error) {
	machines, err := fetchers.ListMachineInfra(ctx, c)
	if err != nil {
		return topology, err
	}

	for _, machine := range machines {
		gvk := machine.GroupVersionKind() // assuming 'machine' is the instance you need the GVK from
		gvr, _ := meta.UnsafeGuessKindToResource(gvk)
		err := fetchOwnerHierarchy(ctx, machine.OwnerReferences, gvr, machine.Namespace, machine.Name)
		if err != nil {
			return topology, err
		}
		fmt.Println("----")
	}

	return topology, err
}
