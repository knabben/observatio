package day2ops

import (
	"fmt"
	"time"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// StalledRolloutGracePeriod is how long an old MachineSet may coexist with a newer one before
// being flagged as stalled rather than a normal in-progress rollout.
const StalledRolloutGracePeriod = 10 * time.Minute

// ComputeStalledRolloutRisk flags a MachineDeployment whose old MachineSet hasn't scaled down
// within the grace period while a newer one exists alongside it (FR-009). oldMachineFinalizers
// lists any finalizers observed on Machines owned by the oldest active MachineSet — a directly
// observable likely cause (research.md R5's finalizer case). The PodDisruptionBudget case would
// require workload-cluster access and is out of scope for this pass (spec.md Assumptions).
func ComputeStalledRolloutRisk(objectRef ObjectRef, machineSets []clusterv1.MachineSet, oldMachineFinalizers []string, now time.Time) *RiskWarning {
	active := make([]clusterv1.MachineSet, 0, len(machineSets))
	for _, ms := range machineSets {
		replicas := int32(1)
		if ms.Spec.Replicas != nil {
			replicas = *ms.Spec.Replicas
		}
		if replicas > 0 || ms.Status.Replicas > 0 {
			active = append(active, ms)
		}
	}
	if len(active) < 2 {
		return nil
	}

	oldest := &active[0]
	for i := range active[1:] {
		if active[i+1].CreationTimestamp.Before(&oldest.CreationTimestamp) {
			oldest = &active[i+1]
		}
	}
	age := now.Sub(oldest.CreationTimestamp.Time)
	if age < StalledRolloutGracePeriod {
		return nil
	}

	likelyCause := ""
	if len(oldMachineFinalizers) > 0 {
		likelyCause = fmt.Sprintf("blocked by finalizer(s) on its Machines: %v", oldMachineFinalizers)
	}
	return &RiskWarning{
		ObjectRef:   objectRef,
		Kind:        RiskStalledRollout,
		Detail:      fmt.Sprintf("MachineSet %s has not scaled down after %s", oldest.Name, age.Round(time.Minute)),
		LikelyCause: likelyCause,
		CheckStatus: RiskCheckEvaluated,
	}
}
