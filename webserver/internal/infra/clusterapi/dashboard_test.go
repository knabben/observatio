package clusterapi

import (
	"context"
	"testing"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/stretchr/testify/assert"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"

	clusterctlv1 "sigs.k8s.io/cluster-api/cmd/clusterctl/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var (
	scheme = runtime.NewScheme()
	_      = clusterctlv1.AddToScheme(scheme) // Register Cluster API types
	_      = clusterv1.AddToScheme(scheme)    // Register Cluster API types
	_      = corev1.AddToScheme(scheme)
)

func Test_GenerateCompVersions(t *testing.T) {
	var providers clusterctlv1.ProviderList
	tests := []struct {
		provider clusterctlv1.Provider
	}{
		{
			provider: clusterctlv1.Provider{
				ProviderName: "bootstrap-kubeadm",
				Type:         "BootstrapProvider",
				Version:      "v1.9.4",
			},
		},
	}

	for _, tt := range tests {
		var c = fake.NewClientBuilder().
			WithScheme(scheme).
			WithRuntimeObjects(&tt.provider).
			WithLists(&providers).
			Build()

		components, err := GenerateComponentVersions(context.Background(), c)
		assert.NoError(t, err)
		assert.Len(t, components, 1)

		for _, component := range components {
			assert.Equal(t, tt.provider.Name, component.Name)
			assert.Equal(t, tt.provider.Type, component.Kind)
			assert.Equal(t, tt.provider.Version, component.Version)
		}
	}
}

func Test_FindServices(t *testing.T) {
	var services corev1.ServiceList
	tests := []struct {
		service corev1.Service
		config  *rest.Config
	}{
		{
			service: corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "service",
					Namespace: "kube-system",
					Labels: map[string]string{
						"kubernetes.io/cluster-service": "true",
					},
				},
				Spec: corev1.ServiceSpec{
					Type: corev1.ServiceTypeClusterIP,
				},
			},
			config: &rest.Config{},
		},
	}
	for _, tt := range tests {
		var c = fake.NewClientBuilder().
			WithScheme(scheme).
			WithRuntimeObjects(&tt.service).
			WithLists(&services).
			Build()
		ctx := context.WithValue(context.Background(), "config", tt.config) // nolint
		services, err := FindServices(ctx, c, tt.service.Namespace)
		assert.NoError(t, err)
		assert.Len(t, services, 2)
		assert.Equal(t, services[0].Name, "control-plane")
		for _, service := range services[1:] {
			assert.Equal(t, tt.service.Name, service.Name)
			assert.NotEmpty(t, service.Path)
		}
	}
}

func Test_ClusterSummary(t *testing.T) {
	var clusters clusterv1.ClusterList
	tests := []struct {
		cluster clusterv1.Cluster
	}{
		{
			cluster: clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster",
					Namespace: "kube-system",
				},
				Spec: clusterv1.ClusterSpec{},
				Status: clusterv1.ClusterStatus{
					Phase: string(clusterv1.ClusterPhase(clusterv1.ClusterPhaseProvisioned)),
				},
			},
		},
	}
	for _, tt := range tests {
		var c = fake.NewClientBuilder().
			WithScheme(scheme).
			WithRuntimeObjects(&tt.cluster).
			WithLists(&clusters).
			Build()
		summary, err := GenerateClusterSummary(context.Background(), c)
		assert.NoError(t, err)
		assert.Equal(t, 0, summary.ClusterFailed)
		assert.Equal(t, 1, summary.ClusterProvisioned)
	}
}
