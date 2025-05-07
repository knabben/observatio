package fetchers

import (
	"context"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/processor"
	"github.com/knabben/observatio/webserver/internal/infra/models"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MachineFetcher[T any] func(ctx context.Context, c client.Client) ([]T, error)
type ProcessMachine[T any, R any] func(machines []T) R

// FetchMachines returns and processes the machine list response
func FetchMachines(ctx context.Context, c client.Client) (response models.MachineResponse, err error) {
	return ProcessMachines[clusterv1.Machine, models.MachineResponse](
		ctx, c, ListMachines, processor.ProcessMachineResponse,
	)
}

func ProcessMachines[T any, R any](
	ctx context.Context, c client.Client, fetcher MachineFetcher[T], process ProcessMachine[T, R],
) (res R, err error) {
	var clusters []T
	if clusters, err = fetcher(ctx, c); err != nil {
		return res, err
	}
	return process(clusters), nil
}

func ListMachines(ctx context.Context, c client.Client) (machines []clusterv1.Machine, err error) {
	var machineList clusterv1.MachineList
	if err = c.List(ctx, &machineList); err != nil {
		return machines, err
	}
	return machineList.Items, nil
}
