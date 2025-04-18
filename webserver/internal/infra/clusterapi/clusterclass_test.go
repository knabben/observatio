package clusterapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func Test_FetchClusterClass(t *testing.T) {
	var clusterClassesList clusterv1.ClusterClassList
	tests := []struct {
		cc clusterv1.ClusterClass
	}{
		{
			cc: clusterv1.ClusterClass{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster-class",
					Namespace: "kube-system",
				},
				Spec: clusterv1.ClusterClassSpec{},
				Status: clusterv1.ClusterClassStatus{
					Conditions: clusterv1.Conditions{
						{
							Type:   "VariablesReconciled",
							Status: corev1.ConditionTrue,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		var c = fake.NewClientBuilder().
			WithScheme(scheme).
			WithRuntimeObjects(&tt.cc).
			WithLists(&clusterClassesList).
			Build()
		ccs, err := FetchClusterClass(context.Background(), c)
		assert.NoError(t, err)
		assert.Len(t, ccs, 1)
	}
}
