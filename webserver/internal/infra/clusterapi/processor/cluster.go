package processor

import (
	"time"

	capv "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"

	"github.com/knabben/observatio/webserver/internal/infra/models"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// ProcessCluster transforms a clusterv1.Cluster object into a models.Cluster representation.
func ProcessCluster(cl clusterv1.Cluster) (cluster models.Cluster) {
	clusterClass := models.ClusterClassType{IsClusterClass: false}
	if cl.Spec.Topology != nil {
		clusterClass = models.ClusterClassType{
			IsClusterClass:       true,
			ClassName:            cl.Spec.Topology.Class,
			ClassNamespace:       cl.Spec.Topology.ClassNamespace,
			KubernetesVersion:    cl.Spec.Topology.Version,
			ControlPlaneReplicas: *cl.Spec.Topology.ControlPlane.Replicas,
		}
		if cl.Spec.Topology.ControlPlane.MachineHealthCheck != nil {
			clusterClass.ControlPlaneMHC = true
		}
		if cl.Spec.Topology.Workers != nil {
			clusterClass.WorkersMachineDeployments = cl.Spec.Topology.Workers.MachineDeployments
		}
	}
	var clusterNetwork clusterv1.ClusterNetwork
	if cl.Spec.ClusterNetwork != nil {
		clusterNetwork = *cl.Spec.ClusterNetwork
	}
	cluster = models.Cluster{
		ObjectMeta:           cl.ObjectMeta,
		Paused:               cl.Spec.Paused,
		Topology:             clusterClass,
		ClusterNetwork:       clusterNetwork,
		ControlPlaneEndpoint: cl.Spec.ControlPlaneEndpoint,
		ControlPlaneRef:      cl.Spec.ControlPlaneRef,
		InfrastructureRef:    cl.Spec.InfrastructureRef,
		Age:                  formatDuration(time.Since(cl.CreationTimestamp.Time)),
		Status:               cl.Status,
	}
	return cluster
}

func ProcessClusterResponse(clusters []clusterv1.Cluster) models.ClusterResponse {
	failedClusterCount := 0
	clusterList := make([]models.Cluster, 0, len(clusters))
	for _, cl := range clusters {
		clusterList = append(clusterList, ProcessCluster(cl))
		if isClusterFailed(cl) {
			failedClusterCount++
		}
	}

	return models.ClusterResponse{
		Total:    len(clusters),
		Clusters: clusterList,
		Failing:  failedClusterCount,
	}
}

func isClusterFailed(cl clusterv1.Cluster) bool {
	return !cl.Status.InfrastructureReady || !cl.Status.ControlPlaneReady
}

// ProcessClusterInfra processes a VSphereCluster object into a ClusterInfra model for consistent infrastructure representation.
func ProcessClusterInfra(cl capv.VSphereCluster) models.ClusterInfra {
	var clusterOwner string
	for _, owner := range cl.OwnerReferences {
		clusterOwner = owner.Name
	}
	return models.ClusterInfra{
		Name:                 cl.Name,
		Namespace:            cl.Namespace,
		Cluster:              clusterOwner,
		Server:               cl.Spec.Server,
		Thumbprint:           cl.Spec.Thumbprint,
		Created:              formatDuration(time.Since(cl.CreationTimestamp.Time)),
		ControlPlaneEndpoint: cl.Spec.ControlPlaneEndpoint.String(),
		Modules:              cl.Spec.ClusterModules,
		Conditions:           cl.Status.Conditions,
		Ready:                cl.Status.Ready,
	}
}

// ProcessClusterInfraResponse generates a response of ClusterInfra models by processing a list of VSphereCluster objects.
func ProcessClusterInfraResponse(clusters []capv.VSphereCluster) models.ClusterInfraResponse {
	failed := 0
	clusterList := make([]models.ClusterInfra, 0, len(clusters))
	for _, cl := range clusters {
		clusterList = append(clusterList, ProcessClusterInfra(cl))
		if !cl.Status.Ready {
			failed++
		}
	}
	return models.ClusterInfraResponse{
		Total:    len(clusters),
		Clusters: clusterList,
		Failing:  failed,
	}
}
