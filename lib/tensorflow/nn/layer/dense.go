// Package layer provides neural network layers, analogous to tf.keras.layers.
package layer

import (
	"fmt"
	"math"

	"github.com/asymmetric-effort/datascience/lib/numgo"
)

// Dense implements a fully connected (dense) layer: output = input @ weights + bias.
// Analogous to tf.keras.layers.Dense.
type Dense struct {
	weights *numgo.NDArray
	bias    *numgo.NDArray
	inSize  int
	outSize int
}

// NewDense creates a new dense layer with the given input and output sizes.
// Weights are initialized using Xavier/Glorot uniform initialization.
func NewDense(inSize, outSize int, rng func() float64) (*Dense, error) {
	if inSize <= 0 || outSize <= 0 {
		return nil, fmt.Errorf("dense layer sizes must be positive: in=%d, out=%d", inSize, outSize)
	}

	// Xavier initialization: scale = sqrt(6 / (fan_in + fan_out))
	scale := math.Sqrt(6.0 / float64(inSize+outSize))

	wData := make([]float64, inSize*outSize)
	for i := range wData {
		wData[i] = (2*rng() - 1) * scale
	}
	weights := numgo.NewNDArray([]int{inSize, outSize}, wData)

	bias := numgo.Zeros(outSize)

	return &Dense{
		weights: weights,
		bias:    bias,
		inSize:  inSize,
		outSize: outSize,
	}, nil
}

// Forward computes the layer output: output = input @ weights + bias.
// Input shape: (batch, inSize). Output shape: (batch, outSize).
func (d *Dense) Forward(input *numgo.NDArray) (*numgo.NDArray, error) {
	shape := input.Shape()
	if len(shape) != 2 || shape[1] != d.inSize {
		return nil, fmt.Errorf("dense forward: expected input shape (batch, %d), got %v", d.inSize, shape)
	}

	batch := shape[0]
	outData := make([]float64, batch*d.outSize)
	inData := input.Data()
	wData := d.weights.Data()
	bData := d.bias.Data()

	// Matrix multiply: input[batch, inSize] @ weights[inSize, outSize] + bias[outSize]
	for b := range batch {
		for o := range d.outSize {
			sum := bData[o]
			for k := range d.inSize {
				sum += inData[b*d.inSize+k] * wData[k*d.outSize+o]
			}
			outData[b*d.outSize+o] = sum
		}
	}

	return numgo.NewNDArray([]int{batch, d.outSize}, outData), nil
}

// Weights returns a copy of the layer's weight tensor.
func (d *Dense) Weights() *numgo.NDArray {
	return d.weights.Copy()
}

// Bias returns a copy of the layer's bias tensor.
func (d *Dense) Bias() *numgo.NDArray {
	return d.bias.Copy()
}

// SetWeights replaces the layer's weights. Shape must be (inSize, outSize).
func (d *Dense) SetWeights(weights *numgo.NDArray) error {
	expectedShape := []int{d.inSize, d.outSize}
	if !shapeEq(weights.Shape(), expectedShape) {
		return fmt.Errorf("expected weight shape %v, got %v", expectedShape, weights.Shape())
	}
	d.weights = weights.Copy()
	return nil
}

// SetBias replaces the layer's bias. Shape must be (outSize).
func (d *Dense) SetBias(bias *numgo.NDArray) error {
	expectedShape := []int{d.outSize}
	if !shapeEq(bias.Shape(), expectedShape) {
		return fmt.Errorf("expected bias shape %v, got %v", expectedShape, bias.Shape())
	}
	d.bias = bias.Copy()
	return nil
}

// InSize returns the input feature size.
func (d *Dense) InSize() int {
	return d.inSize
}

// OutSize returns the output feature size.
func (d *Dense) OutSize() int {
	return d.outSize
}

// shapeEq compares two shapes for equality.
func shapeEq(a, b []int) bool {
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
