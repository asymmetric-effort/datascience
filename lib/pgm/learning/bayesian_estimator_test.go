//go:build unit

package learning

import (
	"math"
	"testing"

	"github.com/asymmetric-effort/datascience/lib/pgm/models"
	"github.com/asymmetric-effort/datascience/lib/tabgo"
)

// helper: build a simple A -> B network with states and data.
func setupSimpleNetwork(t *testing.T) (*models.BayesianNetwork, *tabgo.DataFrame) {
	t.Helper()
	bn := models.NewBayesianNetwork()
	if err := bn.AddNode("A"); err != nil {
		t.Fatal(err)
	}
	if err := bn.AddNode("B"); err != nil {
		t.Fatal(err)
	}
	if err := bn.AddEdge("A", "B"); err != nil {
		t.Fatal(err)
	}
	if err := bn.SetStates("A", []string{"a0", "a1"}); err != nil {
		t.Fatal(err)
	}
	if err := bn.SetStates("B", []string{"b0", "b1"}); err != nil {
		t.Fatal(err)
	}

	// Data: 100 rows. A=a0 60 times, A=a1 40 times.
	// When A=a0: B=b0 50 times, B=b1 10 times.
	// When A=a1: B=b0 10 times, B=b1 30 times.
	rows := make([][]any, 0, 100)
	for i := 0; i < 50; i++ {
		rows = append(rows, []any{"a0", "b0"})
	}
	for i := 0; i < 10; i++ {
		rows = append(rows, []any{"a0", "b1"})
	}
	for i := 0; i < 10; i++ {
		rows = append(rows, []any{"a1", "b0"})
	}
	for i := 0; i < 30; i++ {
		rows = append(rows, []any{"a1", "b1"})
	}

	df := tabgo.NewDataFrameFromRows([]string{"A", "B"}, rows)
	return bn, df
}

// helper: build a no-parent network (single node X with 3 states).
func setupSingleNode(t *testing.T) (*models.BayesianNetwork, *tabgo.DataFrame) {
	t.Helper()
	bn := models.NewBayesianNetwork()
	if err := bn.AddNode("X"); err != nil {
		t.Fatal(err)
	}
	if err := bn.SetStates("X", []string{"x0", "x1", "x2"}); err != nil {
		t.Fatal(err)
	}

	// x0: 70, x1: 20, x2: 10
	rows := make([][]any, 0, 100)
	for i := 0; i < 70; i++ {
		rows = append(rows, []any{"x0"})
	}
	for i := 0; i < 20; i++ {
		rows = append(rows, []any{"x1"})
	}
	for i := 0; i < 10; i++ {
		rows = append(rows, []any{"x2"})
	}
	df := tabgo.NewDataFrameFromRows([]string{"X"}, rows)
	return bn, df
}

func approxEqual(a, b, tol float64) bool {
	return math.Abs(a-b) < tol
}

func TestBayesianEstimator_BDeu_Basic(t *testing.T) {
	bn, df := setupSimpleNetwork(t)
	est := NewBayesianEstimator(bn, df, BDeu, 4.0)

	if err := est.Estimate(); err != nil {
		t.Fatalf("Estimate failed: %v", err)
	}

	// Node A (no parents): counts [60, 40], pseudo = 4/(1*2)=2 each => [62,42]/104
	cpdA, err := est.GetParameters("A")
	if err != nil {
		t.Fatalf("GetParameters(A) failed: %v", err)
	}
	fA := cpdA.ToFactor()
	valsA := fA.Values().Data()
	expectedA0 := 62.0 / 104.0
	expectedA1 := 42.0 / 104.0
	if !approxEqual(valsA[0], expectedA0, 1e-6) {
		t.Errorf("P(A=a0) = %f, want %f", valsA[0], expectedA0)
	}
	if !approxEqual(valsA[1], expectedA1, 1e-6) {
		t.Errorf("P(A=a1) = %f, want %f", valsA[1], expectedA1)
	}

	// Node B (parent A, 2 parent configs): pseudo = 4/(2*2)=1 each
	// parent config 0 (A=a0): counts [50,10]+1 = [51,11]/62
	// parent config 1 (A=a1): counts [10,30]+1 = [11,31]/42
	cpdB, err := est.GetParameters("B")
	if err != nil {
		t.Fatalf("GetParameters(B) failed: %v", err)
	}
	fB := cpdB.ToFactor()
	valsB := fB.Values().Data()
	// Layout: [b0|A=a0, b0|A=a1, b1|A=a0, b1|A=a1]
	if !approxEqual(valsB[0], 51.0/62.0, 1e-6) {
		t.Errorf("P(B=b0|A=a0) = %f, want %f", valsB[0], 51.0/62.0)
	}
	if !approxEqual(valsB[1], 11.0/42.0, 1e-6) {
		t.Errorf("P(B=b0|A=a1) = %f, want %f", valsB[1], 11.0/42.0)
	}
	if !approxEqual(valsB[2], 11.0/62.0, 1e-6) {
		t.Errorf("P(B=b1|A=a0) = %f, want %f", valsB[2], 11.0/62.0)
	}
	if !approxEqual(valsB[3], 31.0/42.0, 1e-6) {
		t.Errorf("P(B=b1|A=a1) = %f, want %f", valsB[3], 31.0/42.0)
	}
}

func TestBayesianEstimator_BDeu_LargeESS_ApproachesUniform(t *testing.T) {
	bn, df := setupSingleNode(t)
	// Very large ESS should make the pseudo-counts dominate and push toward uniform.
	est := NewBayesianEstimator(bn, df, BDeu, 1e8)
	if err := est.Estimate(); err != nil {
		t.Fatalf("Estimate failed: %v", err)
	}

	cpd, err := est.GetParameters("X")
	if err != nil {
		t.Fatal(err)
	}
	vals := cpd.ToFactor().Values().Data()
	uniform := 1.0 / 3.0
	for i, v := range vals {
		if !approxEqual(v, uniform, 1e-3) {
			t.Errorf("state %d: P = %f, want ~%f (uniform)", i, v, uniform)
		}
	}
}

func TestBayesianEstimator_BDeu_SmallESS_ApproachesMLE(t *testing.T) {
	bn, df := setupSingleNode(t)
	// Very small ESS: pseudo-counts negligible, should approach MLE.
	est := NewBayesianEstimator(bn, df, BDeu, 1e-10)
	if err := est.Estimate(); err != nil {
		t.Fatalf("Estimate failed: %v", err)
	}

	cpd, err := est.GetParameters("X")
	if err != nil {
		t.Fatal(err)
	}
	vals := cpd.ToFactor().Values().Data()
	// MLE: 70/100, 20/100, 10/100
	expected := []float64{0.70, 0.20, 0.10}
	for i, v := range vals {
		if !approxEqual(v, expected[i], 1e-3) {
			t.Errorf("state %d: P = %f, want ~%f (MLE)", i, v, expected[i])
		}
	}
}

func TestBayesianEstimator_K2(t *testing.T) {
	bn, df := setupSingleNode(t)
	est := NewBayesianEstimator(bn, df, K2, 0) // ESS unused for K2

	if err := est.Estimate(); err != nil {
		t.Fatalf("Estimate failed: %v", err)
	}

	cpd, err := est.GetParameters("X")
	if err != nil {
		t.Fatal(err)
	}
	vals := cpd.ToFactor().Values().Data()
	// K2: pseudo_count=1 => counts [70+1, 20+1, 10+1] = [71,21,11]/103
	total := 103.0
	expected := []float64{71.0 / total, 21.0 / total, 11.0 / total}
	for i, v := range vals {
		if !approxEqual(v, expected[i], 1e-6) {
			t.Errorf("K2 state %d: P = %f, want %f", i, v, expected[i])
		}
	}
}

func TestBayesianEstimator_UniformPrior(t *testing.T) {
	bn, df := setupSingleNode(t)
	est := NewBayesianEstimator(bn, df, UniformPrior, 0) // ESS unused

	if err := est.Estimate(); err != nil {
		t.Fatalf("Estimate failed: %v", err)
	}

	cpd, err := est.GetParameters("X")
	if err != nil {
		t.Fatal(err)
	}
	vals := cpd.ToFactor().Values().Data()
	// Uniform: pseudo_count = 1/3 => counts [70+1/3, 20+1/3, 10+1/3] = [70.333, 20.333, 10.333]/101
	total := 101.0
	expected := []float64{(70.0 + 1.0/3.0) / total, (20.0 + 1.0/3.0) / total, (10.0 + 1.0/3.0) / total}
	for i, v := range vals {
		if !approxEqual(v, expected[i], 1e-6) {
			t.Errorf("Uniform state %d: P = %f, want %f", i, v, expected[i])
		}
	}
}

func TestBayesianEstimator_K2_vs_Uniform_Differ(t *testing.T) {
	bn1, df1 := setupSingleNode(t)
	bn2, df2 := setupSingleNode(t)

	estK2 := NewBayesianEstimator(bn1, df1, K2, 0)
	estUni := NewBayesianEstimator(bn2, df2, UniformPrior, 0)

	if err := estK2.Estimate(); err != nil {
		t.Fatal(err)
	}
	if err := estUni.Estimate(); err != nil {
		t.Fatal(err)
	}

	cpdK2, _ := estK2.GetParameters("X")
	cpdUni, _ := estUni.GetParameters("X")

	vK2 := cpdK2.ToFactor().Values().Data()
	vUni := cpdUni.ToFactor().Values().Data()

	differ := false
	for i := range vK2 {
		if !approxEqual(vK2[i], vUni[i], 1e-8) {
			differ = true
			break
		}
	}
	if !differ {
		t.Error("K2 and Uniform priors produced identical results; they should differ")
	}
}

func TestBayesianEstimator_ColumnsNormalize(t *testing.T) {
	bn, df := setupSimpleNetwork(t)
	est := NewBayesianEstimator(bn, df, BDeu, 4.0)
	if err := est.Estimate(); err != nil {
		t.Fatal(err)
	}

	// Validate all CPDs pass the built-in validation (columns sum to 1).
	if err := bn.CheckModel(); err != nil {
		t.Fatalf("CheckModel failed after estimation: %v", err)
	}
}

func TestBayesianEstimator_GetParameters_NoEstimate(t *testing.T) {
	bn, df := setupSimpleNetwork(t)
	est := NewBayesianEstimator(bn, df, BDeu, 1.0)

	_, err := est.GetParameters("A")
	if err == nil {
		t.Error("expected error when calling GetParameters before Estimate")
	}
}

func TestBayesianEstimator_MissingStates(t *testing.T) {
	bn := models.NewBayesianNetwork()
	_ = bn.AddNode("A")
	// No states set for A.

	rows := [][]any{{"a0"}, {"a1"}}
	df := tabgo.NewDataFrameFromRows([]string{"A"}, rows)

	est := NewBayesianEstimator(bn, df, K2, 0)
	err := est.Estimate()
	if err == nil {
		t.Error("expected error when node has no states defined")
	}
}

func TestBayesianEstimator_MultipleParents(t *testing.T) {
	bn := models.NewBayesianNetwork()
	for _, n := range []string{"A", "B", "C"} {
		if err := bn.AddNode(n); err != nil {
			t.Fatal(err)
		}
	}
	_ = bn.AddEdge("A", "C")
	_ = bn.AddEdge("B", "C")
	_ = bn.SetStates("A", []string{"a0", "a1"})
	_ = bn.SetStates("B", []string{"b0", "b1"})
	_ = bn.SetStates("C", []string{"c0", "c1"})

	// Build data: deterministic C = XOR(A,B)
	rows := make([][]any, 0, 40)
	for i := 0; i < 10; i++ {
		rows = append(rows, []any{"a0", "b0", "c0"}) // 0 xor 0 = 0
	}
	for i := 0; i < 10; i++ {
		rows = append(rows, []any{"a0", "b1", "c1"}) // 0 xor 1 = 1
	}
	for i := 0; i < 10; i++ {
		rows = append(rows, []any{"a1", "b0", "c1"}) // 1 xor 0 = 1
	}
	for i := 0; i < 10; i++ {
		rows = append(rows, []any{"a1", "b1", "c0"}) // 1 xor 1 = 0
	}
	df := tabgo.NewDataFrameFromRows([]string{"A", "B", "C"}, rows)

	est := NewBayesianEstimator(bn, df, K2, 0)
	if err := est.Estimate(); err != nil {
		t.Fatalf("Estimate failed: %v", err)
	}

	if err := bn.CheckModel(); err != nil {
		t.Fatalf("CheckModel failed: %v", err)
	}

	cpdC, _ := est.GetParameters("C")
	fC := cpdC.ToFactor()
	vals := fC.Values().Data()

	// Parents sorted: A, B. Parent configs in row-major: (a0,b0), (a0,b1), (a1,b0), (a1,b1)
	// C cardinality 2, parent configs 4.
	// Layout: [c0|(a0,b0), c0|(a0,b1), c0|(a1,b0), c0|(a1,b1), c1|(a0,b0), c1|(a0,b1), c1|(a1,b0), c1|(a1,b1)]
	// With K2 (pseudo=1): counts + 1, then normalize
	// (a0,b0): c0=10+1=11, c1=0+1=1 => 11/12, 1/12
	// (a0,b1): c0=0+1=1, c1=10+1=11 => 1/12, 11/12
	// (a1,b0): c0=0+1=1, c1=10+1=11 => 1/12, 11/12
	// (a1,b1): c0=10+1=11, c1=0+1=1 => 11/12, 1/12

	expected := []float64{
		11.0 / 12.0, 1.0 / 12.0, 1.0 / 12.0, 11.0 / 12.0, // c0 row
		1.0 / 12.0, 11.0 / 12.0, 11.0 / 12.0, 1.0 / 12.0, // c1 row
	}
	for i, v := range vals {
		if !approxEqual(v, expected[i], 1e-6) {
			t.Errorf("C vals[%d] = %f, want %f", i, v, expected[i])
		}
	}
}

func TestBayesianEstimator_PriorType_Constants(t *testing.T) {
	// Verify the three prior type constants have distinct values.
	if BDeu == K2 || BDeu == UniformPrior || K2 == UniformPrior {
		t.Error("PriorType constants must be distinct")
	}
}
