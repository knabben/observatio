package day2ops

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func Test_ComputeRollups(t *testing.T) {
	healthyCluster := clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{Name: "healthy", Namespace: "default"},
		Status:     clusterv1.ClusterStatus{InfrastructureReady: true, ControlPlaneReady: true},
	}
	failedCluster := clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{Name: "failed", Namespace: "default"},
		Status:     clusterv1.ClusterStatus{InfrastructureReady: false, ControlPlaneReady: true},
	}
	healthyMD := clusterv1.MachineDeployment{
		ObjectMeta: metav1.ObjectMeta{Name: "healthy-md", Namespace: "default"},
		Status:     clusterv1.MachineDeploymentStatus{Replicas: 3, ReadyReplicas: 3},
	}
	stalledMD := clusterv1.MachineDeployment{
		ObjectMeta: metav1.ObjectMeta{Name: "stalled-md", Namespace: "default"},
		Status:     clusterv1.MachineDeploymentStatus{Replicas: 3, ReadyReplicas: 1},
	}
	healthyMachine := clusterv1.Machine{
		ObjectMeta: metav1.ObjectMeta{Name: "healthy-machine", Namespace: "default"},
		Status:     clusterv1.MachineStatus{InfrastructureReady: true, BootstrapReady: true},
	}
	failedMachine := clusterv1.Machine{
		ObjectMeta: metav1.ObjectMeta{Name: "failed-machine", Namespace: "default"},
		Status:     clusterv1.MachineStatus{InfrastructureReady: true, BootstrapReady: false},
	}

	rollups := ComputeRollups(
		[]clusterv1.Cluster{healthyCluster, failedCluster},
		[]clusterv1.MachineDeployment{healthyMD, stalledMD},
		[]clusterv1.Machine{healthyMachine, failedMachine},
	)

	require := map[Category]HealthRollup{}
	for _, r := range rollups {
		require[r.Category] = r
	}

	assert.Equal(t, HealthRollup{Category: CategoryCluster, Healthy: 1, Failed: 1}, require[CategoryCluster])
	assert.Equal(t, HealthRollup{Category: CategoryMachineDeployment, Healthy: 1, Failed: 1}, require[CategoryMachineDeployment])
	assert.Equal(t, HealthRollup{Category: CategoryMachine, Healthy: 1, Failed: 1}, require[CategoryMachine])
}

func Test_ComputeRollups_AllHealthy(t *testing.T) {
	cluster := clusterv1.Cluster{Status: clusterv1.ClusterStatus{InfrastructureReady: true, ControlPlaneReady: true}}
	rollups := ComputeRollups([]clusterv1.Cluster{cluster}, nil, nil)

	for _, r := range rollups {
		if r.Category == CategoryCluster {
			assert.Equal(t, 1, r.Healthy)
			assert.Equal(t, 0, r.Failed)
			assert.False(t, r.Unavailable)
		}
	}
}

func Test_ComputeRollups_Empty(t *testing.T) {
	rollups := ComputeRollups(nil, nil, nil)
	assert.Len(t, rollups, 3)
	for _, r := range rollups {
		assert.Equal(t, 0, r.Healthy)
		assert.Equal(t, 0, r.Failed)
	}
}
