package day2ops

// Category identifies which object category a HealthRollup summarizes.
type Category string

const (
	CategoryCluster           Category = "cluster"
	CategoryMachineDeployment Category = "machine_deployment"
	CategoryMachine           Category = "machine"
)

// HealthRollup is a per-category healthy/degraded/failed summary shown on the Day-2 Ops landing
// screen (FR-002).
type HealthRollup struct {
	Category    Category `json:"category"`
	Healthy     int      `json:"healthy"`
	Degraded    int      `json:"degraded"`
	Failed      int      `json:"failed"`
	Unavailable bool     `json:"unavailable"`
}

// ObjectRef identifies a specific Kubernetes object by group/version/resource and namespace/name.
type ObjectRef struct {
	Group     string `json:"group"`
	Version   string `json:"version"`
	Resource  string `json:"resource"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

// DebugLayerName identifies one of the four ordered debugging stages (FR-005).
type DebugLayerName string

const (
	LayerConditions         DebugLayerName = "conditions"
	LayerPhase              DebugLayerName = "phase"
	LayerProviderResource   DebugLayerName = "provider_resource"
	LayerControllerActivity DebugLayerName = "controller_activity"
)

// DebugLayerStatus reports whether a given layer explains an object's failure.
type DebugLayerStatus string

const (
	LayerStatusOK           DebugLayerStatus = "ok"
	LayerStatusImplicated   DebugLayerStatus = "implicated"
	LayerStatusInconclusive DebugLayerStatus = "inconclusive"
)

// DebugLayer is one stage of evidence in an object's debugging path (data-model.md DebugLayer).
type DebugLayer struct {
	Layer    DebugLayerName   `json:"layer"`
	Status   DebugLayerStatus `json:"status"`
	Evidence []string         `json:"evidence"`
	Source   string           `json:"source"`
}

// DebugPath is the full ordered layer breakdown for one unhealthy object (FR-004-FR-007).
type DebugPath struct {
	ObjectRef ObjectRef    `json:"objectRef"`
	Layers    []DebugLayer `json:"layers"`
	Summary   string       `json:"summary"`
}

// RiskKind identifies which proactive risk class a RiskWarning reports (US3).
type RiskKind string

const (
	RiskCertExpiry     RiskKind = "cert_expiry"
	RiskStalledRollout RiskKind = "stalled_rollout"
	RiskVersionSkew    RiskKind = "version_skew"
	RiskDrift          RiskKind = "drift"
)

// RiskCheckStatus distinguishes "checked, no issue found" (implicit, no RiskWarning emitted) from
// "the check could not be evaluated" (FR-018).
type RiskCheckStatus string

const (
	RiskCheckEvaluated    RiskCheckStatus = "evaluated"
	RiskCheckNotEvaluable RiskCheckStatus = "not_evaluable"
)

// RiskWarning is a proactively detected issue attached to a specific object (FR-008-FR-011).
type RiskWarning struct {
	ObjectRef   ObjectRef       `json:"objectRef"`
	Kind        RiskKind        `json:"kind"`
	Detail      string          `json:"detail"`
	LikelyCause string          `json:"likelyCause"`
	CheckStatus RiskCheckStatus `json:"checkStatus"`
}

// SeverityLevel is a strictly-increasing urgency classification (FR-012-FR-016).
type SeverityLevel string

const (
	SeveritySelfHealing        SeverityLevel = "self_healing"
	SeverityNeedsInvestigation SeverityLevel = "needs_investigation"
	SeverityProviderDegraded   SeverityLevel = "provider_degraded"
	SeverityManagementCritical SeverityLevel = "management_critical"
)

// FailureSeverity classifies a detected issue's urgency. ObjectRef is nil for cluster-wide
// severities (e.g. management-critical).
type FailureSeverity struct {
	ObjectRef *ObjectRef    `json:"objectRef"`
	Level     SeverityLevel `json:"level"`
	Reason    string        `json:"reason"`
}

// Data is the payload of a Day2OpsEvent broadcast to a connected dashboard (contracts/day2ops-ws-event.md).
type Data struct {
	Rollups           []HealthRollup    `json:"rollups"`
	DebugPaths        []DebugPath       `json:"debugPaths"`
	Risks             []RiskWarning     `json:"risks"`
	Severities        []FailureSeverity `json:"severities"`
	SourceUnavailable bool              `json:"sourceUnavailable"`
}
