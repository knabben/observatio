package kubernetes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/day2ops"
)

func Test_computeMachineDetailPath_ProviderResourceEvidence(t *testing.T) {
	testScheme := runtime.NewScheme()
	require.NoError(t, clusterv1.AddToScheme(testScheme))

	gvrToListKind := map[schema.GroupVersionResource]string{
		{Group: "infrastructure.cluster.x-k8s.io", Version: "v1beta1", Resource: "dockermachines"}: "DockerMachineList",
	}

	machine := &clusterv1.Machine{
		ObjectMeta: metav1.ObjectMeta{Name: "worker-0", Namespace: "default"},
		Spec: clusterv1.MachineSpec{
			InfrastructureRef: corev1.ObjectReference{Kind: "DockerMachine", Name: "worker-0", Namespace: "default"},
		},
		Status: clusterv1.MachineStatus{Phase: "Provisioning"},
	}

	dockerMachine := &unstructured.Unstructured{Object: map[string]interface{}{
		"apiVersion": "infrastructure.cluster.x-k8s.io/v1beta1",
		"kind":       "DockerMachine",
		"metadata":   map[string]interface{}{"name": "worker-0", "namespace": "default"},
		"status": map[string]interface{}{
			"ready": false,
			"conditions": []interface{}{
				map[string]interface{}{"type": "Ready", "status": "False", "reason": "FailedCreate", "message": "VM creation failed"},
			},
		},
	}}

	dyn := dynamicfake.NewSimpleDynamicClientWithCustomListKinds(testScheme, gvrToListKind, machine, dockerMachine)

	obj, err := dyn.Resource(machineGVR).Namespace("default").Get(context.Background(), "worker-0", metav1.GetOptions{})
	require.NoError(t, err)

	objectRef := day2ops.ObjectRef{Resource: "machines", Namespace: "default", Name: "worker-0"}
	path, err := computeMachineDetailPath(context.Background(), dyn, objectRef, obj)
	require.NoError(t, err)

	var providerLayer day2ops.DebugLayer
	for _, l := range path.Layers {
		if l.Layer == day2ops.LayerProviderResource {
			providerLayer = l
		}
	}
	assert.Equal(t, day2ops.LayerStatusImplicated, providerLayer.Status)
	require.Len(t, providerLayer.Evidence, 1)
	assert.Equal(t, "VM creation failed", providerLayer.Evidence[0])
}

func Test_computeMachineDetailPath_NoProviderRef(t *testing.T) {
	testScheme := runtime.NewScheme()
	require.NoError(t, clusterv1.AddToScheme(testScheme))

	machine := &clusterv1.Machine{
		ObjectMeta: metav1.ObjectMeta{Name: "worker-1", Namespace: "default"},
		Status:     clusterv1.MachineStatus{Phase: "Pending"},
	}
	dyn := dynamicfake.NewSimpleDynamicClient(testScheme, machine)

	obj, err := dyn.Resource(machineGVR).Namespace("default").Get(context.Background(), "worker-1", metav1.GetOptions{})
	require.NoError(t, err)

	objectRef := day2ops.ObjectRef{Resource: "machines", Namespace: "default", Name: "worker-1"}
	path, err := computeMachineDetailPath(context.Background(), dyn, objectRef, obj)
	require.NoError(t, err)
	assert.Len(t, path.Layers, 4)
}
