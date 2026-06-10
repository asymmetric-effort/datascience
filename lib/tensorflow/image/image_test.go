package image

import (
	"math"
	"testing"

	"github.com/asymmetric-effort/datascience/lib/numgo"
)

func makeImg(t *testing.T, h, w, c int, data []float64) *numgo.NDArray {
	t.Helper()
	return numgo.NewNDArray([]int{h, w, c}, data)
}

func TestResize(t *testing.T) {
	// 2x2 grayscale -> 4x4
	img := makeImg(t, 2, 2, 1, []float64{0, 1, 0, 1})
	result, err := Resize(img, 4, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedShape := []int{4, 4, 1}
	if !shapeEqual(result.Shape(), expectedShape) {
		t.Errorf("shape = %v, want %v", result.Shape(), expectedShape)
	}
}

func TestResizeWrongRank(t *testing.T) {
	arr := numgo.Zeros(4, 4)
	_, err := Resize(arr, 2, 2)
	if err == nil {
		t.Error("expected error for wrong rank")
	}
}

func TestResizeInvalidDims(t *testing.T) {
	img := makeImg(t, 2, 2, 1, []float64{0, 0, 0, 0})
	_, err := Resize(img, 0, 2)
	if err == nil {
		t.Error("expected error for zero dimension")
	}
}

func TestCropCenter(t *testing.T) {
	img := makeImg(t, 4, 4, 1, []float64{
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16,
	})
	result, err := CropCenter(img, 2, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []float64{6, 7, 10, 11}
	for i, v := range result.Data() {
		if v != want[i] {
			t.Errorf("data[%d] = %f, want %f", i, v, want[i])
		}
	}
}

func TestCropCenterInvalid(t *testing.T) {
	img := makeImg(t, 2, 2, 1, []float64{1, 2, 3, 4})
	_, err := CropCenter(img, 3, 2)
	if err == nil {
		t.Error("expected error for crop larger than image")
	}
}

func TestCropCenterWrongRank(t *testing.T) {
	arr := numgo.Zeros(4)
	_, err := CropCenter(arr, 2, 2)
	if err == nil {
		t.Error("expected error for wrong rank")
	}
}

func TestFlipLeftRight(t *testing.T) {
	img := makeImg(t, 2, 3, 1, []float64{1, 2, 3, 4, 5, 6})
	result, err := FlipLeftRight(img)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []float64{3, 2, 1, 6, 5, 4}
	for i, v := range result.Data() {
		if v != want[i] {
			t.Errorf("data[%d] = %f, want %f", i, v, want[i])
		}
	}
}

func TestFlipLeftRightWrongRank(t *testing.T) {
	arr := numgo.Zeros(4)
	_, err := FlipLeftRight(arr)
	if err == nil {
		t.Error("expected error")
	}
}

func TestFlipUpDown(t *testing.T) {
	img := makeImg(t, 2, 2, 1, []float64{1, 2, 3, 4})
	result, err := FlipUpDown(img)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []float64{3, 4, 1, 2}
	for i, v := range result.Data() {
		if v != want[i] {
			t.Errorf("data[%d] = %f, want %f", i, v, want[i])
		}
	}
}

func TestFlipUpDownWrongRank(t *testing.T) {
	arr := numgo.Zeros(4)
	_, err := FlipUpDown(arr)
	if err == nil {
		t.Error("expected error")
	}
}

func TestAdjustBrightness(t *testing.T) {
	img := makeImg(t, 1, 2, 1, []float64{0.3, 0.7})
	result, err := AdjustBrightness(img, 0.1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data := result.Data()
	if math.Abs(data[0]-0.4) > 1e-10 {
		t.Errorf("data[0] = %f, want 0.4", data[0])
	}
	if math.Abs(data[1]-0.8) > 1e-10 {
		t.Errorf("data[1] = %f, want 0.8", data[1])
	}
}

func TestAdjustBrightnessClamp(t *testing.T) {
	img := makeImg(t, 1, 1, 1, []float64{0.9})
	result, _ := AdjustBrightness(img, 0.5)
	if result.Data()[0] != 1.0 {
		t.Errorf("clamped value = %f, want 1.0", result.Data()[0])
	}
}

func TestAdjustBrightnessWrongRank(t *testing.T) {
	arr := numgo.Zeros(4)
	_, err := AdjustBrightness(arr, 0.1)
	if err == nil {
		t.Error("expected error")
	}
}

func TestAdjustContrast(t *testing.T) {
	img := makeImg(t, 1, 2, 1, []float64{0.4, 0.6})
	result, err := AdjustContrast(img, 2.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data := result.Data()
	// mean=0.5, adjusted: 0.5+2*(0.4-0.5)=0.3, 0.5+2*(0.6-0.5)=0.7
	if math.Abs(data[0]-0.3) > 1e-10 {
		t.Errorf("data[0] = %f, want 0.3", data[0])
	}
}

func TestAdjustContrastWrongRank(t *testing.T) {
	arr := numgo.Zeros(4)
	_, err := AdjustContrast(arr, 1.0)
	if err == nil {
		t.Error("expected error")
	}
}

func TestGrayscale(t *testing.T) {
	img := makeImg(t, 1, 1, 3, []float64{1, 0, 0}) // pure red
	result, err := Grayscale(img)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedShape := []int{1, 1, 1}
	if !shapeEqual(result.Shape(), expectedShape) {
		t.Errorf("shape = %v, want %v", result.Shape(), expectedShape)
	}
	if math.Abs(result.Data()[0]-0.2989) > 1e-4 {
		t.Errorf("gray = %f, want ~0.2989", result.Data()[0])
	}
}

func TestGrayscaleWrongChannels(t *testing.T) {
	img := makeImg(t, 2, 2, 1, []float64{1, 2, 3, 4})
	_, err := Grayscale(img)
	if err == nil {
		t.Error("expected error for non-3 channel input")
	}
}

func TestNormalize(t *testing.T) {
	img := makeImg(t, 1, 3, 1, []float64{50, 100, 200})
	result, err := Normalize(img)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data := result.Data()
	if data[0] != 0 {
		t.Errorf("min should normalize to 0, got %f", data[0])
	}
	if data[2] != 1 {
		t.Errorf("max should normalize to 1, got %f", data[2])
	}
}

func TestNormalizeConstant(t *testing.T) {
	img := makeImg(t, 1, 2, 1, []float64{5, 5})
	result, _ := Normalize(img)
	for _, v := range result.Data() {
		if v != 0 {
			t.Errorf("constant image should normalize to 0, got %f", v)
		}
	}
}

func TestNormalizeWrongRank(t *testing.T) {
	arr := numgo.Zeros(4)
	_, err := Normalize(arr)
	if err == nil {
		t.Error("expected error")
	}
}

func TestRot90(t *testing.T) {
	img := makeImg(t, 2, 3, 1, []float64{1, 2, 3, 4, 5, 6})
	result, err := Rot90(img)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedShape := []int{3, 2, 1}
	if !shapeEqual(result.Shape(), expectedShape) {
		t.Errorf("shape = %v, want %v", result.Shape(), expectedShape)
	}
}

func TestRot90WrongRank(t *testing.T) {
	arr := numgo.Zeros(4)
	_, err := Rot90(arr)
	if err == nil {
		t.Error("expected error")
	}
}

// shapeEqual is a test helper.
func shapeEqual(a, b []int) bool {
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
