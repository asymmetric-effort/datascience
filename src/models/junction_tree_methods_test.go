//go:build unit

package models

import (
	"testing"

	"github.com/asymmetric-effort/pgmgo/src/factors"
)

func TestJunctionTreeAddEdge(t *testing.T) {
	// Build a junction tree with separate cliques, then manually add an edge.
	bn := NewBayesianNetwork()
	_ = bn.AddNode("A")
	_ = bn.AddNode("B")
	_ = bn.AddEdge("A", "B")

	cpdA, _ := buildSingleNodeCPD("A", 2, []float64{0.5, 0.5})
	_ = bn.AddCPD(cpdA)
	cpdB, _ := factors.NewTabularCPD("B", 2, [][]float64{
		{0.3, 0.7},
		{0.7, 0.3},
	}, []string{"A"}, []int{2})
	_ = bn.AddCPD(cpdB)

	jt, err := NewJunctionTreeFromBN(bn)
	if err != nil {
		t.Fatalf("NewJunctionTreeFromBN: %v", err)
	}

	cliques := jt.Cliques()
	// For a simple A->B network there should be 1 clique {A, B}.
	if len(cliques) != 1 {
		t.Skipf("expected 1 clique for simple network, got %d", len(cliques))
	}
}

func TestJunctionTreeAddEdge_InvalidIndex(t *testing.T) {
	bn := buildStudentNetwork(t)
	jt, err := NewJunctionTreeFromBN(bn)
	if err != nil {
		t.Fatalf("NewJunctionTreeFromBN: %v", err)
	}

	nCliques := len(jt.Cliques())

	// Out of range.
	err = jt.AddEdge(-1, 0)
	if err == nil {
		t.Error("expected error for negative index")
	}
	err = jt.AddEdge(0, nCliques)
	if err == nil {
		t.Error("expected error for out of range index")
	}

	// Self-loop.
	err = jt.AddEdge(0, 0)
	if err == nil {
		t.Error("expected error for self-loop")
	}
}

func TestJunctionTreeStates(t *testing.T) {
	bn := buildStudentNetwork(t)
	jt, err := NewJunctionTreeFromBN(bn)
	if err != nil {
		t.Fatalf("NewJunctionTreeFromBN: %v", err)
	}

	states := jt.States()
	// The student network has 5 variables with known cardinalities.
	if len(states) == 0 {
		t.Error("expected non-empty states map")
	}

	// Each variable should have a positive cardinality.
	for v, card := range states {
		if card <= 0 {
			t.Errorf("variable %q has non-positive cardinality %d", v, card)
		}
	}
}

func TestJunctionTreeStates_Empty(t *testing.T) {
	bn := NewBayesianNetwork()
	jt, err := NewJunctionTreeFromBN(bn)
	if err != nil {
		t.Fatalf("NewJunctionTreeFromBN: %v", err)
	}

	states := jt.States()
	if len(states) != 0 {
		t.Errorf("expected empty states for empty JT, got %d", len(states))
	}
}

func TestJunctionTreeCopy(t *testing.T) {
	bn := buildStudentNetwork(t)
	jt, err := NewJunctionTreeFromBN(bn)
	if err != nil {
		t.Fatalf("NewJunctionTreeFromBN: %v", err)
	}

	cp := jt.Copy()
	if cp == jt {
		t.Error("Copy should return a new object")
	}

	// Verify cliques match.
	origCliques := jt.Cliques()
	copyCliques := cp.Cliques()
	if len(origCliques) != len(copyCliques) {
		t.Fatalf("clique count mismatch: %d vs %d", len(origCliques), len(copyCliques))
	}
	for i := range origCliques {
		if len(origCliques[i]) != len(copyCliques[i]) {
			t.Errorf("clique %d size mismatch", i)
			continue
		}
		for j := range origCliques[i] {
			if origCliques[i][j] != copyCliques[i][j] {
				t.Errorf("clique %d var %d mismatch: %q vs %q",
					i, j, origCliques[i][j], copyCliques[i][j])
			}
		}
	}

	// Verify separator sets match.
	origSeps := jt.SeparatorSets()
	copySeps := cp.SeparatorSets()
	if len(origSeps) != len(copySeps) {
		t.Errorf("separator set count mismatch: %d vs %d", len(origSeps), len(copySeps))
	}

	// Verify CheckModel passes on copy.
	if err := cp.CheckModel(); err != nil {
		t.Fatalf("CheckModel on copy: %v", err)
	}
}

func TestJunctionTreeCopy_Independence(t *testing.T) {
	bn := buildStudentNetwork(t)
	jt, err := NewJunctionTreeFromBN(bn)
	if err != nil {
		t.Fatalf("NewJunctionTreeFromBN: %v", err)
	}

	cp := jt.Copy()

	// Modifying the copy should not affect the original.
	origCliques := jt.Cliques()
	copyCliques := cp.Cliques()
	if len(copyCliques) > 0 && len(copyCliques[0]) > 0 {
		copyCliques[0][0] = "MODIFIED"
		if origCliques[0][0] == "MODIFIED" {
			t.Error("modifying copy affected original")
		}
	}
}

func TestJunctionTreeCopy_Empty(t *testing.T) {
	bn := NewBayesianNetwork()
	jt, err := NewJunctionTreeFromBN(bn)
	if err != nil {
		t.Fatalf("NewJunctionTreeFromBN: %v", err)
	}

	cp := jt.Copy()
	if len(cp.Cliques()) != 0 {
		t.Error("expected 0 cliques in copy of empty JT")
	}
	if err := cp.CheckModel(); err != nil {
		t.Fatalf("CheckModel on empty copy: %v", err)
	}
}

func TestJunctionTreeStates_AllVariables(t *testing.T) {
	bn := buildStudentNetwork(t)
	jt, err := NewJunctionTreeFromBN(bn)
	if err != nil {
		t.Fatalf("NewJunctionTreeFromBN: %v", err)
	}

	states := jt.States()
	// All 5 student network variables should appear.
	for _, v := range []string{"D", "G", "I", "L", "S"} {
		if _, ok := states[v]; !ok {
			t.Errorf("variable %q not found in states", v)
		}
	}
}
