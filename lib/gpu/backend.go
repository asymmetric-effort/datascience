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

	// --- Tensor operations ---

	// ElementWiseAdd returns the element-wise sum of a and b.
	ElementWiseAdd(a, b []float64) []float64

	// ElementWiseSub returns the element-wise difference a - b.
	ElementWiseSub(a, b []float64) []float64

	// ElementWiseDiv returns the element-wise quotient a / b.
	ElementWiseDiv(a, b []float64) []float64

	// ScalarMul multiplies every element of a by the scalar s.
	ScalarMul(a []float64, s float64) []float64

	// ScalarAdd adds the scalar s to every element of a.
	ScalarAdd(a []float64, s float64) []float64

	// Exp returns the element-wise exponential of a.
	Exp(a []float64) []float64

	// Log returns the element-wise natural logarithm of a.
	Log(a []float64) []float64

	// Sqrt returns the element-wise square root of a.
	Sqrt(a []float64) []float64

	// Abs returns the element-wise absolute value of a.
	Abs(a []float64) []float64

	// Max returns the maximum element in a.
	Max(a []float64) float64

	// Min returns the minimum element in a.
	Min(a []float64) float64

	// ArgMax returns the index of the maximum element in a.
	ArgMax(a []float64) int

	// ArgMin returns the index of the minimum element in a.
	ArgMin(a []float64) int

	// Dot returns the dot product of a and b.
	Dot(a, b []float64) float64

	// --- Factor operations ---

	// FactorReduce fixes the variable at the given axis to the specified index,
	// returning the reduced values and the new shape with that axis removed.
	FactorReduce(values []float64, shape []int, axis int, index int) ([]float64, []int)

	// FactorMaximize takes the maximum over the given axis,
	// returning the reduced values and the new shape with that axis removed.
	FactorMaximize(values []float64, shape []int, axis int) ([]float64, []int)

	// LogSumExp computes log(sum(exp(a))) in a numerically stable way.
	LogSumExp(a []float64) float64

	// Softmax returns the softmax of a: exp(a_i) / sum(exp(a)).
	Softmax(a []float64) []float64

	// --- Batch operations ---

	// BatchMatMul performs batched matrix multiplication. a and b contain
	// batchSize matrices of dimensions (m x k) and (k x n) respectively,
	// packed contiguously. Returns batchSize result matrices of (m x n).
	BatchMatMul(a, b []float64, batchSize, m, k, n int) []float64

	// BatchNormalize normalizes batchSize vectors of length n packed
	// contiguously in a, so each sub-vector sums to 1.
	BatchNormalize(a []float64, batchSize, n int) []float64

	// --- Memory management ---

	// Alloc allocates a zeroed slice of the given size on the device.
	Alloc(size int) []float64

	// Free releases device memory. For CPU this is a no-op.
	Free(data []float64)

	// CopyToDevice copies data from host to device memory.
	CopyToDevice(data []float64) []float64

	// CopyFromDevice copies data from device to host memory.
	CopyFromDevice(data []float64) []float64

	// --- Device info ---

	// DeviceCount returns the number of available compute devices.
	DeviceCount() int

	// DeviceName returns the name of the device at the given index.
	DeviceName(index int) string

	// MemoryUsed returns the approximate device memory in use (bytes).
	MemoryUsed() int64

	// MemoryTotal returns the total device memory capacity (bytes).
	MemoryTotal() int64
}
