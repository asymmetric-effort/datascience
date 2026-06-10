//go:build unit

package models

import (
	"fmt"
	"math"
	"testing"

	"github.com/asymmetric-effort/datascience/lib/pgm/factors"
	"github.com/asymmetric-effort/datascience/lib/tabgo"
)

// buildDiscreteStudentNetwork constructs the classic Student Bayesian network
// as a DiscreteBayesianNetwork.
func buildDiscreteStudentNetwork(t *testing.T) *DiscreteBayesianNetwork {
	t.Helper()
	dbn := NewDiscreteBayesianNetwork()

	for _, node := range []string{"D", "I", "G", "L", "S"} {
		if err := dbn.AddNode(node); err != nil {
			t.Fatalf("AddNode(%q): %v", node, err)
		}
	}

	edges := [][2]string{{"D", "G"}, {"I", "G"}, {"G", "L"}, {"I", "S"}}
	for _, e := range edges {
		if err := dbn.AddEdge(e[0], e[1]); err != nil {
			t.Fatalf("AddEdge(%q, %q): %v", e[0], e[1], err)
		}
	}

	cpdD, err := factors.NewTabularCPD("D", 2, [][]float64{
		{0.6},
		{0.4},
	}, nil, nil)
	if err != nil {
		t.Fatalf("NewTabularCPD(D): %v", err)
	}

	cpdI, err := factors.NewTabularCPD("I", 2, [][]float64{
		{0.7},
		{0.3},
	}, nil, nil)
	if err != nil {
		t.Fatalf("NewTabularCPD(I): %v", err)
	}

	cpdG, err := factors.NewTabularCPD("G", 3, [][]float64{
		{0.3, 0.05, 0.9, 0.5},
		{0.4, 0.25, 0.08, 0.3},
		{0.3, 0.70, 0.02, 0.2},
	}, []string{"D", "I"}, []int{2, 2})
	if err != nil {
		t.Fatalf("NewTabularCPD(G): %v", err)
	}

	cpdL, err := factors.NewTabularCPD("L", 2, [][]float64{
		{0.1, 0.4, 0.99},
		{0.9, 0.6, 0.01},
	}, []string{"G"}, []int{3})
	if err != nil {
		t.Fatalf("NewTabularCPD(L): %v", err)
	}

	cpdS, err := factors.NewTabularCPD("S", 2, [][]float64{
		{0.95, 0.2},
		{0.05, 0.8},
	}, []string{"I"}, []int{2})
	if err != nil {
		t.Fatalf("NewTabularCPD(S): %v", err)
	}

	for _, cpd := range []*factors.TabularCPD{cpdD, cpdI, cpdG, cpdL, cpdS} {
		if err := dbn.AddCPD(cpd); err != nil {
			t.Fatalf("AddCPD(%q): %v", cpd.Variable(), err)
		}
	}

	return dbn
}

func TestNewDiscreteBayesianNetwork(t *testing.T) {
	dbn := NewDiscreteBayesianNetwork()
	if dbn == nil {
		t.Fatal("NewDiscreteBayesianNetwork returned nil")
	}
	if dbn.BayesianNetwork == nil {
		t.Fatal("embedded BayesianNetwork is nil")
	}
	if len(dbn.Nodes()) != 0 {
		t.Errorf("expected 0 nodes, got %d", len(dbn.Nodes()))
	}
}

func TestDiscreteBayesianNetworkAddCPDNil(t *testing.T) {
	dbn := NewDiscreteBayesianNetwork()
	if err := dbn.AddCPD(nil); err == nil {
		t.Error("expected error for nil CPD")
	}
}

func TestDiscreteBayesianNetworkAddCPDValid(t *testing.T) {
	dbn := NewDiscreteBayesianNetwork()
	_ = dbn.AddNode("X")
	cpd, err := factors.NewTabularCPD("X", 2, [][]float64{{0.5}, {0.5}}, nil, nil)
	if err != nil {
		t.Fatalf("NewTabularCPD: %v", err)
	}
	if err := dbn.AddCPD(cpd); err != nil {
		t.Fatalf("AddCPD: %v", err)
	}
	if dbn.GetCPD("X") == nil {
		t.Error("expected CPD to be set")
	}
}

func TestDiscreteBayesianNetworkAddCPDUnknownNode(t *testing.T) {
	dbn := NewDiscreteBayesianNetwork()
	cpd, _ := factors.NewTabularCPD("Z", 2, [][]float64{{0.5}, {0.5}}, nil, nil)
	if err := dbn.AddCPD(cpd); err == nil {
		t.Error("expected error for CPD with unknown node")
	}
}

func TestDiscreteBayesianNetworkCheckModelValid(t *testing.T) {
	dbn := buildDiscreteStudentNetwork(t)
	if err := dbn.CheckModel(); err != nil {
		t.Fatalf("CheckModel on valid student network: %v", err)
	}
}

func TestDiscreteBayesianNetworkCheckModelMissingCPD(t *testing.T) {
	dbn := buildDiscreteStudentNetwork(t)
	dbn.RemoveCPD("G")
	if err := dbn.CheckModel(); err == nil {
		t.Error("expected error for missing CPD")
	}
}

func TestDiscreteBayesianNetworkCheckModelEvidenceMismatch(t *testing.T) {
	dbn := NewDiscreteBayesianNetwork()
	_ = dbn.AddNode("A")
	_ = dbn.AddNode("B")
	_ = dbn.AddEdge("A", "B")

	cpdA, _ := factors.NewTabularCPD("A", 2, [][]float64{{0.5}, {0.5}}, nil, nil)
	cpdB, _ := factors.NewTabularCPD("B", 2, [][]float64{{0.5}, {0.5}}, nil, nil)
	_ = dbn.BayesianNetwork.AddCPD(cpdA)
	_ = dbn.BayesianNetwork.AddCPD(cpdB)

	if err := dbn.CheckModel(); err == nil {
		t.Error("expected error for evidence/parent mismatch")
	}
}

func TestDiscreteBayesianNetworkCheckModelStateNamesMismatch(t *testing.T) {
	dbn := buildDiscreteStudentNetwork(t)

	// D has cardinality 2, set 3 state names to trigger mismatch.
	_ = dbn.SetStates("D", []string{"easy", "medium", "hard"})

	if err := dbn.CheckModel(); err == nil {
		t.Error("expected error for state names / cardinality mismatch")
	}
}

func TestDiscreteBayesianNetworkCheckModelStateNamesConsistent(t *testing.T) {
	dbn := buildDiscreteStudentNetwork(t)

	// Set consistent state names.
	_ = dbn.SetStates("D", []string{"easy", "hard"})
	_ = dbn.SetStates("I", []string{"low", "high"})
	_ = dbn.SetStates("G", []string{"A", "B", "C"})
	_ = dbn.SetStates("L", []string{"weak", "strong"})
	_ = dbn.SetStates("S", []string{"low", "high"})

	if err := dbn.CheckModel(); err != nil {
		t.Fatalf("CheckModel with consistent state names: %v", err)
	}
}

func TestDiscreteBayesianNetworkCheckModelParentStateNamesMismatch(t *testing.T) {
	dbn := buildDiscreteStudentNetwork(t)

	// I has cardinality 2. S's CPD uses I as evidence with card 2.
	// Set 3 state names for I to create a mismatch.
	_ = dbn.SetStates("I", []string{"low", "medium", "high"})

	if err := dbn.CheckModel(); err == nil {
		t.Error("expected error for parent state names / evidence cardinality mismatch")
	}
}

func TestDiscreteBayesianNetworkCheckModelInvalidCPD(t *testing.T) {
	dbn := NewDiscreteBayesianNetwork()
	_ = dbn.AddNode("X")

	// Columns don't sum to 1.
	cpd, _ := factors.NewTabularCPD("X", 2, [][]float64{{0.3}, {0.3}}, nil, nil)
	_ = dbn.BayesianNetwork.AddCPD(cpd)

	if err := dbn.CheckModel(); err == nil {
		t.Error("expected error for invalid CPD")
	}
}

func TestDiscreteBayesianNetworkCheckModelEmpty(t *testing.T) {
	dbn := NewDiscreteBayesianNetwork()
	if err := dbn.CheckModel(); err != nil {
		t.Errorf("empty network CheckModel should pass: %v", err)
	}
}

func TestDiscreteBayesianNetworkFitWithNilFn(t *testing.T) {
	dbn := NewDiscreteBayesianNetwork()
	df := tabgo.NewDataFrame(map[string]*tabgo.Series{})
	if err := dbn.FitWith(nil, df); err == nil {
		t.Error("expected error for nil estimateFn")
	}
}

func TestDiscreteBayesianNetworkFitWithNilData(t *testing.T) {
	dbn := NewDiscreteBayesianNetwork()
	fn := func(bn *BayesianNetwork, data *tabgo.DataFrame) error {
		return nil
	}
	if err := dbn.FitWith(fn, nil); err == nil {
		t.Error("expected error for nil data")
	}
}

func TestDiscreteBayesianNetworkFitWithSuccess(t *testing.T) {
	dbn := NewDiscreteBayesianNetwork()
	_ = dbn.AddNode("X")

	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", []any{0, 0, 1, 1, 1}),
	})

	// A mock estimator that sets a CPD based on the data.
	mockEstimator := func(bn *BayesianNetwork, df *tabgo.DataFrame) error {
		cpd, err := factors.NewTabularCPD("X", 2, [][]float64{{0.4}, {0.6}}, nil, nil)
		if err != nil {
			return err
		}
		return bn.AddCPD(cpd)
	}

	if err := dbn.FitWith(mockEstimator, data); err != nil {
		t.Fatalf("FitWith: %v", err)
	}

	cpd := dbn.GetCPD("X")
	if cpd == nil {
		t.Fatal("expected CPD to be set after FitWith")
	}
	if cpd.VariableCard() != 2 {
		t.Errorf("expected cardinality 2, got %d", cpd.VariableCard())
	}
}

func TestDiscreteBayesianNetworkFitWithError(t *testing.T) {
	dbn := NewDiscreteBayesianNetwork()
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{})

	errEstimator := func(bn *BayesianNetwork, df *tabgo.DataFrame) error {
		return fmt.Errorf("estimation failed")
	}

	if err := dbn.FitWith(errEstimator, data); err == nil {
		t.Error("expected error from failing estimator")
	}
}

func TestDiscreteBayesianNetworkSimulateValid(t *testing.T) {
	dbn := buildDiscreteStudentNetwork(t)

	n := 1000
	df, err := dbn.Simulate(n, 42)
	if err != nil {
		t.Fatalf("Simulate: %v", err)
	}

	if df.Len() != n {
		t.Errorf("expected %d rows, got %d", n, df.Len())
	}

	// Check that all expected columns exist.
	cols := df.Columns()
	expectedCols := []string{"D", "G", "I", "L", "S"}
	if len(cols) != len(expectedCols) {
		t.Fatalf("expected %d columns, got %d: %v", len(expectedCols), len(cols), cols)
	}
	for i, c := range cols {
		if c != expectedCols[i] {
			t.Errorf("column %d: expected %q, got %q", i, expectedCols[i], c)
		}
	}

	// Check that sampled values are within valid ranges.
	cardMap := map[string]int{"D": 2, "I": 2, "G": 3, "L": 2, "S": 2}
	for varName, card := range cardMap {
		vals := df.Column(varName).Values()
		for j, v := range vals {
			intVal, ok := v.(int)
			if !ok {
				t.Fatalf("sample %d for %q: expected int, got %T", j, varName, v)
			}
			if intVal < 0 || intVal >= card {
				t.Errorf("sample %d for %q: value %d out of range [0, %d)", j, varName, intVal, card)
			}
		}
	}
}

func TestDiscreteBayesianNetworkSimulateDeterministic(t *testing.T) {
	dbn := buildDiscreteStudentNetwork(t)

	df1, err := dbn.Simulate(100, 123)
	if err != nil {
		t.Fatalf("Simulate 1: %v", err)
	}

	df2, err := dbn.Simulate(100, 123)
	if err != nil {
		t.Fatalf("Simulate 2: %v", err)
	}

	// Same seed should produce identical results.
	for _, col := range df1.Columns() {
		v1 := df1.Column(col).Values()
		v2 := df2.Column(col).Values()
		for i := range v1 {
			if v1[i] != v2[i] {
				t.Errorf("sample %d for %q differs: %v vs %v", i, col, v1[i], v2[i])
			}
		}
	}
}

func TestDiscreteBayesianNetworkSimulateDistribution(t *testing.T) {
	dbn := buildDiscreteStudentNetwork(t)

	n := 10000
	df, err := dbn.Simulate(n, 99)
	if err != nil {
		t.Fatalf("Simulate: %v", err)
	}

	// Check that D's marginal distribution is approximately P(D=0)=0.6.
	dVals := df.Column("D").Values()
	count0 := 0
	for _, v := range dVals {
		if v.(int) == 0 {
			count0++
		}
	}
	p0 := float64(count0) / float64(n)
	if math.Abs(p0-0.6) > 0.05 {
		t.Errorf("P(D=0) = %.3f, expected ~0.6", p0)
	}

	// Check that I's marginal distribution is approximately P(I=0)=0.7.
	iVals := df.Column("I").Values()
	iCount0 := 0
	for _, v := range iVals {
		if v.(int) == 0 {
			iCount0++
		}
	}
	pI0 := float64(iCount0) / float64(n)
	if math.Abs(pI0-0.7) > 0.05 {
		t.Errorf("P(I=0) = %.3f, expected ~0.7", pI0)
	}
}

func TestDiscreteBayesianNetworkSimulateInvalidN(t *testing.T) {
	dbn := buildDiscreteStudentNetwork(t)

	_, err := dbn.Simulate(0, 42)
	if err == nil {
		t.Error("expected error for n=0")
	}

	_, err = dbn.Simulate(-5, 42)
	if err == nil {
		t.Error("expected error for negative n")
	}
}

func TestDiscreteBayesianNetworkSimulateInvalidModel(t *testing.T) {
	dbn := buildDiscreteStudentNetwork(t)
	dbn.RemoveCPD("G")

	_, err := dbn.Simulate(10, 42)
	if err == nil {
		t.Error("expected error for invalid model")
	}
}

func TestDiscreteBayesianNetworkCopy(t *testing.T) {
	dbn := buildDiscreteStudentNetwork(t)
	_ = dbn.SetStates("D", []string{"easy", "hard"})

	cp := dbn.Copy()

	// Verify the copy is valid.
	if err := cp.CheckModel(); err != nil {
		t.Fatalf("copied model CheckModel: %v", err)
	}

	// Verify structural equality.
	if len(cp.Nodes()) != len(dbn.Nodes()) {
		t.Errorf("copy has %d nodes, expected %d", len(cp.Nodes()), len(dbn.Nodes()))
	}
	if len(cp.Edges()) != len(dbn.Edges()) {
		t.Errorf("copy has %d edges, expected %d", len(cp.Edges()), len(dbn.Edges()))
	}

	// Verify state names are copied.
	states := cp.GetStates("D")
	if states == nil || len(states) != 2 {
		t.Fatalf("expected 2 state names for D in copy, got %v", states)
	}
	if states[0] != "easy" || states[1] != "hard" {
		t.Errorf("expected state names [easy hard], got %v", states)
	}
}

func TestDiscreteBayesianNetworkCopyIndependence(t *testing.T) {
	dbn := buildDiscreteStudentNetwork(t)
	cp := dbn.Copy()

	// Modify the copy: remove a CPD.
	cp.RemoveCPD("D")
	if dbn.GetCPD("D") == nil {
		t.Error("original CPD was affected by copy modification")
	}

	// Modify the original: remove a different CPD.
	dbn.RemoveCPD("S")
	if cp.GetCPD("S") == nil {
		t.Error("copy CPD was affected by original modification")
	}
}

func TestDiscreteBayesianNetworkCopyIndependenceStates(t *testing.T) {
	dbn := NewDiscreteBayesianNetwork()
	_ = dbn.AddNode("X")
	cpd, _ := factors.NewTabularCPD("X", 2, [][]float64{{0.5}, {0.5}}, nil, nil)
	_ = dbn.AddCPD(cpd)
	_ = dbn.SetStates("X", []string{"a", "b"})

	cp := dbn.Copy()

	// Modify states in original.
	_ = dbn.SetStates("X", []string{"c", "d"})

	origStates := dbn.GetStates("X")
	copyStates := cp.GetStates("X")

	if origStates[0] != "c" || origStates[1] != "d" {
		t.Errorf("original states: expected [c d], got %v", origStates)
	}
	if copyStates[0] != "a" || copyStates[1] != "b" {
		t.Errorf("copy states: expected [a b], got %v", copyStates)
	}
}

func TestDiscreteBayesianNetworkSimulateSingleNode(t *testing.T) {
	dbn := NewDiscreteBayesianNetwork()
	_ = dbn.AddNode("X")
	cpd, _ := factors.NewTabularCPD("X", 3, [][]float64{
		{0.2},
		{0.3},
		{0.5},
	}, nil, nil)
	_ = dbn.AddCPD(cpd)

	df, err := dbn.Simulate(5000, 77)
	if err != nil {
		t.Fatalf("Simulate: %v", err)
	}

	if df.Len() != 5000 {
		t.Errorf("expected 5000 rows, got %d", df.Len())
	}

	// Check approximate distribution.
	vals := df.Column("X").Values()
	counts := [3]int{}
	for _, v := range vals {
		counts[v.(int)]++
	}
	n := float64(5000)
	if math.Abs(float64(counts[0])/n-0.2) > 0.05 {
		t.Errorf("P(X=0) = %.3f, expected ~0.2", float64(counts[0])/n)
	}
	if math.Abs(float64(counts[1])/n-0.3) > 0.05 {
		t.Errorf("P(X=1) = %.3f, expected ~0.3", float64(counts[1])/n)
	}
	if math.Abs(float64(counts[2])/n-0.5) > 0.05 {
		t.Errorf("P(X=2) = %.3f, expected ~0.5", float64(counts[2])/n)
	}
}

func TestDiscreteBayesianNetworkEmbeddingAccess(t *testing.T) {
	dbn := buildDiscreteStudentNetwork(t)

	// Verify we can access BayesianNetwork methods directly.
	nodes := dbn.Nodes()
	if len(nodes) != 5 {
		t.Errorf("expected 5 nodes, got %d", len(nodes))
	}

	edges := dbn.Edges()
	if len(edges) != 4 {
		t.Errorf("expected 4 edges, got %d", len(edges))
	}

	parents := dbn.Parents("G")
	if len(parents) != 2 {
		t.Errorf("expected 2 parents for G, got %d", len(parents))
	}

	children := dbn.Children("I")
	if len(children) != 2 {
		t.Errorf("expected 2 children for I, got %d", len(children))
	}
}

func TestDiscreteBayesianNetworkToMarkovFactors(t *testing.T) {
	dbn := buildDiscreteStudentNetwork(t)

	mf, err := dbn.ToMarkovFactors()
	if err != nil {
		t.Fatalf("ToMarkovFactors: %v", err)
	}
	if len(mf) != 5 {
		t.Errorf("expected 5 factors, got %d", len(mf))
	}
}
