package initializer

import (
	"math"
	"testing"
)

func TestZerosInit(t *testing.T) {
	init := ZerosInit{}
	shape := []int{3, 2}
	result, err := init.Initialize(shape)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, v := range result.Data() {
		if v != 0 {
			t.Errorf("data[%d] = %f, want 0", i, v)
		}
	}
}

func TestOnesInit(t *testing.T) {
	init := OnesInit{}
	shape := []int{2, 3}
	result, err := init.Initialize(shape)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, v := range result.Data() {
		if v != 1 {
			t.Errorf("data[%d] = %f, want 1", i, v)
		}
	}
}

func TestConstantInit(t *testing.T) {
	init := ConstantInit{Value: 42.0}
	shape := []int{4}
	result, err := init.Initialize(shape)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, v := range result.Data() {
		if v != 42.0 {
			t.Errorf("data[%d] = %f, want 42", i, v)
		}
	}
}

func TestRandomNormalInit(t *testing.T) {
	init := RandomNormalInit{Mean: 0, StdDev: 1, Seed: 42}
	shape := []int{1000}
	result, err := init.Initialize(shape)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data := result.Data()
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	mean := sum / float64(len(data))
	if math.Abs(mean) > 0.2 {
		t.Errorf("mean = %f, expected near 0", mean)
	}
}

func TestRandomUniformInit(t *testing.T) {
	init := RandomUniformInit{MinVal: -1, MaxVal: 1, Seed: 42}
	shape := []int{1000}
	result, err := init.Initialize(shape)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, v := range result.Data() {
		if v < -1 || v >= 1 {
			t.Errorf("data[%d] = %f, out of range [-1, 1)", i, v)
		}
	}
}

func TestGlorotUniformInit(t *testing.T) {
	init := GlorotUniformInit{Seed: 42}
	shape := []int{100, 50}
	result, err := init.Initialize(shape)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	limit := math.Sqrt(6.0 / 150.0)
	for i, v := range result.Data() {
		if v < -limit || v >= limit {
			t.Errorf("data[%d] = %f, out of Glorot range [-%f, %f)", i, v, limit, limit)
		}
	}
}

func TestGlorotNormalInit(t *testing.T) {
	init := GlorotNormalInit{Seed: 42}
	shape := []int{100, 50}
	result, err := init.Initialize(shape)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Size() != 5000 {
		t.Errorf("got %d elements, want 5000", result.Size())
	}
}

func TestHeNormalInit(t *testing.T) {
	init := HeNormalInit{Seed: 42}
	shape := []int{100, 50}
	result, err := init.Initialize(shape)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data := result.Data()
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	mean := sum / float64(len(data))
	if math.Abs(mean) > 0.2 {
		t.Errorf("mean = %f, expected near 0", mean)
	}
}

func TestHeUniformInit(t *testing.T) {
	init := HeUniformInit{Seed: 42}
	shape := []int{100, 50}
	result, err := init.Initialize(shape)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	limit := math.Sqrt(6.0 / 100.0)
	for i, v := range result.Data() {
		if v < -limit || v >= limit {
			t.Errorf("data[%d] = %f, out of He range [-%f, %f)", i, v, limit, limit)
		}
	}
}

func TestLecunNormalInit(t *testing.T) {
	init := LecunNormalInit{Seed: 42}
	shape := []int{100, 50}
	result, err := init.Initialize(shape)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Size() != 5000 {
		t.Errorf("got %d elements, want 5000", result.Size())
	}
}

func TestComputeFansScalar(t *testing.T) {
	shape := []int{}
	fanIn, fanOut := computeFans(shape)
	if fanIn != 1 || fanOut != 1 {
		t.Errorf("scalar fans = (%d, %d), want (1, 1)", fanIn, fanOut)
	}
}

func TestComputeFans1D(t *testing.T) {
	shape := []int{10}
	fanIn, fanOut := computeFans(shape)
	if fanIn != 10 || fanOut != 10 {
		t.Errorf("1D fans = (%d, %d), want (10, 10)", fanIn, fanOut)
	}
}

func TestComputeFans2D(t *testing.T) {
	shape := []int{100, 50}
	fanIn, fanOut := computeFans(shape)
	if fanIn != 100 || fanOut != 50 {
		t.Errorf("2D fans = (%d, %d), want (100, 50)", fanIn, fanOut)
	}
}

func TestComputeFans4D(t *testing.T) {
	// Conv kernel: (3, 3, 64, 128)
	shape := []int{3, 3, 64, 128}
	fanIn, fanOut := computeFans(shape)
	// fanIn = 64 * 9 = 576, fanOut = 128 * 9 = 1152
	if fanIn != 576 || fanOut != 1152 {
		t.Errorf("4D fans = (%d, %d), want (576, 1152)", fanIn, fanOut)
	}
}

func TestInitializerInterface(t *testing.T) {
	// Verify all initializers implement the Initializer interface.
	var _ Initializer = ZerosInit{}
	var _ Initializer = OnesInit{}
	var _ Initializer = ConstantInit{}
	var _ Initializer = RandomNormalInit{}
	var _ Initializer = RandomUniformInit{}
	var _ Initializer = GlorotUniformInit{}
	var _ Initializer = GlorotNormalInit{}
	var _ Initializer = HeNormalInit{}
	var _ Initializer = HeUniformInit{}
	var _ Initializer = LecunNormalInit{}
}
