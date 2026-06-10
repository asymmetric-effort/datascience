package activation

import (
	"math"
	"testing"

	"github.com/asymmetric-effort/datascience/lib/numgo"
)

func makeArr(t *testing.T, dims []int, data []float64) *numgo.NDArray {
	t.Helper()
	return numgo.NewNDArray(dims, data)
}

func TestReLU(t *testing.T) {
	input := makeArr(t, []int{5}, []float64{-2, -1, 0, 1, 2})
	result := ReLU(input)
	want := []float64{0, 0, 0, 1, 2}
	for i, v := range result.Data() {
		if v != want[i] {
			t.Errorf("result[%d] = %f, want %f", i, v, want[i])
		}
	}
}

func TestSigmoid(t *testing.T) {
	input := makeArr(t, []int{3}, []float64{0, 100, -100})
	result := Sigmoid(input)
	data := result.Data()
	if math.Abs(data[0]-0.5) > 1e-10 {
		t.Errorf("sigmoid(0) = %f, want 0.5", data[0])
	}
	if math.Abs(data[1]-1.0) > 1e-5 {
		t.Errorf("sigmoid(100) = %f, want ~1.0", data[1])
	}
	if math.Abs(data[2]) > 1e-5 {
		t.Errorf("sigmoid(-100) = %f, want ~0.0", data[2])
	}
}

func TestTanh(t *testing.T) {
	input := makeArr(t, []int{3}, []float64{0, 1, -1})
	result := Tanh(input)
	data := result.Data()
	if math.Abs(data[0]) > 1e-10 {
		t.Errorf("tanh(0) = %f, want 0", data[0])
	}
	if math.Abs(data[1]-math.Tanh(1)) > 1e-10 {
		t.Errorf("tanh(1) = %f, want %f", data[1], math.Tanh(1))
	}
}

func TestSoftmax1D(t *testing.T) {
	input := makeArr(t, []int{3}, []float64{1, 2, 3})
	result := Softmax(input)
	data := result.Data()
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	if math.Abs(sum-1.0) > 1e-10 {
		t.Errorf("softmax sum = %f, want 1.0", sum)
	}
	if data[0] >= data[1] || data[1] >= data[2] {
		t.Errorf("softmax not monotonically increasing: %v", data)
	}
	for i, v := range data {
		if v <= 0 {
			t.Errorf("softmax[%d] = %f, expected positive", i, v)
		}
	}
}

func TestSoftmax2D(t *testing.T) {
	input := makeArr(t, []int{2, 3}, []float64{1, 2, 3, 1, 2, 3})
	result := Softmax(input)
	data := result.Data()
	row1Sum := data[0] + data[1] + data[2]
	row2Sum := data[3] + data[4] + data[5]
	if math.Abs(row1Sum-1.0) > 1e-10 {
		t.Errorf("row 1 sum = %f, want 1.0", row1Sum)
	}
	if math.Abs(row2Sum-1.0) > 1e-10 {
		t.Errorf("row 2 sum = %f, want 1.0", row2Sum)
	}
}

func TestLeakyReLU(t *testing.T) {
	input := makeArr(t, []int{4}, []float64{-2, -1, 0, 1})
	result := LeakyReLU(input, 0.1)
	want := []float64{-0.2, -0.1, 0, 1}
	for i, v := range result.Data() {
		if math.Abs(v-want[i]) > 1e-10 {
			t.Errorf("result[%d] = %f, want %f", i, v, want[i])
		}
	}
}

func TestELU(t *testing.T) {
	input := makeArr(t, []int{3}, []float64{-1, 0, 1})
	result := ELU(input, 1.0)
	data := result.Data()
	wantNeg := 1.0 * (math.Exp(-1) - 1)
	if math.Abs(data[0]-wantNeg) > 1e-10 {
		t.Errorf("elu(-1) = %f, want %f", data[0], wantNeg)
	}
	if data[1] != 0 {
		t.Errorf("elu(0) = %f, want 0", data[1])
	}
	if data[2] != 1 {
		t.Errorf("elu(1) = %f, want 1", data[2])
	}
}

func TestGELU(t *testing.T) {
	input := makeArr(t, []int{3}, []float64{-1, 0, 1})
	result := GELU(input)
	data := result.Data()
	if math.Abs(data[1]) > 1e-10 {
		t.Errorf("gelu(0) = %f, want 0", data[1])
	}
	if data[0] >= 0 {
		t.Errorf("gelu(-1) = %f, expected negative", data[0])
	}
	if data[2] <= 0 || data[2] >= 1 {
		t.Errorf("gelu(1) = %f, expected in (0, 1)", data[2])
	}
}

func TestSwish(t *testing.T) {
	input := makeArr(t, []int{3}, []float64{-1, 0, 1})
	result := Swish(input)
	data := result.Data()
	if math.Abs(data[1]) > 1e-10 {
		t.Errorf("swish(0) = %f, want 0", data[1])
	}
	expected := 1.0 / (1.0 + math.Exp(-1.0))
	if math.Abs(data[2]-expected) > 1e-10 {
		t.Errorf("swish(1) = %f, want %f", data[2], expected)
	}
}

func TestReLUGrad(t *testing.T) {
	input := makeArr(t, []int{4}, []float64{-1, 0, 1, 2})
	gradOut := makeArr(t, []int{4}, []float64{1, 1, 1, 1})
	result := ReLUGrad(input, gradOut)
	want := []float64{0, 0, 1, 1}
	for i, v := range result.Data() {
		if v != want[i] {
			t.Errorf("result[%d] = %f, want %f", i, v, want[i])
		}
	}
}

func TestSigmoidGrad(t *testing.T) {
	sigOut := makeArr(t, []int{1}, []float64{0.5})
	gradOut := makeArr(t, []int{1}, []float64{1.0})
	result := SigmoidGrad(sigOut, gradOut)
	if math.Abs(result.Data()[0]-0.25) > 1e-10 {
		t.Errorf("grad = %f, want 0.25", result.Data()[0])
	}
}
