//go:build unit

package graphgo

import (
	"sort"
	"testing"
)

func TestNewDiGraph(t *testing.T) {
	g := NewDiGraph()
	if g.NumberOfNodes() != 0 {
		t.Fatalf("expected 0 nodes, got %d", g.NumberOfNodes())
	}
	if g.NumberOfEdges() != 0 {
		t.Fatalf("expected 0 edges, got %d", g.NumberOfEdges())
	}
}

func TestDiGraphAddNode(t *testing.T) {
	g := NewDiGraph()
	g.AddNode("A")
	if !g.HasNode("A") {
		t.Fatal("expected node A")
	}
	if g.HasNode("B") {
		t.Fatal("unexpected node B")
	}
	// Idempotent.
	g.AddNode("A")
	if g.NumberOfNodes() != 1 {
		t.Fatalf("expected 1 node after duplicate add, got %d", g.NumberOfNodes())
	}
}

func TestDiGraphAddNodes(t *testing.T) {
	g := NewDiGraph()
	g.AddNodes("A", "B", "C")
	if g.NumberOfNodes() != 3 {
		t.Fatalf("expected 3 nodes, got %d", g.NumberOfNodes())
	}
}

func TestDiGraphRemoveNode(t *testing.T) {
	g := NewDiGraph()
	g.AddEdge("A", "B")
	g.AddEdge("C", "A")
	g.RemoveNode("A")
	if g.HasNode("A") {
		t.Fatal("node A should be removed")
	}
	if g.HasEdge("A", "B") || g.HasEdge("C", "A") {
		t.Fatal("incident edges should be removed")
	}
	if g.NumberOfEdges() != 0 {
		t.Fatalf("expected 0 edges, got %d", g.NumberOfEdges())
	}
	// Removing non-existent node is a no-op.
	g.RemoveNode("Z")
}

func TestDiGraphAddEdge(t *testing.T) {
	g := NewDiGraph()
	g.AddEdge("A", "B")
	if !g.HasEdge("A", "B") {
		t.Fatal("expected edge A->B")
	}
	if g.HasEdge("B", "A") {
		t.Fatal("should not have edge B->A")
	}
	// Nodes auto-created.
	if !g.HasNode("A") || !g.HasNode("B") {
		t.Fatal("nodes should be auto-created")
	}
}

func TestDiGraphRemoveEdge(t *testing.T) {
	g := NewDiGraph()
	g.AddEdge("A", "B")
	if err := g.RemoveEdge("A", "B"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.HasEdge("A", "B") {
		t.Fatal("edge should be removed")
	}
	// Nodes remain.
	if !g.HasNode("A") || !g.HasNode("B") {
		t.Fatal("nodes should remain after edge removal")
	}
	// Removing non-existent edge returns error.
	if err := g.RemoveEdge("X", "Y"); err == nil {
		t.Fatal("expected error for non-existent edge")
	}
}

func TestDiGraphNodesEdges(t *testing.T) {
	g := NewDiGraph()
	g.AddEdge("A", "B")
	g.AddEdge("B", "C")
	g.AddNode("D")

	nodes := g.Nodes()
	sort.Strings(nodes)
	expected := []string{"A", "B", "C", "D"}
	if len(nodes) != len(expected) {
		t.Fatalf("expected %d nodes, got %d", len(expected), len(nodes))
	}
	for i, n := range expected {
		if nodes[i] != n {
			t.Fatalf("expected node %s at index %d, got %s", n, i, nodes[i])
		}
	}

	edges := g.Edges()
	if len(edges) != 2 {
		t.Fatalf("expected 2 edges, got %d", len(edges))
	}
}

func TestDiGraphPredecessorsSuccessors(t *testing.T) {
	g := NewDiGraph()
	g.AddEdge("A", "B")
	g.AddEdge("C", "B")
	g.AddEdge("B", "D")

	preds := g.Predecessors("B")
	sort.Strings(preds)
	if len(preds) != 2 || preds[0] != "A" || preds[1] != "C" {
		t.Fatalf("expected predecessors [A C], got %v", preds)
	}

	succs := g.Successors("B")
	if len(succs) != 1 || succs[0] != "D" {
		t.Fatalf("expected successors [D], got %v", succs)
	}

	// Parents and Children are aliases.
	parents := g.Parents("B")
	sort.Strings(parents)
	if len(parents) != 2 || parents[0] != "A" || parents[1] != "C" {
		t.Fatalf("Parents should equal Predecessors")
	}

	children := g.Children("B")
	if len(children) != 1 || children[0] != "D" {
		t.Fatalf("Children should equal Successors")
	}
}

func TestDiGraphDegree(t *testing.T) {
	g := NewDiGraph()
	g.AddEdge("A", "B")
	g.AddEdge("C", "B")
	g.AddEdge("B", "D")

	if g.InDegree("B") != 2 {
		t.Fatalf("expected InDegree 2, got %d", g.InDegree("B"))
	}
	if g.OutDegree("B") != 1 {
		t.Fatalf("expected OutDegree 1, got %d", g.OutDegree("B"))
	}
	if g.InDegree("A") != 0 {
		t.Fatalf("expected InDegree 0 for A, got %d", g.InDegree("A"))
	}
}

func TestDiGraphAttributes(t *testing.T) {
	g := NewDiGraph()
	g.AddEdge("A", "B")

	g.NodeAttr("A")["color"] = "red"
	if g.NodeAttr("A")["color"] != "red" {
		t.Fatal("node attribute not set")
	}

	g.EdgeAttr("A", "B")["weight"] = 3.14
	if g.EdgeAttr("A", "B")["weight"] != 3.14 {
		t.Fatal("edge attribute not set")
	}

	if g.NodeAttr("Z") != nil {
		t.Fatal("expected nil for non-existent node")
	}
	if g.EdgeAttr("X", "Y") != nil {
		t.Fatal("expected nil for non-existent edge")
	}
}

func TestDiGraphCopy(t *testing.T) {
	g := NewDiGraph()
	g.AddEdge("A", "B")
	g.NodeAttr("A")["color"] = "blue"
	g.EdgeAttr("A", "B")["weight"] = 1

	c := g.Copy()

	// Verify structure.
	if !c.HasEdge("A", "B") {
		t.Fatal("copy should have edge A->B")
	}
	if c.NodeAttr("A")["color"] != "blue" {
		t.Fatal("copy should have node attribute")
	}
	if c.EdgeAttr("A", "B")["weight"] != 1 {
		t.Fatal("copy should have edge attribute")
	}

	// Verify independence.
	c.AddEdge("B", "C")
	if g.HasEdge("B", "C") {
		t.Fatal("original should not be affected by copy mutation")
	}
	c.NodeAttr("A")["color"] = "green"
	if g.NodeAttr("A")["color"] != "blue" {
		t.Fatal("original node attrs should not be affected")
	}
}

func TestDiGraphSubgraph(t *testing.T) {
	g := NewDiGraph()
	g.AddEdge("A", "B")
	g.AddEdge("B", "C")
	g.AddEdge("C", "D")
	g.NodeAttr("A")["val"] = 1

	sub := g.Subgraph([]string{"A", "B", "C"})
	if sub.NumberOfNodes() != 3 {
		t.Fatalf("expected 3 nodes, got %d", sub.NumberOfNodes())
	}
	if !sub.HasEdge("A", "B") || !sub.HasEdge("B", "C") {
		t.Fatal("subgraph should have internal edges")
	}
	if sub.HasEdge("C", "D") {
		t.Fatal("subgraph should not have edge to excluded node")
	}
	if sub.NodeAttr("A")["val"] != 1 {
		t.Fatal("subgraph should copy attributes")
	}

	// Subgraph with non-existent nodes is fine.
	sub2 := g.Subgraph([]string{"A", "Z"})
	if sub2.NumberOfNodes() != 1 {
		t.Fatalf("expected 1 node, got %d", sub2.NumberOfNodes())
	}
}
