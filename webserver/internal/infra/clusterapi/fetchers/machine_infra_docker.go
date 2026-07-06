package fetchers

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	"github.com/knabben/observatio/webserver/internal/infra/models"
)

// dockerMachineGVR identifies the DockerMachine resource, read via the dynamic client since
// no typed Go package for it is imported (see specs/004-detect-infra-adapt-ui/research.md R3).
var dockerMachineGVR = schema.GroupVersionResource{
	Group:    "infrastructure.cluster.x-k8s.io",
	Version:  "v1beta1",
	Resource: "dockermachines",
}

// FetchMachineInfraDocker retrieves and processes all DockerMachine resources into a
// MachineInfraDockerResponse, mirroring FetchMachineInfra for vSphere.
func FetchMachineInfraDocker(ctx context.Context, dyn dynamic.Interface) (models.MachineInfraDockerResponse, error) {
	machines, err := ListMachineInfraDocker(ctx, dyn)
	if err != nil {
		return models.MachineInfraDockerResponse{}, err
	}

	failing := 0
	for _, m := range machines {
		if !m.Ready {
			failing++
		}
	}
	return models.MachineInfraDockerResponse{
		Total:    len(machines),
		Failing:  failing,
		Machines: machines,
	}, nil
}

// ListMachineInfraDocker lists all DockerMachine resources across namespaces via the dynamic
// client and decodes only the fields the Docker infra view needs.
func ListMachineInfraDocker(ctx context.Context, dyn dynamic.Interface) ([]models.MachineInfraDocker, error) {
	list, err := dyn.Resource(dockerMachineGVR).Namespace("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	machines := make([]models.MachineInfraDocker, 0, len(list.Items))
	for _, item := range list.Items {
		machines = append(machines, ProcessDockerMachine(item))
	}
	return machines, nil
}

// ProcessDockerMachine decodes a single unstructured DockerMachine object into a
// models.MachineInfraDocker.
func ProcessDockerMachine(obj unstructured.Unstructured) models.MachineInfraDocker {
	providerID, _, _ := unstructured.NestedString(obj.Object, "spec", "providerID")
	ready, _, _ := unstructured.NestedBool(obj.Object, "status", "ready")

	return models.MachineInfraDocker{
		ObjectMeta: metav1.ObjectMeta{
			Name:              obj.GetName(),
			Namespace:         obj.GetNamespace(),
			CreationTimestamp: obj.GetCreationTimestamp(),
		},
		Age:        formatDuration(time.Since(obj.GetCreationTimestamp().Time)),
		ProviderID: providerID,
		Ready:      ready,
	}
}
