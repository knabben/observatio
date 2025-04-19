package clusterapi

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	utilnet "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/client-go/rest"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterctlv1 "sigs.k8s.io/cluster-api/cmd/clusterctl/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
)

// ClusterSummary defines the summary of cluster states from a kubeconfig.
type ClusterSummary struct {
	// Provisioned stores the number of running clusters.
	Provisioned int `json:"provisioned"`

	// Failed stores the number of failing clusters.
	Failed int `json:"failed"`
}

// GenerateClusterSummary returns the entire cluster summary from a kubeconfig.
func GenerateClusterSummary(ctx context.Context, c client.Client) (summary ClusterSummary, err error) {
	var (
		provisioned, failed int
		clusters            []clusterv1.Cluster
	)
	if clusters, err = fetchClusters(ctx, c); err != nil {
		return summary, err
	}
	for _, cluster := range clusters {
		if clusterv1.ClusterPhase(cluster.Status.Phase) == clusterv1.ClusterPhaseProvisioned {
			provisioned += 1
			continue
		}
		failed += 1
	}
	return ClusterSummary{Provisioned: provisioned, Failed: failed}, nil
}

// Cluster stores the presentation model for a CAPI cluster
type Cluster struct {
	Name        string               `json:"name"`
	HasTopology bool                 `json:"hasTopology"`
	Conditions  clusterv1.Conditions `json:"conditions"`
}

func GenerateClusterList(ctx context.Context, c client.Client) (clusterList []Cluster, err error) {
	var clusters []clusterv1.Cluster
	if clusters, err = fetchClusters(ctx, c); err != nil {
		return clusterList, err
	}
	for _, cl := range clusters {
		hasTopology := false
		if cl.Spec.Topology != nil {
			hasTopology = true
		}
		clusterList = append(clusterList, Cluster{
			Name:        cl.Name,
			HasTopology: hasTopology,
			Conditions:  cl.Status.Conditions,
		})
	}
	return clusterList, err
}

func fetchClusters(ctx context.Context, c client.Client) (clusters []clusterv1.Cluster, err error) {
	var clusterList clusterv1.ClusterList
	if err = c.List(ctx, &clusterList); err != nil {
		return nil, err
	}
	return clusterList.Items, nil
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
			name := svc.ObjectMeta.Name
			if len(svc.Spec.Ports) > 0 {
				port := svc.Spec.Ports[0]
				scheme := ""
				if port.Name == "https" || port.Port == 443 {
					scheme = "https"
				}
				name = utilnet.JoinSchemeNamePort(scheme, svc.ObjectMeta.Name, port.Name)
			}
			if len(svc.GroupVersionKind().Group) == 0 {
				link = cfg.Host + "/api" + svc.GroupVersionKind().Version + "/namespaces/" + svc.ObjectMeta.Namespace + "/services/" + name + "/proxy"
			} else {
				link = cfg.Host + "/api" + svc.GroupVersionKind().Group + "/" + svc.GroupVersionKind().Version + "/namespaces/" + svc.ObjectMeta.Namespace + "/services/" + name + "/proxy"
			}
		}
		name := svc.ObjectMeta.Labels["kubernetes.io/name"]
		if len(name) == 0 {
			name = svc.ObjectMeta.Name
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
