package clusterapi

import (
	"context"
	"sigs.k8s.io/controller-runtime/pkg/client"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// MachineDeployment stores the presentation model for a CAPI machine deployment
type MachineDeployment struct {
	Name            string                           `json:"name"`
	Namespace       string                           `json:"namespace"`
	Cluster         string                           `json:"cluster"`
	Replicas        int32                            `json:"replicas"`
	ReadyReplicas   int32                            `json:"readyReplicas"`
	UpdatedReplicas int32                            `json:"updatedReplicas"`
	Phase           clusterv1.MachineDeploymentPhase `json:"phase"`
}

func FetchMachineDeployments(ctx context.Context, c client.Client) (clusterMDs []MachineDeployment, err error) {
	var machineDeployments []clusterv1.MachineDeployment
	if machineDeployments, err = listMachineDeployments(ctx, c); err != nil {
		return clusterMDs, err
	}
	for _, md := range machineDeployments {
		clusterMDs = append(clusterMDs, MachineDeployment{
			Name:            md.Name,
			Namespace:       md.Namespace,
			Cluster:         md.Spec.ClusterName,
			Replicas:        md.Status.Replicas,
			ReadyReplicas:   md.Status.ReadyReplicas,
			UpdatedReplicas: md.Status.UpdatedReplicas,
			Phase:           clusterv1.MachineDeploymentPhase(md.Status.Phase),
		})
	}
	return clusterMDs, err
}

func listMachineDeployments(ctx context.Context, c client.Client) (machineDeployments []clusterv1.MachineDeployment, err error) {
	var mdsList clusterv1.MachineDeploymentList
	if err = c.List(ctx, &mdsList); err != nil {
		return machineDeployments, err
	}
	return mdsList.Items, nil
}
