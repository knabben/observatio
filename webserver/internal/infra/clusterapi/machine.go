package clusterapi

import (
	"context"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Machine struct {
	Name       string                 `json:"name"`
	Namespace  string                 `json:"namespace"`
	Cluster    string                 `json:"cluster,omitempty"`
	NodeName   string                 `json:"nomeName,omitempty"`
	ProviderID string                 `json:"providerID,omitempty"`
	Version    string                 `json:"version,omitempty"`
	Phase      clusterv1.MachinePhase `json:"phase"`
}

// FetchMachine returns the list of Machines objects from the mgmt cluster.
func FetchMachine(ctx context.Context, c client.Client) (machines []Machine, err error) {
	var machineList clusterv1.MachineList
	if err = c.List(ctx, &machineList); err != nil {
		return machines, err
	}

	for _, m := range machineList.Items {
		machines = append(machines, Machine{
			Name:      m.Name,
			Namespace: m.Namespace,
			Cluster:   m.Spec.ClusterName,
			NodeName:  m.Status.NodeRef.Name,
			Version:   *m.Spec.Version,
			Phase:     clusterv1.MachinePhase(m.Status.Phase),
		})
	}
	return machines, nil
}
