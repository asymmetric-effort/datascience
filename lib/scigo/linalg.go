package scigo

import (
	"errors"
	"math"
)

// ---------------------------------------------------------------------------
// LU Decomposition
// ---------------------------------------------------------------------------

// LU computes the LU decomposition of a square matrix with partial pivoting.
// Returns the permutation matrix p, lower-triangular l, and upper-triangular u
// such that p*a = l*u.
func LU(a [][]float64) (p, l, u [][]float64, err error) {
	n := len(a)
	if n == 0 {
		return nil, nil, nil, errors.New("scigo.LU: empty matrix")
	}
	for _, row := range a {
		if len(row) != n {
			return nil, nil, nil, errors.New("scigo.LU: matrix must be square")
		}
	}

	lu, piv, err := LUFactor(a)
	if err != nil {
		return nil, nil, nil, err
	}

	// Build P, L, U from the compact form.
	p = make([][]float64, n)
	l = make([][]float64, n)
	u = make([][]float64, n)
	for i := 0; i < n; i++ {
		p[i] = make([]float64, n)
		l[i] = make([]float64, n)
		u[i] = make([]float64, n)
	}

	// Build permutation matrix from pivot indices.
	perm := make([]int, n)
	for i := range perm {
		perm[i] = i
	}
	for i, pi := range piv {
		perm[i], perm[pi] = perm[pi], perm[i]
	}
	for i := 0; i < n; i++ {
		p[i][perm[i]] = 1.0
	}

	// Extract L and U.
	for i := 0; i < n; i++ {
		l[i][i] = 1.0
		for j := 0; j < i; j++ {
			l[i][j] = lu[i][j]
		}
		for j := i; j < n; j++ {
			u[i][j] = lu[i][j]
		}
	}

	return p, l, u, nil
}

// LUFactor computes the compact LU factorization of a square matrix with partial pivoting.
// lu stores both L (unit lower) and U (upper) in a single matrix.
// piv[i] indicates the row swapped with row i during pivoting.
func LUFactor(a [][]float64) (lu [][]float64, piv []int, err error) {
	n := len(a)
	if n == 0 {
		return nil, nil, errors.New("scigo.LUFactor: empty matrix")
	}
	for _, row := range a {
		if len(row) != n {
			return nil, nil, errors.New("scigo.LUFactor: matrix must be square")
		}
	}

	// Copy the matrix.
	lu = make([][]float64, n)
	for i := 0; i < n; i++ {
		lu[i] = make([]float64, n)
		copy(lu[i], a[i])
	}
	piv = make([]int, n)
	for i := range piv {
		piv[i] = i
	}

	for col := 0; col < n; col++ {
		// Find pivot.
		maxVal := math.Abs(lu[col][col])
		maxRow := col
		for row := col + 1; row < n; row++ {
			if v := math.Abs(lu[row][col]); v > maxVal {
				maxVal = v
				maxRow = row
			}
		}
		if maxVal < 1e-14 {
			return nil, nil, errors.New("scigo.LUFactor: singular matrix")
		}
		piv[col] = maxRow
		if maxRow != col {
			lu[col], lu[maxRow] = lu[maxRow], lu[col]
		}

		for row := col + 1; row < n; row++ {
			lu[row][col] /= lu[col][col]
			for j := col + 1; j < n; j++ {
				lu[row][j] -= lu[row][col] * lu[col][j]
			}
		}
	}

	return lu, piv, nil
}

// LUSolve solves the system Ax = b given the compact LU factorization and pivot indices.
func LUSolve(lu [][]float64, piv []int, b []float64) ([]float64, error) {
	n := len(lu)
	if n == 0 {
		return nil, errors.New("scigo.LUSolve: empty matrix")
	}
	if len(b) != n {
		return nil, errors.New("scigo.LUSolve: dimension mismatch")
	}

	// Apply permutation to b.
	x := make([]float64, n)
	copy(x, b)
	for i, pi := range piv {
		x[i], x[pi] = x[pi], x[i]
	}

	// Forward substitution (L * y = Pb).
	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			x[i] -= lu[i][j] * x[j]
		}
	}

	// Back substitution (U * x = y).
	for i := n - 1; i >= 0; i-- {
		for j := i + 1; j < n; j++ {
			x[i] -= lu[i][j] * x[j]
		}
		x[i] /= lu[i][i]
	}

	return x, nil
}

// ---------------------------------------------------------------------------
// Cholesky
// ---------------------------------------------------------------------------

// ChoFactor computes the Cholesky factorization of a symmetric positive-definite matrix.
// Returns the lower-triangular matrix L such that A = L * L^T.
func ChoFactor(a [][]float64) ([][]float64, error) {
	n := len(a)
	if n == 0 {
		return nil, errors.New("scigo.ChoFactor: empty matrix")
	}
	for _, row := range a {
		if len(row) != n {
			return nil, errors.New("scigo.ChoFactor: matrix must be square")
		}
	}

	l := make([][]float64, n)
	for i := 0; i < n; i++ {
		l[i] = make([]float64, n)
	}

	for i := 0; i < n; i++ {
		for j := 0; j <= i; j++ {
			sum := 0.0
			for k := 0; k < j; k++ {
				sum += l[i][k] * l[j][k]
			}
			if i == j {
				val := a[i][i] - sum
				if val < 0 {
					return nil, errors.New("scigo.ChoFactor: matrix is not positive definite")
				}
				l[i][j] = math.Sqrt(val)
			} else {
				if math.Abs(l[j][j]) < 1e-14 {
					return nil, errors.New("scigo.ChoFactor: matrix is not positive definite")
				}
				l[i][j] = (a[i][j] - sum) / l[j][j]
			}
		}
	}

	return l, nil
}

// ChoSolve solves the system Ax = b given the Cholesky factor L (lower triangular, A = L*L^T).
func ChoSolve(cho [][]float64, b []float64) ([]float64, error) {
	n := len(cho)
	if n == 0 {
		return nil, errors.New("scigo.ChoSolve: empty matrix")
	}
	if len(b) != n {
		return nil, errors.New("scigo.ChoSolve: dimension mismatch")
	}

	// Forward substitution: L * y = b.
	y := make([]float64, n)
	for i := 0; i < n; i++ {
		sum := b[i]
		for j := 0; j < i; j++ {
			sum -= cho[i][j] * y[j]
		}
		y[i] = sum / cho[i][i]
	}

	// Back substitution: L^T * x = y.
	x := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		sum := y[i]
		for j := i + 1; j < n; j++ {
			sum -= cho[j][i] * x[j]
		}
		x[i] = sum / cho[i][i]
	}

	return x, nil
}

// ---------------------------------------------------------------------------
// Schur and Hessenberg
// ---------------------------------------------------------------------------

// Schur computes a simplified Schur decomposition using the QR algorithm.
// Returns T (quasi-upper-triangular) and Z (unitary) such that A = Z*T*Z^T.
// This is a simplified implementation that works well for symmetric matrices.
func Schur(a [][]float64) (t, z [][]float64, err error) {
	n := len(a)
	if n == 0 {
		return nil, nil, errors.New("scigo.Schur: empty matrix")
	}
	for _, row := range a {
		if len(row) != n {
			return nil, nil, errors.New("scigo.Schur: matrix must be square")
		}
	}

	// Copy matrix.
	t = matCopy(a)
	z = matEye(n)

	maxIter := 200 * n
	for iter := 0; iter < maxIter; iter++ {
		q, r := qrDecomp(t, n)
		t = matMul(r, q, n)
		z = matMul(z, q, n)

		// Check convergence: sub-diagonal elements near zero.
		offDiag := 0.0
		for i := 1; i < n; i++ {
			offDiag += math.Abs(t[i][i-1])
		}
		if offDiag < 1e-10 {
			break
		}
	}

	return t, z, nil
}

// Hessenberg reduces a square matrix to upper Hessenberg form using
// Householder reflections. Returns H and Q such that A = Q*H*Q^T.
func Hessenberg(a [][]float64) (h, q [][]float64, err error) {
	n := len(a)
	if n == 0 {
		return nil, nil, errors.New("scigo.Hessenberg: empty matrix")
	}
	for _, row := range a {
		if len(row) != n {
			return nil, nil, errors.New("scigo.Hessenberg: matrix must be square")
		}
	}

	h = matCopy(a)
	q = matEye(n)

	for k := 0; k < n-2; k++ {
		// Extract the column below the subdiagonal.
		m := n - k - 1
		x := make([]float64, m)
		for i := 0; i < m; i++ {
			x[i] = h[k+1+i][k]
		}

		// Compute Householder vector.
		alpha := vecNorm(x)
		if x[0] > 0 {
			alpha = -alpha
		}
		x[0] -= alpha
		xnorm := vecNorm(x)
		if xnorm < 1e-14 {
			continue
		}
		for i := range x {
			x[i] /= xnorm
		}

		// Apply H * P from the left: H[k+1:n, :] -= 2 * v * (v^T * H[k+1:n, :])
		for j := 0; j < n; j++ {
			dot := 0.0
			for i := 0; i < m; i++ {
				dot += x[i] * h[k+1+i][j]
			}
			for i := 0; i < m; i++ {
				h[k+1+i][j] -= 2 * x[i] * dot
			}
		}

		// Apply P * H from the right: H[:, k+1:n] -= 2 * (H[:, k+1:n] * v) * v^T
		for i := 0; i < n; i++ {
			dot := 0.0
			for j := 0; j < m; j++ {
				dot += h[i][k+1+j] * x[j]
			}
			for j := 0; j < m; j++ {
				h[i][k+1+j] -= 2 * dot * x[j]
			}
		}

		// Accumulate Q: Q[:, k+1:n] -= 2 * (Q[:, k+1:n] * v) * v^T
		for i := 0; i < n; i++ {
			dot := 0.0
			for j := 0; j < m; j++ {
				dot += q[i][k+1+j] * x[j]
			}
			for j := 0; j < m; j++ {
				q[i][k+1+j] -= 2 * dot * x[j]
			}
		}
	}

	return h, q, nil
}

// ---------------------------------------------------------------------------
// Special matrices
// ---------------------------------------------------------------------------

// BlockDiag creates a block diagonal matrix from the provided matrices.
func BlockDiag(matrices ...[][]float64) [][]float64 {
	totalRows := 0
	totalCols := 0
	for _, m := range matrices {
		totalRows += len(m)
		if len(m) > 0 {
			totalCols += len(m[0])
		}
	}

	result := make([][]float64, totalRows)
	for i := range result {
		result[i] = make([]float64, totalCols)
	}

	rowOff := 0
	colOff := 0
	for _, m := range matrices {
		rows := len(m)
		cols := 0
		if rows > 0 {
			cols = len(m[0])
		}
		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {
				result[rowOff+i][colOff+j] = m[i][j]
			}
		}
		rowOff += rows
		colOff += cols
	}

	return result
}

// Companion returns the companion matrix for a polynomial with the given
// coefficients. coeffs[0] is the leading coefficient.
func Companion(coeffs []float64) [][]float64 {
	n := len(coeffs) - 1
	if n < 1 {
		return nil
	}

	result := make([][]float64, n)
	for i := range result {
		result[i] = make([]float64, n)
	}

	lead := coeffs[0]
	for j := 0; j < n; j++ {
		result[0][j] = -coeffs[j+1] / lead
	}
	for i := 1; i < n; i++ {
		result[i][i-1] = 1.0
	}

	return result
}

// Circulant returns the circulant matrix formed by vector c.
func Circulant(c []float64) [][]float64 {
	n := len(c)
	result := make([][]float64, n)
	for i := range result {
		result[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			result[i][j] = c[(j-i+n)%n]
		}
	}
	return result
}

// Hadamard returns a Hadamard matrix of order n. n must be a power of 2.
func Hadamard(n int) [][]float64 {
	if n <= 0 || (n&(n-1)) != 0 {
		return nil
	}

	h := make([][]float64, n)
	for i := range h {
		h[i] = make([]float64, n)
	}
	h[0][0] = 1

	k := 1
	for k < n {
		for i := 0; i < k; i++ {
			for j := 0; j < k; j++ {
				h[i+k][j] = h[i][j]
				h[i][j+k] = h[i][j]
				h[i+k][j+k] = -h[i][j]
			}
		}
		k *= 2
	}

	return h
}

// Hankel returns the Hankel matrix with c as the first column and r as the last row.
// If r is nil, zeros are used for elements below the anti-diagonal.
func Hankel(c, r []float64) [][]float64 {
	n := len(c)
	if n == 0 {
		return nil
	}
	m := n
	if r != nil {
		m = len(r)
	}

	result := make([][]float64, n)
	for i := range result {
		result[i] = make([]float64, m)
		for j := 0; j < m; j++ {
			idx := i + j
			if idx < n {
				result[i][j] = c[idx]
			} else if r != nil && idx-n+1 < m {
				result[i][j] = r[idx-n+1]
			}
		}
	}

	return result
}

// Hilbert returns the n x n Hilbert matrix where H[i][j] = 1/(i+j+1).
func Hilbert(n int) [][]float64 {
	h := make([][]float64, n)
	for i := 0; i < n; i++ {
		h[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			h[i][j] = 1.0 / float64(i+j+1)
		}
	}
	return h
}

// InvHilbert returns the exact inverse of the n x n Hilbert matrix.
func InvHilbert(n int) [][]float64 {
	h := make([][]float64, n)
	for i := 0; i < n; i++ {
		h[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			sign := 1.0
			if (i+j)%2 != 0 {
				sign = -1.0
			}
			h[i][j] = sign * float64(i+j+1) *
				binomial(n+i, n-j-1) * binomial(n+j, n-i-1) *
				binomial(i+j, i) * binomial(i+j, i)
		}
	}
	return h
}

// Pascal returns the n x n Pascal matrix (symmetric).
func Pascal(n int) [][]float64 {
	p := make([][]float64, n)
	for i := 0; i < n; i++ {
		p[i] = make([]float64, n)
		p[i][0] = 1
		p[0][i] = 1
	}
	for i := 1; i < n; i++ {
		for j := 1; j < n; j++ {
			p[i][j] = p[i-1][j] + p[i][j-1]
		}
	}
	return p
}

// Toeplitz returns a Toeplitz matrix with c as the first column and r as the first row.
// If r is nil, c is used as the first row as well (symmetric Toeplitz).
func Toeplitz(c, r []float64) [][]float64 {
	n := len(c)
	if n == 0 {
		return nil
	}
	m := n
	if r != nil {
		m = len(r)
	}

	result := make([][]float64, n)
	for i := range result {
		result[i] = make([]float64, m)
		for j := 0; j < m; j++ {
			if i > j {
				result[i][j] = c[i-j]
			} else if r != nil {
				result[i][j] = r[j-i]
			} else {
				result[i][j] = c[j-i]
			}
		}
	}

	return result
}

// ---------------------------------------------------------------------------
// Matrix Functions
// ---------------------------------------------------------------------------

// Expm computes the matrix exponential using the scaling and squaring method
// with a truncated Taylor series.
func Expm(a [][]float64) ([][]float64, error) {
	n := len(a)
	if n == 0 {
		return nil, errors.New("scigo.Expm: empty matrix")
	}
	for _, row := range a {
		if len(row) != n {
			return nil, errors.New("scigo.Expm: matrix must be square")
		}
	}

	// Compute the 1-norm of a for scaling.
	norm := matNorm1(a)

	// Choose s such that ||A/2^s|| < 0.5.
	s := 0
	for norm > 0.5 {
		norm /= 2
		s++
	}

	// Scale the matrix.
	scale := math.Pow(2, -float64(s))
	as := matScale(a, scale)

	// Taylor series: e^A = I + A + A^2/2! + A^3/3! + ...
	result := matEye(n)
	term := matEye(n)
	for k := 1; k <= 20; k++ {
		term = matScale(matMul(term, as, n), 1.0/float64(k))
		result = matAdd(result, term)
		if matNorm1(term) < 1e-16 {
			break
		}
	}

	// Undo scaling by squaring s times.
	for i := 0; i < s; i++ {
		result = matMul(result, result, n)
	}

	return result, nil
}

// Logm computes the matrix logarithm using the inverse scaling and squaring method.
// Only works for matrices with positive real eigenvalues.
func Logm(a [][]float64) ([][]float64, error) {
	n := len(a)
	if n == 0 {
		return nil, errors.New("scigo.Logm: empty matrix")
	}
	for _, row := range a {
		if len(row) != n {
			return nil, errors.New("scigo.Logm: matrix must be square")
		}
	}

	// Repeated square root: A^{1/2^s} until close to I.
	eye := matEye(n)
	m := matCopy(a)
	s := 0
	maxS := 50

	for s < maxS {
		diff := matNorm1(matSub(m, eye))
		if diff < 0.5 {
			break
		}
		// m = sqrtm(m) — compute matrix square root.
		sq, err := sqrtmInternal(m)
		if err != nil {
			return nil, err
		}
		m = sq
		s++
	}

	// Now m is close to I. Use log(I + X) ≈ X - X^2/2 + X^3/3 - ... (Padé or series).
	x := matSub(m, eye)
	// Use a truncated series: sufficient for small X.
	result := matCopy(x)
	xk := matCopy(x)
	for k := 2; k <= 16; k++ {
		xk = matMul(xk, x, n)
		sign := -1.0
		if k%2 == 0 {
			sign = -1.0
		} else {
			sign = 1.0
		}
		result = matAdd(result, matScale(xk, sign/float64(k)))
	}

	// Undo scaling: log(A) = 2^s * log(A^{1/2^s}).
	result = matScale(result, math.Pow(2, float64(s)))

	return result, nil
}

// Sqrtm computes the matrix square root via the Denman-Beavers iteration.
func Sqrtm(a [][]float64) ([][]float64, error) {
	return sqrtmInternal(a)
}

func sqrtmInternal(a [][]float64) ([][]float64, error) {
	n := len(a)
	if n == 0 {
		return nil, errors.New("scigo.Sqrtm: empty matrix")
	}
	for _, row := range a {
		if len(row) != n {
			return nil, errors.New("scigo.Sqrtm: matrix must be square")
		}
	}

	// Denman-Beavers iteration:
	// Y_{k+1} = 0.5 * (Y_k + Z_k^{-1})
	// Z_{k+1} = 0.5 * (Z_k + Y_k^{-1})
	// Converges to Y = A^{1/2}, Z = A^{-1/2}.
	y := matCopy(a)
	z := matEye(n)

	for iter := 0; iter < 100; iter++ {
		zInv, err := matInverse(z)
		if err != nil {
			return nil, err
		}
		yInv, err := matInverse(y)
		if err != nil {
			return nil, err
		}

		yNew := matScale(matAdd(y, zInv), 0.5)
		zNew := matScale(matAdd(z, yInv), 0.5)

		diff := matNorm1(matSub(yNew, y))
		y = yNew
		z = zNew
		if diff < 1e-12 {
			break
		}
	}

	return y, nil
}

// ---------------------------------------------------------------------------
// Stubs for remaining functions
// ---------------------------------------------------------------------------

// Polar returns an error indicating it is not yet implemented.
func Polar(_ [][]float64) ([][]float64, [][]float64, error) {
	return nil, nil, errors.New("scigo.Polar: not yet implemented")
}

// Fiedler returns an error indicating it is not yet implemented.
func Fiedler(_ []float64) ([][]float64, error) {
	return nil, errors.New("scigo.Fiedler: not yet implemented")
}

// Leslie returns an error indicating it is not yet implemented.
func Leslie(_, _ []float64) ([][]float64, error) {
	return nil, errors.New("scigo.Leslie: not yet implemented")
}

// DFT returns an error indicating it is not yet implemented.
func DFT(_ int) ([][]complex128, error) {
	return nil, errors.New("scigo.DFT: not yet implemented")
}

// LDL returns an error indicating it is not yet implemented.
func LDL(_ [][]float64) ([][]float64, []float64, error) {
	return nil, nil, errors.New("scigo.LDL: not yet implemented")
}

// Interpolative returns an error indicating it is not yet implemented.
func Interpolative(_ [][]float64, _ int) ([]int, [][]float64, error) {
	return nil, nil, errors.New("scigo.Interpolative: not yet implemented")
}

// ---------------------------------------------------------------------------
// Internal matrix helpers ([][]float64 based)
// ---------------------------------------------------------------------------

func matCopy(a [][]float64) [][]float64 {
	n := len(a)
	c := make([][]float64, n)
	for i := range a {
		c[i] = make([]float64, len(a[i]))
		copy(c[i], a[i])
	}
	return c
}

func matEye(n int) [][]float64 {
	m := make([][]float64, n)
	for i := range m {
		m[i] = make([]float64, n)
		m[i][i] = 1.0
	}
	return m
}

func matMul(a, b [][]float64, n int) [][]float64 {
	c := make([][]float64, n)
	for i := range c {
		c[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			s := 0.0
			for k := 0; k < n; k++ {
				s += a[i][k] * b[k][j]
			}
			c[i][j] = s
		}
	}
	return c
}

func matAdd(a, b [][]float64) [][]float64 {
	n := len(a)
	c := make([][]float64, n)
	for i := range c {
		c[i] = make([]float64, n)
		for j := range c[i] {
			c[i][j] = a[i][j] + b[i][j]
		}
	}
	return c
}

func matSub(a, b [][]float64) [][]float64 {
	n := len(a)
	c := make([][]float64, n)
	for i := range c {
		c[i] = make([]float64, n)
		for j := range c[i] {
			c[i][j] = a[i][j] - b[i][j]
		}
	}
	return c
}

func matScale(a [][]float64, s float64) [][]float64 {
	n := len(a)
	c := make([][]float64, n)
	for i := range c {
		c[i] = make([]float64, n)
		for j := range c[i] {
			c[i][j] = a[i][j] * s
		}
	}
	return c
}

func matNorm1(a [][]float64) float64 {
	n := len(a)
	if n == 0 {
		return 0
	}
	maxCol := 0.0
	for j := 0; j < n; j++ {
		s := 0.0
		for i := 0; i < n; i++ {
			s += math.Abs(a[i][j])
		}
		if s > maxCol {
			maxCol = s
		}
	}
	return maxCol
}

func matSolve(a, b [][]float64) ([][]float64, error) {
	n := len(a)
	// Solve A * X = B column by column.
	lu, piv, err := LUFactor(a)
	if err != nil {
		return nil, err
	}
	result := make([][]float64, n)
	for i := range result {
		result[i] = make([]float64, n)
	}
	for j := 0; j < n; j++ {
		col := make([]float64, n)
		for i := 0; i < n; i++ {
			col[i] = b[i][j]
		}
		x, err := LUSolve(lu, piv, col)
		if err != nil {
			return nil, err
		}
		for i := 0; i < n; i++ {
			result[i][j] = x[i]
		}
	}
	return result, nil
}

func matInverse(a [][]float64) ([][]float64, error) {
	return matSolve(a, matEye(len(a)))
}

func vecNorm(x []float64) float64 {
	s := 0.0
	for _, v := range x {
		s += v * v
	}
	return math.Sqrt(s)
}

func qrDecomp(a [][]float64, n int) (q, r [][]float64) {
	q = matEye(n)
	r = matCopy(a)

	for j := 0; j < n-1; j++ {
		// Compute Householder reflection for column j.
		x := make([]float64, n-j)
		for i := j; i < n; i++ {
			x[i-j] = r[i][j]
		}
		alpha := vecNorm(x)
		if x[0] > 0 {
			alpha = -alpha
		}
		x[0] -= alpha
		xn := vecNorm(x)
		if xn < 1e-14 {
			continue
		}
		for i := range x {
			x[i] /= xn
		}

		// Apply to R from left.
		for col := j; col < n; col++ {
			dot := 0.0
			for i := 0; i < len(x); i++ {
				dot += x[i] * r[j+i][col]
			}
			for i := 0; i < len(x); i++ {
				r[j+i][col] -= 2 * x[i] * dot
			}
		}

		// Accumulate Q.
		for row := 0; row < n; row++ {
			dot := 0.0
			for i := 0; i < len(x); i++ {
				dot += q[row][j+i] * x[i]
			}
			for i := 0; i < len(x); i++ {
				q[row][j+i] -= 2 * dot * x[i]
			}
		}
	}

	return q, r
}

// binomial computes "n choose k" as a float64.
func binomial(n, k int) float64 {
	if k < 0 || k > n {
		return 0
	}
	if k == 0 || k == n {
		return 1
	}
	if k > n-k {
		k = n - k
	}
	result := 1.0
	for i := 0; i < k; i++ {
		result *= float64(n - i)
		result /= float64(i + 1)
	}
	return result
}
