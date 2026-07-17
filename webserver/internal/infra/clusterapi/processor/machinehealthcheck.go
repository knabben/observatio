package processor

import (
	"time"

	"github.com/knabben/observatio/webserver/internal/infra/models"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// ProcessMachineHealthCheck transforms a clusterv1.MachineHealthCheck into a models.MachineHealthCheck.
func ProcessMachineHealthCheck(mhc clusterv1.MachineHealthCheck) models.MachineHealthCheck {
	// mhc.Spec.MaxUnhealthy is deprecated in v1beta1 pending CAPI #10722 (no replacement exists yet
	// in this API version), so it's read exactly once here rather than suppressing the linter at
	// every usage site.
	maxUnhealthyThreshold := mhc.Spec.MaxUnhealthy // nolint:staticcheck
	var maxUnhealthy string
	if maxUnhealthyThreshold != nil {
		maxUnhealthy = maxUnhealthyThreshold.String()
	}

	var nodeStartupTimeout string
	if mhc.Spec.NodeStartupTimeout != nil {
		nodeStartupTimeout = mhc.Spec.NodeStartupTimeout.Duration.String()
	}

	return models.MachineHealthCheck{
		ObjectMeta:          mhc.ObjectMeta,
		Age:                 formatDuration(time.Since(mhc.CreationTimestamp.Time)),
		Cluster:             mhc.Spec.ClusterName,
		Selector:            mhc.Spec.Selector,
		MaxUnhealthy:        maxUnhealthy,
		NodeStartupTimeout:  nodeStartupTimeout,
		UnhealthyConditions: mhc.Spec.UnhealthyConditions,
		Status:              mhc.Status,
	}
}
