//go:build unit

package tabgo

import (
	"math"
	"testing"
)

func TestRollingBetaMatchesCovVar(t *testing.T) {
	asset := NewSeries("a", []any{1.0, 2.0, 3.0, 4.0, 5.0})
	bench := NewSeries("b", []any{2.0, 3.0, 5.0, 7.0, 11.0})

	betaResult := RollingBeta(asset, bench, 3)
	vals := betaResult.Values()

	if vals[0] != nil || vals[1] != nil {
		t.Error("RollingBeta: first 2 should be nil")
	}

	// Manually verify: beta = Cov(a,b) / Var(b)
	for i := 2; i < 5; i++ {
		wa := []float64{toFloat64(asset.Values()[i-2]), toFloat64(asset.Values()[i-1]), toFloat64(asset.Values()[i])}
		wb := []float64{toFloat64(bench.Values()[i-2]), toFloat64(bench.Values()[i-1]), toFloat64(bench.Values()[i])}
		cov := covariance(wa, wb)
		varB := variance(wb)
		expectedBeta := cov / varB

		gotBeta := toFloat64(vals[i])
		if !almostEqual(gotBeta, expectedBeta, 1e-10) {
			t.Errorf("RollingBeta[%d] = %v, manual Cov/Var = %v", i, gotBeta, expectedBeta)
		}
	}
}

func TestRollingBetaPerfectCorrelation(t *testing.T) {
	// asset = 2 * benchmark, beta should be 2
	asset := NewSeries("a", []any{2.0, 4.0, 6.0, 8.0, 10.0})
	bench := NewSeries("b", []any{1.0, 2.0, 3.0, 4.0, 5.0})
	result := RollingBeta(asset, bench, 3)
	vals := result.Values()

	for i := 2; i < 5; i++ {
		v := toFloat64(vals[i])
		if !almostEqual(v, 2.0, 1e-10) {
			t.Errorf("RollingBeta[%d] = %v, want 2.0", i, v)
		}
	}
}

func TestRollingBetaZeroVarianceBenchmark(t *testing.T) {
	asset := NewSeries("a", []any{1.0, 2.0, 3.0})
	bench := NewSeries("b", []any{5.0, 5.0, 5.0})
	result := RollingBeta(asset, bench, 3)
	vals := result.Values()
	// Var(b) = 0, beta should be 0
	if v := toFloat64(vals[2]); v != 0.0 {
		t.Errorf("RollingBeta zero var = %v, want 0.0", v)
	}
}

func TestRollingAlpha(t *testing.T) {
	asset := NewSeries("a", []any{3.0, 5.0, 7.0, 9.0, 11.0})
	bench := NewSeries("b", []any{1.0, 2.0, 3.0, 4.0, 5.0})

	result := RollingAlpha(asset, bench, 3)
	vals := result.Values()

	if vals[0] != nil || vals[1] != nil {
		t.Error("RollingAlpha: first 2 should be nil")
	}

	// asset = 2*bench + 1, so alpha should be 1
	for i := 2; i < 5; i++ {
		v := toFloat64(vals[i])
		if !almostEqual(v, 1.0, 1e-10) {
			t.Errorf("RollingAlpha[%d] = %v, want 1.0", i, v)
		}
	}
}

func TestRollingAlphaZeroVariance(t *testing.T) {
	asset := NewSeries("a", []any{3.0, 5.0, 7.0})
	bench := NewSeries("b", []any{5.0, 5.0, 5.0})
	result := RollingAlpha(asset, bench, 3)
	vals := result.Values()
	// Var(b) = 0, alpha = mean(a)
	if v := toFloat64(vals[2]); !almostEqual(v, 5.0, 1e-10) {
		t.Errorf("RollingAlpha zero var = %v, want 5.0 (mean of a)", v)
	}
}

func TestRollingR2(t *testing.T) {
	asset := NewSeries("a", []any{2.0, 4.0, 6.0, 8.0, 10.0})
	bench := NewSeries("b", []any{1.0, 2.0, 3.0, 4.0, 5.0})
	result := RollingR2(asset, bench, 3)
	vals := result.Values()

	if vals[0] != nil || vals[1] != nil {
		t.Error("RollingR2: first 2 should be nil")
	}

	// Perfect correlation -> R2 = 1
	for i := 2; i < 5; i++ {
		v := toFloat64(vals[i])
		if !almostEqual(v, 1.0, 1e-10) {
			t.Errorf("RollingR2[%d] = %v, want 1.0", i, v)
		}
	}
}

func TestRollingR2Range(t *testing.T) {
	asset := NewSeries("a", []any{1.0, 3.0, 2.0, 5.0, 4.0})
	bench := NewSeries("b", []any{5.0, 1.0, 3.0, 2.0, 4.0})
	result := RollingR2(asset, bench, 3)
	vals := result.Values()
	for i := 2; i < 5; i++ {
		v := toFloat64(vals[i])
		if v < -1e-10 || v > 1.0+1e-10 {
			t.Errorf("RollingR2[%d] = %v, out of [0,1]", i, v)
		}
	}
}

func TestRollingResidual(t *testing.T) {
	// Perfect fit: asset = 2*bench + 1, residual should be ~0
	asset := NewSeries("a", []any{3.0, 5.0, 7.0, 9.0, 11.0})
	bench := NewSeries("b", []any{1.0, 2.0, 3.0, 4.0, 5.0})
	result := RollingResidual(asset, bench, 3)
	vals := result.Values()

	if vals[0] != nil || vals[1] != nil {
		t.Error("RollingResidual: first 2 should be nil")
	}

	for i := 2; i < 5; i++ {
		v := toFloat64(vals[i])
		if !almostEqual(v, 0.0, 1e-10) {
			t.Errorf("RollingResidual[%d] = %v, want ~0", i, v)
		}
	}
}

func TestRollingResidualZeroVariance(t *testing.T) {
	asset := NewSeries("a", []any{3.0, 5.0, 7.0})
	bench := NewSeries("b", []any{5.0, 5.0, 5.0})
	result := RollingResidual(asset, bench, 3)
	vals := result.Values()
	// beta=0, alpha=mean(a)=5, residual = a[2] - alpha = 7 - 5 = 2
	if v := toFloat64(vals[2]); !almostEqual(v, 2.0, 1e-10) {
		t.Errorf("RollingResidual zero var = %v, want 2.0", v)
	}
}

func TestRollingMultiBeta(t *testing.T) {
	// asset = 2*f1 + 3*f2 (with independent factors)
	asset := NewSeries("y", []any{
		2*1.0 + 3*10.0,
		2*2.0 + 3*8.0,
		2*3.0 + 3*12.0,
		2*4.0 + 3*7.0,
		2*5.0 + 3*11.0,
	})
	factors := NewDataFrameFromRows(
		[]string{"f1", "f2"},
		[][]any{
			{1.0, 10.0},
			{2.0, 8.0},
			{3.0, 12.0},
			{4.0, 7.0},
			{5.0, 11.0},
		},
	)
	result := RollingMultiBeta(asset, factors, 3)

	f1Vals := result.Column("f1").Values()
	f2Vals := result.Column("f2").Values()

	if f1Vals[0] != nil || f1Vals[1] != nil {
		t.Error("RollingMultiBeta: first 2 should be nil")
	}

	for i := 2; i < 5; i++ {
		b1 := toFloat64(f1Vals[i])
		b2 := toFloat64(f2Vals[i])
		if !almostEqual(b1, 2.0, 0.1) {
			t.Errorf("RollingMultiBeta[%d] f1 beta = %v, want ~2.0", i, b1)
		}
		if !almostEqual(b2, 3.0, 0.1) {
			t.Errorf("RollingMultiBeta[%d] f2 beta = %v, want ~3.0", i, b2)
		}
	}
}

func TestRollingMultiBetaIndependent(t *testing.T) {
	// Use truly independent factors
	asset := NewSeries("y", []any{5.0, 8.0, 11.0, 14.0, 17.0})
	factors := NewDataFrameFromRows(
		[]string{"f1", "f2"},
		[][]any{
			{1.0, 1.0},
			{2.0, 2.0},
			{3.0, 3.0},
			{4.0, 4.0},
			{5.0, 5.0},
		},
	)
	result := RollingMultiBeta(asset, factors, 3)

	// asset = 3*f1 + 0*f2 + 2 (approximately)
	// But f1 and f2 are identical, so multicollinear
	// Just check result is not nil
	f1Vals := result.Column("f1").Values()
	for i := 2; i < 5; i++ {
		if f1Vals[i] == nil {
			t.Errorf("RollingMultiBeta[%d] should not be nil", i)
		}
	}
}

func TestRollingMultiBetaSingleFactor(t *testing.T) {
	// With a single factor, multi-beta should match single beta
	asset := NewSeries("a", []any{2.0, 4.0, 6.0, 8.0, 10.0})
	bench := NewSeries("b", []any{1.0, 2.0, 3.0, 4.0, 5.0})
	factors := NewDataFrameFromRows(
		[]string{"b"},
		[][]any{{1.0}, {2.0}, {3.0}, {4.0}, {5.0}},
	)

	singleBeta := RollingBeta(asset, bench, 3)
	multiBeta := RollingMultiBeta(asset, factors, 3)

	singleVals := singleBeta.Values()
	multiVals := multiBeta.Column("b").Values()

	for i := 2; i < 5; i++ {
		sv := toFloat64(singleVals[i])
		mv := toFloat64(multiVals[i])
		if !almostEqual(sv, mv, 1e-10) {
			t.Errorf("SingleBeta[%d]=%v != MultiBeta[%d]=%v", i, sv, i, mv)
		}
	}
}

func TestRollingBetaUnequalLength(t *testing.T) {
	asset := NewSeries("a", []any{1.0, 2.0, 3.0})
	bench := NewSeries("b", []any{2.0, 4.0})
	result := RollingBeta(asset, bench, 2)
	if result.Len() != 2 {
		t.Errorf("RollingBeta unequal: got len=%d, want 2", result.Len())
	}
}

func TestRollingAlphaUnequalLength(t *testing.T) {
	asset := NewSeries("a", []any{1.0, 2.0, 3.0})
	bench := NewSeries("b", []any{2.0, 4.0})
	result := RollingAlpha(asset, bench, 2)
	if result.Len() != 2 {
		t.Errorf("RollingAlpha unequal: got len=%d, want 2", result.Len())
	}
}

func TestRollingR2UnequalLength(t *testing.T) {
	asset := NewSeries("a", []any{1.0, 2.0, 3.0})
	bench := NewSeries("b", []any{2.0, 4.0})
	result := RollingR2(asset, bench, 2)
	if result.Len() != 2 {
		t.Errorf("RollingR2 unequal: got len=%d, want 2", result.Len())
	}
}

func TestRollingResidualUnequalLength(t *testing.T) {
	asset := NewSeries("a", []any{1.0, 2.0, 3.0})
	bench := NewSeries("b", []any{2.0, 4.0})
	result := RollingResidual(asset, bench, 2)
	if result.Len() != 2 {
		t.Errorf("RollingResidual unequal: got len=%d, want 2", result.Len())
	}
}

func TestRollingBetaNameFormat(t *testing.T) {
	asset := NewSeries("stock", []any{1.0, 2.0, 3.0})
	bench := NewSeries("market", []any{2.0, 4.0, 6.0})
	result := RollingBeta(asset, bench, 2)
	if result.Name() != "stock_beta" {
		t.Errorf("RollingBeta name = %q, want %q", result.Name(), "stock_beta")
	}
}

func TestRollingAlphaNameFormat(t *testing.T) {
	asset := NewSeries("stock", []any{1.0, 2.0, 3.0})
	bench := NewSeries("market", []any{2.0, 4.0, 6.0})
	result := RollingAlpha(asset, bench, 2)
	if result.Name() != "stock_alpha" {
		t.Errorf("RollingAlpha name = %q, want %q", result.Name(), "stock_alpha")
	}
}

func TestGaussianSolveSingular(t *testing.T) {
	// Singular matrix: should return zeros
	A := []float64{0, 0, 0, 0}
	b := []float64{1, 1}
	x := gaussianSolve(A, b, 2)
	for i, v := range x {
		if v != 0 {
			t.Errorf("gaussianSolve singular x[%d] = %v, want 0", i, v)
		}
	}
}

func TestGaussianSolveIdentity(t *testing.T) {
	// Identity * x = b -> x = b
	A := []float64{1, 0, 0, 1}
	b := []float64{3, 7}
	x := gaussianSolve(A, b, 2)
	if !almostEqual(x[0], 3, 1e-10) || !almostEqual(x[1], 7, 1e-10) {
		t.Errorf("gaussianSolve identity: got %v, want [3, 7]", x)
	}
}

func TestGaussianSolveWithPivoting(t *testing.T) {
	// Matrix that needs pivoting: first element is 0
	A := []float64{0, 1, 1, 0}
	b := []float64{4, 5}
	x := gaussianSolve(A, b, 2)
	if !almostEqual(x[0], 5, 1e-10) || !almostEqual(x[1], 4, 1e-10) {
		t.Errorf("gaussianSolve pivot: got %v, want [5, 4]", x)
	}
}

func TestRollingBetaConsistencyWithAlpha(t *testing.T) {
	// Verify: mean(a) = alpha + beta * mean(b)
	asset := NewSeries("a", []any{1.0, 3.0, 2.0, 5.0, 4.0})
	bench := NewSeries("b", []any{5.0, 1.0, 3.0, 2.0, 4.0})

	betas := RollingBeta(asset, bench, 3)
	alphas := RollingAlpha(asset, bench, 3)

	for i := 2; i < 5; i++ {
		beta := toFloat64(betas.Values()[i])
		alpha := toFloat64(alphas.Values()[i])
		wa := []float64{toFloat64(asset.Values()[i-2]), toFloat64(asset.Values()[i-1]), toFloat64(asset.Values()[i])}
		wb := []float64{toFloat64(bench.Values()[i-2]), toFloat64(bench.Values()[i-1]), toFloat64(bench.Values()[i])}
		mA := mean(wa)
		mB := mean(wb)
		predicted := alpha + beta*mB
		if !almostEqual(mA, predicted, 1e-10) {
			t.Errorf("Alpha+Beta*mean(b)[%d] = %v, want mean(a) = %v", i, predicted, mA)
		}
	}
}

func TestRollingResidualSumZero(t *testing.T) {
	// For OLS, sum of residuals in-sample should be 0
	// But we only see the last residual per window.
	// Test: for a window, the predicted value at the last point = actual - residual
	asset := NewSeries("a", []any{1.0, 3.0, 2.0, 5.0, 4.0})
	bench := NewSeries("b", []any{5.0, 1.0, 3.0, 2.0, 4.0})

	residuals := RollingResidual(asset, bench, 3)
	betas := RollingBeta(asset, bench, 3)
	alphas := RollingAlpha(asset, bench, 3)

	for i := 2; i < 5; i++ {
		resid := toFloat64(residuals.Values()[i])
		beta := toFloat64(betas.Values()[i])
		alpha := toFloat64(alphas.Values()[i])
		a := toFloat64(asset.Values()[i])
		b := toFloat64(bench.Values()[i])
		predicted := alpha + beta*b
		expectedResid := a - predicted
		if !almostEqual(resid, expectedResid, 1e-10) {
			t.Errorf("Residual[%d] = %v, want %v", i, resid, expectedResid)
		}
	}
}

func TestRollingBetaWithNils(t *testing.T) {
	asset := NewSeries("a", []any{nil, nil, 3.0})
	bench := NewSeries("b", []any{nil, nil, 6.0})
	result := RollingBeta(asset, bench, 2)
	vals := result.Values()
	if vals[1] != nil {
		t.Errorf("RollingBeta nils: expected nil, got %v", vals[1])
	}
}

func TestRollingAlphaWithNils(t *testing.T) {
	asset := NewSeries("a", []any{nil, nil, 3.0})
	bench := NewSeries("b", []any{nil, nil, 6.0})
	result := RollingAlpha(asset, bench, 2)
	vals := result.Values()
	if vals[1] != nil {
		t.Errorf("RollingAlpha nils: expected nil, got %v", vals[1])
	}
}

func TestRollingR2WithNils(t *testing.T) {
	asset := NewSeries("a", []any{nil, nil, 3.0})
	bench := NewSeries("b", []any{nil, nil, 6.0})
	result := RollingR2(asset, bench, 2)
	vals := result.Values()
	if vals[1] != nil {
		t.Errorf("RollingR2 nils: expected nil, got %v", vals[1])
	}
}

func TestRollingResidualWithNils(t *testing.T) {
	asset := NewSeries("a", []any{nil, nil, 3.0})
	bench := NewSeries("b", []any{nil, nil, 6.0})
	result := RollingResidual(asset, bench, 2)
	vals := result.Values()
	if vals[1] != nil {
		t.Errorf("RollingResidual nils: expected nil, got %v", vals[1])
	}
}

func TestRollingMultiBetaWithNils(t *testing.T) {
	asset := NewSeries("y", []any{nil, nil, 3.0})
	factors := NewDataFrameFromRows(
		[]string{"f1"},
		[][]any{{nil}, {nil}, {6.0}},
	)
	result := RollingMultiBeta(asset, factors, 2)
	vals := result.Column("f1").Values()
	if vals[1] != nil {
		t.Errorf("RollingMultiBeta nils: expected nil, got %v", vals[1])
	}
}

// Dummy usage to avoid import errors
var _ = math.Abs
