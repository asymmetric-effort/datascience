package base

import (
	"fmt"
	"sort"

	"github.com/asymmetric-effort/datascience/lib/graphgo"
)

// UndirectedGraph wraps a graphgo.Graph and provides error-checked operations.
type UndirectedGraph struct {
	g *graphgo.Graph
}

// NewUndirectedGraph creates a new empty undirected graph.
func NewUndirectedGraph() *UndirectedGraph {
	return &UndirectedGraph{g: graphgo.NewGraph()}
}

// AddNode adds a node. Returns an error if the node already exists.
func (u *UndirectedGraph) AddNode(node string) error {
	if u.g.HasNode(node) {
		return fmt.Errorf("base: node %q already exists", node)
	}
	u.g.AddNode(node)
	return nil
}

// AddEdge adds an undirected edge between u and v. Both nodes must already
// exist. Returns an error if either node is missing or the edge already exists.
func (u *UndirectedGraph) AddEdge(a, b string) error {
	if !u.g.HasNode(a) {
		return fmt.Errorf("base: node %q not found", a)
	}
	if !u.g.HasNode(b) {
		return fmt.Errorf("base: node %q not found", b)
	}
	if u.g.HasEdge(a, b) {
		return fmt.Errorf("base: edge (%q, %q) already exists", a, b)
	}
	u.g.AddEdge(a, b)
	return nil
}

// RemoveNode removes a node and all its incident edges.
// Returns an error if the node does not exist.
func (u *UndirectedGraph) RemoveNode(node string) error {
	if !u.g.HasNode(node) {
		return fmt.Errorf("base: node %q not found", node)
	}
	u.g.RemoveNode(node)
	return nil
}

// RemoveEdge removes an undirected edge. Returns an error if it does not exist.
func (u *UndirectedGraph) RemoveEdge(a, b string) error {
	return u.g.RemoveEdge(a, b)
}

// HasNode returns true if the node exists.
func (u *UndirectedGraph) HasNode(node string) bool {
	return u.g.HasNode(node)
}

// HasEdge returns true if the undirected edge exists.
func (u *UndirectedGraph) HasEdge(a, b string) bool {
	return u.g.HasEdge(a, b)
}

// Nodes returns a sorted list of all nodes.
func (u *UndirectedGraph) Nodes() []string {
	nodes := u.g.Nodes()
	sort.Strings(nodes)
	return nodes
}

// Edges returns all undirected edges in canonical form (A < B), sorted lexicographically.
func (u *UndirectedGraph) Edges() []graphgo.UndirectedEdge {
	edges := u.g.Edges()
	// Canonicalize: ensure A < B for each edge.
	for i := range edges {
		if edges[i].A > edges[i].B {
			edges[i].A, edges[i].B = edges[i].B, edges[i].A
		}
	}
	sort.Slice(edges, func(i, j int) bool {
		if edges[i].A != edges[j].A {
			return edges[i].A < edges[j].A
		}
		return edges[i].B < edges[j].B
	})
	return edges
}

// Neighbors returns the sorted neighbors of a node.
func (u *UndirectedGraph) Neighbors(node string) []string {
	n := u.g.Neighbors(node)
	sort.Strings(n)
	return n
}

// Degree returns the number of edges incident to the node.
func (u *UndirectedGraph) Degree(node string) int {
	return u.g.Degree(node)
}

// Copy returns a deep copy of the undirected graph.
func (u *UndirectedGraph) Copy() *UndirectedGraph {
	return &UndirectedGraph{g: u.g.Copy()}
}
