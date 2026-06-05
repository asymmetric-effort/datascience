//go:build unit

package tabgo

import (
	"math"
	"testing"
)

func TestRollingPCABasic(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b"},
		[][]any{
			{1.0, 2.0},
			{2.0, 4.0},
			{3.0, 6.0},
			{4.0, 8.0},
			{5.0, 10.0},
		},
	)
	result := RollingPCA(df, 3, 2)

	// First two should be nil
	if result.Components[0] != nil || result.Components[1] != nil {
		t.Error("RollingPCA: first 2 components should be nil")
	}
	if result.Eigenvalues[0] != nil || result.Eigenvalues[1] != nil {
		t.Error("RollingPCA: first 2 eigenvalues should be nil")
	}
	if result.ExplainedVar[0] != nil || result.ExplainedVar[1] != nil {
		t.Error("RollingPCA: first 2 explained_var should be nil")
	}

	// From index 2 onward, results should exist
	for i := 2; i < 5; i++ {
		if result.Components[i] == nil {
			t.Errorf("RollingPCA.Components[%d] should not be nil", i)
		}
		if result.Eigenvalues[i] == nil {
			t.Errorf("RollingPCA.Eigenvalues[%d] should not be nil", i)
		}
		if result.ExplainedVar[i] == nil {
			t.Errorf("RollingPCA.ExplainedVar[%d] should not be nil", i)
		}
	}
}

func TestRollingPCAEigenvaluesSumToTotalVariance(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b", "c"},
		[][]any{
			{1.0, 5.0, 3.0},
			{2.0, 3.0, 7.0},
			{4.0, 8.0, 1.0},
			{3.0, 2.0, 5.0},
			{5.0, 6.0, 2.0},
		},
	)
	result := RollingPCA(df, 3, 3)

	cols := dfNumericColumns(df)
	colVals := make([][]any, len(cols))
	for ci, c := range cols {
		colVals[ci] = df.Column(c).Values()
	}

	for row := 2; row < 5; row++ {
		// Compute total variance = sum of variances of each column in window
		totalVar := 0.0
		for ci := range cols {
			w := extractWindow(colVals[ci], row, 3)
			totalVar += variance(w)
		}

		// Sum of eigenvalues should equal total variance
		eigVals := result.Eigenvalues[row].Values()
		eigSum := 0.0
		for _, ev := range eigVals {
			eigSum += toFloat64(ev)
		}

		if !almostEqual(eigSum, totalVar, 0.1) {
			t.Errorf("RollingPCA[%d]: eigenvalue sum=%v, total variance=%v", row, eigSum, totalVar)
		}
	}
}

func TestRollingPCAExplainedVarSumsToOne(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b"},
		[][]any{
			{1.0, 5.0},
			{3.0, 2.0},
			{5.0, 8.0},
			{2.0, 1.0},
			{4.0, 6.0},
		},
	)
	result := RollingPCA(df, 3, 2)

	for row := 2; row < 5; row++ {
		evVals := result.ExplainedVar[row].Values()
		sum := 0.0
		for _, v := range evVals {
			sum += toFloat64(v)
		}
		if !almostEqual(sum, 1.0, 0.01) {
			t.Errorf("RollingPCA[%d]: explained var sum=%v, want 1.0", row, sum)
		}
	}
}

func TestRollingPCAComponentsOrthogonal(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b"},
		[][]any{
			{1.0, 5.0},
			{3.0, 2.0},
			{5.0, 8.0},
			{2.0, 1.0},
		},
	)
	result := RollingPCA(df, 3, 2)

	for row := 2; row < 4; row++ {
		compDf := result.Components[row]
		cols := compDf.Columns()
		nCols := len(cols)

		// Extract component vectors
		nComp := compDf.Len()
		if nComp < 2 {
			continue
		}

		// Get first two component vectors
		v1 := make([]float64, nCols)
		v2 := make([]float64, nCols)
		for ci, c := range cols {
			vals := compDf.Column(c).Values()
			v1[ci] = toFloat64(vals[0])
			v2[ci] = toFloat64(vals[1])
		}

		// Dot product should be ~0
		dot := 0.0
		for i := range v1 {
			dot += v1[i] * v2[i]
		}
		if !almostEqual(dot, 0.0, 0.01) {
			t.Errorf("RollingPCA[%d]: components not orthogonal, dot=%v", row, dot)
		}
	}
}

func TestRollingPCANComponentsCapped(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b"},
		[][]any{
			{1.0, 2.0},
			{3.0, 4.0},
			{5.0, 6.0},
		},
	)
	// Request 5 components but only 2 columns
	result := RollingPCA(df, 3, 5)
	eigVals := result.Eigenvalues[2]
	if eigVals.Len() != 2 {
		t.Errorf("RollingPCA nComponents capped: got %d, want 2", eigVals.Len())
	}
}

func TestRollingPCANComponentsZero(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b"},
		[][]any{
			{1.0, 2.0},
			{3.0, 4.0},
			{5.0, 6.0},
		},
	)
	// nComponents=0 should default to nCols
	result := RollingPCA(df, 3, 0)
	eigVals := result.Eigenvalues[2]
	if eigVals.Len() != 2 {
		t.Errorf("RollingPCA nComponents=0: got %d, want 2", eigVals.Len())
	}
}

func TestRollingPCASingleComponent(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b", "c"},
		[][]any{
			{1.0, 2.0, 3.0},
			{4.0, 5.0, 6.0},
			{7.0, 8.0, 9.0},
		},
	)
	result := RollingPCA(df, 3, 1)
	eigVals := result.Eigenvalues[2]
	if eigVals.Len() != 1 {
		t.Errorf("RollingPCA 1 component: got %d eigenvalues", eigVals.Len())
	}

	// First eigenvalue should explain the most variance
	ev := toFloat64(result.ExplainedVar[2].Values()[0])
	if ev < 0.5 {
		t.Errorf("RollingPCA first component explains only %v, want > 0.5", ev)
	}
}

func TestRollingPCAEigenvaluesDescending(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b", "c"},
		[][]any{
			{1.0, 5.0, 3.0},
			{2.0, 3.0, 7.0},
			{4.0, 8.0, 1.0},
			{3.0, 2.0, 5.0},
		},
	)
	result := RollingPCA(df, 3, 3)

	for row := 2; row < 4; row++ {
		eigVals := result.Eigenvalues[row].Values()
		for i := 0; i < len(eigVals)-1; i++ {
			if toFloat64(eigVals[i]) < toFloat64(eigVals[i+1])-1e-10 {
				t.Errorf("RollingPCA[%d]: eigenvalues not descending: %v > %v",
					row, toFloat64(eigVals[i]), toFloat64(eigVals[i+1]))
			}
		}
	}
}

func TestRollingPCAExplainedVarNonNegative(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b"},
		[][]any{
			{1.0, 5.0},
			{3.0, 2.0},
			{5.0, 8.0},
		},
	)
	result := RollingPCA(df, 3, 2)
	evVals := result.ExplainedVar[2].Values()
	for i, v := range evVals {
		if toFloat64(v) < -1e-10 {
			t.Errorf("RollingPCA explained_var[%d] = %v, want >= 0", i, v)
		}
	}
}

func TestJacobiEigenIdentity(t *testing.T) {
	mat := [][]float64{{1, 0}, {0, 1}}
	vals, _ := jacobiEigen(mat, 2)
	if !almostEqual(vals[0], 1.0, 1e-10) || !almostEqual(vals[1], 1.0, 1e-10) {
		t.Errorf("jacobiEigen identity: vals=%v, want [1, 1]", vals)
	}
}

func TestJacobiEigenSymmetric(t *testing.T) {
	mat := [][]float64{{2, 1}, {1, 2}}
	vals, vecs := jacobiEigen(mat, 2)

	// Eigenvalues should be 1 and 3
	sum := vals[0] + vals[1]
	prod := vals[0] * vals[1]
	if !almostEqual(sum, 4.0, 1e-10) {
		t.Errorf("jacobiEigen: eigenvalue sum=%v, want 4", sum)
	}
	if !almostEqual(prod, 3.0, 1e-10) {
		t.Errorf("jacobiEigen: eigenvalue product=%v, want 3", prod)
	}

	// Eigenvectors should be orthogonal
	dot := 0.0
	for i := 0; i < 2; i++ {
		dot += vecs[i][0] * vecs[i][1]
	}
	if !almostEqual(dot, 0.0, 1e-10) {
		t.Errorf("jacobiEigen: eigenvectors not orthogonal, dot=%v", dot)
	}
}

func TestJacobiEigenDiagonal(t *testing.T) {
	mat := [][]float64{{3, 0, 0}, {0, 1, 0}, {0, 0, 2}}
	vals, _ := jacobiEigen(mat, 3)
	// Should return 3, 1, 2 (not necessarily sorted)
	sum := vals[0] + vals[1] + vals[2]
	if !almostEqual(sum, 6.0, 1e-10) {
		t.Errorf("jacobiEigen diagonal: sum=%v, want 6", sum)
	}
}

func TestRollingPCAConstantData(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b"},
		[][]any{
			{1.0, 2.0},
			{1.0, 2.0},
			{1.0, 2.0},
		},
	)
	result := RollingPCA(df, 3, 2)
	// All eigenvalues should be 0 (no variance)
	eigVals := result.Eigenvalues[2].Values()
	for i, v := range eigVals {
		if !almostEqual(toFloat64(v), 0.0, 1e-10) {
			t.Errorf("RollingPCA constant eigenvalue[%d] = %v, want 0", i, toFloat64(v))
		}
	}
}

func TestRollingPCAResultStructure(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b", "c"},
		[][]any{
			{1.0, 2.0, 3.0},
			{4.0, 5.0, 6.0},
			{7.0, 8.0, 9.0},
			{10.0, 11.0, 12.0},
		},
	)
	result := RollingPCA(df, 3, 2)

	if len(result.Components) != 4 {
		t.Errorf("Components length=%d, want 4", len(result.Components))
	}
	if len(result.ExplainedVar) != 4 {
		t.Errorf("ExplainedVar length=%d, want 4", len(result.ExplainedVar))
	}
	if len(result.Eigenvalues) != 4 {
		t.Errorf("Eigenvalues length=%d, want 4", len(result.Eigenvalues))
	}

	// Components at valid rows should have 2 rows (nComponents) and 3 columns (nFeatures)
	for i := 2; i < 4; i++ {
		compDf := result.Components[i]
		if compDf.Len() != 2 {
			t.Errorf("Components[%d] rows=%d, want 2", i, compDf.Len())
		}
		if len(compDf.Columns()) != 3 {
			t.Errorf("Components[%d] cols=%d, want 3", i, len(compDf.Columns()))
		}
	}
}

func TestRollingPCAComponentUnitNorm(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b"},
		[][]any{
			{1.0, 5.0},
			{3.0, 2.0},
			{5.0, 8.0},
			{2.0, 1.0},
		},
	)
	result := RollingPCA(df, 3, 2)

	for row := 2; row < 4; row++ {
		compDf := result.Components[row]
		cols := compDf.Columns()
		nComp := compDf.Len()

		for pc := 0; pc < nComp; pc++ {
			norm := 0.0
			for _, c := range cols {
				v := toFloat64(compDf.Column(c).Values()[pc])
				norm += v * v
			}
			norm = math.Sqrt(norm)
			if !almostEqual(norm, 1.0, 0.01) {
				t.Errorf("RollingPCA[%d] component %d norm=%v, want 1.0", row, pc, norm)
			}
		}
	}
}

func TestRollingPCAThreeByThree(t *testing.T) {
	// 3x3 with known structure: first two vars correlated, third independent
	df := NewDataFrameFromRows(
		[]string{"a", "b", "c"},
		[][]any{
			{1.0, 2.0, 10.0},
			{2.0, 4.0, 5.0},
			{3.0, 6.0, 8.0},
			{4.0, 8.0, 3.0},
			{5.0, 10.0, 7.0},
		},
	)
	result := RollingPCA(df, 3, 3)

	// Verify eigenvalue sum = total variance for each window
	for row := 2; row < 5; row++ {
		cols := dfNumericColumns(df)
		colVals := make([][]any, len(cols))
		for ci, c := range cols {
			colVals[ci] = df.Column(c).Values()
		}

		totalVar := 0.0
		for ci := range cols {
			w := extractWindow(colVals[ci], row, 3)
			totalVar += variance(w)
		}

		eigVals := result.Eigenvalues[row].Values()
		eigSum := 0.0
		for _, ev := range eigVals {
			eigSum += toFloat64(ev)
		}

		if !almostEqual(eigSum, totalVar, 0.2) {
			t.Errorf("RollingPCA 3x3[%d]: eigenvalue sum=%v, total var=%v", row, eigSum, totalVar)
		}
	}
}
