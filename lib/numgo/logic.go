package numgo

import (
	"math"
)

// All returns 1.0 if all elements along the given axes are nonzero, 0.0 otherwise.
// If no axes are given, checks all elements.
func All(a *NDArray, axes ...int) *NDArray {
	if len(axes) == 0 {
		for _, v := range a.data {
			if v == 0 {
				return FromSlice([]float64{0})
			}
		}
		return FromSlice([]float64{1})
	}
	return reduceAxis(a, axes, func(vals []float64) float64 {
		for _, v := range vals {
			if v == 0 {
				return 0
			}
		}
		return 1
	})
}

// Any returns 1.0 if any element along the given axes is nonzero, 0.0 otherwise.
// If no axes are given, checks all elements.
func Any(a *NDArray, axes ...int) *NDArray {
	if len(axes) == 0 {
		for _, v := range a.data {
			if v != 0 {
				return FromSlice([]float64{1})
			}
		}
		return FromSlice([]float64{0})
	}
	return reduceAxis(a, axes, func(vals []float64) float64 {
		for _, v := range vals {
			if v != 0 {
				return 1
			}
		}
		return 0
	})
}

// Isnan returns an NDArray with 1.0 where the element is NaN, 0.0 otherwise.
func Isnan(a *NDArray) *NDArray {
	return unaryOp(a, func(x float64) float64 {
		if math.IsNaN(x) {
			return 1
		}
		return 0
	})
}

// Isinf returns an NDArray with 1.0 where the element is +/-Inf, 0.0 otherwise.
func Isinf(a *NDArray) *NDArray {
	return unaryOp(a, func(x float64) float64 {
		if math.IsInf(x, 0) {
			return 1
		}
		return 0
	})
}

// Isfinite returns an NDArray with 1.0 where the element is finite, 0.0 otherwise.
func Isfinite(a *NDArray) *NDArray {
	return unaryOp(a, func(x float64) float64 {
		if !math.IsInf(x, 0) && !math.IsNaN(x) {
			return 1
		}
		return 0
	})
}

// Isneginf returns an NDArray with 1.0 where the element is -Inf, 0.0 otherwise.
func Isneginf(a *NDArray) *NDArray {
	return unaryOp(a, func(x float64) float64 {
		if math.IsInf(x, -1) {
			return 1
		}
		return 0
	})
}

// Isposinf returns an NDArray with 1.0 where the element is +Inf, 0.0 otherwise.
func Isposinf(a *NDArray) *NDArray {
	return unaryOp(a, func(x float64) float64 {
		if math.IsInf(x, 1) {
			return 1
		}
		return 0
	})
}

// LogicalAnd returns element-wise logical AND. Nonzero values are treated as true.
func LogicalAnd(a, b *NDArray) *NDArray {
	return broadcastElementWise(a, b, func(x, y float64) float64 {
		if x != 0 && y != 0 {
			return 1
		}
		return 0
	})
}

// LogicalOr returns element-wise logical OR. Nonzero values are treated as true.
func LogicalOr(a, b *NDArray) *NDArray {
	return broadcastElementWise(a, b, func(x, y float64) float64 {
		if x != 0 || y != 0 {
			return 1
		}
		return 0
	})
}

// LogicalNot returns element-wise logical NOT. Nonzero values become 0.0, zero becomes 1.0.
func LogicalNot(a *NDArray) *NDArray {
	return unaryOp(a, func(x float64) float64 {
		if x == 0 {
			return 1
		}
		return 0
	})
}

// LogicalXor returns element-wise logical XOR. Nonzero values are treated as true.
func LogicalXor(a, b *NDArray) *NDArray {
	return broadcastElementWise(a, b, func(x, y float64) float64 {
		xb := x != 0
		yb := y != 0
		if xb != yb {
			return 1
		}
		return 0
	})
}

// ArrayEqual returns true if a and b have the same shape and all elements are equal.
func ArrayEqual(a, b *NDArray) bool {
	if !shapeEqual(a.shape, b.shape) {
		return false
	}
	for i := range a.data {
		if a.data[i] != b.data[i] {
			return false
		}
	}
	return true
}

// ArrayEquiv returns true if a and b are equal after broadcasting to a common shape.
func ArrayEquiv(a, b *NDArray) bool {
	resultShape, err := BroadcastShapes(a.shape, b.shape)
	if err != nil {
		return false
	}
	ab, err := BroadcastTo(a, resultShape)
	if err != nil {
		return false
	}
	bb, err := BroadcastTo(b, resultShape)
	if err != nil {
		return false
	}
	for i := range ab.data {
		if ab.data[i] != bb.data[i] {
			return false
		}
	}
	return true
}

// Greater returns 1.0 where a > b, 0.0 otherwise (element-wise with broadcasting).
func Greater(a, b *NDArray) *NDArray {
	return broadcastElementWise(a, b, func(x, y float64) float64 {
		if x > y {
			return 1
		}
		return 0
	})
}

// Less returns 1.0 where a < b, 0.0 otherwise (element-wise with broadcasting).
func Less(a, b *NDArray) *NDArray {
	return broadcastElementWise(a, b, func(x, y float64) float64 {
		if x < y {
			return 1
		}
		return 0
	})
}

// Equal returns 1.0 where a == b, 0.0 otherwise (element-wise with broadcasting).
func Equal(a, b *NDArray) *NDArray {
	return broadcastElementWise(a, b, func(x, y float64) float64 {
		if x == y {
			return 1
		}
		return 0
	})
}

// NotEqual returns 1.0 where a != b, 0.0 otherwise (element-wise with broadcasting).
func NotEqual(a, b *NDArray) *NDArray {
	return broadcastElementWise(a, b, func(x, y float64) float64 {
		if x != y {
			return 1
		}
		return 0
	})
}

// GreaterEqual returns 1.0 where a >= b, 0.0 otherwise (element-wise with broadcasting).
func GreaterEqual(a, b *NDArray) *NDArray {
	return broadcastElementWise(a, b, func(x, y float64) float64 {
		if x >= y {
			return 1
		}
		return 0
	})
}

// LessEqual returns 1.0 where a <= b, 0.0 otherwise (element-wise with broadcasting).
func LessEqual(a, b *NDArray) *NDArray {
	return broadcastElementWise(a, b, func(x, y float64) float64 {
		if x <= y {
			return 1
		}
		return 0
	})
}

// Isclose returns 1.0 where |a-b| <= atol + rtol*|b|, 0.0 otherwise (element-wise with broadcasting).
func Isclose(a, b *NDArray, atol, rtol float64) *NDArray {
	return broadcastElementWise(a, b, func(x, y float64) float64 {
		if math.Abs(x-y) <= atol+rtol*math.Abs(y) {
			return 1
		}
		return 0
	})
}
