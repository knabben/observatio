package clusterapi

import (
	"context"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type {

}

// fetchOwnerHierarchy traverse with bfs and create edges between nodes
func fetchOwnerHierarchy(ctx context.Context, owners []metav1.OwnerReference, gvr schema.GroupVersionResource, namespace, name string) (err error) {
	fmt.Println(name, namespace, gvr)
	cc, _ := NewDynamicClient(ctx)

	for _, owner := range owners {
		ownerGVR := convertGVK(owner)
		fmt.Println(ownerGVR, owner.Name, namespace)

		parentOwners, err := FetchOwnerReferences(cc, ownerGVR, namespace, owner.Name)
		if err != nil {
			return err
		}

		if len(parentOwners) > 0 {
			if err := fetchOwnerHierarchy(ctx, parentOwners, ownerGVR, namespace, owner.Name); err != nil {
				return err
			}
		}
	}
	return nil
}

// FetchOwnerReferences retrieves the owner references of a Kubernetes resource identified by GroupVersionResource, namespace, and name.
// It returns a slice of OwnerReference objects or an error if the retrieval or conversion fails.
func FetchOwnerReferences(c *dynamic.DynamicClient, gvr schema.GroupVersionResource, namespace, name string) (owners []metav1.OwnerReference, err error) {
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
