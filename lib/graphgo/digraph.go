package graphgo

import "fmt"

// Edge represents a directed edge from Src to Dst.
type Edge struct {
	Src string
	Dst string
}

// DiGraph is a directed graph using adjacency lists.
type DiGraph struct {
	succ      map[string]map[string]bool // node -> set of successors
	pred      map[string]map[string]bool // node -> set of predecessors
	nodeAttrs map[string]map[string]any  // node -> attributes
	edgeAttrs map[string]map[string]any  // "src\x00dst" -> attributes
}

// NewDiGraph creates an empty directed graph.
func NewDiGraph() *DiGraph {
	return &DiGraph{
		succ:      make(map[string]map[string]bool),
		pred:      make(map[string]map[string]bool),
		nodeAttrs: make(map[string]map[string]any),
		edgeAttrs: make(map[string]map[string]any),
	}
}

func edgeKey(from, to string) string {
	return from + "\x00" + to
}

// AddNode adds a node to the graph. If the node already exists, this is a no-op.
func (g *DiGraph) AddNode(node string) {
	if _, ok := g.succ[node]; !ok {
		g.succ[node] = make(map[string]bool)
		g.pred[node] = make(map[string]bool)
		g.nodeAttrs[node] = make(map[string]any)
	}
}

// AddNodes adds multiple nodes.
func (g *DiGraph) AddNodes(nodes ...string) {
	for _, n := range nodes {
		g.AddNode(n)
	}
}

// RemoveNode removes a node and all its incident edges.
func (g *DiGraph) RemoveNode(node string) {
	if _, ok := g.succ[node]; !ok {
		return
	}
	// Remove edges from this node.
	for dst := range g.succ[node] {
		delete(g.pred[dst], node)
		delete(g.edgeAttrs, edgeKey(node, dst))
	}
	// Remove edges to this node.
	for src := range g.pred[node] {
		delete(g.succ[src], node)
		delete(g.edgeAttrs, edgeKey(src, node))
	}
	delete(g.succ, node)
	delete(g.pred, node)
	delete(g.nodeAttrs, node)
}

// AddEdge adds a directed edge from -> to. Nodes are created if they don't exist.
func (g *DiGraph) AddEdge(from, to string) {
	g.AddNode(from)
	g.AddNode(to)
	g.succ[from][to] = true
	g.pred[to][from] = true
	k := edgeKey(from, to)
	if _, ok := g.edgeAttrs[k]; !ok {
		g.edgeAttrs[k] = make(map[string]any)
	}
}

// RemoveEdge removes a directed edge. Returns an error if it doesn't exist.
func (g *DiGraph) RemoveEdge(from, to string) error {
	if !g.HasEdge(from, to) {
		return fmt.Errorf("graphgo: edge (%s, %s) not found", from, to)
	}
	delete(g.succ[from], to)
	delete(g.pred[to], from)
	delete(g.edgeAttrs, edgeKey(from, to))
	return nil
}

// HasNode returns true if the node exists.
func (g *DiGraph) HasNode(node string) bool {
	_, ok := g.succ[node]
	return ok
}

// HasEdge returns true if the directed edge exists.
func (g *DiGraph) HasEdge(from, to string) bool {
	if s, ok := g.succ[from]; ok {
		return s[to]
	}
	return false
}

// Nodes returns all nodes in the graph.
func (g *DiGraph) Nodes() []string {
	nodes := make([]string, 0, len(g.succ))
	for n := range g.succ {
		nodes = append(nodes, n)
	}
	return nodes
}

// Edges returns all directed edges in the graph.
func (g *DiGraph) Edges() []Edge {
	var edges []Edge
	for src, dsts := range g.succ {
		for dst := range dsts {
			edges = append(edges, Edge{Src: src, Dst: dst})
		}
	}
	return edges
}

// Predecessors returns the set of nodes with edges pointing to the given node.
func (g *DiGraph) Predecessors(node string) []string {
	preds := g.pred[node]
	out := make([]string, 0, len(preds))
	for p := range preds {
		out = append(out, p)
	}
	return out
}

// Successors returns the set of nodes that the given node has edges to.
func (g *DiGraph) Successors(node string) []string {
	succs := g.succ[node]
	out := make([]string, 0, len(succs))
	for s := range succs {
		out = append(out, s)
	}
	return out
}

// Parents is an alias for Predecessors (Bayesian network terminology).
func (g *DiGraph) Parents(node string) []string {
	return g.Predecessors(node)
}

// Children is an alias for Successors (Bayesian network terminology).
func (g *DiGraph) Children(node string) []string {
	return g.Successors(node)
}

// InDegree returns the number of edges pointing to the node.
func (g *DiGraph) InDegree(node string) int {
	return len(g.pred[node])
}

// OutDegree returns the number of edges from the node.
func (g *DiGraph) OutDegree(node string) int {
	return len(g.succ[node])
}

// NumberOfNodes returns the total number of nodes.
func (g *DiGraph) NumberOfNodes() int {
	return len(g.succ)
}

// NumberOfEdges returns the total number of directed edges.
func (g *DiGraph) NumberOfEdges() int {
	count := 0
	for _, dsts := range g.succ {
		count += len(dsts)
	}
	return count
}

// NodeAttr returns the attribute map for a node. Returns nil if the node doesn't exist.
func (g *DiGraph) NodeAttr(node string) map[string]any {
	return g.nodeAttrs[node]
}

// EdgeAttr returns the attribute map for an edge. Returns nil if the edge doesn't exist.
func (g *DiGraph) EdgeAttr(from, to string) map[string]any {
	return g.edgeAttrs[edgeKey(from, to)]
}

// Copy returns a deep copy of the graph.
func (g *DiGraph) Copy() *DiGraph {
	c := NewDiGraph()
	for n := range g.succ {
		c.AddNode(n)
		for k, v := range g.nodeAttrs[n] {
			c.nodeAttrs[n][k] = v
		}
	}
	for src, dsts := range g.succ {
		for dst := range dsts {
			c.AddEdge(src, dst)
			for k, v := range g.edgeAttrs[edgeKey(src, dst)] {
				c.edgeAttrs[edgeKey(src, dst)][k] = v
			}
		}
	}
	return c
}

// Subgraph returns a new DiGraph containing only the specified nodes and edges between them.
func (g *DiGraph) Subgraph(nodes []string) *DiGraph {
	sub := NewDiGraph()
	nodeSet := make(map[string]bool, len(nodes))
	for _, n := range nodes {
		if g.HasNode(n) {
			nodeSet[n] = true
			sub.AddNode(n)
			for k, v := range g.nodeAttrs[n] {
				sub.nodeAttrs[n][k] = v
			}
		}
	}
	for src := range nodeSet {
		for dst := range g.succ[src] {
			if nodeSet[dst] {
				sub.AddEdge(src, dst)
				for k, v := range g.edgeAttrs[edgeKey(src, dst)] {
					sub.edgeAttrs[edgeKey(src, dst)][k] = v
				}
			}
		}
	}
	return sub
}
