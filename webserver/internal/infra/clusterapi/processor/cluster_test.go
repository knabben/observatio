package processor

import (
	"testing"
	"time"

	"github.com/knabben/observatio/webserver/internal/infra/models"
	"github.com/stretchr/testify/assert"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	capv "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func TestProcessClusterInfra(t *testing.T) {
	created := time.Now()
	tests := []struct {
		name     string
		cluster  capv.VSphereCluster
		expected models.ClusterInfra
	}{
		{
			name: "Cluster with no owner references",
			cluster: capv.VSphereCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster1",
					Namespace: "default",
					CreationTimestamp: metav1.Time{
						Time: created.Add(-10 * time.Minute),
					},
				},
				Spec: capv.VSphereClusterSpec{},
				Status: capv.VSphereClusterStatus{
					Ready: true,
				},
			},
			expected: models.ClusterInfra{
				Name:                 "cluster1",
				Namespace:            "default",
				Cluster:              "",
				Server:               "",
				Thumbprint:           "",
				Created:              "10m",
				ControlPlaneEndpoint: ":0",
				Modules:              nil,
				Conditions:           nil,
				Ready:                true,
			},
		},
		{
			name: "Cluster with owner references",
			cluster: capv.VSphereCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster2",
					Namespace: "test",
					OwnerReferences: []metav1.OwnerReference{
						{Name: "owner-cluster"},
					},
					CreationTimestamp: metav1.Time{
						Time: created.Add(-1 * time.Hour),
					},
				},
				Spec: capv.VSphereClusterSpec{
					Server: "https://server.example.com",
				},
				Status: capv.VSphereClusterStatus{
					Ready: false,
				},
			},
			expected: models.ClusterInfra{
				Name:                 "cluster2",
				Namespace:            "test",
				Cluster:              "owner-cluster",
				Server:               "https://server.example.com",
				Thumbprint:           "",
				Created:              "1h 0m",
				ControlPlaneEndpoint: ":0",
				Modules:              nil,
				Conditions:           nil,
				Ready:                false,
			},
		},
		{
			name: "Cluster with no thumbprint or server specified",
			cluster: capv.VSphereCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster3",
					Namespace: "custom",
					CreationTimestamp: metav1.Time{
						Time: created.Add(-30 * time.Second),
					},
				},
				Spec: capv.VSphereClusterSpec{
					Thumbprint: "",
					Server:     "",
				},
				Status: capv.VSphereClusterStatus{},
			},
			expected: models.ClusterInfra{
				Name:                 "cluster3",
				Namespace:            "custom",
				Cluster:              "",
				Server:               "",
				Thumbprint:           "",
				Created:              "0m",
				ControlPlaneEndpoint: ":0",
				Modules:              nil,
				Conditions:           nil,
				Ready:                false,
			},
		},
		{
			name: "Cluster with empty conditions and modules",
			cluster: capv.VSphereCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster4",
					Namespace: "dev",
					CreationTimestamp: metav1.Time{
						Time: time.Now().Add(-5 * time.Hour),
					},
				},
				Spec: capv.VSphereClusterSpec{
					ClusterModules: []capv.ClusterModule{},
				},
				Status: capv.VSphereClusterStatus{
					Conditions: nil,
					Ready:      true,
				},
			},
			expected: models.ClusterInfra{
				Name:                 "cluster4",
				Namespace:            "dev",
				Cluster:              "",
				Server:               "",
				Thumbprint:           "",
				Created:              "5h 0m",
				ControlPlaneEndpoint: ":0",
				Modules:              []capv.ClusterModule{},
				Conditions:           nil,
				Ready:                true,
			},
		},
		{
			name: "Cluster with populated control plane endpoint",
			cluster: capv.VSphereCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster5",
					Namespace: "prod",
					CreationTimestamp: metav1.Time{
						Time: created.Add(-2 * time.Minute),
					},
				},
				Spec: capv.VSphereClusterSpec{
					ControlPlaneEndpoint: capv.APIEndpoint{
						Host: "192.168.1.1",
						Port: 6443,
					},
				},
				Status: capv.VSphereClusterStatus{
					Ready: false,
				},
			},
			expected: models.ClusterInfra{
				Name:                 "cluster5",
				Namespace:            "prod",
				Cluster:              "",
				Server:               "",
				Thumbprint:           "",
				Created:              "2m",
				ControlPlaneEndpoint: "192.168.1.1:6443",
				Modules:              nil,
				Conditions:           nil,
				Ready:                false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessClusterInfra(tt.cluster)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestProcessCluster(t *testing.T) {
	tests := []struct {
		name     string
		cluster  clusterv1.Cluster
		expected models.Cluster
	}{
		{
			name: "Cluster without topology",
			cluster: clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster1",
					Namespace: "default",
				},
				Spec: clusterv1.ClusterSpec{
					Paused: false,
				},
				Status: clusterv1.ClusterStatus{
					Phase:               "Provisioning",
					InfrastructureReady: false,
					ControlPlaneReady:   false,
				},
			},
			expected: models.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster1",
					Namespace: "default",
				},
				Paused:   false,
				Topology: models.ClusterClassType{IsClusterClass: false},
			},
		},
		{
			name: "Cluster with nil topology but paused",
			cluster: clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster3",
					Namespace: "test",
				},
				Spec: clusterv1.ClusterSpec{
					Paused: true,
				},
				Status: clusterv1.ClusterStatus{
					Phase:               "Paused",
					InfrastructureReady: false,
					ControlPlaneReady:   false,
				},
			},
			expected: models.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster3",
					Namespace: "test",
				},
				Paused:   true,
				Topology: models.ClusterClassType{IsClusterClass: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessCluster(tt.cluster)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.Namespace, result.Namespace)
			assert.Equal(t, tt.expected.Paused, result.Paused)
			assert.Equal(t, tt.expected.Topology, result.Topology)
		})
	}
}

func TestProcessClusterResponse(t *testing.T) {
	tests := []struct {
		name     string
		clusters []clusterv1.Cluster
		expected models.ClusterResponse
	}{
		{
			name:     "No clusters",
			clusters: []clusterv1.Cluster{},
			expected: models.ClusterResponse{
				Total:    0,
				Failing:  0,
				Clusters: []models.Cluster{},
			},
		},
		{
			name: "All clusters healthy",
			clusters: []clusterv1.Cluster{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "healthy-cluster-1",
					},
					Status: clusterv1.ClusterStatus{
						Phase:               "Running",
						InfrastructureReady: true,
						ControlPlaneReady:   true,
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "healthy-cluster-2",
					},
					Status: clusterv1.ClusterStatus{
						Phase:               "Running",
						InfrastructureReady: true,
						ControlPlaneReady:   true,
					},
				},
			},
			expected: models.ClusterResponse{
				Total:   2,
				Failing: 0,
				Clusters: []models.Cluster{
					ProcessCluster(clusterv1.Cluster{
						ObjectMeta: metav1.ObjectMeta{
							Name: "healthy-cluster-1",
						},
						Status: clusterv1.ClusterStatus{
							Phase:               "Running",
							InfrastructureReady: true,
							ControlPlaneReady:   true,
						},
					}),
					ProcessCluster(clusterv1.Cluster{
						ObjectMeta: metav1.ObjectMeta{
							Name: "healthy-cluster-2",
						},
						Status: clusterv1.ClusterStatus{
							Phase:               "Running",
							InfrastructureReady: true,
							ControlPlaneReady:   true,
						},
					}),
				},
			},
		},
		{
			name: "Some clusters failing",
			clusters: []clusterv1.Cluster{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "healthy-cluster",
					},
					Status: clusterv1.ClusterStatus{
						Phase:               "Running",
						InfrastructureReady: true,
						ControlPlaneReady:   true,
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "failed-cluster",
					},
					Status: clusterv1.ClusterStatus{
						Phase:               "Failed",
						InfrastructureReady: false,
					},
				},
			},
			expected: models.ClusterResponse{
				Total:   2,
				Failing: 1,
				Clusters: []models.Cluster{
					ProcessCluster(clusterv1.Cluster{
						ObjectMeta: metav1.ObjectMeta{
							Name: "healthy-cluster",
						},
						Status: clusterv1.ClusterStatus{
							Phase:               "Running",
							InfrastructureReady: true,
							ControlPlaneReady:   true,
						},
					}),
					ProcessCluster(clusterv1.Cluster{
						ObjectMeta: metav1.ObjectMeta{
							Name: "failed-cluster",
						},
						Status: clusterv1.ClusterStatus{
							Phase:               "Failed",
							InfrastructureReady: false,
						},
					}),
				},
			},
		},
		{
			name: "All clusters failing",
			clusters: []clusterv1.Cluster{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "failed-cluster-1",
					},
					Status: clusterv1.ClusterStatus{
						Phase: "Failed",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "failed-cluster-2",
					},
					Status: clusterv1.ClusterStatus{
						Phase: "Failed",
					},
				},
			},
			expected: models.ClusterResponse{
				Total:   2,
				Failing: 2,
				Clusters: []models.Cluster{
					ProcessCluster(clusterv1.Cluster{
						ObjectMeta: metav1.ObjectMeta{
							Name: "failed-cluster-1",
						},
						Status: clusterv1.ClusterStatus{
							Phase: "Failed",
						},
					}),
					ProcessCluster(clusterv1.Cluster{
						ObjectMeta: metav1.ObjectMeta{
							Name: "failed-cluster-2",
						},
						Status: clusterv1.ClusterStatus{
							Phase: "Failed",
						},
					}),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessClusterResponse(tt.clusters)
			assert.Equal(t, tt.expected.Total, result.Total)
			assert.Equal(t, tt.expected.Failing, result.Failing)
			assert.Equal(t, tt.expected.Clusters, result.Clusters)
		})
	}
}

func TestProcessClusterInfraResponse(t *testing.T) {
	tests := []struct {
		name     string
		clusters []capv.VSphereCluster
		expected models.ClusterInfraResponse
	}{
		{
			name:     "No clusters",
			clusters: []capv.VSphereCluster{},
			expected: models.ClusterInfraResponse{
				Total:    0,
				Failing:  0,
				Clusters: []models.ClusterInfra{},
			},
		},
		{
			name: "All clusters ready",
			clusters: []capv.VSphereCluster{
				{
					Status: capv.VSphereClusterStatus{Ready: true},
				},
				{
					Status: capv.VSphereClusterStatus{Ready: true},
				},
			},
			expected: models.ClusterInfraResponse{
				Total:   2,
				Failing: 0,
				Clusters: []models.ClusterInfra{
					ProcessClusterInfra(capv.VSphereCluster{Status: capv.VSphereClusterStatus{Ready: true}}),
					ProcessClusterInfra(capv.VSphereCluster{Status: capv.VSphereClusterStatus{Ready: true}}),
				},
			},
		},
		{
			name: "Some clusters not ready",
			clusters: []capv.VSphereCluster{
				{
					Status: capv.VSphereClusterStatus{Ready: true},
				},
				{
					Status: capv.VSphereClusterStatus{Ready: false},
				},
			},
			expected: models.ClusterInfraResponse{
				Total:   2,
				Failing: 1,
				Clusters: []models.ClusterInfra{
					ProcessClusterInfra(capv.VSphereCluster{Status: capv.VSphereClusterStatus{Ready: true}}),
					ProcessClusterInfra(capv.VSphereCluster{Status: capv.VSphereClusterStatus{Ready: false}}),
				},
			},
		},
		{
			name: "None of the clusters ready",
			clusters: []capv.VSphereCluster{
				{
					Status: capv.VSphereClusterStatus{Ready: false},
				},
				{
					Status: capv.VSphereClusterStatus{Ready: false},
				},
			},
			expected: models.ClusterInfraResponse{
				Total:   2,
				Failing: 2,
				Clusters: []models.ClusterInfra{
					ProcessClusterInfra(capv.VSphereCluster{Status: capv.VSphereClusterStatus{Ready: false}}),
					ProcessClusterInfra(capv.VSphereCluster{Status: capv.VSphereClusterStatus{Ready: false}}),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessClusterInfraResponse(tt.clusters)
			assert.Equal(t, tt.expected.Total, result.Total)
			assert.Equal(t, tt.expected.Failing, result.Failing)
			assert.Equal(t, tt.expected.Clusters, result.Clusters)
		})
	}
}
