package clusterapi

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func NewClient() (client.Client, *rest.Config, error) {
	cfg := config.GetConfigOrDie()
	cli, err := client.New(cfg, client.Options{})
	if err != nil {
		return nil, nil, err
	}
	return cli, cfg, nil
}

func NewClientWithScheme(ctx context.Context, scheme *runtime.Scheme) (client.Client, error) {
	cfg := ctx.Value("config").(*rest.Config)
	return client.New(cfg, client.Options{Scheme: scheme})
}
