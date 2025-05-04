package processor

import (
	"time"

	capv "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"

	"github.com/knabben/observatio/webserver/internal/infra/models"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// ProcessCluster converts a modeling upstream cluster object.
func ProcessCluster(cl clusterv1.Cluster) (cluster models.Cluster, failed bool) {
	clusterClass := models.ClusterClass{IsClusterClass: false}
	if cl.Spec.Topology != nil {
		clusterClass = models.ClusterClass{
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
	cluster = models.Cluster{
		Name:                cl.Name,
		Namespace:           cl.Namespace,
		Paused:              cl.Spec.Paused,
		ClusterClass:        clusterClass,
		Phase:               cl.Status.Phase,
		InfrastructureReady: cl.Status.InfrastructureReady,
		ControlPlaneReady:   cl.Status.ControlPlaneReady,
		Conditions:          cl.Status.Conditions,
		Created:             time.Now().Sub(cl.ObjectMeta.CreationTimestamp.Time).String(),
	}
	if cl.Spec.ClusterNetwork != nil {
		cluster.PodNetwork = cl.Spec.ClusterNetwork.Pods.String()
		cluster.ServiceNetwork = cl.Spec.ClusterNetwork.Services.String()
	}
	if cl.Status.InfrastructureReady || cl.Status.ControlPlaneReady {
		failed = true
	}
	return cluster, failed
}

// ProcessClusterResponse returns the full response and summary of clusters objects.
func ProcessClusterResponse(clusters []clusterv1.Cluster) models.ClusterResponse {
	var (
		failed      int
		clusterList []models.Cluster
	)
	for _, cl := range clusters {
		cluster, clusterFailed := ProcessCluster(cl)
		if clusterFailed {
			failed += 1
		}
		clusterList = append(clusterList, cluster)
	}
	return models.ClusterResponse{
		Total:    len(clusters),
		Clusters: clusterList,
		Failing:  failed,
	}
}

// ProcessClusterInfra converts a upstream CAPV object into internal model.
func ProcessClusterInfra(cl capv.VSphereCluster) (cluster models.ClusterInfra, failed bool) {
	var clusterOwner string
	for _, owner := range cl.OwnerReferences {
		clusterOwner = owner.Name
	}
	cluster = models.ClusterInfra{
		Name:                 cl.Name,
		Cluster:              clusterOwner,
		Server:               cl.Spec.Server,
		Thumbprint:           cl.Spec.Thumbprint,
		Created:              time.Now().Sub(cl.ObjectMeta.CreationTimestamp.Time).String(),
		ControlPlaneEndpoint: cl.Spec.ControlPlaneEndpoint.String(),
		Modules:              cl.Spec.ClusterModules,
		Conditions:           cl.Status.Conditions,
		Ready:                cl.Status.Ready,
	}
	if !cl.Status.Ready {
		failed = true
	}
	return cluster, failed
}

// ProcessClusterInfraResponse returns the full response and summary of CAPV objects.
func ProcessClusterInfraResponse(clusters []capv.VSphereCluster) models.ClusterInfraResponse {
	var (
		clusterList []models.ClusterInfra
		failed      int
	)
	for _, cl := range clusters {
		cluster, clusterFailed := ProcessClusterInfra(cl)
		if clusterFailed {
			failed += 1
		}
		clusterList = append(clusterList, cluster)
	}
	return models.ClusterInfraResponse{
		Total:    len(clusters),
		Clusters: clusterList,
		Failing:  failed,
	}
}
