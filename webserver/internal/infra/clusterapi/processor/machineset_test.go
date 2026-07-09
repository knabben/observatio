package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func Test_ProcessMachineSet(t *testing.T) {
	replicas := int32(3)
	ms := clusterv1.MachineSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: "worker-ms", Namespace: "default",
			Labels: map[string]string{"cluster.x-k8s.io/deployment-name": "worker-md"},
		},
		Spec: clusterv1.MachineSetSpec{
			ClusterName: "capi-workload",
			Replicas:    &replicas,
		},
		Status: clusterv1.MachineSetStatus{
			Replicas: 3, ReadyReplicas: 2, AvailableReplicas: 2,
			Conditions: clusterv1.Conditions{{Type: "Ready", Status: "False"}},
		},
	}

	result := ProcessMachineSet(ms)

	assert.Equal(t, "worker-ms", result.Name)
	assert.Equal(t, "capi-workload", result.Cluster)
	assert.Equal(t, "worker-md", result.MachineDeployment)
	assert.Equal(t, int32(3), *result.Replicas)
	assert.Equal(t, int32(2), result.Status.ReadyReplicas)
	assert.Len(t, result.Status.Conditions, 1)
}

func Test_ProcessMachineSet_Standalone(t *testing.T) {
	ms := clusterv1.MachineSet{
		ObjectMeta: metav1.ObjectMeta{Name: "standalone-ms", Namespace: "default"},
	}

	result := ProcessMachineSet(ms)
	assert.Nil(t, result.Replicas)
	assert.Empty(t, result.MachineDeployment)
}
