//go:build unit

package models

import (
	"testing"

	"github.com/asymmetric-effort/pgmgo/src/factors"
)

// buildLGChain creates a 3-variable chain: X -> Y -> Z with known linear relationships.
//
//	X ~ N(5, 1)
//	Y | X ~ N(2 + 0.5*X, 0.5)
//	Z | Y ~ N(-1 + 1.5*Y, 0.8)
func buildLGChain(t *testing.T) *LinearGaussianBayesianNetwork {
	t.Helper()
	bn := NewLinearGaussianBayesianNetwork()

	for _, node := range []string{"X", "Y", "Z"} {
		if err := bn.AddNode(node); err != nil {
			t.Fatalf("AddNode(%q): %v", node, err)
		}
	}
	if err := bn.AddEdge("X", "Y"); err != nil {
		t.Fatalf("AddEdge(X,Y): %v", err)
	}
	if err := bn.AddEdge("Y", "Z"); err != nil {
		t.Fatalf("AddEdge(Y,Z): %v", err)
	}

	cpdX, err := factors.NewLinearGaussianCPD("X", 5.0, nil, 1.0, nil)
	if err != nil {
		t.Fatalf("NewLinearGaussianCPD(X): %v", err)
	}
	if err := bn.AddLinearGaussianCPD(cpdX); err != nil {
		t.Fatalf("AddLinearGaussianCPD(X): %v", err)
	}

	cpdY, err := factors.NewLinearGaussianCPD("Y", 2.0, []float64{0.5}, 0.5, []string{"X"})
	if err != nil {
		t.Fatalf("NewLinearGaussianCPD(Y): %v", err)
	}
	if err := bn.AddLinearGaussianCPD(cpdY); err != nil {
		t.Fatalf("AddLinearGaussianCPD(Y): %v", err)
	}

	cpdZ, err := factors.NewLinearGaussianCPD("Z", -1.0, []float64{1.5}, 0.8, []string{"Y"})
	if err != nil {
		t.Fatalf("NewLinearGaussianCPD(Z): %v", err)
	}
	if err := bn.AddLinearGaussianCPD(cpdZ); err != nil {
		t.Fatalf("AddLinearGaussianCPD(Z): %v", err)
	}

	return bn
}

func TestNewLinearGaussianBayesianNetwork(t *testing.T) {
	bn := NewLinearGaussianBayesianNetwork()
	if bn == nil {
		t.Fatal("expected non-nil network")
	}
	if bn.BayesianNetwork == nil {
		t.Fatal("expected non-nil embedded BayesianNetwork")
	}
	if len(bn.lgCPDs) != 0 {
		t.Fatalf("expected empty lgCPDs map, got %d", len(bn.lgCPDs))
	}
}

func TestLinearGaussianBN_AddLinearGaussianCPD(t *testing.T) {
	bn := NewLinearGaussianBayesianNetwork()
	_ = bn.AddNode("X")
	_ = bn.AddNode("Y")
	_ = bn.AddEdge("X", "Y")

	// nil CPD should fail.
	if err := bn.AddLinearGaussianCPD(nil); err == nil {
		t.Error("expected error for nil CPD")
	}

	// Variable not in network should fail.
	cpd, _ := factors.NewLinearGaussianCPD("Z", 0, nil, 1, nil)
	if err := bn.AddLinearGaussianCPD(cpd); err == nil {
		t.Error("expected error for variable not in network")
	}

	// Wrong evidence should fail.
	cpdBad, _ := factors.NewLinearGaussianCPD("Y", 0, []float64{1.0}, 1, []string{"Z"})
	if err := bn.AddLinearGaussianCPD(cpdBad); err == nil {
		t.Error("expected error for mismatched evidence")
	}

	// Correct CPD should succeed.
	cpdGood, _ := factors.NewLinearGaussianCPD("Y", 0, []float64{1.0}, 1, []string{"X"})
	if err := bn.AddLinearGaussianCPD(cpdGood); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLinearGaussianBN_GetLinearGaussianCPD(t *testing.T) {
	bn := buildLGChain(t)

	cpd := bn.GetLinearGaussianCPD("X")
	if cpd == nil {
		t.Fatal("expected non-nil CPD for X")
	}
	if cpd.Variable() != "X" {
		t.Errorf("expected variable X, got %q", cpd.Variable())
	}

	cpd = bn.GetLinearGaussianCPD("nonexistent")
	if cpd != nil {
		t.Error("expected nil CPD for nonexistent variable")
	}
}

func TestLinearGaussianBN_CheckModel(t *testing.T) {
	bn := buildLGChain(t)
	if err := bn.CheckModel(); err != nil {
		t.Fatalf("CheckModel: %v", err)
	}
}

func TestLinearGaussianBN_CheckModel_MissingCPD(t *testing.T) {
	bn := NewLinearGaussianBayesianNetwork()
	_ = bn.AddNode("X")
	_ = bn.AddNode("Y")
	_ = bn.AddEdge("X", "Y")

	cpdX, _ := factors.NewLinearGaussianCPD("X", 0, nil, 1, nil)
	_ = bn.AddLinearGaussianCPD(cpdX)

	if err := bn.CheckModel(); err == nil {
		t.Error("expected error for missing CPD on Y")
	}
}

func TestLinearGaussianBN_CheckModel_EvidenceMismatch(t *testing.T) {
	bn := NewLinearGaussianBayesianNetwork()
	_ = bn.AddNode("X")
	_ = bn.AddNode("Y")
	_ = bn.AddEdge("X", "Y")

	cpdX, _ := factors.NewLinearGaussianCPD("X", 0, nil, 1, nil)
	_ = bn.AddLinearGaussianCPD(cpdX)

	// Add Y CPD with no evidence (but Y has parent X).
	cpdY, _ := factors.NewLinearGaussianCPD("Y", 0, nil, 1, nil)
	bn.lgCPDs["Y"] = cpdY // bypass validation to force mismatch

	if err := bn.CheckModel(); err == nil {
		t.Error("expected error for evidence mismatch")
	}
}

func TestLinearGaussianBN_Copy(t *testing.T) {
	bn := buildLGChain(t)
	cp := bn.Copy()

	// Verify it's a distinct object with same structure.
	if cp == bn {
		t.Error("Copy should return a new object")
	}

	if err := cp.CheckModel(); err != nil {
		t.Fatalf("copied model CheckModel: %v", err)
	}

	// Verify nodes match.
	origNodes := bn.Nodes()
	copyNodes := cp.Nodes()
	if len(origNodes) != len(copyNodes) {
		t.Fatalf("node count mismatch: %d vs %d", len(origNodes), len(copyNodes))
	}
	for i := range origNodes {
		if origNodes[i] != copyNodes[i] {
			t.Errorf("node mismatch at %d: %q vs %q", i, origNodes[i], copyNodes[i])
		}
	}

	// Verify CPD parameters match.
	for _, node := range origNodes {
		origCPD := bn.GetLinearGaussianCPD(node)
		copyCPD := cp.GetLinearGaussianCPD(node)
		if origCPD.Mean() != copyCPD.Mean() {
			t.Errorf("mean mismatch for %q", node)
		}
		if origCPD.Variance() != copyCPD.Variance() {
			t.Errorf("variance mismatch for %q", node)
		}
	}

	// Mutating copy should not affect original.
	_ = cp.AddNode("W")
	if len(bn.Nodes()) != 3 {
		t.Error("mutating copy affected original")
	}
}

func TestLinearGaussianBN_NoParentCPD(t *testing.T) {
	bn := NewLinearGaussianBayesianNetwork()
	_ = bn.AddNode("A")

	cpd, _ := factors.NewLinearGaussianCPD("A", 3.0, nil, 2.0, nil)
	if err := bn.AddLinearGaussianCPD(cpd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := bn.CheckModel(); err != nil {
		t.Fatalf("CheckModel: %v", err)
	}

	got := bn.GetLinearGaussianCPD("A")
	if got.Mean() != 3.0 {
		t.Errorf("expected mean 3.0, got %f", got.Mean())
	}
	if got.Variance() != 2.0 {
		t.Errorf("expected variance 2.0, got %f", got.Variance())
	}
}
