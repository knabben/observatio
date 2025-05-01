package processor

import (
	"github.com/knabben/observatio/webserver/internal/infra/models"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"time"
)

func ProcessMachineDeployment(machineDeployments []clusterv1.MachineDeployment) models.MachineDeploymentResponse {
	var clusterMDs []models.MachineDeployment
	var failed int
	for _, md := range machineDeployments {
		clusterMDs = append(clusterMDs, models.MachineDeployment{
			Name:                md.Name,
			Namespace:           md.Namespace,
			Cluster:             md.Spec.ClusterName,
			Replicas:            md.Status.Replicas,
			ReadyReplicas:       md.Status.ReadyReplicas,
			UpdatedReplicas:     md.Status.UpdatedReplicas,
			UnavailableReplicas: md.Status.UnavailableReplicas,
			Created:             time.Now().Sub(md.ObjectMeta.CreationTimestamp.Time).String(),
			Phase:               clusterv1.MachineDeploymentPhase(md.Status.Phase),
		})
		if md.Status.UnavailableReplicas > 0 {
			failed += 1
		}
	}
	return models.MachineDeploymentResponse{
		Total:              len(machineDeployments),
		Failing:            failed,
		MachineDeployments: clusterMDs,
	}
}
