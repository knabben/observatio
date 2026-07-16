package day2ops

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var testClusterRef = ObjectRef{Group: "cluster.x-k8s.io", Version: "v1beta1", Resource: "clusters", Namespace: "default", Name: "capi-workload"}

func completedBackup(name, namespace string, includedNamespaces []string, age time.Duration, now time.Time) BackupInfo {
	ts := now.Add(-age)
	return BackupInfo{Name: name, Namespace: namespace, Phase: "Completed", IncludedNamespaces: includedNamespaces, CompletionTimestamp: &ts}
}

func Test_ComputeClusterBackupCoverage_namespaceMatch(t *testing.T) {
	now := time.Now()
	backups := []BackupInfo{completedBackup("b1", "velero", []string{"default"}, 3*time.Hour, now)}

	coverage := ComputeClusterBackupCoverage(testClusterRef, backups, nil, DefaultRPOThreshold, now)

	assert.True(t, coverage.Covered)
	assert.False(t, coverage.Stale)
	assert.Equal(t, "b1", coverage.MostRecentBackupName)
	assert.NotEmpty(t, coverage.MostRecentBackupAge)
}

func Test_ComputeClusterBackupCoverage_allNamespacesBackupCovers(t *testing.T) {
	now := time.Now()
	backups := []BackupInfo{completedBackup("b-all", "velero", nil, time.Hour, now)}

	coverage := ComputeClusterBackupCoverage(testClusterRef, backups, nil, DefaultRPOThreshold, now)
	assert.True(t, coverage.Covered)
}

func Test_ComputeClusterBackupCoverage_labelSelectorMatch(t *testing.T) {
	now := time.Now()
	b := BackupInfo{
		Name: "b-label", Namespace: "velero", Phase: "Completed",
		IncludedNamespaces:  []string{"other-namespace"},
		LabelSelector:       &metav1.LabelSelector{MatchLabels: map[string]string{"cluster.x-k8s.io/cluster-name": "capi-workload"}},
		CompletionTimestamp: timePtr(now.Add(-time.Hour)),
	}

	coverage := ComputeClusterBackupCoverage(testClusterRef, []BackupInfo{b}, nil, DefaultRPOThreshold, now)
	assert.True(t, coverage.Covered)
}

func Test_ComputeClusterBackupCoverage_noMatch(t *testing.T) {
	now := time.Now()
	backups := []BackupInfo{completedBackup("b1", "velero", []string{"other-namespace"}, time.Hour, now)}

	coverage := ComputeClusterBackupCoverage(testClusterRef, backups, nil, DefaultRPOThreshold, now)
	assert.False(t, coverage.Covered)
	assert.Empty(t, coverage.MostRecentBackupName)
}

func Test_ComputeClusterBackupCoverage_noBackupAtAll(t *testing.T) {
	coverage := ComputeClusterBackupCoverage(testClusterRef, nil, nil, DefaultRPOThreshold, time.Now())
	assert.False(t, coverage.Covered)
	assert.False(t, coverage.Stale)
}

func Test_ComputeClusterBackupCoverage_staleBeyondRPO(t *testing.T) {
	now := time.Now()
	backups := []BackupInfo{completedBackup("old", "default", nil, 30*24*time.Hour, now)}

	coverage := ComputeClusterBackupCoverage(testClusterRef, backups, nil, DefaultRPOThreshold, now)
	assert.True(t, coverage.Covered)
	assert.True(t, coverage.Stale)
}

func Test_ComputeClusterBackupCoverage_partiallyFailedNotCountedAsCoveredButVisible(t *testing.T) {
	now := time.Now()
	b := BackupInfo{
		Name: "pf1", Namespace: "default", Phase: "PartiallyFailed",
		CompletionTimestamp: timePtr(now.Add(-time.Hour)),
	}

	coverage := ComputeClusterBackupCoverage(testClusterRef, []BackupInfo{b}, nil, DefaultRPOThreshold, now)
	assert.False(t, coverage.Covered)
	// Still visible, not hidden (spec.md Edge Cases), even though not counted as a verified recovery point.
	assert.Equal(t, "pf1", coverage.MostRecentBackupName)
}

func Test_ComputeClusterBackupCoverage_mostRecentBackupWins(t *testing.T) {
	now := time.Now()
	backups := []BackupInfo{
		completedBackup("older", "default", nil, 10*time.Hour, now),
		completedBackup("newer", "default", nil, time.Hour, now),
	}

	coverage := ComputeClusterBackupCoverage(testClusterRef, backups, nil, DefaultRPOThreshold, now)
	assert.Equal(t, "newer", coverage.MostRecentBackupName)
}

func Test_ComputeClusterBackupCoverage_restoreInProgress(t *testing.T) {
	now := time.Now()
	backups := []BackupInfo{completedBackup("b1", "default", nil, time.Hour, now)}
	restores := []RestoreInfo{{Name: "r1", BackupName: "b1", Phase: "InProgress"}}

	coverage := ComputeClusterBackupCoverage(testClusterRef, backups, restores, DefaultRPOThreshold, now)
	assert.True(t, coverage.RestoreInProgress)
	assert.Empty(t, coverage.LastRestoreOutcome)
}

func Test_ComputeClusterBackupCoverage_restoreOutcomes(t *testing.T) {
	now := time.Now()
	backups := []BackupInfo{completedBackup("b1", "default", nil, time.Hour, now)}

	succeeded := []RestoreInfo{{Name: "r1", BackupName: "b1", Phase: "Completed", CompletionTimestamp: timePtr(now.Add(-time.Minute))}}
	coverage := ComputeClusterBackupCoverage(testClusterRef, backups, succeeded, DefaultRPOThreshold, now)
	assert.Equal(t, "succeeded", coverage.LastRestoreOutcome)

	failed := []RestoreInfo{{Name: "r2", BackupName: "b1", Phase: "Failed", CompletionTimestamp: timePtr(now.Add(-time.Minute))}}
	coverage = ComputeClusterBackupCoverage(testClusterRef, backups, failed, DefaultRPOThreshold, now)
	assert.Equal(t, "failed", coverage.LastRestoreOutcome)
}

func Test_ComputeClusterBackupCoverage_restoreForUnrelatedBackupIgnored(t *testing.T) {
	now := time.Now()
	backups := []BackupInfo{completedBackup("b1", "other-namespace", []string{"other-namespace"}, time.Hour, now)}
	restores := []RestoreInfo{{Name: "r1", BackupName: "b1", Phase: "InProgress"}}

	coverage := ComputeClusterBackupCoverage(testClusterRef, backups, restores, DefaultRPOThreshold, now)
	assert.False(t, coverage.RestoreInProgress)
}

func Test_ComputeBackupHealth_notAvailableWhenVeleroNotInstalled(t *testing.T) {
	health := ComputeBackupHealth(false, []ObjectRef{testClusterRef}, nil, nil, nil, DefaultRPOThreshold, time.Now())
	assert.False(t, health.Available)
	assert.Empty(t, health.ClusterCoverage)
	assert.Empty(t, health.StorageLocations)
}

func Test_ComputeBackupHealth_reachableAndUnreachableLocations(t *testing.T) {
	locations := []BackupStorageLocationInfo{
		{Name: "default", Namespace: "velero", Phase: "Available", Default: true},
		{Name: "secondary", Namespace: "velero", Phase: "Unavailable"},
	}

	health := ComputeBackupHealth(true, nil, locations, nil, nil, DefaultRPOThreshold, time.Now())
	assert.Len(t, health.StorageLocations, 2)
	byName := map[string]bool{}
	for _, l := range health.StorageLocations {
		byName[l.Name] = l.Reachable
	}
	assert.True(t, byName["default"])
	assert.False(t, byName["secondary"])
}

func Test_ComputeBackupHealth_perClusterCoverageAndAggregateRestoresInProgress(t *testing.T) {
	now := time.Now()
	otherRef := ObjectRef{Namespace: "other-ns", Name: "other-cluster"}
	backups := []BackupInfo{
		completedBackup("b1", "default", []string{"default"}, time.Hour, now),
	}
	restores := []RestoreInfo{{Name: "r1", BackupName: "b1", Phase: "InProgress"}}

	health := ComputeBackupHealth(true, []ObjectRef{testClusterRef, otherRef}, nil, backups, restores, DefaultRPOThreshold, now)
	assert.Len(t, health.ClusterCoverage, 2)
	assert.Equal(t, 1, health.RestoresInProgress)
}

func timePtr(t time.Time) *time.Time { return &t }
