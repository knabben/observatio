package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	controlplanev1 "sigs.k8s.io/cluster-api/controlplane/kubeadm/api/v1beta1"
)

func Test_ProcessKubeadmControlPlane(t *testing.T) {
	replicas := int32(3)
	kcp := controlplanev1.KubeadmControlPlane{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kcp", Namespace: "default",
			Labels: map[string]string{clusterv1.ClusterNameLabel: "capi-mgmt"},
		},
		Spec: controlplanev1.KubeadmControlPlaneSpec{
			Version:  "v1.31.0",
			Replicas: &replicas,
		},
		Status: controlplanev1.KubeadmControlPlaneStatus{
			Replicas: 3, ReadyReplicas: 3, UpdatedReplicas: 3, Initialized: true, Ready: true,
			Conditions: clusterv1.Conditions{{Type: "EtcdClusterHealthy", Status: "True"}},
		},
	}

	result := ProcessKubeadmControlPlane(kcp)

	assert.Equal(t, "kcp", result.Name)
	assert.Equal(t, "capi-mgmt", result.Cluster)
	assert.Equal(t, "v1.31.0", result.Version)
	assert.Equal(t, int32(3), *result.Replicas)
	assert.Equal(t, int32(3), result.Status.ReadyReplicas)
	assert.True(t, result.Status.Ready)
	assert.Len(t, result.Status.Conditions, 1)
}

func Test_ProcessKubeadmControlPlane_NilReplicas(t *testing.T) {
	kcp := controlplanev1.KubeadmControlPlane{
		ObjectMeta: metav1.ObjectMeta{Name: "kcp", Namespace: "default"},
	}

	result := ProcessKubeadmControlPlane(kcp)
	assert.Nil(t, result.Replicas)
	assert.Empty(t, result.Cluster)
}
