//go:build unit

package models

import (
	"strings"
	"testing"

	"github.com/asymmetric-effort/datascience/lib/graphgo"
	"github.com/asymmetric-effort/datascience/lib/pgm/factors"
)

// ---------------------------------------------------------------------------
// AddEdge coverage tests
// ---------------------------------------------------------------------------

func TestAddEdge_InvalidCliqueA_Negative(t *testing.T) {
	jt := buildMultiCliqueJT(t)
	err := jt.AddEdge(-1, 0)
	if err == nil {
		t.Fatal("expected error for negative cliqueA")
	}
	if !strings.Contains(err.Error(), "out of range") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestAddEdge_InvalidCliqueA_TooLarge(t *testing.T) {
	jt := buildMultiCliqueJT(t)
	nCliques := len(jt.Cliques())
	err := jt.AddEdge(nCliques, 0)
	if err == nil {
		t.Fatal("expected error for cliqueA out of range")
	}
}

func TestAddEdge_InvalidCliqueB_Negative(t *testing.T) {
	jt := buildMultiCliqueJT(t)
	err := jt.AddEdge(0, -1)
	if err == nil {
		t.Fatal("expected error for negative cliqueB")
	}
}

func TestAddEdge_InvalidCliqueB_TooLarge(t *testing.T) {
	jt := buildMultiCliqueJT(t)
	nCliques := len(jt.Cliques())
	err := jt.AddEdge(0, nCliques)
	if err == nil {
		t.Fatal("expected error for cliqueB out of range")
	}
}

func TestAddEdge_SelfLoop(t *testing.T) {
	jt := buildMultiCliqueJT(t)
	err := jt.AddEdge(0, 0)
	if err == nil {
		t.Fatal("expected error for self-loop")
	}
	if !strings.Contains(err.Error(), "self-loop") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestAddEdge_AlreadyExists(t *testing.T) {
	jt := buildMultiCliqueJT(t)
	// Add edge between cliques 0 and 1.
	err := jt.AddEdge(0, 1)
	if err != nil {
		t.Fatalf("first AddEdge: %v", err)
	}
	// Try to add the same edge again.
	err = jt.AddEdge(0, 1)
	if err == nil {
		t.Fatal("expected error for already existing edge")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestAddEdge_AlreadyExists_Reverse(t *testing.T) {
	jt := buildMultiCliqueJT(t)
	err := jt.AddEdge(0, 1)
	if err != nil {
		t.Fatalf("first AddEdge: %v", err)
	}
	// Adding in reverse order should also fail (undirected edge).
	err = jt.AddEdge(1, 0)
	if err == nil {
		t.Fatal("expected error for already existing edge (reverse)")
	}
}

func TestAddEdge_ValidEdge_WithOverlap(t *testing.T) {
	jt := buildMultiCliqueJT(t)
	// Cliques share variable B.
	err := jt.AddEdge(0, 1)
	if err != nil {
		t.Fatalf("AddEdge: %v", err)
	}

	// Verify separator was computed correctly.
	seps := jt.SeparatorSets()
	found := false
	for _, sep := range seps {
		for _, v := range sep {
			if v == "B" {
				found = true
			}
		}
	}
	if !found {
		t.Error("expected separator to contain 'B'")
	}
}

func TestAddEdge_ValidEdge_NoOverlap(t *testing.T) {
	jt := buildDisjointJT(t)
	err := jt.AddEdge(0, 1)
	if err != nil {
		t.Fatalf("AddEdge: %v", err)
	}

	// Separator should be empty for disjoint cliques.
	seps := jt.SeparatorSets()
	for _, sep := range seps {
		if len(sep) != 0 {
			t.Errorf("expected empty separator for disjoint cliques, got %v", sep)
		}
	}
}

func TestAddEdge_CanonicalKeyOrdering(t *testing.T) {
	// When cliqueA > cliqueB, the key should still be canonical (smaller first).
	jt := buildMultiCliqueJT(t)
	err := jt.AddEdge(1, 0)
	if err != nil {
		t.Fatalf("AddEdge: %v", err)
	}

	seps := jt.SeparatorSets()
	// The key should be "0-1" not "1-0".
	if _, ok := seps["0-1"]; !ok {
		t.Errorf("expected canonical key '0-1', got keys: %v", seps)
	}
}

// ---------------------------------------------------------------------------
// Helpers to build JTs manually for testing
// ---------------------------------------------------------------------------

// buildMultiCliqueJT creates a JT with 2 cliques that share variable B:
// Clique 0: {A, B}, Clique 1: {B, C}
func buildMultiCliqueJT(t *testing.T) *JunctionTree {
	t.Helper()
	pA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.5, 0.5})
	pBA, _ := factors.NewDiscreteFactor([]string{"B", "A"}, []int{2, 2}, []float64{
		0.3, 0.7, 0.7, 0.3,
	})
	pCB, _ := factors.NewDiscreteFactor([]string{"C", "B"}, []int{2, 2}, []float64{
		0.6, 0.4, 0.4, 0.6,
	})

	tree := graphgo.NewGraph()
	tree.AddNode("0")
	tree.AddNode("1")
	// No edges initially — tests will add them.

	return &JunctionTree{
		cliques:       [][]string{{"A", "B"}, {"B", "C"}},
		tree:          tree,
		separators:    make(map[string][]string),
		cliqueFactors: map[int][]*factors.DiscreteFactor{0: {pA, pBA}, 1: {pCB}},
	}
}

// buildDisjointJT creates a JT with 2 cliques that share no variables:
// Clique 0: {A}, Clique 1: {B}
func buildDisjointJT(t *testing.T) *JunctionTree {
	t.Helper()
	pA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.5, 0.5})
	pB, _ := factors.NewDiscreteFactor([]string{"B"}, []int{2}, []float64{0.5, 0.5})

	tree := graphgo.NewGraph()
	tree.AddNode("0")
	tree.AddNode("1")

	return &JunctionTree{
		cliques:       [][]string{{"A"}, {"B"}},
		tree:          tree,
		separators:    make(map[string][]string),
		cliqueFactors: map[int][]*factors.DiscreteFactor{0: {pA}, 1: {pB}},
	}
}
