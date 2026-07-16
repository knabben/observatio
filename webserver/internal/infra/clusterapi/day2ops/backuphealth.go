package day2ops

import (
	"sort"
	"time"
)

// DefaultRPOThreshold is the default recovery point objective used to classify a cluster's most
// recent backup as on-time or stale when no override is configured (spec.md Assumptions).
const DefaultRPOThreshold = 24 * time.Hour

// clusterNameLabel is the standard CAPI label checked against a Backup's label selector to
// determine whether it was explicitly configured to cover a specific Cluster (research.md R5).
const clusterNameLabel = "cluster.x-k8s.io/cluster-name"

// restoreInProgressPhases are Velero Restore phases that represent an active, not-yet-concluded
// restore.
var restoreInProgressPhases = map[string]bool{
	"New": true, "InProgress": true, "WaitingForPluginOperations": true,
}

// backupCoversNamespace reports whether a Backup's IncludedNamespaces covers the given namespace,
// honoring Velero's "empty or [*] means all namespaces" convention (research.md R5).
func backupCoversNamespace(b BackupInfo, namespace string) bool {
	if len(b.IncludedNamespaces) == 0 {
		return true
	}
	for _, ns := range b.IncludedNamespaces {
		if ns == "*" || ns == namespace {
			return true
		}
	}
	return false
}

// backupCoversClusterLabel reports whether a Backup's label selector was explicitly configured to
// target this cluster (research.md R5).
func backupCoversClusterLabel(b BackupInfo, clusterName string) bool {
	if b.LabelSelector == nil {
		return false
	}
	return b.LabelSelector.MatchLabels[clusterNameLabel] == clusterName
}

// backupCoversCluster reports whether a Backup covers a Cluster via either signal (research.md
// R5) — deliberately permissive (OR, not AND): a false negative here is worse than a false
// positive, since this feature exists to prevent operators from wrongly believing a cluster is
// unrecoverable.
func backupCoversCluster(b BackupInfo, clusterNamespace, clusterName string) bool {
	return backupCoversNamespace(b, clusterNamespace) || backupCoversClusterLabel(b, clusterName)
}

// mostRecentBackup returns the covering backup with the latest CompletionTimestamp, or nil if
// none have completed yet. Backups without a CompletionTimestamp (still running) are ignored for
// this purpose — they aren't a verified recovery point of any kind yet.
func mostRecentBackup(covering []BackupInfo) *BackupInfo {
	var candidates []BackupInfo
	for _, b := range covering {
		if b.CompletionTimestamp != nil {
			candidates = append(candidates, b)
		}
	}
	if len(candidates) == 0 {
		return nil
	}
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].CompletionTimestamp.After(*candidates[j].CompletionTimestamp)
	})
	return &candidates[0]
}

// mostRecentTerminalRestore returns the most recently completed (Completed/Failed/etc.) covering
// restore, or nil if none have concluded.
func mostRecentTerminalRestore(covering []RestoreInfo) *RestoreInfo {
	var candidates []RestoreInfo
	for _, r := range covering {
		if r.CompletionTimestamp != nil {
			candidates = append(candidates, r)
		}
	}
	if len(candidates) == 0 {
		return nil
	}
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].CompletionTimestamp.After(*candidates[j].CompletionTimestamp)
	})
	return &candidates[0]
}

// ComputeClusterBackupCoverage computes one Cluster's backup/restore recoverability from the
// currently known Backups and Restores (008/US1-US3, research.md R5). Always returns a value —
// never omits a cluster, even with zero coverage (spec.md Edge Cases).
func ComputeClusterBackupCoverage(clusterRef ObjectRef, backups []BackupInfo, restores []RestoreInfo, rpo time.Duration, now time.Time) ClusterBackupCoverage {
	coverage := ClusterBackupCoverage{ClusterRef: clusterRef}

	var covering []BackupInfo
	for _, b := range backups {
		if backupCoversCluster(b, clusterRef.Namespace, clusterRef.Name) {
			covering = append(covering, b)
		}
	}

	if recent := mostRecentBackup(covering); recent != nil {
		age := now.Sub(*recent.CompletionTimestamp)
		coverage.MostRecentBackupName = recent.Name
		coverage.MostRecentBackupAge = age.Round(time.Second).String()
		// Only a fully Completed backup counts as a verified recovery point; a PartiallyFailed or
		// otherwise non-Completed backup remains visible above (not hidden) but doesn't mark the
		// cluster as Covered (spec.md Edge Cases).
		if recent.Phase == "Completed" {
			coverage.Covered = true
			coverage.Stale = age > rpo
		}
	}

	var coveringRestores []RestoreInfo
	coveredBackupNames := make(map[string]bool, len(covering))
	for _, b := range covering {
		coveredBackupNames[b.Name] = true
	}
	for _, r := range restores {
		if r.BackupName != "" && coveredBackupNames[r.BackupName] {
			coveringRestores = append(coveringRestores, r)
		}
	}
	for _, r := range coveringRestores {
		if restoreInProgressPhases[r.Phase] {
			coverage.RestoreInProgress = true
		}
	}
	if recent := mostRecentTerminalRestore(coveringRestores); recent != nil {
		if recent.Phase == "Completed" {
			coverage.LastRestoreOutcome = "succeeded"
		} else {
			coverage.LastRestoreOutcome = "failed"
		}
	}

	return coverage
}

// ComputeBackupHealth computes the full Backup Health payload for the Day-2 Ops landing page
// (008/US1, US3). veleroInstalled reflects research.md R8's CRD-existence check — when false,
// every other field is zero-valued and Available is false, driving the "not available" state
// (FR-011) rather than a false-empty "all clear."
func ComputeBackupHealth(
	veleroInstalled bool,
	clusterRefs []ObjectRef,
	storageLocations []BackupStorageLocationInfo,
	backups []BackupInfo,
	restores []RestoreInfo,
	rpo time.Duration,
	now time.Time,
) BackupHealth {
	if !veleroInstalled {
		return BackupHealth{Available: false, RPOThresholdSeconds: int64(rpo.Seconds())}
	}

	locations := make([]BackupStorageLocationStatus, 0, len(storageLocations))
	for _, loc := range storageLocations {
		locations = append(locations, BackupStorageLocationStatus{
			Name: loc.Name, Namespace: loc.Namespace, Default: loc.Default,
			Reachable: loc.Phase == "Available",
		})
	}

	coverage := make([]ClusterBackupCoverage, 0, len(clusterRefs))
	restoresInProgress := 0
	for _, ref := range clusterRefs {
		c := ComputeClusterBackupCoverage(ref, backups, restores, rpo, now)
		if c.RestoreInProgress {
			restoresInProgress++
		}
		coverage = append(coverage, c)
	}

	return BackupHealth{
		Available:           true,
		StorageLocations:    locations,
		ClusterCoverage:     coverage,
		RPOThresholdSeconds: int64(rpo.Seconds()),
		RestoresInProgress:  restoresInProgress,
	}
}
