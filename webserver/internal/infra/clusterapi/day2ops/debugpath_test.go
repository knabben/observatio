package day2ops

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func layerByName(layers []DebugLayer, name DebugLayerName) DebugLayer {
	for _, l := range layers {
		if l.Layer == name {
			return l
		}
	}
	return DebugLayer{}
}

func Test_ComputeMachineDebugPath_ProviderResourceImplicated(t *testing.T) {
	m := clusterv1.Machine{
		ObjectMeta: metav1.ObjectMeta{Name: "worker-0", Namespace: "default"},
		Status: clusterv1.MachineStatus{
			Phase: "Provisioning",
			Conditions: clusterv1.Conditions{
				{Type: clusterv1.InfrastructureReadyCondition, Status: corev1.ConditionFalse, Reason: "WaitingForInfrastructure"},
			},
		},
	}
	provider := &ProviderResourceStatus{Kind: "DockerMachine", Name: "worker-0", Ready: false, Message: "VM creation failed"}

	path := ComputeMachineDebugPath(ObjectRef{Resource: "machines", Namespace: "default", Name: "worker-0"}, m, provider, nil)

	assert.Equal(t, LayerStatusImplicated, layerByName(path.Layers, LayerConditions).Status)
	assert.Equal(t, LayerStatusImplicated, layerByName(path.Layers, LayerPhase).Status)
	assert.Equal(t, LayerStatusImplicated, layerByName(path.Layers, LayerProviderResource).Status)
	// Higher layers already explain the failure: controller_activity must stay inconclusive/empty.
	controllerLayer := layerByName(path.Layers, LayerControllerActivity)
	assert.Equal(t, LayerStatusInconclusive, controllerLayer.Status)
	assert.Empty(t, controllerLayer.Evidence)
	assert.NotEmpty(t, path.Summary)
}

func Test_ComputeMachineDebugPath_ControllerActivityOnlyWhenHigherLayersInconclusive(t *testing.T) {
	m := clusterv1.Machine{
		ObjectMeta: metav1.ObjectMeta{Name: "worker-1", Namespace: "default"},
		Status:     clusterv1.MachineStatus{}, // no phase, no conditions: everything inconclusive
	}

	withoutEvents := ComputeMachineDebugPath(ObjectRef{Name: "worker-1"}, m, nil, nil)
	assert.Equal(t, LayerStatusInconclusive, layerByName(withoutEvents.Layers, LayerControllerActivity).Status)

	withEvents := ComputeMachineDebugPath(ObjectRef{Name: "worker-1"}, m, nil, []string{"Warning FailedCreate: quota exceeded"})
	controllerLayer := layerByName(withEvents.Layers, LayerControllerActivity)
	assert.Equal(t, LayerStatusImplicated, controllerLayer.Status)
	assert.Equal(t, []string{"Warning FailedCreate: quota exceeded"}, controllerLayer.Evidence)
}

func Test_ComputeMachineDebugPath_HealthyLayers(t *testing.T) {
	m := clusterv1.Machine{
		ObjectMeta: metav1.ObjectMeta{Name: "worker-2", Namespace: "default"},
		Status: clusterv1.MachineStatus{
			Phase: "Running",
			Conditions: clusterv1.Conditions{
				{Type: clusterv1.ReadyCondition, Status: corev1.ConditionTrue},
			},
		},
	}
	provider := &ProviderResourceStatus{Kind: "DockerMachine", Name: "worker-2", Ready: true}

	path := ComputeMachineDebugPath(ObjectRef{Name: "worker-2"}, m, provider, nil)

	assert.Equal(t, LayerStatusOK, layerByName(path.Layers, LayerConditions).Status)
	assert.Equal(t, LayerStatusOK, layerByName(path.Layers, LayerPhase).Status)
	assert.Equal(t, LayerStatusOK, layerByName(path.Layers, LayerProviderResource).Status)
}

func Test_ComputeMachineDebugPath_LayersAreOrderedAndLabeled(t *testing.T) {
	path := ComputeMachineDebugPath(ObjectRef{Name: "worker-3"}, clusterv1.Machine{}, nil, nil)
	assert.Len(t, path.Layers, 4)
	assert.Equal(t, []DebugLayerName{LayerConditions, LayerPhase, LayerProviderResource, LayerControllerActivity}, []DebugLayerName{
		path.Layers[0].Layer, path.Layers[1].Layer, path.Layers[2].Layer, path.Layers[3].Layer,
	})
}
