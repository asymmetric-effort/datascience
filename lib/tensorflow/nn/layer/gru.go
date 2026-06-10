package layer

import (
	"fmt"
	"math"

	"github.com/asymmetric-effort/datascience/lib/numgo"
)

// GRU implements a Gated Recurrent Unit layer.
// Input shape: (batch, timeSteps, features).
// Output shape: (batch, timeSteps, units) if returnSequences, else (batch, units).
// Analogous to tf.keras.layers.GRU.
type GRU struct {
	units           int
	inputSize       int
	returnSequences bool
	// W shape: (inputSize, 3*units) — gates: [z (update), r (reset), h (candidate)]
	// U shape: (units, 3*units)
	// bias shape: (3*units)
	wData    []float64
	uData    []float64
	biasData []float64
}

// NewGRU creates a new GRU layer.
func NewGRU(inputSize, units int, returnSequences bool, rng func() float64) (*GRU, error) {
	if inputSize <= 0 || units <= 0 {
		return nil, fmt.Errorf("gru: inputSize and units must be positive, got %d, %d", inputSize, units)
	}

	gateSize := 3 * units
	scale := math.Sqrt(2.0 / float64(inputSize+units))

	wData := make([]float64, inputSize*gateSize)
	for i := range wData {
		wData[i] = (2*rng() - 1) * scale
	}
	uData := make([]float64, units*gateSize)
	for i := range uData {
		uData[i] = (2*rng() - 1) * scale
	}
	biasData := make([]float64, gateSize)

	return &GRU{
		units:           units,
		inputSize:       inputSize,
		returnSequences: returnSequences,
		wData:           wData,
		uData:           uData,
		biasData:        biasData,
	}, nil
}

// Forward runs the GRU over the input sequence.
func (g *GRU) Forward(input *numgo.NDArray) (*numgo.NDArray, error) {
	shape := input.Shape()
	if len(shape) != 3 || shape[2] != g.inputSize {
		return nil, fmt.Errorf("gru: expected (batch, timeSteps, %d), got %v", g.inputSize, shape)
	}
	batch, timeSteps := shape[0], shape[1]
	inData := input.Data()

	gateSize := 3 * g.units

	h := make([]float64, batch*g.units)
	xGates := make([]float64, batch*gateSize)
	hGates := make([]float64, batch*gateSize)

	var allH []float64
	if g.returnSequences {
		allH = make([]float64, batch*timeSteps*g.units)
	}

	for t := range timeSteps {
		// Compute x @ W + bias.
		for i := range xGates {
			xGates[i] = 0
		}
		for b := range batch {
			xOff := (b*timeSteps + t) * g.inputSize
			gOff := b * gateSize
			for gi := range gateSize {
				sum := g.biasData[gi]
				for k := range g.inputSize {
					sum += inData[xOff+k] * g.wData[k*gateSize+gi]
				}
				xGates[gOff+gi] = sum
			}
		}

		// Compute h @ U.
		for i := range hGates {
			hGates[i] = 0
		}
		for b := range batch {
			hOff := b * g.units
			gOff := b * gateSize
			for gi := range gateSize {
				sum := 0.0
				for k := range g.units {
					sum += h[hOff+k] * g.uData[k*gateSize+gi]
				}
				hGates[gOff+gi] = sum
			}
		}

		// Apply gates and update hidden state.
		for b := range batch {
			hOff := b * g.units
			gOff := b * gateSize

			for u := range g.units {
				zGate := sigmoid(xGates[gOff+u] + hGates[gOff+u])                                  // update gate
				rGate := sigmoid(xGates[gOff+g.units+u] + hGates[gOff+g.units+u])                  // reset gate
				hCandidate := math.Tanh(xGates[gOff+2*g.units+u] + rGate*hGates[gOff+2*g.units+u]) // candidate

				h[hOff+u] = (1-zGate)*h[hOff+u] + zGate*hCandidate
			}
		}

		if g.returnSequences {
			for b := range batch {
				srcOff := b * g.units
				dstOff := (b*timeSteps + t) * g.units
				copy(allH[dstOff:dstOff+g.units], h[srcOff:srcOff+g.units])
			}
		}
	}

	if g.returnSequences {
		return numgo.NewNDArray([]int{batch, timeSteps, g.units}, allH), nil
	}
	return numgo.NewNDArray([]int{batch, g.units}, h), nil
}

// Units returns the number of GRU units.
func (g *GRU) Units() int {
	return g.units
}

// ReturnSequences returns whether the layer outputs the full sequence.
func (g *GRU) ReturnSequences() bool {
	return g.returnSequences
}
