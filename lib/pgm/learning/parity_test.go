//go:build unit

package learning

import (
	"math"
	"testing"

	"github.com/asymmetric-effort/datascience/lib/graphgo"
	"github.com/asymmetric-effort/datascience/lib/pgm/models"
	"github.com/asymmetric-effort/datascience/lib/tabgo"
)

// --- MLE.EstimateCPD tests ---

func TestMLE_EstimateCPD_SingleNode(t *testing.T) {
	bn := models.NewBayesianNetwork()
	for _, n := range []string{"A", "B"} {
		if err := bn.AddNode(n); err != nil {
			t.Fatal(err)
		}
	}
	if err := bn.AddEdge("A", "B"); err != nil {
		t.Fatal(err)
	}

	df := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 0, 0, 1, 1, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 0, 1, 1}),
	})

	mle := NewMLE(bn, df)

	cpd, err := mle.EstimateCPD("B")
	if err != nil {
		t.Fatalf("EstimateCPD: %v", err)
	}
	if cpd == nil {
		t.Fatal("EstimateCPD returned nil")
	}
	if cpd.Variable() != "B" {
		t.Errorf("expected variable B, got %s", cpd.Variable())
	}
}

func TestMLE_EstimateCPD_NonexistentNode(t *testing.T) {
	bn := models.NewBayesianNetwork()
	if err := bn.AddNode("A"); err != nil {
		t.Fatal(err)
	}

	df := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1}),
	})

	mle := NewMLE(bn, df)
	_, err := mle.EstimateCPD("MISSING")
	if err == nil {
		t.Error("expected error for nonexistent node")
	}
}

func TestMLE_EstimateCPD_MatchesGetParameters(t *testing.T) {
	bn := models.NewBayesianNetwork()
	if err := bn.AddNode("X"); err != nil {
		t.Fatal(err)
	}

	df := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", []any{0, 0, 1, 1, 1}),
	})

	mle := NewMLE(bn, df)

	cpd1, err := mle.EstimateCPD("X")
	if err != nil {
		t.Fatal(err)
	}

	cpd2, err := mle.GetParameters("X")
	if err != nil {
		t.Fatal(err)
	}

	vals1 := cpd1.ToFactor().Values().Data()
	vals2 := cpd2.ToFactor().Values().Data()
	if len(vals1) != len(vals2) {
		t.Fatalf("different lengths: %d vs %d", len(vals1), len(vals2))
	}
	for i := range vals1 {
		if math.Abs(vals1[i]-vals2[i]) > 1e-10 {
			t.Errorf("value[%d]: %f vs %f", i, vals1[i], vals2[i])
		}
	}
}

// --- BayesianEstimator.EstimateCPD tests ---

func TestBayesianEstimator_EstimateCPD(t *testing.T) {
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

	est := NewBayesianEstimator(bn, df, K2, 0)
	cpd, err := est.EstimateCPD("B")
	if err != nil {
		t.Fatalf("EstimateCPD: %v", err)
	}
	if cpd == nil {
		t.Fatal("EstimateCPD returned nil")
	}
	if cpd.Variable() != "B" {
		t.Errorf("expected variable B, got %s", cpd.Variable())
	}
}

func TestBayesianEstimator_EstimateCPD_NonexistentNode(t *testing.T) {
	bn := models.NewBayesianNetwork()
	if err := bn.AddNode("A"); err != nil {
		t.Fatal(err)
	}
	if err := bn.SetStates("A", []string{"a0", "a1"}); err != nil {
		t.Fatal(err)
	}

	df := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{"a0", "a1"}),
	})

	est := NewBayesianEstimator(bn, df, K2, 0)
	_, err := est.EstimateCPD("MISSING")
	if err == nil {
		t.Error("expected error for nonexistent node")
	}
}

// --- EM.GetParameters tests ---

func TestEM_GetParameters_BeforeEstimate(t *testing.T) {
	bn := models.NewBayesianNetwork()
	if err := bn.AddNode("A"); err != nil {
		t.Fatal(err)
	}
	if err := bn.SetStates("A", []string{"a0", "a1"}); err != nil {
		t.Fatal(err)
	}

	df := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{"a0", "a1"}),
	})

	em := NewEM(bn, df, nil, 10, 1e-4)

	// Before Estimate(), GetParameters should return error (no CPDs).
	_, err := em.GetParameters()
	if err == nil {
		t.Error("expected error when calling GetParameters before Estimate")
	}
}

func TestEM_GetParameters_AfterEstimate(t *testing.T) {
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

	rows := make([][]any, 0, 40)
	for i := 0; i < 20; i++ {
		rows = append(rows, []any{"a0", "b0"})
	}
	for i := 0; i < 20; i++ {
		rows = append(rows, []any{"a1", "b1"})
	}
	df := tabgo.NewDataFrameFromRows([]string{"A", "B"}, rows)

	em := NewEM(bn, df, nil, 10, 1e-4)
	if err := em.Estimate(); err != nil {
		t.Fatalf("Estimate: %v", err)
	}

	params, err := em.GetParameters()
	if err != nil {
		t.Fatalf("GetParameters: %v", err)
	}

	if len(params) != 2 {
		t.Fatalf("expected 2 CPDs, got %d", len(params))
	}
	if params["A"] == nil {
		t.Error("missing CPD for A")
	}
	if params["B"] == nil {
		t.Error("missing CPD for B")
	}
}

func TestEM_GetParameters_NilBN(t *testing.T) {
	em := &ExpectationMaximization{bn: nil}
	_, err := em.GetParameters()
	if err == nil {
		t.Error("expected error for nil BN")
	}
}

// --- SEMEstimator.Fit and GetParameters tests ---

func TestSEMEstimator_Fit(t *testing.T) {
	sem := models.NewSEM()
	if err := sem.AddEquation("X", nil, nil, 0, 1); err != nil {
		t.Fatal(err)
	}
	if err := sem.AddEquation("Y", []string{"X"}, []float64{0}, 0, 1); err != nil {
		t.Fatal(err)
	}

	// Y = 2*X + 1 + noise
	n := 100
	xVals := make([]any, n)
	yVals := make([]any, n)
	for i := 0; i < n; i++ {
		x := float64(i) / float64(n)
		xVals[i] = x
		yVals[i] = 2.0*x + 1.0
	}
	df := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", xVals),
		"Y": tabgo.NewSeries("Y", yVals),
	})

	se := NewSEMEstimator(sem, nil) // initially nil data
	err := se.Fit(df)               // Fit provides data
	if err != nil {
		t.Fatalf("Fit: %v", err)
	}

	betas, intercept, _, err := se.GetCoefficients("Y")
	if err != nil {
		t.Fatalf("GetCoefficients: %v", err)
	}
	if math.Abs(intercept-1.0) > 0.1 {
		t.Errorf("expected intercept ~1.0, got %f", intercept)
	}
	if len(betas) != 1 || math.Abs(betas[0]-2.0) > 0.1 {
		t.Errorf("expected beta ~2.0, got %v", betas)
	}
}

func TestSEMEstimator_GetParameters(t *testing.T) {
	sem := models.NewSEM()
	if err := sem.AddEquation("X", nil, nil, 0, 1); err != nil {
		t.Fatal(err)
	}
	if err := sem.AddEquation("Y", []string{"X"}, []float64{0}, 0, 1); err != nil {
		t.Fatal(err)
	}

	n := 50
	xVals := make([]any, n)
	yVals := make([]any, n)
	for i := 0; i < n; i++ {
		x := float64(i)
		xVals[i] = x
		yVals[i] = 3.0*x + 5.0
	}
	df := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", xVals),
		"Y": tabgo.NewSeries("Y", yVals),
	})

	se := NewSEMEstimator(sem, df)
	if err := se.Estimate(); err != nil {
		t.Fatal(err)
	}

	params, err := se.GetParameters()
	if err != nil {
		t.Fatalf("GetParameters: %v", err)
	}

	if len(params) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(params))
	}

	yParams := params["Y"]
	if yParams == nil {
		t.Fatal("missing parameters for Y")
	}
	// yParams should be [coefficient, intercept, variance]
	if len(yParams) != 3 {
		t.Fatalf("expected 3 params for Y, got %d", len(yParams))
	}
	// Coefficient should be ~3.0.
	if math.Abs(yParams[0]-3.0) > 0.1 {
		t.Errorf("expected coefficient ~3.0, got %f", yParams[0])
	}
	// Intercept should be ~5.0.
	if math.Abs(yParams[1]-5.0) > 0.1 {
		t.Errorf("expected intercept ~5.0, got %f", yParams[1])
	}
}

func TestSEMEstimator_GetParameters_NilSEM(t *testing.T) {
	se := &SEMEstimator{sem: nil}
	_, err := se.GetParameters()
	if err == nil {
		t.Error("expected error for nil SEM")
	}
}

// --- IVEstimator.Estimate tests ---

func TestIVEstimator_Estimate(t *testing.T) {
	// Simple IV scenario: Z -> X -> Y with Z as instrument.
	n := 200
	zVals := make([]any, n)
	xVals := make([]any, n)
	yVals := make([]any, n)

	for i := 0; i < n; i++ {
		z := float64(i) / float64(n)
		x := 2.0*z + 0.5
		y := 3.0*x + 1.0
		zVals[i] = z
		xVals[i] = x
		yVals[i] = y
	}

	df := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"Z": tabgo.NewSeries("Z", zVals),
		"X": tabgo.NewSeries("X", xVals),
		"Y": tabgo.NewSeries("Y", yVals),
	})

	iv := NewIVEstimator("X", "Y", []string{"Z"})

	ate, err := iv.Estimate(df)
	if err != nil {
		t.Fatalf("Estimate: %v", err)
	}
	if math.Abs(ate-3.0) > 0.1 {
		t.Errorf("expected ATE ~3.0, got %f", ate)
	}
	if !iv.Fitted() {
		t.Error("expected Fitted() to return true after Estimate")
	}
}

func TestIVEstimator_Estimate_NilData(t *testing.T) {
	iv := NewIVEstimator("X", "Y", []string{"Z"})
	_, err := iv.Estimate(nil)
	if err == nil {
		t.Error("expected error for nil data")
	}
}

func TestIVEstimator_Fitted_BeforeFit(t *testing.T) {
	iv := NewIVEstimator("X", "Y", []string{"Z"})
	if iv.Fitted() {
		t.Error("expected Fitted() to return false before fitting")
	}
}

// --- PC.BuildSkeleton tests ---

func TestPC_BuildSkeleton(t *testing.T) {
	// Three variables: A, B, C where A and B are independent given empty set,
	// and both connected to C.
	n := 200
	aVals := make([]any, n)
	bVals := make([]any, n)
	cVals := make([]any, n)

	for i := 0; i < n; i++ {
		a := i % 2
		b := (i / 2) % 2
		c := (a + b) % 2
		aVals[i] = a
		bVals[i] = b
		cVals[i] = c
	}

	df := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", aVals),
		"B": tabgo.NewSeries("B", bVals),
		"C": tabgo.NewSeries("C", cVals),
	})

	// Simple CI test: always says independent for (A,B) and dependent otherwise.
	ciTest := func(x, y string, z []string, data *tabgo.DataFrame, sig float64) (float64, float64, bool) {
		if (x == "A" && y == "B") || (x == "B" && y == "A") {
			if len(z) == 0 {
				return 0, 1.0, true // independent
			}
		}
		return 10.0, 0.001, false // dependent
	}

	pc := NewPC(df, ciTest, 0.05)

	skeleton, sepSets, err := pc.BuildSkeleton()
	if err != nil {
		t.Fatalf("BuildSkeleton: %v", err)
	}

	// A-B edge should be removed (they're independent).
	if skeleton.HasUndirectedEdge("A", "B") {
		t.Error("expected A-B edge to be removed in skeleton")
	}

	// A-C and B-C edges should remain.
	if !skeleton.HasUndirectedEdge("A", "C") {
		t.Error("expected A-C edge to remain in skeleton")
	}
	if !skeleton.HasUndirectedEdge("B", "C") {
		t.Error("expected B-C edge to remain in skeleton")
	}

	// Separating set for (A, B) should be empty.
	key := sepSetKey("A", "B")
	ss, ok := sepSets[key]
	if !ok {
		t.Error("expected separating set entry for (A, B)")
	}
	if len(ss) != 0 {
		t.Errorf("expected empty separating set, got %v", ss)
	}
}

func TestPC_BuildSkeleton_TooFewVariables(t *testing.T) {
	df := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1}),
	})

	ciTest := func(x, y string, z []string, data *tabgo.DataFrame, sig float64) (float64, float64, bool) {
		return 0, 1, true
	}

	pc := NewPC(df, ciTest, 0.05)
	_, _, err := pc.BuildSkeleton()
	if err == nil {
		t.Error("expected error for fewer than 2 variables")
	}
}

// --- HillClimbSearch.LegalOperations tests ---

func TestHillClimbSearch_LegalOperations(t *testing.T) {
	df := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 0, 1, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 1, 0, 1, 0, 1}),
	})

	// Score function that prefers A -> B.
	scoreFn := func(variable string, parents []string, data *tabgo.DataFrame) float64 {
		if variable == "B" && len(parents) == 1 && parents[0] == "A" {
			return 10.0
		}
		if variable == "A" && len(parents) == 1 && parents[0] == "B" {
			return 5.0
		}
		return 0.0
	}

	hc := NewHillClimbSearch(df, scoreFn)

	// Build an empty graph.
	g := newTestDiGraph([]string{"A", "B"})

	ops := hc.LegalOperations(g, []string{"A", "B"})
	if len(ops) == 0 {
		t.Fatal("expected at least one legal operation")
	}

	// Should find add operations for both directions.
	foundAddAB := false
	foundAddBA := false
	for _, op := range ops {
		if op.Type == "add" && op.From == "A" && op.To == "B" {
			foundAddAB = true
			if op.Delta <= 0 {
				t.Errorf("expected positive delta for add A->B, got %f", op.Delta)
			}
		}
		if op.Type == "add" && op.From == "B" && op.To == "A" {
			foundAddBA = true
		}
	}
	if !foundAddAB {
		t.Error("expected to find add A->B operation")
	}
	if !foundAddBA {
		t.Error("expected to find add B->A operation")
	}
}

func TestHillClimbSearch_LegalOperations_WithBlackList(t *testing.T) {
	df := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 0, 1, 1}),
		"B": tabgo.NewSeries("B", []any{0, 1, 0, 1}),
	})

	scoreFn := func(variable string, parents []string, data *tabgo.DataFrame) float64 {
		if len(parents) > 0 {
			return 5.0
		}
		return 0.0
	}

	hc := NewHillClimbSearch(df, scoreFn, WithBlackList([][2]string{{"A", "B"}}))
	g := newTestDiGraph([]string{"A", "B"})

	ops := hc.LegalOperations(g, []string{"A", "B"})
	for _, op := range ops {
		if op.Type == "add" && op.From == "A" && op.To == "B" {
			t.Error("blacklisted edge A->B should not appear in legal operations")
		}
	}
}

// --- ConstraintBasedEstimator tests ---

func TestConstraintBasedEstimator_Estimate(t *testing.T) {
	df := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 0, 1, 1}),
		"B": tabgo.NewSeries("B", []any{0, 1, 0, 1}),
	})

	ciTest := func(x, y string, z []string, data *tabgo.DataFrame, sig float64) (float64, float64, bool) {
		return 5.0, 0.01, false // always dependent
	}

	cbe := NewConstraintBasedEstimator(df, ciTest, 0.05)
	pdag, err := cbe.Estimate()
	if err != nil {
		t.Fatalf("Estimate: %v", err)
	}
	if pdag == nil {
		t.Fatal("Estimate returned nil")
	}
}

func TestConstraintBasedEstimator_EstimateBN(t *testing.T) {
	df := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 0, 1, 1}),
		"B": tabgo.NewSeries("B", []any{0, 1, 0, 1}),
	})

	ciTest := func(x, y string, z []string, data *tabgo.DataFrame, sig float64) (float64, float64, bool) {
		return 5.0, 0.01, false
	}

	cbe := NewConstraintBasedEstimator(df, ciTest, 0.05)
	bn, err := cbe.EstimateBN()
	if err != nil {
		t.Fatalf("EstimateBN: %v", err)
	}
	if bn == nil {
		t.Fatal("EstimateBN returned nil")
	}
	nodes := bn.Nodes()
	if len(nodes) != 2 {
		t.Errorf("expected 2 nodes, got %d", len(nodes))
	}
}

func TestConstraintBasedEstimator_BuildSkeleton(t *testing.T) {
	df := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 0, 1, 1}),
		"B": tabgo.NewSeries("B", []any{0, 1, 0, 1}),
	})

	ciTest := func(x, y string, z []string, data *tabgo.DataFrame, sig float64) (float64, float64, bool) {
		return 0, 1.0, true // always independent
	}

	cbe := NewConstraintBasedEstimator(df, ciTest, 0.05)
	skeleton, sepSets, err := cbe.BuildSkeleton()
	if err != nil {
		t.Fatalf("BuildSkeleton: %v", err)
	}
	if skeleton == nil {
		t.Fatal("BuildSkeleton returned nil skeleton")
	}
	if sepSets == nil {
		t.Fatal("BuildSkeleton returned nil sepSets")
	}
	// A and B should have been separated.
	if skeleton.HasUndirectedEdge("A", "B") {
		t.Error("expected A-B edge to be removed")
	}
}

// --- Scoring function tests ---

func TestBICScore_Basic(t *testing.T) {
	df := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 0, 0, 1, 1, 1, 0, 0, 1, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 0, 1, 1, 0, 1, 0, 1}),
	})

	bic := BICScore()

	// Score of B with no parents vs B with parent A.
	scoreNoParent := bic("B", nil, df)
	scoreWithParent := bic("B", []string{"A"}, df)

	// Both should be finite numbers.
	if math.IsNaN(scoreNoParent) || math.IsInf(scoreNoParent, 0) {
		t.Errorf("BIC score with no parents is not finite: %f", scoreNoParent)
	}
	if math.IsNaN(scoreWithParent) || math.IsInf(scoreWithParent, 0) {
		t.Errorf("BIC score with parent A is not finite: %f", scoreWithParent)
	}
}

func TestK2Score_Basic(t *testing.T) {
	df := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 0, 1, 1}),
		"B": tabgo.NewSeries("B", []any{0, 1, 0, 1}),
	})

	k2 := K2Score()
	score := k2("B", nil, df)
	if math.IsNaN(score) || math.IsInf(score, 0) {
		t.Errorf("K2 score is not finite: %f", score)
	}
}

func TestBDeuScore_Basic(t *testing.T) {
	df := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 0, 1, 1}),
		"B": tabgo.NewSeries("B", []any{0, 1, 0, 1}),
	})

	bdeu := BDeuScore(1.0)
	score := bdeu("B", []string{"A"}, df)
	if math.IsNaN(score) || math.IsInf(score, 0) {
		t.Errorf("BDeu score is not finite: %f", score)
	}
}

func TestAICScore_Basic(t *testing.T) {
	df := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 0, 0, 1, 1, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 0, 1, 1}),
	})

	aic := AICScore()
	score := aic("B", nil, df)
	if math.IsNaN(score) || math.IsInf(score, 0) {
		t.Errorf("AIC score is not finite: %f", score)
	}
}

func TestBICScore_PenalizesComplexity(t *testing.T) {
	// With truly independent data (B is constant pattern regardless of A),
	// BIC should prefer the simpler model (no parent).
	// We use a large sample so the penalty dominates any sampling noise.
	n := 1000
	aVals := make([]any, n)
	bVals := make([]any, n)
	for i := 0; i < n; i++ {
		aVals[i] = i % 2
		// B is a fixed repeating pattern independent of A.
		bVals[i] = i % 2 // same marginal as A but assigned independently
	}
	// Shuffle B values relative to A to break correlation.
	// Use a deterministic "shuffle" that keeps 50/50 but breaks A-B correlation.
	for i := 0; i < n; i++ {
		// B = 0 for first half, 1 for second half (independent of A's alternating).
		if i < n/2 {
			bVals[i] = 0
		} else {
			bVals[i] = 1
		}
	}

	df := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", aVals),
		"B": tabgo.NewSeries("B", bVals),
	})

	bic := BICScore()
	scoreNoParent := bic("B", nil, df)
	scoreWithParent := bic("B", []string{"A"}, df)

	// For independent data with large n, the simpler model should score higher.
	if scoreWithParent > scoreNoParent {
		t.Errorf("BIC should penalize unnecessary parent: no_parent=%f, with_parent=%f",
			scoreNoParent, scoreWithParent)
	}
}

// --- Helper for creating test DiGraph ---

func newTestDiGraph(nodes []string) *graphgo.DiGraph {
	g := graphgo.NewDiGraph()
	for _, n := range nodes {
		g.AddNode(n)
	}
	return g
}
