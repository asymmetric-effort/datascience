//go:build unit

package graphgo

import (
	"sort"
	"testing"
)

func TestNewGraph(t *testing.T) {
	g := NewGraph()
	if len(g.Nodes()) != 0 {
		t.Fatal("expected 0 nodes")
	}
	if len(g.Edges()) != 0 {
		t.Fatal("expected 0 edges")
	}
}

func TestGraphAddNode(t *testing.T) {
	g := NewGraph()
	g.AddNode("A")
	if !g.HasNode("A") {
		t.Fatal("expected node A")
	}
	// Idempotent.
	g.AddNode("A")
	if len(g.Nodes()) != 1 {
		t.Fatal("duplicate add should not create extra node")
	}
}

func TestGraphAddEdge(t *testing.T) {
	g := NewGraph()
	g.AddEdge("A", "B")
	if !g.HasEdge("A", "B") {
		t.Fatal("expected edge A-B")
	}
	if !g.HasEdge("B", "A") {
		t.Fatal("undirected: B-A should also exist")
	}
	if !g.HasNode("A") || !g.HasNode("B") {
		t.Fatal("nodes should be auto-created")
	}
}

func TestGraphEdgesNoDuplicates(t *testing.T) {
	g := NewGraph()
	g.AddEdge("A", "B")
	g.AddEdge("B", "C")
	edges := g.Edges()
	if len(edges) != 2 {
		t.Fatalf("expected 2 edges, got %d", len(edges))
	}
}

func TestGraphNeighbors(t *testing.T) {
	g := NewGraph()
	g.AddEdge("A", "B")
	g.AddEdge("A", "C")
	g.AddEdge("B", "D")

	neighbors := g.Neighbors("A")
	sort.Strings(neighbors)
	if len(neighbors) != 2 || neighbors[0] != "B" || neighbors[1] != "C" {
		t.Fatalf("expected neighbors [B C], got %v", neighbors)
	}
}

func TestGraphDegree(t *testing.T) {
	g := NewGraph()
	g.AddEdge("A", "B")
	g.AddEdge("A", "C")
	if g.Degree("A") != 2 {
		t.Fatalf("expected degree 2, got %d", g.Degree("A"))
	}
	if g.Degree("B") != 1 {
		t.Fatalf("expected degree 1, got %d", g.Degree("B"))
	}
}

func TestGraphCopy(t *testing.T) {
	g := NewGraph()
	g.AddEdge("A", "B")
	g.AddEdge("B", "C")

	c := g.Copy()
	if !c.HasEdge("A", "B") || !c.HasEdge("B", "C") {
		t.Fatal("copy should have edges")
	}
	if len(c.Nodes()) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(c.Nodes()))
	}

	// Verify independence.
	c.AddEdge("C", "D")
	if g.HasNode("D") {
		t.Fatal("original should not be affected by copy")
	}
}

func TestGraphIsolatedNode(t *testing.T) {
	g := NewGraph()
	g.AddNode("X")
	g.AddEdge("A", "B")

	if g.Degree("X") != 0 {
		t.Fatal("isolated node should have degree 0")
	}
	neighbors := g.Neighbors("X")
	if len(neighbors) != 0 {
		t.Fatal("isolated node should have no neighbors")
	}
}

func TestGraphHasEdgeNonExistent(t *testing.T) {
	g := NewGraph()
	if g.HasEdge("X", "Y") {
		t.Fatal("should not have edge in empty graph")
	}
}
