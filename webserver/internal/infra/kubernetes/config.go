package kubernetes

import (
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
