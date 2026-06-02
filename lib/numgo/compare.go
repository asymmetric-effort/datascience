package numgo

import (
	"math"
)

// AllClose returns true if two arrays have the same shape and all corresponding
// elements satisfy |a-b| <= atol + rtol*|b|.
func AllClose(a, b *NDArray, atol, rtol float64) bool {
	if !shapeEqual(a.shape, b.shape) {
		return false
	}
	for i := range a.data {
		diff := math.Abs(a.data[i] - b.data[i])
		if diff > atol+rtol*math.Abs(b.data[i]) {
			return false
		}
	}
	return true
}
