package day2ops

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ComputeDriftRisk_ObservedGenerationBehind(t *testing.T) {
	provider := ProviderResourceStatus{Kind: "DockerMachine", Name: "worker-0", Generation: 3, ObservedGeneration: 2}

	risk := ComputeDriftRisk(ObjectRef{Name: "worker-0"}, provider)

	require.NotNil(t, risk)
	assert.Equal(t, RiskDrift, risk.Kind)
	assert.Contains(t, risk.Detail, "DockerMachine/worker-0")
}

func Test_ComputeDriftRisk_UpToDate(t *testing.T) {
	provider := ProviderResourceStatus{Kind: "DockerMachine", Name: "worker-0", Generation: 3, ObservedGeneration: 3}
	assert.Nil(t, ComputeDriftRisk(ObjectRef{Name: "worker-0"}, provider))
}

func Test_ComputeDriftRisk_ZeroObservedGenerationIgnored(t *testing.T) {
	// A provider that doesn't populate observedGeneration at all (0) must not be flagged as drift
	// just because it's less than a nonzero generation - that would be a false positive on every
	// object from a provider that simply doesn't report this field.
	provider := ProviderResourceStatus{Kind: "DockerMachine", Name: "worker-0", Generation: 3, ObservedGeneration: 0}
	assert.Nil(t, ComputeDriftRisk(ObjectRef{Name: "worker-0"}, provider))
}
