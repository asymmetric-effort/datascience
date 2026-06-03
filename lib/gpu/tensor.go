package gpu

// Tensor is a multi-dimensional array backed by a compute Backend.
// It stores data as a flat float64 slice in row-major order along with
// the shape describing its dimensionality.
type Tensor struct {
	data   []float64
	shape  []int
	device Backend
}

// NewTensor creates a Tensor with the given shape and data on the specified backend.
// If data is nil, the tensor is zero-initialized. If data is non-nil its length
// must equal the product of shape dimensions.
func NewTensor(shape []int, data []float64, backend Backend) *Tensor {
	size := tensorSize(shape)
	var d []float64
	if data == nil {
		d = backend.Alloc(size)
	} else {
		d = backend.CopyToDevice(data)
	}
	s := make([]int, len(shape))
	copy(s, shape)
	return &Tensor{data: d, shape: s, device: backend}
}

// Shape returns the dimensions of the tensor.
func (t *Tensor) Shape() []int {
	s := make([]int, len(t.shape))
	copy(s, t.shape)
	return s
}

// Data returns a copy of the tensor's underlying data on the host.
func (t *Tensor) Data() []float64 {
	return t.device.CopyFromDevice(t.data)
}

// Size returns the total number of elements in the tensor.
func (t *Tensor) Size() int {
	return tensorSize(t.shape)
}

// ToDevice copies this tensor to the target backend, returning a new Tensor.
func (t *Tensor) ToDevice(backend Backend) *Tensor {
	hostData := t.device.CopyFromDevice(t.data)
	return NewTensor(t.shape, hostData, backend)
}

// Add performs element-wise addition with another tensor.
// Both tensors must have the same shape.
func (t *Tensor) Add(other *Tensor) *Tensor {
	result := t.device.ElementWiseAdd(t.data, other.data)
	return &Tensor{data: result, shape: copyShape(t.shape), device: t.device}
}

// Sub performs element-wise subtraction (t - other).
func (t *Tensor) Sub(other *Tensor) *Tensor {
	result := t.device.ElementWiseSub(t.data, other.data)
	return &Tensor{data: result, shape: copyShape(t.shape), device: t.device}
}

// Mul performs element-wise multiplication with another tensor.
func (t *Tensor) Mul(other *Tensor) *Tensor {
	result := t.device.ElementWiseMul(t.data, other.data)
	return &Tensor{data: result, shape: copyShape(t.shape), device: t.device}
}

// MatMul performs matrix multiplication (t @ other).
// t must be 2-D (m x k), other must be 2-D (k x n).
func (t *Tensor) MatMul(other *Tensor) *Tensor {
	m := t.shape[0]
	k := t.shape[1]
	n := other.shape[1]
	result := t.device.MatMul(t.data, other.data, m, k, n)
	return &Tensor{data: result, shape: []int{m, n}, device: t.device}
}

// ScalarMul multiplies every element by a scalar.
func (t *Tensor) ScalarMul(s float64) *Tensor {
	result := t.device.ScalarMul(t.data, s)
	return &Tensor{data: result, shape: copyShape(t.shape), device: t.device}
}

// Sum reduces the tensor along the given axis by summation.
// axis must be in [0, ndim). Returns a tensor with that axis removed.
func (t *Tensor) Sum(axis int) *Tensor {
	result, newShape := t.device.Marginalize(t.data, t.shape, axis)
	return &Tensor{data: result, shape: newShape, device: t.device}
}

// Max reduces the tensor along the given axis by taking the maximum.
// axis must be in [0, ndim). Returns a tensor with that axis removed.
func (t *Tensor) Max(axis int) *Tensor {
	result, newShape := t.device.FactorMaximize(t.data, t.shape, axis)
	return &Tensor{data: result, shape: newShape, device: t.device}
}

// Reshape returns a new tensor with the same data but a different shape.
// The total number of elements must remain the same.
func (t *Tensor) Reshape(shape []int) *Tensor {
	return &Tensor{data: t.data, shape: copyShape(shape), device: t.device}
}

// Clone returns a deep copy of the tensor on the same device.
func (t *Tensor) Clone() *Tensor {
	d := t.device.CopyToDevice(t.device.CopyFromDevice(t.data))
	return &Tensor{data: d, shape: copyShape(t.shape), device: t.device}
}

// Normalize returns a new tensor whose elements sum to 1.
func (t *Tensor) Normalize() *Tensor {
	result := t.device.Normalize(t.data)
	return &Tensor{data: result, shape: copyShape(t.shape), device: t.device}
}

// Exp returns the element-wise exponential.
func (t *Tensor) Exp() *Tensor {
	result := t.device.Exp(t.data)
	return &Tensor{data: result, shape: copyShape(t.shape), device: t.device}
}

// Log returns the element-wise natural logarithm.
func (t *Tensor) Log() *Tensor {
	result := t.device.Log(t.data)
	return &Tensor{data: result, shape: copyShape(t.shape), device: t.device}
}

// tensorSize returns the product of shape dimensions.
func tensorSize(shape []int) int {
	n := 1
	for _, d := range shape {
		n *= d
	}
	return n
}

// copyShape returns a copy of a shape slice.
func copyShape(s []int) []int {
	c := make([]int, len(s))
	copy(c, s)
	return c
}
