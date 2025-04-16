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
	// Running stores the number of running clusters.
	Running int `json:"running"`

	// Failed stores the number of failing clusters.
	Failed int `json:"failed"`
}

// GenerateClusterSummary returns the entire cluster summary from a kubeconfig.
func GenerateClusterSummary(ctx context.Context, c client.Client) (summary ClusterSummary, err error) {
	var (
		running, failed int
		clusterList     clusterv1.ClusterList
	)

	if err = c.List(ctx, &clusterList); err != nil {
		return summary, err
	}
	for _, cluster := range clusterList.Items {
		if clusterv1.MachinePhase(cluster.Status.Phase) == clusterv1.MachinePhaseRunning {
			running += 1
			continue
		}
		failed += 1
	}

	return ClusterSummary{Running: running, Failed: failed}, nil
}

type Services struct {
	Name string `json:"name"`
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

type Components struct {
	Name    string `json:"name"`
	Kind    string `json:"kind"`
	Version string `json:"version"`
}

// GenerateComponentVersions return a list of ClusterAPI components and versions.
func GenerateComponentVersions(ctx context.Context, c client.Client) (components []Components, err error) {
	var providers clusterctlv1.ProviderList
	if err := c.List(ctx, &providers); err != nil {
		return components, err
	}
	for _, r := range providers.Items {
		components = append(components, Components{Name: r.Name, Kind: r.Type, Version: r.Version})
	}
	return components, nil
}
