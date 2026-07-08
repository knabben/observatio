package clusterapi

import (
	"context"
	"fmt"
	"strings"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	"github.com/knabben/observatio/webserver/internal/infra/models"
)

// machineGVR identifies the Machine resource, the provider-agnostic seed for the topology walk
// (see GenerateClusterTopology's doc comment for why this replaced a VSphereMachine-only seed).
var machineGVR = schema.GroupVersionResource{Group: "cluster.x-k8s.io", Version: "v1beta1", Resource: "machines"}

// ErrOwnerHierarchyFetch indicates an error occurred while fetching the owner hierarchy
type ErrOwnerHierarchyFetch struct {
	Msg string
	Err error
}

func (e *ErrOwnerHierarchyFetch) Error() string {
	return fmt.Sprintf("failed to fetch owner hierarchy: %s: %v", e.Msg, e.Err)
}

// processOwnerHierarchy processes a list of Machines to build the owner-reference hierarchy for a
// cluster, provider-agnostically (Docker, vSphere, or any other infrastructure provider).
func processOwnerHierarchy(ctx context.Context, machines []models.Machine) (topology ClusterTopology, err error) {
	var wg sync.WaitGroup
	topology = NewClusterTopology()
	processMachine := func(wg *sync.WaitGroup, idx int, machine models.Machine) error {
		defer wg.Done()
		return topology.fetchOwnerHierarchy(ctx, machine.OwnerReferences, ObjectInfo{
			GVR:       machineGVR,
			Namespace: machine.Namespace,
			Name:      machine.Name,
			PositionX: int32(idx) * 150,
			PositionY: 0,
			Failed:    !machine.Status.InfrastructureReady || !machine.Status.BootstrapReady,
		})
	}

	for idx, machine := range machines {
		wg.Add(1)
		go func() {
			err := processMachine(&wg, idx, machine)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
	wg.Wait()
	return topology, nil
}

// fetchOwnerHierarchy traverses and constructs the ownership hierarchy of a Kubernetes resource for visualization or analysis.
func (cl *ClusterTopology) fetchOwnerHierarchy(ctx context.Context, owners []metav1.OwnerReference, currentResource ObjectInfo) error {
	dynamicClient, err := NewDynamicClient(ctx)
	if err != nil {
		return &ErrOwnerHierarchyFetch{Msg: "failed to create dynamic client", Err: err}
	}

	current := cl.AddNode(currentResource)
	for _, owner := range owners {
		ownerInfo := ObjectInfo{
			GVR:       convertGVK(owner),
			Namespace: currentResource.Namespace,
			Name:      owner.Name,
			PositionX: currentResource.PositionX,
			PositionY: currentResource.PositionY + 150,
		}
		// Fetch the owner's object once: it tells us both its own health (for status coloring)
		// and its further owner references (to keep walking up the hierarchy).
		parentOwners, failed, err := cl.fetchOwnerReferences(dynamicClient, ownerInfo)
		if err != nil {
			return &ErrOwnerHierarchyFetch{Msg: "failed to process owner", Err: err}
		}
		ownerInfo.Failed = failed
		cl.AddEdge(current, cl.AddNode(ownerInfo))
		if len(parentOwners) > 0 {
			if err := cl.fetchOwnerHierarchy(ctx, parentOwners, ownerInfo); err != nil {
				return err
			}
		}
	}
	return nil
}

// fetchOwnerReferences retrieves a resource's owner references and a generic health signal
// (status.ready, or a Ready condition), reading both from the same fetched object.
func (cl *ClusterTopology) fetchOwnerReferences(c *dynamic.DynamicClient, currentOwner ObjectInfo) (owners []metav1.OwnerReference, failed bool, err error) {
	var resource *unstructured.Unstructured
	if resource, err = c.Resource(currentOwner.GVR).Namespace(currentOwner.Namespace).Get(
		context.TODO(), currentOwner.Name, metav1.GetOptions{},
	); err != nil {
		return owners, false, err
	}
	failed = isUnstructuredUnhealthy(resource)

	ownerList, exists, err := unstructured.NestedSlice(resource.Object, "metadata", "ownerReferences")
	if err != nil || !exists {
		return nil, failed, nil
	}
	var ownerRefs []metav1.OwnerReference
	for _, ownerRef := range ownerList {
		data, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&ownerRef)
		if err != nil {
			return nil, failed, err
		}
		var ref metav1.OwnerReference
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(data, &ref); err != nil {
			return nil, failed, err
		}
		ownerRefs = append(ownerRefs, ref)
	}
	return ownerRefs, failed, nil
}

// isUnstructuredUnhealthy reports whether a generic CAPI object looks unhealthy, checked
// generically (no per-provider/per-kind branching, Constitution Principle III) via
// `status.ready` first, falling back to a `Ready` condition. An object exposing neither is
// treated as healthy rather than guessed at, to avoid false-positive failure coloring.
func isUnstructuredUnhealthy(obj *unstructured.Unstructured) bool {
	if ready, found, _ := unstructured.NestedBool(obj.Object, "status", "ready"); found {
		return !ready
	}
	conditions, _, _ := unstructured.NestedSlice(obj.Object, "status", "conditions")
	for _, raw := range conditions {
		cond, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		condType, _, _ := unstructured.NestedString(cond, "type")
		if condType != "Ready" {
			continue
		}
		condStatus, _, _ := unstructured.NestedString(cond, "status")
		return condStatus != "True"
	}
	return false
}

// convertGVK converts a metav1.OwnerReference object into a schema.GroupVersionResource representation.
// It extracts the API version to determine the Group and Version, and derives the Resource name from the Kind.
func convertGVK(obj metav1.OwnerReference) schema.GroupVersionResource {
	splits := strings.Split(obj.APIVersion, "/")
	return schema.GroupVersionResource{
		Group:    splits[0],
		Version:  splits[1],
		Resource: strings.ToLower(obj.Kind) + "s", // assuming it's following the usual convention of kind to resource conversion
	}
}
