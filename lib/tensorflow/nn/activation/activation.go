// Package activation provides activation functions for neural network layers,
// analogous to tf.nn.relu, tf.nn.sigmoid, tf.nn.softmax, tf.nn.tanh.
package activation

import (
	"math"

	"github.com/asymmetric-effort/datascience/lib/numgo"
)

// ReLU applies the rectified linear unit function element-wise: max(0, x).
func ReLU(a *numgo.NDArray) *numgo.NDArray {
	data := a.Data()
	for i, v := range data {
		if v <= 0 {
			data[i] = 0
		}
	}
	return numgo.NewNDArray(a.Shape(), data)
}

// Sigmoid applies the sigmoid function element-wise: 1 / (1 + exp(-x)).
func Sigmoid(a *numgo.NDArray) *numgo.NDArray {
	data := a.Data()
	for i, v := range data {
		data[i] = 1.0 / (1.0 + math.Exp(-v))
	}
	return numgo.NewNDArray(a.Shape(), data)
}

// Tanh applies the hyperbolic tangent function element-wise.
func Tanh(a *numgo.NDArray) *numgo.NDArray {
	return numgo.Tanh(a)
}

// Softmax applies the softmax function along the last axis.
// For a 1D array, softmax(x)_i = exp(x_i) / sum(exp(x_j)).
// For a 2D array, softmax is applied independently to each row.
func Softmax(a *numgo.NDArray) *numgo.NDArray {
	data := a.Data()
	result := make([]float64, len(data))

	if a.Ndim() <= 1 {
		softmaxSlice(data, result)
	} else {
		shape := a.Shape()
		lastDim := shape[len(shape)-1]
		numSlices := len(data) / lastDim
		for s := range numSlices {
			offset := s * lastDim
			softmaxSlice(data[offset:offset+lastDim], result[offset:offset+lastDim])
		}
	}

	return numgo.NewNDArray(a.Shape(), result)
}

// LeakyReLU applies the leaky ReLU function: x if x > 0, else alpha * x.
func LeakyReLU(a *numgo.NDArray, alpha float64) *numgo.NDArray {
	data := a.Data()
	for i, v := range data {
		if v <= 0 {
			data[i] = alpha * v
		}
	}
	return numgo.NewNDArray(a.Shape(), data)
}

// ELU applies the exponential linear unit: x if x > 0, else alpha * (exp(x) - 1).
func ELU(a *numgo.NDArray, alpha float64) *numgo.NDArray {
	data := a.Data()
	for i, v := range data {
		if v <= 0 {
			data[i] = alpha * (math.Exp(v) - 1)
		}
	}
	return numgo.NewNDArray(a.Shape(), data)
}

// GELU applies the Gaussian Error Linear Unit approximation.
func GELU(a *numgo.NDArray) *numgo.NDArray {
	data := a.Data()
	for i, v := range data {
		data[i] = 0.5 * v * (1.0 + math.Tanh(math.Sqrt(2.0/math.Pi)*(v+0.044715*v*v*v)))
	}
	return numgo.NewNDArray(a.Shape(), data)
}

// Swish applies the swish activation: x * sigmoid(x).
func Swish(a *numgo.NDArray) *numgo.NDArray {
	data := a.Data()
	for i, v := range data {
		data[i] = v / (1.0 + math.Exp(-v))
	}
	return numgo.NewNDArray(a.Shape(), data)
}

// ReLUGrad computes the gradient of ReLU with respect to its input.
func ReLUGrad(input *numgo.NDArray, gradOutput *numgo.NDArray) *numgo.NDArray {
	inData := input.Data()
	gradData := gradOutput.Data()
	result := make([]float64, len(inData))
	for i := range inData {
		if inData[i] > 0 {
			result[i] = gradData[i]
		}
	}
	return numgo.NewNDArray(input.Shape(), result)
}

// SigmoidGrad computes the gradient of sigmoid. Given sigmoid output s, grad = s * (1 - s) * gradOutput.
func SigmoidGrad(sigmoidOutput *numgo.NDArray, gradOutput *numgo.NDArray) *numgo.NDArray {
	sData := sigmoidOutput.Data()
	gData := gradOutput.Data()
	result := make([]float64, len(sData))
	for i := range sData {
		result[i] = sData[i] * (1.0 - sData[i]) * gData[i]
	}
	return numgo.NewNDArray(sigmoidOutput.Shape(), result)
}

func softmaxSlice(input, output []float64) {
	maxVal := input[0]
	for _, v := range input[1:] {
		if v > maxVal {
			maxVal = v
		}
	}
	sumExp := 0.0
	for i, v := range input {
		output[i] = math.Exp(v - maxVal)
		sumExp += output[i]
	}
	for i := range output {
		output[i] /= sumExp
	}
}
