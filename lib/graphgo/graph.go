package graphgo

import "fmt"

// UndirectedEdge represents an undirected edge.
type UndirectedEdge struct {
	A string
	B string
}

// Graph is an undirected graph using adjacency lists.
type Graph struct {
	adj       map[string]map[string]bool // node -> set of neighbors
	nodeAttrs map[string]map[string]any
	edgeAttrs map[string]map[string]any // canonical key -> attributes
}

// NewGraph creates an empty undirected graph.
func NewGraph() *Graph {
	return &Graph{
		adj:       make(map[string]map[string]bool),
		nodeAttrs: make(map[string]map[string]any),
		edgeAttrs: make(map[string]map[string]any),
	}
}

// undirectedEdgeKey returns a canonical key for an undirected edge.
func undirectedEdgeKey(a, b string) string {
	if a < b {
		return a + "\x00" + b
	}
	return b + "\x00" + a
}

// AddNode adds a node to the graph.
func (g *Graph) AddNode(node string) {
	if _, ok := g.adj[node]; !ok {
		g.adj[node] = make(map[string]bool)
		g.nodeAttrs[node] = make(map[string]any)
	}
}

// AddEdge adds an undirected edge. Nodes are created if they don't exist.
func (g *Graph) AddEdge(a, b string) {
	g.AddNode(a)
	g.AddNode(b)
	g.adj[a][b] = true
	g.adj[b][a] = true
	k := undirectedEdgeKey(a, b)
	if _, ok := g.edgeAttrs[k]; !ok {
		g.edgeAttrs[k] = make(map[string]any)
	}
}

// RemoveNode removes a node and all its incident edges.
func (g *Graph) RemoveNode(node string) {
	if _, ok := g.adj[node]; !ok {
		return
	}
	for neighbor := range g.adj[node] {
		delete(g.adj[neighbor], node)
		delete(g.edgeAttrs, undirectedEdgeKey(node, neighbor))
	}
	delete(g.adj, node)
	delete(g.nodeAttrs, node)
}

// RemoveEdge removes an undirected edge. Returns an error if it doesn't exist.
func (g *Graph) RemoveEdge(a, b string) error {
	if !g.HasEdge(a, b) {
		return fmt.Errorf("graphgo: edge (%s, %s) not found", a, b)
	}
	delete(g.adj[a], b)
	delete(g.adj[b], a)
	delete(g.edgeAttrs, undirectedEdgeKey(a, b))
	return nil
}

// HasNode returns true if the node exists.
func (g *Graph) HasNode(node string) bool {
	_, ok := g.adj[node]
	return ok
}

// HasEdge returns true if the undirected edge exists.
func (g *Graph) HasEdge(a, b string) bool {
	if neighbors, ok := g.adj[a]; ok {
		return neighbors[b]
	}
	return false
}

// Nodes returns all nodes.
func (g *Graph) Nodes() []string {
	nodes := make([]string, 0, len(g.adj))
	for n := range g.adj {
		nodes = append(nodes, n)
	}
	return nodes
}

// Edges returns all undirected edges (each edge appears once).
func (g *Graph) Edges() []UndirectedEdge {
	seen := make(map[string]bool)
	var edges []UndirectedEdge
	for a, neighbors := range g.adj {
		for b := range neighbors {
			k := undirectedEdgeKey(a, b)
			if !seen[k] {
				seen[k] = true
				edges = append(edges, UndirectedEdge{A: a, B: b})
			}
		}
	}
	return edges
}

// Neighbors returns all neighbors of the given node.
func (g *Graph) Neighbors(node string) []string {
	neighbors := g.adj[node]
	out := make([]string, 0, len(neighbors))
	for n := range neighbors {
		out = append(out, n)
	}
	return out
}

// Degree returns the number of edges incident to the node.
func (g *Graph) Degree(node string) int {
	return len(g.adj[node])
}

// Copy returns a deep copy of the graph.
func (g *Graph) Copy() *Graph {
	c := NewGraph()
	for n := range g.adj {
		c.AddNode(n)
		for k, v := range g.nodeAttrs[n] {
			c.nodeAttrs[n][k] = v
		}
	}
	for a, neighbors := range g.adj {
		for b := range neighbors {
			if a < b { // add each edge once
				c.AddEdge(a, b)
				k := undirectedEdgeKey(a, b)
				for ak, av := range g.edgeAttrs[k] {
					c.edgeAttrs[k][ak] = av
				}
			}
		}
	}
	return c
}
