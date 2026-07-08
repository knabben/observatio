package day2ops

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Test_ExtractProviderResourceStatus_Ready(t *testing.T) {
	u := &unstructured.Unstructured{Object: map[string]interface{}{
		"kind":     "DockerMachine",
		"metadata": map[string]interface{}{"name": "worker-0"},
		"status":   map[string]interface{}{"ready": true},
	}}

	status := ExtractProviderResourceStatus(u)
	assert.Equal(t, "DockerMachine", status.Kind)
	assert.Equal(t, "worker-0", status.Name)
	assert.True(t, status.Ready)
	assert.Empty(t, status.Message)
}

func Test_ExtractProviderResourceStatus_NotReadyWithConditionMessage(t *testing.T) {
	u := &unstructured.Unstructured{Object: map[string]interface{}{
		"kind":     "DockerMachine",
		"metadata": map[string]interface{}{"name": "worker-1"},
		"status": map[string]interface{}{
			"ready": false,
			"conditions": []interface{}{
				map[string]interface{}{"type": "Ready", "status": "False", "reason": "FailedCreate", "message": "VM creation failed"},
			},
		},
	}}

	status := ExtractProviderResourceStatus(u)
	assert.False(t, status.Ready)
	assert.Equal(t, "VM creation failed", status.Message)
}

func Test_ExtractProviderResourceStatus_NoConditions(t *testing.T) {
	u := &unstructured.Unstructured{Object: map[string]interface{}{
		"kind":     "VSphereMachine",
		"metadata": map[string]interface{}{"name": "worker-2"},
		"status":   map[string]interface{}{"ready": false},
	}}

	status := ExtractProviderResourceStatus(u)
	assert.False(t, status.Ready)
	assert.Empty(t, status.Message)
}
