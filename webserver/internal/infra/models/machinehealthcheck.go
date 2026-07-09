package models

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// MachineHealthCheck stores the details for a CAPI MachineHealthCheck — the remediation policy
// (target selector, timeouts, maxUnhealthy threshold) behind the Day-2 Ops dashboard's
// self-healing/needs-investigation severity classification (006/US4).
type MachineHealthCheck struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Age represents the creation time of the MachineHealthCheck, formatted as a duration string.
	Age string `json:"age"`

	// Cluster represents the name of the cluster this MachineHealthCheck belongs to.
	Cluster string `json:"cluster"`

	// Selector is the label selector matching the Machines this MachineHealthCheck targets.
	Selector metav1.LabelSelector `json:"selector"`

	// MaxUnhealthy is the maximum number/percentage of unhealthy Machines allowed before
	// remediation is paused. Empty when unset (the deprecated field was not configured).
	MaxUnhealthy string `json:"maxUnhealthy,omitempty"`

	// NodeStartupTimeout is the formatted duration a Node has to appear before being considered
	// unhealthy.
	NodeStartupTimeout string `json:"nodeStartupTimeout,omitempty"`

	// UnhealthyConditions are the Node conditions that, combined with an OR, mark a Machine
	// unhealthy.
	UnhealthyConditions []clusterv1.UnhealthyCondition `json:"unhealthyConditions,omitempty"`

	// Status represents the current remediation status (expected vs. currently healthy Machines).
	Status clusterv1.MachineHealthCheckStatus `json:"status"`
}
