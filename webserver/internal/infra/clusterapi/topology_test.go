package clusterapi

import (
	"testing"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestClusterTopology_AddNode(t *testing.T) {
	tp := NewClusterTopology()
	obj := ObjectInfo{
		Name: "node1",
		GVR: schema.GroupVersionResource{
			Group: "group1", Version: "v1", Resource: "res1",
		},
		PositionY: 1,
		PositionX: 1,
	}
	node := tp.AddNode(obj)

	if len(tp.Nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(tp.Nodes))
	}
	if tp.Nodes[0].Id != node.Id {
		t.Errorf("Expected node ID %s, got %s", node.Id, tp.Nodes[0].Id)
	}
}

func TestClusterTopology_AddEdge(t *testing.T) {
	tp := NewClusterTopology()
	node1 := tp.AddNode(ObjectInfo{
		Name:      "node1",
		GVR:       schema.GroupVersionResource{Group: "group1", Version: "v1", Resource: "res1"},
		PositionY: 1,
		PositionX: 1,
	})
	node2 := tp.AddNode(ObjectInfo{
		Name:      "node2",
		GVR:       schema.GroupVersionResource{Group: "group1", Version: "v1", Resource: "res1"},
		PositionY: 2,
		PositionX: 2,
	})

	tp.AddEdge(node1, node2)

	if len(tp.Edges) != 1 {
		t.Errorf("Expected 1 edge, got %d", len(tp.Edges))
	}
	if tp.Edges[0].Id != node1.Id+node2.Id {
		t.Errorf("Expected edge ID %s, got %s", node1.Id+node2.Id, tp.Edges[0].Id)
	}
}

func TestClusterTopology_Find(t *testing.T) {
	tp := NewClusterTopology()
	node := tp.AddNode(ObjectInfo{Name: "node1", GVR: schema.GroupVersionResource{Group: "group1", Version: "v1", Resource: "res1"}, Index: 1})

	found := tp.Find(&node)
	if !found {
		t.Errorf("Expected to find the node, but it was not found")
	}
}
