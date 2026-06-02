package numgo

import "math"

// unaryOp applies fn element-wise to a, returning a new NDArray of the same shape.
func unaryOp(a *NDArray, fn func(float64) float64) *NDArray {
	data := make([]float64, a.Size())
	for i, v := range a.data {
		data[i] = fn(v)
	}
	return NewNDArray(a.shape, data)
}

// Sin returns the element-wise sine of the array.
func Sin(a *NDArray) *NDArray {
	return unaryOp(a, math.Sin)
}

// Cos returns the element-wise cosine of the array.
func Cos(a *NDArray) *NDArray {
	return unaryOp(a, math.Cos)
}

// Tan returns the element-wise tangent of the array.
func Tan(a *NDArray) *NDArray {
	return unaryOp(a, math.Tan)
}

// Arcsin returns the element-wise arcsine of the array.
func Arcsin(a *NDArray) *NDArray {
	return unaryOp(a, math.Asin)
}

// Arccos returns the element-wise arccosine of the array.
func Arccos(a *NDArray) *NDArray {
	return unaryOp(a, math.Acos)
}

// Arctan returns the element-wise arctangent of the array.
func Arctan(a *NDArray) *NDArray {
	return unaryOp(a, math.Atan)
}

// Arctan2 returns the element-wise two-argument arctangent of y/x with broadcasting.
func Arctan2(y, x *NDArray) *NDArray {
	return broadcastElementWise(y, x, math.Atan2)
}

// Sinh returns the element-wise hyperbolic sine of the array.
func Sinh(a *NDArray) *NDArray {
	return unaryOp(a, math.Sinh)
}

// Cosh returns the element-wise hyperbolic cosine of the array.
func Cosh(a *NDArray) *NDArray {
	return unaryOp(a, math.Cosh)
}

// Tanh returns the element-wise hyperbolic tangent of the array.
func Tanh(a *NDArray) *NDArray {
	return unaryOp(a, math.Tanh)
}

// Exp returns the element-wise exponential (e^x) of the array.
func Exp(a *NDArray) *NDArray {
	return unaryOp(a, math.Exp)
}

// Exp2 returns the element-wise 2^x of the array.
func Exp2(a *NDArray) *NDArray {
	return unaryOp(a, math.Exp2)
}

// Expm1 returns the element-wise exp(x)-1 of the array.
func Expm1(a *NDArray) *NDArray {
	return unaryOp(a, math.Expm1)
}

// Log returns the element-wise natural logarithm of the array.
func Log(a *NDArray) *NDArray {
	return unaryOp(a, math.Log)
}

// Log2 returns the element-wise base-2 logarithm of the array.
func Log2(a *NDArray) *NDArray {
	return unaryOp(a, math.Log2)
}

// Log10 returns the element-wise base-10 logarithm of the array.
func Log10(a *NDArray) *NDArray {
	return unaryOp(a, math.Log10)
}

// Log1p returns the element-wise log(1+x) of the array.
func Log1p(a *NDArray) *NDArray {
	return unaryOp(a, math.Log1p)
}

// Power returns the element-wise base**exp with broadcasting.
func Power(base, exp *NDArray) *NDArray {
	return broadcastElementWise(base, exp, math.Pow)
}

// Sqrt returns the element-wise square root of the array.
func Sqrt(a *NDArray) *NDArray {
	return unaryOp(a, math.Sqrt)
}

// Cbrt returns the element-wise cube root of the array.
func Cbrt(a *NDArray) *NDArray {
	return unaryOp(a, math.Cbrt)
}

// Square returns the element-wise x*x of the array.
func Square(a *NDArray) *NDArray {
	return unaryOp(a, func(x float64) float64 { return x * x })
}

// Absolute returns the element-wise absolute value of the array.
func Absolute(a *NDArray) *NDArray {
	return unaryOp(a, math.Abs)
}

// Fabs returns the element-wise absolute value of the array (same as Absolute for float64).
func Fabs(a *NDArray) *NDArray {
	return unaryOp(a, math.Abs)
}

// Sign returns the element-wise sign of the array: -1 for negative, 0 for zero, 1 for positive.
func Sign(a *NDArray) *NDArray {
	return unaryOp(a, func(x float64) float64 {
		switch {
		case x < 0:
			return -1
		case x > 0:
			return 1
		default:
			return 0
		}
	})
}

// Heaviside computes the Heaviside step function element-wise with broadcasting.
// Returns 0 where x < 0, h0 where x == 0, and 1 where x > 0.
func Heaviside(x, h0 *NDArray) *NDArray {
	return broadcastElementWise(x, h0, func(xv, h0v float64) float64 {
		switch {
		case xv < 0:
			return 0
		case xv == 0:
			return h0v
		default:
			return 1
		}
	})
}

// Fmod returns the element-wise floating-point remainder (math.Mod) with broadcasting.
func Fmod(x, y *NDArray) *NDArray {
	return broadcastElementWise(x, y, math.Mod)
}

// Modf returns the fractional and integer parts of each element.
// Both returned arrays have the same shape as the input.
func Modf(x *NDArray) (frac, integer *NDArray) {
	fracData := make([]float64, x.Size())
	intData := make([]float64, x.Size())
	for i, v := range x.data {
		intPart, fracPart := math.Modf(v)
		fracData[i] = fracPart
		intData[i] = intPart
	}
	return NewNDArray(x.shape, fracData), NewNDArray(x.shape, intData)
}

// Remainder returns the element-wise IEEE 754 remainder (math.Remainder) with broadcasting.
func Remainder(x, y *NDArray) *NDArray {
	return broadcastElementWise(x, y, math.Remainder)
}

// Clip clamps every element of a to the range [min, max].
func Clip(a *NDArray, min, max float64) *NDArray {
	return unaryOp(a, func(x float64) float64 {
		if x < min {
			return min
		}
		if x > max {
			return max
		}
		return x
	})
}

// Around rounds every element to the given number of decimal places.
func Around(a *NDArray, decimals int) *NDArray {
	shift := math.Pow(10, float64(decimals))
	return unaryOp(a, func(x float64) float64 {
		return math.RoundToEven(x*shift) / shift
	})
}

// Rint rounds every element to the nearest integer (using banker's rounding).
func Rint(a *NDArray) *NDArray {
	return unaryOp(a, math.RoundToEven)
}

// Floor returns the element-wise floor of the array.
func Floor(a *NDArray) *NDArray {
	return unaryOp(a, math.Floor)
}

// Ceil returns the element-wise ceiling of the array.
func Ceil(a *NDArray) *NDArray {
	return unaryOp(a, math.Ceil)
}
