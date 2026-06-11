package readwrite

import (
	"fmt"
	"math"
)

// maxCPDElements is the maximum number of elements allowed in a single CPD
// table. This prevents denial-of-service from crafted model files declaring
// unrealistically large cardinality products.
const maxCPDElements = 100_000_000 // 100M elements (~800 MB of float64)

// safeParentConfigs computes the product of parent cardinalities with overflow
// and size limit checks. Returns an error if the product overflows int or
// exceeds maxCPDElements.
func safeParentConfigs(cardinalities []int) (int, error) {
	p := 1
	for _, c := range cardinalities {
		if c <= 0 {
			return 0, fmt.Errorf("readwrite: invalid cardinality %d", c)
		}
		if p > math.MaxInt/c {
			return 0, fmt.Errorf("readwrite: parent config product overflows int")
		}
		p *= c
		if p > maxCPDElements {
			return 0, fmt.Errorf("readwrite: parent config product %d exceeds limit %d", p, maxCPDElements)
		}
	}
	return p, nil
}
