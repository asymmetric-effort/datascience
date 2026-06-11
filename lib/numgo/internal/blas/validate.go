package blas

import "fmt"

// requiredLen returns the minimum slice length needed for n elements
// with the given increment. For inc=1, this is simply n.
// For general inc, the last accessed index is (n-1)*inc, so length is (n-1)*inc + 1.
func requiredLen(n, inc int) int {
	if inc == 1 {
		return n
	}
	return (n-1)*inc + 1
}

// validateVector panics if x is too short for n elements at the given increment.
func validateVector(name string, x []float64, n, inc int) {
	if need := requiredLen(n, inc); len(x) < need {
		panic(fmt.Sprintf("blas: %s: slice length %d < required %d (n=%d, inc=%d)", name, len(x), need, n, inc))
	}
}

// validateMatrix panics if a is too short for an m×n matrix with leading dimension lda.
func validateMatrix(name string, a []float64, m, n, lda int) {
	if m <= 0 || n <= 0 {
		return
	}
	need := (m-1)*lda + n
	if len(a) < need {
		panic(fmt.Sprintf("blas: %s: slice length %d < required %d (m=%d, n=%d, lda=%d)", name, len(a), need, m, n, lda))
	}
}
