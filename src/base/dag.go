package base

import (
	"fmt"
	"sort"

	"github.com/asymmetric-effort/pgmgo/lib/graphgo"
)

// DAG is a directed acyclic graph. It wraps a graphgo.DiGraph and enforces
// acyclicity on every edge addition.
type DAG struct {
	g *graphgo.DiGraph
}

// NewDAG creates a new empty DAG.
func NewDAG() *DAG {
	return &DAG{g: graphgo.NewDiGraph()}
}

// AddNode adds a node to the DAG. Returns an error if the node already exists.
func (d *DAG) AddNode(node string) error {
	if d.g.HasNode(node) {
		return fmt.Errorf("base: node %q already exists", node)
	}
	d.g.AddNode(node)
	return nil
}

// AddNodes adds multiple nodes. Returns an error if any node already exists;
// nodes added before the error are retained.
func (d *DAG) AddNodes(nodes ...string) error {
	for _, n := range nodes {
		if err := d.AddNode(n); err != nil {
			return err
		}
	}
	return nil
}

// RemoveNode removes a node and all its incident edges.
// Returns an error if the node does not exist.
func (d *DAG) RemoveNode(node string) error {
	if !d.g.HasNode(node) {
		return fmt.Errorf("base: node %q not found", node)
	}
	d.g.RemoveNode(node)
	return nil
}

// AddEdge adds a directed edge from -> to. Both nodes must already exist.
// Returns an error if the edge would create a cycle, if either node does not
// exist, or if the edge already exists.
func (d *DAG) AddEdge(from, to string) error {
	if !d.g.HasNode(from) {
		return fmt.Errorf("base: node %q not found", from)
	}
	if !d.g.HasNode(to) {
		return fmt.Errorf("base: node %q not found", to)
	}
	if d.g.HasEdge(from, to) {
		return fmt.Errorf("base: edge (%q, %q) already exists", from, to)
	}

	// Temporarily add the edge and check acyclicity.
	d.g.AddEdge(from, to)
	if !graphgo.IsDAG(d.g) {
		// Revert: remove the edge we just added.
		_ = d.g.RemoveEdge(from, to)
		return fmt.Errorf("base: edge (%q, %q) would create a cycle", from, to)
	}
	return nil
}

// RemoveEdge removes a directed edge. Returns an error if it does not exist.
func (d *DAG) RemoveEdge(from, to string) error {
	return d.g.RemoveEdge(from, to)
}

// HasNode returns true if the node exists in the DAG.
func (d *DAG) HasNode(node string) bool {
	return d.g.HasNode(node)
}

// HasEdge returns true if the directed edge exists.
func (d *DAG) HasEdge(from, to string) bool {
	return d.g.HasEdge(from, to)
}

// Nodes returns a sorted list of all nodes.
func (d *DAG) Nodes() []string {
	nodes := d.g.Nodes()
	sort.Strings(nodes)
	return nodes
}

// Edges returns all directed edges, sorted lexicographically.
func (d *DAG) Edges() []graphgo.Edge {
	edges := d.g.Edges()
	sort.Slice(edges, func(i, j int) bool {
		if edges[i].Src != edges[j].Src {
			return edges[i].Src < edges[j].Src
		}
		return edges[i].Dst < edges[j].Dst
	})
	return edges
}

// Parents returns the sorted parents (predecessors) of a node.
func (d *DAG) Parents(node string) []string {
	p := d.g.Parents(node)
	sort.Strings(p)
	return p
}

// Children returns the sorted children (successors) of a node.
func (d *DAG) Children(node string) []string {
	c := d.g.Children(node)
	sort.Strings(c)
	return c
}

// GetRoots returns nodes with no parents, sorted.
func (d *DAG) GetRoots() []string {
	var roots []string
	for _, n := range d.Nodes() {
		if d.g.InDegree(n) == 0 {
			roots = append(roots, n)
		}
	}
	return roots
}

// GetLeaves returns nodes with no children, sorted.
func (d *DAG) GetLeaves() []string {
	var leaves []string
	for _, n := range d.Nodes() {
		if d.g.OutDegree(n) == 0 {
			leaves = append(leaves, n)
		}
	}
	return leaves
}

// TopologicalOrder returns a topological ordering of the DAG.
func (d *DAG) TopologicalOrder() ([]string, error) {
	return graphgo.TopologicalSort(d.g)
}

// Copy returns a deep copy of the DAG.
func (d *DAG) Copy() *DAG {
	return &DAG{g: d.g.Copy()}
}
