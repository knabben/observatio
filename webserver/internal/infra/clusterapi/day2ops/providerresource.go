package day2ops

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

// ExtractProviderResourceStatus reads `.status.ready` and the first non-ready condition's
// reason/message from an infra-machine object (DockerMachine, VSphereMachine, ...) generically via
// unstructured accessors — no per-provider Go type is needed, keeping provider specifics opaque to
// the core domain (Constitution Principle III; mirrors the existing convention already used by
// fetchers.ProcessDockerMachine).
func ExtractProviderResourceStatus(u *unstructured.Unstructured) ProviderResourceStatus {
	ready, _, _ := unstructured.NestedBool(u.Object, "status", "ready")
	observedGeneration, _, _ := unstructured.NestedInt64(u.Object, "status", "observedGeneration")
	status := ProviderResourceStatus{
		Kind: u.GetKind(), Name: u.GetName(), Ready: ready,
		Generation:         u.GetGeneration(),
		ObservedGeneration: observedGeneration,
	}

	conditions, _, _ := unstructured.NestedSlice(u.Object, "status", "conditions")
	for _, raw := range conditions {
		cond, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		condStatus, _, _ := unstructured.NestedString(cond, "status")
		if condStatus == "True" {
			continue
		}
		message, _, _ := unstructured.NestedString(cond, "message")
		reason, _, _ := unstructured.NestedString(cond, "reason")
		status.Message = firstNonEmpty(message, reason)
		break
	}
	return status
}
