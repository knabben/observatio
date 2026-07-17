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

// RecoveryInfo augments a CA-secret-missing FailureSeverity with backup-based recoverability
// (008/US2). Nil for every severity kind other than CA-secret-missing, and nil (not a false
// Recoverable) when Velero isn't installed — recoverability is genuinely unknown in that case,
// which is a different fact than "known to have no covering backup."
type RecoveryInfo struct {
	Recoverable       bool   `json:"recoverable"`
	CoveringBackupAge string `json:"coveringBackupAge,omitempty"`
}

// FailureSeverity classifies a detected issue's urgency. ObjectRef is nil for cluster-wide
// severities (e.g. management-critical).
type FailureSeverity struct {
	ObjectRef    *ObjectRef    `json:"objectRef"`
	Level        SeverityLevel `json:"level"`
	Reason       string        `json:"reason"`
	RecoveryInfo *RecoveryInfo `json:"recoveryInfo,omitempty"`
}

// BackupStorageLocationStatus is a normalized summary of one Velero BackupStorageLocation's
// reachability (008/US1).
type BackupStorageLocationStatus struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Reachable bool   `json:"reachable"`
	Default   bool   `json:"default"`
}

// ClusterBackupCoverage is the per-cluster recoverability computed from known Backups/Restores
// (008/US1-US3, research.md R5). One entry per known Cluster — always present, even with no
// covering backup, so a cluster is never silently omitted (spec.md Edge Cases).
type ClusterBackupCoverage struct {
	ClusterRef           ObjectRef `json:"clusterRef"`
	Covered              bool      `json:"covered"`
	MostRecentBackupAge  string    `json:"mostRecentBackupAge,omitempty"`
	MostRecentBackupName string    `json:"mostRecentBackupName,omitempty"`
	Stale                bool      `json:"stale"`
	RestoreInProgress    bool      `json:"restoreInProgress"`
	LastRestoreOutcome   string    `json:"lastRestoreOutcome"` // "", "succeeded", or "failed"
}

// BackupHealth is the full Backup Health payload for the Day-2 Ops landing page (008/US1,
// FR-001-FR-005, FR-009, FR-011).
type BackupHealth struct {
	Available           bool                          `json:"available"`
	StorageLocations    []BackupStorageLocationStatus `json:"storageLocations"`
	ClusterCoverage     []ClusterBackupCoverage       `json:"clusterCoverage"`
	RPOThresholdSeconds int64                         `json:"rpoThresholdSeconds"`
	RestoresInProgress  int                           `json:"restoresInProgress"`
}

// Data is the payload of a Day2OpsEvent broadcast to a connected dashboard (contracts/day2ops-ws-event.md).
type Data struct {
	Rollups           []HealthRollup    `json:"rollups"`
	DebugPaths        []DebugPath       `json:"debugPaths"`
	Risks             []RiskWarning     `json:"risks"`
	Severities        []FailureSeverity `json:"severities"`
	SourceUnavailable bool              `json:"sourceUnavailable"`
	BackupHealth      BackupHealth      `json:"backupHealth"`
}
