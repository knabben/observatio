package fetchers

import (
	"context"

	capv "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/processor"
	"github.com/knabben/observatio/webserver/internal/infra/models"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MachineInput is a type constraint that allows types clusterv1.Machine or capv.VSphereMachine.
type MachineInput interface {
	clusterv1.Machine | capv.VSphereMachine
}

// MachineResponse represents a type constraint that encompasses machine response types like MachineResponse and MachineInfraResponse.
type MachineResponse interface {
	models.MachineResponse | models.MachineInfraResponse
}

// MachineFetcher is a function type that retrieves a list of machines from a client within a given context and returns any errors.
type MachineFetcher[T MachineInput] func(ctx context.Context, c client.Client) ([]T, error)

// ProcessMachine is a function type that processes a list of machines of type T and returns a response of type R.
type ProcessMachine[T MachineInput, R MachineResponse] func(machines []T) R

// FetchMachines retrieves a summary of machines in the cluster including their total count, failures, and detailed info.
func FetchMachines(ctx context.Context, c client.Client) (response models.MachineResponse, err error) {
	return ProcessMachines[clusterv1.Machine, models.MachineResponse](
		ctx, c, ListMachines, processor.ProcessMachineResponse,
	)
}

// FetchMachineInfra retrieves and processes infrastructure data for VSphere machines in the cluster context.
func FetchMachineInfra(ctx context.Context, c client.Client) (response models.MachineInfraResponse, err error) {
	return ProcessMachines[capv.VSphereMachine, models.MachineInfraResponse](
		ctx, c, ListMachineInfra, processor.ProcessMachineInfraResponse,
	)
}

// ProcessMachines retrieves a list of machines using the provided fetcher and processes them using the given processor.
func ProcessMachines[T MachineInput, R MachineResponse](
	ctx context.Context, c client.Client, fetcher MachineFetcher[T], process ProcessMachine[T, R],
) (res R, err error) {
	var clusters []T
	if clusters, err = fetcher(ctx, c); err != nil {
		return res, err
	}
	return process(clusters), nil
}

// ListMachines retrieves a list of Machine objects from the Kubernetes cluster using the provided client.
func ListMachines(ctx context.Context, c client.Client) (machines []clusterv1.Machine, err error) {
	var machineList clusterv1.MachineList
	if err = c.List(ctx, &machineList); err != nil {
		return machines, err
	}
	return machineList.Items, nil
}

// ListMachineInfra retrieves a list of VSphereMachine resources using the provided client and context.
func ListMachineInfra(ctx context.Context, c client.Client) (machines []capv.VSphereMachine, err error) {
	var vsphereMachines capv.VSphereMachineList
	if err = c.List(ctx, &vsphereMachines); err != nil {
		return nil, err
	}
	return vsphereMachines.Items, nil
}
