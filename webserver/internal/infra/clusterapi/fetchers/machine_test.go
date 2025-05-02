package fetchers

import (
	"context"
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi"
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
			WithScheme(clusterapi.scheme).
			WithRuntimeObjects(&tt.m).
			WithLists(&machineList).
			Build()
		machines, err := FetchMachine(context.Background(), c)
		assert.NoError(t, err)
		assert.Len(t, machines, 1)
	}
}
