//go:build unit

package learning

import (
	"math"
	"math/rand"
	"testing"

	"github.com/asymmetric-effort/datascience/lib/pgm/models"
	"github.com/asymmetric-effort/datascience/lib/tabgo"
)

// buildLGBNChain creates a 3-variable chain: X -> Y -> Z.
func buildLGBNChain() *models.LinearGaussianBayesianNetwork {
	bn := models.NewLinearGaussianBayesianNetwork()
	_ = bn.AddNode("X")
	_ = bn.AddNode("Y")
	_ = bn.AddNode("Z")
	_ = bn.AddEdge("X", "Y")
	_ = bn.AddEdge("Y", "Z")
	return bn
}

// generateChainData generates synthetic data from the chain X -> Y -> Z
// with known parameters:
//
//	X ~ N(muX, varX)
//	Y | X ~ N(interceptY + betaYX * X, varY)
//	Z | Y ~ N(interceptZ + betaZY * Y, varZ)
func generateChainData(n int, seed int64) *tabgo.DataFrame {
	rng := rand.New(rand.NewSource(seed))

	// True parameters.
	muX := 5.0
	stdX := 1.0
	interceptY := 2.0
	betaYX := 0.5
	stdY := math.Sqrt(0.5)
	interceptZ := -1.0
	betaZY := 1.5
	stdZ := math.Sqrt(0.8)

	xVals := make([]any, n)
	yVals := make([]any, n)
	zVals := make([]any, n)

	for i := 0; i < n; i++ {
		x := muX + stdX*rng.NormFloat64()
		y := interceptY + betaYX*x + stdY*rng.NormFloat64()
		z := interceptZ + betaZY*y + stdZ*rng.NormFloat64()
		xVals[i] = x
		yVals[i] = y
		zVals[i] = z
	}

	return tabgo.NewDataFrame(map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", xVals),
		"Y": tabgo.NewSeries("Y", yVals),
		"Z": tabgo.NewSeries("Z", zVals),
	})
}

func TestNewLinearGaussianMLE(t *testing.T) {
	bn := buildLGBNChain()
	data := generateChainData(100, 42)
	mle := NewLinearGaussianMLE(bn, data)
	if mle == nil {
		t.Fatal("expected non-nil MLE")
	}
}

func TestLinearGaussianMLE_Estimate(t *testing.T) {
	bn := buildLGBNChain()
	data := generateChainData(10000, 42)
	mle := NewLinearGaussianMLE(bn, data)

	if err := mle.Estimate(); err != nil {
		t.Fatalf("Estimate: %v", err)
	}

	// Verify the model checks out.
	if err := bn.CheckModel(); err != nil {
		t.Fatalf("CheckModel after Estimate: %v", err)
	}

	// Check estimated parameters are close to true values.
	// True: X ~ N(5, 1)
	cpdX := bn.GetLinearGaussianCPD("X")
	if cpdX == nil {
		t.Fatal("expected CPD for X")
	}
	assertClose(t, "X mean", cpdX.Mean(), 5.0, 0.1)
	assertClose(t, "X variance", cpdX.Variance(), 1.0, 0.1)
	if len(cpdX.Betas()) != 0 {
		t.Errorf("X should have no betas, got %v", cpdX.Betas())
	}

	// True: Y | X ~ N(2 + 0.5*X, 0.5)
	cpdY := bn.GetLinearGaussianCPD("Y")
	if cpdY == nil {
		t.Fatal("expected CPD for Y")
	}
	assertClose(t, "Y intercept", cpdY.Mean(), 2.0, 0.2)
	yBetas := cpdY.Betas()
	if len(yBetas) != 1 {
		t.Fatalf("Y should have 1 beta, got %d", len(yBetas))
	}
	assertClose(t, "Y beta(X)", yBetas[0], 0.5, 0.05)
	assertClose(t, "Y variance", cpdY.Variance(), 0.5, 0.1)

	// True: Z | Y ~ N(-1 + 1.5*Y, 0.8)
	cpdZ := bn.GetLinearGaussianCPD("Z")
	if cpdZ == nil {
		t.Fatal("expected CPD for Z")
	}
	assertClose(t, "Z intercept", cpdZ.Mean(), -1.0, 0.3)
	zBetas := cpdZ.Betas()
	if len(zBetas) != 1 {
		t.Fatalf("Z should have 1 beta, got %d", len(zBetas))
	}
	assertClose(t, "Z beta(Y)", zBetas[0], 1.5, 0.05)
	assertClose(t, "Z variance", cpdZ.Variance(), 0.8, 0.15)
}

func TestLinearGaussianMLE_GetParameters(t *testing.T) {
	bn := buildLGBNChain()
	data := generateChainData(5000, 99)
	mle := NewLinearGaussianMLE(bn, data)

	cpd, err := mle.GetParameters("Y")
	if err != nil {
		t.Fatalf("GetParameters(Y): %v", err)
	}
	if cpd.Variable() != "Y" {
		t.Errorf("expected variable Y, got %q", cpd.Variable())
	}
	assertClose(t, "Y intercept", cpd.Mean(), 2.0, 0.3)
	betas := cpd.Betas()
	if len(betas) != 1 {
		t.Fatalf("expected 1 beta, got %d", len(betas))
	}
	assertClose(t, "Y beta(X)", betas[0], 0.5, 0.1)
}

func TestLinearGaussianMLE_GetParameters_NotFound(t *testing.T) {
	bn := buildLGBNChain()
	data := generateChainData(100, 1)
	mle := NewLinearGaussianMLE(bn, data)

	_, err := mle.GetParameters("W")
	if err == nil {
		t.Error("expected error for nonexistent node")
	}
}

func TestLinearGaussianMLE_Estimate_NilBN(t *testing.T) {
	mle := NewLinearGaussianMLE(nil, nil)
	if err := mle.Estimate(); err == nil {
		t.Error("expected error for nil BN")
	}
}

func TestLinearGaussianMLE_Estimate_NilData(t *testing.T) {
	bn := buildLGBNChain()
	mle := NewLinearGaussianMLE(bn, nil)
	if err := mle.Estimate(); err == nil {
		t.Error("expected error for nil data")
	}
}

func TestLinearGaussianMLE_Estimate_MissingColumn(t *testing.T) {
	bn := buildLGBNChain()
	// DataFrame missing column Z.
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", []any{1.0, 2.0}),
		"Y": tabgo.NewSeries("Y", []any{1.0, 2.0}),
	})
	mle := NewLinearGaussianMLE(bn, data)
	if err := mle.Estimate(); err == nil {
		t.Error("expected error for missing column")
	}
}

func TestLinearGaussianMLE_RootNode(t *testing.T) {
	// A single root node should estimate mean and variance from data.
	bn := models.NewLinearGaussianBayesianNetwork()
	_ = bn.AddNode("A")

	rng := rand.New(rand.NewSource(42))
	n := 5000
	vals := make([]any, n)
	for i := range vals {
		vals[i] = 10.0 + 3.0*rng.NormFloat64()
	}
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", vals),
	})

	mle := NewLinearGaussianMLE(bn, data)
	if err := mle.Estimate(); err != nil {
		t.Fatalf("Estimate: %v", err)
	}

	cpd := bn.GetLinearGaussianCPD("A")
	assertClose(t, "A mean", cpd.Mean(), 10.0, 0.2)
	assertClose(t, "A variance", cpd.Variance(), 9.0, 0.5)
}

func TestLinearGaussianMLE_MultipleParents(t *testing.T) {
	// Test with a node that has 2 parents: A, B -> C
	// C | A, B ~ N(1.0 + 2.0*A + 3.0*B, 0.5)
	bn := models.NewLinearGaussianBayesianNetwork()
	_ = bn.AddNode("A")
	_ = bn.AddNode("B")
	_ = bn.AddNode("C")
	_ = bn.AddEdge("A", "C")
	_ = bn.AddEdge("B", "C")

	rng := rand.New(rand.NewSource(123))
	n := 10000
	aVals := make([]any, n)
	bVals := make([]any, n)
	cVals := make([]any, n)
	for i := 0; i < n; i++ {
		a := rng.NormFloat64()
		b := 2.0 + rng.NormFloat64()
		c := 1.0 + 2.0*a + 3.0*b + math.Sqrt(0.5)*rng.NormFloat64()
		aVals[i] = a
		bVals[i] = b
		cVals[i] = c
	}
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", aVals),
		"B": tabgo.NewSeries("B", bVals),
		"C": tabgo.NewSeries("C", cVals),
	})

	mle := NewLinearGaussianMLE(bn, data)
	if err := mle.Estimate(); err != nil {
		t.Fatalf("Estimate: %v", err)
	}

	cpdC := bn.GetLinearGaussianCPD("C")
	if cpdC == nil {
		t.Fatal("expected CPD for C")
	}
	assertClose(t, "C intercept", cpdC.Mean(), 1.0, 0.2)
	betas := cpdC.Betas()
	// Parents are sorted: A, B
	if len(betas) != 2 {
		t.Fatalf("expected 2 betas, got %d", len(betas))
	}
	assertClose(t, "C beta(A)", betas[0], 2.0, 0.1)
	assertClose(t, "C beta(B)", betas[1], 3.0, 0.1)
	assertClose(t, "C variance", cpdC.Variance(), 0.5, 0.1)
}

func TestSolveLinearSystem(t *testing.T) {
	// Solve: 2x + y = 5, x + 3y = 7 => x=1.6, y=1.8
	A := []float64{2, 1, 1, 3}
	b := []float64{5, 7}
	x, err := solveLinearSystem(A, b, 2)
	if err != nil {
		t.Fatalf("solveLinearSystem: %v", err)
	}
	assertClose(t, "x[0]", x[0], 1.6, 1e-10)
	assertClose(t, "x[1]", x[1], 1.8, 1e-10)
}

func TestSolveLinearSystem_Singular(t *testing.T) {
	// Singular matrix.
	A := []float64{1, 2, 2, 4}
	b := []float64{3, 6}
	_, err := solveLinearSystem(A, b, 2)
	if err == nil {
		t.Error("expected error for singular matrix")
	}
}

func TestSolveLinearSystem_1x1(t *testing.T) {
	A := []float64{4}
	b := []float64{8}
	x, err := solveLinearSystem(A, b, 1)
	if err != nil {
		t.Fatalf("solveLinearSystem: %v", err)
	}
	assertClose(t, "x[0]", x[0], 2.0, 1e-10)
}

// assertClose checks that got is within tol of want.
func assertClose(t *testing.T, name string, got, want, tol float64) {
	t.Helper()
	if math.Abs(got-want) > tol {
		t.Errorf("%s: got %f, want %f (tol %f)", name, got, want, tol)
	}
}

func TestLinearGaussianMLE_PerfectFit(t *testing.T) {
	// Perfect linear relationship: Y = 3 + 2*X (no noise).
	bn := models.NewLinearGaussianBayesianNetwork()
	_ = bn.AddNode("X")
	_ = bn.AddNode("Y")
	_ = bn.AddEdge("X", "Y")

	n := 100
	xVals := make([]any, n)
	yVals := make([]any, n)
	for i := 0; i < n; i++ {
		x := float64(i)
		xVals[i] = x
		yVals[i] = 3.0 + 2.0*x
	}
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", xVals),
		"Y": tabgo.NewSeries("Y", yVals),
	})

	mle := NewLinearGaussianMLE(bn, data)
	if err := mle.Estimate(); err != nil {
		t.Fatalf("Estimate: %v", err)
	}

	cpdY := bn.GetLinearGaussianCPD("Y")
	assertClose(t, "Y intercept", cpdY.Mean(), 3.0, 1e-8)
	betas := cpdY.Betas()
	if len(betas) != 1 {
		t.Fatalf("expected 1 beta, got %d", len(betas))
	}
	assertClose(t, "Y beta(X)", betas[0], 2.0, 1e-8)
	// Variance should be very small (perfect fit, guarded to 1e-10).
	if cpdY.Variance() > 1e-6 {
		t.Errorf("expected near-zero variance for perfect fit, got %f", cpdY.Variance())
	}
}

func TestLinearGaussianMLE_Estimate_EmptyNetwork(t *testing.T) {
	bn := models.NewLinearGaussianBayesianNetwork()
	data := generateChainData(100, 1)
	mle := NewLinearGaussianMLE(bn, data)
	if err := mle.Estimate(); err == nil {
		t.Error("expected error for empty network")
	}
}

func TestLinearGaussianMLE_GetParameters_NilBN(t *testing.T) {
	mle := NewLinearGaussianMLE(nil, nil)
	_, err := mle.GetParameters("X")
	if err == nil {
		t.Error("expected error for nil BN")
	}
}

func TestLinearGaussianMLE_GetParameters_MissingParentColumn(t *testing.T) {
	bn := models.NewLinearGaussianBayesianNetwork()
	_ = bn.AddNode("X")
	_ = bn.AddNode("Y")
	_ = bn.AddEdge("X", "Y")

	// Data only has Y, missing parent X.
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"Y": tabgo.NewSeries("Y", []any{1.0, 2.0}),
	})
	mle := NewLinearGaussianMLE(bn, data)
	_, err := mle.GetParameters("Y")
	if err == nil {
		t.Error("expected error for missing parent column")
	}
}

func TestLinearGaussianMLE_InsufficientData(t *testing.T) {
	// Node with 1 parent needs at least 2 data points (intercept + 1 beta).
	// Provide only 1 row.
	bn := models.NewLinearGaussianBayesianNetwork()
	_ = bn.AddNode("X")
	_ = bn.AddNode("Y")
	_ = bn.AddEdge("X", "Y")

	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", []any{1.0}),
		"Y": tabgo.NewSeries("Y", []any{2.0}),
	})
	mle := NewLinearGaussianMLE(bn, data)
	err := mle.Estimate()
	if err == nil {
		t.Error("expected error for insufficient data")
	}
}
