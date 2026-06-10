package layer

import (
	"fmt"
	"math"

	"github.com/asymmetric-effort/datascience/lib/numgo"
)

// MaxPool2D implements 2D max pooling.
// Input shape: (batch, height, width, channels) — NHWC format.
// Analogous to tf.keras.layers.MaxPool2D.
type MaxPool2D struct {
	poolH   int
	poolW   int
	strideH int
	strideW int
}

// NewMaxPool2D creates a new MaxPool2D layer.
func NewMaxPool2D(poolH, poolW, strideH, strideW int) (*MaxPool2D, error) {
	if poolH <= 0 || poolW <= 0 || strideH <= 0 || strideW <= 0 {
		return nil, fmt.Errorf("maxpool2d: all dimensions must be positive")
	}
	return &MaxPool2D{poolH: poolH, poolW: poolW, strideH: strideH, strideW: strideW}, nil
}

// Forward performs 2D max pooling.
func (m *MaxPool2D) Forward(input *numgo.NDArray) (*numgo.NDArray, error) {
	shape := input.Shape()
	if len(shape) != 4 {
		return nil, fmt.Errorf("maxpool2d: expected rank-4 input (NHWC), got rank %d", len(shape))
	}
	batch, inH, inW, channels := shape[0], shape[1], shape[2], shape[3]
	outH := (inH-m.poolH)/m.strideH + 1
	outW := (inW-m.poolW)/m.strideW + 1

	if outH <= 0 || outW <= 0 {
		return nil, fmt.Errorf("maxpool2d: output dimensions non-positive (%d, %d)", outH, outW)
	}

	outSize := batch * outH * outW * channels
	outData := make([]float64, outSize)
	inData := input.Data()

	for b := range batch {
		for oh := range outH {
			for ow := range outW {
				for c := range channels {
					maxVal := math.Inf(-1)
					for kh := range m.poolH {
						for kw := range m.poolW {
							ih := oh*m.strideH + kh
							iw := ow*m.strideW + kw
							idx := b*inH*inW*channels + ih*inW*channels + iw*channels + c
							if inData[idx] > maxVal {
								maxVal = inData[idx]
							}
						}
					}
					outIdx := b*outH*outW*channels + oh*outW*channels + ow*channels + c
					outData[outIdx] = maxVal
				}
			}
		}
	}

	return numgo.NewNDArray([]int{batch, outH, outW, channels}, outData), nil
}
