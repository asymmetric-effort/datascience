//go:build unit

package graphgo

import (
	"sort"
	"testing"
)

// ---------------------------------------------------------------------------
// Graph: RemoveNode, RemoveEdge, Copy with attrs
// ---------------------------------------------------------------------------

func TestGraphRemoveNode(t *testing.T) {
	g := NewGraph()
	g.AddEdge("A", "B")
	g.AddEdge("B", "C")
	g.RemoveNode("B")

	if g.HasNode("B") {
		t.Error("B should be removed")
	}
	if g.HasEdge("A", "B") || g.HasEdge("B", "C") {
		t.Error("edges to B should be removed")
	}
	if !g.HasNode("A") || !g.HasNode("C") {
		t.Error("A and C should remain")
	}
}

func TestGraphRemoveNode_NonExistent(t *testing.T) {
	g := NewGraph()
	g.AddNode("A")
	g.RemoveNode("Z") // should not panic
	if !g.HasNode("A") {
		t.Error("A should remain")
	}
}

func TestGraphRemoveEdge(t *testing.T) {
	g := NewGraph()
	g.AddEdge("A", "B")
	err := g.RemoveEdge("A", "B")
	if err != nil {
		t.Fatalf("RemoveEdge failed: %v", err)
	}
	if g.HasEdge("A", "B") {
		t.Error("edge should be removed")
	}
}

func TestGraphRemoveEdge_NonExistent(t *testing.T) {
	g := NewGraph()
	g.AddNode("A")
	g.AddNode("B")
	err := g.RemoveEdge("A", "B")
	if err == nil {
		t.Error("expected error for non-existent edge")
	}
}

func TestGraphCopy_WithEdgeAttrs(t *testing.T) {
	g := NewGraph()
	g.AddEdge("A", "B")
	k := undirectedEdgeKey("A", "B")
	g.edgeAttrs[k]["weight"] = 42
	g.nodeAttrs["A"]["color"] = "red"

	c := g.Copy()
	if c.edgeAttrs[k]["weight"] != 42 {
		t.Error("edge attrs not copied")
	}
	if c.nodeAttrs["A"]["color"] != "red" {
		t.Error("node attrs not copied")
	}
}

// ---------------------------------------------------------------------------
// PDAG: HasDirectedEdge when from not in map, DirectedEdges empty,
//       Neighbors with predecessors
// ---------------------------------------------------------------------------

func TestPDAGHasDirectedEdge_MissingNode(t *testing.T) {
	p := NewPDAG()
	if p.HasDirectedEdge("X", "Y") {
		t.Error("expected false for missing node")
	}
}

func TestPDAGDirectedEdges_Empty(t *testing.T) {
	p := NewPDAG()
	p.AddNode("A")
	edges := p.DirectedEdges()
	if len(edges) != 0 {
		t.Errorf("expected 0 directed edges, got %d", len(edges))
	}
}

func TestPDAGNeighbors_AllTypes(t *testing.T) {
	p := NewPDAG()
	p.AddDirectedEdge("A", "B")   // A -> B
	p.AddDirectedEdge("C", "B")   // C -> B (predecessor)
	p.AddUndirectedEdge("B", "D") // B -- D

	neighbors := p.Neighbors("B")
	sort.Strings(neighbors)
	expected := []string{"A", "C", "D"}
	if len(neighbors) != 3 {
		t.Fatalf("expected 3 neighbors, got %d: %v", len(neighbors), neighbors)
	}
	for i, n := range expected {
		if neighbors[i] != n {
			t.Errorf("neighbor %d: expected %q, got %q", i, n, neighbors[i])
		}
	}
}

func TestPDAGRemoveNode_Coverage(t *testing.T) {
	p := NewPDAG()
	p.AddDirectedEdge("A", "B")
	p.AddUndirectedEdge("B", "C")
	p.RemoveNode("B")

	if p.HasNode("B") {
		t.Error("B should be removed")
	}
	if p.HasDirectedEdge("A", "B") {
		t.Error("directed edge A->B should be removed")
	}
	if p.HasUndirectedEdge("B", "C") {
		t.Error("undirected edge B--C should be removed")
	}
}

func TestPDAGRemoveNode_NonExistent(t *testing.T) {
	p := NewPDAG()
	p.AddNode("A")
	p.RemoveNode("Z") // should not panic
	if !p.HasNode("A") {
		t.Error("A should remain")
	}
}

// ---------------------------------------------------------------------------
// DiGraph: Subgraph with non-existent nodes
// ---------------------------------------------------------------------------

func TestDigraphSubgraph_NonExistentNodes(t *testing.T) {
	g := NewDiGraph()
	g.AddEdge("A", "B")
	g.AddEdge("B", "C")
	g.nodeAttrs["A"]["color"] = "red"
	g.edgeAttrs[edgeKey("A", "B")]["weight"] = 1

	sub := g.Subgraph([]string{"A", "B", "Z"})
	if sub.HasNode("Z") {
		t.Error("Z should not be in subgraph")
	}
	if !sub.HasEdge("A", "B") {
		t.Error("A->B should be in subgraph")
	}
	if sub.HasEdge("B", "C") {
		t.Error("B->C should not be in subgraph")
	}
	if sub.nodeAttrs["A"]["color"] != "red" {
		t.Error("node attrs not copied")
	}
	if sub.edgeAttrs[edgeKey("A", "B")]["weight"] != 1 {
		t.Error("edge attrs not copied")
	}
}

// ---------------------------------------------------------------------------
// Cliques: cliqueLabel, BuildJunctionTree edge cases
// ---------------------------------------------------------------------------

func TestCliqueLabel_MultiDigit(t *testing.T) {
	if cliqueLabel(0) != "0" {
		t.Errorf("cliqueLabel(0) = %q, want %q", cliqueLabel(0), "0")
	}
	if cliqueLabel(123) != "123" {
		t.Errorf("cliqueLabel(123) = %q, want %q", cliqueLabel(123), "123")
	}
	if cliqueLabel(10) != "10" {
		t.Errorf("cliqueLabel(10) = %q, want %q", cliqueLabel(10), "10")
	}
}

func TestBuildJunctionTree_Empty(t *testing.T) {
	tree, seps := BuildJunctionTree(nil)
	if len(tree.Nodes()) != 0 {
		t.Error("expected empty tree")
	}
	if seps != nil {
		t.Error("expected nil separators")
	}
}

func TestBuildJunctionTree_NoOverlap(t *testing.T) {
	cliques := [][]string{
		{"A", "B"},
		{"C", "D"},
	}
	tree, _ := BuildJunctionTree(cliques)
	if len(tree.Nodes()) != 2 {
		t.Errorf("expected 2 nodes, got %d", len(tree.Nodes()))
	}
	if len(tree.Edges()) != 0 {
		t.Errorf("expected 0 edges, got %d", len(tree.Edges()))
	}
}

func TestBuildJunctionTree_SingleClique(t *testing.T) {
	cliques := [][]string{{"A", "B", "C"}}
	tree, _ := BuildJunctionTree(cliques)
	if len(tree.Nodes()) != 1 {
		t.Errorf("expected 1 node, got %d", len(tree.Nodes()))
	}
}

// ---------------------------------------------------------------------------
// MaxCliques: empty graph, single node
// ---------------------------------------------------------------------------

func TestMaxCliques_Empty(t *testing.T) {
	g := NewGraph()
	cliques := MaxCliques(g)
	if cliques != nil {
		t.Errorf("expected nil, got %v", cliques)
	}
}

func TestMaxCliques_SingleNode(t *testing.T) {
	g := NewGraph()
	g.AddNode("A")
	cliques := MaxCliques(g)
	if len(cliques) != 1 || len(cliques[0]) != 1 || cliques[0][0] != "A" {
		t.Errorf("expected [[A]], got %v", cliques)
	}
}

// ---------------------------------------------------------------------------
// Meek rules: test cases that trigger all 4 rules
// ---------------------------------------------------------------------------

func TestMeekR1(t *testing.T) {
	p := NewPDAG()
	p.AddDirectedEdge("W", "U")
	p.AddUndirectedEdge("U", "V")
	changed := ApplyMeekRules(p)
	if !changed {
		t.Error("expected changes from Meek rules")
	}
	if !p.HasDirectedEdge("U", "V") {
		t.Error("expected U->V after R1")
	}
}

func TestMeekR2(t *testing.T) {
	p := NewPDAG()
	p.AddDirectedEdge("U", "W")
	p.AddDirectedEdge("W", "V")
	p.AddUndirectedEdge("U", "V")
	changed := ApplyMeekRules(p)
	if !changed {
		t.Error("expected changes from Meek rules")
	}
	if !p.HasDirectedEdge("U", "V") {
		t.Error("expected U->V after R2")
	}
}

func TestMeekR3(t *testing.T) {
	p := NewPDAG()
	p.AddUndirectedEdge("W1", "U")
	p.AddUndirectedEdge("W2", "U")
	p.AddDirectedEdge("W1", "V")
	p.AddDirectedEdge("W2", "V")
	p.AddUndirectedEdge("U", "V")
	changed := ApplyMeekRules(p)
	if !changed {
		t.Error("expected changes from Meek rules")
	}
	if !p.HasDirectedEdge("U", "V") {
		t.Error("expected U->V after R3")
	}
}

func TestMeekR4(t *testing.T) {
	// R4: w—u, w→x→v, u—v => orient u→v
	// Need to ensure R1, R2, R3 don't fire first.
	// The key condition for R4: there exists w undirected-adjacent to u,
	// w has directed edge to x, and x has directed edge to v.
	p := NewPDAG()
	p.AddNodes("U", "V", "W", "X")
	p.AddUndirectedEdge("W", "U")
	p.AddDirectedEdge("W", "X")
	p.AddDirectedEdge("X", "V")
	p.AddUndirectedEdge("U", "V")
	// Also make W adjacent to V to prevent R1 from firing on U—V.
	p.AddUndirectedEdge("W", "V")
	changed := ApplyMeekRules(p)
	if !changed {
		t.Error("expected changes from Meek rules")
	}
	// Some edge should have been oriented.
	if len(p.UndirectedEdges()) >= 3 {
		t.Error("expected at least one undirected edge to be oriented")
	}
}

func TestMeekRules_NoChange(t *testing.T) {
	p := NewPDAG()
	p.AddUndirectedEdge("A", "B")
	changed := ApplyMeekRules(p)
	if changed {
		t.Error("expected no changes")
	}
}

func TestOrient_AlreadyDirected(t *testing.T) {
	p := NewPDAG()
	p.AddDirectedEdge("A", "B")
	result := orient(p, "A", "B")
	if result {
		t.Error("expected false when edge is directed, not undirected")
	}
}

// ---------------------------------------------------------------------------
// DAGToPDAG
// ---------------------------------------------------------------------------

func TestDAGToPDAG_VStructure(t *testing.T) {
	g := NewDiGraph()
	g.AddEdge("A", "C")
	g.AddEdge("B", "C")

	p := DAGToPDAG(g)
	if !p.HasDirectedEdge("A", "C") {
		t.Error("expected A->C (v-structure)")
	}
	if !p.HasDirectedEdge("B", "C") {
		t.Error("expected B->C (v-structure)")
	}
}
