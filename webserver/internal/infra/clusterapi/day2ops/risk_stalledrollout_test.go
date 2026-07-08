package day2ops

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func int32Ptr(i int32) *int32 { return &i }

func Test_ComputeStalledRolloutRisk_StalledPastGracePeriod(t *testing.T) {
	now := time.Now()
	oldMS := clusterv1.MachineSet{
		ObjectMeta: metav1.ObjectMeta{Name: "md-old", CreationTimestamp: metav1.NewTime(now.Add(-30 * time.Minute))},
		Spec:       clusterv1.MachineSetSpec{Replicas: int32Ptr(1)},
	}
	newMS := clusterv1.MachineSet{
		ObjectMeta: metav1.ObjectMeta{Name: "md-new", CreationTimestamp: metav1.NewTime(now.Add(-5 * time.Minute))},
		Spec:       clusterv1.MachineSetSpec{Replicas: int32Ptr(3)},
	}

	risk := ComputeStalledRolloutRisk(ObjectRef{Name: "md"}, []clusterv1.MachineSet{oldMS, newMS}, nil, now)

	require.NotNil(t, risk)
	assert.Equal(t, RiskStalledRollout, risk.Kind)
	assert.Contains(t, risk.Detail, "md-old")
}

func Test_ComputeStalledRolloutRisk_WithinGracePeriod(t *testing.T) {
	now := time.Now()
	oldMS := clusterv1.MachineSet{
		ObjectMeta: metav1.ObjectMeta{Name: "md-old", CreationTimestamp: metav1.NewTime(now.Add(-2 * time.Minute))},
		Spec:       clusterv1.MachineSetSpec{Replicas: int32Ptr(1)},
	}
	newMS := clusterv1.MachineSet{
		ObjectMeta: metav1.ObjectMeta{Name: "md-new", CreationTimestamp: metav1.NewTime(now)},
		Spec:       clusterv1.MachineSetSpec{Replicas: int32Ptr(3)},
	}

	risk := ComputeStalledRolloutRisk(ObjectRef{Name: "md"}, []clusterv1.MachineSet{oldMS, newMS}, nil, now)
	assert.Nil(t, risk, "a rollout still within the grace period is normal, not stalled")
}

func Test_ComputeStalledRolloutRisk_OnlyOneActiveMachineSet(t *testing.T) {
	now := time.Now()
	ms := clusterv1.MachineSet{
		ObjectMeta: metav1.ObjectMeta{Name: "md-only", CreationTimestamp: metav1.NewTime(now.Add(-time.Hour))},
		Spec:       clusterv1.MachineSetSpec{Replicas: int32Ptr(3)},
	}
	risk := ComputeStalledRolloutRisk(ObjectRef{Name: "md"}, []clusterv1.MachineSet{ms}, nil, now)
	assert.Nil(t, risk)
}

func Test_ComputeStalledRolloutRisk_LikelyCauseFromFinalizers(t *testing.T) {
	now := time.Now()
	oldMS := clusterv1.MachineSet{
		ObjectMeta: metav1.ObjectMeta{Name: "md-old", CreationTimestamp: metav1.NewTime(now.Add(-time.Hour))},
		Spec:       clusterv1.MachineSetSpec{Replicas: int32Ptr(1)},
	}
	newMS := clusterv1.MachineSet{
		ObjectMeta: metav1.ObjectMeta{Name: "md-new", CreationTimestamp: metav1.NewTime(now)},
		Spec:       clusterv1.MachineSetSpec{Replicas: int32Ptr(3)},
	}

	risk := ComputeStalledRolloutRisk(ObjectRef{Name: "md"}, []clusterv1.MachineSet{oldMS, newMS}, []string{"my.finalizer/protect"}, now)
	require.NotNil(t, risk)
	assert.Contains(t, risk.LikelyCause, "my.finalizer/protect")
}
