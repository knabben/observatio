package clusterapi

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

// ObjectInfo holds the information needed to process an object
type ObjectInfo struct {
	GVR       schema.GroupVersionResource
	Namespace string
	Name      string
	PositionX int32
	PositionY int32
}

// NodeData represents the data content of a Node
type NodeData struct {
	Label string `json:"label"`
}

// NodePosition represents the position of a Node in the graph
type NodePosition struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

type NodeStyle struct {
	BackgroundColor string `json:"backgroundColor"`
	Color           string `json:"color"`
}

// Node represents a graph node with an identifier,
// metadata, and its positional coordinates.
type Node struct {
	Id       string       `json:"id"`
	Data     NodeData     `json:"data"`
	Position NodePosition `json:"position"`
	Style    NodeStyle    `json:"style"`
}

var styles = map[int32]NodeStyle{
	0: NodeStyle{
		BackgroundColor: "#47556A",
		Color:           "#ffffff",
	},
	150: NodeStyle{
		BackgroundColor: "#DC8525",
		Color:           "#ffffff",
	},
	300: NodeStyle{
		BackgroundColor: "#72BA34",
		Color:           "#ffffff",
	},
	350: NodeStyle{
		BackgroundColor: "#018978",
		Color:           "#ffffff",
	},
	450: NodeStyle{
		BackgroundColor: "#0E779A",
		Color:           "#ffffff",
	},
}

// generateNodeID creates a unique identifier for a node
func generateNodeID(info ObjectInfo) string {
	return info.Name + info.GVR.String()
}

func generateStyle(info ObjectInfo) NodeStyle {
	if _, ok := styles[info.PositionY]; ok {
		return styles[info.PositionY]
	}
	return NodeStyle{}
}

// formatNodeLabel creates a formatted label for a node
func formatNodeLabel(info ObjectInfo) string {
	return fmt.Sprintf("%s\n%s", info.Name, info.GVR.Resource)
}

// NewNode creates a new Node instance with the given ObjectInfo
func NewNode(info ObjectInfo) Node {
	return Node{
		Id: generateNodeID(info),
		Data: NodeData{
			Label: formatNodeLabel(info),
		},
		Position: NodePosition{
			X: info.PositionX,
			Y: info.PositionY,
		},
		Style: generateStyle(info),
	}
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

// NewClusterTopology creates and returns a new initialized ClusterTopology instance
func NewClusterTopology() ClusterTopology {
	return ClusterTopology{
		Nodes: make([]Node, 0),
		Edges: make([]Edge, 0),
	}
}

// AddNode adds a new node to the ClusterTopology if it does not already exist and returns the created or existing node.
func (cl *ClusterTopology) AddNode(objectInfo ObjectInfo) (node Node) {
	node = NewNode(objectInfo)
	if !cl.Find(&node) {
		cl.Nodes = append(cl.Nodes, node)
	}
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
