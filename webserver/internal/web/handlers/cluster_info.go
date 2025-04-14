package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"

	utilnet "k8s.io/apimachinery/pkg/util/net"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Services struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func handleClusterInfo(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		services []Services
	)

	cli := r.Context().Value("client").(client.Client)
	if services, err = findServices(r.Context(), cli, "kube-system"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if response, err := json.Marshal(&services); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

func findServices(ctx context.Context, cli client.Client, namespace string) ([]Services, error) {
	var svcList = corev1.ServiceList{}

	cfg := ctx.Value("config").(*rest.Config)
	services := []Services{{Name: "control-plane", Path: cfg.Host}}
	labels := client.MatchingLabels{"kubernetes.io/cluster-service": "true"}
	if err := cli.List(ctx, &svcList, client.InNamespace(namespace), labels); err != nil {
		return nil, err
	}

	for _, svc := range svcList.Items {
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
