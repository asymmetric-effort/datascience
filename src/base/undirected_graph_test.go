//go:build unit

package base

import (
	"testing"
)

func TestNewUndirectedGraph(t *testing.T) {
	g := NewUndirectedGraph()
	if g == nil {
		t.Fatal("NewUndirectedGraph returned nil")
	}
	if len(g.Nodes()) != 0 {
		t.Errorf("new graph should have 0 nodes, got %d", len(g.Nodes()))
	}
	if len(g.Edges()) != 0 {
		t.Errorf("new graph should have 0 edges, got %d", len(g.Edges()))
	}
}

func TestUndirectedAddNode(t *testing.T) {
	g := NewUndirectedGraph()
	if err := g.AddNode("A"); err != nil {
		t.Fatalf("AddNode failed: %v", err)
	}
	if !g.HasNode("A") {
		t.Error("HasNode returned false for added node")
	}
}

func TestUndirectedAddNodeDuplicate(t *testing.T) {
	g := NewUndirectedGraph()
	_ = g.AddNode("A")
	err := g.AddNode("A")
	if err == nil {
		t.Error("expected error when adding duplicate node")
	}
}

func TestUndirectedRemoveNode(t *testing.T) {
	g := NewUndirectedGraph()
	_ = g.AddNode("A")
	_ = g.AddNode("B")
	_ = g.AddEdge("A", "B")

	if err := g.RemoveNode("A"); err != nil {
		t.Fatalf("RemoveNode failed: %v", err)
	}
	if g.HasNode("A") {
		t.Error("removed node should not exist")
	}
	if g.HasEdge("A", "B") {
		t.Error("edges incident to removed node should be gone")
	}
	if g.HasEdge("B", "A") {
		t.Error("reverse edge incident to removed node should be gone")
	}
}

func TestUndirectedRemoveNodeNotFound(t *testing.T) {
	g := NewUndirectedGraph()
	err := g.RemoveNode("X")
	if err == nil {
		t.Error("expected error when removing non-existent node")
	}
}

func TestUndirectedAddEdge(t *testing.T) {
	g := NewUndirectedGraph()
	_ = g.AddNode("A")
	_ = g.AddNode("B")
	if err := g.AddEdge("A", "B"); err != nil {
		t.Fatalf("AddEdge failed: %v", err)
	}
	if !g.HasEdge("A", "B") {
		t.Error("HasEdge(A,B) returned false for added edge")
	}
	if !g.HasEdge("B", "A") {
		t.Error("HasEdge(B,A) returned false; undirected edge should be symmetric")
	}
}

func TestUndirectedAddEdgeMissingNode(t *testing.T) {
	g := NewUndirectedGraph()
	_ = g.AddNode("A")
	err := g.AddEdge("A", "B")
	if err == nil {
		t.Error("expected error when node B does not exist")
	}
	err = g.AddEdge("X", "A")
	if err == nil {
		t.Error("expected error when node X does not exist")
	}
}

func TestUndirectedAddEdgeDuplicate(t *testing.T) {
	g := NewUndirectedGraph()
	_ = g.AddNode("A")
	_ = g.AddNode("B")
	_ = g.AddEdge("A", "B")
	err := g.AddEdge("A", "B")
	if err == nil {
		t.Error("expected error when adding duplicate edge")
	}
	// Also test reverse direction.
	err = g.AddEdge("B", "A")
	if err == nil {
		t.Error("expected error when adding duplicate edge in reverse direction")
	}
}

func TestUndirectedRemoveEdge(t *testing.T) {
	g := NewUndirectedGraph()
	_ = g.AddNode("A")
	_ = g.AddNode("B")
	_ = g.AddEdge("A", "B")

	if err := g.RemoveEdge("A", "B"); err != nil {
		t.Fatalf("RemoveEdge failed: %v", err)
	}
	if g.HasEdge("A", "B") {
		t.Error("edge should be removed")
	}
	if g.HasEdge("B", "A") {
		t.Error("reverse edge should also be removed")
	}
}

func TestUndirectedRemoveEdgeNotFound(t *testing.T) {
	g := NewUndirectedGraph()
	_ = g.AddNode("A")
	_ = g.AddNode("B")
	err := g.RemoveEdge("A", "B")
	if err == nil {
		t.Error("expected error when removing non-existent edge")
	}
}

func TestUndirectedNodesSorted(t *testing.T) {
	g := NewUndirectedGraph()
	_ = g.AddNode("C")
	_ = g.AddNode("A")
	_ = g.AddNode("B")
	nodes := g.Nodes()
	if len(nodes) != 3 || nodes[0] != "A" || nodes[1] != "B" || nodes[2] != "C" {
		t.Errorf("Nodes() should be sorted, got %v", nodes)
	}
}

func TestUndirectedEdgesSorted(t *testing.T) {
	g := NewUndirectedGraph()
	_ = g.AddNode("A")
	_ = g.AddNode("B")
	_ = g.AddNode("C")
	_ = g.AddEdge("B", "C")
	_ = g.AddEdge("A", "C")
	_ = g.AddEdge("A", "B")

	edges := g.Edges()
	if len(edges) != 3 {
		t.Fatalf("expected 3 edges, got %d", len(edges))
	}
	// Sorted: (A,B), (A,C), (B,C)
	if edges[0].A != "A" || edges[0].B != "B" {
		t.Errorf("edges[0] = (%s,%s), want (A,B)", edges[0].A, edges[0].B)
	}
	if edges[1].A != "A" || edges[1].B != "C" {
		t.Errorf("edges[1] = (%s,%s), want (A,C)", edges[1].A, edges[1].B)
	}
	if edges[2].A != "B" || edges[2].B != "C" {
		t.Errorf("edges[2] = (%s,%s), want (B,C)", edges[2].A, edges[2].B)
	}
}

func TestUndirectedNeighbors(t *testing.T) {
	g := NewUndirectedGraph()
	_ = g.AddNode("A")
	_ = g.AddNode("B")
	_ = g.AddNode("C")
	_ = g.AddNode("D")
	_ = g.AddEdge("B", "A")
	_ = g.AddEdge("B", "C")
	_ = g.AddEdge("B", "D")

	neighbors := g.Neighbors("B")
	if len(neighbors) != 3 || neighbors[0] != "A" || neighbors[1] != "C" || neighbors[2] != "D" {
		t.Errorf("Neighbors(B) = %v, want [A C D]", neighbors)
	}
}

func TestUndirectedDegree(t *testing.T) {
	g := NewUndirectedGraph()
	_ = g.AddNode("A")
	_ = g.AddNode("B")
	_ = g.AddNode("C")
	_ = g.AddEdge("A", "B")
	_ = g.AddEdge("A", "C")

	if g.Degree("A") != 2 {
		t.Errorf("Degree(A) = %d, want 2", g.Degree("A"))
	}
	if g.Degree("B") != 1 {
		t.Errorf("Degree(B) = %d, want 1", g.Degree("B"))
	}
}

func TestUndirectedCopy(t *testing.T) {
	g := NewUndirectedGraph()
	_ = g.AddNode("A")
	_ = g.AddNode("B")
	_ = g.AddNode("C")
	_ = g.AddEdge("A", "B")
	_ = g.AddEdge("B", "C")

	c := g.Copy()

	if len(c.Nodes()) != 3 {
		t.Errorf("copy should have 3 nodes, got %d", len(c.Nodes()))
	}
	if !c.HasEdge("A", "B") || !c.HasEdge("B", "C") {
		t.Error("copy should have the same edges")
	}

	// Modifying the copy should not affect the original.
	_ = c.AddNode("D")
	_ = c.AddEdge("C", "D")
	if g.HasNode("D") {
		t.Error("original should not have node D after modifying copy")
	}
	if g.HasEdge("C", "D") {
		t.Error("original should not have edge C-D after modifying copy")
	}
}

func TestUndirectedGraphComplex(t *testing.T) {
	// Build a triangle: A-B-C-A with an extra node D connected to B.
	g := NewUndirectedGraph()
	for _, n := range []string{"A", "B", "C", "D"} {
		if err := g.AddNode(n); err != nil {
			t.Fatalf("AddNode(%s) failed: %v", n, err)
		}
	}
	for _, e := range [][2]string{{"A", "B"}, {"B", "C"}, {"C", "A"}, {"B", "D"}} {
		if err := g.AddEdge(e[0], e[1]); err != nil {
			t.Fatalf("AddEdge(%s, %s) failed: %v", e[0], e[1], err)
		}
	}

	if len(g.Nodes()) != 4 {
		t.Errorf("expected 4 nodes, got %d", len(g.Nodes()))
	}
	if len(g.Edges()) != 4 {
		t.Errorf("expected 4 edges, got %d", len(g.Edges()))
	}
	if g.Degree("B") != 3 {
		t.Errorf("Degree(B) = %d, want 3", g.Degree("B"))
	}
	if g.Degree("D") != 1 {
		t.Errorf("Degree(D) = %d, want 1", g.Degree("D"))
	}

	// Remove node C; edges A-C and B-C should disappear.
	if err := g.RemoveNode("C"); err != nil {
		t.Fatalf("RemoveNode(C) failed: %v", err)
	}
	if g.HasNode("C") {
		t.Error("C should be removed")
	}
	if g.HasEdge("A", "C") || g.HasEdge("C", "A") {
		t.Error("edges to/from C should be removed")
	}
	if len(g.Edges()) != 2 {
		t.Errorf("expected 2 edges after removing C, got %d", len(g.Edges()))
	}
}
