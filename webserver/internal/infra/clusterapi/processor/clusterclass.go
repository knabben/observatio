package processor

import (
	"github.com/knabben/observatio/webserver/internal/infra/models"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// ProcessClusterClass transforms a clusterv1.ClusterClass into a models.ClusterClass by mapping its metadata and conditions.
func ProcessClusterClass(cc clusterv1.ClusterClass) models.ClusterClass {
	return models.ClusterClass{
		Name:       cc.Name,
		Namespace:  cc.Namespace,
		Generation: cc.Status.ObservedGeneration,
		Conditions: cc.Status.Conditions,
	}
}

// ProcessClusterClassResponse converts a list of ClusterClass objects into a ClusterClassResponse model.
func ProcessClusterClassResponse(ccs []clusterv1.ClusterClass) models.ClusterClassResponse {
	clusterClasses := make([]models.ClusterClass, 0, len(ccs))
	for _, cc := range ccs {
		clusterClasses = append(clusterClasses, ProcessClusterClass(cc))
	}
	return models.ClusterClassResponse{ClusterClasses: clusterClasses}
}
