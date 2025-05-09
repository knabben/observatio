package processor

import (
	"testing"

	"github.com/knabben/observatio/webserver/internal/infra/models"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func TestProcessCluster(t *testing.T) {
	tests := []struct {
		name     string
		cluster  clusterv1.Cluster
		expected models.Cluster
	}{
		{
			name: "Cluster without topology",
			cluster: clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster1",
					Namespace: "default",
				},
				Spec: clusterv1.ClusterSpec{
					Paused: false,
				},
				Status: clusterv1.ClusterStatus{
					Phase:               "Provisioning",
					InfrastructureReady: false,
					ControlPlaneReady:   false,
				},
			},
			expected: models.Cluster{
				Name:                "cluster1",
				Namespace:           "default",
				Paused:              false,
				ClusterClass:        models.ClusterClassType{IsClusterClass: false},
				Phase:               "Provisioning",
				InfrastructureReady: false,
				ControlPlaneReady:   false,
			},
		},
		{
			name: "Cluster with nil topology but paused",
			cluster: clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster3",
					Namespace: "test",
				},
				Spec: clusterv1.ClusterSpec{
					Paused: true,
				},
				Status: clusterv1.ClusterStatus{
					Phase:               "Paused",
					InfrastructureReady: false,
					ControlPlaneReady:   false,
				},
			},
			expected: models.Cluster{
				Name:                "cluster3",
				Namespace:           "test",
				Paused:              true,
				ClusterClass:        models.ClusterClassType{IsClusterClass: false},
				Phase:               "Paused",
				InfrastructureReady: false,
				ControlPlaneReady:   false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessCluster(tt.cluster)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.Namespace, result.Namespace)
			assert.Equal(t, tt.expected.Paused, result.Paused)
			assert.Equal(t, tt.expected.ClusterClass, result.ClusterClass)
			assert.Equal(t, tt.expected.Phase, result.Phase)
			assert.Equal(t, tt.expected.InfrastructureReady, result.InfrastructureReady)
			assert.Equal(t, tt.expected.ControlPlaneReady, result.ControlPlaneReady)
		})
	}
}
