package clusterapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Test_isUnstructuredUnhealthy_StatusReady(t *testing.T) {
	healthy := &unstructured.Unstructured{Object: map[string]interface{}{"status": map[string]interface{}{"ready": true}}}
	assert.False(t, isUnstructuredUnhealthy(healthy))

	unhealthy := &unstructured.Unstructured{Object: map[string]interface{}{"status": map[string]interface{}{"ready": false}}}
	assert.True(t, isUnstructuredUnhealthy(unhealthy))
}

func Test_isUnstructuredUnhealthy_ReadyCondition(t *testing.T) {
	obj := &unstructured.Unstructured{Object: map[string]interface{}{
		"status": map[string]interface{}{
			"conditions": []interface{}{
				map[string]interface{}{"type": "Ready", "status": "False"},
			},
		},
	}}
	assert.True(t, isUnstructuredUnhealthy(obj))
}

func Test_isUnstructuredUnhealthy_NoSignal(t *testing.T) {
	obj := &unstructured.Unstructured{Object: map[string]interface{}{"status": map[string]interface{}{}}}
	assert.False(t, isUnstructuredUnhealthy(obj), "no ready/condition signal must not be treated as a failure")
}
