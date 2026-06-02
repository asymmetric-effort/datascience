//go:build unit

package numgo

import (
	"math"
	"testing"
)

func approxEqual(a, b, tol float64) bool {
	return math.Abs(a-b) < tol
}

func matApproxEqual(a, b *NDArray, tol float64) bool {
	if a.Size() != b.Size() {
		return false
	}
	for i := 0; i < a.Size(); i++ {
		if !approxEqual(a.data[i], b.data[i], tol) {
			return false
		}
	}
	return true
}

// --- Dot ---

func TestDot1D(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3})
	b := FromSlice([]float64{4, 5, 6})
	r, err := Dot(a, b)
	if err != nil {
		t.Fatal(err)
	}
	if !approxEqual(r.data[0], 32, 1e-10) {
		t.Fatalf("expected 32, got %f", r.data[0])
	}
}

func TestDot2D(t *testing.T) {
	a := NewNDArray([]int{2, 2}, []float64{1, 2, 3, 4})
	b := NewNDArray([]int{2, 2}, []float64{5, 6, 7, 8})
	r, err := Dot(a, b)
	if err != nil {
		t.Fatal(err)
	}
	// [[19,22],[43,50]]
	expected := NewNDArray([]int{2, 2}, []float64{19, 22, 43, 50})
	if !matApproxEqual(r, expected, 1e-10) {
		t.Fatalf("unexpected result: %v", r.Data())
	}
}

func TestDot2D1D(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	b := FromSlice([]float64{1, 0, 1})
	r, err := Dot(a, b)
	if err != nil {
		t.Fatal(err)
	}
	// [4, 10]
	if !approxEqual(r.data[0], 4, 1e-10) || !approxEqual(r.data[1], 10, 1e-10) {
		t.Fatalf("unexpected result: %v", r.Data())
	}
}

// --- Matmul ---

func TestMatmul(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	b := NewNDArray([]int{3, 2}, []float64{7, 8, 9, 10, 11, 12})
	r, err := Matmul(a, b)
	if err != nil {
		t.Fatal(err)
	}
	expected := NewNDArray([]int{2, 2}, []float64{58, 64, 139, 154})
	if !matApproxEqual(r, expected, 1e-10) {
		t.Fatalf("unexpected: %v", r.Data())
	}
}

func TestMatmulDimensionMismatch(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	b := NewNDArray([]int{2, 2}, []float64{1, 2, 3, 4})
	_, err := Matmul(a, b)
	if err == nil {
		t.Fatal("expected error for dimension mismatch")
	}
}

// --- Inner and Outer ---

func TestInner(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3})
	b := FromSlice([]float64{4, 5, 6})
	r, err := Inner(a, b)
	if err != nil {
		t.Fatal(err)
	}
	if !approxEqual(r.data[0], 32, 1e-10) {
		t.Fatalf("expected 32, got %f", r.data[0])
	}
}

func TestOuter(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3})
	b := FromSlice([]float64{4, 5})
	r := Outer(a, b)
	expected := NewNDArray([]int{3, 2}, []float64{4, 5, 8, 10, 12, 15})
	if !matApproxEqual(r, expected, 1e-10) {
		t.Fatalf("unexpected: %v", r.Data())
	}
}

// --- Tensordot ---

func TestTensordot(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	b := NewNDArray([]int{3, 2}, []float64{7, 8, 9, 10, 11, 12})
	r, err := Tensordot(a, b, 1)
	if err != nil {
		t.Fatal(err)
	}
	// Same as matmul for 2D with axes=1.
	expected := NewNDArray([]int{2, 2}, []float64{58, 64, 139, 154})
	if !matApproxEqual(r, expected, 1e-10) {
		t.Fatalf("unexpected: %v", r.Data())
	}
}

// --- Solve ---

func TestSolve2x2(t *testing.T) {
	// 2x + y = 5, x + 3y = 7 => x=1.6, y=1.8
	a := NewNDArray([]int{2, 2}, []float64{2, 1, 1, 3})
	b := FromSlice([]float64{5, 7})
	x, err := Solve(a, b)
	if err != nil {
		t.Fatal(err)
	}
	if !approxEqual(x.data[0], 1.6, 1e-10) || !approxEqual(x.data[1], 1.8, 1e-10) {
		t.Fatalf("unexpected solution: %v", x.Data())
	}
}

func TestSolve3x3(t *testing.T) {
	// System: x=1, y=2, z=3
	a := NewNDArray([]int{3, 3}, []float64{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	})
	b := FromSlice([]float64{1, 2, 3})
	x, err := Solve(a, b)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		if !approxEqual(x.data[i], float64(i+1), 1e-10) {
			t.Fatalf("unexpected solution at %d: %f", i, x.data[i])
		}
	}
}

func TestSolve3x3NonTrivial(t *testing.T) {
	a := NewNDArray([]int{3, 3}, []float64{
		2, 1, -1,
		-3, -1, 2,
		-2, 1, 2,
	})
	b := FromSlice([]float64{8, -11, -3})
	x, err := Solve(a, b)
	if err != nil {
		t.Fatal(err)
	}
	// Verify Ax = b.
	for i := 0; i < 3; i++ {
		sum := 0.0
		for j := 0; j < 3; j++ {
			sum += a.Get(i, j) * x.data[j]
		}
		if !approxEqual(sum, b.data[i], 1e-10) {
			t.Fatalf("Ax != b at row %d: got %f, expected %f", i, sum, b.data[i])
		}
	}
}

// --- Det ---

func TestDetIdentity(t *testing.T) {
	for n := 1; n <= 4; n++ {
		d, err := Det(Eye(n))
		if err != nil {
			t.Fatal(err)
		}
		if !approxEqual(d, 1, 1e-10) {
			t.Fatalf("det(I_%d) = %f, expected 1", n, d)
		}
	}
}

func TestDet2x2(t *testing.T) {
	a := NewNDArray([]int{2, 2}, []float64{1, 2, 3, 4})
	d, err := Det(a)
	if err != nil {
		t.Fatal(err)
	}
	if !approxEqual(d, -2, 1e-10) {
		t.Fatalf("expected -2, got %f", d)
	}
}

func TestDetSingular(t *testing.T) {
	a := NewNDArray([]int{2, 2}, []float64{1, 2, 2, 4})
	d, err := Det(a)
	if err != nil {
		t.Fatal(err)
	}
	if !approxEqual(d, 0, 1e-10) {
		t.Fatalf("expected 0, got %f", d)
	}
}

// --- Inv ---

func TestInv2x2(t *testing.T) {
	a := NewNDArray([]int{2, 2}, []float64{4, 7, 2, 6})
	inv, err := Inv(a)
	if err != nil {
		t.Fatal(err)
	}
	// A * A^-1 should be identity.
	product, err := Matmul(a, inv)
	if err != nil {
		t.Fatal(err)
	}
	identity := Eye(2)
	if !matApproxEqual(product, identity, 1e-10) {
		t.Fatalf("A*inv(A) != I: %v", product.Data())
	}
}

func TestInv3x3(t *testing.T) {
	a := NewNDArray([]int{3, 3}, []float64{
		1, 2, 3,
		0, 1, 4,
		5, 6, 0,
	})
	inv, err := Inv(a)
	if err != nil {
		t.Fatal(err)
	}
	product, err := Matmul(a, inv)
	if err != nil {
		t.Fatal(err)
	}
	identity := Eye(3)
	if !matApproxEqual(product, identity, 1e-9) {
		t.Fatalf("A*inv(A) != I: %v", product.Data())
	}
}

func TestInvSingular(t *testing.T) {
	a := NewNDArray([]int{2, 2}, []float64{1, 2, 2, 4})
	_, err := Inv(a)
	if err == nil {
		t.Fatal("expected error for singular matrix")
	}
}

// --- Eig ---

func TestEigDiagonal(t *testing.T) {
	// Diagonal matrix: eigenvalues are diagonal elements.
	a := NewNDArray([]int{2, 2}, []float64{3, 0, 0, 5})
	vals, _, err := Eig(a)
	if err != nil {
		t.Fatal(err)
	}
	// Sort eigenvalues.
	v := vals.Data()
	if len(v) != 2 {
		t.Fatalf("expected 2 eigenvalues, got %d", len(v))
	}
	// They should be 3 and 5 (order may vary).
	if !(approxEqual(v[0], 3, 1e-6) && approxEqual(v[1], 5, 1e-6)) &&
		!(approxEqual(v[0], 5, 1e-6) && approxEqual(v[1], 3, 1e-6)) {
		t.Fatalf("unexpected eigenvalues: %v", v)
	}
}

func TestEigSymmetric(t *testing.T) {
	// Symmetric 2x2: eigenvalues of [[2,1],[1,2]] are 1 and 3.
	a := NewNDArray([]int{2, 2}, []float64{2, 1, 1, 2})
	vals, _, err := Eig(a)
	if err != nil {
		t.Fatal(err)
	}
	v := vals.Data()
	if !(approxEqual(v[0], 1, 1e-6) && approxEqual(v[1], 3, 1e-6)) &&
		!(approxEqual(v[0], 3, 1e-6) && approxEqual(v[1], 1, 1e-6)) {
		t.Fatalf("unexpected eigenvalues: %v", v)
	}
}

// --- Eigvals ---

func TestEigvals(t *testing.T) {
	a := NewNDArray([]int{2, 2}, []float64{3, 0, 0, 7})
	vals, err := Eigvals(a)
	if err != nil {
		t.Fatal(err)
	}
	v := vals.Data()
	if !(approxEqual(v[0], 3, 1e-6) && approxEqual(v[1], 7, 1e-6)) &&
		!(approxEqual(v[0], 7, 1e-6) && approxEqual(v[1], 3, 1e-6)) {
		t.Fatalf("unexpected eigenvalues: %v", v)
	}
}

// --- SVD ---

func TestSVDIdentity(t *testing.T) {
	a := Eye(3)
	u, s, vt, err := SVD(a)
	if err != nil {
		t.Fatal(err)
	}
	// Singular values of identity are all 1.
	for i := 0; i < 3; i++ {
		if !approxEqual(s.data[i], 1, 1e-6) {
			t.Fatalf("expected singular value 1, got %f", s.data[i])
		}
	}
	// U and Vt should be orthogonal.
	_ = u
	_ = vt
}

func TestSVDReconstruction(t *testing.T) {
	a := NewNDArray([]int{2, 2}, []float64{3, 0, 0, 5})
	u, s, vt, err := SVD(a)
	if err != nil {
		t.Fatal(err)
	}
	// Reconstruct: A = U * diag(S) * Vt
	sData := s.Data()
	// Build diagonal matrix.
	sDiag := Zeros(2, 2)
	for i := 0; i < 2; i++ {
		sDiag.Set(sData[i], i, i)
	}
	us, _ := Matmul(u, sDiag)
	recon, _ := Matmul(us, vt)
	if !matApproxEqual(recon, a, 1e-6) {
		t.Fatalf("SVD reconstruction failed: %v vs %v", recon.Data(), a.Data())
	}
}

// --- Cholesky ---

func TestCholesky(t *testing.T) {
	// SPD matrix: [[4,2],[2,3]]
	a := NewNDArray([]int{2, 2}, []float64{4, 2, 2, 3})
	l, err := Cholesky(a)
	if err != nil {
		t.Fatal(err)
	}
	// Verify L * L^T = A.
	lt := l.T()
	product, err := Matmul(l, lt)
	if err != nil {
		t.Fatal(err)
	}
	if !matApproxEqual(product, a, 1e-10) {
		t.Fatalf("L*L^T != A: %v", product.Data())
	}
}

func TestCholeskyNotPD(t *testing.T) {
	a := NewNDArray([]int{2, 2}, []float64{1, 2, 2, 1})
	_, err := Cholesky(a)
	if err == nil {
		t.Fatal("expected error for non-positive-definite matrix")
	}
}

// --- QR ---

func TestQR(t *testing.T) {
	a := NewNDArray([]int{3, 3}, []float64{
		12, -51, 4,
		6, 167, -68,
		-4, 24, -41,
	})
	q, r, err := QR(a)
	if err != nil {
		t.Fatal(err)
	}
	// Q*R should equal A.
	product, err := Matmul(q, r)
	if err != nil {
		t.Fatal(err)
	}
	if !matApproxEqual(product, a, 1e-8) {
		t.Fatalf("Q*R != A")
	}
	// Q should be orthogonal: Q^T * Q = I.
	qt := q.T()
	qtq, _ := Matmul(qt, q)
	identity := Eye(3)
	if !matApproxEqual(qtq, identity, 1e-8) {
		t.Fatalf("Q^T*Q != I: %v", qtq.Data())
	}
}

// --- Lstsq ---

func TestLstsq(t *testing.T) {
	// Overdetermined system: fit y = mx + c to (0,1), (1,2), (2,3).
	a := NewNDArray([]int{3, 2}, []float64{0, 1, 1, 1, 2, 1})
	b := FromSlice([]float64{1, 2, 3})
	x, err := Lstsq(a, b)
	if err != nil {
		t.Fatal(err)
	}
	// Expect m=1, c=1.
	if !approxEqual(x.data[0], 1, 1e-8) || !approxEqual(x.data[1], 1, 1e-8) {
		t.Fatalf("unexpected solution: %v", x.Data())
	}
}

// --- Norm ---

func TestNormL2(t *testing.T) {
	a := FromSlice([]float64{3, 4})
	r, err := Norm(a, 2, -1)
	if err != nil {
		t.Fatal(err)
	}
	if !approxEqual(r.data[0], 5, 1e-10) {
		t.Fatalf("expected 5, got %f", r.data[0])
	}
}

func TestNormL1(t *testing.T) {
	a := FromSlice([]float64{-3, 4})
	r, err := Norm(a, 1, -1)
	if err != nil {
		t.Fatal(err)
	}
	if !approxEqual(r.data[0], 7, 1e-10) {
		t.Fatalf("expected 7, got %f", r.data[0])
	}
}

// --- Cond ---

func TestCondIdentity(t *testing.T) {
	c, err := Cond(Eye(3))
	if err != nil {
		t.Fatal(err)
	}
	if !approxEqual(c, 1, 1e-6) {
		t.Fatalf("expected cond(I)=1, got %f", c)
	}
}

// --- MatrixRank ---

func TestMatrixRank(t *testing.T) {
	a := Eye(3)
	r, err := MatrixRank(a)
	if err != nil {
		t.Fatal(err)
	}
	if r != 3 {
		t.Fatalf("expected rank 3, got %d", r)
	}
}

func TestMatrixRankDeficient(t *testing.T) {
	// Rank 1 matrix.
	a := NewNDArray([]int{2, 2}, []float64{1, 2, 2, 4})
	r, err := MatrixRank(a)
	if err != nil {
		t.Fatal(err)
	}
	if r != 1 {
		t.Fatalf("expected rank 1, got %d", r)
	}
}

// --- MatrixPower ---

func TestMatrixPower(t *testing.T) {
	a := NewNDArray([]int{2, 2}, []float64{1, 1, 0, 1})
	r, err := MatrixPower(a, 3)
	if err != nil {
		t.Fatal(err)
	}
	// [[1,1],[0,1]]^3 = [[1,3],[0,1]]
	expected := NewNDArray([]int{2, 2}, []float64{1, 3, 0, 1})
	if !matApproxEqual(r, expected, 1e-10) {
		t.Fatalf("unexpected: %v", r.Data())
	}
}

func TestMatrixPowerZero(t *testing.T) {
	a := NewNDArray([]int{2, 2}, []float64{5, 6, 7, 8})
	r, err := MatrixPower(a, 0)
	if err != nil {
		t.Fatal(err)
	}
	if !matApproxEqual(r, Eye(2), 1e-10) {
		t.Fatalf("A^0 should be identity")
	}
}

func TestMatrixPowerNegative(t *testing.T) {
	a := NewNDArray([]int{2, 2}, []float64{4, 7, 2, 6})
	r, err := MatrixPower(a, -1)
	if err != nil {
		t.Fatal(err)
	}
	// A^-1 * A should be identity.
	product, _ := Matmul(r, a)
	if !matApproxEqual(product, Eye(2), 1e-9) {
		t.Fatalf("A^-1 * A != I")
	}
}

// --- Pinv ---

func TestPinv(t *testing.T) {
	a := NewNDArray([]int{2, 2}, []float64{1, 2, 3, 4})
	p, err := Pinv(a)
	if err != nil {
		t.Fatal(err)
	}
	// A * pinv(A) * A ≈ A
	ap, _ := Matmul(a, p)
	apa, _ := Matmul(ap, a)
	if !matApproxEqual(apa, a, 1e-6) {
		t.Fatalf("A*pinv(A)*A != A: %v", apa.Data())
	}
}

// --- Slogdet ---

func TestSlogdet(t *testing.T) {
	a := NewNDArray([]int{2, 2}, []float64{1, 2, 3, 4})
	sign, logdet, err := Slogdet(a)
	if err != nil {
		t.Fatal(err)
	}
	// det = -2, so sign=-1, logdet=ln(2)
	if !approxEqual(sign, -1, 1e-10) {
		t.Fatalf("expected sign=-1, got %f", sign)
	}
	if !approxEqual(logdet, math.Log(2), 1e-10) {
		t.Fatalf("expected logdet=%f, got %f", math.Log(2), logdet)
	}
}

// --- Trace ---

func TestTrace(t *testing.T) {
	a := NewNDArray([]int{3, 3}, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9})
	tr, err := Trace(a)
	if err != nil {
		t.Fatal(err)
	}
	if !approxEqual(tr, 15, 1e-10) {
		t.Fatalf("expected 15, got %f", tr)
	}
}

func TestTraceIdentity(t *testing.T) {
	tr, err := Trace(Eye(5))
	if err != nil {
		t.Fatal(err)
	}
	if !approxEqual(tr, 5, 1e-10) {
		t.Fatalf("expected 5, got %f", tr)
	}
}

// --- Cross ---

func TestCross(t *testing.T) {
	a := FromSlice([]float64{1, 0, 0})
	b := FromSlice([]float64{0, 1, 0})
	r, err := Cross(a, b)
	if err != nil {
		t.Fatal(err)
	}
	expected := FromSlice([]float64{0, 0, 1})
	if !matApproxEqual(r, expected, 1e-10) {
		t.Fatalf("unexpected: %v", r.Data())
	}
}

func TestCrossAnticommutative(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3})
	b := FromSlice([]float64{4, 5, 6})
	ab, _ := Cross(a, b)
	ba, _ := Cross(b, a)
	// a x b = -(b x a)
	for i := 0; i < 3; i++ {
		if !approxEqual(ab.data[i], -ba.data[i], 1e-10) {
			t.Fatalf("cross product not anti-commutative at %d", i)
		}
	}
}
