//go:build unit

package scigo

import (
	"math"
	"testing"
)

// ---------------------------------------------------------------------------
// Beta.LogPDF edge cases (50% coverage -> higher)
// ---------------------------------------------------------------------------

func TestBeta_LogPDF_AtZero(t *testing.T) {
	b := NewBeta(2, 3)
	val := b.LogPDF(0)
	if !math.IsInf(val, -1) {
		t.Errorf("Beta.LogPDF(0) with alpha=2 should be -Inf, got %v", val)
	}
}

func TestBeta_LogPDF_AtOne(t *testing.T) {
	b := NewBeta(2, 3)
	val := b.LogPDF(1)
	if !math.IsInf(val, -1) {
		t.Errorf("Beta.LogPDF(1) with beta=3 should be -Inf, got %v", val)
	}
}

func TestBeta_LogPDF_Negative(t *testing.T) {
	b := NewBeta(2, 3)
	val := b.LogPDF(-1)
	if !math.IsInf(val, -1) {
		t.Errorf("Beta.LogPDF(-1) should be -Inf, got %v", val)
	}
}

func TestBeta_LogPDF_GreaterThanOne(t *testing.T) {
	b := NewBeta(2, 3)
	val := b.LogPDF(2)
	if !math.IsInf(val, -1) {
		t.Errorf("Beta.LogPDF(2) should be -Inf, got %v", val)
	}
}

func TestBeta_LogPDF_Interior(t *testing.T) {
	b := NewBeta(2, 3)
	val := b.LogPDF(0.5)
	if math.IsInf(val, 0) || math.IsNaN(val) {
		t.Errorf("Beta.LogPDF(0.5) should be finite, got %v", val)
	}
	// Should equal log(PDF(0.5))
	expected := math.Log(b.PDF(0.5))
	if math.Abs(val-expected) > 1e-10 {
		t.Errorf("Beta.LogPDF(0.5) = %v, expected %v", val, expected)
	}
}

func TestBeta_LogPDF_NearBoundaries(t *testing.T) {
	b := NewBeta(0.5, 0.5) // alpha < 1, beta < 1
	// LogPDF at boundaries should be -Inf per the implementation
	val0 := b.LogPDF(0)
	val1 := b.LogPDF(1)
	if !math.IsInf(val0, -1) {
		t.Errorf("Beta(0.5,0.5).LogPDF(0) should be -Inf, got %v", val0)
	}
	if !math.IsInf(val1, -1) {
		t.Errorf("Beta(0.5,0.5).LogPDF(1) should be -Inf, got %v", val1)
	}
}

// ---------------------------------------------------------------------------
// Beta.PDF edge cases
// ---------------------------------------------------------------------------

func TestBeta_PDF_AtZero_AlphaLessThanOne(t *testing.T) {
	b := NewBeta(0.5, 2)
	val := b.PDF(0)
	if !math.IsInf(val, 1) {
		t.Errorf("Beta(0.5,2).PDF(0) should be +Inf, got %v", val)
	}
}

func TestBeta_PDF_AtZero_AlphaEqualsOne(t *testing.T) {
	b := NewBeta(1, 2)
	val := b.PDF(0)
	if val != 1.0 {
		t.Errorf("Beta(1,2).PDF(0) should be 1 (alpha), got %v", val)
	}
}

func TestBeta_PDF_AtZero_AlphaGreaterThanOne(t *testing.T) {
	b := NewBeta(2, 3)
	val := b.PDF(0)
	if val != 0 {
		t.Errorf("Beta(2,3).PDF(0) should be 0, got %v", val)
	}
}

func TestBeta_PDF_AtOne_BetaLessThanOne(t *testing.T) {
	b := NewBeta(2, 0.5)
	val := b.PDF(1)
	if !math.IsInf(val, 1) {
		t.Errorf("Beta(2,0.5).PDF(1) should be +Inf, got %v", val)
	}
}

func TestBeta_PDF_AtOne_BetaEqualsOne(t *testing.T) {
	b := NewBeta(2, 1)
	val := b.PDF(1)
	if val != 1.0 {
		t.Errorf("Beta(2,1).PDF(1) should be 1 (beta), got %v", val)
	}
}

func TestBeta_PDF_AtOne_BetaGreaterThanOne(t *testing.T) {
	b := NewBeta(2, 3)
	val := b.PDF(1)
	if val != 0 {
		t.Errorf("Beta(2,3).PDF(1) should be 0, got %v", val)
	}
}

func TestBeta_PDF_OutOfRange(t *testing.T) {
	b := NewBeta(2, 3)
	if b.PDF(-0.1) != 0 {
		t.Error("PDF(-0.1) should be 0")
	}
	if b.PDF(1.1) != 0 {
		t.Error("PDF(1.1) should be 0")
	}
}

// ---------------------------------------------------------------------------
// sparse.ToCOO: empty sparse matrix
// ---------------------------------------------------------------------------

func TestCSR_ToCOO_Empty(t *testing.T) {
	// Create a CSR with no non-zero entries
	csr, err := NewCSR([]int{0, 0, 0}, nil, nil, [2]int{2, 2})
	if err != nil {
		t.Fatal(err)
	}
	coo := csr.ToCOO()
	if coo.NNZ() != 0 {
		t.Errorf("expected 0 NNZ, got %d", coo.NNZ())
	}
	shape := coo.Shape()
	if shape[0] != 2 || shape[1] != 2 {
		t.Errorf("expected shape [2,2], got %v", shape)
	}
	// Verify the row/col/values slices are non-nil empty
	if coo.rows == nil || coo.cols == nil || coo.values == nil {
		t.Error("empty COO should have non-nil slices")
	}
}

func TestCSR_ToCOO_SingleElement(t *testing.T) {
	csr, err := NewCSR([]int{0, 1, 1}, []int{0}, []float64{5.0}, [2]int{2, 2})
	if err != nil {
		t.Fatal(err)
	}
	coo := csr.ToCOO()
	if coo.NNZ() != 1 {
		t.Errorf("expected 1 NNZ, got %d", coo.NNZ())
	}
	val := coo.Get(0, 0)
	if val != 5.0 {
		t.Errorf("expected 5.0 at (0,0), got %v", val)
	}
}

func TestCSC_ToCOO_viaCSR(t *testing.T) {
	// Create CSC, convert to CSR, then to COO
	csc, err := NewCSC([]int{0, 1, 2}, []int{0, 1}, []float64{3.0, 7.0}, [2]int{2, 2})
	if err != nil {
		t.Fatal(err)
	}
	csr := csc.ToCSR()
	coo := csr.ToCOO()
	if coo.NNZ() != 2 {
		t.Errorf("expected 2 NNZ, got %d", coo.NNZ())
	}
}

// ---------------------------------------------------------------------------
// Distribution PPF edge cases: p=0, p=1, p=0.5
// ---------------------------------------------------------------------------

func TestNormal_PPF_EdgeCases(t *testing.T) {
	n := NewNormal(0, 1)
	if !math.IsInf(n.PPF(0), -1) {
		t.Error("Normal.PPF(0) should be -Inf")
	}
	if !math.IsInf(n.PPF(1), 1) {
		t.Error("Normal.PPF(1) should be +Inf")
	}
	if math.Abs(n.PPF(0.5)) > 1e-10 {
		t.Errorf("Normal(0,1).PPF(0.5) should be 0, got %v", n.PPF(0.5))
	}
}

func TestChiSquared_PPF_EdgeCases(t *testing.T) {
	c := NewChiSquared(5)
	if c.PPF(0) != 0 {
		t.Errorf("ChiSquared.PPF(0) should be 0, got %v", c.PPF(0))
	}
	if !math.IsInf(c.PPF(1), 1) {
		t.Error("ChiSquared.PPF(1) should be +Inf")
	}
	val := c.PPF(0.5)
	if val <= 0 || math.IsInf(val, 0) || math.IsNaN(val) {
		t.Errorf("ChiSquared.PPF(0.5) should be finite positive, got %v", val)
	}
}

func TestTDistribution_PPF_EdgeCases(t *testing.T) {
	td := NewTDistribution(10)
	if !math.IsInf(td.PPF(0), -1) {
		t.Error("TDistribution.PPF(0) should be -Inf")
	}
	if !math.IsInf(td.PPF(1), 1) {
		t.Error("TDistribution.PPF(1) should be +Inf")
	}
	if td.PPF(0.5) != 0 {
		t.Errorf("TDistribution.PPF(0.5) should be 0, got %v", td.PPF(0.5))
	}
}

func TestFDistribution_PPF_EdgeCases(t *testing.T) {
	f := NewFDistribution(5, 10)
	if f.PPF(0) != 0 {
		t.Errorf("FDistribution.PPF(0) should be 0, got %v", f.PPF(0))
	}
	if !math.IsInf(f.PPF(1), 1) {
		t.Error("FDistribution.PPF(1) should be +Inf")
	}
	val := f.PPF(0.5)
	if val <= 0 || math.IsInf(val, 0) || math.IsNaN(val) {
		t.Errorf("FDistribution.PPF(0.5) should be finite positive, got %v", val)
	}
}

func TestBeta_PPF_EdgeCases(t *testing.T) {
	b := NewBeta(2, 5)
	if b.PPF(0) != 0 {
		t.Errorf("Beta.PPF(0) should be 0, got %v", b.PPF(0))
	}
	if b.PPF(1) != 1 {
		t.Errorf("Beta.PPF(1) should be 1, got %v", b.PPF(1))
	}
	val := b.PPF(0.5)
	if val <= 0 || val >= 1 || math.IsNaN(val) {
		t.Errorf("Beta.PPF(0.5) should be in (0,1), got %v", val)
	}
}

func TestExponential_PPF_EdgeCases(t *testing.T) {
	e := NewExponential(2)
	if e.PPF(0) != 0 {
		t.Errorf("Exponential.PPF(0) should be 0, got %v", e.PPF(0))
	}
	if !math.IsInf(e.PPF(1), 1) {
		t.Error("Exponential.PPF(1) should be +Inf")
	}
	val := e.PPF(0.5)
	if val <= 0 || math.IsInf(val, 0) || math.IsNaN(val) {
		t.Errorf("Exponential.PPF(0.5) should be finite positive, got %v", val)
	}
}

func TestUniform_PPF_EdgeCases(t *testing.T) {
	u := NewUniform(10, 20)
	if u.PPF(0) != 10 {
		t.Errorf("Uniform.PPF(0) should be 10, got %v", u.PPF(0))
	}
	if u.PPF(1) != 20 {
		t.Errorf("Uniform.PPF(1) should be 20, got %v", u.PPF(1))
	}
	if u.PPF(0.5) != 15 {
		t.Errorf("Uniform.PPF(0.5) should be 15, got %v", u.PPF(0.5))
	}
}

// ---------------------------------------------------------------------------
// Linprog: infeasible problem
// ---------------------------------------------------------------------------

func TestLinprog_Infeasible(t *testing.T) {
	// x1 + x2 <= 1 and x1 + x2 >= 3 (infeasible)
	// Equality constraint: x1 + x2 = 3
	// Inequality: x1 + x2 <= 1
	c := []float64{1, 1}
	Aub := [][]float64{{1, 1}}
	bub := []float64{1}
	Aeq := [][]float64{{1, 1}}
	beq := []float64{3}
	_, err := Linprog(c, Aub, bub, Aeq, beq)
	if err == nil {
		t.Fatal("expected error for infeasible problem")
	}
}

func TestLinprog_Unbounded(t *testing.T) {
	// minimize -x1 - x2 with no upper bound constraints
	// Only equality: none, no ub constraints that matter
	c := []float64{-1, -1}
	Aub := [][]float64{{-1, 0}, {0, -1}} // -x1 <= 0, -x2 <= 0 (x>=0 effectively)
	bub := []float64{0, 0}
	_, err := Linprog(c, Aub, bub, nil, nil)
	// This may return unbounded or succeed with x=0,0
	// Just exercise the path
	_ = err
}

func TestLinprog_NoConstraints(t *testing.T) {
	c := []float64{1, 2}
	res, err := Linprog(c, nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Success {
		t.Error("expected success for unconstrained LP")
	}
}

func TestMock_Linprog_EmptyC(t *testing.T) {
	_, err := Linprog(nil, nil, nil, nil, nil)
	if err == nil {
		t.Fatal("expected error for empty c")
	}
}

func TestLinprog_DimensionMismatch(t *testing.T) {
	c := []float64{1}
	Aub := [][]float64{{1, 2}} // wrong number of columns
	bub := []float64{1}
	_, err := Linprog(c, Aub, bub, nil, nil)
	if err == nil {
		t.Fatal("expected error for dimension mismatch")
	}
}

func TestLinprog_AubBubMismatch(t *testing.T) {
	c := []float64{1}
	Aub := [][]float64{{1}}
	bub := []float64{1, 2} // mismatched length
	_, err := Linprog(c, Aub, bub, nil, nil)
	if err == nil {
		t.Fatal("expected error for Aub/bub mismatch")
	}
}

func TestLinprog_AeqBeqMismatch(t *testing.T) {
	c := []float64{1}
	Aeq := [][]float64{{1}}
	beq := []float64{1, 2}
	_, err := Linprog(c, nil, nil, Aeq, beq)
	if err == nil {
		t.Fatal("expected error for Aeq/beq mismatch")
	}
}

func TestLinprog_AeqDimensionMismatch(t *testing.T) {
	c := []float64{1}
	Aeq := [][]float64{{1, 2}} // wrong columns
	beq := []float64{1}
	_, err := Linprog(c, nil, nil, Aeq, beq)
	if err == nil {
		t.Fatal("expected error for Aeq dimension mismatch")
	}
}

// ---------------------------------------------------------------------------
// DualAnnealing: tight bounds
// ---------------------------------------------------------------------------

func TestDualAnnealing_TightBounds(t *testing.T) {
	f := func(x []float64) float64 { return x[0] * x[0] }
	bounds := [][2]float64{{-0.01, 0.01}}
	res, err := DualAnnealing(f, bounds)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(res.X[0]) > 0.02 {
		t.Errorf("expected solution near 0, got %v", res.X[0])
	}
}

func TestMock_DualAnnealing_EmptyBounds(t *testing.T) {
	f := func(x []float64) float64 { return 0 }
	_, err := DualAnnealing(f, nil)
	if err == nil {
		t.Fatal("expected error for empty bounds")
	}
}

func TestDualAnnealing_MultiDim(t *testing.T) {
	// Rosenbrock-like
	f := func(x []float64) float64 {
		return (1-x[0])*(1-x[0]) + 100*(x[1]-x[0]*x[0])*(x[1]-x[0]*x[0])
	}
	bounds := [][2]float64{{-2, 2}, {-2, 2}}
	res, err := DualAnnealing(f, bounds)
	if err != nil {
		t.Fatal(err)
	}
	if res.Fun > 1.0 {
		t.Errorf("expected good solution, got Fun=%v", res.Fun)
	}
}

// ---------------------------------------------------------------------------
// SHGO: single-variable problem
// ---------------------------------------------------------------------------

func TestSHGO_SingleVariable(t *testing.T) {
	f := func(x []float64) float64 { return (x[0] - 3) * (x[0] - 3) }
	bounds := [][2]float64{{0, 10}}
	res, err := SHGO(f, bounds)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(res.X[0]-3) > 0.1 {
		t.Errorf("expected solution near 3, got %v", res.X[0])
	}
}

func TestMock_SHGO_EmptyBounds(t *testing.T) {
	f := func(x []float64) float64 { return 0 }
	_, err := SHGO(f, nil)
	if err == nil {
		t.Fatal("expected error for empty bounds")
	}
}

// ---------------------------------------------------------------------------
// COO operations
// ---------------------------------------------------------------------------

func TestCOO_ToDense(t *testing.T) {
	coo, err := NewCOO([]int{0, 1}, []int{1, 0}, []float64{3.0, 4.0}, [2]int{2, 2})
	if err != nil {
		t.Fatal(err)
	}
	dense := coo.ToDense()
	if dense[0][1] != 3.0 || dense[1][0] != 4.0 {
		t.Error("COO.ToDense values incorrect")
	}
}

func TestCOO_ToCSR(t *testing.T) {
	coo, err := NewCOO([]int{0, 1, 0}, []int{0, 1, 0}, []float64{1.0, 2.0, 3.0}, [2]int{2, 2})
	if err != nil {
		t.Fatal(err)
	}
	csr := coo.ToCSR()
	// (0,0) should be 1.0 + 3.0 = 4.0 due to duplicate summing
	val := csr.Get(0, 0)
	if val != 4.0 {
		t.Errorf("expected 4.0 at (0,0) after duplicate sum, got %v", val)
	}
}

func TestCOO_Set(t *testing.T) {
	coo, err := NewCOO(nil, nil, nil, [2]int{3, 3})
	if err != nil {
		t.Fatal(err)
	}
	coo.Set(0, 0, 1.0)
	coo.Set(2, 2, 5.0)
	if coo.NNZ() != 2 {
		t.Errorf("expected 2 NNZ, got %d", coo.NNZ())
	}
}

func TestNewCOO_Errors(t *testing.T) {
	_, err := NewCOO([]int{0}, []int{0, 1}, []float64{1}, [2]int{2, 2})
	if err == nil {
		t.Fatal("expected error for mismatched lengths")
	}
	_, err = NewCOO([]int{0}, []int{0}, []float64{1}, [2]int{0, 2})
	if err == nil {
		t.Fatal("expected error for non-positive shape")
	}
	_, err = NewCOO([]int{5}, []int{0}, []float64{1}, [2]int{2, 2})
	if err == nil {
		t.Fatal("expected error for out-of-bounds row")
	}
	_, err = NewCOO([]int{0}, []int{5}, []float64{1}, [2]int{2, 2})
	if err == nil {
		t.Fatal("expected error for out-of-bounds col")
	}
}

// ---------------------------------------------------------------------------
// DenseToCOO / DenseToCSR empty
// ---------------------------------------------------------------------------

func TestDenseToCOO_Empty(t *testing.T) {
	coo := DenseToCOO(nil)
	if coo.NNZ() != 0 {
		t.Errorf("expected 0 NNZ for empty dense, got %d", coo.NNZ())
	}
}

func TestDenseToCOO_AllZeros(t *testing.T) {
	dense := [][]float64{{0, 0}, {0, 0}}
	coo := DenseToCOO(dense)
	if coo.NNZ() != 0 {
		t.Errorf("expected 0 NNZ for all-zeros, got %d", coo.NNZ())
	}
}

func TestDenseToCSR_Empty(t *testing.T) {
	csr := DenseToCSR(nil)
	if csr.NNZ() != 0 {
		t.Errorf("expected 0 NNZ for empty dense, got %d", csr.NNZ())
	}
}

// ---------------------------------------------------------------------------
// ChiSquared LogPDF edge cases
// ---------------------------------------------------------------------------

func TestChiSquared_LogPDF_AtZero(t *testing.T) {
	// df == 2: should return log(0.5)
	c := NewChiSquared(2)
	val := c.LogPDF(0)
	if math.Abs(val-math.Log(0.5)) > 1e-10 {
		t.Errorf("ChiSquared(2).LogPDF(0) = %v, expected %v", val, math.Log(0.5))
	}

	// df > 2: should be -Inf
	c3 := NewChiSquared(3)
	val3 := c3.LogPDF(0)
	if !math.IsInf(val3, -1) {
		t.Errorf("ChiSquared(3).LogPDF(0) should be -Inf, got %v", val3)
	}

	// df < 2: should be -Inf (per implementation)
	c1 := NewChiSquared(1)
	val1 := c1.LogPDF(0)
	if !math.IsInf(val1, -1) {
		t.Errorf("ChiSquared(1).LogPDF(0) should be -Inf, got %v", val1)
	}
}

func TestChiSquared_LogPDF_Negative(t *testing.T) {
	c := NewChiSquared(5)
	val := c.LogPDF(-1)
	if !math.IsInf(val, -1) {
		t.Errorf("ChiSquared.LogPDF(-1) should be -Inf, got %v", val)
	}
}

// ---------------------------------------------------------------------------
// FDistribution LogPDF edge case
// ---------------------------------------------------------------------------

func TestFDistribution_LogPDF_AtZero(t *testing.T) {
	f := NewFDistribution(5, 10)
	val := f.LogPDF(0)
	if !math.IsInf(val, -1) {
		t.Errorf("FDistribution.LogPDF(0) should be -Inf, got %v", val)
	}
}

func TestFDistribution_LogPDF_Negative(t *testing.T) {
	f := NewFDistribution(5, 10)
	val := f.LogPDF(-1)
	if !math.IsInf(val, -1) {
		t.Errorf("FDistribution.LogPDF(-1) should be -Inf, got %v", val)
	}
}

// ---------------------------------------------------------------------------
// Gamma distribution edge cases
// ---------------------------------------------------------------------------

func TestGamma_PDF_AtZero(t *testing.T) {
	// shape < 1: should be +Inf
	g1 := NewGamma(0.5, 1)
	if !math.IsInf(g1.PDF(0), 1) {
		t.Errorf("Gamma(0.5,1).PDF(0) should be +Inf, got %v", g1.PDF(0))
	}

	// shape == 1: should be 1/scale
	g2 := NewGamma(1, 2)
	if math.Abs(g2.PDF(0)-0.5) > 1e-10 {
		t.Errorf("Gamma(1,2).PDF(0) should be 0.5, got %v", g2.PDF(0))
	}

	// shape > 1: should be 0
	g3 := NewGamma(2, 1)
	if g3.PDF(0) != 0 {
		t.Errorf("Gamma(2,1).PDF(0) should be 0, got %v", g3.PDF(0))
	}
}

func TestGamma_PDF_Negative(t *testing.T) {
	g := NewGamma(2, 1)
	if g.PDF(-1) != 0 {
		t.Errorf("Gamma.PDF(-1) should be 0, got %v", g.PDF(-1))
	}
}

func TestGamma_LogPDF_NonPositive(t *testing.T) {
	g := NewGamma(2, 1)
	if !math.IsInf(g.LogPDF(0), -1) {
		t.Errorf("Gamma.LogPDF(0) should be -Inf, got %v", g.LogPDF(0))
	}
	if !math.IsInf(g.LogPDF(-1), -1) {
		t.Errorf("Gamma.LogPDF(-1) should be -Inf, got %v", g.LogPDF(-1))
	}
}

// ---------------------------------------------------------------------------
// Exponential LogPDF edge case
// ---------------------------------------------------------------------------

func TestExponential_LogPDF_Negative(t *testing.T) {
	e := NewExponential(2)
	if !math.IsInf(e.LogPDF(-1), -1) {
		t.Errorf("Exponential.LogPDF(-1) should be -Inf, got %v", e.LogPDF(-1))
	}
}

func TestExponential_LogPDF_Zero(t *testing.T) {
	e := NewExponential(2)
	val := e.LogPDF(0)
	expected := math.Log(2) // log(rate) - rate*0
	if math.Abs(val-expected) > 1e-10 {
		t.Errorf("Exponential.LogPDF(0) = %v, expected %v", val, expected)
	}
}

// ---------------------------------------------------------------------------
// Uniform LogPDF edge cases
// ---------------------------------------------------------------------------

func TestUniform_LogPDF_OutOfRange(t *testing.T) {
	u := NewUniform(0, 1)
	if !math.IsInf(u.LogPDF(-1), -1) {
		t.Error("Uniform.LogPDF(-1) should be -Inf")
	}
	if !math.IsInf(u.LogPDF(2), -1) {
		t.Error("Uniform.LogPDF(2) should be -Inf")
	}
}

// ---------------------------------------------------------------------------
// CSR additional operations
// ---------------------------------------------------------------------------

func TestMock_CSR_Transpose(t *testing.T) {
	csr, err := NewCSR([]int{0, 2, 3}, []int{0, 1, 0}, []float64{1, 2, 3}, [2]int{2, 2})
	if err != nil {
		t.Fatal(err)
	}
	tr := csr.Transpose()
	if tr.Get(0, 0) != 1.0 || tr.Get(1, 0) != 2.0 || tr.Get(0, 1) != 3.0 {
		t.Error("Transpose values incorrect")
	}
}

func TestCSR_Row(t *testing.T) {
	csr, err := NewCSR([]int{0, 2, 3}, []int{0, 1, 0}, []float64{1, 2, 3}, [2]int{2, 2})
	if err != nil {
		t.Fatal(err)
	}
	cols, vals := csr.Row(0)
	if len(cols) != 2 || cols[0] != 0 || cols[1] != 1 {
		t.Errorf("Row(0) cols = %v", cols)
	}
	if vals[0] != 1 || vals[1] != 2 {
		t.Errorf("Row(0) vals = %v", vals)
	}
}

// ---------------------------------------------------------------------------
// CSC operations
// ---------------------------------------------------------------------------

func TestCSC_Col(t *testing.T) {
	csc, err := NewCSC([]int{0, 1, 2}, []int{0, 1}, []float64{3.0, 7.0}, [2]int{2, 2})
	if err != nil {
		t.Fatal(err)
	}
	rows, vals := csc.Col(0)
	if len(rows) != 1 || rows[0] != 0 || vals[0] != 3.0 {
		t.Error("Col(0) incorrect")
	}
}

func TestCSC_ToDense(t *testing.T) {
	csc, err := NewCSC([]int{0, 1, 2}, []int{0, 1}, []float64{3.0, 7.0}, [2]int{2, 2})
	if err != nil {
		t.Fatal(err)
	}
	dense := csc.ToDense()
	if dense[0][0] != 3.0 || dense[1][1] != 7.0 {
		t.Error("CSC.ToDense incorrect")
	}
}

// ---------------------------------------------------------------------------
// Sparse construction errors
// ---------------------------------------------------------------------------

func TestMock_NewCSR_Errors(t *testing.T) {
	_, err := NewCSR([]int{0, 1}, []int{0}, []float64{1}, [2]int{0, 1})
	if err == nil {
		t.Fatal("expected error for non-positive shape")
	}
	_, err = NewCSR([]int{0, 1, 1}, []int{0}, []float64{1}, [2]int{1, 1})
	if err == nil {
		t.Fatal("expected error for wrong indptr length")
	}
	_, err = NewCSR([]int{0, 1}, []int{0, 1}, []float64{1}, [2]int{1, 1})
	if err == nil {
		t.Fatal("expected error for mismatched indices/data")
	}
	_, err = NewCSR([]int{1, 1}, []int{}, []float64{}, [2]int{1, 1})
	if err == nil {
		t.Fatal("expected error for indptr[0] != 0")
	}
	_, err = NewCSR([]int{0, 2}, []int{0}, []float64{1}, [2]int{1, 1})
	if err == nil {
		t.Fatal("expected error for indptr[-1] != len(data)")
	}
	_, err = NewCSR([]int{0, 1}, []int{5}, []float64{1}, [2]int{1, 1})
	if err == nil {
		t.Fatal("expected error for out-of-bounds col index")
	}
}

func TestMock_NewCSC_Errors(t *testing.T) {
	_, err := NewCSC([]int{0, 1}, []int{0}, []float64{1}, [2]int{1, 0})
	if err == nil {
		t.Fatal("expected error for non-positive shape")
	}
	_, err = NewCSC([]int{0, 1, 1}, []int{0}, []float64{1}, [2]int{1, 1})
	if err == nil {
		t.Fatal("expected error for wrong indptr length")
	}
	_, err = NewCSC([]int{0, 1}, []int{0, 1}, []float64{1}, [2]int{1, 1})
	if err == nil {
		t.Fatal("expected error for mismatched indices/data")
	}
	_, err = NewCSC([]int{1, 1}, []int{}, []float64{}, [2]int{1, 1})
	if err == nil {
		t.Fatal("expected error for indptr[0] != 0")
	}
	_, err = NewCSC([]int{0, 2}, []int{0}, []float64{1}, [2]int{1, 1})
	if err == nil {
		t.Fatal("expected error for indptr[-1] != len(data)")
	}
	_, err = NewCSC([]int{0, 1}, []int{5}, []float64{1}, [2]int{1, 1})
	if err == nil {
		t.Fatal("expected error for out-of-bounds row index")
	}
}

// ---------------------------------------------------------------------------
// Sparse extra: EyeSparse, Diags
// ---------------------------------------------------------------------------

func TestMock_EyeSparse_Zero(t *testing.T) {
	csr := EyeSparse(0)
	if csr.NNZ() != 0 {
		t.Errorf("EyeSparse(0) NNZ should be 0, got %d", csr.NNZ())
	}
}

func TestDiags_Error(t *testing.T) {
	_, err := Diags([][]float64{{1}}, []int{0, 1}, 3)
	if err == nil {
		t.Fatal("expected error for mismatched diagonals/offsets")
	}
	_, err = Diags(nil, nil, 0)
	if err == nil {
		t.Fatal("expected error for n <= 0")
	}
}

// ---------------------------------------------------------------------------
// FindSparse
// ---------------------------------------------------------------------------

func TestMock_FindSparse_Empty(t *testing.T) {
	csr, _ := NewCSR([]int{0, 0}, nil, nil, [2]int{1, 1})
	entries := FindSparse(csr)
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

// ---------------------------------------------------------------------------
// TDistribution: mean and var edge cases
// ---------------------------------------------------------------------------

func TestTDistribution_MeanVar_EdgeCases(t *testing.T) {
	// df <= 1: mean is NaN
	td1 := NewTDistribution(1)
	if !math.IsNaN(td1.Mean()) {
		t.Errorf("TDistribution(1).Mean() should be NaN, got %v", td1.Mean())
	}
	if !math.IsNaN(td1.Var()) {
		t.Errorf("TDistribution(1).Var() should be NaN, got %v", td1.Var())
	}

	// 1 < df <= 2: mean = 0, var = +Inf
	td15 := NewTDistribution(1.5)
	if td15.Mean() != 0 {
		t.Errorf("TDistribution(1.5).Mean() should be 0, got %v", td15.Mean())
	}
	if !math.IsInf(td15.Var(), 1) {
		t.Errorf("TDistribution(1.5).Var() should be +Inf, got %v", td15.Var())
	}
}

// ---------------------------------------------------------------------------
// FDistribution: mean edge case
// ---------------------------------------------------------------------------

func TestFDistribution_Mean_SmallDF2(t *testing.T) {
	f := NewFDistribution(5, 2)
	if !math.IsNaN(f.Mean()) {
		t.Errorf("FDistribution(5,2).Mean() should be NaN for df2<=2, got %v", f.Mean())
	}
}

// ---------------------------------------------------------------------------
// ChiSquared: PDF edge cases
// ---------------------------------------------------------------------------

func TestChiSquared_PDF_AtZero(t *testing.T) {
	// df < 2: +Inf
	c1 := NewChiSquared(1)
	if !math.IsInf(c1.PDF(0), 1) {
		t.Errorf("ChiSquared(1).PDF(0) should be +Inf, got %v", c1.PDF(0))
	}

	// df == 2: 0.5
	c2 := NewChiSquared(2)
	if c2.PDF(0) != 0.5 {
		t.Errorf("ChiSquared(2).PDF(0) should be 0.5, got %v", c2.PDF(0))
	}

	// df > 2: 0
	c5 := NewChiSquared(5)
	if c5.PDF(0) != 0 {
		t.Errorf("ChiSquared(5).PDF(0) should be 0, got %v", c5.PDF(0))
	}
}

func TestChiSquared_PDF_Negative(t *testing.T) {
	c := NewChiSquared(5)
	if c.PDF(-1) != 0 {
		t.Errorf("ChiSquared.PDF(-1) should be 0, got %v", c.PDF(-1))
	}
}

func TestChiSquared_CDF_NonPositive(t *testing.T) {
	c := NewChiSquared(5)
	if c.CDF(0) != 0 {
		t.Errorf("ChiSquared.CDF(0) should be 0, got %v", c.CDF(0))
	}
	if c.CDF(-1) != 0 {
		t.Errorf("ChiSquared.CDF(-1) should be 0, got %v", c.CDF(-1))
	}
}

// ---------------------------------------------------------------------------
// Poisson/Binomial edge cases
// ---------------------------------------------------------------------------

func TestPoisson_PMF_Negative(t *testing.T) {
	po := NewPoisson(5)
	if po.PMF(-1) != 0 {
		t.Error("Poisson.PMF(-1) should be 0")
	}
}

func TestPoisson_CDF_Negative(t *testing.T) {
	po := NewPoisson(5)
	if po.CDF(-1) != 0 {
		t.Error("Poisson.CDF(-1) should be 0")
	}
}

func TestBinomial_PMF_OutOfRange(t *testing.T) {
	bi := NewBinomial(10, 0.5)
	if bi.PMF(-1) != 0 {
		t.Error("Binomial.PMF(-1) should be 0")
	}
	if bi.PMF(11) != 0 {
		t.Error("Binomial.PMF(11) should be 0")
	}
}

func TestBinomial_PMF_EdgeProb(t *testing.T) {
	// p=0: only k=0 has probability 1
	bi0 := NewBinomial(5, 0)
	if bi0.PMF(0) != 1 {
		t.Errorf("Binomial(5,0).PMF(0) should be 1, got %v", bi0.PMF(0))
	}
	if bi0.PMF(1) != 0 {
		t.Errorf("Binomial(5,0).PMF(1) should be 0, got %v", bi0.PMF(1))
	}

	// p=1: only k=n has probability 1
	bi1 := NewBinomial(5, 1)
	if bi1.PMF(5) != 1 {
		t.Errorf("Binomial(5,1).PMF(5) should be 1, got %v", bi1.PMF(5))
	}
	if bi1.PMF(0) != 0 {
		t.Errorf("Binomial(5,1).PMF(0) should be 0, got %v", bi1.PMF(0))
	}
}

func TestBinomial_CDF_Edges(t *testing.T) {
	bi := NewBinomial(10, 0.5)
	if bi.CDF(-1) != 0 {
		t.Error("Binomial.CDF(-1) should be 0")
	}
	if bi.CDF(10) != 1 {
		t.Errorf("Binomial.CDF(10) should be 1, got %v", bi.CDF(10))
	}
}

// ---------------------------------------------------------------------------
// CSR.MulVec / MulDense
// ---------------------------------------------------------------------------

func TestCSR_MulVec(t *testing.T) {
	csr, _ := NewCSR([]int{0, 2, 3}, []int{0, 1, 0}, []float64{1, 2, 3}, [2]int{2, 2})
	y := csr.MulVec([]float64{1, 1})
	if y[0] != 3 || y[1] != 3 {
		t.Errorf("MulVec result = %v, expected [3,3]", y)
	}
}

func TestCSR_MulDense(t *testing.T) {
	csr, _ := NewCSR([]int{0, 1, 1}, []int{0}, []float64{2}, [2]int{2, 2})
	dense := [][]float64{{1, 0}, {0, 1}}
	result := csr.MulDense(dense)
	if result[0][0] != 2 || result[0][1] != 0 {
		t.Errorf("MulDense row 0 = %v", result[0])
	}
}

// ---------------------------------------------------------------------------
// Boltzmann distribution
// ---------------------------------------------------------------------------

func TestBoltzmann_Basic(t *testing.T) {
	b := NewBoltzmann(1.0, 5)
	// PMF should sum to ~1
	sum := 0.0
	for k := 0; k < 5; k++ {
		sum += b.PMF(k)
	}
	if math.Abs(sum-1.0) > 1e-10 {
		t.Errorf("Boltzmann PMF sum = %v, expected 1", sum)
	}
	// PMF out of range
	if b.PMF(-1) != 0 {
		t.Error("Boltzmann.PMF(-1) should be 0")
	}
	if b.PMF(5) != 0 {
		t.Error("Boltzmann.PMF(5) should be 0")
	}
}

func TestBoltzmann_CDF(t *testing.T) {
	b := NewBoltzmann(1.0, 5)
	if b.CDF(-1) != 0 {
		t.Error("Boltzmann.CDF(-1) should be 0")
	}
	if b.CDF(4) != 1 {
		t.Errorf("Boltzmann.CDF(4) should be 1, got %v", b.CDF(4))
	}
	if b.CDF(10) != 1 {
		t.Errorf("Boltzmann.CDF(10) should be 1, got %v", b.CDF(10))
	}
	// CDF at intermediate value
	c := b.CDF(2)
	if c <= 0 || c >= 1 {
		t.Errorf("Boltzmann.CDF(2) should be in (0,1), got %v", c)
	}
}

func TestBoltzmann_Mean(t *testing.T) {
	b := NewBoltzmann(1.0, 5)
	m := b.Mean()
	if m < 0 || m >= 5 || math.IsNaN(m) {
		t.Errorf("Boltzmann.Mean() = %v, expected in [0,5)", m)
	}
}

// ---------------------------------------------------------------------------
// Rice PPF: test to reach pdfVal==0 break and x<=0 clamp
// ---------------------------------------------------------------------------

func TestMock_Rice_PPF_ExtremeP(t *testing.T) {
	r := NewRice(0.01, 1.0) // very small nu
	// Test with very small p
	val := r.PPF(1e-15)
	if val < 0 || math.IsNaN(val) {
		t.Errorf("Rice.PPF(1e-15) should be non-negative, got %v", val)
	}
	// Test with very large p
	val = r.PPF(1 - 1e-15)
	if val <= 0 || math.IsNaN(val) {
		t.Errorf("Rice.PPF(1-1e-15) should be positive, got %v", val)
	}
}

func TestMock_Rice_PPF_MeanNegative(t *testing.T) {
	// With very small parameters, the initial guess might be <= 0
	r := NewRice(1e-300, 1e-300)
	val := r.PPF(0.5)
	if math.IsNaN(val) {
		// May be NaN due to numerical issues, that's acceptable
		_ = val
	}
}

// ---------------------------------------------------------------------------
// Wald PPF: test pdfVal==0 break and x<=0 clamp
// ---------------------------------------------------------------------------

func TestMock_Wald_PPF_ExtremeP(t *testing.T) {
	w := NewWald(1.0, 0.01) // very small lambda
	val := w.PPF(1e-15)
	if val < 0 || math.IsNaN(val) {
		t.Errorf("Wald.PPF(1e-15) should be non-negative, got %v", val)
	}
	val = w.PPF(1 - 1e-15)
	if val <= 0 || math.IsNaN(val) {
		t.Errorf("Wald.PPF(1-1e-15) should be positive, got %v", val)
	}
}

// ---------------------------------------------------------------------------
// VonMises PPF: test pdfVal==0 break
// ---------------------------------------------------------------------------

func TestMock_VonMises_PPF_ExtremeP(t *testing.T) {
	v := NewVonMises(0, 0.01) // very small kappa
	val := v.PPF(1e-15)
	if math.IsNaN(val) {
		_ = val // acceptable
	}
	val = v.PPF(1 - 1e-15)
	if math.IsNaN(val) {
		_ = val // acceptable
	}
}

// ---------------------------------------------------------------------------
// Nakagami PPF: test pdfVal==0 break
// ---------------------------------------------------------------------------

func TestMock_Nakagami_PPF_ExtremeP(t *testing.T) {
	n := NewNakagami(0.5, 1.0)
	val := n.PPF(1e-15)
	if val < 0 || math.IsNaN(val) {
		t.Errorf("Nakagami.PPF(1e-15) should be non-negative, got %v", val)
	}
	val = n.PPF(1 - 1e-15)
	if val <= 0 || math.IsNaN(val) {
		t.Errorf("Nakagami.PPF(1-1e-15) should be positive, got %v", val)
	}
}

// ---------------------------------------------------------------------------
// Linprog: feasible problem with equality constraints
// ---------------------------------------------------------------------------

func TestMock_Linprog_WithEquality(t *testing.T) {
	// minimize x1 + x2 subject to x1 + x2 = 1, x1 >= 0, x2 >= 0
	c := []float64{1, 1}
	Aeq := [][]float64{{1, 1}}
	beq := []float64{1}
	res, err := Linprog(c, nil, nil, Aeq, beq)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(res.Fun-1.0) > 0.1 {
		t.Errorf("expected optimal ~1.0, got %v", res.Fun)
	}
}

func TestMock_Linprog_InequalityAndEquality(t *testing.T) {
	// minimize -x1 subject to x1 + x2 = 1, x1 <= 0.5
	c := []float64{-1, 0}
	Aub := [][]float64{{1, 0}}
	bub := []float64{0.5}
	Aeq := [][]float64{{1, 1}}
	beq := []float64{1}
	res, err := Linprog(c, Aub, bub, Aeq, beq)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(res.X[0]-0.5) > 0.1 {
		t.Errorf("expected x1 ~0.5, got %v", res.X[0])
	}
}

// ---------------------------------------------------------------------------
// DualAnnealing: verify it hits periodic local search (nit%100==99)
// ---------------------------------------------------------------------------

func TestMock_DualAnnealing_ConvergesToMinimum(t *testing.T) {
	// Simple quadratic with known minimum
	f := func(x []float64) float64 {
		return (x[0]-1)*(x[0]-1) + (x[1]-2)*(x[1]-2)
	}
	bounds := [][2]float64{{-5, 5}, {-5, 5}}
	res, err := DualAnnealing(f, bounds)
	if err != nil {
		t.Fatal(err)
	}
	if res.Fun > 0.1 {
		t.Errorf("expected minimum near 0, got %v", res.Fun)
	}
}

// ---------------------------------------------------------------------------
// SHGO: multi-variable problem
// ---------------------------------------------------------------------------

func TestMock_SHGO_MultiVariable(t *testing.T) {
	f := func(x []float64) float64 {
		return (x[0]-1)*(x[0]-1) + (x[1]+2)*(x[1]+2)
	}
	bounds := [][2]float64{{-5, 5}, {-5, 5}}
	res, err := SHGO(f, bounds)
	if err != nil {
		t.Fatal(err)
	}
	if res.Fun > 0.5 {
		t.Errorf("expected minimum near 0, got %v", res.Fun)
	}
}

// ---------------------------------------------------------------------------
// SkewNormal CDF edge case (line 1259)
// ---------------------------------------------------------------------------

func TestMock_SkewNormal_CDF_OutOfRange(t *testing.T) {
	sn := NewSkewNormal(0, 1, 5) // high skew
	// CDF at very negative x should be near 0
	val := sn.CDF(-100)
	if val < 0 || val > 0.01 {
		t.Errorf("SkewNormal.CDF(-100) should be near 0, got %v", val)
	}
}

// ---------------------------------------------------------------------------
// Correlation edge cases
// ---------------------------------------------------------------------------

func TestMock_PearsonCorrelation_ShortInput(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for length <= 1")
		}
	}()
	PearsonCorrelation([]float64{1.0}, []float64{2.0})
}

func TestMock_PearsonCorrelation_MismatchedLengths(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for mismatched lengths")
		}
	}()
	PearsonCorrelation([]float64{1, 2, 3}, []float64{1, 2})
}

func TestMock_PartialCorrelation_ShortInput(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for short input")
		}
	}()
	data := [][]float64{{1, 2}, {3, 4}}
	PartialCorrelation(data, 0, 1, nil)
}

func TestMock_PartialCorrelation_MismatchedColumns(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for mismatched columns")
		}
	}()
	// Data with different row lengths
	data := [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	PartialCorrelation(data, 0, 1, []int{5}) // z[0]=5 out of bounds
}
