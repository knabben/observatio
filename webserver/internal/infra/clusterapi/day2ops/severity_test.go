package day2ops

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ComputeMachineHealthCheckSeverity_SelfHealing(t *testing.T) {
	mhc := MachineHealthCheckStatus{Name: "mhc", ExpectedMachines: 3, CurrentHealthy: 2, RemediationsAllowed: 1}

	severity := ComputeMachineHealthCheckSeverity(ObjectRef{Name: "mhc"}, mhc)

	require.NotNil(t, severity)
	assert.Equal(t, SeveritySelfHealing, severity.Level)
}

func Test_ComputeMachineHealthCheckSeverity_MaxUnhealthyBreached(t *testing.T) {
	mhc := MachineHealthCheckStatus{Name: "mhc", ExpectedMachines: 3, CurrentHealthy: 1, RemediationsAllowed: 0}

	severity := ComputeMachineHealthCheckSeverity(ObjectRef{Name: "mhc"}, mhc)

	require.NotNil(t, severity)
	assert.Equal(t, SeverityNeedsInvestigation, severity.Level)
	assert.Contains(t, severity.Reason, "network partition")
}

func Test_ComputeMachineHealthCheckSeverity_AllHealthy(t *testing.T) {
	mhc := MachineHealthCheckStatus{Name: "mhc", ExpectedMachines: 3, CurrentHealthy: 3, RemediationsAllowed: 1}
	assert.Nil(t, ComputeMachineHealthCheckSeverity(ObjectRef{Name: "mhc"}, mhc))
}

func Test_ComputeProviderControllerSeverity_CrashLooping(t *testing.T) {
	pod := ControllerPodStatus{Namespace: "capd-system", PodName: "capd-controller-manager-abc", Ready: false, WaitingReason: "CrashLoopBackOff"}

	severity := ComputeProviderControllerSeverity(pod)

	require.NotNil(t, severity)
	assert.Equal(t, SeverityProviderDegraded, severity.Level)
	assert.Contains(t, severity.Reason, "CrashLoopBackOff")
}

func Test_ComputeProviderControllerSeverity_Ready(t *testing.T) {
	pod := ControllerPodStatus{Namespace: "capd-system", PodName: "capd-controller-manager-abc", Ready: true}
	assert.Nil(t, ComputeProviderControllerSeverity(pod))
}

func Test_ComputeManagementClusterSeverity(t *testing.T) {
	require.NotNil(t, ComputeManagementClusterSeverity(true))
	assert.Nil(t, ComputeManagementClusterSeverity(false))
}

func Test_ComputeCASecretMissingSeverity(t *testing.T) {
	severity := ComputeCASecretMissingSeverity(ObjectRef{Name: "prod-1"}, false)
	require.NotNil(t, severity)
	assert.Equal(t, SeverityManagementCritical, severity.Level)

	assert.Nil(t, ComputeCASecretMissingSeverity(ObjectRef{Name: "prod-1"}, true))
}
