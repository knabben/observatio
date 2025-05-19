package processor

import (
	"testing"
	"time"

	"github.com/knabben/observatio/webserver/internal/infra/models"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	capv "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func TestProcessMachine(t *testing.T) {
	tests := []struct {
		name   string
		input  clusterv1.Machine
		expect models.Machine
	}{
		{
			name: "Complete machine details",
			input: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "test-machine",
					Namespace:         "test-namespace",
					CreationTimestamp: metav1.Time{Time: time.Now().Add(-5 * time.Minute)},
					OwnerReferences: []metav1.OwnerReference{{
						Name: "owner-1",
					}},
				},
				Spec: clusterv1.MachineSpec{
					ClusterName: "test-cluster",
					Version:     strPointer("v1.2.3"),
					ProviderID:  strPointer("provider-id"),
					Bootstrap: clusterv1.Bootstrap{
						ConfigRef: &corev1.ObjectReference{
							Name: "bootstrap-config",
						},
					},
				},
				Status: clusterv1.MachineStatus{
					NodeRef: &corev1.ObjectReference{
						Name: "node-ref",
					},
					BootstrapReady:      true,
					InfrastructureReady: true,
					Phase:               "Running",
				},
			},
			expect: models.Machine{
				Name:                "test-machine",
				Namespace:           "test-namespace",
				Owner:               "owner-1",
				Cluster:             "test-cluster",
				NodeName:            "node-ref",
				ProviderID:          "provider-id",
				Version:             "v1.2.3",
				BootstrapReady:      true,
				InfrastructureReady: true,
				Bootstrap:           "bootstrap-config",
				Phase:               "Running",
			},
		},
		{
			name: "Missing optional details",
			input: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "test-machine",
					Namespace:         "test-namespace",
					CreationTimestamp: metav1.Time{Time: time.Now().Add(-10 * time.Minute)},
				},
				Spec: clusterv1.MachineSpec{
					ClusterName: "test-cluster",
				},
				Status: clusterv1.MachineStatus{},
			},
			expect: models.Machine{
				Name:                "test-machine",
				Namespace:           "test-namespace",
				Cluster:             "test-cluster",
				BootstrapReady:      false,
				InfrastructureReady: false,
			},
		},
		{
			name: "Multiple owner references, last one selected",
			input: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "test-machine",
					Namespace:         "test-namespace",
					CreationTimestamp: metav1.Time{Time: time.Now().Add(-15 * time.Minute)},
					OwnerReferences: []metav1.OwnerReference{
						{Name: "owner-1"},
						{Name: "owner-2"},
					},
				},
			},
			expect: models.Machine{
				Name:      "test-machine",
				Namespace: "test-namespace",
				Owner:     "owner-2",
			},
		},
		{
			name: "No OwnerReferences or NodeRef",
			input: clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "machine-no-owners",
					Namespace:         "default",
					CreationTimestamp: metav1.Time{Time: time.Now()},
				},
			},
			expect: models.Machine{
				Name:      "machine-no-owners",
				Namespace: "default",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessMachine(tt.input)
			assert.Equal(t, tt.expect.Name, result.Name)
			assert.Equal(t, tt.expect.Namespace, result.Namespace)
			assert.Equal(t, tt.expect.Owner, result.Owner)
			assert.Equal(t, tt.expect.Cluster, result.Cluster)
			assert.Equal(t, tt.expect.NodeName, result.NodeName)
			assert.Equal(t, tt.expect.ProviderID, result.ProviderID)
			assert.Equal(t, tt.expect.Version, result.Version)
			assert.Equal(t, tt.expect.BootstrapReady, result.BootstrapReady)
			assert.Equal(t, tt.expect.InfrastructureReady, result.InfrastructureReady)
		})
	}
}

func TestProcessMachineInfra(t *testing.T) {
	tests := []struct {
		name   string
		input  capv.VSphereMachine
		expect models.MachineInfra
	}{
		{
			name: "All fields populated",
			input: capv.VSphereMachine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "infra-machine-1",
					Namespace: "test-namespace",
				},
				Spec: capv.VSphereMachineSpec{
					ProviderID:    strPointer("provider-id"),
					FailureDomain: strPointer("failure-domain-1"),
					PowerOffMode:  capv.VirtualMachinePowerOpModeTrySoft,
					VirtualMachineCloneSpec: capv.VirtualMachineCloneSpec{
						CloneMode:         capv.LinkedClone,
						NumCPUs:           4,
						NumCoresPerSocket: 2,
						MemoryMiB:         8192,
						DiskGiB:           100,
					},
				},
				Status: capv.VSphereMachineStatus{
					Ready: true,
				},
			},
			expect: models.MachineInfra{
				ObjectMeta:        metav1.ObjectMeta{Name: "infra-machine-1", Namespace: "test-namespace"},
				ProviderID:        "provider-id",
				FailureDomain:     "failure-domain-1",
				PowerOffMode:      capv.VirtualMachinePowerOpModeTrySoft,
				CloneMode:         capv.LinkedClone,
				NumCPUs:           4,
				NumCoresPerSocket: 2,
				MemoryMiB:         8192,
				DiskGiB:           100,
				Status: capv.VSphereMachineStatus{
					Ready: true,
				},
			},
		},
		{
			name: "Minimal fields populated",
			input: capv.VSphereMachine{
				ObjectMeta: metav1.ObjectMeta{
					Name: "infra-machine-2",
				},
				Spec: capv.VSphereMachineSpec{},
			},
			expect: models.MachineInfra{
				ObjectMeta: metav1.ObjectMeta{Name: "infra-machine-2"},
			},
		},
		{
			name: "ProviderID and FailureDomain only",
			input: capv.VSphereMachine{
				ObjectMeta: metav1.ObjectMeta{
					Name: "infra-machine-3",
				},
				Spec: capv.VSphereMachineSpec{
					ProviderID:    strPointer("provider-id-3"),
					FailureDomain: strPointer("failure-domain-3"),
				},
			},
			expect: models.MachineInfra{
				ObjectMeta:    metav1.ObjectMeta{Name: "infra-machine-3"},
				ProviderID:    "provider-id-3",
				FailureDomain: "failure-domain-3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessMachineInfra(tt.input)

			assert.Equal(t, tt.expect.ObjectMeta, result.ObjectMeta)
			assert.Equal(t, tt.expect.ProviderID, result.ProviderID)
			assert.Equal(t, tt.expect.FailureDomain, result.FailureDomain)
			assert.Equal(t, tt.expect.PowerOffMode, result.PowerOffMode)
			assert.Equal(t, tt.expect.CloneMode, result.CloneMode)
			assert.Equal(t, tt.expect.NumCPUs, result.NumCPUs)
			assert.Equal(t, tt.expect.NumCoresPerSocket, result.NumCoresPerSocket)
			assert.Equal(t, tt.expect.MemoryMiB, result.MemoryMiB)
			assert.Equal(t, tt.expect.DiskGiB, result.DiskGiB)
			assert.Equal(t, tt.expect.Status, result.Status)
		})
	}
}

func strPointer(s string) *string {
	return &s
}
