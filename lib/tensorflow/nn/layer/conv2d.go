package layer

import (
	"fmt"
	"math"

	"github.com/asymmetric-effort/datascience/lib/numgo"
)

// Conv2D implements a 2D convolution layer.
// Input shape: (batch, height, width, inChannels) — NHWC format.
// Output shape: (batch, outHeight, outWidth, outChannels).
// Analogous to tf.keras.layers.Conv2D.
type Conv2D struct {
	filters     *numgo.NDArray // shape: (kernelH, kernelW, inChannels, outChannels)
	bias        *numgo.NDArray // shape: (outChannels)
	kernelH     int
	kernelW     int
	inChannels  int
	outChannels int
	strideH     int
	strideW     int
	padSame     bool
}

// NewConv2D creates a new Conv2D layer.
// kernelSize is (height, width), stride is (strideH, strideW).
// padSame: if true, output spatial dims match input (zero-padded).
func NewConv2D(inChannels, outChannels, kernelH, kernelW, strideH, strideW int, padSame bool, rng func() float64) (*Conv2D, error) {
	if inChannels <= 0 || outChannels <= 0 || kernelH <= 0 || kernelW <= 0 {
		return nil, fmt.Errorf("conv2d: all dimensions must be positive")
	}
	if strideH <= 0 || strideW <= 0 {
		return nil, fmt.Errorf("conv2d: strides must be positive")
	}

	// He initialization for filters.
	fanIn := float64(kernelH * kernelW * inChannels)
	scale := math.Sqrt(2.0 / fanIn)

	fData := make([]float64, kernelH*kernelW*inChannels*outChannels)
	for i := range fData {
		fData[i] = (2*rng() - 1) * scale
	}
	filters := numgo.NewNDArray([]int{kernelH, kernelW, inChannels, outChannels}, fData)

	bias := numgo.Zeros(outChannels)

	return &Conv2D{
		filters:     filters,
		bias:        bias,
		kernelH:     kernelH,
		kernelW:     kernelW,
		inChannels:  inChannels,
		outChannels: outChannels,
		strideH:     strideH,
		strideW:     strideW,
		padSame:     padSame,
	}, nil
}

// Forward performs the 2D convolution.
// Input shape: (batch, height, width, inChannels).
func (c *Conv2D) Forward(input *numgo.NDArray) (*numgo.NDArray, error) {
	shape := input.Shape()
	if len(shape) != 4 || shape[3] != c.inChannels {
		return nil, fmt.Errorf("conv2d forward: expected (batch, H, W, %d), got %v", c.inChannels, shape)
	}
	batch, inH, inW := shape[0], shape[1], shape[2]

	var padH, padW int
	if c.padSame {
		padH = (c.kernelH - 1) / 2
		padW = (c.kernelW - 1) / 2
	}

	outH := (inH+2*padH-c.kernelH)/c.strideH + 1
	outW := (inW+2*padW-c.kernelW)/c.strideW + 1

	outSize := batch * outH * outW * c.outChannels
	outData := make([]float64, outSize)
	inData := input.Data()
	fData := c.filters.Data()
	bData := c.bias.Data()

	for b := range batch {
		for oh := range outH {
			for ow := range outW {
				for oc := range c.outChannels {
					sum := bData[oc]
					for kh := range c.kernelH {
						for kw := range c.kernelW {
							ih := oh*c.strideH + kh - padH
							iw := ow*c.strideW + kw - padW
							if ih < 0 || ih >= inH || iw < 0 || iw >= inW {
								continue
							}
							for ic := range c.inChannels {
								inIdx := b*inH*inW*c.inChannels + ih*inW*c.inChannels + iw*c.inChannels + ic
								fIdx := kh*c.kernelW*c.inChannels*c.outChannels + kw*c.inChannels*c.outChannels + ic*c.outChannels + oc
								sum += inData[inIdx] * fData[fIdx]
							}
						}
					}
					outIdx := b*outH*outW*c.outChannels + oh*outW*c.outChannels + ow*c.outChannels + oc
					outData[outIdx] = sum
				}
			}
		}
	}

	return numgo.NewNDArray([]int{batch, outH, outW, c.outChannels}, outData), nil
}

// Filters returns a copy of the filter tensor.
func (c *Conv2D) Filters() *numgo.NDArray {
	return c.filters.Copy()
}

// OutChannels returns the number of output channels.
func (c *Conv2D) OutChannels() int {
	return c.outChannels
}
