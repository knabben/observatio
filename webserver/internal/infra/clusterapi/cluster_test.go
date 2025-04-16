package clusterapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"k8s.io/apimachinery/pkg/runtime"
	clusterctlv1 "sigs.k8s.io/cluster-api/cmd/clusterctl/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var (
	scheme = runtime.NewScheme()
	_      = clusterctlv1.AddToScheme(scheme) // Register Cluster API types
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
