package fetchers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterctlv1 "sigs.k8s.io/cluster-api/cmd/clusterctl/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var (
	scheme = runtime.NewScheme()
	_      = clusterctlv1.AddToScheme(scheme) // Register Cluster API types
	_      = clusterv1.AddToScheme(scheme)    // Register Cluster API types
	_      = corev1.AddToScheme(scheme)
)

func Test_ClusterList(t *testing.T) {
	var clusters clusterv1.ClusterList
	replicas := int32(1)
	enabled := true
	tests := []struct {
		cluster clusterv1.Cluster
	}{
		{
			cluster: clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster",
					Namespace: "kube-system",
				},
				Spec: clusterv1.ClusterSpec{
					Topology: &clusterv1.Topology{
						Class: "fake-clusterclass",
						ControlPlane: clusterv1.ControlPlaneTopology{
							Replicas:           &replicas,
							MachineHealthCheck: &clusterv1.MachineHealthCheckTopology{Enable: &enabled},
						},
					},
				},
				Status: clusterv1.ClusterStatus{
					Phase: string(clusterv1.ClusterPhase(clusterv1.ClusterPhaseProvisioned)),
				},
			},
		},
		{
			cluster: clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster-2",
					Namespace: "kube-system",
				},
				Spec: clusterv1.ClusterSpec{
					Topology: &clusterv1.Topology{
						Class: "fake-clusterclass",
						ControlPlane: clusterv1.ControlPlaneTopology{
							Replicas:           &replicas,
							MachineHealthCheck: &clusterv1.MachineHealthCheckTopology{Enable: &enabled},
						},
					},
				},
				Status: clusterv1.ClusterStatus{
					Phase: string(clusterv1.ClusterPhase(clusterv1.ClusterPhaseProvisioned)),
				},
			},
		},
	}
	for _, tt := range tests {
		var c = fake.NewClientBuilder().
			WithScheme(scheme).
			WithRuntimeObjects(&tt.cluster).
			WithLists(&clusters).
			Build()
		response, err := FetchClusters(context.Background(), c)
		assert.NoError(t, err)
		assert.Equal(t, 1, response.Total, 1)
		assert.Equal(t, 1, response.Failing)
		assert.Len(t, response.Clusters, 1)
		for _, cl := range response.Clusters {
			assert.Equal(t, tt.cluster.Name, cl.Name)
			assert.Equal(t, tt.cluster.Status.Phase, cl.Phase)
		}
	}
}
