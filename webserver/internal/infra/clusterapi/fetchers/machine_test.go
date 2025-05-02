package fetchers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func Test_FetchMachine(t *testing.T) {
	var machineList clusterv1.MachineList
	tests := []struct {
		m clusterv1.Machine
	}{
		{
			m: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "machine",
					Namespace: "default",
				},
				Spec: clusterv1.MachineSpec{},
				Status: clusterv1.MachineStatus{
					InfrastructureReady: true,
					BootstrapReady:      true,
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
			WithRuntimeObjects(&tt.m).
			WithLists(&machineList).
			Build()
		machines, err := FetchMachines(context.Background(), c)
		assert.NoError(t, err)
		assert.Equal(t, machines.Total, 1)
		assert.Equal(t, machines.Failing, 0)
		assert.Len(t, machines.Machines, 1)
	}
}
