package fetchers

import (
	"context"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/processor"
	"github.com/knabben/observatio/webserver/internal/infra/models"
)

// FetchMachineDeployment returns and process the MD list response
func FetchMachineDeployment(ctx context.Context, c client.Client) (response models.MachineDeploymentResponse, err error) {
	var mds []clusterv1.MachineDeployment
	if mds, err = ListMachineDeployment(ctx, c); err != nil {
		return response, err
	}
	return processor.ProcessMachineDeployment(mds), nil
}

func ListMachineDeployment(ctx context.Context, c client.Client) (machineDeployments []clusterv1.MachineDeployment, err error) {
	var mdsList clusterv1.MachineDeploymentList
	if err = c.List(ctx, &mdsList); err != nil {
		return machineDeployments, err
	}
	return mdsList.Items, nil
}
