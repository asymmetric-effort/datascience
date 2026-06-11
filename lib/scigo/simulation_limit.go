package scigo

import "fmt"

// MaxSimulationElements is the default maximum total number of elements
// (nPaths * nSteps) that simulation functions will allocate. This prevents
// denial-of-service via unbounded memory consumption from large parameter
// combinations. Use the *WithLimit variants for larger simulations.
const MaxSimulationElements = 100_000_000 // 100M elements (~800 MB of float64)

// checkSimulationSize validates that the total allocation nPaths * nSteps
// does not exceed maxElements. Panics with a descriptive message if exceeded.
func checkSimulationSize(nPaths, nSteps, maxElements int) {
	total := nPaths * nSteps
	// Check for overflow: if nSteps > 0 and the division doesn't round-trip,
	// the multiplication overflowed.
	if nSteps > 0 && total/nSteps != nPaths {
		panic(fmt.Sprintf("scigo: simulation size overflows int: nPaths=%d * nSteps=%d", nPaths, nSteps))
	}
	if total > maxElements {
		panic(fmt.Sprintf("scigo: simulation size %d exceeds limit %d (nPaths=%d, nSteps=%d); use WithLimit variant for larger simulations", total, maxElements, nPaths, nSteps))
	}
}
