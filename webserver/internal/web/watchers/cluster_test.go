package watchers

import (
	"context"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/knabben/observatio/webserver/internal/infra/clusterapi/processor"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func TestWatchClusters(t *testing.T) {
	tests := []struct {
		name        string
		objType     string
		mockConn    *websocket.Conn
		mockEvent   runtime.Object
		expectedErr bool
	}{
		{
			name:     "Valid Cluster Object",
			objType:  "clusters",
			mockConn: &websocket.Conn{},
			mockEvent: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "cluster.x-k8s.io/v1beta1",
					"kind":       "Cluster",
					"metadata": map[string]interface{}{
						"name":      "test-cluster",
						"namespace": "default",
					},
				},
			},
			expectedErr: false,
		},
		{
			name:        "Invalid Object Type",
			objType:     "invalid-type",
			mockConn:    &websocket.Conn{},
			mockEvent:   &unstructured.Unstructured{},
			expectedErr: true,
		},
		{
			name:     "Missing API Version and Kind",
			objType:  "clusters",
			mockConn: &websocket.Conn{},
			mockEvent: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"metadata": map[string]interface{}{
						"name":      "test-cluster",
						"namespace": "default",
					},
				},
			},
			expectedErr: true,
		},
		{
			name:        "Nil Event Object",
			objType:     "clusters",
			mockConn:    &websocket.Conn{},
			mockEvent:   nil,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConverter := func(event runtime.Object) (any, error) {
				if event == nil {
					return nil, runtime.NewMissingKindErr("")
				}

				var cluster clusterv1.Cluster
				unstructuredObj := event.(*unstructured.Unstructured)
				err := runtime.DefaultUnstructuredConverter.FromUnstructured(
					unstructuredObj.UnstructuredContent(), &cluster,
				)
				if err != nil {
					return nil, err
				}
				return processor.ProcessCluster(cluster), nil
			}

			mockGVR := schema.GroupVersionResource{
				Group:    "cluster.x-k8s.io",
				Version:  "v1beta1",
				Resource: "clusters",
			}

			config := WebSocketWatchConfig{
				ObjectType: tt.objType,
				Conn:       tt.mockConn,
				Converter:  mockConverter,
				GVR:        mockGVR,
			}

			err := WatchResourceViaWebSocket(context.Background(), config)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
