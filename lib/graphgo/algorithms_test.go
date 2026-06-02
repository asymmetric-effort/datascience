//go:build unit

package graphgo

import (
	"sort"
	"testing"
)

func TestIsDAG(t *testing.T) {
	g := NewDiGraph()
	g.AddEdge("A", "B")
	g.AddEdge("B", "C")
	if !IsDAG(g) {
		t.Fatal("expected DAG")
	}

	g.AddEdge("C", "A")
	if IsDAG(g) {
		t.Fatal("expected cycle detected")
	}
}

func TestIsDAGEmpty(t *testing.T) {
	g := NewDiGraph()
	if !IsDAG(g) {
		t.Fatal("empty graph should be a DAG")
	}
}

func TestIsDAGSingleNode(t *testing.T) {
	g := NewDiGraph()
	g.AddNode("A")
	if !IsDAG(g) {
		t.Fatal("single node should be a DAG")
	}
}

func TestTopologicalSort(t *testing.T) {
	g := NewDiGraph()
	g.AddEdge("A", "B")
	g.AddEdge("A", "C")
	g.AddEdge("B", "D")
	g.AddEdge("C", "D")

	order, err := TopologicalSort(g)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 4 {
		t.Fatalf("expected 4 nodes, got %d", len(order))
	}

	// Verify topological property: for each edge u->v, u appears before v.
	pos := make(map[string]int)
	for i, n := range order {
		pos[n] = i
	}
	for _, e := range g.Edges() {
		if pos[e.Src] >= pos[e.Dst] {
			t.Fatalf("topological violation: %s (pos %d) should come before %s (pos %d)",
				e.Src, pos[e.Src], e.Dst, pos[e.Dst])
		}
	}
}

func TestTopologicalSortCycle(t *testing.T) {
	g := NewDiGraph()
	g.AddEdge("A", "B")
	g.AddEdge("B", "C")
	g.AddEdge("C", "A")

	_, err := TopologicalSort(g)
	if err == nil {
		t.Fatal("expected error for cyclic graph")
	}
}

func TestTopologicalSortDisconnected(t *testing.T) {
	g := NewDiGraph()
	g.AddEdge("A", "B")
	g.AddNode("C")
	g.AddNode("D")

	order, err := TopologicalSort(g)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 4 {
		t.Fatalf("expected 4 nodes, got %d", len(order))
	}
}

func TestAncestors(t *testing.T) {
	g := NewDiGraph()
	g.AddEdge("A", "B")
	g.AddEdge("B", "C")
	g.AddEdge("D", "C")

	anc := Ancestors(g, "C")
	expected := map[string]bool{"A": true, "B": true, "D": true}
	if len(anc) != len(expected) {
		t.Fatalf("expected %d ancestors, got %d", len(expected), len(anc))
	}
	for k := range expected {
		if !anc[k] {
			t.Fatalf("expected ancestor %s", k)
		}
	}
}

func TestAncestorsRoot(t *testing.T) {
	g := NewDiGraph()
	g.AddEdge("A", "B")
	anc := Ancestors(g, "A")
	if len(anc) != 0 {
		t.Fatalf("root should have no ancestors, got %v", anc)
	}
}

func TestDescendants(t *testing.T) {
	g := NewDiGraph()
	g.AddEdge("A", "B")
	g.AddEdge("A", "C")
	g.AddEdge("B", "D")

	desc := Descendants(g, "A")
	got := make([]string, 0, len(desc))
	for k := range desc {
		got = append(got, k)
	}
	sort.Strings(got)
	expected := []string{"B", "C", "D"}
	if len(got) != len(expected) {
		t.Fatalf("expected %d descendants, got %d", len(expected), len(got))
	}
	for i, n := range expected {
		if got[i] != n {
			t.Fatalf("expected %s, got %s", n, got[i])
		}
	}
}

func TestDescendantsLeaf(t *testing.T) {
	g := NewDiGraph()
	g.AddEdge("A", "B")
	desc := Descendants(g, "B")
	if len(desc) != 0 {
		t.Fatalf("leaf should have no descendants, got %v", desc)
	}
}

func TestAncestorsDescendantsDiamond(t *testing.T) {
	// Diamond: A -> B, A -> C, B -> D, C -> D
	g := NewDiGraph()
	g.AddEdge("A", "B")
	g.AddEdge("A", "C")
	g.AddEdge("B", "D")
	g.AddEdge("C", "D")

	anc := Ancestors(g, "D")
	if len(anc) != 3 {
		t.Fatalf("expected 3 ancestors of D, got %d: %v", len(anc), anc)
	}
	for _, n := range []string{"A", "B", "C"} {
		if !anc[n] {
			t.Fatalf("expected %s in ancestors of D", n)
		}
	}

	desc := Descendants(g, "A")
	if len(desc) != 3 {
		t.Fatalf("expected 3 descendants of A, got %d: %v", len(desc), desc)
	}
}
