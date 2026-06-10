//go:build unit

package models

import (
	"testing"

	"github.com/asymmetric-effort/datascience/lib/pgm/factors"
)

// buildStudentNetwork constructs the classic Student Bayesian network:
//
//	D -> G <- I
//	G -> L
//	I -> S
//
// Variables and cardinalities:
//
//	D (Difficulty):   2 states {easy=0, hard=1}
//	I (Intelligence): 2 states {low=0, high=1}
//	G (Grade):        3 states {A=0, B=1, C=2}
//	L (Letter):       2 states {weak=0, strong=1}
//	S (SAT):          2 states {low=0, high=1}
func buildStudentNetwork(t *testing.T) *BayesianNetwork {
	t.Helper()
	bn := NewBayesianNetwork()

	// Add nodes.
	for _, node := range []string{"D", "I", "G", "L", "S"} {
		if err := bn.AddNode(node); err != nil {
			t.Fatalf("AddNode(%q): %v", node, err)
		}
	}

	// Add edges.
	edges := [][2]string{{"D", "G"}, {"I", "G"}, {"G", "L"}, {"I", "S"}}
	for _, e := range edges {
		if err := bn.AddEdge(e[0], e[1]); err != nil {
			t.Fatalf("AddEdge(%q, %q): %v", e[0], e[1], err)
		}
	}

	// CPD for D (no parents): P(D=easy)=0.6, P(D=hard)=0.4
	cpdD, err := factors.NewTabularCPD("D", 2, [][]float64{
		{0.6},
		{0.4},
	}, nil, nil)
	if err != nil {
		t.Fatalf("NewTabularCPD(D): %v", err)
	}

	// CPD for I (no parents): P(I=low)=0.7, P(I=high)=0.3
	cpdI, err := factors.NewTabularCPD("I", 2, [][]float64{
		{0.7},
		{0.3},
	}, nil, nil)
	if err != nil {
		t.Fatalf("NewTabularCPD(I): %v", err)
	}

	// CPD for G (parents: D, I):
	// Columns: (D=0,I=0), (D=0,I=1), (D=1,I=0), (D=1,I=1)
	cpdG, err := factors.NewTabularCPD("G", 3, [][]float64{
		{0.3, 0.05, 0.9, 0.5},
		{0.4, 0.25, 0.08, 0.3},
		{0.3, 0.70, 0.02, 0.2},
	}, []string{"D", "I"}, []int{2, 2})
	if err != nil {
		t.Fatalf("NewTabularCPD(G): %v", err)
	}

	// CPD for L (parent: G):
	// Columns: G=0(A), G=1(B), G=2(C)
	cpdL, err := factors.NewTabularCPD("L", 2, [][]float64{
		{0.1, 0.4, 0.99},
		{0.9, 0.6, 0.01},
	}, []string{"G"}, []int{3})
	if err != nil {
		t.Fatalf("NewTabularCPD(L): %v", err)
	}

	// CPD for S (parent: I):
	// Columns: I=0(low), I=1(high)
	cpdS, err := factors.NewTabularCPD("S", 2, [][]float64{
		{0.95, 0.2},
		{0.05, 0.8},
	}, []string{"I"}, []int{2})
	if err != nil {
		t.Fatalf("NewTabularCPD(S): %v", err)
	}

	for _, cpd := range []*factors.TabularCPD{cpdD, cpdI, cpdG, cpdL, cpdS} {
		if err := bn.AddCPD(cpd); err != nil {
			t.Fatalf("AddCPD(%q): %v", cpd.Variable(), err)
		}
	}

	return bn
}

func TestNewBayesianNetwork(t *testing.T) {
	bn := NewBayesianNetwork()
	if bn == nil {
		t.Fatal("NewBayesianNetwork returned nil")
	}
	if len(bn.Nodes()) != 0 {
		t.Errorf("expected 0 nodes, got %d", len(bn.Nodes()))
	}
	if len(bn.Edges()) != 0 {
		t.Errorf("expected 0 edges, got %d", len(bn.Edges()))
	}
}

func TestAddNode(t *testing.T) {
	bn := NewBayesianNetwork()
	if err := bn.AddNode("A"); err != nil {
		t.Fatalf("AddNode: %v", err)
	}
	nodes := bn.Nodes()
	if len(nodes) != 1 || nodes[0] != "A" {
		t.Errorf("expected [A], got %v", nodes)
	}

	// Duplicate node should fail.
	if err := bn.AddNode("A"); err == nil {
		t.Error("expected error for duplicate node")
	}
}

func TestAddEdge(t *testing.T) {
	bn := NewBayesianNetwork()
	_ = bn.AddNode("A")
	_ = bn.AddNode("B")

	if err := bn.AddEdge("A", "B"); err != nil {
		t.Fatalf("AddEdge: %v", err)
	}

	edges := bn.Edges()
	if len(edges) != 1 || edges[0] != [2]string{"A", "B"} {
		t.Errorf("expected [[A B]], got %v", edges)
	}
}

func TestAddEdgeCycleDetection(t *testing.T) {
	bn := NewBayesianNetwork()
	_ = bn.AddNode("A")
	_ = bn.AddNode("B")
	_ = bn.AddNode("C")
	_ = bn.AddEdge("A", "B")
	_ = bn.AddEdge("B", "C")

	// C -> A would create a cycle.
	if err := bn.AddEdge("C", "A"); err == nil {
		t.Error("expected error for cycle-creating edge")
	}
}

func TestAddEdgeNonexistentNode(t *testing.T) {
	bn := NewBayesianNetwork()
	_ = bn.AddNode("A")

	if err := bn.AddEdge("A", "Z"); err == nil {
		t.Error("expected error for nonexistent destination node")
	}
	if err := bn.AddEdge("Z", "A"); err == nil {
		t.Error("expected error for nonexistent source node")
	}
}

func TestParentsChildren(t *testing.T) {
	bn := buildStudentNetwork(t)

	// G has parents D and I.
	parents := bn.Parents("G")
	if len(parents) != 2 || parents[0] != "D" || parents[1] != "I" {
		t.Errorf("expected parents [D I], got %v", parents)
	}

	// D has child G.
	children := bn.Children("D")
	if len(children) != 1 || children[0] != "G" {
		t.Errorf("expected children [G], got %v", children)
	}

	// I has children G and S.
	children = bn.Children("I")
	if len(children) != 2 || children[0] != "G" || children[1] != "S" {
		t.Errorf("expected children [G S], got %v", children)
	}

	// D has no parents.
	parents = bn.Parents("D")
	if len(parents) != 0 {
		t.Errorf("expected no parents for D, got %v", parents)
	}
}

func TestStudentNetworkStructure(t *testing.T) {
	bn := buildStudentNetwork(t)

	nodes := bn.Nodes()
	if len(nodes) != 5 {
		t.Fatalf("expected 5 nodes, got %d", len(nodes))
	}

	edges := bn.Edges()
	if len(edges) != 4 {
		t.Fatalf("expected 4 edges, got %d", len(edges))
	}
}

func TestAddCPD(t *testing.T) {
	bn := NewBayesianNetwork()
	_ = bn.AddNode("X")

	cpd, err := factors.NewTabularCPD("X", 2, [][]float64{{0.5}, {0.5}}, nil, nil)
	if err != nil {
		t.Fatalf("NewTabularCPD: %v", err)
	}

	if err := bn.AddCPD(cpd); err != nil {
		t.Fatalf("AddCPD: %v", err)
	}

	got := bn.GetCPD("X")
	if got == nil {
		t.Fatal("GetCPD returned nil")
	}
	if got.Variable() != "X" {
		t.Errorf("expected variable X, got %s", got.Variable())
	}
}

func TestAddCPDNil(t *testing.T) {
	bn := NewBayesianNetwork()
	if err := bn.AddCPD(nil); err == nil {
		t.Error("expected error for nil CPD")
	}
}

func TestAddCPDUnknownNode(t *testing.T) {
	bn := NewBayesianNetwork()
	cpd, _ := factors.NewTabularCPD("Z", 2, [][]float64{{0.5}, {0.5}}, nil, nil)
	if err := bn.AddCPD(cpd); err == nil {
		t.Error("expected error for CPD with unknown node")
	}
}

func TestRemoveCPD(t *testing.T) {
	bn := NewBayesianNetwork()
	_ = bn.AddNode("X")
	cpd, _ := factors.NewTabularCPD("X", 2, [][]float64{{0.5}, {0.5}}, nil, nil)
	_ = bn.AddCPD(cpd)

	bn.RemoveCPD("X")
	if bn.GetCPD("X") != nil {
		t.Error("expected nil after RemoveCPD")
	}

	// Removing a nonexistent CPD should not panic.
	bn.RemoveCPD("nonexistent")
}

func TestGetCPDs(t *testing.T) {
	bn := buildStudentNetwork(t)
	cpds := bn.GetCPDs()
	if len(cpds) != 5 {
		t.Fatalf("expected 5 CPDs, got %d", len(cpds))
	}

	// Should be sorted by variable name.
	expected := []string{"D", "G", "I", "L", "S"}
	for i, cpd := range cpds {
		if cpd.Variable() != expected[i] {
			t.Errorf("cpd[%d]: expected variable %q, got %q", i, expected[i], cpd.Variable())
		}
	}
}

func TestCheckModel(t *testing.T) {
	bn := buildStudentNetwork(t)
	if err := bn.CheckModel(); err != nil {
		t.Fatalf("CheckModel: %v", err)
	}
}

func TestCheckModelMissingCPD(t *testing.T) {
	bn := buildStudentNetwork(t)
	bn.RemoveCPD("G")
	if err := bn.CheckModel(); err == nil {
		t.Error("expected error for missing CPD")
	}
}

func TestCheckModelEvidenceMismatch(t *testing.T) {
	bn := NewBayesianNetwork()
	_ = bn.AddNode("A")
	_ = bn.AddNode("B")
	_ = bn.AddEdge("A", "B")

	// B has parent A, but we give it a CPD with no evidence.
	cpdA, _ := factors.NewTabularCPD("A", 2, [][]float64{{0.5}, {0.5}}, nil, nil)
	cpdB, _ := factors.NewTabularCPD("B", 2, [][]float64{{0.5}, {0.5}}, nil, nil)
	_ = bn.AddCPD(cpdA)
	_ = bn.AddCPD(cpdB)

	if err := bn.CheckModel(); err == nil {
		t.Error("expected error for evidence/parent mismatch")
	}
}

func TestCheckModelWrongEvidence(t *testing.T) {
	bn := NewBayesianNetwork()
	_ = bn.AddNode("A")
	_ = bn.AddNode("B")
	_ = bn.AddNode("C")
	_ = bn.AddEdge("A", "C")

	// C has parent A, but CPD says evidence is B.
	cpdA, _ := factors.NewTabularCPD("A", 2, [][]float64{{0.5}, {0.5}}, nil, nil)
	cpdB, _ := factors.NewTabularCPD("B", 2, [][]float64{{0.5}, {0.5}}, nil, nil)
	cpdC, _ := factors.NewTabularCPD("C", 2, [][]float64{
		{0.3, 0.7},
		{0.7, 0.3},
	}, []string{"B"}, []int{2})
	_ = bn.AddCPD(cpdA)
	_ = bn.AddCPD(cpdB)
	_ = bn.AddCPD(cpdC)

	if err := bn.CheckModel(); err == nil {
		t.Error("expected error for wrong evidence variables")
	}
}

func TestCheckModelInvalidCPD(t *testing.T) {
	bn := NewBayesianNetwork()
	_ = bn.AddNode("X")

	// CPD columns don't sum to 1.
	cpd, _ := factors.NewTabularCPD("X", 2, [][]float64{{0.3}, {0.3}}, nil, nil)
	_ = bn.AddCPD(cpd)

	if err := bn.CheckModel(); err == nil {
		t.Error("expected error for invalid CPD")
	}
}

func TestCopy(t *testing.T) {
	bn := buildStudentNetwork(t)
	bn2 := bn.Copy()

	// Verify the copy is independent.
	if err := bn2.CheckModel(); err != nil {
		t.Fatalf("copied model CheckModel: %v", err)
	}

	// Modify the copy and ensure the original is unaffected.
	bn2.RemoveCPD("D")
	if bn.GetCPD("D") == nil {
		t.Error("original CPD was affected by copy modification")
	}

	// Check structural equality.
	if len(bn2.Nodes()) != len(bn.Nodes()) {
		t.Errorf("copy has %d nodes, expected %d", len(bn2.Nodes()), len(bn.Nodes()))
	}
	if len(bn2.Edges()) != len(bn.Edges()) {
		t.Errorf("copy has %d edges, expected %d", len(bn2.Edges()), len(bn.Edges()))
	}
}

func TestToMarkovFactors(t *testing.T) {
	bn := buildStudentNetwork(t)

	markovFactors, err := bn.ToMarkovFactors()
	if err != nil {
		t.Fatalf("ToMarkovFactors: %v", err)
	}

	if len(markovFactors) != 5 {
		t.Fatalf("expected 5 factors, got %d", len(markovFactors))
	}

	// Factors should be sorted by variable name (same as Nodes() order).
	// Each factor's first variable should match the CPD's variable.
	expectedVars := []string{"D", "G", "I", "L", "S"}
	for i, f := range markovFactors {
		vars := f.Variables()
		if vars[0] != expectedVars[i] {
			t.Errorf("factor[%d] first variable: expected %q, got %q", i, expectedVars[i], vars[0])
		}
	}

	// Check that factor for D has 2 values summing to 1.
	dFactor := markovFactors[0]
	dVals := dFactor.Values().Data()
	sum := 0.0
	for _, v := range dVals {
		sum += v
	}
	if sum < 0.999 || sum > 1.001 {
		t.Errorf("D factor values sum to %f, expected 1.0", sum)
	}

	// Check that factor for G has variables [G, D, I] and 3*2*2=12 values.
	gFactor := markovFactors[1]
	gVars := gFactor.Variables()
	if len(gVars) != 3 {
		t.Errorf("G factor has %d variables, expected 3", len(gVars))
	}
	gCard := gFactor.Cardinality()
	gSize := 1
	for _, c := range gCard {
		gSize *= c
	}
	if gSize != 12 {
		t.Errorf("G factor size: expected 12, got %d", gSize)
	}
}

func TestToMarkovFactorsInvalidModel(t *testing.T) {
	bn := buildStudentNetwork(t)
	bn.RemoveCPD("S")

	_, err := bn.ToMarkovFactors()
	if err == nil {
		t.Error("expected error from ToMarkovFactors with invalid model")
	}
}

func TestEmptyNetworkCheckModel(t *testing.T) {
	bn := NewBayesianNetwork()
	// Empty network is valid (no nodes, no CPDs needed).
	if err := bn.CheckModel(); err != nil {
		t.Errorf("empty network CheckModel should pass: %v", err)
	}
}

func TestSingleNodeNetwork(t *testing.T) {
	bn := NewBayesianNetwork()
	_ = bn.AddNode("X")
	cpd, _ := factors.NewTabularCPD("X", 3, [][]float64{
		{0.2},
		{0.3},
		{0.5},
	}, nil, nil)
	_ = bn.AddCPD(cpd)

	if err := bn.CheckModel(); err != nil {
		t.Fatalf("single node CheckModel: %v", err)
	}

	markovFactors, err := bn.ToMarkovFactors()
	if err != nil {
		t.Fatalf("ToMarkovFactors: %v", err)
	}
	if len(markovFactors) != 1 {
		t.Errorf("expected 1 factor, got %d", len(markovFactors))
	}
}

func TestAddCPDOverwrite(t *testing.T) {
	bn := NewBayesianNetwork()
	_ = bn.AddNode("X")

	cpd1, _ := factors.NewTabularCPD("X", 2, [][]float64{{0.3}, {0.7}}, nil, nil)
	_ = bn.AddCPD(cpd1)

	cpd2, _ := factors.NewTabularCPD("X", 2, [][]float64{{0.9}, {0.1}}, nil, nil)
	_ = bn.AddCPD(cpd2)

	got := bn.GetCPD("X")
	// Should reflect the second CPD's values.
	f := got.ToFactor()
	vals := f.Values().Data()
	if vals[0] != 0.9 || vals[1] != 0.1 {
		t.Errorf("expected overwritten CPD values [0.9 0.1], got %v", vals)
	}
}

func TestGetCPDNonexistent(t *testing.T) {
	bn := NewBayesianNetwork()
	if bn.GetCPD("nonexistent") != nil {
		t.Error("expected nil for nonexistent CPD")
	}
}
