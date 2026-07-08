package day2ops

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ComputeVersionSkewRisk_StoredVersionNoLongerServed(t *testing.T) {
	info := CRDVersionInfo{
		Name:           "machines.cluster.x-k8s.io",
		ServedVersions: []string{"v1beta1"},
		StoredVersions: []string{"v1alpha4", "v1beta1"},
	}

	risk := ComputeVersionSkewRisk(ObjectRef{Name: "machines.cluster.x-k8s.io"}, info)

	require.NotNil(t, risk)
	assert.Equal(t, RiskVersionSkew, risk.Kind)
	assert.Contains(t, risk.Detail, "v1alpha4")
}

func Test_ComputeVersionSkewRisk_AllStoredVersionsServed(t *testing.T) {
	info := CRDVersionInfo{
		Name:           "machines.cluster.x-k8s.io",
		ServedVersions: []string{"v1beta1"},
		StoredVersions: []string{"v1beta1"},
	}
	assert.Nil(t, ComputeVersionSkewRisk(ObjectRef{Name: "machines.cluster.x-k8s.io"}, info))
}
