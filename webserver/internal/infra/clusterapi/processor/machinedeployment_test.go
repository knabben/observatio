package processor

import (
	"testing"

	"github.com/knabben/observatio/webserver/internal/infra/models"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func TestProcessMachineDeployment(t *testing.T) {
	tests := []struct {
		name     string
		input    clusterv1.MachineDeployment
		expected models.MachineDeployment
	}{
		{
			name: "fully populated",
			input: clusterv1.MachineDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-deployment",
					Namespace: "test-namespace",
				},
				Spec: clusterv1.MachineDeploymentSpec{
					ClusterName: "test-cluster",
				},
				Status: clusterv1.MachineDeploymentStatus{
					Replicas:        5,
					ReadyReplicas:   3,
					UpdatedReplicas: 4,
					Phase:           string(clusterv1.MachineDeploymentPhaseScalingUp),
				},
			},
			expected: models.MachineDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-deployment",
					Namespace: "test-namespace",
				},
				Cluster:  "test-cluster",
				Replicas: 5,
			},
		},
		{
			name: "zero values",
			input: clusterv1.MachineDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "empty-deployment",
					Namespace: "default",
				},
				Spec:   clusterv1.MachineDeploymentSpec{},
				Status: clusterv1.MachineDeploymentStatus{},
			},
			expected: models.MachineDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "empty-deployment",
					Namespace: "default",
				},
				Cluster:  "",
				Replicas: 0,
			},
		},
		{
			name: "nil phase",
			input: clusterv1.MachineDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "nil-phase",
					Namespace: "ns1",
				},
				Spec: clusterv1.MachineDeploymentSpec{
					ClusterName: "cluster1",
				},
				Status: clusterv1.MachineDeploymentStatus{
					Replicas: 6,
				},
			},
			expected: models.MachineDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "nil-phase",
					Namespace: "ns1",
				},
				Cluster:  "cluster1",
				Replicas: 6,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessMachineDeployment(tt.input)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.Namespace, result.Namespace)
			assert.Equal(t, tt.expected.Cluster, result.Cluster)
			assert.Equal(t, tt.expected.Replicas, result.Replicas)
			//assert.Equal(t, tt.expected.ReadyReplicas, result.ReadyReplicas)
			//assert.Equal(t, tt.expected.UpdatedReplicas, result.UpdatedReplicas)
			//assert.Equal(t, tt.expected.UnavailableReplicas, result.UnavailableReplicas)
			//assert.Equal(t, tt.expected.Phase, result.Phase)
		})
	}
}
