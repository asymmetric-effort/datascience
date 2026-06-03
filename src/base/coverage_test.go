//go:build unit

package base

import (
	"testing"
)

// ---------------------------------------------------------------------------
// DAG: GetImmoralities, LocalIndependencies, IsIEquivalent, ActiveTrailNodes
// ---------------------------------------------------------------------------

func TestGetImmoralities_NoImmoralities(t *testing.T) {
	d := NewDAG()
	mustE(t, d.AddNode("A"))
	mustE(t, d.AddNode("B"))
	mustE(t, d.AddNode("C"))
	mustE(t, d.AddEdge("A", "B"))
	mustE(t, d.AddEdge("B", "C"))
	imm := d.GetImmoralities()
	if len(imm) != 0 {
		t.Errorf("expected 0 immoralities, got %d", len(imm))
	}
}

func TestGetImmoralities_WithImmorality(t *testing.T) {
	d := NewDAG()
	mustE(t, d.AddNode("A"))
	mustE(t, d.AddNode("B"))
	mustE(t, d.AddNode("C"))
	mustE(t, d.AddEdge("A", "C"))
	mustE(t, d.AddEdge("B", "C"))
	imm := d.GetImmoralities()
	if len(imm) != 1 {
		t.Fatalf("expected 1 immorality, got %d", len(imm))
	}
	if imm[0][1] != "C" {
		t.Errorf("expected child C, got %s", imm[0][1])
	}
}

func TestGetImmoralities_ParentsAdjacent(t *testing.T) {
	d := NewDAG()
	mustE(t, d.AddNode("A"))
	mustE(t, d.AddNode("B"))
	mustE(t, d.AddNode("C"))
	mustE(t, d.AddEdge("A", "C"))
	mustE(t, d.AddEdge("B", "C"))
	mustE(t, d.AddEdge("A", "B"))
	imm := d.GetImmoralities()
	if len(imm) != 0 {
		t.Errorf("expected 0 immoralities, got %d", len(imm))
	}
}

func TestLocalIndependencies_EmptyResult(t *testing.T) {
	d := NewDAG()
	mustE(t, d.AddNode("A"))
	mustE(t, d.AddNode("B"))
	mustE(t, d.AddEdge("A", "B"))
	ind := d.LocalIndependencies("A")
	if len(ind) != 0 {
		t.Errorf("expected 0 independencies for A, got %d", len(ind))
	}
}

func TestIsIEquivalent_DifferentNodeCount(t *testing.T) {
	d1 := NewDAG()
	mustE(t, d1.AddNode("A"))
	mustE(t, d1.AddNode("B"))
	d2 := NewDAG()
	mustE(t, d2.AddNode("A"))
	if d1.IsIEquivalent(d2) {
		t.Error("expected not I-equivalent with different node counts")
	}
}

func TestIsIEquivalent_DifferentNodes(t *testing.T) {
	d1 := NewDAG()
	mustE(t, d1.AddNode("A"))
	mustE(t, d1.AddNode("B"))
	d2 := NewDAG()
	mustE(t, d2.AddNode("A"))
	mustE(t, d2.AddNode("C"))
	if d1.IsIEquivalent(d2) {
		t.Error("expected not I-equivalent with different nodes")
	}
}

func TestIsIEquivalent_DifferentSkeleton(t *testing.T) {
	d1 := NewDAG()
	mustE(t, d1.AddNode("A"))
	mustE(t, d1.AddNode("B"))
	mustE(t, d1.AddEdge("A", "B"))
	d2 := NewDAG()
	mustE(t, d2.AddNode("A"))
	mustE(t, d2.AddNode("B"))
	if d1.IsIEquivalent(d2) {
		t.Error("expected not I-equivalent with different skeletons")
	}
}

func TestIsIEquivalent_SameEquiv(t *testing.T) {
	// Both DAGs: A->B->C  -- same skeleton, same immoralities (none).
	d1 := NewDAG()
	mustE(t, d1.AddNode("A"))
	mustE(t, d1.AddNode("B"))
	mustE(t, d1.AddNode("C"))
	mustE(t, d1.AddEdge("A", "B"))
	mustE(t, d1.AddEdge("B", "C"))

	d2 := NewDAG()
	mustE(t, d2.AddNode("A"))
	mustE(t, d2.AddNode("B"))
	mustE(t, d2.AddNode("C"))
	mustE(t, d2.AddEdge("B", "A"))
	mustE(t, d2.AddEdge("B", "C"))
	// d2: B->A, B->C. Same skeleton (A-B-C). No immoralities in either.
	if !d1.IsIEquivalent(d2) {
		t.Error("expected I-equivalent (same skeleton, no immoralities)")
	}
}

func TestActiveTrailNodes_Observed(t *testing.T) {
	d := NewDAG()
	mustE(t, d.AddNode("A"))
	mustE(t, d.AddNode("B"))
	mustE(t, d.AddNode("C"))
	mustE(t, d.AddEdge("A", "B"))
	mustE(t, d.AddEdge("B", "C"))
	active := d.ActiveTrailNodes("A", []string{"B"})
	for _, n := range active {
		if n == "C" {
			t.Error("C should not be reachable from A when B is observed (chain)")
		}
	}
}

// ---------------------------------------------------------------------------
// FromLavaan error paths
// ---------------------------------------------------------------------------

func TestFromLavaan_EmptyChild(t *testing.T) {
	_, err := FromLavaan(" ~ X")
	if err == nil {
		t.Error("expected error for empty child")
	}
}

func TestFromLavaan_EmptyParent(t *testing.T) {
	_, err := FromLavaan("Y ~ ")
	if err == nil {
		t.Error("expected error for empty parent list")
	}
}

func TestFromLavaan_Valid(t *testing.T) {
	d, err := FromLavaan("Y ~ X1 + X2\nZ ~ Y")
	if err != nil {
		t.Fatalf("FromLavaan failed: %v", err)
	}
	if !d.HasEdge("X1", "Y") || !d.HasEdge("X2", "Y") || !d.HasEdge("Y", "Z") {
		t.Error("missing expected edges")
	}
}

// ---------------------------------------------------------------------------
// FromDagitty error paths
// ---------------------------------------------------------------------------

func TestFromDagitty_EmptyDag(t *testing.T) {
	_, err := FromDagitty("dag { }")
	if err == nil {
		t.Error("expected error for empty dag body")
	}
}

func TestFromDagitty_NoEdges(t *testing.T) {
	_, err := FromDagitty("dag { X Y Z }")
	if err == nil {
		t.Error("expected error for no -> edges")
	}
}

func TestFromDagitty_Valid(t *testing.T) {
	d, err := FromDagitty("dag { X -> Y; Y -> Z }")
	if err != nil {
		t.Fatalf("FromDagitty failed: %v", err)
	}
	if !d.HasEdge("X", "Y") || !d.HasEdge("Y", "Z") {
		t.Error("missing expected edges")
	}
}

func TestFromDagitty_ChainedArrows(t *testing.T) {
	d, err := FromDagitty("dag { X -> Y -> Z }")
	if err != nil {
		t.Fatalf("FromDagitty failed: %v", err)
	}
	if !d.HasEdge("X", "Y") || !d.HasEdge("Y", "Z") {
		t.Error("missing expected edges")
	}
}

func TestFromDagitty_NoDagPrefix(t *testing.T) {
	d, err := FromDagitty("{ X -> Y }")
	if err != nil {
		t.Fatalf("FromDagitty failed: %v", err)
	}
	if !d.HasEdge("X", "Y") {
		t.Error("missing expected edge")
	}
}

// ---------------------------------------------------------------------------
// ADMG: AddBidirectedEdge error paths, BidirectedEdges
// ---------------------------------------------------------------------------

func TestADMG_AddBidirectedEdge_MissingNode(t *testing.T) {
	a := NewADMG()
	mustE(t, a.AddNode("A"))
	err := a.AddBidirectedEdge("A", "Z")
	if err == nil {
		t.Error("expected error for missing node Z")
	}
	err = a.AddBidirectedEdge("Z", "A")
	if err == nil {
		t.Error("expected error for missing node Z")
	}
}

func TestADMG_AddBidirectedEdge_Duplicate(t *testing.T) {
	a := NewADMG()
	mustE(t, a.AddNode("A"))
	mustE(t, a.AddNode("B"))
	mustE(t, a.AddBidirectedEdge("A", "B"))
	err := a.AddBidirectedEdge("A", "B")
	if err == nil {
		t.Error("expected error for duplicate bidirected edge")
	}
}

func TestADMG_BidirectedEdges(t *testing.T) {
	a := NewADMG()
	mustE(t, a.AddNode("A"))
	mustE(t, a.AddNode("B"))
	mustE(t, a.AddNode("C"))
	mustE(t, a.AddBidirectedEdge("A", "B"))
	mustE(t, a.AddBidirectedEdge("B", "C"))
	edges := a.BidirectedEdges()
	if len(edges) != 2 {
		t.Errorf("expected 2 bidirected edges, got %d", len(edges))
	}
}

// ---------------------------------------------------------------------------
// PDAG: AddUndirectedEdge errors, ToDAG
// ---------------------------------------------------------------------------

func TestPDAG_AddUndirectedEdge_MissingNode(t *testing.T) {
	pd := NewPDAG()
	mustE(t, pd.AddNode("A"))
	err := pd.AddUndirectedEdge("A", "Z")
	if err == nil {
		t.Error("expected error for missing node Z")
	}
	err = pd.AddUndirectedEdge("Z", "A")
	if err == nil {
		t.Error("expected error for missing node Z")
	}
}

func TestPDAG_AddUndirectedEdge_Duplicate(t *testing.T) {
	pd := NewPDAG()
	mustE(t, pd.AddNode("A"))
	mustE(t, pd.AddNode("B"))
	mustE(t, pd.AddUndirectedEdge("A", "B"))
	err := pd.AddUndirectedEdge("A", "B")
	if err == nil {
		t.Error("expected error for duplicate undirected edge")
	}
}

func TestPDAG_AddDirectedEdge_MissingNode(t *testing.T) {
	pd := NewPDAG()
	mustE(t, pd.AddNode("A"))
	err := pd.AddDirectedEdge("A", "Z")
	if err == nil {
		t.Error("expected error for missing node Z")
	}
	err = pd.AddDirectedEdge("Z", "A")
	if err == nil {
		t.Error("expected error for missing node Z")
	}
}

func TestPDAG_AddDirectedEdge_Duplicate(t *testing.T) {
	pd := NewPDAG()
	mustE(t, pd.AddNode("A"))
	mustE(t, pd.AddNode("B"))
	mustE(t, pd.AddDirectedEdge("A", "B"))
	err := pd.AddDirectedEdge("A", "B")
	if err == nil {
		t.Error("expected error for duplicate directed edge")
	}
}

func TestPDAG_ToDAG_AllDirected(t *testing.T) {
	pd := NewPDAG()
	mustE(t, pd.AddNode("A"))
	mustE(t, pd.AddNode("B"))
	mustE(t, pd.AddDirectedEdge("A", "B"))
	dag, err := pd.ToDAG()
	if err != nil {
		t.Fatalf("ToDAG failed: %v", err)
	}
	if !dag.HasEdge("A", "B") {
		t.Error("expected edge A->B in DAG")
	}
}

func TestPDAG_ToDAG_WithUndirected(t *testing.T) {
	pd := NewPDAG()
	mustE(t, pd.AddNode("A"))
	mustE(t, pd.AddNode("B"))
	mustE(t, pd.AddUndirectedEdge("A", "B"))
	dag, err := pd.ToDAG()
	if err != nil {
		t.Fatalf("ToDAG failed: %v", err)
	}
	if !dag.HasEdge("A", "B") && !dag.HasEdge("B", "A") {
		t.Error("expected an edge in either direction")
	}
}

// ---------------------------------------------------------------------------
// MAG: FromADMG, MSeparation
// ---------------------------------------------------------------------------

func TestMAG_FromADMG_Nil(t *testing.T) {
	_, err := FromADMG(nil)
	if err == nil {
		t.Error("expected error for nil ADMG")
	}
}

func TestMAG_FromADMG_Basic(t *testing.T) {
	a := NewADMG()
	mustE(t, a.AddNode("A"))
	mustE(t, a.AddNode("B"))
	mustE(t, a.AddNode("C"))
	mustE(t, a.AddDirectedEdge("A", "B"))
	mustE(t, a.AddDirectedEdge("B", "C"))
	mag, err := FromADMG(a)
	if err != nil {
		t.Fatalf("FromADMG failed: %v", err)
	}
	if !mag.HasNode("A") || !mag.HasNode("B") || !mag.HasNode("C") {
		t.Error("missing nodes")
	}
}

func TestMAG_MSeparation(t *testing.T) {
	mag := NewMAG()
	mustE(t, mag.AddNode("A"))
	mustE(t, mag.AddNode("B"))
	mustE(t, mag.AddNode("C"))
	mustE(t, mag.AddDirectedEdge("A", "B"))
	mustE(t, mag.AddDirectedEdge("B", "C"))
	sep := mag.MSeparation(
		map[string]bool{"A": true},
		map[string]bool{"C": true},
		map[string]bool{"B": true},
	)
	if !sep {
		t.Error("expected A and C to be m-separated given B")
	}
}

// ---------------------------------------------------------------------------
// SimpleCausalModel: Sample
// ---------------------------------------------------------------------------

func TestSimpleCausalModel_Sample_NoEquation(t *testing.T) {
	d := NewDAG()
	mustE(t, d.AddNode("X"))
	mustE(t, d.AddNode("Y"))
	mustE(t, d.AddEdge("X", "Y"))
	m := NewSimpleCausalModel(d)
	m.SetEquation("Y", func(parents map[string]float64) float64 {
		return parents["X"] * 2
	})
	vals, err := m.Sample(map[string]float64{"X": 3.0})
	if err != nil {
		t.Fatalf("Sample failed: %v", err)
	}
	if vals["Y"] != 6.0 {
		t.Errorf("expected Y=6.0, got %f", vals["Y"])
	}
	if vals["X"] != 3.0 {
		t.Errorf("expected X=3.0, got %f", vals["X"])
	}
}

// ---------------------------------------------------------------------------
// MinimalDSeparator
// ---------------------------------------------------------------------------

func TestMinimalDSeparator_NotSeparable(t *testing.T) {
	d := NewDAG()
	mustE(t, d.AddNode("A"))
	mustE(t, d.AddNode("B"))
	mustE(t, d.AddEdge("A", "B"))
	_, ok := d.MinimalDSeparator("A", "B")
	if ok {
		t.Error("A and B should not be d-separable (direct edge)")
	}
}

func TestMinimalDSeparator_Chain(t *testing.T) {
	d := NewDAG()
	mustE(t, d.AddNode("A"))
	mustE(t, d.AddNode("B"))
	mustE(t, d.AddNode("C"))
	mustE(t, d.AddEdge("A", "B"))
	mustE(t, d.AddEdge("B", "C"))
	sep, ok := d.MinimalDSeparator("A", "C")
	if !ok {
		t.Fatal("expected A and C to be d-separable")
	}
	if len(sep) != 1 || sep[0] != "B" {
		t.Errorf("expected {B}, got %v", sep)
	}
}

// ---------------------------------------------------------------------------
// Helper
// ---------------------------------------------------------------------------

func mustE(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
