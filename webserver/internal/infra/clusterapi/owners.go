package clusterapi

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	capv "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
)

// ObjectInfo holds the information needed to process an object
type ObjectInfo struct {
	GVR       schema.GroupVersionResource
	Namespace string
	Name      string
	Index     int
}

// ErrOwnerHierarchyFetch indicates an error occurred while fetching the owner hierarchy
type ErrOwnerHierarchyFetch struct {
	Msg string
	Err error
}

func (e *ErrOwnerHierarchyFetch) Error() string {
	return fmt.Sprintf("failed to fetch owner hierarchy: %s: %v", e.Msg, e.Err)
}

// processOwnerHierarchy processes a list of VSphereMachine objects to build the owner-reference hierarchy for a cluster.
func processOwnerHierarchy(ctx context.Context, machines []capv.VSphereMachine) (topology ClusterTopology, err error) {
	var wg sync.WaitGroup
	// In the processOwnerHierarchy function, replace the direct initialization with:
	topology = NewClusterTopology()
	processMachine := func(wg *sync.WaitGroup, idx int, machine capv.VSphereMachine) error {
		defer wg.Done()
		gvr, _ := meta.UnsafeGuessKindToResource(machine.GroupVersionKind())
		return topology.fetchOwnerHierarchy(ctx, machine.OwnerReferences, ObjectInfo{
			GVR:       gvr,
			Namespace: machine.Namespace,
			Name:      machine.Name,
			Index:     idx,
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

	var current = cl.AddNode(currentResource)
	for idx, owner := range owners {
		// Adding owner node and edge - creates a directed edge from the current resource to its owner,
		// with the owner node being added to the graph if it doesn't already exist. This establishes
		// the parent-child relationship in the ownership hierarchy visualization.
		currentOwner := ObjectInfo{
			GVR:       convertGVK(owner),
			Namespace: currentResource.Namespace,
			Name:      owner.Name,
			Index:     idx,
		}
		cl.AddEdge(current, cl.AddNode(currentOwner))
		if err := cl.processParentOwner(ctx, dynamicClient, currentOwner); err != nil {
			return &ErrOwnerHierarchyFetch{Msg: "failed to process owner", Err: err}
		}
	}
	return nil
}

// processParentOwner handles the processing of a single owner reference
func (cl *ClusterTopology) processParentOwner(ctx context.Context, client *dynamic.DynamicClient, currentOwner ObjectInfo) error {
	parentOwners, err := cl.fetchOwnerReferences(client, currentOwner)
	if err != nil {
		return err
	}
	if len(parentOwners) > 0 {
		return cl.fetchOwnerHierarchy(ctx, parentOwners, currentOwner)
	}
	return nil
}

// FetchOwnerReferences retrieves the owner references of a Kubernetes resource identified by GroupVersionResource, namespace, and name.
// It returns a slice of OwnerReference objects or an error if the retrieval or conversion fails.
func (cl *ClusterTopology) fetchOwnerReferences(c *dynamic.DynamicClient, currentOwner ObjectInfo) (owners []metav1.OwnerReference, err error) {
	var resource *unstructured.Unstructured
	if resource, err = c.Resource(currentOwner.GVR).Namespace(currentOwner.Namespace).Get(
		context.TODO(), currentOwner.Name, metav1.GetOptions{},
	); err != nil {
		return owners, err
	}
	ownerList, exists, err := unstructured.NestedSlice(resource.Object, "metadata", "ownerReferences")
	if err != nil || !exists {
		return nil, err
	}
	var ownerRefs []metav1.OwnerReference
	for _, ownerRef := range ownerList {
		data, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&ownerRef)
		if err != nil {
			return nil, err
		}
		var ref metav1.OwnerReference
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(data, &ref); err != nil {
			return nil, err
		}
		ownerRefs = append(ownerRefs, ref)
	}
	return ownerRefs, nil
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
