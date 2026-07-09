package processor

import (
	"time"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/knabben/observatio/webserver/internal/infra/models"
)

// ProcessMachineSet transforms a clusterv1.MachineSet into a models.MachineSet.
func ProcessMachineSet(ms clusterv1.MachineSet) models.MachineSet {
	return models.MachineSet{
		ObjectMeta:        ms.ObjectMeta,
		Age:               formatDuration(time.Since(ms.CreationTimestamp.Time)),
		Cluster:           ms.Spec.ClusterName,
		MachineDeployment: ms.Labels["cluster.x-k8s.io/deployment-name"],
		Replicas:          ms.Spec.Replicas,
		Status:            ms.Status,
	}
}
