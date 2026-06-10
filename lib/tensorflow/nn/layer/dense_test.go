package layer

import (
	"math"
	"testing"

	"github.com/asymmetric-effort/datascience/lib/numgo"
)

// deterministicRNG returns a simple deterministic RNG for testing.
func deterministicRNG() func() float64 {
	state := 0.5
	return func() float64 {
		state = math.Mod(state*1.1+0.3, 1.0)
		return state
	}
}

func TestNewDense(t *testing.T) {
	layer, err := NewDense(4, 3, deterministicRNG())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if layer.InSize() != 4 {
		t.Errorf("InSize() = %d, want 4", layer.InSize())
	}
	if layer.OutSize() != 3 {
		t.Errorf("OutSize() = %d, want 3", layer.OutSize())
	}
}

func TestNewDenseInvalidSize(t *testing.T) {
	_, err := NewDense(0, 3, deterministicRNG())
	if err == nil {
		t.Error("expected error for zero input size")
	}
	_, err = NewDense(3, -1, deterministicRNG())
	if err == nil {
		t.Error("expected error for negative output size")
	}
}

func TestDenseForward(t *testing.T) {
	layer, _ := NewDense(3, 2, deterministicRNG())

	// Set known weights for deterministic test.
	weights := numgo.NewNDArray([]int{3, 2}, []float64{
		1, 0,
		0, 1,
		1, 1,
	})
	bias := numgo.NewNDArray([]int{2}, []float64{0.5, -0.5})

	_ = layer.SetWeights(weights)
	_ = layer.SetBias(bias)

	// Input: batch=2, features=3
	input := numgo.NewNDArray([]int{2, 3}, []float64{
		1, 2, 3, // row 0: (1*1+2*0+3*1)+0.5=4.5, (1*0+2*1+3*1)-0.5=4.5
		4, 5, 6, // row 1: (4*1+5*0+6*1)+0.5=10.5, (4*0+5*1+6*1)-0.5=10.5
	})

	output, err := layer.Forward(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedShape := []int{2, 2}
	if !shapeEq(output.Shape(), expectedShape) {
		t.Errorf("output shape = %v, want %v", output.Shape(), expectedShape)
	}

	want := []float64{4.5, 4.5, 10.5, 10.5}
	data := output.Data()
	for i, v := range data {
		if math.Abs(v-want[i]) > 1e-10 {
			t.Errorf("output[%d] = %f, want %f", i, v, want[i])
		}
	}
}

func TestDenseForwardBadInput(t *testing.T) {
	layer, _ := NewDense(3, 2, deterministicRNG())

	// Wrong rank.
	input1 := numgo.NewNDArray([]int{6}, []float64{1, 2, 3, 4, 5, 6})
	_, err := layer.Forward(input1)
	if err == nil {
		t.Error("expected error for wrong rank")
	}

	// Wrong feature size.
	input2 := numgo.NewNDArray([]int{2, 4}, []float64{1, 2, 3, 4, 5, 6, 7, 8})
	_, err = layer.Forward(input2)
	if err == nil {
		t.Error("expected error for wrong feature size")
	}
}

func TestDenseWeightsAndBias(t *testing.T) {
	layer, _ := NewDense(2, 3, deterministicRNG())

	w := layer.Weights()
	expectedWShape := []int{2, 3}
	if !shapeEq(w.Shape(), expectedWShape) {
		t.Errorf("weights shape = %v, want %v", w.Shape(), expectedWShape)
	}

	b := layer.Bias()
	expectedBShape := []int{3}
	if !shapeEq(b.Shape(), expectedBShape) {
		t.Errorf("bias shape = %v, want %v", b.Shape(), expectedBShape)
	}
	// Bias should be initialized to zeros.
	for i, v := range b.Data() {
		if v != 0 {
			t.Errorf("bias[%d] = %f, want 0", i, v)
		}
	}
}

func TestDenseSetWeightsWrongShape(t *testing.T) {
	layer, _ := NewDense(2, 3, deterministicRNG())
	badWeights := numgo.Zeros(3, 3)
	err := layer.SetWeights(badWeights)
	if err == nil {
		t.Error("expected error for wrong weight shape")
	}
}

func TestDenseSetBiasWrongShape(t *testing.T) {
	layer, _ := NewDense(2, 3, deterministicRNG())
	badBias := numgo.Zeros(2)
	err := layer.SetBias(badBias)
	if err == nil {
		t.Error("expected error for wrong bias shape")
	}
}
