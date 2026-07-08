package day2ops

import "fmt"

// ComputeDriftRisk flags a provider-infra object whose controller has not finished reconciling
// the current spec (metadata.generation ahead of status.observedGeneration) — a provider-agnostic
// proxy for drift and other stuck-reconciliation cases that requires no cloud-provider-side
// introspection (research.md R3). Best-effort: a truly instantaneous mismatch during normal
// reconciliation would also match this check; there is no grace-period tracking in this pass.
func ComputeDriftRisk(objectRef ObjectRef, provider ProviderResourceStatus) *RiskWarning {
	if provider.ObservedGeneration <= 0 || provider.ObservedGeneration >= provider.Generation {
		return nil
	}
	return &RiskWarning{
		ObjectRef: objectRef,
		Kind:      RiskDrift,
		Detail: fmt.Sprintf(
			"%s/%s: observedGeneration (%d) is behind generation (%d) — the controller has not finished reconciling the current spec",
			provider.Kind, provider.Name, provider.ObservedGeneration, provider.Generation,
		),
		CheckStatus: RiskCheckEvaluated,
	}
}
