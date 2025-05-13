package clusterapi

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
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
