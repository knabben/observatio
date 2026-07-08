package watchers

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/day2ops"
)

func toUnstructuredEvent(t *testing.T, eventType watch.EventType, obj interface{}) day2opsEvent {
	t.Helper()
	content, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	require.NoError(t, err)
	return day2opsEvent{
		event: watch.Event{Type: eventType, Object: &unstructured.Unstructured{Object: content}},
	}
}

func Test_day2opsStore_apply_upsertsAndDeletes(t *testing.T) {
	store := newDay2opsStore()

	cluster := &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{Name: "c1", Namespace: "default"},
		Status:     clusterv1.ClusterStatus{InfrastructureReady: true, ControlPlaneReady: true},
	}
	evt := toUnstructuredEvent(t, watch.Added, cluster)
	evt.gvr = clusterGVR
	require.NoError(t, store.apply(evt))

	clusters, _, _ := store.snapshot()
	assert.Len(t, clusters, 1)
	assert.Equal(t, "c1", clusters[0].Name)

	delEvt := toUnstructuredEvent(t, watch.Deleted, cluster)
	delEvt.gvr = clusterGVR
	require.NoError(t, store.apply(delEvt))

	clusters, _, _ = store.snapshot()
	assert.Len(t, clusters, 0)
}

func Test_day2opsStore_apply_machineDeploymentAndMachine(t *testing.T) {
	store := newDay2opsStore()

	md := &clusterv1.MachineDeployment{
		ObjectMeta: metav1.ObjectMeta{Name: "md1", Namespace: "default"},
		Status:     clusterv1.MachineDeploymentStatus{Replicas: 3, ReadyReplicas: 3},
	}
	mdEvt := toUnstructuredEvent(t, watch.Added, md)
	mdEvt.gvr = machineDeploymentGVR
	require.NoError(t, store.apply(mdEvt))

	machine := &clusterv1.Machine{
		ObjectMeta: metav1.ObjectMeta{Name: "m1", Namespace: "default"},
		Status:     clusterv1.MachineStatus{InfrastructureReady: true, BootstrapReady: true},
	}
	mEvt := toUnstructuredEvent(t, watch.Added, machine)
	mEvt.gvr = machineGVR
	require.NoError(t, store.apply(mEvt))

	_, machineDeployments, machines := store.snapshot()
	assert.Len(t, machineDeployments, 1)
	assert.Len(t, machines, 1)
}

func Test_assembleData_reflectsSnapshotAndUnavailable(t *testing.T) {
	store := newDay2opsStore()
	cluster := &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{Name: "c1", Namespace: "default"},
		Status:     clusterv1.ClusterStatus{InfrastructureReady: false, ControlPlaneReady: true},
	}
	evt := toUnstructuredEvent(t, watch.Added, cluster)
	evt.gvr = clusterGVR
	require.NoError(t, store.apply(evt))

	ctx := context.Background()
	data := assembleData(ctx, nil, nil, store, false)
	assert.False(t, data.SourceUnavailable)
	var clusterRollup day2ops.HealthRollup
	for _, r := range data.Rollups {
		if r.Category == day2ops.CategoryCluster {
			clusterRollup = r
		}
	}
	assert.Equal(t, 1, clusterRollup.Failed)

	unavailableData := assembleData(ctx, nil, nil, store, true)
	assert.True(t, unavailableData.SourceUnavailable)
}

func Test_assembleData_computesDebugPathForFailedMachine(t *testing.T) {
	store := newDay2opsStore()
	machine := &clusterv1.Machine{
		ObjectMeta: metav1.ObjectMeta{Name: "worker-0", Namespace: "default"},
		Status: clusterv1.MachineStatus{
			InfrastructureReady: false,
			BootstrapReady:      true,
			Phase:               "Provisioning",
		},
	}
	evt := toUnstructuredEvent(t, watch.Added, machine)
	evt.gvr = machineGVR
	require.NoError(t, store.apply(evt))

	data := assembleData(context.Background(), nil, nil, store, false)

	require.Len(t, data.DebugPaths, 1)
	assert.Equal(t, "worker-0", data.DebugPaths[0].ObjectRef.Name)
}

func Test_day2opsStore_apply_providerResource(t *testing.T) {
	store := newDay2opsStore()

	dockerMachine := &unstructured.Unstructured{Object: map[string]interface{}{
		"kind":     "DockerMachine",
		"metadata": map[string]interface{}{"name": "infra-0", "namespace": "default"},
		"status":   map[string]interface{}{"ready": true},
	}}
	require.NoError(t, store.apply(day2opsEvent{gvr: dockerMachineGVR, event: watch.Event{Type: watch.Added, Object: dockerMachine}}))

	machine := clusterv1.Machine{
		ObjectMeta: metav1.ObjectMeta{Name: "worker-0", Namespace: "default"},
		Spec: clusterv1.MachineSpec{
			InfrastructureRef: corev1.ObjectReference{Name: "infra-0", Namespace: "default"},
		},
	}
	status := store.providerResourceFor(machine)
	require.NotNil(t, status)
	assert.True(t, status.Ready)
}

func Test_stalledRolloutRisks_flagsOldMachineSet(t *testing.T) {
	store := newDay2opsStore()
	now := metav1.Now()

	oldMS := &clusterv1.MachineSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: "md-old", Namespace: "default",
			Labels:            map[string]string{"cluster.x-k8s.io/deployment-name": "md"},
			CreationTimestamp: metav1.NewTime(now.Add(-time.Hour)),
		},
		Spec: clusterv1.MachineSetSpec{Replicas: int32Ptr(1)},
	}
	newMS := &clusterv1.MachineSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: "md-new", Namespace: "default",
			Labels:            map[string]string{"cluster.x-k8s.io/deployment-name": "md"},
			CreationTimestamp: now,
		},
		Spec: clusterv1.MachineSetSpec{Replicas: int32Ptr(3)},
	}
	require.NoError(t, store.apply(day2opsEvent{gvr: machineSetGVR, event: watch.Event{Type: watch.Added, Object: mustToUnstructured(t, oldMS)}}))
	require.NoError(t, store.apply(day2opsEvent{gvr: machineSetGVR, event: watch.Event{Type: watch.Added, Object: mustToUnstructured(t, newMS)}}))

	md := clusterv1.MachineDeployment{ObjectMeta: metav1.ObjectMeta{Name: "md", Namespace: "default"}}
	risks := stalledRolloutRisks(store, []clusterv1.MachineDeployment{md})

	require.Len(t, risks, 1)
	assert.Equal(t, day2ops.RiskStalledRollout, risks[0].Kind)
}

func Test_driftRisks_flagsGenerationMismatch(t *testing.T) {
	store := newDay2opsStore()
	dockerMachine := &unstructured.Unstructured{Object: map[string]interface{}{
		"kind":     "DockerMachine",
		"metadata": map[string]interface{}{"name": "infra-0", "namespace": "default", "generation": int64(3)},
		"status":   map[string]interface{}{"ready": true, "observedGeneration": int64(1)},
	}}
	require.NoError(t, store.apply(day2opsEvent{gvr: dockerMachineGVR, event: watch.Event{Type: watch.Added, Object: dockerMachine}}))

	risks := driftRisks(store)
	require.Len(t, risks, 1)
	assert.Equal(t, day2ops.RiskDrift, risks[0].Kind)
}

func Test_machineHealthCheckSeverities_flagsMaxUnhealthyBreach(t *testing.T) {
	store := newDay2opsStore()
	mhc := &clusterv1.MachineHealthCheck{
		ObjectMeta: metav1.ObjectMeta{Name: "mhc", Namespace: "default"},
		Status: clusterv1.MachineHealthCheckStatus{
			ExpectedMachines: 3, CurrentHealthy: 1, RemediationsAllowed: 0,
		},
	}
	require.NoError(t, store.apply(day2opsEvent{gvr: machineHealthCheckGVR, event: watch.Event{Type: watch.Added, Object: mustToUnstructured(t, mhc)}}))

	severities := machineHealthCheckSeverities(store)
	require.Len(t, severities, 1)
	assert.Equal(t, day2ops.SeverityNeedsInvestigation, severities[0].Level)
}

func Test_ComputeManagementClusterSeverity_wiredIntoAssembleData(t *testing.T) {
	store := newDay2opsStore()
	data := assembleData(context.Background(), nil, nil, store, true)

	require.Len(t, data.Severities, 1)
	assert.Equal(t, day2ops.SeverityManagementCritical, data.Severities[0].Level)
}

func int32Ptr(i int32) *int32 { return &i }

func mustToUnstructured(t *testing.T, obj interface{}) *unstructured.Unstructured {
	t.Helper()
	content, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	require.NoError(t, err)
	return &unstructured.Unstructured{Object: content}
}
