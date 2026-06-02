package numgo

import (
	"fmt"
	"math"
)

// elementWise applies an operation element-wise to two arrays of the same shape.
func elementWise(a, b *NDArray, op func(float64, float64) float64) *NDArray {
	if !shapeEqual(a.shape, b.shape) {
		panic(fmt.Sprintf("numgo: shape mismatch %v vs %v", a.shape, b.shape))
	}
	data := make([]float64, a.Size())
	for i := range data {
		data[i] = op(a.data[i], b.data[i])
	}
	return NewNDArray(a.shape, data)
}

// scalarOp applies an operation between an array and a scalar.
func scalarOp(a *NDArray, s float64, op func(float64, float64) float64) *NDArray {
	data := make([]float64, a.Size())
	for i := range data {
		data[i] = op(a.data[i], s)
	}
	return NewNDArray(a.shape, data)
}

func shapeEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Add returns the element-wise sum of two arrays.
// Arrays with compatible shapes are broadcast to a common shape before the operation.
func Add(a, b *NDArray) *NDArray {
	return broadcastElementWise(a, b, func(x, y float64) float64 { return x + y })
}

// AddScalar adds a scalar to every element.
func AddScalar(a *NDArray, s float64) *NDArray {
	return scalarOp(a, s, func(x, y float64) float64 { return x + y })
}

// Sub returns the element-wise difference of two arrays.
// Arrays with compatible shapes are broadcast to a common shape before the operation.
func Sub(a, b *NDArray) *NDArray {
	return broadcastElementWise(a, b, func(x, y float64) float64 { return x - y })
}

// SubScalar subtracts a scalar from every element.
func SubScalar(a *NDArray, s float64) *NDArray {
	return scalarOp(a, s, func(x, y float64) float64 { return x - y })
}

// Mul returns the element-wise product of two arrays.
// Arrays with compatible shapes are broadcast to a common shape before the operation.
func Mul(a, b *NDArray) *NDArray {
	return broadcastElementWise(a, b, func(x, y float64) float64 { return x * y })
}

// MulScalar multiplies every element by a scalar.
func MulScalar(a *NDArray, s float64) *NDArray {
	return scalarOp(a, s, func(x, y float64) float64 { return x * y })
}

// Div returns the element-wise quotient of two arrays.
// Arrays with compatible shapes are broadcast to a common shape before the operation.
func Div(a, b *NDArray) *NDArray {
	return broadcastElementWise(a, b, func(x, y float64) float64 { return x / y })
}

// DivScalar divides every element by a scalar.
func DivScalar(a *NDArray, s float64) *NDArray {
	return scalarOp(a, s, func(x, y float64) float64 { return x / y })
}

// Sum reduces the array by summing along the given axes.
// If no axes are given, it sums over all elements and returns a scalar (1-D, length-1) array.
func Sum(a *NDArray, axes ...int) *NDArray {
	if len(axes) == 0 {
		total := 0.0
		for _, v := range a.data {
			total += v
		}
		return FromSlice([]float64{total})
	}
	return reduceAxis(a, axes, func(vals []float64) float64 {
		s := 0.0
		for _, v := range vals {
			s += v
		}
		return s
	})
}

// Prod reduces the array by multiplying along the given axes.
func Prod(a *NDArray, axes ...int) *NDArray {
	if len(axes) == 0 {
		total := 1.0
		for _, v := range a.data {
			total *= v
		}
		return FromSlice([]float64{total})
	}
	return reduceAxis(a, axes, func(vals []float64) float64 {
		p := 1.0
		for _, v := range vals {
			p *= v
		}
		return p
	})
}

// Max reduces the array by taking the maximum along the given axes.
func Max(a *NDArray, axes ...int) *NDArray {
	if len(axes) == 0 {
		m := math.Inf(-1)
		for _, v := range a.data {
			if v > m {
				m = v
			}
		}
		return FromSlice([]float64{m})
	}
	return reduceAxis(a, axes, func(vals []float64) float64 {
		m := math.Inf(-1)
		for _, v := range vals {
			if v > m {
				m = v
			}
		}
		return m
	})
}

// ArgMax returns the index of the maximum value along the given axis.
// Only a single axis is supported. Returns an NDArray of float64 indices.
func ArgMax(a *NDArray, axis int) *NDArray {
	return reduceAxis(a, []int{axis}, func(vals []float64) float64 {
		best := 0
		for i, v := range vals {
			if v > vals[best] {
				best = i
			}
		}
		return float64(best)
	})
}

// reduceAxis reduces along the specified axes using the given function.
// For simplicity, this implementation supports reducing one axis at a time
// and chains reductions for multiple axes.
func reduceAxis(a *NDArray, axes []int, fn func([]float64) float64) *NDArray {
	// Validate axes.
	for _, ax := range axes {
		if ax < 0 || ax >= a.Ndim() {
			panic(fmt.Sprintf("numgo: axis %d out of range for %d dimensions", ax, a.Ndim()))
		}
	}

	// Reduce one axis at a time (from highest to lowest to keep indices valid).
	sorted := make([]int, len(axes))
	copy(sorted, axes)
	// Simple insertion sort (axes list is tiny).
	for i := 1; i < len(sorted); i++ {
		for j := i; j > 0 && sorted[j] > sorted[j-1]; j-- {
			sorted[j], sorted[j-1] = sorted[j-1], sorted[j]
		}
	}

	result := a
	for _, ax := range sorted {
		result = reduceSingleAxis(result, ax, fn)
	}
	return result
}

// reduceSingleAxis reduces a single axis.
func reduceSingleAxis(a *NDArray, axis int, fn func([]float64) float64) *NDArray {
	// Build output shape (remove the axis).
	newShape := make([]int, 0, a.Ndim()-1)
	for i, s := range a.shape {
		if i != axis {
			newShape = append(newShape, s)
		}
	}
	if len(newShape) == 0 {
		newShape = []int{1}
	}

	outSize := product(newShape)
	outData := make([]float64, outSize)
	axisLen := a.shape[axis]
	indices := make([]int, a.Ndim())
	buf := make([]float64, axisLen)

	for outIdx := 0; outIdx < outSize; outIdx++ {
		// Decompose outIdx into output indices.
		outStrides := computeStrides(newShape)
		rem := outIdx
		oi := 0
		for d := 0; d < a.Ndim(); d++ {
			if d == axis {
				continue
			}
			indices[d] = rem / outStrides[oi]
			rem %= outStrides[oi]
			oi++
		}

		// Gather values along the axis.
		for k := 0; k < axisLen; k++ {
			indices[axis] = k
			buf[k] = a.data[a.flatIndex(indices)]
		}
		outData[outIdx] = fn(buf)
	}

	return NewNDArray(newShape, outData)
}
