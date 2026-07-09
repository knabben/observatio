package processor

import (
	"time"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	controlplanev1 "sigs.k8s.io/cluster-api/controlplane/kubeadm/api/v1beta1"

	"github.com/knabben/observatio/webserver/internal/infra/models"
)

// ProcessKubeadmControlPlane transforms a controlplanev1.KubeadmControlPlane into a
// models.KubeadmControlPlane.
func ProcessKubeadmControlPlane(kcp controlplanev1.KubeadmControlPlane) models.KubeadmControlPlane {
	return models.KubeadmControlPlane{
		ObjectMeta: kcp.ObjectMeta,
		Age:        formatDuration(time.Since(kcp.CreationTimestamp.Time)),
		Cluster:    kcp.Labels[clusterv1.ClusterNameLabel],
		Version:    kcp.Spec.Version,
		Replicas:   kcp.Spec.Replicas,
		Status:     kcp.Status,
	}
}
