package processor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func Test_ProcessMachineHealthCheck(t *testing.T) {
	maxUnhealthy := intstr.FromString("40%")
	mhc := clusterv1.MachineHealthCheck{
		ObjectMeta: metav1.ObjectMeta{Name: "mhc", Namespace: "default"},
		Spec: clusterv1.MachineHealthCheckSpec{
			ClusterName:  "capi-workload",
			Selector:     metav1.LabelSelector{MatchLabels: map[string]string{"role": "worker"}},
			MaxUnhealthy: &maxUnhealthy,
			UnhealthyConditions: []clusterv1.UnhealthyCondition{
				{Type: corev1.NodeReady, Status: corev1.ConditionFalse, Timeout: metav1.Duration{Duration: 5 * time.Minute}},
			},
			NodeStartupTimeout: &metav1.Duration{Duration: 10 * time.Minute},
		},
		Status: clusterv1.MachineHealthCheckStatus{
			ExpectedMachines: 3, CurrentHealthy: 2, RemediationsAllowed: 1,
		},
	}

	result := ProcessMachineHealthCheck(mhc)

	assert.Equal(t, "mhc", result.Name)
	assert.Equal(t, "capi-workload", result.Cluster)
	assert.Equal(t, "40%", result.MaxUnhealthy)
	assert.Equal(t, "worker", result.Selector.MatchLabels["role"])
	assert.Len(t, result.UnhealthyConditions, 1)
	assert.Equal(t, int32(3), result.Status.ExpectedMachines)
	assert.NotEmpty(t, result.NodeStartupTimeout)
}

func Test_ProcessMachineHealthCheck_NilMaxUnhealthy(t *testing.T) {
	mhc := clusterv1.MachineHealthCheck{
		ObjectMeta: metav1.ObjectMeta{Name: "mhc", Namespace: "default"},
	}

	result := ProcessMachineHealthCheck(mhc)
	assert.Empty(t, result.MaxUnhealthy)
}
