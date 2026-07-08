package day2ops

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// ProviderResourceStatus is a normalized, provider-agnostic summary of an infra-machine object's
// readiness (e.g. DockerMachine, VSphereMachine), extracted generically so provider specifics stay
// opaque to the core domain (Constitution Principle III).
type ProviderResourceStatus struct {
	Kind    string
	Name    string
	Ready   bool
	Message string
	// Generation/ObservedGeneration back the drift heuristic (research.md R3): a controller that
	// hasn't finished reconciling the current spec leaves ObservedGeneration behind Generation.
	Generation         int64
	ObservedGeneration int64
}

var layerFriendlyName = map[DebugLayerName]string{
	LayerConditions:         "object conditions",
	LayerPhase:              "machine phase",
	LayerProviderResource:   "infrastructure provisioning",
	LayerControllerActivity: "controller reconciliation",
}

// ComputeMachineDebugPath synthesizes the ordered, labeled layer breakdown for one Machine
// (FR-004-FR-007): object conditions -> Machine phase -> provider-infra resource -> controller
// reconciliation activity. provider is nil when no matching provider-infra object is known.
// events are recent controller-activity evidence (research.md R2); they are only surfaced when
// the three higher layers are all inconclusive (FR-007).
func ComputeMachineDebugPath(objectRef ObjectRef, m clusterv1.Machine, provider *ProviderResourceStatus, events []string) DebugPath {
	conditionsLayer := debugConditionsLayer(m)
	phaseLayer := debugPhaseLayer(m)
	providerLayer := debugProviderResourceLayer(provider)

	higherLayersInconclusive := conditionsLayer.Status != LayerStatusImplicated &&
		phaseLayer.Status != LayerStatusImplicated &&
		providerLayer.Status != LayerStatusImplicated

	controllerLayer := DebugLayer{Layer: LayerControllerActivity, Status: LayerStatusInconclusive, Evidence: []string{}, Source: ""}
	if higherLayersInconclusive && len(events) > 0 {
		controllerLayer = DebugLayer{
			Layer:    LayerControllerActivity,
			Status:   LayerStatusImplicated,
			Evidence: events,
			Source:   fmt.Sprintf("Machine/%s events", m.Name),
		}
	}

	layers := []DebugLayer{conditionsLayer, phaseLayer, providerLayer, controllerLayer}
	return DebugPath{
		ObjectRef: objectRef,
		Layers:    layers,
		Summary:   summarizeDebugPath(layers),
	}
}

func debugConditionsLayer(m clusterv1.Machine) DebugLayer {
	layer := DebugLayer{Layer: LayerConditions, Status: LayerStatusInconclusive, Evidence: []string{}, Source: fmt.Sprintf("Machine/%s", m.Name)}
	if len(m.Status.Conditions) == 0 {
		return layer
	}
	allTrue := true
	for _, c := range m.Status.Conditions {
		if c.Status != corev1.ConditionTrue {
			allTrue = false
			layer.Evidence = append(layer.Evidence, fmt.Sprintf("%s=%s: %s", c.Type, c.Status, firstNonEmpty(c.Reason, c.Message)))
		}
	}
	if allTrue {
		layer.Status = LayerStatusOK
	} else {
		layer.Status = LayerStatusImplicated
	}
	return layer
}

func debugPhaseLayer(m clusterv1.Machine) DebugLayer {
	layer := DebugLayer{Layer: LayerPhase, Status: LayerStatusInconclusive, Evidence: []string{}, Source: fmt.Sprintf("Machine/%s", m.Name)}
	switch m.Status.Phase {
	case "":
		return layer
	case string(clusterv1.MachinePhaseRunning):
		layer.Status = LayerStatusOK
		return layer
	default:
		layer.Status = LayerStatusImplicated
		layer.Evidence = []string{fmt.Sprintf("Phase=%s", m.Status.Phase)}
		return layer
	}
}

func debugProviderResourceLayer(provider *ProviderResourceStatus) DebugLayer {
	layer := DebugLayer{Layer: LayerProviderResource, Status: LayerStatusInconclusive, Evidence: []string{}}
	if provider == nil {
		return layer
	}
	layer.Source = fmt.Sprintf("%s/%s", provider.Kind, provider.Name)
	if provider.Ready {
		layer.Status = LayerStatusOK
		return layer
	}
	layer.Status = LayerStatusImplicated
	layer.Evidence = []string{firstNonEmpty(provider.Message, "Ready=False")}
	return layer
}

// ShouldFetchControllerActivityEvents reports whether it's worth fetching Events for this
// Machine's controller-activity layer at all — true only when conditions and phase (the two
// cheapest, already-in-memory layers) are both inconclusive, so callers can skip an extra API
// call in the common case where those higher layers already explain the failure (research.md R2,
// FR-007). Shared by the live WS watcher and the on-demand REST detail handler so the "when to
// bother fetching Events" decision lives in one place.
func ShouldFetchControllerActivityEvents(m clusterv1.Machine) bool {
	for _, c := range m.Status.Conditions {
		if c.Status != corev1.ConditionTrue {
			return false
		}
	}
	return m.Status.Phase == "" || m.Status.Phase == string(clusterv1.MachinePhaseRunning)
}

func summarizeDebugPath(layers []DebugLayer) string {
	for _, l := range layers {
		if l.Status == LayerStatusImplicated && len(l.Evidence) > 0 {
			return fmt.Sprintf("Waiting on %s (%s: %s)", layerFriendlyName[l.Layer], firstNonEmpty(l.Source, string(l.Layer)), l.Evidence[0])
		}
	}
	return "No specific cause identified from conditions, phase, provider resource, or recent controller activity"
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
