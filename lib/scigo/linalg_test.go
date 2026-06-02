//go:build unit

package scigo

import (
	"math"
	"testing"
)

// ---------------------------------------------------------------------------
// LU Decomposition Tests
// ---------------------------------------------------------------------------

func TestLU(t *testing.T) {
	a := [][]float64{
		{2, 1, 1},
		{4, 3, 3},
		{8, 7, 9},
	}
	p, l, u, err := LU(a)
	if err != nil {
		t.Fatalf("LU: unexpected error: %v", err)
	}

	// Verify P*A = L*U.
	n := len(a)
	pa := matMul(p, a, n)
	lu := matMul(l, u, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if !approxEqual(pa[i][j], lu[i][j], 1e-10) {
				t.Errorf("LU: P*A[%d][%d] = %v, L*U[%d][%d] = %v", i, j, pa[i][j], i, j, lu[i][j])
			}
		}
	}

	// Verify L is lower-triangular with ones on diagonal.
	for i := 0; i < n; i++ {
		if !approxEqual(l[i][i], 1, 1e-14) {
			t.Errorf("LU: L[%d][%d] = %v, want 1", i, i, l[i][i])
		}
		for j := i + 1; j < n; j++ {
			if l[i][j] != 0 {
				t.Errorf("LU: L[%d][%d] = %v, want 0", i, j, l[i][j])
			}
		}
	}

	// Verify U is upper-triangular.
	for i := 0; i < n; i++ {
		for j := 0; j < i; j++ {
			if !approxEqual(u[i][j], 0, 1e-10) {
				t.Errorf("LU: U[%d][%d] = %v, want 0", i, j, u[i][j])
			}
		}
	}
}

func TestLUFactor(t *testing.T) {
	a := [][]float64{
		{1, 2},
		{3, 4},
	}
	lu, piv, err := LUFactor(a)
	if err != nil {
		t.Fatalf("LUFactor: unexpected error: %v", err)
	}
	if lu == nil || piv == nil {
		t.Fatal("LUFactor: nil result")
	}
}

func TestLUSolve(t *testing.T) {
	a := [][]float64{
		{2, 1, 1},
		{4, 3, 3},
		{8, 7, 9},
	}
	b := []float64{1, 1, 1}
	lu, piv, err := LUFactor(a)
	if err != nil {
		t.Fatalf("LUFactor: %v", err)
	}
	x, err := LUSolve(lu, piv, b)
	if err != nil {
		t.Fatalf("LUSolve: %v", err)
	}

	// Verify A*x = b.
	for i := 0; i < 3; i++ {
		sum := 0.0
		for j := 0; j < 3; j++ {
			sum += a[i][j] * x[j]
		}
		if !approxEqual(sum, b[i], 1e-10) {
			t.Errorf("LUSolve: A*x[%d] = %v, want %v", i, sum, b[i])
		}
	}
}

func TestLUSingular(t *testing.T) {
	a := [][]float64{
		{1, 2},
		{2, 4},
	}
	_, _, err := LUFactor(a)
	if err == nil {
		t.Error("LUFactor: expected error for singular matrix")
	}
}

// ---------------------------------------------------------------------------
// Cholesky Tests
// ---------------------------------------------------------------------------

func TestChoFactor(t *testing.T) {
	// Symmetric positive-definite matrix.
	a := [][]float64{
		{4, 2},
		{2, 3},
	}
	l, err := ChoFactor(a)
	if err != nil {
		t.Fatalf("ChoFactor: %v", err)
	}

	// Verify L * L^T = A.
	n := 2
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			sum := 0.0
			for k := 0; k < n; k++ {
				sum += l[i][k] * l[j][k]
			}
			if !approxEqual(sum, a[i][j], 1e-10) {
				t.Errorf("ChoFactor: L*L^T[%d][%d] = %v, want %v", i, j, sum, a[i][j])
			}
		}
	}
}

func TestChoFactorNotPD(t *testing.T) {
	a := [][]float64{
		{1, 2},
		{2, 1},
	}
	_, err := ChoFactor(a)
	if err == nil {
		t.Error("ChoFactor: expected error for non-positive-definite matrix")
	}
}

func TestChoSolve(t *testing.T) {
	a := [][]float64{
		{4, 2},
		{2, 3},
	}
	b := []float64{1, 2}
	cho, err := ChoFactor(a)
	if err != nil {
		t.Fatalf("ChoFactor: %v", err)
	}
	x, err := ChoSolve(cho, b)
	if err != nil {
		t.Fatalf("ChoSolve: %v", err)
	}

	// Verify A*x = b.
	for i := 0; i < 2; i++ {
		sum := 0.0
		for j := 0; j < 2; j++ {
			sum += a[i][j] * x[j]
		}
		if !approxEqual(sum, b[i], 1e-10) {
			t.Errorf("ChoSolve: A*x[%d] = %v, want %v", i, sum, b[i])
		}
	}
}

// ---------------------------------------------------------------------------
// Schur Tests
// ---------------------------------------------------------------------------

func TestSchur(t *testing.T) {
	// Symmetric matrix.
	a := [][]float64{
		{2, 1},
		{1, 3},
	}
	tt, z, err := Schur(a)
	if err != nil {
		t.Fatalf("Schur: %v", err)
	}

	// Verify A = Z*T*Z^T.
	n := 2
	zt := make([][]float64, n)
	for i := 0; i < n; i++ {
		zt[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			zt[i][j] = z[j][i]
		}
	}
	ztt := matMul(z, tt, n)
	result := matMul(ztt, zt, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if !approxEqual(result[i][j], a[i][j], 1e-6) {
				t.Errorf("Schur: Z*T*Z^T[%d][%d] = %v, want %v", i, j, result[i][j], a[i][j])
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Hessenberg Tests
// ---------------------------------------------------------------------------

func TestHessenberg(t *testing.T) {
	a := [][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	h, q, err := Hessenberg(a)
	if err != nil {
		t.Fatalf("Hessenberg: %v", err)
	}

	// Verify A = Q*H*Q^T.
	n := 3
	qt := make([][]float64, n)
	for i := 0; i < n; i++ {
		qt[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			qt[i][j] = q[j][i]
		}
	}
	qh := matMul(q, h, n)
	result := matMul(qh, qt, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if !approxEqual(result[i][j], a[i][j], 1e-10) {
				t.Errorf("Hessenberg: Q*H*Q^T[%d][%d] = %v, want %v", i, j, result[i][j], a[i][j])
			}
		}
	}

	// Verify H is upper Hessenberg (zeros below first subdiagonal).
	for i := 2; i < n; i++ {
		for j := 0; j < i-1; j++ {
			if !approxEqual(h[i][j], 0, 1e-10) {
				t.Errorf("Hessenberg: H[%d][%d] = %v, want 0", i, j, h[i][j])
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Special Matrix Tests
// ---------------------------------------------------------------------------

func TestBlockDiag(t *testing.T) {
	a := [][]float64{{1, 2}, {3, 4}}
	b := [][]float64{{5}}
	result := BlockDiag(a, b)
	expected := [][]float64{
		{1, 2, 0},
		{3, 4, 0},
		{0, 0, 5},
	}
	for i := range expected {
		for j := range expected[i] {
			if result[i][j] != expected[i][j] {
				t.Errorf("BlockDiag[%d][%d] = %v, want %v", i, j, result[i][j], expected[i][j])
			}
		}
	}
}

func TestCompanion(t *testing.T) {
	// Polynomial: x^3 + 2x^2 + 3x + 4 => coeffs = [1, 2, 3, 4]
	c := Companion([]float64{1, 2, 3, 4})
	if len(c) != 3 || len(c[0]) != 3 {
		t.Fatalf("Companion: wrong size")
	}
	if c[0][0] != -2 || c[0][1] != -3 || c[0][2] != -4 {
		t.Errorf("Companion: first row = %v, want [-2, -3, -4]", c[0])
	}
	if c[1][0] != 1 || c[2][1] != 1 {
		t.Error("Companion: subdiagonal should be 1")
	}
}

func TestCirculant(t *testing.T) {
	c := Circulant([]float64{1, 2, 3})
	expected := [][]float64{
		{1, 2, 3},
		{3, 1, 2},
		{2, 3, 1},
	}
	for i := range expected {
		for j := range expected[i] {
			if c[i][j] != expected[i][j] {
				t.Errorf("Circulant[%d][%d] = %v, want %v", i, j, c[i][j], expected[i][j])
			}
		}
	}
}

func TestHadamard(t *testing.T) {
	h := Hadamard(4)
	if h == nil {
		t.Fatal("Hadamard(4): nil")
	}
	// Verify H * H^T = 4 * I.
	n := 4
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			sum := 0.0
			for k := 0; k < n; k++ {
				sum += h[i][k] * h[j][k]
			}
			expected := 0.0
			if i == j {
				expected = float64(n)
			}
			if !approxEqual(sum, expected, 1e-10) {
				t.Errorf("Hadamard: H*H^T[%d][%d] = %v, want %v", i, j, sum, expected)
			}
		}
	}
}

func TestHadamardInvalidSize(t *testing.T) {
	if Hadamard(3) != nil {
		t.Error("Hadamard(3): should return nil")
	}
	if Hadamard(0) != nil {
		t.Error("Hadamard(0): should return nil")
	}
}

func TestHilbert(t *testing.T) {
	h := Hilbert(3)
	expected := [][]float64{
		{1, 0.5, 1.0 / 3},
		{0.5, 1.0 / 3, 0.25},
		{1.0 / 3, 0.25, 0.2},
	}
	for i := range expected {
		for j := range expected[i] {
			if !approxEqual(h[i][j], expected[i][j], 1e-14) {
				t.Errorf("Hilbert[%d][%d] = %v, want %v", i, j, h[i][j], expected[i][j])
			}
		}
	}
}

func TestInvHilbert(t *testing.T) {
	n := 3
	h := Hilbert(n)
	hinv := InvHilbert(n)

	// Verify H * H^{-1} = I.
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			sum := 0.0
			for k := 0; k < n; k++ {
				sum += h[i][k] * hinv[k][j]
			}
			expected := 0.0
			if i == j {
				expected = 1.0
			}
			if !approxEqual(sum, expected, 1e-6) {
				t.Errorf("InvHilbert: H*Hinv[%d][%d] = %v, want %v", i, j, sum, expected)
			}
		}
	}
}

func TestPascal(t *testing.T) {
	p := Pascal(4)
	expected := [][]float64{
		{1, 1, 1, 1},
		{1, 2, 3, 4},
		{1, 3, 6, 10},
		{1, 4, 10, 20},
	}
	for i := range expected {
		for j := range expected[i] {
			if p[i][j] != expected[i][j] {
				t.Errorf("Pascal[%d][%d] = %v, want %v", i, j, p[i][j], expected[i][j])
			}
		}
	}
}

func TestToeplitz(t *testing.T) {
	c := []float64{1, 2, 3}
	r := []float64{1, 4, 5}
	tp := Toeplitz(c, r)
	expected := [][]float64{
		{1, 4, 5},
		{2, 1, 4},
		{3, 2, 1},
	}
	for i := range expected {
		for j := range expected[i] {
			if tp[i][j] != expected[i][j] {
				t.Errorf("Toeplitz[%d][%d] = %v, want %v", i, j, tp[i][j], expected[i][j])
			}
		}
	}
}

func TestHankel(t *testing.T) {
	c := []float64{1, 2, 3}
	r := []float64{3, 4, 5}
	h := Hankel(c, r)
	expected := [][]float64{
		{1, 2, 3},
		{2, 3, 4},
		{3, 4, 5},
	}
	for i := range expected {
		for j := range expected[i] {
			if h[i][j] != expected[i][j] {
				t.Errorf("Hankel[%d][%d] = %v, want %v", i, j, h[i][j], expected[i][j])
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Matrix Function Tests
// ---------------------------------------------------------------------------

func TestExpm(t *testing.T) {
	// e^(0) = I.
	zero := [][]float64{{0, 0}, {0, 0}}
	result, err := Expm(zero)
	if err != nil {
		t.Fatalf("Expm: %v", err)
	}
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			expected := 0.0
			if i == j {
				expected = 1.0
			}
			if !approxEqual(result[i][j], expected, 1e-10) {
				t.Errorf("Expm(0)[%d][%d] = %v, want %v", i, j, result[i][j], expected)
			}
		}
	}

	// e^(I) should have e on diagonal.
	eye := [][]float64{{1, 0}, {0, 1}}
	result, err = Expm(eye)
	if err != nil {
		t.Fatalf("Expm: %v", err)
	}
	if !approxEqual(result[0][0], math.E, 1e-8) {
		t.Errorf("Expm(I)[0][0] = %v, want %v", result[0][0], math.E)
	}
	if !approxEqual(result[1][1], math.E, 1e-8) {
		t.Errorf("Expm(I)[1][1] = %v, want %v", result[1][1], math.E)
	}
	if !approxEqual(result[0][1], 0, 1e-10) {
		t.Errorf("Expm(I)[0][1] = %v, want 0", result[0][1])
	}
}

func TestSqrtm(t *testing.T) {
	// sqrt(I) = I.
	eye := [][]float64{{1, 0}, {0, 1}}
	result, err := Sqrtm(eye)
	if err != nil {
		t.Fatalf("Sqrtm: %v", err)
	}
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			expected := 0.0
			if i == j {
				expected = 1.0
			}
			if !approxEqual(result[i][j], expected, 1e-10) {
				t.Errorf("Sqrtm(I)[%d][%d] = %v, want %v", i, j, result[i][j], expected)
			}
		}
	}

	// sqrt(4*I) = 2*I.
	a := [][]float64{{4, 0}, {0, 4}}
	result, err = Sqrtm(a)
	if err != nil {
		t.Fatalf("Sqrtm: %v", err)
	}
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			expected := 0.0
			if i == j {
				expected = 2.0
			}
			if !approxEqual(result[i][j], expected, 1e-8) {
				t.Errorf("Sqrtm(4I)[%d][%d] = %v, want %v", i, j, result[i][j], expected)
			}
		}
	}
}

func TestLogm(t *testing.T) {
	// log(I) = 0.
	eye := [][]float64{{1, 0}, {0, 1}}
	result, err := Logm(eye)
	if err != nil {
		t.Fatalf("Logm: %v", err)
	}
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			if !approxEqual(result[i][j], 0, 1e-8) {
				t.Errorf("Logm(I)[%d][%d] = %v, want 0", i, j, result[i][j])
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Stub Tests
// ---------------------------------------------------------------------------

func TestStubs(t *testing.T) {
	_, _, err := Polar(nil)
	if err == nil {
		t.Error("Polar: expected error")
	}
	_, err = Fiedler(nil)
	if err == nil {
		t.Error("Fiedler: expected error")
	}
	_, err = Leslie(nil, nil)
	if err == nil {
		t.Error("Leslie: expected error")
	}
	_, err = DFT(0)
	if err == nil {
		t.Error("DFT: expected error")
	}
	_, _, err = LDL(nil)
	if err == nil {
		t.Error("LDL: expected error")
	}
	_, _, err = Interpolative(nil, 0)
	if err == nil {
		t.Error("Interpolative: expected error")
	}
}
