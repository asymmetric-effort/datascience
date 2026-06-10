package layer

import (
	"fmt"
	"math"

	"github.com/asymmetric-effort/datascience/lib/numgo"
)

// LSTM implements a Long Short-Term Memory recurrent layer.
// Input shape: (batch, timeSteps, features).
// Output shape: (batch, timeSteps, units) if returnSequences, else (batch, units).
// Analogous to tf.keras.layers.LSTM.
type LSTM struct {
	units           int
	inputSize       int
	returnSequences bool
	// Combined weight matrices: W for input, U for recurrent.
	// W shape: (inputSize, 4*units) — gates: [i, f, c, o]
	// U shape: (units, 4*units)
	// bias shape: (4*units)
	wData    []float64
	uData    []float64
	biasData []float64
}

// NewLSTM creates a new LSTM layer.
func NewLSTM(inputSize, units int, returnSequences bool, rng func() float64) (*LSTM, error) {
	if inputSize <= 0 || units <= 0 {
		return nil, fmt.Errorf("lstm: inputSize and units must be positive, got %d, %d", inputSize, units)
	}

	gateSize := 4 * units
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
	// Initialize forget gate bias to 1 for better gradient flow.
	for i := units; i < 2*units; i++ {
		biasData[i] = 1.0
	}

	return &LSTM{
		units:           units,
		inputSize:       inputSize,
		returnSequences: returnSequences,
		wData:           wData,
		uData:           uData,
		biasData:        biasData,
	}, nil
}

// Forward runs the LSTM over the input sequence.
func (l *LSTM) Forward(input *numgo.NDArray) (*numgo.NDArray, error) {
	shape := input.Shape()
	if len(shape) != 3 || shape[2] != l.inputSize {
		return nil, fmt.Errorf("lstm: expected (batch, timeSteps, %d), got %v", l.inputSize, shape)
	}
	batch, timeSteps := shape[0], shape[1]
	inData := input.Data()

	gateSize := 4 * l.units

	// State buffers.
	h := make([]float64, batch*l.units) // hidden state
	c := make([]float64, batch*l.units) // cell state
	gates := make([]float64, batch*gateSize)

	var allH []float64
	if l.returnSequences {
		allH = make([]float64, batch*timeSteps*l.units)
	}

	for t := range timeSteps {
		// Clear gates.
		for i := range gates {
			gates[i] = 0
		}

		// Compute gates = x @ W + h @ U + bias.
		for b := range batch {
			xOff := (b*timeSteps + t) * l.inputSize
			gOff := b * gateSize

			// x @ W
			for g := range gateSize {
				sum := l.biasData[g]
				for k := range l.inputSize {
					sum += inData[xOff+k] * l.wData[k*gateSize+g]
				}
				gates[gOff+g] = sum
			}

			// + h @ U
			hOff := b * l.units
			for g := range gateSize {
				sum := 0.0
				for k := range l.units {
					sum += h[hOff+k] * l.uData[k*gateSize+g]
				}
				gates[gOff+g] += sum
			}

			// Apply activations and update states.
			for u := range l.units {
				iGate := sigmoid(gates[gOff+u])                  // input gate
				fGate := sigmoid(gates[gOff+l.units+u])          // forget gate
				cCandidate := math.Tanh(gates[gOff+2*l.units+u]) // cell candidate
				oGate := sigmoid(gates[gOff+3*l.units+u])        // output gate

				c[hOff+u] = fGate*c[hOff+u] + iGate*cCandidate
				h[hOff+u] = oGate * math.Tanh(c[hOff+u])
			}
		}

		if l.returnSequences {
			for b := range batch {
				srcOff := b * l.units
				dstOff := (b*timeSteps + t) * l.units
				copy(allH[dstOff:dstOff+l.units], h[srcOff:srcOff+l.units])
			}
		}
	}

	if l.returnSequences {
		return numgo.NewNDArray([]int{batch, timeSteps, l.units}, allH), nil
	}
	return numgo.NewNDArray([]int{batch, l.units}, h), nil
}

// Units returns the number of LSTM units.
func (l *LSTM) Units() int {
	return l.units
}

// ReturnSequences returns whether the layer outputs the full sequence.
func (l *LSTM) ReturnSequences() bool {
	return l.returnSequences
}

func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}
