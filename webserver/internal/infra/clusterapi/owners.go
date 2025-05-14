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

// Node represents a graph node with an identifier,
// metadata, and its positional coordinates.
type Node struct {
	Id   string `json:"id"`
	Data struct {
		Label string `json:"label"`
	} `json:"data"`
	Position struct {
		X int32 `json:"x"`
		Y int32 `json:"y"`
	} `json:"position"`
}

// Edge represents a single connection between two nodes in a graph,
// identified by an Id with source and destination nodes.
type Edge struct {
	Id     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

type ClusterTopology struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

// AddNode adds a new node to the ClusterTopology if it does not already exist and returns the created or existing node.
func (cl *ClusterTopology) AddNode(gvr schema.GroupVersionResource, namespace string, name string) (node Node) {
	id := name + gvr.Group + gvr.Version + namespace
	node = Node{
		Id: id, Data: struct {
			Label string `json:"label"`
		}{Label: name},
	}
	if !cl.Find(&node) {
		cl.Nodes = append(cl.Nodes, node)
	}
	fmt.Println("ADDING: ", id, name)
	return node
}

// AddEdge creates a directed edge between the current and owner nodes and appends it to the cluster topology's edges.
func (cl *ClusterTopology) AddEdge(current, owner Node) {
	cl.Edges = append(cl.Edges, Edge{
		Id:     current.Id + owner.Id,
		Source: current.Id,
		Target: owner.Id,
	})
}

// Find checks if a given node exists in the ClusterTopology's Nodes slice by comparing the node's Id. Returns true if found.
func (cl *ClusterTopology) Find(node *Node) bool {
	for _, n := range cl.Nodes {
		if n.Id == node.Id {
			return true
		}
	}
	return false
}

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
	topology = ClusterTopology{Nodes: make([]Node, 0), Edges: make([]Edge, 0)}

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

	var current = cl.AddNode(currentResource.GVR, currentResource.Namespace, currentResource.Name)
	for idx, owner := range owners {
		// Adding owner node and edge - creates a directed edge from the current resource to its owner,
		// with the owner node being added to the graph if it doesn't already exist. This establishes
		// the parent-child relationship in the ownership hierarchy visualization.
		cl.AddEdge(current, cl.AddNode(convertGVK(owner), currentResource.Namespace, owner.Name))
		if err := cl.processParentOwner(ctx, dynamicClient, idx, owner, currentResource.Namespace); err != nil {
			return &ErrOwnerHierarchyFetch{Msg: "failed to process owner", Err: err}
		}
	}
	return nil
}

// processParentOwner handles the processing of a single owner reference
func (cl *ClusterTopology) processParentOwner(ctx context.Context, client *dynamic.DynamicClient, idx int, owner metav1.OwnerReference, namespace string) error {
	parentOwners, err := cl.fetchOwnerReferences(client, convertGVK(owner), namespace, owner.Name)
	if err != nil {
		return err
	}

	if len(parentOwners) > 0 {
		currentOwner := ObjectInfo{
			GVR:       convertGVK(owner),
			Namespace: namespace,
			Name:      owner.Name,
			Index:     idx,
		}
		return cl.fetchOwnerHierarchy(ctx, parentOwners, currentOwner)
	}
	return nil
}

// FetchOwnerReferences retrieves the owner references of a Kubernetes resource identified by GroupVersionResource, namespace, and name.
// It returns a slice of OwnerReference objects or an error if the retrieval or conversion fails.
func (cl *ClusterTopology) fetchOwnerReferences(
	c *dynamic.DynamicClient, gvr schema.GroupVersionResource, namespace, name string,
) (owners []metav1.OwnerReference, err error) {
	var resource *unstructured.Unstructured
	if resource, err = c.Resource(gvr).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{}); err != nil {
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
