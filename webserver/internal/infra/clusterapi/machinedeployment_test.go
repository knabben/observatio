package clusterapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func Test_FetchMachineDeployment(t *testing.T) {
	var machineDeploymentList clusterv1.MachineDeploymentList
	tests := []struct {
		d clusterv1.MachineDeployment
	}{
		{
			d: clusterv1.MachineDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "machine-deployment",
					Namespace: "default",
				},
				Spec: clusterv1.MachineDeploymentSpec{},
				Status: clusterv1.MachineDeploymentStatus{
					Conditions: clusterv1.Conditions{
						{
							Type:   "InfrastructureReady",
							Status: corev1.ConditionTrue,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		var c = fake.NewClientBuilder().
			WithScheme(scheme).
			WithRuntimeObjects(&tt.d).
			WithLists(&machineDeploymentList).
			Build()
		mds, err := FetchMachineDeployments(context.Background(), c)
		assert.NoError(t, err)
		assert.Len(t, mds, 1)
	}
}
