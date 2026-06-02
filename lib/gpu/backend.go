package gpu

// Backend defines the compute backend interface for accelerated factor operations.
// Implementations may use CPU, CUDA, OpenCL, or other compute backends.
type Backend interface {
	// Name returns the backend identifier (e.g. "cpu", "cuda", "opencl").
	Name() string

	// IsAvailable reports whether this backend can run on the current system.
	IsAvailable() bool

	// MatMul performs matrix multiplication of a (m x k) and b (k x n),
	// returning the result as a flat slice of length m*n in row-major order.
	MatMul(a, b []float64, m, k, n int) []float64

	// ElementWiseMul returns the element-wise product of a and b.
	// The slices must have equal length.
	ElementWiseMul(a, b []float64) []float64

	// Sum returns the sum of all elements in a.
	Sum(a []float64) float64

	// Normalize divides each element by the sum of all elements.
	// Returns a new slice; does not modify the input.
	Normalize(a []float64) []float64

	// FactorProduct computes the product of two discrete factors.
	// aValues/bValues are the flat value arrays; aShape/bShape describe their
	// dimensionality; resultShape is the shape of the output factor.
	FactorProduct(aValues []float64, aShape []int, bValues []float64, bShape []int, resultShape []int) []float64

	// Marginalize sums out the given axis from a multi-dimensional array,
	// returning the reduced values and the new shape with that axis removed.
	Marginalize(values []float64, shape []int, axis int) ([]float64, []int)

	// Close releases any resources held by the backend.
	Close() error
}
