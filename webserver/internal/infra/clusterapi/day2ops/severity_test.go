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
	severity := ComputeCASecretMissingSeverity(ObjectRef{Name: "prod-1"}, false, nil)
	require.NotNil(t, severity)
	assert.Equal(t, SeverityManagementCritical, severity.Level)

	assert.Nil(t, ComputeCASecretMissingSeverity(ObjectRef{Name: "prod-1"}, true, nil))
}

func Test_ComputeCASecretMissingSeverity_RecoverableWithCoveringBackup(t *testing.T) {
	coverage := &ClusterBackupCoverage{Covered: true, MostRecentBackupAge: "3h0m0s"}

	severity := ComputeCASecretMissingSeverity(ObjectRef{Name: "prod-1"}, false, coverage)

	require.NotNil(t, severity)
	require.NotNil(t, severity.RecoveryInfo)
	assert.True(t, severity.RecoveryInfo.Recoverable)
	assert.Equal(t, "3h0m0s", severity.RecoveryInfo.CoveringBackupAge)
	assert.Contains(t, severity.Reason, "3h0m0s")
	assert.Contains(t, severity.Reason, "recovery is straightforward")
}

func Test_ComputeCASecretMissingSeverity_NoCoveringBackup(t *testing.T) {
	coverage := &ClusterBackupCoverage{Covered: false}

	severity := ComputeCASecretMissingSeverity(ObjectRef{Name: "prod-1"}, false, coverage)

	require.NotNil(t, severity)
	require.NotNil(t, severity.RecoveryInfo)
	assert.False(t, severity.RecoveryInfo.Recoverable)
	assert.Contains(t, severity.Reason, "unrecoverable")
}

func Test_ComputeCASecretMissingSeverity_NoCoverageData_OmitsRecoveryInfo(t *testing.T) {
	// Velero not installed (or no coverage data available) — recoverability is genuinely
	// unknown, which must not be conflated with "known to have no covering backup."
	severity := ComputeCASecretMissingSeverity(ObjectRef{Name: "prod-1"}, false, nil)

	require.NotNil(t, severity)
	assert.Nil(t, severity.RecoveryInfo)
}
