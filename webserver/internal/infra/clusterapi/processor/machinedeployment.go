package processor

import (
	"time"

	"github.com/knabben/observatio/webserver/internal/infra/models"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// ProcessMachineDeployment converts a clusterv1.MachineDeployment object into a models.MachineDeployment structure.
func ProcessMachineDeployment(md clusterv1.MachineDeployment) models.MachineDeployment {
	return models.MachineDeployment{
		Name:                md.Name,
		Namespace:           md.Namespace,
		Cluster:             md.Spec.ClusterName,
		Replicas:            md.Status.Replicas,
		ReadyReplicas:       md.Status.ReadyReplicas,
		UpdatedReplicas:     md.Status.UpdatedReplicas,
		UnavailableReplicas: md.Status.UnavailableReplicas, // nolint
		Created:             time.Since(md.ObjectMeta.CreationTimestamp.Time).String(),
		Phase:               clusterv1.MachineDeploymentPhase(md.Status.Phase),
	}
}

func ProcessMachineDeploymentResponse(machineDeployments []clusterv1.MachineDeployment) models.MachineDeploymentResponse {
	failed := 0
	clusterMDs := make([]models.MachineDeployment, 0, len(machineDeployments))
	for _, md := range machineDeployments {
		clusterMDs = append(clusterMDs, ProcessMachineDeployment(md))
		if md.Status.UnavailableReplicas > 0 { // nolint
			failed += 1
		}
	}
	return models.MachineDeploymentResponse{
		Total:              len(machineDeployments),
		Failing:            failed,
		MachineDeployments: clusterMDs,
	}
}
