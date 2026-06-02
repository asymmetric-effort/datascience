//go:build unit

package graphgo

import (
	"sort"
	"testing"
)

// buildStudentNetwork creates the classic "Student" Bayesian network:
//
//	Difficulty -> Grade <- Intelligence
//	Intelligence -> SAT
//	Grade -> Letter
//
// Nodes: D, I, G, S, L
func buildStudentNetwork() *DiGraph {
	g := NewDiGraph()
	g.AddNodes("D", "I", "G", "S", "L")
	g.AddEdge("D", "G")
	g.AddEdge("I", "G")
	g.AddEdge("I", "S")
	g.AddEdge("G", "L")
	return g
}

func TestDSeparation_IndependentNoEvidence(t *testing.T) {
	g := buildStudentNetwork()

	// D and I are d-separated given empty set (no common observed descendant).
	// D -> G <- I is a collider; without observing G, D _|_ I.
	x := map[string]bool{"D": true}
	y := map[string]bool{"I": true}
	z := map[string]bool{}

	if !DSeparation(g, x, y, z) {
		t.Error("D and I should be d-separated given empty evidence")
	}
}

func TestDSeparation_ColliderOpened(t *testing.T) {
	g := buildStudentNetwork()

	// Observing G (the collider) opens the path D -> G <- I.
	x := map[string]bool{"D": true}
	y := map[string]bool{"I": true}
	z := map[string]bool{"G": true}

	if DSeparation(g, x, y, z) {
		t.Error("D and I should NOT be d-separated when G is observed (collider opened)")
	}
}

func TestDSeparation_ColliderDescendantOpened(t *testing.T) {
	g := buildStudentNetwork()

	// Observing L (descendant of collider G) also opens D -> G <- I.
	x := map[string]bool{"D": true}
	y := map[string]bool{"I": true}
	z := map[string]bool{"L": true}

	if DSeparation(g, x, y, z) {
		t.Error("D and I should NOT be d-separated when L (descendant of collider G) is observed")
	}
}

func TestDSeparation_ChainBlocked(t *testing.T) {
	g := buildStudentNetwork()

	// D -> G -> L is a chain. Observing G blocks the path.
	x := map[string]bool{"D": true}
	y := map[string]bool{"L": true}
	z := map[string]bool{"G": true}

	if !DSeparation(g, x, y, z) {
		t.Error("D and L should be d-separated when G is observed (chain blocked)")
	}
}

func TestDSeparation_ChainOpen(t *testing.T) {
	g := buildStudentNetwork()

	// D -> G -> L without observing G: path is open.
	x := map[string]bool{"D": true}
	y := map[string]bool{"L": true}
	z := map[string]bool{}

	if DSeparation(g, x, y, z) {
		t.Error("D and L should NOT be d-separated with no evidence (chain open)")
	}
}

func TestDSeparation_ForkBlocked(t *testing.T) {
	g := buildStudentNetwork()

	// I -> G and I -> S: fork at I. Observing I blocks the path G-I-S.
	x := map[string]bool{"G": true}
	y := map[string]bool{"S": true}
	z := map[string]bool{"I": true}

	if !DSeparation(g, x, y, z) {
		t.Error("G and S should be d-separated when I is observed (fork blocked)")
	}
}

func TestDSeparation_ForkOpen(t *testing.T) {
	g := buildStudentNetwork()

	// I -> G and I -> S without observing I: path is open.
	x := map[string]bool{"G": true}
	y := map[string]bool{"S": true}
	z := map[string]bool{}

	if DSeparation(g, x, y, z) {
		t.Error("G and S should NOT be d-separated with no evidence (fork open)")
	}
}

func TestDSeparation_SameNode(t *testing.T) {
	g := buildStudentNetwork()

	// A node is trivially not d-separated from itself.
	x := map[string]bool{"D": true}
	y := map[string]bool{"D": true}
	z := map[string]bool{}

	if DSeparation(g, x, y, z) {
		t.Error("a node should not be d-separated from itself")
	}
}

func TestDSeparation_DisconnectedNodes(t *testing.T) {
	g := buildStudentNetwork()

	// D and S: only path is D -> G <- I -> S.
	// G is a collider on this path, unobserved, so path is blocked.
	x := map[string]bool{"D": true}
	y := map[string]bool{"S": true}
	z := map[string]bool{}

	if !DSeparation(g, x, y, z) {
		t.Error("D and S should be d-separated with no evidence (collider G blocks)")
	}
}

func TestDSeparation_MultiplePaths(t *testing.T) {
	// Create a graph with two paths: A->C->D and A->B->D (fork at A, merge at D).
	g := NewDiGraph()
	g.AddNodes("A", "B", "C", "D")
	g.AddEdge("A", "B")
	g.AddEdge("A", "C")
	g.AddEdge("B", "D")
	g.AddEdge("C", "D")

	// A and D: paths open without evidence.
	x := map[string]bool{"A": true}
	y := map[string]bool{"D": true}

	if DSeparation(g, x, y, map[string]bool{}) {
		t.Error("A and D should NOT be d-separated with no evidence")
	}

	// Observing both B and C blocks all paths.
	z := map[string]bool{"B": true, "C": true}
	if !DSeparation(g, x, y, z) {
		t.Error("A and D should be d-separated when both B and C are observed")
	}
}

// --- MarkovBlanket tests ---

func TestMarkovBlanket_Grade(t *testing.T) {
	g := buildStudentNetwork()

	// G's parents: D, I. G's children: L. Parents of L (besides G): none.
	// So MB(G) = {D, I, L}.
	mb := MarkovBlanket(g, "G")
	expected := map[string]bool{"D": true, "I": true, "L": true}

	if len(mb) != len(expected) {
		t.Fatalf("MarkovBlanket(G) has %d elements, want %d: %v", len(mb), len(expected), mb)
	}
	for k := range expected {
		if !mb[k] {
			t.Errorf("expected %s in MarkovBlanket(G)", k)
		}
	}
}

func TestMarkovBlanket_Intelligence(t *testing.T) {
	g := buildStudentNetwork()

	// I's parents: none. I's children: G, S.
	// Parents of G (besides I): D. Parents of S (besides I): none.
	// So MB(I) = {G, S, D}.
	mb := MarkovBlanket(g, "I")
	expected := map[string]bool{"G": true, "S": true, "D": true}

	if len(mb) != len(expected) {
		t.Fatalf("MarkovBlanket(I) has %d elements, want %d: %v", len(mb), len(expected), mb)
	}
	for k := range expected {
		if !mb[k] {
			t.Errorf("expected %s in MarkovBlanket(I)", k)
		}
	}
}

func TestMarkovBlanket_Leaf(t *testing.T) {
	g := buildStudentNetwork()

	// L's parents: G. L's children: none.
	// MB(L) = {G}.
	mb := MarkovBlanket(g, "L")
	if len(mb) != 1 || !mb["G"] {
		t.Errorf("MarkovBlanket(L) = %v, want {G}", mb)
	}
}

func TestMarkovBlanket_Root(t *testing.T) {
	g := buildStudentNetwork()

	// D's parents: none. D's children: G.
	// Parents of G (besides D): I.
	// MB(D) = {G, I}.
	mb := MarkovBlanket(g, "D")
	expected := map[string]bool{"G": true, "I": true}

	if len(mb) != len(expected) {
		t.Fatalf("MarkovBlanket(D) has %d elements, want %d: %v", len(mb), len(expected), mb)
	}
	for k := range expected {
		if !mb[k] {
			t.Errorf("expected %s in MarkovBlanket(D)", k)
		}
	}
}

func TestMarkovBlanket_Isolated(t *testing.T) {
	g := NewDiGraph()
	g.AddNode("X")

	mb := MarkovBlanket(g, "X")
	if len(mb) != 0 {
		t.Errorf("MarkovBlanket of isolated node should be empty, got %v", mb)
	}
}

func TestMarkovBlanket_CoParents(t *testing.T) {
	// A->C, B->C, C->D: MB(A) should include B (co-parent at C), C (child), but not D.
	g := NewDiGraph()
	g.AddNodes("A", "B", "C", "D")
	g.AddEdge("A", "C")
	g.AddEdge("B", "C")
	g.AddEdge("C", "D")

	mb := MarkovBlanket(g, "A")
	expected := []string{"B", "C"}
	sort.Strings(expected)

	got := make([]string, 0, len(mb))
	for k := range mb {
		got = append(got, k)
	}
	sort.Strings(got)

	if len(got) != len(expected) {
		t.Fatalf("MarkovBlanket(A) = %v, want %v", got, expected)
	}
	for i := range expected {
		if got[i] != expected[i] {
			t.Errorf("MarkovBlanket(A)[%d] = %s, want %s", i, got[i], expected[i])
		}
	}
}
