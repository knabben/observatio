package models

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// MachineDeploymentResponse returns the list of available machine
// deployments.
type MachineDeploymentResponse struct {
	Total              int                 `json:"total"`
	Failing            int                 `json:"failing"`
	MachineDeployments []MachineDeployment `json:"machineDeployments"`
}

// MachineDeployment stores the details for CAPI machine deployments objects
type MachineDeployment struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Replicas represents the desired number of machine instances for the machine deployment.
	Replicas int32 `json:"replicas"`

	// Cluster represents the name of the cluster associated with the machine deployment.
	Cluster string `json:"cluster"`

	// Created represents the timestamp when the machine deployment was created.
	Age string `json:"created"`

	// bootstrap is a reference to a local struct which encapsulates
	// fields to configure the Machine’s bootstrapping mechanism.
	TemplateBootstrap clusterv1.Bootstrap `json:"templateBootstrap"`

	// TemplateInfrastructureRef is a required reference to a custom resource
	// offered by an infrastructure provider.
	TemplateInfrastructureRef corev1.ObjectReference `json:"templateInfrastructureRef"`

	// TemplateVersion defines the desired Kubernetes version.
	// This field is meant to be optionally used by bootstrap providers.
	// +optional
	TemplateVersion *string `json:"templateversion,omitempty"`

	// Status represents the current status details of a machine deployment's infrastructure.
	Status clusterv1.MachineDeploymentStatus `json:"status"`
}
