package day2ops

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// BackupInfo is a lightweight, decoded projection of a Velero Backup — only the fields this
// feature needs, extracted generically via unstructured accessors so no typed Velero Go
// dependency is required (008, research.md R1; mirrors the existing ExtractProviderResourceStatus
// convention).
type BackupInfo struct {
	Name                string
	Namespace           string
	Phase               string
	IncludedNamespaces  []string
	LabelSelector       *metav1.LabelSelector
	StorageLocation     string
	CompletionTimestamp *time.Time
}

// RestoreInfo is a lightweight, decoded projection of a Velero Restore.
type RestoreInfo struct {
	Name                string
	Namespace           string
	Phase               string
	BackupName          string
	CompletionTimestamp *time.Time
}

// BackupStorageLocationInfo is a lightweight, decoded projection of a Velero
// BackupStorageLocation.
type BackupStorageLocationInfo struct {
	Name      string
	Namespace string
	Phase     string
	Default   bool
}

// ScheduleInfo is a lightweight, decoded projection of a Velero Schedule. Watched and stored
// alongside the other three Velero kinds (spec.md's Key Entities), though no functional
// requirement in this iteration computes from it yet — kept minimal rather than building unused
// computation ahead of a concrete need.
type ScheduleInfo struct {
	Name      string
	Namespace string
}

// ExtractScheduleInfo decodes a Schedule's identity via unstructured accessors.
func ExtractScheduleInfo(u *unstructured.Unstructured) ScheduleInfo {
	return ScheduleInfo{Name: u.GetName(), Namespace: u.GetNamespace()}
}

// ExtractBackupInfo decodes a Backup's fields via unstructured accessors (008/US1).
func ExtractBackupInfo(u *unstructured.Unstructured) BackupInfo {
	phase, _, _ := unstructured.NestedString(u.Object, "status", "phase")
	includedNamespaces, _, _ := unstructured.NestedStringSlice(u.Object, "spec", "includedNamespaces")
	storageLocation, _, _ := unstructured.NestedString(u.Object, "spec", "storageLocation")

	info := BackupInfo{
		Name:               u.GetName(),
		Namespace:          u.GetNamespace(),
		Phase:              phase,
		IncludedNamespaces: includedNamespaces,
		StorageLocation:    storageLocation,
	}

	if selMap, found, _ := unstructured.NestedMap(u.Object, "spec", "labelSelector"); found {
		var sel metav1.LabelSelector
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(selMap, &sel); err == nil {
			info.LabelSelector = &sel
		}
	}

	if ts := parseNestedTimestamp(u.Object, "status", "completionTimestamp"); ts != nil {
		info.CompletionTimestamp = ts
	}

	return info
}

// ExtractRestoreInfo decodes a Restore's fields via unstructured accessors (008/US3).
func ExtractRestoreInfo(u *unstructured.Unstructured) RestoreInfo {
	phase, _, _ := unstructured.NestedString(u.Object, "status", "phase")
	backupName, _, _ := unstructured.NestedString(u.Object, "spec", "backupName")

	info := RestoreInfo{
		Name:       u.GetName(),
		Namespace:  u.GetNamespace(),
		Phase:      phase,
		BackupName: backupName,
	}

	if ts := parseNestedTimestamp(u.Object, "status", "completionTimestamp"); ts != nil {
		info.CompletionTimestamp = ts
	}

	return info
}

// ExtractBackupStorageLocationInfo decodes a BackupStorageLocation's fields via unstructured
// accessors (008/US1). Reachable is derived from Phase == "Available" by the caller
// (ComputeBackupHealth), not baked in here, so the raw phase remains available for diagnostics.
func ExtractBackupStorageLocationInfo(u *unstructured.Unstructured) BackupStorageLocationInfo {
	phase, _, _ := unstructured.NestedString(u.Object, "status", "phase")
	isDefault, _, _ := unstructured.NestedBool(u.Object, "spec", "default")
	return BackupStorageLocationInfo{
		Name: u.GetName(), Namespace: u.GetNamespace(), Phase: phase, Default: isDefault,
	}
}

// parseNestedTimestamp reads an RFC3339 timestamp string at the given path, returning nil if
// absent, empty, or unparseable rather than erroring — a malformed timestamp on one object must
// not break the whole Backup Health computation.
func parseNestedTimestamp(obj map[string]interface{}, fields ...string) *time.Time {
	s, found, _ := unstructured.NestedString(obj, fields...)
	if !found || s == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil
	}
	return &t
}
