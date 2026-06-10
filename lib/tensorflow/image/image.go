// Package image provides image manipulation operations on NDArrays,
// analogous to tf.image.
package image

import (
	"fmt"
	"math"

	"github.com/asymmetric-effort/datascience/lib/numgo"
)

// Resize resizes an image array using bilinear interpolation.
// Input shape: (height, width, channels). Output shape: (newH, newW, channels).
// Analogous to tf.image.resize.
func Resize(img *numgo.NDArray, newH, newW int) (*numgo.NDArray, error) {
	shape := img.Shape()
	if len(shape) != 3 {
		return nil, fmt.Errorf("resize: expected rank-3 array (H, W, C), got rank %d", len(shape))
	}
	if newH <= 0 || newW <= 0 {
		return nil, fmt.Errorf("resize: dimensions must be positive")
	}
	h, w, c := shape[0], shape[1], shape[2]
	data := img.Data()
	outData := make([]float64, newH*newW*c)

	yScale := float64(h) / float64(newH)
	xScale := float64(w) / float64(newW)

	for oy := range newH {
		for ox := range newW {
			srcY := (float64(oy)+0.5)*yScale - 0.5
			srcX := (float64(ox)+0.5)*xScale - 0.5

			y0 := int(math.Floor(srcY))
			x0 := int(math.Floor(srcX))
			y1 := y0 + 1
			x1 := x0 + 1

			dy := srcY - float64(y0)
			dx := srcX - float64(x0)

			y0 = clamp(y0, 0, h-1)
			y1 = clamp(y1, 0, h-1)
			x0 = clamp(x0, 0, w-1)
			x1 = clamp(x1, 0, w-1)

			for ch := range c {
				v00 := data[(y0*w+x0)*c+ch]
				v01 := data[(y0*w+x1)*c+ch]
				v10 := data[(y1*w+x0)*c+ch]
				v11 := data[(y1*w+x1)*c+ch]
				val := v00*(1-dx)*(1-dy) + v01*dx*(1-dy) + v10*(1-dx)*dy + v11*dx*dy
				outData[(oy*newW+ox)*c+ch] = val
			}
		}
	}
	return numgo.NewNDArray([]int{newH, newW, c}, outData), nil
}

// CropCenter crops the center region of an image.
// Analogous to tf.image.central_crop.
func CropCenter(img *numgo.NDArray, cropH, cropW int) (*numgo.NDArray, error) {
	shape := img.Shape()
	if len(shape) != 3 {
		return nil, fmt.Errorf("crop: expected rank-3 array (H, W, C), got rank %d", len(shape))
	}
	h, w, c := shape[0], shape[1], shape[2]
	if cropH > h || cropW > w || cropH <= 0 || cropW <= 0 {
		return nil, fmt.Errorf("crop: invalid crop size (%d, %d) for image (%d, %d)", cropH, cropW, h, w)
	}

	startY := (h - cropH) / 2
	startX := (w - cropW) / 2
	data := img.Data()
	outData := make([]float64, cropH*cropW*c)

	for y := range cropH {
		for x := range cropW {
			srcOff := ((startY+y)*w + (startX + x)) * c
			dstOff := (y*cropW + x) * c
			copy(outData[dstOff:dstOff+c], data[srcOff:srcOff+c])
		}
	}
	return numgo.NewNDArray([]int{cropH, cropW, c}, outData), nil
}

// FlipLeftRight flips an image horizontally.
// Analogous to tf.image.flip_left_right.
func FlipLeftRight(img *numgo.NDArray) (*numgo.NDArray, error) {
	shape := img.Shape()
	if len(shape) != 3 {
		return nil, fmt.Errorf("flip_lr: expected rank-3 array")
	}
	h, w, c := shape[0], shape[1], shape[2]
	data := img.Data()
	outData := make([]float64, len(data))
	for y := range h {
		for x := range w {
			srcOff := (y*w + x) * c
			dstOff := (y*w + (w - 1 - x)) * c
			copy(outData[dstOff:dstOff+c], data[srcOff:srcOff+c])
		}
	}
	return numgo.NewNDArray(shape, outData), nil
}

// FlipUpDown flips an image vertically.
// Analogous to tf.image.flip_up_down.
func FlipUpDown(img *numgo.NDArray) (*numgo.NDArray, error) {
	shape := img.Shape()
	if len(shape) != 3 {
		return nil, fmt.Errorf("flip_ud: expected rank-3 array")
	}
	h, w, c := shape[0], shape[1], shape[2]
	data := img.Data()
	outData := make([]float64, len(data))
	rowSize := w * c
	for y := range h {
		srcOff := y * rowSize
		dstOff := (h - 1 - y) * rowSize
		copy(outData[dstOff:dstOff+rowSize], data[srcOff:srcOff+rowSize])
	}
	return numgo.NewNDArray(shape, outData), nil
}

// AdjustBrightness adds a delta to all pixel values.
// Analogous to tf.image.adjust_brightness.
func AdjustBrightness(img *numgo.NDArray, delta float64) (*numgo.NDArray, error) {
	shape := img.Shape()
	if len(shape) != 3 {
		return nil, fmt.Errorf("adjust_brightness: expected rank-3 array")
	}
	data := img.Data()
	for i := range data {
		data[i] = math.Max(0, math.Min(1, data[i]+delta))
	}
	return numgo.NewNDArray(shape, data), nil
}

// AdjustContrast adjusts the contrast by the given factor.
// Analogous to tf.image.adjust_contrast.
func AdjustContrast(img *numgo.NDArray, factor float64) (*numgo.NDArray, error) {
	shape := img.Shape()
	if len(shape) != 3 {
		return nil, fmt.Errorf("adjust_contrast: expected rank-3 array")
	}
	data := img.Data()
	// Compute mean.
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	mean := sum / float64(len(data))

	for i := range data {
		data[i] = math.Max(0, math.Min(1, mean+factor*(data[i]-mean)))
	}
	return numgo.NewNDArray(shape, data), nil
}

// Grayscale converts an RGB image to grayscale using luminosity method.
// Input shape: (H, W, 3). Output shape: (H, W, 1).
// Analogous to tf.image.rgb_to_grayscale.
func Grayscale(img *numgo.NDArray) (*numgo.NDArray, error) {
	shape := img.Shape()
	if len(shape) != 3 || shape[2] != 3 {
		return nil, fmt.Errorf("grayscale: expected rank-3 array with 3 channels, got %v", shape)
	}
	h, w := shape[0], shape[1]
	data := img.Data()
	outData := make([]float64, h*w)
	for i := range h * w {
		r := data[i*3]
		g := data[i*3+1]
		b := data[i*3+2]
		outData[i] = 0.2989*r + 0.5870*g + 0.1140*b
	}
	return numgo.NewNDArray([]int{h, w, 1}, outData), nil
}

// Normalize normalizes image values to [0, 1] based on min/max.
func Normalize(img *numgo.NDArray) (*numgo.NDArray, error) {
	shape := img.Shape()
	if len(shape) != 3 {
		return nil, fmt.Errorf("normalize: expected rank-3 array")
	}
	data := img.Data()
	minVal := data[0]
	maxVal := data[0]
	for _, v := range data[1:] {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}
	rng := maxVal - minVal
	if rng == 0 {
		for i := range data {
			data[i] = 0
		}
	} else {
		for i := range data {
			data[i] = (data[i] - minVal) / rng
		}
	}
	return numgo.NewNDArray(shape, data), nil
}

// Rot90 rotates an image 90 degrees counter-clockwise.
// Analogous to tf.image.rot90.
func Rot90(img *numgo.NDArray) (*numgo.NDArray, error) {
	shape := img.Shape()
	if len(shape) != 3 {
		return nil, fmt.Errorf("rot90: expected rank-3 array")
	}
	h, w, c := shape[0], shape[1], shape[2]
	data := img.Data()
	outData := make([]float64, len(data))
	for y := range h {
		for x := range w {
			srcOff := (y*w + x) * c
			// 90 CCW: (x, y) -> (h-1-y, x) but new dims are (w, h)
			// new row = w-1-x, new col = y
			dstOff := ((w-1-x)*h + y) * c
			copy(outData[dstOff:dstOff+c], data[srcOff:srcOff+c])
		}
	}
	return numgo.NewNDArray([]int{w, h, c}, outData), nil
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
