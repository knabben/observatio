package processor

import (
	"time"

	"github.com/knabben/observatio/webserver/internal/infra/models"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// ProcessMachineDeployment converts a clusterv1.MachineDeployment object into a models.MachineDeployment structure.
func ProcessMachineDeployment(md clusterv1.MachineDeployment) models.MachineDeployment {
	return models.MachineDeployment{
		ObjectMeta:                md.ObjectMeta,
		Replicas:                  md.Status.Replicas,
		Cluster:                   md.Spec.ClusterName,
		Age:                       formatDuration(time.Since(md.CreationTimestamp.Time)),
		TemplateBootstrap:         md.Spec.Template.Spec.Bootstrap,
		TemplateInfrastructureRef: md.Spec.Template.Spec.InfrastructureRef,
		TemplateVersion:           md.Spec.Template.Spec.Version,
		Status:                    md.Status,
	}
}

func ProcessMachineDeploymentResponse(machineDeployments []clusterv1.MachineDeployment) models.MachineDeploymentResponse {
	failed := 0
	clusterMDs := make([]models.MachineDeployment, 0, len(machineDeployments))
	for _, md := range machineDeployments {
		clusterMDs = append(clusterMDs, ProcessMachineDeployment(md))
		if md.Status.ReadyReplicas != md.Status.Replicas { // nolint
			failed += 1
		}
	}
	return models.MachineDeploymentResponse{
		Total:              len(machineDeployments),
		Failing:            failed,
		MachineDeployments: clusterMDs,
	}
}
