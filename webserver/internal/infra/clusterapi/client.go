package clusterapi

import (
	"context"

	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// NewClient returns a client and configuration for cluster accessibility.
func NewClient() (client.Client, *rest.Config, error) {
	cfg := config.GetConfigOrDie()
	cli, err := client.New(cfg, client.Options{})
	if err != nil {
		return nil, nil, err
	}
	return cli, cfg, nil
}

// NewClientWithScheme returns the client with scheme.
func NewClientWithScheme(ctx context.Context, scheme *runtime.Scheme) (client.Client, error) {
	cfg := ctx.Value("config").(*rest.Config)
	return client.New(cfg, client.Options{Scheme: scheme})
}

// NewDynamicClient returns a dynamic client for unstructured object access.
func NewDynamicClient(ctx context.Context) (*dynamic.DynamicClient, error) {
	return dynamic.NewForConfig(ctx.Value("config").(*rest.Config))
}

// NewDiscoveryClient returns a DiscoveryClient configured with provided rest.Config.
func NewDiscoveryClient(ctx context.Context) (*discovery.DiscoveryClient, error) {
	return discovery.NewDiscoveryClientForConfig(ctx.Value("config").(*rest.Config))
}

// NewAPIExtensionsClient returns a clientset for reading CustomResourceDefinitions, used by the
// Day-2 Ops dashboard's version-skew detection (specs/006-day2-ops-dashboard/research.md R6).
func NewAPIExtensionsClient(ctx context.Context) (*apiextensionsclientset.Clientset, error) {
	return apiextensionsclientset.NewForConfig(ctx.Value("config").(*rest.Config))
}

// NewClientset returns the typed core Kubernetes clientset, used for the Pod-log subresource
// (the same one `kubectl logs` uses) backing the Day-2 Ops Logs view (research.md R10).
func NewClientset(ctx context.Context) (*kubernetes.Clientset, error) {
	return kubernetes.NewForConfig(ctx.Value("config").(*rest.Config))
}
