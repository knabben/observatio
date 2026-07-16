package day2ops

import "fmt"

// MachineHealthCheckStatus is a normalized summary of one MachineHealthCheck's remediation state.
type MachineHealthCheckStatus struct {
	Name                string
	ExpectedMachines    int32
	CurrentHealthy      int32
	RemediationsAllowed int32
}

// ComputeMachineHealthCheckSeverity classifies a MachineHealthCheck's current remediation state
// (FR-012, FR-013): self-healing while unhealthy Machines are being remediated, or
// needs-investigation once remediation has paused because the maxUnhealthy threshold was breached
// — which may indicate a network partition rather than independent node failures. Returns nil
// when every target Machine is currently healthy (nothing to report).
func ComputeMachineHealthCheckSeverity(objectRef ObjectRef, mhc MachineHealthCheckStatus) *FailureSeverity {
	unhealthy := mhc.ExpectedMachines - mhc.CurrentHealthy
	if unhealthy <= 0 {
		return nil
	}
	if mhc.RemediationsAllowed > 0 {
		return &FailureSeverity{
			ObjectRef: &objectRef,
			Level:     SeveritySelfHealing,
			Reason:    fmt.Sprintf("MachineHealthCheck %s is remediating %d unhealthy machine(s)", mhc.Name, unhealthy),
		}
	}
	return &FailureSeverity{
		ObjectRef: &objectRef,
		Level:     SeverityNeedsInvestigation,
		Reason: fmt.Sprintf(
			"MachineHealthCheck %s has %d unhealthy machine(s) but remediation is paused (maxUnhealthy threshold breached) — may indicate a network partition rather than independent node failures",
			mhc.Name, unhealthy,
		),
	}
}

// ControllerPodStatus is a normalized summary of one controller Pod's readiness, used to detect
// crash-looping/degraded controllers (FR-014).
type ControllerPodStatus struct {
	Namespace     string
	PodName       string
	Ready         bool
	WaitingReason string // e.g. "CrashLoopBackOff", empty if not currently waiting
}

// ComputeProviderControllerSeverity flags a not-ready controller Pod as a provider-level (not
// merely object-level) failure (FR-014). Returns nil for a ready Pod.
func ComputeProviderControllerSeverity(pod ControllerPodStatus) *FailureSeverity {
	if pod.Ready {
		return nil
	}
	reason := fmt.Sprintf("Controller Pod %s/%s is not ready", pod.Namespace, pod.PodName)
	if pod.WaitingReason != "" {
		reason = fmt.Sprintf("Controller Pod %s/%s is %s", pod.Namespace, pod.PodName, pod.WaitingReason)
	}
	return &FailureSeverity{Level: SeverityProviderDegraded, Reason: reason}
}

// ComputeManagementClusterSeverity flags the management cluster itself as degraded (FR-015) when
// the aggregator's own connection to its API server has been lost — approximated from the
// existing SourceUnavailable signal rather than a separate etcd-quorum probe, which isn't
// accessible read-only through the Kubernetes API (research.md R8). This is a cluster-wide
// severity: ObjectRef is nil.
func ComputeManagementClusterSeverity(sourceUnavailable bool) *FailureSeverity {
	if !sourceUnavailable {
		return nil
	}
	return &FailureSeverity{
		Level:  SeverityManagementCritical,
		Reason: "Management cluster API server is unreachable — all lifecycle operations (scaling, upgrades, certificate rotation) are blocked",
	}
}

// ComputeCASecretMissingSeverity flags a cluster's missing CA secret as the highest-severity
// warning available for that cluster (FR-016): certificate issuance/rotation is blocked, and the
// original CA cannot be substituted. caSecretFound is false when the cert-expiry fetch succeeded
// (no permissions error) but found no "<cluster>-ca" Secret among the results.
//
// coverage is this cluster's backup recoverability (008/US2, computed separately by
// ComputeClusterBackupCoverage) — nil when that data isn't available at all (e.g. Velero isn't
// installed), in which case RecoveryInfo is left nil rather than falsely reporting
// Recoverable: false; recoverability is genuinely unknown in that case, not known-absent.
func ComputeCASecretMissingSeverity(objectRef ObjectRef, caSecretFound bool, coverage *ClusterBackupCoverage) *FailureSeverity {
	if caSecretFound {
		return nil
	}
	reason := fmt.Sprintf(
		"Cluster %s's CA secret is missing or inaccessible — certificate issuance/rotation is blocked for its nodes; the original CA cannot be substituted.",
		objectRef.Name,
	)
	var recoveryInfo *RecoveryInfo
	if coverage != nil {
		if coverage.Covered {
			reason += fmt.Sprintf(" A backup completed %s ago covers this cluster — recovery is straightforward.", coverage.MostRecentBackupAge)
			recoveryInfo = &RecoveryInfo{Recoverable: true, CoveringBackupAge: coverage.MostRecentBackupAge}
		} else {
			reason += " No covering backup exists — this cluster's data is unrecoverable."
			recoveryInfo = &RecoveryInfo{Recoverable: false}
		}
	}
	return &FailureSeverity{
		ObjectRef:    &objectRef,
		Level:        SeverityManagementCritical,
		Reason:       reason,
		RecoveryInfo: recoveryInfo,
	}
}
