package fetchers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Test_isPodReady(t *testing.T) {
	ready := unstructured.Unstructured{Object: map[string]interface{}{
		"status": map[string]interface{}{
			"conditions": []interface{}{
				map[string]interface{}{"type": "Ready", "status": "True"},
			},
		},
	}}
	assert.True(t, isPodReady(ready))

	notReady := unstructured.Unstructured{Object: map[string]interface{}{
		"status": map[string]interface{}{
			"conditions": []interface{}{
				map[string]interface{}{"type": "Ready", "status": "False"},
			},
		},
	}}
	assert.False(t, isPodReady(notReady))

	noConditions := unstructured.Unstructured{Object: map[string]interface{}{"status": map[string]interface{}{}}}
	assert.False(t, isPodReady(noConditions))
}

func Test_podWaitingReason(t *testing.T) {
	crashLooping := unstructured.Unstructured{Object: map[string]interface{}{
		"status": map[string]interface{}{
			"containerStatuses": []interface{}{
				map[string]interface{}{"state": map[string]interface{}{"waiting": map[string]interface{}{"reason": "CrashLoopBackOff"}}}},
		},
	}}
	assert.Equal(t, "CrashLoopBackOff", podWaitingReason(crashLooping))

	running := unstructured.Unstructured{Object: map[string]interface{}{
		"status": map[string]interface{}{
			"containerStatuses": []interface{}{
				map[string]interface{}{"state": map[string]interface{}{"running": map[string]interface{}{}}}},
		},
	}}
	assert.Equal(t, "", podWaitingReason(running))
}
